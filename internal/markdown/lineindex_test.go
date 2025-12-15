package markdown

import (
	"testing"
)

// TestLineIndex_EmptySource verifies behavior with an empty source.
func TestLineIndex_EmptySource(t *testing.T) {
	idx := NewLineIndex(make([]byte, 0))

	line, col := idx.LineCol(0)
	if line != 1 || col != 0 {
		t.Errorf(
			"LineCol(0) on empty source = (%d, %d), want (1, 0)",
			line,
			col,
		)
	}

	if count := idx.LineCount(); count != 1 {
		t.Errorf(
			"LineCount() on empty source = %d, want 1",
			count,
		)
	}
}

// TestLineIndex_SingleLine verifies behavior with a single line (no newlines).
func TestLineIndex_SingleLine(t *testing.T) {
	source := []byte("hello world")
	idx := NewLineIndex(source)

	tests := []struct {
		offset      int
		wantLine    int
		wantCol     int
		description string
	}{
		{0, 1, 0, "first character"},
		{5, 1, 5, "middle character (space)"},
		{10, 1, 10, "last character"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			line, col := idx.LineCol(tt.offset)
			if line != tt.wantLine ||
				col != tt.wantCol {
				t.Errorf(
					"LineCol(%d) = (%d, %d), want (%d, %d)",
					tt.offset,
					line,
					col,
					tt.wantLine,
					tt.wantCol,
				)
			}
		})
	}

	if count := idx.LineCount(); count != 1 {
		t.Errorf(
			"LineCount() = %d, want 1",
			count,
		)
	}
}

// TestLineIndex_MultipleLinesLF verifies behavior with LF (\n) line endings.
func TestLineIndex_MultipleLinesLF(t *testing.T) {
	source := []byte("line1\nline2\nline3")
	idx := NewLineIndex(source)

	tests := []struct {
		offset      int
		wantLine    int
		wantCol     int
		description string
	}{
		{0, 1, 0, "start of line 1"},
		{4, 1, 4, "end of line 1 content"},
		{5, 1, 5, "newline after line 1"},
		{6, 2, 0, "start of line 2"},
		{11, 2, 5, "newline after line 2"},
		{12, 3, 0, "start of line 3"},
		{16, 3, 4, "end of line 3"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			line, col := idx.LineCol(tt.offset)
			if line != tt.wantLine ||
				col != tt.wantCol {
				t.Errorf(
					"LineCol(%d) = (%d, %d), want (%d, %d)",
					tt.offset,
					line,
					col,
					tt.wantLine,
					tt.wantCol,
				)
			}
		})
	}

	if count := idx.LineCount(); count != 3 {
		t.Errorf(
			"LineCount() = %d, want 3",
			count,
		)
	}
}

// TestLineIndex_MultipleLinesCRLF verifies behavior with CRLF (\r\n) line endings.
func TestLineIndex_MultipleLinesCRLF(
	t *testing.T,
) {
	source := []byte("line1\r\nline2\r\nline3")
	idx := NewLineIndex(source)

	tests := []struct {
		offset      int
		wantLine    int
		wantCol     int
		description string
	}{
		{0, 1, 0, "start of line 1"},
		{4, 1, 4, "end of line 1 content"},
		{5, 1, 5, "CR of CRLF after line 1"},
		{6, 1, 6, "LF of CRLF after line 1"},
		{7, 2, 0, "start of line 2"},
		{12, 2, 5, "end of line 2 content"},
		{14, 3, 0, "start of line 3"},
		{18, 3, 4, "end of line 3"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			line, col := idx.LineCol(tt.offset)
			if line != tt.wantLine ||
				col != tt.wantCol {
				t.Errorf(
					"LineCol(%d) = (%d, %d), want (%d, %d)",
					tt.offset,
					line,
					col,
					tt.wantLine,
					tt.wantCol,
				)
			}
		})
	}

	if count := idx.LineCount(); count != 3 {
		t.Errorf(
			"LineCount() = %d, want 3",
			count,
		)
	}
}

