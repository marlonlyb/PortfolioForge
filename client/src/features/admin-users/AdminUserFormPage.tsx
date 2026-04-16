import { useEffect, useState, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import { deleteAdminUser, fetchAdminUserById, updateAdminUser } from './api';
import { useSession } from '../../app/providers/SessionProvider';
import { AppError } from '../../shared/api/errors';
import type { AdminUserDetail } from '../../shared/types/admin-user';

function formatDate(value?: string): string {
  if (!value) {
    return '—';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return '—';
  }

  return date.toLocaleString();
}

export function AdminUserFormPage() {
  const { id } = useParams<{ id: string }>();
  const { user: currentUser } = useSession();
  const navigate = useNavigate();

  const [user, setUser] = useState<AdminUserDetail | null>(null);
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) {
      setError('Missing user identifier.');
      setLoading(false);
      return;
    }

    let cancelled = false;
    setLoading(true);
    setError(null);

    fetchAdminUserById(id)
      .then((response) => {
        if (!cancelled) {
          setUser(response);
          setIsAdmin(response.is_admin);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Failed to load user.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [id]);

  const isSelf = currentUser?.id === user?.id;
  const isProtectedAdmin = Boolean(user?.is_admin);
  const deleteDisabled = isProtectedAdmin || isSelf || deleting;
  const saveDisabled = submitting || user === null || isProtectedAdmin;
  const restrictionMessage = isProtectedAdmin
    ? 'Existing admins are read-only in this flow.'
    : isSelf
      ? 'Self-management is out of scope for this admin flow.'
      : null;

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (!id || !user || isProtectedAdmin) {
      return;
    }

    setSubmitting(true);
    setError(null);

    try {
      const updatedUser = await updateAdminUser(id, { is_admin: isAdmin });
      setUser(updatedUser);
      setIsAdmin(updatedUser.is_admin);
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : 'Failed to update user.');
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete() {
    if (!id || !user || deleteDisabled) {
      return;
    }
    if (!window.confirm(`Soft-delete user "${user.email}"?`)) {
      return;
    }

    setDeleting(true);
    setError(null);

    try {
      await deleteAdminUser(id);
      navigate('/admin/users', { replace: true });
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : 'Failed to delete user.');
      setDeleting(false);
    }
  }

  if (loading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading user…</p>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <Link className="detail__back" to="/admin/users">
          ← Back to users
        </Link>

        <p className="eyebrow">Admin</p>
        <h2>User detail</h2>

        {error ? <div className="admin__error" role="alert">{error}</div> : null}
        {restrictionMessage ? <p className="admin__helper-copy">{restrictionMessage}</p> : null}

        {user ? (
          <form className="admin__form" onSubmit={handleSubmit}>
            <div className="admin__form-section">
              <h3>Identity</h3>

              <label className="admin__label">
                Email
                <input className="admin__input" type="text" value={user.email} readOnly />
              </label>

              <div className="admin__form-row">
                <label className="admin__label">
                  Full name
                  <input className="admin__input" type="text" value={user.full_name || '—'} readOnly />
                </label>

                <label className="admin__label">
                  Company
                  <input className="admin__input" type="text" value={user.company || '—'} readOnly />
                </label>
              </div>

              <div className="admin__form-row">
                <label className="admin__label">
                  Provider
                  <input className="admin__input" type="text" value={user.auth_provider} readOnly />
                </label>

                <label className="admin__label">
                  Email verified
                  <input className="admin__input" type="text" value={user.email_verified ? 'Yes' : 'No'} readOnly />
                </label>
              </div>

              <div className="admin__form-row">
                <label className="admin__label">
                  Created at
                  <input className="admin__input" type="text" value={formatDate(user.created_at)} readOnly />
                </label>

                <label className="admin__label">
                  Last login
                  <input className="admin__input" type="text" value={formatDate(user.last_login_at)} readOnly />
                </label>
              </div>
            </div>

            <div className="admin__form-section">
              <h3>Permissions</h3>

              <label className="admin__label admin__label--checkbox">
                <input
                  checked={isAdmin}
                  disabled={isProtectedAdmin || submitting}
                  type="checkbox"
                  onChange={(event) => setIsAdmin(event.target.checked)}
                />
                Grant admin access
              </label>

              <p className="admin__helper-copy">
                This flow only mutates <code>is_admin</code>. Profile and identity fields are read-only.
              </p>
            </div>

            <div className="admin__form-actions">
              <button className="btn btn--primary" type="submit" disabled={saveDisabled}>
                {submitting ? 'Saving…' : 'Save changes'}
              </button>
              <button className="btn btn--ghost" type="button" disabled={deleteDisabled} onClick={handleDelete}>
                {deleting ? 'Deleting…' : 'Delete user'}
              </button>
            </div>
          </form>
        ) : null}
      </article>
    </section>
  );
}
