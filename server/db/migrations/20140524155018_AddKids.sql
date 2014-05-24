
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE kids (
  id    SERIAL PRIMARY KEY,
  parent_id integer REFERENCES users NOT NULL,
  name varchar(32) NOT NULL,
  birthday timestamp NOT NULL,
  is_boy boolean NOT NULL,
  picture bytea,
  created_at timestamp NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE kids;
