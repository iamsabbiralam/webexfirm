-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id                      VARCHAR(100) PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name              VARCHAR(100) NOT NULL DEFAULT '',
    last_name               VARCHAR(100) NOT NULL DEFAULT '',
    email                   VARCHAR(100) UNIQUE NOT NULL DEFAULT '',
    password                VARCHAR(100) NOT NULL DEFAULT '',
    status                  SMALLINT DEFAULT 0,
    created_at              TIMESTAMP   DEFAULT current_timestamp,
    created_by              VARCHAR(100) NOT NULL DEFAULT '',
    updated_at              TIMESTAMP DEFAULT current_timestamp,
    updated_by              VARCHAR(100) NOT NULL DEFAULT '',
    deleted_at              TIMESTAMP  DEFAULT NULL,
    deleted_by              VARCHAR(100) NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users ;
-- +goose StatementEnd
