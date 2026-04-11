import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import { fetchAdminTechnologies, deleteTechnology } from './api';
import type { Technology } from '../../shared/types/project';
import { AppError } from '../../shared/api/errors';

export function AdminTechnologyListPage() {
  const [technologies, setTechnologies] = useState<Technology[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const loadTechnologies = () => {
    let cancelled = false;
    setLoading(true);
    setError(null);

    fetchAdminTechnologies()
      .then((response) => {
        if (!cancelled) {
          setTechnologies(response.items);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Failed to load technologies.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  };

  useEffect(loadTechnologies, []);

  async function handleDelete(technology: Technology) {
    if (!window.confirm(`Delete technology "${technology.name}"?`)) return;

    setDeletingId(technology.id);

    try {
      await deleteTechnology(technology.id);
      setTechnologies((prev) => prev.filter((item) => item.id !== technology.id));
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : 'Failed to delete technology.');
    } finally {
      setDeletingId(null);
    }
  }

  if (loading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading technologies…</p>
      </section>
    );
  }

  if (error) {
    return (
      <section className="card-stack">
        <article className="card">
          <div className="admin__error" role="alert">{error}</div>
          <button className="btn btn--ghost" type="button" onClick={loadTechnologies}>
            Retry
          </button>
        </article>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <div className="admin__header">
          <div>
            <p className="eyebrow">Admin</p>
            <h2>Technologies</h2>
          </div>
          <div className="admin__header-actions">
            <Link className="btn btn--primary" to="/admin/technologies/new">
              New Technology
            </Link>
          </div>
        </div>
      </article>

      {technologies.length === 0 ? (
        <article className="card card--muted">
          <p>No technologies found yet. Create the first one.</p>
        </article>
      ) : (
        <div className="admin__table-wrap">
          <table className="admin__table">
            <thead>
              <tr>
                <th>Technology</th>
                <th>Category</th>
                <th>Color</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {technologies.map((tech) => (
                <tr key={tech.id}>
                  <td>
                    <Link className="admin__link" to={`/admin/technologies/${tech.id}`}>
                      {tech.name}
                    </Link>
                  </td>
                  <td>{tech.category}</td>
                  <td>
                    {tech.color ? (
                      <span style={{ color: tech.color }}>{tech.color}</span>
                    ) : (
                      '—'
                    )}
                  </td>
                  <td className="admin__actions">
                    <Link className="btn btn--small btn--ghost" to={`/admin/technologies/${tech.id}`}>
                      Edit
                    </Link>
                    <button
                      className="btn btn--small btn--primary"
                      type="button"
                      disabled={deletingId === tech.id}
                      onClick={() => handleDelete(tech)}
                    >
                      {deletingId === tech.id ? '…' : 'Delete'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}