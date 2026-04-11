import { httpPost } from '../../shared/api/http';
import type { SessionUser } from '../../app/providers/SessionProvider';

// ─── Types ────────────────────────────────────────────────────────────

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  confirm_password: string;
}

export interface LoginResponse {
  user: SessionUser;
  token: string;
  expires_in: number;
}

export interface RegisterResponse {
  user: SessionUser;
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Authenticate with email + password.
 * POST /api/v1/public/login → { data: { user, token, expires_in } }
 */
export function login(payload: LoginRequest): Promise<LoginResponse> {
  return httpPost<LoginResponse>('/api/v1/public/login', payload);
}

/**
 * Register a new customer account.
 * POST /api/v1/public/register → { data: { user } }
 */
export function register(payload: RegisterRequest): Promise<RegisterResponse> {
  return httpPost<RegisterResponse>('/api/v1/public/register', payload);
}
