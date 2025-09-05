package output

import (
	"bytes"
	"testing"
	"time"

	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeEntry(key, value string, tags []string) vault.Entry {
	now := time.Now().UTC()
	return vault.Entry{
		Key:       key,
		Value:     value,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestJSONFormatter_ShowValues(t *testing.T) {
	entries := []vault.Entry{
		makeEntry("API_KEY", "supersecret", []string{"production"}),
	}

	f := &JSONFormatter{}
	var buf bytes.Buffer
	err := f.WriteEntries(&buf, entries, true)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "supersecret")
	assert.Contains(t, buf.String(), "API_KEY")
	assert.Contains(t, buf.String(), "production")
}

func TestJSONFormatter_MaskedValues(t *testing.T) {
	entries := []vault.Entry{
		makeEntry("API_KEY", "supersecret", nil),
	}

	f := &JSONFormatter{}
	var buf bytes.Buffer
	err := f.WriteEntries(&buf, entries, false)
	require.NoError(t, err)
	assert.NotContains(t, buf.String(), "supersecret")
	assert.Contains(t, buf.String(), "****")
}

func TestJSONFormatter_EmptyEntries(t *testing.T) {
	f := &JSONFormatter{}
	var buf bytes.Buffer
	err := f.WriteEntries(&buf, []vault.Entry{}, true)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "[]")
}

func TestTableFormatter_ShowValues(t *testing.T) {
	entries := []vault.Entry{
		makeEntry("DB_URL", "postgres://localhost/db", []string{"prod"}),
	}

	f := &TableFormatter{}
	var buf bytes.Buffer
	err := f.WriteEntries(&buf, entries, true)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "DB_URL")
	assert.Contains(t, buf.String(), "postgres://localhost/db")
}

func TestTableFormatter_MaskedValues(t *testing.T) {
	entries := []vault.Entry{
		makeEntry("DB_URL", "postgres://localhost/db", nil),
	}

	f := &TableFormatter{}
	var buf bytes.Buffer
	err := f.WriteEntries(&buf, entries, false)
	require.NoError(t, err)
	assert.NotContains(t, buf.String(), "postgres://localhost/db")
	assert.Contains(t, buf.String(), "****")
}

func TestGet_JSONFormatter(t *testing.T) {
	f := Get(FormatJSON)
	_, ok := f.(*JSONFormatter)
	assert.True(t, ok)
}

func TestGet_TableFormatter(t *testing.T) {
	f := Get(FormatTable)
	_, ok := f.(*TableFormatter)
	assert.True(t, ok)
}

func TestGet_UnknownDefaultsToTable(t *testing.T) {
	f := Get("unknown")
	_, ok := f.(*TableFormatter)
	assert.True(t, ok)
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"supersecret", "****cret"},
		{"abcd", "****"},
		{"abc", "****"},
		{"", "****"},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			result := maskValue(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
