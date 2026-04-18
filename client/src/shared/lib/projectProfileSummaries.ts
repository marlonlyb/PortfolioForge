import type { ProjectProfilePrimitive, ProjectProfileStructuredList } from '../types/project';

export const PROJECT_PROFILE_LIST_KIND = {
  GENERIC: 'generic',
  INTEGRATIONS: 'integrations',
  TECHNICAL_DECISIONS: 'technicalDecisions',
  CHALLENGES: 'challenges',
  RESULTS: 'results',
  TIMELINE: 'timeline',
} as const;

export type ProjectProfileListKind = (typeof PROJECT_PROFILE_LIST_KIND)[keyof typeof PROJECT_PROFILE_LIST_KIND];

interface StructuredProfileField {
  key: string;
  value: string;
}

interface StructuredProfileItem {
  rawText: string;
  fields: StructuredProfileField[];
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function isPrimitiveValue(value: unknown): value is ProjectProfilePrimitive {
  return typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean';
}

function normalizeWhitespace(value: string): string {
  return value.replace(/\s+/g, ' ').trim();
}

function normalizeKey(value: string): string {
  return normalizeWhitespace(value).toLowerCase().replace(/[\s-]+/g, '_');
}

function formatPrimitiveValue(value: ProjectProfilePrimitive): string {
  return typeof value === 'string' ? normalizeWhitespace(value) : String(value);
}

function capitalizeFirst(value: string): string {
  if (value.length === 0) return value;
  return `${value.charAt(0).toUpperCase()}${value.slice(1)}`;
}

function ensureSentence(value: string): string {
  const normalized = normalizeWhitespace(value);
  if (normalized.length === 0) return '';
  const capitalized = capitalizeFirst(normalized);
  return /[.!?]$/.test(capitalized) ? capitalized : `${capitalized}.`;
}

function buildExcerpt(value: string, maxLength = 180): string {
  const normalized = normalizeWhitespace(value);
  if (normalized.length <= maxLength) {
    return normalized;
  }

  return `${normalized.slice(0, maxLength).trimEnd()}…`;
}

function buildFieldsFromText(value: string): StructuredProfileField[] {
  return value
    .split('|')
    .map((segment) => segment.trim())
    .flatMap((segment) => {
      const separatorIndex = segment.indexOf(':');
      if (separatorIndex <= 0) return [];

      const rawKey = segment.slice(0, separatorIndex).trim();
      const rawValue = segment.slice(separatorIndex + 1).trim();
      if (!rawKey || !rawValue) return [];
      if (!/^[A-Za-z][A-Za-z0-9_\-/ ]{0,40}$/.test(rawKey)) return [];

      return [{
        key: normalizeKey(rawKey),
        value: normalizeWhitespace(rawValue),
      }];
    });
}

function buildFieldsFromRecord(value: Record<string, unknown>): StructuredProfileField[] {
  return Object.entries(value)
    .flatMap(([entryKey, entryValue]) => (isPrimitiveValue(entryValue)
      ? [{
        key: normalizeKey(entryKey),
        value: formatPrimitiveValue(entryValue),
      }]
      : []));
}

function parseStructuredProfileItem(value: unknown): StructuredProfileItem | null {
  if (isPrimitiveValue(value)) {
    const rawText = formatPrimitiveValue(value);
    if (!rawText) return null;

    return {
      rawText,
      fields: buildFieldsFromText(rawText),
    };
  }

  if (isRecord(value)) {
    const fields = buildFieldsFromRecord(value);
    if (fields.length === 0) return null;

    return {
      rawText: fields.map((field) => field.value).join(' '),
      fields,
    };
  }

  return null;
}

function getFieldValue(item: StructuredProfileItem, ...keys: string[]): string | null {
  const normalizedKeys = keys.map((key) => normalizeKey(key));
  const field = item.fields.find((entry) => normalizedKeys.includes(entry.key));
  return field?.value ?? null;
}

function buildSummary(parts: Array<string | null | undefined>, maxLength = 180): string | null {
  const uniqueParts = parts
    .map((part) => (part ? ensureSentence(part) : null))
    .filter((part): part is string => Boolean(part))
    .filter((part, index, collection) => collection.findIndex((candidate) => candidate.toLowerCase() === part.toLowerCase()) === index);

  if (uniqueParts.length === 0) return null;
  return buildExcerpt(uniqueParts.join(' '), maxLength);
}

function buildStatusSummary(status: string | null): string | null {
  if (!status) return null;

  const normalized = normalizeKey(status);
  if (normalized.includes('partial')) return 'Quedó parcialmente resuelto.';
  if (normalized.includes('parcial')) return 'Quedó parcialmente resuelto.';
  if (normalized.includes('pend')) return 'Quedó pendiente.';
  if (normalized.includes('block')) return 'Quedó bloqueado.';
  if (normalized.includes('risk')) return 'Siguió siendo un riesgo abierto.';
  return null;
}

function summarizeStructuredProfileItem(item: StructuredProfileItem, kind: ProjectProfileListKind): string | null {
  const rawText = normalizeWhitespace(item.rawText);

  switch (kind) {
    case PROJECT_PROFILE_LIST_KIND.INTEGRATIONS:
      return buildSummary([
        getFieldValue(item, 'name') ?? rawText,
      ], 140);
    case PROJECT_PROFILE_LIST_KIND.TECHNICAL_DECISIONS:
      return buildSummary([
        getFieldValue(item, 'decision') ?? rawText,
      ], 170);
    case PROJECT_PROFILE_LIST_KIND.CHALLENGES:
      return buildSummary([
        getFieldValue(item, 'challenge') ?? rawText,
        buildStatusSummary(getFieldValue(item, 'status')),
      ], 180);
    case PROJECT_PROFILE_LIST_KIND.RESULTS:
      return buildSummary([
        getFieldValue(item, 'result') ?? rawText,
        getFieldValue(item, 'impact'),
      ], 190);
    case PROJECT_PROFILE_LIST_KIND.TIMELINE:
      return buildSummary([
        getFieldValue(item, 'phase') ?? rawText,
        getFieldValue(item, 'outcome'),
      ], 190);
    case PROJECT_PROFILE_LIST_KIND.GENERIC:
    default:
      return buildSummary(
        item.fields.length > 0
          ? item.fields.slice(0, 2).map((field) => field.value)
          : [rawText],
        180,
      );
  }
}

export function summarizeProjectProfileList(
  value: ProjectProfileStructuredList | unknown,
  kind: ProjectProfileListKind = PROJECT_PROFILE_LIST_KIND.GENERIC,
): string[] {
  if (!Array.isArray(value)) return [];

  return value
    .map((item) => parseStructuredProfileItem(item))
    .flatMap((item) => (item ? [summarizeStructuredProfileItem(item, kind)] : []))
    .filter((item): item is string => Boolean(item?.trim()));
}
