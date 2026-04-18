import { test, type Page, type Route } from '@playwright/test';

import { AdminSettingsPage } from './admin-settings-page';

const adminUser = {
	id: 'admin-1',
	email: 'admin@portfolioforge.local',
	is_admin: true,
	auth_provider: 'local',
	email_verified: true,
	profile_completed: true,
	assistant_eligible: false,
	can_use_project_assistant: false,
	created_at: '2026-04-18T00:00:00Z',
} as const;

const baseSource = {
	allowed_root: '/safe/root',
	requested_path: '/safe/root/90. dev_portfolioforge/demo',
	normalized_path: '/safe/root/90. dev_portfolioforge/demo',
	canonical_root_path: '/safe/root/90. dev_portfolioforge',
	canonical_directory: '/safe/root/90. dev_portfolioforge/demo',
	canonical_markdown_path: '/safe/root/90. dev_portfolioforge/demo/demo.md',
	slug: 'demo',
} as const;

function ok(data: unknown) {
	return {
		status: 200,
		contentType: 'application/json',
		body: JSON.stringify({ data }),
	};
}

function created(data: unknown) {
	return {
		status: 201,
		contentType: 'application/json',
		body: JSON.stringify({ data }),
	};
}

async function fulfill(route: Route, payload: ReturnType<typeof ok>): Promise<void> {
	await route.fulfill(payload);
}

function buildHappyRuns() {
	const started = {
		id: 'run-happy',
		status: 'awaiting_confirmation',
		source: baseSource,
		options: { run_localization_backfill: false, run_reembed: false, locales: [] },
		steps: [
			{ run_id: 'run-happy', step: 'resolve_source', status: 'succeeded', requires_confirmation: false, attempt_count: 1, output: { slug: 'demo' } },
			{ run_id: 'run-happy', step: 'publish_canonical', status: 'awaiting_confirmation', requires_confirmation: true, attempt_count: 0 },
			{ run_id: 'run-happy', step: 'import_or_update_project', status: 'blocked', requires_confirmation: true, attempt_count: 0 },
			{ run_id: 'run-happy', step: 'localization_backfill', status: 'skipped', requires_confirmation: false, attempt_count: 0 },
			{ run_id: 'run-happy', step: 'reembed', status: 'skipped', requires_confirmation: false, attempt_count: 0 },
		],
		created_at: '2026-04-18T00:00:00Z',
		updated_at: '2026-04-18T00:00:00Z',
		generation_scope: {
			canonical_generation_available: false,
			canonical_generation_message: 'MVP starts from an existing canonical source.',
		},
	};

	const publishPending = {
		...started,
		steps: [
			started.steps[0],
			{ ...started.steps[1], status: 'pending', confirmation_granted_at: '2026-04-18T00:01:00Z' },
			started.steps[2],
			started.steps[3],
			started.steps[4],
		],
	};

	const importAwaiting = {
		...started,
		canonical_url: 'https://example.com/demo/demo.md',
		steps: [
			started.steps[0],
			{ ...started.steps[1], status: 'succeeded', attempt_count: 1, output: { canonical_url: 'https://example.com/demo/demo.md' } },
			{ ...started.steps[2], status: 'awaiting_confirmation' },
			started.steps[3],
			started.steps[4],
		],
	};

	const importPending = {
		...importAwaiting,
		steps: [
			importAwaiting.steps[0],
			importAwaiting.steps[1],
			{ ...importAwaiting.steps[2], status: 'pending', confirmation_granted_at: '2026-04-18T00:02:00Z' },
			importAwaiting.steps[3],
			importAwaiting.steps[4],
		],
	};

	const finished = {
		...importAwaiting,
		status: 'succeeded',
		project_id: 'project-42',
		steps: [
			importAwaiting.steps[0],
			importAwaiting.steps[1],
			{ ...importAwaiting.steps[2], status: 'succeeded', attempt_count: 1, output: { project_id: 'project-42' } },
			importAwaiting.steps[3],
			importAwaiting.steps[4],
		],
	};

	return { started, publishPending, importAwaiting, importPending, finished };
}

