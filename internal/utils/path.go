package utils

import (
	"os"
	"path/filepath"
)

const (
	defaultDir  = ".vaultenv"
	defaultFile = "vault.enc"
)

// DefaultVaultPath returns the default vault file path: ~/.vaultenv/vault.enc
func DefaultVaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultDir, defaultFile), nil
}

// DefaultVaultDir returns the default vault directory: ~/.vaultenv
func DefaultVaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultDir), nil
}

// EnsureDir creates the directory at path if it does not exist.
func EnsureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0700)
}
