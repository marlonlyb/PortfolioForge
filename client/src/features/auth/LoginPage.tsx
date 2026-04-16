import { useEffect, useRef, useState, type FormEvent } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';
import { useSession, type SessionUser } from '../../app/providers/SessionProvider';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';
import {
  loginWithGoogle,
  publicLogin,
  publicSignup,
  type EmailVerificationDispatchResponse,
} from './api';

interface LoginLocationState {
  from?: string;
  notice?: string;
}

const PUBLIC_AUTH_VARIANT = {
  LOGIN: 'login',
  SIGNUP: 'signup',
} as const;

type PublicAuthVariant = (typeof PUBLIC_AUTH_VARIANT)[keyof typeof PUBLIC_AUTH_VARIANT];

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
  variant?: PublicAuthVariant;
}

interface SignupSuccessState {
  email: string;
  response: EmailVerificationDispatchResponse;
}

function formatValidationDetail(field?: string, issue?: string): string | null {
  if (field === 'email' && issue === 'required') return 'Email is required.';
  if (field === 'password' && issue === 'invalid') return 'Password must contain at least 8 characters.';
  if (field === 'confirm_password' && issue === 'required') return 'Confirm password is required.';
  if (field === 'confirm_password' && issue === 'mismatch') return 'Passwords must match.';
  return null;
}

