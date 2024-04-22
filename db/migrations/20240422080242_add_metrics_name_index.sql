-- +goose Up
-- +goose StatementBegin
CREATE INDEX name_idx ON metrics (name)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX name_idx
-- +goose StatementEnd
