CREATE TABLE IF NOT EXISTS opds (
    opd_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    address TEXT,
    phone VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT opds_public_unique UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_opds_code ON opds (code);

CREATE INDEX IF NOT EXISTS idx_opds_active ON opds (is_active);

CREATE INDEX IF NOT EXISTS idx_opds_deleted_at ON opds (deleted_at);

CREATE TRIGGER update_opds_updated_at BEFORE UPDATE ON opds FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();