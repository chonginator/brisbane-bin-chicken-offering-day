-- +goose Up
CREATE TABLE suburbs (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE suburbs;