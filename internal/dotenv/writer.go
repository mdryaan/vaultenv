package dotenv

import (
	"fmt"
	"io"
	"strings"
)

// Write writes key=value pairs to w in .env file format.
// Values containing spaces, special chars, or quotes are double-quoted.
func Write(w io.Writer, entries []Entry) error {
	for _, e := range entries {
		value := e.Value
		if needsQuoting(value) {
			value = `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, value); err != nil {
			return err
		}
	}
	return nil
}

// WriteString formats entries as a .env string.
func WriteString(entries []Entry) (string, error) {
	var sb strings.Builder
	if err := Write(&sb, entries); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func needsQuoting(value string) bool {
	if value == "" {
		return false
	}
	for _, ch := range value {
		if ch == ' ' || ch == '\t' || ch == '"' || ch == '\'' || ch == '#' || ch == '$' || ch == '\\' {
			return true
		}
	}
	return false
}
