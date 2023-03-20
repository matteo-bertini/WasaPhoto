package database

func (db *appdbimpl) GetBanned(id string) (*[]Database_banned, error) {
	table_name := "\"" + id + "_bans" + "\""
	query1 := "SELECT * FROM " + table_name
	rows, err := db.c.Query(query1)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err
	} else {
		var bannedid string
		bannedusers := []Database_banned{}
		for rows.Next() == true {
			err = rows.Scan(&bannedid)
			// Si è verificato un errore nella scan
			if err != nil {
				return nil, err
			} else {
				username, err := db.UsernameFromId(bannedid)
				if err != nil {
					return nil, err
				} else {
					// Username estratto correttamente dal followerid
					banned := Database_banned{BannedId: *username}
					bannedusers = append(bannedusers, banned)

				}
			}

		}
		if rows.Err() != nil {
			// Si è verificato un errore durante l'iterazione delle righe o nella loro chiusura
			return nil, rows.Err()
		} else {
			return &bannedusers, nil
		}

	}

}
