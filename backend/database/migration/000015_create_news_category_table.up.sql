CREATE TABLE IF NOT EXISTS news_categories (
    category_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT news_categories_public_unique UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_news_categories_slug ON news_categories (slug);

CREATE INDEX IF NOT EXISTS idx_news_categories_active ON news_categories (is_active);

CREATE INDEX IF NOT EXISTS idx_news_categories_order ON news_categories (display_order);

CREATE INDEX IF NOT EXISTS idx_news_categories_deleted_at ON news_categories (deleted_at);

CREATE TRIGGER update_news_categories_updated_at
    BEFORE UPDATE ON news_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();