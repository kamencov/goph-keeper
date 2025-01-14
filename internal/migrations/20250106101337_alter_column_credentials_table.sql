-- +goose Up
-- +goose StatementBegin
ALTER TABLE credentials ADD COLUMN version INT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN version FROM credentials;
-- +goose StatementEnd
