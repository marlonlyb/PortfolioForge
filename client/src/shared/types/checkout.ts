/**
 * Checkout types aligned to the PayPal backend-driven flow
 * and the API Contract (see docs/store-mvp/API-Contract-TiendaRopa.md).
 */

import type { CartItem } from '../storage';

// ─── Checkout request ─────────────────────────────────────────────────

export interface CheckoutItem {
  variant_id: string;
  quantity: number;
}

export interface CheckoutRequest {
  items: CheckoutItem[];
}

// ─── Checkout response ────────────────────────────────────────────────

export interface CheckoutOrder {
  id: string;
  user_id: string;
  status: string;
  payment_provider: string;
  payment_status: string;
  currency: string;
  subtotal: number;
  total: number;
  paypal_order_id: string;
  created_at: string;
  items: Array<{
    id: string;
    product_id: string;
    product_name: string;
    variant_id: string;
    variant_sku: string;
    color: string;
    size: string;
    unit_price: number;
    quantity: number;
    line_total: number;
  }>;
}

export interface CheckoutPayPalInfo {
  order_id: string;
}

export interface CheckoutResponse {
  order: CheckoutOrder;
  paypal: CheckoutPayPalInfo;
}

// ─── Capture response ─────────────────────────────────────────────────

export interface CaptureResponse {
  order: CheckoutOrder;
}

// ─── Checkout state machine ───────────────────────────────────────────

export const CHECKOUT_STATES = {
  IDLE: 'idle',
  CREATING: 'creating',
  AWAITING_APPROVAL: 'awaiting_approval',
  CAPTURING: 'capturing',
  SUCCESS: 'success',
  ERROR: 'error',
} as const;

export type CheckoutState = (typeof CHECKOUT_STATES)[keyof typeof CHECKOUT_STATES];

// ─── Helper: cart items → checkout items ───────────────────────────────

export function toCheckoutItems(cartItems: CartItem[]): CheckoutItem[] {
  return cartItems.map((item) => ({
    variant_id: item.variant_id,
    quantity: item.quantity,
  }));
}
