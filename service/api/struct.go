package api

import "WasaPhoto/service/database"

// User struct represent an user in every data exchange with the external world via REST API. JSON tags have been
// added to the struct to conform to the OpenAPI specifications regarding JSON key names.
// Note: there is a similar struct in the database package. See Fountain.FromDatabase (below) to understand why.
type User struct {
	Username       string `json:"Username"`
	Followers      int    `json:"Followers"`
	Following      int    `json:"Following"`
	Numberofphotos int    `json:"Numberofphotos"`
	//UploadedPhotos []string `json:"UploadedpPhotos"`
}

// FromDatabase populates the struct with data from the database, overwriting all values.
// You might think this is code duplication, which is correct. However, it's "good" code duplication because it allows
// us to uncouple the database and API packages.
// Suppose we were using the "database.Fountain" struct inside the API package; in that case, we were forced to conform
// either the API specifications to the database package or the other way around. However, very often, the database
// structure is different from the structure of the REST API.
// Also, in this way the database package is freely usable by other packages without the assumption that structs from
// the database should somehow be JSON-serializable (or, in general, serializable).
func (u *User) FromDatabase(user database.User) {
	u.Username = user.Username
	u.Followers = user.Followers
	u.Following = user.Following
	u.Numberofphotos = user.Numberofphotos
}

// ToDatabase returns the user in a database-compatible representation
func (u *User) ToDatabase() database.User {
	return database.User{
		Username:       u.Username,
		Followers:      u.Followers,
		Following:      u.Following,
		Numberofphotos: u.Numberofphotos,
	}
}

// IsValid checks the validity of the content. In particular, coordinates should be in their range of validity, and the
// status should be either FountainStatusGood or FountainStatusFaulty. Note that the ID is not checked, as fountains
// read from requests have zero IDs as the user won't send us the ID in that way.
func (u User) IsValid() bool {
	// Checking username length
	length := len(u.Username)
	if length < 3 || length > 16 {
		return false
	} else {
		// Checking struct fields
		if u.Followers == 0 && u.Following == 0 && u.Numberofphotos == 0 {
			return true
		} else {
			return false
		}
	}

}

// doLogin operation structs
type dologinRequestBody struct {
	Username string `json:"Username"`
}
type doLoginResponseBody struct {
	Identifier string `json:"Identifier"`
}
