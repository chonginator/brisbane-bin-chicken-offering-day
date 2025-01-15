-- +goose Up
CREATE TABLE streets (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  suburb_id UUID NOT NULL,
  FOREIGN KEY (suburb_id) REFERENCES suburbs(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE streets;