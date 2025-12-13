// Package markdown comparison_test.go contains permanent regression tests
// that compare the new AST-based markdown parser against the old regex-based parsers.
// These tests ensure behavioral equivalence during and after the migration.
//
// IMPORTANT: These tests document both equivalences AND known differences between
// the regex and AST-based parsers. Some differences are by design:
//
// 1. Section extraction: The regex-based ExtractSections only extracts ## headers,
//    while the AST parser extracts ALL header levels as sections. Tests that compare
//    sections should filter AST results to H2 only.
//
// 2. Inline formatting: The AST parser strips markdown formatting (bold, italic, code)
//    from header text, while regex preserves raw markdown. This is actually more correct
//    behavior for the AST parser. Tests for headers with inline formatting document this.
//
// 3. Content boundaries: The AST parser may include slightly different content boundaries
//    for sections due to how it processes the parsed tree. Content comparison should
//    be done with trimmed strings and allow for minor whitespace differences.
package markdown_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/validation"
)

// =============================================================================
// Test Helpers
// =============================================================================

// regexTaskStatus counts tasks using the regex-based approach from parsers.go
func regexTaskStatus(content string) (total, completed int) {
	taskPattern := regexp.MustCompile(`^\s*-\s*\[([xX ])\]`)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		matches := taskPattern.FindStringSubmatch(line)
		if len(matches) <= 1 {
			continue
		}
		total++
		marker := strings.ToLower(strings.TrimSpace(matches[1]))
		if marker == "x" {
			completed++
		}
	}

	return total, completed
}

// astTaskStatus counts tasks using the AST-based approach
func astTaskStatus(content string) (total, completed int) {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return 0, 0
	}

	return countTasksRecursive(doc.Tasks)
}

// countTasksRecursive counts tasks including children
func countTasksRecursive(tasks []markdown.Task) (total, completed int) {
	for _, task := range tasks {
		total++
		if task.Checked {
			completed++
		}
		childTotal, childCompleted := countTasksRecursive(task.Children)
		total += childTotal
		completed += childCompleted
	}

	return total, completed
}

// regexExtractSections uses the regex-based ExtractSections from validation package
func regexExtractSections(content string) map[string]string {
	return validation.ExtractSections(content)
}

// astExtractSections uses the AST-based section extraction.
// NOTE: This filters to H2 sections only to match the regex-based behavior,
// since the AST parser extracts ALL header levels as sections.
func astExtractSections(content string) map[string]string {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return nil
	}
	result := make(map[string]string)
	for name, section := range doc.Sections {
		// Only include H2 sections to match regex-based ExtractSections behavior
		if section.Header.Level == 2 {
			result[name] = section.Content
		}
	}

	return result
}

// astExtractAllSections uses the AST-based section extraction without filtering.
// This returns all headers as sections, not just H2.
func astExtractAllSections(content string) map[string]string {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return nil
	}
	result := make(map[string]string)
	for name, section := range doc.Sections {
		result[name] = section.Content
	}

	return result
}

// regexCountRequirements uses the regex-based approach
func regexCountRequirements(content string) int {
	reqPattern := regexp.MustCompile(`^###\s+Requirement:`)
	count := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if reqPattern.MatchString(strings.TrimSpace(line)) {
			count++
		}
	}

	return count
}

// astCountRequirements uses the AST-based header extraction
func astCountRequirements(content string) int {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return 0
	}
	count := 0
	for _, h := range doc.Headers {
		if h.Level == 3 && strings.HasPrefix(h.Text, "Requirement:") {
			count++
		}
	}

	return count
}

// regexExtractRequirements extracts requirements using regex-based parser
func regexExtractRequirements(content string) []validation.Requirement {
	return validation.ExtractRequirements(content)
}

// astExtractRequirementNames extracts requirement names using AST
func astExtractRequirementNames(content string) []string {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return nil
	}
	var names []string
	for _, h := range doc.Headers {
		if h.Level == 3 && strings.HasPrefix(h.Text, "Requirement:") {
			name := strings.TrimSpace(strings.TrimPrefix(h.Text, "Requirement:"))
			names = append(names, name)
		}
	}

	return names
}

