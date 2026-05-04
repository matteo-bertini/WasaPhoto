package database

import (
	"WasaPhoto/service/models"
)

// GetFollowing retrieves the list of usernames followed by a specific user ID, checking for mutual bans
func (db *appdbimpl) GetFollowing(targetId string, requesterId string) ([]models.UserResponse, bool, error) {
	query := `
		SELECT a.username 
		FROM follows f
		JOIN accounts a ON f.followed_id = a.user_id
		WHERE f.follower_id = ? 
		AND NOT EXISTS (
			SELECT 1 FROM bans 
			WHERE (banner_id = ? AND banned_id = ?) 
			   OR (banner_id = ? AND banned_id = ?)
		)`

	rows, err := db.c.Query(query, targetId, requesterId, targetId, targetId, requesterId)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	following := []models.UserResponse{}

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, false, err
		}
		following = append(following, models.UserResponse{Username: username})
	}

	if err = rows.Err(); err != nil {
		return nil, false, err
	}

	// Check if empty list is due to a ban or simply 0 following
	if len(following) == 0 {
		var banExists int
		banCheckQuery := `SELECT COUNT(*) FROM bans WHERE (banner_id = ? AND banned_id = ?) OR (banner_id = ? AND banned_id = ?)`
		err = db.c.QueryRow(banCheckQuery, requesterId, targetId, targetId, requesterId).Scan(&banExists)
		if err != nil {
			return nil, false, err
		}
		if banExists > 0 {
			return nil, true, nil // Access denied
		}
	}

	// Always return an empty slice instead of nil for JSON consistency

	return following, false, nil
}

// GetFollowers checks if a ban exists between requester and target, then returns the followers list
func (db *appdbimpl) GetFollowers(targetId string, requesterId string) ([]models.UserResponse, bool, error) {
	query := `
		SELECT a.username 
		FROM follows f
		JOIN accounts a ON f.follower_id = a.user_id
		WHERE f.followed_id = ? 
		AND NOT EXISTS (
			SELECT 1 FROM bans
			WHERE (banner_id = ? AND banned_id = ?) 
			   OR (banner_id = ? AND banned_id = ?)
		)`

	rows, err := db.c.Query(query, targetId, requesterId, targetId, targetId, requesterId)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	followers := []models.UserResponse{}
	var follower models.UserResponse

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, false, err
		}
		follower.Username = username
		followers = append(followers, follower)
	}

	if err = rows.Err(); err != nil {
		return nil, false, err
	}

	// Double check for ban: if list is empty, we verify if it's because of a ban
	// or just a user with 0 followers.
	if len(followers) == 0 {
		var banExists int
		banCheckQuery := `SELECT COUNT(*) FROM bans WHERE (banner_id = ? AND banned_id = ?) OR (banner_id = ? AND banned_id = ?)`
		err = db.c.QueryRow(banCheckQuery, requesterId, targetId, targetId, requesterId).Scan(&banExists)
		if err != nil {
			return nil, false, err
		}
		if banExists > 0 {
			return nil, true, nil // Access denied due to ban
		}
	}

	return followers, false, nil
}

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

// GetBannedUsers retrieves the list of usernames banned by a specific user ID
func (db *appdbimpl) GetBannedUsers(bannerId string) ([]models.UserResponse, error) {
	var bannedUsers []models.UserResponse
	bannedUsers = []models.UserResponse{}

	query := `
		SELECT a.username 
		FROM bans b
		JOIN accounts a ON b.banned_id = a.user_id
		WHERE b.banner_id = ?`

	rows, err := db.c.Query(query, bannerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u models.UserResponse
	for rows.Next() {
		var bannedName string
		if err := rows.Scan(&bannedName); err != nil {
			return nil, err
		}
		u.Username = bannedName
		bannedUsers = append(bannedUsers, u)
	}

	return bannedUsers, nil
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
