import { describe, expect, it } from 'vitest';

import { getOrderedProjectMedia, getProjectMediaFull, getProjectMediaMedium, getProjectMediaThumbnail } from './projectMedia';

describe('projectMedia helpers', () => {
  it('preserves card, gallery, and lightbox variant preferences', () => {
    const media = {
      id: 'media-1',
      project_id: 'project-1',
      media_type: 'image',
      low_url: 'https://cdn.example.com/project-low.webp',
      medium_url: 'https://cdn.example.com/project-medium.webp',
      high_url: 'https://cdn.example.com/project-high.webp',
      fallback_url: 'https://cdn.example.com/project-original.jpg',
      sort_order: 0,
      featured: true,
    };

    expect(getProjectMediaThumbnail(media)).toBe(media.low_url);
    expect(getProjectMediaMedium(media)).toBe(media.medium_url);
    expect(getProjectMediaFull(media)).toBe(media.high_url);
  });

  it('rebuilds legacy image-only projects using canonical variant keys', () => {
    const media = getOrderedProjectMedia({
      images: ['https://cdn.example.com/project-image.webp'],
      media: [],
    });

    expect(media).toEqual([
      expect.objectContaining({
        low_url: 'https://cdn.example.com/project-image.webp',
        medium_url: 'https://cdn.example.com/project-image.webp',
        high_url: 'https://cdn.example.com/project-image.webp',
        fallback_url: 'https://cdn.example.com/project-image.webp',
      }),
    ]);
  });
});
