import { useEffect, useState } from 'react';
import { Link, useLocation, useOutletContext, useParams } from 'react-router-dom';
import useEmblaCarousel from 'embla-carousel-react';

import { useLocale } from '../../app/providers/LocaleProvider';
import { useSession } from '../../app/providers/SessionProvider';
import type { StoreHeaderContent, StoreLayoutOutletContext } from '../../app/layouts/StoreLayout';
import { fetchProjectBySlug } from './api';
import { searchProjects } from '../search/api';
import type { Project } from '../../shared/types/project';
import { AppError } from '../../shared/api/errors';
import { fetchAdminProjectById } from '../admin-projects/api';
import { ProjectAssistantChat } from './ProjectAssistantChat';
import {
  buildSearchResultsPath,
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
  getProjectMediaFull,
  getProjectMediaMedium,
} from '../../shared/lib/projectMedia';
import {
  PROJECT_PROFILE_LIST_KIND,
  summarizeProjectProfileList,
  type ProjectProfileListKind,
} from '../../shared/lib/projectProfileSummaries';

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

interface DetailLayerData {
  title: string;
  textSections: DetailSectionData[];
  listSections?: DetailListSectionData[];
  metrics?: KeyValueEntry[];
  metricsTitle?: string;
}

interface GalleryImage {
  id: string;
  preview: string;
  full: string;
  alt: string;
  caption?: string;
  featured: boolean;
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

function formatTimestamp(timestamp: number | undefined, locale: string): string | null {
  if (!timestamp) return null;

  return new Intl.DateTimeFormat(locale, {
    month: 'short',
    year: 'numeric',
  }).format(new Date(timestamp * 1000));
}

function buildProjectHeaderContent(
  project: Project,
  fallbackCaption: string,
): StoreHeaderContent {
  const captionParts = [project.category, project.client_name]
    .map((value) => value?.trim() ?? '')
    .filter((value) => value.length > 0);

  return {
    title: project.name.trim(),
    summary: '',
    caption: captionParts.join(' · ') || fallbackCaption,
  };
}

function buildExcerpt(value?: string | null, maxLength = 280): string {
	const trimmed = value?.trim() ?? '';
	if (trimmed.length <= maxLength) {
		return trimmed;
	}
	return `${trimmed.slice(0, maxLength).trimEnd()}…`;
}

function getRenderableList(value: unknown, kind: ProjectProfileListKind = PROJECT_PROFILE_LIST_KIND.GENERIC): string[] {
  return summarizeProjectProfileList(value, kind);
}

function getRenderableEntries(value: unknown): KeyValueEntry[] {
  if (!isRecord(value)) return [];

  return Object.entries(value)
    .flatMap(([entryKey, entryValue]) => (isPrimitiveValue(entryValue)
      ? [{
        label: formatLabel(entryKey),
        value: formatPrimitiveValue(entryValue),
      }]
      : []))
    .filter((entry) => entry.value.trim().length > 0);
}

function CaseStudySection({ title, content, eyebrow }: CaseStudySectionProps & { eyebrow: string }) {
  return (
    <article className="detail__story-section">
      <p className="eyebrow">{eyebrow}</p>
      <h4 className="detail__section-title">{title}</h4>
      <p className="detail__section-copy">{content}</p>
    </article>
  );
}

function DetailLayer({ title, textSections, listSections = [], metrics = [], metricsTitle, caseStudyEyebrow }: DetailLayerData & { caseStudyEyebrow: string }) {
  if (textSections.length === 0 && listSections.length === 0 && metrics.length === 0) {
    return null;
  }

  return (
    <article className="card detail__layer">
      <div className="detail__column-intro">
        <p className="eyebrow">{caseStudyEyebrow}</p>
        <h3 className="detail__column-title">{title}</h3>
      </div>

        <div className="detail__main-stack">
          {textSections.map((section) => (
          <CaseStudySection key={section.title} title={section.title} content={buildExcerpt(section.content, 300)} eyebrow={caseStudyEyebrow} />
        ))}

        {listSections.map((section) => (
          <article
            key={section.title}
            className={section.accent === 'highlight'
              ? 'detail__aside-panel detail__aside-panel--results'
              : 'detail__aside-panel'}
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
          <article className="detail__aside-panel detail__aside-panel--metrics">
            <p className="eyebrow">{metricsTitle}</p>
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
    </article>
  );
}

export function ProductDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const location = useLocation();
  const { locale, t } = useLocale();
  const { user, loading: sessionLoading } = useSession();
  const outletContext = useOutletContext<StoreLayoutOutletContext | undefined>();
  const setHeaderContent = outletContext?.setHeaderContent;
  const [galleryViewportRef, galleryViewportApi] = useEmblaCarousel({
    align: 'start',
    loop: false,
    dragFree: false,
    containScroll: 'trimSnaps',
  });

  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedGalleryIndex, setSelectedGalleryIndex] = useState(0);
  const [lightboxIndex, setLightboxIndex] = useState<number | null>(null);
  const [canScrollGalleryPrev, setCanScrollGalleryPrev] = useState(false);
  const [canScrollGalleryNext, setCanScrollGalleryNext] = useState(false);
  const locationState = (location.state as ProjectDetailLocationState | null) ?? null;
  const activeSearchQuery = locationState?.activeSearchQuery?.trim() ?? '';
  const activeSearchCategory = locationState?.activeSearchCategory?.trim() ?? '';
  const activeSearchClient = locationState?.activeSearchClient?.trim() ?? '';
  const activeSearchTechnologies = locationState?.activeSearchTechnologies ?? [];
  const activeSearchTechnologiesParam = activeSearchTechnologies.join(',');
  const locationMatchContext = locationState?.searchMatchContext;
  const [resolvedSearchMatchContext, setResolvedSearchMatchContext] = useState<SearchMatchContext | undefined>(
    locationMatchContext,
  );
  const [adminSourceMarkdownURL, setAdminSourceMarkdownURL] = useState<string | null>(null);

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
    setSelectedGalleryIndex(0);
    setLightboxIndex(null);

    fetchProjectBySlug(slug, locale)
      .then((data) => {
        if (!cancelled) {
          setProject(data);
          setAdminSourceMarkdownURL(null);
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
    if (!project?.id || !sessionStorage.getItem('auth_token') || sessionLoading || !user?.is_admin) {
      setAdminSourceMarkdownURL(null);
      return;
    }

    let cancelled = false;
    fetchAdminProjectById(project.id)
      .then((adminProject) => {
        if (!cancelled) {
          setAdminSourceMarkdownURL(adminProject.source_markdown_url?.trim() || null);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setAdminSourceMarkdownURL(null);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [project?.id, sessionLoading, user?.is_admin]);

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
      client: activeSearchClient || undefined,
      technologies: activeSearchTechnologiesParam || undefined,
      pageSize: 100,
      lang: locale,
    })
      .then((response) => {
        if (cancelled) return;

        const remoteMatch = response.data.find((candidate) => candidate.slug === slug);
        if (remoteMatch) {
          setResolvedSearchMatchContext(
            buildSearchMatchContext(remoteMatch) ?? buildProjectSearchMatchContext(project, activeSearchQuery, t),
          );
          return;
        }

        setResolvedSearchMatchContext(buildProjectSearchMatchContext(project, activeSearchQuery, t));
      })
      .catch(() => {
        if (cancelled) return;
        setResolvedSearchMatchContext(buildProjectSearchMatchContext(project, activeSearchQuery, t));
      });

    return () => {
      cancelled = true;
    };
  }, [activeSearchCategory, activeSearchClient, activeSearchQuery, activeSearchTechnologiesParam, locale, locationMatchContext, project, slug, t]);

  const searchBackPath = activeSearchQuery
    ? buildSearchResultsPath(activeSearchQuery, {
      category: activeSearchCategory || null,
      client: activeSearchClient || null,
      technologies: activeSearchTechnologies,
    })
    : '/';
  const searchBackState = activeSearchQuery ? locationState ?? undefined : undefined;

  const projectName = project?.name ?? '';
  const technologies = project?.technologies ?? [];
  const galleryMedia = project ? getOrderedProjectMedia(project) : [];
  const galleryImages: GalleryImage[] = galleryMedia.flatMap((media) => {
    const preview = getProjectMediaMedium(media);
    const full = getProjectMediaFull(media);

    if (!preview || !full) {
      return [];
    }

    return [{
      id: media.id,
      preview,
      full,
      alt: media.alt_text?.trim() || media.caption?.trim() || projectName,
      caption: media.caption?.trim(),
      featured: media.featured,
    }];
  });
  const selectedGalleryImage = galleryImages[selectedGalleryIndex] ?? null;
  const activeLightboxImage = lightboxIndex !== null ? galleryImages[lightboxIndex] ?? null : null;
  const businessGoal = project?.profile?.business_goal?.trim();
  const problemStatement = project?.profile?.problem_statement?.trim();
  const solutionSummary = project?.profile?.solution_summary?.trim();
  const deliveryScope = project?.profile?.delivery_scope?.trim();
  const responsibilityScope = project?.profile?.responsibility_scope?.trim();
  const architecture = project?.profile?.architecture?.trim();
  const aiUsage = project?.profile?.ai_usage?.trim();
  const integrations = getRenderableList(project?.profile?.integrations, PROJECT_PROFILE_LIST_KIND.INTEGRATIONS);
  const technicalDecisions = getRenderableList(project?.profile?.technical_decisions, PROJECT_PROFILE_LIST_KIND.TECHNICAL_DECISIONS);
  const challenges = getRenderableList(project?.profile?.challenges, PROJECT_PROFILE_LIST_KIND.CHALLENGES);
  const results = getRenderableList(project?.profile?.results, PROJECT_PROFILE_LIST_KIND.RESULTS);
  const timeline = getRenderableList(project?.profile?.timeline, PROJECT_PROFILE_LIST_KIND.TIMELINE);
  const metrics = getRenderableEntries(project?.profile?.metrics);
  const searchMatchContext = resolvedSearchMatchContext;
  const searchEvidence = searchMatchContext?.evidence ?? [];
  const hasExplanation = Boolean(searchMatchContext?.explanation?.trim());
  const hasEvidence = searchEvidence.length > 0;
  const showSearchMatchContext = hasExplanation || hasEvidence;
  const lastUpdated = formatTimestamp(project?.updated_at, locale);
  const showAssistantChat = Boolean(project?.assistant_available && user?.can_use_project_assistant);
  const showAssistantAccessCard = Boolean(project?.assistant_available && user) && !sessionLoading && !showAssistantChat;
  const heroFacts = [
    hasText(project?.client_name) ? { label: t.detailClient, value: project.client_name.trim() } : null,
    lastUpdated ? { label: t.detailUpdated, value: lastUpdated } : null,
  ].filter((item): item is KeyValueEntry => Boolean(item));
  const overviewItems = [
    { label: t.detailCategory, value: project?.category ?? '' },
    { label: t.detailClient, value: project?.client_name ?? t.detailIndependent },
    { label: t.detailUpdated, value: lastUpdated ?? t.detailRecentlyCurated },
  ].filter((item) => item.value.trim().length > 0);
  const strategySections: DetailSectionData[] = [
    hasText(businessGoal) ? { title: t.detailBusinessGoal, content: businessGoal } : null,
    hasText(problemStatement) ? { title: t.detailProblem, content: problemStatement } : null,
    hasText(solutionSummary) ? { title: t.detailSolution, content: solutionSummary } : null,
  ].filter((section): section is DetailSectionData => Boolean(section));
  const executionTextSections: DetailSectionData[] = [
    hasText(deliveryScope) ? { title: t.detailDeliveryScope, content: deliveryScope } : null,
    hasText(responsibilityScope) ? { title: t.detailResponsibilityScope, content: responsibilityScope } : null,
  ].filter((section): section is DetailSectionData => Boolean(section));
  const executionListSections: DetailListSectionData[] = [
    challenges.length > 0 ? { title: t.detailChallenges, items: challenges } : null,
    results.length > 0 ? { title: t.detailResults, items: results, accent: 'highlight' } : null,
    timeline.length > 0 ? { title: t.detailTimeline, items: timeline } : null,
  ].filter((section): section is DetailListSectionData => Boolean(section));
  const technicalTextSections: DetailSectionData[] = [
    hasText(architecture) ? { title: t.detailArchitecture, content: architecture } : null,
    hasText(aiUsage) ? { title: t.detailAIUsage, content: aiUsage } : null,
  ].filter((section): section is DetailSectionData => Boolean(section));
  const technicalListSections: DetailListSectionData[] = [
    integrations.length > 0 ? { title: t.detailIntegrations, items: integrations } : null,
    technicalDecisions.length > 0 ? { title: t.detailTechnicalDecisions, items: technicalDecisions } : null,
  ].filter((section): section is DetailListSectionData => Boolean(section));
  const detailLayers: DetailLayerData[] = [
    { title: t.detailStrategyLayer, textSections: strategySections },
    { title: t.detailExecutionLayer, textSections: executionTextSections, listSections: executionListSections },
    { title: t.detailTechnicalLayer, textSections: technicalTextSections, listSections: technicalListSections, metrics, metricsTitle: t.detailMetrics },
  ];

  useEffect(() => {
    setHeaderContent?.(null);

    return () => {
      setHeaderContent?.(null);
    };
  }, [locale, setHeaderContent, slug]);

  useEffect(() => {
    if (!project || loading || error) {
      setHeaderContent?.(null);
      return;
    }

    if (!setHeaderContent) {
      return;
    }

    setHeaderContent(buildProjectHeaderContent(project, t.headerCaption));
  }, [error, loading, project, setHeaderContent, t.headerCaption]);

  useEffect(() => {
    if (!galleryViewportApi) return undefined;

    const viewportApi = galleryViewportApi;

    function syncGalleryState() {
      setSelectedGalleryIndex(viewportApi.selectedScrollSnap());
      setCanScrollGalleryPrev(viewportApi.canScrollPrev());
      setCanScrollGalleryNext(viewportApi.canScrollNext());
    }

    syncGalleryState();
    viewportApi.on('select', syncGalleryState);
    viewportApi.on('reInit', syncGalleryState);

    return () => {
      viewportApi.off('select', syncGalleryState);
      viewportApi.off('reInit', syncGalleryState);
    };
  }, [galleryViewportApi]);

  useEffect(() => {
    if (!galleryViewportApi || galleryImages.length === 0) return;

    galleryViewportApi.reInit();
    galleryViewportApi.scrollTo(0, true);
    setSelectedGalleryIndex(0);
  }, [galleryImages.length, galleryViewportApi, project?.id]);

  useEffect(() => {
    if (lightboxIndex === null) return;
    if (lightboxIndex < galleryImages.length) return;
    setLightboxIndex(galleryImages.length > 0 ? galleryImages.length - 1 : null);
  }, [galleryImages.length, lightboxIndex]);

  useEffect(() => {
    if (lightboxIndex === null) return undefined;

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setLightboxIndex(null);
      }

      if (event.key === 'ArrowRight') {
        setLightboxIndex((current) => {
          if (current === null) return current;
          return Math.min(current + 1, galleryImages.length - 1);
        });
      }

      if (event.key === 'ArrowLeft') {
        setLightboxIndex((current) => {
          if (current === null) return current;
          return Math.max(current - 1, 0);
        });
      }
    }

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [galleryImages.length, lightboxIndex]);

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
          <p className="eyebrow">{t.detailErrorEyebrow}</p>
          <p>{error ?? t.detailNotFound}</p>
          <Link to={searchBackPath} state={searchBackState}>{t.detailBack}</Link>
        </div>
      </section>
    );
  }

