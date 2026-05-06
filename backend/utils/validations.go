package utils

import (
	"errors"
	"regexp"
	"strings"
)

var (
	usernameRe = regexp.MustCompile(`^[a-z0-9_.-]{3,24}$`)
	emailRe    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

func ValidateUsername(username string) error {
	if !usernameRe.MatchString(strings.ToLower(username)) {
		return errors.New("username debe tener 3-24 caracteres (a-z, 0-9, ._-)")
	}
	return nil
}

func ValidateEmail(email string) error {
	if !emailRe.MatchString(email) {
		return errors.New("email no tiene formato válido")
	}
	return nil
}
