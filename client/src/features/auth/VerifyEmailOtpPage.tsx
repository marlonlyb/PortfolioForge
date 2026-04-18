import { useEffect, useState, type FormEvent } from 'react';
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';
import { useSession } from '../../app/providers/SessionProvider';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';
import { resendEmailVerification, verifyEmailVerification } from './api';

interface VerifyEmailLocationState {
  email?: string;
  from?: string;
  cooldownSeconds?: number;
}

export function VerifyEmailOtpPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const { t } = useLocale();
  const { user, refreshSession, setUser } = useSession();
  const state = (location.state as VerifyEmailLocationState | null) ?? null;

  const [email, setEmail] = useState(state?.email ?? user?.email ?? '');
  const [code, setCode] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [serverNotice, setServerNotice] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [resending, setResending] = useState(false);
  const [resendRemaining, setResendRemaining] = useState(state?.cooldownSeconds ?? 60);

  useEffect(() => {
    if (resendRemaining <= 0) return undefined;

    const timeoutId = window.setTimeout(() => {
      setResendRemaining((current) => Math.max(current - 1, 0));
    }, 1000);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [resendRemaining]);

  if (user?.is_admin) {
    return <Navigate replace to="/admin/projects" />;
  }

  async function handleVerify(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const response = await verifyEmailVerification({ email: email.trim(), code: code.trim() });

      if (user) {
        setUser(response.user);
        const refreshedUser = await refreshSession();
        const nextUser = refreshedUser ?? response.user;

        if (!nextUser.profile_completed) {
          navigate('/complete-profile', { replace: true, state: { from: state?.from ?? '/' } });
          return;
        }

        navigate(state?.from ?? '/', { replace: true });
        return;
      }

        navigate('/login', {
          replace: true,
          state: { notice: t.authVerifyLoginNotice },
        });
    } catch (err) {
      if (err instanceof AppError) {
        if (err.code === API_ERROR_CODES.OTP_INVALID) {
          setError(t.authVerifyInvalidCode);
        } else if (err.code === API_ERROR_CODES.OTP_EXPIRED) {
          setError(t.authVerifyExpiredCode);
        } else {
          setError(err.message);
        }
      } else {
        setError(t.authVerifyUnableToVerify);
      }
    } finally {
      setSubmitting(false);
    }
  }

  async function handleResend() {
    setResending(true);
    setError(null);

    try {
      const response = await resendEmailVerification({ email: email.trim() });
      setServerNotice(response.message);
      setResendRemaining(response.cooldown_seconds);
    } catch (err) {
      if (err instanceof AppError && err.code === API_ERROR_CODES.VALIDATION_ERROR) {
        setError(err.message);
      } else {
        setError(t.authVerifyUnableToResend);
      }
    } finally {
      setResending(false);
    }
  }

  const notice = serverNotice ?? t.authVerifyDefaultNotice;

  return (
      <section className="auth-page">
      <article className="card">
        <p className="eyebrow">{t.authVerifyEyebrow}</p>
        <h2>{t.authVerifyTitle}</h2>
        <p className="auth-page__redirect-note">{notice}</p>

        {error ? <div className="auth-page__error" role="alert">{error}</div> : null}

        <form className="auth-page__form" onSubmit={handleVerify}>
          <label className="auth-page__label">
            {t.authFieldEmail}
            <input
              type="email"
              className="auth-page__input"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
              autoComplete="email"
              disabled={submitting || resending}
            />
          </label>

          <label className="auth-page__label">
            {t.authVerifyCodeLabel}
            <input
              type="text"
              inputMode="numeric"
              pattern="[0-9]{6}"
              maxLength={6}
              className="auth-page__input"
              value={code}
              onChange={(event) => setCode(event.target.value.replace(/\D/g, '').slice(0, 6))}
              required
              disabled={submitting}
            />
          </label>

          <button type="submit" className="btn btn--primary" disabled={submitting}>
            {submitting ? t.authVerifySubmitting : t.authVerifySubmit}
          </button>
        </form>

        <button type="button" className="btn btn--secondary" onClick={handleResend} disabled={resending || resendRemaining > 0}>
          {resending ? t.authVerifyResending : resendRemaining > 0 ? `${t.authVerifyResendIn} ${resendRemaining}s` : t.authVerifyResend}
        </button>

        <p className="auth-page__alt">
          <Link to="/login">{t.authVerifyBackLogin}</Link>
        </p>
      </article>
    </section>
  );
}
