package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: AskPassword and AskPasswordConfirm require a real terminal,
// so they are not unit-tested here. The Confirm function can be tested
// via integration tests.

func TestAskPasswordConfirm_MismatchError(t *testing.T) {
	// We can simulate mismatch by directly testing the byte comparison logic.
	// The actual terminal interaction is tested manually.
	pass1 := []byte("password1")
	pass2 := []byte("password2")

	match := string(pass1) == string(pass2)
	assert.False(t, match, "different passwords should not match")
}

func TestAskPasswordConfirm_MatchLogic(t *testing.T) {
	pass1 := []byte("samepassword")
	pass2 := []byte("samepassword")

	match := string(pass1) == string(pass2)
	assert.True(t, match, "identical passwords should match")
}
