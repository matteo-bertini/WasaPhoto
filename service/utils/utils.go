package utils

import (
	"strings"
)

// checkUsername checks the validity of the username
func CheckUsername(username string) bool {
	length := len(username)
	if length < 3 || length > 30 {
		return false
	}
	return true
}

// parseAuthToken parses the Bearer token from the header passed (e.g parseAuthToken("Bearer fnekbk") = "fnekbk")
func ParseAuthToken(auth_header string) *string {
	// Authentication token extraction
	bearer_auth := strings.Split(auth_header, "Bearer ")
	if bearer_auth[0] == "Bearer" {
		return nil
	} else {
		return &(bearer_auth[1])

	}

}
func CheckPresence(list []string, value string) bool {
	for _, a := range list {
		if a == value {
			return true

		}

	}
	return false

}
