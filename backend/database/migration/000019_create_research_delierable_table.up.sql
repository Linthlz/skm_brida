CREATE TABLE IF NOT EXISTS research_deliverables (
    deliverable_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    research_internal_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    proof_type VARCHAR(20) NOT NULL DEFAULT 'checkbox',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    file_path TEXT,
    completed_by_internal_id BIGINT,
    completed_at TIMESTAMPTZ,
    display_order INT NOT NULL DEFAULT 0,
    created_by_internal_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT research_deliverables_public_unique UNIQUE (public_id),
    CONSTRAINT research_deliverables_research_fk FOREIGN KEY (research_internal_id) REFERENCES researches (research_id) ON DELETE CASCADE,
    CONSTRAINT research_deliverables_completer_fk FOREIGN KEY (completed_by_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL,
    CONSTRAINT research_deliverables_creator_fk FOREIGN KEY (created_by_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL,
    CONSTRAINT research_deliverables_proof_type_check CHECK (
        proof_type IN ('checkbox', 'file')
    ),
    CONSTRAINT research_deliverables_status_check CHECK (status IN ('pending', 'done'))
);

CREATE INDEX IF NOT EXISTS idx_research_deliverables_research ON research_deliverables (research_internal_id);

CREATE INDEX IF NOT EXISTS idx_research_deliverables_status ON research_deliverables (status);

CREATE INDEX IF NOT EXISTS idx_research_deliverables_order ON research_deliverables (
    research_internal_id,
    display_order
);

CREATE INDEX IF NOT EXISTS idx_research_deliverables_deleted_at ON research_deliverables (deleted_at);

CREATE TRIGGER update_research_deliverables_updated_at
    BEFORE UPDATE ON research_deliverables
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();