package parser

import (
	"strings"
	"testing"
)

// TestExtractRequirements_NoScenarios tests extracting requirements without scenarios.
func TestExtractRequirements_NoScenarios(t *testing.T) {
	input := `## Requirements

### Requirement: User Authentication
The system SHALL authenticate users with credentials.

### Requirement: Session Management
The system SHALL manage user sessions.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(reqs))
	}

	// Check first requirement
	if reqs[0].Name != "User Authentication" {
		t.Errorf("expected name 'User Authentication', got '%s'", reqs[0].Name)
	}
	if !strings.Contains(reqs[0].Content, "SHALL authenticate") {
		t.Errorf("expected content to contain 'SHALL authenticate', got '%s'", reqs[0].Content)
	}
	if len(reqs[0].Scenarios) != 0 {
		t.Errorf("expected 0 scenarios, got %d", len(reqs[0].Scenarios))
	}

	// Check second requirement
	if reqs[1].Name != "Session Management" {
		t.Errorf("expected name 'Session Management', got '%s'", reqs[1].Name)
	}
	if !strings.Contains(reqs[1].Content, "SHALL manage") {
		t.Errorf("expected content to contain 'SHALL manage', got '%s'", reqs[1].Content)
	}
	if len(reqs[1].Scenarios) != 0 {
		t.Errorf("expected 0 scenarios, got %d", len(reqs[1].Scenarios))
	}
}

// TestExtractRequirements_WithScenarios tests extracting requirements with scenarios.
func TestExtractRequirements_WithScenarios(t *testing.T) {
	input := `## Requirements

### Requirement: User Authentication
The system SHALL authenticate users with credentials.

#### Scenario: Valid credentials
- **WHEN** valid credentials provided
- **THEN** user is authenticated

#### Scenario: Invalid credentials
- **WHEN** invalid credentials provided
- **THEN** error is returned

### Requirement: Session Management
The system SHALL manage user sessions.

#### Scenario: Session creation
- **WHEN** user logs in
- **THEN** session is created
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(reqs))
	}

	// Check first requirement has 2 scenarios
	if len(reqs[0].Scenarios) != 2 {
		t.Fatalf("expected 2 scenarios in first requirement, got %d", len(reqs[0].Scenarios))
	}

	if reqs[0].Scenarios[0].Name != "Valid credentials" {
		t.Errorf("expected scenario name 'Valid credentials', got '%s'", reqs[0].Scenarios[0].Name)
	}
	if !strings.Contains(reqs[0].Scenarios[0].Content, "**WHEN**") {
		t.Errorf(
			"expected scenario content to contain '**WHEN**', got '%s'",
			reqs[0].Scenarios[0].Content,
		)
	}

	if reqs[0].Scenarios[1].Name != "Invalid credentials" {
		t.Errorf(
			"expected scenario name 'Invalid credentials', got '%s'",
			reqs[0].Scenarios[1].Name,
		)
	}

	// Check second requirement has 1 scenario
	if len(reqs[1].Scenarios) != 1 {
		t.Fatalf("expected 1 scenario in second requirement, got %d", len(reqs[1].Scenarios))
	}

	if reqs[1].Scenarios[0].Name != "Session creation" {
		t.Errorf("expected scenario name 'Session creation', got '%s'", reqs[1].Scenarios[0].Name)
	}
}

