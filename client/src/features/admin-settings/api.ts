import { httpGet, httpPost } from '../../shared/api/http';
import { API_ERROR_CODES, AppError } from '../../shared/api/errors';

export type CaseStudyWorkflowStepName =
  | 'resolve_source'
  | 'publish_canonical'
  | 'import_or_update_project'
  | 'localization_backfill'
  | 'reembed';

export type CaseStudyWorkflowStatus =
  | 'pending'
  | 'blocked'
  | 'awaiting_confirmation'
  | 'running'
  | 'succeeded'
  | 'failed'
  | 'skipped';

export interface CaseStudyWorkflowStep {
  run_id: string;
  step: CaseStudyWorkflowStepName;
  status: CaseStudyWorkflowStatus;
  requires_confirmation: boolean;
  confirmation_granted_at?: string;
  started_at?: string;
  finished_at?: string;
  attempt_count: number;
  error_message?: string;
  output?: Record<string, unknown>;
}

export interface CaseStudyWorkflowRun {
  id: string;
  status: CaseStudyWorkflowStatus;
  source: {
    allowed_root: string;
    requested_path: string;
    normalized_path: string;
    canonical_root_path: string;
    canonical_directory: string;
    canonical_markdown_path: string;
    slug: string;
  };
  options: {
    run_localization_backfill: boolean;
    run_reembed: boolean;
    locales: string[];
  };
  canonical_url?: string;
  project_id?: string;
  steps: CaseStudyWorkflowStep[];
  created_at: string;
  updated_at: string;
  last_error?: string;
  generation_scope: {
    canonical_generation_available: boolean;
    canonical_generation_message: string;
  };
}

export interface CaseStudyWorkflowLogEntry {
  id: number;
  run_id: string;
  step: CaseStudyWorkflowStepName;
  level: 'info' | 'warn' | 'error';
  message: string;
  created_at: string;
}

export interface CaseStudyWorkflowAvailability {
  configured: boolean;
  reason?: string;
}

export interface StartCaseStudyWorkflowPayload {
  source_path: string;
  slug?: string;
  run_localization_backfill?: boolean;
  run_reembed?: boolean;
  locales?: string[];
}

export function fetchCaseStudyWorkflowAvailability(): Promise<CaseStudyWorkflowAvailability> {
  return httpGet<CaseStudyWorkflowAvailability>('/api/v1/admin/settings/case-study-workflow');
}

export function isWorkflowUnavailableError(error: unknown): error is AppError {
  return error instanceof AppError && error.code === API_ERROR_CODES.WORKFLOW_UNAVAILABLE;
}

export function startCaseStudyWorkflowRun(
  payload: StartCaseStudyWorkflowPayload,
): Promise<CaseStudyWorkflowRun> {
  return httpPost<CaseStudyWorkflowRun>('/api/v1/admin/settings/case-study-runs', payload);
}

export function fetchCaseStudyWorkflowRun(id: string): Promise<CaseStudyWorkflowRun> {
  return httpGet<CaseStudyWorkflowRun>(`/api/v1/admin/settings/case-study-runs/${id}`);
}

export function fetchCaseStudyWorkflowLogs(id: string): Promise<{ items: CaseStudyWorkflowLogEntry[] }> {
  return httpGet<{ items: CaseStudyWorkflowLogEntry[] }>(`/api/v1/admin/settings/case-study-runs/${id}/logs`);
}

export function confirmCaseStudyWorkflowStep(
  id: string,
  step: CaseStudyWorkflowStepName,
): Promise<CaseStudyWorkflowRun> {
  return httpPost<CaseStudyWorkflowRun>(`/api/v1/admin/settings/case-study-runs/${id}/steps/${step}/confirm`, {});
}

export function startCaseStudyWorkflowStep(
  id: string,
  step: CaseStudyWorkflowStepName,
): Promise<CaseStudyWorkflowRun> {
  return httpPost<CaseStudyWorkflowRun>(`/api/v1/admin/settings/case-study-runs/${id}/steps/${step}/start`, {});
}

export function retryCaseStudyWorkflowStep(
  id: string,
  step: CaseStudyWorkflowStepName,
): Promise<CaseStudyWorkflowRun> {
  return httpPost<CaseStudyWorkflowRun>(`/api/v1/admin/settings/case-study-runs/${id}/steps/${step}/retry`, {});
}

export function resumeCaseStudyWorkflowRun(id: string): Promise<CaseStudyWorkflowRun> {
  return httpPost<CaseStudyWorkflowRun>(`/api/v1/admin/settings/case-study-runs/${id}/resume`, {});
}
