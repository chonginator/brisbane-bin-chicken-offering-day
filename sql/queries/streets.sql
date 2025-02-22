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
SELECT streets.name, streets.suburb_id
FROM streets
INNER JOIN suburbs
ON streets.suburb_id = suburbs.id
WHERE suburbs.name = :name;

-- name: GetStreetsWithSuburb :many
SELECT streets.*, streets.name AS street_name, suburbs.*, suburbs.name AS suburb_name
FROM streets
INNER JOIN suburbs
ON streets.suburb_id = suburbs.id;