// TestExtractSections tests extracting level-2 sections.
func TestExtractSections(t *testing.T) {
	input := `# Spec

## ADDED Requirements

### Requirement: New Feature
The system SHALL provide new feature.

## MODIFIED Requirements

### Requirement: Updated Feature
The system SHALL update feature.

## REMOVED Requirements

### Requirement: Old Feature
**Reason**: No longer needed
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	sections, err := ExtractSections(doc)
	if err != nil {
		t.Fatalf("ExtractSections failed: %v", err)
	}

	if len(sections) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(sections))
	}

	// Check section names
	expectedNames := []string{"ADDED Requirements", "MODIFIED Requirements", "REMOVED Requirements"}
	for i, expected := range expectedNames {
		if sections[i].Name != expected {
			t.Errorf("expected section name '%s', got '%s'", expected, sections[i].Name)
		}
	}

	// Check that content contains requirements
	if !strings.Contains(sections[0].Content, "New Feature") {
		t.Errorf("expected ADDED section to contain 'New Feature', got '%s'", sections[0].Content)
	}
	if !strings.Contains(sections[1].Content, "Updated Feature") {
		t.Errorf(
			"expected MODIFIED section to contain 'Updated Feature', got '%s'",
			sections[1].Content,
		)
	}
	if !strings.Contains(sections[2].Content, "Old Feature") {
		t.Errorf("expected REMOVED section to contain 'Old Feature', got '%s'", sections[2].Content)
	}
}

// TestExtractRequirements_HierarchyValidation tests scenario hierarchy validation.
func TestExtractRequirements_HierarchyValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
		errMsg    string
	}{
		{
			name: "scenario without requirement",
			input: `## Requirements

#### Scenario: Orphan scenario
- **WHEN** something happens
- **THEN** error
`,
			expectErr: true,
			errMsg:    "but no requirements defined",
		},
		{
			name: "scenario before requirement",
			input: `## Requirements

#### Scenario: Early scenario
- **WHEN** something happens

### Requirement: Late requirement
The system SHALL do something.
`,
			expectErr: true,
			errMsg:    "is not within a requirement",
		},
		{
			name: "valid hierarchy",
			input: `## Requirements

### Requirement: Valid
The system SHALL work.

#### Scenario: Valid scenario
- **WHEN** it works
- **THEN** success
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

			_, err = ExtractRequirements(doc)

			if tt.expectErr && err == nil {
				t.Errorf("expected error containing '%s', got nil", tt.errMsg)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if tt.expectErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

// TestExtractRequirements_ContentExtraction tests correct text extraction.
func TestExtractRequirements_ContentExtraction(t *testing.T) {
	input := `## Requirements

### Requirement: Complete Feature
The system SHALL provide complete functionality.

This includes multiple paragraphs
of descriptive text.

#### Scenario: First scenario
- **WHEN** user performs action
- **THEN** system responds

Some more text after the scenario.

#### Scenario: Second scenario
- **WHEN** another action
- **THEN** another response
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(reqs))
	}

	req := reqs[0]

	// Check content includes multiple paragraphs
	if !strings.Contains(req.Content, "complete functionality") {
		t.Errorf("expected content to include 'complete functionality', got '%s'", req.Content)
	}
	if !strings.Contains(req.Content, "multiple paragraphs") {
		t.Errorf("expected content to include 'multiple paragraphs', got '%s'", req.Content)
	}

	// Check scenarios
	if len(req.Scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(req.Scenarios))
	}

	// Check first scenario content
	if !strings.Contains(req.Scenarios[0].Content, "**WHEN**") {
		t.Errorf(
			"expected scenario content to contain '**WHEN**', got '%s'",
			req.Scenarios[0].Content,
		)
	}
	if !strings.Contains(req.Scenarios[0].Content, "system responds") {
		t.Errorf(
			"expected scenario content to contain 'system responds', got '%s'",
			req.Scenarios[0].Content,
		)
	}
	if !strings.Contains(req.Scenarios[0].Content, "Some more text") {
		t.Errorf(
			"expected scenario content to include text after list, got '%s'",
			req.Scenarios[0].Content,
		)
	}
}

// TestExtractRequirements_SkipCodeBlocks tests that code blocks are skipped.
func TestExtractRequirements_SkipCodeBlocks(t *testing.T) {
	input := `## Requirements

### Requirement: Real Requirement
The system SHALL do real work.

Here is some example code:

` + "```markdown" + `
### Requirement: Fake Requirement
This is in a code block and should be ignored.

#### Scenario: Fake scenario
Should not be extracted.
` + "```" + `

#### Scenario: Real scenario
- **WHEN** real action
- **THEN** real result
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	// Should only find the real requirement, not the one in the code block
	if len(reqs) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(reqs))
	}

	if reqs[0].Name != "Real Requirement" {
		t.Errorf("expected 'Real Requirement', got '%s'", reqs[0].Name)
	}

	// Should only find the real scenario
	if len(reqs[0].Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(reqs[0].Scenarios))
	}

	if reqs[0].Scenarios[0].Name != "Real scenario" {
		t.Errorf("expected 'Real scenario', got '%s'", reqs[0].Scenarios[0].Name)
	}

	// Code block content should not appear in requirement content
	if strings.Contains(reqs[0].Content, "Fake Requirement") {
		t.Errorf(
			"requirement content should not include code block text, got '%s'",
			reqs[0].Content,
		)
	}
}

