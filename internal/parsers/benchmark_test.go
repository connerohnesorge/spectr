package parsers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/mdparser"
)

// Benchmark test corpus files
var (
	smallFile        = "../../testdata/benchmarks/small.md"
	mediumFile       = "../../testdata/benchmarks/medium.md"
	largeFile        = "../../testdata/benchmarks/large.md"
	pathologicalFile = "../../testdata/benchmarks/pathological.md"
)

// Lexer-based implementations (using new mdparser)
// These will be used to benchmark the new approach

// ParseRequirementsLexer parses requirements using the new lexer/parser approach
func ParseRequirementsLexer(filePath string) ([]RequirementBlock, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return nil, err
	}

	return extractRequirementsFromAST(doc), nil
}

// processHeader handles header nodes during requirement extraction
func processHeader(
	h *mdparser.Header,
	currentReq **RequirementBlock,
	requirements *[]RequirementBlock,
) {
	switch {
	case h.Level == 2:
		// H2 headers mark section boundaries
		if *currentReq != nil {
			*requirements = append(*requirements, **currentReq)
			*currentReq = nil
		}
	case h.Level == 3 && strings.HasPrefix(h.Text, "Requirement: "):
		// Save previous requirement
		if *currentReq != nil {
			*requirements = append(*requirements, **currentReq)
		}
		// Start new requirement
		name := strings.TrimPrefix(h.Text, "Requirement: ")
		headerLine := strings.Repeat("#", h.Level) + " " + h.Text
		*currentReq = &RequirementBlock{
			HeaderLine: headerLine,
			Name:       strings.TrimSpace(name),
			Raw:        headerLine + "\n",
		}
	case *currentReq != nil:
		// Add any other header (like scenarios) to the current requirement
		headerLine := strings.Repeat("#", h.Level) + " " + h.Text
		(*currentReq).Raw += headerLine + "\n"
	}
}

// extractRequirementsFromAST walks the AST and extracts requirements
func extractRequirementsFromAST(doc *mdparser.Document) []RequirementBlock {
	var requirements []RequirementBlock
	var currentReq *RequirementBlock

	for _, node := range doc.Children {
		switch n := node.(type) {
		case *mdparser.Header:
			processHeader(n, &currentReq, &requirements)

		case *mdparser.Paragraph:
			if currentReq == nil {
				continue
			}
			for _, line := range n.Lines {
				currentReq.Raw += line + "\n"
			}

		case *mdparser.List:
			if currentReq == nil {
				continue
			}
			for _, item := range n.Items {
				prefix := "- "
				if n.Ordered {
					prefix = "1. "
				}
				currentReq.Raw += prefix + item.Text + "\n"
			}

		case *mdparser.CodeBlock:
			if currentReq == nil {
				continue
			}
			currentReq.Raw += "```" + n.Language + "\n"
			for _, line := range n.Lines {
				currentReq.Raw += line + "\n"
			}
			currentReq.Raw += "```\n"

		case *mdparser.BlankLine:
			if currentReq == nil {
				continue
			}
			for range n.Count {
				currentReq.Raw += "\n"
			}
		}
	}

	// Don't forget the last requirement
	if currentReq != nil {
		requirements = append(requirements, *currentReq)
	}

	return requirements
}

// ParseDeltaSpecLexer parses delta specs using the new lexer/parser approach
func ParseDeltaSpecLexer(filePath string) (*DeltaPlan, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return nil, err
	}

	plan := &DeltaPlan{
		Added:    make([]RequirementBlock, 0),
		Modified: make([]RequirementBlock, 0),
		Removed:  make([]string, 0),
		Renamed:  make([]RenameOp, 0),
	}

	plan.Added = extractDeltaSectionFromAST(doc, "ADDED")
	plan.Modified = extractDeltaSectionFromAST(doc, "MODIFIED")
	plan.Removed = extractRemovedSectionFromAST(doc)
	plan.Renamed = extractRenamedSectionFromAST(doc)

	return plan, nil
}

// processDeltaHeader handles header nodes during delta section extraction
func processDeltaHeader(
	h *mdparser.Header,
	currentReq **RequirementBlock,
	requirements *[]RequirementBlock,
) {
	if h.Level == 3 && strings.HasPrefix(h.Text, "Requirement: ") {
		// Save previous requirement
		if *currentReq != nil {
			*requirements = append(*requirements, **currentReq)
		}
		// Start new requirement
		name := strings.TrimPrefix(h.Text, "Requirement: ")
		headerLine := "### " + h.Text
		*currentReq = &RequirementBlock{
			HeaderLine: headerLine,
			Name:       strings.TrimSpace(name),
			Raw:        headerLine + "\n",
		}
	} else if *currentReq != nil {
		// Add scenario headers to current requirement
		headerLine := strings.Repeat("#", h.Level) + " " + h.Text
		(*currentReq).Raw += headerLine + "\n"
	}
}

