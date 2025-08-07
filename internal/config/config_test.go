package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve_DefaultPath(t *testing.T) {
	os.Unsetenv(EnvVaultPath)
	os.Unsetenv(EnvVaultPassword)

	cfg, err := Resolve("")
	require.NoError(t, err)
	assert.NotEmpty(t, cfg.VaultPath)
	assert.Contains(t, cfg.VaultPath, "vault.enc")
	assert.Empty(t, cfg.Password)
}

func TestResolve_FlagPath(t *testing.T) {
	os.Unsetenv(EnvVaultPath)
	cfg, err := Resolve("/tmp/custom.enc")
	require.NoError(t, err)
	assert.Equal(t, "/tmp/custom.enc", cfg.VaultPath)
}

func TestResolve_EnvPath(t *testing.T) {
	os.Setenv(EnvVaultPath, "/tmp/env.enc")
	defer os.Unsetenv(EnvVaultPath)

	cfg, err := Resolve("")
	require.NoError(t, err)
	assert.Equal(t, "/tmp/env.enc", cfg.VaultPath)
}

func TestResolve_FlagOverridesEnv(t *testing.T) {
	os.Setenv(EnvVaultPath, "/tmp/env.enc")
	defer os.Unsetenv(EnvVaultPath)

	cfg, err := Resolve("/tmp/flag.enc")
	require.NoError(t, err)
	assert.Equal(t, "/tmp/flag.enc", cfg.VaultPath)
}

func TestResolve_Password(t *testing.T) {
	os.Unsetenv(EnvVaultPath)
	os.Setenv(EnvVaultPassword, "supersecret")
	defer os.Unsetenv(EnvVaultPassword)

	cfg, err := Resolve("")
	require.NoError(t, err)
	assert.Equal(t, "supersecret", cfg.Password)
}
