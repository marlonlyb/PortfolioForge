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

export interface PublicEmailRequest {
  email: string;
}

export interface PublicLoginRequest {
  email: string;
  password: string;
}

export interface PublicSignupRequest {
  email: string;
  password: string;
  confirm_password: string;
}

export interface EmailVerificationVerifyRequest {
  email: string;
  code: string;
}

export interface EmailVerificationDispatchResponse {
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

export function publicLogin(payload: PublicLoginRequest): Promise<LoginResponse> {
  return httpPost<LoginResponse>('/api/v1/public/login', payload);
}

export function publicSignup(payload: PublicSignupRequest): Promise<EmailVerificationDispatchResponse> {
  return httpPost<EmailVerificationDispatchResponse>('/api/v1/public/signup', payload);
}

export function requestEmailVerification(payload: PublicEmailRequest): Promise<EmailVerificationDispatchResponse> {
  return httpPost<EmailVerificationDispatchResponse>('/api/v1/public/email-verification/request', payload);
}

export function resendEmailVerification(payload: PublicEmailRequest): Promise<EmailVerificationDispatchResponse> {
  return httpPost<EmailVerificationDispatchResponse>('/api/v1/public/email-verification/resend', payload);
}

export function verifyEmailVerification(payload: EmailVerificationVerifyRequest): Promise<{ user: SessionUser }> {
  return httpPost<{ user: SessionUser }>('/api/v1/public/email-verification/verify', payload);
}

export function updateMyProfile(payload: UpdateProfileRequest): Promise<UpdateProfileResponse> {
  return httpPut<UpdateProfileResponse>('/api/v1/private/me/profile', payload);
}
