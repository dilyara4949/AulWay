ALTER TABLE tickets
    ADD COLUMN payment_id VARCHAR(50) REFERENCES payments(id);
