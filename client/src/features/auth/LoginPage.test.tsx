import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { SessionProvider, type SessionUser } from '../../app/providers/SessionProvider';
import { LoginPage } from './LoginPage';
import { VerifyEmailOtpPage } from './VerifyEmailOtpPage';
import { adminLogin, loginWithGoogle, requestEmailLogin } from './api';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
    adminLogin: vi.fn(),
    loginWithGoogle: vi.fn(),
    requestEmailLogin: vi.fn(),
  };
});

const mockedAdminLogin = vi.mocked(adminLogin);
const mockedLoginWithGoogle = vi.mocked(loginWithGoogle);
const mockedRequestEmailLogin = vi.mocked(requestEmailLogin);

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

function GoogleDestination() {
  return <p>google destination</p>;
}

function PublicDestination() {
  return <p>public destination</p>;
}

function AdminDestination() {
  return <p>admin destination</p>;
}

function CompleteProfileDestination() {
  return <p>complete profile destination</p>;
}

function VerifyEmailDestination() {
  return <VerifyEmailOtpPage />;
}

interface RenderLoginPageOptions {
  mode?: 'public' | 'admin';
  initialEntries?: Array<string | { pathname: string; state?: { from?: string } }>;
}

function renderLoginPage(options: RenderLoginPageOptions = {}) {
  const { mode = 'public', initialEntries = ['/login'] } = options;

  return render(
    <MemoryRouter initialEntries={initialEntries}>
      <SessionProvider>
        <LocaleProvider>
          <Routes>
            <Route path="/login" element={<LoginPage mode={mode} />} />
            <Route path="/admin/login" element={<LoginPage mode="admin" />} />
            <Route path="/projects/portfolioforge" element={<GoogleDestination />} />
            <Route path="/dashboard" element={<PublicDestination />} />
            <Route path="/admin/projects" element={<AdminDestination />} />
            <Route path="/verify-email" element={<VerifyEmailDestination />} />
            <Route path="/complete-profile" element={<CompleteProfileDestination />} />
          </Routes>
        </LocaleProvider>
      </SessionProvider>
    </MemoryRouter>,
  );
}

describe('LoginPage', () => {
  beforeEach(() => {
    mockedAdminLogin.mockReset();
    mockedLoginWithGoogle.mockReset();
    mockedRequestEmailLogin.mockReset();
    vi.stubEnv('VITE_GOOGLE_CLIENT_ID', 'google-client-id');
    window.localStorage.clear();
    window.sessionStorage.clear();
    window.localStorage.setItem('portfolioforge.locale', 'en');
    window.google = {
      accounts: {
        id: {
          initialize: vi.fn(),
          renderButton: vi.fn(),
        },
      },
    };
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllEnvs();
    delete window.google;
  });

  it('routes public email login into the OTP verification flow', async () => {
    mockedRequestEmailLogin.mockResolvedValue({
      verification_required: true,
      message: 'If the account is eligible, a verification code will be sent shortly.',
      cooldown_seconds: 60,
    });

    renderLoginPage();

    expect(screen.getByRole('heading', { name: 'Login to PortfolioForge' })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: 'Continue with email' }));

    expect(await screen.findByRole('heading', { name: 'Verify your email code' })).toBeInTheDocument();
    expect(screen.getByDisplayValue('ada@example.com')).toBeInTheDocument();

    expect(mockedRequestEmailLogin).toHaveBeenCalledWith({ email: 'ada@example.com' });
  });

  it('keeps the user on login with actionable feedback when the email request is rejected', async () => {
    mockedRequestEmailLogin.mockRejectedValue(new AppError(400, {
      code: API_ERROR_CODES.UNEXPECTED_ERROR,
      message: 'Unable to process the email login request',
    }));

    renderLoginPage();

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: 'Continue with email' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Unable to process the email login request');
    expect(screen.getByRole('heading', { name: 'Login to PortfolioForge' })).toBeInTheDocument();
  });

  it('redirects public email login into verification with the requested destination', async () => {
    mockedRequestEmailLogin.mockResolvedValue({
      verification_required: true,
      message: 'If the account is eligible, a verification code will be sent shortly.',
      cooldown_seconds: 60,
    });

    renderLoginPage({ initialEntries: [{ pathname: '/login', state: { from: '/dashboard' } }] });

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: 'Continue with email' }));

    expect(await screen.findByRole('heading', { name: 'Verify your email code' })).toBeInTheDocument();
  });

  it('does not show public password fields', () => {
    renderLoginPage();

    expect(screen.queryByLabelText('Password')).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Sign up' })).not.toBeInTheDocument();
  });

  it('keeps admin mode hidden from public auth options and redirects admins to admin projects', async () => {
		mockedAdminLogin.mockResolvedValue(buildLoginResponse({ is_admin: true }));

    renderLoginPage({ mode: 'admin', initialEntries: ['/admin/login'] });

    expect(screen.getByRole('heading', { name: 'Admin access' })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Sign up' })).not.toBeInTheDocument();
    expect(screen.queryByText('Continue with Google')).not.toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'admin@example.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret-123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));

    expect(await screen.findByText('admin destination')).toBeInTheDocument();
  });

  it('sends incomplete Google users through complete profile before returning to public pages', async () => {
    const initializeGoogle = vi.fn();
    const renderGoogleButton = vi.fn();
    mockedLoginWithGoogle.mockResolvedValue(buildLoginResponse({
      auth_provider: 'google',
      email_verified: true,
      profile_completed: false,
    }));
    window.google = {
      accounts: {
        id: {
          initialize: initializeGoogle,
          renderButton: renderGoogleButton,
        },
      },
    };

    renderLoginPage({ initialEntries: [{ pathname: '/login', state: { from: '/projects/portfolioforge' } }] });

    await waitFor(() => {
      expect(initializeGoogle).toHaveBeenCalled();
    });

    const callback = initializeGoogle.mock.calls[0]?.[0]?.callback as ((response: { credential?: string }) => Promise<void>) | undefined;
    expect(callback).toBeTypeOf('function');

    await callback?.({ credential: 'google-token' });

    expect(await screen.findByText('complete profile destination')).toBeInTheDocument();
    expect(mockedLoginWithGoogle).toHaveBeenCalledWith({ id_token: 'google-token' });
    expect(renderGoogleButton).toHaveBeenCalled();
  });
});
