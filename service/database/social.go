package database

import (
	"WasaPhoto/service/models"
)

// FollowUser creates a follow relationship if no bans exist and it's not a self-follow.
// It ensures idempotency and enforces ban restrictions in a single atomic database trip using a transaction.
func (db *appdbimpl) FollowUser(followerID, targetID string) error {
	if followerID == targetID {
		return models.ErrSelfFollow
	}

	// We start a transaction to ensure atomicity on SQLite
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	// Defer rollback in case of error; it's a no-op if tx.Commit() is called
	defer tx.Rollback()

	// 1. Check for an existing ban in both directions
	var isBanned bool
	checkBanQuery := `
		SELECT EXISTS (
			SELECT 1 FROM bans 
			WHERE (banner_id = ? AND banned_id = ?) 
			   OR (banner_id = ? AND banned_id = ?)
		)`

	err = tx.QueryRow(checkBanQuery, followerID, targetID, targetID, followerID).Scan(&isBanned)
	if err != nil {
		return err
	}

	if isBanned {
		// This will be mapped to 403 Forbidden in the handler
		return models.ErrForbiddenAction
	}

	// 2. Attempt to insert the follow record.
	// ON CONFLICT prevents errors if the relationship already exists (idempotency).
	insertQuery := `
		INSERT INTO follows (follower_id, followed_id)
		VALUES (?, ?)
		ON CONFLICT (follower_id, followed_id) DO NOTHING`

	_, err = tx.Exec(insertQuery, followerID, targetID)
	if err != nil {
		return err
	}

	// 3. Commit the transaction to finalize changes
	if err := tx.Commit(); err != nil {
		return err
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

// CheckBanStatus verifies if a ban exists between the requester and the target user.
// Returns true if access should be denied (403).
func (db *appdbimpl) CheckBanStatus(requesterID string, targetUsername string) (bool, error) {
	var count int

	// We check if requesterID banned targetUsername OR targetUsername banned requesterID.
	// We use a subquery to get the target's ID from their username.
	query := `
		SELECT COUNT(*) 
		FROM bans 
		WHERE (banner_id = ? AND banned_id = (SELECT user_id FROM accounts WHERE username = ?))
		   OR (banner_id = (SELECT user_id FROM accounts WHERE username = ?) AND banned_id = ?)`

	err := db.c.QueryRow(query, requesterID, targetUsername, targetUsername, requesterID).Scan(&count)
	if err != nil {
		return false, err
	}

	// If count > 0, the relationship is "Forbidden"
	return count > 0, nil
}
