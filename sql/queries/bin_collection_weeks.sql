-- name: CreateBinCollectionWeek :one
INSERT INTO bin_collection_weeks (
  id,
  created_at,
  updated_at,
  week_start_date,
  zone
)
VALUES (
  :id,
  DATETIME('now', 'utc'),
  DATETIME('now', 'utc'),
  :week_start_date,
  :zone
)
RETURNING *;