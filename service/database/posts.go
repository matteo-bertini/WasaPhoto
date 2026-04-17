package database

import (
	"WasaPhoto/service/models"
	"fmt"
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
		// Return the error with context to allow the handler to perform cleanup (os.Remove)
		return fmt.Errorf("database.UploadPost: failed to execute insert for post %s: %w", post.PostId, err)
	}

	return nil
}
