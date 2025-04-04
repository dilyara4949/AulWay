ALTER TABLE payments
    ADD COLUMN ticket_id VARCHAR(50) UNIQUE;

ALTER TABLE payments
    ADD CONSTRAINT payments_ticket_id_fkey
        FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE;
