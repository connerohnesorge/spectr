package tui

import (
	"encoding/base64"
	"fmt"

	"github.com/atotto/clipboard"
)

const (
	// EllipsisMinLength is the minimum string length before
	// truncation adds ellipsis.
	EllipsisMinLength = 3
)

// TruncateString truncates a string and adds ellipsis if needed.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= EllipsisMinLength {
		return s[:maxLen]
	}

	return s[:maxLen-EllipsisMinLength] + "..."
}

// CopyToClipboard copies text to clipboard using native or OSC 52.
func CopyToClipboard(text string) error {
	// Try native clipboard first
	err := clipboard.WriteAll(text)
	if err == nil {
		return nil
	}

	// Fallback to OSC 52 for SSH sessions
	encoded := base64.StdEncoding.EncodeToString(
		[]byte(text),
	)
	osc52 := "\x1b]52;c;" + encoded + "\x07"
	fmt.Print(osc52)

	// OSC 52 doesn't report errors, consider it successful
	return nil
}
