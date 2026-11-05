CREATE TABLE researches (
    research_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    slug VARCHAR(255) UNIQUE,
    title VARCHAR(255) NOT NULL,
    purpose TEXT,
    category VARCHAR(255),
    funding_scheme VARCHAR(255),
    fund_amount VARCHAR(255),
    location VARCHAR(255),
    address TEXT,
    rpjmd TEXT,
    promised_output TEXT,
    research_significance TEXT,
    research_time_range VARCHAR(255),
    partnership_target_selection TEXT,
    partner_name VARCHAR(255),
    estimated_number_of_direct_beneficiaries VARCHAR(255),
    research_proposal_file TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    uploader_internal_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT research_public_id_unique UNIQUE (public_id),
    CONSTRAINT research_uploader_fk FOREIGN KEY (uploader_internal_id) REFERENCES users (internal_id) ON DELETE RESTRICT
);

CREATE INDEX idx_researches_uploader_internal_id ON researches (uploader_internal_id);

CREATE INDEX idx_researches_deleted_at ON researches (deleted_at);