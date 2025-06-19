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

-- name: GetCollectionSchedulesByPropertyID :many
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

-- WITH address_data AS (
--   SELECT 
--     a.property_id,
--     a.collection_day,
--     a.zone,
--     CAST(
--       CASE 
--         WHEN a.unit_number IS NOT NULL AND a.unit_number != '' 
--         THEN a.unit_number || '/' || a.house_number 
--         ELSE a.house_number 
--       END || 
--       CASE 
--         WHEN a.house_number_suffix IS NOT NULL AND a.house_number_suffix != '' 
--         THEN a.house_number_suffix 
--         ELSE '' 
--       END || ' ' || s.name || ', ' || sub.name
--     AS VARCHAR) AS formatted_address
--   FROM addresses a
--   JOIN streets s ON a.street_id = s.id
--   JOIN suburbs sub ON s.suburb_id = sub.id
--   WHERE 
--     (
--       COALESCE(a.unit_number, '') || 
--       '/' ||
--       COALESCE(a.house_number, '') || 
--       COALESCE(a.house_number_suffix, '') || 
--       ' ' ||
--       s.name || 
--       ', ' ||
--       sub.name
--     ) LIKE LOWER('%' || :query || '%')
-- )
-- SELECT 
--   property_id,
--   collection_day,
--   zone,
--   formatted_address
-- FROM address_data
-- LIMIT :limit;
