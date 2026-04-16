import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { SessionProvider, useSession, type SessionUser } from '../../app/providers/SessionProvider';
import { VerifyEmailOtpPage } from './VerifyEmailOtpPage';
import { requestEmailLogin, verifyEmailLogin } from './api';

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
    requestEmailLogin: vi.fn(),
    verifyEmailLogin: vi.fn(),
  };
});

const mockedRequestEmailLogin = vi.mocked(requestEmailLogin);
const mockedVerifyEmailLogin = vi.mocked(verifyEmailLogin);

function buildSessionUser(overrides: Partial<SessionUser> = {}): SessionUser {
  return {
    id: 'user-1',
    email: 'ada@example.com',
    is_admin: false,
    auth_provider: 'local',
    email_verified: false,
    full_name: 'Ada Lovelace',
    company: 'Analytical Engines',
    profile_completed: true,
    assistant_eligible: false,
    can_use_project_assistant: false,
    created_at: '2026-04-15T00:00:00Z',
    ...overrides,
  };
}

function buildLoginResponse(overrides: Partial<SessionUser> = {}) {
  return {
    user: buildSessionUser(overrides),
    token: 'token',
    expires_in: 3600,
  };
}

function AssistantGateDestination() {
  const { user } = useSession();

  return (
    <>
      <p>{user?.email_verified ? 'verified state' : 'unverified state'}</p>
      <p>{user?.can_use_project_assistant ? 'assistant unlocked' : 'assistant locked'}</p>
    </>
  );
}

function renderVerifyPage(initialEntry: string | { pathname: string; state?: { email?: string; from?: string; cooldownSeconds?: number } }) {
  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <SessionProvider>
        <LocaleProvider>
          <Routes>
            <Route path="/verify-email" element={<VerifyEmailOtpPage />} />
            <Route path="/complete-profile" element={<p>complete profile destination</p>} />
            <Route path="/dashboard" element={<p>dashboard destination</p>} />
            <Route path="/assistant" element={<AssistantGateDestination />} />
            <Route path="/login" element={<p>login destination</p>} />
          </Routes>
        </LocaleProvider>
      </SessionProvider>
    </MemoryRouter>,
  );
}

describe('VerifyEmailOtpPage', () => {
  beforeEach(() => {
    mockedRequestEmailLogin.mockReset();
    mockedVerifyEmailLogin.mockReset();
    window.localStorage.clear();
    window.sessionStorage.clear();
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('prefills email from navigation state and verifies signed-out users back to login', async () => {
    mockedVerifyEmailLogin.mockResolvedValue({ user: buildSessionUser({ email_verified: true, profile_completed: false }), token: 'token', expires_in: 3600 });

    renderVerifyPage({ pathname: '/verify-email', state: { email: 'ada@example.com', from: '/dashboard', cooldownSeconds: 0 } });

    fireEvent.change(screen.getByLabelText('Verification code'), { target: { value: '123456' } });
    fireEvent.click(screen.getByRole('button', { name: 'Complete sign in' }));

    expect(await screen.findByText('complete profile destination')).toBeInTheDocument();
  });

  it('allows resend when the cooldown is clear', async () => {
    mockedRequestEmailLogin.mockResolvedValue({
      verification_required: true,
      message: 'If the account is eligible, a verification code will be sent shortly.',
      cooldown_seconds: 60,
    });

    renderVerifyPage({ pathname: '/verify-email', state: { email: 'ada@example.com', cooldownSeconds: 0 } });

    const resendButton = await screen.findByRole('button', { name: 'Resend code' });
    fireEvent.click(resendButton);

    await waitFor(() => {
      expect(mockedRequestEmailLogin).toHaveBeenCalledWith({ email: 'ada@example.com' });
    });
  });

  it('refreshes the signed-in session after OTP verification and unlocks assistant gating', async () => {
    const initialUser = buildSessionUser({ email_verified: false, assistant_eligible: false, can_use_project_assistant: false });
    const refreshedUser = buildSessionUser({ email_verified: true, assistant_eligible: true, can_use_project_assistant: true });
    let privateMeCalls = 0;

    window.sessionStorage.setItem('auth_token', 'token');
    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
      const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
      if (url.endsWith('/api/v1/private/me')) {
        privateMeCalls += 1;
        return new Response(JSON.stringify({ data: privateMeCalls === 1 ? initialUser : refreshedUser }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        });
      }

      throw new Error(`Unhandled fetch: ${url}`);
    }));

    mockedVerifyEmailLogin.mockResolvedValue(buildLoginResponse({ email_verified: true, profile_completed: false, assistant_eligible: false, can_use_project_assistant: false }));

    renderVerifyPage({ pathname: '/verify-email', state: { email: 'ada@example.com', from: '/assistant', cooldownSeconds: 0 } });

    await waitFor(() => {
      expect(screen.getByDisplayValue('ada@example.com')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Verification code'), { target: { value: '123456' } });
    fireEvent.click(screen.getByRole('button', { name: 'Complete sign in' }));

    expect(await screen.findByText('verified state')).toBeInTheDocument();
    expect(screen.getByText('assistant unlocked')).toBeInTheDocument();
    expect(privateMeCalls).toBe(2);
  });
});
