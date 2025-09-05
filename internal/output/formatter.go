package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/mdryaan/vaultenv/internal/vault"
)

// Format is the output format type.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatEnv   Format = "env"
)

// Formatter writes entries to the given writer.
type Formatter interface {
	WriteEntries(w io.Writer, entries []vault.Entry, showValues bool) error
}

// Get returns the appropriate Formatter for the given format string.
func Get(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	default:
		return &TableFormatter{}
	}
}

// JSONFormatter outputs entries as JSON.
type JSONFormatter struct{}

type jsonEntry struct {
	Key       string   `json:"key"`
	Value     string   `json:"value,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func (f *JSONFormatter) WriteEntries(w io.Writer, entries []vault.Entry, showValues bool) error {
	out := make([]jsonEntry, len(entries))
	for i, e := range entries {
		je := jsonEntry{
			Key:       e.Key,
			Tags:      e.Tags,
			CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if showValues {
			je.Value = e.Value
		} else {
			je.Value = maskValue(e.Value)
		}
		out[i] = je
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func maskValue(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return "****" + value[len(value)-4:]
}
