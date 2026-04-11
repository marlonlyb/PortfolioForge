import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react';
import {
  getCartItems,
  setCartItems,
  clearCart as clearCartStorage,
  type CartItem,
} from '../../shared/storage';

// ─── Types ────────────────────────────────────────────────────────────

interface CartContextValue {
  items: CartItem[];
  itemCount: number;
  addItem: (item: CartItem) => void;
  updateItemQuantity: (variantId: string, quantity: number) => void;
  removeItem: (variantId: string) => void;
  clearCart: () => void;
}

// ─── Context ──────────────────────────────────────────────────────────

const CartContext = createContext<CartContextValue | null>(null);

// ─── Provider ─────────────────────────────────────────────────────────

export function CartProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<CartItem[]>(() => getCartItems());

  // Sync to localStorage on every change
  useEffect(() => {
    setCartItems(items);
  }, [items]);

  const itemCount = items.reduce((sum, item) => sum + item.quantity, 0);

  const addItem = useCallback((incoming: CartItem) => {
    setItems((prev) => {
      const idx = prev.findIndex((i) => i.variant_id === incoming.variant_id);
      if (idx !== -1) {
        const updated = [...prev];
        const newQuantity = Math.min(
          updated[idx]!.quantity + incoming.quantity,
          incoming.available_stock,
        );
        updated[idx] = { ...updated[idx]!, quantity: newQuantity };
        return updated;
      }
      return [...prev, incoming];
    });
  }, []);

  const updateItemQuantity = useCallback((variantId: string, quantity: number) => {
    setItems((prev) => {
      if (quantity < 1) return prev.filter((i) => i.variant_id !== variantId);
      return prev.map((i) =>
        i.variant_id === variantId
          ? { ...i, quantity: Math.min(quantity, i.available_stock) }
          : i,
      );
    });
  }, []);

  const removeItem = useCallback((variantId: string) => {
    setItems((prev) => prev.filter((i) => i.variant_id !== variantId));
  }, []);

  const clearCart = useCallback(() => {
    setItems([]);
    clearCartStorage();
  }, []);

  return (
    <CartContext.Provider value={{ items, itemCount, addItem, updateItemQuantity, removeItem, clearCart }}>
      {children}
    </CartContext.Provider>
  );
}

// ─── Hook ─────────────────────────────────────────────────────────────

export function useCart(): CartContextValue {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return ctx;
}
