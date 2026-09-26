package repositories

import (
	"database/sql"
	"time"
)

// deletes user sessions
func DeleteUserSessions(tx *sql.Tx, userID int) error {
	_, err := tx.Exec("DELETE FROM user_sessions WHERE user_id = ?", userID)

	return err
}

// inserts a session
func CreateSession(tx *sql.Tx, userID int, sessionToken string, expiresAt time.Time) error {
	_, err := tx.Exec("INSERT INTO user_sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, sessionToken, expiresAt)

	return err
}

// deletes session by token
func DeleteSessionByToken(tx *sql.Tx, sessionToken string) error {
	_, err := tx.Exec("DELETE FROM user_sessions WHERE session_token = ?", sessionToken)

	return err
}
