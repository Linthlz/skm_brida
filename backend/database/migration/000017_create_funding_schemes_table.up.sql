CREATE TABLE IF NOT EXISTS funding_schemes (
    funding_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    name VARCHAR(255) NOT NULL,
    organizer VARCHAR(255) NOT NULL,
    scheme_code VARCHAR(50) NOT NULL,
    ceiling BIGINT NOT NULL DEFAULT 0,
    quota INT NOT NULL DEFAULT 0,
    deadline DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    target_fields TEXT NOT NULL DEFAULT '',
    requirements TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    website VARCHAR(500) NOT NULL DEFAULT '',
    website_name VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT funding_schemes_public_unique UNIQUE (public_id),
    CONSTRAINT funding_schemes_status_check CHECK (status IN ('open', 'closing', 'soon'))
);

CREATE INDEX IF NOT EXISTS idx_funding_schemes_status ON funding_schemes (status);

CREATE INDEX IF NOT EXISTS idx_funding_schemes_deleted_at ON funding_schemes (deleted_at);

CREATE TRIGGER update_funding_schemes_updated_at
    BEFORE UPDATE ON funding_schemes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
