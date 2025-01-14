-- +goose Up
-- +goose StatementBegin
ALTER TABLE text_data ADD COLUMN deleted INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN deleted FROM text_data;
-- +goose StatementEnd