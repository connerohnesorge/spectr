package markdown

import (
	"testing"
)

const testValidation = "validation"

func TestParseSpec(t *testing.T) {
	content := []byte(`# Title

## Overview

This is an overview section.

## Requirements

### Requirement: User Authentication

Users must be able to authenticate.

#### Scenario: Valid credentials

- **WHEN** user provides valid credentials
- **THEN** user is authenticated

#### Scenario: Invalid credentials

- **WHEN** user provides invalid credentials
- **THEN** authentication fails

### Requirement: Session Management

Sessions should be managed securely.
`)

	spec, errors := ParseSpec(content)

	if len(errors) > 0 {
		t.Errorf(
			"ParseSpec returned errors: %v",
			errors,
		)
	}

	if spec.Root == nil {
		t.Fatal("ParseSpec returned nil root")
	}

	// Check sections
	if len(spec.Sections) == 0 {
		t.Error("No sections extracted")
	}

	if _, found := spec.Sections["Overview"]; !found {
		t.Error("Overview section not found")
	}

	if _, found := spec.Sections["Requirements"]; !found {
		t.Error("Requirements section not found")
	}

	// Check requirements
	if len(spec.Requirements) != 2 {
		t.Errorf(
			"Expected 2 requirements, got %d",
			len(spec.Requirements),
		)
	}

	if len(spec.Requirements) == 0 {
		return
	}

	req := spec.Requirements[0]
	if req.Name != "User Authentication" {
		t.Errorf(
			"Expected requirement name 'User Authentication', got '%s'",
			req.Name,
		)
	}

	// Check scenarios
	if len(req.Scenarios) != 2 {
		t.Errorf(
			"Expected 2 scenarios for User Authentication, got %d",
			len(req.Scenarios),
		)
	}

	if len(req.Scenarios) > 0 &&
		req.Scenarios[0].Name != "Valid credentials" {
		t.Errorf(
			"Expected scenario name 'Valid credentials', got '%s'",
			req.Scenarios[0].Name,
		)
	}
}

func TestExtractSections(t *testing.T) {
	content := []byte(`# Document

## First Section

Content of first section.

## Second Section

Content of second section.

### Subsection

Subsection content.
`)

	sections := ExtractSections(content)

	if len(sections) < 2 {
		t.Errorf(
			"Expected at least 2 sections, got %d",
			len(sections),
		)
	}

	if section, found := sections["First Section"]; found {
		if section.Level != 2 {
			t.Errorf(
				"Expected First Section level 2, got %d",
				section.Level,
			)
		}
		if section.Start < 0 {
			t.Error(
				"Section start should be >= 0",
			)
		}
		if section.End <= section.Start {
			t.Error(
				"Section end should be > start",
			)
		}
	} else {
		t.Error("First Section not found")
	}

	if _, found := sections["Second Section"]; !found {
		t.Error("Second Section not found")
	}
}

func TestExtractRequirements(t *testing.T) {
	content := []byte(`## Functional Requirements

### Requirement: Data Validation

All data must be validated.

#### Scenario: Valid input

- **WHEN** valid data is provided
- **THEN** data is accepted

### Requirement: Error Handling

Errors must be handled gracefully.

## Non-Functional Requirements

### Requirement: Performance

System must respond within 100ms.
`)

	requirements := ExtractRequirements(content)

	if len(requirements) != 3 {
		t.Errorf(
			"Expected 3 requirements, got %d",
			len(requirements),
		)
	}

	// Check first requirement
	if len(requirements) > 0 {
		req := requirements[0]
		if req.Name != "Data Validation" {
			t.Errorf(
				"Expected 'Data Validation', got '%s'",
				req.Name,
			)
		}
		if req.Section != "Functional Requirements" {
			t.Errorf(
				"Expected section 'Functional Requirements', got '%s'",
				req.Section,
			)
		}
		if len(req.Scenarios) != 1 {
			t.Errorf(
				"Expected 1 scenario, got %d",
				len(req.Scenarios),
			)
		}
	}

	// Check last requirement is in different section
	if len(requirements) <= 2 {
		return
	}

	req := requirements[2]
	if req.Section != "Non-Functional Requirements" {
		t.Errorf(
			"Expected section 'Non-Functional Requirements', got '%s'",
			req.Section,
		)
	}
}

