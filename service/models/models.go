package models

import (
	"regexp"
	"time"
)

// usernamePattern and passwordPattern mirror exactly the `pattern` constraints
// declared in doc/api.yaml for the Username and LoginRequest.password schemas.
var (
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	passwordPattern = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=\-]{8,72}$`)
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
type UserResponse struct {
	Username string `json:"username"`
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

// Post represents a single post entry
type StreamPost struct {
	PostId         string    `json:"postId"`
	Username       string    `json:"username"`
	Caption        string    `json:"caption"`
	LikesNumber    int       `json:"likesNumber"`
	CommentsNumber int       `json:"commentsNumber"`
	DateOfUpload   time.Time `json:"dateOfUpload"`
	IsLikedByMe    bool      `json:"isLikedByMe"`
	IsSuggested    bool      `json:"isSuggested"`
}

type Comment struct {
	CommentID      int64     `json:"commentId,string"`
	PostID         string    `json:"postId"`
	AuthorID       string    `json:"-"` // internal only, not part of the documented Comment schema
	AuthorUsername string    `json:"authorUsername"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CheckUsername validates the username against the exact pattern declared in
// doc/api.yaml (`^[a-zA-Z0-9_]{3,30}$`), enforcing both length and charset.
func CheckUsername(username string) bool {
	return usernamePattern.MatchString(username)
}

// CheckPassword validates the password against the exact pattern declared in
// doc/api.yaml (`^[a-zA-Z0-9!@#$%^&*()_+=-]{8,72}$`), enforcing both length and charset.
func CheckPassword(password string) bool {
	return passwordPattern.MatchString(password)
}
