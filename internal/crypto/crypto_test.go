package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		plaintext string
	}{
		{"basic", "password123", "hello world"},
		{"empty plaintext", "password", ""},
		{"complex password", "P@ssw0rd!#$%", "secret value"},
		{"long data", "pass", `{"version":1,"entries":{"KEY":"value","ANOTHER_KEY":"another_value"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt([]byte(tt.password), []byte(tt.plaintext))
			require.NoError(t, err)
			assert.True(t, len(encrypted) > SaltLen+NonceLen)

			decrypted, err := Decrypt([]byte(tt.password), encrypted)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, string(decrypted))
		})
	}
}

func TestEncrypt_UniqueSalts(t *testing.T) {
	password := []byte("password")
	plaintext := []byte("data")

	enc1, err := Encrypt(password, plaintext)
	require.NoError(t, err)

	enc2, err := Encrypt(password, plaintext)
	require.NoError(t, err)

	// Salts must differ
	assert.False(t, bytes.Equal(enc1[:SaltLen], enc2[:SaltLen]))
}

func TestDecrypt_WrongPassword(t *testing.T) {
	encrypted, err := Encrypt([]byte("correctpassword"), []byte("secret"))
	require.NoError(t, err)

	_, err = Decrypt([]byte("wrongpassword"), encrypted)
	assert.Error(t, err)
}

func TestDecrypt_TooShort(t *testing.T) {
	_, err := Decrypt([]byte("password"), []byte("tooshort"))
	assert.Error(t, err)
}

func TestZeroBytes(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5}
	ZeroBytes(b)
	for _, v := range b {
		assert.Equal(t, byte(0), v)
	}
}
