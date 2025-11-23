package parser

import (
	"strings"
	"testing"
)

// TestNodeType verifies NodeType string representation.
func TestNodeType(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		expected string
	}{
		{NodeDocument, "Document"},
		{NodeHeader, "Header"},
		{NodeParagraph, "Paragraph"},
		{NodeCodeBlock, "CodeBlock"},
		{NodeList, "List"},
		{NodeBlankLine, "BlankLine"},
		{NodeType(999), "Unknown(999)"},
	}

	for _, tt := range tests {
		got := tt.nodeType.String()
		if got != tt.expected {
			t.Errorf("NodeType.String() = %q, want %q", got, tt.expected)
		}
	}
}

// TestDocument verifies Document node creation and interface implementation.
func TestDocument(t *testing.T) {
	doc := NewDocument()

	// Test type
	if doc.Type() != NodeDocument {
		t.Errorf("Document.Type() = %v, want NodeDocument", doc.Type())
	}

	// Test position (should be 1:1)
	pos := doc.Pos()
	if pos.Line != 1 || pos.Column != 1 || pos.Offset != 0 {
		t.Errorf("Document.Pos() = %v, want Position{1, 1, 0}", pos)
	}

	// Test String representation
	str := doc.String()
	if !strings.Contains(str, "Document") || !strings.Contains(str, "0 nodes") {
		t.Errorf("Document.String() = %q, want to contain 'Document' and '0 nodes'", str)
	}

	// Test with children
	doc.Children = append(doc.Children, NewHeader(1, "Test", Position{Line: 1, Column: 1}))
	str = doc.String()
	if !strings.Contains(str, "1 nodes") {
		t.Errorf("Document.String() = %q, want to contain '1 nodes'", str)
	}
}

// TestHeader verifies Header node creation and interface implementation.
func TestHeader(t *testing.T) {
	pos := Position{Line: 5, Column: 3, Offset: 42}
	header := NewHeader(2, "Test Header", pos)

	// Test type
	if header.Type() != NodeHeader {
		t.Errorf("Header.Type() = %v, want NodeHeader", header.Type())
	}

	// Test position
	if header.Pos() != pos {
		t.Errorf("Header.Pos() = %v, want %v", header.Pos(), pos)
	}

	// Test fields
	if header.Level != 2 {
		t.Errorf("Header.Level = %d, want 2", header.Level)
	}
	if header.Text != "Test Header" {
		t.Errorf("Header.Text = %q, want %q", header.Text, "Test Header")
	}

	// Test String representation (short text)
	str := header.String()
	if !strings.Contains(str, "Header") || !strings.Contains(str, "Level: 2") ||
		!strings.Contains(str, "Test Header") || !strings.Contains(str, "5:3") {
		t.Errorf(
			"Header.String() = %q, want to contain 'Header', 'Level: 2', 'Test Header', and '5:3'",
			str,
		)
	}

	// Test String representation (long text - should truncate)
	longHeader := NewHeader(
		1,
		"This is a very long header text that should be truncated when displayed",
		pos,
	)
	str = longHeader.String()
	if !strings.Contains(str, "...") {
		t.Errorf("Header.String() with long text should contain '...', got %q", str)
	}
	if len(str) > 200 {
		t.Errorf("Header.String() too long: %d chars, expected truncation", len(str))
	}
}

// TestParagraph verifies Paragraph node creation and interface implementation.
func TestParagraph(t *testing.T) {
	pos := Position{Line: 10, Column: 1, Offset: 100}
	para := NewParagraph("This is a test paragraph.", pos)

	// Test type
	if para.Type() != NodeParagraph {
		t.Errorf("Paragraph.Type() = %v, want NodeParagraph", para.Type())
	}

	// Test position
	if para.Pos() != pos {
		t.Errorf("Paragraph.Pos() = %v, want %v", para.Pos(), pos)
	}

	// Test fields
	if para.Text != "This is a test paragraph." {
		t.Errorf("Paragraph.Text = %q, want %q", para.Text, "This is a test paragraph.")
	}

	// Test String representation (short text)
	str := para.String()
	if !strings.Contains(str, "Paragraph") || !strings.Contains(str, "test paragraph") ||
		!strings.Contains(str, "10:1") {
		t.Errorf(
			"Paragraph.String() = %q, want to contain 'Paragraph', 'test paragraph', and '10:1'",
			str,
		)
	}

	// Test String representation (long text - should truncate)
	longPara := NewParagraph(
		"This is a very long paragraph with lots of text that should be truncated",
		pos,
	)
	str = longPara.String()
	if !strings.Contains(str, "...") {
		t.Errorf("Paragraph.String() with long text should contain '...', got %q", str)
	}

	// Test with whitespace trimming in String
	paraTrim := NewParagraph("  \n  Text with whitespace  \n  ", pos)
	str = paraTrim.String()
	if strings.Contains(str, "\n") {
		t.Errorf("Paragraph.String() should trim whitespace, got %q", str)
	}
}

