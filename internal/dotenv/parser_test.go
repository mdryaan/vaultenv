package dotenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Entry
		wantErr  bool
	}{
		{
			name:  "simple key=value",
			input: "DB_URL=postgres://localhost/db",
			expected: []Entry{
				{Key: "DB_URL", Value: "postgres://localhost/db"},
			},
		},
		{
			name:  "double quoted value",
			input: `API_KEY="mysecretkey"`,
			expected: []Entry{
				{Key: "API_KEY", Value: "mysecretkey"},
			},
		},
		{
			name:  "single quoted value",
			input: `SECRET='my secret'`,
			expected: []Entry{
				{Key: "SECRET", Value: "my secret"},
			},
		},
		{
			name:  "skip comment line",
			input: "# This is a comment\nDB_HOST=localhost",
			expected: []Entry{
				{Key: "DB_HOST", Value: "localhost"},
			},
		},
		{
			name:     "skip blank lines",
			input:    "\n\nDB_HOST=localhost\n\n",
			expected: []Entry{{Key: "DB_HOST", Value: "localhost"}},
		},
		{
			name: "multiple entries",
			input: `DB_URL=postgres://localhost/db
API_KEY=secret123
PORT=5432`,
			expected: []Entry{
				{Key: "DB_URL", Value: "postgres://localhost/db"},
				{Key: "API_KEY", Value: "secret123"},
				{Key: "PORT", Value: "5432"},
			},
		},
		{
			name:  "empty value",
			input: "EMPTY=",
			expected: []Entry{
				{Key: "EMPTY", Value: ""},
			},
		},
		{
			name:  "inline comment stripped",
			input: "PORT=5432 # default port",
			expected: []Entry{
				{Key: "PORT", Value: "5432"},
			},
		},
		{
			name:    "missing equals",
			input:   "NOEQUALS",
			wantErr: true,
		},
		{
			name:    "empty key",
			input:   "=value",
			wantErr: true,
		},
		{
			name:    "unterminated double quote",
			input:   `KEY="unterminated`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := ParseString(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, entries)
		})
	}
}
