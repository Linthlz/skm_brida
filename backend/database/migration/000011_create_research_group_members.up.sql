CREATE TABLE IF NOT EXISTS research_group_members (
    member_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    group_internal_id BIGINT NOT NULL,
    user_internal_id BIGINT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT research_group_members_public_unique UNIQUE (public_id),
    CONSTRAINT research_group_members_group_fk FOREIGN KEY (group_internal_id) REFERENCES research_groups (group_id) ON DELETE CASCADE,
    CONSTRAINT research_group_members_user_fk FOREIGN KEY (user_internal_id) REFERENCES users (internal_id) ON DELETE CASCADE,
    CONSTRAINT research_group_members_unique UNIQUE (
        group_internal_id,
        user_internal_id
    )
);

CREATE INDEX IF NOT EXISTS idx_research_group_members_group ON research_group_members (group_internal_id);

CREATE INDEX IF NOT EXISTS idx_research_group_members_user ON research_group_members (user_internal_id);

CREATE INDEX IF NOT EXISTS idx_research_group_members_deleted_at ON research_group_members (deleted_at);