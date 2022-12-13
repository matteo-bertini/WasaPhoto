/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// User struct represent a user in every API call between this package and the outside world.
// Note that the internal representation of user in the database might be different.
type User struct {
	Username       string
	Followers      int
	Following      int
	Numberofphotos int
	//UploadedPhotos []string
}

var ErrUserDoesNotExist error = errors.New("user does not exist")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// DoLogin
	DoLogin(username string) (*string, error)
	// AddUser adds a user in the database
	AddUser(u User) (User, error)
	AddUser_Authcheck(username string, authstring string) (*bool, error)

	// GetUserProfile gets a user profile searched via username
	GetUserProfile(username string) (*User, error)
	GetUserProfile_Authcheck(authstring string) (*bool, error)

	// Deleteuser deletes a user from de system
	DeleteUser(username string) error
	DeleteUser_Authcheck(username string, authstring string) (*bool, error)

	// SetMyUsername modifies the username of the specified user
	SetMyUsername(username string, new_username string) (*string, error)
	SetMyUsername_Authcheck(string, authstring string) (*bool, error)

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error
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
		sqlStmt := `CREATE TABLE users (username TEXT NOT NULL PRIMARY KEY,followers INTEGER NOT NULL,following INTEGER NOT NULL,numberofphotos INTEGER NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		sqlStmt = `CREATE TABLE authstrings (username TEXT NOT NULL PRIMARY KEY,authentication TEXT);`
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
