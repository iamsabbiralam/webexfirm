-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS circular
(
    id                          VARCHAR(100) PRIMARY KEY DEFAULT uuid_generate_v4(),
    title                       VARCHAR(255) NOT NULL DEFAULT '',
    description                 TEXT         NOT NULL DEFAULT '',
    status                      SMALLINT              DEFAULT 0,
    position                    INT                   DEFAULT 0,
    start_date                  TIMESTAMP             DEFAULT current_timestamp,
    end_date                    TIMESTAMP             DEFAULT NULL,
    apply_limit                 INT                   DEFAULT 0,
    experience                  TEXT         NOT NULL DEFAULT '',
    no_vacancy                  INT                   DEFAULT 0,
    gender_allowed              varchar(150) NOT NULL DEFAULT '',
    age_limit                   VARCHAR(100) NOT NULL DEFAULT '',
    salary_limit                VARCHAR(100) NOT NULL DEFAULT '',
    job_type_id                 VARCHAR(100) NOT NULL DEFAULT '',
    circular_category_id        VARCHAR(100) NOT NULL DEFAULT '',
    is_online                   BOOLEAN               DEFAULT false,
    created_at                  TIMESTAMP             DEFAULT current_timestamp,
    created_by                  VARCHAR(100) NOT NULL DEFAULT '',
    updated_at                  TIMESTAMP             DEFAULT current_timestamp,
    updated_by                  VARCHAR(100) NOT NULL DEFAULT '',
    deleted_at                  TIMESTAMP             DEFAULT NULL,
    deleted_by                  VARCHAR(100) NOT NULL DEFAULT ''
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS circular;
