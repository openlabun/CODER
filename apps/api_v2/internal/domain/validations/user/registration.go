package user_validations

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateRegistrationData(email, name, password string) error {
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)

	// Check email format is "{text}@{text}.{text}"
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	// Check password strength
	// At least 8 characters, one uppercase, one lowercase, one number, one special character (!, @, #, $, _, -, .)
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$_.\-]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character (!, @, #, $, _, -, .)")
	}

	// Check name is not empty
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	return nil

}
