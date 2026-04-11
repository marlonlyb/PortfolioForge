/**
 * Storage wrappers for JWT token (sessionStorage) and cart (localStorage).
 *
 * These centralise access so the rest of the app never touches
 * the raw Storage API directly.
 */

// ─── Auth token (sessionStorage) ─────────────────────────────────────

const AUTH_TOKEN_KEY = 'auth_token';

export function getToken(): string | null {
  return sessionStorage.getItem(AUTH_TOKEN_KEY);
}

export function setToken(token: string): void {
  sessionStorage.setItem(AUTH_TOKEN_KEY, token);
}

export function removeToken(): void {
  sessionStorage.removeItem(AUTH_TOKEN_KEY);
}

// ─── Cart (localStorage) ─────────────────────────────────────────────

export interface CartItem {
  product_id: string;
  product_name: string;
  product_image: string;
  variant_id: string;
  variant_sku: string;
  color: string;
  size: string;
  unit_price: number;
  quantity: number;
  available_stock: number;
}

const CART_KEY = 'store_cart';

export function getCartItems(): CartItem[] {
  try {
    const raw = localStorage.getItem(CART_KEY);
    if (!raw) return [];
    const parsed: unknown = JSON.parse(raw);
    return Array.isArray(parsed) ? (parsed as CartItem[]) : [];
  } catch {
    return [];
  }
}

export function setCartItems(items: CartItem[]): void {
  localStorage.setItem(CART_KEY, JSON.stringify(items));
}

export function clearCart(): void {
  localStorage.removeItem(CART_KEY);
}
