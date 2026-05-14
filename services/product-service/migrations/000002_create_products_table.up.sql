-- 000002_create_products_table.up.sql

CREATE TABLE IF NOT EXISTS products (
    id          UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255)   NOT NULL,
    description TEXT           NOT NULL DEFAULT '',
    price       NUMERIC(12, 2) NOT NULL CHECK (price > 0),
    stock       INTEGER        NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category_id UUID           NOT NULL REFERENCES categories (id) ON DELETE RESTRICT,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Speed up filtered list queries (name ILIKE search).
CREATE INDEX IF NOT EXISTS idx_products_name        ON products (name);
-- Speed up joins/filters by category.
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products (category_id);
-- Speed up pagination by creation time.
CREATE INDEX IF NOT EXISTS idx_products_created_at  ON products (created_at DESC);
