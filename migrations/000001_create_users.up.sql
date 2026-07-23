CREATE TABLE users (
    id            BIGSERIAL    PRIMARY KEY,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Case-insensitive uniqueness; login must query WHERE LOWER(email) = LOWER($1)
CREATE UNIQUE INDEX idx_users_email_lower ON users (LOWER(email));