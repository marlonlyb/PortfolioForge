import { slugifyProjectName } from './sourceMarkdownUrl';

export const PROJECT_MEDIA_BASE_URL = 'https://mlbautomation.com/dev/portfolioforge';
export const PROJECT_MEDIA_FALLBACK_URL = `${PROJECT_MEDIA_BASE_URL}/imagen_fallback/Logo_500_500.png`;

type MediaUrlFields = {
  low_url?: string;
  medium_url?: string;
  high_url?: string;
  fallback_url?: string;
};

export function buildDefaultProjectMediaUrls(slugOrName: string, imageNumber: number): Required<MediaUrlFields> {
  const slug = slugifyProjectName(slugOrName);
  const suffix = `imagen${String(Math.max(1, imageNumber)).padStart(2, '0')}`;

  if (!slug) {
    return {
      low_url: '',
      medium_url: '',
      high_url: '',
      fallback_url: PROJECT_MEDIA_FALLBACK_URL,
    };
  }

  const base = `${PROJECT_MEDIA_BASE_URL}/${slug}/${suffix}`;
  return {
    low_url: `${base}_low.webp`,
    medium_url: `${base}_medium.webp`,
    high_url: `${base}_high.webp`,
    fallback_url: PROJECT_MEDIA_FALLBACK_URL,
  };
}

export function applyDefaultProjectMediaUrls<T extends MediaUrlFields>(
  item: T,
  slugOrName: string,
  imageNumber: number,
  previousSlugOrName = '',
): T {
  const currentDefaults = buildDefaultProjectMediaUrls(slugOrName, imageNumber);
  const previousDefaults = buildDefaultProjectMediaUrls(previousSlugOrName, imageNumber);

  return {
    ...item,
    low_url: shouldReplaceWithDefault(item.low_url, previousDefaults.low_url) ? currentDefaults.low_url : (item.low_url ?? ''),
    medium_url: shouldReplaceWithDefault(item.medium_url, previousDefaults.medium_url)
      ? currentDefaults.medium_url
      : (item.medium_url ?? ''),
    high_url: shouldReplaceWithDefault(item.high_url, previousDefaults.high_url) ? currentDefaults.high_url : (item.high_url ?? ''),
    fallback_url: shouldReplaceWithDefault(item.fallback_url, previousDefaults.fallback_url)
      ? currentDefaults.fallback_url
      : (item.fallback_url ?? ''),
  };
}

function shouldReplaceWithDefault(currentValue: string | undefined, previousDefault: string): boolean {
  const normalized = currentValue?.trim() ?? '';
  return normalized === '' || normalized === previousDefault;
}
