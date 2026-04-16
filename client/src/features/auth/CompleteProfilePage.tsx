import { useState, type FormEvent } from 'react';
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom';

import { useSession } from '../../app/providers/SessionProvider';
import { AppError } from '../../shared/api/errors';
import { updateMyProfile } from './api';

interface CompleteProfileLocationState {
  from?: string;
}

export function CompleteProfilePage() {
  const { token, user, refreshSession, setUser } = useSession();
  const location = useLocation();
  const navigate = useNavigate();
  const state = location.state as CompleteProfileLocationState | null;

  const [fullName, setFullName] = useState(user?.full_name ?? '');
  const [company, setCompany] = useState(user?.company ?? '');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  if (!token) {
    return <Navigate replace to="/login" state={{ from: state?.from ?? '/' }} />;
  }

  if (user?.profile_completed) {
    return <Navigate replace to={state?.from ?? '/'} />;
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      const response = await updateMyProfile({ full_name: fullName, company });
      setUser(response.user);
      await refreshSession();
      navigate(state?.from ?? '/', { replace: true });
    } catch (err) {
      if (err instanceof AppError) {
        setError(err.message);
      } else {
        setError('Unable to save your profile right now.');
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="auth-page">
      <article className="card">
        <p className="eyebrow">Complete your profile</p>
        <h2>Unlock the project assistant</h2>
        <p className="auth-page__redirect-note">Add your full name and company to continue with project-specific chat.</p>

        {error ? <div className="auth-page__error" role="alert">{error}</div> : null}

        <form className="auth-page__form" onSubmit={handleSubmit}>
          <label className="auth-page__label">
            Full name
            <input
              type="text"
              className="auth-page__input"
              value={fullName}
              onChange={(event) => setFullName(event.target.value)}
              required
              disabled={submitting}
              autoComplete="name"
            />
          </label>

          <label className="auth-page__label">
            Company
            <input
              type="text"
              className="auth-page__input"
              value={company}
              onChange={(event) => setCompany(event.target.value)}
              required
              disabled={submitting}
              autoComplete="organization"
            />
          </label>

          <button type="submit" className="btn btn--primary" disabled={submitting}>
            {submitting ? 'Saving…' : 'Save profile'}
          </button>
        </form>

        <p className="auth-page__alt">
          <Link to={state?.from ?? '/'}>Back to project</Link>
        </p>
      </article>
    </section>
  );
}
