import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { useLocale } from '../../app/providers/LocaleProvider';
import type { SearchResponse } from '../../shared/types/search';
import { searchProjects } from './api';
import { SearchResultsPage } from './SearchResultsPage';

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
    searchProjects: vi.fn(),
  };
});

const mockedSearchProjects = vi.mocked(searchProjects);

function buildSearchResponse(overrides: Partial<SearchResponse> = {}): SearchResponse {
  return {
    data: [
      {
        id: 'project-1',
        slug: 'portfolioforge',
        title: 'PortfolioForge',
        category: 'platform',
        client_name: 'Acme',
        summary: 'Search-ready portfolio platform',
        technologies: [{ id: 'tech-1', name: 'React', slug: 'react', color: '#61dafb' }],
        hero_image: 'https://img/portfolioforge.png',
        score: 0.95,
        explanation: null,
        evidence: [],
      },
    ],
    meta: {
      total: 1,
      page_size: 10,
      cursor: null,
      query: 'portfolio',
      filters_applied: {
        category: 'platform',
        client: 'Acme',
        technologies: ['react'],
      },
    },
    ...overrides,
  };
}

function renderSearchPage(initialEntry = '/search?q=portfolio&category=platform&client=Acme&technologies=react') {
  function LocaleControls() {
    const { setLocale } = useLocale();
    return (
      <div>
        <button type="button" onClick={() => setLocale('es')}>locale-es</button>
        <button type="button" onClick={() => setLocale('en')}>locale-en</button>
      </div>
    );
  }

  function DetailLocationStateProbe() {
    const location = useLocation();
    return <pre data-testid="detail-location-state">{JSON.stringify(location.state)}</pre>;
  }

  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <LocaleProvider>
        <LocaleControls />
        <Routes>
          <Route path="/" element={<p>catalog destination</p>} />
          <Route path="/projects/:slug" element={<DetailLocationStateProbe />} />
          <Route path="/search" element={<SearchResultsPage />} />
        </Routes>
      </LocaleProvider>
    </MemoryRouter>,
  );
}

