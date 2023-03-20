package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

type Database_photostream_component struct {
	Username       string
	PhotoId        string
	LikesNumber    int
	CommentsNumber int
	DateOfUpload   string
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
	// se l'username non è registrato verrà creato e restituito un nuovo id,altrimenti verrò resituito quello esistente //
	DoLogin(username string) (*string, error)

	// AddUser crea ed aggiunge il profilo dell'username //
	AddUser(username string, id string) error

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

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		// Creazione della tabella authstrings
		// authstrings memorizza per ogni username registrato l'id univoco che riconosce l'utente nel sistema e nelle richieste
		sqlStmt := `CREATE TABLE authstrings (username TEXT NOT NULL PRIMARY KEY,id TEXT NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("Errore nella creazione della tabella authstrings: %w", err)
		}
		// Creazione della tabella users
		// users memorizza  il profilo per ogni username registrato in authstrings
		sqlStmt = `CREATE TABLE users (username TEXT NOT NULL PRIMARY KEY,followers INTEGER NOT NULL,following INTEGER NOT NULL,numberofphotos INTEGER NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
