CREATE TABLE IF NOT EXISTS monev_batches (
    batch_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    year INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_by_internal_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT monev_batches_public_unique UNIQUE (public_id),
    CONSTRAINT monev_batches_year_unique UNIQUE (year),
    CONSTRAINT monev_batches_creator_fk FOREIGN KEY (created_by_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL,
    CONSTRAINT monev_batches_status_check CHECK (
        status IN ('draft', 'open', 'closed')
    )
);

CREATE INDEX IF NOT EXISTS idx_monev_batches_year ON monev_batches (year);

CREATE INDEX IF NOT EXISTS idx_monev_batches_status ON monev_batches (status);

CREATE INDEX IF NOT EXISTS idx_monev_batches_deleted_at ON monev_batches (deleted_at);

CREATE TRIGGER update_monev_batches_updated_at BEFORE UPDATE ON monev_batches FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();