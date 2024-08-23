package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"forumProject/internal/database"
	"forumProject/internal/functions"
	"forumProject/internal/models"
)

func AuthenticateUser(uname string, psw string) (int, error) {
	users, err := database.GetUsers()
	if err != nil {
		return 0, err
	}

	for _, user := range users {
		if user.Username == uname {
			if functions.CheckPasswordHash(psw, user.Password) {
				return user.ID, nil
			}
		}
	}

	return 0, errors.New("invalid credentials")
}

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data models.User
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			RenderErrorPage(w, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", err))
			return
		}
		user, err := AuthenticateUser(data.Username, data.Password)
		if err != nil {
			response := Response{Message: "Unauthorized"}
			jsonResponse(w, response, http.StatusUnauthorized)
			return
		} else {
			// Delete any existing sessions for this user
			err = database.DeleteUserSessions(user)
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
				return
			}

			// Create new session
			cookie, err := SetCookie(user)
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
				return
			}
			http.SetCookie(w, &cookie)
			response := Response{Message: "Login successful"}
			jsonResponse(w, response, http.StatusCreated)
		}
	} else {
		response := Response{Message: "Invalid request method"}
		jsonResponse(w, response, http.StatusMethodNotAllowed)
	}
}

func SignupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data models.User
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}

		if data.Password != data.PasswordRep {
			response := Response{Message: "Passwords do not match"}
			jsonResponse(w, response, http.StatusBadRequest)
			return
		}

		pass, errmsg := functions.ValidUserData(data.Username, data.Email)
		if !pass {
			response := Response{Message: errmsg}
			jsonResponse(w, response, http.StatusBadRequest)
			return
		}

		users, err := database.GetUsers()
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		for _, user := range users {
			if user.Username == data.Username {
				response := Response{Message: "Username is already taken"}
				jsonResponse(w, response, http.StatusBadRequest)
				return
			}
			if user.Email == data.Email {
				response := Response{Message: "Email is already registered"}
				jsonResponse(w, response, http.StatusBadRequest)
				return
			}
		}

		hash, err := functions.HashPassword(data.Password)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}

		data.Password = hash

		err = database.CreateUser(data)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		response := Response{Message: "Signup successful"}
		jsonResponse(w, response, http.StatusCreated)
	} else {
		response := Response{Message: "Invalid request method"}
		jsonResponse(w, response, http.StatusMethodNotAllowed)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Session_token")
	if err != nil {
		RenderErrorPage(w, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", err))
		return
	}

	sessionID := cookie.Value
	err = database.DeleteSession(sessionID)
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	response := Response{Message: "Logout successful"}
	jsonResponse(w, response, http.StatusOK)
}
