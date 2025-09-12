package dotenv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteString(t *testing.T) {
	tests := []struct {
		name     string
		entries  []Entry
		expected string
	}{
		{
			name:     "simple",
			entries:  []Entry{{Key: "KEY", Value: "value"}},
			expected: "KEY=value\n",
		},
		{
			name:     "value with space",
			entries:  []Entry{{Key: "MSG", Value: "hello world"}},
			expected: `MSG="hello world"` + "\n",
		},
		{
			name:     "empty value",
			entries:  []Entry{{Key: "EMPTY", Value: ""}},
			expected: "EMPTY=\n",
		},
		{
			name: "multiple entries",
			entries: []Entry{
				{Key: "A", Value: "1"},
				{Key: "B", Value: "2"},
			},
			expected: "A=1\nB=2\n",
		},
		{
			name:     "value with hash",
			entries:  []Entry{{Key: "KEY", Value: "val#ue"}},
			expected: `KEY="val#ue"` + "\n",
		},
		{
			name:     "value with dollar sign",
			entries:  []Entry{{Key: "KEY", Value: "val$ue"}},
			expected: `KEY="val$ue"` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WriteString(tt.entries)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWrite_RoundTrip(t *testing.T) {
	original := []Entry{
		{Key: "DB_URL", Value: "postgres://localhost/db"},
		{Key: "API_KEY", Value: "secret123"},
		{Key: "MSG", Value: "hello world"},
	}

	written, err := WriteString(original)
	require.NoError(t, err)

	parsed, err := ParseString(written)
	require.NoError(t, err)
	assert.Equal(t, original, parsed)
}

func TestNeedsQuoting(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"simple", false},
		{"with space", true},
		{"with#hash", true},
		{"with$dollar", true},
		{`with"quote`, true},
		{"with'quote", true},
		{"", false},
		{"no-special-chars", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.expected, needsQuoting(tt.value))
		})
	}
}

func TestWrite_EmptyEntries(t *testing.T) {
	result, err := WriteString([]Entry{})
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestWrite_IOWriter(t *testing.T) {
	var sb strings.Builder
	err := Write(&sb, []Entry{{Key: "X", Value: "y"}})
	require.NoError(t, err)
	assert.Equal(t, "X=y\n", sb.String())
}
