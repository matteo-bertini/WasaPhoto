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
		INSERT INTO follows (follower_id, followed_id)
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
	query := `DELETE FROM follows WHERE follower_id = ? AND followed_id = ?`
	_, err := db.c.Exec(query, followerID, targetID)
	return err
}

// BanUser creates a ban and removes any existing follows between the two users
func (db *appdbimpl) BanUser(bannerID string, bannedID string) error {
	if bannerID == bannedID {
		return models.ErrSelfBan
	}

	// Start a transaction to ensure atomicity
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	// 1. Insert the ban relationship
	// INSERT OR IGNORE avoids errors if the ban already exists
	_, err = tx.Exec("INSERT OR IGNORE INTO bans (banner_id, banned_id) VALUES (?, ?)", bannerID, bannedID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// 2. Remove any follow relationship in BOTH directions
	// If User A bans User B, they shouldn't follow each other anymore
	_, err = tx.Exec(`
		DELETE FROM follows 
		WHERE (follower_id = ? AND followed_id = ?) 
		   OR (follower_id = ? AND followed_id = ?)`,
		bannerID, bannedID, bannedID, bannerID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// UnbanUser removes the ban record
func (db *appdbimpl) UnbanUser(bannerID string, bannedID string) error {
	_, err := db.c.Exec("DELETE FROM bans WHERE banner_id = ? AND banned_id = ?", bannerID, bannedID)
	return err
}
