import { httpGet, httpPost } from '../../shared/api/http';
import type { PublicLocale } from '../../shared/i18n/config';
import type {
  Project,
  ProjectAssistantRequest,
  ProjectAssistantResponse,
  ProjectListResponse,
} from '../../shared/types/project';

/**
 * Fetch all public projects using the portfolio contract.
 */
export function fetchProjects(locale: PublicLocale): Promise<ProjectListResponse> {
  return httpGet<ProjectListResponse>(`/api/v1/public/projects?lang=${encodeURIComponent(locale)}`);
}

/**
 * Fetch a single public project by slug using the rich case-study contract.
 */
export function fetchProjectBySlug(slug: string, locale: PublicLocale): Promise<Project> {
  return httpGet<Project>(`/api/v1/public/projects/${slug}?lang=${encodeURIComponent(locale)}`);
}

export function sendProjectAssistantMessage(
  slug: string,
  payload: ProjectAssistantRequest,
): Promise<ProjectAssistantResponse> {
  return httpPost<ProjectAssistantResponse>(`/api/v1/public/projects/${slug}/assistant/messages`, payload);
}
