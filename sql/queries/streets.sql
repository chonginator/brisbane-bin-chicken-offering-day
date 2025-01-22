-- name: CreateStreet :one
INSERT INTO streets (id, created_at, updated_at, name, suburb_id)
VALUES (
  :id,
  DATETIME('now', 'utc'),
  DATETIME('now', 'utc'),
  :name,
  :suburb_id
)
RETURNING *;

-- name: GetStreetsBySuburbName :many
SELECT streets.*
FROM streets
INNER JOIN suburbs
ON streets.suburb_id = suburbs.id
WHERE suburbs.name = :name;