// TestExtractRequirements_CriticalCase tests requirements in code blocks are NOT extracted.
func TestExtractRequirements_CriticalCase(t *testing.T) {
	input := `## Example Documentation

Here's how to write a requirement:

` + "```markdown" + `
### Requirement: Example Requirement
This is just an example in documentation.

#### Scenario: Example scenario
This should not be extracted as a real requirement.
` + "```" + `

## Real Requirements

### Requirement: Actual Requirement
The system SHALL do actual work.

#### Scenario: Actual scenario
- **WHEN** actual work is done
- **THEN** actual result
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	// Should only extract the actual requirement, not the example in code block
	if len(reqs) != 1 {
		t.Fatalf("expected 1 requirement (code block should be ignored), got %d", len(reqs))
	}

	if reqs[0].Name != "Actual Requirement" {
		t.Errorf("expected 'Actual Requirement', got '%s'", reqs[0].Name)
	}

	// Check no example content leaked into requirements
	if strings.Contains(reqs[0].Content, "Example Requirement") {
		t.Error("should not extract requirements from code blocks")
	}

	// Check scenarios
	if len(reqs[0].Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(reqs[0].Scenarios))
	}

	if reqs[0].Scenarios[0].Name != "Actual scenario" {
		t.Errorf("expected 'Actual scenario', got '%s'", reqs[0].Scenarios[0].Name)
	}
}

// TestExtractRequirements_EmptyDocument tests empty document handling.
func TestExtractRequirements_EmptyDocument(t *testing.T) {
	input := `# Empty Spec

No requirements here.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 0 {
		t.Errorf("expected 0 requirements in empty document, got %d", len(reqs))
	}
}

// TestExtractRequirements_MultilineContent tests extracting multi-line content.
func TestExtractRequirements_MultilineContent(t *testing.T) {
	input := `## Requirements

### Requirement: Complex Feature
The system SHALL provide:
- Feature A
- Feature B
- Feature C

Additional context paragraph.

#### Scenario: Multi-step scenario
- **GIVEN** initial state
- **WHEN** action performed
- **THEN** expected result
- **AND** additional verification
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(reqs))
	}

	// Check that list items are captured
	content := reqs[0].Content
	if !strings.Contains(content, "Feature A") || !strings.Contains(content, "Feature B") {
		t.Errorf("expected content to include list items, got '%s'", content)
	}

	// Check scenario has all steps
	if len(reqs[0].Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(reqs[0].Scenarios))
	}

	scenarioContent := reqs[0].Scenarios[0].Content
	expectedSteps := []string{"**GIVEN**", "**WHEN**", "**THEN**", "**AND**"}
	for _, step := range expectedSteps {
		if !strings.Contains(scenarioContent, step) {
			t.Errorf("expected scenario to contain '%s', got '%s'", step, scenarioContent)
		}
	}
}

// TestNormalizeRequirementName tests requirement name normalization.
func TestNormalizeRequirementName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User Authentication", "user authentication"},
		{"  Whitespace  ", "whitespace"},
		{"UPPERCASE", "uppercase"},
		{"Mixed Case Name", "mixed case name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeRequirementName(tt.input)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestExtractSections_EmptyDocument tests section extraction from empty document.
func TestExtractSections_EmptyDocument(t *testing.T) {
	input := `# Title

Just some content, no sections.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	sections, err := ExtractSections(doc)
	if err != nil {
		t.Fatalf("ExtractSections failed: %v", err)
	}

	if len(sections) != 0 {
		t.Errorf("expected 0 sections (no ## headers), got %d", len(sections))
	}
}

