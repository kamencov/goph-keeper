-- +goose Up
-- +goose StatementBegin
ALTER TABLE binary_data ADD COLUMN version INT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN version FROM binary_data;
-- +goose StatementEnd
