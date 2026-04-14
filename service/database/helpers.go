package database

import (
	"WasaPhoto/service/models"
	"database/sql"
	"errors"
)

// GetUserIDByToken retrieves the user_id associated with a specific session_token.
// It returns an error if the token is not found or has been invalidated (set to NULL).
func (db *appdbimpl) GetUserIDByToken(token string) (string, error) {
	var userID string
	// Check if the token exists in the accounts table
	query := `SELECT user_id FROM accounts WHERE session_token = ?`

	err := db.c.QueryRow(query, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Token is missing or session has been cleared
			return "", models.ErrInvalidToken
		}
		return "", err
	}
	return userID, nil
}

// GetIDByUsername retrieves the UserID for a given username.
func (db *appdbimpl) GetIDByUsername(username string) (string, error) {
	var id string
	query := `SELECT user_id FROM accounts WHERE username = ?`
	err := db.c.QueryRow(query, username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrUserNotFound
		}
		return "", err
	}
	return id, nil
}
