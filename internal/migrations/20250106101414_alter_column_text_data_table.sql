-- +goose Up
-- +goose StatementBegin
ALTER TABLE text_data ADD COLUMN version INT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN version FROM text_data;
-- +goose StatementEnd
