package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type Score struct {
	ID      uint `gorm:"primary_key"`
	Names   string
	Score   uint
	Created int64 `gorm:"autoCreateTime"`
}

type RankedScore struct {
	ID      uint   `json:"ID"`
	Names   string `json:"Names"`
	Score   uint   `json:"Score"`
	Rank    uint   `json:"Rank"`
	Created int64  `json:"Created"`
}

type ScoreRequest struct {
	Names string `json:"Names"`
	Score uint   `json:"Score"`
}

type ScoreResponse struct {
	Success bool         `json:"Success"`
	Score   *RankedScore `json:"Score"`
}

type ScoreDeletionResponse struct {
	Success bool `json:"Success"`
}

// Add a new game's score to the database
func addScore(w http.ResponseWriter, r *http.Request) {
	var scoreRequest ScoreRequest

	err := json.NewDecoder(r.Body).Decode(&scoreRequest)
	if err != nil {
		sendErrorResponse(w, r, 400, "Bad request")
		return
	}

	score := Score{
		Names: scoreRequest.Names,
		Score: scoreRequest.Score,
	}

	result := db.Create(&score)
	if result.Error != nil {
		sendErrorResponse(w, r, 500, "Could not create score")
	}

	// Get ranking of new score
	var pos uint
	result = db.Raw("SELECT count(*) FROM scores WHERE score < (SELECT score from scores where id = ?)", score.ID).Scan(&pos)
	if result.Error != nil {
		sendErrorResponse(w, r, 500, "Could not get ranking")
	}

	render.JSON(w, r, ScoreResponse{
		Success: true,
		Score: &RankedScore{
			ID:    score.ID,
			Names: score.Names,
			Score: score.Score,
			// Make rank 1 indexed, not 0 indexed
			Rank:    pos + 1,
			Created: score.Created,
		},
	})
}

func deleteScore(w http.ResponseWriter, r *http.Request) {
	scoreId := chi.URLParam(r, "id")

	result := db.Delete(&Score{}, scoreId)
	if result.Error != nil {
		sendErrorResponse(w, r, 500, "Could not delete score")
	}

	render.JSON(w, r, ScoreDeletionResponse{
		Success: true,
	})
}
