package handlers

import (
	"encoding/json"
	"fmt"
	"forumProject/internal/database"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID, err := SessionActive(r)
		if err != nil {
			response := Response{Message: "Login to add likes!", Error: err.Error()}
			jsonResponse(w, response, http.StatusMethodNotAllowed)
			return
		}
		var data struct {
			Type       string `json:"type"`
			ID         int    `json:"id"`
			EntityType string `json:"entityType"`
		}
		json.NewDecoder(r.Body).Decode(&data)
		likes, err := database.Liking(data.EntityType, data.Type, data.ID, userID)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		response := Response{
			Message: fmt.Sprintf("%d,%d", likes[0], likes[1]),
		}
		jsonResponse(w, response, http.StatusOK)
	} else {
		response := Response{Message: "Invalid request method"}
		jsonResponse(w, response, http.StatusMethodNotAllowed)
	}
}