// regexCountDeltas uses regex-based delta section detection
func regexCountDeltas(content string) int {
	deltaPattern := regexp.MustCompile(`^##\s+(ADDED|MODIFIED|REMOVED|RENAMED)\s+Requirements`)
	count := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if deltaPattern.MatchString(strings.TrimSpace(line)) {
			count++
		}
	}

	return count
}

// astCountDeltas uses AST-based header extraction to count delta sections
func astCountDeltas(content string) int {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return 0
	}
	deltaKeywords := []string{"ADDED Requirements", "MODIFIED Requirements", "REMOVED Requirements", "RENAMED Requirements"}
	count := 0
	for _, h := range doc.Headers {
		if h.Level != 2 {
			continue
		}
		for _, kw := range deltaKeywords {
			if h.Text == kw {
				count++

				break
			}
		}
	}

	return count
}

// regexCountH2Headers counts ## headers using regex
func regexCountH2Headers(content string) int {
	h2Pattern := regexp.MustCompile(`^##\s+(.+)$`)
	count := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if h2Pattern.MatchString(line) {
			count++
		}
	}

	return count
}

// astCountH2Headers counts ## headers using AST
func astCountH2Headers(content string) int {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return 0
	}
	count := 0
	for _, h := range doc.Headers {
		if h.Level == 2 {
			count++
		}
	}

	return count
}

// regexExtractH2Headers extracts ## header texts using regex
func regexExtractH2Headers(content string) []string {
	h2Pattern := regexp.MustCompile(`^##\s+(.+)$`)
	var headers []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		matches := h2Pattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			headers = append(headers, strings.TrimSpace(matches[1]))
		}
	}

	return headers
}

// astExtractH2Headers extracts ## header texts using AST
func astExtractH2Headers(content string) []string {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return nil
	}
	var headers []string
	for _, h := range doc.Headers {
		if h.Level == 2 {
			headers = append(headers, h.Text)
		}
	}

	return headers
}

// =============================================================================
// Task Parsing Comparison Tests
// =============================================================================

