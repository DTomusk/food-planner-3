package user

import (
	"errors"
)

var (
	ErrEmailInUse       = errors.New("email is already in use")
	ErrUsernameInUse    = errors.New("username is already in use")
	ErrUsernameRequired = errors.New("username is required")
	ErrUsernameTooLong  = errors.New("username must be at most 50 characters long")
)
