import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import type { AdminProjectDetail } from '../../shared/types/admin-project';
import type { Project } from '../../shared/types/project';
import { fetchAdminProjectById } from '../admin-projects/api';
import { ProductDetailPage } from './ProductDetailPage';
import { fetchProjectBySlug } from './api';

vi.mock('embla-carousel-react', () => ({
  default: vi.fn(() => [
    () => undefined,
    {
      selectedScrollSnap: () => 0,
      canScrollPrev: () => false,
      canScrollNext: () => false,
      on: () => undefined,
      off: () => undefined,
      reInit: () => undefined,
      scrollTo: () => undefined,
      scrollPrev: () => undefined,
      scrollNext: () => undefined,
    },
  ]),
}));

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
    fetchProjectBySlug: vi.fn(),
    sendProjectAssistantMessage: vi.fn(),
  };
});

vi.mock('../admin-projects/api', async () => {
  const actual = await vi.importActual<typeof import('../admin-projects/api')>('../admin-projects/api');
  return {
    ...actual,
    fetchAdminProjectById: vi.fn(),
  };
});

vi.mock('../search/api', async () => {
  const actual = await vi.importActual<typeof import('../search/api')>('../search/api');
  return {
    ...actual,
    searchProjects: vi.fn().mockResolvedValue({ data: [] }),
  };
});

const mockedFetchProjectBySlug = vi.mocked(fetchProjectBySlug);
const mockedFetchAdminProjectById = vi.mocked(fetchAdminProjectById);

function buildProject(overrides: Partial<Project> = {}): Project {
  return {
    id: 'project-1',
    name: 'PortfolioForge',
    slug: 'portfolioforge',
    description: 'Detailed project description.',
    category: 'platform',
    status: 'published',
    featured: false,
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

function renderDetailPage() {
  return render(
    <MemoryRouter initialEntries={['/projects/portfolioforge']}>
      <LocaleProvider>
        <Routes>
          <Route path="/projects/:slug" element={<ProductDetailPage />} />
        </Routes>
      </LocaleProvider>
    </MemoryRouter>,
  );
}

describe('ProductDetailPage', () => {
  beforeEach(() => {
    mockedFetchProjectBySlug.mockReset();
    mockedFetchAdminProjectById.mockReset();
    window.localStorage.clear();
    window.sessionStorage.clear();
  });

  afterEach(() => {
    cleanup();
  });

  it('renders the public project detail and assistant without admin source leak', async () => {
    mockedFetchProjectBySlug.mockResolvedValue(buildProject());

    renderDetailPage();

    expect(await screen.findByRole('heading', { name: 'PortfolioForge' })).toBeInTheDocument();
    expect(screen.getByText('Detailed project description.')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Ask project assistant' })).toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Admin markdown source' })).not.toBeInTheDocument();
    expect(mockedFetchAdminProjectById).not.toHaveBeenCalled();
  });

  it('shows the admin markdown source only for authenticated admin reads', async () => {
    mockedFetchProjectBySlug.mockResolvedValue(buildProject());
    mockedFetchAdminProjectById.mockResolvedValue({
      id: 'project-1',
      name: 'PortfolioForge',
      slug: 'portfolioforge',
      description: 'Detailed project description.',
      category: 'platform',
      images: [],
      variants: [],
      active: true,
      source_markdown_url: 'https://mlbautomation.com/docs.md',
    } as AdminProjectDetail);
    window.sessionStorage.setItem('auth_token', 'token');

    renderDetailPage();

    const adminLink = await screen.findByRole('link', { name: 'Admin markdown source' });
    expect(adminLink).toHaveAttribute('href', 'https://mlbautomation.com/docs.md');

    await waitFor(() => {
      expect(mockedFetchAdminProjectById).toHaveBeenCalledWith('project-1');
    });
  });

	it('hides the assistant entrypoint when markdown is absent or cleared', async () => {
	  mockedFetchProjectBySlug.mockResolvedValue(buildProject({ assistant_available: false }));

	  renderDetailPage();

	  expect(await screen.findByRole('heading', { name: 'PortfolioForge' })).toBeInTheDocument();
	  expect(screen.queryByRole('button', { name: 'Ask project assistant' })).not.toBeInTheDocument();
	  expect(mockedFetchAdminProjectById).not.toHaveBeenCalled();
	});
});
