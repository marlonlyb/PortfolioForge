import { expect, type Page } from '@playwright/test';

import { BasePage } from '../base-page';

export class AdminSettingsPage extends BasePage {
	constructor(page: Page) {
		super(page);
	}

	async bootstrapSession(): Promise<void> {
		await this.page.addInitScript(() => {
			window.sessionStorage.setItem('auth_token', 'playwright-admin-token');
		});
	}

	async gotoSettings(): Promise<void> {
		await this.goto('/admin/settings/case-studies');
	}

	async gotoWorkflow(runId?: string): Promise<void> {
		const suffix = runId ? `?run=${runId}` : '';
		await this.goto(`/admin/settings/case-studies${suffix}`);
	}

	async startWorkflow(sourcePath: string): Promise<void> {
		await this.page.getByLabel('Canonical source path').fill(sourcePath);
		await this.page.getByRole('checkbox', { name: 'Run localization backfill after import' }).uncheck();
		await this.page.getByRole('checkbox', { name: 'Refresh search document after import/localization' }).uncheck();
		await this.page.getByRole('button', { name: 'Start workflow run' }).click();
	}

	async confirmStep(): Promise<void> {
		await this.page.getByRole('button', { name: 'Confirm' }).click();
	}

	async runStep(): Promise<void> {
		await this.page.getByRole('button', { name: 'Run step' }).click();
	}

	async resumeRun(): Promise<void> {
		await this.page.getByRole('button', { name: 'Continue from latest checkpoint' }).click();
	}

	async expectBrandingFormVisible(): Promise<void> {
		await expect(this.page.getByRole('heading', { name: 'Public branding' })).toBeVisible();
		await expect(this.page.getByRole('button', { name: 'Save settings' })).toBeVisible();
	}

	async expectWorkflowIsolation(): Promise<void> {
		await expect(this.page).toHaveURL(/\/admin\/settings\/case-studies/);
		await expect(this.page.getByRole('heading', { name: 'Case-study workflow' })).toBeVisible();
		await expect(this.page.getByRole('button', { name: 'Save settings' })).toBeVisible();
	}

	async expectGenerationMessage(): Promise<void> {
		await expect(this.page.getByText('Raw folder → canonical generation is intentionally out of scope for this MVP.')).toBeVisible();
	}

	async expectRunStatus(status: string): Promise<void> {
		await expect(this.page.locator('p.admin__helper-copy').filter({ hasText: new RegExp(`Status:\\s*${status}`, 'i') }).first()).toBeVisible();
	}

	async expectLogMessage(message: string): Promise<void> {
		await expect(this.page.locator('ul.admin__list li').filter({ hasText: message }).first()).toBeVisible();
	}

	async expectProjectId(projectId: string): Promise<void> {
		await expect(this.page.getByText(projectId, { exact: true }).first()).toBeVisible();
	}

	async expectCanonicalUrl(url: string): Promise<void> {
		await expect(this.page.getByText(url, { exact: true }).first()).toBeVisible();
	}
}
