-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Events table
CREATE TABLE IF NOT EXISTS events (
{{[- if .Storage.Postgres ]}}
	id   UUID DEFAULT gen_random_uuid() PRIMARY KEY,
{{[- end ]}}
{{[- if .Storage.MySQL ]}}
	id   VARCHAR(36) PRIMARY KEY,
{{[- end ]}}
	name VARCHAR(128) NOT NULL
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS events;
