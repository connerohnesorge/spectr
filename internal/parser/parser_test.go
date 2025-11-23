package parser

import (
	"strings"
	"testing"
)

// TestParseTokens_EmptyInput tests parsing an empty token stream.
func TestParseTokens_EmptyInput(t *testing.T) {
	doc, err := ParseTokens(make([]Token, 0))
	if err != nil {
		t.Fatalf("ParseTokens failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	if len(doc.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(doc.Children))
	}
}

// TestParseTokens_SingleHeader tests parsing a single header token.
func TestParseTokens_SingleHeader(t *testing.T) {
	tokens := []Token{
		{Type: TokenHeader, Value: "### Requirement: Test", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 1, Column: 24}},
	}

	doc, err := ParseTokens(tokens)
	if err != nil {
		t.Fatalf("ParseTokens failed: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(doc.Children))
	}

	header, ok := doc.Children[0].(*Header)
	if !ok {
		t.Fatalf("Expected Header node, got %T", doc.Children[0])
	}

	if header.Level != 3 {
		t.Errorf("Expected level 3, got %d", header.Level)
	}

	expectedText := "Requirement: Test"
	if header.Text != expectedText {
		t.Errorf("Expected text %q, got %q", expectedText, header.Text)
	}

	if header.Pos().Line != 1 || header.Pos().Column != 1 {
		t.Errorf("Expected position 1:1, got %s", header.Pos())
	}
}

// TestParseTokens_SingleParagraph tests parsing a single text token.
func TestParseTokens_SingleParagraph(t *testing.T) {
	tokens := []Token{
		{Type: TokenText, Value: "This is some text.\n", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 2, Column: 1}},
	}

	doc, err := ParseTokens(tokens)
	if err != nil {
		t.Fatalf("ParseTokens failed: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(doc.Children))
	}

	para, ok := doc.Children[0].(*Paragraph)
	if !ok {
		t.Fatalf("Expected Paragraph node, got %T", doc.Children[0])
	}

	expectedText := "This is some text.\n"
	if para.Text != expectedText {
		t.Errorf("Expected text %q, got %q", expectedText, para.Text)
	}
}

// TestParseTokens_CodeBlock tests parsing a code block token.
func TestParseTokens_CodeBlock(t *testing.T) {
	// Simulate a code block token as the lexer would produce it
	codeBlockValue := "```go\nfunc main() {\n}\n```"

	tokens := []Token{
		{Type: TokenCodeBlock, Value: codeBlockValue, Pos: Position{Line: 1, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 5, Column: 1}},
	}

	doc, err := ParseTokens(tokens)
	if err != nil {
		t.Fatalf("ParseTokens failed: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(doc.Children))
	}

	code, ok := doc.Children[0].(*CodeBlock)
	if !ok {
		t.Fatalf("Expected CodeBlock node, got %T", doc.Children[0])
	}

	if code.Language != "go" {
		t.Errorf("Expected language 'go', got %q", code.Language)
	}

	expectedContent := "func main() {\n}"
	if code.Content != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, code.Content)
	}
}

// TestParseTokens_ListItem tests parsing a list item token.
func TestParseTokens_ListItem(t *testing.T) {
	tokens := []Token{
		{Type: TokenListItem, Value: "- First item", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 2, Column: 1}},
	}

	doc, err := ParseTokens(tokens)
	if err != nil {
		t.Fatalf("ParseTokens failed: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(doc.Children))
	}

	list, ok := doc.Children[0].(*List)
	if !ok {
		t.Fatalf("Expected List node, got %T", doc.Children[0])
	}

	if len(list.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(list.Items))
	}

	expectedText := "First item"
	if list.Items[0] != expectedText {
		t.Errorf("Expected item %q, got %q", expectedText, list.Items[0])
	}
}

// TestParseTokens_BlankLine tests parsing a blank line token.
func TestParseTokens_BlankLine(t *testing.T) {
	tokens := []Token{
		{Type: TokenBlankLine, Value: "\n\n", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 3, Column: 1}},
	}

	doc, err := ParseTokens(tokens)
	if err != nil {
		t.Fatalf("ParseTokens failed: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(doc.Children))
	}

	blank, ok := doc.Children[0].(*BlankLine)
	if !ok {
		t.Fatalf("Expected BlankLine node, got %T", doc.Children[0])
	}

	if blank.Count != 2 {
		t.Errorf("Expected count 2, got %d", blank.Count)
	}
}

// TestParseTokens_MixedContent tests parsing various token types together.
func TestParseTokens_MixedContent(t *testing.T) {
	tokens := []Token{
		{Type: TokenHeader, Value: "## Title", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenBlankLine, Value: "\n", Pos: Position{Line: 2, Column: 1}},
		{Type: TokenText, Value: "Some text.\n", Pos: Position{Line: 3, Column: 1}},
		{Type: TokenListItem, Value: "- Item", Pos: Position{Line: 4, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 5, Column: 1}},
	}

	doc, err := ParseTokens(tokens)
	if err != nil {
		t.Fatalf("ParseTokens failed: %v", err)
	}

	if len(doc.Children) != 4 {
		t.Fatalf("Expected 4 children, got %d", len(doc.Children))
	}

	// Verify node types in order
	expectedTypes := []NodeType{NodeHeader, NodeBlankLine, NodeParagraph, NodeList}
	for i, expectedType := range expectedTypes {
		if doc.Children[i].Type() != expectedType {
			t.Errorf("Child %d: expected type %s, got %s",
				i, expectedType, doc.Children[i].Type())
		}
	}
}

