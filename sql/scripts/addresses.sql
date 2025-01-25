-- Delete duplicate address records
DELETE FROM addresses
WHERE id NOT IN (
  SELECT MIN(id)
  FROM addresses
  GROUP BY property_id,
          unit_number,
          house_number,
          house_number_suffix,
          street_id,
          collection_day,
          zone
);