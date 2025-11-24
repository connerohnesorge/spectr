package parsers

import (
	"testing"

	"github.com/connerohnesorge/spectr/internal/mdparser"
)

func TestExtractRequirements_SingleRequirement(t *testing.T) {
	content := `# Test Spec

## Requirements

### Requirement: Feature One
The system SHALL do something.

#### Scenario: Success case
- **WHEN** action occurs
- **THEN** result happens
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement, got %d", len(reqs))
	}

	req := reqs[0]
	if req.Name != "Feature One" {
		t.Errorf("Expected name 'Feature One', got %q", req.Name)
	}

	if req.HeaderLine != "### Requirement: Feature One" {
		t.Errorf("Expected header line '### Requirement: Feature One', got %q", req.HeaderLine)
	}

	if req.Raw == "" {
		t.Error("Expected non-empty raw content")
	}
}

func TestExtractRequirements_MultipleRequirements(t *testing.T) {
	content := `# Test Spec

## Requirements

### Requirement: Feature One
Description one.

#### Scenario: Case 1
- **WHEN** something
- **THEN** result

### Requirement: Feature Two
Description two.

#### Scenario: Case 2
- **WHEN** other thing
- **THEN** other result
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 2 {
		t.Fatalf("Expected 2 requirements, got %d", len(reqs))
	}

	if reqs[0].Name != "Feature One" {
		t.Errorf("Expected first requirement 'Feature One', got %q", reqs[0].Name)
	}

	if reqs[1].Name != "Feature Two" {
		t.Errorf("Expected second requirement 'Feature Two', got %q", reqs[1].Name)
	}
}

func TestExtractRequirements_NoScenarios_Error(t *testing.T) {
	content := `# Test Spec

## Requirements

### Requirement: Missing Scenarios
This requirement has no scenarios.
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = ExtractRequirements(doc)
	if err == nil {
		t.Fatal("Expected error for requirement without scenarios, got nil")
	}

	if !containsString(err.Error(), "no scenarios") {
		t.Errorf("Expected error about missing scenarios, got: %v", err)
	}
}

func TestExtractRequirements_PreservesListMarkers(t *testing.T) {
	content := `# Test Spec

## Requirements

### Requirement: List Formatting
1. First item
4. Fourth item

* Star bullet
- Dash bullet

#### Scenario: Example
- **WHEN** something happens
- **THEN** result
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement, got %d", len(reqs))
	}

	expected := `### Requirement: List Formatting
1. First item
4. Fourth item

* Star bullet
- Dash bullet

#### Scenario: Example
- **WHEN** something happens
- **THEN** result
`

	if reqs[0].Raw != expected {
		t.Fatalf("Raw requirement did not preserve list markers.\nExpected:\n%s\nGot:\n%s", expected, reqs[0].Raw)
	}
}

func TestExtractRequirements_CodeBlockIgnored(t *testing.T) {
	content := `# Test Spec

## Requirements

### Requirement: Real Feature
The system SHALL work.

#### Scenario: Example
- **WHEN** something happens
- **THEN** result occurs

` + "```markdown" + `
### Requirement: Fake Feature
This is in a code block and should be ignored.

#### Scenario: Fake scenario
- **WHEN** fake
- **THEN** fake
` + "```" + `
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	// Should only find the real requirement, not the one in code block
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement (code block should be ignored), got %d", len(reqs))
	}

	if reqs[0].Name != "Real Feature" {
		t.Errorf("Expected 'Real Feature', got %q", reqs[0].Name)
	}
}

func TestExtractScenarios_SingleScenario(t *testing.T) {
	content := `# Test Spec

### Requirement: Test
Description.

#### Scenario: Success case
- **WHEN** action
- **THEN** result
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Find the requirement header
	var reqHeader *mdparser.Header
	var startIdx int
	//nolint:revive // early-return would make this search logic less clear
	for i, node := range doc.Children {
		if header, ok := node.(*mdparser.Header); ok && header.Level == 3 {
			reqHeader = header
			startIdx = i + 1

			break
		}
	}

	if reqHeader == nil {
		t.Fatal("Could not find requirement header")
	}

	siblings := getSiblingsUntilNextHeader(doc.Children, startIdx, 3)
	scenarios, err := ExtractScenarios(reqHeader, siblings)
	if err != nil {
		t.Fatalf("ExtractScenarios failed: %v", err)
	}

	if len(scenarios) != 1 {
		t.Fatalf("Expected 1 scenario, got %d", len(scenarios))
	}

	if scenarios[0].Name != "Success case" {
		t.Errorf("Expected scenario name 'Success case', got %q", scenarios[0].Name)
	}

	if len(scenarios[0].Steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(scenarios[0].Steps))
	}
}

