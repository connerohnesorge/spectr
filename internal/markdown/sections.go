package markdown

import (
	"bytes"

	"github.com/russross/blackfriday/v2"
)

// sectionBoundary represents the position of a header in the source.
type sectionBoundary struct {
	header    Header
	startByte int // Byte offset of the header line start
}

// extractSections extracts section content between headers.
// A section is the content between one header and the next same-level header.
// Returns sections keyed by header text.
func extractSections(
	_ *blackfriday.Node,
	source []byte,
	headers []Header,
	_ *lineIndex,
) map[string]Section {
	sections := make(map[string]Section)

	if len(headers) == 0 {
		return sections
	}

	// Find byte positions for each header
	boundaries := findSectionBoundaries(source, headers)

	// Extract content for each section
	for i, boundary := range boundaries {
		// Content starts after the header line
		contentStart := findLineEnd(source, boundary.startByte) + 1

		// Content ends at the start of next same/higher level header, or EOF.
		// Lower-level headers (H4 under H3) are INCLUDED in section content.
		contentEnd := len(source)
		for j := i + 1; j < len(boundaries); j++ {
			// Section ends at next header of same or higher level only
			if boundaries[j].header.Level <= boundary.header.Level {
				contentEnd = boundaries[j].startByte

				break
			}
			// Lower-level headers are included, so continue looking
		}

		// Extract and trim the content
		if contentStart < contentEnd && contentStart < len(source) {
			content := bytes.TrimSpace(source[contentStart:contentEnd])
			sections[boundary.header.Text] = Section{
				Header:  boundary.header,
				Content: string(content),
			}
		} else {
			// Empty section
			sections[boundary.header.Text] = Section{
				Header:  boundary.header,
				Content: "",
			}
		}
	}

	return sections
}

// findSectionBoundaries locates headers in source and returns byte positions.
func findSectionBoundaries(source []byte, headers []Header) []sectionBoundary {
	boundaries := make([]sectionBoundary, 0, len(headers))

	// Track which headers we've found to avoid duplicates
	foundHeaders := make(map[int]bool)

	for _, header := range headers {
		offset := findHeaderInSource(source, header, foundHeaders)
		if offset >= 0 {
			foundHeaders[offset] = true
			boundaries = append(boundaries, sectionBoundary{
				header:    header,
				startByte: offset,
			})
		}
	}

	return boundaries
}

// findHeaderInSource searches for a header in source.
// Returns the byte offset or -1 if not found.
func findHeaderInSource(
	source []byte,
	header Header,
	foundHeaders map[int]bool,
) int {
	pattern := buildHeaderPattern(header.Level)
	searchStart := 0

	for {
		idx := bytes.Index(source[searchStart:], pattern)
		if idx < 0 {
			return -1
		}

		pos := searchStart + idx
		if shouldSkipPosition(source, pos, header.Level, foundHeaders) {
			searchStart = pos + 1

			continue
		}

		if matchesHeaderText(source, pos, header) {
			return pos
		}

		searchStart = pos + 1
	}
}

// buildHeaderPattern creates the hash pattern for a header level.
func buildHeaderPattern(level int) []byte {
	pattern := make([]byte, level)
	for i := range level {
		pattern[i] = '#'
	}

	return pattern
}

// shouldSkipPosition checks if we should skip this position.
func shouldSkipPosition(
	source []byte,
	pos, level int,
	foundHeaders map[int]bool,
) bool {
	if foundHeaders[pos] {
		return true
	}
	if pos > 0 && source[pos-1] != '\n' {
		return true
	}
	hashEnd := pos + level
	if hashEnd < len(source) && source[hashEnd] == '#' {
		return true
	}

	return false
}

// matchesHeaderText checks if the header at pos matches the expected text.
func matchesHeaderText(source []byte, pos int, header Header) bool {
	hashEnd := pos + header.Level
	textStart := skipWS(source, hashEnd)
	lineEnd := findNextNewline(source, textStart)
	headerText := bytes.TrimSpace(source[textStart:lineEnd])

	return string(headerText) == header.Text
}

// skipWS advances past whitespace.
func skipWS(source []byte, startPos int) int {
	pos := startPos
	for pos < len(source) && (source[pos] == ' ' || source[pos] == '\t') {
		pos++
	}

	return pos
}

// findNextNewline finds the next newline position.
func findNextNewline(source []byte, startPos int) int {
	pos := startPos
	for pos < len(source) && source[pos] != '\n' {
		pos++
	}

	return pos
}

// findLineEnd finds the byte offset of the newline character ending the line
// that contains the given offset, or len(source)-1 if no newline.
func findLineEnd(source []byte, offset int) int {
	for i := offset; i < len(source); i++ {
		if source[i] == '\n' {
			return i
		}
	}

	return len(source) - 1
}
