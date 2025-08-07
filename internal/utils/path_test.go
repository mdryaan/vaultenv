package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultVaultPath(t *testing.T) {
	path, err := DefaultVaultPath()
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(path, "vault.enc"))
	assert.True(t, strings.Contains(path, ".vaultenv"))
}

func TestDefaultVaultDir(t *testing.T) {
	dir, err := DefaultVaultDir()
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(dir, ".vaultenv"))
}

func TestEnsureDir(t *testing.T) {
	t.TempDir()
	path := t.TempDir() + "/subdir/vault.enc"
	err := EnsureDir(path)
	require.NoError(t, err)
}
