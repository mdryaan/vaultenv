package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	Argon2Time    = 3
	Argon2Memory  = 64 * 1024 // 64 MB
	Argon2Threads = 4
	Argon2KeyLen  = 32
	SaltLen       = 16
)

// GenerateSalt creates a cryptographically random 16-byte salt.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}
	return salt, nil
}

// DeriveKey derives a 32-byte AES key from password and salt using Argon2id.
// The returned key slice must be zeroed by the caller after use.
func DeriveKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)
}
