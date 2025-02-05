-- name: CreateAddress :one
INSERT INTO addresses (
  id,
  created_at,
  updated_at,
  property_id,
  unit_number,
  house_number,
  house_number_suffix,
  street_id,
  collection_day,
  zone
)
VALUES (
  :id,
  DATETIME('now', 'utc'),
  DATETIME('now', 'utc'),
  :property_id,
  :unit_number,
  :house_number,
  :house_number_suffix,
  :street_id,
  :collection_day,
  :zone
)
RETURNING *;

-- name: GetAddressesByStreetName :many
SELECT addresses.*
FROM addresses
INNER JOIN streets
ON addresses.street_id = streets.id
WHERE streets.name = :name;

-- name: GetAddressBatch :many
SELECT *
FROM addresses
LIMIT :batch_size
OFFSET :offset;
