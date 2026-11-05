CREATE TABLE IF NOT EXISTS documentation_images (
    image_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    documentation_internal_id BIGINT NOT NULL,
    file_path TEXT NOT NULL,
    caption VARCHAR(500),
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT documentation_images_public_unique UNIQUE (public_id),
    CONSTRAINT documentation_images_documentation_fk FOREIGN KEY (documentation_internal_id) REFERENCES group_documentations (documentation_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_documentation_images_doc ON documentation_images (documentation_internal_id);

CREATE INDEX IF NOT EXISTS idx_documentation_images_order ON documentation_images (
    documentation_internal_id,
    display_order
);

CREATE INDEX IF NOT EXISTS idx_documentation_images_deleted_at ON documentation_images (deleted_at);