-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS text_data (
id SERIAL PRIMARY KEY,
user_id INT NOT NULL,
text TEXT NOT NULL,
updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE text_data
-- +goose StatementEnd
