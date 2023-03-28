package database

func (db *appdbimpl) GetLikes(photoid string) (*[]Database_like, error) {
	table_name := "\"" + photoid + "_likes" + "\""
	query0 := "SELECT * FROM " + table_name
	rows, err := db.c.Query(query0)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err

	} else {
		var likeid string
		likes := []Database_like{}
		for rows.Next() {
			err = rows.Scan(&likeid)
			// Si è verificato un errore nella scan
			if err != nil {
				return nil, err
			} else {
				username, err := db.UsernameFromId(likeid)
				if err != nil {
					return nil, err
				} else {
					// Username estratto correttamente dal likeid
					like := Database_like{Username: *username}
					likes = append(likes, like)

				}
			}

		}
		if rows.Err() != nil {
			// Si è verificato un errore durante l'iterazione delle righe o nella loro chiusura
			return nil, rows.Err()
		} else {
			return &likes, nil
		}

	}

}
