ALTER TABLE users
    ADD COLUMN IF NOT EXISTS password_hash VARCHAR(72);

UPDATE users
SET password_hash = password
WHERE password_hash IS NULL
  AND password IS NOT NULL;

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS name VARCHAR(160),
    ADD COLUMN IF NOT EXISTS slug VARCHAR(180),
    ADD COLUMN IF NOT EXISTS category VARCHAR(80),
    ADD COLUMN IF NOT EXISTS brand VARCHAR(80),
    ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT TRUE;

UPDATE products
SET name = COALESCE(name, product_name),
    slug = COALESCE(slug, regexp_replace(lower(product_name), '[^a-z0-9]+', '-', 'g')),
    category = COALESCE(category, 'general')
WHERE product_name IS NOT NULL;

CREATE INDEX IF NOT EXISTS ix_products_slug ON products (slug);
CREATE INDEX IF NOT EXISTS ix_products_category ON products (category);
CREATE INDEX IF NOT EXISTS ix_products_active ON products (active);

CREATE TABLE IF NOT EXISTS product_variants (
    id UUID NOT NULL,
    product_id UUID NOT NULL,
    sku VARCHAR(120) NOT NULL,
    color VARCHAR(60) NOT NULL,
    size VARCHAR(30) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT product_variants_id_pk PRIMARY KEY (id),
    CONSTRAINT product_variants_product_id_fk FOREIGN KEY (product_id)
        REFERENCES products (id) ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT product_variants_sku_uk UNIQUE (sku),
    CONSTRAINT product_variants_product_color_size_uk UNIQUE (product_id, color, size),
    CONSTRAINT product_variants_price_ck CHECK (price >= 0),
    CONSTRAINT product_variants_stock_ck CHECK (stock >= 0)
);

CREATE INDEX IF NOT EXISTS ix_product_variants_product_id ON product_variants (product_id);

CREATE TABLE IF NOT EXISTS orders (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending_payment',
    payment_provider VARCHAR(24) NOT NULL DEFAULT 'paypal',
    payment_status VARCHAR(24) NOT NULL DEFAULT 'pending',
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    subtotal NUMERIC(10,2) NOT NULL,
    total NUMERIC(10,2) NOT NULL,
    paypal_order_id VARCHAR(64),
    paypal_capture_id VARCHAR(64),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT orders_id_pk PRIMARY KEY (id),
    CONSTRAINT orders_user_id_fk FOREIGN KEY (user_id)
        REFERENCES users (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT orders_status_ck CHECK (status IN ('pending_payment', 'paid', 'payment_failed', 'cancelled', 'refunded')),
    CONSTRAINT orders_payment_provider_ck CHECK (payment_provider IN ('paypal')),
    CONSTRAINT orders_payment_status_ck CHECK (payment_status IN ('pending', 'approved', 'captured', 'failed', 'refunded')),
    CONSTRAINT orders_subtotal_ck CHECK (subtotal >= 0),
    CONSTRAINT orders_total_ck CHECK (total >= 0),
    CONSTRAINT orders_paypal_order_id_uk UNIQUE (paypal_order_id),
    CONSTRAINT orders_paypal_capture_id_uk UNIQUE (paypal_capture_id)
);

CREATE INDEX IF NOT EXISTS ix_orders_user_id_created_at ON orders (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS ix_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS ix_orders_payment_status ON orders (payment_status);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID NOT NULL,
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    product_name VARCHAR(160) NOT NULL,
    variant_sku VARCHAR(120) NOT NULL,
    color VARCHAR(60) NOT NULL,
    size VARCHAR(30) NOT NULL,
    unit_price NUMERIC(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    line_total NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT order_items_id_pk PRIMARY KEY (id),
    CONSTRAINT order_items_order_id_fk FOREIGN KEY (order_id)
        REFERENCES orders (id) ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT order_items_product_id_fk FOREIGN KEY (product_id)
        REFERENCES products (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT order_items_variant_id_fk FOREIGN KEY (variant_id)
        REFERENCES product_variants (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT order_items_quantity_ck CHECK (quantity > 0),
    CONSTRAINT order_items_unit_price_ck CHECK (unit_price >= 0),
    CONSTRAINT order_items_line_total_ck CHECK (line_total >= 0)
);

CREATE INDEX IF NOT EXISTS ix_order_items_order_id ON order_items (order_id);
CREATE INDEX IF NOT EXISTS ix_order_items_variant_id ON order_items (variant_id);
