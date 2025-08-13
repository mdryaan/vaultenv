package vault

import (
	"time"
)

// Entry represents a single secret stored in the vault.
type Entry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HasTag reports whether the entry has the given tag.
func (e *Entry) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// HasAnyTag reports whether the entry has any of the given tags.
func (e *Entry) HasAnyTag(tags []string) bool {
	for _, tag := range tags {
		if e.HasTag(tag) {
			return true
		}
	}
	return false
}

// NewEntry creates a new Entry with the current timestamp.
func NewEntry(key, value string, tags []string) Entry {
	now := time.Now().UTC()
	return Entry{
		Key:       key,
		Value:     value,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
