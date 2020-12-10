package auth

import (
	"fmt"
	"strings"
	"unicode"
)

func Validate(username string, password string) error {
	// username [a-zA-Z0-9-_]
	const minUsernameLength = 2
	const maxUsernameLength = 20
	const minPasswordLength = 2
	const maxPasswordLength = 50
  
  if (len(username) < minUsernameLength || len(username) > maxUsernameLength) || (len(password) < minPasswordLength || len(password) > maxPasswordLength) {
    return fmt.Errorf("username or password is invalid")
  }
	for _, ch := range username {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && (ch != '_') && (ch != '-') {
			return fmt.Errorf("username or password is invalid")
		}
	}
	// password
	for _, ch = range password {
		if ch < '!' || ch > '~' {
			return fmt.Errorf("username or password is invalid")
		}
  }
	return nil
}