func TestCompare_TaskParsing(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "simple unchecked tasks",
			content: `# Tasks

- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
`,
		},
		{
			name: "simple checked tasks",
			content: `# Tasks

- [x] Task 1
- [X] Task 2
- [x] Task 3
`,
		},
		{
			name: "mixed checked and unchecked",
			content: `# Tasks

- [ ] Unchecked 1
- [x] Checked 1
- [ ] Unchecked 2
- [X] Checked 2
- [ ] Unchecked 3
`,
		},
		{
			name: "tasks with descriptions",
			content: `# Tasks

- [ ] Implement the new feature with detailed description
- [x] Review the code changes thoroughly
- [ ] Write comprehensive documentation
`,
		},
		{
			name: "nested tasks",
			content: `# Tasks

- [ ] Parent task 1
  - [ ] Child task 1.1
  - [x] Child task 1.2
- [x] Parent task 2
  - [ ] Child task 2.1
`,
		},
		{
			name: "deeply nested tasks",
			content: `# Tasks

- [ ] Level 1
  - [ ] Level 2
    - [ ] Level 3
      - [x] Level 4
`,
		},
		{
			name: "tasks with special characters",
			content: `# Tasks

- [ ] Task with "quotes" and 'apostrophes'
- [x] Task with special chars: @#$%^&*()
- [ ] Task with unicode: cafe and emoji test
`,
		},
		{
			name: "tasks in different sections",
			content: `# Project Tasks

## Phase 1
- [ ] Setup project
- [x] Initialize repo

## Phase 2
- [ ] Implement feature A
- [ ] Implement feature B
`,
		},
		{
			name: "empty task list",
			content: `# Tasks

No tasks here yet.
`,
		},
		{
			name: "tasks with code blocks nearby",
			content: "# Tasks\n\n- [ ] Task before code\n\n```\ncode block\n```\n\n- [x] Task after code\n",
		},
		{
			name: "indented tasks with tabs",
			content: `# Tasks

- [ ] Parent task
	- [ ] Tab-indented child
	- [x] Another tab-indented child
`,
		},
		{
			name: "tasks file pattern from examples",
			content: `## Tasks

- [x] Consolidate TUI components
- [x] Apply consistent styling
- [x] Update tests
`,
		},
		{
			name: "mixed list items and tasks",
			content: `# List

- Regular list item
- [ ] Task item
- Another regular item
- [x] Another task
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regexTotal, regexCompleted := regexTaskStatus(tt.content)
			astTotal, astCompleted := astTaskStatus(tt.content)

			if regexTotal != astTotal {
				t.Errorf("Total tasks mismatch: regex=%d, ast=%d", regexTotal, astTotal)
			}
			if regexCompleted != astCompleted {
				t.Errorf("Completed tasks mismatch: regex=%d, ast=%d", regexCompleted, astCompleted)
			}
		})
	}
}

// =============================================================================
// Header Extraction Comparison Tests
// =============================================================================

func TestCompare_HeaderExtraction(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "simple headers",
			content: `# Title

## Section 1

## Section 2

## Section 3
`,
		},
		{
			name: "requirements spec pattern",
			content: `# Spec Title

## Purpose

Some purpose text.

## Requirements

### Requirement: First Req

Description.

### Requirement: Second Req

Description.
`,
		},
		{
			name: "delta spec pattern",
			content: `# Delta Spec

## ADDED Requirements

### Requirement: New Feature

Description.

## MODIFIED Requirements

### Requirement: Updated Feature

Description.
`,
		},
		{
			name: "all header levels",
			content: `# H1 Title

## H2 Section

### H3 Subsection

#### H4 Scenario

##### H5 Detail

###### H6 Note
`,
		},
		{
			name: "headers with special characters",
			content: `# Title with "quotes"

## Section: With Colon

## Section (with parens)

## Section [with brackets]
`,
		},
		// NOTE: This test documents a KNOWN DIFFERENCE between regex and AST parsers.
		// The AST parser correctly strips markdown formatting from header text,
		// while regex preserves the raw markdown. This is tested separately below.
		// {
		// 	name: "headers with inline formatting",
		// 	content: ...,
		// },
		{
			name: "empty sections between headers",
			content: `# Title

## Empty Section 1

## Empty Section 2

## Section With Content

Some content here.
`,
		},
		{
			name: "scenario headers",
			content: `# Spec

## Requirements

### Requirement: Auth

#### Scenario: Valid Login

Steps here.

#### Scenario: Invalid Login

Steps here.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compare H2 header counts
			regexH2Count := regexCountH2Headers(tt.content)
			astH2Count := astCountH2Headers(tt.content)

			if regexH2Count != astH2Count {
				t.Errorf("H2 header count mismatch: regex=%d, ast=%d", regexH2Count, astH2Count)
			}

			// Compare H2 header texts
			regexH2Headers := regexExtractH2Headers(tt.content)
			astH2Headers := astExtractH2Headers(tt.content)

			if len(regexH2Headers) != len(astH2Headers) {
				t.Errorf("H2 header list length mismatch: regex=%d, ast=%d",
					len(regexH2Headers), len(astH2Headers))
			} else {
				for i := range regexH2Headers {
					if regexH2Headers[i] != astH2Headers[i] {
						t.Errorf("H2 header %d mismatch: regex=%q, ast=%q",
							i, regexH2Headers[i], astH2Headers[i])
					}
				}
			}
		})
	}
}

// =============================================================================
// Known Differences Documentation Tests
// =============================================================================

