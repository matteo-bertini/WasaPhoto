package database

func (db *appdbimpl) GetUserProfile_Authcheck(authstring string) (*bool, error) {
	const query = `SELECT * FROM authstrings WHERE authentication = ?`
	var ret bool
	// Issue the query
	rows, err := db.c.Query(query, authstring)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Read the resultset
	if rows.Next() == true {
		ret = true
		return &ret, nil
	}
	ret = false
	return &ret, nil
}
