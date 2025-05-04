package validation

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrUsernameTooShort   = errors.New("username must be at least 5 characters long")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrPasswordNoUpper    = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower    = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoSpecial  = errors.New("password must contain at least one special character")
)

type UserValidator struct {
	minUsernameLength int
	minPasswordLength int
	emailRegex        *regexp.Regexp
	specialChars      string
	version           *regexp.Regexp
}

func NewUserValidator() *UserValidator {
	return &UserValidator{
		minUsernameLength: 5,
		minPasswordLength: 8,
		emailRegex:        regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		specialChars:      "!@#$%^&*()-_=+[]{}|;:',.<>/?",
		version:           regexp.MustCompile(`0-9`),
	}
}

func (v *UserValidator) ValidateRegistration(username, email, password string) error {
	if err := v.validateUsername(username); err != nil {
		return err
	}

	if err := v.ValidateEmail(email); err != nil {
		return err
	}

	return v.validatePassword(password)
}

func (v *UserValidator) ValidateLogin(email, password string) error {
	if err := v.ValidateEmail(email); err != nil {
		return err
	}

	return v.validatePassword(password)
}

func (v *UserValidator) validateUsername(username string) error {
	if len(username) < v.minUsernameLength {
		return ErrUsernameTooShort
	}
	return nil
}

func (v *UserValidator) ValidateEmail(email string) error {
	if !v.emailRegex.MatchString(email) {
		return ErrInvalidEmailFormat
	}
	return nil
}

func (v *UserValidator) validatePassword(password string) error {
	if len(password) < v.minPasswordLength {
		return ErrPasswordTooShort
	}

	if !containsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return ErrPasswordNoUpper
	}

	if !containsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return ErrPasswordNoLower
	}

	if !containsAny(password, v.specialChars) {
		return ErrPasswordNoSpecial
	}

	return nil
}

func containsAny(s, chars string) bool {
	for _, c := range chars {
		if strings.ContainsRune(s, c) {
			return true
		}
	}
	return false
}
