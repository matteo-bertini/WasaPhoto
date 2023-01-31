package database

// The function AddUser adds a user u (database.User) in the database
func (db *appdbimpl) AddUser(u User) (User, error) {

	_, err := db.c.Exec(`INSERT INTO users (username,followers,following,numberofphotos) VALUES (?, ?, ?, ?)`,
		u.Username, u.Followers, u.Following, u.Numberofphotos)
	if err != nil {
		return u, err
	}
	_, err = db.c.Exec("INSERT into following_followers VALUES (?,?,?)", u.Username, "NO", "NO")
	if err != nil {
		return u, err
	}
	return u, nil
}
