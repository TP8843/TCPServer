package main

import (
	"github.com/go-chi/render"
	"net/http"
)

type LeaderboardResponse struct {
	Success bool     `json:"Success"`
	Scores  *[]Score `json:"Scores"`
}

// Add a new game's score to the database
func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	var scores []Score

	result := db.Order("score desc").Find(&scores)
	if result.Error != nil {
		sendErrorResponse(w, r, 500, "Could not get scores from database")
	}

	render.JSON(w, r, LeaderboardResponse{
		Success: true,
		Scores:  &scores,
	})
}
