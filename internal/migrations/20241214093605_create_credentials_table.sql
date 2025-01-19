-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS credentials (
id SERIAL PRIMARY KEY,
user_id INT NOT NULL,
resource VARCHAR(255) NOT NULL,
login VARCHAR(255) NOT NULL,
password VARCHAR(255) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE credentials;
-- +goose StatementEnd
