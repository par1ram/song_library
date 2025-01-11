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

func (cfg *ApiConfig) InsertSong(w http.ResponseWriter, r *http.Request) {
	cfg.Logger.Info("InsertSong called")

	var req struct {
		GroupName string `json:"group"`
		SongName  string `json:"song"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.WithError(err).Error("Invalid request payload")
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cfg.Logger.WithFields(logrus.Fields{
		"group": req.GroupName,
		"song":  req.SongName,
	}).Debug("Decoded request payload")

	groupID, err := cfg.DB.GetGroupIDByGroupName(r.Context(), req.GroupName)
	if err != nil {
		cfg.Logger.WithError(err).WithField("group", req.GroupName).Error("Group not found")
		common.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}

	externalApiURL := common.GetExternalApiURL()
	reqURL := externalApiURL + "/info?group=" + req.GroupName + "&song=" + req.SongName
	cfg.Logger.WithField("url", reqURL).Debug("Fetching external API details")

	resp, err := cfg.HTTPClient.Get(reqURL)
	if err != nil {
		cfg.Logger.WithError(err).Error("Failed to fetch song details from external API")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch song details from external API")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		cfg.Logger.WithField("status_code", resp.StatusCode).Error("External API returned an error")
		common.RespondWithError(w, http.StatusInternalServerError, "External API returned an error")
		return
	}

	var songDetails struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&songDetails); err != nil {
		cfg.Logger.WithError(err).Error("Failed to parse external API response")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to parse external API response")
		return
	}

	cfg.Logger.WithFields(logrus.Fields{
		"release_date": songDetails.ReleaseDate,
		"text":         songDetails.Text,
		"link":         songDetails.Link,
	}).Debug("Parsed external API response")

	id, err := cfg.DB.InsertSong(r.Context(), database.InsertSongParams{
		GroupID:     groupID,
		SongName:    req.SongName,
		ReleaseDate: sql.NullTime{Time: parseDate(songDetails.ReleaseDate), Valid: true},
		Text:        sql.NullString{String: songDetails.Text, Valid: true},
		Link:        sql.NullString{String: songDetails.Link, Valid: true},
	})
	if err != nil {
		cfg.Logger.WithError(err).Error("Failed to insert song")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to insert song")
		return
	}

	cfg.Logger.WithField("song_id", id).Info("Song inserted successfully")
	common.RespondWithJSON(w, http.StatusCreated, map[string]int32{"id": id})
}

func parseDate(dateStr string) time.Time {
	layout := "02.01.2006"
	parsedDate, _ := time.Parse(layout, dateStr)
	return parsedDate
}

func (cfg *ApiConfig) UpdateSong(w http.ResponseWriter, r *http.Request) {
	cfg.Logger.Info("UpdateSong called")
	var req struct {
		ID          int32  `json:"id"`
		GroupID     int32  `json:"group_id"`
		SongName    string `json:"song_name"`
		Text        string `json:"text"`
		ReleaseDate string `json:"release_date"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.WithError(err).Error("Invalid request payload")
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cfg.Logger.WithFields(logrus.Fields{
		"id":           req.ID,
		"group_id":     req.GroupID,
		"song_name":    req.SongName,
		"release_date": req.ReleaseDate,
	}).Debug("Decoded request payload for UpdateSong")

	var releaseDate sql.NullTime
	if req.ReleaseDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			cfg.Logger.WithError(err).Error("Invalid date format")
			common.RespondWithError(w, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
			return
		}
		releaseDate = sql.NullTime{Time: parsedDate, Valid: true}
	}

	err := cfg.DB.UpdateSong(r.Context(), database.UpdateSongParams{
		ID:          req.ID,
		GroupID:     req.GroupID,
		SongName:    req.SongName,
		Text:        sql.NullString{String: req.Text, Valid: req.Text != ""},
		ReleaseDate: releaseDate,
		Link:        sql.NullString{String: req.Link, Valid: req.Link != ""},
	})
	if err != nil {
		cfg.Logger.WithError(err).Error("Failed to update song")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to update song")
		return
	}

	cfg.Logger.WithField("song_id", req.ID).Info("Song updated successfully")
	w.WriteHeader(http.StatusNoContent)
}

type PatchSongRequest struct {
	GroupID     *int32  `json:"group_id,omitempty"`
	SongName    *string `json:"song_name,omitempty"`
	Text        *string `json:"text,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty"`
	Link        *string `json:"link,omitempty"`
}

func (cfg *ApiConfig) PatchSong(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID          int32   `json:"id"`
		GroupID     *int32  `json:"group_id,omitempty"`
		SongName    *string `json:"song_name,omitempty"`
		Text        *string `json:"text,omitempty"`
		ReleaseDate *string `json:"release_date,omitempty"`
		Link        *string `json:"link,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.WithError(err).Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID <= 0 {
		cfg.Logger.Error("Invalid song ID")
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	params := database.UpdateSongPartialParams{
		ID: req.ID,
	}

	if req.GroupID != nil {
		params.GroupID = *req.GroupID
	}
	if req.SongName != nil {
		params.SongName = *req.SongName
	}
	if req.Text != nil {
		params.Text = sql.NullString{String: *req.Text, Valid: true}
	}
	if req.ReleaseDate != nil {
		params.ReleaseDate = sql.NullTime{Time: parseDate(*req.ReleaseDate), Valid: true}
	}
	if req.Link != nil {
		params.Link = sql.NullString{String: *req.Link, Valid: true}
	}

	if err := cfg.DB.UpdateSongPartial(r.Context(), params); err != nil {
		cfg.Logger.WithError(err).Error("Failed to update song")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cfg *ApiConfig) DeleteSong(w http.ResponseWriter, r *http.Request) {
	cfg.Logger.Info("Starting to process delete song request")

	songID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || songID <= 0 {
		cfg.Logger.WithFields(logrus.Fields{
			"song_id": songID,
			"error":   err,
		}).Error("Received invalid song ID")
		common.RespondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	cfg.Logger.WithField("song_id", songID).Debug("Attempting to delete song")

	err = cfg.DB.DeleteSong(r.Context(), int32(songID))
	if err != nil {
		if err == sql.ErrNoRows {
			cfg.Logger.WithField("song_id", songID).Warn("Song not found")
			common.RespondWithError(w, http.StatusNotFound, "Song not found")
			return
		}
		cfg.Logger.WithFields(logrus.Fields{
			"song_id": songID,
			"error":   err,
		}).Error("Error deleting song from database")
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to delete song")
		return
	}

	cfg.Logger.WithField("song_id", songID).Info("Song successfully deleted")
	common.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Song successfully deleted"})
}