// TestExtractRequirements_PositionTracking tests that positions are captured.
func TestExtractRequirements_PositionTracking(t *testing.T) {
	input := `## Requirements

### Requirement: First
Content here.

### Requirement: Second
More content.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs, err := ExtractRequirements(doc)
	if err != nil {
		t.Fatalf("ExtractRequirements failed: %v", err)
	}

	if len(reqs) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(reqs))
	}

	// Check that positions are tracked
	if reqs[0].Position.Line == 0 {
		t.Error("expected position to be set for first requirement")
	}
	if reqs[1].Position.Line == 0 {
		t.Error("expected position to be set for second requirement")
	}

	// Second requirement should have a higher line number
	if reqs[1].Position.Line <= reqs[0].Position.Line {
		t.Errorf("expected second requirement to have higher line number, got %d <= %d",
			reqs[1].Position.Line, reqs[0].Position.Line)
	}
}

// TestExtractDeltas_ADDED tests extracting ADDED requirements.
func TestExtractDeltas_ADDED(t *testing.T) {
	input := `# Delta Spec

## ADDED Requirements

### Requirement: New Feature
The system SHALL provide new feature.

#### Scenario: Feature works
- **WHEN** feature is used
- **THEN** it works
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Added) != 1 {
		t.Fatalf("expected 1 ADDED requirement, got %d", len(delta.Added))
	}

	if delta.Added[0].Name != "New Feature" {
		t.Errorf("expected 'New Feature', got '%s'", delta.Added[0].Name)
	}

	if len(delta.Added[0].Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(delta.Added[0].Scenarios))
	}

	if delta.Added[0].Scenarios[0].Name != "Feature works" {
		t.Errorf("expected 'Feature works', got '%s'", delta.Added[0].Scenarios[0].Name)
	}

	// Other sections should be empty
	if len(delta.Modified) != 0 || len(delta.Removed) != 0 || len(delta.Renamed) != 0 {
		t.Error("expected other delta sections to be empty")
	}
}

// TestExtractDeltas_MODIFIED tests extracting MODIFIED requirements.
func TestExtractDeltas_MODIFIED(t *testing.T) {
	input := `# Delta Spec

## MODIFIED Requirements

### Requirement: Updated Feature
The system SHALL update existing feature.

#### Scenario: Update works
- **WHEN** update applied
- **THEN** feature is updated
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Modified) != 1 {
		t.Fatalf("expected 1 MODIFIED requirement, got %d", len(delta.Modified))
	}

	if delta.Modified[0].Name != "Updated Feature" {
		t.Errorf("expected 'Updated Feature', got '%s'", delta.Modified[0].Name)
	}

	if len(delta.Modified[0].Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(delta.Modified[0].Scenarios))
	}
}

// TestExtractDeltas_REMOVED tests extracting REMOVED requirements.
func TestExtractDeltas_REMOVED(t *testing.T) {
	input := `# Delta Spec

## REMOVED Requirements

### Requirement: Old Feature
**Reason**: No longer needed
**Migration**: Use new feature instead
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Removed) != 1 {
		t.Fatalf("expected 1 REMOVED requirement, got %d", len(delta.Removed))
	}

	if delta.Removed[0].Name != "Old Feature" {
		t.Errorf("expected 'Old Feature', got '%s'", delta.Removed[0].Name)
	}

	// Check content includes reason and migration
	if !strings.Contains(delta.Removed[0].Content, "No longer needed") {
		t.Error("expected content to contain reason")
	}
	if !strings.Contains(delta.Removed[0].Content, "Use new feature instead") {
		t.Error("expected content to contain migration")
	}
}

