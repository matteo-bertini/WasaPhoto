package database

import (
	"WasaPhoto/service/models"
	"database/sql"
	"errors"
	"fmt"
)

type Database_photostream_component struct {
	PhotoStreamComponentUsername       string
	PhotoStreamComponentPhotoId        string
	PhotoStreamComponentLikesNumber    int
	PhotoStreamComponentCommentsNumber int
	PhotoStreamComponentDateOfUpload   string
}
type Database_photo struct {
	PhotoId        string
	LikesNumber    int
	CommentsNumber int
	DateOfUpload   string
}
type Database_user struct {
	Username       string
	Followers      int
	Following      int
	Numberofphotos int
	UploadedPhotos []Database_photo
}

type Database_follower struct {
	FollowerId string
}

type Database_following struct {
	Username string
}

type Database_banned struct {
	BannedId string
}
type Database_like struct {
	Username string
}
type Database_comment struct {
	CommentId     string
	CommentAuthor string
	CommentText   string
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// session
	DoLogin(username string, password string, IsSignUp bool) (*string, error)
	DoLogout(userID string) error

	// user
	GetUserProfile(targetUsername string, requestingUserID string) (UserProfile models.UserProfile, err error)

	// social
	FollowUser(followerID, targetID string) error
	UnfollowUser(followerID, targetID string) error
	BanUser(bannerID string, bannedID string) error
	UnbanUser(bannerID string, bannedID string) error

	// posts
	UploadPost(post models.Post) error
	// helpers
	GetUserIDByToken(token string) (string, error)
	GetIDByUsername(username string) (string, error)

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New creates a new instance of AppDatabase and initializes the database schema.
// It enforces referential integrity through foreign keys and sets up the
// unified table structure for accounts, profiles, posts, and follows.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Enable foreign key support in SQLite.
	// This ensures that deleting an account automatically cleans up
	// related profiles, posts, and follows (ON DELETE CASCADE).
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, fmt.Errorf("failed to enable foreign key support: %w", err)
	}

	// Define the unified database schema.

	// Table: accounts - Security and authentication.
	accountsTable := `
	CREATE TABLE IF NOT EXISTS accounts (
		user_id        TEXT NOT NULL PRIMARY KEY,
		username       TEXT NOT NULL UNIQUE,
		password_hash  TEXT NOT NULL,
		session_token  TEXT UNIQUE
	);`

	// Table: profiles - User metadata and denormalized counters.
	profilesTable := `
	CREATE TABLE IF NOT EXISTS profiles (
		user_id        TEXT NOT NULL PRIMARY KEY,
		bio            TEXT DEFAULT '',
		avatar_path TEXT,
		CONSTRAINT fk_profile_account 
			FOREIGN KEY (user_id) 
			REFERENCES accounts(user_id) 
			ON DELETE CASCADE
	);`

	// Table: posts - Unified content storage.
	postsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		post_id      TEXT NOT NULL PRIMARY KEY,
		author_id    TEXT NOT NULL,
		caption      TEXT,
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_post_author 
			FOREIGN KEY (author_id) 
			REFERENCES accounts(user_id) 
			ON DELETE CASCADE
	);`

	// Table: follows - Social relationships (Follower -> Followed).
	// Primary Key is a composite of both IDs to prevent duplicate follows.
	followsTable := `
	CREATE TABLE IF NOT EXISTS follows (
		follower_id  TEXT NOT NULL,
		followed_id  TEXT NOT NULL,
		PRIMARY KEY (follower_id, followed_id),
		CONSTRAINT fk_follower 
			FOREIGN KEY (follower_id) REFERENCES accounts(user_id) ON DELETE CASCADE,
		CONSTRAINT fk_followed 
			FOREIGN KEY (followed_id) REFERENCES accounts(user_id) ON DELETE CASCADE
	);`
	likesTable := `
    CREATE TABLE IF NOT EXISTS likes (
        post_id    INTEGER NOT NULL,
        user_id    TEXT NOT NULL,
        PRIMARY KEY (post_id, user_id),
        CONSTRAINT fk_post_liked 
            FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
        CONSTRAINT fk_user_liker 
            FOREIGN KEY (user_id) REFERENCES accounts(user_id) ON DELETE CASCADE
    );`

	commentsTable := `
    CREATE TABLE IF NOT EXISTS comments (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id    INTEGER NOT NULL,
        author_id  TEXT NOT NULL,
        content    TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT fk_post_commented 
            FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
        CONSTRAINT fk_comment_author 
            FOREIGN KEY (author_id) REFERENCES accounts(user_id) ON DELETE CASCADE
    );`

	bansTable := `
    CREATE TABLE IF NOT EXISTS bans (
        banner_id  TEXT NOT NULL,
        banned_id  TEXT NOT NULL,
        PRIMARY KEY (banner_id, banned_id),
        CONSTRAINT fk_banner 
            FOREIGN KEY (banner_id) REFERENCES accounts(user_id) ON DELETE CASCADE,
        CONSTRAINT fk_banned 
            FOREIGN KEY (banned_id) REFERENCES accounts(user_id) ON DELETE CASCADE
    );`

	// Execution block for all tables.
	queries := []string{accountsTable, profilesTable, postsTable, followsTable, likesTable, commentsTable, bansTable}
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return nil, fmt.Errorf("failed to initialize database schema: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
