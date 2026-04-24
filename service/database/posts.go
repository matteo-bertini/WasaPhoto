package database

import (
	"WasaPhoto/service/models"
	"database/sql"
	"errors"
)

// UploadPost inserts a new post record into the database.
// It maps the Post model to the 'posts' table schema.
// Note: image_path is intentionally omitted from the DB insert as per the
// deterministic naming convention (post_id = filename).
func (db *appdbimpl) UploadPost(post models.Post) error {
	// We use a subquery to resolve the author_id from the username provided in the struct.
	// The created_at column is populated with the DateOfUpload from the Post model
	// to ensure consistency between the application logic and the storage.
	query := `
		INSERT INTO posts (post_id, author_id, caption, created_at)
		VALUES (?, (SELECT user_id FROM accounts WHERE username = ?), ?, ?)`

	_, err := db.c.Exec(query,
		post.PostId,
		post.Username,
		post.Caption,
		post.DateOfUpload,
	)

	if err != nil {

		return err
	}

	return nil
}

// DeletePost removes a post from the DB.
// Returns models.ErrResourceNotFound if the post doesn't exist.
// Returns models.ErrForbiddenAction if the requester is not the owner.
func (db *appdbimpl) DeletePost(postID string, requesterID string) error {
	var ownerID string

	// 1. Check existence and ownership
	err := db.c.QueryRow("SELECT author_id FROM posts WHERE post_id = ?", postID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrResourceNotFound
		}
		return err
	}

	if ownerID != requesterID {
		return models.ErrForbiddenAction
	}

	// 2. Delete the post (cascading will handle likes and comments)
	_, err = db.c.Exec("DELETE FROM posts WHERE post_id = ?", postID)
	return err
}

// LikePost executes the entire validation and insertion logic in a single database round-trip.
// It returns specific domain errors for bans or missing resources, while remaining idempotent.
func (db *appdbimpl) LikePost(postID string, actorID string, targetUsername string) error {
	// We use a CTE (Common Table Expression) to gather all necessary facts in one go.
	query := `
		WITH constants AS (
			SELECT ? AS post_id, ? AS actor_id, ? AS target_username
		),
		post_info AS (
			-- Verify if the post exists and belongs to the target username
			SELECT p.author_id 
			FROM posts p 
			JOIN accounts a ON p.author_id = a.user_id
			WHERE p.post_id = (SELECT post_id FROM constants) 
			  AND a.username = (SELECT target_username FROM constants)
		),
		ban_check AS (
			-- Check for bidirectional bans between actor and author
			SELECT 1 FROM bans 
			WHERE (banner_id = (SELECT actor_id FROM constants) AND banned_id = (SELECT author_id FROM post_info))
			   OR (banner_id = (SELECT author_id FROM post_info) AND banned_id = (SELECT actor_id FROM constants))
		),
		insertion AS (
			-- Try to insert if post exists and no ban is found
			INSERT INTO likes (post_id, user_id)
			SELECT post_id, actor_id FROM constants
			WHERE EXISTS (SELECT 1 FROM post_info) 
			  AND NOT EXISTS (SELECT 1 FROM ban_check)
			ON CONFLICT DO NOTHING
			RETURNING 1
		)
		-- Final report: tells the Go code exactly what happened
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM post_info) THEN 'not_found'
				WHEN EXISTS (SELECT 1 FROM ban_check) THEN 'forbidden'
				ELSE 'ok'
			END AS status`

	var status string
	err := db.c.QueryRow(query, postID, actorID, targetUsername).Scan(&status)
	if err != nil {
		return err // Real DB error (e.g. connection lost)
	}

	// Simple switch to map the DB status to your models/errors
	switch status {
	case "not_found":
		return models.ErrResourceNotFound
	case "forbidden":
		return models.ErrForbiddenAction
	case "ok":
		return nil // Success or already liked (idempotent)
	default:
		return nil
	}
}

