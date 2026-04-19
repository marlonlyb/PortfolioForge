/**
 * Search types aligned to the backend search API contract
 * and the evidence-based project search spec.
 */

import type { TechnologySummary } from './project';

// ─── Evidence Trace ───────────────────────────────────────────────────

export const MATCH_TYPE = {
  FTS: "fts",
  FUZZY: "fuzzy",
  SEMANTIC: "semantic",
  STRUCTURED: "structured",
} as const;

export type MatchType = (typeof MATCH_TYPE)[keyof typeof MATCH_TYPE];

export interface EvidenceField {
  field: string;
  matched_text: string;
  match_type: MatchType;
  score: number;
}

// ─── Search Result ────────────────────────────────────────────────────

export interface SearchResult {
  id: string;
  slug: string;
  title: string;
  category: string;
  client_name: string | null;
  industry_type?: string | null;
  final_product?: string | null;
  summary: string | null;
  technologies: TechnologySummary[];
  hero_image: string | null;
  score: number;
  explanation: string | null;
  evidence: EvidenceField[];
}

// ─── Search Meta ──────────────────────────────────────────────────────

export interface SearchMeta {
  total: number;
  page_size: number;
  cursor: string | null;
  query: string;
  filters_applied: {
    category: string | null;
    client: string | null;
    technologies: string[];
  };
}

// ─── Search Response ──────────────────────────────────────────────────

export interface SearchResponse {
  data: SearchResult[];
  meta: SearchMeta;
}

// ─── Search Filters ───────────────────────────────────────────────────

export interface SearchFilters {
  category: string | null;
  client: string | null;
  technologies: string[];
}
