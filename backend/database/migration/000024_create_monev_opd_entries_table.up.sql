CREATE TABLE IF NOT EXISTS monev_opd_entries (
    entry_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    batch_internal_id BIGINT NOT NULL,
    opd_internal_id BIGINT NOT NULL,
    representative_name VARCHAR(255),
    representative_nip VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    filled_by_internal_id BIGINT,
    submitted_at TIMESTAMPTZ,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT monev_opd_entries_public_unique UNIQUE (public_id),
    CONSTRAINT monev_opd_entries_batch_fk FOREIGN KEY (batch_internal_id) REFERENCES monev_batches (batch_id) ON DELETE CASCADE,
    CONSTRAINT monev_opd_entries_opd_fk FOREIGN KEY (opd_internal_id) REFERENCES opds (opd_id) ON DELETE CASCADE,
    CONSTRAINT monev_opd_entries_filler_fk FOREIGN KEY (filled_by_internal_id) REFERENCES users (internal_id) ON DELETE SET NULL,
    CONSTRAINT monev_opd_entries_unique UNIQUE (
        batch_internal_id,
        opd_internal_id
    ),
    CONSTRAINT monev_opd_entries_status_check CHECK (
        status IN (
            'pending',
            'submitted',
            'verified'
        )
    )
);

CREATE INDEX IF NOT EXISTS idx_monev_opd_entries_batch ON monev_opd_entries (batch_internal_id);

CREATE INDEX IF NOT EXISTS idx_monev_opd_entries_opd ON monev_opd_entries (opd_internal_id);

CREATE INDEX IF NOT EXISTS idx_monev_opd_entries_status ON monev_opd_entries (status);

CREATE INDEX IF NOT EXISTS idx_monev_opd_entries_deleted_at ON monev_opd_entries (deleted_at);

CREATE TRIGGER update_monev_opd_entries_updated_at BEFORE UPDATE ON monev_opd_entries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();