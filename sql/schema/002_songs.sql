-- +goose Up
CREATE TABLE songs (
  id SERIAL PRIMARY KEY,
  song_name TEXT NOT NULL,
  release_date DATE,
  text TEXT,
  link TEXT,
  group_id INTEGER NOT NULL,
  FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE INDEX idx_songs_group_id ON songs (group_id);
CREATE INDEX idx_songs_song_name ON songs (song_name);
CREATE INDEX idx_songs_release_date ON songs (release_date);

-- +goose Down
DROP TABLE IF EXISTS songs;
