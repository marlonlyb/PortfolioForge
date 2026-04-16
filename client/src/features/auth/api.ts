import { httpPost, httpPut } from '../../shared/api/http';
import type { SessionUser } from '../../app/providers/SessionProvider';

export interface AdminLoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: SessionUser;
  token: string;
  expires_in: number;
}

export interface EmailLoginRequest {
  email: string;
}

export interface EmailLoginVerifyRequest {
  email: string;
  code: string;
}

export interface EmailLoginDispatchResponse {
  verification_required: boolean;
  message: string;
  cooldown_seconds: number;
}

export interface GoogleLoginRequest {
  id_token: string;
}

export interface UpdateProfileRequest {
  full_name: string;
  company: string;
}

export interface UpdateProfileResponse {
  user: SessionUser;
}

export function adminLogin(payload: AdminLoginRequest): Promise<LoginResponse> {
  return httpPost<LoginResponse>('/api/v1/admin/login', payload);
}

export function loginWithGoogle(payload: GoogleLoginRequest): Promise<LoginResponse> {
  return httpPost<LoginResponse>('/api/v1/public/login/google', payload);
}

export function requestEmailLogin(payload: EmailLoginRequest): Promise<EmailLoginDispatchResponse> {
  return httpPost<EmailLoginDispatchResponse>('/api/v1/public/login/email/request', payload);
}

export function verifyEmailLogin(payload: EmailLoginVerifyRequest): Promise<LoginResponse> {
  return httpPost<LoginResponse>('/api/v1/public/login/email/verify', payload);
}

export function updateMyProfile(payload: UpdateProfileRequest): Promise<UpdateProfileResponse> {
  return httpPut<UpdateProfileResponse>('/api/v1/private/me/profile', payload);
}
