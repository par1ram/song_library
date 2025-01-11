-- name: GetSongWithFiltersAndPagination :many
SELECT s.id, g.group_name, s.song_name, s.release_date, s.link
FROM songs s
JOIN groups g ON s.group_id = g.id
WHERE 
  (g.group_name ILIKE '%' || $1 || '%' OR $1 IS NULL) AND
  (s.song_name ILIKE '%' || $2 || '%' OR $2 IS NULL) AND
  (s.release_date = $3 OR $3 IS NULL)
ORDER BY s.release_date DESC
LIMIT $4 OFFSET $5;

-- name: GetSongVersesWithPagination :one
WITH verses AS (
  SELECT unnest(string_to_array(text, E'\n\n')) AS verse
  FROM songs
  WHERE id = $1
)
SELECT verse
FROM verses
LIMIT $2 OFFSET $3;