  return (
    <>
      <section className="detail">
        <Link className="detail__back" to={searchBackPath} state={searchBackState}>
          {t.detailBack}
        </Link>

        <article className="detail__hero card">
          <div className="detail__hero-content">
            <div className="detail__hero-intro">
              <div className="detail__hero-heading">
                {project.category ? <p className="eyebrow detail__hero-eyebrow">{project.category}</p> : null}
                {hasText(project.client_name) ? <p className="detail__context detail__context--hero">{t.detailClientContext}</p> : null}
              </div>

              {hasText(project.description) ? (
                <p className="detail__summary detail__summary--hero">{project.description.trim()}</p>
              ) : null}

              {heroFacts.length > 0 ? (
                <dl className="detail__hero-facts" aria-label={t.detailProjectHighlightsAria}>
                  {heroFacts.map((item) => (
                    <div key={item.label} className="detail__hero-fact">
                      <dt>{item.label}</dt>
                      <dd>{item.value}</dd>
                    </div>
                  ))}
                </dl>
              ) : null}
            </div>

            {adminSourceMarkdownURL ? (
              <p className="detail__admin-source">
                <a href={adminSourceMarkdownURL} target="_blank" rel="noreferrer">{t.detailAdminMarkdownSource}</a>
              </p>
            ) : null}

            {technologies.length > 0 ? (
              <div className="detail__hero-tech">
                <p className="detail__hero-tech-label">{t.detailTechnologies}</p>
                <div className="detail__chips detail__chips--hero" aria-label={t.detailTechnologiesUsedAria}>
                  {technologies.map((technology) => (
                    <span key={technology.id} className="detail__chip">
                      {technology.name}
                    </span>
                  ))}
                </div>
              </div>
            ) : null}
          </div>

          <div className="detail__hero-media">
            {galleryImages.length > 0 ? (
              <div className="detail__gallery-panel detail__gallery-panel--hero" aria-label={t.detailHeroGalleryAria}>
                {selectedGalleryImage ? (
                  <div className="detail__gallery-toolbar detail__gallery-toolbar--hero">
                    <p className="detail__gallery-counter">
                      {selectedGalleryIndex + 1} / {galleryImages.length}
                    </p>

                    {selectedGalleryImage.featured ? <span className="detail__gallery-toolbar-badge">{t.detailGalleryFeatured}</span> : null}
                  </div>
                ) : null}

                <div className="detail__gallery-carousel">
                  <div className="detail__gallery-stage-shell detail__gallery-stage-shell--hero">
                    <div className="detail__gallery-viewport" ref={galleryViewportRef}>
                      <div className="detail__gallery-track">
                        {galleryImages.map((image, index) => (
                          <div key={image.id} className="detail__gallery-slide">
                            <button
                              type="button"
                              className="detail__gallery-stage"
                              onClick={() => setLightboxIndex(index)}
                              aria-label={`${t.detailGalleryOpenImage} ${index + 1}`}
                            >
                              <img className="detail__gallery-stage-image" src={image.preview} alt={image.alt} loading="lazy" />
                            </button>
                          </div>
                        ))}
                      </div>
                    </div>

                    <div className="detail__gallery-nav" aria-label={t.detailGalleryControlsAria}>
                      <button
                        type="button"
                        className="detail__gallery-nav-button detail__gallery-nav-button--prev"
                        onClick={() => galleryViewportApi?.scrollPrev()}
                        disabled={!canScrollGalleryPrev}
                        aria-label={t.detailGalleryPreviousImage}
                      >
                        <span aria-hidden="true">←</span>
                      </button>

                      <button
                        type="button"
                        className="detail__gallery-nav-button detail__gallery-nav-button--next"
                        onClick={() => galleryViewportApi?.scrollNext()}
                        disabled={!canScrollGalleryNext}
                        aria-label={t.detailGalleryNextImage}
                      >
                        <span aria-hidden="true">→</span>
                      </button>
                    </div>

                    {selectedGalleryImage ? (
                      <div className="detail__gallery-stage-meta">
                        <div className="detail__gallery-stage-copy">
                          {selectedGalleryImage.caption ? (
                            <p className="detail__gallery-stage-caption">{selectedGalleryImage.caption}</p>
                          ) : (
                            <p className="detail__gallery-stage-caption detail__gallery-stage-caption--muted">
                              {t.detailGalleryFallbackCaption} {selectedGalleryIndex + 1}
                            </p>
                          )}
                        </div>

                        <button
                          type="button"
                          className="detail__gallery-stage-cta"
                          onClick={() => setLightboxIndex(selectedGalleryIndex)}
                        >
                          {t.detailGalleryViewFull}
                        </button>
                      </div>
                    ) : null}
                  </div>

                </div>
              </div>
            ) : (
              <div className="detail__gallery-empty detail__hero-image--placeholder">
                {t.detailVisualUnavailable}
              </div>
            )}
          </div>
        </article>

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
            <div className="detail__main-stack">
              {detailLayers.map((layer) => (
                <DetailLayer
                  key={layer.title}
                  title={layer.title}
                  textSections={layer.textSections}
                  listSections={layer.listSections}
                  metrics={layer.metrics}
                  metricsTitle={layer.metricsTitle}
                  caseStudyEyebrow={t.detailCaseStudyEyebrow}
                />
              ))}

              {showSearchMatchContext ? (
                <article className="detail__aside-panel detail__match-context card" aria-label={t.searchContextTitle}>
                  <p className="eyebrow">{t.searchContextTitle}</p>

                  {hasExplanation && searchMatchContext?.explanation ? (
                    <p className="detail__match-explanation">{searchMatchContext.explanation}</p>
                  ) : null}

                  {hasEvidence ? (
                    <div className="detail__match-evidence">
                      <p className="detail__match-evidence-title">{t.searchContextEvidenceTitle}</p>
                      <ul className="detail__match-evidence-list">
                        {searchEvidence.map((evidence, index) => (
                          <li
                            key={`${project.id}-${evidence.field}-${evidence.match_type}-${index}`}
                            className="detail__match-evidence-item"
                          >
                            <div className="detail__match-evidence-meta">
                              <span className="detail__match-evidence-field">{formatEvidenceField(evidence.field, t)}</span>
                              <span className="detail__match-evidence-type">{formatMatchType(evidence.match_type, t)}</span>
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
            </div>
          </div>
        </div>
        </section>
      {showAssistantAccessCard ? (
        <section className="card assistant-chat__panel" aria-label={t.detailAssistantAccessRequirementsAria}>
          <div className="assistant-chat__header">
	            <p className="eyebrow">{t.detailAssistantEyebrow}</p>
			{!user ? (
				<>
				  <p className="assistant-chat__copy">{t.detailAssistantLoginPrompt}</p>
				  <Link className="btn btn--primary" to="/login" state={{ from: location.pathname }}>
				    {t.detailAssistantLoginCta}
				  </Link>
				</>
			) : user.auth_provider === 'local' && !user.email_verified ? (
				<>
				  <p className="assistant-chat__copy">{t.detailAssistantVerifyPrompt}</p>
				  <Link className="btn btn--primary" to="/verify-email" state={{ from: location.pathname, email: user.email }}>
				    {t.detailAssistantVerifyCta}
				  </Link>
				</>
			) : !user.profile_completed ? (
              <>
                <p className="assistant-chat__copy">{t.detailAssistantCompleteProfilePrompt}</p>
                <Link className="btn btn--primary" to="/complete-profile" state={{ from: location.pathname }}>
                  {t.detailAssistantCompleteProfileCta}
                </Link>
              </>
            ) : !user.email_verified && user.auth_provider === 'google' ? (
              <p className="assistant-chat__copy">{t.detailAssistantGoogleRestriction}</p>
            ) : (
              <p className="assistant-chat__copy">{t.detailAssistantLocalRestriction}</p>
            )}
          </div>
        </section>
      ) : null}
      <ProjectAssistantChat slug={project.slug} enabled={showAssistantChat} lang={locale} />
      {activeLightboxImage ? (
        <div className="detail__lightbox" role="dialog" aria-modal="true" aria-label={t.detailLightboxAria} onClick={() => setLightboxIndex(null)}>
          <button type="button" className="detail__lightbox-close" aria-label={t.detailLightboxClose} onClick={() => setLightboxIndex(null)}>
            ×
          </button>

          {galleryImages.length > 1 ? (
            <>
              <button
                type="button"
                className="detail__lightbox-nav detail__lightbox-nav--prev"
                onClick={(event) => {
                  event.stopPropagation();
                  setLightboxIndex((current) => {
                    if (current === null) return current;
                    return Math.max(current - 1, 0);
                  });
                }}
                disabled={lightboxIndex === 0}
                aria-label={t.detailGalleryPreviousImage}
              >
                <span aria-hidden="true">←</span>
              </button>

              <button
                type="button"
                className="detail__lightbox-nav detail__lightbox-nav--next"
                onClick={(event) => {
                  event.stopPropagation();
                  setLightboxIndex((current) => {
                    if (current === null) return current;
                    return Math.min(current + 1, galleryImages.length - 1);
                  });
                }}
                disabled={lightboxIndex === galleryImages.length - 1}
                aria-label={t.detailGalleryNextImage}
              >
                <span aria-hidden="true">→</span>
              </button>
            </>
          ) : null}

          <div className="detail__lightbox-content" onClick={(event) => event.stopPropagation()}>
            <img className="detail__lightbox-image" src={activeLightboxImage.full} alt={activeLightboxImage.alt} />

            <div className="detail__lightbox-meta">
              <span className="detail__lightbox-counter">
                {(lightboxIndex ?? 0) + 1} / {galleryImages.length}
              </span>
              {activeLightboxImage.caption ? <p>{activeLightboxImage.caption}</p> : null}
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}
