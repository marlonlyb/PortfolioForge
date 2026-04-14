import { NavLink, Outlet, useNavigate } from 'react-router-dom';

import { useSession } from '../providers/SessionProvider';

const adminLinks = [
  { to: '/admin/projects', label: 'Projects' },
  { to: '/admin/technologies', label: 'Technologies' },
  { to: '/admin/settings', label: 'Settings' },
] as const;

export function AdminLayout() {
  const { user, logout } = useSession();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/', { replace: true });
  };

  return (
    <div className="app-shell">
      <header className="app-header admin-header">
        <div className="admin-header__top-row">
          <p className="eyebrow">PortfolioForge</p>
          {user ? <span className="admin-header__user">{user.email}</span> : null}
        </div>

        <div className="admin-header__main-row">
          <h1>Admin console</h1>

          <nav className="nav-list" aria-label="Admin">
            {adminLinks.map((link) => (
              <NavLink
                key={link.to}
                className={({ isActive }) =>
                  isActive ? 'nav-link nav-link--active' : 'nav-link'
                }
                to={link.to}
              >
                {link.label}
              </NavLink>
            ))}

            <NavLink
              className={({ isActive }) =>
                isActive ? 'nav-link nav-link--active' : 'nav-link'
              }
              to="/"
            >
              Portfolio
            </NavLink>

            <button className="nav-link nav-link--logout" onClick={handleLogout} type="button">
              Logout
            </button>
          </nav>
        </div>
      </header>

      <main className="app-content">
        <Outlet />
      </main>
    </div>
  );
}
