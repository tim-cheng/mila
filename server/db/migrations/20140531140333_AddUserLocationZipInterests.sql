
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE users ADD COLUMN location varchar(255);
ALTER TABLE users ADD COLUMN zip varchar(32);
ALTER TABLE users ADD COLUMN interests varchar(255);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE users DROP COLUMN location;
ALTER TABLE users DROP COLUMN zip;
ALTER TABLE users DROP COLUMN interests;
