import type { ProductDetail, ProductSummary, ProductVariant } from './product';

// Canonical admin naming during the products → projects transition.
// The runtime payload still matches the legacy store contract backed by `products`.
export type AdminProjectVariant = ProductVariant;
export type AdminProjectSummary = ProductSummary;
export type AdminProjectDetail = ProductDetail;

export interface AdminProjectListResponse {
  items: AdminProjectDetail[];
}
