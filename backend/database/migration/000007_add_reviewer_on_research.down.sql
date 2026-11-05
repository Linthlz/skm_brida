DROP INDEX IF EXISTS idx_researches_status;

DROP INDEX IF EXISTS idx_researches_reviewer;

ALTER TABLE researches
DROP COLUMN IF EXISTS review_finished_at,
DROP COLUMN IF EXISTS review_started_at,
DROP COLUMN IF EXISTS reviewer_internal_id;