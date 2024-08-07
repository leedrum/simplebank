package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

const (
	MinUsernameLength = 3
	MaxUsernameLength = 25
	MinPasswordLength = 6
	MaxPasswordLength = 50
	MinEmailLength    = 3
	MaxEmailLength    = 200
	MinFullNameLength = 3
	MaxFullNameLength = 100
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(str string, minLength int, maxLenth int) error {
	n := len(str)
	if n < minLength || n > maxLenth {
		return fmt.Errorf("string length must be between %d and %d", minLength, maxLenth)
	}

	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, MinUsernameLength, MaxUsernameLength); err != nil {
		return err
	}

	if !isValidUsername(username) {
		return fmt.Errorf("username can only contain lowercase letters, digits and underscores")
	}

	return nil
}

func ValidateFullName(full_name string) error {
	if err := ValidateString(full_name, MinFullNameLength, MaxFullNameLength); err != nil {
		return err
	}

	if !isValidFullName(full_name) {
		return fmt.Errorf("full name can only contain letters and spaces")
	}

	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, MinPasswordLength, MaxPasswordLength)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, MinEmailLength, MaxEmailLength); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email address")
	}

	return nil
}