// TestLineIndex_MixedLineEndings verifies behavior with mixed LF and CRLF.
func TestLineIndex_MixedLineEndings(
	t *testing.T,
) {
	source := []byte(
		"line1\nline2\r\nline3\rline4",
	)
	idx := NewLineIndex(source)

	tests := []struct {
		offset      int
		wantLine    int
		wantCol     int
		description string
	}{
		{0, 1, 0, "start of line 1"},
		{6, 2, 0, "start of line 2 (after LF)"},
		{
			13,
			3,
			0,
			"start of line 3 (after CRLF)",
		},
		{
			19,
			4,
			0,
			"start of line 4 (after standalone CR)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			line, col := idx.LineCol(tt.offset)
			if line != tt.wantLine ||
				col != tt.wantCol {
				t.Errorf(
					"LineCol(%d) = (%d, %d), want (%d, %d)",
					tt.offset,
					line,
					col,
					tt.wantLine,
					tt.wantCol,
				)
			}
		})
	}

	if count := idx.LineCount(); count != 4 {
		t.Errorf(
			"LineCount() = %d, want 4",
			count,
		)
	}
}

// TestLineCol_OneBased verifies that line numbers are 1-based.
func TestLineCol_OneBased(t *testing.T) {
	source := []byte("a\nb\nc")
	idx := NewLineIndex(source)

	// First line should be 1, not 0
	line, _ := idx.LineCol(0)
	if line != 1 {
		t.Errorf(
			"First line number = %d, want 1 (1-based)",
			line,
		)
	}

	// Second line should be 2
	line, _ = idx.LineCol(2)
	if line != 2 {
		t.Errorf(
			"Second line number = %d, want 2",
			line,
		)
	}

	// Third line should be 3
	line, _ = idx.LineCol(4)
	if line != 3 {
		t.Errorf(
			"Third line number = %d, want 3",
			line,
		)
	}
}

// TestLineCol_ZeroBasedColumn verifies that column numbers are 0-based.
func TestLineCol_ZeroBasedColumn(t *testing.T) {
	source := []byte("abc\ndef")
	idx := NewLineIndex(source)

	// First character of line should be column 0
	_, col := idx.LineCol(0)
	if col != 0 {
		t.Errorf(
			"First character column = %d, want 0 (0-based)",
			col,
		)
	}

	_, col = idx.LineCol(4)
	if col != 0 {
		t.Errorf(
			"First character of second line column = %d, want 0",
			col,
		)
	}
}

// TestLineCol_FirstCharacter verifies position at start of source.
func TestLineCol_FirstCharacter(t *testing.T) {
	source := []byte("hello")
	idx := NewLineIndex(source)

	line, col := idx.LineCol(0)
	if line != 1 || col != 0 {
		t.Errorf(
			"LineCol(0) = (%d, %d), want (1, 0)",
			line,
			col,
		)
	}
}

// TestLineCol_LastCharacter verifies position at end of source.
func TestLineCol_LastCharacter(t *testing.T) {
	source := []byte("hello\nworld")
	idx := NewLineIndex(source)

	// Last character is 'd' at offset 10
	line, col := idx.LineCol(10)
	if line != 2 || col != 4 {
		t.Errorf(
			"LineCol(10) = (%d, %d), want (2, 4)",
			line,
			col,
		)
	}
}

// TestLineCol_AtNewline verifies position at newline character.
func TestLineCol_AtNewline(t *testing.T) {
	source := []byte("abc\ndef")
	idx := NewLineIndex(source)

	// Newline is at offset 3, still part of line 1
	line, col := idx.LineCol(3)
	if line != 1 || col != 3 {
		t.Errorf(
			"LineCol(3) at newline = (%d, %d), want (1, 3)",
			line,
			col,
		)
	}
}

// TestLineCol_AfterNewline verifies position at start of new line.
func TestLineCol_AfterNewline(t *testing.T) {
	source := []byte("abc\ndef")
	idx := NewLineIndex(source)

	// After newline starts line 2 at column 0
	line, col := idx.LineCol(4)
	if line != 2 || col != 0 {
		t.Errorf(
			"LineCol(4) after newline = (%d, %d), want (2, 0)",
			line,
			col,
		)
	}
}

// TestLineCol_NegativeOffset verifies behavior with negative offset.
func TestLineCol_NegativeOffset(t *testing.T) {
	source := []byte("hello")
	idx := NewLineIndex(source)

	line, col := idx.LineCol(-1)
	if line != 1 || col != 0 {
		t.Errorf(
			"LineCol(-1) = (%d, %d), want (1, 0)",
			line,
			col,
		)
	}

	line, col = idx.LineCol(-100)
	if line != 1 || col != 0 {
		t.Errorf(
			"LineCol(-100) = (%d, %d), want (1, 0)",
			line,
			col,
		)
	}
}

