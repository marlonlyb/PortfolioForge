import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

import { AdminCaseStudyWorkflowPage } from './AdminCaseStudyWorkflowPage';
import { AdminSiteSettingsPage } from './AdminSiteSettingsPage';
import {
  fetchAdminSiteSettings,
  updateAdminSiteSettings,
} from '../../shared/api/siteSettings';
import {
	confirmCaseStudyWorkflowStep,
  fetchCaseStudyWorkflowLogs,
	fetchCaseStudyWorkflowRun,
	resumeCaseStudyWorkflowRun,
	retryCaseStudyWorkflowStep,
	startCaseStudyWorkflowStep,
  startCaseStudyWorkflowRun,
	type CaseStudyWorkflowRun,
} from './api';

vi.mock('../../shared/api/siteSettings', async () => {
  const actual = await vi.importActual<typeof import('../../shared/api/siteSettings')>('../../shared/api/siteSettings');
  return {
    ...actual,
    fetchAdminSiteSettings: vi.fn(),
    updateAdminSiteSettings: vi.fn(),
  };
});

vi.mock('./api', async () => {
  const actual = await vi.importActual<typeof import('./api')>('./api');
  return {
    ...actual,
		confirmCaseStudyWorkflowStep: vi.fn(),
    fetchCaseStudyWorkflowLogs: vi.fn(),
		fetchCaseStudyWorkflowRun: vi.fn(),
		resumeCaseStudyWorkflowRun: vi.fn(),
		retryCaseStudyWorkflowStep: vi.fn(),
		startCaseStudyWorkflowStep: vi.fn(),
    startCaseStudyWorkflowRun: vi.fn(),
  };
});

const mockedFetchAdminSiteSettings = vi.mocked(fetchAdminSiteSettings);
const mockedUpdateAdminSiteSettings = vi.mocked(updateAdminSiteSettings);
const mockedConfirmCaseStudyWorkflowStep = vi.mocked(confirmCaseStudyWorkflowStep);
const mockedFetchCaseStudyWorkflowLogs = vi.mocked(fetchCaseStudyWorkflowLogs);
const mockedFetchCaseStudyWorkflowRun = vi.mocked(fetchCaseStudyWorkflowRun);
const mockedResumeCaseStudyWorkflowRun = vi.mocked(resumeCaseStudyWorkflowRun);
const mockedRetryCaseStudyWorkflowStep = vi.mocked(retryCaseStudyWorkflowStep);
const mockedStartCaseStudyWorkflowStep = vi.mocked(startCaseStudyWorkflowStep);
const mockedStartCaseStudyWorkflowRun = vi.mocked(startCaseStudyWorkflowRun);

const awaitingConfirmationRun: CaseStudyWorkflowRun = {
	id: 'run-1',
	status: 'awaiting_confirmation',
	source: {
		allowed_root: '/safe/root',
		requested_path: '/safe/root/90. dev_portfolioforge/demo',
		normalized_path: '/safe/root/90. dev_portfolioforge/demo',
		canonical_root_path: '/safe/root/90. dev_portfolioforge',
		canonical_directory: '/safe/root/90. dev_portfolioforge/demo',
		canonical_markdown_path: '/safe/root/90. dev_portfolioforge/demo/demo.md',
		slug: 'demo',
	},
	options: { run_localization_backfill: true, run_reembed: true, locales: ['ca', 'en'] },
	steps: [
		{
			run_id: 'run-1',
			step: 'resolve_source',
			status: 'succeeded',
			requires_confirmation: false,
			attempt_count: 1,
			output: { canonical_directory: '/safe/root/90. dev_portfolioforge/demo' },
		},
		{
			run_id: 'run-1',
			step: 'publish_canonical',
			status: 'awaiting_confirmation',
			requires_confirmation: true,
			attempt_count: 0,
		},
	],
	created_at: '2026-04-18T00:00:00Z',
	updated_at: '2026-04-18T00:00:00Z',
	generation_scope: {
		canonical_generation_available: false,
		canonical_generation_message: 'MVP starts from an existing canonical source.',
	},
};

const failedImportRun: CaseStudyWorkflowRun = {
	...awaitingConfirmationRun,
	id: 'run-failed',
	status: 'failed',
	canonical_url: 'https://example.com/demo/demo.md',
	steps: [
		awaitingConfirmationRun.steps[0]!,
		{
			run_id: 'run-failed',
			step: 'publish_canonical',
			status: 'succeeded',
			requires_confirmation: true,
			attempt_count: 1,
			output: { canonical_url: 'https://example.com/demo/demo.md' },
		},
		{
			run_id: 'run-failed',
			step: 'import_or_update_project',
			status: 'failed',
			requires_confirmation: true,
			attempt_count: 1,
			error_message: 'Import failed',
		},
	],
};

