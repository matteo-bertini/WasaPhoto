package database

func (db *appdbimpl) UncommentPhoto(userid string, photoid string, commentid string, commentauthor string) error {
	table_name := "\"" + photoid + "_comments" + "\""
	query1 := "DELETE FROM " + table_name + " WHERE comment_id = ? AND comment_author = ?"
	sql_result, err := db.c.Exec(query1, commentid, commentauthor)
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
				query1 := "UPDATE " + table_name + " SET commentsnumber=commentsnumber-1 WHERE photoid = ?"
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