describe('SearchResultsPage', () => {
  beforeEach(() => {
    mockedSearchProjects.mockReset();
    window.localStorage.clear();
  });

  afterEach(() => {
    cleanup();
  });

  it('round-trips client and technology slug filters through the public search API', async () => {
    mockedSearchProjects.mockResolvedValue(buildSearchResponse());
    renderSearchPage();

    await waitFor(() => {
      expect(mockedSearchProjects).toHaveBeenCalledWith(expect.objectContaining({
        q: 'portfolio',
        category: 'platform',
        client: 'Acme',
        technologies: 'react',
      }));
    });

    expect(await screen.findByText('PortfolioForge')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'React' })).toBeInTheDocument();
  });

  it('keeps load-more filters stable and resets pagination when a filter changes', async () => {
    mockedSearchProjects
      .mockResolvedValueOnce(buildSearchResponse({
        data: [
          {
            id: 'project-1',
            slug: 'portfolioforge',
            title: 'PortfolioForge',
            category: 'platform',
            client_name: 'Acme',
            summary: 'First page',
            technologies: [{ id: 'tech-1', name: 'React', slug: 'react', color: '#61dafb' }],
            hero_image: 'https://img/portfolioforge.png',
            score: 0.95,
            explanation: null,
            evidence: [],
          },
          {
            id: 'project-2',
            slug: 'beta',
            title: 'Beta',
            category: 'platform',
            client_name: 'Contoso',
            summary: 'First page second result',
            technologies: [{ id: 'tech-2', name: 'Go', slug: 'go', color: '#00add8' }],
            hero_image: 'https://img/beta.png',
            score: 0.8,
            explanation: null,
            evidence: [],
          },
        ],
        meta: {
          total: 3,
          page_size: 10,
          cursor: '2',
          query: 'portfolio',
          filters_applied: { category: 'platform', client: null, technologies: [] },
        },
      }))
      .mockResolvedValueOnce(buildSearchResponse({
        data: [
          {
            id: 'project-3',
            slug: 'gamma',
            title: 'Gamma',
            category: 'platform',
            client_name: 'Acme',
            summary: 'Second page result',
            technologies: [{ id: 'tech-3', name: 'PostgreSQL', slug: 'postgresql', color: '#336791' }],
            hero_image: 'https://img/gamma.png',
            score: 0.7,
            explanation: null,
            evidence: [],
          },
        ],
        meta: {
          total: 3,
          page_size: 10,
          cursor: null,
          query: 'portfolio',
          filters_applied: { category: 'platform', client: null, technologies: [] },
        },
      }))
      .mockResolvedValueOnce(buildSearchResponse({
        meta: {
          total: 1,
          page_size: 10,
          cursor: null,
          query: 'portfolio',
          filters_applied: { category: 'platform', client: 'Acme', technologies: [] },
        },
      }));

    renderSearchPage('/search?q=portfolio&category=platform');

    expect(await screen.findByText('PortfolioForge')).toBeInTheDocument();
    fireEvent.click(await screen.findByRole('button', { name: 'Cargar más' }));

    await waitFor(() => {
      expect(mockedSearchProjects).toHaveBeenNthCalledWith(2, expect.objectContaining({
        q: 'portfolio',
        category: 'platform',
        client: undefined,
        cursor: '2',
      }));
    });

    expect(await screen.findByText('Gamma')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Acme' }));

    await waitFor(() => {
      expect(mockedSearchProjects).toHaveBeenNthCalledWith(3, expect.objectContaining({
        q: 'portfolio',
        category: 'platform',
        client: 'Acme',
        cursor: undefined,
      }));
    });

    expect(await screen.findByText('PortfolioForge')).toBeInTheDocument();
    expect(screen.queryByText('Gamma')).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Cargar más' })).not.toBeInTheDocument();
  });

  it('passes the active search context and pagination snapshot into detail navigation', async () => {
    mockedSearchProjects
      .mockResolvedValueOnce(buildSearchResponse({
        data: [
          {
            id: 'project-1',
            slug: 'portfolioforge',
            title: 'PortfolioForge',
            category: 'platform',
            client_name: 'Acme',
            summary: 'First page',
            technologies: [{ id: 'tech-1', name: 'React', slug: 'react', color: '#61dafb' }],
            hero_image: 'https://img/portfolioforge.png',
            score: 0.95,
            explanation: null,
            evidence: [],
          },
        ],
        meta: {
          total: 2,
          page_size: 10,
          cursor: 'cursor-2',
          query: 'portfolio',
          filters_applied: { category: 'platform', client: 'Acme', technologies: ['react'] },
        },
      }))
      .mockResolvedValueOnce(buildSearchResponse({
        data: [
          {
            id: 'project-2',
            slug: 'beta',
            title: 'Beta',
            category: 'platform',
            client_name: 'Acme',
            summary: 'Second page',
            technologies: [{ id: 'tech-2', name: 'Go', slug: 'go', color: '#00add8' }],
            hero_image: 'https://img/beta.png',
            score: 0.7,
            explanation: null,
            evidence: [],
          },
        ],
        meta: {
          total: 2,
          page_size: 10,
          cursor: null,
          query: 'portfolio',
          filters_applied: { category: 'platform', client: 'Acme', technologies: ['react'] },
        },
      }));

    renderSearchPage();

    expect(await screen.findByText('PortfolioForge')).toBeInTheDocument();
    fireEvent.click(await screen.findByRole('button', { name: 'Cargar más' }));
    expect(await screen.findByText('Beta')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('link', { name: /PortfolioForge/i }));

    await waitFor(() => {
      const locationState = JSON.parse(screen.getByTestId('detail-location-state').textContent ?? 'null');

      expect(locationState).toEqual(expect.objectContaining({
        activeSearchQuery: 'portfolio',
        activeSearchCategory: 'platform',
        activeSearchClient: 'Acme',
        activeSearchTechnologies: ['react'],
        searchResultsSnapshot: expect.objectContaining({
          total: 2,
          cursor: null,
          results: expect.arrayContaining([
            expect.objectContaining({ slug: 'portfolioforge' }),
            expect.objectContaining({ slug: 'beta' }),
          ]),
        }),
      }));
    });
  });

  it('updates public search copy when the locale changes', async () => {
    mockedSearchProjects.mockResolvedValue(buildSearchResponse({
      meta: {
        total: 1,
        page_size: 10,
        cursor: 'next-page',
        query: 'portfolio',
        filters_applied: { category: null, client: null, technologies: [] },
      },
    }));
    renderSearchPage('/search?q=portfolio');

    expect(await screen.findByText('Búsqueda pública')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Filtros' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'locale-en' }));

    expect(await screen.findByText('Public search')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Filters' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Load more' })).toBeInTheDocument();
  });
});