// TestExtractDeltas_RENAMED tests extracting RENAMED requirements.
func TestExtractDeltas_RENAMED(t *testing.T) {
	input := `# Delta Spec

## RENAMED Requirements

- FROM: ` + "`### Requirement: Old Name`" + `
- TO: ` + "`### Requirement: New Name`" + `
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Renamed) != 1 {
		t.Fatalf("expected 1 RENAMED requirement, got %d", len(delta.Renamed))
	}

	if delta.Renamed[0].From != "Old Name" {
		t.Errorf("expected From='Old Name', got '%s'", delta.Renamed[0].From)
	}

	if delta.Renamed[0].To != "New Name" {
		t.Errorf("expected To='New Name', got '%s'", delta.Renamed[0].To)
	}
}

// TestExtractDeltas_MultipleOperations tests multiple delta operations in one spec.
func TestExtractDeltas_MultipleOperations(t *testing.T) {
	input := `# Delta Spec

## ADDED Requirements

### Requirement: First Addition
New feature one.

### Requirement: Second Addition
New feature two.

## MODIFIED Requirements

### Requirement: Modified One
Updated feature.

## REMOVED Requirements

### Requirement: Removed One
Old feature.

## RENAMED Requirements

- FROM: ` + "`### Requirement: Old Name`" + `
- TO: ` + "`### Requirement: New Name`" + `
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	// Check counts
	if len(delta.Added) != 2 {
		t.Errorf("expected 2 ADDED requirements, got %d", len(delta.Added))
	}
	if len(delta.Modified) != 1 {
		t.Errorf("expected 1 MODIFIED requirement, got %d", len(delta.Modified))
	}
	if len(delta.Removed) != 1 {
		t.Errorf("expected 1 REMOVED requirement, got %d", len(delta.Removed))
	}
	if len(delta.Renamed) != 1 {
		t.Errorf("expected 1 RENAMED requirement, got %d", len(delta.Renamed))
	}

	// Check names
	if delta.Added[0].Name != "First Addition" {
		t.Errorf("expected 'First Addition', got '%s'", delta.Added[0].Name)
	}
	if delta.Added[1].Name != "Second Addition" {
		t.Errorf("expected 'Second Addition', got '%s'", delta.Added[1].Name)
	}
}

// TestExtractDeltas_WithScenarios tests delta requirements with scenarios.
func TestExtractDeltas_WithScenarios(t *testing.T) {
	input := `# Delta Spec

## ADDED Requirements

### Requirement: Feature With Scenarios
The system SHALL provide feature.

#### Scenario: First scenario
- **WHEN** first action
- **THEN** first result

#### Scenario: Second scenario
- **WHEN** second action
- **THEN** second result
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Added) != 1 {
		t.Fatalf("expected 1 ADDED requirement, got %d", len(delta.Added))
	}

	req := delta.Added[0]
	if len(req.Scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(req.Scenarios))
	}

	if req.Scenarios[0].Name != "First scenario" {
		t.Errorf("expected 'First scenario', got '%s'", req.Scenarios[0].Name)
	}
	if req.Scenarios[1].Name != "Second scenario" {
		t.Errorf("expected 'Second scenario', got '%s'", req.Scenarios[1].Name)
	}
}

// TestExtractDeltas_EmptyDocument tests empty document handling.
func TestExtractDeltas_EmptyDocument(t *testing.T) {
	input := `# Spec

No delta operations here.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Added) != 0 || len(delta.Modified) != 0 || len(delta.Removed) != 0 ||
		len(delta.Renamed) != 0 {
		t.Error("expected all delta sections to be empty")
	}
}

