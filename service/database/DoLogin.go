package database

import (
	"WasaPhoto/service/utils"
	"database/sql"
	"errors"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

func (db *appdbimpl) DoLogin(username string, password string, IsSignUp bool) (*string, error) {

	var id, hashedpassword string

	// 1) Searching the user.
	err := db.c.QueryRow("SELECT id,hashedpassword FROM authstrings WHERE username = ?", username).Scan(&id, &hashedpassword)

	// 2) User not found.
	if errors.Is(err, sql.ErrNoRows) {

		// 2.1) User sign up.
		if IsSignUp {
			id = ksuid.New().String()
			newhashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			// Insertion into database.
			_, err = db.c.Exec("INSERT INTO authstrings (username,hashedpassword,id) VALUES (?,?,?)", username, string(newhashedpassword), id)
			if err != nil {
				return nil, err

			}
			_, err = db.c.Exec("INSERT INTO users (username,followers,following,numberofphotos) VALUES (?,0,0,0)", username)
			if err != nil {
				return nil, err
			}

			// TO ADD -> Other table creation
			//  Successful sign up.
			return &id, nil

		}
		// 2.2) Login attempt.
		return nil, utils.ErrUserDoesNotExist

	}
	if err != nil {
		// 3) Generic database error.
		return nil, err

	}
	// 4) User found.
	if IsSignUp {
		// 4.1) Registration attempt.
		return nil, utils.ErrUserAlreadyExists
	}
	// 4.2) Login attempt.
	err = bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}
	// Successful login.
	return &id, nil

}
