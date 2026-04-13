import { httpGet } from '../../shared/api/http';
import type { Project, ProjectListResponse } from '../../shared/types/project';

/**
 * Fetch all public projects using the portfolio contract.
 */
export function fetchProjects(): Promise<ProjectListResponse> {
  return httpGet<ProjectListResponse>('/api/v1/public/projects');
}

/**
 * Fetch a single public project by slug using the rich case-study contract.
 */
export function fetchProjectBySlug(slug: string): Promise<Project> {
  return httpGet<Project>(`/api/v1/public/projects/${slug}`);
}
