import { httpGet } from '../../shared/api/http';
import type { ProductListResponse, ProductDetail } from '../../shared/types/product';

/**
 * Fetch all active products (public endpoint).
 * The backend currently returns `{ data: { items: StoreProduct[] } }`
 * without server-side pagination or filters.
 */
export function fetchProducts(): Promise<ProductListResponse> {
  return httpGet<ProductListResponse>('/api/v1/public/products');
}

/**
 * Fetch a single product by ID with full detail and variants.
 * Returns `{ data: StoreProduct }` which includes the `variants` array.
 */
export function fetchProductById(id: string): Promise<ProductDetail> {
  return httpGet<ProductDetail>(`/api/v1/public/products/${id}`);
}
