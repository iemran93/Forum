package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forumProject/internal/database"
	"forumProject/internal/models"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	// add the comment
	if r.Method == http.MethodPost {
		userID, err := SessionActive(r)
		if err != nil {
			respones := Response{Message: "Please login to add comments!"}
			jsonResponse(w, respones, http.StatusBadRequest)
			return
		}
		// can add a comment
		var data struct {
			Content string `json:"content"`
			PostID  string `json:"post_id"`
		}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			response := Response{Message: "Invalid request data"}
			jsonResponse(w, response, http.StatusBadRequest)
			return
		}
		postID, err := strconv.Atoi(data.PostID)
		if err != nil {
			response := Response{Message: "Internal server error"}
			jsonResponse(w, response, http.StatusInternalServerError)
			return
		}
		comment := models.Comment{
			UserID:  userID,
			PostID:  postID,
			Content: data.Content,
		}
		err = database.CreateComment(comment)
		if err != nil {
			response := Response{Message: "Internal server error"}
			jsonResponse(w, response, http.StatusInternalServerError)
			return
		}

		response := Response{Message: "Comment added successfully"}
		jsonResponse(w, response, http.StatusCreated)

	} else {
		response := Response{Message: "Invalid request method"}
		jsonResponse(w, response, http.StatusMethodNotAllowed)
	}
}