// TestCodeBlock verifies CodeBlock node creation and interface implementation.
func TestCodeBlock(t *testing.T) {
	pos := Position{Line: 15, Column: 1, Offset: 150}

	// Test without language
	code := NewCodeBlock("", "func main() {\n\tfmt.Println(\"Hello\")\n}", pos)

	// Test type
	if code.Type() != NodeCodeBlock {
		t.Errorf("CodeBlock.Type() = %v, want NodeCodeBlock", code.Type())
	}

	// Test position
	if code.Pos() != pos {
		t.Errorf("CodeBlock.Pos() = %v, want %v", code.Pos(), pos)
	}

	// Test fields
	if code.Language != "" {
		t.Errorf("CodeBlock.Language = %q, want empty string", code.Language)
	}
	if !strings.Contains(code.Content, "func main()") {
		t.Errorf("CodeBlock.Content = %q, want to contain 'func main()'", code.Content)
	}

	// Test String representation (no language)
	str := code.String()
	if !strings.Contains(str, "CodeBlock") || !strings.Contains(str, "15:1") {
		t.Errorf("CodeBlock.String() = %q, want to contain 'CodeBlock' and '15:1'", str)
	}
	if strings.Contains(str, "Lang:") {
		t.Errorf(
			"CodeBlock.String() should not contain 'Lang:' when language is empty, got %q",
			str,
		)
	}

	// Test with language
	codeWithLang := NewCodeBlock("go", "func test() {}", pos)
	str = codeWithLang.String()
	if !strings.Contains(str, "Lang: go") {
		t.Errorf("CodeBlock.String() with language should contain 'Lang: go', got %q", str)
	}

	// Test truncation for long content
	longContent := strings.Repeat("very long code content ", 20)
	longCode := NewCodeBlock("python", longContent, pos)
	str = longCode.String()
	if !strings.Contains(str, "...") {
		t.Errorf("CodeBlock.String() with long content should contain '...', got %q", str)
	}
}

// TestList verifies List node creation and interface implementation.
func TestList(t *testing.T) {
	pos := Position{Line: 20, Column: 1, Offset: 200}
	list := NewList("First item", pos)

	// Test type
	if list.Type() != NodeList {
		t.Errorf("List.Type() = %v, want NodeList", list.Type())
	}

	// Test position
	if list.Pos() != pos {
		t.Errorf("List.Pos() = %v, want %v", list.Pos(), pos)
	}

	// Test fields
	if len(list.Items) != 1 {
		t.Errorf("List.Items length = %d, want 1", len(list.Items))
	}
	if list.Items[0] != "First item" {
		t.Errorf("List.Items[0] = %q, want %q", list.Items[0], "First item")
	}

	// Test String representation (single item)
	str := list.String()
	if !strings.Contains(str, "List") || !strings.Contains(str, "First item") ||
		!strings.Contains(str, "20:1") {
		t.Errorf("List.String() = %q, want to contain 'List', 'First item', and '20:1'", str)
	}

	// Test empty list
	emptyList := &List{Items: make([]string, 0), position: pos}
	str = emptyList.String()
	if !strings.Contains(str, "Empty") {
		t.Errorf("Empty List.String() should contain 'Empty', got %q", str)
	}

	// Test multiple items
	multiList := &List{Items: []string{"Item 1", "Item 2", "Item 3"}, position: pos}
	str = multiList.String()
	if !strings.Contains(str, "Items: 3") {
		t.Errorf("Multi-item List.String() should contain 'Items: 3', got %q", str)
	}

	// Test truncation for long item
	longItem := "This is a very long list item that should be truncated in the string representation"
	longList := NewList(longItem, pos)
	str = longList.String()
	if !strings.Contains(str, "...") {
		t.Errorf("List.String() with long item should contain '...', got %q", str)
	}
}

