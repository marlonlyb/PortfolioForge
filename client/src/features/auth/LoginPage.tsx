import { useEffect, useRef, useState, type FormEvent } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';
import { useSession, type SessionUser } from '../../app/providers/SessionProvider';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';
import {
  adminLogin,
  loginWithGoogle,
  publicLogin,
  publicSignup,
  type EmailVerificationDispatchResponse,
} from './api';

interface LoginLocationState {
  from?: string;
  notice?: string;
}

const AUTH_VARIANT = {
  LOGIN: 'login',
  SIGNUP: 'signup',
  ADMIN: 'admin',
} as const;

type AuthVariant = (typeof AUTH_VARIANT)[keyof typeof AUTH_VARIANT];

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
  variant?: AuthVariant;
}

interface SignupSuccessState {
  email: string;
  response: EmailVerificationDispatchResponse;
}

function formatValidationDetail(
  field: string | undefined,
  issue: string | undefined,
  t: ReturnType<typeof useLocale>['t'],
): string | null {
  if (field === 'email' && issue === 'required') return t.authValidationEmailRequired;
  if (field === 'password' && issue === 'invalid') return t.authValidationPasswordInvalid;
  if (field === 'confirm_password' && issue === 'required') return t.authValidationConfirmRequired;
  if (field === 'confirm_password' && issue === 'mismatch') return t.authValidationConfirmMismatch;
  return null;
}

