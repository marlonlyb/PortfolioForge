/**
 * API error types aligned to the backend envelope:
 * { "error": { "code", "message", "details?", "request_id?" } }
 */

export interface ApiErrorDetail {
  field?: string;
  issue: string;
}

export interface ApiErrorPayload {
  code: string;
  message: string;
  details?: ApiErrorDetail[];
  request_id?: string;
}

export class AppError extends Error {
  readonly status: number;
  readonly code: string;
  readonly details: ApiErrorDetail[];
  readonly requestId?: string;

  constructor(status: number, payload: ApiErrorPayload) {
    super(payload.message);
    this.name = 'AppError';
    this.status = status;
    this.code = payload.code;
    this.details = payload.details ?? [];
    this.requestId = payload.request_id;
  }
}

/** Recognised error codes from the API contract. */
export const API_ERROR_CODES = {
  VALIDATION_ERROR: 'validation_error',
  AUTHENTICATION_REQUIRED: 'authentication_required',
  INVALID_CREDENTIALS: 'invalid_credentials',
  FORBIDDEN: 'forbidden',
  NOT_FOUND: 'not_found',
  PRODUCT_INACTIVE: 'product_inactive',
  STOCK_INSUFFICIENT: 'stock_insufficient',
  UNEXPECTED_ERROR: 'unexpected_error',
} as const;

export type ApiErrorCode = (typeof API_ERROR_CODES)[keyof typeof API_ERROR_CODES];

/** Type guard to narrow an AppError by code. */
export function isAppErrorWithCode(err: unknown, code: ApiErrorCode): err is AppError {
  return err instanceof AppError && err.code === code;
}