// TestParseTokens_ErrorToken tests that error tokens are reported.
func TestParseTokens_ErrorToken(t *testing.T) {
	tokens := []Token{
		{Type: TokenText, Value: "Normal text\n", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenError, Value: "unclosed code block", Pos: Position{Line: 2, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 2, Column: 1}},
	}

	_, err := ParseTokens(tokens)
	if err == nil {
		t.Fatal("Expected error for TokenError, got nil")
	}

	if !strings.Contains(err.Error(), "lexer error") {
		t.Errorf("Expected error message to contain 'lexer error', got: %v", err)
	}

	if !strings.Contains(err.Error(), "unclosed code block") {
		t.Errorf("Expected error message to contain error text, got: %v", err)
	}
}

// TestParse_EndToEnd tests the convenience Parse() function end-to-end.
func TestParse_EndToEnd(t *testing.T) {
	input := `## Header

Some text here.

- List item 1
- List item 2
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Should have: header, blank, paragraph, blank, list, list
	if len(doc.Children) < 3 {
		t.Fatalf("Expected at least 3 children, got %d", len(doc.Children))
	}

	// First should be header
	if doc.Children[0].Type() != NodeHeader {
		t.Errorf("Expected first child to be Header, got %s", doc.Children[0].Type())
	}
}

// TestParse_CriticalCodeBlockCase tests the critical case where headers
// appear inside code blocks and should not be parsed as headers.
func TestParse_CriticalCodeBlockCase(t *testing.T) {
	input := "```go\n### Requirement: This is inside a code block\n```\n### Requirement: This is a real requirement\n"

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Should have exactly 2 children: CodeBlock and Header
	if len(doc.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(doc.Children))
	}

	// First child should be CodeBlock
	code, ok := doc.Children[0].(*CodeBlock)
	if !ok {
		t.Fatalf("Expected first child to be CodeBlock, got %T", doc.Children[0])
	}

	// Code block should contain the fake requirement
	if !strings.Contains(code.Content, "### Requirement: This is inside a code block") {
		t.Errorf("Code block should contain fake requirement, got: %q", code.Content)
	}

	// Second child should be Header
	header, ok := doc.Children[1].(*Header)
	if !ok {
		t.Fatalf("Expected second child to be Header, got %T", doc.Children[1])
	}

	// Header should be the real requirement
	expectedText := "Requirement: This is a real requirement"
	if header.Text != expectedText {
		t.Errorf("Expected header text %q, got %q", expectedText, header.Text)
	}

	if header.Level != 3 {
		t.Errorf("Expected header level 3, got %d", header.Level)
	}
}

// TestParse_CodeBlockWithoutLanguage tests code blocks without language specifier.
func TestParse_CodeBlockWithoutLanguage(t *testing.T) {
	input := "```\nplain code\n```\n"

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(doc.Children))
	}

	code, ok := doc.Children[0].(*CodeBlock)
	if !ok {
		t.Fatalf("Expected CodeBlock, got %T", doc.Children[0])
	}

	if code.Language != "" {
		t.Errorf("Expected empty language, got %q", code.Language)
	}

	if code.Content != "plain code" {
		t.Errorf("Expected content 'plain code', got %q", code.Content)
	}
}

// TestParse_MultipleHeaders tests parsing multiple headers with different levels.
func TestParse_MultipleHeaders(t *testing.T) {
	input := `# Level 1
## Level 2
### Level 3
#### Level 4
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Count headers
	headerCount := 0
	for _, child := range doc.Children {
		if _, ok := child.(*Header); ok {
			headerCount++
		}
	}

	if headerCount != 4 {
		t.Fatalf("Expected 4 headers, got %d", headerCount)
	}

	// Verify levels
	expectedLevels := []int{1, 2, 3, 4}
	headerIdx := 0
	for _, child := range doc.Children {
		if header, ok := child.(*Header); ok {
			if header.Level != expectedLevels[headerIdx] {
				t.Errorf("Header %d: expected level %d, got %d",
					headerIdx, expectedLevels[headerIdx], header.Level)
			}
			headerIdx++
		}
	}
}

// TestParse_PreservesOrder tests that nodes appear in source order.
func TestParse_PreservesOrder(t *testing.T) {
	input := `First paragraph
## Header
Second paragraph
- List item
Third paragraph
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Expected order: Paragraph, Header, Paragraph, List, Paragraph
	expectedTypes := []NodeType{NodeParagraph, NodeHeader, NodeParagraph, NodeList, NodeParagraph}

	actualTypes := make([]NodeType, 0)
	for _, child := range doc.Children {
		actualTypes = append(actualTypes, child.Type())
	}

	// Filter out blank lines for comparison
	filteredTypes := make([]NodeType, 0)
	for _, typ := range actualTypes {
		if typ != NodeBlankLine {
			filteredTypes = append(filteredTypes, typ)
		}
	}

	if len(filteredTypes) != len(expectedTypes) {
		t.Fatalf("Expected %d non-blank nodes, got %d", len(expectedTypes), len(filteredTypes))
	}

	for i, expected := range expectedTypes {
		if filteredTypes[i] != expected {
			t.Errorf("Position %d: expected %s, got %s", i, expected, filteredTypes[i])
		}
	}
}

