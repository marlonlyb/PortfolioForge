import type { Project, ProjectMedia } from '../types/project';

function firstNonEmpty(...values: Array<string | undefined>): string | undefined {
  return values.find((value) => Boolean(value?.trim()))?.trim();
}

export function getProjectMediaThumbnail(media: ProjectMedia): string | undefined {
  return firstNonEmpty(media.low_url, media.medium_url, media.high_url, media.fallback_url);
}

export function getProjectMediaMedium(media: ProjectMedia): string | undefined {
  return firstNonEmpty(media.medium_url, media.high_url, media.low_url, media.fallback_url);
}

export function getProjectMediaFull(media: ProjectMedia): string | undefined {
  return firstNonEmpty(media.high_url, media.medium_url, media.low_url, media.fallback_url);
}

export function getOrderedProjectMedia(project: Pick<Project, 'images' | 'media'>): ProjectMedia[] {
  const media = [...(project.media ?? [])].sort((left, right) => {
    if (left.featured !== right.featured) return left.featured ? -1 : 1;
    return left.sort_order - right.sort_order;
  });

  if (media.length > 0) return media;

  return (project.images ?? []).map((image, index) => ({
    id: `legacy-${index}`,
    project_id: '',
    media_type: 'image',
    fallback_url: image,
    low_url: image,
    medium_url: image,
    high_url: image,
    sort_order: index,
    featured: index === 0,
  }));
}

export function getProjectCardImage(project: Pick<Project, 'images' | 'media'>): string | undefined {
  const featured = getOrderedProjectMedia(project)[0];
  return featured ? getProjectMediaThumbnail(featured) : project.images?.[0];
}

export function getProjectHeroImage(project: Pick<Project, 'images' | 'media'>): string | undefined {
  const featured = getOrderedProjectMedia(project)[0];
  return featured ? firstNonEmpty(getProjectMediaMedium(featured), getProjectMediaFull(featured)) : project.images?.[0];
}
