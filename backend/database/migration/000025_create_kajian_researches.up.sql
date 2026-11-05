CREATE TABLE IF NOT EXISTS kajian_researches (
    kajian_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    opd_entry_internal_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT kajian_researches_public_unique UNIQUE (public_id),
    CONSTRAINT kajian_researches_entry_fk FOREIGN KEY (opd_entry_internal_id) REFERENCES monev_opd_entries (entry_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_kajian_researches_entry ON kajian_researches (opd_entry_internal_id);

CREATE INDEX IF NOT EXISTS idx_kajian_researches_order ON kajian_researches (
    opd_entry_internal_id,
    display_order
);

CREATE INDEX IF NOT EXISTS idx_kajian_researches_deleted_at ON kajian_researches (deleted_at);

CREATE TRIGGER update_kajian_researches_updated_at BEFORE UPDATE ON kajian_researches FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();