package models

import "errors"

var (
	ErrMissingFields             = errors.New("all fields are required: username, password, and isSignUp")
	ErrInvalidUsernameOrPassword = errors.New("username must be 3-16 characters and alphanumeric and  password must be at least 8 characters long")
	ErrUnauthorized              = errors.New("unauthorized: invalid or missing session token")
	ErrUserNotFound              = errors.New("user not found")
	ErrProfileAccessForbidden    = errors.New("access to this profile is restricted")
	ErrForbiddenAction           = errors.New("forbidden action")
	ErrSelfFollow                = errors.New("users cannot follow themselves")
	ErrSelfBan                   = errors.New("users cannot ban themselves")
	ErrResourceNotFound          = errors.New("resource not found")
	ErrUsernameTaken             = errors.New("username taken")
)

// GetUserIDByToken
var ErrInvalidToken error = errors.New("invalid SessionToken.")

// DoLogin //
var ErrUserAlreadyExists error = errors.New("user already exists.")
var ErrInvalidCredentials = errors.New("invalid username or password")

// CheckAuthorization //
var ErrAuthorizationNotSpecified error = errors.New("Authorization non specificata nell'header.")
var ErrBearerTokenNotSpecifiedWell error = errors.New("Bearer Token non specificato correttamente nel campo Authorization dell'header.")
