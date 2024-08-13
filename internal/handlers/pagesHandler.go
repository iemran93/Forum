package handlers

import (
	"log"
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
	var userID int

	userID, err := SessionActive(r)
	if err == nil {
		loggedIn = true
	}

	// to js/html logged
	log.Println(userID)

	// get all the posts
	posts, err := database.GetPosts(0, "ALL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := database.GetCategories() // Implement this function
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
	}

	pd := PageData{
		Posts:      posts,
		LoggedIn:   loggedIn,
		Categories: categories,
	}
	// serve the template with the data
	t, err := template.ParseFiles("web/index.html", "web/base.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, pd)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// render the login page
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, err := SessionActive(r)
		if err == nil {
			http.Error(w, "User logged in", http.StatusBadRequest)
			return
		}
		t, _ := template.ParseFiles("web/login.html", "web/base.html")
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
			http.Error(w, "User logged in", http.StatusBadRequest)
			return
		}
		t, _ := template.ParseFiles("web/signup.html", "web/base.html")
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

func PostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		categories, err := database.GetCategories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Categories []models.Category
			LoggedIn   bool
		}{
			Categories: categories,
			LoggedIn:   true, // Since this is a secured route, the user should be logged in
		}

		t, _ := template.ParseFiles("web/postform.html", "web/base.html")
		t.Execute(w, data)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
