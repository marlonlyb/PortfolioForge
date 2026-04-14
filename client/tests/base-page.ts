import { expect, type Page } from '@playwright/test';

export class BasePage {
  constructor(protected readonly page: Page) {}

  async goto(path: string): Promise<void> {
    await this.page.goto(path, { waitUntil: 'networkidle' });
  }

  async expectPath(pathname: RegExp | string): Promise<void> {
    await expect(this.page).toHaveURL(pathname);
  }
}
