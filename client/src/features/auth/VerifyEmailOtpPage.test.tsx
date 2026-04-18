import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { useLocale } from '../../app/providers/LocaleProvider';
import { SessionProvider, type SessionUser } from '../../app/providers/SessionProvider';
import { VerifyEmailOtpPage } from './VerifyEmailOtpPage';
import { resendEmailVerification, verifyEmailVerification } from './api';

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
    resendEmailVerification: vi.fn(),
    verifyEmailVerification: vi.fn(),
  };
});

const mockedResendEmailVerification = vi.mocked(resendEmailVerification);
const mockedVerifyEmailVerification = vi.mocked(verifyEmailVerification);

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

function renderVerifyPage(initialEntry: string | { pathname: string; state?: { email?: string; from?: string; cooldownSeconds?: number } }) {
  function LocaleControls() {
    const { setLocale } = useLocale();
    return (
      <div>
        <button type="button" onClick={() => setLocale('es')}>locale-es</button>
        <button type="button" onClick={() => setLocale('en')}>locale-en</button>
      </div>
    );
  }

  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <SessionProvider>
        <LocaleProvider>
          <LocaleControls />
          <Routes>
            <Route path="/verify-email" element={<VerifyEmailOtpPage />} />
            <Route path="/login" element={<p>login destination</p>} />
            <Route path="/complete-profile" element={<p>complete profile destination</p>} />
            <Route path="/dashboard" element={<p>dashboard destination</p>} />
          </Routes>
        </LocaleProvider>
      </SessionProvider>
    </MemoryRouter>,
  );
}

describe('VerifyEmailOtpPage', () => {
  beforeEach(() => {
    mockedResendEmailVerification.mockReset();
    mockedVerifyEmailVerification.mockReset();
    window.localStorage.clear();
    window.sessionStorage.clear();
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('verifies signed-out users and routes them back to login', async () => {
    mockedVerifyEmailVerification.mockResolvedValue({ user: buildSessionUser({ email_verified: true }) });
    renderVerifyPage({ pathname: '/verify-email', state: { email: 'ada@example.com', cooldownSeconds: 0 } });

    fireEvent.click(screen.getByRole('button', { name: 'locale-en' }));

    fireEvent.change(screen.getByLabelText('Verification code'), { target: { value: '123456' } });
    fireEvent.click(screen.getByRole('button', { name: 'Verify email' }));

    expect(await screen.findByText('login destination')).toBeInTheDocument();
  });

  it('allows resend when the cooldown is clear', async () => {
    mockedResendEmailVerification.mockResolvedValue({
      verification_required: true,
      message: 'A fresh verification code is on the way.',
      cooldown_seconds: 60,
    });

    renderVerifyPage({ pathname: '/verify-email', state: { email: 'ada@example.com', cooldownSeconds: 0 } });
    fireEvent.click(screen.getByRole('button', { name: 'locale-en' }));
    fireEvent.click(await screen.findByRole('button', { name: 'Resend code' }));

    await waitFor(() => {
      expect(mockedResendEmailVerification).toHaveBeenCalledWith({ email: 'ada@example.com' });
    });
  });

  it('refreshes the signed-in session after verification', async () => {
    window.sessionStorage.setItem('auth_token', 'token');
    let privateMeCalls = 0;
    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
      const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
      if (url.endsWith('/api/v1/private/me')) {
        privateMeCalls += 1;
        return new Response(JSON.stringify({ data: buildSessionUser({ email_verified: true, assistant_eligible: true, can_use_project_assistant: true }) }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        });
      }
      throw new Error(`Unhandled fetch: ${url}`);
    }));

    mockedVerifyEmailVerification.mockResolvedValue({ user: buildSessionUser({ email_verified: true }) });
    renderVerifyPage({ pathname: '/verify-email', state: { email: 'ada@example.com', from: '/dashboard', cooldownSeconds: 0 } });
    fireEvent.click(screen.getByRole('button', { name: 'locale-en' }));

    await waitFor(() => {
      expect(screen.getByDisplayValue('ada@example.com')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Verification code'), { target: { value: '123456' } });
    fireEvent.click(screen.getByRole('button', { name: 'Verify email' }));

    expect(await screen.findByText('dashboard destination')).toBeInTheDocument();
    expect(privateMeCalls).toBeGreaterThan(0);
  });

  it('updates verification copy when the locale changes', async () => {
    renderVerifyPage({ pathname: '/verify-email', state: { email: 'ada@example.com', cooldownSeconds: 0 } });

    expect(await screen.findByText('Verifica tu email')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'locale-en' }));

    expect(await screen.findByText('Verify your email')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Verify email' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Resend code' })).toBeInTheDocument();
  });
});
