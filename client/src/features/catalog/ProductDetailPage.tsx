import { useEffect, useState } from 'react';
import { Link, useLocation, useParams } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';
import { fetchProjectBySlug } from './api';
import { searchProjects } from '../search/api';
import type { Project } from '../../shared/types/project';
import { AppError } from '../../shared/api/errors';
import {
  buildProjectSearchMatchContext,
  buildSearchMatchContext,
  formatEvidenceField,
  formatMatchType,
  hasMatchedText,
  type SearchMatchContext,
  truncateEvidenceText,
  type ProjectDetailLocationState,
} from '../search/matchContext';
import {
  getOrderedProjectMedia,
  getProjectHeroImage,
  getProjectMediaFull,
  getProjectMediaMedium,
} from '../../shared/lib/projectMedia';

interface CaseStudySectionProps {
  title: string;
  content: string;
}

interface KeyValueEntry {
  label: string;
  value: string;
}

interface DetailSectionData {
  title: string;
  content: string;
}

interface DetailListSectionData {
  title: string;
  items: string[];
  accent?: 'default' | 'highlight';
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

function formatTimestamp(timestamp?: number): string | null {
  if (!timestamp) return null;

  return new Intl.DateTimeFormat('en', {
    month: 'short',
    year: 'numeric',
  }).format(new Date(timestamp * 1000));
}

function formatTechnologySummary(technologies: Project['technologies']): string | null {
  if (!technologies || technologies.length === 0) return null;

  const [first, second, third, ...rest] = technologies;
  const visible = [first, second, third].filter((technology): technology is NonNullable<typeof technology> => Boolean(technology));
  const visibleNames = visible.map((technology) => technology.name);

  if (rest.length === 0) {
    return visibleNames.join(' · ');
  }

  return `${visibleNames.join(' · ')} +${rest.length}`;
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
    <article className="detail__story-section">
      <p className="eyebrow">{title}</p>
      <p className="detail__section-copy">{content}</p>
    </article>
  );
}

