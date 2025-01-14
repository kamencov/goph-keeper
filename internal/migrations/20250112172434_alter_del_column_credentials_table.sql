-- +goose Up
-- +goose StatementBegin
ALTER TABLE credentials ADD COLUMN deleted INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN deleted FROM credentials;
-- +goose StatementEnd
