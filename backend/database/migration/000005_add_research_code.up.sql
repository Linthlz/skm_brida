ALTER TABLE researches ADD COLUMN research_code VARCHAR(50);

CREATE UNIQUE INDEX idx_researches_research_code ON researches (research_code);