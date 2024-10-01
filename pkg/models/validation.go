package models

import "strings"

func ValidPassword(password string) (bool, string) {
	// check length
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	// check digits
	hasDigits := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasDigits = true
			break
		}
	}
	if !hasDigits {
		return false, "Password must contain at least one digit"
	}
	// check special chars
	hasSpecialChars := false
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",.<>?/`~"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			hasSpecialChars = true
			break
		}
	}
	if !hasSpecialChars {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}

func ValidName(name string) (bool, string) {
	// check length
	if len(name) < 2 {
		return false, "Name must be at least 2 characters"
	}
	if len(name) > 16 {
		return false, "Name must be less than 16 characters"
	}
	// check special chars
	hasSpecialChars := false
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",.<>?/`~"
	for _, char := range name {
		if strings.ContainsRune(specialChars, char) {
			hasSpecialChars = true
			break
		}
	}
	if hasSpecialChars {
		return false, "Name must not contain special characters"
	}
	return true, ""
}
