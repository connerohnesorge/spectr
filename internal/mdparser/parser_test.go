package mdparser

import (
	"testing"
)

func TestParser_Headers(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantHeaders int
		wantLevel   int
		wantText    string
	}{
		{
			name:        "single header",
			input:       "# Title",
			wantHeaders: 1,
			wantLevel:   1,
			wantText:    "Title",
		},
		{
			name:        "level 3 header",
			input:       "### Requirement: Login",
			wantHeaders: 1,
			wantLevel:   3,
			wantText:    "Requirement: Login",
		},
		{
			name:        "level 4 header",
			input:       "#### Scenario: Success",
			wantHeaders: 1,
			wantLevel:   4,
			wantText:    "Scenario: Success",
		},
		{
			name:        "multiple headers",
			input:       "# H1\n## H2\n### H3",
			wantHeaders: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			headerCount := 0
			for _, node := range doc.Children {
				h, ok := node.(*Header)
				if !ok {
					continue
				}
				headerCount++
				if tt.wantHeaders != 1 {
					continue
				}
				// Check specific header details
				if h.Level != tt.wantLevel {
					t.Errorf("header level = %d, want %d", h.Level, tt.wantLevel)
				}
				if h.Text != tt.wantText {
					t.Errorf("header text = %q, want %q", h.Text, tt.wantText)
				}
			}

			if headerCount != tt.wantHeaders {
				t.Errorf("got %d headers, want %d", headerCount, tt.wantHeaders)
			}
		})
	}
}

func TestParser_Paragraphs(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLines int
		wantText  string
	}{
		{
			name:      "single line paragraph",
			input:     "This is a paragraph.",
			wantLines: 1,
			wantText:  "This is a paragraph.",
		},
		{
			name:      "multi-line paragraph",
			input:     "Line 1\nLine 2\nLine 3",
			wantLines: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if len(doc.Children) == 0 {
				t.Fatal("expected at least one node")
			}

			para, ok := doc.Children[0].(*Paragraph)
			if !ok {
				t.Fatalf("expected Paragraph node, got %T", doc.Children[0])
			}

			if len(para.Lines) != tt.wantLines {
				t.Errorf(
					"paragraph has %d lines, want %d",
					len(para.Lines),
					tt.wantLines,
				)
			}

			if tt.wantLines == 1 && para.Lines[0] != tt.wantText {
				t.Errorf("paragraph text = %q, want %q", para.Lines[0], tt.wantText)
			}
		})
	}
}

func TestParser_CodeBlocks(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantLanguage string
		wantLines    int
	}{
		{
			name:         "simple code block",
			input:        "```\ncode\n```",
			wantLanguage: "",
			wantLines:    1,
		},
		{
			name:         "code block with language",
			input:        "```go\nfunc main() {}\n```",
			wantLanguage: "go",
			wantLines:    1,
		},
		{
			name:         "code block with multiple lines",
			input:        "```python\nline 1\nline 2\nline 3\n```",
			wantLanguage: "python",
			wantLines:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if len(doc.Children) == 0 {
				t.Fatal("expected at least one node")
			}

			code, ok := doc.Children[0].(*CodeBlock)
			if !ok {
				t.Fatalf("expected CodeBlock node, got %T", doc.Children[0])
			}

			if code.Language != tt.wantLanguage {
				t.Errorf("language = %q, want %q", code.Language, tt.wantLanguage)
			}

			if len(code.Lines) != tt.wantLines {
				t.Errorf(
					"code block has %d lines, want %d",
					len(code.Lines),
					tt.wantLines,
				)
			}
		})
	}
}

func TestParser_CodeBlocksWithMarkdown(t *testing.T) {
	input := `# Real Header

` + "```markdown" + `
### Requirement: Not Real
#### Scenario: Not Real
- Not a real list
` + "```" + `

## Another Real Header
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Count headers (should only be 2 real headers)
	headerCount := 0
	codeBlockCount := 0
	for _, node := range doc.Children {
		switch node.(type) {
		case *Header:
			headerCount++
		case *CodeBlock:
			codeBlockCount++
		}
	}

	if headerCount != 2 {
		t.Errorf(
			"got %d headers, want 2",
			headerCount,
		)
	}

	if codeBlockCount != 1 {
		t.Errorf("got %d code blocks, want 1", codeBlockCount)
	}

	// Verify the code block contains the markdown-like content
	for _, node := range doc.Children {
		code, ok := node.(*CodeBlock)
		if !ok {
			continue
		}
		if len(code.Lines) < 3 {
			t.Errorf(
				"code block should have at least 3 lines, got %d",
				len(code.Lines),
			)
		}
		// The content should be preserved as-is
		if code.Lines[0] != "### Requirement: Not Real" {
			t.Errorf(
				"code line 0 = %q, want '### Requirement: Not Real'",
				code.Lines[0],
			)
		}
	}
}

func TestParser_Lists(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantOrdered bool
		wantItems   int
	}{
		{
			name:        "unordered list",
			input:       "- Item 1\n- Item 2",
			wantOrdered: false,
			wantItems:   2,
		},
		{
			name:        "ordered list",
			input:       "1. First\n2. Second",
			wantOrdered: true,
			wantItems:   2,
		},
		{
			name:        "asterisk list",
			input:       "* Item A\n* Item B\n* Item C",
			wantOrdered: false,
			wantItems:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if len(doc.Children) == 0 {
				t.Fatal("expected at least one node")
			}

			list, ok := doc.Children[0].(*List)
			if !ok {
				t.Fatalf("expected List node, got %T", doc.Children[0])
			}

			if list.Ordered != tt.wantOrdered {
				t.Errorf("list.Ordered = %v, want %v", list.Ordered, tt.wantOrdered)
			}

			if len(list.Items) != tt.wantItems {
				t.Errorf("list has %d items, want %d", len(list.Items), tt.wantItems)
			}
		})
	}
}

func TestParser_BlankLines(t *testing.T) {
	input := "Text 1\n\n\nText 2"

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should have: Paragraph, BlankLine, Paragraph
	if len(doc.Children) < 3 {
		t.Fatalf("expected at least 3 nodes, got %d", len(doc.Children))
	}

	// Check for blank line node
	foundBlankLine := false
	for _, node := range doc.Children {
		bl, ok := node.(*BlankLine)
		if !ok {
			continue
		}
		foundBlankLine = true
		if bl.Count != 2 {
			t.Errorf("blank line count = %d, want 2", bl.Count)
		}
	}

	if !foundBlankLine {
		t.Error("expected to find BlankLine node")
	}
}

func TestParser_MixedContent(t *testing.T) {
	input := `# Title

Some paragraph text.

## Section

- Item 1
- Item 2

` + "```go" + `
code here
` + "```" + `

Final text.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Count node types
	counts := make(map[string]int)
	for _, node := range doc.Children {
		switch node.(type) {
		case *Header:
			counts["header"]++
		case *Paragraph:
			counts["paragraph"]++
		case *List:
			counts["list"]++
		case *CodeBlock:
			counts["code"]++
		case *BlankLine:
			counts["blank"]++
		}
	}

	if counts["header"] < 2 {
		t.Errorf("expected at least 2 headers, got %d", counts["header"])
	}
	if counts["paragraph"] < 1 {
		t.Errorf("expected at least 1 paragraph, got %d", counts["paragraph"])
	}
	if counts["list"] < 1 {
		t.Errorf("expected at least 1 list, got %d", counts["list"])
	}
	if counts["code"] < 1 {
		t.Errorf("expected at least 1 code block, got %d", counts["code"])
	}
}

