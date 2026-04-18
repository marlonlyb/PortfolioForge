import { test } from '@playwright/test';

import {
  LANDING_PROJECT_NAME,
  LANDING_PROJECT_SLUG,
  LANDING_PROMPT_LABEL,
  LANDING_PROMPT_QUERY,
  LANDING_SUGGESTION,
  mockPublicStorefrontApi,
  TEST_PROJECT_NAME,
  TEST_PROJECT_SLUG,
} from '../helpers';
import { StorePage } from './store-page';

const RESPONSIVE_VIEWPORTS = [
  { name: 'mobile', width: 390, height: 844 },
  { name: 'tablet', width: 768, height: 1024 },
  { name: 'desktop', width: 1280, height: 900 },
] as const;

test.describe('Storefront smoke review', () => {
  test.beforeEach(async ({ page }) => {
    await mockPublicStorefrontApi(page);
  });

  test('opens the home page and exposes the catalog fixture', async ({ page }) => {
    const storePage = new StorePage(page);

    await storePage.gotoHome();
    await storePage.expectCatalogProjectVisible(TEST_PROJECT_NAME);
  });

  test('opens the project detail page from the catalog path', async ({ page }) => {
    const storePage = new StorePage(page);

    await storePage.gotoHome();
    await storePage.openProjectFromCatalog(TEST_PROJECT_NAME);
    await storePage.expectProjectDetailLoaded(TEST_PROJECT_NAME, TEST_PROJECT_SLUG);
  });

  test('keeps the landing prompt, suggestions, and results usable across responsive breakpoints', async ({ page }) => {
    const storePage = new StorePage(page);

    await storePage.setEnglishLocale();

    for (const viewport of RESPONSIVE_VIEWPORTS) {
      await test.step(`landing search remains usable on ${viewport.name}`, async () => {
        await page.setViewportSize({ width: viewport.width, height: viewport.height });
        await storePage.gotoHome();

        await storePage.expectLandingSearchReady(LANDING_PROMPT_LABEL);
        await storePage.triggerLandingPrompt(LANDING_PROMPT_LABEL, LANDING_PROMPT_QUERY);
        await storePage.expectLandingResultVisible(LANDING_PROJECT_NAME, LANDING_PROJECT_SLUG);
        await storePage.expectLandingClearVisible();

        await storePage.searchFromLanding('Print');
        await storePage.expectLandingSuggestionVisible(LANDING_SUGGESTION);
        await storePage.selectLandingSuggestion(LANDING_SUGGESTION, LANDING_SUGGESTION);
        await storePage.expectLandingResultVisible(LANDING_PROJECT_NAME, LANDING_PROJECT_SLUG);
      });
    }
  });
});