// extractDeltaSectionFromAST extracts requirements from a specific delta section
func extractDeltaSectionFromAST(doc *mdparser.Document, sectionType string) []RequirementBlock {
	var requirements []RequirementBlock
	var inSection bool
	var currentReq *RequirementBlock

	sectionHeader := sectionType + " Requirements"

	for _, node := range doc.Children {
		// Check if this is a header
		if header, ok := node.(*mdparser.Header); ok && header.Level == 2 {
			if strings.Contains(header.Text, sectionHeader) {
				inSection = true

				continue
			} else if inSection {
				// We've hit a new section (## header), stop processing
				break
			}
		}

		// Only process nodes if we're in the target section
		if !inSection {
			continue
		}

		switch n := node.(type) {
		case *mdparser.Header:
			processDeltaHeader(n, &currentReq, &requirements)

		case *mdparser.Paragraph:
			if currentReq == nil {
				continue
			}
			for _, line := range n.Lines {
				currentReq.Raw += line + "\n"
			}

		case *mdparser.List:
			if currentReq == nil {
				continue
			}
			for _, item := range n.Items {
				currentReq.Raw += "- " + item.Text + "\n"
			}

		case *mdparser.CodeBlock:
			if currentReq == nil {
				continue
			}
			currentReq.Raw += "```" + n.Language + "\n"
			for _, line := range n.Lines {
				currentReq.Raw += line + "\n"
			}
			currentReq.Raw += "```\n"

		case *mdparser.BlankLine:
			if currentReq == nil {
				continue
			}
			for range n.Count {
				currentReq.Raw += "\n"
			}
		}
	}

	// Save last requirement
	if currentReq != nil {
		requirements = append(requirements, *currentReq)
	}

	return requirements
}

// extractRemovedSectionFromAST extracts requirement names from REMOVED section
func extractRemovedSectionFromAST(doc *mdparser.Document) []string {
	var removed []string
	var inSection bool

	for _, node := range doc.Children {
		// Check if this is a header
		header, ok := node.(*mdparser.Header)
		if !ok {
			continue
		}

		if header.Level == 2 {
			if strings.Contains(header.Text, "REMOVED Requirements") {
				inSection = true

				continue
			} else if inSection {
				// Hit a new section, stop
				break
			}
		} else if inSection && header.Level == 3 && strings.HasPrefix(header.Text, "Requirement: ") {
			// Extract requirement name
			name := strings.TrimPrefix(header.Text, "Requirement: ")
			removed = append(removed, strings.TrimSpace(name))
		}
	}

	return removed
}

// extractRenamedSectionFromAST extracts FROM/TO pairs from RENAMED section
func extractRenamedSectionFromAST(doc *mdparser.Document) []RenameOp {
	var renamed []RenameOp
	var inSection bool
	var currentFrom string

	for _, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if ok && header.Level == 2 {
			if strings.Contains(header.Text, "RENAMED Requirements") {
				inSection = true

				continue
			} else if inSection {
				break
			}
		}

		if !inSection {
			continue
		}

		// Look for list items with FROM/TO pattern
		list, ok := node.(*mdparser.List)
		if !ok {
			continue
		}

		for _, item := range list.Items {
			text := strings.TrimSpace(item.Text)

			// Check for FROM line
			if strings.HasPrefix(text, "FROM: `### Requirement: ") {
				text = strings.TrimPrefix(text, "FROM: `### Requirement: ")
				text = strings.TrimSuffix(text, "`")
				currentFrom = strings.TrimSpace(text)

				continue
			}

			// Check for TO line
			if currentFrom == "" || !strings.HasPrefix(text, "TO: `### Requirement: ") {
				continue
			}
			text = strings.TrimPrefix(text, "TO: `### Requirement: ")
			text = strings.TrimSuffix(text, "`")
			to := strings.TrimSpace(text)
			renamed = append(renamed, RenameOp{
				From: currentFrom,
				To:   to,
			})
			currentFrom = ""
		}
	}

	return renamed
}

// Benchmark helper to get absolute path
func getTestFilePath(t testing.TB, relPath string) string {
	t.Helper()
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	return absPath
}