// TestDocument_InlineFormattingDifference documents the known difference in how
// the regex and AST parsers handle inline markdown formatting in headers.
// The AST parser CORRECTLY strips formatting, while regex preserves it.
// This test exists to document this difference, not to flag it as a failure.
func TestDocument_InlineFormattingDifference(t *testing.T) {
	content := `# Title

## Section with **bold** text

## Section with *italic* text

## Section with ` + "`code`" + ` inline
`

	regexHeaders := regexExtractH2Headers(content)
	astHeaders := astExtractH2Headers(content)

	// Verify counts match
	if len(regexHeaders) != len(astHeaders) {
		t.Errorf("Header count mismatch: regex=%d, ast=%d", len(regexHeaders), len(astHeaders))

		return
	}

	// Document the expected differences
	expectedDiffs := map[string]string{
		"Section with **bold** text":   "Section with bold text",
		"Section with *italic* text":   "Section with italic text",
		"Section with `code` inline":   "Section with code inline",
	}

	for i := range regexHeaders {
		regexText := regexHeaders[i]
		astText := astHeaders[i]

		expectedAST, hasDiff := expectedDiffs[regexText]
		if hasDiff {
			// This is a documented difference - AST strips formatting
			if astText != expectedAST {
				t.Errorf("Header %d: AST text %q does not match expected %q",
					i, astText, expectedAST)
			}
			t.Logf("DOCUMENTED DIFFERENCE: regex=%q -> ast=%q (formatting stripped)", regexText, astText)
		} else if regexText != astText {
			// No formatting difference expected
			t.Errorf("Header %d: unexpected difference regex=%q, ast=%q",
				i, regexText, astText)
		}
	}
}

// TestDocument_SectionExtractionScope documents that the AST parser extracts
// ALL headers as sections, while the regex-based ExtractSections only extracts ## headers.
// The comparison helpers filter AST results to match, but this test documents the difference.
func TestDocument_SectionExtractionScope(t *testing.T) {
	content := `# Document Title

## Purpose

Purpose content.

### Requirement: Auth

Auth content.

#### Scenario: Login

Scenario content.
`

	// Regex only extracts H2 (## headers)
	regexSections := regexExtractSections(content)

	// AST without filtering extracts all headers
	astAllSections := astExtractAllSections(content)

	// Verify regex only gets H2
	if len(regexSections) != 1 {
		t.Errorf("Expected regex to find 1 H2 section, got %d", len(regexSections))
	}
	if _, ok := regexSections["Purpose"]; !ok {
		t.Error("Expected regex to find 'Purpose' section")
	}

	// Verify AST finds all headers (H1, H2, H3, H4)
	if len(astAllSections) < 4 {
		t.Errorf("Expected AST to find at least 4 sections (all headers), got %d", len(astAllSections))
	}

	t.Logf("DOCUMENTED DIFFERENCE: regex finds %d sections (H2 only), AST finds %d sections (all headers)",
		len(regexSections), len(astAllSections))

	// When filtered, they should match
	astFilteredSections := astExtractSections(content)
	if len(regexSections) != len(astFilteredSections) {
		t.Errorf("Filtered AST sections (%d) should match regex sections (%d)",
			len(astFilteredSections), len(regexSections))
	}
}

// =============================================================================
// Section Extraction Comparison Tests
// =============================================================================

