import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import { fetchAdminProducts, updateProductStatus, reembedStale } from './api';
import type { ProductDetail } from '../../shared/types/product';
import { AppError } from '../../shared/api/errors';

export function AdminProductListPage() {
  const [projects, setProjects] = useState<ProductDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [togglingId, setTogglingId] = useState<string | null>(null);
  const [reembedLoading, setReembedLoading] = useState(false);
  const [reembedMessage, setReembedMessage] = useState<string | null>(null);

  const loadProjects = () => {
    let cancelled = false;
    setLoading(true);
    setError(null);

    fetchAdminProducts()
      .then((response) => {
        if (!cancelled) {
          setProjects(response.items);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Failed to load projects.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  };

  useEffect(loadProjects, []);

  async function handleToggleStatus(project: ProductDetail) {
    setTogglingId(project.id);

    try {
      await updateProductStatus(project.id, { active: !project.active });
      setProjects((prev) =>
        prev.map((item) => (item.id === project.id ? { ...item, active: !item.active } : item)),
      );
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : 'Failed to update project status.');
    } finally {
      setTogglingId(null);
    }
  }

  async function handleBatchReembed() {
    if (!window.confirm('¿Actualizar los documentos de búsqueda de todos los proyectos?')) return;

    setReembedLoading(true);
    setReembedMessage(null);

    try {
      const result = await reembedStale();
      setReembedMessage(result.message);
    } catch (err: unknown) {
      setReembedMessage(
        err instanceof AppError ? err.message : 'Error al actualizar documentos de búsqueda.',
      );
    } finally {
      setReembedLoading(false);
    }
  }

  if (loading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading projects…</p>
      </section>
    );
  }

  if (error) {
    return (
      <section className="card-stack">
        <article className="card">
          <div className="admin__error" role="alert">{error}</div>
          <button className="btn btn--ghost" type="button" onClick={loadProjects}>
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
            <h2>Projects</h2>
          </div>
          <div className="admin__header-actions">
            <button
              className="btn btn--ghost"
              type="button"
              disabled={reembedLoading}
              onClick={handleBatchReembed}
            >
              {reembedLoading ? 'Actualizando búsqueda…' : 'Actualizar búsqueda general'}
            </button>
            <Link className="btn btn--primary" to="/admin/projects/new">
              New Project
            </Link>
          </div>
        </div>

        {reembedMessage && (
          <p className="admin__reembed-msg" role="status">{reembedMessage}</p>
        )}
      </article>

      {projects.length === 0 ? (
        <article className="card card--muted">
          <p>No projects found yet. Create the first one.</p>
        </article>
      ) : (
        <div className="admin__table-wrap">
          <table className="admin__table">
            <thead>
              <tr>
                <th>Project</th>
                <th>Category</th>
                <th>Client / Context</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {projects.map((project) => (
                <tr key={project.id}>
                  <td>
                    <Link className="admin__link" to={`/admin/projects/${project.id}`}>
                      {project.name}
                    </Link>
                  </td>
                  <td>{project.category}</td>
                  <td>{project.brand || '—'}</td>
                  <td>
                    <span
                      className={`admin__badge ${project.active ? 'admin__badge--active' : 'admin__badge--inactive'}`}
                    >
                      {project.active ? 'Published' : 'Draft'}
                    </span>
                  </td>
                  <td className="admin__actions">
                    <Link className="btn btn--small btn--ghost" to={`/admin/projects/${project.id}`}>
                      Edit
                    </Link>
                    <button
                      className={`btn btn--small ${project.active ? 'btn--ghost' : 'btn--primary'}`}
                      type="button"
                      disabled={togglingId === project.id}
                      onClick={() => handleToggleStatus(project)}
                    >
                      {togglingId === project.id ? '…' : project.active ? 'Unpublish' : 'Publish'}
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
