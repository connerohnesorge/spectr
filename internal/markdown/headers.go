//nolint:revive // file-length-limit - header extraction requires many helpers
package markdown

import (
	"strings"

	bf "github.com/russross/blackfriday/v2"
)

// ExtractHeaders extracts all headers from the markdown AST.
// Returns headers in document order with their level (1-6) and text.
func ExtractHeaders(node *bf.Node) []Header {
	var headers []Header
	if node == nil {
		return headers
	}

	node.Walk(func(n *bf.Node, entering bool) bf.WalkStatus {
		if !entering || n.Type != bf.Heading {
			return bf.GoToNext
		}

		headers = append(headers, Header{
			Level: n.Level,
			Text:  strings.TrimSpace(extractText(n)),
			Line:  nodeLineNumber(n),
		})

		return bf.GoToNext
	})

	return headers
}

// ExtractH1Title extracts the first H1 title from the document.
// Returns empty string if no H1 is found.
// This is compatible with the existing ExtractTitle behavior.
func ExtractH1Title(node *bf.Node) string {
	if node == nil {
		return ""
	}

	var title string
	node.Walk(func(n *bf.Node, entering bool) bf.WalkStatus {
		if !entering || n.Type != bf.Heading {
			return bf.GoToNext
		}

		if n.Level == 1 {
			title = strings.TrimSpace(extractText(n))

			return bf.Terminate // Stop after first H1
		}

		return bf.GoToNext
	})

	return title
}

// ExtractH1TitleClean extracts the first H1 title and removes common
// prefixes. Removes "Change:" and "Spec:" prefixes, matching existing
// ExtractTitle behavior.
func ExtractH1TitleClean(node *bf.Node) string {
	title := ExtractH1Title(node)
	if title == "" {
		return ""
	}

	// Remove "Change:" or "Spec:" prefix (matching parsers.ExtractTitle)
	title = strings.TrimPrefix(title, "Change:")
	title = strings.TrimPrefix(title, "Spec:")

	return strings.TrimSpace(title)
}

// ExtractH2Sections extracts all H2 sections and returns a map of header
// text to content. Content includes everything from after the H2 header
// until the next H2 or H1 header.
func ExtractH2Sections(node *bf.Node) map[string]string {
	return extractSectionsAtLevel(node, 2)
}

// extractSectionsAtLevel extracts sections at a specific header level.
// Returns a map of header text to section content.
//
//nolint:revive // cognitive-complexity - AST walking is inherently complex
func extractSectionsAtLevel(node *bf.Node, level int) map[string]string {
	sections := make(map[string]string)
	if node == nil {
		return sections
	}

	// We need to walk through document nodes and track content between headers
	var currentHeader string
	var contentBuilder strings.Builder
	var inSection bool

	// Get all document children (block-level nodes)
	for child := node.FirstChild; child != nil; child = child.Next {
		if child.Type == bf.Heading {
			headingLevel := child.Level

			// If we hit a header at our target level, save previous and start new section
			if headingLevel == level {
				if inSection && currentHeader != "" {
					sections[currentHeader] = strings.TrimSpace(contentBuilder.String())
				}
				currentHeader = strings.TrimSpace(extractText(child))
				contentBuilder.Reset()
				inSection = true

				continue
			}

			// If we hit a higher-level header (smaller number), end current section
			if headingLevel < level {
				if inSection && currentHeader != "" {
					sections[currentHeader] = strings.TrimSpace(contentBuilder.String())
				}
				currentHeader = ""
				inSection = false

				continue
			}
		}

		// If we're in a section, append the content
		if inSection {
			contentBuilder.WriteString(renderNode(child))
		}
	}

	// Save final section
	if inSection && currentHeader != "" {
		sections[currentHeader] = strings.TrimSpace(contentBuilder.String())
	}

	return sections
}

// FindDeltaSection finds a delta section header.
// Matches ADDED/MODIFIED/REMOVED/RENAMED Requirements headers.
// Returns the heading node if found, nil otherwise.
//
//nolint:revive // modifies-parameter is intentional for normalization
func FindDeltaSection(node *bf.Node, sectionType string) *bf.Node {
	if node == nil {
		return nil
	}

	// Normalize section type
	sectionType = strings.ToUpper(strings.TrimSpace(sectionType))

	var result *bf.Node
	node.Walk(func(n *bf.Node, entering bool) bf.WalkStatus {
		if !entering || n.Type != bf.Heading {
			return bf.GoToNext
		}

		// Only look at H2 headers
		if n.Level != 2 {
			return bf.GoToNext
		}

		headerText := strings.TrimSpace(extractText(n))
		// Match pattern: "ADDED Requirements", "MODIFIED Requirements", etc.
		expectedPrefix := sectionType + " Requirements"
		if strings.HasPrefix(headerText, expectedPrefix) {
			result = n

			return bf.Terminate
		}

		return bf.GoToNext
	})

	return result
}

