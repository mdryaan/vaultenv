package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate_Open(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	password := []byte("testpassword")

	v, err := Create(path, password)
	require.NoError(t, err)
	assert.NotNil(t, v)

	_, err = os.Stat(path)
	require.NoError(t, err, "vault file should exist")

	v2, err := Open(path, password)
	require.NoError(t, err)
	assert.NotNil(t, v2)
}

func TestOpen_WrongPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")

	_, err := Create(path, []byte("correctpassword"))
	require.NoError(t, err)

	_, err = Open(path, []byte("wrongpassword"))
	assert.Error(t, err)
}

func TestOpen_NonExistent(t *testing.T) {
	_, err := Open("/nonexistent/path/vault.enc", []byte("password"))
	assert.Error(t, err)
}

func TestVault_SetGet(t *testing.T) {
	v := newTestVault(t)

	err := v.Set("DB_URL", "postgres://localhost/mydb", []string{"production"})
	require.NoError(t, err)

	e, err := v.Get("DB_URL")
	require.NoError(t, err)
	assert.Equal(t, "postgres://localhost/mydb", e.Value)
	assert.Equal(t, []string{"production"}, e.Tags)
}

func TestVault_Get_Missing(t *testing.T) {
	v := newTestVault(t)
	_, err := v.Get("NONEXISTENT")
	assert.Error(t, err)
}

func TestVault_Delete(t *testing.T) {
	v := newTestVault(t)

	err := v.Set("KEY", "value", nil)
	require.NoError(t, err)

	err = v.Delete("KEY")
	require.NoError(t, err)

	_, err = v.Get("KEY")
	assert.Error(t, err)
}

func TestVault_List(t *testing.T) {
	v := newTestVault(t)

	require.NoError(t, v.Set("BETA", "b", nil))
	require.NoError(t, v.Set("ALPHA", "a", []string{"prod"}))

	entries := v.List(nil)
	require.Len(t, entries, 2)
	assert.Equal(t, "ALPHA", entries[0].Key)
	assert.Equal(t, "BETA", entries[1].Key)
}

func TestVault_Persist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	password := []byte("password")

	v, err := Create(path, password)
	require.NoError(t, err)
	require.NoError(t, v.Set("PERSISTED_KEY", "persisted_value", nil))

	v2, err := Open(path, password)
	require.NoError(t, err)

	e, err := v2.Get("PERSISTED_KEY")
	require.NoError(t, err)
	assert.Equal(t, "persisted_value", e.Value)
}

func TestVault_Rotate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	oldPass := []byte("oldpassword")
	newPass := []byte("newpassword")

	v, err := Create(path, oldPass)
	require.NoError(t, err)
	require.NoError(t, v.Set("KEY", "value", nil))

	err = v.Rotate(newPass)
	require.NoError(t, err)

	// Should open with new password
	v2, err := Open(path, newPass)
	require.NoError(t, err)
	e, err := v2.Get("KEY")
	require.NoError(t, err)
	assert.Equal(t, "value", e.Value)

	// Should NOT open with old password
	_, err = Open(path, oldPass)
	assert.Error(t, err)
}

func newTestVault(t *testing.T) *Vault {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.enc")
	v, err := Create(path, []byte("testpass"))
	require.NoError(t, err)
	return v
}
