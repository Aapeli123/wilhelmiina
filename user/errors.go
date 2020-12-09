package user

import "errors"

var (
	// ErrUsernameTaken is thrown if username is already taken while creating an user
	ErrUsernameTaken = errors.New("This username is already taken")
	// ErrUserNotFound is thrown if user is not found when searching users from database
	ErrUserNotFound = errors.New("User not found")
)
