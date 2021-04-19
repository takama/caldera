-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Setup events
INSERT INTO
    events (id, name)
VALUES
    (
        '00000000-0000-0000-0000-000000000000',
        'First event'
    ) ON CONFLICT(id) DO NOTHING;

INSERT INTO
    events (id, name)
VALUES
    (
        '9800a9d0-3af9-4c7f-8059-44eeb4876592',
        'Second event'
    ) ON CONFLICT(id) DO NOTHING;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
TRUNCATE TABLE events;
