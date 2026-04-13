import { useEffect, useState, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import {
  createProduct,
  fetchAdminProductById,
  updateProduct,
  fetchProjectReadiness,
  reembedProject,
  updateProjectEnrichment,
} from './api';
import type {
  CreateProductPayload,
  UpdateProductPayload,
  ProjectReadiness,
  UpdateProjectEnrichmentProfilePayload,
} from './api';
import { fetchAdminTechnologies } from '../admin-technologies/api';
import type { Technology } from '../../shared/types/project';
import { AppError } from '../../shared/api/errors';
import {
  parseProfileList,
  parseProfileMetrics,
  serializeProfileList,
  serializeProfileMetrics,
} from './profileFormSerializers';

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
  const [images, setImages] = useState('');
  const [active, setActive] = useState(true);

  // Enrichment
  const [availableTechnologies, setAvailableTechnologies] = useState<Technology[]>([]);
  const [selectedTechIds, setSelectedTechIds] = useState<string[]>([]);
  const [businessGoal, setBusinessGoal] = useState('');
  const [problemStatement, setProblemStatement] = useState('');
  const [solutionSummary, setSolutionSummary] = useState('');
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

    fetchAdminProductById(id)
      .then((project) => {
        if (!cancelled) {
          setName(project.name);
          setDescription(project.description);
          setCategory(project.category);
          setBrand(project.brand ?? '');
          setImages(project.images.join(', '));
          setActive(project.active);

          if (project.profile) {
            setBusinessGoal(project.profile.business_goal ?? '');
            setProblemStatement(project.profile.problem_statement ?? '');
            setSolutionSummary(project.profile.solution_summary ?? '');
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

    const imageList = images
      .split(',')
      .map((item) => item.trim())
      .filter((item) => item.length > 0);

    try {
      const enrichmentProfile: UpdateProjectEnrichmentProfilePayload = {
        business_goal: businessGoal || undefined,
        problem_statement: problemStatement || undefined,
        solution_summary: solutionSummary || undefined,
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
        const payload: UpdateProductPayload = {
          name,
          description,
          category,
          brand: brand || undefined,
          images: imageList,
          active,
        };
        await updateProduct(id, payload);
      } else {
        const payload: CreateProductPayload = {
          name,
          description,
          category,
          brand: brand || undefined,
          images: imageList,
          active,
        };
        const created = await createProduct(payload);
        projectId = created.id;
      }

      if (projectId) {
        await updateProjectEnrichment(projectId, {
          profile: enrichmentProfile,
          technology_ids: selectedTechIds,
        });
      }

      if (isEdit) {
        refreshReadiness();
      }

      navigate('/admin/projects');
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
              Main images (comma-separated URLs)
              <input
                className="admin__input"
                type="text"
                placeholder="https://example.com/hero.jpg, https://example.com/secondary.jpg"
                value={images}
                onChange={(e) => setImages(e.target.value)}
              />
            </label>

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
                placeholder={'Microservicios para checkout\nWorkers asíncronos para webhooks'}
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
                placeholder={'Menor tiempo de respuesta\nMayor conversión en checkout'}
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
