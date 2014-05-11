
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE posts ADD COLUMN bg_color varchar(8);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE posts DROP COLUMN bg_color;
