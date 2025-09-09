package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a parsed key-value pair from a .env file.
type Entry struct {
	Key   string
	Value string
}

// Parse reads a .env file and returns all key=value pairs.
// It handles:
//   - Comments (# ...)
//   - Blank lines
//   - Quoted values ("..." and '...')
//   - Inline comments after unquoted values
func Parse(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("line %d: missing '=' in %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		raw := line[idx+1:]

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		value, err := parseValue(raw)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		entries = append(entries, Entry{Key: key, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading env file: %w", err)
	}

	return entries, nil
}

// ParseString parses a .env formatted string.
func ParseString(s string) ([]Entry, error) {
	return Parse(strings.NewReader(s))
}

func parseValue(raw string) (string, error) {
	raw = strings.TrimLeft(raw, " \t")

	if len(raw) == 0 {
		return "", nil
	}

	switch raw[0] {
	case '"':
		return parseQuoted(raw, '"')
	case '\'':
		return parseQuoted(raw, '\'')
	default:
		// Strip inline comment
		if idx := strings.Index(raw, " #"); idx >= 0 {
			raw = raw[:idx]
		}
		return strings.TrimRight(raw, " \t"), nil
	}
}

func parseQuoted(raw string, quote byte) (string, error) {
	if len(raw) < 2 {
		return "", fmt.Errorf("unterminated quoted value")
	}

	end := strings.IndexByte(raw[1:], quote)
	if end < 0 {
		return "", fmt.Errorf("unterminated quoted value starting with %c", quote)
	}

	return raw[1 : end+1], nil
}
