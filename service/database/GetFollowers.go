package database

func (db *appdbimpl) GetFollowers(id string) (*[]Database_follower, error) {
	table_name := "\"" + id + "_followers" + "\""
	query1 := "SELECT * FROM " + table_name
	rows, err := db.c.Query(query1)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err
	} else {
		var followerid string
		followers := []Database_follower{}
		for rows.Next() == true {
			err = rows.Scan(&followerid)
			// Si è verificato un errore nella scan
			if err != nil {
				return nil, err
			} else {
				username, err := db.UsernameFromId(followerid)
				if err != nil {
					return nil, err
				} else {
					// Username estratto correttamente dal followerid
					follower := Database_follower{FollowerId: *username}
					followers = append(followers, follower)

				}
			}

		}
		if rows.Err() != nil {
			// Si è verificato un errore durante l'iterazione delle righe o nella loro chiusura
			return nil, rows.Err()
		} else {
			return &followers, nil
		}

	}

}
