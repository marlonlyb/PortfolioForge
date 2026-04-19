import type { Project } from '../../shared/types/project';
import type { getMessages } from '../../shared/i18n/messages';
import {
  MATCH_TYPE,
  type EvidenceField,
  type MatchType,
  type SearchFilters,
  type SearchResult,
} from '../../shared/types/search';
import { normalizeEditorialMetadataText } from '../../shared/lib/projectMetadata';

const MAX_EVIDENCE_TEXT_LENGTH = 160;
const MAX_LOCAL_EVIDENCE_ITEMS = 3;

type SearchMatchMessages = ReturnType<typeof getMessages>;

export interface SearchMatchContext {
  explanation: string | null;
  evidence: EvidenceField[];
}

export interface ProjectDetailLocationState {
  searchMatchContext?: SearchMatchContext;
  activeSearchQuery?: string;
  activeSearchCategory?: string;
  activeSearchClient?: string;
  activeSearchTechnologies?: string[];
  searchResultsSnapshot?: SearchResultsSnapshot;
}

export interface SearchResultsSnapshot {
  results: SearchResult[];
  total: number;
  cursor: string | null;
}

function normalizeActiveSearchValue(value?: string | null): string | undefined {
  const trimmed = value?.trim() ?? '';
  return trimmed || undefined;
}

export function normalizeActiveSearchTechnologies(technologies?: string[]): string[] | undefined {
  const normalized = technologies
    ?.map((technology) => technology.trim())
    .filter((technology) => technology.length > 0);

  return normalized && normalized.length > 0 ? normalized : undefined;
}

export function buildSearchResultsLocationState(
  query: string,
  filters: SearchFilters,
  snapshot?: SearchResultsSnapshot,
): ProjectDetailLocationState | undefined {
  const activeSearchQuery = normalizeActiveSearchValue(query);
  const activeSearchCategory = normalizeActiveSearchValue(filters.category);
  const activeSearchClient = normalizeActiveSearchValue(filters.client);
  const activeSearchTechnologies = normalizeActiveSearchTechnologies(filters.technologies);

  if (!activeSearchQuery && !snapshot) {
    return undefined;
  }

  return {
    activeSearchQuery,
    activeSearchCategory,
    activeSearchClient,
    activeSearchTechnologies,
    searchResultsSnapshot: snapshot,
  };
}

export function buildSearchResultsPath(query: string, filters: SearchFilters): string {
  const searchParams = new URLSearchParams();
  const trimmedQuery = query.trim();

  if (trimmedQuery) searchParams.set('q', trimmedQuery);
  if (filters.category) searchParams.set('category', filters.category);
  if (filters.client) searchParams.set('client', filters.client);
  if (filters.technologies.length > 0) {
    searchParams.set('technologies', filters.technologies.join(','));
  }

  const serializedSearch = searchParams.toString();
  return serializedSearch ? `/search?${serializedSearch}` : '/search';
}

export function matchesActiveSearchState(
  locationState: ProjectDetailLocationState | null,
  query: string,
  filters: SearchFilters,
): boolean {
  if (!locationState?.searchResultsSnapshot) {
    return false;
  }

  return locationState.activeSearchQuery === normalizeActiveSearchValue(query)
    && locationState.activeSearchCategory === normalizeActiveSearchValue(filters.category)
    && locationState.activeSearchClient === normalizeActiveSearchValue(filters.client)
    && JSON.stringify(locationState.activeSearchTechnologies ?? []) === JSON.stringify(filters.technologies);
}

interface ProjectSearchField {
  field: string;
  value: string;
}

function normalizeSearchValue(value?: string): string {
  return (value ?? '')
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim();
}

