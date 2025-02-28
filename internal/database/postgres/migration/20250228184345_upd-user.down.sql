-- Add back the firebase_uid column (if needed)
ALTER TABLE users ADD COLUMN firebase_uid VARCHAR(128) UNIQUE;

-- Revert fields to nullable
ALTER TABLE users
    ALTER COLUMN email DROP NOT NULL,
ALTER COLUMN phone DROP NOT NULL,
    ALTER COLUMN first_name DROP NOT NULL,
    ALTER COLUMN last_name DROP NOT NULL;
