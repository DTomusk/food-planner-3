package auth

import (
	"net/mail"
	"strings"
)

func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return ErrInvalidEmail
	}

	domain := parts[1]

	if !strings.Contains(domain, ".") {
		return ErrInvalidEmail
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	if len(password) > 64 {
		return ErrInvalidPassword
	}
	// TODO: add checks for characters that aren't allowed
	return nil
}