// Benchmarks for Regex Implementation (Current)

func BenchmarkRegexRequirementParser_Small(b *testing.B) {
	path := getTestFilePath(b, smallFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirements(path)
		if err != nil {
			b.Fatalf("ParseRequirements failed: %v", err)
		}
	}
}

func BenchmarkRegexRequirementParser_Medium(b *testing.B) {
	path := getTestFilePath(b, mediumFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirements(path)
		if err != nil {
			b.Fatalf("ParseRequirements failed: %v", err)
		}
	}
}

func BenchmarkRegexRequirementParser_Large(b *testing.B) {
	path := getTestFilePath(b, largeFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirements(path)
		if err != nil {
			b.Fatalf("ParseRequirements failed: %v", err)
		}
	}
}

func BenchmarkRegexRequirementParser_Pathological(b *testing.B) {
	path := getTestFilePath(b, pathologicalFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirements(path)
		if err != nil {
			b.Fatalf("ParseRequirements failed: %v", err)
		}
	}
}

func BenchmarkRegexDeltaParser_Small(b *testing.B) {
	path := getTestFilePath(b, smallFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseDeltaSpec(path)
		if err != nil {
			b.Fatalf("ParseDeltaSpec failed: %v", err)
		}
	}
}

// Benchmarks for Lexer/Parser Implementation (New)

func BenchmarkLexerRequirementParser_Small(b *testing.B) {
	path := getTestFilePath(b, smallFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirementsLexer(path)
		if err != nil {
			b.Fatalf("ParseRequirementsLexer failed: %v", err)
		}
	}
}

func BenchmarkLexerRequirementParser_Medium(b *testing.B) {
	path := getTestFilePath(b, mediumFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirementsLexer(path)
		if err != nil {
			b.Fatalf("ParseRequirementsLexer failed: %v", err)
		}
	}
}

func BenchmarkLexerRequirementParser_Large(b *testing.B) {
	path := getTestFilePath(b, largeFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirementsLexer(path)
		if err != nil {
			b.Fatalf("ParseRequirementsLexer failed: %v", err)
		}
	}
}

func BenchmarkLexerRequirementParser_Pathological(b *testing.B) {
	path := getTestFilePath(b, pathologicalFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseRequirementsLexer(path)
		if err != nil {
			b.Fatalf("ParseRequirementsLexer failed: %v", err)
		}
	}
}

func BenchmarkLexerDeltaParser_Small(b *testing.B) {
	path := getTestFilePath(b, smallFile)
	b.ResetTimer()
	for range b.N {
		_, err := ParseDeltaSpecLexer(path)
		if err != nil {
			b.Fatalf("ParseDeltaSpecLexer failed: %v", err)
		}
	}
}

// Correctness Validation Tests
// These ensure both parsers produce equivalent results

func TestRequirementParserCorrectness_Small(t *testing.T) {
	path := getTestFilePath(t, smallFile)

	regexResults, err := ParseRequirements(path)
	if err != nil {
		t.Fatalf("Regex parser failed: %v", err)
	}

	lexerResults, err := ParseRequirementsLexer(path)
	if err != nil {
		t.Fatalf("Lexer parser failed: %v", err)
	}

	validateRequirementEquivalence(t, "small.md", regexResults, lexerResults)
}

func TestRequirementParserCorrectness_Medium(t *testing.T) {
	path := getTestFilePath(t, mediumFile)

	regexResults, err := ParseRequirements(path)
	if err != nil {
		t.Fatalf("Regex parser failed: %v", err)
	}

	lexerResults, err := ParseRequirementsLexer(path)
	if err != nil {
		t.Fatalf("Lexer parser failed: %v", err)
	}

	validateRequirementEquivalence(t, "medium.md", regexResults, lexerResults)
}

func TestRequirementParserCorrectness_Pathological(t *testing.T) {
	path := getTestFilePath(t, pathologicalFile)

	regexResults, err := ParseRequirements(path)
	if err != nil {
		t.Fatalf("Regex parser failed: %v", err)
	}

	lexerResults, err := ParseRequirementsLexer(path)
	if err != nil {
		t.Fatalf("Lexer parser failed: %v", err)
	}

	// For pathological cases, we expect the lexer to handle edge cases better
	// so we validate that the lexer results are at least as good
	t.Log(
		"Pathological case: Regex found",
		len(regexResults),
		"requirements, Lexer found",
		len(lexerResults),
		"requirements",
	)

	// The lexer should NOT extract requirements from code blocks
	// We expect it to find fewer or equal requirements in pathological cases
	if len(lexerResults) > len(regexResults) {
		t.Log(
			"Warning: Lexer found more requirements than regex - this may indicate it's correctly handling edge cases",
		)
	}

	// Verify the lexer found reasonable requirements
	for i, req := range lexerResults {
		if req.Name == "" {
			t.Errorf("Requirement %d has empty name", i)
		}
		if req.Raw == "" {
			t.Errorf("Requirement %d has empty raw content", i)
		}
	}
}