// TestBlankLine verifies BlankLine node creation and interface implementation.
func TestBlankLine(t *testing.T) {
	pos := Position{Line: 25, Column: 1, Offset: 250}

	// Test single blank line
	blank := NewBlankLine(1, pos)

	// Test type
	if blank.Type() != NodeBlankLine {
		t.Errorf("BlankLine.Type() = %v, want NodeBlankLine", blank.Type())
	}

	// Test position
	if blank.Pos() != pos {
		t.Errorf("BlankLine.Pos() = %v, want %v", blank.Pos(), pos)
	}

	// Test fields
	if blank.Count != 1 {
		t.Errorf("BlankLine.Count = %d, want 1", blank.Count)
	}

	// Test String representation (single blank line)
	str := blank.String()
	if !strings.Contains(str, "BlankLine") || !strings.Contains(str, "25:1") {
		t.Errorf("BlankLine.String() = %q, want to contain 'BlankLine' and '25:1'", str)
	}
	// For single blank line, should not show count
	if strings.Contains(str, "Count:") {
		t.Errorf("Single BlankLine.String() should not contain 'Count:', got %q", str)
	}

	// Test multiple blank lines
	multiBlank := NewBlankLine(3, pos)
	str = multiBlank.String()
	if !strings.Contains(str, "Count: 3") {
		t.Errorf("Multi BlankLine.String() should contain 'Count: 3', got %q", str)
	}
}

// TestWalk verifies the Walk function for traversing the AST.
func TestWalk(t *testing.T) {
	// Build a simple document
	doc := NewDocument()
	doc.Children = append(doc.Children,
		NewHeader(1, "Title", Position{Line: 1, Column: 1}),
		NewParagraph("Text", Position{Line: 2, Column: 1}),
		NewCodeBlock("go", "code", Position{Line: 4, Column: 1}),
		NewBlankLine(1, Position{Line: 7, Column: 1}),
	)

	// Walk and count nodes
	nodeCount := 0
	Walk(doc, func(_ Node) bool {
		nodeCount++

		return true
	})

	// Should visit: Document + 4 children = 5 nodes
	expectedCount := 5
	if nodeCount != expectedCount {
		t.Errorf("Walk visited %d nodes, want %d", nodeCount, expectedCount)
	}

	// Test early termination
	visitedTypes := make([]NodeType, 0)
	Walk(doc, func(n Node) bool {
		visitedTypes = append(visitedTypes, n.Type())
		// Stop after visiting 3 nodes (Document + first 2 children)
		return len(visitedTypes) < 3
	})

	// Should visit: Document, Header, Paragraph (then stop)
	expectedVisitCount := 3
	if len(visitedTypes) != expectedVisitCount {
		t.Errorf(
			"Walk with early termination visited %d nodes, want %d",
			len(visitedTypes),
			expectedVisitCount,
		)
	}
	if len(visitedTypes) >= 3 {
		// Check the sequence
		expectedSequence := []NodeType{NodeDocument, NodeHeader, NodeParagraph}
		//nolint:intrange // Keep compatible with Go <1.22
		for i := 0; i < 3; i++ {
			if visitedTypes[i] != expectedSequence[i] {
				t.Errorf(
					"Walk visited node type %v at position %d, want %v",
					visitedTypes[i],
					i,
					expectedSequence[i],
				)
			}
		}
	}

	// Test Walk with nil node
	nilCount := 0
	Walk(nil, func(_ Node) bool {
		nilCount++

		return true
	})
	if nilCount != 0 {
		t.Errorf("Walk with nil node visited %d nodes, want 0", nilCount)
	}
}

