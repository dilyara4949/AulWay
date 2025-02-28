-- Remove firebase_uid column
ALTER TABLE users DROP COLUMN IF EXISTS firebase_uid;

-- Ensure no NULL values before setting NOT NULL
UPDATE users SET phone = 'unknown' WHERE phone IS NULL;
UPDATE users SET email = 'unknown@example.com' WHERE email IS NULL;
UPDATE users SET first_name = 'Unknown' WHERE first_name IS NULL;
UPDATE users SET last_name = 'Unknown' WHERE last_name IS NULL;

-- Make all fields NOT NULL
ALTER TABLE users
    ALTER COLUMN email SET NOT NULL,
ALTER COLUMN phone SET NOT NULL,
    ALTER COLUMN first_name SET NOT NULL,
    ALTER COLUMN last_name SET NOT NULL;
