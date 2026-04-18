import { expect, type Page } from '@playwright/test';

import { BasePage } from '../base-page';

export class StorePage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async gotoHome(): Promise<void> {
    await this.goto('/');
  }

  async setEnglishLocale(): Promise<void> {
    await this.page.addInitScript(() => {
      window.localStorage.setItem('portfolioforge.locale', 'en');
      document.documentElement.lang = 'en';
    });
  }

  async gotoProjectDetail(slug: string): Promise<void> {
    await this.goto(`/projects/${slug}`);
  }

  async expectCatalogProjectVisible(projectName: string): Promise<void> {
    await expect(this.page.getByRole('link', { name: new RegExp(projectName, 'i') }).first()).toBeVisible();
  }

  async expectLandingSearchReady(promptLabel: string): Promise<void> {
    await expect(this.page.getByRole('heading', { name: /How can I help/i })).toBeVisible();
    await expect(this.page.getByRole('textbox', { name: 'Search' })).toBeVisible();
    await expect(this.page.getByRole('button', { name: promptLabel })).toBeVisible();
  }

  async triggerLandingPrompt(promptLabel: string, expectedQuery: string): Promise<void> {
    await this.page.getByRole('button', { name: promptLabel }).click();
    await expect(this.page.getByRole('textbox', { name: 'Search' })).toHaveValue(expectedQuery);
  }

  async expectLandingResultVisible(projectName: string, slug: string): Promise<void> {
    const projectLink = this.page.getByRole('link', { name: new RegExp(projectName, 'i') }).first();
    await expect(projectLink).toBeVisible();
    await expect(projectLink).toHaveAttribute('href', `/projects/${slug}`);
  }

  async searchFromLanding(query: string): Promise<void> {
    await this.page.getByRole('textbox', { name: 'Search' }).fill(query);
  }

  async expectLandingSuggestionVisible(suggestion: string): Promise<void> {
    await expect(this.page.getByRole('button', { name: suggestion, exact: true })).toBeVisible();
  }

  async selectLandingSuggestion(suggestion: string, expectedQuery: string): Promise<void> {
    await this.page.getByRole('button', { name: suggestion, exact: true }).dispatchEvent('click');
    await expect(this.page.getByRole('textbox', { name: 'Search' })).toHaveValue(expectedQuery);
  }

  async expectLandingClearVisible(): Promise<void> {
    await expect(this.page.getByRole('button', { name: 'Clear search' })).toBeVisible();
  }

  async openProjectFromCatalog(projectName: string): Promise<void> {
    await this.page.getByRole('link', { name: new RegExp(projectName, 'i') }).first().click();
  }

  async expectProjectDetailLoaded(projectName: string, slug: string): Promise<void> {
    await this.expectPath(new RegExp(`/projects/${slug}$`));
    await expect(this.page.getByRole('heading', { name: projectName })).toBeVisible();
    await expect(this.page.getByRole('link', { name: /volver a proyectos|back to projects/i })).toBeVisible();
  }
}
