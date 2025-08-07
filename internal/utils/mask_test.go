package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"long value", "supersecretpassword", "****word"},
		{"exactly 4 chars", "abcd", "****"},
		{"less than 4 chars", "abc", "***"},
		{"empty", "", ""},
		{"exactly 5 chars", "hello", "****ello"},
		{"last 4 shown", "postgres://localhost/mydb", "****mydb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskValue(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateSecret(t *testing.T) {
	secret, err := GenerateSecret(32)
	require.NoError(t, err)
	assert.Len(t, secret, 64) // 32 bytes = 64 hex chars
	assert.NotEmpty(t, secret)

	// Generate two secrets and ensure they differ
	secret2, err := GenerateSecret(32)
	require.NoError(t, err)
	assert.NotEqual(t, secret, secret2)
}

func TestGenerateSecret_SmallSize(t *testing.T) {
	secret, err := GenerateSecret(8)
	require.NoError(t, err)
	assert.Len(t, secret, 16) // 8 bytes = 16 hex chars
}
