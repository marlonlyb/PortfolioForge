import { useState, type FormEvent } from 'react';
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';
import { useSession } from '../../app/providers/SessionProvider';
import { AppError } from '../../shared/api/errors';
import { updateMyProfile } from './api';

interface CompleteProfileLocationState {
  from?: string;
}

export function CompleteProfilePage() {
  const { t } = useLocale();
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
        setError(t.authCompleteProfileUnableToSave);
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
      <section className="auth-page">
      <article className="card">
        <p className="eyebrow">{t.authCompleteProfileEyebrow}</p>
        <h2>{t.authCompleteProfileTitle}</h2>
        <p className="auth-page__redirect-note">{t.authCompleteProfileDescription}</p>

        {error ? <div className="auth-page__error" role="alert">{error}</div> : null}

        <form className="auth-page__form" onSubmit={handleSubmit}>
          <label className="auth-page__label">
            {t.authCompleteProfileFullName}
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
            {t.authCompleteProfileCompany}
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
            {submitting ? t.authCompleteProfileSaving : t.authCompleteProfileSave}
          </button>
        </form>

        <p className="auth-page__alt">
          <Link to={state?.from ?? '/'}>{t.authCompleteProfileBack}</Link>
        </p>
      </article>
    </section>
  );
}
