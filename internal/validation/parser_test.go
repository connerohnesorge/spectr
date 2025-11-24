package validation

import (
	"reflect"
	"testing"
)

func TestExtractSections(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
	}{
		{
			name: "valid markdown with multiple sections",
			content: `## Purpose
This is the purpose section.
It has multiple lines.

## Context
This is the context section.

## Goals
- Goal 1
- Goal 2`,
			expected: map[string]string{
				"Purpose": "This is the purpose section.\nIt has multiple lines.",
				"Context": "This is the context section.",
				"Goals":   "- Goal 1\n- Goal 2",
			},
		},
		{
			name: "single section",
			content: `## Purpose
This is the only section.`,
			expected: map[string]string{
				"Purpose": "This is the only section.",
			},
		},
		{
			name:     "empty content",
			content:  "",
			expected: make(map[string]string),
		},
		{
			name: "no sections",
			content: `This is just text
without any headers`,
			expected: make(map[string]string),
		},
		{
			name: "section with trailing whitespace",
			content: `## Purpose
Content here`,
			expected: map[string]string{
				"Purpose": "Content here",
			},
		},
		{
			name: "section stops at next section",
			content: `## First
Content 1
## Second
Content 2`,
			expected: map[string]string{
				"First":  "Content 1",
				"Second": "Content 2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSections(tt.content)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractSections() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractRequirements(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Requirement
	}{
		{
			name: "single requirement with scenario",
			content: `### Requirement: User Authentication
The system SHALL authenticate users.

#### Scenario: Valid credentials
- WHEN user provides valid credentials
- THEN system grants access`,
			expected: []Requirement{
				{
					Name: "User Authentication",
					Content: `The system SHALL authenticate users.

#### Scenario: Valid credentials
- WHEN user provides valid credentials
- THEN system grants access`,
					Scenarios: []string{
						`#### Scenario: Valid credentials
- WHEN user provides valid credentials
- THEN system grants access`,
					},
				},
			},
		},
		{
			name: "multiple requirements",
			content: `### Requirement: Login
User login functionality.

#### Scenario: Success
User logs in successfully.

### Requirement: Logout
User logout functionality.

#### Scenario: Success
User logs out successfully.`,
			expected: []Requirement{
				{
					Name:    "Login",
					Content: "User login functionality.\n\n#### Scenario: Success\nUser logs in successfully.",
					Scenarios: []string{
						"#### Scenario: Success\nUser logs in successfully.",
					},
				},
				{
					Name:    "Logout",
					Content: "User logout functionality.\n\n#### Scenario: Success\nUser logs out successfully.",
					Scenarios: []string{
						"#### Scenario: Success\nUser logs out successfully.",
					},
				},
			},
		},
		{
			name: "requirement without scenario",
			content: `### Requirement: Basic Feature
This is a basic feature without scenarios.`,
			expected: []Requirement{
				{
					Name:      "Basic Feature",
					Content:   "This is a basic feature without scenarios.",
					Scenarios: make([]string, 0),
				},
			},
		},
		{
			name:     "empty content",
			content:  "",
			expected: make([]Requirement, 0),
		},
		{
			name: "requirement with multiple scenarios",
			content: `### Requirement: Multi-scenario
Feature with multiple scenarios.

#### Scenario: First
First scenario.

#### Scenario: Second
Second scenario.`,
			expected: []Requirement{
				{
					Name: "Multi-scenario",
					Content: `Feature with multiple scenarios.

#### Scenario: First
First scenario.

#### Scenario: Second
Second scenario.`,
					Scenarios: []string{
						"#### Scenario: First\nFirst scenario.",
						"#### Scenario: Second\nSecond scenario.",
					},
				},
			},
		},
		{
			name: "requirement stops at section boundary",
			content: `### Requirement: First
Content here.

## Next Section
Not part of requirement.`,
			expected: []Requirement{
				{
					Name:      "First",
					Content:   "Content here.",
					Scenarios: make([]string, 0),
				},
			},
		},
		{
			name: "malformed requirement header ignored",
			content: `### NotARequirement
This should be ignored.

### Requirement: Valid
This is valid.`,
			expected: []Requirement{
				{
					Name:      "Valid",
					Content:   "This is valid.",
					Scenarios: make([]string, 0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractRequirements(tt.content)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractRequirements() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractScenarios(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "single scenario",
			content: `#### Scenario: Success case
- WHEN user performs action
- THEN expected result`,
			expected: []string{
				`#### Scenario: Success case
- WHEN user performs action
- THEN expected result`,
			},
		},
		{
			name: "multiple scenarios",
			content: `#### Scenario: First
Content 1

#### Scenario: Second
Content 2`,
			expected: []string{
				"#### Scenario: First\nContent 1",
				"#### Scenario: Second\nContent 2",
			},
		},
		{
			name:     "no scenarios",
			content:  "Just some text without scenarios",
			expected: make([]string, 0),
		},
		{
			name:     "empty content",
			content:  "",
			expected: make([]string, 0),
		},
		{
			name: "scenario with complex content",
			content: `#### Scenario: Complex
- **WHEN** user does something
- **AND** another condition
- **THEN** result happens
- **AND** another result`,
			expected: []string{
				`#### Scenario: Complex
- **WHEN** user does something
- **AND** another condition
- **THEN** result happens
- **AND** another result`,
			},
		},
		{
			name: "scenario stops at requirement boundary",
			content: `#### Scenario: First
Content here
### Requirement: Next
Not part of scenario`,
			expected: []string{
				"#### Scenario: First\nContent here",
			},
		},
		{
			name: "scenario with trailing whitespace",
			content: `#### Scenario: Test
Content`,
			expected: []string{
				"#### Scenario: Test\nContent",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractScenarios(tt.content)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractScenarios() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContainsShallOrMust(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "contains SHALL",
			text:     "The system SHALL authenticate users",
			expected: true,
		},
		{
			name:     "contains MUST",
			text:     "The system MUST validate input",
			expected: true,
		},
		{
			name:     "contains shall lowercase",
			text:     "The system shall authenticate users",
			expected: true,
		},
		{
			name:     "contains must lowercase",
			text:     "The system must validate input",
			expected: true,
		},
		{
			name:     "contains SHALL mixed case",
			text:     "The system ShAlL authenticate users",
			expected: true,
		},
		{
			name:     "contains neither",
			text:     "The system should authenticate users",
			expected: false,
		},
		{
			name:     "empty text",
			text:     "",
			expected: false,
		},
		{
			name:     "partial match not counted",
			text:     "marshall and mustard are not keywords",
			expected: false,
		},
		{
			name:     "ignores partial match but catches later whole word",
			text:     "mustard MUST comply",
			expected: true,
		},
		{
			name:     "word boundary respected",
			text:     "The system SHALL do something",
			expected: true,
		},
		{
			name:     "multiple occurrences",
			text:     "SHALL and MUST both present",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsShallOrMust(tt.text)
			if result != tt.expected {
				t.Errorf("ContainsShallOrMust(%q) = %v, want %v", tt.text, result, tt.expected)
			}
		})
	}
}

func TestNormalizeRequirementName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic normalization",
			input:    "User Authentication",
			expected: "user authentication",
		},
		{
			name:     "leading whitespace",
			input:    "  User Authentication",
			expected: "user authentication",
		},
		{
			name:     "trailing whitespace",
			input:    "User Authentication  ",
			expected: "user authentication",
		},
		{
			name:     "multiple spaces",
			input:    "User    Authentication",
			expected: "user authentication",
		},
		{
			name:     "mixed whitespace",
			input:    "  User    Authentication  ",
			expected: "user authentication",
		},
		{
			name:     "already normalized",
			input:    "user authentication",
			expected: "user authentication",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: "",
		},
		{
			name:     "special characters preserved",
			input:    "User-Authentication & Login",
			expected: "user-authentication & login",
		},
		{
			name:     "tabs and newlines",
			input:    "User\t\nAuthentication",
			expected: "user authentication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeRequirementName(tt.input)
			if result != tt.expected {
				t.Errorf(
					"NormalizeRequirementName(%q) = %q, want %q",
					tt.input,
					result,
					tt.expected,
				)
			}
		})
	}
}

// Test edge cases and integration scenarios
func TestFullSpecParsing(t *testing.T) {
	content := `## Purpose
Authentication system specification.

## Requirements

### Requirement: User Login
The system SHALL allow users to login.

#### Scenario: Valid credentials
- WHEN user provides valid credentials
- THEN system grants access

#### Scenario: Invalid credentials
- WHEN user provides invalid credentials
- THEN system denies access

### Requirement: Password Reset
The system MUST support password reset.

#### Scenario: Email sent
- WHEN user requests password reset
- THEN system sends email`

	sections := ExtractSections(content)
	if len(sections) != 2 {
		t.Errorf("Expected 2 sections, got %d", len(sections))
	}

	requirements := ExtractRequirements(content)
	if len(requirements) != 2 {
		t.Errorf("Expected 2 requirements, got %d", len(requirements))
	}

	if len(requirements[0].Scenarios) != 2 {
		t.Errorf(
			"Expected 2 scenarios for first requirement, got %d",
			len(requirements[0].Scenarios),
		)
	}

	if len(requirements[1].Scenarios) != 1 {
		t.Errorf(
			"Expected 1 scenario for second requirement, got %d",
			len(requirements[1].Scenarios),
		)
	}

	if !ContainsShallOrMust(requirements[0].Content) {
		t.Error("Expected first requirement to contain SHALL/MUST")
	}

	if !ContainsShallOrMust(requirements[1].Content) {
		t.Error("Expected second requirement to contain SHALL/MUST")
	}
}

func TestDuplicateDetectionViaNormalization(t *testing.T) {
	names := []string{
		"User Authentication",
		"  User Authentication  ",
		"user    authentication",
		"USER AUTHENTICATION",
	}

	normalized := make(map[string]bool)
	for _, name := range names {
		norm := NormalizeRequirementName(name)
		if normalized[norm] {
			// Expected - all normalize to same value
			continue
		}
		normalized[norm] = true
	}

	if len(normalized) != 1 {
		t.Errorf(
			"Expected all names to normalize to same value, got %d unique",
			len(normalized),
		)
	}
}

func TestMalformedMarkdownGracefulHandling(t *testing.T) {
	content := `## Section
### Requirement: Valid
Content here

#### Scenario: Test
Scenario content

##### Extra level ignored
### Another non-requirement header
### Requirement: Another Valid
More content`

	requirements := ExtractRequirements(content)
	if len(requirements) != 2 {
		t.Errorf(
			"Expected 2 requirements despite malformed content, got %d",
			len(requirements),
		)
	}
}
