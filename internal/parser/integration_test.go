package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseRealSpecFiles tests parsing actual spec files from the project.
// This ensures the parser works with real-world Spectr specs.
func TestParseRealSpecFiles(t *testing.T) {
	// Find all spec.md files in the project
	specDir := filepath.Join("..", "..", "spectr", "specs")

	// Check if the directory exists (it may not in test environments)
	if _, err := os.Stat(specDir); os.IsNotExist(err) {
		t.Skip("Spec directory not found, skipping integration test")

		return
	}

	var specFiles []string
	err := filepath.Walk(specDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "spec.md" {
			specFiles = append(specFiles, path)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk spec directory: %v", err)
	}

	if len(specFiles) == 0 {
		t.Skip("No spec files found, skipping integration test")

		return
	}

	// Parse each spec file
	for _, specFile := range specFiles {
		t.Run(filepath.Base(filepath.Dir(specFile)), func(t *testing.T) {
			content, err := os.ReadFile(specFile)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", specFile, err)
			}

			doc, err := Parse(string(content))
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", specFile, err)
			}

			if doc == nil {
				t.Fatal("Expected non-nil document")
			}

			// Verify we can extract requirements
			reqs, err := ExtractRequirements(doc)
			if err != nil {
				t.Fatalf("Failed to extract requirements from %s: %v", specFile, err)
			}

			// Most specs should have at least one requirement
			t.Logf("Parsed %s: %d children, %d requirements",
				specFile, len(doc.Children), len(reqs))

			// Verify sections can be extracted
			sections, err := ExtractSections(doc)
			if err != nil {
				t.Fatalf("Failed to extract sections from %s: %v", specFile, err)
			}

			t.Logf("  - %d sections extracted", len(sections))
		})
	}
}

// TestLexer_LargeDocument tests lexing a very large document.
func TestLexer_LargeDocument(t *testing.T) {
	// Generate a large document with many headers, code blocks, and lists
	var sb strings.Builder

	sb.WriteString("# Main Title\n\n")
	sb.WriteString("This is an introduction paragraph.\n\n")

	for range 100 {
		sb.WriteString("## Section ")
		sb.WriteString(strings.Repeat("X", 1000)) // Long header
		sb.WriteString("\n\n")

		sb.WriteString("Some content paragraph with text.\n\n")

		sb.WriteString("```go\n")
		sb.WriteString("func example")
		sb.WriteString(strings.Repeat("X", 500)) // Long code
		sb.WriteString("() {}\n```\n\n")

		sb.WriteString("- List item one\n")
		sb.WriteString("- List item two\n\n")
	}

	input := sb.String()

	// Lex the large document
	l := NewLexer(input)
	tokens := l.Lex()

	// Should have many tokens
	if len(tokens) < 500 {
		t.Errorf("Expected at least 500 tokens, got %d", len(tokens))
	}

	// Should end with EOF
	if tokens[len(tokens)-1].Type != TokenEOF {
		t.Error("Expected last token to be EOF")
	}

	// Should have no errors
	for _, tok := range tokens {
		if tok.Type == TokenError {
			t.Errorf("Unexpected error token: %s", tok.Value)
		}
	}
}

// TestParser_LargeDocument tests parsing a very large document.
func TestParser_LargeDocument(t *testing.T) {
	var sb strings.Builder

	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 1000; i++ {
		sb.WriteString("### Requirement: Feature ")
		sb.WriteString(strings.Repeat("A", i%10+1))
		sb.WriteString("\n\nThe system SHALL provide feature.\n\n")

		sb.WriteString("#### Scenario: Test case\n")
		sb.WriteString("- **WHEN** action\n")
		sb.WriteString("- **THEN** result\n\n")
	}

	input := sb.String()

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse large document: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Should have many children
	if len(doc.Children) < 1000 {
		t.Errorf("Expected at least 1000 children, got %d", len(doc.Children))
	}
}

