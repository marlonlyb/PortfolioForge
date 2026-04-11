import { AppError } from './errors';

/**
 * Base URL resolved from Vite env at build-time.
 * Falls back to localhost:8080 for local development.
 */
const BASE_URL: string = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

/** Return the current JWT from sessionStorage, if any. */
function getAuthToken(): string | null {
  return sessionStorage.getItem('auth_token');
}

/** Build the headers map, attaching Bearer when a token exists. */
function buildHeaders(extra?: Record<string, string>): HeadersInit {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    Accept: 'application/json',
    ...extra,
  };

  const token = getAuthToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  return headers;
}

/**
 * Parse the API error envelope.
 *
 * The new contract uses: `{ "error": { code, message, details?, request_id? } }`
 * The legacy envelope uses: `{ "errors": [{ code, message }] }`
 *
 * We normalise both into `ApiErrorPayload`.
 */
function parseErrorEnvelope(status: number, body: unknown): AppError {
  if (isRecord(body)) {
    // New envelope: { error: { code, message, ... } }
    const errorObj = body['error'];
    if (isRecord(errorObj) && typeof errorObj['code'] === 'string') {
      return new AppError(status, {
        code: errorObj['code'],
        message: typeof errorObj['message'] === 'string' ? errorObj['message'] : 'Unknown error',
        details: Array.isArray(errorObj['details']) ? errorObj['details'] : undefined,
        request_id: typeof errorObj['request_id'] === 'string' ? errorObj['request_id'] : undefined,
      });
    }

    // Legacy envelope: { errors: [{ code, message }] }
    const errors = body['errors'];
    if (Array.isArray(errors) && errors.length > 0 && isRecord(errors[0])) {
      const first = errors[0];
      return new AppError(status, {
        code: typeof first['code'] === 'string' ? first['code'] : 'unexpected_error',
        message: typeof first['message'] === 'string' ? first['message'] : 'Unknown error',
      });
    }
  }

  return new AppError(status, {
    code: 'unexpected_error',
    message: 'An unexpected error occurred',
  });
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

/**
 * Core fetch wrapper.
 *
 * - Prepends BASE_URL to relative paths.
 * - Attaches Authorization header when a token exists.
 * - Parses the error envelope on non-2xx responses.
 * - Returns the raw `data` field on success (new envelope).
 * - Falls back to returning the full body for legacy endpoints.
 */
export async function http<T>(path: string, init?: RequestInit): Promise<T> {
  const url = path.startsWith('http') ? path : `${BASE_URL}${path}`;

  const response = await fetch(url, {
    ...init,
    headers: buildHeaders(init?.headers as Record<string, string> | undefined),
  });

  if (!response.ok) {
    let body: unknown;
    try {
      body = await response.json();
    } catch {
      body = null;
    }
    throw parseErrorEnvelope(response.status, body);
  }

  // 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  const json: unknown = await response.json();

  // New envelope: { data: ... }
  if (isRecord(json) && 'data' in json) {
    return json['data'] as T;
  }

  // Legacy or raw response — return as-is
  return json as T;
}

/** Convenience helpers */

export function httpGet<T>(path: string): Promise<T> {
  return http<T>(path);
}

export function httpPost<T>(path: string, body: unknown): Promise<T> {
  return http<T>(path, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export function httpPut<T>(path: string, body: unknown): Promise<T> {
  return http<T>(path, {
    method: 'PUT',
    body: JSON.stringify(body),
  });
}

export function httpPatch<T>(path: string, body: unknown): Promise<T> {
  return http<T>(path, {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}

export function httpDelete<T>(path: string): Promise<T> {
  return http<T>(path, { method: 'DELETE' });
}
