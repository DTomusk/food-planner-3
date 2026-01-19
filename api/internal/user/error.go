package users

import (
	"errors"
)

var (
	ErrEmailInUse = errors.New("email is already in use")
)
