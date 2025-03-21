package main

import (
	"github.com/go-chi/render"
	"net/http"
)

type LeaderboardResponse struct {
	Success bool           `json:"Success"`
	Scores  *[]RankedScore `json:"Scores"`
}

// Add a new game's score to the database
func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	var scores []Score

	result := db.Order("score desc").Find(&scores)
	if result.Error != nil {
		sendErrorResponse(w, r, 500, "Could not get scores from database")
	}

	var rankedScores []RankedScore
	var currentRank uint = 0
	var currentScore uint = 0
	for _, score := range scores {
		if currentScore > score.Score || currentScore == 0 {
			currentScore = score.Score
			currentRank += 1
		}

		rankedScores = append(rankedScores, RankedScore{
			ID:      score.ID,
			Score:   score.Score,
			Names:   score.Names,
			Created: score.Created,
			Rank:    currentRank,
		})
	}

	render.JSON(w, r, LeaderboardResponse{
		Success: true,
		Scores:  &rankedScores,
	})
}
