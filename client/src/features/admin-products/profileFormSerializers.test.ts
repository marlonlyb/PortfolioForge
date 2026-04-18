import { describe, expect, it } from 'vitest';

import { parseProfileList, parseProfileMetrics, serializeProfileList, serializeProfileMetrics } from './profileFormSerializers';

describe('profileFormSerializers', () => {
  it('round-trips structured profile lists without collapsing object fields', () => {
    const value = [
      {
        name: 'CAN Bus',
        type: 'fieldbus',
        note: 'backbone entre la medición existente y la estación de monitoreo',
      },
      'Pantallas de operador para supervisión en piso',
    ];

    const serialized = serializeProfileList(value);

    expect(serialized).toContain('"name": "CAN Bus"');
    expect(parseProfileList(serialized)).toEqual(value);
  });

  it('round-trips metrics preserving primitive types', () => {
    const value = {
      users_impacted: 1200,
      verified: true,
      note: 'retrofit completado',
    };

    const serialized = serializeProfileMetrics(value);

    expect(serialized).toContain('"users_impacted": 1200');
    expect(parseProfileMetrics(serialized)).toEqual(value);
  });
});
