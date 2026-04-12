import { useState, type FormEvent } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { useSession } from '../../app/providers/SessionProvider';
import { login as apiLogin } from './api';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';

interface LoginLocationState {
  from?: string;
}

export function LoginPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const { login } = useSession();
  const state = location.state as LoginLocationState | null;

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setSubmitting(true);

    try {
      const response = await apiLogin({ email, password });
      login(response);
      navigate(state?.from ?? '/admin/projects', { replace: true });
    } catch (err) {
      if (err instanceof AppError) {
        if (err.code === API_ERROR_CODES.INVALID_CREDENTIALS) {
          setError('Invalid email or password.');
        } else if (err.code === API_ERROR_CODES.VALIDATION_ERROR) {
          setError('Please review the submitted fields.');
        } else {
          setError(err.message);
        }
      } else {
        setError('An unexpected error occurred. Please try again.');
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="auth-page">
      <article className="card">
        <p className="eyebrow">Admin sign in</p>
        <h2>Manage PortfolioForge</h2>

        {state?.from ? (
          <p className="auth-page__redirect-note">You need to sign in to access the admin console.</p>
        ) : null}

        {error ? <div className="auth-page__error" role="alert">{error}</div> : null}

        <form className="auth-page__form" onSubmit={handleSubmit}>
          <label className="auth-page__label">
            Email
            <input
              type="email"
              className="auth-page__input"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
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
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete="current-password"
              disabled={submitting}
            />
          </label>

          <button type="submit" className="btn btn--primary" disabled={submitting}>
            {submitting ? 'Signing in…' : 'Sign in'}
          </button>
        </form>

        <p className="auth-page__alt">
          <Link to="/">Back to portfolio</Link>
        </p>
      </article>
    </section>
  );
}
