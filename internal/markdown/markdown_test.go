// Package markdown provides comprehensive tests for AST-based markdown parsing.
// These tests verify that the blackfriday-based parser produces equivalent
// output to the previous regex-based implementation.
//
//nolint:revive // early-return, unused-parameter - test file patterns
package markdown

import (
	"testing"
)

// testMarkdown is a comprehensive test document covering all parsing scenarios
const testMarkdown = `# Change: Example Change

## Why
This is why we make the change.

## Requirements

### Requirement: First Requirement
This requirement has content.

The requirement MUST do something.

#### Scenario: Happy Path
When something happens.

#### Scenario: Error Case
When something goes wrong.

### Requirement: Second Requirement
Another requirement.

## ADDED Requirements

### Requirement: New Feature
Added requirement content.

## MODIFIED Requirements

### Requirement: Updated Feature
Modified requirement content.

## Tasks
- [ ] 1.1 First task
- [x] 1.2 Completed task
- [ ] 2.1 Another task
`

// TestExtractH1Title tests H1 title extraction
func TestExtractH1Title(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "change prefix",
			content:  "# Change: Example Change\n\nSome content.",
			expected: "Change: Example Change",
		},
		{
			name:     "spec prefix",
			content:  "# Spec: My Spec\n\nContent here.",
			expected: "Spec: My Spec",
		},
		{
			name:     "no prefix",
			content:  "# Simple Title\n\nMore content.",
			expected: "Simple Title",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "no h1",
			content:  "## Only H2\n\nContent.",
			expected: "",
		},
		{
			name:     "multiple h1 - returns first",
			content:  "# First Title\n\n# Second Title\n",
			expected: "First Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Parse([]byte(tt.content))
			result := ExtractH1Title(node)
			if result != tt.expected {
				t.Errorf("ExtractH1Title() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestExtractH1TitleClean tests cleaned H1 title extraction (removes prefixes)
func TestExtractH1TitleClean(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "change prefix removed",
			content:  "# Change: Example Change\n\nSome content.",
			expected: "Example Change",
		},
		{
			name:     "spec prefix removed",
			content:  "# Spec: My Spec\n\nContent here.",
			expected: "My Spec",
		},
		{
			name:     "no prefix unchanged",
			content:  "# Simple Title\n\nMore content.",
			expected: "Simple Title",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "comprehensive test",
			content:  testMarkdown,
			expected: "Example Change",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Parse([]byte(tt.content))
			result := ExtractH1TitleClean(node)
			if result != tt.expected {
				t.Errorf("ExtractH1TitleClean() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestExtractHeaders tests header extraction at all levels
func TestExtractHeaders(t *testing.T) {
	content := `# H1 Title

## H2 Section

### H3 Subsection

#### H4 Detail

## Another H2
`

	node := Parse([]byte(content))
	headers := ExtractHeaders(node)

	expected := []struct {
		level int
		text  string
	}{
		{1, "H1 Title"},
		{2, "H2 Section"},
		{3, "H3 Subsection"},
		{4, "H4 Detail"},
		{2, "Another H2"},
	}

	if len(headers) != len(expected) {
		t.Fatalf("ExtractHeaders() returned %d headers, want %d", len(headers), len(expected))
	}

	for i, exp := range expected {
		if headers[i].Level != exp.level {
			t.Errorf("headers[%d].Level = %d, want %d", i, headers[i].Level, exp.level)
		}
		if headers[i].Text != exp.text {
			t.Errorf("headers[%d].Text = %q, want %q", i, headers[i].Text, exp.text)
		}
	}
}

// TestExtractH2Sections tests H2 section extraction
func TestExtractH2Sections(t *testing.T) {
	node := Parse([]byte(testMarkdown))
	sections := ExtractH2Sections(node)

	// Verify expected sections exist
	expectedSections := []string{
		"Why",
		"Requirements",
		"ADDED Requirements",
		"MODIFIED Requirements",
		"Tasks",
	}
	for _, name := range expectedSections {
		if _, ok := sections[name]; !ok {
			t.Errorf("ExtractH2Sections() missing section %q", name)
		}
	}

	// Verify "Why" section content
	if why, ok := sections["Why"]; ok {
		if why == "" {
			t.Error("ExtractH2Sections() 'Why' section is empty")
		}
	}
}

// TestExtractRequirements tests requirement block extraction
func TestExtractRequirements(t *testing.T) {
	node := Parse([]byte(testMarkdown))
	reqs := ExtractRequirements(node)

	// Should find 4 requirements
	expectedNames := []string{
		"First Requirement",
		"Second Requirement",
		"New Feature",
		"Updated Feature",
	}
	if len(reqs) != len(expectedNames) {
		t.Fatalf(
			"ExtractRequirements() returned %d requirements, want %d",
			len(reqs),
			len(expectedNames),
		)
	}

	for i, name := range expectedNames {
		if reqs[i].Name != name {
			t.Errorf("reqs[%d].Name = %q, want %q", i, reqs[i].Name, name)
		}
	}

	// First requirement should have 2 scenarios
	if len(reqs[0].Scenarios) != 2 {
		t.Errorf("First requirement has %d scenarios, want 2", len(reqs[0].Scenarios))
	}
}

// TestExtractRequirementsFromContent tests requirement extraction from string
func TestExtractRequirementsFromContent(t *testing.T) {
	content := `### Requirement: Single Req
Content here.

#### Scenario: Test Scenario
Scenario content.
`
	reqs := ExtractRequirementsFromContent(content)

	if len(reqs) != 1 {
		t.Fatalf("ExtractRequirementsFromContent() returned %d requirements, want 1", len(reqs))
	}

	if reqs[0].Name != "Single Req" {
		t.Errorf("Name = %q, want %q", reqs[0].Name, "Single Req")
	}

	if len(reqs[0].Scenarios) != 1 {
		t.Errorf("Scenarios count = %d, want 1", len(reqs[0].Scenarios))
	}
}

// TestExtractScenarios tests scenario extraction from requirement content
func TestExtractScenarios(t *testing.T) {
	content := `### Requirement: Test
Some content.

#### Scenario: Happy Path
When everything works.

#### Scenario: Error Case
When things fail.

#### Scenario: Edge Case
When edge conditions occur.
`

	scenarios := ExtractScenarios(content)

	if len(scenarios) != 3 {
		t.Fatalf("ExtractScenarios() returned %d scenarios, want 3", len(scenarios))
	}

	// Each scenario should contain its header
	expectedPrefixes := []string{
		"#### Scenario: Happy Path",
		"#### Scenario: Error Case",
		"#### Scenario: Edge Case",
	}
	for i, prefix := range expectedPrefixes {
		if len(scenarios[i]) < len(prefix) {
			t.Errorf("scenarios[%d] is too short", i)
		}
	}
}

// TestExtractScenarioNames tests scenario name extraction
func TestExtractScenarioNames(t *testing.T) {
	content := `### Requirement: Test
Content here.

#### Scenario: Happy Path
When everything works.

#### Scenario: Error Case
When things fail.
`

	names := ExtractScenarioNames(content)

	expected := []string{"Happy Path", "Error Case"}
	if len(names) != len(expected) {
		t.Fatalf("ExtractScenarioNames() returned %d names, want %d", len(names), len(expected))
	}

	for i, exp := range expected {
		if names[i] != exp {
			t.Errorf("names[%d] = %q, want %q", i, names[i], exp)
		}
	}
}

// TestCountTasks tests task counting
func TestCountTasks(t *testing.T) {
	node := Parse([]byte(testMarkdown))
	total, completed := CountTasks(node)

	if total != 3 {
		t.Errorf("CountTasks() total = %d, want 3", total)
	}
	if completed != 1 {
		t.Errorf("CountTasks() completed = %d, want 1", completed)
	}
}

// TestCountTasksFromContent tests task counting from string content
func TestCountTasksFromContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedTotal int
		expectedDone  int
	}{
		{
			name: "mixed tasks",
			content: `# Tasks
- [ ] First task
- [x] Second task (done)
- [ ] Third task
- [X] Fourth task (also done)
`,
			expectedTotal: 4,
			expectedDone:  2,
		},
		{
			name:          "no tasks",
			content:       "# Just a title\n\nSome content.",
			expectedTotal: 0,
			expectedDone:  0,
		},
		{
			name: "all unchecked",
			content: `- [ ] Task 1
- [ ] Task 2
`,
			expectedTotal: 2,
			expectedDone:  0,
		},
		{
			name: "all checked",
			content: `- [x] Task 1
- [X] Task 2
`,
			expectedTotal: 2,
			expectedDone:  2,
		},
		{
			name:          "empty content",
			content:       "",
			expectedTotal: 0,
			expectedDone:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, completed := CountTasksFromContent(tt.content)
			if total != tt.expectedTotal {
				t.Errorf("total = %d, want %d", total, tt.expectedTotal)
			}
			if completed != tt.expectedDone {
				t.Errorf("completed = %d, want %d", completed, tt.expectedDone)
			}
		})
	}
}

// TestExtractTasks tests task extraction
func TestExtractTasks(t *testing.T) {
	content := `# Tasks
- [ ] First task
- [x] Completed task
- [ ] Another task
- [X] Also completed
`

	tasks := ExtractTasksFromContent(content)

	if len(tasks) != 4 {
		t.Fatalf("ExtractTasksFromContent() returned %d tasks, want 4", len(tasks))
	}

	expected := []struct {
		text    string
		checked bool
	}{
		{"First task", false},
		{"Completed task", true},
		{"Another task", false},
		{"Also completed", true},
	}

	for i, exp := range expected {
		if tasks[i].Text != exp.text {
			t.Errorf("tasks[%d].Text = %q, want %q", i, tasks[i].Text, exp.text)
		}
		if tasks[i].Checked != exp.checked {
			t.Errorf("tasks[%d].Checked = %v, want %v", i, tasks[i].Checked, exp.checked)
		}
	}
}

// TestFindDeltaSection tests delta section finding
func TestFindDeltaSection(t *testing.T) {
	node := Parse([]byte(testMarkdown))

	tests := []struct {
		sectionType string
		shouldFind  bool
	}{
		{"ADDED", true},
		{"MODIFIED", true},
		{"REMOVED", false},
		{"RENAMED", false},
		{"added", true},    // case insensitive
		{"Modified", true}, // case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.sectionType, func(t *testing.T) {
			result := FindDeltaSection(node, tt.sectionType)
			found := result != nil
			if found != tt.shouldFind {
				t.Errorf(
					"FindDeltaSection(%q) found = %v, want %v",
					tt.sectionType,
					found,
					tt.shouldFind,
				)
			}
		})
	}
}

// TestFindAllDeltaSections tests finding all delta sections
func TestFindAllDeltaSections(t *testing.T) {
	node := Parse([]byte(testMarkdown))
	sections := FindAllDeltaSections(node)

	// Should find ADDED and MODIFIED
	if _, ok := sections[DeltaAdded]; !ok {
		t.Error("FindAllDeltaSections() missing ADDED section")
	}
	if _, ok := sections[DeltaModified]; !ok {
		t.Error("FindAllDeltaSections() missing MODIFIED section")
	}
	if _, ok := sections[DeltaRemoved]; ok {
		t.Error("FindAllDeltaSections() found REMOVED section that doesn't exist")
	}
	if _, ok := sections[DeltaRenamed]; ok {
		t.Error("FindAllDeltaSections() found RENAMED section that doesn't exist")
	}
}

// TestExtractDeltaSectionContent tests delta section content extraction
func TestExtractDeltaSectionContent(t *testing.T) {
	node := Parse([]byte(testMarkdown))

	addedContent := ExtractDeltaSectionContent(node, "ADDED")
	if addedContent == "" {
		t.Error("ExtractDeltaSectionContent(ADDED) returned empty string")
	}

	// Should contain the requirement
	if !containsSubstring(addedContent, "New Feature") {
		t.Error("ADDED section should contain 'New Feature' requirement")
	}

	modifiedContent := ExtractDeltaSectionContent(node, "MODIFIED")
	if modifiedContent == "" {
		t.Error("ExtractDeltaSectionContent(MODIFIED) returned empty string")
	}
}

// TestExtractSectionContent tests specific section content extraction
func TestExtractSectionContent(t *testing.T) {
	node := Parse([]byte(testMarkdown))

	whyContent := ExtractSectionContent(node, "Why", 2)
	if whyContent == "" {
		t.Error("ExtractSectionContent(Why) returned empty string")
	}
	if !containsSubstring(whyContent, "why we make the change") {
		t.Errorf(
			"Why section content = %q, expected to contain 'why we make the change'",
			whyContent,
		)
	}
}

// TestExtractRequirementsSection tests requirements section extraction
func TestExtractRequirementsSection(t *testing.T) {
	node := Parse([]byte(testMarkdown))

	reqSection := ExtractRequirementsSection(node)
	if reqSection == "" {
		t.Error("ExtractRequirementsSection() returned empty string")
	}

	// Should contain First Requirement and Second Requirement
	if !containsSubstring(reqSection, "First Requirement") {
		t.Error("Requirements section should contain 'First Requirement'")
	}
	if !containsSubstring(reqSection, "Second Requirement") {
		t.Error("Requirements section should contain 'Second Requirement'")
	}
}

// TestExtractTasksWithIDs tests extraction of tasks with section context and IDs
func TestExtractTasksWithIDs(t *testing.T) {
	content := `# Tasks

## 1. Setup
- [ ] 1.1 First task
- [x] 1.2 Second task

## 2. Implementation
- [ ] 2.1 Third task
- [x] 2.2 Fourth task
`

	tasks := ExtractTasksWithIDs(content)

	if len(tasks) != 4 {
		t.Fatalf("ExtractTasksWithIDs() returned %d tasks, want 4", len(tasks))
	}

	expected := []struct {
		id      string
		section string
		desc    string
		checked bool
	}{
		{"1.1", "Setup", "First task", false},
		{"1.2", "Setup", "Second task", true},
		{"2.1", "Implementation", "Third task", false},
		{"2.2", "Implementation", "Fourth task", true},
	}

	for i, exp := range expected {
		if tasks[i].ID != exp.id {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, exp.id)
		}
		if tasks[i].Section != exp.section {
			t.Errorf("tasks[%d].Section = %q, want %q", i, tasks[i].Section, exp.section)
		}
		if tasks[i].Description != exp.desc {
			t.Errorf("tasks[%d].Description = %q, want %q", i, tasks[i].Description, exp.desc)
		}
		if tasks[i].Checked != exp.checked {
			t.Errorf("tasks[%d].Checked = %v, want %v", i, tasks[i].Checked, exp.checked)
		}
	}
}

