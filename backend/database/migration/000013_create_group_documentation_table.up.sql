CREATE TABLE IF NOT EXISTS group_documentations (
    documentation_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    group_internal_id BIGINT NOT NULL,
    author_internal_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    narrative TEXT,
    location VARCHAR(255),
    video_url VARCHAR(500),
    event_date DATE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT group_documentations_public_unique UNIQUE (public_id),
    CONSTRAINT group_documentations_group_fk FOREIGN KEY (group_internal_id) REFERENCES research_groups (group_id) ON DELETE CASCADE,
    CONSTRAINT group_documentations_author_fk FOREIGN KEY (author_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_group_documentations_group ON group_documentations (group_internal_id);

CREATE INDEX IF NOT EXISTS idx_group_documentations_author ON group_documentations (author_internal_id);

CREATE INDEX IF NOT EXISTS idx_group_documentations_event_date ON group_documentations (event_date);

CREATE INDEX IF NOT EXISTS idx_group_documentations_deleted_at ON group_documentations (deleted_at);

CREATE TRIGGER update_group_documentations_updated_at
    BEFORE UPDATE ON group_documentations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();