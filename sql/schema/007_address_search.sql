-- +goose Up
CREATE VIRTUAL TABLE address_search USING fts5(
  property_id UNINDEXED,
  search_text
);

INSERT INTO address_search (property_id, search_text)
SELECT
  a.property_id,
  (
    COALESCE(a.unit_number || '/', '') ||
    a.house_number ||
    COALESCE(a.house_number_suffix, '') ||
    ' ' || s.name || ', ' || sub.name
  )
FROM addresses a
JOIN streets s ON a.street_id = s.id
JOIN suburbs sub ON s.suburb_id = sub.id;

-- +goose Down
DROP TABLE address_search;