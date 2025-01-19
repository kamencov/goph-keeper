-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS binary_data (
id SERIAL PRIMARY KEY,
user_id INT NOT NULL,
binary_data BYTEA NOT NULL,
updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS binary_data;
-- +goose StatementEnd
