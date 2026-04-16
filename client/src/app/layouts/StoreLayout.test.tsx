import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../providers/LocaleProvider';
import { SessionProvider } from '../providers/SessionProvider';
import { StoreLayout } from './StoreLayout';
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

    expect(screen.getByRole('link', { name: 'Login' })).toHaveAttribute('href', '/login');
    expect(screen.queryByRole('link', { name: 'Admin' })).not.toBeInTheDocument();

    await waitFor(() => {
      expect(mockedFetchPublicSiteSettings).toHaveBeenCalled();
    });
  });
});
