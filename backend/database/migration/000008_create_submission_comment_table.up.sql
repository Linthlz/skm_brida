CREATE TABLE IF NOT EXISTS submission_comments (
    comment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    submission_internal_id BIGINT NOT NULL,
    author_internal_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    decision VARCHAR(50),
    round INT,
    parent_comment_id BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT submission_comments_public_unique UNIQUE (public_id),
    CONSTRAINT submission_comments_submission_fk FOREIGN KEY (submission_internal_id) REFERENCES researches (research_id) ON DELETE CASCADE,
    CONSTRAINT submission_comments_author_fk FOREIGN KEY (author_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL,
    CONSTRAINT submission_comments_parent_fk FOREIGN KEY (parent_comment_id) REFERENCES submission_comments (comment_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_submission_comments_submission ON submission_comments (submission_internal_id);

CREATE INDEX IF NOT EXISTS idx_submission_comments_author ON submission_comments (author_internal_id);

CREATE INDEX IF NOT EXISTS idx_submission_comments_decision ON submission_comments (
    submission_internal_id,
    decision
)
WHERE
    decision IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_submission_comments_created_at ON submission_comments (
    submission_internal_id,
    created_at
);

CREATE INDEX IF NOT EXISTS idx_submission_comments_deleted_at ON submission_comments (deleted_at);

-- CREATE TRIGGER update_submission_comments_updated_at
--     BEFORE UPDATE ON submission_comments
--     FOR EACH ROW
--     EXECUTE FUNCTION update_updated_at_column();