function buildRetryRuns() {
	const failed = {
		id: 'run-retry',
		status: 'failed',
		source: baseSource,
		options: { run_localization_backfill: false, run_reembed: false, locales: [] },
		canonical_url: 'https://example.com/demo/demo.md',
		steps: [
			{ run_id: 'run-retry', step: 'resolve_source', status: 'succeeded', requires_confirmation: false, attempt_count: 1 },
			{ run_id: 'run-retry', step: 'publish_canonical', status: 'succeeded', requires_confirmation: true, attempt_count: 1, output: { canonical_url: 'https://example.com/demo/demo.md' } },
			{ run_id: 'run-retry', step: 'import_or_update_project', status: 'failed', requires_confirmation: true, attempt_count: 1, error_message: 'Import failed' },
			{ run_id: 'run-retry', step: 'localization_backfill', status: 'skipped', requires_confirmation: false, attempt_count: 0 },
			{ run_id: 'run-retry', step: 'reembed', status: 'skipped', requires_confirmation: false, attempt_count: 0 },
		],
		created_at: '2026-04-18T00:00:00Z',
		updated_at: '2026-04-18T00:00:00Z',
		generation_scope: {
			canonical_generation_available: false,
			canonical_generation_message: 'MVP starts from an existing canonical source.',
		},
	};

	const resumed = {
		...failed,
		status: 'succeeded',
		project_id: 'project-42',
		steps: [
			failed.steps[0],
			failed.steps[1],
			{ run_id: 'run-retry', step: 'import_or_update_project', status: 'succeeded', requires_confirmation: true, attempt_count: 2, output: { project_id: 'project-42' } },
			failed.steps[3],
			failed.steps[4],
		],
	};

	return { failed, resumed };
}

async function mockAdminApi(page: Page, scenario: 'happy' | 'retry'): Promise<void> {
	const happy = buildHappyRuns();
	const retry = buildRetryRuns();
	let currentRun: any = scenario === 'happy' ? null : retry.failed;
	let logs: any[] = scenario === 'happy'
		? []
		: [
			{ id: 1, run_id: 'run-retry', step: 'publish_canonical', level: 'info', message: 'Publish completed.', created_at: '2026-04-18T00:00:00Z' },
			{ id: 2, run_id: 'run-retry', step: 'import_or_update_project', level: 'error', message: 'Import failed', created_at: '2026-04-18T00:01:00Z' },
		];

	await page.route('**/api/v1/private/me', (route) => fulfill(route, ok(adminUser)));
	await page.route('**/api/v1/admin/site-settings', async (route) => {
		if (route.request().method() === 'PUT') {
			await fulfill(route, ok({ public_hero_logo_url: 'https://cdn.example.com/logo.svg', public_hero_logo_alt: 'PortfolioForge' }));
			return;
		}
		await fulfill(route, ok({ public_hero_logo_url: '', public_hero_logo_alt: '' }));
	});
	await page.route('**/api/v1/admin/settings/case-study-runs', async (route) => {
		currentRun = happy.started;
		logs = [
			{ id: 1, run_id: 'run-happy', step: 'resolve_source', level: 'info', message: 'Canonical source resolved.', created_at: '2026-04-18T00:00:00Z' },
		];
		await fulfill(route, created(currentRun));
	});
	await page.route('**/api/v1/admin/settings/case-study-runs/**', async (route) => {
		const url = new URL(route.request().url());
		const pathname = url.pathname;
		const method = route.request().method();

		if (pathname.endsWith('/logs')) {
			await fulfill(route, ok({ items: logs }));
			return;
		}
		if (pathname.endsWith('/resume')) {
			currentRun = retry.resumed;
			logs = [...logs, { id: 3, run_id: 'run-retry', step: 'import_or_update_project', level: 'info', message: 'Resume completed import.', created_at: '2026-04-18T00:02:00Z' }];
			await fulfill(route, ok(currentRun));
			return;
		}
		if (pathname.endsWith('/steps/publish_canonical/confirm')) {
			currentRun = happy.publishPending;
			logs = [...logs, { id: 2, run_id: 'run-happy', step: 'publish_canonical', level: 'warn', message: 'Operator confirmed the step.', created_at: '2026-04-18T00:01:00Z' }];
			await fulfill(route, ok(currentRun));
			return;
		}
		if (pathname.endsWith('/steps/publish_canonical/start')) {
			currentRun = happy.importAwaiting;
			logs = [...logs, { id: 3, run_id: 'run-happy', step: 'publish_canonical', level: 'info', message: 'Step completed successfully.', created_at: '2026-04-18T00:02:00Z' }];
			await fulfill(route, ok(currentRun));
			return;
		}
		if (pathname.endsWith('/steps/import_or_update_project/confirm')) {
			currentRun = happy.importPending;
			logs = [...logs, { id: 4, run_id: 'run-happy', step: 'import_or_update_project', level: 'warn', message: 'Operator confirmed the step.', created_at: '2026-04-18T00:03:00Z' }];
			await fulfill(route, ok(currentRun));
			return;
		}
		if (pathname.endsWith('/steps/import_or_update_project/start')) {
			currentRun = happy.finished;
			logs = [...logs, { id: 5, run_id: 'run-happy', step: 'import_or_update_project', level: 'info', message: 'Step completed successfully.', created_at: '2026-04-18T00:04:00Z' }];
			await fulfill(route, ok(currentRun));
			return;
		}
		if (method === 'GET') {
			await fulfill(route, ok(currentRun ?? retry.failed));
			return;
		}

		await route.abort();
	});
}

