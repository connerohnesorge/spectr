//nolint:revive // empty-block: benchmark loops are intentionally empty
package markdown

import (
	"bytes"
	"testing"
)

// Test document sizes for benchmarks

// smallDoc is a minimal Spectr-style document (~100 bytes)
var smallDoc = []byte(`# Spec: Test

## Requirements

### Requirement: Feature

- **WHEN** condition
- **THEN** result
`)

// mediumDoc is a typical Spectr document (~1KB)
var mediumDoc = []byte(
	`# Spec: User Authentication

## Overview

This specification defines the user authentication requirements.

## Requirements

### Requirement: Login

Users must be able to log into the system using credentials.

#### Scenario: Successful Login

- **WHEN** user enters valid credentials
- **AND** user clicks login button
- **THEN** system authenticates user
- **AND** redirects to dashboard

#### Scenario: Failed Login

- **WHEN** user enters invalid credentials
- **AND** user clicks login button
- **THEN** system shows error message
- **AND** keeps user on login page

### Requirement: Password Reset

Users must be able to reset forgotten passwords.

#### Scenario: Request Reset

- **WHEN** user clicks forgot password
- **THEN** system prompts for email
- **AND** sends reset link

## Technical Details

The authentication system uses:

` + "```go" + `
type AuthService struct {
    db Database
    cache Cache
}

func (a *AuthService) Login(username, password string) error {
    // implementation
    return nil
}
` + "```" + `

> Note: All passwords must be hashed using bcrypt.

See [[Authentication API]] for more details.
`,
)

// largeDoc is a comprehensive Spectr document (~10KB)
var largeDoc = generateLargeDoc()

func generateLargeDoc() []byte {
	var buf bytes.Buffer

	buf.WriteString(
		"# Spec: E-Commerce Platform\n\n",
	)
	buf.WriteString("## Overview\n\n")
	buf.WriteString(
		"This specification defines the complete e-commerce platform requirements.\n\n",
	)

	// Generate multiple requirement sections
	requirements := []string{
		"Product Catalog",
		"Shopping Cart",
		"Checkout Process",
		"Payment Processing",
		"Order Management",
		"User Accounts",
		"Search Functionality",
		"Reviews and Ratings",
		"Inventory Management",
		"Shipping Integration",
	}

	for _, req := range requirements {
		buf.WriteString("## " + req + "\n\n")
		buf.WriteString(
			"### Requirement: " + req + "\n\n",
		)
		buf.WriteString(
			"This requirement covers all aspects of " + req + ".\n\n",
		)

		// Multiple scenarios per requirement
		scenarios := []string{
			"Basic Flow",
			"Edge Cases",
			"Error Handling",
			"Performance",
		}
		for _, scenario := range scenarios {
			buf.WriteString(
				"#### Scenario: " + scenario + "\n\n",
			)
			buf.WriteString(
				"- **WHEN** user performs initial action\n",
			)
			buf.WriteString(
				"- **AND** system is in ready state\n",
			)
			buf.WriteString(
				"- **AND** all preconditions are met\n",
			)
			buf.WriteString(
				"- **THEN** system responds appropriately\n",
			)
			buf.WriteString(
				"- **AND** state is updated correctly\n",
			)
			buf.WriteString(
				"- **AND** user receives feedback\n\n",
			)
		}

		// Add code block
		buf.WriteString("```go\n")
		buf.WriteString(
			"type " + req + "Service struct {\n",
		)
		buf.WriteString("    repo Repository\n")
		buf.WriteString("    cache Cache\n")
		buf.WriteString("}\n")
		buf.WriteString("```\n\n")

		// Add blockquote
		buf.WriteString(
			"> Important: " + req + " must be implemented with care.\n\n",
		)

		// Add wikilink
		buf.WriteString(
			"See [[" + req + " API]] for implementation details.\n\n",
		)
	}

	return buf.Bytes()
}

// =============================================================================
// Lexer Benchmarks (Task 21.2)
// =============================================================================

