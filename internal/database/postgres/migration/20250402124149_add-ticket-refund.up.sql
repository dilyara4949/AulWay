ALTER TABLE tickets
DROP CONSTRAINT tickets_payment_status_check;

ALTER TABLE tickets
    ADD CONSTRAINT tickets_payment_status_check
        CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded'));
