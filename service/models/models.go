package models

import (
	"strings"
	"time"
)

// CheckUsername controlla che l'Username passato sia una stringa conforme alle specifiche dichiarate
// La funzione ritorna true quando l'Username passato è valido,false altrimenti.
// UserProfile represents the full profile data returned to the frontend
type UserProfile struct {
	Username       string `json:"Username"`
	Bio            string `json:"Bio"`
	FollowersCount int    `json:"FollowersCount"`
	FollowingCount int    `json:"FollowingCount"`
	PostsCount     int    `json:"PostsCount"`
	IsFollowing    bool   `json:"IsFollowing"`
	IsBannedByMe   bool   `json:"IsBannedByMe"`
	UserPosts      []Post `json:"UserPosts"`
}

// Post represents a single post entry
type Post struct {
	PostId         string    `json:"PostId"`
	Username       string    `json:"Username"`
	Caption        string    `json:"Caption"`
	LikesNumber    int       `json:"LikesNumber"`
	CommentsNumber int       `json:"CommentsNumber"`
	DateOfUpload   time.Time `json:"DateOfUpload"`
	IsLikedByMe    bool      `json:"IsLikedByMe"`
}

// doLogin operation //
type DoLoginRequestBody struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	IsSignUp *bool  `json:"IsSignUp"`
}
type DoLoginResponseBody struct {
	SessionToken string `json:"SessionToken"`
}

func CheckUsername(username string) bool {

	// L'Username passato è composto solo da spazi bianchi quindi non è valido
	if strings.TrimSpace(username) == "" {
		return false
	} else {
		len := len(username)
		if len > 30 || len < 3 {
			return false
		} else {
			return true
		}
	}
}

func CheckPassword(password string) bool {
	// La password passata è composta solo da spazi bianchi quindi non è valida
	if strings.TrimSpace(password) == "" {
		return false
	} else {
		len := len(password)
		if len < 8 {
			return false
		} else {
			return true
		}
	}

}