const resumedRun: CaseStudyWorkflowRun = {
	...failedImportRun,
	status: 'succeeded',
	project_id: 'project-1',
	steps: [
		failedImportRun.steps[0]!,
		failedImportRun.steps[1]!,
		{
			run_id: 'run-failed',
			step: 'import_or_update_project',
			status: 'succeeded',
			requires_confirmation: true,
			attempt_count: 2,
			output: { project_id: 'project-1' },
		},
	],
};

describe('Admin settings workflow pages', () => {
  beforeEach(() => {
    mockedFetchAdminSiteSettings.mockReset();
    mockedUpdateAdminSiteSettings.mockReset();
    mockedConfirmCaseStudyWorkflowStep.mockReset();
    mockedFetchCaseStudyWorkflowLogs.mockReset();
    mockedFetchCaseStudyWorkflowRun.mockReset();
    mockedResumeCaseStudyWorkflowRun.mockReset();
    mockedRetryCaseStudyWorkflowStep.mockReset();
    mockedStartCaseStudyWorkflowStep.mockReset();
    mockedStartCaseStudyWorkflowRun.mockReset();
    sessionStorage.clear();
  });

  afterEach(() => {
    cleanup();
  });

  it('shows the case-study workflow entry on the settings hub', async () => {
    mockedFetchAdminSiteSettings.mockResolvedValue({ public_hero_logo_url: '', public_hero_logo_alt: '' });

    render(
      <MemoryRouter>
        <AdminSiteSettingsPage />
      </MemoryRouter>,
    );

    expect(await screen.findByRole('heading', { name: 'Case-study workflow' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Open workflow' })).toHaveAttribute('href', '/admin/settings/case-studies');
  });

  it('starts a workflow run from the dedicated settings subpage', async () => {
		mockedStartCaseStudyWorkflowRun.mockResolvedValue(awaitingConfirmationRun);
		mockedFetchCaseStudyWorkflowLogs.mockResolvedValue({ items: [] });

    render(
      <MemoryRouter initialEntries={['/admin/settings/case-studies']}>
        <Routes>
          <Route path="/admin/settings" element={<p>settings</p>} />
          <Route path="/admin/settings/case-studies" element={<AdminCaseStudyWorkflowPage />} />
        </Routes>
      </MemoryRouter>,
    );

    fireEvent.change(screen.getByLabelText('Canonical source path'), {
      target: { value: '/safe/root/90. dev_portfolioforge/demo' },
    });
    fireEvent.change(screen.getByLabelText('Localization locales (optional comma-separated subset)'), {
      target: { value: 'ca,en' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Start workflow run' }));

    await waitFor(() => {
      expect(mockedStartCaseStudyWorkflowRun).toHaveBeenCalledWith(
        expect.objectContaining({
          source_path: '/safe/root/90. dev_portfolioforge/demo',
          locales: ['ca', 'en'],
        }),
      );
    });
		await waitFor(() => {
			expect(sessionStorage.getItem('admin.case-study-workflow.last-run-id')).toBe('run-1');
		});
		expect(screen.getByText('MVP starts from an existing canonical source.')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Confirm' })).toBeInTheDocument();
		expect(screen.queryByRole('button', { name: 'Save settings' })).not.toBeInTheDocument();
	});

	it('keeps branding save isolated from workflow actions on the settings hub', async () => {
		mockedFetchAdminSiteSettings.mockResolvedValue({
			public_hero_logo_url: 'https://cdn.example.com/logo.svg',
			public_hero_logo_alt: 'PortfolioForge',
		});
		mockedUpdateAdminSiteSettings.mockResolvedValue({
			public_hero_logo_url: 'https://cdn.example.com/brand.svg',
			public_hero_logo_alt: 'New logo',
		});

		render(
			<MemoryRouter>
				<AdminSiteSettingsPage />
			</MemoryRouter>,
		);

		fireEvent.change(await screen.findByLabelText('Public logo URL'), {
			target: { value: 'https://cdn.example.com/brand.svg' },
		});
		fireEvent.change(screen.getByLabelText('Alt text'), {
			target: { value: 'New logo' },
		});
		fireEvent.click(screen.getByRole('button', { name: 'Save settings' }));

		await waitFor(() => {
			expect(mockedUpdateAdminSiteSettings).toHaveBeenCalledWith({
				public_hero_logo_url: 'https://cdn.example.com/brand.svg',
				public_hero_logo_alt: 'New logo',
			});
		});
		expect(mockedStartCaseStudyWorkflowRun).not.toHaveBeenCalled();
		expect(await screen.findByRole('status')).toHaveTextContent('Public hero logo updated.');
	});

	it('supports confirm, start, retry, and resume semantics on the workflow page', async () => {
		mockedFetchCaseStudyWorkflowRun.mockResolvedValue(failedImportRun);
		mockedFetchCaseStudyWorkflowLogs.mockResolvedValue({
			items: [
				{
					id: 1,
					run_id: 'run-failed',
					step: 'import_or_update_project',
					level: 'error',
					message: 'Import failed',
					created_at: '2026-04-18T00:00:00Z',
				},
			],
		});
		mockedConfirmCaseStudyWorkflowStep.mockResolvedValue({
			...awaitingConfirmationRun,
			steps: [
				awaitingConfirmationRun.steps[0]!,
				{
					...awaitingConfirmationRun.steps[1]!,
					status: 'pending',
					confirmation_granted_at: '2026-04-18T00:01:00Z',
				},
			],
		});
		mockedStartCaseStudyWorkflowStep.mockResolvedValue({
			...awaitingConfirmationRun,
			status: 'awaiting_confirmation',
			canonical_url: 'https://example.com/demo/demo.md',
			steps: [
				awaitingConfirmationRun.steps[0]!,
				{
					run_id: 'run-1',
					step: 'publish_canonical',
					status: 'succeeded',
					requires_confirmation: true,
					attempt_count: 1,
					output: { canonical_url: 'https://example.com/demo/demo.md' },
				},
				{
					run_id: 'run-1',
					step: 'import_or_update_project',
					status: 'awaiting_confirmation',
					requires_confirmation: true,
					attempt_count: 0,
				},
			],
		});
		mockedRetryCaseStudyWorkflowStep.mockResolvedValue(resumedRun);
		mockedResumeCaseStudyWorkflowRun.mockResolvedValue(resumedRun);

		render(
			<MemoryRouter initialEntries={['/admin/settings/case-studies?run=run-failed']}>
				<Routes>
					<Route path="/admin/settings" element={<p>settings</p>} />
					<Route path="/admin/settings/case-studies" element={<AdminCaseStudyWorkflowPage />} />
				</Routes>
			</MemoryRouter>,
		);

		expect(await screen.findByRole('button', { name: 'Continue from latest checkpoint' })).toBeInTheDocument();
		expect(await screen.findByRole('button', { name: 'Retry step' })).toBeInTheDocument();
		expect(screen.getByText('Import failed')).toBeInTheDocument();

		fireEvent.click(screen.getByRole('button', { name: 'Retry step' }));
		await waitFor(() => {
			expect(mockedRetryCaseStudyWorkflowStep).toHaveBeenCalledWith('run-failed', 'import_or_update_project');
		});

		cleanup();
		mockedResumeCaseStudyWorkflowRun.mockClear();
		mockedFetchCaseStudyWorkflowRun.mockResolvedValue(failedImportRun);
		mockedFetchCaseStudyWorkflowLogs.mockResolvedValue({
			items: [
				{
					id: 1,
					run_id: 'run-failed',
					step: 'import_or_update_project',
					level: 'error',
					message: 'Import failed',
					created_at: '2026-04-18T00:00:00Z',
				},
			],
		});
		render(
			<MemoryRouter initialEntries={['/admin/settings/case-studies?run=run-failed']}>
				<Routes>
					<Route path="/admin/settings" element={<p>settings</p>} />
					<Route path="/admin/settings/case-studies" element={<AdminCaseStudyWorkflowPage />} />
				</Routes>
			</MemoryRouter>,
		);

		fireEvent.click(await screen.findByRole('button', { name: 'Continue from latest checkpoint' }));
		await waitFor(() => {
			expect(mockedResumeCaseStudyWorkflowRun).toHaveBeenCalledWith('run-failed');
		});

		cleanup();
		mockedConfirmCaseStudyWorkflowStep.mockClear();
		mockedStartCaseStudyWorkflowStep.mockClear();
		mockedFetchCaseStudyWorkflowRun.mockResolvedValue(awaitingConfirmationRun);
		mockedFetchCaseStudyWorkflowLogs.mockResolvedValue({ items: [] });
		render(
			<MemoryRouter initialEntries={['/admin/settings/case-studies?run=run-1']}>
				<Routes>
					<Route path="/admin/settings" element={<p>settings</p>} />
					<Route path="/admin/settings/case-studies" element={<AdminCaseStudyWorkflowPage />} />
				</Routes>
			</MemoryRouter>,
		);

		fireEvent.click(await screen.findByRole('button', { name: 'Confirm' }));
		await waitFor(() => {
			expect(mockedConfirmCaseStudyWorkflowStep).toHaveBeenCalledWith('run-1', 'publish_canonical');
		});
		fireEvent.click(await screen.findByRole('button', { name: 'Run step' }));
		await waitFor(() => {
			expect(mockedStartCaseStudyWorkflowStep).toHaveBeenCalledWith('run-1', 'publish_canonical');
		});
  });
});