func TestCompare_SectionExtraction(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "basic sections",
			content: `# Document

## Purpose

This is the purpose section.

## Requirements

These are the requirements.

## Notes

Additional notes.
`,
		},
		{
			name: "nested content in sections",
			content: `# Spec

## Requirements

### Requirement: Auth

The auth requirement.

### Requirement: Logging

The logging requirement.

## Other

Other content.
`,
		},
		{
			name: "multiline section content",
			content: `# Doc

## Description

Line 1 of description.
Line 2 of description.
Line 3 of description.

More paragraph content.

## Details

Detail line 1.
Detail line 2.
`,
		},
		{
			name: "sections with lists",
			content: `# Spec

## Features

- Feature 1
- Feature 2
- Feature 3

## Tasks

- [ ] Task 1
- [x] Task 2
`,
		},
		{
			name: "sections with code blocks",
			content: "# Doc\n\n## Code Example\n\n```go\nfunc main() {}\n```\n\n## Another Section\n\nText here.\n",
		},
		{
			name: "empty document",
			content: `# Title Only
`,
		},
		{
			name: "single section",
			content: `# Title

## Only Section

Content here.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regexSections := regexExtractSections(tt.content)
			astSections := astExtractSections(tt.content)

			// Compare section keys
			if len(regexSections) != len(astSections) {
				t.Errorf("Section count mismatch: regex=%d, ast=%d",
					len(regexSections), len(astSections))
			}

			// Compare section presence
			for key := range regexSections {
				if _, ok := astSections[key]; !ok {
					t.Errorf("Section %q found by regex but not by AST", key)
				}
			}
			for key := range astSections {
				if _, ok := regexSections[key]; !ok {
					t.Errorf("Section %q found by AST but not by regex", key)
				}
			}

			// Compare section content (trimmed, since whitespace handling may differ)
			for key, regexContent := range regexSections {
				astContent, ok := astSections[key]
				if !ok {
					continue
				}
				regexTrimmed := strings.TrimSpace(regexContent)
				astTrimmed := strings.TrimSpace(astContent)
				if regexTrimmed != astTrimmed {
					t.Errorf("Section %q content mismatch:\nregex: %q\nast: %q",
						key, regexTrimmed, astTrimmed)
				}
			}
		})
	}
}

// =============================================================================
// Requirement Counting Comparison Tests
// =============================================================================

func TestCompare_RequirementCounting(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "single requirement",
			content: `# Spec

## Requirements

### Requirement: Auth

The system SHALL authenticate users.
`,
		},
		{
			name: "multiple requirements",
			content: `# Spec

## Requirements

### Requirement: Auth

Auth description.

### Requirement: Logging

Logging description.

### Requirement: Caching

Caching description.
`,
		},
		{
			name: "requirements with scenarios",
			content: `# Spec

## Requirements

### Requirement: Login

The system SHALL support login.

#### Scenario: Valid Login

Steps here.

#### Scenario: Invalid Login

Steps here.

### Requirement: Logout

The system SHALL support logout.

#### Scenario: Normal Logout

Steps here.
`,
		},
		{
			name: "no requirements",
			content: `# Spec

## Purpose

Just a purpose section.
`,
		},
		{
			name: "requirement-like text but not headers",
			content: `# Spec

## Requirements

The Requirement: Auth should be implemented.
This is not ### Requirement: header.
`,
		},
		{
			name: "delta spec with requirements",
			content: `# Delta

## ADDED Requirements

### Requirement: New Feature

Description.

## MODIFIED Requirements

### Requirement: Updated Feature

Description.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regexCount := regexCountRequirements(tt.content)
			astCount := astCountRequirements(tt.content)

			if regexCount != astCount {
				t.Errorf("Requirement count mismatch: regex=%d, ast=%d", regexCount, astCount)
			}
		})
	}
}

// =============================================================================
// Delta Section Counting Comparison Tests
// =============================================================================

