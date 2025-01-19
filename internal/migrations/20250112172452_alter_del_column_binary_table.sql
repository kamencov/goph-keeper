-- +goose Up
-- +goose StatementBegin
ALTER TABLE binary_data ADD COLUMN deleted INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN deleted FROM binary_data;
-- +goose StatementEnd