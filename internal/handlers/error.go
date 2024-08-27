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

var errorTemplate *template.Template

func init() {
	var err error
	errorTemplate, err = template.ParseFiles("web/error.html")
	if err != nil {
		log.Fatalf("Failed to parse error template: %v", err)
	}
}

func RenderErrorPage(w http.ResponseWriter, statusCode int, message string) {
	errorData := ErrorData{
		StatusCode: statusCode,
		Message:    message,
	}

	w.WriteHeader(statusCode)

	err := errorTemplate.Execute(w, errorData)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		return
	}
}
