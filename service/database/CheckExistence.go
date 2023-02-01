package database

func (db *appdbimpl) CheckExistence(username string) (error, *bool) {
	query := "SELECT * FROM users WHERE username = ?"
	rows, err := db.c.Query(query, username)
	if err != nil {
		return err, nil
	} else {
		ris := rows.Next()
		err = rows.Close()
		if err != nil {
			return err, nil
		} else {
			if ris == true {
				to_ret := true
				return nil, &to_ret
			} else {
				if rows.Err() != nil {
					return rows.Err(), nil
				} else {
					to_ret := false
					return nil, &to_ret
				}
			}
		}

	}

}
