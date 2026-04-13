import { useEffect } from 'react';
import { NavLink, Outlet, useLocation } from 'react-router-dom';

import { useSession } from '../providers/SessionProvider';

export function StoreLayout() {
  const { user, logout } = useSession();
  const location = useLocation();
  const isLanding = location.pathname === '/' || location.pathname.startsWith('/projects/');

  useEffect(() => {
    window.scrollTo({ top: 0, left: 0, behavior: 'auto' });
  }, [location.pathname]);

  return (
    <div className={isLanding ? 'app-shell app-shell--landing' : 'app-shell'}>
      <header className="app-header">
        <div>
          <p className="eyebrow">PortfolioForge</p>
          <h1>Interactive project portfolio</h1>
        </div>

        <nav className="nav-list" aria-label="Primary">
          <NavLink
            className={({ isActive }) =>
              isActive ? 'nav-link nav-link--active' : 'nav-link'
            }
            to="/"
            end
          >
            Home
          </NavLink>

          {user ? (
            <>
              {user.is_admin && (
                <NavLink
                  className={({ isActive }) =>
                    isActive ? 'nav-link nav-link--active' : 'nav-link'
                  }
                  to="/admin/projects"
                >
                  Admin
                </NavLink>
              )}
              <button className="nav-link nav-link--logout" onClick={logout} type="button">
                Logout
              </button>
            </>
          ) : (
            <>
              <NavLink
                className={({ isActive }) =>
                  isActive ? 'nav-link nav-link--active' : 'nav-link'
                }
                to="/login"
              >
                Admin
              </NavLink>
            </>
          )}
        </nav>
      </header>

      <main className="app-content">
        <Outlet />
      </main>
    </div>
  );
}
