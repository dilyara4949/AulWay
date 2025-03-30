ALTER TABLE routes
DROP COLUMN IF EXISTS departure_location,
DROP COLUMN IF EXISTS destination_location;
