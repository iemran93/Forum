package models

import "time"

type SessionData struct {
	Expiration time.Time
	UserID     int
}
