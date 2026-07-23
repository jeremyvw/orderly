CREATE TABLE order_logs (
    id         BIGSERIAL    PRIMARY KEY,
    order_id   BIGINT       NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    event_id   VARCHAR(255) NOT NULL UNIQUE,
    event_type VARCHAR(50)  NOT NULL,
    payload    JSONB        NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_logs_order_id ON order_logs (order_id);