import { useEffect, useState, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import {
  createTechnology,
  fetchAdminTechnologyById,
  updateTechnology,
} from './api';
import type { CreateTechnologyPayload, UpdateTechnologyPayload } from './api';
import { AppError } from '../../shared/api/errors';

export function AdminTechnologyFormPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEdit = Boolean(id);

  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [name, setName] = useState('');
  const [category, setCategory] = useState('');
  const [icon, setIcon] = useState('');
  const [color, setColor] = useState('');

  useEffect(() => {
    if (!id) return;

    let cancelled = false;

    fetchAdminTechnologyById(id)
      .then((tech) => {
        if (!cancelled) {
          setName(tech.name);
          setCategory(tech.category);
          setIcon(tech.icon || '');
          setColor(tech.color || '');
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Failed to load technology.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [id]);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setSubmitting(true);

    try {
      if (isEdit && id) {
        const payload: UpdateTechnologyPayload = {
          name,
          category,
          icon: icon || undefined,
          color: color || undefined,
        };
        await updateTechnology(id, payload);
      } else {
        const payload: CreateTechnologyPayload = {
          name,
          category,
          icon: icon || undefined,
          color: color || undefined,
        };
        await createTechnology(payload);
      }

      navigate('/admin/technologies');
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : 'Failed to save technology.');
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading technology…</p>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <Link className="detail__back" to="/admin/technologies">
          ← Back to technologies
        </Link>

        <p className="eyebrow">Admin</p>
        <h2>{isEdit ? 'Edit Technology' : 'New Technology'}</h2>

        {error ? <div className="admin__error" role="alert">{error}</div> : null}

        <form className="admin__form" onSubmit={handleSubmit}>
          <div className="admin__form-section">
            <h3>Technology Details</h3>

            <label className="admin__label">
              Name
              <input
                className="admin__input"
                type="text"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Category
              <input
                className="admin__input"
                type="text"
                required
                value={category}
                onChange={(e) => setCategory(e.target.value)}
              />
            </label>

            <div className="admin__form-row">
              <label className="admin__label">
                Icon (CSS class or SVG path)
                <input
                  className="admin__input"
                  type="text"
                  value={icon}
                  onChange={(e) => setIcon(e.target.value)}
                />
              </label>

              <label className="admin__label">
                Brand Color (Hex or CSS var)
                <input
                  className="admin__input"
                  type="text"
                  placeholder="#FFFFFF"
                  value={color}
                  onChange={(e) => setColor(e.target.value)}
                />
              </label>
            </div>
          </div>

          <div className="admin__form-actions">
            <button className="btn btn--primary" type="submit" disabled={submitting}>
              {submitting ? 'Saving…' : isEdit ? 'Update Technology' : 'Create Technology'}
            </button>
          </div>
        </form>
      </article>
    </section>
  );
}