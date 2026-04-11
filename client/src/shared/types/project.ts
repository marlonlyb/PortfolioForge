/**
 * Project types aligned to the backend Project / ProjectProfile models
 * and the search API contract.
 */

// ─── Technology ────────────────────────────────────────────────────────

export interface Technology {
  id: string;
  name: string;
  slug: string;
  category: string;
  icon?: string;
  color?: string;
}

// ─── Technology Summary (used in search results) ──────────────────────

export interface TechnologySummary {
  id: string;
  name: string;
  slug: string;
  color?: string;
}

// ─── Project Profile ──────────────────────────────────────────────────

export interface ProjectProfile {
  project_id: string;
  business_goal?: string;
  problem_statement?: string;
  solution_summary?: string;
  architecture?: string;
  integrations: unknown[];
  ai_usage?: string;
  technical_decisions: unknown[];
  challenges: unknown[];
  results: unknown[];
  metrics: Record<string, unknown>;
  timeline: unknown[];
  updated_at: number;
}

// ─── Project ──────────────────────────────────────────────────────────

export const PROJECT_STATUS = {
  DRAFT: 'draft',
  PUBLISHED: 'published',
  ARCHIVED: 'archived',
} as const;

export type ProjectStatus = (typeof PROJECT_STATUS)[keyof typeof PROJECT_STATUS];

export interface Project {
  id: string;
  name: string;
  slug: string;
  description: string;
  category: string;
  client_name?: string;
  status: ProjectStatus;
  featured: boolean;
  active: boolean;
  images: string[];
  created_at: number;
  updated_at: number;
  profile?: ProjectProfile;
  technologies?: Technology[];
}
