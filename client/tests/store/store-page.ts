import { expect, type Page } from '@playwright/test';

import { BasePage } from '../base-page';

export class StorePage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async gotoHome(): Promise<void> {
    await this.goto('/');
  }

  async gotoProjectDetail(slug: string): Promise<void> {
    await this.goto(`/projects/${slug}`);
  }

  async expectCatalogProjectVisible(projectName: string): Promise<void> {
    await expect(this.page.getByRole('link', { name: new RegExp(projectName, 'i') }).first()).toBeVisible();
  }

  async openProjectFromCatalog(projectName: string): Promise<void> {
    await this.page.getByRole('link', { name: new RegExp(projectName, 'i') }).first().click();
  }

  async expectProjectDetailLoaded(projectName: string, slug: string): Promise<void> {
    await this.expectPath(new RegExp(`/projects/${slug}$`));
    await expect(this.page.getByRole('heading', { name: projectName })).toBeVisible();
    await expect(this.page.getByRole('button', { name: 'Ask project assistant' })).toBeVisible();
  }
}
