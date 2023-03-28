package database

import "WasaPhoto/service/utils"

func (db *appdbimpl) LikePhoto(userid string, photoid string, likeid string) error {
	table_name := "\"" + photoid + "_likes" + "\""
	query0 := "SELECT * FROM " + table_name + " WHERE like_id = ?"
	rows, err := db.c.Query(query0, likeid)
	if err != nil {
		return err
	} else {
		found := rows.Next()
		if found {
			err = rows.Close()
			if err != nil {
				return err
			} else {
				return utils.ErrLikeAlreadyExists

			}

		} else {
			if rows.Err() != nil {
				return rows.Err()
			}
			query1 := "INSERT INTO " + table_name + " VALUES (?)"
			_, err := db.c.Exec(query1, likeid)
			if err != nil {
				return err
			} else {
				// Aggiorno il numero di likes della foto
				table_name := "\"" + userid + "_photos" + "\""
				query1 := "UPDATE " + table_name + " SET likesnumber=likesnumber+1 WHERE photoid = ?"
				_, err = db.c.Exec(query1, photoid)
				if err != nil {
					return err
				} else {
					return nil
				}

			}

		}
	}

}
