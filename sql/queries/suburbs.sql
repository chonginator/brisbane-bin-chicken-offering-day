-- name: CreateSuburb :one
INSERT INTO suburbs (id, created_at, updated_at, name)
VALUES (
  :id,
  DATETIME('now', 'utc'),
  DATETIME('now', 'utc'),
  :name
)
RETURNING *;

-- name: GetSuburbs :many
SELECT * FROM suburbs;

-- name: GetSuburbIdByName :one
SELECT * FROM suburbs
WHERE name = :name;