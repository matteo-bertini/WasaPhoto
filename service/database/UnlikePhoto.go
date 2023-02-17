package database

func (db *appdbimpl) UnlikePhoto(userid string, photoid string, likeid string) error {
	table_name := "\"" + photoid + "_likes" + "\""
	query1 := "DELETE FROM " + table_name + " WHERE like_id = ?"
	sql_result, err := db.c.Exec(query1, likeid)
	if err != nil {
		return err
	} else {
		// Aggiorno il numero di likes della foto
		x, err := sql_result.RowsAffected()
		if err != nil {
			return err
		} else {
			if x == 1 {
				table_name := "\"" + userid + "_photos" + "\""
				query1 := "UPDATE " + table_name + " SET likesnumber=likesnumber-1 WHERE photoid = ?"
				_, err = db.c.Exec(query1, photoid)
				if err != nil {
					return err
				} else {
					return nil
				}

			} else {
				return nil
			}

		}

	}

}
