package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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

func (cfg *ApiConfig) GetSongVersesWithPagination(w http.ResponseWriter, r *http.Request) {
	cfg.Logger.Info("GetSongVersesWithPagination called")

	songID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || songID <= 0 {
		cfg.Logger.WithError(err).Error("Invalid song ID")
		common.RespondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		cfg.Logger.Debug("Limit not provided or invalid, setting default to 10")
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		cfg.Logger.Debug("Offset not provided or invalid, setting default to 0")
		offset = 0
	}

	cfg.Logger.WithFields(logrus.Fields{
		"song_id": songID,
		"limit":   limit,
		"offset":  offset,
	}).Debug("Extracted query parameters")

	cfg.Logger.Debug("Querying database for song verses")
	dbVerses, err := cfg.DB.GetSongVersesWithPagination(r.Context(), database.GetSongVersesWithPaginationParams{
		ID:     int32(songID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		cfg.Logger.WithError(err).Error("Failed to fetch song verses from database")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch song verses")
		return
	}

	verses, ok := dbVerses.([]string)
	if !ok {
		cfg.Logger.Error("Unexpected data format for verses")
		common.RespondWithError(w, http.StatusInternalServerError, "Unexpected data format for verses")
		return
	}

	cfg.Logger.WithField("verse_count", len(verses)).Info("Fetched verses successfully")

	result := struct {
		ID     int32    `json:"id"`
		Verses []string `json:"verses"`
	}{
		ID:     int32(songID),
		Verses: verses,
	}

	common.RespondWithJSON(w, http.StatusOK, result)
}
