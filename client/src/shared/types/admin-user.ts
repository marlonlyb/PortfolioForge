export interface AdminUserSummary {
  id: string;
  email: string;
  is_admin: boolean;
  auth_provider: 'local' | 'google';
  email_verified: boolean;
  full_name?: string;
  company?: string;
  created_at: string;
  updated_at?: string;
  last_login_at?: string;
  deleted_at?: string;
}

export interface AdminUserDetail extends AdminUserSummary {}

export interface AdminUserListResponse {
  items: AdminUserSummary[];
}

export interface AdminUserUpdateRequest {
  is_admin: boolean;
}