// UnlikePost removes a like record. It is idempotent: if the like
// doesn't exist, it returns nil to signify the desired state is reached.
func (db *appdbimpl) UnlikePost(postID string, actorID string) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE post_id = ? AND user_id = ?", postID, actorID)
	return err

}

// GetLikes retrieves usernames who liked a post, ensuring the post exists
// and no ban is active. Returns ErrResourceNotFound or ErrForbiddenAction accordingly.
func (db *appdbimpl) GetLikes(postID string, requestingUserID string) ([]string, error) {
	query := `
		WITH post_data AS (
			-- Check if the post exists and get its author
			SELECT author_id FROM posts WHERE post_id = ?
		),
		ban_check AS (
			-- Check for bidirectional bans
			SELECT EXISTS (
				SELECT 1 FROM bans 
				WHERE (banner_id = (SELECT author_id FROM post_data) AND banned_id = ?)
				   OR (banner_id = ? AND banned_id = (SELECT author_id FROM post_data))
			) AS is_banned
		)
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM post_data) THEN 'not_found'
				WHEN (SELECT is_banned FROM ban_check) THEN 'forbidden'
				ELSE 'ok'
			END AS status,
			a.username
		FROM (SELECT 1) -- Dummy row to ensure we get a result even if zero likes
		LEFT JOIN likes l ON l.post_id = ? AND (SELECT is_banned FROM ban_check) = 0
		LEFT JOIN accounts a ON l.user_id = a.user_id
		WHERE (SELECT author_id FROM post_data) IS NOT NULL OR NOT EXISTS (SELECT 1 FROM post_data);`

	rows, err := db.c.Query(query, postID, requestingUserID, requestingUserID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	var status string

	for rows.Next() {
		var username sql.NullString
		if err := rows.Scan(&status, &username); err != nil {
			return nil, err
		}

		// Handle logical errors returned by the DB
		switch status {
		case "not_found":
			return nil, models.ErrResourceNotFound
		case "forbidden":
			return nil, models.ErrForbiddenAction
		}

		// If status is 'ok' and we have a username, add it
		if username.Valid {
			usernames = append(usernames, username.String)
		}
	}

	// Safety: ensure we return [] instead of nil for JSON
	if usernames == nil {
		usernames = []string{}
	}

	return usernames, nil
}

// AddComment inserts a new comment and returns the full Comment object.
// It strictly follows the ban and existence rules in a single atomic-like flow.
func (db *appdbimpl) AddComment(postID string, authorID string, content string) (models.Comment, error) {
	// First, we check if the post exists and if there is a ban relationship.
	// Then we insert. We use a transaction to ensure we get the full object back safely.
	tx, err := db.c.Begin()
	if err != nil {
		return models.Comment{}, err
	}
	defer tx.Rollback()

	var status string
	var postAuthor string

	// Check post existence and ban status
	err = tx.QueryRow(`
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM posts WHERE post_id = ?) THEN 'not_found'
				WHEN EXISTS (
					SELECT 1 FROM bans 
					WHERE (banner_id = (SELECT author_id FROM posts WHERE post_id = ?) AND banned_id = ?)
					   OR (banner_id = ? AND banned_id = (SELECT author_id FROM posts WHERE post_id = ?))
				) THEN 'forbidden'
				ELSE 'ok'
			END AS status,
			(SELECT author_id FROM posts WHERE post_id = ?)`,
		postID, postID, authorID, authorID, postID, postID).Scan(&status, &postAuthor)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || status == "not_found" {
			return models.Comment{}, models.ErrResourceNotFound
		}
		return models.Comment{}, err
	}

	if status == "forbidden" {
		return models.Comment{}, models.ErrForbiddenAction
	}

	// Insert the comment
	res, err := tx.Exec(`
		INSERT INTO comments (post_id, author_id, content, created_at) 
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)`, postID, authorID, content)
	if err != nil {
		return models.Comment{}, err
	}

	lastID, _ := res.LastInsertId()

	// Get the full object including the AuthorUsername for the response
	var c models.Comment
	err = tx.QueryRow(`
		SELECT c.comment_id, c.post_id, c.author_id, a.username, c.content, c.created_at
		FROM comments c
		JOIN accounts a ON c.author_id = a.user_id
		WHERE c.comment_id = ?`, lastID).Scan(
		&c.CommentID, &c.PostID, &c.AuthorID, &c.AuthorUsername, &c.Content, &c.CreatedAt)

	if err != nil {
		return models.Comment{}, err
	}

	if err = tx.Commit(); err != nil {
		return models.Comment{}, err
	}

	return c, nil
}

