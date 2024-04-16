-- +goose Up
-- +goose StatementBegin
CREATE TYPE metric_type AS ENUM ('counter', 'gauge');

CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    mtype metric_type NOT NULL,
    delta BIGINT NULL,
    value DOUBLE PRECISION NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics;
DROP TYPE IF EXISTS metric_type;
-- +goose StatementEnd
