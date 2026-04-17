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

export type ProjectProfilePrimitive = string | number | boolean;
export type ProjectProfileStructuredItem =
  | ProjectProfilePrimitive
  | Record<string, ProjectProfilePrimitive>;
export type ProjectProfileStructuredList = ProjectProfileStructuredItem[];
export type ProjectProfileMetrics = Record<string, ProjectProfilePrimitive>;

export interface ProjectProfile {
  project_id: string;
  business_goal?: string;
  problem_statement?: string;
  solution_summary?: string;
  delivery_scope?: string;
  responsibility_scope?: string;
  architecture?: string;
  integrations: ProjectProfileStructuredList;
  ai_usage?: string;
  technical_decisions: ProjectProfileStructuredList;
  challenges: ProjectProfileStructuredList;
  results: ProjectProfileStructuredList;
  metrics: ProjectProfileMetrics;
  timeline: ProjectProfileStructuredList;
  updated_at: number;
}

// ─── Project ──────────────────────────────────────────────────────────

export const PROJECT_STATUS = {
  DRAFT: 'draft',
  PUBLISHED: 'published',
  ARCHIVED: 'archived',
} as const;

export type ProjectStatus = (typeof PROJECT_STATUS)[keyof typeof PROJECT_STATUS];

export interface ProjectMedia {
  id: string;
  project_id: string;
  media_type: string;
  fallback_url?: string;
  low_url?: string;
  medium_url?: string;
  high_url?: string;
  caption?: string;
  alt_text?: string;
  sort_order: number;
  featured: boolean;
}

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
  assistant_available: boolean;
  images: string[];
  media?: ProjectMedia[];
  created_at: number;
  updated_at: number;
  profile?: ProjectProfile;
  technologies?: Technology[];
}

export interface ProjectListResponse {
  items: Project[];
}

export interface ProjectAssistantMessage {
  role: 'user' | 'assistant';
  content: string;
}

export interface ProjectAssistantRequest {
  question: string;
  history: ProjectAssistantMessage[];
  lang: string;
}

export interface ProjectAssistantResponse {
  answer: string;
}
