import type { Page, Route } from '@playwright/test';

const TEST_PROJECT_SLUG = 'portfolioforge';
const TEST_PROJECT_NAME = 'PortfolioForge Control Hub';
const LANDING_PROJECT_SLUG = 'printer-05-plc-migration';
const LANDING_PROJECT_NAME = 'Printer 05 PLC Migration';
const LANDING_PROMPT_LABEL = 'Show me the Printer 05 PLC migration';
const LANDING_PROMPT_QUERY = 'Printer 05';
const LANDING_SUGGESTION = 'Printer 05';

function createSvgDataUrl(label: string): string {
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900" viewBox="0 0 1600 900">
      <defs>
        <linearGradient id="bg" x1="0" x2="1" y1="0" y2="1">
          <stop offset="0%" stop-color="#111827" />
          <stop offset="100%" stop-color="#1d4ed8" />
        </linearGradient>
      </defs>
      <rect width="1600" height="900" fill="url(#bg)" rx="32" />
      <text x="120" y="420" fill="#f9fafb" font-size="84" font-family="Arial, sans-serif">PortfolioForge</text>
      <text x="120" y="520" fill="#bfdbfe" font-size="46" font-family="Arial, sans-serif">${label}</text>
    </svg>
  `;

  return `data:image/svg+xml,${encodeURIComponent(svg)}`;
}

const projectImage = createSvgDataUrl('Visual smoke review fixture');

export const smokeProject = {
  id: 'project-portfolioforge',
  name: TEST_PROJECT_NAME,
  slug: TEST_PROJECT_SLUG,
  description: 'Industrial portfolio platform used as a stable fixture for local browser smoke review.',
  category: 'platform',
  client_name: 'Internal Lab',
  status: 'published',
  featured: true,
  active: true,
  assistant_available: true,
  images: [],
  media: [
    {
      id: 'media-hero',
      project_id: 'project-portfolioforge',
      media_type: 'image',
      fallback_url: projectImage,
      low_url: projectImage,
      medium_url: projectImage,
      high_url: projectImage,
      caption: 'Fixture preview used for Playwright smoke review.',
      alt_text: 'PortfolioForge smoke fixture preview',
      sort_order: 0,
      featured: true,
    },
  ],
  created_at: 1710000000,
  updated_at: 1710000000,
  technologies: [
    { id: 'tech-react', name: 'React', slug: 'react', category: 'frontend' },
    { id: 'tech-playwright', name: 'Playwright', slug: 'playwright', category: 'qa' },
  ],
  profile: {
    project_id: 'project-portfolioforge',
    business_goal: 'Validate the public storefront visually before release.',
    problem_statement: 'The local environment had no browser automation available for quick regression review.',
    solution_summary: 'Playwright adds a lightweight browser smoke path against the local Vite app.',
    architecture: 'Vite frontend with mocked public API responses during smoke automation.',
    integrations: ['Vite dev server', 'Playwright Chromium'],
    ai_usage: 'Not applicable for this smoke fixture.',
    technical_decisions: ['Use route interception instead of requiring a live API.'],
    challenges: ['Keep the setup isolated from existing Vitest configuration.'],
    results: ['Local visual review can cover home and detail flows with stable fixtures.'],
    metrics: {
      smoke_paths: 2,
      browser: 'chromium',
    },
    timeline: ['Install Playwright', 'Mock public API', 'Verify storefront routes'],
    updated_at: 1710000000,
  },
} as const;

const landingProjectImage = createSvgDataUrl('Printer 05 PLC migration fixture');

const landingProject = {
  id: 'project-printer-05',
  name: LANDING_PROJECT_NAME,
  slug: LANDING_PROJECT_SLUG,
  description: 'CompactLogix migration for Printer 05 with Ethernet/IP commissioning and operator handoff.',
  category: 'automation',
  client_name: 'Printer 05',
  status: 'published',
  featured: true,
  active: true,
  assistant_available: true,
  images: [],
  media: [
    {
      id: 'media-printer-05',
      project_id: 'project-printer-05',
      media_type: 'image',
      fallback_url: landingProjectImage,
      low_url: landingProjectImage,
      medium_url: landingProjectImage,
      high_url: landingProjectImage,
      caption: 'Printer 05 Playwright landing-search fixture.',
      alt_text: 'Printer 05 PLC migration fixture preview',
      sort_order: 0,
      featured: true,
    },
  ],
  created_at: 1710000000,
  updated_at: 1710000000,
  technologies: [
    { id: 'tech-compactlogix', name: 'CompactLogix', slug: 'compactlogix', category: 'plc' },
    { id: 'tech-ethernet-ip', name: 'Ethernet/IP', slug: 'ethernet-ip', category: 'network' },
  ],
  profile: {
    project_id: 'project-printer-05',
    business_goal: 'Modernize Printer 05 controls without changing the operator workflow.',
    problem_statement: 'The line needed a deterministic PLC migration and cleaner diagnostics for maintenance teams.',
    solution_summary: 'Delivered a CompactLogix migration with Ethernet/IP commissioning and startup support.',
    architecture: 'Allen-Bradley PLC architecture with Ethernet/IP-connected devices and staged commissioning.',
    integrations: ['Allen-Bradley CompactLogix', 'Ethernet/IP', 'Factory floor HMI'],
    ai_usage: 'None.',
    technical_decisions: ['Preserve deterministic controls behavior during cutover.'],
    challenges: ['Coordinate commissioning without extending downtime.'],
    results: ['Printer 05 resumed production with modernized controls and better diagnostics.'],
    metrics: {
      migration_scope: 1,
      commissioning_days: 3,
    },
    timeline: ['Audit legacy PLC', 'Implement CompactLogix replacement', 'Commission Printer 05'],
    updated_at: 1710000000,
  },
} as const;

function fulfillJson(route: Route, body: unknown): Promise<void> {
  return route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify(body),
  });
}

function buildSearchResponse(query: string) {
  return {
    data: [
      {
        id: 'search-result-printer-05',
        slug: LANDING_PROJECT_SLUG,
        title: LANDING_PROJECT_NAME,
        category: 'automation',
        client_name: 'Printer 05',
        summary: 'CompactLogix migration for Printer 05 with Ethernet/IP commissioning and operator handoff.',
        technologies: [
          { id: 'tech-compactlogix', name: 'CompactLogix', slug: 'compactlogix' },
          { id: 'tech-ethernet-ip', name: 'Ethernet/IP', slug: 'ethernet-ip' },
        ],
        hero_image: landingProjectImage,
        score: 0.99,
        explanation: null,
        evidence: [],
      },
    ],
    meta: {
      total: 1,
      page_size: 1,
      cursor: null,
      query,
      filters_applied: {
        category: null,
        client: null,
        technologies: [],
      },
    },
  };
}

export async function mockPublicStorefrontApi(page: Page): Promise<void> {
  await page.route('**/api/v1/public/site-settings', (route) => fulfillJson(route, { data: {} }));
  await page.route('**/api/v1/public/projects?*', (route) => fulfillJson(route, { data: { items: [smokeProject, landingProject] } }));
  await page.route(`**/api/v1/public/projects/${TEST_PROJECT_SLUG}?*`, (route) => fulfillJson(route, { data: smokeProject }));
  await page.route(`**/api/v1/public/projects/${LANDING_PROJECT_SLUG}?*`, (route) => fulfillJson(route, { data: landingProject }));
  await page.route('**/api/v1/public/search?*', async (route) => {
    const url = new URL(route.request().url());
    const query = url.searchParams.get('q') ?? '';
    await fulfillJson(route, buildSearchResponse(query));
  });
}

export {
  LANDING_PROJECT_NAME,
  LANDING_PROJECT_SLUG,
  LANDING_PROMPT_LABEL,
  LANDING_PROMPT_QUERY,
  LANDING_SUGGESTION,
  TEST_PROJECT_NAME,
  TEST_PROJECT_SLUG,
};
