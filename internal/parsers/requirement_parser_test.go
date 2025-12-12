package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestParseRequirements(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int // number of requirements expected
	}{
		{
			name: "Single requirement with scenario",
			content: `# Test Spec

## Requirements

### Requirement: Feature One
The system SHALL do something.

#### Scenario: Success case
- **WHEN** action occurs
- **THEN** result happens
`,
			expected: 1,
		},
		{
			name: "Multiple requirements",
			content: `# Test Spec

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
`,
			expected: 2,
		},
		{
			name: "Requirements with multiple scenarios",
			content: `# Test Spec

## Requirements

### Requirement: Complex Feature
The system SHALL handle complexity.

#### Scenario: Happy path
- **WHEN** normal operation
- **THEN** success

#### Scenario: Error path
- **WHEN** error occurs
- **THEN** handle gracefully
`,
			expected: 1,
		},
		{
			name: "No requirements",
			content: `# Test Spec

## Requirements

No actual requirement blocks here.
`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(
				tmpDir,
				"spec.md",
			)
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			reqs, err := ParseRequirements(
				filePath,
			)
			if err != nil {
				t.Fatalf(
					"ParseRequirements failed: %v",
					err,
				)
			}

			if len(reqs) != tt.expected {
				t.Errorf(
					"Expected %d requirements, got %d",
					tt.expected,
					len(reqs),
				)
			}

			// Verify that each requirement has a name and raw content
			for i, req := range reqs {
				if req.Name == "" {
					t.Errorf(
						"Requirement %d has empty name",
						i,
					)
				}
				if req.Raw == "" {
					t.Errorf(
						"Requirement %d has empty raw content",
						i,
					)
				}
				if req.HeaderLine == "" {
					t.Errorf(
						"Requirement %d has empty header line",
						i,
					)
				}
			}
		})
	}
}

func TestParseRequirements_MultipleH2Sections(
	t *testing.T,
) {
	content := `# Test Spec

## Requirements

### Requirement: First Feature
Content for first feature.

## Another Section

Some content here.

## More Requirements

### Requirement: Second Feature
Content for second feature.
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	reqs, err := ParseRequirements(filePath)
	if err != nil {
		t.Fatalf(
			"ParseRequirements failed: %v",
			err,
		)
	}

	// Should parse requirements from all sections
	if len(reqs) != 2 {
		t.Errorf(
			"Expected 2 requirements, got %d",
			len(reqs),
		)
	}

	if len(reqs) > 0 &&
		reqs[0].Name != "First Feature" {
		t.Errorf(
			"Expected first requirement name 'First Feature', got %q",
			reqs[0].Name,
		)
	}

	if len(reqs) > 1 &&
		reqs[1].Name != "Second Feature" {
		t.Errorf(
			"Expected second requirement name 'Second Feature', got %q",
			reqs[1].Name,
		)
	}
}

func TestParseScenarios(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "Single scenario",
			content: `### Requirement: Test
Description.

#### Scenario: Success case
- **WHEN** action
- **THEN** result
`,
			expected: []string{"Success case"},
		},
		{
			name: "Multiple scenarios",
			content: `### Requirement: Test
Description.

#### Scenario: Happy path
- **WHEN** normal
- **THEN** success

#### Scenario: Error path
- **WHEN** error
- **THEN** handle
`,
			expected: []string{
				"Happy path",
				"Error path",
			},
		},
		{
			name: "No scenarios",
			content: `### Requirement: Test
Description without scenarios.
`,
			expected: nil,
		},
		{
			name: "Scenario with extra whitespace",
			content: `### Requirement: Test

####    Scenario:    Whitespace Test
- **WHEN** something
- **THEN** result
`,
			expected: []string{"Whitespace Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scenarios := ParseScenarios(
				tt.content,
			)

			if len(
				scenarios,
			) != len(
				tt.expected,
			) {
				t.Errorf(
					"Expected %d scenarios, got %d",
					len(tt.expected),
					len(scenarios),
				)
			}

			for i, expected := range tt.expected {
				if i >= len(scenarios) {
					break
				}
				if scenarios[i] != expected {
					t.Errorf(
						"Scenario %d: expected %q, got %q",
						i,
						expected,
						scenarios[i],
					)
				}
			}
		})
	}
}

func TestNormalizeRequirementName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Feature One", "feature one"},
		{"  Feature One  ", "feature one"},
		{"FEATURE ONE", "feature one"},
		{"feature one", "feature one"},
		{"\tFeature\tOne\t", "feature\tone"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeRequirementName(
				tt.input,
			)
			if result != tt.expected {
				t.Errorf(
					"Expected %q, got %q",
					tt.expected,
					result,
				)
			}
		})
	}
}
