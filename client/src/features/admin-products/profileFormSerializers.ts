import type {
  ProjectProfileMetrics,
  ProjectProfilePrimitive,
  ProjectProfileStructuredItem,
  ProjectProfileStructuredList,
} from '../../shared/types/project';
import type { ProjectProfileListKind } from '../../shared/lib/projectProfileSummaries';

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function isPrimitiveValue(value: unknown): value is ProjectProfilePrimitive {
  return typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean';
}

function isStructuredItem(value: unknown): value is ProjectProfileStructuredItem {
  if (isPrimitiveValue(value)) {
    return true;
  }

  return isRecord(value) && Object.values(value).every((entryValue) => isPrimitiveValue(entryValue));
}

function normalizeStructuredItem(value: unknown): ProjectProfileStructuredItem {
  if (isPrimitiveValue(value)) {
    return value;
  }

  if (!isRecord(value)) {
    throw new Error('Cada elemento debe ser un valor primitivo o un objeto plano con valores primitivos.');
  }

  const normalizedEntries = Object.entries(value).flatMap(([entryKey, entryValue]) => (
    isPrimitiveValue(entryValue) ? [[entryKey, entryValue] as const] : []
  ));

  if (normalizedEntries.length !== Object.keys(value).length) {
    throw new Error('Cada elemento debe ser un valor primitivo o un objeto plano con valores primitivos.');
  }

  return Object.fromEntries(normalizedEntries) as Record<string, ProjectProfilePrimitive>;
}

function normalizePrimitiveRecord(value: Record<string, unknown>): ProjectProfileMetrics {
  return Object.fromEntries(
    Object.entries(value).map(([entryKey, entryValue]) => [entryKey, entryValue as ProjectProfilePrimitive]),
  ) as ProjectProfileMetrics;
}

export function serializeProfileList(
  value: unknown,
  _kind?: ProjectProfileListKind,
): string {
  if (!Array.isArray(value) || value.length === 0) {
    return '';
  }

  if (!value.every((item) => isStructuredItem(item))) {
    return '';
  }

  return JSON.stringify(value, null, 2);
}

export function parseProfileList(value: string): ProjectProfileStructuredList {
  const trimmed = value.trim();
  if (trimmed.length === 0) {
    return [];
  }

  let parsed: unknown;
  try {
    parsed = JSON.parse(trimmed);
  } catch {
    throw new Error('El campo debe usar un array JSON válido para preservar la estructura original.');
  }

  if (!Array.isArray(parsed)) {
    throw new Error('El campo debe usar un array JSON válido para preservar la estructura original.');
  }

  return parsed.map((item) => normalizeStructuredItem(item));
}

export function serializeProfileMetrics(value: unknown): string {
  if (!isRecord(value)) {
    return '';
  }

  if (!Object.values(value).every((entryValue) => isPrimitiveValue(entryValue))) {
    return '';
  }

  return JSON.stringify(value, null, 2);
}

export function parseProfileMetrics(value: string): ProjectProfileMetrics {
  const trimmed = value.trim();
  if (trimmed.length === 0) {
    return {};
  }

  let parsed: unknown;
  try {
    parsed = JSON.parse(trimmed);
  } catch {
    throw new Error('El campo de métricas debe usar un objeto JSON válido para preservar tipos y claves.');
  }

  if (!isRecord(parsed) || !Object.values(parsed).every((entryValue) => isPrimitiveValue(entryValue))) {
    throw new Error('El campo de métricas debe usar un objeto JSON plano con valores primitivos.');
  }

  return normalizePrimitiveRecord(parsed);
}
