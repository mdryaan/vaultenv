package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// MaskValue returns **** + last 4 characters of value.
// If the value is 4 chars or fewer, returns all asterisks.
func MaskValue(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return "****" + value[len(value)-4:]
}

// GenerateSecret generates a cryptographically random hex string of the given byte length.
func GenerateSecret(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
