package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"valid simple", "DATABASE_URL", false},
		{"valid single", "A", false},
		{"valid with digits", "API2_KEY", false},
		{"valid underscore", "MY_VAR_123", false},
		{"empty", "", true},
		{"lowercase", "database_url", true},
		{"starts with digit", "1_KEY", true},
		{"starts with underscore", "_KEY", true},
		{"has hyphen", "MY-KEY", true},
		{"has space", "MY KEY", true},
		{"mixed case", "MyKey", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKey(tt.key)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
