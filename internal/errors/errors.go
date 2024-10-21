package errors

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
	ErrInvalidPassword = errors.New("Incorrect password")
	// Other errors
)