func TestFindSection(t *testing.T) {
	content := []byte(`## Overview

Overview content.

## Requirements

Requirements content.

## Implementation Notes

Implementation details.
`)

	tests := []struct {
		name      string
		search    string
		wantFound bool
		wantLevel int
	}{
		{"exact match", "Overview", true, 2},
		{"case insensitive", "OVERVIEW", true, 2},
		{"mixed case", "overview", true, 2},
		{
			"another section",
			"Requirements",
			true,
			2,
		},
		{"not found", "Nonexistent", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section, found := FindSection(
				content,
				tt.search,
			)

			if found != tt.wantFound {
				t.Errorf(
					"FindSection(%q) found = %v, want %v",
					tt.search,
					found,
					tt.wantFound,
				)
			}

			if found &&
				section.Level != tt.wantLevel {
				t.Errorf(
					"FindSection(%q) level = %d, want %d",
					tt.search,
					section.Level,
					tt.wantLevel,
				)
			}
		})
	}
}

func TestExtractWikilinks(t *testing.T) {
	content := []byte(`## References

See [[validation]] for details.

Also check [[changes/my-change|My Change]] and [[spec#section]].

For more info: [[another-spec|Display Text#anchor-name]].
`)

	wikilinks := ExtractWikilinks(content)

	if len(wikilinks) != 4 {
		t.Errorf(
			"Expected 4 wikilinks, got %d",
			len(wikilinks),
		)
	}

	// Test first wikilink (simple)
	if len(wikilinks) > 0 {
		wl := wikilinks[0]
		if wl.Target != testValidation {
			t.Errorf(
				"Expected target 'validation', got '%s'",
				wl.Target,
			)
		}
		if wl.Display != "" {
			t.Errorf(
				"Expected empty display, got '%s'",
				wl.Display,
			)
		}
		if wl.Anchor != "" {
			t.Errorf(
				"Expected empty anchor, got '%s'",
				wl.Anchor,
			)
		}
	}

	// Test second wikilink (with display text)
	if len(wikilinks) > 1 {
		wl := wikilinks[1]
		if wl.Target != "changes/my-change" {
			t.Errorf(
				"Expected target 'changes/my-change', got '%s'",
				wl.Target,
			)
		}
		if wl.Display != "My Change" {
			t.Errorf(
				"Expected display 'My Change', got '%s'",
				wl.Display,
			)
		}
	}

	// Test third wikilink (with anchor)
	if len(wikilinks) > 2 {
		wl := wikilinks[2]
		if wl.Target != "spec" {
			t.Errorf(
				"Expected target 'spec', got '%s'",
				wl.Target,
			)
		}
		if wl.Anchor != "section" {
			t.Errorf(
				"Expected anchor 'section', got '%s'",
				wl.Anchor,
			)
		}
	}

	// Test fourth wikilink (full format)
	if len(wikilinks) <= 3 {
		return
	}

	wl := wikilinks[3]
	if wl.Target != "another-spec" {
		t.Errorf(
			"Expected target 'another-spec', got '%s'",
			wl.Target,
		)
	}
	if wl.Display != "Display Text" {
		t.Errorf(
			"Expected display 'Display Text', got '%s'",
			wl.Display,
		)
	}
	if wl.Anchor != "anchor-name" {
		t.Errorf(
			"Expected anchor 'anchor-name', got '%s'",
			wl.Anchor,
		)
	}
}

func TestFindRequirement(t *testing.T) {
	content := []byte(`## Requirements

### Requirement: User Login

Users can login.

### Requirement: User Logout

Users can logout.
`)

	tests := []struct {
		name      string
		search    string
		wantFound bool
		wantName  string
	}{
		{
			"exact match",
			"User Login",
			true,
			"User Login",
		},
		{
			"case insensitive",
			"user login",
			true,
			"User Login",
		},
		{
			"another requirement",
			"User Logout",
			true,
			"User Logout",
		},
		{"not found", "Nonexistent", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, found := FindRequirement(
				content,
				tt.search,
			)

			if found != tt.wantFound {
				t.Errorf(
					"FindRequirement(%q) found = %v, want %v",
					tt.search,
					found,
					tt.wantFound,
				)
			}

			if found && req.Name != tt.wantName {
				t.Errorf(
					"FindRequirement(%q) name = %s, want %s",
					tt.search,
					req.Name,
					tt.wantName,
				)
			}
		})
	}
}

