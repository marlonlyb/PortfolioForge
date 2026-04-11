import { Navigate, Outlet } from 'react-router-dom';

import { useSession } from '../../app/providers/SessionProvider';

/**
 * Guard that verifies the authenticated user has admin privileges.
 * Must be rendered INSIDE a RequireAuth subtree so `user` is guaranteed.
 */
export function RequireAdmin() {
  const { user, loading } = useSession();

  if (loading) {
    return null;
  }

  if (!user?.is_admin) {
    return <Navigate replace to="/" />;
  }

  return <Outlet />;
}
