package handlers

import (
	"errors"
	"net/http"
	"time"

	"forumProject/internal/database"

	"github.com/gofrs/uuid/v5"
)

func GenerateSessionID() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func SetCookie(userID int) (http.Cookie, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return http.Cookie{}, err
	}

	expiration := time.Now().Add(12 * time.Hour)
	cookie := http.Cookie{
		Name:     "Session_token",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	}

	err = database.StoreSession(sessionID, userID, expiration)
	if err != nil {
		return http.Cookie{}, err
	}

	return cookie, nil
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Session_token")
		if err != nil {
			http.Error(w, "Forbidden: No session token", http.StatusForbidden)
			return
		}

		sessionID := cookie.Value
		sessionData, exists, err := database.GetSession(sessionID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !exists || time.Now().After(sessionData.Expiration) {
			http.Error(w, "Forbidden: Not exist or expired", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SessionActive(r *http.Request) (int, error) {
	cookie, err := r.Cookie("Session_token")
	if err != nil {
		return 0, errors.New("forbidden: No session token")
	}
	sessionID := cookie.Value
	sessionData, exists, err := database.GetSession(sessionID)
	if err != nil {
		return 0, err
	}
	if !exists || time.Now().After(sessionData.Expiration) {
		return 0, errors.New("forbidden: Not exist or expired")
	}

	return sessionData.UserID, nil
}
