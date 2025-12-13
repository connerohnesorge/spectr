package markdown

import (
	"bytes"

	"github.com/connerohnesorge/spectr/internal/specterrs"
	"github.com/russross/blackfriday/v2"
)

// lineIndex maps byte offsets to line numbers for efficient lookups.
type lineIndex struct {
	source     []byte
	lineStarts []int // Byte offsets where each line starts
}

// newLineIndex builds a line index from source content.
func newLineIndex(source []byte) *lineIndex {
	starts := []int{0} // Line 0 starts at byte 0

	for i, b := range source {
		if b == '\n' && i+1 < len(source) {
			starts = append(starts, i+1)
		}
	}

	return &lineIndex{
		source:     source,
		lineStarts: starts,
	}
}

// lineAt returns the 1-indexed line number for the given byte offset.
func (li *lineIndex) lineAt(offset int) int {
	if li == nil || len(li.lineStarts) == 0 {
		return 1
	}

	// Binary search for the line containing this offset
	low, high := 0, len(li.lineStarts)-1
	for low < high {
		mid := (low + high + 1) / 2
		if li.lineStarts[mid] <= offset {
			low = mid
		} else {
			high = mid - 1
		}
	}

	return low + 1 // Convert to 1-indexed
}

// ParseDocument parses markdown content and extracts all structural elements.
// Returns error for invalid input (empty content, binary data).
// Blackfriday internals are never exposed; all data is package-defined.
func ParseDocument(content []byte) (*Document, error) {
	// Check for empty content
	if len(bytes.TrimSpace(content)) == 0 {
		return nil, &specterrs.EmptyContentError{}
	}

	// Check for binary content (contains null bytes)
	if bytes.ContainsRune(content, '\x00') {
		return nil, &specterrs.BinaryContentError{}
	}

	// Build line index for position lookups
	lineIndex := newLineIndex(content)

	// Parse with blackfriday using common extensions
	extensions := blackfriday.CommonExtensions
	parser := blackfriday.New(blackfriday.WithExtensions(extensions))
	node := parser.Parse(content)

	// Extract all structural elements
	headers := extractHeaders(node, lineIndex)
	sections := extractSections(node, content, headers, lineIndex)
	tasks := extractTasks(node, content, lineIndex)

	return &Document{
		Headers:  headers,
		Sections: sections,
		Tasks:    tasks,
	}, nil
}
