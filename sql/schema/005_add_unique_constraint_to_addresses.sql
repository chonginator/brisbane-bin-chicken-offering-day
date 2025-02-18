-- +goose Up
CREATE UNIQUE INDEX unique_address
ON addresses (
  property_id,
  COALESCE(unit_number, ''),
  house_number,
  COALESCE(house_number_suffix, ''),
  street_id,
  collection_day,
  zone
);

-- +goose Down
DROP INDEX unique_address;