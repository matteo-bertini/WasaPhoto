package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
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

	// DoLogin resitituisce l'id relativo all'username passato come argomento. //
	// se l'username non è registrato verrà creato e restituito un nuovo id,altrimenti verrà resituito quello esistente //
	DoLogin(username string, password string, IsSignUp bool) (*string, error)

	// GetUserProfile gets a user profile searched via username //
	GetUserProfile(username string) (*Database_user, error)

	// GetMyStream //
	GetMyStream(id string) (*[]Database_photostream_component, error)

	// Deleteuser elimina completamente un utente dal sistema //
	DeleteUser(id string, username string) error

	// SetMyUsername modifica l'username dell'user con username passato come argomento //
	SetMyUsername(old_username string, new_username string) error

	// GetFollowers //
	GetFollowers(id string) (*[]Database_follower, error)

	// GetFollowing //
	GetFollowing(id string) (*[]Database_following, error)

	// GetBanned //
	GetBanned(id string) (*[]Database_banned, error)

	// FollowUser //
	FollowUser(to_add_username string, to_add_id string, username string, id string) error

	// UnfollowUser //
	UnfollowUser(username string, id string, to_del string, to_del_id string) error

	// BanUser //
	BanUser(username string, id string, to_ban_username string, to_ban_id string) error

	// UnbanUser //
	UnbanUser(id string, to_del_id string) error

	// UploadPhoto //
	UploadPhoto(photo Database_photo, id string) error

	// DeletePhoto //
	DeletePhoto(userid string, photoid string) error

	// GetLikes //
	GetLikes(photoid string) (*[]Database_like, error)

	// LikePhoto //
	LikePhoto(userid string, photoid string, likeid string) error

	// UnlikePhoto //
	UnlikePhoto(userid string, photoid string, likeid string) error

	// GetComments //
	GetComments(photoid string) (*[]Database_comment, error)

	// CommentPhoto //
	CommentPhoto(userid string, photoid string, commentid string, commentauthor string, commenttext string) error

	// UncommentPhoto //
	UncommentPhoto(userid string, photoid string, commentid string, commentauthor string) error

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error

	// Funzioni ausiliarie definite in database_utilities
	CheckAuthorization(request *http.Request, username string) error
	CheckUserExistence(username string) error
	IdFromUsername(username string) (*string, error)
	UsernameFromId(id string) (*string, error)
	IsAllowed(id1 string, id2 string) error
	CheckPhotoExistence(user_id string, photoid string) error
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
		num_followers  INTEGER DEFAULT 0,
		num_following  INTEGER DEFAULT 0,
		num_posts      INTEGER DEFAULT 0,
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
		image_path   TEXT NOT NULL,
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

	// Execution block for all tables.
	queries := []string{accountsTable, profilesTable, postsTable, followsTable}
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
