package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	e := NewEntry("DB_URL", "postgres://localhost/db", []string{"production"})
	assert.Equal(t, "DB_URL", e.Key)
	assert.Equal(t, "postgres://localhost/db", e.Value)
	assert.Equal(t, []string{"production"}, e.Tags)
	assert.False(t, e.CreatedAt.IsZero())
	assert.False(t, e.UpdatedAt.IsZero())
	assert.Equal(t, e.CreatedAt, e.UpdatedAt)
}

func TestEntry_HasTag(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		tag      string
		expected bool
	}{
		{"has tag", []string{"production", "backend"}, "production", true},
		{"no tag", []string{"production"}, "staging", false},
		{"empty tags", []string{}, "any", false},
		{"nil tags", nil, "any", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{Tags: tt.tags}
			assert.Equal(t, tt.expected, e.HasTag(tt.tag))
		})
	}
}

func TestEntry_HasAnyTag(t *testing.T) {
	tests := []struct {
		name     string
		entryTags []string
		searchTags []string
		expected bool
	}{
		{"one match", []string{"production", "backend"}, []string{"staging", "production"}, true},
		{"no match", []string{"production"}, []string{"staging", "dev"}, false},
		{"empty search", []string{"production"}, []string{}, false},
		{"both empty", []string{}, []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{Tags: tt.entryTags}
			assert.Equal(t, tt.expected, e.HasAnyTag(tt.searchTags))
		})
	}
}