// TestParse_PositionPreservation tests that position information is preserved.
func TestParse_PositionPreservation(t *testing.T) {
	input := `Line 1
Line 2
## Header on line 3
Line 4
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find the header
	var header *Header
	for _, child := range doc.Children {
		if h, ok := child.(*Header); ok {
			header = h

			break
		}
	}

	if header == nil {
		t.Fatal("Expected to find header node")
	}

	// Header should be on line 3
	if header.Pos().Line != 3 {
		t.Errorf("Expected header on line 3, got line %d", header.Pos().Line)
	}
}

// TestParser_Helpers tests the parser helper methods.
func TestParser_Helpers(t *testing.T) {
	tokens := []Token{
		{Type: TokenText, Value: "text", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenHeader, Value: "# header", Pos: Position{Line: 2, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 3, Column: 1}},
	}

	p := NewParser(tokens)

	// Test current()
	if p.current().Type != TokenText {
		t.Errorf("Expected current to be TokenText, got %s", p.current().Type)
	}

	// Test peek()
	if p.peek().Type != TokenHeader {
		t.Errorf("Expected peek to be TokenHeader, got %s", p.peek().Type)
	}

	// Current should still be TokenText
	if p.current().Type != TokenText {
		t.Errorf("Peek should not advance, current is %s", p.current().Type)
	}

	// Test advance()
	p.advance()
	if p.current().Type != TokenHeader {
		t.Errorf("After advance, expected TokenHeader, got %s", p.current().Type)
	}

	// Test consume()
	token := p.consume()
	if token.Type != TokenHeader {
		t.Errorf("Consume should return TokenHeader, got %s", token.Type)
	}
	if p.current().Type != TokenEOF {
		t.Errorf("After consume, expected TokenEOF, got %s", p.current().Type)
	}

	// Test atEnd() - we're at the EOF token, not past it
	if p.atEnd() {
		t.Error("Should not be at end while still on EOF token")
	}

	// Consume EOF and then we should be at end
	p.advance()
	if !p.atEnd() {
		t.Error("Should be at end after consuming EOF token")
	}
}

// TestParser_Expect tests the expect() method.
func TestParser_Expect(t *testing.T) {
	tokens := []Token{
		{Type: TokenText, Value: "text", Pos: Position{Line: 1, Column: 1}},
		{Type: TokenEOF, Value: "", Pos: Position{Line: 2, Column: 1}},
	}

	p := NewParser(tokens)

	// Test successful expect
	token, err := p.expect(TokenText)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token.Type != TokenText {
		t.Errorf("Expected TokenText, got %s", token.Type)
	}

	// Test failed expect
	p2 := NewParser(tokens)
	_, err = p2.expect(TokenHeader)
	if err == nil {
		t.Error("Expected error for mismatched token type")
	}
	if !strings.Contains(err.Error(), "expected") {
		t.Errorf("Error should mention 'expected', got: %v", err)
	}
}

// TestParse_ComplexDocument tests a realistic spec document.
func TestParse_ComplexDocument(t *testing.T) {
	input := `# Specification

## Requirements

### Requirement: User Authentication

The system SHALL authenticate users.

#### Scenario: Valid login
- **WHEN** user provides valid credentials
- **THEN** system grants access

#### Scenario: Invalid login
- **WHEN** user provides invalid credentials
- **THEN** system denies access

### Requirement: Session Management

The system SHALL manage user sessions.

## Implementation

` + "```go\nfunc Login() {}\n```\n"

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Verify we have headers
	headers := FindHeaders(doc, func(_ *Header) bool {
		return true
	})

	if len(headers) < 5 {
		t.Errorf("Expected at least 5 headers, got %d", len(headers))
	}

	// Verify we have a code block
	hasCodeBlock := false
	Walk(doc, func(n Node) bool {
		if _, ok := n.(*CodeBlock); ok {
			hasCodeBlock = true

			return false
		}

		return true
	})

	if !hasCodeBlock {
		t.Error("Expected document to contain a code block")
	}

	// Verify we have list items
	hasLists := false
	Walk(doc, func(n Node) bool {
		if _, ok := n.(*List); ok {
			hasLists = true

			return false
		}

		return true
	})

	if !hasLists {
		t.Error("Expected document to contain list items")
	}
}

// TestParse_EmptyString tests parsing an empty string.
func TestParse_EmptyString(t *testing.T) {
	doc, err := Parse("")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Empty input should produce empty document (just EOF)
	if len(doc.Children) != 0 {
		t.Errorf("Expected 0 children for empty input, got %d", len(doc.Children))
	}
}
