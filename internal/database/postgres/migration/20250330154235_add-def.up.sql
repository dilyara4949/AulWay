UPDATE routes
SET departure_location = ''
WHERE departure_location IS NULL;

UPDATE routes
SET destination_location = ''
WHERE destination_location IS NULL;