-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards ADD COLUMN version INT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN version FROM cards;
-- +goose StatementEnd
