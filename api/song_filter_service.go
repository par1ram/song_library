package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/par1ram/song-library/common"
	"github.com/par1ram/song-library/internal/database"
	"github.com/sirupsen/logrus"
)

func (cfg *ApiConfig) GetSongWithFiltersAndPagination(w http.ResponseWriter, r *http.Request) {
	cfg.Logger.Info("GetSongWithFiltersAndPagination called")

	var req struct {
		Group       string `json:"group"`
		Song        string `json:"song"`
		ReleaseDate string `json:"release_date"`
		Limit       int32  `json:"limit"`
		Offset      int32  `json:"offset"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.WithError(err).Error("Invalid request payload")
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cfg.Logger.WithFields(logrus.Fields{
		"group":        req.Group,
		"song":         req.Song,
		"release_date": req.ReleaseDate,
		"limit":        req.Limit,
		"offset":       req.Offset,
	}).Debug("Decoded request payload")

	var releaseDate sql.NullTime
	if req.ReleaseDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			cfg.Logger.WithError(err).Error("Invalid date format")
			common.RespondWithError(w, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
			return
		}
		releaseDate = sql.NullTime{Time: parsedDate, Valid: true}
		cfg.Logger.Debug("Parsed release date successfully")
	}

	cfg.Logger.Debug("Querying database with filters")
	songs, err := cfg.DB.GetSongWithFiltersAndPagination(r.Context(), database.GetSongWithFiltersAndPaginationParams{
		Column1:     sql.NullString{String: req.Group, Valid: req.Group != ""},
		Column2:     sql.NullString{String: req.Song, Valid: req.Song != ""},
		ReleaseDate: releaseDate,
		Limit:       req.Limit,
		Offset:      req.Offset,
	})
	if err != nil {
		cfg.Logger.WithError(err).Error("Failed to fetch songs from database")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch songs")
		return
	}

	cfg.Logger.WithField("song_count", len(songs)).Info("Fetched songs successfully")
	common.RespondWithJSON(w, http.StatusOK, songs)
}

type SongVersesRequest struct {
	ID     int32 `json:"id"`
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
}

func (cfg *ApiConfig) GetSongVersesWithPagination(w http.ResponseWriter, r *http.Request) {
	cfg.Logger.Info("GetSongVersesWithPagination called")

	var req SongVersesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.WithError(err).Error("Invalid request payload")
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.ID <= 0 {
		cfg.Logger.Error("Invalid song ID")
		common.RespondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	limit := req.Limit
	if limit <= 0 {
		cfg.Logger.Debug("Limit not provided or invalid, setting default to 10")
		limit = 10
	}

	offset := req.Offset
	if offset < 0 {
		cfg.Logger.Debug("Offset not provided or invalid, setting default to 0")
		offset = 0
	}

	cfg.Logger.WithFields(logrus.Fields{
		"song_id": req.ID,
		"limit":   limit,
		"offset":  offset,
	}).Debug("Extracted query parameters")

	cfg.Logger.Debug("Querying database for song verses")
	dbVerses, err := cfg.DB.GetSongVersesWithPagination(r.Context(), database.GetSongVersesWithPaginationParams{
		ID:     req.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		cfg.Logger.WithError(err).Error("Failed to fetch song verses from database")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch song verses")
		return
	}

	dbVerse := dbVerses.(string)
	verses := []string{dbVerse}
	cfg.Logger.WithField("verse_count", len(verses)).Info("Fetched verses successfully")

	result := struct {
		ID     int32    `json:"id"`
		Verses []string `json:"verses"`
	}{
		ID:     req.ID,
		Verses: verses,
	}

	common.RespondWithJSON(w, http.StatusOK, result)
}
