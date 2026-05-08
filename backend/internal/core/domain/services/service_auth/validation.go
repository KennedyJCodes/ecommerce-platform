// Package service_auth provides implementations of domain validators.
// It ensures that user-provided data meets the security and business requirements before it reaches the core services.
package service_auth

import (
	"fmt"
	"unicode"
)

const (
	// minUserNameLength defines the lower bound for username characters.
	minUserNameLength = 5

	// minPasswordLength defines the minimum entropy required for passwords.
	minPasswordLength = 10
)

// UserNameValidator implements the input.Validator interface for username strings.
type UserNameValidator struct{}

// Validate checks the username against length and nullability constraints.
// Parameters:
//   - input: expected to be a string.
// Returns an error with a descriptive message if validation fails.
func (c *UserNameValidator) Validate(input interface{}) error {
	username := input.(string)
	if username == "" {
		return fmt.Errorf("you cannot enter empty fields")
	}

	if len(username) < minUserNameLength {
		return fmt.Errorf("you cannot enter a name that is less than 5 characters")
	}

	return nil
}

// PasswordValidator implements complex rule validation for user passwords.
// It checks for length, casing, digits, and special characters.
type PasswordValidator struct{}

// Validate ensures the password meets the minimum security policy:
//  - Non-empty and at least 10 characters.
//  - At least one uppercase letter.
//  - At least one numeric digit.
//  - At least one special character (symbol or punctuation).
func (p *PasswordValidator) Validate(input interface{}) error {
	password := input.(string)
	if password == "" {
		return fmt.Errorf("you cannot enter empty fields")
	}
	if len(password) < minPasswordLength {
		return fmt.Errorf("you cannot enter a password that is less than 10 characters")
	}

	var hasUppercase bool
	var hasDigit bool
	var hasSpecialCharacter bool
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecialCharacter = true
		}

		// Early exit if all conditions are met
		if hasUppercase && hasDigit && hasSpecialCharacter {
			break
		}
	}

	if !hasUppercase {
		return fmt.Errorf("the password must have at least one uppercase letter")
	}

	if !hasDigit {
		return fmt.Errorf("the password must have at least one number")
	}

	if !hasSpecialCharacter {
		return fmt.Errorf("the password must have some special character")
	}
	return nil
}
