package parsers

import (
	"testing"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantNil bool
	}{
		{
			name:    "simple markdown",
			content: "# Hello\n\nWorld",
			wantNil: false,
		},
		{
			name:    "empty content",
			content: "",
			wantNil: false, // goldmark returns empty document, not nil
		},
		{
			name:    "complex markdown",
			content: "# Title\n\n## Section\n\n- List item\n- Another item\n\n```go\ncode\n```",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := ParseMarkdown([]byte(tt.content))
			if tt.wantNil && node != nil {
				t.Errorf("ParseMarkdown() expected nil, got %v", node)
			}
			if !tt.wantNil && node == nil {
				t.Error("ParseMarkdown() returned nil, expected node")
			}
		})
	}
}

func TestFindHeading(t *testing.T) {
	tests := []struct {
		name      string
		markdown  string
		level     int
		text      string
		wantFound bool
		wantText  string
	}{
		{
			name:      "find H1",
			markdown:  "# Hello World\n\nSome content",
			level:     1,
			text:      "Hello World",
			wantFound: true,
			wantText:  "Hello World",
		},
		{
			name:      "find H2",
			markdown:  "# Title\n\n## Section One\n\n## Section Two",
			level:     2,
			text:      "Section One",
			wantFound: true,
			wantText:  "Section One",
		},
		{
			name:      "not found - wrong level",
			markdown:  "# Title\n\n## Section",
			level:     3,
			text:      "Section",
			wantFound: false,
		},
		{
			name:      "not found - wrong text",
			markdown:  "# Title\n\n## Section",
			level:     2,
			text:      "Different",
			wantFound: false,
		},
		{
			name:      "find with whitespace trimming",
			markdown:  "#   Title   \n\nContent",
			level:     1,
			text:      "Title",
			wantFound: true,
			wantText:  "Title",
		},
		{
			name:      "find delta section",
			markdown:  "# Spec\n\n## ADDED Requirements\n\n### Requirement: Test",
			level:     2,
			text:      "ADDED Requirements",
			wantFound: true,
			wantText:  "ADDED Requirements",
		},
		{
			name:      "multiple headings - find first",
			markdown:  "# Title\n\n## Section\n\nContent\n\n## Section\n\nMore",
			level:     2,
			text:      "Section",
			wantFound: true,
			wantText:  "Section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := ParseMarkdown([]byte(tt.markdown))
			if doc == nil {
				t.Fatal("ParseMarkdown returned nil")
			}

			result := FindHeading(doc, tt.level, tt.text)

			if tt.wantFound && result == nil {
				t.Errorf("FindHeading() expected to find heading, got nil")
				return
			}
			if !tt.wantFound && result != nil {
				t.Errorf("FindHeading() expected nil, found heading")
				return
			}

			if tt.wantFound {
				// Verify it's actually a heading
				heading, ok := result.(*ast.Heading)
				if !ok {
					t.Errorf("FindHeading() returned non-heading node")
					return
				}
				if heading.Level != tt.level {
					t.Errorf("FindHeading() level = %d, want %d", heading.Level, tt.level)
				}

				// Verify the text content
				text := ExtractTextContent(result)
				if text != tt.wantText {
					t.Errorf("FindHeading() text = %q, want %q", text, tt.wantText)
				}
			}
		})
	}
}

func TestSegmentToLocation(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		segmentStart   int
		segmentEnd     int
		wantLine       int
		wantColumn     int
		wantByteOffset int
	}{
		{
			name:           "first line, first character",
			source:         "Hello World",
			segmentStart:   0,
			segmentEnd:     5,
			wantLine:       1,
			wantColumn:     1,
			wantByteOffset: 0,
		},
		{
			name:           "first line, middle",
			source:         "Hello World",
			segmentStart:   6,
			segmentEnd:     11,
			wantLine:       1,
			wantColumn:     7,
			wantByteOffset: 6,
		},
		{
			name:           "second line, start",
			source:         "Line 1\nLine 2",
			segmentStart:   7,
			segmentEnd:     13,
			wantLine:       2,
			wantColumn:     1,
			wantByteOffset: 7,
		},
		{
			name:           "second line, middle",
			source:         "Line 1\nLine 2",
			segmentStart:   10,
			segmentEnd:     13,
			wantLine:       2,
			wantColumn:     4,
			wantByteOffset: 10,
		},
		{
			name:           "third line",
			source:         "Line 1\nLine 2\nLine 3",
			segmentStart:   14,
			segmentEnd:     20,
			wantLine:       3,
			wantColumn:     1,
			wantByteOffset: 14,
		},
		{
			name:           "CRLF line endings",
			source:         "Line 1\r\nLine 2\r\nLine 3",
			segmentStart:   8,
			segmentEnd:     14,
			wantLine:       2,
			wantColumn:     1,
			wantByteOffset: 8,
		},
		{
			name:           "empty source",
			source:         "",
			segmentStart:   0,
			segmentEnd:     0,
			wantLine:       1,
			wantColumn:     1,
			wantByteOffset: 0,
		},
		{
			name:           "heading on line 3",
			source:         "# Title\n\n## Section",
			segmentStart:   10,
			segmentEnd:     20,
			wantLine:       3,
			wantColumn:     2,
			wantByteOffset: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a text.Segment
			segment := text.NewSegment(tt.segmentStart, tt.segmentEnd)

			loc := SegmentToLocation([]byte(tt.source), segment)

			if loc.LineNumber != tt.wantLine {
				t.Errorf("SegmentToLocation() line = %d, want %d", loc.LineNumber, tt.wantLine)
			}
			if loc.Column != tt.wantColumn {
				t.Errorf("SegmentToLocation() column = %d, want %d", loc.Column, tt.wantColumn)
			}
			if loc.ByteOffset != tt.wantByteOffset {
				t.Errorf(
					"SegmentToLocation() byteOffset = %d, want %d",
					loc.ByteOffset,
					tt.wantByteOffset,
				)
			}
		})
	}
}