function splitSearchTokens(query: string): string[] {
  return normalizeSearchValue(query)
    .split(/[^a-z0-9+#.]+/)
    .filter(Boolean);
}

function getProjectSearchFields(project: Project): ProjectSearchField[] {
  return [
    { field: 'title', value: project.name },
    { field: 'description', value: project.description },
    { field: 'category', value: project.category },
    { field: 'client_name', value: project.client_name ?? '' },
    { field: 'industry_type', value: normalizeEditorialMetadataText(project.industry_type) },
    { field: 'final_product', value: project.final_product ?? '' },
    { field: 'business_goal', value: project.profile?.business_goal ?? '' },
    { field: 'solution_summary', value: project.profile?.solution_summary ?? '' },
    { field: 'architecture', value: project.profile?.architecture ?? '' },
    { field: 'ai_usage', value: project.profile?.ai_usage ?? '' },
    ...((project.technologies ?? []).map((technology) => ({ field: 'technology', value: technology.name })) ?? []),
  ].filter((entry) => entry.value.trim().length > 0);
}

function buildLocalEvidence(project: Project, query: string): EvidenceField[] {
  const normalizedQuery = normalizeSearchValue(query);
  const tokens = splitSearchTokens(query);

  if (!normalizedQuery && tokens.length === 0) {
    return [];
  }

  const evidence: EvidenceField[] = [];
  const seen = new Set<string>();

  getProjectSearchFields(project).forEach((entry) => {
    const normalizedValue = normalizeSearchValue(entry.value);
    const directMatch = normalizedQuery.length > 0 && normalizedValue.includes(normalizedQuery);
    const tokenMatch = tokens.some((token) => normalizedValue.includes(token));

    if (!directMatch && !tokenMatch) {
      return;
    }

    const key = `${entry.field}:${entry.value}`;
    if (seen.has(key) || evidence.length >= MAX_LOCAL_EVIDENCE_ITEMS) {
      return;
    }

    seen.add(key);
    evidence.push({
      field: entry.field,
      matched_text: entry.value.trim(),
      match_type: MATCH_TYPE.STRUCTURED,
      score: 1,
    });
  });

  return evidence;
}

function joinEvidenceLabels(labels: string[], t: SearchMatchMessages): string {
  if (labels.length <= 1) return labels[0] ?? t.searchContextRelevantProjectFields;
  if (labels.length === 2) return `${labels[0]} ${t.searchContextAnd} ${labels[1]}`;
  return `${labels.slice(0, -1).join(', ')} ${t.searchContextAnd} ${labels.at(-1)}`;
}

export function buildSearchMatchContext(result: SearchResult): SearchMatchContext | undefined {
  const explanation = result.explanation?.trim() ?? '';
  const evidence = result.evidence ?? [];

  if (!explanation && evidence.length === 0) {
    return undefined;
  }

  return {
    explanation: explanation || null,
    evidence,
  };
}

export function buildProjectSearchMatchContext(
  project: Project,
  query: string,
  t: SearchMatchMessages,
): SearchMatchContext | undefined {
  const trimmedQuery = query.trim();
  if (!trimmedQuery) return undefined;

  const evidence = buildLocalEvidence(project, trimmedQuery);
  if (evidence.length === 0) return undefined;

  const evidenceLabels = [...new Set(evidence.map((item) => formatEvidenceField(item.field, t)))];

  return {
    explanation: `${t.searchContextExplanationPrefix} “${trimmedQuery}” ${t.searchContextExplanationConnector} ${joinEvidenceLabels(evidenceLabels, t)}.`,
    evidence,
  };
}

export function formatEvidenceField(field: string, t: SearchMatchMessages): string {
  const fieldLabels: Record<string, string> = {
    title: t.searchEvidenceFieldTitle,
    summary: t.searchEvidenceFieldSummary,
    description: t.searchEvidenceFieldDescription,
    client_name: t.searchEvidenceFieldClient,
    category: t.searchEvidenceFieldCategory,
    technology: t.searchEvidenceFieldTechnology,
    technologies: t.searchEvidenceFieldTechnologies,
    solution_summary: t.searchEvidenceFieldSolution,
    architecture: t.searchEvidenceFieldArchitecture,
    business_goal: t.searchEvidenceFieldBusinessGoal,
    ai_usage: t.searchEvidenceFieldAIUsage,
    industry_type: t.searchEvidenceFieldIndustry,
    final_product: t.searchEvidenceFieldFinalProduct,
    technical_decisions: t.searchEvidenceFieldTechnicalDecisions,
    results: t.searchEvidenceFieldResults,
  };

  return fieldLabels[field] ?? field.replaceAll('_', ' ').replace(/\b\w/g, (letter) => letter.toUpperCase());
}

export function formatMatchType(matchType: MatchType, t: SearchMatchMessages): string {
  const matchTypeLabels: Record<MatchType, string> = {
    fts: t.searchMatchTypeFTS,
    fuzzy: t.searchMatchTypeFuzzy,
    semantic: t.searchMatchTypeSemantic,
    structured: t.searchMatchTypeStructured,
  };

  return matchTypeLabels[matchType] ?? matchType;
}

export function truncateEvidenceText(text: string): string {
  const cleaned = text.replace(/\s+/g, ' ').trim();
  if (cleaned.length <= MAX_EVIDENCE_TEXT_LENGTH) return cleaned;
  return `${cleaned.slice(0, MAX_EVIDENCE_TEXT_LENGTH)}…`;
}

export function hasMatchedText(evidence: EvidenceField): boolean {
  return evidence.matched_text.trim().length > 0;
}