// GetComments retrieves all comments for a post, checking for bans.
// It performs an atomic check for post existence and bidirectional bans.
// Returns ErrResourceNotFound if the post is missing, or ErrForbiddenAction if a ban is active.
func (db *appdbimpl) GetComments(postID string, requestingUserID string) ([]models.Comment, error) {
	query := `
		WITH post_data AS (
			-- Verify post existence and identify the author
			SELECT author_id FROM posts WHERE post_id = ?
		),
		ban_check AS (
			-- Check if either the requester or the author has blocked the other
			SELECT EXISTS (
				SELECT 1 FROM bans 
				WHERE (banner_id = (SELECT author_id FROM post_data) AND banned_id = ?)
				   OR (banner_id = ? AND banned_id = (SELECT author_id FROM post_data))
			) AS is_banned
		)
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM post_data) THEN 'not_found'
				WHEN (SELECT is_banned FROM ban_check) THEN 'forbidden'
				ELSE 'ok'
			END AS status,
			c.comment_id, c.post_id, c.author_id, a.username, c.content, c.created_at
		FROM (SELECT 1) -- Guaranteed single row to capture status even with 0 comments
		LEFT JOIN comments c ON c.post_id = ? AND (SELECT is_banned FROM ban_check) = 0
		LEFT JOIN accounts a ON c.author_id = a.user_id
		WHERE (SELECT author_id FROM post_data) IS NOT NULL OR NOT EXISTS (SELECT 1 FROM post_data)
		ORDER BY c.created_at ASC;`

	rows, err := db.c.Query(query, postID, requestingUserID, requestingUserID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var status string
		var cID sql.NullInt64
		var pID, aID, aName, cont sql.NullString
		var cAt sql.NullTime

		if err := rows.Scan(&status, &cID, &pID, &aID, &aName, &cont, &cAt); err != nil {
			return nil, err
		}

		// Evaluate the business logic status returned by the query
		switch status {
		case "not_found":
			return nil, models.ErrResourceNotFound
		case "forbidden":
			return nil, models.ErrForbiddenAction
		}

		// If the row contains an actual comment, map it to the model
		if cID.Valid {
			cm := models.Comment{
				CommentID:      cID.Int64,
				PostID:         pID.String,
				AuthorID:       aID.String,
				AuthorUsername: aName.String,
				Content:        cont.String,
				CreatedAt:      cAt.Time,
			}
			comments = append(comments, cm)
		}
	}

	// Ensure we return an empty slice instead of nil for consistent JSON encoding
	if comments == nil {
		comments = []models.Comment{}
	}
	return comments, nil
}

// DeleteComment removes a comment.
// It returns ErrResourceNotFound if the comment doesn't exist,
// and ErrForbiddenAction if the user is neither the author nor the post owner.
func (db *appdbimpl) DeleteComment(commentID int, requesterID string) error {
	var authorID, postAuthorID string

	// 1. Get comment author and post author in one shot
	queryCheck := `
		SELECT c.author_id, p.author_id 
		FROM comments c
		JOIN posts p ON c.post_id = p.post_id
		WHERE c.comment_id = ?`

	err := db.c.QueryRow(queryCheck, commentID).Scan(&authorID, &postAuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrResourceNotFound
		}
		return err
	}

	// 2. Permission check
	if authorID != requesterID && postAuthorID != requesterID {
		return models.ErrForbiddenAction
	}

	// 3. Actual deletion
	_, err = db.c.Exec("DELETE FROM comments WHERE comment_id = ?", commentID)
	return err
}
