import { cleanup, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import type { Project } from '../../shared/types/project';
import { searchProjects } from '../search/api';
import { fetchProjects } from './api';
import { CatalogPage } from './CatalogPage';

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
    fetchProjects: vi.fn(),
  };
});

vi.mock('../search/api', async () => {
  const actual = await vi.importActual<typeof import('../search/api')>('../search/api');
  return {
    ...actual,
    searchProjects: vi.fn().mockResolvedValue({ data: [] }),
  };
});

const mockedFetchProjects = vi.mocked(fetchProjects);
const mockedSearchProjects = vi.mocked(searchProjects);

function buildProject(overrides: Partial<Project> = {}): Project {
  return {
    id: 'project-1',
    name: 'PortfolioForge',
    slug: 'portfolioforge',
    description: 'Detailed project description.',
    category: 'platform',
    industry_type: 'automatización industrial',
    final_product: 'Panel HMI para diagnóstico y monitoreo',
    status: 'published',
    featured: true,
    active: true,
    assistant_available: true,
    images: [],
    media: [],
    created_at: 1710000000,
    updated_at: 1710000000,
    technologies: [],
    ...overrides,
  };
}

function renderCatalogPage() {
  return render(
    <MemoryRouter>
      <LocaleProvider>
        <CatalogPage />
      </LocaleProvider>
    </MemoryRouter>,
  );
}

describe('CatalogPage', () => {
  beforeEach(() => {
    mockedFetchProjects.mockReset();
    mockedSearchProjects.mockClear();
    window.localStorage.clear();
  });

  afterEach(() => {
    cleanup();
  });

  it('renders low_url first for catalog card images', async () => {
    mockedFetchProjects.mockResolvedValue({
      items: [
        buildProject({
          media: [
            {
              id: 'media-1',
              project_id: 'project-1',
              media_type: 'image',
              low_url: 'https://cdn.example.com/project-low.webp',
              medium_url: 'https://cdn.example.com/project-medium.webp',
              high_url: 'https://cdn.example.com/project-high.webp',
              fallback_url: 'https://cdn.example.com/project-original.jpg',
              sort_order: 0,
              featured: true,
            },
          ],
        }),
      ],
    });

    renderCatalogPage();

    const image = await screen.findByRole('img', { name: 'PortfolioForge' });
    expect(image).toHaveAttribute('src', 'https://cdn.example.com/project-low.webp');
    expect(screen.getByText('automatización industrial')).toBeInTheDocument();
    expect(screen.getByText('Panel HMI para diagnóstico y monitoreo')).toBeInTheDocument();
  });
});
