CREATE TABLE IF NOT EXISTS research_groups (
    group_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    research_internal_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE,
    join_code VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    leader_internal_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT research_groups_public_unique UNIQUE (public_id),
    CONSTRAINT research_groups_research_fk FOREIGN KEY (research_internal_id) REFERENCES researches (research_id) ON DELETE CASCADE,
    CONSTRAINT research_groups_leader_fk FOREIGN KEY (leader_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_research_groups_research ON research_groups (research_internal_id);

CREATE INDEX IF NOT EXISTS idx_research_groups_code ON research_groups (join_code);

CREATE INDEX IF NOT EXISTS idx_research_groups_leader ON research_groups (leader_internal_id);

CREATE INDEX IF NOT EXISTS idx_research_groups_deleted_at ON research_groups (deleted_at);

CREATE TRIGGER update_research_groups_updated_at
    BEFORE UPDATE ON research_groups
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();