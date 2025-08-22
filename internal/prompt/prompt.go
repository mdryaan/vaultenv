package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// AskPassword prompts the user for a password without echo.
func AskPassword(msg string) ([]byte, error) {
	return ReadPassword(msg)
}

// AskPasswordConfirm prompts for a password twice and returns an error if they differ.
func AskPasswordConfirm(msg, confirmMsg string) ([]byte, error) {
	pass1, err := ReadPassword(msg)
	if err != nil {
		return nil, err
	}

	pass2, err := ReadPassword(confirmMsg)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Zero the confirmation copy
		for i := range pass2 {
			pass2[i] = 0
		}
	}()

	if !bytes.Equal(pass1, pass2) {
		// Zero pass1 on mismatch
		for i := range pass1 {
			pass1[i] = 0
		}
		return nil, fmt.Errorf("passwords do not match")
	}

	return pass1, nil
}

// Confirm asks the user a yes/no question and returns true for "y" or "yes".
func Confirm(msg string) (bool, error) {
	fmt.Fprint(os.Stderr, msg+" [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false, scanner.Err()
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

// AskValue prompts the user for a secret value without echo.
func AskValue(key string) (string, error) {
	raw, err := ReadPassword(fmt.Sprintf("Enter value for %s: ", key))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