// TestExtractor_LargeDocument tests extracting from a very large document.
func TestExtractor_LargeDocument(t *testing.T) {
	var sb strings.Builder

	sb.WriteString("## Requirements\n\n")

	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 500; i++ {
		sb.WriteString("### Requirement: Feature ")
		sb.WriteString(strings.Repeat("B", i%5+1))
		sb.WriteString("\n\nThe system SHALL provide feature.\n\n")

		for j := 0; j < 3; j++ {
			sb.WriteString("#### Scenario: Test case ")
			sb.WriteString(strings.Repeat("C", j+1))
			sb.WriteString("\n- **WHEN** action\n")
			sb.WriteString("- **THEN** result\n\n")
		}
	}

	input := sb.String()

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("Failed to extract requirements: %v", err)
	}

	// Should extract all requirements
	if len(reqs) != 500 {
		t.Errorf("Expected 500 requirements, got %d", len(reqs))
	}

	// Each should have 3 scenarios
	for i, req := range reqs {
		if len(req.Scenarios) != 3 {
			t.Errorf("Requirement %d: expected 3 scenarios, got %d", i, len(req.Scenarios))

			break
		}
	}
}

// TestLexer_SpecialCharacters tests lexing with special characters.
func TestLexer_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unicode characters",
			input: "# 你好世界\n\nThis has unicode: 🚀 ✨ 🎉\n",
		},
		{
			name:  "special markdown chars",
			input: "# Header with **bold** and *italic*\n\n```\nCode with `backticks` inside\n```\n",
		},
		{
			name:  "URLs and paths",
			input: "# Links\n\nSee https://example.com and /path/to/file.txt\n",
		},
		{
			name:  "HTML entities",
			input: "# Title\n\n&lt;html&gt; &amp; other entities\n",
		},
		{
			name:  "tabs and special whitespace",
			input: "# Title\n\n\tTabbed content\n\t- List with tabs\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := l.Lex()

			// Should complete without errors
			hasError := false
			for _, tok := range tokens {
				if tok.Type == TokenError {
					hasError = true
					t.Errorf("Unexpected error: %s", tok.Value)
				}
			}

			if hasError {
				t.Fatal("Lexing failed for input with special characters")
			}

			// Should end with EOF
			if len(tokens) > 0 && tokens[len(tokens)-1].Type != TokenEOF {
				t.Error("Expected last token to be EOF")
			}
		})
	}
}

// TestParser_EdgeCases tests parser edge cases.
func TestParser_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "only blank lines",
			input:     "\n\n\n\n\n",
			expectErr: false,
		},
		{
			name:      "only code block",
			input:     "```\ncode\n```",
			expectErr: false,
		},
		{
			name:      "only lists",
			input:     "- Item 1\n- Item 2\n- Item 3",
			expectErr: false,
		},
		{
			name:      "mixed newlines",
			input:     "text\r\n# header\r\nmore text",
			expectErr: false,
		},
		{
			name:      "very long line",
			input:     "# " + strings.Repeat("X", 10000) + "\n",
			expectErr: false,
		},
		{
			name:      "nested code blocks (invalid)",
			input:     "```\n```\ncode\n```\n```",
			expectErr: true, // Lexer will report unclosed code block error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectErr && doc == nil {
				t.Error("Expected non-nil document")
			}
		})
	}
}

// TestExtractor_EdgeCases tests extractor edge cases.
func TestExtractor_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantReqs  int
		expectErr bool
	}{
		{
			name: "requirement with empty content",
			input: `## Requirements

### Requirement: Empty

### Requirement: Also Empty
`,
			wantReqs:  2,
			expectErr: false,
		},
		{
			name: "scenario with no content",
			input: `## Requirements

### Requirement: Has Empty Scenario

#### Scenario: Empty
`,
			wantReqs:  1,
			expectErr: false,
		},
		{
			name: "requirement with only whitespace name",
			input: `## Requirements

### Requirement:
Content here.
`,
			wantReqs:  1,
			expectErr: false,
		},
		{
			name: "many nested headers",
			input: `## Requirements

### Requirement: Level 3
#### Scenario: Level 4
##### Not a scenario (level 5)
###### Also not a scenario (level 6)
`,
			wantReqs:  1,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			reqs, err := ExtractRequirements(doc)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectErr && len(reqs) != tt.wantReqs {
				t.Errorf("Expected %d requirements, got %d", tt.wantReqs, len(reqs))
			}
		})
	}
}

