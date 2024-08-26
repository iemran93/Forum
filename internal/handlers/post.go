package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"forumProject/internal/database"
	"forumProject/internal/models"
)

func PostSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data models.Post
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		userID, err := SessionActive(r)
		if err != nil {
			response := Response{Message: err.Error()}
			jsonResponse(w, response, http.StatusForbidden)
			return
		}

		// handle empty title or content
		title := strings.TrimSpace(data.Title)
		content := strings.TrimSpace(data.Content)
		if title == "" || content == "" {
			response := Response{Message: "Title and content cannot be empty"}
			jsonResponse(w, response, http.StatusBadRequest)
			return
		}

		data.UserID = userID
		postID, err := database.CreatePost(data)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		postIDs := strconv.Itoa(postID)

		response := Response{Message: postIDs}
		jsonResponse(w, response, http.StatusCreated)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")

	i, err := strconv.Atoi(postID)
	if err != nil {
		RenderErrorPage(w, http.StatusNotFound, fmt.Sprintf("Post not found: %v", err))
		return
	}

	post, err := database.GetPosts(i, "SINGLE")
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}
	comments, err := database.GetComments(i)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}

	loggedIn := false
	_, err = SessionActive(r)
	if err == nil {
		loggedIn = true
	}

	pd := PageData{
		Posts:    post,
		Comments: comments,
		LoggedIn: loggedIn,
	}

	t, err := template.ParseFiles("web/post.html", "web/base.html")
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}
	t.Execute(w, pd)
}
