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

// LikePost adds a like record if the post exists, belongs to the targetUsername, and no ban is active.
func (db *appdbimpl) LikePost(postID string, actorID string, targetUsername string) error {
	// Atomic insert: checking for post existence, target ownership, and double-sided ban.
	query := `
		INSERT INTO likes (post_id, user_id)
		SELECT ?, ?
		WHERE EXISTS (
			SELECT 1 FROM posts p 
			JOIN accounts a ON p.author_id = a.user_id
			WHERE p.post_id = ? AND a.username = ?
		)
		AND NOT EXISTS (
			SELECT 1 FROM bans 
			WHERE (banner_id = ? AND banned_id = (SELECT author_id FROM posts WHERE post_id = ?))
			   OR (banner_id = (SELECT author_id FROM posts WHERE post_id = ?) AND banned_id = ?)
		);`

	res, err := db.c.Exec(query, postID, actorID, postID, targetUsername, actorID, postID, postID, actorID)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		// Differentiate between 404 (post/user not found) and 403 (forbidden by ban)
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM posts p JOIN accounts a ON p.author_id = a.user_id WHERE p.post_id = ? AND a.username = ?)`
		_ = db.c.QueryRow(checkQuery, postID, targetUsername).Scan(&exists)
		if !exists {
			return models.ErrResourceNotFound
		}
		return models.ErrForbiddenAction
	}
	return nil
}

// UnlikePost removes a like and returns ErrNotFound if the record didn't exist.
func (db *appdbimpl) UnlikePost(postID string, actorID string) error {
	res, err := db.c.Exec("DELETE FROM likes WHERE post_id = ? AND user_id = ?", postID, actorID)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff == 0 {
		return models.ErrResourceNotFound
	}
	return nil
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
