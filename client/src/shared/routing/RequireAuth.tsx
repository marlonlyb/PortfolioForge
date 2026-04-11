import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useSession } from '../../app/providers/SessionProvider';

export function RequireAuth() {
  const { token, loading } = useSession();
  const location = useLocation();

  if (loading) {
    // Session is being restored — render nothing (or a spinner later)
    return null;
  }

  if (!token) {
    const returnTo = `${location.pathname}${location.search}${location.hash}`;
    return <Navigate replace to="/login" state={{ from: returnTo }} />;
  }

  return <Outlet />;
}
