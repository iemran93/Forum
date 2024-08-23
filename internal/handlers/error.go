package handlers

import (
	"log"
	"net/http"
	"text/template"
)

type ErrorData struct {
	StatusCode int
	Message    string
}

func RenderErrorPage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorData := ErrorData{
		StatusCode: statusCode,
		Message:    message,
	}

	t, err := template.ParseFiles("web/error.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, errorData)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
