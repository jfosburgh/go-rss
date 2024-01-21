-- +goose Up
CREATE TABLE posts(
  id UUID NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  description TEXT,
  published_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULL,
  CONSTRAINT fk_users
  FOREIGN KEY (feed_id)
  REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
