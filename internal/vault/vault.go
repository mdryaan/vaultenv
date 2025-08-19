package vault

import (
	"fmt"
	"os"

	"github.com/mdryaan/vaultenv/internal/crypto"
	"github.com/mdryaan/vaultenv/internal/utils"
)

// Vault provides high-level operations over an encrypted vault file.
type Vault struct {
	path     string
	password []byte
	data     *VaultData
}

// Open reads and decrypts the vault at path using password.
func Open(path string, password []byte) (*Vault, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading vault file: %w", err)
	}

	plaintext, err := crypto.Decrypt(password, raw)
	if err != nil {
		return nil, fmt.Errorf("wrong password or corrupted vault: %w", err)
	}

	data, err := UnmarshalVaultData(plaintext)
	if err != nil {
		return nil, fmt.Errorf("parsing vault: %w", err)
	}

	return &Vault{path: path, password: password, data: data}, nil
}

// Create initializes a new empty vault and saves it to disk.
func Create(path string, password []byte) (*Vault, error) {
	if err := utils.EnsureDir(path); err != nil {
		return nil, fmt.Errorf("creating vault directory: %w", err)
	}

	v := &Vault{
		path:     path,
		password: password,
		data:     NewVaultData(),
	}

	if err := v.Save(); err != nil {
		return nil, err
	}
	return v, nil
}

// Save marshals and encrypts the vault, then writes it atomically.
func (v *Vault) Save() error {
	plaintext, err := v.data.Marshal()
	if err != nil {
		return fmt.Errorf("marshaling vault: %w", err)
	}

	ciphertext, err := crypto.Encrypt(v.password, plaintext)
	if err != nil {
		return fmt.Errorf("encrypting vault: %w", err)
	}

	if err := writeAtomic(v.path, ciphertext, 0600); err != nil {
		return fmt.Errorf("writing vault: %w", err)
	}
	return nil
}

// Set adds or updates a secret.
func (v *Vault) Set(key, value string, tags []string) error {
	v.data.Set(key, value, tags)
	return v.Save()
}

// Get retrieves a secret by key.
func (v *Vault) Get(key string) (Entry, error) {
	e, ok := v.data.Get(key)
	if !ok {
		return Entry{}, fmt.Errorf("key %q not found", key)
	}
	return e, nil
}

// Delete removes a secret by key.
func (v *Vault) Delete(key string) error {
	if err := v.data.Delete(key); err != nil {
		return err
	}
	return v.Save()
}

// List returns all entries optionally filtered by tags.
func (v *Vault) List(tags []string) []Entry {
	return v.data.List(tags)
}

// Rotate replaces the master password with newPassword by re-encrypting the vault.
func (v *Vault) Rotate(newPassword []byte) error {
	oldPassword := v.password
	defer crypto.ZeroBytes(oldPassword)

	v.password = newPassword
	if err := v.Save(); err != nil {
		v.password = oldPassword
		return fmt.Errorf("rotating password: %w", err)
	}
	return nil
}

// Data returns the underlying vault data (read-only for export).
func (v *Vault) Data() *VaultData {
	return v.data
}

// Path returns the vault file path.
func (v *Vault) Path() string {
	return v.path
}
