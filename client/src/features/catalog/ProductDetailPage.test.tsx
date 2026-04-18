import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import type { ReactNode } from 'react';
import { Link, MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { SessionProvider } from '../../app/providers/SessionProvider';
import { StoreLayout } from '../../app/layouts/StoreLayout';
import type { AdminProjectDetail } from '../../shared/types/admin-project';
import { AppError, API_ERROR_CODES } from '../../shared/api/errors';
import type { Project } from '../../shared/types/project';
import { fetchAdminProjectById } from '../admin-projects/api';
import { CompleteProfilePage } from '../auth/CompleteProfilePage';
import { ProductDetailPage } from './ProductDetailPage';
import { fetchProjectBySlug } from './api';
import type { SessionUser } from '../../app/providers/SessionProvider';

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

function buildSessionUser(overrides: Partial<SessionUser> = {}): SessionUser {
  return {
    id: 'user-1',
    email: 'ada@example.com',
    is_admin: false,
    auth_provider: 'google',
    email_verified: true,
    full_name: 'Ada Lovelace',
    company: 'Analytical Engines',
    profile_completed: true,
    assistant_eligible: true,
    can_use_project_assistant: true,
    created_at: '2026-04-15T00:00:00Z',
    ...overrides,
  };
}

function mockPrivateMe(user: SessionUser) {
  vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
    const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
    if (url.endsWith('/api/v1/private/me')) {
      return new Response(JSON.stringify({ data: user }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      });
    }

    throw new Error(`Unhandled fetch: ${url}`);
  }));
}

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
  return renderDetailRoute();
}

function renderDetailRoute(detailElement: ReactNode = <ProductDetailPage />) {
  return render(
    <MemoryRouter initialEntries={['/projects/portfolioforge']}>
      <SessionProvider>
        <LocaleProvider>
          <Routes>
            <Route path="/" element={<StoreLayout />}>
              <Route index element={<div>landing content</div>} />
              <Route path="projects/:slug" element={detailElement} />
              <Route path="search" element={<div>search page</div>} />
              <Route path="login" element={<div>login page</div>} />
            </Route>
          </Routes>
        </LocaleProvider>
      </SessionProvider>
    </MemoryRouter>,
  );
}

function renderAssistantFlow() {
  return render(
    <MemoryRouter initialEntries={['/projects/portfolioforge']}>
      <SessionProvider>
        <LocaleProvider>
          <Routes>
            <Route path="/" element={<StoreLayout />}>
              <Route index element={<div>landing content</div>} />
              <Route path="projects/:slug" element={<ProductDetailPage />} />
              <Route path="complete-profile" element={<CompleteProfilePage />} />
              <Route path="verify-email" element={<p>verify email destination</p>} />
            </Route>
          </Routes>
        </LocaleProvider>
      </SessionProvider>
    </MemoryRouter>,
  );
}

function createDeferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((res, rej) => {
    resolve = res;
    reject = rej;
  });

  return { promise, resolve, reject };
}

function mockProfileCompletionFlow(initialUser: SessionUser, refreshedUser: SessionUser) {
  let currentUser = initialUser;
  const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
    const method = init?.method ?? 'GET';

    if (url.endsWith('/api/v1/private/me') && method === 'GET') {
      return new Response(JSON.stringify({ data: currentUser }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      });
    }

    if (url.endsWith('/api/v1/private/me/profile') && method === 'PUT') {
      currentUser = refreshedUser;
      return new Response(JSON.stringify({ data: { user: refreshedUser } }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      });
    }

    throw new Error(`Unhandled fetch: ${method} ${url}`);
  });

  vi.stubGlobal('fetch', fetchMock);

  return fetchMock;
}