func BenchmarkLexerSmall(b *testing.B) {
	b.SetBytes(int64(len(smallDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		lex := newLexer(smallDoc)
		for {
			tok := lex.Next()
			if tok.Type == TokenEOF {
				break
			}
		}
	}
}

func BenchmarkLexerMedium(b *testing.B) {
	b.SetBytes(int64(len(mediumDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		lex := newLexer(mediumDoc)
		for {
			tok := lex.Next()
			if tok.Type == TokenEOF {
				break
			}
		}
	}
}

func BenchmarkLexerLarge(b *testing.B) {
	b.SetBytes(int64(len(largeDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		lex := newLexer(largeDoc)
		for {
			tok := lex.Next()
			if tok.Type == TokenEOF {
				break
			}
		}
	}
}

func BenchmarkLexerAll(b *testing.B) {
	b.Run("Small", BenchmarkLexerSmall)
	b.Run("Medium", BenchmarkLexerMedium)
	b.Run("Large", BenchmarkLexerLarge)
}

// BenchmarkLexerTokenTypes benchmarks lexing documents with specific token patterns
func BenchmarkLexerTokenTypes(b *testing.B) {
	// Document heavy on headers
	headerDoc := []byte(
		"# H1\n## H2\n### H3\n#### H4\n##### H5\n###### H6\n",
	)
	for range 10 {
		headerDoc = append(
			headerDoc,
			headerDoc...)
	}

	// Document heavy on lists
	listDoc := []byte(
		"- item 1\n- item 2\n- item 3\n* star item\n+ plus item\n",
	)
	for range 10 {
		listDoc = append(listDoc, listDoc...)
	}

	// Document heavy on inline formatting
	inlineDoc := []byte(
		"**bold** *italic* ~~strike~~ `code` [link](url) [[wikilink]]\n",
	)
	for range 10 {
		inlineDoc = append(
			inlineDoc,
			inlineDoc...)
	}

	b.Run("Headers", func(b *testing.B) {
		b.SetBytes(int64(len(headerDoc)))
		b.ReportAllocs()
		for range b.N {
			lex := newLexer(headerDoc)
			for lex.Next().Type != TokenEOF {
				// intentionally empty for benchmarking
			}
		}
	})

	b.Run("Lists", func(b *testing.B) {
		b.SetBytes(int64(len(listDoc)))
		b.ReportAllocs()
		for range b.N {
			lex := newLexer(listDoc)
			for lex.Next().Type != TokenEOF {
				// intentionally empty for benchmarking
			}
		}
	})

	b.Run("Inline", func(b *testing.B) {
		b.SetBytes(int64(len(inlineDoc)))
		b.ReportAllocs()
		for range b.N {
			lex := newLexer(inlineDoc)
			for lex.Next().Type != TokenEOF {
				// intentionally empty for benchmarking
			}
		}
	})
}

// =============================================================================
// Parser Benchmarks (Task 21.3)
// =============================================================================

func BenchmarkParseSmall(b *testing.B) {
	b.SetBytes(int64(len(smallDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = Parse(smallDoc)
	}
}

func BenchmarkParseMedium(b *testing.B) {
	b.SetBytes(int64(len(mediumDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = Parse(mediumDoc)
	}
}

func BenchmarkParseLarge(b *testing.B) {
	b.SetBytes(int64(len(largeDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = Parse(largeDoc)
	}
}

func BenchmarkParseAll(b *testing.B) {
	b.Run("Small", BenchmarkParseSmall)
	b.Run("Medium", BenchmarkParseMedium)
	b.Run("Large", BenchmarkParseLarge)
}

// BenchmarkParseDocumentTypes benchmarks parsing specific document structures
func BenchmarkParseDocumentTypes(b *testing.B) {
	// Deep nested list document
	nestedListDoc := []byte(`
- Level 1
  - Level 2
    - Level 3
      - Level 4
- Another Level 1
  - Another Level 2
`)

	// Code-heavy document
	codeDoc := []byte(
		"```go\npackage main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```\n\n",
	)
	for range 5 {
		codeDoc = append(codeDoc, codeDoc...)
	}

	// Wikilink-heavy document
	wikilinkDoc := []byte(
		"See [[Page1]], [[Page2|Display]], [[Page3#anchor]], and [[Page4|Text#section]].\n",
	)
	for range 10 {
		wikilinkDoc = append(
			wikilinkDoc,
			wikilinkDoc...)
	}

	b.Run("NestedLists", func(b *testing.B) {
		b.SetBytes(int64(len(nestedListDoc)))
		b.ReportAllocs()
		for range b.N {
			_, _ = Parse(nestedListDoc)
		}
	})

	b.Run("CodeBlocks", func(b *testing.B) {
		b.SetBytes(int64(len(codeDoc)))
		b.ReportAllocs()
		for range b.N {
			_, _ = Parse(codeDoc)
		}
	})

	b.Run("Wikilinks", func(b *testing.B) {
		b.SetBytes(int64(len(wikilinkDoc)))
		b.ReportAllocs()
		for range b.N {
			_, _ = Parse(wikilinkDoc)
		}
	})
}

// =============================================================================
// Incremental Parsing Benchmarks (Task 21.4)
// =============================================================================

func BenchmarkParseIncrementalSmallEdit(
	b *testing.B,
) {
	// Parse the original document
	oldTree, _ := Parse(mediumDoc)

	// Create a small edit (insert a word)
	editPos := 50
	newDoc := make([]byte, len(mediumDoc)+5)
	copy(newDoc[:editPos], mediumDoc[:editPos])
	copy(
		newDoc[editPos:editPos+5],
		[]byte("test "),
	)
	copy(newDoc[editPos+5:], mediumDoc[editPos:])

	b.SetBytes(int64(len(newDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = ParseIncremental(
			oldTree,
			mediumDoc,
			newDoc,
		)
	}
}

func BenchmarkParseIncrementalLargeEdit(
	b *testing.B,
) {
	// Parse the original document
	oldTree, _ := Parse(mediumDoc)

	// Create a larger edit (insert a new section)
	newSection := []byte(
		"\n### Requirement: New Feature\n\n- **WHEN** something happens\n- **THEN** something else\n",
	)
	editPos := len(mediumDoc) / 2
	newDoc := make(
		[]byte,
		len(mediumDoc)+len(newSection),
	)
	copy(newDoc[:editPos], mediumDoc[:editPos])
	copy(
		newDoc[editPos:editPos+len(newSection)],
		newSection,
	)
	copy(
		newDoc[editPos+len(newSection):],
		mediumDoc[editPos:],
	)

	b.SetBytes(int64(len(newDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = ParseIncremental(
			oldTree,
			mediumDoc,
			newDoc,
		)
	}
}

func BenchmarkParseIncrementalWithState(
	b *testing.B,
) {
	// Parse the original document and create state
	oldTree, _ := Parse(mediumDoc)
	state := NewIncrementalParseState(
		oldTree,
		mediumDoc,
	)

	// Create a small edit
	editPos := 50
	newDoc := make([]byte, len(mediumDoc)+5)
	copy(newDoc[:editPos], mediumDoc[:editPos])
	copy(
		newDoc[editPos:editPos+5],
		[]byte("test "),
	)
	copy(newDoc[editPos+5:], mediumDoc[editPos:])

	b.SetBytes(int64(len(newDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _, _ = ParseIncrementalWithState(
			state,
			oldTree,
			mediumDoc,
			newDoc,
		)
	}
}

func BenchmarkParseFullVsIncremental(
	b *testing.B,
) {
	oldTree, _ := Parse(mediumDoc)

	// Small edit
	editPos := 50
	newDoc := make([]byte, len(mediumDoc)+5)
	copy(newDoc[:editPos], mediumDoc[:editPos])
	copy(
		newDoc[editPos:editPos+5],
		[]byte("test "),
	)
	copy(newDoc[editPos+5:], mediumDoc[editPos:])

	b.Run("FullParse", func(b *testing.B) {
		b.SetBytes(int64(len(newDoc)))
		b.ReportAllocs()
		for range b.N {
			_, _ = Parse(newDoc)
		}
	})

	b.Run("IncrementalParse", func(b *testing.B) {
		b.SetBytes(int64(len(newDoc)))
		b.ReportAllocs()
		for range b.N {
			_, _ = ParseIncremental(
				oldTree,
				mediumDoc,
				newDoc,
			)
		}
	})
}

// =============================================================================
// Position Index Benchmarks (Task 21.5)
// =============================================================================

func BenchmarkPositionIndexBuild(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		idx := NewPositionIndex(tree, largeDoc)
		// Force build by calling a query method
		_ = idx.NodeAt(0)
	}
}

func BenchmarkPositionIndexNodeAt(b *testing.B) {
	tree, _ := Parse(largeDoc)
	idx := NewPositionIndex(tree, largeDoc)
	// Pre-build the index
	_ = idx.NodeAt(0)

	// Test positions at beginning, middle, and end
	positions := []int{
		0,
		len(largeDoc) / 4,
		len(largeDoc) / 2,
		3 * len(largeDoc) / 4,
		len(largeDoc) - 1,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		for _, pos := range positions {
			_ = idx.NodeAt(pos)
		}
	}
}

func BenchmarkPositionIndexNodesAt(b *testing.B) {
	tree, _ := Parse(largeDoc)
	idx := NewPositionIndex(tree, largeDoc)
	// Pre-build the index
	_ = idx.NodeAt(0)

	midPos := len(largeDoc) / 2

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = idx.NodesAt(midPos)
	}
}

func BenchmarkPositionIndexNodesInRange(
	b *testing.B,
) {
	tree, _ := Parse(largeDoc)
	idx := NewPositionIndex(tree, largeDoc)
	// Pre-build the index
	_ = idx.NodeAt(0)

	// Test various range sizes
	b.Run("SmallRange", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = idx.NodesInRange(0, 100)
		}
	})

	b.Run("MediumRange", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = idx.NodesInRange(
				0,
				len(largeDoc)/2,
			)
		}
	})

	b.Run("LargeRange", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = idx.NodesInRange(0, len(largeDoc))
		}
	})
}

func BenchmarkPositionIndexEnclosingSection(
	b *testing.B,
) {
	tree, _ := Parse(largeDoc)
	idx := NewPositionIndex(tree, largeDoc)
	// Pre-build
	_ = idx.NodeAt(0)

	midPos := len(largeDoc) / 2

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = idx.EnclosingSection(midPos)
	}
}

// =============================================================================
// Line Index Benchmarks (Task 21.5 continued)
// =============================================================================

func BenchmarkLineIndexBuild(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		idx := NewLineIndex(largeDoc)
		// Force build by calling LineCol
		_, _ = idx.LineCol(0)
	}
}

func BenchmarkLineIndexLineCol(b *testing.B) {
	idx := NewLineIndex(largeDoc)
	// Pre-build
	_, _ = idx.LineCol(0)

	// Test positions at various points
	positions := []int{
		0,
		len(largeDoc) / 4,
		len(largeDoc) / 2,
		3 * len(largeDoc) / 4,
		len(largeDoc) - 1,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		for _, pos := range positions {
			_, _ = idx.LineCol(pos)
		}
	}
}

func BenchmarkLineIndexOffsetAt(b *testing.B) {
	idx := NewLineIndex(largeDoc)
	// Pre-build and get line count
	lineCount := idx.LineCount()

	// Test various line numbers
	lines := []int{
		1,
		lineCount / 4,
		lineCount / 2,
		3 * lineCount / 4,
		lineCount,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		for _, line := range lines {
			_ = idx.OffsetAt(line, 0)
		}
	}
}

// =============================================================================
// Pool Benchmarks (Task 21.6)
// =============================================================================

func BenchmarkPoolTokenGetPut(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		tok := GetToken()
		tok.Type = TokenText
		tok.Start = 0
		tok.End = 10
		PutToken(tok)
	}
}

func BenchmarkPoolNodeGetPut(b *testing.B) {
	b.Run("Document", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			doc := GetDocument()
			PutNode(doc)
		}
	})

	b.Run("Section", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			sec := GetSection()
			PutNode(sec)
		}
	})

	b.Run("Paragraph", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			p := GetParagraph()
			PutNode(p)
		}
	})

	b.Run("Text", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			t := GetText()
			PutNode(t)
		}
	})
}

