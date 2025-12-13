// Package markdown provides AST-based markdown parsing using blackfriday.
package markdown

import (
	"github.com/russross/blackfriday/v2"
)

// extractHeaders walks the AST and extracts all H1-H4 headers.
// Returns headers in document order with 1-indexed line numbers.
func extractHeaders(node *blackfriday.Node, lineIndex *lineIndex) []Header {
	var headers []Header

	node.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		// Only process headings when entering the node
		if !entering || n.Type != blackfriday.Heading {
			return blackfriday.GoToNext
		}

		level := n.Level
		// Only extract H1-H4 as per spec
		if level < 1 || level > 4 {
			return blackfriday.GoToNext
		}

		// Extract header text from children
		text := extractNodeText(n)

		// Get line number from byte offset
		// Blackfriday doesn't expose source positions,
		// so we search for the header in source
		line := 1
		if lineIndex != nil {
			offset := findNodeOffset(n, lineIndex.source)
			line = lineIndex.lineAt(offset)
		}

		headers = append(headers, Header{
			Level: level,
			Text:  text,
			Line:  line,
		})

		return blackfriday.GoToNext
	})

	return headers
}

// extractNodeText extracts all text content from a node and its children.
func extractNodeText(node *blackfriday.Node) string {
	var text []byte

	node.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.GoToNext
		}

		if n.Type == blackfriday.Text || n.Type == blackfriday.Code {
			text = append(text, n.Literal...)
		}

		return blackfriday.GoToNext
	})

	return string(text)
}

// findNodeOffset attempts to find the byte offset of a node in the source.
// This is needed because blackfriday doesn't directly expose source positions.
func findNodeOffset(node *blackfriday.Node, source []byte) int {
	if node.Type != blackfriday.Heading {
		return 0
	}

	text := extractNodeText(node)
	if len(text) == 0 {
		return 0
	}

	return findHeaderOffset(source, text, node.Level)
}

// findHeaderOffset searches for a header pattern in source.
func findHeaderOffset(source []byte, text string, level int) int {
	for i := range source {
		if source[i] != '#' {
			continue
		}

		offset := matchHeaderAtPosition(source, i, text, level)
		if offset >= 0 {
			return offset
		}
	}

	return 0
}

// matchHeaderAtPosition checks if there's a matching header at position i.
// Returns i if matched, -1 otherwise.
func matchHeaderAtPosition(source []byte, i int, text string, level int) int {
	hashCount, j := countHashes(source, i)
	if hashCount != level {
		return -1
	}

	j = skipWhitespace(source, j)
	if !headerTextMatches(source, j, text) {
		return -1
	}

	return i
}

// countHashes counts consecutive # characters starting at position i.
func countHashes(source []byte, i int) (count, endPos int) {
	j := i
	for j < len(source) && source[j] == '#' {
		count++
		j++
	}

	return count, j
}

// skipWhitespace advances past spaces and tabs.
func skipWhitespace(source []byte, start int) int {
	pos := start
	for pos < len(source) && (source[pos] == ' ' || source[pos] == '\t') {
		pos++
	}

	return pos
}

// headerTextMatches checks if the text at position j matches the header text.
func headerTextMatches(source []byte, j int, text string) bool {
	return j+len(text) <= len(source) && string(source[j:j+len(text)]) == text
}
