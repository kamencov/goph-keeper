-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards ADD COLUMN deleted INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN deleted FROM cards;
-- +goose StatementEnd