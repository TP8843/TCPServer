package main

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
)

type Score struct {
	ID      uint `gorm:"primary_key"`
	Names   string
	Score   uint
	Created int64 `gorm:"autoCreateTime"`
}

type ScoreRequest struct {
	Names string `json:"Names"`
	Score uint   `json:"Score"`
}

type ScoreResponse struct {
	Success bool   `json:"Success"`
	Rank    uint   `json:"Rank"`
	Score   *Score `json:"Score"`
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

	render.JSON(w, r, ScoreResponse{
		Success: true,
		Score:   &score,
	})
}
