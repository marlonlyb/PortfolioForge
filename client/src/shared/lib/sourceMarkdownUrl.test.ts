import { describe, expect, it } from 'vitest';

import { buildCanonicalMarkdownURLFromName, buildCanonicalMarkdownURLFromSlug, slugifyProjectName } from './sourceMarkdownUrl';

describe('sourceMarkdownUrl helpers', () => {
  it('slugifies project names using the same simple convention as runtime', () => {
    expect(slugifyProjectName('  CAN Bus Crane Monitoring  ')).toBe('can-bus-crane-monitoring');
    expect(slugifyProjectName('Printer 05   Controls Migration')).toBe('printer-05-controls-migration');
  });

  it('builds canonical markdown URLs from slug', () => {
    expect(buildCanonicalMarkdownURLFromSlug('can-bus-crane-monitoring')).toBe(
      'https://mlbautomation.com/dev/portfolioforge/can-bus-crane-monitoring/can-bus-crane-monitoring.md',
    );
  });

  it('builds canonical markdown URLs from project name', () => {
    expect(buildCanonicalMarkdownURLFromName('Printer 05 Controls Migration')).toBe(
      'https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/printer-05-controls-migration.md',
    );
  });

  it('returns empty string when slug cannot be inferred', () => {
    expect(buildCanonicalMarkdownURLFromSlug('   ')).toBe('');
  });
});
