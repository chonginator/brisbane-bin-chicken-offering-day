-- +goose Up
CREATE TABLE seed_progress (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  last_processed_index INTEGER NOT NULL
);

-- +goose Down
DROP TABLE seed_progress;