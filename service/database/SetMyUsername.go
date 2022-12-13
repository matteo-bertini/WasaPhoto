package database

func (db *appdbimpl) SetMyUsername(username string, new_username string) (*string, error) {
	res, err := db.c.Exec(`UPDATE authstrings SET username = ? WHERE username=?`, new_username, username)
	if err != nil {
		return nil, err
	}
	res, err = db.c.Exec(`UPDATE users SET username = ? WHERE username = ?`, new_username, username)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	} else if affected == 0 {
		// If we didn't update any row, then the user didn't exist
		return nil, ErrUserDoesNotExist
	}
	return &new_username, nil
}
