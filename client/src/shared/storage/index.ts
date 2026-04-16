/**
 * Storage wrappers for frontend session state.
 *
 * These centralise access so the rest of the app never touches
 * the raw Storage API directly.
 */

// ─── Auth token (sessionStorage) ─────────────────────────────────────

const AUTH_TOKEN_KEY = 'auth_token';
const ASSISTANT_HISTORY_PREFIX = 'assistant_history:';

export function getToken(): string | null {
  return sessionStorage.getItem(AUTH_TOKEN_KEY);
}

export function setToken(token: string): void {
  sessionStorage.setItem(AUTH_TOKEN_KEY, token);
}

export function removeToken(): void {
  sessionStorage.removeItem(AUTH_TOKEN_KEY);
}

export interface StoredAssistantMessage {
  role: 'user' | 'assistant';
  content: string;
}

export function getAssistantHistory(projectSlug: string): StoredAssistantMessage[] {
  const rawValue = sessionStorage.getItem(`${ASSISTANT_HISTORY_PREFIX}${projectSlug}`);
  if (!rawValue) {
    return [];
  }

  try {
    const parsed = JSON.parse(rawValue) as unknown;
    if (!Array.isArray(parsed)) {
      return [];
    }

    return parsed.flatMap((entry) => {
      if (
        typeof entry === 'object'
        && entry !== null
        && (entry as { role?: unknown }).role
        && (entry as { content?: unknown }).content
      ) {
        const role = (entry as { role: unknown }).role;
        const content = (entry as { content: unknown }).content;
        if ((role === 'user' || role === 'assistant') && typeof content === 'string') {
          return [{ role, content } satisfies StoredAssistantMessage];
        }
      }

      return [];
    });
  } catch {
    return [];
  }
}

export function setAssistantHistory(projectSlug: string, history: StoredAssistantMessage[]): void {
  sessionStorage.setItem(`${ASSISTANT_HISTORY_PREFIX}${projectSlug}`, JSON.stringify(history));
}
