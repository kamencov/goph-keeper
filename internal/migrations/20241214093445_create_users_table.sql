-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
login TEXT NOT NULL,
password TEXT NOT NULL,
token TEXT,
updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