export function LoginPage({ variant = PUBLIC_AUTH_VARIANT.LOGIN }: LoginPageProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { t } = useLocale();
  const { login } = useSession();
  const state = (location.state as LoginLocationState | null) ?? null;
  const googleButtonRef = useRef<HTMLDivElement | null>(null);
  const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined;
  const isSignupVariant = variant === PUBLIC_AUTH_VARIANT.SIGNUP;

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [formError, setFormError] = useState<string | null>(null);
  const [formNotice, setFormNotice] = useState<string | null>(state?.notice ?? null);
  const [googleError, setGoogleError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [googleSubmitting, setGoogleSubmitting] = useState(false);
  const [googleReady, setGoogleReady] = useState(false);
  const [signupSuccess, setSignupSuccess] = useState<SignupSuccessState | null>(null);

  useEffect(() => {
    setFormError(null);
    setGoogleError(null);
    setSignupSuccess(null);
  }, [isSignupVariant]);

  useEffect(() => {
    if (!googleClientId) return;

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
  }, [googleClientId]);

  useEffect(() => {
    if (!googleClientId || !googleReady || !googleButtonRef.current || !window.google?.accounts.id) {
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
        setFormError(null);
        setGoogleSubmitting(true);

        try {
          const session = await loginWithGoogle({ id_token: response.credential });
          handleSuccessfulLogin(session.user, session);
        } catch (err) {
          setGoogleError(err instanceof AppError ? err.message : 'Unable to complete Google sign-in.');
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
      text: isSignupVariant ? 'signup_with' : 'continue_with',
      shape: 'pill',
    });
  }, [googleClientId, googleReady, isSignupVariant, login, navigate, state?.from]);

  function resolvePostLoginPath(user: SessionUser): string {
    if (state?.from) {
      return state.from;
    }

    return user.is_admin ? '/admin/projects' : '/';
  }

  function handleSuccessfulLogin(user: SessionUser, response: Parameters<typeof login>[0]) {
    login(response);

    if (!user.profile_completed) {
      navigate('/complete-profile', { replace: true, state: { from: resolvePostLoginPath(user) } });
      return;
    }

    navigate(resolvePostLoginPath(user), { replace: true });
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setFormError(null);
    setGoogleError(null);
    setFormNotice(null);

    try {
      if (isSignupVariant) {
        const response = await publicSignup({ email: email.trim(), password, confirm_password: confirmPassword });
        setSignupSuccess({ email: email.trim(), response });
        return;
      }

      const response = await publicLogin({ email: email.trim(), password });
      handleSuccessfulLogin(response.user, response);
    } catch (err) {
      if (err instanceof AppError) {
        if (err.code === API_ERROR_CODES.VALIDATION_ERROR) {
          const detail = err.details[0];
          setFormError(formatValidationDetail(detail?.field, detail?.issue) ?? err.message);
        } else if (err.code === API_ERROR_CODES.INVALID_CREDENTIALS) {
          setFormError('Invalid email or password.');
        } else if (err.code === API_ERROR_CODES.PASSWORD_SETUP_REQUIRED) {
          setFormError('This account still needs a password setup or reset before it can sign in.');
        } else {
          setFormError(err.message);
        }
      } else {
        setFormError('An unexpected error occurred. Please try again.');
      }
    } finally {
      setSubmitting(false);
    }
  }

  function renderGoogleBlock() {
    return (
      <>
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
      </>
    );
  }

  function renderPublicSuccess() {
    if (!signupSuccess) return null;

    return (
      <section className="auth-page__panel auth-page__panel--primary">
        <p className="eyebrow">Public sign up</p>
        <h2>Check your email</h2>
        <p className="auth-page__redirect-note">{signupSuccess.response.message}</p>
        <p className="auth-page__helper">Verify the code first, then log in with your new password.</p>

        <div className="auth-page__form">
          <Link
            className="btn btn--primary"
            to="/verify-email"
            state={{ email: signupSuccess.email, cooldownSeconds: signupSuccess.response.cooldown_seconds }}
          >
            Verify email
          </Link>
          <Link className="btn btn--secondary" to="/login">
            Back to login
          </Link>
        </div>
      </section>
    );
  }

  function renderPublicForm() {
    if (signupSuccess) return renderPublicSuccess();

    return (
      <section className="auth-page__panel auth-page__panel--primary">
        <p className="eyebrow">{t.authPublicEyebrow}</p>
        <h2>{isSignupVariant ? t.authPublicSignupTitle : t.authPublicLoginTitle}</h2>
        <p className="auth-page__redirect-note">{isSignupVariant ? t.authPublicSignupDescription : t.authPublicLoginDescription}</p>

        {renderGoogleBlock()}

        <p className="auth-page__helper">{t.authPublicLocalRestriction}</p>
        {formNotice ? <p className="auth-page__alt">{formNotice}</p> : null}
        {formError ? <div className="auth-page__error" role="alert">{formError}</div> : null}

        <form className="auth-page__form" onSubmit={handleSubmit}>
          <label className="auth-page__label">
            Email
            <input
              type="email"
              className="auth-page__input"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
              autoComplete="email"
              disabled={submitting}
            />
          </label>

          <label className="auth-page__label">
            Password
            <input
              type="password"
              className="auth-page__input"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
              autoComplete={isSignupVariant ? 'new-password' : 'current-password'}
              disabled={submitting}
            />
          </label>

          {isSignupVariant ? (
            <label className="auth-page__label">
              Confirm password
              <input
                type="password"
                className="auth-page__input"
                value={confirmPassword}
                onChange={(event) => setConfirmPassword(event.target.value)}
                required
                autoComplete="new-password"
                disabled={submitting}
              />
            </label>
          ) : null}

          <button type="submit" className="btn btn--primary" disabled={submitting}>
            {submitting ? (isSignupVariant ? 'Creating account…' : 'Signing in…') : (isSignupVariant ? 'Create account' : 'Sign in')}
          </button>
        </form>

        <p className="auth-page__alt">
          {isSignupVariant ? (
            <Link to="/login">Already have an account? Log in</Link>
          ) : (
            <Link to="/signup">Need an account? Sign up</Link>
          )}
        </p>
      </section>
    );
  }

  return (
    <section className="auth-page">
      <article className="card auth-page__shell">
        <div className="auth-page__stack">
          {renderPublicForm()}
        </div>

        <p className="auth-page__alt auth-page__alt--footer">
          <Link to="/">{t.authBackToPortfolio}</Link>
        </p>
      </article>
    </section>
  );
}
