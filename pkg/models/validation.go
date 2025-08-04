package models

import (
	"fmt"
	"strings"
)

func ValidPassword(password string) (bool, string) {
	minPathLength := 8
	msg := fmt.Sprintf("Password must be at least %v characters long, contain at least one digit and one special character", &minPathLength)
	// check length
	if len(password) < minPathLength {
		return false, msg
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
		return false, msg
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
		return false, msg
	}

	return true, ""
}

func ValidName(name string) (bool, string) {
	minLen := 2
	maxLen := 16
	msg := fmt.Sprintf("Name must be at least %d and less than %d characters", minLen, maxLen)
	// check length
	if len(name) < 2 {
		return false, msg
	}
	if len(name) > 16 {
		return false, msg
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
