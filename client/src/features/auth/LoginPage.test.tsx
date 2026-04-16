import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { SessionProvider, type SessionUser } from '../../app/providers/SessionProvider';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';
import { LoginPage } from './LoginPage';
import { adminLogin, loginWithGoogle, publicLogin, publicSignup } from './api';

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
    adminLogin: vi.fn(),
    loginWithGoogle: vi.fn(),
    publicLogin: vi.fn(),
    publicSignup: vi.fn(),
  };
});

const mockedAdminLogin = vi.mocked(adminLogin);
const mockedLoginWithGoogle = vi.mocked(loginWithGoogle);
const mockedPublicLogin = vi.mocked(publicLogin);
const mockedPublicSignup = vi.mocked(publicSignup);

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

function renderLoginPage({ mode = 'public', variant = 'login', initialEntries = ['/login'] }: { mode?: 'public' | 'admin'; variant?: 'login' | 'signup'; initialEntries?: Array<string | { pathname: string; state?: { from?: string; notice?: string } }>; } = {}) {
  return render(
    <MemoryRouter initialEntries={initialEntries}>
      <SessionProvider>
        <LocaleProvider>
          <Routes>
            <Route path="/login" element={<LoginPage mode={mode} variant={variant} />} />
            <Route path="/signup" element={<LoginPage mode="public" variant="signup" />} />
            <Route path="/admin/login" element={<LoginPage mode="admin" />} />
            <Route path="/verify-email" element={<p>verify email destination</p>} />
            <Route path="/complete-profile" element={<p>complete profile destination</p>} />
            <Route path="/admin/projects" element={<p>admin destination</p>} />
            <Route path="/dashboard" element={<p>dashboard destination</p>} />
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
    mockedPublicLogin.mockReset();
    mockedPublicSignup.mockReset();
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

  it('logs public users in with email and password', async () => {
    mockedPublicLogin.mockResolvedValue(buildLoginResponse({ email_verified: true }));
    renderLoginPage();

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret-123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));

    await waitFor(() => {
      expect(mockedPublicLogin).toHaveBeenCalledWith({ email: 'ada@example.com', password: 'secret-123' });
    });
  });

  it('shows password-setup guidance for migrated local accounts', async () => {
    mockedPublicLogin.mockRejectedValue(new AppError(409, {
      code: API_ERROR_CODES.PASSWORD_SETUP_REQUIRED,
      message: 'This account still needs a password setup or reset before it can sign in.',
    }));

    renderLoginPage();
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret-123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('This account still needs a password setup or reset before it can sign in.');
  });

  it('creates a public local account and shows signup success', async () => {
    mockedPublicSignup.mockResolvedValue({
      verification_required: true,
      message: 'Account created. Check your email for the verification code.',
      cooldown_seconds: 60,
    });

    renderLoginPage({ variant: 'signup', initialEntries: ['/signup'] });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'ada@example.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret-123' } });
    fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'secret-123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create account' }));

    expect(await screen.findByRole('heading', { name: 'Check your email' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Verify email' })).toHaveAttribute('href', '/verify-email');
  });

  it('keeps admin mode isolated from public signup options', async () => {
    mockedAdminLogin.mockResolvedValue(buildLoginResponse({ is_admin: true }));
    renderLoginPage({ mode: 'admin', initialEntries: ['/admin/login'] });

    expect(screen.getByRole('heading', { name: 'Admin access' })).toBeInTheDocument();
    expect(screen.queryByText('Need an account? Sign up')).not.toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'admin@example.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'secret-123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));

    expect(await screen.findByText('admin destination')).toBeInTheDocument();
  });

  it('sends incomplete Google users through complete profile before returning to public pages', async () => {
    const initializeGoogle = vi.fn();
    mockedLoginWithGoogle.mockResolvedValue(buildLoginResponse({ auth_provider: 'google', email_verified: true, profile_completed: false }));
    window.google = {
      accounts: {
        id: {
          initialize: initializeGoogle,
          renderButton: vi.fn(),
        },
      },
    };

    renderLoginPage({ initialEntries: [{ pathname: '/login', state: { from: '/dashboard' } }] });
    await waitFor(() => expect(initializeGoogle).toHaveBeenCalled());
    const callback = initializeGoogle.mock.calls[0]?.[0]?.callback as ((response: { credential?: string }) => Promise<void>) | undefined;
    await callback?.({ credential: 'google-token' });
    expect(await screen.findByText('complete profile destination')).toBeInTheDocument();
  });
});
