export const INDUSTRY_TYPE_MAX_LENGTH = 160;

export function normalizeEditorialMetadataText(value: string | null | undefined): string {
  return value?.trim() ?? '';
}