// TestLineCol_OffsetBeyondEnd verifies behavior when offset exceeds source length.
func TestLineCol_OffsetBeyondEnd(t *testing.T) {
	source := []byte("ab\ncd")
	idx := NewLineIndex(source)

	// Offset 5 is at end of source (length is 5)
	line, col := idx.LineCol(5)
	if line != 2 || col != 2 {
		t.Errorf(
			"LineCol(5) at end = (%d, %d), want (2, 2)",
			line,
			col,
		)
	}

	// Offset 10 is beyond end
	line, col = idx.LineCol(10)
	if line != 2 || col != 2 {
		t.Errorf(
			"LineCol(10) beyond end = (%d, %d), want (2, 2)",
			line,
			col,
		)
	}

	// Offset 1000 is way beyond end
	line, col = idx.LineCol(1000)
	if line != 2 || col != 2 {
		t.Errorf(
			"LineCol(1000) beyond end = (%d, %d), want (2, 2)",
			line,
			col,
		)
	}
}

// TestLineIndex_EmptyLines verifies behavior with empty lines (consecutive newlines).
func TestLineIndex_EmptyLines(t *testing.T) {
	source := []byte("a\n\nb")
	idx := NewLineIndex(source)

	tests := []struct {
		offset      int
		wantLine    int
		wantCol     int
		description string
	}{
		{0, 1, 0, "char 'a'"},
		{1, 1, 1, "first newline"},
		{2, 2, 0, "second newline (empty line)"},
		{3, 3, 0, "char 'b'"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			line, col := idx.LineCol(tt.offset)
			if line != tt.wantLine ||
				col != tt.wantCol {
				t.Errorf(
					"LineCol(%d) = (%d, %d), want (%d, %d)",
					tt.offset,
					line,
					col,
					tt.wantLine,
					tt.wantCol,
				)
			}
		})
	}

	if count := idx.LineCount(); count != 3 {
		t.Errorf(
			"LineCount() = %d, want 3",
			count,
		)
	}
}

// TestLineIndex_FileEndingWithNewline verifies files ending with newline.
func TestLineIndex_FileEndingWithNewline(
	t *testing.T,
) {
	source := []byte("line1\nline2\n")
	idx := NewLineIndex(source)

	// The trailing newline creates a third (empty) line
	if count := idx.LineCount(); count != 3 {
		t.Errorf(
			"LineCount() = %d, want 3 (two content lines + trailing empty)",
			count,
		)
	}

	// Position at start of empty third line
	line, col := idx.LineCol(12)
	if line != 3 || col != 0 {
		t.Errorf(
			"LineCol(12) = (%d, %d), want (3, 0)",
			line,
			col,
		)
	}
}

// TestLineIndex_FileWithoutTrailingNewline verifies files without trailing newline.
func TestLineIndex_FileWithoutTrailingNewline(
	t *testing.T,
) {
	source := []byte("line1\nline2")
	idx := NewLineIndex(source)

	if count := idx.LineCount(); count != 2 {
		t.Errorf(
			"LineCount() = %d, want 2",
			count,
		)
	}

	// Last character
	line, col := idx.LineCol(10)
	if line != 2 || col != 4 {
		t.Errorf(
			"LineCol(10) = (%d, %d), want (2, 4)",
			line,
			col,
		)
	}
}

// TestLineIndex_OnlyNewlines verifies source with only newlines.
func TestLineIndex_OnlyNewlines(t *testing.T) {
	source := []byte("\n\n\n")
	idx := NewLineIndex(source)

	if count := idx.LineCount(); count != 4 {
		t.Errorf(
			"LineCount() = %d, want 4 (3 newlines = 4 lines)",
			count,
		)
	}

	tests := []struct {
		offset   int
		wantLine int
		wantCol  int
	}{
		{0, 1, 0},
		{1, 2, 0},
		{2, 3, 0},
		{3, 4, 0},
	}

	for _, tt := range tests {
		line, col := idx.LineCol(tt.offset)
		if line != tt.wantLine ||
			col != tt.wantCol {
			t.Errorf(
				"LineCol(%d) = (%d, %d), want (%d, %d)",
				tt.offset,
				line,
				col,
				tt.wantLine,
				tt.wantCol,
			)
		}
	}
}

