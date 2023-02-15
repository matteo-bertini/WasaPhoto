package database

func (db *appdbimpl) SetMyUsername(old_username string, new_username string) error {
	query1 := "UPDATE authstrings SET username = ? WHERE username = ?"
	_, err := db.c.Exec(query1, new_username, old_username)
	if err != nil {
		return err
	} else {
		query2 := "UPDATE users SET username = ? WHERE username = ?"
		_, err = db.c.Exec(query2, new_username, old_username)
		if err != nil {
			return err
		} else {
			return nil
		}
	}

}