func BenchmarkPoolChildrenGetPut(b *testing.B) {
	b.Run("Small", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			c := GetChildren(4)
			PutChildren(c)
		}
	})

	b.Run("Medium", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			c := GetChildren(16)
			PutChildren(c)
		}
	})

	b.Run("Large", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			c := GetChildren(64)
			PutChildren(c)
		}
	})
}

// BenchmarkParseWithPoolStats benchmarks parsing with pool statistics enabled
func BenchmarkParseWithPoolStats(b *testing.B) {
	b.Run("StatsDisabled", func(b *testing.B) {
		DisablePoolStats()
		b.SetBytes(int64(len(mediumDoc)))
		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			_, _ = Parse(mediumDoc)
		}
	})

	b.Run("StatsEnabled", func(b *testing.B) {
		EnablePoolStats()
		defer DisablePoolStats()

		b.SetBytes(int64(len(mediumDoc)))
		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			_, _ = Parse(mediumDoc)
		}
	})
}

// =============================================================================
// Printer Benchmarks
// =============================================================================

func BenchmarkPrintSmall(b *testing.B) {
	tree, _ := Parse(smallDoc)

	b.SetBytes(int64(len(smallDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = Print(tree)
	}
}

func BenchmarkPrintMedium(b *testing.B) {
	tree, _ := Parse(mediumDoc)

	b.SetBytes(int64(len(mediumDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = Print(tree)
	}
}

func BenchmarkPrintLarge(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.SetBytes(int64(len(largeDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = Print(tree)
	}
}

func BenchmarkPrintToBuffer(b *testing.B) {
	tree, _ := Parse(largeDoc)
	var buf bytes.Buffer
	buf.Grow(
		len(largeDoc) * 2,
	) // Pre-allocate buffer

	b.SetBytes(int64(len(largeDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		buf.Reset()
		_ = PrintTo(&buf, tree)
	}
}

// =============================================================================
// Query Benchmarks
// =============================================================================

func BenchmarkFind(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.Run("AllSections", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Find(tree, IsType[*NodeSection]())
		}
	})

	b.Run("AllRequirements", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Find(
				tree,
				IsType[*NodeRequirement](),
			)
		}
	})

	b.Run("AllListItems", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Find(
				tree,
				IsType[*NodeListItem](),
			)
		}
	})

	b.Run("AllTextNodes", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Find(tree, IsType[*NodeText]())
		}
	})
}

