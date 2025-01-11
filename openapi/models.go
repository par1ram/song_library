package database

import (
	"database/sql"
	"time"
)

type Group struct {
	ID        int32  `json:"id"`
	GroupName string `json:"group_name"`
}

type Song struct {
	ID          int32          `json:"id"`
	SongName    string         `json:"song_name"`
	ReleaseDate *time.Time     `json:"release_date,omitempty"`
	Text        sql.NullString `json:"text,omitempty"`
	Link        sql.NullString `json:"link,omitempty"`
	GroupID     int32          `json:"group_id"`
	GroupName   string         `json:"group_name"`
}

type InsertSongParams struct {
	GroupID     int32          `json:"group_id"`
	SongName    string         `json:"song_name"`
	ReleaseDate *time.Time     `json:"release_date,omitempty"`
	Text        sql.NullString `json:"text,omitempty"`
	Link        sql.NullString `json:"link,omitempty"`
}

type UpdateSongParams struct {
	ID          int32          `json:"id"`
	GroupID     int32          `json:"group_id"`
	SongName    string         `json:"song_name"`
	Text        sql.NullString `json:"text,omitempty"`
	ReleaseDate *time.Time     `json:"release_date,omitempty"`
	Link        sql.NullString `json:"link,omitempty"`
}

type GetSongsParams struct {
	GroupName string `json:"group_name,omitempty"`
	SongName  string `json:"song_name,omitempty"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type GetSongWithFiltersAndPaginationParams struct {
	Column1     sql.NullString `json:"group_name,omitempty"`
	Column2     sql.NullString `json:"song_name,omitempty"`
	ReleaseDate sql.NullTime   `json:"release_date,omitempty"`
	Limit       int32          `json:"limit"`
	Offset      int32          `json:"offset"`
}

type VerseResponse struct {
	ID     int32    `json:"id"`
	Verses []string `json:"verses"`
}

type GetSongVersesWithPaginationParams struct {
	ID     int32 `json:"id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
