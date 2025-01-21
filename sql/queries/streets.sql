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
SELECT * FROM streets
WHERE suburb_id = (
  SELECT id FROM suburbs
  WHERE suburbs.name = :name
);