export function ProductDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const location = useLocation();
  const { locale, t } = useLocale();

  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [failedImage, setFailedImage] = useState<string | null>(null);
  const [lightboxImage, setLightboxImage] = useState<string | null>(null);
  const locationState = (location.state as ProjectDetailLocationState | null) ?? null;
  const activeSearchQuery = locationState?.activeSearchQuery?.trim() ?? '';
  const activeSearchCategory = locationState?.activeSearchCategory?.trim() ?? '';
  const locationMatchContext = locationState?.searchMatchContext;
  const [resolvedSearchMatchContext, setResolvedSearchMatchContext] = useState<SearchMatchContext | undefined>(
    locationMatchContext,
  );

  useEffect(() => {
    if (!slug) {
      setLoading(false);
        setError(t.detailNotFound);
        return;
      }

    let cancelled = false;
    setLoading(true);
    setError(null);
    setProject(null);
    setFailedImage(null);
    setLightboxImage(null);

    fetchProjectBySlug(slug, locale)
      .then((data) => {
        if (!cancelled) {
          setProject(data);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : t.detailNotFound);
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [locale, slug, t.detailNotFound]);

  useEffect(() => {
    if (locationMatchContext) {
      setResolvedSearchMatchContext(locationMatchContext);
      return;
    }

    if (!project || !slug || !activeSearchQuery) {
      setResolvedSearchMatchContext(undefined);
      return;
    }

    let cancelled = false;

    searchProjects({
      q: activeSearchQuery,
      category: activeSearchCategory || undefined,
      pageSize: 100,
      lang: locale,
    })
      .then((response) => {
        if (cancelled) return;

        const remoteMatch = response.data.find((candidate) => candidate.slug === slug);
        if (remoteMatch) {
          setResolvedSearchMatchContext(
            buildSearchMatchContext(remoteMatch) ?? buildProjectSearchMatchContext(project, activeSearchQuery),
          );
          return;
        }

        setResolvedSearchMatchContext(buildProjectSearchMatchContext(project, activeSearchQuery));
      })
      .catch(() => {
        if (cancelled) return;
        setResolvedSearchMatchContext(buildProjectSearchMatchContext(project, activeSearchQuery));
      });

    return () => {
      cancelled = true;
    };
  }, [activeSearchCategory, activeSearchQuery, locale, locationMatchContext, project, slug]);

  useEffect(() => {
    if (!lightboxImage) return undefined;

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setLightboxImage(null);
      }
    }

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [lightboxImage]);

  if (loading) {
    return (
      <section className="detail">
        <p className="catalog__loading">{t.detailLoading}</p>
      </section>
    );
  }

  if (error || !project) {
    return (
      <section className="detail">
        <div className="card card--muted">
          <p className="eyebrow">Error</p>
          <p>{error ?? t.detailNotFound}</p>
          <Link to="/">{t.detailBack}</Link>
        </div>
      </section>
    );
  }

  const technologies = project.technologies ?? [];
  const galleryMedia = getOrderedProjectMedia(project);
  const mainImage = getProjectHeroImage(project) ?? null;
  const galleryImages = galleryMedia
    .map((media) => ({
      id: media.id,
      preview: getProjectMediaMedium(media),
      full: getProjectMediaFull(media),
      alt: media.alt_text?.trim() || media.caption?.trim() || project.name,
      caption: media.caption?.trim(),
      featured: media.featured,
    }))
    .filter((item): item is { id: string; preview: string; full: string; alt: string; caption?: string; featured: boolean } =>
      Boolean(item.preview && item.full),
    );
  const visibleMainImage = failedImage && mainImage === failedImage
    ? galleryImages.find((item) => item.preview !== failedImage)?.preview ?? null
    : mainImage;
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
  const searchMatchContext = resolvedSearchMatchContext;
  const searchEvidence = searchMatchContext?.evidence ?? [];
  const hasExplanation = Boolean(searchMatchContext?.explanation?.trim());
  const hasEvidence = searchEvidence.length > 0;
  const showSearchMatchContext = hasExplanation || hasEvidence;
  const lastUpdated = formatTimestamp(project.updated_at);
  const technologySummary = formatTechnologySummary(technologies);
  const overviewItems = [
    { label: t.detailCategory, value: project.category },
    { label: t.detailClient, value: project.client_name ?? t.detailIndependent },
    { label: t.detailUpdated, value: lastUpdated ?? t.detailRecentlyCurated },
    { label: t.detailTechnologies, value: technologySummary ?? t.detailNotSpecified },
  ].filter((item) => item.value.trim().length > 0);
  const narrativeSections: DetailSectionData[] = [
    hasText(businessGoal) ? { title: t.detailBusinessGoal, content: businessGoal } : null,
    hasText(problemStatement) ? { title: t.detailProblem, content: problemStatement } : null,
    hasText(solutionSummary) ? { title: t.detailSolution, content: solutionSummary } : null,
    hasText(architecture) ? { title: t.detailArchitecture, content: architecture } : null,
    hasText(aiUsage) ? { title: t.detailAIUsage, content: aiUsage } : null,
  ].filter((section): section is DetailSectionData => Boolean(section));
  const sidebarSections: DetailListSectionData[] = [
    integrations.length > 0 ? { title: t.detailIntegrations, items: integrations } : null,
    technicalDecisions.length > 0 ? { title: t.detailTechnicalDecisions, items: technicalDecisions } : null,
    challenges.length > 0 ? { title: t.detailChallenges, items: challenges } : null,
    results.length > 0 ? { title: t.detailResults, items: results, accent: 'highlight' } : null,
    timeline.length > 0 ? { title: t.detailTimeline, items: timeline } : null,
  ].filter((section): section is DetailListSectionData => Boolean(section));

  return (
    <>
      <section className="detail">
        <Link className="detail__back" to="/">
          {t.detailBack}
        </Link>

        <article className="detail__hero card">
          <div className="detail__hero-content">
            {project.category ? <p className="eyebrow">{project.category}</p> : null}
            <h2 className="detail__title">{project.name}</h2>
            {hasText(project.client_name) ? (
              <p className="detail__context">{t.detailClientContext} · {project.client_name}</p>
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
            {visibleMainImage ? (
              <img
                className="detail__hero-image"
                src={visibleMainImage}
                alt={project.name}
                onError={() => setFailedImage(visibleMainImage)}
              />
            ) : (
              <div className="detail__hero-image detail__hero-image--placeholder">
                {t.detailVisualUnavailable}
              </div>
            )}
          </div>
        </article>

        {galleryImages.length > 0 ? (
          <article className="detail__section card detail__gallery-section">
            <div className="detail__gallery-header">
              <div>
                <p className="eyebrow">Gallery</p>
                <h3 className="detail__gallery-title">Project visuals</h3>
              </div>
            </div>

            <div className="detail__gallery-grid">
              {galleryImages.map((image) => (
                <button
                  key={image.id}
                  type="button"
                  className={image.featured ? 'detail__gallery-item detail__gallery-item--featured' : 'detail__gallery-item'}
                  onClick={() => setLightboxImage(image.full)}
                >
                  <img className="detail__gallery-image" src={image.preview} alt={image.alt} loading="lazy" />
                  {image.caption ? <span className="detail__gallery-caption">{image.caption}</span> : null}
                </button>
              ))}
            </div>
          </article>
        ) : null}

        <article className="detail__overview-strip card">
          <div className="detail__overview-strip-header">
            <p className="eyebrow">{t.detailProjectOverview}</p>
          </div>

          <dl className="detail__overview-list detail__overview-list--strip">
            {overviewItems.map((item) => (
              <div key={item.label} className="detail__overview-item">
                <dt>{item.label}</dt>
                <dd>{item.value}</dd>
              </div>
            ))}
          </dl>
        </article>

        <div className="detail__content-layout">
          <div className="detail__main-column">
            <div className="detail__column-intro">
              <p className="eyebrow">Case study</p>
              <h3 className="detail__column-title">Narrativa principal</h3>
            </div>

            <div className="detail__main-stack">
              {narrativeSections.map((section) => (
                <CaseStudySection key={section.title} title={section.title} content={section.content} />
              ))}
            </div>
          </div>

          <aside className="detail__side-column">
            <div className="detail__column-intro detail__column-intro--side">
              <p className="eyebrow">Evidence</p>
              <h3 className="detail__column-title">Analítica y soporte</h3>
            </div>

            <div className="detail__side-stack">
              {showSearchMatchContext ? (
                <article className="detail__aside-panel detail__match-context card" aria-label="Por qué coincide con la búsqueda">
                  <p className="eyebrow">Por qué coincide</p>

                  {hasExplanation && searchMatchContext?.explanation ? (
                    <p className="detail__match-explanation">{searchMatchContext.explanation}</p>
                  ) : null}

                  {hasEvidence ? (
                    <div className="detail__match-evidence">
                      <p className="detail__match-evidence-title">Evidencia utilizada</p>
                      <ul className="detail__match-evidence-list">
                        {searchEvidence.map((evidence, index) => (
                          <li
                            key={`${project.id}-${evidence.field}-${evidence.match_type}-${index}`}
                            className="detail__match-evidence-item"
                          >
                            <div className="detail__match-evidence-meta">
                              <span className="detail__match-evidence-field">{formatEvidenceField(evidence.field)}</span>
                              <span className="detail__match-evidence-type">{formatMatchType(evidence.match_type)}</span>
                            </div>

                            {hasMatchedText(evidence) ? (
                              <p className="detail__match-evidence-text">“{truncateEvidenceText(evidence.matched_text)}”</p>
                            ) : null}
                          </li>
                        ))}
                      </ul>
                    </div>
                  ) : null}
                </article>
              ) : null}

              {sidebarSections.map((section) => (
                <article
                  key={section.title}
                  className={section.accent === 'highlight'
                    ? 'detail__aside-panel detail__aside-panel--results card'
                    : 'detail__aside-panel card'}
                >
                  <p className="eyebrow">{section.title}</p>
                  <ul className={section.accent === 'highlight' ? 'detail__list detail__list--results' : 'detail__list'}>
                    {section.items.map((item, index) => (
                      <li key={`${section.title}-${index}`}>{item}</li>
                    ))}
                  </ul>
                </article>
              ))}

              {metrics.length > 0 ? (
                <article className="detail__aside-panel detail__aside-panel--metrics card">
                  <p className="eyebrow">{t.detailMetrics}</p>
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
            </div>
          </aside>
        </div>
        </section>
      {lightboxImage ? (
        <div className="detail__lightbox" role="dialog" aria-modal="true" aria-label="Image preview" onClick={() => setLightboxImage(null)}>
          <button type="button" className="detail__lightbox-close" aria-label="Close image preview" onClick={() => setLightboxImage(null)}>
            ×
          </button>
          <img className="detail__lightbox-image" src={lightboxImage} alt={project.name} onClick={(event) => event.stopPropagation()} />
        </div>
      ) : null}
    </>
  );
}
