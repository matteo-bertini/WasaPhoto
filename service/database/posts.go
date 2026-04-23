package database

import (
	"WasaPhoto/service/models"
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

// GetLikes retrieves the list of usernames who liked a specific post.
func (db *appdbimpl) GetLikes(postID string, requestingUserID string) ([]string, error) {
	// Ban check: requester cannot see likes if a ban relationship exists with the post author.
	var banned bool
	queryBan := `
		SELECT EXISTS (
			SELECT 1 FROM bans 
			WHERE (banner_id = ? AND banned_id = (SELECT author_id FROM posts WHERE post_id = ?))
			   OR (banner_id = (SELECT author_id FROM posts WHERE post_id = ?) AND banned_id = ?)
		)`
	err := db.c.QueryRow(queryBan, requestingUserID, postID, postID, requestingUserID).Scan(&banned)
	if err != nil {
		return nil, err
	}
	if banned {
		return nil, models.ErrForbiddenAction
	}

	rows, err := db.c.Query(`
		SELECT a.username 
		FROM likes l 
		JOIN accounts a ON l.user_id = a.user_id 
		WHERE l.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		usernames = append(usernames, u)
	}
	return usernames, nil
}

// AddComment adds a comment if no ban exists and returns the new comment object.
func (db *appdbimpl) AddComment(postID string, authorID string, content string) (models.Comment, error) {
	// Check for bans between the commenter and the post owner.
	var banned bool
	queryBan := `
		SELECT EXISTS (
			SELECT 1 FROM bans 
			WHERE (banner_id = ? AND banned_id = (SELECT author_id FROM posts WHERE post_id = ?))
			   OR (banner_id = (SELECT author_id FROM posts WHERE post_id = ?) AND banned_id = ?)
		)`
	err := db.c.QueryRow(queryBan, authorID, postID, postID, authorID).Scan(&banned)
	if err != nil {
		return models.Comment{}, err
	}
	if banned {
		return models.Comment{}, models.ErrForbiddenAction
	}

	// Insert the comment.
	res, err := db.c.Exec(`
		INSERT INTO comments (post_id, author_id, content, created_at) 
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)`, postID, authorID, content)
	if err != nil {
		return models.Comment{}, err
	}

	lastID, _ := res.LastInsertId()

	// Get the full object to return to the frontend.
	var c models.Comment
	err = db.c.QueryRow(`
		SELECT c.comment_id, c.post_id, c.author_id, a.username, c.content, c.created_at
		FROM comments c
		JOIN accounts a ON c.author_id = a.user_id
		WHERE c.comment_id = ?`, lastID).Scan(&c.CommentID, &c.PostID, &c.AuthorID, &c.AuthorUsername, &c.Content, &c.CreatedAt)

	return c, err
}

// GetComments retrieves all comments for a post, checking for bans.
func (db *appdbimpl) GetComments(postID string, requestingUserID string) ([]models.Comment, error) {
	var banned bool
	queryBan := `
		SELECT EXISTS (
			SELECT 1 FROM bans 
			WHERE (banner_id = ? AND banned_id = (SELECT author_id FROM posts WHERE post_id = ?))
			   OR (banner_id = (SELECT author_id FROM posts WHERE post_id = ?) AND banned_id = ?)
		)`
	_ = db.c.QueryRow(queryBan, requestingUserID, postID, postID, requestingUserID).Scan(&banned)
	if banned {
		return nil, models.ErrForbiddenAction
	}

	rows, err := db.c.Query(`
		SELECT c.comment_id, c.post_id, c.author_id, a.username, c.content, c.created_at
		FROM comments c
		JOIN accounts a ON c.author_id = a.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var cm models.Comment
		if err := rows.Scan(&cm.CommentID, &cm.PostID, &cm.AuthorID, &cm.AuthorUsername, &cm.Content, &cm.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, cm)
	}
	return comments, nil
}

// DeleteComment removes a comment if the user is the author of the comment or the post owner.
func (db *appdbimpl) DeleteComment(commentID int, requesterID string) error {
	res, err := db.c.Exec(`
		DELETE FROM comments 
		WHERE comment_id = ? 
		AND (author_id = ? OR post_id IN (SELECT post_id FROM posts WHERE author_id = ?))`,
		commentID, requesterID, requesterID)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff == 0 {
		return models.ErrResourceNotFound
	}
	return nil
}