func TestFindScenario(t *testing.T) {
	content := []byte(`### Requirement: Test

#### Scenario: Happy Path

Test passes.

#### Scenario: Error Case

Test fails.
`)

	tests := []struct {
		name      string
		search    string
		wantFound bool
		wantName  string
	}{
		{
			"exact match",
			"Happy Path",
			true,
			"Happy Path",
		},
		{
			"case insensitive",
			"happy path",
			true,
			"Happy Path",
		},
		{
			"another scenario",
			"Error Case",
			true,
			"Error Case",
		},
		{"not found", "Nonexistent", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scenario, found := FindScenario(
				content,
				tt.search,
			)

			if found != tt.wantFound {
				t.Errorf(
					"FindScenario(%q) found = %v, want %v",
					tt.search,
					found,
					tt.wantFound,
				)
			}

			if found &&
				scenario.Name != tt.wantName {
				t.Errorf(
					"FindScenario(%q) name = %s, want %s",
					tt.search,
					scenario.Name,
					tt.wantName,
				)
			}
		})
	}
}

func TestHasRequirement(t *testing.T) {
	content := []byte(`### Requirement: Existing

Content.
`)

	if !HasRequirement(content, "Existing") {
		t.Error(
			"HasRequirement should return true for existing requirement",
		)
	}

	if !HasRequirement(content, "existing") {
		t.Error(
			"HasRequirement should be case insensitive",
		)
	}

	if HasRequirement(content, "Nonexistent") {
		t.Error(
			"HasRequirement should return false for nonexistent requirement",
		)
	}
}

func TestHasSection(t *testing.T) {
	content := []byte(`## Existing Section

Content.
`)

	if !HasSection(content, "Existing Section") {
		t.Error(
			"HasSection should return true for existing section",
		)
	}

	if !HasSection(content, "existing section") {
		t.Error(
			"HasSection should be case insensitive",
		)
	}

	if HasSection(content, "Nonexistent") {
		t.Error(
			"HasSection should return false for nonexistent section",
		)
	}
}

func TestGetSectionContent(t *testing.T) {
	content := []byte(`## Test Section

This is the content.
`)

	result := GetSectionContent(
		content,
		"Test Section",
	)

	if result == "" {
		t.Error(
			"GetSectionContent should return non-empty content",
		)
	}

	result = GetSectionContent(
		content,
		"Nonexistent",
	)
	if result != "" {
		t.Error(
			"GetSectionContent should return empty string for nonexistent section",
		)
	}
}

func TestGetRequirementNames(t *testing.T) {
	content := []byte(`### Requirement: First

### Requirement: Second

### Requirement: Third
`)

	names := GetRequirementNames(content)

	if len(names) != 3 {
		t.Errorf(
			"Expected 3 names, got %d",
			len(names),
		)
	}

	expected := []string{
		"First",
		"Second",
		"Third",
	}
	for i, exp := range expected {
		if i < len(names) && names[i] != exp {
			t.Errorf(
				"names[%d] = %s, want %s",
				i,
				names[i],
				exp,
			)
		}
	}
}

func TestGetSectionNames(t *testing.T) {
	content := []byte(`## Section A

## Section B

## Section C
`)

	names := GetSectionNames(content)

	if len(names) != 3 {
		t.Errorf(
			"Expected 3 section names, got %d",
			len(names),
		)
	}

	// Check all expected sections are present (order may vary due to map)
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	for _, expected := range []string{"Section A", "Section B", "Section C"} {
		if !nameSet[expected] {
			t.Errorf(
				"Missing section '%s'",
				expected,
			)
		}
	}
}

