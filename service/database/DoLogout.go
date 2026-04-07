package database

// DoLogout clears the session_token for a specific user, effectively logging them out.
func (db *appdbimpl) DoLogout(userID string) error {
	// We nullify the token so future GetUserIDByToken calls for this token will fail
	query := `UPDATE accounts SET session_token = NULL WHERE user_id = ?`

	_, err := db.c.Exec(query, userID)
	return err
}
