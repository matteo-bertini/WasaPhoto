package database

import (
	"WasaPhoto/service/models"
)

// FollowUser creates a follow relationship if no bans exist and it's not a self-follow.
func (db *appdbimpl) FollowUser(followerID, targetID string) error {
	if followerID == targetID {
		return models.ErrSelfFollow
	}

	// The query inserts only if a ban record does not exist between these two users.
	query := `
		INSERT INTO followers (follower_id, followed_id)
		SELECT ?, ?
		WHERE NOT EXISTS (
			SELECT 1 FROM bans 
			WHERE (banner_id = ? AND banned_id = ?) 
			   OR (banner_id = ? AND banned_id = ?)
		)`

	result, err := db.c.Exec(query, followerID, targetID, followerID, targetID, targetID, followerID)
	if err != nil {
		return err // General DB error
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		// If no rows were affected, it means the WHERE NOT EXISTS failed (a ban exists)
		return models.ErrForbiddenAction
	}

	return nil
}

// UnfollowUser removes a follow relationship.
func (db *appdbimpl) UnfollowUser(followerID, targetID string) error {
	query := `DELETE FROM followers WHERE follower_id = ? AND followed_id = ?`
	_, err := db.c.Exec(query, followerID, targetID)
	return err
}
