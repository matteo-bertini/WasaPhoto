package database

func (db *appdbimpl) UploadPhoto(photo Database_photo, id string) error {
	// Inserisco la nuova foto all'interno del database
	table_name := "\"" + id + "_photos" + "\""
	query1 := "INSERT INTO " + table_name + " VALUES (?,?,?,?)"
	_, err := db.c.Exec(query1, photo.PhotoId, photo.LikesNumber, photo.CommentsNumber, photo.DateOfUpload)
	if err != nil {
		return err
	} else {
		// Aggiorno il numero di foto dell'user
		query2 := "UPDATE users SET numberofphotos = numberofphotos+1 WHERE username = ?"
		username, err := db.UsernameFromId(id)
		if err != nil {
			return err
		} else {
			_, err = db.c.Exec(query2, username)
			if err != nil {
				return err
			} else {
				// La foto è stata inserita con successo nel database,creo le tabelle likes e comments per la foto
				table_name = "\"" + photo.PhotoId + "_likes" + "\""
				sqlStmt := "CREATE TABLE " + table_name + " (like_id TEXT NOT NULL PRIMARY KEY)"
				_, err = db.c.Exec(sqlStmt)
				if err != nil {
					return err
				} else {
					table_name = "\"" + photo.PhotoId + "_comments" + "\""
					sqlStmt := "CREATE TABLE " + table_name + " (user_id TEXT NOT NULL PRIMARY KEY,comment_id TEXT NOT NULL,comment_text TEXT NOT NULL)"
					_, err = db.c.Exec(sqlStmt)
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