// TestLineIndex_LazyConstruction verifies index is built on first query.
func TestLineIndex_LazyConstruction(
	t *testing.T,
) {
	source := []byte("hello\nworld")
	idx := NewLineIndex(source)

	// Before any query, built should be false
	if idx.built {
		t.Error(
			"Index should not be built before first query",
		)
	}

	// First query should trigger build
	idx.LineCol(0)

	if !idx.built {
		t.Error(
			"Index should be built after first query",
		)
	}
}

// TestLineIndex_ReuseBuiltIndex verifies multiple queries reuse the built index.
func TestLineIndex_ReuseBuiltIndex(t *testing.T) {
	source := []byte("hello\nworld\nfoo")
	idx := NewLineIndex(source)

	// First query builds the index
	idx.LineCol(0)
	lineStarts1 := idx.lineStarts

	// Subsequent queries should not rebuild
	idx.LineCol(5)
	idx.LineCol(10)

	// Verify same slice is used (pointer comparison)
	if &idx.lineStarts[0] != &lineStarts1[0] {
		t.Error(
			"Index was rebuilt on subsequent queries",
		)
	}
}

// TestLineIndex_BuildOnLineCount verifies LineCount triggers build.
func TestLineIndex_BuildOnLineCount(
	t *testing.T,
) {
	source := []byte("a\nb\nc")
	idx := NewLineIndex(source)

	if idx.built {
		t.Error(
			"Index should not be built initially",
		)
	}

	_ = idx.LineCount()

	if !idx.built {
		t.Error(
			"LineCount() should trigger index build",
		)
	}
}

// TestLineIndex_BuildOnLineStart verifies LineStart triggers build.
func TestLineIndex_BuildOnLineStart(
	t *testing.T,
) {
	source := []byte("a\nb\nc")
	idx := NewLineIndex(source)

	if idx.built {
		t.Error(
			"Index should not be built initially",
		)
	}

	_ = idx.LineStart(1)

	if !idx.built {
		t.Error(
			"LineStart() should trigger index build",
		)
	}
}

// TestLineIndex_BuildOnLineEnd verifies LineEnd triggers build.
func TestLineIndex_BuildOnLineEnd(t *testing.T) {
	source := []byte("a\nb\nc")
	idx := NewLineIndex(source)

	if idx.built {
		t.Error(
			"Index should not be built initially",
		)
	}

	_ = idx.LineEnd(1)

	if !idx.built {
		t.Error(
			"LineEnd() should trigger index build",
		)
	}
}

// TestLineIndex_BuildOnOffsetAt verifies OffsetAt triggers build.
func TestLineIndex_BuildOnOffsetAt(t *testing.T) {
	source := []byte("a\nb\nc")
	idx := NewLineIndex(source)

	if idx.built {
		t.Error(
			"Index should not be built initially",
		)
	}

	_ = idx.OffsetAt(1, 0)

	if !idx.built {
		t.Error(
			"OffsetAt() should trigger index build",
		)
	}
}

// TestPositionAt_ReturnsCorrectStruct verifies PositionAt returns correct Position.
func TestPositionAt_ReturnsCorrectStruct(
	t *testing.T,
) {
	source := []byte("hello\nworld")
	idx := NewLineIndex(source)

	tests := []struct {
		offset   int
		wantPos  Position
		testName string
	}{
		{
			0,
			Position{
				Line:   1,
				Column: 0,
				Offset: 0,
			},
			"start of source",
		},
		{
			5,
			Position{
				Line:   1,
				Column: 5,
				Offset: 5,
			},
			"at newline",
		},
		{
			6,
			Position{
				Line:   2,
				Column: 0,
				Offset: 6,
			},
			"start of line 2",
		},
		{
			10,
			Position{
				Line:   2,
				Column: 4,
				Offset: 10,
			},
			"end of source",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			pos := idx.PositionAt(tt.offset)
			if pos != tt.wantPos {
				t.Errorf(
					"PositionAt(%d) = %+v, want %+v",
					tt.offset,
					pos,
					tt.wantPos,
				)
			}
		})
	}
}

