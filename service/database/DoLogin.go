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
	err := db.c.QueryRow("SELECT user_id, password_hash FROM accounts WHERE username = ?", username).Scan(&id, &hashedpassword)

	// 2) User not found.
	if errors.Is(err, sql.ErrNoRows) {

		// 2.1) User sign up.
		if IsSignUp {
			id = ksuid.New().String()
			newhashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}

			// Start transaction for sign up.
			tx, err := db.c.Begin()
			if err != nil {
				return nil, err
			}

			// Insertion into database.
			_, err = tx.Exec("INSERT INTO accounts (user_id, username, password_hash) VALUES (?, ?, ?)", id, username, string(newhashedpassword))
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			// Profile creation into the profiles table.
			_, err = tx.Exec("INSERT INTO profiles (user_id, num_followers, num_following, num_posts) VALUES (?, 0, 0, 0)", id)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			// Session token generation.
			sessionToken := ksuid.New().String()
			_, err = tx.Exec("UPDATE accounts SET session_token = ? WHERE user_id = ?", sessionToken, id)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			// Commit transaction.
			err = tx.Commit()
			if err != nil {
				return nil, err
			}

			//  Successful sign up.
			return &sessionToken, nil

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

	// Generazione session_token per il login
	sessionToken := ksuid.New().String()
	_, err = db.c.Exec("UPDATE accounts SET session_token = ? WHERE user_id = ?", sessionToken, id)
	if err != nil {
		return nil, err
	}

	// Successful login.
	return &sessionToken, nil

}
