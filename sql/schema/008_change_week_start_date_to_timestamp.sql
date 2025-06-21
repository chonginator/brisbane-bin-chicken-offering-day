-- +goose Up
ALTER TABLE bin_collection_weeks RENAME TO old_bin_collection_weeks;

CREATE TABLE bin_collection_weeks (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  week_start_date TIMESTAMP NOT NULL,
  zone TEXT NOT NULL,
  CHECK (zone IN ('Zone 1', 'Zone 2'))
);

INSERT INTO bin_collection_weeks (id, created_at, updated_at, week_start_date, zone)
SELECT id, created_at, updated_at, strftime('%Y-%m-%d', week_start_date) || ' 00:00:00', zone
FROM old_bin_collection_weeks;

DROP TABLE old_bin_collection_weeks;

-- +goose Down
ALTER TABLE bin_collection_weeks RENAME TO old_bin_collection_weeks;

CREATE TABLE bin_collection_weeks (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  week_start_date DATE NOT NULL,
  zone TEXT NOT NULL,
  CHECK (zone IN ('Zone 1', 'Zone 2'))
);

INSERT INTO bin_collection_weeks (id, created_at, updated_at, week_start_date, zone)
SELECT id, created_at, updated_at, CAST(week_start_date AS DATE), zone
FROM old_bin_collection_weeks;

DROP TABLE old_bin_collection_weeks;