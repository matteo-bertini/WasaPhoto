package database

import (
	"WasaPhoto/service/models"
	"database/sql"
	"errors"
	"strings"
)

func (db *appdbimpl) GetUserProfile(targetUsername string, requestingUserID string) (userProfile models.UserProfile, err error) {

	// 1. Get Profile Metadata, Stats, and Relationship status
	// We use subqueries to get counts and relationship flags in a single row for maximum efficiency.
	query := `
        SELECT 
            a.username, 
            p.bio,
            (SELECT COUNT(*) FROM follows WHERE followed_id = a.user_id) AS followers_count,
            (SELECT COUNT(*) FROM follows WHERE follower_id = a.user_id) AS following_count,
            (SELECT COUNT(*) FROM posts WHERE author_id = a.user_id) AS posts_count,
            EXISTS(SELECT 1 FROM follows WHERE follower_id = ? AND followed_id = a.user_id) AS is_following,
            EXISTS(SELECT 1 FROM bans WHERE banner_id = ? AND banned_id = a.user_id) AS is_banned_by_me,
            EXISTS(SELECT 1 FROM bans WHERE banner_id = a.user_id AND banned_id = ?) AS has_banned_me,
            a.user_id
        FROM accounts a
        JOIN profiles p ON a.user_id = p.user_id
        WHERE a.username = ?`

	var targetUserID string
	var hasBannedMe bool

	err = db.c.QueryRow(query, requestingUserID, requestingUserID, requestingUserID, targetUsername).Scan(
		&userProfile.Username,
		&userProfile.Bio,
		&userProfile.FollowersCount,
		&userProfile.FollowingCount,
		&userProfile.PostsCount,
		&userProfile.IsFollowing,
		&userProfile.IsBannedByMe,
		&hasBannedMe,
		&targetUserID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserProfile{}, models.ErrUserNotFound
		}
		return models.UserProfile{}, err
	}

	// 2. Security Check: if the target user has banned the requester, return a custom error.
	// The API handler will map this specific error to a 403 Forbidden response.
	if hasBannedMe || userProfile.IsBannedByMe {
		return models.UserProfile{}, models.ErrProfileAccessForbidden
	}

	// 3. Get User's Posts including the Caption
	// We join with the accounts table to retrieve the username for each post
	// and use subqueries to count likes and comments per post.
	postQuery := `
        SELECT 
            p.post_id, 
            a.username, 
            p.caption,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.post_id) AS likes_number,
            (SELECT COUNT(*) FROM comments WHERE post_id = p.post_id) AS comments_number,
            p.created_at,
            EXISTS(SELECT 1 FROM likes WHERE post_id = p.post_id AND user_id = ?) AS is_liked_by_me
        FROM posts p
        JOIN accounts a ON p.author_id = a.user_id
        WHERE p.author_id = ?
        ORDER BY p.created_at DESC`

	rows, err := db.c.Query(postQuery, requestingUserID, targetUserID)
	if err != nil {
		return models.UserProfile{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err = rows.Scan(
			&post.PostId,
			&post.Username,
			&post.Caption,
			&post.LikesNumber,
			&post.CommentsNumber,
			&post.DateOfUpload,
			&post.IsLikedByMe,
		)
		if err != nil {
			return models.UserProfile{}, err
		}
		userProfile.UserPosts = append(userProfile.UserPosts, post)
	}

	return userProfile, nil
}

// UpdateUsername updates the username for a specific user ID.
// It returns models.ErrUsernameTaken if the new username is already in use.
func (db *appdbimpl) UpdateUsername(userID string, newUsername string) error {
	// 1. Execute the update query on the accounts table.
	// We only need to change the 'username' column where 'user_id' matches.
	query := `UPDATE accounts SET username = ? WHERE user_id = ?`

	res, err := db.c.Exec(query, newUsername, userID)
	if err != nil {
		// 2. Check if the error is a UNIQUE constraint violation (SQLite specific check)
		// This happens if another user already has the newUsername.
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return models.ErrUsernameTaken
		}
		// Return any other generic database error
		return err
	}

	// 3. Verify that a row was actually updated
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// If no rows were changed, it means the userID doesn't exist
		return models.ErrUserNotFound
	}

	return nil
}

// DeleteUser removes a user account from the database.
// Due to the ON DELETE CASCADE constraints defined in the schema,
// deleting the account will automatically remove all associated:
// profiles, posts, comments, likes, follows, and bans.
func (db *appdbimpl) DeleteUser(userID string) error {
	// Execute the DELETE statement on the accounts table
	res, err := db.c.Exec("DELETE FROM accounts WHERE user_id = ?", userID)
	if err != nil {
		// Return the database error to the handler
		return err
	}

	// Check if any row was actually deleted
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, it means the user_id didn't exist
	if affected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}