// FindAllDeltaSections finds all delta section headers in the document.
// Returns a map of delta type to the header node.
func FindAllDeltaSections(node *bf.Node) map[DeltaType]*bf.Node {
	sections := make(map[DeltaType]*bf.Node)
	if node == nil {
		return sections
	}

	node.Walk(func(n *bf.Node, entering bool) bf.WalkStatus {
		if !entering || n.Type != bf.Heading {
			return bf.GoToNext
		}

		if n.Level != 2 {
			return bf.GoToNext
		}

		headerText := strings.TrimSpace(extractText(n))
		for _, dt := range ValidDeltaTypes() {
			if strings.HasPrefix(headerText, dt+" Requirements") {
				sections[DeltaType(dt)] = n

				break
			}
		}

		return bf.GoToNext
	})

	return sections
}

// renderNode renders a single AST node back to markdown-like text.
// This is a simplified renderer for extracting section content.
//
//nolint:revive // cognitive-complexity, function-length - switch per node type
func renderNode(node *bf.Node) string {
	if node == nil {
		return ""
	}

	var result strings.Builder

	//nolint:exhaustive // default case handles all other node types
	switch node.Type {
	case bf.Document:
		for child := node.FirstChild; child != nil; child = child.Next {
			result.WriteString(renderNode(child))
		}

	case bf.Heading:
		level := node.Level
		result.WriteString(strings.Repeat("#", level))
		result.WriteString(" ")
		result.WriteString(renderInlineMarkdown(node))
		result.WriteString("\n")

	case bf.Paragraph:
		result.WriteString(renderInlineMarkdown(node))
		result.WriteString("\n\n")

	case bf.List:
		for item := node.FirstChild; item != nil; item = item.Next {
			result.WriteString(renderNode(item))
		}
		result.WriteString("\n") // Add blank line after list for proper markdown parsing

	case bf.Item:
		// Check if this is a task list item
		result.WriteString("- ")
		for child := node.FirstChild; child != nil; child = child.Next {
			if child.Type == bf.Paragraph {
				result.WriteString(renderInlineMarkdown(child))
			} else {
				result.WriteString(renderNode(child))
			}
		}
		result.WriteString("\n")

	case bf.CodeBlock:
		result.WriteString("```")
		if len(node.Info) > 0 {
			result.WriteString(string(node.Info))
		}
		result.WriteString("\n")
		result.WriteString(string(node.Literal))
		result.WriteString("```\n\n")

	case bf.BlockQuote:
		for child := node.FirstChild; child != nil; child = child.Next {
			lines := strings.Split(renderNode(child), "\n")
			for _, line := range lines {
				if line != "" {
					result.WriteString("> ")
					result.WriteString(line)
					result.WriteString("\n")
				}
			}
		}

	case bf.Text:
		result.WriteString(string(node.Literal))

	case bf.Softbreak, bf.Hardbreak:
		result.WriteString("\n")

	default:
		// For other node types, try to extract text content
		result.WriteString(extractText(node))
	}

	return result.String()
}

// renderInlineMarkdown renders inline content preserving markdown formatting
// like bold (**text**) and emphasis (*text*).
func renderInlineMarkdown(node *bf.Node) string {
	if node == nil {
		return ""
	}

	var result strings.Builder
	for child := node.FirstChild; child != nil; child = child.Next {
		//nolint:exhaustive // default case handles all other node types
		switch child.Type {
		case bf.Text:
			result.WriteString(string(child.Literal))
		case bf.Code:
			result.WriteString("`")
			result.WriteString(string(child.Literal))
			result.WriteString("`")
		case bf.Strong:
			result.WriteString("**")
			result.WriteString(renderInlineMarkdown(child))
			result.WriteString("**")
		case bf.Emph:
			result.WriteString("*")
			result.WriteString(renderInlineMarkdown(child))
			result.WriteString("*")
		case bf.Link:
			result.WriteString("[")
			result.WriteString(renderInlineMarkdown(child))
			result.WriteString("](")
			result.WriteString(string(child.Destination))
			result.WriteString(")")
		case bf.Softbreak:
			result.WriteString(" ")
		case bf.Hardbreak:
			result.WriteString("\n")
		default:
			// Recursively render for any other inline type
			result.WriteString(renderInlineMarkdown(child))
		}
	}

	return result.String()
}