func BenchmarkFindFirst(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.Run("FirstSection", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = FindFirst(
				tree,
				IsType[*NodeSection](),
			)
		}
	})

	b.Run("FirstRequirement", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = FindFirst(
				tree,
				IsType[*NodeRequirement](),
			)
		}
	})

	b.Run("FirstCodeBlock", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = FindFirst(
				tree,
				IsType[*NodeCodeBlock](),
			)
		}
	})
}

func BenchmarkFindByType(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.Run("Sections", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = FindByType[*NodeSection](tree)
		}
	})

	b.Run("Requirements", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = FindByType[*NodeRequirement](tree)
		}
	})

	b.Run("Paragraphs", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = FindByType[*NodeParagraph](tree)
		}
	})
}

func BenchmarkCount(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.Run("AllNodes", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Count(
				tree,
				func(_ Node) bool { return true },
			)
		}
	})

	b.Run("Sections", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Count(
				tree,
				IsType[*NodeSection](),
			)
		}
	})
}

func BenchmarkExists(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.Run("ExistingType", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Exists(
				tree,
				IsType[*NodeCodeBlock](),
			)
		}
	})

	b.Run("NonExistingType", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = Exists(
				tree,
				func(_ Node) bool { return false },
			)
		}
	})
}

// =============================================================================
// Visitor Benchmarks
// =============================================================================

