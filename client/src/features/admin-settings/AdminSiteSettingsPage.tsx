import { useEffect, useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';

import { AppError } from '../../shared/api/errors';
import { fetchAdminSiteSettings, updateAdminSiteSettings } from '../../shared/api/siteSettings';

export function AdminSiteSettingsPage() {
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [publicHeroLogoUrl, setPublicHeroLogoUrl] = useState('');
  const [publicHeroLogoAlt, setPublicHeroLogoAlt] = useState('');

  useEffect(() => {
    let cancelled = false;

    fetchAdminSiteSettings()
      .then((settings) => {
        if (cancelled) return;
        setPublicHeroLogoUrl(settings.public_hero_logo_url ?? '');
        setPublicHeroLogoAlt(settings.public_hero_logo_alt ?? '');
        setLoading(false);
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        setError(err instanceof AppError ? err.message : 'Failed to load site settings.');
        setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setMessage(null);

    try {
      const settings = await updateAdminSiteSettings({
        public_hero_logo_url: publicHeroLogoUrl,
        public_hero_logo_alt: publicHeroLogoAlt,
      });

      setPublicHeroLogoUrl(settings.public_hero_logo_url ?? '');
      setPublicHeroLogoAlt(settings.public_hero_logo_alt ?? '');
      setMessage('Public hero logo updated.');
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : 'Failed to save site settings.');
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading site settings…</p>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <p className="eyebrow">Settings hub</p>
        <h2>Case-study workflow</h2>
        <p className="admin__helper-copy">
          Publish an existing canonical case-study source, create/update the admin project,
          then run localization and re-embed with persisted status and logs.
        </p>

        <div className="admin__form-actions">
          <Link className="btn btn--secondary" to="/admin/settings/case-studies">
            Open workflow
          </Link>
        </div>
      </article>

      <article className="card">
        <p className="eyebrow">Admin</p>
        <h2>Public branding</h2>
        <p className="admin__helper-copy">
          Configure the image displayed in the logo slot of the public landing hero.
        </p>

        {error ? <div className="admin__error" role="alert">{error}</div> : null}
        {message ? <p className="admin__success" role="status">{message}</p> : null}

        <form className="admin__form" onSubmit={handleSubmit}>
          <div className="admin__form-section">
            <h3>Hero logo</h3>

            <label className="admin__label">
              Public logo URL
              <input
                className="admin__input"
                type="url"
                placeholder="https://cdn.example.com/brand/logo.svg"
                value={publicHeroLogoUrl}
                onChange={(event) => setPublicHeroLogoUrl(event.target.value)}
              />
            </label>

            <label className="admin__label">
              Alt text
              <input
                className="admin__input"
                type="text"
                placeholder="Portfolio logo"
                value={publicHeroLogoAlt}
                onChange={(event) => setPublicHeroLogoAlt(event.target.value)}
              />
            </label>
          </div>

          <div className="admin__form-actions">
            <button className="btn btn--primary" type="submit" disabled={submitting}>
              {submitting ? 'Saving…' : 'Save settings'}
            </button>
          </div>
        </form>
      </article>
    </section>
  );
}
