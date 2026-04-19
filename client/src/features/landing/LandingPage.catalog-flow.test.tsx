import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import type { Project } from '../../shared/types/project';
import type { SearchResult } from '../../shared/types/search';
import { fetchProjects } from '../catalog/api';
import { LandingPage } from './LandingPage';
import { searchProjects } from '../search/api';

vi.mock('../catalog/api', async () => {
  const actual = await vi.importActual<typeof import('../catalog/api')>('../catalog/api');
  return {
    ...actual,
    fetchProjects: vi.fn(),
  };
});

vi.mock('../search/api', async () => {
  const actual = await vi.importActual<typeof import('../search/api')>('../search/api');
  return {
    ...actual,
    searchProjects: vi.fn(),
  };
});

const mockedFetchProjects = vi.mocked(fetchProjects);
const mockedSearchProjects = vi.mocked(searchProjects);

function buildProject(overrides: Partial<Project> = {}): Project {
  return {
    id: 'project-1',
    name: 'Printer 05 PLC Migration',
    slug: 'printer-05-plc-migration',
    description: 'CompactLogix migration with Ethernet/IP commissioning.',
    category: 'automation',
    client_name: 'Printer 05',
    status: 'published',
    featured: true,
    active: true,
    assistant_available: true,
    images: [],
    media: [],
    created_at: 1710000000,
    updated_at: 1710000000,
    technologies: [
      { id: 'tech-1', name: 'CompactLogix', slug: 'compactlogix', category: 'plc' },
      { id: 'tech-2', name: 'Ethernet/IP', slug: 'ethernet-ip', category: 'network' },
    ],
    ...overrides,
  };
}

function buildSearchResult(overrides: Partial<SearchResult> = {}): SearchResult {
  return {
    id: 'result-1',
    slug: 'printer-05-plc-migration',
    title: 'Printer 05 PLC Migration',
    category: 'automation',
    client_name: 'Printer 05',
    summary: 'CompactLogix migration with Ethernet/IP commissioning.',
    technologies: [{ id: 'tech-1', name: 'CompactLogix', slug: 'compactlogix' }],
    hero_image: null,
    score: 0.98,
    explanation: null,
    evidence: [],
    ...overrides,
  };
}

describe('LandingPage catalog flow', () => {
  beforeEach(() => {
    window.localStorage.clear();
    window.localStorage.setItem('portfolioforge.locale', 'en');
    document.documentElement.lang = 'en';
    Object.defineProperty(HTMLElement.prototype, 'scrollIntoView', {
      configurable: true,
      value: vi.fn(),
    });
    mockedFetchProjects.mockReset();
    mockedSearchProjects.mockReset();
    mockedFetchProjects.mockResolvedValue({
      items: [buildProject()],
    });
    mockedSearchProjects.mockResolvedValue({
      data: [buildSearchResult()],
      meta: {
        total: 1,
        page_size: 1,
        cursor: null,
        query: 'Printer 05',
        filters_applied: {
          category: null,
          client: null,
          technologies: [],
        },
      },
    });
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });

  it('keeps landing prompt interactions on the existing catalog suggestion and result flow', async () => {
    render(
      <MemoryRouter>
        <LocaleProvider>
          <LandingPage />
        </LocaleProvider>
      </MemoryRouter>,
    );

    expect(await screen.findByRole('heading', { level: 2, name: 'Selected case studies' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'CompactLogix' }));

    expect(screen.getByRole('textbox', { name: 'Search' })).toHaveValue('CompactLogix');
    expect(HTMLElement.prototype.scrollIntoView).toHaveBeenCalledTimes(1);

    await waitFor(() => {
      expect(mockedSearchProjects).toHaveBeenCalledWith({
        q: 'CompactLogix',
        category: undefined,
        lang: 'en',
      });
    }, { timeout: 2000 });

    expect(screen.getByRole('link', { name: /Printer 05 PLC Migration/i })).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 3, name: 'Printer 05 PLC Migration' }).closest('a')).toHaveAttribute(
      'href',
      '/projects/printer-05-plc-migration',
    );
  });
});
