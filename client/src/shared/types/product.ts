/**
 * Product types aligned to the backend StoreProduct / StoreProductVariant models
 * and the API Contract (see docs/store-mvp/API-Contract-TiendaRopa.md).
 */

// ─── Variant ──────────────────────────────────────────────────────────

import type { ProjectMedia, ProjectProfile, Technology } from './project';

export interface ProductVariant {
  id: string;
  product_id: string;
  sku: string;
  color: string;
  size: string;
  price: number;
  stock: number;
  image_url?: string;
}

// ─── Product Summary (list endpoint) ──────────────────────────────────

export interface ProductSummary {
  id: string;
  name: string;
  slug: string;
  description?: string;
  category: string;
  brand?: string;
  images: string[];
  media?: ProjectMedia[];
  active: boolean;
  price_from?: number;
  available_colors?: string[];
  available_sizes?: string[];
  /** Variants may be present from the backend even on list endpoints (omitempty). */
  variants?: ProductVariant[];
  profile?: ProjectProfile;
  technologies?: Technology[];
}

// ─── Product Detail (detail endpoint) ─────────────────────────────────

export interface ProductDetail extends ProductSummary {
  description: string;
  variants: ProductVariant[];
}

// ─── List response ────────────────────────────────────────────────────

export interface ProductListResponse {
  items: ProductSummary[];
}

// ─── Sort options (for future server-side support) ────────────────────

export const SORT_OPTIONS = {
  PRICE_ASC: 'price_asc',
  PRICE_DESC: 'price_desc',
  NEWEST: 'newest',
} as const;

export type SortOption = (typeof SORT_OPTIONS)[keyof typeof SORT_OPTIONS];