// TestLexer_BackupEdgeCases tests backup over complex boundaries.
func TestLexer_BackupEdgeCases(t *testing.T) {
	l := NewLexer("abc\n\ndef")

	// Consume multiple characters including newlines
	l.next() // 'a' at 1:2
	l.next() // 'b' at 1:3
	l.next() // 'c' at 1:4
	l.next() // '\n' at 2:1
	l.next() // '\n' at 3:1
	l.next() // 'd' at 3:2

	// Backup should work (backs up one character)
	l.backup() // back from 'd' at 3:2 to '\n' at 3:1

	// Verify position - should be back to the second newline
	if l.line != 3 || l.column != 1 {
		t.Errorf("After backup from 'd', expected 3:1, got %d:%d", l.line, l.column)
	}

	// Read again to verify it's the second newline
	r := l.next()
	if r != 'd' {
		t.Errorf("After backup and next, expected 'd', got %c", r)
	}
}

// TestParser_UnexpectedTokenSequences tests handling of unexpected token sequences.
func TestParser_UnexpectedTokenSequences(t *testing.T) {
	// Create token sequences manually that might not occur naturally
	tests := []struct {
		name   string
		tokens []Token
	}{
		{
			name: "multiple EOF tokens",
			tokens: []Token{
				{Type: TokenText, Value: "text", Pos: Position{Line: 1, Column: 1}},
				{Type: TokenEOF, Value: "", Pos: Position{Line: 1, Column: 5}},
				{Type: TokenEOF, Value: "", Pos: Position{Line: 1, Column: 5}},
			},
		},
		{
			name: "blank lines only",
			tokens: []Token{
				{Type: TokenBlankLine, Value: "\n", Pos: Position{Line: 1, Column: 1}},
				{Type: TokenBlankLine, Value: "\n", Pos: Position{Line: 2, Column: 1}},
				{Type: TokenBlankLine, Value: "\n", Pos: Position{Line: 3, Column: 1}},
				{Type: TokenEOF, Value: "", Pos: Position{Line: 4, Column: 1}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := ParseTokens(tt.tokens)
			if err != nil {
				t.Errorf("ParseTokens failed: %v", err)
			}
			if doc == nil {
				t.Error("Expected non-nil document")
			}
		})
	}
}

// TestExtractDeltas_MalformedRENAMED tests handling of malformed RENAMED entries.
func TestExtractDeltas_MalformedRENAMED(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name: "missing TO",
			input: `## RENAMED Requirements

- FROM: ` + "`### Requirement: Old Name`" + `
`,
			expectErr: false, // Should handle gracefully, just won't find the rename
		},
		{
			name: "missing FROM",
			input: `## RENAMED Requirements

- TO: ` + "`### Requirement: New Name`" + `
`,
			expectErr: false, // Should handle gracefully
		},
		{
			name: "invalid format",
			input: `## RENAMED Requirements

Just some text, not a rename entry.
`,
			expectErr: false,
		},
		{
			name: "empty requirement name",
			input: `## RENAMED Requirements

- FROM: ` + "`### Requirement: `" + `
- TO: ` + "`### Requirement: `" + `
`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			delta, err := ExtractDeltas(doc)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// In most cases, malformed entries should just be ignored
			if !tt.expectErr && delta == nil {
				t.Error("Expected non-nil delta")
			}
		})
	}
}

// TestNodeType_String tests NodeType.String() for all types including unknown.
func TestNodeType_String(t *testing.T) {
	tests := []struct {
		nodeType NodeType
		want     string
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
		if got != tt.want {
			t.Errorf("NodeType(%d).String() = %q, want %q", tt.nodeType, got, tt.want)
		}
	}
}