// TestExtractOrderedRequirementNames tests ordered requirement name extraction
func TestExtractOrderedRequirementNames(t *testing.T) {
	content := `## Requirements

### Requirement: Alpha
Content.

### Requirement: Beta
Content.

### Requirement: Gamma
Content.
`

	names := ExtractOrderedRequirementNames(content)

	expected := []string{"Alpha", "Beta", "Gamma"}
	if len(names) != len(expected) {
		t.Fatalf(
			"ExtractOrderedRequirementNames() returned %d names, want %d",
			len(names),
			len(expected),
		)
	}

	for i, exp := range expected {
		if names[i] != exp {
			t.Errorf("names[%d] = %q, want %q", i, names[i], exp)
		}
	}
}

// TestSplitSpec tests spec splitting into preamble, requirements, and after
func TestSplitSpec(t *testing.T) {
	content := `# Spec: Test Spec

## Purpose
This is the purpose.

## Requirements

### Requirement: First
Content here.

## Notes
Some notes at the end.
`

	preamble, requirements, after := SplitSpec(content)

	// Preamble should contain title and purpose, plus the ## Requirements header
	if !containsSubstring(preamble, "Test Spec") {
		t.Error("Preamble should contain title")
	}
	if !containsSubstring(preamble, "Purpose") {
		t.Error("Preamble should contain Purpose section")
	}

	// Requirements should contain the requirement
	if !containsSubstring(requirements, "First") {
		t.Error("Requirements should contain 'First' requirement")
	}

	// After should contain notes
	if !containsSubstring(after, "Notes") {
		t.Error("After should contain Notes section")
	}
}