// TestFindHeaders verifies the FindHeaders convenience function.
func TestFindHeaders(t *testing.T) {
	// Build a document with various headers
	doc := NewDocument()
	doc.Children = append(doc.Children,
		NewHeader(1, "Title", Position{Line: 1, Column: 1}),
		NewParagraph("Text", Position{Line: 2, Column: 1}),
		NewHeader(2, "ADDED Requirements", Position{Line: 4, Column: 1}),
		NewHeader(3, "Requirement: Feature", Position{Line: 6, Column: 1}),
		NewHeader(4, "Scenario: Test case", Position{Line: 8, Column: 1}),
	)

	// Find all level-2 headers
	level2Headers := FindHeaders(doc, func(h *Header) bool {
		return h.Level == 2
	})
	if len(level2Headers) != 1 {
		t.Errorf("FindHeaders found %d level-2 headers, want 1", len(level2Headers))
	}
	if len(level2Headers) > 0 && level2Headers[0].Text != "ADDED Requirements" {
		t.Errorf(
			"FindHeaders found header with text %q, want 'ADDED Requirements'",
			level2Headers[0].Text,
		)
	}

	// Find headers starting with "Requirement:"
	reqHeaders := FindHeaders(doc, func(h *Header) bool {
		return strings.HasPrefix(h.Text, "Requirement:")
	})
	if len(reqHeaders) != 1 {
		t.Errorf("FindHeaders found %d 'Requirement:' headers, want 1", len(reqHeaders))
	}

	// Find headers starting with "Scenario:"
	scenarioHeaders := FindHeaders(doc, func(h *Header) bool {
		return strings.HasPrefix(h.Text, "Scenario:")
	})
	if len(scenarioHeaders) != 1 {
		t.Errorf("FindHeaders found %d 'Scenario:' headers, want 1", len(scenarioHeaders))
	}

	// Find no matches
	noMatch := FindHeaders(doc, func(h *Header) bool {
		return h.Level == 6
	})
	if len(noMatch) != 0 {
		t.Errorf("FindHeaders found %d level-6 headers, want 0", len(noMatch))
	}

	// Test with empty document
	emptyDoc := NewDocument()
	emptyResults := FindHeaders(emptyDoc, func(_ *Header) bool {
		return true
	})
	if len(emptyResults) != 0 {
		t.Errorf("FindHeaders on empty document found %d headers, want 0", len(emptyResults))
	}
}

// TestNodesBetween verifies the NodesBetween function for extracting content ranges.
func TestNodesBetween(t *testing.T) {
	// Build a document
	doc := NewDocument()
	doc.Children = append(doc.Children,
		NewHeader(1, "Title", Position{Line: 1, Column: 1, Offset: 0}),
		NewParagraph("Intro", Position{Line: 2, Column: 1, Offset: 10}),
		NewHeader(2, "Section", Position{Line: 4, Column: 1, Offset: 20}),
		NewParagraph("Content", Position{Line: 5, Column: 1, Offset: 30}),
		NewCodeBlock("go", "code", Position{Line: 7, Column: 1, Offset: 40}),
		NewHeader(2, "Next Section", Position{Line: 10, Column: 1, Offset: 50}),
	)

	// Get nodes between "Section" header and "Next Section" header
	startPos := Position{Offset: 20}
	endPos := Position{Offset: 50}
	nodes := NodesBetween(doc, startPos, endPos)

	// Should get: Paragraph("Content") and CodeBlock
	expectedCount := 2
	if len(nodes) != expectedCount {
		t.Errorf("NodesBetween returned %d nodes, want %d", len(nodes), expectedCount)
	}

	if len(nodes) >= 2 {
		if nodes[0].Type() != NodeParagraph {
			t.Errorf("NodesBetween node[0] type = %v, want NodeParagraph", nodes[0].Type())
		}
		if nodes[1].Type() != NodeCodeBlock {
			t.Errorf("NodesBetween node[1] type = %v, want NodeCodeBlock", nodes[1].Type())
		}
	}

	// Test with no nodes in range
	emptyNodes := NodesBetween(doc, Position{Offset: 100}, Position{Offset: 200})
	if len(emptyNodes) != 0 {
		t.Errorf(
			"NodesBetween with out-of-range positions returned %d nodes, want 0",
			len(emptyNodes),
		)
	}

	// Test with empty document
	emptyDoc := NewDocument()
	emptyResults := NodesBetween(emptyDoc, Position{Offset: 0}, Position{Offset: 100})
	if len(emptyResults) != 0 {
		t.Errorf("NodesBetween on empty document returned %d nodes, want 0", len(emptyResults))
	}
}

// TestNodeInterfaceImplementation verifies all nodes implement the Node interface correctly.
func TestNodeInterfaceImplementation(t *testing.T) {
	pos := Position{Line: 1, Column: 1, Offset: 0}

	nodes := []Node{
		NewDocument(),
		NewHeader(1, "Test", pos),
		NewParagraph("Test", pos),
		NewCodeBlock("", "Test", pos),
		NewList("Test", pos),
		NewBlankLine(1, pos),
	}

	for _, node := range nodes {
		// Test Type method
		nodeType := node.Type()
		if nodeType < NodeDocument || nodeType > NodeBlankLine {
			t.Errorf("Node %T has invalid Type() = %v", node, nodeType)
		}

		// Test Pos method
		nodePos := node.Pos()
		if nodePos.Line < 0 || nodePos.Column < 0 || nodePos.Offset < 0 {
			t.Errorf("Node %T has invalid Pos() = %v", node, nodePos)
		}

		// Test String method
		str := node.String()
		if str == "" {
			t.Errorf("Node %T has empty String() output", node)
		}
	}
}
