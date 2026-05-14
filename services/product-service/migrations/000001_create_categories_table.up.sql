-- 000001_create_categories_table.up.sql
--
-- Categories are created first because products hold a foreign key reference.
-- In a future multi-service world this table moves to a dedicated category-service
-- and the FK becomes a soft reference (UUID only, no DB constraint).

CREATE TABLE IF NOT EXISTS categories (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_categories_name ON categories (name);