func TestParser_ErrorRecovery(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unclosed code block",
			input: "```\ncode without closing",
		},
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "only whitespace",
			input: "   \n\n  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf(
					"Parse() should not error on malformed input, got: %v",
					err,
				)
			}
			if doc == nil {
				t.Fatal("Parse() returned nil document")
			}
		})
	}
}

func TestParser_NestedStructures(t *testing.T) {
	input := `# Main

## Subsection

### Requirement: Feature

The system SHALL do something.

#### Scenario: Success case

- **WHEN** user acts
- **THEN** result happens

#### Scenario: Failure case

- **WHEN** error occurs
- **THEN** error handled
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Count headers by level
	levelCounts := make(map[int]int)
	for _, node := range doc.Children {
		h, ok := node.(*Header)
		if !ok {
			continue
		}
		levelCounts[h.Level]++
	}

	if levelCounts[1] != 1 {
		t.Errorf("expected 1 level-1 header, got %d", levelCounts[1])
	}
	if levelCounts[2] != 1 {
		t.Errorf("expected 1 level-2 header, got %d", levelCounts[2])
	}
	if levelCounts[3] != 1 {
		t.Errorf("expected 1 level-3 header, got %d", levelCounts[3])
	}
	if levelCounts[4] != 2 {
		t.Errorf("expected 2 level-4 headers, got %d", levelCounts[4])
	}
}

func TestParser_SpectrFormat(t *testing.T) {
	// Test parsing a realistic Spectr spec file
	input := `## ADDED Requirements

### Requirement: User Authentication

The system SHALL authenticate users via OAuth2.

#### Scenario: Successful login

- **WHEN** valid credentials are provided
- **THEN** JWT token is returned
- **AND** session is created

#### Scenario: Invalid credentials

- **WHEN** invalid credentials are provided
- **THEN** authentication fails
- **AND** error message is returned

## MODIFIED Requirements

### Requirement: Data Validation

The system SHALL validate all input data.

#### Scenario: Valid input

- **WHEN** valid data is submitted
- **THEN** data is processed
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should have headers at different levels
	foundLevel2 := false
	foundLevel3 := false
	foundLevel4 := false

	for _, node := range doc.Children {
		h, ok := node.(*Header)
		if !ok {
			continue
		}
		switch h.Level {
		case 2:
			foundLevel2 = true
		case 3:
			foundLevel3 = true
		case 4:
			foundLevel4 = true
		}
	}

	if !foundLevel2 {
		t.Error("expected to find level-2 headers (delta sections)")
	}
	if !foundLevel3 {
		t.Error("expected to find level-3 headers (requirements)")
	}
	if !foundLevel4 {
		t.Error("expected to find level-4 headers (scenarios)")
	}
}

func TestParser_EmptyDocument(t *testing.T) {
	doc, err := Parse("")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if len(doc.Children) != 0 {
		t.Errorf("empty input should have 0 children, got %d", len(doc.Children))
	}
}

func TestParser_ListItemText(t *testing.T) {
	input := "- **WHEN** user logs in\n- **THEN** session created"

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	list, ok := doc.Children[0].(*List)
	if !ok {
		t.Fatalf("expected List node, got %T", doc.Children[0])
	}

	if len(list.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list.Items))
	}

	// Check that item text is properly extracted
	if list.Items[0].Text != "**WHEN** user logs in" {
		t.Errorf(
			"item 0 text = %q, want '**WHEN** user logs in'",
			list.Items[0].Text,
		)
	}
	if list.Items[1].Text != "**THEN** session created" {
		t.Errorf(
			"item 1 text = %q, want '**THEN** session created'",
			list.Items[1].Text,
		)
	}
}
