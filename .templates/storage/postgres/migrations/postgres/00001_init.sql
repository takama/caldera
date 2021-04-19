-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Create pgcrypto extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
