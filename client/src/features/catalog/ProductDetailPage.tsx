import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';

import { fetchProjectBySlug } from './api';
import type { Project } from '../../shared/types/project';
import { AppError } from '../../shared/api/errors';

interface CaseStudySectionProps {
  title: string;
  content: string;
}

interface KeyValueEntry {
  label: string;
  value: string;
}

function hasText(value?: string | null): value is string {
  return Boolean(value?.trim());
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function isPrimitiveValue(value: unknown): value is string | number | boolean {
  return typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean';
}

function formatLabel(value: string): string {
  return value
    .replace(/[_-]+/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function formatPrimitiveValue(value: string | number | boolean): string {
  return typeof value === 'string' ? value.trim() : String(value);
}

function getRenderableList(value: unknown): string[] {
  if (!Array.isArray(value)) return [];

  return value
    .map((item) => {
      if (isPrimitiveValue(item)) {
        return formatPrimitiveValue(item);
      }

      if (isRecord(item)) {
        const primitiveEntries = Object.entries(item)
          .filter(([, entryValue]) => isPrimitiveValue(entryValue))
          .map(([entryKey, entryValue]) => `${formatLabel(entryKey)}: ${formatPrimitiveValue(entryValue)}`);

        return primitiveEntries.length > 0 ? primitiveEntries.join(' · ') : null;
      }

      return null;
    })
    .filter((item): item is string => Boolean(item?.trim()));
}

function getRenderableEntries(value: unknown): KeyValueEntry[] {
  if (!isRecord(value)) return [];

  return Object.entries(value)
    .filter(([, entryValue]) => isPrimitiveValue(entryValue))
    .map(([entryKey, entryValue]) => ({
      label: formatLabel(entryKey),
      value: formatPrimitiveValue(entryValue),
    }))
    .filter((entry) => entry.value.trim().length > 0);
}

function CaseStudySection({ title, content }: CaseStudySectionProps) {
  return (
    <article className="detail__section card">
      <p className="eyebrow">{title}</p>
      <p className="detail__section-copy">{content}</p>
    </article>
  );
}

export function ProductDetailPage() {
  const { slug } = useParams<{ slug: string }>();

  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [failedImage, setFailedImage] = useState<string | null>(null);

  useEffect(() => {
    if (!slug) {
      setLoading(false);
      setError('Project not found.');
      return;
    }

    let cancelled = false;
    setLoading(true);
    setError(null);
    setProject(null);
    setFailedImage(null);

    fetchProjectBySlug(slug)
      .then((data) => {
        if (!cancelled) {
          setProject(data);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Could not load the project.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [slug]);

  if (loading) {
    return (
      <section className="detail">
        <p className="catalog__loading">Loading project…</p>
      </section>
    );
  }

  if (error || !project) {
    return (
      <section className="detail">
        <div className="card card--muted">
          <p className="eyebrow">Error</p>
          <p>{error ?? 'Project not found.'}</p>
          <Link to="/">Back to projects</Link>
        </div>
      </section>
    );
  }

  const technologies = project.technologies ?? [];
  const mainImage = project.images.find((image) => image !== failedImage) ?? null;
  const businessGoal = project.profile?.business_goal?.trim();
  const problemStatement = project.profile?.problem_statement?.trim();
  const solutionSummary = project.profile?.solution_summary?.trim();
  const architecture = project.profile?.architecture?.trim();
  const aiUsage = project.profile?.ai_usage?.trim();
  const integrations = getRenderableList(project.profile?.integrations);
  const technicalDecisions = getRenderableList(project.profile?.technical_decisions);
  const challenges = getRenderableList(project.profile?.challenges);
  const results = getRenderableList(project.profile?.results);
  const timeline = getRenderableList(project.profile?.timeline);
  const metrics = getRenderableEntries(project.profile?.metrics);

  return (
    <section className="detail">
      <Link className="detail__back" to="/">
        &larr; Back to projects
      </Link>

      <article className="detail__hero card">
        <div className="detail__hero-content">
          {project.category ? <p className="eyebrow">{project.category}</p> : null}
          <h2 className="detail__title">{project.name}</h2>
          {hasText(project.client_name) ? (
            <p className="detail__context">Client context · {project.client_name}</p>
          ) : null}
          <p className="detail__summary">{project.description}</p>

          {technologies.length > 0 ? (
            <div className="detail__chips" aria-label="Technologies used">
              {technologies.map((technology) => (
                <span key={technology.id} className="detail__chip">
                  {technology.name}
                </span>
              ))}
            </div>
          ) : null}
        </div>

        <div className="detail__hero-media">
          {mainImage ? (
            <img
              className="detail__hero-image"
              src={mainImage}
              alt={project.name}
              onError={() => setFailedImage(mainImage)}
            />
          ) : (
            <div className="detail__hero-image detail__hero-image--placeholder">
              Case study visual unavailable
            </div>
          )}
        </div>
      </article>

      <div className="detail__sections">
        {hasText(businessGoal) ? (
          <CaseStudySection title="Business context / Goal" content={businessGoal} />
        ) : null}

        {hasText(problemStatement) ? (
          <CaseStudySection title="Problem" content={problemStatement} />
        ) : null}

        {hasText(solutionSummary) ? (
          <CaseStudySection title="Solution" content={solutionSummary} />
        ) : null}

        {hasText(architecture) ? (
          <CaseStudySection title="Architecture" content={architecture} />
        ) : null}

        {hasText(aiUsage) ? (
          <CaseStudySection title="AI Usage" content={aiUsage} />
        ) : null}

        {integrations.length > 0 ? (
          <article className="detail__section card">
            <p className="eyebrow">Integrations</p>
            <ul className="detail__list">
              {integrations.map((item, index) => (
                <li key={`integration-${index}`}>{item}</li>
              ))}
            </ul>
          </article>
        ) : null}

        {technicalDecisions.length > 0 ? (
          <article className="detail__section card">
            <p className="eyebrow">Technical decisions</p>
            <ul className="detail__list">
              {technicalDecisions.map((item, index) => (
                <li key={`decision-${index}`}>{item}</li>
              ))}
            </ul>
          </article>
        ) : null}

        {challenges.length > 0 ? (
          <article className="detail__section card">
            <p className="eyebrow">Challenges</p>
            <ul className="detail__list">
              {challenges.map((item, index) => (
                <li key={`challenge-${index}`}>{item}</li>
              ))}
            </ul>
          </article>
        ) : null}

        {results.length > 0 ? (
          <article className="detail__section card">
            <p className="eyebrow">Results</p>
            <ul className="detail__list">
              {results.map((item, index) => (
                <li key={`result-${index}`}>{item}</li>
              ))}
            </ul>
          </article>
        ) : null}

        {metrics.length > 0 ? (
          <article className="detail__section card">
            <p className="eyebrow">Metrics</p>
            <dl className="detail__metrics">
              {metrics.map((metric) => (
                <div key={metric.label} className="detail__metric">
                  <dt>{metric.label}</dt>
                  <dd>{metric.value}</dd>
                </div>
              ))}
            </dl>
          </article>
        ) : null}

        {timeline.length > 0 ? (
          <article className="detail__section card">
            <p className="eyebrow">Timeline</p>
            <ul className="detail__list">
              {timeline.map((item, index) => (
                <li key={`timeline-${index}`}>{item}</li>
              ))}
            </ul>
          </article>
        ) : null}
      </div>
    </section>
  );
}