describe('ProductDetailPage', () => {
  beforeEach(() => {
    mockedFetchProjectBySlug.mockReset();
    mockedFetchAdminProjectById.mockReset();
    mockedFetchAdminProjectById.mockResolvedValue({
      id: 'project-1',
      name: 'PortfolioForge',
      slug: 'portfolioforge',
      description: 'Detailed project description.',
      category: 'platform',
      images: [],
      variants: [],
      active: true,
      source_markdown_url: '',
    } as AdminProjectDetail);
    vi.unstubAllGlobals();
    Object.defineProperty(Element.prototype, 'scrollIntoView', {
      configurable: true,
      value: vi.fn(),
    });
    window.localStorage.clear();
    window.sessionStorage.clear();
    window.localStorage.setItem('portfolioforge.locale', 'en');
    window.scrollTo = vi.fn();
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it('renders the public project detail without assistant UI for signed-out visitors and without admin source leak', async () => {
    mockedFetchProjectBySlug.mockResolvedValue(buildProject());

    renderDetailPage();

    expect(await screen.findByText('Detailed project description.')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Ask project assistant' })).not.toBeInTheDocument();
    expect(screen.queryByText('Project assistant')).not.toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Sign in with Google' })).not.toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Admin markdown source' })).not.toBeInTheDocument();
    expect(mockedFetchAdminProjectById).not.toHaveBeenCalled();
  });

  it('keeps assistant chat hidden for authenticated users with incomplete profiles', async () => {
    mockedFetchProjectBySlug.mockResolvedValue(buildProject());
    mockPrivateMe(buildSessionUser({
      full_name: '',
      company: '',
      profile_completed: false,
      assistant_eligible: false,
      can_use_project_assistant: false,
    }));
    window.sessionStorage.setItem('auth_token', 'token');

    renderDetailPage();

    expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Ask project assistant' })).not.toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Complete profile' })).toBeInTheDocument();
  });

  it('renders assistant chat entry for eligible authenticated users', async () => {
    mockedFetchProjectBySlug.mockResolvedValue(buildProject());
    mockPrivateMe(buildSessionUser());
    window.sessionStorage.setItem('auth_token', 'token');

    renderDetailPage();

    expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
    expect(await screen.findByRole('button', { name: 'Ask project assistant' })).toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Complete profile' })).not.toBeInTheDocument();
  });

  it('routes unverified local accounts toward email verification', async () => {
    mockedFetchProjectBySlug.mockResolvedValue(buildProject());
    mockPrivateMe(buildSessionUser({
      auth_provider: 'local',
      email_verified: false,
      assistant_eligible: false,
      can_use_project_assistant: false,
    }));
    window.sessionStorage.setItem('auth_token', 'token');

    renderDetailPage();

    expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
    expect(screen.getByText('Verify your email to keep your local account eligible for the assistant.')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Verify email' })).toBeInTheDocument();
  });

  it('restores assistant continuity after profile completion returns to the same project', async () => {
    mockedFetchProjectBySlug.mockResolvedValue(buildProject());
    const incompleteUser = buildSessionUser({
      full_name: '',
      company: '',
      profile_completed: false,
      assistant_eligible: false,
      can_use_project_assistant: false,
    });
    const eligibleUser = buildSessionUser();
    const fetchMock = mockProfileCompletionFlow(incompleteUser, eligibleUser);
    window.sessionStorage.setItem('auth_token', 'token');
    window.sessionStorage.setItem('assistant_history:portfolioforge', JSON.stringify([
      { role: 'assistant', content: 'Restored answer.' },
    ]));

    renderAssistantFlow();

    expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('link', { name: 'Complete profile' }));

    expect(await screen.findByRole('heading', { name: 'Unlock the project assistant' })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Full name'), { target: { value: 'Ada Lovelace' } });
    fireEvent.change(screen.getByLabelText('Company'), { target: { value: 'Analytical Engines' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save profile' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();

    fireEvent.click(await screen.findByRole('button', { name: 'Ask project assistant' }));

    expect(screen.getByText('Restored answer.')).toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Complete profile' })).not.toBeInTheDocument();

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/private/me/profile',
        expect.objectContaining({ method: 'PUT' }),
      );
    });
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
    mockPrivateMe(buildSessionUser({
      is_admin: true,
      can_use_project_assistant: false,
    }));
    window.sessionStorage.setItem('auth_token', 'token');

    renderDetailPage();

    expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
    const adminLink = await screen.findByRole('link', { name: 'Admin markdown source' });
    expect(adminLink).toHaveAttribute('href', 'https://mlbautomation.com/docs.md');

    await waitFor(() => {
      expect(mockedFetchAdminProjectById).toHaveBeenCalledWith('project-1');
    });
  });

	it('renders medium variants in the gallery and high variants in the lightbox', async () => {
		mockedFetchProjectBySlug.mockResolvedValue(buildProject({
			media: [{
				id: 'media-1',
				project_id: 'project-1',
				media_type: 'image',
				low_url: 'https://cdn.example.com/project-low.webp',
				medium_url: 'https://cdn.example.com/project-medium.webp',
				high_url: 'https://cdn.example.com/project-high.webp',
				fallback_url: 'https://cdn.example.com/project-fallback.webp',
				alt_text: 'Project hero',
				caption: 'Hero image',
				sort_order: 0,
				featured: true,
			}],
		}));

		renderDetailPage();

		await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' });
		const galleryImage = screen.getByRole('img', { name: 'Project hero' });
		expect(galleryImage).toHaveAttribute('src', 'https://cdn.example.com/project-medium.webp');

		fireEvent.click(screen.getByRole('button', { name: 'Ver completa' }));

		const dialog = await screen.findByRole('dialog', { name: 'Image preview' });
		expect(within(dialog).getByRole('img', { name: 'Project hero' })).toHaveAttribute(
			'src',
			'https://cdn.example.com/project-high.webp',
		);
	});

   	it('hides the assistant entrypoint when markdown is absent or cleared', async () => {
	  mockedFetchProjectBySlug.mockResolvedValue(buildProject({ assistant_available: false }));

	  renderDetailPage();

	  expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
	  expect(screen.queryByRole('button', { name: 'Ask project assistant' })).not.toBeInTheDocument();
	  expect(screen.queryByRole('link', { name: 'Sign in with Google' })).not.toBeInTheDocument();
	  expect(mockedFetchAdminProjectById).not.toHaveBeenCalled();
	});

	it('shows the generic store header while loading and promotes project context after resolution', async () => {
	  const deferred = createDeferred<Project>();
	  mockedFetchProjectBySlug.mockReturnValue(deferred.promise);

	  renderDetailPage();

	  expect(screen.getByRole('heading', { level: 1, name: 'Project portfolio' })).toBeInTheDocument();
	  expect(screen.getByText('Strategy, execution, and technical judgment.')).toBeInTheDocument();
	  expect(screen.getByText('Curated public archive')).toBeInTheDocument();
	  expect(screen.getByText('Loading project…')).toBeInTheDocument();
	  await waitFor(() => {
	    expect(mockedFetchProjectBySlug).toHaveBeenCalled();
	  });
	  const initialFetchCount = mockedFetchProjectBySlug.mock.calls.length;

	  deferred.resolve(buildProject({
	    client_name: 'Analytical Engines',
	    category: 'platform',
	  }));

	  expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
	  expect(mockedFetchProjectBySlug).toHaveBeenCalledTimes(initialFetchCount);
	  expect(mockedFetchProjectBySlug).toHaveBeenCalledWith('portfolioforge', 'en');
	  expect(screen.getByText('Detailed project description.', { selector: '.detail__summary--hero' })).toBeInTheDocument();
	  expect(screen.queryByText('Detailed project description.', { selector: '.app-header__summary' })).not.toBeInTheDocument();
	  expect(screen.getByText('platform · Analytical Engines')).toBeInTheDocument();
	});

	it('falls back to the generic store header when the detail request fails', async () => {
	  mockedFetchProjectBySlug.mockRejectedValue(
	    new AppError(404, { code: API_ERROR_CODES.NOT_FOUND, message: 'Project not found' }),
	  );

	  renderDetailPage();

	  expect(await screen.findByText('Project not found')).toBeInTheDocument();
	  expect(screen.getByRole('heading', { level: 1, name: 'Project portfolio' })).toBeInTheDocument();
	  expect(screen.getByText('Curated public archive')).toBeInTheDocument();
	});

	it('moves the descriptive summary into the hero and keeps technologies there as the authoritative surface', async () => {
	  mockedFetchProjectBySlug.mockResolvedValue(buildProject({
	    client_name: 'Analytical Engines',
	    technologies: [
	      { id: 'tech-1', name: 'React', slug: 'react', category: 'frontend' },
	      { id: 'tech-2', name: 'TypeScript', slug: 'typescript', category: 'language' },
	      { id: 'tech-3', name: 'Vite', slug: 'vite', category: 'tooling' },
	      { id: 'tech-4', name: 'Vitest', slug: 'vitest', category: 'testing' },
	      { id: 'tech-5', name: 'Playwright', slug: 'playwright', category: 'testing' },
	    ],
	    media: [{
	      id: 'media-1',
	      project_id: 'project-1',
	      media_type: 'image',
	      low_url: 'https://cdn.example.com/project-low.webp',
	      medium_url: 'https://cdn.example.com/project-medium.webp',
	      high_url: 'https://cdn.example.com/project-high.webp',
	      fallback_url: 'https://cdn.example.com/project-fallback.webp',
	      alt_text: 'Project hero',
	      caption: 'Hero image',
	      sort_order: 0,
	      featured: true,
	    }],
	  }));

	  renderDetailPage();

	  expect(await screen.findByRole('heading', { level: 1, name: 'PortfolioForge' })).toBeInTheDocument();
	  const heroFacts = screen.getByLabelText('Project highlights');
	  const heroCard = heroFacts.closest('article');
	  const heroSummary = screen.getByText('Detailed project description.', { selector: '.detail__summary--hero' });
	  expect(heroCard).not.toBeNull();
	  expect(screen.queryByRole('heading', { level: 2, name: 'PortfolioForge' })).not.toBeInTheDocument();
	  expect(screen.queryByText('Detailed project description.', { selector: '.app-header__summary' })).not.toBeInTheDocument();
	  expect(within(heroCard as HTMLElement).getByText('Detailed project description.')).toBeInTheDocument();
	  expect(
	    heroSummary.compareDocumentPosition(heroFacts) & Node.DOCUMENT_POSITION_FOLLOWING,
	  ).toBeTruthy();
	  expect(screen.getByRole('link', { name: '← Back to projects' })).toBeInTheDocument();
	  expect(screen.getByLabelText('Project highlights')).toBeInTheDocument();
	  expect(screen.getByRole('button', { name: 'Ver completa' })).toBeInTheDocument();
	  expect(screen.getByText('React')).toBeInTheDocument();
	  expect(screen.getByText('TypeScript')).toBeInTheDocument();
	  expect(screen.getByText('Vite')).toBeInTheDocument();
	  expect(screen.getByText('Vitest')).toBeInTheDocument();
	  expect(screen.getByText('Playwright')).toBeInTheDocument();
	  expect(screen.queryByText('+1')).not.toBeInTheDocument();
	  expect(screen.queryAllByText('Technologies').length).toBe(1);
	});

	it('clears stale project header content during slug and route transitions', async () => {
	  mockedFetchProjectBySlug.mockImplementation(async (slug) => {
	    if (slug === 'alpha') {
	      return buildProject({
	        id: 'project-alpha',
	        name: 'Alpha',
	        slug: 'alpha',
	        description: 'Alpha summary.',
	        category: 'alpha category',
	      });
	    }

	    if (slug === 'beta') {
	      return buildProject({
	        id: 'project-beta',
	        name: 'Beta',
	        slug: 'beta',
	        description: 'Beta summary.',
	        category: 'beta category',
	      });
	    }

	    throw new Error(`Unexpected slug: ${slug}`);
	  });

	  function DetailHarness() {
	    return (
	      <>
	        <nav>
	          <Link to="/projects/beta">Go to beta</Link>
	          <Link to="/login">Leave detail</Link>
	        </nav>
	        <ProductDetailPage />
	      </>
	    );
	  }

	  render(
	    <MemoryRouter initialEntries={['/projects/alpha']}>
	      <SessionProvider>
	        <LocaleProvider>
	          <Routes>
	            <Route path="/" element={<StoreLayout />}>
	              <Route path="projects/:slug" element={<DetailHarness />} />
	              <Route path="login" element={<div>login page</div>} />
	            </Route>
	          </Routes>
	        </LocaleProvider>
	      </SessionProvider>
	    </MemoryRouter>,
	  );

	  expect(await screen.findByRole('heading', { level: 1, name: 'Alpha' })).toBeInTheDocument();

	  fireEvent.click(screen.getByRole('link', { name: 'Go to beta' }));

	  await waitFor(() => {
	    expect(screen.getByRole('heading', { level: 1, name: 'Project portfolio' })).toBeInTheDocument();
	  });

	  expect(await screen.findByRole('heading', { level: 1, name: 'Beta' })).toBeInTheDocument();

	  fireEvent.click(screen.getByRole('link', { name: 'Leave detail' }));

	  await waitFor(() => {
	    expect(screen.getByRole('heading', { level: 1, name: 'Project portfolio' })).toBeInTheDocument();
	  });
	  expect(screen.getByText('login page')).toBeInTheDocument();
	});
});