test.describe('Admin settings workflow', () => {
	test('navigates from settings to workflow and completes the canonical-first happy path', async ({ page }) => {
		const adminSettingsPage = new AdminSettingsPage(page);
		await mockAdminApi(page, 'happy');
		await adminSettingsPage.bootstrapSession();

		await adminSettingsPage.gotoSettings();
		await adminSettingsPage.expectBrandingFormVisible();
		await adminSettingsPage.openWorkflowFromSettings();
		await adminSettingsPage.expectWorkflowIsolation();
		await adminSettingsPage.expectGenerationMessage();

		await adminSettingsPage.startWorkflow('/safe/root/90. dev_portfolioforge/demo');
		await adminSettingsPage.expectRunStatus('Needs confirmation');
		await adminSettingsPage.confirmStep();
		await adminSettingsPage.runStep();
		await adminSettingsPage.expectCanonicalUrl('https://example.com/demo/demo.md');
		await adminSettingsPage.confirmStep();
		await adminSettingsPage.runStep();

		await adminSettingsPage.expectRunStatus('Done');
		await adminSettingsPage.expectProjectId('project-42');
		await adminSettingsPage.expectLogMessage('Operator confirmed the step.');
		await adminSettingsPage.expectLogMessage('Step completed successfully.');
	});

	test('resumes a failed import run without losing published canonical evidence', async ({ page }) => {
		const adminSettingsPage = new AdminSettingsPage(page);
		await mockAdminApi(page, 'retry');
		await adminSettingsPage.bootstrapSession();

		await adminSettingsPage.gotoWorkflow('run-retry');
		await adminSettingsPage.expectWorkflowIsolation();
		await adminSettingsPage.expectRunStatus('Failed');
		await adminSettingsPage.expectCanonicalUrl('https://example.com/demo/demo.md');
		await adminSettingsPage.expectLogMessage('Import failed');
		await adminSettingsPage.resumeRun();

		await adminSettingsPage.expectRunStatus('Done');
		await adminSettingsPage.expectProjectId('project-42');
		await adminSettingsPage.expectLogMessage('Resume completed import.');
	});
});