// benchmarkVisitor is a simple visitor that counts nodes
type benchmarkVisitor struct {
	BaseVisitor
	count int
}

func (v *benchmarkVisitor) VisitDocument(
	*NodeDocument,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitSection(
	*NodeSection,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitRequirement(
	*NodeRequirement,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitScenario(
	*NodeScenario,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitParagraph(
	*NodeParagraph,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitList(
	*NodeList,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitListItem(
	*NodeListItem,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitCodeBlock(
	*NodeCodeBlock,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitBlockquote(
	*NodeBlockquote,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitText(
	*NodeText,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitStrong(
	*NodeStrong,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitEmphasis(
	*NodeEmphasis,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitStrikethrough(
	*NodeStrikethrough,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitCode(
	*NodeCode,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitLink(
	*NodeLink,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitLinkDef(
	*NodeLinkDef,
) error {
	v.count++

	return nil
}

func (v *benchmarkVisitor) VisitWikilink(
	*NodeWikilink,
) error {
	v.count++

	return nil
}

// skipCodeBlockVisitor skips code block children
type skipCodeBlockVisitor struct {
	BaseVisitor
	count int
}

func (v *skipCodeBlockVisitor) VisitDocument(
	*NodeDocument,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitSection(
	*NodeSection,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitRequirement(
	*NodeRequirement,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitScenario(
	*NodeScenario,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitParagraph(
	*NodeParagraph,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitList(
	*NodeList,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitListItem(
	*NodeListItem,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitCodeBlock(
	*NodeCodeBlock,
) error {
	v.count++

	return SkipChildren
}

func (v *skipCodeBlockVisitor) VisitBlockquote(
	*NodeBlockquote,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitText(
	*NodeText,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitStrong(
	*NodeStrong,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitEmphasis(
	*NodeEmphasis,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitStrikethrough(
	*NodeStrikethrough,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitCode(
	*NodeCode,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitLink(
	*NodeLink,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitLinkDef(
	*NodeLinkDef,
) error {
	v.count++

	return nil
}

func (v *skipCodeBlockVisitor) VisitWikilink(
	*NodeWikilink,
) error {
	v.count++

	return nil
}

func BenchmarkWalk(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		v := &benchmarkVisitor{}
		_ = Walk(tree, v)
	}
}

func BenchmarkVisitorPattern(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.Run("WalkAll", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			v := &benchmarkVisitor{}
			_ = Walk(tree, v)
		}
	})

	b.Run("WalkWithSkip", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			v := &skipCodeBlockVisitor{}
			_ = Walk(tree, v)
		}
	})
}

// =============================================================================
// Hash Computation Benchmarks
// =============================================================================

func BenchmarkHashComputation(b *testing.B) {
	tree, _ := Parse(largeDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = tree.Hash()
	}
}

func BenchmarkNodeEquality(b *testing.B) {
	tree1, _ := Parse(largeDoc)
	tree2, _ := Parse(largeDoc)

	b.Run("SameTrees", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = tree1.Equal(tree2)
		}
	})

	// Modify one tree
	modifiedDoc := make([]byte, len(largeDoc))
	copy(modifiedDoc, largeDoc)
	modifiedDoc[100] = 'X'
	tree3, _ := Parse(modifiedDoc)

	b.Run("DifferentTrees", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = tree1.Equal(tree3)
		}
	})
}

// =============================================================================
// End-to-End Benchmarks
// =============================================================================

func BenchmarkParseAndPrint(b *testing.B) {
	b.SetBytes(int64(len(mediumDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		tree, _ := Parse(mediumDoc)
		_ = Print(tree)
	}
}

func BenchmarkParseAndQuery(b *testing.B) {
	b.SetBytes(int64(len(mediumDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		tree, _ := Parse(mediumDoc)
		_ = FindByType[*NodeRequirement](tree)
		_ = FindByType[*NodeScenario](tree)
	}
}

func BenchmarkParseAndIndex(b *testing.B) {
	b.SetBytes(int64(len(mediumDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		tree, _ := Parse(mediumDoc)
		idx := NewPositionIndex(tree, mediumDoc)
		_ = idx.NodeAt(len(mediumDoc) / 2)
	}
}

// =============================================================================
// Memory Allocation Benchmarks
// =============================================================================

func BenchmarkParseAllocations(b *testing.B) {
	b.Run("Small", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, _ = Parse(smallDoc)
		}
	})

	b.Run("Medium", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, _ = Parse(mediumDoc)
		}
	})

	b.Run("Large", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, _ = Parse(largeDoc)
		}
	})
}

// BenchmarkConcurrentParse tests thread safety and concurrent performance
func BenchmarkConcurrentParse(b *testing.B) {
	b.SetBytes(int64(len(mediumDoc)))
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = Parse(mediumDoc)
		}
	})
}

// BenchmarkConcurrentPoolAccess tests pool contention under concurrent access
func BenchmarkConcurrentPoolAccess(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tok := GetToken()
			tok.Type = TokenText
			PutToken(tok)

			c := GetChildren(8)
			PutChildren(c)
		}
	})
}