func TestExtractScenarios_MultipleScenarios(t *testing.T) {
	content := `# Test Spec

### Requirement: Test
Description.

#### Scenario: Happy path
- **WHEN** normal
- **THEN** success

#### Scenario: Error path
- **WHEN** error
- **THEN** handle
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var reqHeader *mdparser.Header
	var startIdx int
	//nolint:revive // early-return would make this search logic less clear
	for i, node := range doc.Children {
		if header, ok := node.(*mdparser.Header); ok && header.Level == 3 {
			reqHeader = header
			startIdx = i + 1

			break
		}
	}

	if reqHeader == nil {
		t.Fatal("Could not find requirement header")
	}

	siblings := getSiblingsUntilNextHeader(doc.Children, startIdx, 3)
	scenarios, err := ExtractScenarios(reqHeader, siblings)
	if err != nil {
		t.Fatalf("ExtractScenarios failed: %v", err)
	}

	if len(scenarios) != 2 {
		t.Fatalf("Expected 2 scenarios, got %d", len(scenarios))
	}

	if scenarios[0].Name != "Happy path" {
		t.Errorf("Expected first scenario 'Happy path', got %q", scenarios[0].Name)
	}

	if scenarios[1].Name != "Error path" {
		t.Errorf("Expected second scenario 'Error path', got %q", scenarios[1].Name)
	}
}

func TestExtractDeltaSections_AllOperations(t *testing.T) {
	content := `# Delta Spec

## ADDED Requirements

### Requirement: New Feature
Content here.

#### Scenario: Test
- **WHEN** something
- **THEN** result

## MODIFIED Requirements

### Requirement: Updated Feature
Modified content.

#### Scenario: Modified test
- **WHEN** action
- **THEN** new result

## REMOVED Requirements

### Requirement: Old Feature
**Reason**: Deprecated
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	sections, err := ExtractDeltaSections(doc)
	if err != nil {
		t.Fatalf("ExtractDeltaSections failed: %v", err)
	}

	if len(sections["ADDED"]) != 1 {
		t.Errorf("Expected 1 ADDED requirement, got %d", len(sections["ADDED"]))
	}

	if len(sections["MODIFIED"]) != 1 {
		t.Errorf("Expected 1 MODIFIED requirement, got %d", len(sections["MODIFIED"]))
	}

	if len(sections["REMOVED"]) != 1 {
		t.Errorf("Expected 1 REMOVED requirement, got %d", len(sections["REMOVED"]))
	}

	if sections["ADDED"][0].Name != "New Feature" {
		t.Errorf("Expected ADDED 'New Feature', got %q", sections["ADDED"][0].Name)
	}

	if sections["MODIFIED"][0].Name != "Updated Feature" {
		t.Errorf("Expected MODIFIED 'Updated Feature', got %q", sections["MODIFIED"][0].Name)
	}

	if sections["REMOVED"][0].Name != "Old Feature" {
		t.Errorf("Expected REMOVED 'Old Feature', got %q", sections["REMOVED"][0].Name)
	}
}

func TestExtractDeltaSections_EmptySections(t *testing.T) {
	content := `# Delta Spec

## ADDED Requirements

## MODIFIED Requirements

## REMOVED Requirements
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	sections, err := ExtractDeltaSections(doc)
	if err != nil {
		t.Fatalf("ExtractDeltaSections failed: %v", err)
	}

	if len(sections["ADDED"]) != 0 {
		t.Errorf("Expected 0 ADDED requirements, got %d", len(sections["ADDED"]))
	}

	if len(sections["MODIFIED"]) != 0 {
		t.Errorf("Expected 0 MODIFIED requirements, got %d", len(sections["MODIFIED"]))
	}

	if len(sections["REMOVED"]) != 0 {
		t.Errorf("Expected 0 REMOVED requirements, got %d", len(sections["REMOVED"]))
	}
}

func TestExtractDeltaSections_NoSections(t *testing.T) {
	content := `# Regular Spec

## Requirements

### Requirement: Regular Requirement
This is not a delta.

#### Scenario: Test
- **WHEN** something
- **THEN** result
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	sections, err := ExtractDeltaSections(doc)
	if err != nil {
		t.Fatalf("ExtractDeltaSections failed: %v", err)
	}

	if len(sections) != 0 {
		t.Errorf("Expected empty sections map, got %d entries", len(sections))
	}
}

