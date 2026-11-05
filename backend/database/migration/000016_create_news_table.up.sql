CREATE TABLE IF NOT EXISTS news (
    news_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    slug VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    category_internal_id BIGINT,
    content TEXT NOT NULL,
    cover_image TEXT NOT NULL,
    source VARCHAR(255),
    author_internal_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT news_public_unique UNIQUE (public_id),
    CONSTRAINT news_author_fk FOREIGN KEY (author_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL,
    CONSTRAINT news_category_fk FOREIGN KEY (category_internal_id) REFERENCES news_categories (category_id) ON DELETE SET NULL,
    CONSTRAINT news_status_check CHECK (
        status IN ('draft', 'published')
    )
);

CREATE INDEX IF NOT EXISTS idx_news_slug ON news (slug);

CREATE INDEX IF NOT EXISTS idx_news_category ON news (category_internal_id);

CREATE INDEX IF NOT EXISTS idx_news_status ON news (status);

CREATE INDEX IF NOT EXISTS idx_news_published_at ON news (published_at DESC);

CREATE INDEX IF NOT EXISTS idx_news_deleted_at ON news (deleted_at);

CREATE TRIGGER update_news_updated_at
    BEFORE UPDATE ON news
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();