// TestEmptyContent tests behavior with empty content
func TestEmptyContent(t *testing.T) {
	node := Parse([]byte(""))

	// All extraction functions should handle empty content gracefully
	if title := ExtractH1Title(node); title != "" {
		t.Errorf("ExtractH1Title(empty) = %q, want empty", title)
	}

	headers := ExtractHeaders(node)
	if len(headers) != 0 {
		t.Errorf("ExtractHeaders(empty) returned %d headers, want 0", len(headers))
	}

	sections := ExtractH2Sections(node)
	if len(sections) != 0 {
		t.Errorf("ExtractH2Sections(empty) returned %d sections, want 0", len(sections))
	}

	reqs := ExtractRequirements(node)
	if len(reqs) != 0 {
		t.Errorf("ExtractRequirements(empty) returned %d requirements, want 0", len(reqs))
	}

	tasks := ExtractTasks(node)
	if len(tasks) != 0 {
		t.Errorf("ExtractTasks(empty) returned %d tasks, want 0", len(tasks))
	}

	total, completed := CountTasks(node)
	if total != 0 || completed != 0 {
		t.Errorf("CountTasks(empty) = (%d, %d), want (0, 0)", total, completed)
	}
}

// TestNilNode tests behavior with nil node
func TestNilNode(t *testing.T) {
	// All extraction functions should handle nil node gracefully
	if title := ExtractH1Title(nil); title != "" {
		t.Errorf("ExtractH1Title(nil) = %q, want empty", title)
	}

	if title := ExtractH1TitleClean(nil); title != "" {
		t.Errorf("ExtractH1TitleClean(nil) = %q, want empty", title)
	}

	headers := ExtractHeaders(nil)
	if len(headers) != 0 {
		t.Errorf("ExtractHeaders(nil) returned %d headers, want 0", len(headers))
	}

	sections := ExtractH2Sections(nil)
	if len(sections) != 0 {
		t.Errorf("ExtractH2Sections(nil) returned %d sections, want 0", len(sections))
	}

	reqs := ExtractRequirements(nil)
	if len(reqs) != 0 {
		t.Errorf("ExtractRequirements(nil) returned %d requirements, want 0", len(reqs))
	}

	tasks := ExtractTasks(nil)
	if len(tasks) != 0 {
		t.Errorf("ExtractTasks(nil) returned %d tasks, want 0", len(tasks))
	}

	total, completed := CountTasks(nil)
	if total != 0 || completed != 0 {
		t.Errorf("CountTasks(nil) = (%d, %d), want (0, 0)", total, completed)
	}

	if node := FindDeltaSection(nil, "ADDED"); node != nil {
		t.Error("FindDeltaSection(nil) should return nil")
	}

	deltaSections := FindAllDeltaSections(nil)
	if len(deltaSections) != 0 {
		t.Errorf("FindAllDeltaSections(nil) returned %d sections, want 0", len(deltaSections))
	}
}

