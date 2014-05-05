
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users (
  id    SERIAL PRIMARY KEY,
  created_at timestamp,
  email varchar(128) UNIQUE NOT NULL,
  password varchar(128) NOT NULL,
  type varchar(10) NOT NULL,
  admin boolean NOT NULL,
  first_name varchar(32) NOT NULL,
  last_name varchar(32) NOT NULL,
  num_degree1 integer,
  num_degree2 integer,
  description varchar(128),
  picture bytea
);

CREATE TABLE connections (
  user1_id integer REFERENCES users,
  user2_id integer REFERENCES users,
  created_at timestamp,
  PRIMARY KEY (user1_id, user2_id)
);

CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  created_at timestamp,
  user_id integer REFERENCES users,
  body text,
  picture bytea
);

CREATE TABLE stars (
  post_id integer REFERENCES posts,
  user_id integer REFERENCES users,
  PRIMARY KEY (post_id, user_id)
);

CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  created_at timestamp,
  post_id integer REFERENCES posts,
  user_id integer REFERENCES users,
  body text
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE users, connections, posts, stars, comments;
