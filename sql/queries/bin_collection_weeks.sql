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

-- name: GetZoneForCurrentWeek :one
SELECT zone, week_start_date
FROM bin_collection_weeks
WHERE week_start_date <= DATE('now', 'utc')
ORDER BY week_start_date DESC
LIMIT 1;