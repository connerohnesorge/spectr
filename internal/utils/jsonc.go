// Package utils provides utility functions for spectr.
package utils

import "strings"

// StripJSONCComments removes single-line comments from JSONC content.
func StripJSONCComments(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	var cleaned []string

	for _, line := range lines {
		line = stripFullLineComment(line)
		line = stripInlineComment(line)

		if strings.TrimSpace(line) != "" {
			cleaned = append(cleaned, line)
		}
	}

	return []byte(strings.Join(cleaned, "\n"))
}

// stripFullLineComment skips lines that are entirely comments.
func stripFullLineComment(line string) string {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "//") {
		return ""
	}

	return line
}

// stripInlineComment removes inline comments if they're not inside quotes.
func stripInlineComment(line string) string {
	idx := strings.Index(line, "//")
	if idx == -1 {
		return line
	}

	if isInsideQuotes(line, idx) {
		return line
	}

	return strings.TrimRight(line[:idx], " \t")
}

// isInsideQuotes checks if position idx in line is inside quotes.
func isInsideQuotes(line string, idx int) bool {
	quoteCount := 0
	for i := range idx {
		if line[i] == '"' && (i == 0 || line[i-1] != '\\') {
			quoteCount++
		}
	}

	return quoteCount%2 != 0
}
