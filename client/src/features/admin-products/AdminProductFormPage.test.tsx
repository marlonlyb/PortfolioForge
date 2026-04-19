import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import {
  PUBLIC_CONTENT_FIELDS,
  PUBLIC_LOCALE,
  TRANSLATION_MODE,
  type PublicContentFieldKey,
  type TranslationMode,
} from '../../shared/i18n/config';
import type { AdminProjectDetail } from '../../shared/types/admin-project';
import type { ProjectProfileStructuredList } from '../../shared/types/project';
import { fetchAdminTechnologies } from '../admin-technologies/api';
import {
  fetchAdminProjectById,
  fetchProjectLocalizations,
  fetchProjectReadiness,
  updateAdminProject,
  updateProjectEnrichment,
} from '../admin-projects/api';
import { AdminProductFormPage } from './AdminProductFormPage';

vi.mock('../admin-projects/api', async () => {
  const actual = await vi.importActual<typeof import('../admin-projects/api')>('../admin-projects/api');
  return {
    ...actual,
    fetchAdminProjectById: vi.fn(),
    fetchProjectLocalizations: vi.fn(),
    fetchProjectReadiness: vi.fn(),
    updateAdminProject: vi.fn(),
    updateProjectEnrichment: vi.fn(),
  };
});

vi.mock('../admin-technologies/api', async () => {
  const actual = await vi.importActual<typeof import('../admin-technologies/api')>('../admin-technologies/api');
  return {
    ...actual,
    fetchAdminTechnologies: vi.fn(),
  };
});

const mockedFetchAdminProjectById = vi.mocked(fetchAdminProjectById);
const mockedFetchProjectLocalizations = vi.mocked(fetchProjectLocalizations);
const mockedFetchProjectReadiness = vi.mocked(fetchProjectReadiness);
const mockedUpdateAdminProject = vi.mocked(updateAdminProject);
const mockedUpdateProjectEnrichment = vi.mocked(updateProjectEnrichment);
const mockedFetchAdminTechnologies = vi.mocked(fetchAdminTechnologies);

const STRUCTURED_INTEGRATIONS: ProjectProfileStructuredList = [
  {
    name: 'CAN Bus',
    type: 'fieldbus',
    note: 'backbone entre la medición existente y la estación de monitoreo',
  },
  {
    name: 'Ethernet/UDP',
    purpose: 'consumo externo futuro',
  },
];

const STRUCTURED_TECHNICAL_DECISIONS: ProjectProfileStructuredList = [
  {
    decision: 'Exponer datos por Ethernet/UDP para consumo externo futuro',
    why: 'los documentos muestran preparación de interfaz',
    tradeoff: 'arquitectura lista para integrarse sin sobreactuar el rollout',
  },
];

const STRUCTURED_RESULTS: ProjectProfileStructuredList = [
  {
    result: 'Se instaló visualización de operador en dos pantallas',
    impact: 'la información pasó a supervisarse con mayor claridad',
    evidence: 'informe de instalación y propuesta de retrofit',
  },
];

function buildLocalizationFields(): Record<PublicContentFieldKey, { value: unknown; mode: TranslationMode }> {
  return PUBLIC_CONTENT_FIELDS.reduce<Record<PublicContentFieldKey, { value: unknown; mode: TranslationMode }>>((accumulator, fieldKey) => {
    accumulator[fieldKey] = { value: '', mode: TRANSLATION_MODE.AUTO };
    return accumulator;
  }, {} as Record<PublicContentFieldKey, { value: unknown; mode: TranslationMode }>);
}

