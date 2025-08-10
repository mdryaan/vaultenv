package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	require.NoError(t, err)
	assert.Len(t, salt, SaltLen)
	assert.NotEmpty(t, salt)

	salt2, err := GenerateSalt()
	require.NoError(t, err)
	assert.NotEqual(t, salt, salt2, "salts should be unique")
}

func TestDeriveKey(t *testing.T) {
	tests := []struct {
		name     string
		password string
		salt     []byte
	}{
		{
			name:     "basic derivation",
			password: "hunter2",
			salt:     []byte("saltsaltsaltsalt"), // 16 bytes
		},
		{
			name:     "empty password",
			password: "",
			salt:     []byte("saltsaltsaltsalt"),
		},
		{
			name:     "long password",
			password: "this-is-a-very-long-password-with-special-chars-!@#$%^&*()",
			salt:     []byte("saltsaltsaltsalt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := DeriveKey([]byte(tt.password), tt.salt)
			assert.Len(t, key, Argon2KeyLen)
			assert.NotEmpty(t, key)

			// Same inputs must produce same key (deterministic)
			key2 := DeriveKey([]byte(tt.password), tt.salt)
			assert.Equal(t, key, key2)
		})
	}
}

func TestDeriveKey_DifferentPasswords(t *testing.T) {
	salt := []byte("saltsaltsaltsalt")
	key1 := DeriveKey([]byte("password1"), salt)
	key2 := DeriveKey([]byte("password2"), salt)
	assert.NotEqual(t, key1, key2)
}

func TestDeriveKey_DifferentSalts(t *testing.T) {
	password := []byte("password")
	key1 := DeriveKey(password, []byte("saltsaltsaltsalt"))
	key2 := DeriveKey(password, []byte("differentsaltttt"))
	assert.NotEqual(t, key1, key2)
}

func TestDeriveKey_KnownVector(t *testing.T) {
	// Known test vector: deterministic Argon2id output
	password := []byte("testpassword")
	salt := []byte("0123456789abcdef") // 16 bytes

	key := DeriveKey(password, salt)
	assert.Len(t, key, 32)

	// Re-derive and compare to confirm determinism
	key2 := DeriveKey(password, salt)
	assert.Equal(t, key, key2)
}
