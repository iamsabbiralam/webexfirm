-- +goose Up
-- SQL in this section is executed when the migration is applied.

INSERT INTO users (id, first_name, last_name, email, password, status, created_at ,created_by , updated_at)
VALUES ('b6ddbe32-3d7e-4828-b2d7-da9927846e6b','Super', 'Admin', 'superadmin@gmail.com','$2a$10$JfrJhMaA34LPMbKl8G6Kiu0Q3EtKZnyVvehwHSn8mFX2eLXi7cVgy', 1, '2021-12-29 06:30:19.7526','abb6b1c7-c050-4dad-8d73-e246982545aa','2021-12-29 06:30:19.7526') ON CONFLICT (email) DO NOTHING;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
