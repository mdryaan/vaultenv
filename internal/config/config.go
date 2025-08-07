package config

import (
	"os"

	"github.com/mdryaan/vaultenv/internal/utils"
)

const (
	EnvVaultPath     = "VAULTENV_PATH"
	EnvVaultPassword = "VAULTENV_PASSWORD"
)

// Config holds the resolved runtime configuration.
type Config struct {
	VaultPath string
	Password  string
}

// Resolve returns a Config with vault path and optional password resolved
// from flags, environment variables, and defaults (in that priority order).
func Resolve(flagPath string) (*Config, error) {
	cfg := &Config{}

	switch {
	case flagPath != "":
		cfg.VaultPath = flagPath
	case os.Getenv(EnvVaultPath) != "":
		cfg.VaultPath = os.Getenv(EnvVaultPath)
	default:
		p, err := utils.DefaultVaultPath()
		if err != nil {
			return nil, err
		}
		cfg.VaultPath = p
	}

	cfg.Password = os.Getenv(EnvVaultPassword)
	return cfg, nil
}
