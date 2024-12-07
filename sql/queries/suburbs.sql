-- name: CreateSuburb :one
INSERT INTO suburbs (id, created_at, updated_at, name)
VALUES (
  :id,
  DATETIME('now', 'utc'),
  DATETIME('now', 'utc'),
  :name
)
RETURNING *;