// TestMalformedMarkdown tests behavior with malformed markdown
func TestMalformedMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "unclosed code block",
			content: "# Title\n\n```go\nfunc main() {\n",
		},
		{
			name:    "only hashes",
			content: "#\n##\n###\n",
		},
		{
			name:    "mixed indentation",
			content: "# Title\n  ## Indented H2\n    ### More indented\n",
		},
		{
			name:    "task without checkbox",
			content: "- Not a task\n- Also not a task\n",
		},
		{
			name:    "requirement without name",
			content: "### Requirement:\n\nContent.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These should not panic
			node := Parse([]byte(tt.content))
			_ = ExtractH1Title(node)
			_ = ExtractHeaders(node)
			_ = ExtractH2Sections(node)
			_ = ExtractRequirements(node)
			_ = ExtractTasks(node)
			_, _ = CountTasks(node)
		})
	}
}

// TestNestedContent tests extraction of nested content
func TestNestedContent(t *testing.T) {
	content := `# Title

## Section

### Requirement: Test Req
The requirement text.

- List item 1
- List item 2

#### Scenario: Nested Scenario
Scenario content.

> Quote in scenario

Some more text.
`

	node := Parse([]byte(content))

	reqs := ExtractRequirements(node)
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement, got %d", len(reqs))
	}

	// Requirement should have the scenario
	if len(reqs[0].Scenarios) != 1 {
		t.Errorf("Expected 1 scenario, got %d", len(reqs[0].Scenarios))
	}
}

