package output

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warnColor    = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan)
	keyColor     = color.New(color.FgMagenta, color.Bold)
	valueColor   = color.New(color.FgWhite)
)

// Success prints a green success message.
func Success(w io.Writer, format string, args ...interface{}) {
	successColor.Fprintf(w, "✓ "+format+"\n", args...)
}

// Error prints a red error message.
func Error(w io.Writer, format string, args ...interface{}) {
	errorColor.Fprintf(w, "✗ "+format+"\n", args...)
}

// Warn prints a yellow warning message.
func Warn(w io.Writer, format string, args ...interface{}) {
	warnColor.Fprintf(w, "⚠ "+format+"\n", args...)
}

// Info prints a cyan info message.
func Info(w io.Writer, format string, args ...interface{}) {
	infoColor.Fprintf(w, "ℹ "+format+"\n", args...)
}

// KeyValue prints a key=value pair with color.
func KeyValue(w io.Writer, key, value string) {
	keyColor.Fprintf(w, "%s", key)
	fmt.Fprint(w, "=")
	valueColor.Fprintf(w, "%s\n", value)
}
