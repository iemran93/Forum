package database

import (
	"database/sql"
	"forumProject/internal/models"
	"time"
)

func StoreSession(sessionID string, userID int, expiration time.Time) error {
	sqlStm := `INSERT OR REPLACE INTO sessions (id, user_id, expiration) VALUES (?, ?, ?)`
	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(sessionID, userID, expiration)
	return err
}

func GetSession(sessionID string) (models.SessionData, bool, error) {
	sqlStm := `SELECT user_id, expiration FROM sessions WHERE id = ?`
	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return models.SessionData{}, false, err
	}
	defer stm.Close()

	var sessionData models.SessionData
	err = stm.QueryRow(sessionID).Scan(&sessionData.UserID, &sessionData.Expiration)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SessionData{}, false, nil
		}
		return models.SessionData{}, false, err
	}

	return sessionData, true, nil
}

func DeleteSession(sessionID string) error {
	sqlStm := `DELETE FROM sessions WHERE id = ?`
	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(sessionID)
	return err
}

func DeleteUserSessions(userID int) error {
	sqlStm := `DELETE FROM sessions WHERE user_id = ?`
	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(userID)
	return err
}
