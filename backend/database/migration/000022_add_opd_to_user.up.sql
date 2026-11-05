ALTER TABLE users ADD COLUMN IF NOT EXISTS opd_internal_id BIGINT;

ALTER TABLE users
ADD CONSTRAINT users_opd_fk FOREIGN KEY (opd_internal_id) REFERENCES opds (opd_id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_users_opd ON users (opd_internal_id);