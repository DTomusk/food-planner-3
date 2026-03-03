package user

import (
	"errors"
)

var (
	ErrEmailInUse    = errors.New("email is already in use")
	ErrUsernameInUse = errors.New("username is already in use")
)