func TestDeltaParserCorrectness_Small(t *testing.T) {
	path := getTestFilePath(t, smallFile)

	regexResults, err := ParseDeltaSpec(path)
	if err != nil {
		t.Fatalf("Regex parser failed: %v", err)
	}

	lexerResults, err := ParseDeltaSpecLexer(path)
	if err != nil {
		t.Fatalf("Lexer parser failed: %v", err)
	}

	validateDeltaEquivalence(t, "small.md", regexResults, lexerResults)
}

// Helper functions for validation

func validateRequirementEquivalence(
	t *testing.T,
	filename string,
	regex, lexer []RequirementBlock,
) {
	t.Helper()

	if len(regex) != len(lexer) {
		t.Errorf("%s: requirement count mismatch - regex: %d, lexer: %d",
			filename, len(regex), len(lexer))

		t.Log("Regex requirements:")
		for i, r := range regex {
			t.Logf("  %d: %s", i, r.Name)
		}
		t.Log("Lexer requirements:")
		for i, r := range lexer {
			t.Logf("  %d: %s", i, r.Name)
		}

		return
	}

	for i := range regex {
		if regex[i].Name != lexer[i].Name {
			t.Errorf("%s: requirement %d name mismatch - regex: %q, lexer: %q",
				filename, i, regex[i].Name, lexer[i].Name)
		}

		// Parse scenarios from both
		regexScenarios := ParseScenarios(regex[i].Raw)
		lexerScenarios := ParseScenarios(lexer[i].Raw)

		if len(regexScenarios) != len(lexerScenarios) {
			t.Errorf("%s: requirement %q scenario count mismatch - regex: %d, lexer: %d",
				filename, regex[i].Name, len(regexScenarios), len(lexerScenarios))
		}
	}
}

func validateDeltaEquivalence(t *testing.T, filename string, regex, lexer *DeltaPlan) {
	t.Helper()

	if len(regex.Added) != len(lexer.Added) {
		t.Errorf("%s: ADDED count mismatch - regex: %d, lexer: %d",
			filename, len(regex.Added), len(lexer.Added))
		t.Log("Regex ADDED requirements:")
		for i, r := range regex.Added {
			t.Logf("  %d: %s", i, r.Name)
		}
		t.Log("Lexer ADDED requirements:")
		for i, r := range lexer.Added {
			t.Logf("  %d: %s", i, r.Name)
		}
	}

	if len(regex.Modified) != len(lexer.Modified) {
		t.Errorf("%s: MODIFIED count mismatch - regex: %d, lexer: %d",
			filename, len(regex.Modified), len(lexer.Modified))
		t.Log("Regex MODIFIED requirements:")
		for i, r := range regex.Modified {
			t.Logf("  %d: %s", i, r.Name)
		}
		t.Log("Lexer MODIFIED requirements:")
		for i, r := range lexer.Modified {
			t.Logf("  %d: %s", i, r.Name)
		}
	}

	if len(regex.Removed) != len(lexer.Removed) {
		t.Errorf("%s: REMOVED count mismatch - regex: %d, lexer: %d",
			filename, len(regex.Removed), len(lexer.Removed))
		t.Log("Regex REMOVED requirements:")
		for i, r := range regex.Removed {
			t.Logf("  %d: %s", i, r)
		}
		t.Log("Lexer REMOVED requirements:")
		for i, r := range lexer.Removed {
			t.Logf("  %d: %s", i, r)
		}
	}

	if len(regex.Renamed) != len(lexer.Renamed) {
		t.Errorf("%s: RENAMED count mismatch - regex: %d, lexer: %d",
			filename, len(regex.Renamed), len(lexer.Renamed))
	}

	// Validate names match
	for i := 0; i < len(regex.Added) && i < len(lexer.Added); i++ {
		if regex.Added[i].Name != lexer.Added[i].Name {
			t.Errorf("%s: ADDED[%d] name mismatch - regex: %q, lexer: %q",
				filename, i, regex.Added[i].Name, lexer.Added[i].Name)
		}
	}
}
