package crypto

import (
	"fmt"
)

// Encrypt derives a key from password using Argon2id with a fresh random salt,
// then encrypts plaintext with AES-256-GCM.
// Output format: salt (16 bytes) | nonce (12 bytes) | ciphertext+tag
func Encrypt(password, plaintext []byte) ([]byte, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}

	key := DeriveKey(password, salt)
	defer ZeroBytes(key)

	encrypted, err := EncryptAES(key, plaintext)
	if err != nil {
		return nil, fmt.Errorf("encrypting: %w", err)
	}

	// Prepend salt to the encrypted payload
	result := make([]byte, SaltLen+len(encrypted))
	copy(result[:SaltLen], salt)
	copy(result[SaltLen:], encrypted)

	return result, nil
}

// Decrypt extracts the salt from the payload, re-derives the key, and decrypts.
// Input must be salt (16 bytes) | nonce (12 bytes) | ciphertext+tag
func Decrypt(password, data []byte) ([]byte, error) {
	if len(data) < SaltLen+NonceLen {
		return nil, fmt.Errorf("data too short to contain salt and nonce")
	}

	salt := data[:SaltLen]
	encrypted := data[SaltLen:]

	key := DeriveKey(password, salt)
	defer ZeroBytes(key)

	plaintext, err := DecryptAES(key, encrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypting vault: %w", err)
	}

	return plaintext, nil
}

// ZeroBytes overwrites a byte slice with zeros to remove sensitive data from memory.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