// TestAST_StringMethods tests all AST node String() methods.
func TestAST_StringMethods(t *testing.T) {
	// Test Document.String()
	doc := NewDocument()
	doc.Children = append(doc.Children, NewHeader(1, "Test", Position{Line: 1, Column: 1}))
	if !strings.Contains(doc.String(), "1 nodes") {
		t.Errorf("Document.String() should contain node count, got: %s", doc.String())
	}

	// Test Header.String() with long text
	longHeader := NewHeader(2, strings.Repeat("X", 100), Position{Line: 5, Column: 1})
	if !strings.Contains(longHeader.String(), "...") {
		t.Errorf("Long Header.String() should truncate, got: %s", longHeader.String())
	}

	// Test Paragraph.String() with long text
	longPara := NewParagraph(strings.Repeat("Y", 100), Position{Line: 10, Column: 1})
	if !strings.Contains(longPara.String(), "...") {
		t.Errorf("Long Paragraph.String() should truncate, got: %s", longPara.String())
	}

	// Test CodeBlock.String() with language
	codeWithLang := NewCodeBlock("go", "func main() {}", Position{Line: 15, Column: 1})
	if !strings.Contains(codeWithLang.String(), "Lang: go") {
		t.Errorf("CodeBlock.String() should include language, got: %s", codeWithLang.String())
	}

	// Test CodeBlock.String() without language
	codeNoLang := NewCodeBlock("", "code", Position{Line: 20, Column: 1})
	if strings.Contains(codeNoLang.String(), "Lang:") {
		t.Errorf(
			"CodeBlock.String() should not include Lang when empty, got: %s",
			codeNoLang.String(),
		)
	}

	// Test List.String() with empty items
	emptyList := NewList("", Position{Line: 25, Column: 1})
	emptyList.Items = make([]string, 0)
	if !strings.Contains(emptyList.String(), "Empty") {
		t.Errorf("Empty List.String() should say Empty, got: %s", emptyList.String())
	}

	// Test List.String() with multiple items
	multiList := &List{
		Items:    []string{"item1", "item2", "item3"},
		position: Position{Line: 30, Column: 1},
	}
	if !strings.Contains(multiList.String(), "3") {
		t.Errorf("Multi-item List.String() should show count, got: %s", multiList.String())
	}

	// Test BlankLine.String() with single line
	singleBlank := NewBlankLine(1, Position{Line: 35, Column: 1})
	if strings.Contains(singleBlank.String(), "Count:") {
		t.Errorf("Single BlankLine.String() should not show count, got: %s", singleBlank.String())
	}

	// Test BlankLine.String() with multiple lines
	multiBlank := NewBlankLine(5, Position{Line: 40, Column: 1})
	if !strings.Contains(multiBlank.String(), "Count: 5") {
		t.Errorf("Multi-line BlankLine.String() should show count, got: %s", multiBlank.String())
	}
}

// TestWalk_EarlyTermination tests that Walk stops when visitor returns false.
func TestWalk_EarlyTermination(t *testing.T) {
	doc := NewDocument()
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 10; i++ {
		doc.Children = append(doc.Children, NewParagraph("text", Position{Line: i + 1, Column: 1}))
	}

	count := 0
	Walk(doc, func(_ Node) bool {
		count++
		// Stop after visiting 3 nodes (document + 2 children)
		return count < 3
	})

	if count != 3 {
		t.Errorf("Expected Walk to stop after 3 nodes, visited %d", count)
	}
}

// TestNodesBetween_EdgeCases tests NodesBetween with edge cases.
func TestNodesBetween_EdgeCases(t *testing.T) {
	doc := NewDocument()
	doc.Children = append(doc.Children,
		NewHeader(1, "First", Position{Line: 1, Column: 1, Offset: 0}),
		NewParagraph("Content", Position{Line: 2, Column: 1, Offset: 10}),
		NewHeader(2, "Second", Position{Line: 3, Column: 1, Offset: 20}),
		NewParagraph("More", Position{Line: 4, Column: 1, Offset: 30}),
	)

	// Test with same start and end (should return nothing)
	nodes := NodesBetween(doc, Position{Offset: 10}, Position{Offset: 10})
	if len(nodes) != 0 {
		t.Errorf("NodesBetween with same start/end should return 0 nodes, got %d", len(nodes))
	}

	// Test with end before start (should return nothing)
	nodes = NodesBetween(doc, Position{Offset: 30}, Position{Offset: 10})
	if len(nodes) != 0 {
		t.Errorf("NodesBetween with end before start should return 0 nodes, got %d", len(nodes))
	}

	// Test spanning entire document
	nodes = NodesBetween(doc, Position{Offset: 0}, Position{Offset: 1000})
	if len(nodes) != 3 { // All children except the start node
		t.Errorf("NodesBetween spanning document should return 3 nodes, got %d", len(nodes))
	}
}