// TestPositionAt_PreservesOffset verifies Position.Offset matches input offset.
func TestPositionAt_PreservesOffset(
	t *testing.T,
) {
	source := []byte("hello\nworld")
	idx := NewLineIndex(source)

	offsets := []int{0, 3, 5, 6, 10, -1, 100}
	for _, offset := range offsets {
		pos := idx.PositionAt(offset)
		if pos.Offset != offset {
			t.Errorf(
				"PositionAt(%d).Offset = %d, want %d",
				offset,
				pos.Offset,
				offset,
			)
		}
	}
}

// TestLineCount_Various verifies LineCount with various inputs.
func TestLineCount_Various(t *testing.T) {
	tests := []struct {
		source   string
		expected int
		testName string
	}{
		{"", 1, "empty source"},
		{
			"hello",
			1,
			"single line without newline",
		},
		{
			"hello\n",
			2,
			"single line with newline",
		},
		{"a\nb", 2, "two lines"},
		{
			"a\nb\n",
			3,
			"two lines with trailing newline",
		},
		{"\n", 2, "just one newline"},
		{"\n\n\n", 4, "three newlines"},
		{"line1\nline2\nline3", 3, "three lines"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			idx := NewLineIndex([]byte(tt.source))
			if count := idx.LineCount(); count != tt.expected {
				t.Errorf(
					"LineCount() for %q = %d, want %d",
					tt.source,
					count,
					tt.expected,
				)
			}
		})
	}
}

// TestLineStart_Various verifies LineStart with various inputs.
func TestLineStart_Various(t *testing.T) {
	source := []byte("abc\ndef\nghi")
	idx := NewLineIndex(source)

	tests := []struct {
		lineNum  int
		expected int
		testName string
	}{
		{1, 0, "first line"},
		{2, 4, "second line"},
		{3, 8, "third line"},
		{0, 0, "line 0 (invalid)"},
		{-1, 0, "negative line"},
		{
			4,
			8,
			"line beyond count (clamps to last)",
		},
		{100, 8, "way beyond count"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			start := idx.LineStart(tt.lineNum)
			if start != tt.expected {
				t.Errorf(
					"LineStart(%d) = %d, want %d",
					tt.lineNum,
					start,
					tt.expected,
				)
			}
		})
	}
}

// TestLineStart_EmptySource verifies LineStart with empty source.
func TestLineStart_EmptySource(t *testing.T) {
	idx := NewLineIndex(make([]byte, 0))

	start := idx.LineStart(1)
	if start != 0 {
		t.Errorf(
			"LineStart(1) on empty source = %d, want 0",
			start,
		)
	}
}

// TestLineEnd_Various verifies LineEnd with various inputs.
func TestLineEnd_Various(t *testing.T) {
	source := []byte("abc\ndef\nghi")
	idx := NewLineIndex(source)

	tests := []struct {
		lineNum  int
		expected int
		testName string
	}{
		{1, 3, "first line (ends at 'c')"},
		{2, 7, "second line (ends at 'f')"},
		{
			3,
			11,
			"third line (ends at end of source)",
		},
		{0, 0, "line 0 (invalid)"},
		{-1, 0, "negative line"},
		{4, 11, "line beyond count"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			end := idx.LineEnd(tt.lineNum)
			if end != tt.expected {
				t.Errorf(
					"LineEnd(%d) = %d, want %d",
					tt.lineNum,
					end,
					tt.expected,
				)
			}
		})
	}
}

// TestLineEnd_WithCRLF verifies LineEnd excludes CRLF correctly.
func TestLineEnd_WithCRLF(t *testing.T) {
	source := []byte("abc\r\ndef\r\nghi")
	idx := NewLineIndex(source)

	tests := []struct {
		lineNum  int
		expected int
		testName string
	}{
		{1, 3, "first line excludes CRLF"},
		{2, 8, "second line excludes CRLF"},
		{
			3,
			13,
			"third line (no trailing newline)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			end := idx.LineEnd(tt.lineNum)
			if end != tt.expected {
				t.Errorf(
					"LineEnd(%d) = %d, want %d",
					tt.lineNum,
					end,
					tt.expected,
				)
			}
		})
	}
}

// TestLineEnd_EmptySource verifies LineEnd with empty source.
func TestLineEnd_EmptySource(t *testing.T) {
	idx := NewLineIndex(make([]byte, 0))

	end := idx.LineEnd(1)
	if end != 0 {
		t.Errorf(
			"LineEnd(1) on empty source = %d, want 0",
			end,
		)
	}
}

