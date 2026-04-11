import { httpGet, httpPost, httpPut, httpPatch } from '../../shared/api/http';
import type { ProductDetail } from '../../shared/types/product';

// ─── Types ────────────────────────────────────────────────────────────

export interface AdminProductListResponse {
  items: ProductDetail[];
}

export interface ProjectReadiness {
  project_id: string;
  level: 'incomplete' | 'basic' | 'complete';
  missing_fields: string[];
  has_name: boolean;
  has_description: boolean;
  has_category: boolean;
  has_technologies: boolean;
  has_solution_summary: boolean;
}

export interface ReembedResponse {
  message: string;
  project_id?: string;
}

export interface CreateProductPayload {
  name: string;
  description: string;
  category: string;
  brand?: string;
  images: string[];
  active: boolean;
  variants?: CreateVariantPayload[];
}

export interface CreateVariantPayload {
  sku: string;
  color: string;
  size: string;
  price: number;
  stock: number;
  image_url?: string;
}

export interface UpdateProductPayload {
  name: string;
  description: string;
  category: string;
  brand?: string;
  images: string[];
  active: boolean;
  variants?: UpdateVariantPayload[];
}

export interface UpdateVariantPayload {
  id?: string;
  sku: string;
  color: string;
  size: string;
  price: number;
  stock: number;
  image_url?: string;
}

export interface UpdateProductStatusPayload {
  active: boolean;
}

export interface UpdateProjectEnrichmentPayload {
  profile: {
    solution_summary?: string;
    architecture?: string;
    business_goal?: string;
    problem_statement?: string;
    ai_usage?: string;
  };
  technology_ids: string[];
}

export function updateProjectEnrichment(
  id: string,
  payload: UpdateProjectEnrichmentPayload,
): Promise<void> {
  return httpPut<void>(`/api/v1/admin/projects/${id}/enrichment`, payload);
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Fetch all products (admin view — includes inactive, with variants).
 * GET /api/v1/admin/products → { data: { items: ProductDetail[] } }
 */
export function fetchAdminProducts(): Promise<AdminProductListResponse> {
  return httpGet<AdminProductListResponse>('/api/v1/admin/products');
}

/**
 * Fetch a single product by ID (admin view, with variants).
 * GET /api/v1/admin/products/:id → { data: ProductDetail }
 */
export function fetchAdminProductById(id: string): Promise<ProductDetail> {
  return httpGet<ProductDetail>(`/api/v1/admin/products/${id}`);
}

/**
 * Create a new product with variants.
 * POST /api/v1/admin/products → { data: ProductDetail }
 */
export function createProduct(payload: CreateProductPayload): Promise<ProductDetail> {
  return httpPost<ProductDetail>('/api/v1/admin/products', payload);
}

/**
 * Update an existing product and its variants.
 * PUT /api/v1/admin/products/:id → { data: ProductDetail }
 */
export function updateProduct(id: string, payload: UpdateProductPayload): Promise<ProductDetail> {
  return httpPut<ProductDetail>(`/api/v1/admin/products/${id}`, payload);
}

/**
 * Toggle product active status.
 * PATCH /api/v1/admin/products/:id/status → { data: ProductDetail }
 */
export function updateProductStatus(
  id: string,
  payload: UpdateProductStatusPayload,
): Promise<ProductDetail> {
  return httpPatch<ProductDetail>(`/api/v1/admin/products/${id}/status`, payload);
}

/**
 * Fetch search readiness for a project.
 * GET /api/v1/admin/projects/:id/readiness → { data: ProjectReadiness }
 */
export function fetchProjectReadiness(id: string): Promise<ProjectReadiness> {
  return httpGet<ProjectReadiness>(`/api/v1/admin/projects/${id}/readiness`);
}

/**
 * Re-embed a single project's search document.
 * POST /api/v1/admin/projects/:id/reembed → { data: ReembedResponse }
 */
export function reembedProject(id: string): Promise<ReembedResponse> {
  return httpPost<ReembedResponse>(`/api/v1/admin/projects/${id}/reembed`, {});
}

/**
 * Batch re-embed all stale search documents.
 * POST /api/v1/admin/projects/reembed-stale → { data: ReembedResponse }
 */
export function reembedStale(): Promise<ReembedResponse> {
  return httpPost<ReembedResponse>('/api/v1/admin/projects/reembed-stale', {});
}
