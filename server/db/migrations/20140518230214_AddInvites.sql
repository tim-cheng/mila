
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE invites (
  user1_id integer REFERENCES users,
  user2_id integer REFERENCES users,
  created_at timestamp,
  PRIMARY KEY (user1_id, user2_id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE invites;
