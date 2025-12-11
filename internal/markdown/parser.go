package markdown

import (
	"os"

	bf "github.com/russross/blackfriday/v2"
)

// Parse parses markdown content and returns the root AST node.
// The returned node can be traversed using Node.Walk() or Node.FirstChild/Next.
func Parse(content []byte) *bf.Node {
	// Use common extensions but disable HTML parsing for spec files
	extensions := bf.CommonExtensions | bf.NoIntraEmphasis
	parser := bf.New(bf.WithExtensions(extensions))

	return parser.Parse(content)
}

// ParseFile reads a file and parses its markdown content.
// Returns the root AST node or an error if the file cannot be read.
func ParseFile(path string) (*bf.Node, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(content), nil
}

// extractText recursively extracts text content from a node and its
// children. This is useful for getting the text content of headers
// and other inline elements.
func extractText(node *bf.Node) string {
	if node == nil {
		return ""
	}

	var result string
	for child := node.FirstChild; child != nil; child = child.Next {
		//nolint:exhaustive // default case handles all other node types
		switch child.Type {
		case bf.Text:
			result += string(child.Literal)
		case bf.Code:
			result += string(child.Literal)
		case bf.Softbreak, bf.Hardbreak:
			result += " "
		default:
			// Recursively extract text from inline elements
			result += extractText(child)
		}
	}

	return result
}

// nodeLineNumber attempts to get the line number for a node.
// Returns 0 if line number is not available.
//
//nolint:revive // unused-parameter - kept for API compatibility
func nodeLineNumber(_ *bf.Node) int {
	// Blackfriday doesn't track line numbers in the AST directly,
	// so we return 0 for now. This could be enhanced by tracking
	// line numbers during parsing if needed.
	return 0
}
