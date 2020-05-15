-- +goose Up
-- SQL in this section is executed when the migration is applied.

{{[- if .Storage.Postgres ]}}

-- Create pgcrypto extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;
{{[- end ]}}

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
