import type { ProjectProfileMetrics, ProjectProfilePrimitive } from '../../shared/types/project';

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function isPrimitiveValue(value: unknown): value is ProjectProfilePrimitive {
  return typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean';
}

function formatPrimitiveValue(value: ProjectProfilePrimitive): string {
  return typeof value === 'string' ? value.trim() : String(value);
}

export function serializeProfileList(value: unknown): string {
  if (!Array.isArray(value)) return '';

  return value
    .map((item) => {
      if (isPrimitiveValue(item)) {
        return formatPrimitiveValue(item);
      }

      if (isRecord(item)) {
        const primitiveEntries = Object.entries(item)
          .flatMap(([entryKey, entryValue]) => (isPrimitiveValue(entryValue)
            ? [`${entryKey}: ${formatPrimitiveValue(entryValue)}`]
            : []));

        return primitiveEntries.length > 0 ? primitiveEntries.join(' · ') : null;
      }

      return null;
    })
    .filter((item): item is string => Boolean(item?.trim()))
    .join('\n');
}

export function parseProfileList(value: string): string[] {
  return value
    .split('\n')
    .map((item) => item.trim())
    .filter((item) => item.length > 0);
}

export function serializeProfileMetrics(value: unknown): string {
  if (!isRecord(value)) return '';

  return Object.entries(value)
    .flatMap(([entryKey, entryValue]) => (isPrimitiveValue(entryValue)
      ? [`${entryKey}: ${formatPrimitiveValue(entryValue)}`]
      : []))
    .join('\n');
}

export function parseProfileMetrics(value: string): ProjectProfileMetrics {
  return value
    .split('\n')
    .map((line) => line.trim())
    .filter((line) => line.length > 0)
    .reduce<ProjectProfileMetrics>((metrics, line) => {
      const separatorIndex = line.indexOf(':');

      if (separatorIndex === -1) {
        throw new Error(`La línea de métrica "${line}" debe usar el formato clave: valor.`);
      }

      const key = line.slice(0, separatorIndex).trim();
      const rawValue = line.slice(separatorIndex + 1).trim();

      if (!key) {
        throw new Error('Cada métrica debe tener una clave antes de los dos puntos.');
      }

      if (!rawValue) {
        throw new Error(`La métrica "${key}" debe tener un valor después de los dos puntos.`);
      }

      metrics[key] = rawValue;
      return metrics;
    }, {});
}