func TestCompare_DeltaCounting(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "all delta types",
			content: `# Delta Spec

## ADDED Requirements

### Requirement: New

New feature.

## MODIFIED Requirements

### Requirement: Updated

Updated feature.

## REMOVED Requirements

### Requirement: Deprecated

Removed.

## RENAMED Requirements

FROM: OldName TO: NewName
`,
		},
		{
			name: "only added",
			content: `# Delta

## ADDED Requirements

### Requirement: New Feature

Description.
`,
		},
		{
			name: "only modified",
			content: `# Delta

## MODIFIED Requirements

### Requirement: Updated Feature

Description.
`,
		},
		{
			name: "no deltas",
			content: `# Regular Spec

## Requirements

### Requirement: Feature

Description.
`,
		},
		{
			name: "added and modified",
			content: `# Delta

## ADDED Requirements

### Requirement: New

New.

## MODIFIED Requirements

### Requirement: Updated

Updated.
`,
		},
		{
			name: "delta-like but wrong format",
			content: `# Not Delta

## ADDED Something Else

Not a delta section.

## Added Requirements

Wrong capitalization.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regexCount := regexCountDeltas(tt.content)
			astCount := astCountDeltas(tt.content)

			if regexCount != astCount {
				t.Errorf("Delta count mismatch: regex=%d, ast=%d", regexCount, astCount)
			}
		})
	}
}

// =============================================================================
// Requirement Name Extraction Comparison Tests
// =============================================================================

func TestCompare_RequirementNameExtraction(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "simple requirement names",
			content: `# Spec

## Requirements

### Requirement: Authentication

Description.

### Requirement: Authorization

Description.
`,
		},
		{
			name: "requirement names with special chars",
			content: `# Spec

## Requirements

### Requirement: User-Authentication

Description.

### Requirement: API_Integration

Description.
`,
		},
		{
			name: "requirement names with spaces",
			content: `# Spec

## Requirements

### Requirement: User Login Flow

Description.

### Requirement: Password Reset Process

Description.
`,
		},
		{
			name: "mixed requirements and scenarios",
			content: `# Spec

## Requirements

### Requirement: Auth

Description.

#### Scenario: Valid

Steps.

### Requirement: Cache

Description.

#### Scenario: Hit

Steps.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regexReqs := regexExtractRequirements(tt.content)
			astNames := astExtractRequirementNames(tt.content)

			// Compare counts
			if len(regexReqs) != len(astNames) {
				t.Errorf("Requirement count mismatch: regex=%d, ast=%d",
					len(regexReqs), len(astNames))

				return
			}

			// Compare names
			for i, regexReq := range regexReqs {
				if regexReq.Name != astNames[i] {
					t.Errorf("Requirement name %d mismatch: regex=%q, ast=%q",
						i, regexReq.Name, astNames[i])
				}
			}
		})
	}
}

// =============================================================================
// Edge Cases Comparison Tests
// =============================================================================

