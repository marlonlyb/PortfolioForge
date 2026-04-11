import { NavLink, Outlet, useNavigate } from 'react-router-dom';

import { useSession } from '../providers/SessionProvider';

const accountLinks = [
  { to: '/profile', label: 'Profile' },
  { to: '/profile/orders', label: 'Orders' },
  { to: '/checkout', label: 'Checkout' },
] as const;

export function AccountLayout() {
  const { user, logout } = useSession();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/', { replace: true });
  };

  return (
    <div className="app-shell">
      <header className="app-header">
        <div>
          <p className="eyebrow">Protected area</p>
          <h1>Customer account shell</h1>
        </div>

        <nav className="nav-list" aria-label="Account">
          {accountLinks.map((link) => (
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

          {user ? (
            <span className="nav-link nav-link--user">{user.email}</span>
          ) : null}

          <button className="nav-link nav-link--logout" onClick={handleLogout} type="button">
            Logout
          </button>
        </nav>
      </header>

      <main className="app-content">
        <Outlet />
      </main>
    </div>
  );
}
