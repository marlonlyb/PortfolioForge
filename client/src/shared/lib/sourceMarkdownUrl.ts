const CANONICAL_MARKDOWN_BASE_URL = 'https://mlbautomation.com/dev/portfolioforge';

export function slugifyProjectName(value: string): string {
  let normalized = value.trim().toLowerCase().replaceAll(' ', '-');
  while (normalized.includes('--')) {
    normalized = normalized.replaceAll('--', '-');
  }
  return normalized.replace(/^-+|-+$/g, '');
}

export function buildCanonicalMarkdownURLFromSlug(slug: string): string {
  const normalizedSlug = slugifyProjectName(slug);
  if (!normalizedSlug) {
    return '';
  }
  return `${CANONICAL_MARKDOWN_BASE_URL}/${normalizedSlug}/${normalizedSlug}.md`;
}

export function buildCanonicalMarkdownURLFromName(name: string): string {
  return buildCanonicalMarkdownURLFromSlug(name);
}
