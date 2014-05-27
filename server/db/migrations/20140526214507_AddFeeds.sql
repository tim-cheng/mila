
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE feeds (
  user_id integer REFERENCES users,
  post_id integer REFERENCES posts,
  ref_user_id integer REFERENCES users,
  PRIMARY KEY (user_id, post_id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE feeds;
