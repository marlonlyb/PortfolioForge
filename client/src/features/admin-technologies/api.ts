import { httpGet, httpPost, httpPut, httpDelete } from '../../shared/api/http';
import type { Technology } from '../../shared/types/project';

export interface AdminTechnologyListResponse {
  items: Technology[];
}

export interface CreateTechnologyPayload {
  name: string;
  category: string;
  icon?: string;
  color?: string;
}

export interface UpdateTechnologyPayload {
  name: string;
  category: string;
  icon?: string;
  color?: string;
}

export function fetchAdminTechnologies(): Promise<AdminTechnologyListResponse> {
  return httpGet<AdminTechnologyListResponse>('/api/v1/admin/technologies');
}

export function fetchAdminTechnologyById(id: string): Promise<Technology> {
  return httpGet<Technology>(`/api/v1/admin/technologies/${id}`);
}

export function createTechnology(payload: CreateTechnologyPayload): Promise<Technology> {
  return httpPost<Technology>('/api/v1/admin/technologies', payload);
}

export function updateTechnology(id: string, payload: UpdateTechnologyPayload): Promise<Technology> {
  return httpPut<Technology>(`/api/v1/admin/technologies/${id}`, payload);
}

export function deleteTechnology(id: string): Promise<void> {
  return httpDelete(`/api/v1/admin/technologies/${id}`);
}
