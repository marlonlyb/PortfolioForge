import { httpGet, httpPatch, httpPost, httpPut } from '../../shared/api/http';
import type { AdminProjectDetail, AdminProjectListResponse } from '../../shared/types/admin-project';
import type { ProjectMedia, ProjectProfileMetrics } from '../../shared/types/project';
import type { PublicContentFieldKey, PublicLocale, TranslationMode } from '../../shared/i18n/config';

export interface ProjectReadiness {
  project_id: string;
  level: 'incomplete' | 'basic' | 'complete';
  missing_fields: string[];
  has_name: boolean;
  has_description: boolean;
  has_category: boolean;
  has_technologies: boolean;
  has_solution_summary: boolean;
}

export interface ReembedResponse {
  message: string;
  project_id?: string;
}

export interface CreateAdminProjectPayload {
  name: string;
  description: string;
  category: string;
  brand?: string;
  source_markdown_url?: string;
  images: string[];
  media?: ProjectMedia[];
  active: boolean;
  variants?: CreateAdminProjectVariantPayload[];
}

export interface CreateAdminProjectVariantPayload {
  sku: string;
  color: string;
  size: string;
  price: number;
  stock: number;
  image_url?: string;
}

export interface UpdateAdminProjectPayload {
  name: string;
  description: string;
  category: string;
  brand?: string;
  source_markdown_url?: string;
  images: string[];
  media?: ProjectMedia[];
  active: boolean;
  variants?: UpdateAdminProjectVariantPayload[];
}

export interface UpdateAdminProjectVariantPayload {
  id?: string;
  sku: string;
  color: string;
  size: string;
  price: number;
  stock: number;
  image_url?: string;
}

export interface UpdateAdminProjectStatusPayload {
  active: boolean;
}

export interface UpdateProjectEnrichmentProfilePayload {
  solution_summary?: string;
  delivery_scope?: string;
  responsibility_scope?: string;
  architecture?: string;
  business_goal?: string;
  problem_statement?: string;
  ai_usage?: string;
  integrations?: string[];
  technical_decisions?: string[];
  challenges?: string[];
  results?: string[];
  metrics?: ProjectProfileMetrics;
  timeline?: string[];
}

export interface UpdateProjectEnrichmentPayload {
  profile: UpdateProjectEnrichmentProfilePayload;
  technology_ids: string[];
}

export interface LocalizedAdminField {
  value: unknown;
  mode: TranslationMode;
}

export interface AdminProjectLocalizationLocale {
  locale: PublicLocale;
  fields: Record<PublicContentFieldKey, LocalizedAdminField>;
}

export interface AdminProjectLocalizationsResponse {
  project_id: string;
  base: Record<PublicContentFieldKey, unknown>;
  locales: Record<string, AdminProjectLocalizationLocale>;
}

export interface SaveProjectLocalizationsPayload {
  fields: Partial<Record<PublicContentFieldKey, unknown>>;
}

export function updateProjectEnrichment(
  id: string,
  payload: UpdateProjectEnrichmentPayload,
): Promise<void> {
  return httpPut<void>(`/api/v1/admin/projects/${id}/enrichment`, payload);
}

export function fetchProjectLocalizations(id: string): Promise<AdminProjectLocalizationsResponse> {
  return httpGet<AdminProjectLocalizationsResponse>(`/api/v1/admin/projects/${id}/localizations`);
}

export function saveProjectLocalizations(
  id: string,
  locale: PublicLocale,
  payload: SaveProjectLocalizationsPayload,
): Promise<void> {
  return httpPut<void>(`/api/v1/admin/projects/${id}/localizations/${locale}`, payload);
}

export function fetchAdminProjects(): Promise<AdminProjectListResponse> {
  return httpGet<AdminProjectListResponse>('/api/v1/admin/projects');
}

export function fetchAdminProjectById(id: string): Promise<AdminProjectDetail> {
  return httpGet<AdminProjectDetail>(`/api/v1/admin/projects/${id}`);
}

export function createAdminProject(payload: CreateAdminProjectPayload): Promise<AdminProjectDetail> {
  return httpPost<AdminProjectDetail>('/api/v1/admin/projects', payload);
}

export function updateAdminProject(
  id: string,
  payload: UpdateAdminProjectPayload,
): Promise<AdminProjectDetail> {
  return httpPut<AdminProjectDetail>(`/api/v1/admin/projects/${id}`, payload);
}

export function updateAdminProjectStatus(
  id: string,
  payload: UpdateAdminProjectStatusPayload,
): Promise<AdminProjectDetail> {
  return httpPatch<AdminProjectDetail>(`/api/v1/admin/projects/${id}/status`, payload);
}

export function fetchProjectReadiness(id: string): Promise<ProjectReadiness> {
  return httpGet<ProjectReadiness>(`/api/v1/admin/projects/${id}/readiness`);
}

export function reembedProject(id: string): Promise<ReembedResponse> {
  return httpPost<ReembedResponse>(`/api/v1/admin/projects/${id}/reembed`, {});
}

export function reembedStale(): Promise<ReembedResponse> {
  return httpPost<ReembedResponse>('/api/v1/admin/projects/reembed-stale', {});
}
