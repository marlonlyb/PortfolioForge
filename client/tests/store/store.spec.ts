import { test } from '@playwright/test';

import { mockPublicStorefrontApi, TEST_PROJECT_NAME, TEST_PROJECT_SLUG } from '../helpers';
import { StorePage } from './store-page';

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
});
