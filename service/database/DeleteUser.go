package database

func (db *appdbimpl) DeleteUser(username string) error {
	res, err := db.c.Exec(`DELETE FROM users WHERE username=?`, username)
	if err != nil {
		return err
	}
	res, err = db.c.Exec(`DELETE FROM authstrings WHERE username=?`, username)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the user didn't exist
		return ErrUserDoesNotExist
	}
	return nil
}
