/**
 * Order types aligned to the backend Order / OrderItem models
 * and the API Contract (see docs/store-mvp/API-Contract-TiendaRopa.md).
 */

// ─── Order Item ────────────────────────────────────────────────────────

export interface OrderItem {
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
}

// ─── Order ─────────────────────────────────────────────────────────────

export const ORDER_STATUSES = {
  PENDING_PAYMENT: 'pending_payment',
  PAID: 'paid',
  PAYMENT_FAILED: 'payment_failed',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const;

export type OrderStatus = (typeof ORDER_STATUSES)[keyof typeof ORDER_STATUSES];

export const PAYMENT_STATUSES = {
  PENDING: 'pending',
  APPROVED: 'approved',
  CAPTURED: 'captured',
  FAILED: 'failed',
  REFUNDED: 'refunded',
} as const;

export type PaymentStatus = (typeof PAYMENT_STATUSES)[keyof typeof PAYMENT_STATUSES];

export interface Order {
  id: string;
  user_id: string;
  status: OrderStatus;
  payment_provider: string;
  payment_status: PaymentStatus;
  currency: string;
  subtotal: number;
  total: number;
  paypal_order_id?: string;
  paypal_capture_id?: string;
  paid_at?: string;
  created_at: string;
  items: OrderItem[];
}

// ─── List response ────────────────────────────────────────────────────

export interface OrderListResponse {
  items: Order[];
}
