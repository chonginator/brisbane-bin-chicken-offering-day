-- +goose Up
CREATE TABLE bin_collection_weeks (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  week_start_date DATE NOT NULL,
  zone TEXT NOT NULL,
  CHECK (zone IN ('Zone 1', 'Zone 2'))
);

-- +goose Down
DROP TABLE bin_collection_weeks;