func TestCountRequirements(t *testing.T) {
	content := []byte(`### Requirement: One

### Requirement: Two

### Requirement: Three
`)

	count := CountRequirements(content)

	if count != 3 {
		t.Errorf(
			"CountRequirements = %d, want 3",
			count,
		)
	}
}

func TestCountScenarios(t *testing.T) {
	content := []byte(`### Requirement: Test

#### Scenario: One

#### Scenario: Two

### Requirement: Another

#### Scenario: Three
`)

	count := CountScenarios(content)

	if count != 3 {
		t.Errorf(
			"CountScenarios = %d, want 3",
			count,
		)
	}
}

func TestCountWikilinks(t *testing.T) {
	content := []byte(
		`See [[one]], [[two]], and [[three]].
`,
	)

	count := CountWikilinks(content)

	if count != 3 {
		t.Errorf(
			"CountWikilinks = %d, want 3",
			count,
		)
	}
}

func TestParseSpecEmptyContent(t *testing.T) {
	content := []byte("")

	spec, errors := ParseSpec(content)

	if spec == nil {
		t.Fatal(
			"ParseSpec should not return nil spec for empty content",
		)
	}

	if len(errors) > 0 {
		t.Errorf(
			"ParseSpec should not return errors for empty content, got %v",
			errors,
		)
	}

	if len(spec.Sections) != 0 {
		t.Error(
			"Empty content should have no sections",
		)
	}

	if len(spec.Requirements) != 0 {
		t.Error(
			"Empty content should have no requirements",
		)
	}
}

func TestExtractSectionsNested(t *testing.T) {
	content := []byte(`# Top Level

## Second Level

### Third Level

Content.
`)

	sections := ExtractSections(content)

	if section, found := sections["Second Level"]; found {
		if section.Level != 2 {
			t.Errorf(
				"Second Level should have level 2, got %d",
				section.Level,
			)
		}
	} else {
		t.Error("Second Level section not found")
	}

	if section, found := sections["Third Level"]; found {
		if section.Level != 3 {
			t.Errorf(
				"Third Level should have level 3, got %d",
				section.Level,
			)
		}
	} else {
		t.Error("Third Level section not found")
	}
}

func TestSpecSectionNode(t *testing.T) {
	content := []byte(`## My Section

Content here.
`)

	spec, _ := ParseSpec(content)

	section, found := spec.Sections["My Section"]
	if !found {
		t.Fatal("Section not found")
	}

	if section.Node == nil {
		t.Error("Section.Node should not be nil")
	}

	if section.Node.Level() != 2 {
		t.Errorf(
			"Section.Node.Level() = %d, want 2",
			section.Node.Level(),
		)
	}
}

func TestRequirementNode(t *testing.T) {
	content := []byte(
		`### Requirement: My Requirement

Content here.
`,
	)

	spec, _ := ParseSpec(content)

	if len(spec.Requirements) == 0 {
		t.Fatal("No requirements found")
	}

	req := spec.Requirements[0]
	if req.Node == nil {
		t.Error(
			"Requirement.Node should not be nil",
		)
	}

	if req.Node.Name() != "My Requirement" {
		t.Errorf(
			"Requirement.Node.Name() = %s, want 'My Requirement'",
			req.Node.Name(),
		)
	}
}

func TestWikilinkNode(t *testing.T) {
	content := []byte(
		`See [[target|display#anchor]].
`,
	)

	wikilinks := ExtractWikilinks(content)

	if len(wikilinks) == 0 {
		t.Fatal("No wikilinks found")
	}

	wl := wikilinks[0]
	if wl.Node == nil {
		t.Error("Wikilink.Node should not be nil")
	}

	if string(wl.Node.Target()) != "target" {
		t.Errorf(
			"Wikilink.Node.Target() = %s, want 'target'",
			wl.Node.Target(),
		)
	}

	if string(wl.Node.Display()) != "display" {
		t.Errorf(
			"Wikilink.Node.Display() = %s, want 'display'",
			wl.Node.Display(),
		)
	}

	if string(wl.Node.Anchor()) != "anchor" {
		t.Errorf(
			"Wikilink.Node.Anchor() = %s, want 'anchor'",
			wl.Node.Anchor(),
		)
	}
}