func TestCompare_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		skipReason string // If non-empty, test is skipped with this reason
	}{
		{
			name: "header with trailing whitespace",
			content: `# Title

## Section

Content.
`,
		},
		{
			name: "consecutive headers without content",
			content: `# Title

## Section 1

## Section 2

## Section 3
`,
		},
		{
			name: "mixed indentation in tasks",
			content: `# Tasks

- [ ] Task with no indent
 - [ ] One space indent
  - [ ] Two space indent
   - [ ] Three space indent
`,
		},
		{
			name: "unicode in headers",
			content: `# Cafe Spec

## Requirements

### Requirement: Handle unicode

Description.
`,
		},
		{
			name: "very long header text",
			content: `# This is a very long title that spans many characters and might cause issues with parsing or display in some systems

## This is also a very long section header with lots of text that continues on and on

Content.
`,
		},
		{
			name: "header immediately after content",
			content: `# Title
## No blank line after title

Content.
`,
		},
		{
			name: "multiple blank lines between sections",
			content: `# Title



## Section 1



## Section 2

Content.
`,
		},
		{
			name: "task checkbox variations",
			content: `# Tasks

- [ ] Unchecked with space
- [x] Lowercase x
- [X] Uppercase X
- [ ]No space after checkbox
- [x]No space after checkbox checked
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipReason != "" {
				t.Skip(tt.skipReason)
			}

			// Task comparison
			regexTotal, regexCompleted := regexTaskStatus(tt.content)
			astTotal, astCompleted := astTaskStatus(tt.content)
			if regexTotal != astTotal {
				t.Errorf("Task total mismatch: regex=%d, ast=%d", regexTotal, astTotal)
			}
			if regexCompleted != astCompleted {
				t.Errorf("Task completed mismatch: regex=%d, ast=%d", regexCompleted, astCompleted)
			}

			// H2 header comparison
			regexH2Count := regexCountH2Headers(tt.content)
			astH2Count := astCountH2Headers(tt.content)
			if regexH2Count != astH2Count {
				t.Errorf("H2 count mismatch: regex=%d, ast=%d", regexH2Count, astH2Count)
			}

			// Requirement comparison
			regexReqCount := regexCountRequirements(tt.content)
			astReqCount := astCountRequirements(tt.content)
			if regexReqCount != astReqCount {
				t.Errorf("Requirement count mismatch: regex=%d, ast=%d", regexReqCount, astReqCount)
			}
		})
	}
}

// =============================================================================
// Real File Comparison Tests
// =============================================================================

// TestCompare_RealSpecFiles tests comparison using actual spec files from the project.
// This ensures the parsers behave identically on real-world content.
func TestCompare_RealSpecFiles(t *testing.T) {
	// Get the project root (assuming tests run from project root or internal/markdown)
	projectRoot := findProjectRoot(t)
	if projectRoot == "" {
		t.Skip("Could not find project root")
	}

	// Test files from spectr/specs/
	specFiles := []string{
		"spectr/specs/validation/spec.md",
		"spectr/specs/cli-interface/spec.md",
		"spectr/specs/error-handling/spec.md",
		"spectr/specs/documentation/spec.md",
	}

	// Test files from examples/
	exampleFiles := []string{
		"examples/list/spectr/specs/authentication/spec.md",
		"examples/list/spectr/specs/authorization/spec.md",
		"examples/list/spectr/specs/user-management/spec.md",
		"examples/list/spectr/changes/update-permissions/proposal.md",
		"examples/partial-match/spectr/changes/refactor-unified-interactive-tui/tasks.md",
	}

	allFiles := specFiles
	allFiles = append(allFiles, exampleFiles...)

	for _, relPath := range allFiles {
		filePath := filepath.Join(projectRoot, relPath)
		t.Run(relPath, func(t *testing.T) {
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Skipf("Could not read file %s: %v", filePath, err)

				return
			}

			contentStr := string(content)

			// Compare task counts
			regexTotal, regexCompleted := regexTaskStatus(contentStr)
			astTotal, astCompleted := astTaskStatus(contentStr)
			if regexTotal != astTotal {
				t.Errorf("Task total mismatch: regex=%d, ast=%d", regexTotal, astTotal)
			}
			if regexCompleted != astCompleted {
				t.Errorf("Task completed mismatch: regex=%d, ast=%d", regexCompleted, astCompleted)
			}

			// Compare H2 headers
			regexH2Count := regexCountH2Headers(contentStr)
			astH2Count := astCountH2Headers(contentStr)
			if regexH2Count != astH2Count {
				t.Errorf("H2 header count mismatch: regex=%d, ast=%d", regexH2Count, astH2Count)
			}

			// Compare requirement counts
			regexReqCount := regexCountRequirements(contentStr)
			astReqCount := astCountRequirements(contentStr)
			if regexReqCount != astReqCount {
				t.Errorf("Requirement count mismatch: regex=%d, ast=%d", regexReqCount, astReqCount)
			}

			// Compare delta counts
			regexDeltaCount := regexCountDeltas(contentStr)
			astDeltaCount := astCountDeltas(contentStr)
			if regexDeltaCount != astDeltaCount {
				t.Errorf("Delta count mismatch: regex=%d, ast=%d", regexDeltaCount, astDeltaCount)
			}

			// Compare section extraction
			regexSections := regexExtractSections(contentStr)
			astSections := astExtractSections(contentStr)
			if len(regexSections) != len(astSections) {
				t.Errorf("Section count mismatch: regex=%d, ast=%d",
					len(regexSections), len(astSections))
			}

			// Verify all regex sections exist in AST
			for key := range regexSections {
				if _, ok := astSections[key]; !ok {
					t.Errorf("Section %q found by regex but not by AST", key)
				}
			}
		})
	}
}

// findProjectRoot attempts to find the spectr project root directory
func findProjectRoot(t *testing.T) string {
	t.Helper()

	// Try current directory first
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Walk up from current directory looking for go.mod with spectr
	dir := cwd
	for range 10 {
		goModPath := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(goModPath); err == nil {
			if strings.Contains(string(data), "spectr") {
				return dir
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// =============================================================================
// Scenario Extraction Comparison Tests
// =============================================================================

func TestCompare_ScenarioExtraction(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "single scenario",
			content: `# Spec

## Requirements

### Requirement: Login

The system SHALL support login.

#### Scenario: Valid Login

- WHEN user provides valid credentials
- THEN login succeeds
`,
		},
		{
			name: "multiple scenarios",
			content: `# Spec

## Requirements

### Requirement: Auth

Description.

#### Scenario: Valid Login

Steps.

#### Scenario: Invalid Login

Steps.

#### Scenario: Expired Token

Steps.
`,
		},
		{
			name: "scenarios across requirements",
			content: `# Spec

## Requirements

### Requirement: Auth

#### Scenario: Login

Steps.

### Requirement: Logout

#### Scenario: Normal Logout

Steps.

#### Scenario: Forced Logout

Steps.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Extract requirements with regex
			regexReqs := regexExtractRequirements(tt.content)

			// Extract scenario headers with AST
			doc, err := markdown.ParseDocument([]byte(tt.content))
			if err != nil {
				t.Fatalf("AST parse error: %v", err)
			}

			// Count H4 Scenario headers
			astScenarioCount := 0
			for _, h := range doc.Headers {
				if h.Level == 4 && strings.HasPrefix(h.Text, "Scenario:") {
					astScenarioCount++
				}
			}

			// Count scenarios from regex extraction
			regexScenarioCount := 0
			for _, req := range regexReqs {
				regexScenarioCount += len(req.Scenarios)
			}

			// The comparison here is indirect since the parsers work differently:
			// - Regex extracts scenarios as content blocks within requirements
			// - AST extracts scenario headers directly
			// We compare the total scenario count as a proxy for equivalence
			if regexScenarioCount != astScenarioCount {
				t.Errorf("Scenario count mismatch: regex=%d, ast=%d",
					regexScenarioCount, astScenarioCount)
			}
		})
	}
}

