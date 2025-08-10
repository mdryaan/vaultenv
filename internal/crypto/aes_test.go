package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptAES(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
	}{
		{"short message", "hello"},
		{"empty message", ""},
		{"long message", "this is a longer message with various characters: !@#$%^&*()"},
		{"unicode", "secret: 日本語テスト"},
		{"json payload", `{"key":"value","nested":{"a":1}}`},
	}

	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := EncryptAES(key, []byte(tt.plaintext))
			require.NoError(t, err)
			assert.NotNil(t, ciphertext)

			plaintext, err := DecryptAES(key, ciphertext)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, string(plaintext))
		})
	}
}

func TestEncryptAES_UniqueNonces(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("same message")

	ct1, err := EncryptAES(key, plaintext)
	require.NoError(t, err)

	ct2, err := EncryptAES(key, plaintext)
	require.NoError(t, err)

	// Nonces should be different
	assert.False(t, bytes.Equal(ct1[:NonceLen], ct2[:NonceLen]), "nonces must be unique")
	// Ciphertexts should differ due to different nonces
	assert.False(t, bytes.Equal(ct1, ct2))
}

func TestDecryptAES_WrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	key2[0] = 0xFF

	ciphertext, err := EncryptAES(key1, []byte("secret"))
	require.NoError(t, err)

	_, err = DecryptAES(key2, ciphertext)
	assert.Error(t, err)
}

func TestDecryptAES_TamperedCiphertext(t *testing.T) {
	key := make([]byte, 32)
	ciphertext, err := EncryptAES(key, []byte("secret data"))
	require.NoError(t, err)

	// Tamper with the ciphertext
	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err = DecryptAES(key, ciphertext)
	assert.Error(t, err)
}

func TestDecryptAES_TooShort(t *testing.T) {
	key := make([]byte, 32)
	_, err := DecryptAES(key, []byte("short"))
	assert.Error(t, err)
}

func TestEncryptAES_InvalidKeyLength(t *testing.T) {
	key := make([]byte, 10) // invalid key length
	_, err := EncryptAES(key, []byte("data"))
	assert.Error(t, err)
}
