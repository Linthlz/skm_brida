DROP INDEX IF EXISTS idx_researches_group;

ALTER TABLE researches DROP CONSTRAINT IF EXISTS researches_group_fk;

ALTER TABLE researches DROP COLUMN IF EXISTS group_internal_id;