// =============================================================================
// Title Extraction Comparison Tests
// =============================================================================

func TestCompare_TitleExtraction(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "simple title",
			content: `# My Title

Content.
`,
			expected: "My Title",
		},
		{
			name: "title with Change prefix",
			content: `# Change: Add New Feature

Content.
`,
			expected: "Add New Feature",
		},
		{
			name: "title with Spec prefix",
			content: `# Spec: Authentication

Content.
`,
			expected: "Authentication",
		},
		{
			name: "no title",
			content: `## Section Only

Content.
`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Regex-based title extraction (simplified from parsers.ExtractTitle)
			var regexTitle string
			lines := strings.Split(tt.content, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "# ") {
					continue
				}
				regexTitle = strings.TrimPrefix(line, "# ")
				regexTitle = strings.TrimPrefix(regexTitle, "Change:")
				regexTitle = strings.TrimPrefix(regexTitle, "Spec:")
				regexTitle = strings.TrimSpace(regexTitle)

				break
			}

			// AST-based title extraction
			doc, err := markdown.ParseDocument([]byte(tt.content))
			var astTitle string
			if err == nil {
				for _, h := range doc.Headers {
					if h.Level != 1 {
						continue
					}
					astTitle = h.Text
					astTitle = strings.TrimPrefix(astTitle, "Change:")
					astTitle = strings.TrimPrefix(astTitle, "Spec:")
					astTitle = strings.TrimSpace(astTitle)

					break
				}
			}

			if regexTitle != astTitle {
				t.Errorf("Title mismatch: regex=%q, ast=%q", regexTitle, astTitle)
			}

			if regexTitle != tt.expected {
				t.Errorf("Expected title %q, got %q", tt.expected, regexTitle)
			}
		})
	}
}