func TestExtractTextContent(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		nodeType string // "heading", "paragraph", "document"
		want     string
	}{
		{
			name:     "simple heading",
			markdown: "# Hello World",
			nodeType: "heading",
			want:     "Hello World",
		},
		{
			name:     "heading with inline code",
			markdown: "# Title with `code`",
			nodeType: "heading",
			want:     "Title with code",
		},
		{
			name:     "paragraph",
			markdown: "This is a paragraph.",
			nodeType: "paragraph",
			want:     "This is a paragraph.",
		},
		{
			name:     "paragraph with soft line break",
			markdown: "Line one\nLine two",
			nodeType: "paragraph",
			want:     "Line one Line two", // soft line break becomes space
		},
		{
			name:     "entire document",
			markdown: "# Title\n\nParagraph text.",
			nodeType: "document",
			want:     "Title Paragraph text.",
		},
		{
			name:     "requirement header",
			markdown: "### Requirement: User Authentication",
			nodeType: "heading",
			want:     "Requirement: User Authentication",
		},
		{
			name:     "scenario header",
			markdown: "#### Scenario: Login success",
			nodeType: "heading",
			want:     "Scenario: Login success",
		},
		{
			name:     "nil node",
			markdown: "",
			nodeType: "nil",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node ast.Node

			if tt.nodeType == "nil" {
				node = nil
			} else {
				doc := ParseMarkdown([]byte(tt.markdown))
				if doc == nil {
					t.Fatal("ParseMarkdown returned nil")
				}

				switch tt.nodeType {
				case "heading":
					// Find first heading
					node = findFirstNode(doc, func(n ast.Node) bool {
						_, ok := n.(*ast.Heading)
						return ok
					})
				case "paragraph":
					// Find first paragraph
					node = findFirstNode(doc, func(n ast.Node) bool {
						_, ok := n.(*ast.Paragraph)
						return ok
					})
				case "document":
					node = doc
				}
			}

			got := ExtractTextContent(node)
			if got != tt.want {
				t.Errorf("ExtractTextContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Helper types and functions for tests

// findFirstNode walks the AST and returns the first node matching the predicate
func findFirstNode(root ast.Node, predicate func(ast.Node) bool) ast.Node {
	var result ast.Node
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && predicate(n) {
			result = n
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return result
}

// Benchmark tests

func BenchmarkParseMarkdown(b *testing.B) {
	content := []byte(`# Specification

## Section 1

Some content here.

### Requirement: Test Feature

The system SHALL do something.

#### Scenario: Success case
- **WHEN** condition
- **THEN** result

## ADDED Requirements

### Requirement: New Feature

More content.
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseMarkdown(content)
	}
}

func BenchmarkFindHeading(b *testing.B) {
	content := []byte(`# Specification

## Section 1
## Section 2
## Section 3
## ADDED Requirements
## Section 5
`)

	doc := ParseMarkdown(content)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindHeading(doc, 2, "ADDED Requirements")
	}
}

func BenchmarkExtractTextContent(b *testing.B) {
	content := []byte("# This is a heading with some text content")
	doc := ParseMarkdown(content)
	heading := findFirstNode(doc, func(n ast.Node) bool {
		_, ok := n.(*ast.Heading)
		return ok
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ExtractTextContent(heading)
	}
}

func BenchmarkSegmentToLocation(b *testing.B) {
	source := []byte(`Line 1
Line 2
Line 3
Line 4
Line 5
`)
	segment := text.NewSegment(14, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SegmentToLocation(source, segment)
	}
}
