import { httpDelete, httpGet, httpPatch } from '../../shared/api/http';
import type {
  AdminUserDetail,
  AdminUserListResponse,
  AdminUserUpdateRequest,
} from '../../shared/types/admin-user';

export function fetchAdminUsers(): Promise<AdminUserListResponse> {
  return httpGet<AdminUserListResponse>('/api/v1/admin/users');
}

export function fetchAdminUserById(id: string): Promise<AdminUserDetail> {
  return httpGet<AdminUserDetail>(`/api/v1/admin/users/${id}`);
}

export function updateAdminUser(id: string, payload: AdminUserUpdateRequest): Promise<AdminUserDetail> {
  return httpPatch<AdminUserDetail>(`/api/v1/admin/users/${id}`, payload);
}

export function deleteAdminUser(id: string): Promise<void> {
  return httpDelete(`/api/v1/admin/users/${id}`);
}
