package database

import "errors"

func (db *appdbimpl) Authcheck(username string, authstring string) (*bool, error) {
	query := `SELECT authentication FROM authstrings WHERE username = ?`
	rows, err := db.c.Query(query, username)
	defer rows.Close()
	if err != nil {
		// Si è verificato un errore nell'esecuzione della query
		return nil, err
	} else {
		if rows.Next() == true {
			var cmp string
			err = rows.Scan(&cmp)
			if err != nil {
				return nil, err
			} else {
				if cmp == authstring {
					ris := true
					return &ris, nil
				} else {
					ris := false
					return &ris, nil
				}

			}
		} else {
			if rows.Err() != nil {
				return nil, rows.Err()
			} else {
				return nil, errors.New("The username searched does not exists.")
			}
		}
	}

}
