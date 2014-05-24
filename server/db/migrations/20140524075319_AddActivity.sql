
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE activities (
  id    SERIAL PRIMARY KEY,
  user_id integer REFERENCES users NOT NULL,
  friend_id integer REFERENCES users,
  type smallint NOT NULL,
  message varchar(255),
  post_id integer REFERENCES posts,
  created_at timestamp NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE activities;
