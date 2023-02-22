package database

import "WasaPhoto/service/utils"

func (db *appdbimpl) FollowUser(to_add_username string, to_add_id string, username string, id string) error {
	table_name := "\"" + id + "_followers" + "\""
	query0 := "SELECT * FROM " + table_name + " WHERE id = ?"
	rows, err := db.c.Query(query0, to_add_id)
	if err != nil {
		return err
	} else {
		if rows.Next() == false {
			if rows.Err() != nil {
				return rows.Err()
			} else {
				query1 := "INSERT INTO " + table_name + " VALUES (?)"
				_, err := db.c.Exec(query1, to_add_id)
				if err != nil {
					return err
				} else {
					query2 := "UPDATE users SET followers = followers+1 WHERE username = ?"
					_, err := db.c.Exec(query2, username)
					if err != nil {
						return err
					} else {
						table_name := "\"" + to_add_id + "_following" + "\""
						query3 := "INSERT INTO " + table_name + " VALUES (?)"
						_, err := db.c.Exec(query3, id)
						if err != nil {
							return err
						} else {
							query4 := "UPDATE users SET following = following+1 WHERE username = ?"
							_, err := db.c.Exec(query4, to_add_username)
							if err != nil {
								return err
							} else {
								return nil
							}

						}

					}

				}

			}
		} else {
			err = rows.Close()
			if err != nil {
				return err
			} else {
				return utils.ErrFollowerAlreadyAdded
			}
		}
	}

}