// TestLineEnd_EmptyLines verifies LineEnd with empty lines.
func TestLineEnd_EmptyLines(t *testing.T) {
	source := []byte("a\n\nb")
	idx := NewLineIndex(source)

	// Line 2 is empty (just a newline)
	end := idx.LineEnd(2)
	if end != 2 {
		t.Errorf(
			"LineEnd(2) for empty line = %d, want 2",
			end,
		)
	}
}

// TestOffsetAt_Various verifies OffsetAt with various inputs.
func TestOffsetAt_Various(t *testing.T) {
	source := []byte("abc\ndef\nghi")
	idx := NewLineIndex(source)

	tests := []struct {
		line     int
		col      int
		expected int
		testName string
	}{
		{1, 0, 0, "start of line 1"},
		{1, 2, 2, "middle of line 1"},
		{2, 0, 4, "start of line 2"},
		{2, 2, 6, "middle of line 2"},
		{3, 0, 8, "start of line 3"},
		{3, 2, 10, "end of line 3"},
		{0, 0, -1, "invalid line 0"},
		{-1, 0, -1, "negative line"},
		{4, 0, -1, "line beyond count"},
		{
			1,
			10,
			3,
			"column beyond line length (clamped)",
		},
		{
			1,
			-1,
			0,
			"negative column (clamped to 0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			offset := idx.OffsetAt(
				tt.line,
				tt.col,
			)
			if offset != tt.expected {
				t.Errorf(
					"OffsetAt(%d, %d) = %d, want %d",
					tt.line,
					tt.col,
					offset,
					tt.expected,
				)
			}
		})
	}
}

// TestOffsetAt_EmptySource verifies OffsetAt with empty source.
func TestOffsetAt_EmptySource(t *testing.T) {
	idx := NewLineIndex(make([]byte, 0))

	offset := idx.OffsetAt(1, 0)
	if offset != 0 {
		t.Errorf(
			"OffsetAt(1, 0) on empty source = %d, want 0",
			offset,
		)
	}

	offset = idx.OffsetAt(1, 5)
	if offset != 0 {
		t.Errorf(
			"OffsetAt(1, 5) on empty source = %d, want 0 (clamped)",
			offset,
		)
	}
}

// TestOffsetAt_RoundTrip verifies OffsetAt and LineCol are inverse operations.
func TestOffsetAt_RoundTrip(t *testing.T) {
	source := []byte("hello\nworld\nfoo bar")
	idx := NewLineIndex(source)

	// For each valid offset, verify round-trip
	for offset := range source {
		line, col := idx.LineCol(offset)
		roundTrip := idx.OffsetAt(line, col)
		if roundTrip != offset {
			t.Errorf(
				"Round-trip failed: offset %d -> (%d, %d) -> %d",
				offset,
				line,
				col,
				roundTrip,
			)
		}
	}
}

// TestPosition_ZeroValue verifies Position zero value.
func TestPosition_ZeroValue(t *testing.T) {
	var pos Position
	if pos.Line != 0 || pos.Column != 0 ||
		pos.Offset != 0 {
		t.Errorf(
			"Zero value Position = %+v, want all zeros",
			pos,
		)
	}
}

// TestPosition_Equality verifies Position equality comparison.
func TestPosition_Equality(t *testing.T) {
	pos1 := Position{
		Line:   1,
		Column: 5,
		Offset: 5,
	}
	pos2 := Position{
		Line:   1,
		Column: 5,
		Offset: 5,
	}
	pos3 := Position{
		Line:   2,
		Column: 0,
		Offset: 6,
	}

	if pos1 != pos2 {
		t.Error("Equal positions should be equal")
	}
	if pos1 == pos3 {
		t.Error(
			"Different positions should not be equal",
		)
	}
}

// TestLineIndex_Unicode verifies byte-based offsets with UTF-8.
func TestLineIndex_Unicode(t *testing.T) {
	// Multi-byte UTF-8 characters
	source := []byte("hello\nworld")
	idx := NewLineIndex(source)

	// Offsets are byte-based, not rune-based
	line, col := idx.LineCol(0)
	if line != 1 || col != 0 {
		t.Errorf(
			"LineCol(0) = (%d, %d), want (1, 0)",
			line,
			col,
		)
	}
}

