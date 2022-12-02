package database

import "github.com/gofrs/uuid"

func (db *appdbimpl) DoLogin(username string) (*string, error) {
	query := "SELECT authentication FROM authstrings WHERE username = ?"
	var authstring string
	// Issue the query to search the username authstring
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	// if the user is registered (an authstring already exists)
	if rows.Next() == true {
		err = rows.Scan(&authstring)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, err

		}
		rows.Close()
		return &authstring, nil
		// the user is not registered
	} else {
		query = "INSERT INTO authstrings (username,authentication) VALUES (?,?)"
		id, _ := uuid.NewV4()
		authstring = id.String()
		_, err := db.c.Exec(query, username, authstring)
		if err != nil {
			return nil, err
		}
		return &authstring, nil

	}

}
