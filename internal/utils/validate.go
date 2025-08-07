package utils

import (
	"fmt"
	"regexp"
)

var keyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// ValidateKey checks that a secret key matches [A-Z][A-Z0-9_]*
func ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	if !keyPattern.MatchString(key) {
		return fmt.Errorf("invalid key %q: must match [A-Z][A-Z0-9_]* (uppercase letters, digits, underscores, starting with a letter)", key)
	}
	return nil
}
