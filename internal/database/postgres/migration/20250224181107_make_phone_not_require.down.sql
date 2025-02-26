UPDATE users SET phone = '' WHERE phone IS NULL;

ALTER TABLE users
    ALTER COLUMN phone SET NOT NULL;
