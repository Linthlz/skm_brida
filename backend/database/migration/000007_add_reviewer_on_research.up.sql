ALTER TABLE researches
ADD COLUMN reviewer_internal_id BIGINT REFERENCES users (internal_id) ON DELETE SET NULL,
ADD COLUMN review_started_at TIMESTAMPTZ,
ADD COLUMN review_finished_at TIMESTAMPTZ;

CREATE INDEX idx_researches_reviewer ON researches (reviewer_internal_id);

CREATE INDEX idx_researches_status ON researches (status);