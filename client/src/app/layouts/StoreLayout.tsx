import { useEffect, useState } from 'react';
import { NavLink, Outlet, useLocation } from 'react-router-dom';

import { useSession } from '../providers/SessionProvider';
import { useLocale } from '../providers/LocaleProvider';
import { PUBLIC_LOCALE_LABELS, type PublicLocale } from '../../shared/i18n/config';
import { fetchPublicSiteSettings } from '../../shared/api/siteSettings';
import type { SiteSettings } from '../../shared/types/siteSettings';

export interface StoreHeaderContent {
  title: string;
  summary: string;
  caption: string;
}

export interface StoreLayoutOutletContext {
  setHeaderContent: (value: StoreHeaderContent | null) => void;
}

const PROJECT_DETAIL_ROUTE_PATTERN = /^\/projects\/[^/]+\/?$/;

export function StoreLayout() {
  const { user, logout } = useSession();
  const { locale, setLocale, t } = useLocale();
  const location = useLocation();
  const isLanding = location.pathname === '/' || location.pathname.startsWith('/projects/');
  const isEditorialLanding = location.pathname === '/';
  const isProjectDetailRoute = PROJECT_DETAIL_ROUTE_PATTERN.test(location.pathname);
  const [siteSettings, setSiteSettings] = useState<SiteSettings | null>(null);
  const [contextualHeader, setContextualHeader] = useState<StoreHeaderContent | null>(null);

  const defaultHeader: StoreHeaderContent = {
    title: t.headerTitle,
    summary: t.headerSummary,
    caption: t.headerCaption,
  };
  const hasContextualDetailHeader = isProjectDetailRoute && contextualHeader !== null;
  const activeHeader = hasContextualDetailHeader ? contextualHeader : defaultHeader;
  const showHeaderSummary = activeHeader.summary.trim().length > 0;
  const showHeaderCaption = activeHeader.caption.trim().length > 0;
  const headerClassName = hasContextualDetailHeader
    ? 'app-header app-header--store app-header--detail-compact'
    : 'app-header app-header--store';

  useEffect(() => {
    window.scrollTo({ top: 0, left: 0, behavior: 'auto' });
  }, [location.pathname]);

  useEffect(() => {
    if (!isEditorialLanding) {
      setSiteSettings(null);
      return;
    }

    let cancelled = false;

    fetchPublicSiteSettings()
      .then((settings) => {
        if (!cancelled) {
          setSiteSettings(settings);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setSiteSettings(null);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [isEditorialLanding]);

  useEffect(() => {
    if (!isProjectDetailRoute) {
      setContextualHeader(null);
    }
  }, [isProjectDetailRoute]);

  const heroLogoUrl = siteSettings?.public_hero_logo_url?.trim();
  const heroLogoAlt = siteSettings?.public_hero_logo_alt?.trim() || 'Public portfolio logo';

  function renderPrimaryNav() {
    return (
      <nav className="nav-list" aria-label="Primary">
        <NavLink
          className={({ isActive }) =>
            isActive ? 'nav-link nav-link--active' : 'nav-link'
          }
          to="/"
          end
        >
          {t.navHome}
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
                {t.navAdmin}
              </NavLink>
            )}
            <button className="nav-link nav-link--logout" onClick={logout} type="button">
              {t.navLogout}
            </button>
          </>
        ) : (
          <NavLink
            className={({ isActive }) =>
              isActive ? 'nav-link nav-link--active' : 'nav-link'
            }
            to="/login"
          >
            {t.navLogin}
          </NavLink>
        )}
      </nav>
    );
  }

  return (
    <div className={isLanding ? 'app-shell app-shell--landing' : 'app-shell'}>
      {isEditorialLanding ? (
        <div className="landing-composition">
          <div className="app-header__brand landing-composition__brand">
            <p className="app-header__caption app-header__caption--landing">{t.headerCaption}</p>
            <NavLink className="app-header__home" to="/">
              <h1>{t.headerTitle}</h1>
            </NavLink>
            <p className="app-header__summary">{t.headerSummary}</p>
            <div className="locale-switcher" aria-label="Language selector">
              {(Object.keys(PUBLIC_LOCALE_LABELS) as PublicLocale[]).map((option) => (
                <button
                  key={option}
                  className={option === locale ? 'locale-switcher__button locale-switcher__button--active' : 'locale-switcher__button'}
                  type="button"
                  onClick={() => setLocale(option)}
                >
                  {PUBLIC_LOCALE_LABELS[option]}
                </button>
              ))}
            </div>
          </div>

          <div className="app-header__actions landing-composition__actions">
            <div className="app-header__toolbar">
              {renderPrimaryNav()}
            </div>
          </div>

          <aside className="card landing-hero landing-hero--logo landing-composition__logo" aria-label="Brand logo slot">
            <div className="landing-hero__logo-slot landing-hero__logo-slot--standalone">
              {heroLogoUrl ? (
                <img className="landing-hero__logo-image" src={heroLogoUrl} alt={heroLogoAlt} loading="lazy" />
              ) : (
                <>
                  <div className="landing-hero__logo-badge">PF</div>
                  <div className="landing-hero__logo-copy">
                    <strong>Logo slot</strong>
                    <span>Prepared for your brand image or personal mark.</span>
                  </div>
                </>
              )}
            </div>
          </aside>

          <Outlet context={{ setHeaderContent: setContextualHeader }} />
        </div>
      ) : (
        <>
          <header className={headerClassName}>
            <div className="app-header__brand">
              <NavLink className="app-header__home" to="/">
                <h1>{activeHeader.title}</h1>
              </NavLink>
              {showHeaderSummary ? <p className="app-header__summary">{activeHeader.summary}</p> : null}
              {hasContextualDetailHeader && showHeaderCaption ? (
                <p className="app-header__caption app-header__caption--detail">{activeHeader.caption}</p>
              ) : null}
            </div>

            <div className="app-header__actions">
              <div className="app-header__toolbar">
                <div className="locale-switcher" aria-label="Language selector">
                  {(Object.keys(PUBLIC_LOCALE_LABELS) as PublicLocale[]).map((option) => (
                    <button
                      key={option}
                      className={option === locale ? 'locale-switcher__button locale-switcher__button--active' : 'locale-switcher__button'}
                      type="button"
                      onClick={() => setLocale(option)}
                    >
                      {PUBLIC_LOCALE_LABELS[option]}
                    </button>
                  ))}
                </div>

                {renderPrimaryNav()}
              </div>

              {!hasContextualDetailHeader && showHeaderCaption ? (
                <p className="app-header__caption">{activeHeader.caption}</p>
              ) : null}
            </div>
          </header>

          <main className={isLanding ? 'app-content app-content--landing' : 'app-content'}>
            <Outlet context={{ setHeaderContent: setContextualHeader }} />
          </main>
        </>
      )}
    </div>
  );
}
