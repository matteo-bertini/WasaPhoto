package models

import (
	"strings"
	"time"
)

// DoLoginRequestBody represents the 'LoginRequest' schema defined in the OpenAPI specification.
type DoLoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsSignUp *bool  `json:"isSignUp"`
}

// DoLoginResponseBody represents the 'LoginResponse' schema defined in the OpenAPI specification.
type DoLoginResponseBody struct {
	SessionToken string `json:"sessionToken"`
}

// UserProfile represents the full profile data returned to the frontend
type UserProfile struct {
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	FollowersCount int    `json:"followersCount"`
	FollowingCount int    `json:"followingCount"`
	PostsCount     int    `json:"postsCount"`
	IsFollowing    bool   `json:"isFollowing"`
	IsBannedByMe   bool   `json:"isBannedByMe"`
	UserPosts      []Post `json:"userPosts"`
}

// Post represents a single post entry
type Post struct {
	PostId         string    `json:"postId"`
	Username       string    `json:"username"`
	Caption        string    `json:"caption"`
	LikesNumber    int       `json:"likesNumber"`
	CommentsNumber int       `json:"commentsNumber"`
	DateOfUpload   time.Time `json:"dateOfUpload"`
	IsLikedByMe    bool      `json:"isLikedByMe"`
}
type Comment struct {
	CommentID      int64     `json:"commentId,string"`
	PostID         string    `json:"postId"`
	AuthorID       string    `json:"authorId"`
	AuthorUsername string    `json:"authorUsername"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
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
		if len < 8 || len > 72 {
			return false
		} else {
			return true
		}
	}

}
