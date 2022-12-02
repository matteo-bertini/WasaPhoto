package database

func (db *appdbimpl) GetUserProfile(username string) (*User, error) {
	const query = `SELECT * FROM users WHERE username = ?`
	var ret User
	// Issue the query
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Read all users in the resultset
	rows.Next()
	err = rows.Scan(&ret.Username, &ret.Followers, &ret.Following, &ret.Numberofphotos)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	return &ret, nil
}
