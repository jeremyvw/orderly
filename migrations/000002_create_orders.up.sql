CREATE TABLE orders (
    id              BIGSERIAL     PRIMARY KEY,
    user_id         BIGINT        NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    description     TEXT          NOT NULL,
    total_amount    NUMERIC(14,2) NOT NULL,
    status          VARCHAR(20)   NOT NULL DEFAULT 'pending',
    transaction_id  VARCHAR(255)  UNIQUE,
    idempotency_key VARCHAR(255)  UNIQUE,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_status_check
        CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    CONSTRAINT orders_total_amount_positive
        CHECK (total_amount > 0)
);

CREATE INDEX idx_orders_user_id_created_at
    ON orders (user_id, created_at DESC, id DESC);