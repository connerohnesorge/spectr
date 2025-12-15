//nolint:revive // file-length-limit: line index requires comprehensive position tracking
package markdown

import "sort"

// Position represents a location in the source as line/column coordinates.
// Line numbers are 1-based (line 1 is the first line).
// Column numbers are 0-based byte offsets from the start of the line.
type Position struct {
	Line   int // 1-based line number
	Column int // 0-based byte offset within line
	Offset int // Original byte offset in source
}

// LineIndex provides efficient conversion from byte offsets to line/column positions.
// It uses lazy construction - the line index is only built on the first query.
// The index supports O(log n) lookups via binary search.
type LineIndex struct {
	source     []byte
	lineStarts []int // Byte offsets of each line start (0-indexed internally)
	built      bool  // Whether the index has been built
}

// NewLineIndex creates a new LineIndex for the given source.
// The index is lazily constructed on the first query.
func NewLineIndex(source []byte) *LineIndex {
	return &LineIndex{
		source: source,
		built:  false,
	}
}

// build constructs the line-start-offsets index.
// It scans the source for newlines and records the start position of each line.
// Handles both LF (\n) and CRLF (\r\n) line endings.
func (idx *LineIndex) build() {
	if idx.built {
		return
	}

	// First line always starts at offset 0
	idx.lineStarts = []int{0}

	i := 0
	for i < len(idx.source) {
		b := idx.source[i]
		switch b {
		case '\n':
			// LF: next line starts at i+1
			idx.lineStarts = append(
				idx.lineStarts,
				i+1,
			)
			i++
		case '\r':
			// Check for CRLF
			if i+1 < len(idx.source) &&
				idx.source[i+1] == '\n' {
				// CRLF: next line starts at i+2
				idx.lineStarts = append(
					idx.lineStarts,
					i+2,
				)
				i += 2
			} else {
				// Standalone CR: treat as line ending
				idx.lineStarts = append(idx.lineStarts, i+1)
				i++
			}
		default:
			i++
		}
	}

	idx.built = true
}

// LineCol returns the 1-based line number and 0-based column for a byte offset.
// If the offset is beyond the source length, returns the position at end of source.
// If the offset is negative, returns line 1, column 0.
func (idx *LineIndex) LineCol(
	offset int,
) (line, col int) {
	// Ensure the index is built
	idx.build()

	// Handle edge cases
	if offset < 0 {
		return 1, 0
	}
	if offset >= len(idx.source) {
		// Beyond source: return position at end
		if len(idx.lineStarts) == 0 {
			return 1, 0
		}
		lastLine := len(idx.lineStarts)
		lastLineStart := idx.lineStarts[lastLine-1]

		return lastLine, len(
			idx.source,
		) - lastLineStart
	}

	// Binary search to find the line containing the offset
	// We're looking for the largest lineStart that is <= offset
	lineIdx := sort.Search(
		len(idx.lineStarts),
		func(i int) bool {
			return idx.lineStarts[i] > offset
		},
	)

	// lineIdx is now the first line whose start is > offset
	// So the line we want is lineIdx - 1
	if lineIdx > 0 {
		lineIdx--
	}

	// Line numbers are 1-based
	line = lineIdx + 1
	// Column is the byte offset from the start of the line
	col = offset - idx.lineStarts[lineIdx]

	return line, col
}

// PositionAt returns a Position struct for the given byte offset.
// This combines the line, column, and original offset into a single struct.
func (idx *LineIndex) PositionAt(
	offset int,
) Position {
	line, col := idx.LineCol(offset)

	return Position{
		Line:   line,
		Column: col,
		Offset: offset,
	}
}

// LineCount returns the total number of lines in the source.
// This triggers index construction if not already built.
func (idx *LineIndex) LineCount() int {
	idx.build()

	return len(idx.lineStarts)
}

// LineStart returns the byte offset of the start of the given 1-based line number.
// Returns 0 for line <= 0 or if source is empty.
// Returns the start of the last line if lineNum exceeds the total line count.
func (idx *LineIndex) LineStart(lineNum int) int {
	idx.build()

	if lineNum <= 0 || len(idx.lineStarts) == 0 {
		return 0
	}

	// Convert 1-based to 0-based index
	lineIdx := lineNum - 1
	if lineIdx >= len(idx.lineStarts) {
		lineIdx = len(idx.lineStarts) - 1
	}

	return idx.lineStarts[lineIdx]
}

// LineEnd returns the byte offset of the end of the given 1-based line number.
// The end is the offset of the newline character (or end of source for last line).
// Returns 0 for line <= 0 or if source is empty.
func (idx *LineIndex) LineEnd(lineNum int) int {
	idx.build()

	if lineNum <= 0 || len(idx.lineStarts) == 0 {
		return 0
	}

	// Convert 1-based to 0-based index
	lineIdx := lineNum - 1
	if lineIdx >= len(idx.lineStarts) {
		// Beyond last line, return end of source
		return len(idx.source)
	}

	// If not the last line, end is just before the next line's start
	if lineIdx+1 < len(idx.lineStarts) {
		end := idx.lineStarts[lineIdx+1]
		// Exclude the newline character(s)
		if end > 0 && idx.source[end-1] == '\n' {
			end--
			// Also exclude \r in CRLF
			if end > 0 &&
				idx.source[end-1] == '\r' {
				end--
			}
		} else if end > 0 && idx.source[end-1] == '\r' {
			end--
		}

		return end
	}

	// Last line ends at end of source
	return len(idx.source)
}

// OffsetAt converts a 1-based line and 0-based column to a byte offset.
// Returns -1 if the line number is invalid.
// Clamps the column to the line's length if it exceeds it.
func (idx *LineIndex) OffsetAt(
	line, col int,
) int {
	idx.build()

	if line <= 0 || len(idx.lineStarts) == 0 {
		return -1
	}

	// Convert 1-based to 0-based index
	lineIdx := line - 1
	if lineIdx >= len(idx.lineStarts) {
		return -1
	}

	lineStart := idx.lineStarts[lineIdx]

	// Determine line length (excluding newline)
	var lineLen int
	if lineIdx+1 < len(idx.lineStarts) {
		lineEnd := idx.lineStarts[lineIdx+1]
		lineLen = lineEnd - lineStart
		// Exclude newline character(s) from length
		if lineLen > 0 &&
			idx.source[lineStart+lineLen-1] == '\n' {
			lineLen--
			if lineLen > 0 &&
				idx.source[lineStart+lineLen-1] == '\r' {
				lineLen--
			}
		} else if lineLen > 0 && idx.source[lineStart+lineLen-1] == '\r' {
			lineLen--
		}
	} else {
		// Last line
		lineLen = len(idx.source) - lineStart
	}

	// Clamp column to line length
	if col < 0 {
		col = 0
	}
	if col > lineLen {
		col = lineLen
	}

	return lineStart + col
}