// TestRequirementBlockRaw tests that RequirementBlock.Raw contains full content
func TestRequirementBlockRaw(t *testing.T) {
	content := `### Requirement: Test
Content line 1.
Content line 2.

#### Scenario: Test Scenario
Scenario content.
`

	reqs := ExtractRequirementsFromContent(content)
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement, got %d", len(reqs))
	}

	// Raw should start with the header
	if !containsSubstring(reqs[0].Raw, "### Requirement: Test") {
		t.Error("Raw should contain requirement header")
	}

	// Raw should contain the content
	if !containsSubstring(reqs[0].Raw, "Content line 1") {
		t.Error("Raw should contain content")
	}

	// Raw should contain the scenario
	if !containsSubstring(reqs[0].Raw, "Scenario: Test Scenario") {
		t.Error("Raw should contain scenario")
	}
}

// TestValidDeltaTypes tests the ValidDeltaTypes helper function
func TestValidDeltaTypes(t *testing.T) {
	types := ValidDeltaTypes()

	expected := []string{"ADDED", "MODIFIED", "REMOVED", "RENAMED"}
	if len(types) != len(expected) {
		t.Fatalf("ValidDeltaTypes() returned %d types, want %d", len(types), len(expected))
	}

	for i, exp := range expected {
		if types[i] != exp {
			t.Errorf("types[%d] = %q, want %q", i, types[i], exp)
		}
	}
}

