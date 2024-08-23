package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"forumProject/internal/database"
	"forumProject/internal/models"
)

type PageData struct {
	Posts      []models.Post
	Comments   []models.Comment
	LoggedIn   bool
	Categories []models.Category // Add this line
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// show/hide depend on session
	loggedIn := false

	_, err := SessionActive(r)
	if err == nil {
		loggedIn = true
	}

	// get all the posts
	posts, err := database.GetPosts(0, "ALL")
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}

	categories, err := database.GetCategories() // Implement this function
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}

	pd := PageData{
		Posts:      posts,
		LoggedIn:   loggedIn,
		Categories: categories,
	}
	// serve the template with the data
	t, err := template.ParseFiles("web/index.html", "web/base.html")
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}
	err = t.Execute(w, pd)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, err := SessionActive(r)
		if err == nil {
			RenderErrorPage(w, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", "User logged in"))
			return
		}
		t, err := template.ParseFiles("web/login.html", "web/base.html")
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		data := struct {
			LoggedIn bool
		}{
			LoggedIn: false,
		}
		t.Execute(w, data)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func SignupFormHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, err := SessionActive(r)
		if err == nil {
			RenderErrorPage(w, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", "User logged in"))
			return
		}
		data := struct {
			LoggedIn bool
		}{
			LoggedIn: false,
		}
		t, err := template.ParseFiles("web/signup.html", "web/base.html")
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		t.Execute(w, data)
	} else {
		RenderErrorPage(w, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", "Invalid request method"))
		return
	}
}

func PostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		categories, err := database.GetCategories()
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}

		data := struct {
			Categories []models.Category
			LoggedIn   bool
		}{
			Categories: categories,
			LoggedIn:   true, // Since this is a secured route, the user should be logged in. Really ? No way!
		}

		t, err := template.ParseFiles("web/postform.html", "web/base.html")
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		t.Execute(w, data)
	} else {
		RenderErrorPage(w, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", "Invalid request method"))
		return
	}
}
