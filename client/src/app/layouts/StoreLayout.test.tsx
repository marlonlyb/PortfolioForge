import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { useEffect } from 'react';
import { MemoryRouter, Navigate, Route, Routes, useOutletContext } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../providers/LocaleProvider';
import { SessionProvider } from '../providers/SessionProvider';
import { StoreLayout, type StoreLayoutOutletContext } from './StoreLayout';
import { fetchPublicSiteSettings } from '../../shared/api/siteSettings';
import type { SiteSettings } from '../../shared/types/siteSettings';

vi.mock('../../shared/api/siteSettings', async () => {
  const actual = await vi.importActual<typeof import('../../shared/api/siteSettings')>('../../shared/api/siteSettings');
  return {
    ...actual,
    fetchPublicSiteSettings: vi.fn(),
  };
});

const mockedFetchPublicSiteSettings = vi.mocked(fetchPublicSiteSettings);
const defaultSiteSettings: SiteSettings = {};

function DetailHeaderProbe() {
  const { setHeaderContent } = useOutletContext<StoreLayoutOutletContext>();

  useEffect(() => {
    setHeaderContent({
      title: 'PortfolioForge',
      summary: '',
      caption: 'platform · Analytical Engines',
    });

    return () => setHeaderContent(null);
  }, [setHeaderContent]);

  return <div>detail content</div>;
}

describe('StoreLayout', () => {
  beforeEach(() => {
    mockedFetchPublicSiteSettings.mockReset();
    mockedFetchPublicSiteSettings.mockResolvedValue(defaultSiteSettings);
    window.localStorage.clear();
    window.sessionStorage.clear();
    window.localStorage.setItem('portfolioforge.locale', 'en');
    window.scrollTo = vi.fn();
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });

  it('shows only the public Login entry for signed-out visitors', async () => {
    render(
      <MemoryRouter initialEntries={['/']}>
        <SessionProvider>
          <LocaleProvider>
            <Routes>
              <Route path="/" element={<StoreLayout />}>
                <Route index element={<div>landing content</div>} />
              </Route>
            </Routes>
          </LocaleProvider>
        </SessionProvider>
      </MemoryRouter>,
    );

    const primaryNav = screen.getByRole('navigation', { name: 'Primary' });

    expect(screen.getByRole('link', { name: 'Login' })).toHaveAttribute('href', '/login');
    expect(screen.queryByRole('link', { name: 'Sign up' })).not.toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Admin' })).not.toBeInTheDocument();
    expect(primaryNav).not.toHaveTextContent('Sign up');

    await waitFor(() => {
      expect(mockedFetchPublicSiteSettings).toHaveBeenCalled();
    });
  });

  it('keeps the generic header copy on non-detail routes after adding outlet context', () => {
    render(
      <MemoryRouter initialEntries={['/search']}>
        <SessionProvider>
          <LocaleProvider>
            <Routes>
              <Route path="/" element={<StoreLayout />}>
                <Route path="search" element={<div>search content</div>} />
              </Route>
            </Routes>
          </LocaleProvider>
        </SessionProvider>
      </MemoryRouter>,
    );

    expect(screen.getByRole('heading', { level: 1, name: 'Project portfolio' })).toBeInTheDocument();
    expect(screen.getByText('Strategy, execution, and technical judgment.')).toBeInTheDocument();
    expect(screen.getByText('Marlon Ly Bellido · Engineer')).toBeInTheDocument();
    expect(screen.getByText('search content')).toBeInTheDocument();
    expect(mockedFetchPublicSiteSettings).not.toHaveBeenCalled();
  });

  it('keeps the generic header copy after the /projects redirect resolves to landing', async () => {
    render(
      <MemoryRouter initialEntries={['/projects']}>
        <SessionProvider>
          <LocaleProvider>
            <Routes>
              <Route path="/" element={<StoreLayout />}>
                <Route index element={<div>landing content</div>} />
                <Route path="projects" element={<Navigate replace to="/" />} />
              </Route>
            </Routes>
          </LocaleProvider>
        </SessionProvider>
      </MemoryRouter>,
    );

    expect(await screen.findByText('landing content')).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 1, name: 'Project portfolio' })).toBeInTheDocument();
    expect(screen.getByText('Strategy, execution, and technical judgment.')).toBeInTheDocument();
    expect(screen.getByText('Marlon Ly Bellido · Engineer')).toBeInTheDocument();
  });

  it('renders compact contextual detail headers without summary copy', async () => {
    render(
      <MemoryRouter initialEntries={['/projects/portfolioforge']}>
        <SessionProvider>
          <LocaleProvider>
            <Routes>
              <Route path="/" element={<StoreLayout />}>
                <Route path="projects/:slug" element={<DetailHeaderProbe />} />
              </Route>
            </Routes>
          </LocaleProvider>
        </SessionProvider>
      </MemoryRouter>,
    );

    expect(await screen.findByText('detail content')).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
    expect(screen.getByText('platform · Analytical Engines')).toBeInTheDocument();
    expect(screen.queryByText('Strategy, execution, and technical judgment.')).not.toBeInTheDocument();
    expect(document.querySelector('.app-header--detail-compact')).not.toBeNull();
    expect(document.querySelector('.app-header__summary')).toBeNull();
  });
});
