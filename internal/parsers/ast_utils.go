// Package parsers provides utilities for extracting and counting
// information from markdown specification files.
package parsers

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

// SourceLocation represents a position in the source markdown file
// with line number, column, and byte offset information.
type SourceLocation struct {
	LineNumber int // 1-indexed line number
	Column     int // 1-indexed column number
	ByteOffset int // 0-indexed byte offset from start of file
}

// ParseMarkdown parses markdown content and returns the root AST node.
// The source bytes are stored in the document's metadata under the "source" key
// for later use by functions like ExtractTextContent.
//
// Uses goldmark with GFM (GitHub Flavored Markdown) extensions to support
// task lists, tables, and other common markdown features.
//
// Parameters:
//   - content: The markdown content to parse
//
// Returns the parsed document AST node with source stored in metadata.
func ParseMarkdown(content []byte) ast.Node {
	parser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	doc := parser.Parser().Parse(text.NewReader(content))

	// Store source bytes in document metadata for later text extraction
	// This follows goldmark's pattern of passing source to segment.Value()
	if doc != nil {
		doc.(*ast.Document).AddMeta("source", content)
	}

	return doc
}

// FindHeading searches the AST for a heading node with the specified level
// and text content. The text comparison is case-sensitive and whitespace-sensitive.
//
// Parameters:
//   - node: Root AST node to start search from
//   - level: Heading level (1 for H1, 2 for H2, etc.)
//   - targetText: Text content to match (e.g., "Requirements" for "## Requirements")
//
// Returns the first matching heading node, or nil if not found.
func FindHeading(node ast.Node, level int, targetText string) ast.Node {
	var result ast.Node

	// Walk the AST tree
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		// Only process nodes when entering (not when leaving)
		if !entering {
			return ast.WalkContinue, nil
		}

		// Check if this is a heading with the right level
		heading, ok := n.(*ast.Heading)
		if !ok || heading.Level != level {
			return ast.WalkContinue, nil
		}

		// Extract the heading text
		headingText := ExtractTextContent(heading)

		// Check if the text matches (after trimming whitespace)
		if strings.TrimSpace(headingText) == strings.TrimSpace(targetText) {
			result = heading
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	return result
}

// SegmentToLocation converts a goldmark text.Segment (byte offsets) to a
// SourceLocation with line number, column, and byte offset.
//
// Line numbers and columns are 1-indexed for human readability.
// Handles both LF (\n) and CRLF (\r\n) line endings correctly.
//
// Parameters:
//   - source: The original markdown content as bytes
//   - segment: The text segment from a goldmark AST node
//
// Returns a SourceLocation with line, column, and byte offset.
func SegmentToLocation(source []byte, segment text.Segment) SourceLocation {
	return SourceLocation{
		LineNumber: segmentToLineNumber(source, segment),
		Column:     segmentToColumn(source, segment),
		ByteOffset: segment.Start,
	}
}

// ExtractTextContent extracts all text content from an AST node and its children.
// This traverses the entire subtree and concatenates text from all Text nodes.
//
// Useful for extracting heading text, paragraph content, or any text within
// a node hierarchy. Preserves spaces and newlines as they appear in the source.
//
// When extracting from containers with multiple block-level children (like extracting
// from an entire document), spaces are added between blocks for readability.
//
// The source bytes are retrieved from the document's metadata (stored by ParseMarkdown).
// If no source is available, returns an empty string.
//
// Parameters:
//   - node: The AST node to extract text from
//
// Returns the concatenated text content as a string.
func ExtractTextContent(node ast.Node) string {
	if node == nil {
		return ""
	}

	var buf bytes.Buffer

	// Get source from document metadata
	var source []byte
	doc := node.OwnerDocument()
	if doc != nil {
		if meta := doc.Meta(); meta != nil {
			if src, ok := meta["source"].([]byte); ok {
				source = src
			}
		}
	}

	var hasSeenText bool

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		// Only process when entering nodes
		if !entering {
			return ast.WalkContinue, nil
		}

		// When we encounter a block element (heading, paragraph, etc.) after already
		// having extracted text, add a space for separation
		if n.Type() == ast.TypeBlock && n != node && hasSeenText && buf.Len() > 0 {
			buf.WriteByte(' ')
			hasSeenText = false // Reset for next block's text
		}

		// Extract text from Text nodes using goldmark's segment.Value(source) pattern
		if textNode, ok := n.(*ast.Text); ok {
			if len(source) > 0 {
				buf.Write(textNode.Segment.Value(source))
			}
			// Add soft line breaks as spaces
			if textNode.SoftLineBreak() {
				buf.WriteByte(' ')
			}
			hasSeenText = true
		}

		return ast.WalkContinue, nil
	})

	return buf.String()
}

// segmentToLineNumber calculates the 1-indexed line number for a text segment.
// Counts newlines from the start of the source to segment.Start.
// Handles both LF and CRLF line endings.
func segmentToLineNumber(source []byte, segment text.Segment) int {
	// Handle edge case: empty source or invalid segment
	if len(source) == 0 || segment.Start < 0 {
		return 1
	}

	// Clamp segment.Start to source length
	start := segment.Start
	if start > len(source) {
		start = len(source)
	}

	// Count newlines before the segment start
	// Line numbers are 1-indexed, so start at 1
	return bytes.Count(source[:start], []byte("\n")) + 1
}

// segmentToColumn calculates the 1-indexed column number for a text segment.
// Finds the last newline before segment.Start and calculates the offset.
// Handles both LF and CRLF line endings correctly.
func segmentToColumn(source []byte, segment text.Segment) int {
	// Handle edge case: empty source or invalid segment
	if len(source) == 0 || segment.Start <= 0 {
		return 1
	}

	// Clamp segment.Start to source length
	start := segment.Start
	if start > len(source) {
		start = len(source)
	}

	// Find the last newline before segment.Start
	lastNewline := bytes.LastIndex(source[:start], []byte("\n"))

	// If no newline found, we're on the first line
	if lastNewline == -1 {
		// Columns are 1-indexed, so add 1
		return start + 1
	}

	// Calculate offset from last newline
	// Columns are 1-indexed, so the character after \n is column 1
	return start - lastNewline
}
