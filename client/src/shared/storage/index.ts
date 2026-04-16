/**
 * Storage wrappers for frontend session state.
 *
 * These centralise access so the rest of the app never touches
 * the raw Storage API directly.
 */

// ─── Auth token (sessionStorage) ─────────────────────────────────────

const AUTH_TOKEN_KEY = 'auth_token';
const ASSISTANT_HISTORY_PREFIX = 'assistant_history:';

export const ASSISTANT_HISTORY_LIMIT = 8;

const ASSISTANT_MESSAGE_ROLE = {
  USER: 'user',
  ASSISTANT: 'assistant',
} as const;

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
  role: (typeof ASSISTANT_MESSAGE_ROLE)[keyof typeof ASSISTANT_MESSAGE_ROLE];
  content: string;
}

function sanitizeAssistantMessage(entry: unknown): StoredAssistantMessage | null {
  if (typeof entry !== 'object' || entry === null) {
    return null;
  }

  const role = 'role' in entry ? entry.role : undefined;
  const content = 'content' in entry ? entry.content : undefined;

  if (
    (role !== ASSISTANT_MESSAGE_ROLE.USER && role !== ASSISTANT_MESSAGE_ROLE.ASSISTANT)
    || typeof content !== 'string'
  ) {
    return null;
  }

  const trimmedContent = content.trim();
  if (trimmedContent.length === 0) {
    return null;
  }

  return { role, content: trimmedContent };
}

export function normalizeAssistantHistory(history: StoredAssistantMessage[]): StoredAssistantMessage[] {
  if (history.length === 0) {
    return [];
  }

  const normalized = history.flatMap((entry) => {
    const sanitized = sanitizeAssistantMessage(entry);
    return sanitized ? [sanitized] : [];
  });

  if (normalized.length <= ASSISTANT_HISTORY_LIMIT) {
    return normalized;
  }

  return normalized.slice(-ASSISTANT_HISTORY_LIMIT);
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

    return normalizeAssistantHistory(parsed as StoredAssistantMessage[]);
  } catch {
    return [];
  }
}

export function setAssistantHistory(projectSlug: string, history: StoredAssistantMessage[]): void {
  const normalized = normalizeAssistantHistory(history);
  if (normalized.length === 0) {
    sessionStorage.removeItem(`${ASSISTANT_HISTORY_PREFIX}${projectSlug}`);
    return;
  }

  sessionStorage.setItem(`${ASSISTANT_HISTORY_PREFIX}${projectSlug}`, JSON.stringify(normalized));
}

export function clearAssistantHistory(projectSlug: string): void {
  sessionStorage.removeItem(`${ASSISTANT_HISTORY_PREFIX}${projectSlug}`);
}