// TestExtractDeltas_CaseInsensitive tests case-insensitive section matching.
func TestExtractDeltas_CaseInsensitive(t *testing.T) {
	input := `# Delta Spec

## Added Requirements

### Requirement: Feature One
Added feature.

## Modified Requirements

### Requirement: Feature Two
Modified feature.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Added) != 1 {
		t.Errorf("expected 1 ADDED requirement (case-insensitive), got %d", len(delta.Added))
	}
	if len(delta.Modified) != 1 {
		t.Errorf("expected 1 MODIFIED requirement (case-insensitive), got %d", len(delta.Modified))
	}
}

// TestExtractDeltas_HierarchyValidation tests scenario hierarchy validation in deltas.
func TestExtractDeltas_HierarchyValidation(t *testing.T) {
	input := `# Delta Spec

## ADDED Requirements

#### Scenario: Orphan scenario
- **WHEN** no requirement
- **THEN** should error
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = ExtractDeltas(doc)
	if err == nil {
		t.Error("expected error for scenario without requirement, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "but no requirements defined") {
		t.Errorf("expected error about missing requirements, got: %v", err)
	}
}

// TestExtractDeltas_CodeBlocksIgnored tests that deltas in code blocks are ignored.
func TestExtractDeltas_CodeBlocksIgnored(t *testing.T) {
	input := `# Delta Spec

Example of delta format:

` + "```markdown" + `
## ADDED Requirements

### Requirement: Fake Requirement
This is just an example and should not be extracted.
` + "```" + `

## ADDED Requirements

### Requirement: Real Requirement
This should be extracted.
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	// Should only extract the real requirement, not the example in code block
	if len(delta.Added) != 1 {
		t.Fatalf("expected 1 ADDED requirement (code block ignored), got %d", len(delta.Added))
	}

	if delta.Added[0].Name != "Real Requirement" {
		t.Errorf("expected 'Real Requirement', got '%s'", delta.Added[0].Name)
	}

	// Check no fake content leaked
	if strings.Contains(delta.Added[0].Content, "Fake Requirement") {
		t.Error("should not extract requirements from code blocks")
	}
}

// TestExtractDeltas_MultipleRENAMED tests multiple RENAMED operations.
func TestExtractDeltas_MultipleRENAMED(t *testing.T) {
	input := `# Delta Spec

## RENAMED Requirements

- FROM: ` + "`### Requirement: Old Name One`" + `
- TO: ` + "`### Requirement: New Name One`" + `
- FROM: ` + "`### Requirement: Old Name Two`" + `
- TO: ` + "`### Requirement: New Name Two`" + `
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	if len(delta.Renamed) != 2 {
		t.Fatalf("expected 2 RENAMED operations, got %d", len(delta.Renamed))
	}

	if delta.Renamed[0].From != "Old Name One" {
		t.Errorf("expected 'Old Name One', got '%s'", delta.Renamed[0].From)
	}
	if delta.Renamed[0].To != "New Name One" {
		t.Errorf("expected 'New Name One', got '%s'", delta.Renamed[0].To)
	}

	if delta.Renamed[1].From != "Old Name Two" {
		t.Errorf("expected 'Old Name Two', got '%s'", delta.Renamed[1].From)
	}
	if delta.Renamed[1].To != "New Name Two" {
		t.Errorf("expected 'New Name Two', got '%s'", delta.Renamed[1].To)
	}
}

// TestExtractDeltas_RENAMEDVariations tests different RENAMED format variations.
func TestExtractDeltas_RENAMEDVariations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFrom string
		wantTo   string
	}{
		{
			name: "with backticks",
			input: `## RENAMED Requirements
- FROM: ` + "`### Requirement: Old Name`" + `
- TO: ` + "`### Requirement: New Name`",
			wantFrom: "Old Name",
			wantTo:   "New Name",
		},
		{
			name: "without backticks",
			input: `## RENAMED Requirements
- FROM: ### Requirement: Old Name
- TO: ### Requirement: New Name`,
			wantFrom: "Old Name",
			wantTo:   "New Name",
		},
		{
			name: "with extra whitespace",
			input: `## RENAMED Requirements
- FROM:   ` + "`  ### Requirement:   Old Name  `" + `
- TO:   ` + "`  ### Requirement:   New Name  `",
			wantFrom: "Old Name",
			wantTo:   "New Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			delta, err := ExtractDeltas(doc)
			if err != nil {
				t.Fatalf("ExtractDeltas failed: %v", err)
			}

			if len(delta.Renamed) != 1 {
				t.Fatalf("expected 1 RENAMED operation, got %d", len(delta.Renamed))
			}

			if delta.Renamed[0].From != tt.wantFrom {
				t.Errorf("expected From='%s', got '%s'", tt.wantFrom, delta.Renamed[0].From)
			}
			if delta.Renamed[0].To != tt.wantTo {
				t.Errorf("expected To='%s', got '%s'", tt.wantTo, delta.Renamed[0].To)
			}
		})
	}
}

