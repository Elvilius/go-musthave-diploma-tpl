-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE orders (
    number VARCHAR(50) PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    accrual FLOAT8 DEFAULT 0.0,
    uploaded_at TIMESTAMP NOT NULL,
    user_id INTEGER REFERENCES users(id)
);
CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE TABLE balances (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    current_balance FLOAT8 DEFAULT 0.0,
    withdrawn FLOAT8 DEFAULT 0.0,
    UNIQUE(user_id)
);

CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL,
    sum FLOAT8 DEFAULT 0.0,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id)
);
CREATE INDEX idx_withdrawals_user_id ON withdrawals(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_user_id;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS balances;
DROP INDEX IF EXISTS idx_withdrawals_user_id;
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd