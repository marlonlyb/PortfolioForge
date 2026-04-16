import { useEffect, useState, type FormEvent } from 'react';
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom';

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
  const { user, refreshSession, setUser } = useSession();
  const state = (location.state as VerifyEmailLocationState | null) ?? null;

  const [email, setEmail] = useState(state?.email ?? user?.email ?? '');
  const [code, setCode] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string>('Enter the 6-digit verification code we sent to your email.');
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
        state: { notice: 'Email verified. Log in with your password to continue.' },
      });
    } catch (err) {
      if (err instanceof AppError) {
        if (err.code === API_ERROR_CODES.OTP_INVALID) {
          setError('The verification code is invalid. Please try again.');
        } else if (err.code === API_ERROR_CODES.OTP_EXPIRED) {
          setError('This verification code expired. Request a new code to continue.');
        } else {
          setError(err.message);
        }
      } else {
        setError('Unable to verify the code right now.');
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
      setNotice(response.message);
      setResendRemaining(response.cooldown_seconds);
    } catch (err) {
      if (err instanceof AppError && err.code === API_ERROR_CODES.VALIDATION_ERROR) {
        setError(err.message);
      } else {
        setError('Unable to resend the verification code right now.');
      }
    } finally {
      setResending(false);
    }
  }

  return (
    <section className="auth-page">
      <article className="card">
        <p className="eyebrow">Email verification</p>
        <h2>Verify your email</h2>
        <p className="auth-page__redirect-note">{notice}</p>

        {error ? <div className="auth-page__error" role="alert">{error}</div> : null}

        <form className="auth-page__form" onSubmit={handleVerify}>
          <label className="auth-page__label">
            Email
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
            Verification code
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
            {submitting ? 'Verifying…' : 'Verify email'}
          </button>
        </form>

        <button type="button" className="btn btn--secondary" onClick={handleResend} disabled={resending || resendRemaining > 0}>
          {resending ? 'Sending…' : resendRemaining > 0 ? `Resend code in ${resendRemaining}s` : 'Resend code'}
        </button>

        <p className="auth-page__alt">
          <Link to="/login">Back to login</Link>
        </p>
      </article>
    </section>
  );
}
