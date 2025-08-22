package prompt

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ReadPassword reads a password from the terminal without echoing.
// It writes the prompt to stderr and reads from the controlling terminal.
func ReadPassword(prompt string) ([]byte, error) {
	fmt.Fprint(os.Stderr, prompt)
	fd := int(os.Stdin.Fd())
	password, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("reading password: %w", err)
	}
	return password, nil
}
