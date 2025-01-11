-- +goose Up
CREATE TABLE groups (
  id SERIAL PRIMARY KEY,
  group_name TEXT UNIQUE NOT NULL
);

CREATE INDEX idx_groups_group_name ON groups (group_name);

-- +goose Down
DROP TABLE IF EXISTS groups;