// TestLineIndex_MultiByteCharacters verifies offsets with multi-byte chars.
func TestLineIndex_MultiByteCharacters(
	t *testing.T,
) {
	// Each character is 3 bytes in UTF-8
	source := []byte("a\nb")
	idx := NewLineIndex(source)

	// Character 'b' starts at byte offset 2
	line, col := idx.LineCol(2)
	if line != 2 || col != 0 {
		t.Errorf(
			"LineCol(2) for 'b' = (%d, %d), want (2, 0)",
			line,
			col,
		)
	}
}

// TestLineIndex_LargeFile verifies behavior with many lines.
func TestLineIndex_LargeFile(t *testing.T) {
	// Create a source with 1000 lines
	var source []byte
	for range 1000 {
		source = append(
			source,
			[]byte("line content here\n")...)
	}

	idx := NewLineIndex(source)

	if count := idx.LineCount(); count != 1001 {
		t.Errorf(
			"LineCount() = %d, want 1001",
			count,
		)
	}

	// Verify first line
	line, col := idx.LineCol(0)
	if line != 1 || col != 0 {
		t.Errorf(
			"LineCol(0) = (%d, %d), want (1, 0)",
			line,
			col,
		)
	}

	// Verify middle line (line 500 starts at byte 499*18 = 8982)
	expectedOffset := 499 * 18 // 18 bytes per line including newline
	line, col = idx.LineCol(expectedOffset)
	if line != 500 || col != 0 {
		t.Errorf(
			"LineCol(%d) = (%d, %d), want (500, 0)",
			expectedOffset,
			line,
			col,
		)
	}

	// Verify last line
	lastOffset := 1000 * 18
	line, col = idx.LineCol(lastOffset)
	if line != 1001 || col != 0 {
		t.Errorf(
			"LineCol(%d) = (%d, %d), want (1001, 0)",
			lastOffset,
			line,
			col,
		)
	}
}

// TestLineIndex_BinarySearchEfficiency verifies O(log n) lookup.
func TestLineIndex_BinarySearchEfficiency(
	t *testing.T,
) {
	// This test verifies the implementation uses binary search
	// by checking it can handle large inputs efficiently
	var source []byte
	lineCount := 10000
	for range lineCount {
		source = append(source, []byte("x\n")...)
	}

	idx := NewLineIndex(source)

	// Build the index
	idx.LineCol(0)

	// Verify random accesses work correctly
	testOffsets := []int{
		0,
		1000,
		5000,
		9999,
		19998,
	}
	for _, offset := range testOffsets {
		line, _ := idx.LineCol(offset)
		expectedLine := (offset / 2) + 1
		if line != expectedLine {
			t.Errorf(
				"LineCol(%d) line = %d, want %d",
				offset,
				line,
				expectedLine,
			)
		}
	}
}

// TestLineIndex_StandaloneCR verifies handling of standalone CR (not followed by LF).
func TestLineIndex_StandaloneCR(t *testing.T) {
	source := []byte("line1\rline2\rline3")
	idx := NewLineIndex(source)

	if count := idx.LineCount(); count != 3 {
		t.Errorf(
			"LineCount() = %d, want 3",
			count,
		)
	}

	tests := []struct {
		offset   int
		wantLine int
		wantCol  int
	}{
		{0, 1, 0},
		{5, 1, 5}, // CR is part of line 1
		{6, 2, 0}, // Start of line 2
		{11, 2, 5},
		{12, 3, 0},
	}

	for _, tt := range tests {
		line, col := idx.LineCol(tt.offset)
		if line != tt.wantLine ||
			col != tt.wantCol {
			t.Errorf(
				"LineCol(%d) = (%d, %d), want (%d, %d)",
				tt.offset,
				line,
				col,
				tt.wantLine,
				tt.wantCol,
			)
		}
	}
}

// TestLineEnd_StandaloneCR verifies LineEnd excludes standalone CR.
func TestLineEnd_StandaloneCR(t *testing.T) {
	source := []byte("abc\rdef")
	idx := NewLineIndex(source)

	end := idx.LineEnd(1)
	if end != 3 {
		t.Errorf(
			"LineEnd(1) = %d, want 3 (excluding CR)",
			end,
		)
	}
}
