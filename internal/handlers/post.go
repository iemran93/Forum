package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	"forumProject/internal/database"
	"forumProject/internal/models"
)

func PostSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data models.Post
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			response := Response{Message: "Internal server error"}
			jsonResponse(w, response, http.StatusInternalServerError)
			return
		}
		userID, err := SessionActive(r)
		if err != nil {
			response := Response{Message: err.Error()}
			jsonResponse(w, response, http.StatusForbidden)
			return
		}

		data.UserID = userID
		postID, err := database.CreatePost(data)
		if err != nil {
			response := Response{Message: "Internal server error"}
			jsonResponse(w, response, http.StatusInternalServerError)
			return
		}
		postIDs := strconv.Itoa(postID)

		response := Response{Message: postIDs}
		jsonResponse(w, response, http.StatusCreated)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")

	i, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	post, err := database.GetPosts(i, "SINGLE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	comments, err := database.GetComments(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pd := PageData{
		Posts:    post,
		Comments: comments,
	}

	t, _ := template.ParseFiles("web/post.html", "web/base.html")
	t.Execute(w, pd)
}
