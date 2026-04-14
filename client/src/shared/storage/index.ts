/**
 * Storage wrappers for frontend session state.
 *
 * These centralise access so the rest of the app never touches
 * the raw Storage API directly.
 */

// ─── Auth token (sessionStorage) ─────────────────────────────────────

const AUTH_TOKEN_KEY = 'auth_token';

export function getToken(): string | null {
  return sessionStorage.getItem(AUTH_TOKEN_KEY);
}

export function setToken(token: string): void {
  sessionStorage.setItem(AUTH_TOKEN_KEY, token);
}

export function removeToken(): void {
  sessionStorage.removeItem(AUTH_TOKEN_KEY);
}
