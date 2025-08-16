package vault

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVaultData(t *testing.T) {
	v := NewVaultData()
	assert.Equal(t, VaultVersion, v.Version)
	assert.NotNil(t, v.Entries)
	assert.False(t, v.CreatedAt.IsZero())
	assert.False(t, v.UpdatedAt.IsZero())
}

func TestVaultData_SetGet(t *testing.T) {
	v := NewVaultData()
	v.Set("DB_URL", "postgres://localhost/db", []string{"production"})

	e, ok := v.Get("DB_URL")
	require.True(t, ok)
	assert.Equal(t, "DB_URL", e.Key)
	assert.Equal(t, "postgres://localhost/db", e.Value)
	assert.Equal(t, []string{"production"}, e.Tags)
}

func TestVaultData_Set_Update(t *testing.T) {
	v := NewVaultData()
	v.Set("KEY", "old_value", nil)

	e, _ := v.Get("KEY")
	createdAt := e.CreatedAt

	time.Sleep(2 * time.Millisecond)
	v.Set("KEY", "new_value", []string{"tag1"})

	e, ok := v.Get("KEY")
	require.True(t, ok)
	assert.Equal(t, "new_value", e.Value)
	assert.Equal(t, []string{"tag1"}, e.Tags)
	assert.Equal(t, createdAt, e.CreatedAt, "created_at should not change on update")
	assert.True(t, e.UpdatedAt.After(createdAt))
}

func TestVaultData_Get_Missing(t *testing.T) {
	v := NewVaultData()
	_, ok := v.Get("NONEXISTENT")
	assert.False(t, ok)
}

func TestVaultData_Delete(t *testing.T) {
	v := NewVaultData()
	v.Set("KEY", "value", nil)

	err := v.Delete("KEY")
	require.NoError(t, err)

	_, ok := v.Get("KEY")
	assert.False(t, ok)
}

func TestVaultData_Delete_Missing(t *testing.T) {
	v := NewVaultData()
	err := v.Delete("NONEXISTENT")
	assert.Error(t, err)
}

func TestVaultData_List(t *testing.T) {
	v := NewVaultData()
	v.Set("BETA", "val_b", []string{"staging"})
	v.Set("ALPHA", "val_a", []string{"production"})
	v.Set("GAMMA", "val_g", []string{"staging", "production"})

	all := v.List(nil)
	require.Len(t, all, 3)
	assert.Equal(t, "ALPHA", all[0].Key)
	assert.Equal(t, "BETA", all[1].Key)
	assert.Equal(t, "GAMMA", all[2].Key)
}

func TestVaultData_List_FilterByTag(t *testing.T) {
	v := NewVaultData()
	v.Set("DB_URL", "val_db", []string{"production"})
	v.Set("CACHE_URL", "val_cache", []string{"staging"})
	v.Set("API_KEY", "val_api", []string{"production", "staging"})

	prod := v.List([]string{"production"})
	require.Len(t, prod, 2)
	keys := []string{prod[0].Key, prod[1].Key}
	assert.Contains(t, keys, "DB_URL")
	assert.Contains(t, keys, "API_KEY")
}

func TestVaultData_Marshal_Unmarshal(t *testing.T) {
	v := NewVaultData()
	v.Set("TEST_KEY", "test_value", []string{"tag1"})

	data, err := v.Marshal()
	require.NoError(t, err)

	v2, err := UnmarshalVaultData(data)
	require.NoError(t, err)

	e, ok := v2.Get("TEST_KEY")
	require.True(t, ok)
	assert.Equal(t, "test_value", e.Value)
}

func TestWriteAtomic(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/vault.enc"

	err := writeAtomic(path, []byte("hello world"), 0600)
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(data))
}

func TestWriteAtomic_NoTmpLeft(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/vault.enc"

	err := writeAtomic(path, []byte("content"), 0600)
	require.NoError(t, err)

	_, err = os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err), "tmp file should be removed after atomic write")
}
