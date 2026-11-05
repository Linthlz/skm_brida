ALTER TABLE researches
ADD COLUMN IF NOT EXISTS group_internal_id BIGINT;

ALTER TABLE researches
ADD CONSTRAINT researches_group_fk FOREIGN KEY (group_internal_id) REFERENCES research_groups (group_id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_researches_group ON researches (group_internal_id);