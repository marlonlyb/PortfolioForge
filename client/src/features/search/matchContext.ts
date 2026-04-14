import type { Project } from '../../shared/types/project';
import { MATCH_TYPE, type EvidenceField, type MatchType, type SearchResult } from '../../shared/types/search';

const MAX_EVIDENCE_TEXT_LENGTH = 160;
const MAX_LOCAL_EVIDENCE_ITEMS = 3;

const FIELD_LABELS: Record<string, string> = {
  title: 'Título del proyecto',
  summary: 'Resumen',
  description: 'Descripción',
  client_name: 'Cliente',
  category: 'Categoría',
  technology: 'Tecnología',
  technologies: 'Tecnologías',
  solution_summary: 'Solución implementada',
  architecture: 'Arquitectura',
  business_goal: 'Objetivo de negocio',
  ai_usage: 'Uso de IA',
  technical_decisions: 'Decisiones técnicas',
  results: 'Resultados',
};

const MATCH_TYPE_LABELS: Record<MatchType, string> = {
  fts: 'Texto coincidente',
  fuzzy: 'Coincidencia aproximada',
  semantic: 'Coincidencia semántica',
  structured: 'Coincidencia estructurada',
};

export interface SearchMatchContext {
  explanation: string | null;
  evidence: EvidenceField[];
}

export interface ProjectDetailLocationState {
  searchMatchContext?: SearchMatchContext;
  activeSearchQuery?: string;
  activeSearchCategory?: string;
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

function joinEvidenceLabels(labels: string[]): string {
  if (labels.length <= 1) return labels[0] ?? 'campos relevantes del proyecto';
  if (labels.length === 2) return `${labels[0]} y ${labels[1]}`;
  return `${labels.slice(0, -1).join(', ')} y ${labels.at(-1)}`;
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

export function buildProjectSearchMatchContext(project: Project, query: string): SearchMatchContext | undefined {
  const trimmedQuery = query.trim();
  if (!trimmedQuery) return undefined;

  const evidence = buildLocalEvidence(project, trimmedQuery);
  if (evidence.length === 0) return undefined;

  const evidenceLabels = [...new Set(evidence.map((item) => formatEvidenceField(item.field)))];

  return {
    explanation: `Coincide con la búsqueda “${trimmedQuery}” en ${joinEvidenceLabels(evidenceLabels)}.`,
    evidence,
  };
}

export function formatEvidenceField(field: string): string {
  return FIELD_LABELS[field] ?? field.replaceAll('_', ' ').replace(/\b\w/g, (letter) => letter.toUpperCase());
}

export function formatMatchType(matchType: MatchType): string {
  return MATCH_TYPE_LABELS[matchType] ?? matchType;
}

export function truncateEvidenceText(text: string): string {
  const cleaned = text.replace(/\s+/g, ' ').trim();
  if (cleaned.length <= MAX_EVIDENCE_TEXT_LENGTH) return cleaned;
  return `${cleaned.slice(0, MAX_EVIDENCE_TEXT_LENGTH)}…`;
}

export function hasMatchedText(evidence: EvidenceField): boolean {
  return evidence.matched_text.trim().length > 0;
}
