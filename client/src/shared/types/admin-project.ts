import type { ProductDetail, ProductSummary, ProductVariant } from './product';

// Canonical admin naming during the products → projects transition.
// The runtime payload still matches the legacy store contract backed by `products`.
export type AdminProjectVariant = ProductVariant;

export interface AdminProjectSummary extends ProductSummary {
  source_markdown_url?: string;
}

export interface AdminProjectDetail extends ProductDetail {
  source_markdown_url?: string;
}

export interface AdminProjectListResponse {
  items: AdminProjectDetail[];
}
