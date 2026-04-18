import { useEffect, useState, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import {
  createAdminProject,
  fetchAdminProjectById,
  fetchProjectLocalizations,
  fetchProjectReadiness,
  reembedProject,
  saveProjectLocalizations,
  updateAdminProject,
  updateProjectEnrichment,
} from '../admin-projects/api';
import type {
  AdminProjectLocalizationsResponse,
  CreateAdminProjectPayload,
  ProjectReadiness,
  UpdateAdminProjectPayload,
  UpdateProjectEnrichmentProfilePayload,
} from '../admin-projects/api';
import { fetchAdminTechnologies } from '../admin-technologies/api';
import type { ProjectMedia, Technology } from '../../shared/types/project';
import { AppError } from '../../shared/api/errors';
import {
  parseProfileList,
  parseProfileMetrics,
  serializeProfileList,
  serializeProfileMetrics,
} from './profileFormSerializers';
import { buildCanonicalMarkdownURLFromName, buildCanonicalMarkdownURLFromSlug } from '../../shared/lib/sourceMarkdownUrl';
import {
  PUBLIC_CONTENT_FIELDS,
  PUBLIC_LOCALE,
  TRANSLATION_MODE,
  type PublicContentFieldKey,
  type PublicLocale,
} from '../../shared/i18n/config';

const TRANSLATABLE_LOCALES = [PUBLIC_LOCALE.CA, PUBLIC_LOCALE.EN, PUBLIC_LOCALE.DE] as const;

type MediaFormItem = ProjectMedia;

function createEmptyMediaItem(sortOrder = 0): MediaFormItem {
  return {
    id: `new-${sortOrder}-${Date.now()}`,
    project_id: '',
    media_type: 'image',
    low_url: '',
    medium_url: '',
    high_url: '',
    fallback_url: '',
    caption: '',
    alt_text: '',
    sort_order: sortOrder,
    featured: sortOrder === 0,
  };
}

function createFallbackMediaItems(images: string[]): MediaFormItem[] {
  const items = images
    .filter((image) => image.trim().length > 0)
    .map((image, index) => ({
      ...createEmptyMediaItem(index),
      fallback_url: image,
      low_url: image,
      medium_url: image,
      high_url: image,
      featured: index === 0,
    }));

  return items.length > 0 ? items : [createEmptyMediaItem(0)];
}

function normalizeMediaItems(items: MediaFormItem[]): MediaFormItem[] {
  const normalized = items
    .map((item, index) => ({
      ...item,
      media_type: item.media_type?.trim() || 'image',
      fallback_url: item.fallback_url?.trim() || '',
      low_url: item.low_url?.trim() || '',
      medium_url: item.medium_url?.trim() || '',
      high_url: item.high_url?.trim() || '',
      caption: item.caption?.trim() || '',
      alt_text: item.alt_text?.trim() || '',
      sort_order: Number.isFinite(item.sort_order) ? item.sort_order : index,
    }))
    .filter((item) => item.low_url || item.medium_url || item.high_url || item.fallback_url);

  if (normalized.length === 0) {
    return [];
  }

  const featuredIndex = normalized.findIndex((item) => item.featured);

  return normalized.map((item, index) => ({
    ...item,
    featured: featuredIndex >= 0 ? featuredIndex === index : index === 0,
    sort_order: item.sort_order,
  }));
}

function buildLegacyImagesFromMedia(items: MediaFormItem[]): string[] {
  return normalizeMediaItems(items)
    .map((item) => item.medium_url || item.high_url || item.low_url || item.fallback_url || '')
    .filter((item) => item.length > 0);
}

const TRANSLATION_FIELD_LABELS: Record<PublicContentFieldKey, string> = {
  name: 'Título',
  description: 'Descripción',
  category: 'Categoría',
  client_name: 'Cliente / Contexto',
  business_goal: 'Business Goal',
  problem_statement: 'Problem Statement',
  solution_summary: 'Solution Summary',
  delivery_scope: 'Delivery Scope',
  responsibility_scope: 'Responsibility Scope',
  architecture: 'Architecture',
  ai_usage: 'AI Usage',
  integrations: 'Integrations',
  technical_decisions: 'Technical Decisions',
  challenges: 'Challenges',
  results: 'Results',
  metrics: 'Metrics',
  timeline: 'Timeline',
};

const TRANSLATION_LOCALE_LABELS: Record<(typeof TRANSLATABLE_LOCALES)[number], string> = {
  ca: 'Català',
  en: 'English',
  de: 'Deutsch',
};

type TranslationEditorState = Record<PublicContentFieldKey, string>;
type TranslationModesState = Record<PublicContentFieldKey, string>;

function createEmptyTranslationEditorState(): TranslationEditorState {
  return PUBLIC_CONTENT_FIELDS.reduce<TranslationEditorState>((accumulator, fieldKey) => {
    accumulator[fieldKey] = '';
    return accumulator;
  }, {} as TranslationEditorState);
}

function createEmptyTranslationModesState(): TranslationModesState {
  return PUBLIC_CONTENT_FIELDS.reduce<TranslationModesState>((accumulator, fieldKey) => {
    accumulator[fieldKey] = TRANSLATION_MODE.AUTO;
    return accumulator;
  }, {} as TranslationModesState);
}

function serializeTranslationField(fieldKey: PublicContentFieldKey, value: unknown): string {
  switch (fieldKey) {
    case 'integrations':
    case 'technical_decisions':
    case 'challenges':
    case 'results':
    case 'timeline':
      return serializeProfileList(value);
    case 'metrics':
      return serializeProfileMetrics(value);
    default:
      return typeof value === 'string' ? value : '';
  }
}

function parseTranslationField(fieldKey: PublicContentFieldKey, value: string): unknown {
  switch (fieldKey) {
    case 'integrations':
    case 'technical_decisions':
    case 'challenges':
    case 'results':
    case 'timeline':
      return parseProfileList(value);
    case 'metrics':
      return parseProfileMetrics(value);
    default:
      return value.trim();
  }
}

function buildTranslationEditorState(response: AdminProjectLocalizationsResponse, locale: PublicLocale): TranslationEditorState {
  const localeData = response.locales[locale];
  return PUBLIC_CONTENT_FIELDS.reduce<TranslationEditorState>((accumulator, fieldKey) => {
    accumulator[fieldKey] = serializeTranslationField(fieldKey, localeData?.fields[fieldKey]?.value ?? response.base[fieldKey]);
    return accumulator;
  }, {} as TranslationEditorState);
}

function buildTranslationModesState(response: AdminProjectLocalizationsResponse, locale: PublicLocale): TranslationModesState {
  const localeData = response.locales[locale];
  return PUBLIC_CONTENT_FIELDS.reduce<TranslationModesState>((accumulator, fieldKey) => {
    accumulator[fieldKey] = localeData?.fields[fieldKey]?.mode ?? TRANSLATION_MODE.AUTO;
    return accumulator;
  }, {} as TranslationModesState);
}

// ─── Readiness badge helpers ──────────────────────────────────────────

const READINESS_LABELS: Record<ProjectReadiness['level'], string> = {
  incomplete: 'Incompleto',
  basic: 'Básico',
  complete: 'Completo',
};

const MISSING_FIELD_PROMPTS: Record<string, string> = {
  description: 'Agrega una descripción para que este proyecto aparezca en búsquedas.',
  category: 'Agrega una categoría para mejorar la búsqueda por filtros.',
  technologies: 'Vincula tecnologías para mejorar la búsqueda fuzzy y semántica.',
  solution_summary:
    'Agrega un resumen de solución para mejorar la búsqueda de este proyecto.',
};

function ReadinessBadge({ readiness }: { readiness: ProjectReadiness | null }) {
  if (!readiness) return null;

  const missingFields = readiness.missing_fields ?? [];

  const levelClass =
    readiness.level === 'complete'
      ? 'readiness-badge--complete'
      : readiness.level === 'basic'
        ? 'readiness-badge--basic'
        : 'readiness-badge--incomplete';

  return (
    <div className="readiness-badge">
      <span className={`readiness-badge__dot ${levelClass}`} />
      <span className="readiness-badge__label">
        Disponibilidad de búsqueda: <strong>{READINESS_LABELS[readiness.level]}</strong>
      </span>
      {missingFields.length > 0 && (
        <ul className="readiness-badge__missing">
          {missingFields.map((field) => (
            <li key={field}>{MISSING_FIELD_PROMPTS[field] ?? `Agrega ${field} para mejorar la búsqueda.`}</li>
          ))}
        </ul>
      )}
    </div>
  );
}

function MissingFieldHint({ fieldName, readiness }: { fieldName: string; readiness: ProjectReadiness | null }) {
  if (!readiness) return null;

  const missingFields = readiness.missing_fields ?? [];
  if (!missingFields.includes(fieldName)) return null;

  return (
    <p className="admin__field-hint">{MISSING_FIELD_PROMPTS[fieldName] ?? `Agrega ${fieldName} para mejorar la búsqueda.`}</p>
  );
}

// ─── Main component ───────────────────────────────────────────────────

export function AdminProductFormPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEdit = Boolean(id);

  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('');
  const [brand, setBrand] = useState('');
  const [sourceMarkdownURL, setSourceMarkdownURL] = useState('');
  const [sourceMarkdownURLManuallyEdited, setSourceMarkdownURLManuallyEdited] = useState(false);
  const [mediaItems, setMediaItems] = useState<MediaFormItem[]>([createEmptyMediaItem(0)]);
  const [active, setActive] = useState(true);

  // Enrichment
  const [availableTechnologies, setAvailableTechnologies] = useState<Technology[]>([]);
  const [selectedTechIds, setSelectedTechIds] = useState<string[]>([]);
  const [businessGoal, setBusinessGoal] = useState('');
  const [problemStatement, setProblemStatement] = useState('');
  const [solutionSummary, setSolutionSummary] = useState('');
  const [deliveryScope, setDeliveryScope] = useState('');
  const [responsibilityScope, setResponsibilityScope] = useState('');
  const [architecture, setArchitecture] = useState('');
  const [aiUsage, setAiUsage] = useState('');
  const [integrations, setIntegrations] = useState('');
  const [technicalDecisions, setTechnicalDecisions] = useState('');
  const [challenges, setChallenges] = useState('');
  const [results, setResults] = useState('');
  const [metrics, setMetrics] = useState('');
  const [timeline, setTimeline] = useState('');

  // Search readiness
  const [readiness, setReadiness] = useState<ProjectReadiness | null>(null);
  const [reembedLoading, setReembedLoading] = useState(false);
  const [reembedMessage, setReembedMessage] = useState<string | null>(null);
  const [activeTranslationLocale, setActiveTranslationLocale] = useState<(typeof TRANSLATABLE_LOCALES)[number]>(PUBLIC_LOCALE.CA);
  const [translationEditors, setTranslationEditors] = useState<Record<(typeof TRANSLATABLE_LOCALES)[number], TranslationEditorState>>({
    ca: createEmptyTranslationEditorState(),
    en: createEmptyTranslationEditorState(),
    de: createEmptyTranslationEditorState(),
  });
  const [translationModes, setTranslationModes] = useState<Record<(typeof TRANSLATABLE_LOCALES)[number], TranslationModesState>>({
    ca: createEmptyTranslationModesState(),
    en: createEmptyTranslationModesState(),
    de: createEmptyTranslationModesState(),
  });
  const [translationSnapshot, setTranslationSnapshot] = useState<Record<(typeof TRANSLATABLE_LOCALES)[number], TranslationEditorState>>({
    ca: createEmptyTranslationEditorState(),
    en: createEmptyTranslationEditorState(),
    de: createEmptyTranslationEditorState(),
  });

  useEffect(() => {
    let cancelled = false;

    fetchAdminTechnologies()
      .then((res) => {
        if (!cancelled) setAvailableTechnologies(res.items);
      })
      .catch(() => {
        // Silently ignore or log
      });

    if (!id) {
      setLoading(false);
      return () => {
        cancelled = true;
      };
    }

    fetchAdminProjectById(id)
      .then((project) => {
        if (!cancelled) {
          const inferredMarkdownURL = buildCanonicalMarkdownURLFromSlug(project.slug || project.name);
          const explicitMarkdownURL = project.source_markdown_url?.trim() ?? '';

          setName(project.name);
          setDescription(project.description);
          setCategory(project.category);
          setBrand(project.brand ?? '');
          setSourceMarkdownURL(explicitMarkdownURL || inferredMarkdownURL);
          setSourceMarkdownURLManuallyEdited(Boolean(explicitMarkdownURL) && explicitMarkdownURL !== inferredMarkdownURL);
          setMediaItems(project.media && project.media.length > 0 ? project.media : createFallbackMediaItems(project.images ?? []));
          setActive(project.active);

          if (project.profile) {
            setBusinessGoal(project.profile.business_goal ?? '');
            setProblemStatement(project.profile.problem_statement ?? '');
            setSolutionSummary(project.profile.solution_summary ?? '');
            setDeliveryScope(project.profile.delivery_scope ?? '');
            setResponsibilityScope(project.profile.responsibility_scope ?? '');
            setArchitecture(project.profile.architecture ?? '');
            setAiUsage(project.profile.ai_usage ?? '');
            setIntegrations(serializeProfileList(project.profile.integrations));
            setTechnicalDecisions(serializeProfileList(project.profile.technical_decisions));
            setChallenges(serializeProfileList(project.profile.challenges));
            setResults(serializeProfileList(project.profile.results));
            setMetrics(serializeProfileMetrics(project.profile.metrics));
            setTimeline(serializeProfileList(project.profile.timeline));
          }

          if (project.technologies) {
            setSelectedTechIds(project.technologies.map((t) => t.id));
          }

          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Failed to load project.');
          setLoading(false);
        }
      });

    fetchProjectLocalizations(id)
      .then((response) => {
        if (cancelled) return;

        setTranslationEditors({
          ca: buildTranslationEditorState(response, PUBLIC_LOCALE.CA),
          en: buildTranslationEditorState(response, PUBLIC_LOCALE.EN),
          de: buildTranslationEditorState(response, PUBLIC_LOCALE.DE),
        });
        setTranslationModes({
          ca: buildTranslationModesState(response, PUBLIC_LOCALE.CA),
          en: buildTranslationModesState(response, PUBLIC_LOCALE.EN),
          de: buildTranslationModesState(response, PUBLIC_LOCALE.DE),
        });
        setTranslationSnapshot({
          ca: buildTranslationEditorState(response, PUBLIC_LOCALE.CA),
          en: buildTranslationEditorState(response, PUBLIC_LOCALE.EN),
          de: buildTranslationEditorState(response, PUBLIC_LOCALE.DE),
        });
      })
      .catch(() => {
        // ignore translation errors for now
      });

    // Fetch readiness in parallel
    fetchProjectReadiness(id)
      .then((data) => {
        if (!cancelled) setReadiness(data);
      })
      .catch(() => {
        // Readiness is optional — silently ignore errors
      });

    return () => {
      cancelled = true;
    };
  }, [id]);

  useEffect(() => {
    if (sourceMarkdownURLManuallyEdited) {
      return;
    }

    setSourceMarkdownURL(buildCanonicalMarkdownURLFromName(name));
  }, [name, sourceMarkdownURLManuallyEdited]);

  // Refresh readiness after save
  function refreshReadiness() {
    if (!id) return;
    fetchProjectReadiness(id)
      .then((data) => setReadiness(data))
      .catch(() => {
        // Silently ignore
      });
  }

  async function handleReembed() {
    if (!id) return;
    if (!window.confirm('¿Actualizar el documento de búsqueda para este proyecto?')) return;

    setReembedLoading(true);
    setReembedMessage(null);

    try {
      const result = await reembedProject(id);
      setReembedMessage(result.message);
    } catch (err: unknown) {
      setReembedMessage(
        err instanceof AppError ? err.message : 'Error al actualizar documento de búsqueda.',
      );
    } finally {
      setReembedLoading(false);
    }
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setSubmitting(true);

    const normalizedMedia = normalizeMediaItems(mediaItems);
    const imageList = buildLegacyImagesFromMedia(normalizedMedia);

    try {
      const enrichmentProfile: UpdateProjectEnrichmentProfilePayload = {
        business_goal: businessGoal || undefined,
        problem_statement: problemStatement || undefined,
        solution_summary: solutionSummary || undefined,
        delivery_scope: deliveryScope || undefined,
        responsibility_scope: responsibilityScope || undefined,
        architecture: architecture || undefined,
        ai_usage: aiUsage || undefined,
        integrations: parseProfileList(integrations),
        technical_decisions: parseProfileList(technicalDecisions),
        challenges: parseProfileList(challenges),
        results: parseProfileList(results),
        metrics: parseProfileMetrics(metrics),
        timeline: parseProfileList(timeline),
      };

      let projectId = id;

      if (isEdit && id) {
        const payload: UpdateAdminProjectPayload = {
          name,
          description,
          category,
          brand: brand || undefined,
          source_markdown_url: sourceMarkdownURL.trim(),
          images: imageList,
          media: normalizedMedia,
          active,
        };
        await updateAdminProject(id, payload);
      } else {
        const payload: CreateAdminProjectPayload = {
          name,
          description,
          category,
          brand: brand || undefined,
          source_markdown_url: sourceMarkdownURL.trim(),
          images: imageList,
          media: normalizedMedia,
          active,
        };
        const created = await createAdminProject(payload);
        projectId = created.id;
      }

      if (projectId) {
        await updateProjectEnrichment(projectId, {
          profile: enrichmentProfile,
          technology_ids: selectedTechIds,
        });

        if (isEdit) {
          for (const locale of TRANSLATABLE_LOCALES) {
            const fields = translationEditors[locale];
            const snapshot = translationSnapshot[locale];
            const changedFields = Object.entries(fields).reduce<Partial<Record<PublicContentFieldKey, unknown>>>((accumulator, [fieldKey, fieldValue]) => {
              const typedFieldKey = fieldKey as PublicContentFieldKey;
              if (fieldValue.trim() === snapshot[typedFieldKey].trim()) {
                return accumulator;
              }
              accumulator[typedFieldKey] = parseTranslationField(typedFieldKey, fieldValue);
              return accumulator;
            }, {});

            if (Object.keys(changedFields).length > 0) {
              await saveProjectLocalizations(projectId, locale, { fields: changedFields });
            }
          }
        }
      }

      if (isEdit) {
        refreshReadiness();
      }

      navigate(isEdit ? '/admin/projects' : `/admin/projects/${projectId}`);
    } catch (err: unknown) {
      setError(err instanceof AppError || err instanceof Error ? err.message : 'Failed to save project.');
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading project…</p>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <Link className="detail__back" to="/admin/projects">
          ← Back to projects
        </Link>

        <p className="eyebrow">Admin</p>
        <h2>{isEdit ? 'Edit Project' : 'New Project'}</h2>

        {error ? <div className="admin__error" role="alert">{error}</div> : null}

        {isEdit && <ReadinessBadge readiness={readiness} />}

        <form className="admin__form" onSubmit={handleSubmit}>
          <div className="admin__form-section">
            <h3>Project profile</h3>

            <label className="admin__label">
              Title
              <input
                className="admin__input"
                type="text"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Summary / Description
              <textarea
                className="admin__textarea"
                required
                rows={6}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </label>
            <MissingFieldHint fieldName="description" readiness={readiness} />

            <div className="admin__form-row">
              <label className="admin__label">
                Category
                <input
                  className="admin__input"
                  type="text"
                  required
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                />
              </label>
              <MissingFieldHint fieldName="category" readiness={readiness} />

              <label className="admin__label">
                Client / Context
                <input
                  className="admin__input"
                  type="text"
                  value={brand}
                  onChange={(e) => setBrand(e.target.value)}
                />
              </label>
            </div>

            <label className="admin__label">
              Markdown source URL (private)
              <input
                className="admin__input"
                type="url"
                placeholder="https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md"
                value={sourceMarkdownURL}
                onChange={(event) => {
                  setSourceMarkdownURL(event.target.value);
                  setSourceMarkdownURLManuallyEdited(true);
                }}
              />
            </label>
            <p className="admin__field-hint">
              Se completa automáticamente desde el slug usando la convención remota oficial. Puedes sobrescribirla o dejarla vacía para deshabilitar el chat público.
            </p>

            <div className="admin__label">
              <div className="admin__media-header">
                <span>Optimized media</span>
                <button
                  className="btn btn--ghost"
                  type="button"
                  onClick={() => setMediaItems((prev) => [...prev, createEmptyMediaItem(prev.length)])}
                >
                  Add image
                </button>
              </div>

              <span className="admin__field-hint">
                Usa tres variantes por imagen: _low para catálogo, _medium para galería y _high para vista ampliada. Marca una imagen como principal.
              </span>

              <div className="admin__media-stack">
                {mediaItems.map((item, index) => (
                  <article key={item.id || `media-${index}`} className="admin__media-card">
                    <div className="admin__media-card-header">
                      <strong>Image #{index + 1}</strong>
                      <button
                        className="btn btn--ghost"
                        type="button"
                        disabled={mediaItems.length === 1}
                        onClick={() =>
                          setMediaItems((prev) => {
                            const next = prev.filter((_, currentIndex) => currentIndex !== index);
                            return next.length > 0 ? next : [createEmptyMediaItem(0)];
                          })
                        }
                      >
                        Remove
                      </button>
                    </div>

                    <div className="admin__form-row">
                      <label className="admin__label">
                        Sort order
                        <input
                          className="admin__input"
                          type="number"
                          value={item.sort_order}
                          onChange={(event) => {
                            const value = Number(event.target.value);
                            setMediaItems((prev) => prev.map((entry, entryIndex) => (
                              entryIndex === index ? { ...entry, sort_order: Number.isFinite(value) ? value : index } : entry
                            )));
                          }}
                        />
                      </label>

                      <label className="admin__label admin__label--checkbox">
                        <input
                          type="checkbox"
                          checked={item.featured}
                          onChange={() =>
                            setMediaItems((prev) => prev.map((entry, entryIndex) => ({
                              ...entry,
                              featured: entryIndex === index,
                            })))
                          }
                        />
                        Featured image
                      </label>
                    </div>

                    <label className="admin__label">
                      Low / catálogo URL
                      <input
                        className="admin__input"
                        type="url"
                        placeholder="https://cdn.example.com/project_low.webp"
                        value={item.low_url ?? ''}
                        onChange={(event) => setMediaItems((prev) => prev.map((entry, entryIndex) => (
                          entryIndex === index ? { ...entry, low_url: event.target.value } : entry
                        )))}
                      />
                    </label>

                    <label className="admin__label">
                      Medium / galería URL
                      <input
                        className="admin__input"
                        type="url"
                        placeholder="https://cdn.example.com/project_medium.webp"
                        value={item.medium_url ?? ''}
                        onChange={(event) => setMediaItems((prev) => prev.map((entry, entryIndex) => (
                          entryIndex === index ? { ...entry, medium_url: event.target.value } : entry
                        )))}
                      />
                    </label>

                    <label className="admin__label">
                      High / ampliada URL
                      <input
                        className="admin__input"
                        type="url"
                        placeholder="https://cdn.example.com/project_high.webp"
                        value={item.high_url ?? ''}
                        onChange={(event) => setMediaItems((prev) => prev.map((entry, entryIndex) => (
                          entryIndex === index ? { ...entry, high_url: event.target.value } : entry
                        )))}
                      />
                    </label>

                    <label className="admin__label">
                      Fallback URL (optional)
                      <input
                        className="admin__input"
                        type="url"
                        placeholder="https://cdn.example.com/project-original.jpg"
                        value={item.fallback_url ?? ''}
                        onChange={(event) => setMediaItems((prev) => prev.map((entry, entryIndex) => (
                          entryIndex === index ? { ...entry, fallback_url: event.target.value } : entry
                        )))}
                      />
                    </label>

                    <label className="admin__label">
                      Caption
                      <input
                        className="admin__input"
                        type="text"
                        value={item.caption ?? ''}
                        onChange={(event) => setMediaItems((prev) => prev.map((entry, entryIndex) => (
                          entryIndex === index ? { ...entry, caption: event.target.value } : entry
                        )))}
                      />
                    </label>

                    <label className="admin__label">
                      Alt text
                      <input
                        className="admin__input"
                        type="text"
                        value={item.alt_text ?? ''}
                        onChange={(event) => setMediaItems((prev) => prev.map((entry, entryIndex) => (
                          entryIndex === index ? { ...entry, alt_text: event.target.value } : entry
                        )))}
                      />
                    </label>
                  </article>
                ))}
              </div>
            </div>

            <label className="admin__label admin__label--checkbox">
              <input type="checkbox" checked={active} onChange={(e) => setActive(e.target.checked)} />
              Published
            </label>

            <MissingFieldHint fieldName="solution_summary" readiness={readiness} />
            <MissingFieldHint fieldName="technologies" readiness={readiness} />
          </div>

          <div className="admin__form-section">
            <h3>Rich Profile (Search Enrichment)</h3>

            <label className="admin__label">
              Technologies
              <select
                className="admin__input"
                multiple
                size={5}
                value={selectedTechIds}
                onChange={(e) => {
                  const options = Array.from(e.target.selectedOptions, (option) => option.value);
                  setSelectedTechIds(options);
                }}
              >
                {availableTechnologies.map((tech) => (
                  <option key={tech.id} value={tech.id}>
                    {tech.name} ({tech.category})
                  </option>
                ))}
              </select>
              <span className="admin__field-hint">Hold Ctrl/Cmd to select multiple.</span>
            </label>

            <label className="admin__label">
              Business Goal
              <textarea
                className="admin__textarea"
                rows={3}
                value={businessGoal}
                onChange={(e) => setBusinessGoal(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Problem Statement
              <textarea
                className="admin__textarea"
                rows={3}
                value={problemStatement}
                onChange={(e) => setProblemStatement(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Solution Summary
              <textarea
                className="admin__textarea"
                rows={4}
                value={solutionSummary}
                onChange={(e) => setSolutionSummary(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Delivery Scope
              <textarea
                className="admin__textarea"
                rows={3}
                value={deliveryScope}
                onChange={(e) => setDeliveryScope(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Responsibility Scope
              <textarea
                className="admin__textarea"
                rows={3}
                value={responsibilityScope}
                onChange={(e) => setResponsibilityScope(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Architecture
              <textarea
                className="admin__textarea"
                rows={4}
                value={architecture}
                onChange={(e) => setArchitecture(e.target.value)}
              />
            </label>

            <label className="admin__label">
              AI Usage
              <textarea
                className="admin__textarea"
                rows={3}
                value={aiUsage}
                onChange={(e) => setAiUsage(e.target.value)}
              />
            </label>

            <label className="admin__label">
              Integrations
              <textarea
                className="admin__textarea"
                rows={4}
                placeholder={'Stripe\nHubSpot\nOpenAI API'}
                value={integrations}
                onChange={(e) => setIntegrations(e.target.value)}
              />
              <span className="admin__field-hint">Una línea por integración.</span>
            </label>

            <label className="admin__label">
              Technical Decisions
              <textarea
                className="admin__textarea"
                rows={4}
                placeholder={'Arquitectura orientada a eventos\nWorkers asíncronos para integraciones'}
                value={technicalDecisions}
                onChange={(e) => setTechnicalDecisions(e.target.value)}
              />
              <span className="admin__field-hint">Una línea por decisión técnica.</span>
            </label>

            <label className="admin__label">
              Challenges
              <textarea
                className="admin__textarea"
                rows={4}
                placeholder={'Latencia en integraciones externas\nNormalización de datos legacy'}
                value={challenges}
                onChange={(e) => setChallenges(e.target.value)}
              />
              <span className="admin__field-hint">Una línea por desafío.</span>
            </label>

            <label className="admin__label">
              Results
              <textarea
                className="admin__textarea"
                rows={4}
                placeholder={'Menor tiempo de respuesta\nMayor activación de usuarios'}
                value={results}
                onChange={(e) => setResults(e.target.value)}
              />
              <span className="admin__field-hint">Una línea por resultado.</span>
            </label>

            <label className="admin__label">
              Metrics
              <textarea
                className="admin__textarea"
                rows={5}
                placeholder={'conversion_rate: +18%\nresponse_time: -42%\nusers_impacted: 1200'}
                value={metrics}
                onChange={(e) => setMetrics(e.target.value)}
              />
              <span className="admin__field-hint">Formato: clave: valor, una métrica por línea.</span>
            </label>

            <label className="admin__label">
              Timeline
              <textarea
                className="admin__textarea"
                rows={4}
                placeholder={'Discovery\nMVP\nRollout\nOptimización'}
                value={timeline}
                onChange={(e) => setTimeline(e.target.value)}
              />
              <span className="admin__field-hint">Una línea por hito o etapa.</span>
            </label>
          </div>

          <div className="admin__form-section admin__translations-section">
            <div className="admin__translations-header">
              <div>
                <h3>Traducciones persistidas</h3>
                <p className="admin__field-hint">
                  Castellano sigue siendo la fuente. Los campos en otros idiomas muestran si están en modo automático o manual.
                </p>
              </div>
              <div className="admin__translation-tabs">
                {TRANSLATABLE_LOCALES.map((locale) => (
                  <button
                    key={locale}
                    className={activeTranslationLocale === locale ? 'admin__translation-tab admin__translation-tab--active' : 'admin__translation-tab'}
                    type="button"
                    onClick={() => setActiveTranslationLocale(locale)}
                  >
                    {TRANSLATION_LOCALE_LABELS[locale]}
                  </button>
                ))}
              </div>
            </div>

            {!isEdit ? (
              <p className="admin__field-hint">Guarda primero el proyecto en castellano para generar traducciones automáticas.</p>
            ) : (
              <div className="admin__translation-grid">
                {PUBLIC_CONTENT_FIELDS.map((fieldKey) => {
                  const value = translationEditors[activeTranslationLocale][fieldKey];
                  const baseMode = translationModes[activeTranslationLocale][fieldKey];
                  const isDirty = value.trim() !== translationSnapshot[activeTranslationLocale][fieldKey].trim();
                  const effectiveMode = isDirty ? TRANSLATION_MODE.MANUAL : baseMode;
                  const isLargeField = !['name', 'category', 'client_name'].includes(fieldKey);

                  return (
                    <label key={`${activeTranslationLocale}-${fieldKey}`} className="admin__label">
                      <span className="admin__translation-label-row">
                        <span>{TRANSLATION_FIELD_LABELS[fieldKey]}</span>
                        <span className={effectiveMode === TRANSLATION_MODE.MANUAL ? 'admin__translation-mode admin__translation-mode--manual' : 'admin__translation-mode admin__translation-mode--auto'}>
                          {effectiveMode === TRANSLATION_MODE.MANUAL ? 'manual' : 'auto'}
                        </span>
                      </span>
                      {isLargeField ? (
                        <textarea
                          className="admin__textarea"
                          rows={fieldKey === 'description' ? 5 : 4}
                          value={value}
                          onChange={(event) =>
                            setTranslationEditors((prev) => ({
                              ...prev,
                              [activeTranslationLocale]: {
                                ...prev[activeTranslationLocale],
                                [fieldKey]: event.target.value,
                              },
                            }))
                          }
                        />
                      ) : (
                        <input
                          className="admin__input"
                          type="text"
                          value={value}
                          onChange={(event) =>
                            setTranslationEditors((prev) => ({
                              ...prev,
                              [activeTranslationLocale]: {
                                ...prev[activeTranslationLocale],
                                [fieldKey]: event.target.value,
                              },
                            }))
                          }
                        />
                      )}
                    </label>
                  );
                })}
              </div>
            )}
          </div>

          <div className="admin__form-actions">
            <button className="btn btn--primary" type="submit" disabled={submitting}>
              {submitting ? 'Saving…' : isEdit ? 'Update Project' : 'Create Project'}
            </button>

            {isEdit && (
              <button
                className="btn btn--ghost"
                type="button"
                disabled={reembedLoading}
                onClick={handleReembed}
              >
                {reembedLoading ? 'Actualizando búsqueda…' : 'Actualizar búsqueda'}
              </button>
            )}
          </div>

          {reembedMessage && (
            <p className="admin__reembed-msg" role="status">{reembedMessage}</p>
          )}
        </form>
      </article>
    </section>
  );
}
