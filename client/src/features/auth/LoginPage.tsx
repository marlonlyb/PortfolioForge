import { useEffect, useRef, useState, type FormEvent } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';
import { useSession } from '../../app/providers/SessionProvider';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';
import { adminLogin, loginWithGoogle, requestEmailLogin } from './api';

interface LoginLocationState {
  from?: string;
}

const LOGIN_PAGE_MODE = {
  PUBLIC: 'public',
  ADMIN: 'admin',
} as const;

type LoginPageMode = (typeof LOGIN_PAGE_MODE)[keyof typeof LOGIN_PAGE_MODE];

interface GoogleCredentialResponse {
  credential?: string;
}

interface GoogleAccounts {
  id: {
    initialize: (config: {
      client_id: string;
      callback: (response: GoogleCredentialResponse) => void;
    }) => void;
    renderButton: (element: HTMLElement, options: Record<string, string>) => void;
  };
}

declare global {
  interface Window {
    google?: {
      accounts: GoogleAccounts;
    };
  }
}

interface LoginPageProps {
  mode?: LoginPageMode;
}

function formatValidationDetail(field?: string, issue?: string): string | null {
  if (field === 'email' && issue === 'required') {
    return 'Email is required.';
  }

  return null;
}

export function LoginPage({ mode = LOGIN_PAGE_MODE.PUBLIC }: LoginPageProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { t } = useLocale();
  const { login } = useSession();
  const state = location.state as LoginLocationState | null;
  const googleButtonRef = useRef<HTMLDivElement | null>(null);
  const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined;
  const isPublicMode = mode === LOGIN_PAGE_MODE.PUBLIC;

  const [loginEmail, setLoginEmail] = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [loginError, setLoginError] = useState<string | null>(null);
  const [googleError, setGoogleError] = useState<string | null>(null);
  const [loginSubmitting, setLoginSubmitting] = useState(false);
  const [googleSubmitting, setGoogleSubmitting] = useState(false);
  const [googleReady, setGoogleReady] = useState(false);

  useEffect(() => {
    setLoginError(null);
    setGoogleError(null);
  }, [isPublicMode]);

  useEffect(() => {
    if (!isPublicMode || !googleClientId) {
      return;
    }

    if (window.google?.accounts.id) {
      setGoogleReady(true);
      return;
    }

    const existingScript = document.querySelector<HTMLScriptElement>('script[data-google-identity="true"]');
    if (existingScript) {
      existingScript.addEventListener('load', () => setGoogleReady(true), { once: true });
      return;
    }

    const script = document.createElement('script');
    script.src = 'https://accounts.google.com/gsi/client';
    script.async = true;
    script.defer = true;
    script.dataset.googleIdentity = 'true';
    script.onload = () => setGoogleReady(true);
    document.head.appendChild(script);
  }, [googleClientId, isPublicMode]);

  useEffect(() => {
    const isPublicLoginView = isPublicMode;

    if (!isPublicLoginView || !googleClientId || !googleReady || !googleButtonRef.current || !window.google?.accounts.id) {
      return;
    }

    window.google.accounts.id.initialize({
      client_id: googleClientId,
      callback: async (response) => {
        if (!response.credential) {
          setGoogleError('Google sign-in did not return a valid credential.');
          return;
        }

        setGoogleError(null);
        setLoginError(null);
        setGoogleSubmitting(true);

        try {
          const session = await loginWithGoogle({ id_token: response.credential });
          login(session);

          if (!session.user.profile_completed) {
            navigate('/complete-profile', { replace: true, state: { from: state?.from ?? '/' } });
            return;
          }

          navigate(state?.from ?? '/', { replace: true });
        } catch (err) {
          if (err instanceof AppError) {
            setGoogleError(err.message);
          } else {
            setGoogleError('Unable to complete Google sign-in.');
          }
        } finally {
          setGoogleSubmitting(false);
        }
      },
    });

    googleButtonRef.current.innerHTML = '';
    window.google.accounts.id.renderButton(googleButtonRef.current, {
      theme: 'outline',
      size: 'large',
      type: 'standard',
      text: 'continue_with',
      shape: 'pill',
    });
  }, [googleClientId, googleReady, isPublicMode, login, navigate, state?.from]);

  async function handleLoginSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoginError(null);
    setGoogleError(null);
    setLoginSubmitting(true);

    try {
      if (isPublicMode) {
        const response = await requestEmailLogin({ email: loginEmail.trim() });
        navigate('/verify-email', {
          replace: true,
          state: {
            email: loginEmail.trim(),
            from: state?.from ?? '/',
            cooldownSeconds: response.cooldown_seconds,
          },
        });
        return;
      }

      const response = await adminLogin({ email: loginEmail, password: loginPassword });
      login(response);
      navigate('/admin/projects', { replace: true });
      return;

    } catch (err) {
      if (err instanceof AppError) {
        if (isPublicMode && err.code === API_ERROR_CODES.VALIDATION_ERROR) {
          const detail = err.details[0];
          setLoginError(formatValidationDetail(detail?.field, detail?.issue) ?? err.message);
        } else if (!isPublicMode && err.code === API_ERROR_CODES.INVALID_CREDENTIALS) {
          setLoginError('Invalid email or password.');
        } else {
          setLoginError(err.message);
        }
      } else {
        setLoginError('An unexpected error occurred. Please try again.');
      }
    } finally {
      setLoginSubmitting(false);
    }
  }

  function renderPublicAuthContent() {
    return (
      <section className="auth-page__panel auth-page__panel--primary">
        <p className="eyebrow">{t.authPublicEyebrow}</p>

        <h2>{t.authPublicLoginTitle}</h2>
        <p className="auth-page__redirect-note">{t.authPublicLoginDescription}</p>

        <div className="auth-page__google-block">
          <div className="auth-page__google-button-wrap" ref={googleButtonRef} aria-live="polite" />
          <span className="auth-page__google-copy">Continue with Google</span>
        </div>

        {!googleClientId ? (
          <p className="auth-page__alt">Google sign-in is not configured in this environment.</p>
        ) : !googleReady ? (
          <p className="auth-page__alt">Loading Google sign-in…</p>
        ) : null}
        {googleSubmitting ? <p className="auth-page__alt">Completing Google sign-in…</p> : null}
        {googleError ? <div className="auth-page__error" role="alert">{googleError}</div> : null}

        <p className="auth-page__helper">{t.authPublicLocalRestriction}</p>

        {loginError ? <div className="auth-page__error" role="alert">{loginError}</div> : null}

        <form className="auth-page__form" onSubmit={handleLoginSubmit}>
          <label className="auth-page__label">
            Email
            <input
              type="email"
              className="auth-page__input"
              value={loginEmail}
              onChange={(event) => setLoginEmail(event.target.value)}
              required
              autoComplete="email"
              disabled={loginSubmitting}
            />
          </label>

          <button type="submit" className="btn btn--primary" disabled={loginSubmitting}>
            {loginSubmitting ? 'Sending code…' : 'Continue with email'}
          </button>
        </form>
      </section>
    );
  }

  return (
    <section className="auth-page">
      <article className="card auth-page__shell">
        <div className="auth-page__stack">
          {isPublicMode ? renderPublicAuthContent() : null}

          {!isPublicMode ? (
            <section className="auth-page__panel auth-page__panel--primary">
              <p className="eyebrow">{t.authAdminEyebrow}</p>
              <h2>{t.authAdminTitle}</h2>
              <p className="auth-page__redirect-note">{t.authAdminDescription}</p>

              {loginError ? <div className="auth-page__error" role="alert">{loginError}</div> : null}

              <form className="auth-page__form" onSubmit={handleLoginSubmit}>
                <label className="auth-page__label">
                  Email
                  <input
                    type="email"
                    className="auth-page__input"
                    value={loginEmail}
                    onChange={(event) => setLoginEmail(event.target.value)}
                    required
                    autoComplete="email"
                    disabled={loginSubmitting}
                  />
                </label>

                <label className="auth-page__label">
                  Password
                  <input
                    type="password"
                    className="auth-page__input"
                    value={loginPassword}
                    onChange={(event) => setLoginPassword(event.target.value)}
                    required
                    autoComplete="current-password"
                    disabled={loginSubmitting}
                  />
                </label>

                <button type="submit" className="btn btn--primary" disabled={loginSubmitting}>
                  {loginSubmitting ? 'Signing in…' : 'Sign in'}
                </button>
              </form>
            </section>
          ) : null}
        </div>

        <p className="auth-page__alt auth-page__alt--footer">
          <Link to={isPublicMode ? '/' : '/login'}>{isPublicMode ? t.authBackToPortfolio : t.authBackToPublicLogin}</Link>
        </p>
      </article>
    </section>
  );
}
