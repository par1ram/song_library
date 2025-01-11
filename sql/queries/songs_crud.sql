-- name: InsertSong :one
INSERT INTO songs (group_id, song_name, release_date, text, link)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateSong :exec
UPDATE songs SET group_id = $2, song_name = $3, text = $4, release_date = $5, link = $6 WHERE id = $1;

-- name: DeleteSong :exec
DELETE FROM songs WHERE id = $1;

-- name: UpdateSongPartial :exec
UPDATE songs
SET
    group_id = COALESCE($2, group_id),
    song_name = COALESCE($3, song_name),
    text = COALESCE($4, text),
    release_date = COALESCE($5, release_date),
    link = COALESCE($6, link)
WHERE id = $1;

-- name: GetSongsFiltered :many
SELECT *
FROM songs
WHERE ($1 IS NULL OR group_id = $1)
  AND ($2 IS NULL OR song_name ILIKE '%' || $2 || '%')
  AND ($3 IS NULL OR text ILIKE '%' || $3 || '%')
  AND ($4 IS NULL OR release_date = $4)
  AND ($5 IS NULL OR link ILIKE '%' || $5 || '%')
ORDER BY id
LIMIT $6 OFFSET $7;
