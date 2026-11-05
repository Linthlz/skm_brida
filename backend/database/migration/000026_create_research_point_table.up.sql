CREATE TABLE IF NOT EXISTS research_points (
    point_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    kajian_internal_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    monitoring_file_path TEXT,
    has_monitoring_result BOOLEAN NOT NULL DEFAULT FALSE,
    monitoring_result_description TEXT,
    kendala_pelaksanaan_rekomendasi TEXT,
    benefits TEXT,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT research_points_public_unique UNIQUE (public_id),
    CONSTRAINT research_points_kajian_fk FOREIGN KEY (kajian_internal_id) REFERENCES kajian_researches (kajian_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_research_points_kajian ON research_points (kajian_internal_id);

CREATE INDEX IF NOT EXISTS idx_research_points_order ON research_points (
    kajian_internal_id,
    display_order
);

CREATE INDEX IF NOT EXISTS idx_research_points_deleted_at ON research_points (deleted_at);

CREATE TRIGGER update_research_points_updated_at BEFORE UPDATE ON research_points FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();