export function LoginPage({ variant = AUTH_VARIANT.LOGIN }: LoginPageProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { t } = useLocale();
  const { login } = useSession();
  const state = (location.state as LoginLocationState | null) ?? null;
  const googleButtonRef = useRef<HTMLDivElement | null>(null);
  const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined;
  const isSignupVariant = variant === AUTH_VARIANT.SIGNUP;
  const isAdminVariant = variant === AUTH_VARIANT.ADMIN;

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
  }, [isAdminVariant, isSignupVariant]);

  useEffect(() => {
    if (isAdminVariant || !googleClientId) return;

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
  }, [googleClientId, isAdminVariant]);

  useEffect(() => {
    if (isAdminVariant || !googleClientId || !googleReady || !googleButtonRef.current || !window.google?.accounts.id) {
      return;
    }

    window.google.accounts.id.initialize({
      client_id: googleClientId,
      callback: async (response) => {
        if (!response.credential) {
          setGoogleError(t.authGoogleCredentialMissing);
          return;
        }

        setGoogleError(null);
        setFormError(null);
        setGoogleSubmitting(true);

        try {
          const session = await loginWithGoogle({ id_token: response.credential });
          handleSuccessfulLogin(session.user, session);
        } catch (err) {
          setGoogleError(err instanceof AppError ? err.message : t.authGoogleSignInFailed);
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
  }, [googleClientId, googleReady, isAdminVariant, isSignupVariant, login, navigate, state?.from]);

  function resolvePostLoginPath(user: SessionUser): string {
    if (state?.from) {
      return state.from;
    }

    return user.is_admin ? '/admin/projects' : '/';
  }

  function handleSuccessfulLogin(user: SessionUser, response: Parameters<typeof login>[0]) {
    login(response);

    if (!isAdminVariant && !user.is_admin && !user.profile_completed) {
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

      const response = isAdminVariant
        ? await adminLogin({ email: email.trim(), password })
        : await publicLogin({ email: email.trim(), password });
      handleSuccessfulLogin(response.user, response);
    } catch (err) {
      if (err instanceof AppError) {
        if (err.code === API_ERROR_CODES.VALIDATION_ERROR) {
          const detail = err.details[0];
          setFormError(formatValidationDetail(detail?.field, detail?.issue, t) ?? err.message);
        } else if (err.code === API_ERROR_CODES.INVALID_CREDENTIALS) {
          setFormError(t.authInvalidCredentials);
        } else if (err.code === API_ERROR_CODES.FORBIDDEN && isAdminVariant) {
          setFormError(t.authForbiddenAdmin);
        } else if (err.code === API_ERROR_CODES.PASSWORD_SETUP_REQUIRED) {
          setFormError(t.authPasswordSetupRequired);
        } else {
          setFormError(err.message);
        }
      } else {
        setFormError(t.authUnexpectedError);
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
          <span className="auth-page__google-copy">{t.authGoogleContinue}</span>
        </div>

        {!googleClientId ? (
          <p className="auth-page__alt">{t.authGoogleNotConfigured}</p>
        ) : !googleReady ? (
          <p className="auth-page__alt">{t.authGoogleLoading}</p>
        ) : null}
        {googleSubmitting ? <p className="auth-page__alt">{t.authGoogleCompleting}</p> : null}
        {googleError ? <div className="auth-page__error" role="alert">{googleError}</div> : null}
      </>
    );
  }

  function renderPublicSuccess() {
    if (!signupSuccess) return null;

    return (
      <section className="auth-page__panel auth-page__panel--primary">
        <p className="eyebrow">{t.authSignupSuccessEyebrow}</p>
        <h2>{t.authSignupSuccessTitle}</h2>
        <p className="auth-page__redirect-note">{signupSuccess.response.message}</p>
        <p className="auth-page__helper">{t.authSignupSuccessHelper}</p>

        <div className="auth-page__form">
          <Link
            className="btn btn--primary"
            to="/verify-email"
            state={{ email: signupSuccess.email, cooldownSeconds: signupSuccess.response.cooldown_seconds }}
          >
            {t.authSignupSuccessVerifyCta}
          </Link>
          <Link className="btn btn--secondary" to="/login">
            {t.authSignupSuccessBackLogin}
          </Link>
        </div>
      </section>
    );
  }

  function renderPublicForm() {
    if (signupSuccess) return renderPublicSuccess();

    const eyebrow = isAdminVariant ? t.authAdminEyebrow : t.authPublicEyebrow;
    const title = isAdminVariant
      ? t.authAdminTitle
      : (isSignupVariant ? t.authPublicSignupTitle : t.authPublicLoginTitle);
    const description = isAdminVariant
      ? t.authAdminDescription
      : (isSignupVariant ? t.authPublicSignupDescription : t.authPublicLoginDescription);

    return (
      <section className="auth-page__panel auth-page__panel--primary">
        <p className="eyebrow">{eyebrow}</p>
        <h2>{title}</h2>
        <p className="auth-page__redirect-note">{description}</p>

        {isAdminVariant ? null : renderGoogleBlock()}

        {isAdminVariant ? null : <p className="auth-page__helper">{t.authPublicLocalRestriction}</p>}
        {formNotice ? <p className="auth-page__alt">{formNotice}</p> : null}
        {formError ? <div className="auth-page__error" role="alert">{formError}</div> : null}

        <form className="auth-page__form" onSubmit={handleSubmit}>
          <label className="auth-page__label">
            {t.authFieldEmail}
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
            {t.authFieldPassword}
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
              {t.authFieldConfirmPassword}
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
            {submitting
              ? (isSignupVariant ? t.authSubmitCreatingAccount : t.authSubmitSigningIn)
              : (isSignupVariant ? t.authSubmitCreateAccount : t.authSubmitSignIn)}
          </button>
        </form>

        {isAdminVariant ? null : (
          <p className="auth-page__alt">
            {isSignupVariant ? (
              <Link to="/login">{t.authAltAlreadyHaveAccount}</Link>
            ) : (
              <Link to="/signup">{t.authAltNeedAccount}</Link>
            )}
          </p>
        )}
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
          <Link to={isAdminVariant ? '/login' : '/'}>
            {isAdminVariant ? t.authBackToPublicLogin : t.authBackToPortfolio}
          </Link>
        </p>
      </article>
    </section>
  );
}