// TestExtractDeltas_RealWorldExample tests a comprehensive real-world delta spec.
func TestExtractDeltas_RealWorldExample(t *testing.T) {
	input := `# Change: Add Two-Factor Authentication

## ADDED Requirements

### Requirement: Two-Factor Authentication
Users MUST provide a second factor during login.

#### Scenario: OTP required
- **WHEN** valid credentials are provided
- **THEN** an OTP challenge is required

#### Scenario: OTP verification
- **WHEN** correct OTP is entered
- **THEN** user is authenticated

## MODIFIED Requirements

### Requirement: User Authentication
Updated to support 2FA.

#### Scenario: Login with 2FA
- **WHEN** user logs in with 2FA enabled
- **THEN** OTP is requested

## REMOVED Requirements

### Requirement: Simple Password Login
**Reason**: Security requirements changed
**Migration**: Use 2FA authentication

## RENAMED Requirements

- FROM: ` + "`### Requirement: Login`" + `
- TO: ` + "`### Requirement: User Authentication`" + `
`

	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	delta, err := ExtractDeltas(doc)
	if err != nil {
		t.Fatalf("ExtractDeltas failed: %v", err)
	}

	// Verify ADDED requirements
	if len(delta.Added) != 1 {
		t.Fatalf("expected 1 ADDED requirement, got %d", len(delta.Added))
	}
	if delta.Added[0].Name != "Two-Factor Authentication" {
		t.Errorf("expected 'Two-Factor Authentication', got '%s'", delta.Added[0].Name)
	}
	if len(delta.Added[0].Scenarios) != 2 {
		t.Fatalf("expected 2 scenarios in ADDED requirement, got %d", len(delta.Added[0].Scenarios))
	}
	if delta.Added[0].Scenarios[0].Name != "OTP required" {
		t.Errorf("expected scenario 'OTP required', got '%s'", delta.Added[0].Scenarios[0].Name)
	}

	// Verify MODIFIED requirements
	if len(delta.Modified) != 1 {
		t.Fatalf("expected 1 MODIFIED requirement, got %d", len(delta.Modified))
	}
	if delta.Modified[0].Name != "User Authentication" {
		t.Errorf("expected 'User Authentication', got '%s'", delta.Modified[0].Name)
	}
	if len(delta.Modified[0].Scenarios) != 1 {
		t.Fatalf(
			"expected 1 scenario in MODIFIED requirement, got %d",
			len(delta.Modified[0].Scenarios),
		)
	}

	// Verify REMOVED requirements
	if len(delta.Removed) != 1 {
		t.Fatalf("expected 1 REMOVED requirement, got %d", len(delta.Removed))
	}
	if delta.Removed[0].Name != "Simple Password Login" {
		t.Errorf("expected 'Simple Password Login', got '%s'", delta.Removed[0].Name)
	}
	if !strings.Contains(delta.Removed[0].Content, "Security requirements changed") {
		t.Error("expected REMOVED content to contain reason")
	}

	// Verify RENAMED requirements
	if len(delta.Renamed) != 1 {
		t.Fatalf("expected 1 RENAMED requirement, got %d", len(delta.Renamed))
	}
	if delta.Renamed[0].From != "Login" {
		t.Errorf("expected From='Login', got '%s'", delta.Renamed[0].From)
	}
	if delta.Renamed[0].To != "User Authentication" {
		t.Errorf("expected To='User Authentication', got '%s'", delta.Renamed[0].To)
	}
}
