package database

func (db *appdbimpl) UnfollowUser(username string, id string, to_del string, to_del_id string) error {
	table_name := "\"" + id + "_followers" + "\""
	query0 := "SELECT * FROM " + table_name + " WHERE id = ?"
	rows, err := db.c.Query(query0, to_del_id)
	if err != nil {
		return err
	} else {
		found := rows.Next()
		if !found {
			if rows.Err() != nil {
				return rows.Err()
			} else {
				return nil
			}
		} else {
			err = rows.Close()
			if err != nil {
				return err
			}
			query1 := "DELETE FROM " + table_name + " WHERE id = ?"
			_, err = db.c.Exec(query1, to_del_id)
			if err != nil {
				return err
			} else {
				query2 := "UPDATE users SET followers = followers -1 WHERE username = ?"
				_, err = db.c.Exec(query2, username)
				if err != nil {
					return err
				} else {
					table_name := "\"" + to_del_id + "_following" + "\""
					query3 := "DELETE FROM " + table_name + " WHERE id = ?"
					_, err := db.c.Exec(query3, id)
					if err != nil {
						return err
					} else {
						query4 := "UPDATE users SET following = following -1 WHERE username = ?"
						_, err = db.c.Exec(query4, to_del)
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

}
