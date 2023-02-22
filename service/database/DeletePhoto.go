package database

func (db *appdbimpl) DeletePhoto(userid string, photoid string) error {
	// Eliminazione della foto dalla tabella photos dell'user
	table_name := "\"" + userid + "_photos" + "\""
	query0 := "DELETE FROM " + table_name + " WHERE photoid = ?"
	_, err := db.c.Exec(query0, photoid)
	if err != nil {
		return err
	} else {
		// Decremento del numero di foto dell'user
		query1 := "UPDATE users SET numberofphotos = numberofphotos-1 WHERE username = ?"
		username, err := db.UsernameFromId(userid)
		if err != nil {
			return err
		} else {
			_, err = db.c.Exec(query1, username)
			if err != nil {
				return err
			} else {
				// Eliminazione della tabella likes e comments della foto
				table_name := "\"" + photoid + "_likes" + "\""

				stmt1 := "DROP TABLE " + table_name
				_, err = db.c.Exec(stmt1)
				if err != nil {
					return err
				} else {
					table_name := "\"" + photoid + "_comments" + "\""

					stmt1 := "DROP TABLE " + table_name
					_, err = db.c.Exec(stmt1)
					if err != nil {
						return err
					} else {
						return nil

					}
				}

			}

		}

	}

}
