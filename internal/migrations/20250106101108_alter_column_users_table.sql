-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN version INT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN version FROM users;
-- +goose StatementEnd
