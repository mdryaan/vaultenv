package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

const VaultVersion = 1

// VaultData is the in-memory representation of the vault's contents.
type VaultData struct {
	Version   int              `json:"version"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Entries   map[string]Entry `json:"entries"`
}

// NewVaultData returns an empty vault with the current timestamp.
func NewVaultData() *VaultData {
	now := time.Now().UTC()
	return &VaultData{
		Version:   VaultVersion,
		CreatedAt: now,
		UpdatedAt: now,
		Entries:   make(map[string]Entry),
	}
}

// Set adds or updates a secret entry in the vault.
func (v *VaultData) Set(key, value string, tags []string) {
	existing, ok := v.Entries[key]
	now := time.Now().UTC()

	if ok {
		existing.Value = value
		existing.Tags = tags
		existing.UpdatedAt = now
		v.Entries[key] = existing
	} else {
		v.Entries[key] = Entry{
			Key:       key,
			Value:     value,
			Tags:      tags,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	v.UpdatedAt = now
}

// Get retrieves a secret by key. Returns the entry and a boolean indicating existence.
func (v *VaultData) Get(key string) (Entry, bool) {
	e, ok := v.Entries[key]
	return e, ok
}

// Delete removes a secret by key. Returns error if the key does not exist.
func (v *VaultData) Delete(key string) error {
	if _, ok := v.Entries[key]; !ok {
		return fmt.Errorf("key %q not found", key)
	}
	delete(v.Entries, key)
	v.UpdatedAt = time.Now().UTC()
	return nil
}

// List returns all entries, optionally filtered by tags (empty filter = all).
func (v *VaultData) List(tags []string) []Entry {
	entries := make([]Entry, 0, len(v.Entries))
	for _, e := range v.Entries {
		if len(tags) == 0 || e.HasAnyTag(tags) {
			entries = append(entries, e)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// Marshal serializes the vault data to JSON.
func (v *VaultData) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalVaultData deserializes JSON into a VaultData struct.
func UnmarshalVaultData(data []byte) (*VaultData, error) {
	var vd VaultData
	if err := json.Unmarshal(data, &vd); err != nil {
		return nil, fmt.Errorf("unmarshaling vault: %w", err)
	}
	if vd.Entries == nil {
		vd.Entries = make(map[string]Entry)
	}
	return &vd, nil
}

// writeAtomic writes data to path atomically: write to path+".tmp" then os.Rename.
func writeAtomic(path string, data []byte, perm os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("renaming temp file: %w", err)
	}
	return nil
}