func TestExtractRenamedRequirements_SingleRename(t *testing.T) {
	content := `# Delta Spec

## RENAMED Requirements

- FROM: ` + "`" + `### Requirement: Old Name` + "`" + `
- TO: ` + "`" + `### Requirement: New Name` + "`" + `
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	renamed, err := ExtractRenamedRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRenamedRequirements failed: %v", err)
	}

	if len(renamed) != 1 {
		t.Fatalf("Expected 1 renamed requirement, got %d", len(renamed))
	}

	if renamed[0].From != "Old Name" {
		t.Errorf("Expected from 'Old Name', got %q", renamed[0].From)
	}

	if renamed[0].To != "New Name" {
		t.Errorf("Expected to 'New Name', got %q", renamed[0].To)
	}
}

func TestExtractRenamedRequirements_MultipleRenames(t *testing.T) {
	content := `# Delta Spec

## RENAMED Requirements

- FROM: ` + "`" + `### Requirement: Old Name` + "`" + `
- TO: ` + "`" + `### Requirement: New Name` + "`" + `

- FROM: ` + "`" + `### Requirement: Another Old Name` + "`" + `
- TO: ` + "`" + `### Requirement: Another New Name` + "`" + `
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	renamed, err := ExtractRenamedRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRenamedRequirements failed: %v", err)
	}

	if len(renamed) != 2 {
		t.Fatalf("Expected 2 renamed requirements, got %d", len(renamed))
	}

	if renamed[0].From != "Old Name" {
		t.Errorf("Expected first from 'Old Name', got %q", renamed[0].From)
	}

	if renamed[0].To != "New Name" {
		t.Errorf("Expected first to 'New Name', got %q", renamed[0].To)
	}

	if renamed[1].From != "Another Old Name" {
		t.Errorf("Expected second from 'Another Old Name', got %q", renamed[1].From)
	}

	if renamed[1].To != "Another New Name" {
		t.Errorf("Expected second to 'Another New Name', got %q", renamed[1].To)
	}
}

func TestExtractRenamedRequirements_NoRenamedSection(t *testing.T) {
	content := `# Delta Spec

## ADDED Requirements

### Requirement: New Feature
Content.

#### Scenario: Test
- **WHEN** something
- **THEN** result
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	renamed, err := ExtractRenamedRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRenamedRequirements failed: %v", err)
	}

	if len(renamed) != 0 {
		t.Errorf("Expected 0 renamed requirements, got %d", len(renamed))
	}
}

func TestExtractRenamedRequirements_EmptyRenamedSection(t *testing.T) {
	content := `# Delta Spec

## RENAMED Requirements

No actual renames here.
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	renamed, err := ExtractRenamedRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRenamedRequirements failed: %v", err)
	}

	if len(renamed) != 0 {
		t.Errorf("Expected 0 renamed requirements, got %d", len(renamed))
	}
}

func TestExtractRequirements_MalformedHeaders(t *testing.T) {
	content := `# Test Spec

## Requirements

### Feature One
This is not a requirement (missing "Requirement:" prefix).

#### Scenario: Test
- **WHEN** something
- **THEN** result

### Requirement: Valid Feature
This is a valid requirement.

#### Scenario: Valid test
- **WHEN** action
- **THEN** result
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	// Should only find the valid requirement
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement (malformed header ignored), got %d", len(reqs))
	}

	if reqs[0].Name != "Valid Feature" {
		t.Errorf("Expected 'Valid Feature', got %q", reqs[0].Name)
	}
}

func TestExtractScenarios_MalformedScenarioHeaders(t *testing.T) {
	content := `# Test Spec

### Requirement: Test
Description.

#### Success case
This is not a scenario (missing "Scenario:" prefix).
- **WHEN** something
- **THEN** result

#### Scenario: Valid scenario
- **WHEN** action
- **THEN** valid result
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var reqHeader *mdparser.Header
	var startIdx int
	//nolint:revive // early-return would make this search logic less clear
	for i, node := range doc.Children {
		if header, ok := node.(*mdparser.Header); ok && header.Level == 3 {
			reqHeader = header
			startIdx = i + 1

			break
		}
	}

	if reqHeader == nil {
		t.Fatal("Could not find requirement header")
	}

	siblings := getSiblingsUntilNextHeader(doc.Children, startIdx, 3)
	scenarios, err := ExtractScenarios(reqHeader, siblings)
	if err != nil {
		t.Fatalf("ExtractScenarios failed: %v", err)
	}

	// Should only find the valid scenario
	if len(scenarios) != 1 {
		t.Fatalf("Expected 1 scenario (malformed header ignored), got %d", len(scenarios))
	}

	if scenarios[0].Name != "Valid scenario" {
		t.Errorf("Expected 'Valid scenario', got %q", scenarios[0].Name)
	}
}

func TestExtractRequirements_MultipleScenarios(t *testing.T) {
	content := `# Test Spec

## Requirements

### Requirement: Complex Feature
The system SHALL handle complexity.

#### Scenario: Happy path
- **WHEN** normal operation
- **THEN** success

#### Scenario: Error path
- **WHEN** error occurs
- **THEN** handle gracefully

#### Scenario: Edge case
- **WHEN** boundary condition
- **THEN** handle correctly
`

	doc, err := mdparser.Parse(content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 1 {
		t.Fatalf("Expected 1 requirement, got %d", len(reqs))
	}

	// The requirement should be valid (has scenarios)
	if reqs[0].Name != "Complex Feature" {
		t.Errorf("Expected 'Complex Feature', got %q", reqs[0].Name)
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
