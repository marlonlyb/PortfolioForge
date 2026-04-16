import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useSession } from '../../app/providers/SessionProvider';

/**
 * Guard that verifies the authenticated user has admin privileges.
 * Must be rendered INSIDE a RequireAuth subtree so `user` is guaranteed.
 */
export function RequireAdmin() {
  const { user, loading } = useSession();
  const location = useLocation();

  if (loading) {
    return null;
  }

  if (!user) {
    return <Navigate replace to="/login" state={{ from: `${location.pathname}${location.search}${location.hash}` }} />;
  }

  if (!user.is_admin) {
    return <Navigate replace to="/" />;
  }

  return <Outlet />;
}
