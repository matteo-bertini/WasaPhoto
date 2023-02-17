package database

func (db *appdbimpl) CommentPhoto(userid string, photoid string, commentid string, commentauthor string, commenttext string) error {
	// Aggiungo il commento alla lista dei commenti della foto
	table_name := "\"" + photoid + "_comments" + "\""
	query1 := "INSERT INTO " + table_name + " VALUES (?,?,?)"
	_, err := db.c.Exec(query1, commentid, commentauthor, commenttext)
	if err != nil {
		return err
	} else {
		// Aggiorno il numero di commenti della foto
		table_name := "\"" + userid + "_photos" + "\""
		query2 := "UPDATE " + table_name + " SET commentsnumber = commentsnumber +1 WHERE photoid = ?"
		_, err = db.c.Exec(query2, photoid)
		if err != nil {
			return err
		} else {
			return nil
		}

	}

}
