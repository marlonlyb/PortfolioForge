import { describe, expect, it } from 'vitest';

import {
  applyDefaultProjectMediaUrls,
  buildDefaultProjectMediaUrls,
  PROJECT_MEDIA_FALLBACK_URL,
} from './projectMediaUrlDefaults';

describe('projectMediaUrlDefaults', () => {
  it('builds zero-padded default project media URLs from slug', () => {
    expect(buildDefaultProjectMediaUrls('can-bus-crane-monitoring', 1)).toEqual({
      low_url: 'https://mlbautomation.com/dev/portfolioforge/can-bus-crane-monitoring/imagen01_low.webp',
      medium_url: 'https://mlbautomation.com/dev/portfolioforge/can-bus-crane-monitoring/imagen01_medium.webp',
      high_url: 'https://mlbautomation.com/dev/portfolioforge/can-bus-crane-monitoring/imagen01_high.webp',
      fallback_url: PROJECT_MEDIA_FALLBACK_URL,
    });
  });

  it('returns fallback-only defaults when slug cannot be inferred yet', () => {
    expect(buildDefaultProjectMediaUrls('   ', 1)).toEqual({
      low_url: '',
      medium_url: '',
      high_url: '',
      fallback_url: PROJECT_MEDIA_FALLBACK_URL,
    });
  });

  it('fills blank media fields with generated defaults', () => {
    expect(applyDefaultProjectMediaUrls({ low_url: '', medium_url: '', high_url: '', fallback_url: '' }, 'Printer 05 Controls Migration', 2)).toEqual({
      low_url: 'https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/imagen02_low.webp',
      medium_url: 'https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/imagen02_medium.webp',
      high_url: 'https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/imagen02_high.webp',
      fallback_url: PROJECT_MEDIA_FALLBACK_URL,
    });
  });

  it('updates generated URLs when the project slug changes', () => {
    expect(
      applyDefaultProjectMediaUrls(
        {
          low_url: 'https://mlbautomation.com/dev/portfolioforge/old-slug/imagen01_low.webp',
          medium_url: 'https://mlbautomation.com/dev/portfolioforge/old-slug/imagen01_medium.webp',
          high_url: 'https://mlbautomation.com/dev/portfolioforge/old-slug/imagen01_high.webp',
          fallback_url: PROJECT_MEDIA_FALLBACK_URL,
        },
        'new-slug',
        1,
        'old-slug',
      ),
    ).toEqual({
      low_url: 'https://mlbautomation.com/dev/portfolioforge/new-slug/imagen01_low.webp',
      medium_url: 'https://mlbautomation.com/dev/portfolioforge/new-slug/imagen01_medium.webp',
      high_url: 'https://mlbautomation.com/dev/portfolioforge/new-slug/imagen01_high.webp',
      fallback_url: PROJECT_MEDIA_FALLBACK_URL,
    });
  });

  it('preserves manual overrides', () => {
    expect(
      applyDefaultProjectMediaUrls(
        {
          low_url: 'https://cdn.example.com/custom-low.webp',
          medium_url: 'https://cdn.example.com/custom-medium.webp',
          high_url: 'https://cdn.example.com/custom-high.webp',
          fallback_url: 'https://cdn.example.com/custom-fallback.png',
        },
        'can-bus-crane-monitoring',
        1,
      ),
    ).toEqual({
      low_url: 'https://cdn.example.com/custom-low.webp',
      medium_url: 'https://cdn.example.com/custom-medium.webp',
      high_url: 'https://cdn.example.com/custom-high.webp',
      fallback_url: 'https://cdn.example.com/custom-fallback.png',
    });
  });
});
