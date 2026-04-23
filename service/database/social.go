package database

import (
	"WasaPhoto/service/models"
)

// FollowUser creates a follow relationship if no bans exist and it's not a self-follow.
// It ensures idempotency and enforces ban restrictions in a single atomic database trip.
func (db *appdbimpl) FollowUser(followerID, targetID string) error {
	if followerID == targetID {
		return models.ErrSelfFollow
	}

	// This single query performs the following:
	// 1. Checks for an existing ban.
	// 2. If no ban exists, attempts to insert the follow record.
	// 3. If the record exists, 'ON CONFLICT' prevents an error (idempotency).
	// 4. Returns 'banned' if the WHERE clause failed due to a ban, 'inserted' otherwise.
	query := `
		WITH check_ban AS (
			SELECT EXISTS (
				SELECT 1 FROM bans 
				WHERE (banner_id = $1 AND banned_id = $2) 
				   OR (banner_id = $2 AND banned_id = $1)
			) AS is_banned
		),
		insertion AS (
			INSERT INTO follows (follower_id, followed_id)
			SELECT $1, $2
			WHERE NOT (SELECT is_banned FROM check_ban)
			ON CONFLICT (follower_id, followed_id) DO NOTHING
			RETURNING 1
		)
		SELECT 
			CASE 
				WHEN (SELECT is_banned FROM check_ban) THEN 'banned'
				ELSE 'ok'
			END AS result`

	var result string
	err := db.c.QueryRow(query, followerID, targetID).Scan(&result)
	if err != nil {
		return err // Real database error
	}

	if result == "banned" {
		// This will be mapped to 403 Forbidden in the handler
		return models.ErrForbiddenAction
	}

	// If result is 'ok', it means either a new row was inserted
	// or it already existed (idempotency).
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