function buildProjectDetail(): AdminProjectDetail {
  return {
    id: 'project-1',
    name: 'PortfolioForge',
    slug: 'portfolioforge',
    description: 'Original description',
    category: 'platform',
    industry_type: 'automatización industrial',
    final_product: 'Panel HMI para diagnóstico y monitoreo',
    images: [],
    media: [],
    variants: [],
    active: true,
    source_markdown_url: 'https://example.com/portfolioforge.md',
    profile: {
      project_id: 'project-1',
      business_goal: 'Improve observability',
      problem_statement: 'Signals were fragmented',
      solution_summary: 'Unified the monitoring flow',
      delivery_scope: 'Retrofit and rollout',
      responsibility_scope: 'Architecture and integration',
      architecture: 'Go API and PLC bridge',
      ai_usage: 'none',
      integrations: STRUCTURED_INTEGRATIONS,
      technical_decisions: STRUCTURED_TECHNICAL_DECISIONS,
      challenges: [],
      results: STRUCTURED_RESULTS,
      metrics: {
        users_impacted: 1200,
        verified: true,
      },
      timeline: [
        {
          phase: 'Installation',
          outcome: 'System delivered',
        },
      ],
      updated_at: 1710000000,
    },
    technologies: [],
  };
}

function renderPage() {
  return render(
    <MemoryRouter initialEntries={['/admin/projects/project-1']}>
      <Routes>
        <Route path="/admin/projects/:id" element={<AdminProductFormPage />} />
        <Route path="/admin/projects" element={<div>admin list</div>} />
      </Routes>
    </MemoryRouter>,
  );
}

describe('AdminProductFormPage', () => {
  beforeEach(() => {
    mockedFetchAdminTechnologies.mockResolvedValue({ items: [] });
    mockedFetchAdminProjectById.mockResolvedValue(buildProjectDetail());
    mockedFetchProjectReadiness.mockResolvedValue({
      project_id: 'project-1',
      level: 'complete',
      missing_fields: [],
      has_name: true,
      has_description: true,
      has_category: true,
      has_technologies: true,
      has_solution_summary: true,
    });
    mockedFetchProjectLocalizations.mockResolvedValue({
      project_id: 'project-1',
      base: Object.fromEntries(PUBLIC_CONTENT_FIELDS.map((fieldKey) => [fieldKey, ''])) as Record<PublicContentFieldKey, unknown>,
      locales: {
        [PUBLIC_LOCALE.CA]: { locale: PUBLIC_LOCALE.CA, fields: buildLocalizationFields() },
        [PUBLIC_LOCALE.EN]: { locale: PUBLIC_LOCALE.EN, fields: buildLocalizationFields() },
        [PUBLIC_LOCALE.DE]: { locale: PUBLIC_LOCALE.DE, fields: buildLocalizationFields() },
      },
    });
    mockedUpdateAdminProject.mockResolvedValue(buildProjectDetail());
    mockedUpdateProjectEnrichment.mockResolvedValue(undefined);
  });

  afterEach(() => {
    cleanup();
    vi.clearAllMocks();
  });

  it('preserves structured project profile arrays when saving an unrelated edit', async () => {
    renderPage();

    const integrationsField = await screen.findByDisplayValue((value) => typeof value === 'string' && value.includes('CAN Bus') && value.includes('fieldbus'));
    expect((integrationsField as HTMLTextAreaElement).value).toContain('"name": "CAN Bus"');

    fireEvent.change(screen.getByLabelText('Summary / Description'), {
      target: { value: 'Updated description only' },
    });

    fireEvent.click(screen.getByRole('button', { name: 'Update Project' }));

    await waitFor(() => {
      expect(mockedUpdateAdminProject).toHaveBeenCalledWith(
        'project-1',
        expect.objectContaining({
          description: 'Updated description only',
          industry_type: 'automatización industrial',
          final_product: 'Panel HMI para diagnóstico y monitoreo',
        }),
      );
    });

    expect(screen.getByLabelText('Industry type')).toHaveAttribute('maxlength', '160');

    await waitFor(() => {
      expect(mockedUpdateProjectEnrichment).toHaveBeenCalledWith(
        'project-1',
        expect.objectContaining({
          profile: expect.objectContaining({
            integrations: STRUCTURED_INTEGRATIONS,
            technical_decisions: STRUCTURED_TECHNICAL_DECISIONS,
            results: STRUCTURED_RESULTS,
            metrics: {
              users_impacted: 1200,
              verified: true,
            },
            timeline: [
              {
                phase: 'Installation',
                outcome: 'System delivered',
              },
            ],
          }),
        }),
      );
    });
  });
});
