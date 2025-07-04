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
SELECT DISTINCT
  addresses.property_id, 
  addresses.unit_number, 
  addresses.house_number, 
  addresses.house_number_suffix
FROM addresses
INNER JOIN streets
ON addresses.street_id = streets.id
WHERE streets.name = :name;

-- name: GetCollectionScheduleByPropertyID :one
SELECT collection_day, zone
FROM addresses
WHERE property_id = :property_id;

-- name: GetAddressBatch :many
SELECT *
FROM addresses
LIMIT :batch_size
OFFSET :offset;

-- name: SearchAddresses :many
SELECT property_id, search_text AS formatted_address
FROM address_search
WHERE search_text MATCH :query
LIMIT :limit;
