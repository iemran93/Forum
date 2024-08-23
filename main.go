package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"forumProject/internal/database"
	"forumProject/internal/handlers"
)

func main() {
	fmt.Print("\nHi, go to http://localhost:8080/ to view the site!\n")

	// initiate the database
	database.InitDB()

	// handlers routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginFormHandler)
	http.HandleFunc("/login/submit", handlers.LoginSubmitHandler)
	http.HandleFunc("/signup", handlers.SignupFormHanlder)
	http.HandleFunc("/signup/submit", handlers.SignupSubmitHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/comment", handlers.CommentHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)
	// http.HandleFunc("/error", handlers.ErrorHandler)

	// secured routes
	http.Handle("/postform", handlers.SessionMiddleware(http.HandlerFunc(handlers.PostFormHandler)))
	http.Handle("/postform/submit", handlers.SessionMiddleware(http.HandlerFunc(handlers.PostSubmitHandler)))

	// serve css files
	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("css", r.URL.Path[len("/css/"):])
		http.ServeFile(w, r, filePath)
	})

	// serve js files
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("js", r.URL.Path[len("/js/"):])
		http.ServeFile(w, r, filePath)
	})

	// listen and serve
	http.ListenAndServe(":8080", nil)

}
