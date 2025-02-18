-- +goose Up
CREATE TABLE addresses (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  property_id TEXT NOT NULL,
  unit_number TEXT,
  house_number TEXT NOT NULL,
  house_number_suffix TEXT,
  street_id UUID NOT NULL,
  collection_day TEXT NOT NULL,
  zone TEXT NOT NULL,
  FOREIGN KEY (street_id) REFERENCES streets(id) ON DELETE CASCADE,
  CHECK (collection_day IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday')),
  CHECK (zone IN ('Zone 1', 'Zone 2'))
);

-- +goose Down
DROP TABLE addresses;