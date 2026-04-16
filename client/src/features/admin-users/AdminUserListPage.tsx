import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import { deleteAdminUser, fetchAdminUsers } from './api';
import { AppError } from '../../shared/api/errors';
import type { AdminUserSummary } from '../../shared/types/admin-user';
import { useSession } from '../../app/providers/SessionProvider';

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

export function AdminUserListPage() {
  const { user: currentUser } = useSession();
  const [users, setUsers] = useState<AdminUserSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const loadUsers = () => {
    let cancelled = false;
    setLoading(true);
    setError(null);

    fetchAdminUsers()
      .then((response) => {
        if (!cancelled) {
          setUsers(response.items);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Failed to load users.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  };

  useEffect(loadUsers, []);

  async function handleDelete(targetUser: AdminUserSummary) {
    if (targetUser.is_admin) {
      setError('Admin accounts cannot be deleted from this flow.');
      return;
    }
    if (currentUser?.id === targetUser.id) {
      setError('You cannot delete your own admin account.');
      return;
    }
    if (!window.confirm(`Soft-delete user "${targetUser.email}"?`)) {
      return;
    }

    setDeletingId(targetUser.id);
    setError(null);

    try {
      await deleteAdminUser(targetUser.id);
      setUsers((current) => current.filter((item) => item.id !== targetUser.id));
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : 'Failed to delete user.');
    } finally {
      setDeletingId(null);
    }
  }

  if (loading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading users…</p>
      </section>
    );
  }

  if (error) {
    return (
      <section className="card-stack">
        <article className="card">
          <div className="admin__error" role="alert">{error}</div>
          <button className="btn btn--ghost" type="button" onClick={loadUsers}>
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
            <h2>Users</h2>
          </div>
        </div>

        <p className="admin__helper-copy">
          Only active users are listed. Deleted identities remain reserved and lose access immediately.
        </p>
      </article>

      {users.length === 0 ? (
        <article className="card card--muted">
          <p>No active users found.</p>
        </article>
      ) : (
        <div className="admin__table-wrap">
          <table className="admin__table">
            <thead>
              <tr>
                <th>User</th>
                <th>Provider</th>
                <th>Role</th>
                <th>Last login</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((targetUser) => {
                const deleteBlocked = targetUser.is_admin || currentUser?.id === targetUser.id;
                const deleteBlockedReason = targetUser.is_admin
                  ? 'Admins cannot be deleted from this flow.'
                  : currentUser?.id === targetUser.id
                    ? 'You cannot delete your own admin account.'
                    : undefined;

                return (
                  <tr key={targetUser.id}>
                    <td>
                      <Link className="admin__link" to={`/admin/users/${targetUser.id}`}>
                        {targetUser.email}
                      </Link>
                      <div>{targetUser.full_name || '—'}</div>
                    </td>
                    <td>{targetUser.auth_provider}</td>
                    <td>
                      <span className={`admin__badge ${targetUser.is_admin ? 'admin__badge--active' : 'admin__badge--inactive'}`}>
                        {targetUser.is_admin ? 'Admin' : 'Standard'}
                      </span>
                    </td>
                    <td>{formatDate(targetUser.last_login_at)}</td>
                    <td className="admin__actions">
                      <Link className="btn btn--small btn--ghost" to={`/admin/users/${targetUser.id}`}>
                        View
                      </Link>
                      <button
                        className="btn btn--small btn--primary"
                        type="button"
                        disabled={deleteBlocked || deletingId === targetUser.id}
                        onClick={() => handleDelete(targetUser)}
                        title={deleteBlockedReason}
                      >
                        {deletingId === targetUser.id ? '…' : 'Delete'}
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
