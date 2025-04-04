ALTER TABLE payments
DROP CONSTRAINT IF EXISTS payments_ticket_id_fkey;

ALTER TABLE payments
DROP COLUMN IF EXISTS ticket_id;