// TestDeltaTypeConstants tests delta type constants
func TestDeltaTypeConstants(t *testing.T) {
	if string(DeltaAdded) != "ADDED" {
		t.Errorf("DeltaAdded = %q, want %q", DeltaAdded, "ADDED")
	}
	if string(DeltaModified) != "MODIFIED" {
		t.Errorf("DeltaModified = %q, want %q", DeltaModified, "MODIFIED")
	}
	if string(DeltaRemoved) != "REMOVED" {
		t.Errorf("DeltaRemoved = %q, want %q", DeltaRemoved, "REMOVED")
	}
	if string(DeltaRenamed) != "RENAMED" {
		t.Errorf("DeltaRenamed = %q, want %q", DeltaRenamed, "RENAMED")
	}
}

// TestCodeBlockInRequirement tests that code blocks are handled in requirements
func TestCodeBlockInRequirement(t *testing.T) {
	content := "### Requirement: Code Example\nDescription.\n\n```go\nfunc main() {\n    fmt.Println(\"hello\")\n}\n```\n\nMore text.\n"

	reqs := ExtractRequirementsFromContent(content)
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement, got %d", len(reqs))
	}

	// The Raw content should include the code block
	if !containsSubstring(reqs[0].Raw, "func main()") {
		t.Error("Requirement Raw should contain code block content")
	}
}

// TestInlineCodeInHeader tests headers with inline code
func TestInlineCodeInHeader(t *testing.T) {
	content := "# Title with `code`\n\n## Section with `more code`\n"

	node := Parse([]byte(content))

	title := ExtractH1Title(node)
	if !containsSubstring(title, "code") {
		t.Errorf("Title should contain 'code', got %q", title)
	}

	sections := ExtractH2Sections(node)
	found := false
	for name := range sections {
		if containsSubstring(name, "code") {
			found = true

			break
		}
	}
	if !found {
		t.Error("Should find section with 'code' in name")
	}
}

// containsSubstring is a helper for checking if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
