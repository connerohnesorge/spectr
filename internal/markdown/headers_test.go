package markdown

import (
	"testing"
)

func TestMatchH2SectionHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "simple section",
			input:    "## Requirements",
			wantName: "Requirements",
			wantOk:   true,
		},
		{
			name:     "section with multiple words",
			input:    "## Purpose and Goals",
			wantName: "Purpose and Goals",
			wantOk:   true,
		},
		{
			name:     "delta section",
			input:    "## ADDED Requirements",
			wantName: "ADDED Requirements",
			wantOk:   true,
		},
		{
			name:     "section with trailing spaces",
			input:    "## Notes   ",
			wantName: "Notes",
			wantOk:   true,
		},
		{
			name:     "not a section - H3",
			input:    "### Requirement: Name",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a section - H1",
			input:    "# Title",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a section - plain text",
			input:    "This is regular text",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "empty line",
			input:    "",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "H2 without space",
			input:    "##NoSpace",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchH2SectionHeader(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchH2SectionHeader() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchH2SectionHeader() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchH2DeltaSection(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantDeltaType string
		wantOk        bool
	}{
		{
			name:          "ADDED",
			input:         "## ADDED Requirements",
			wantDeltaType: "ADDED",
			wantOk:        true,
		},
		{
			name:          "MODIFIED",
			input:         "## MODIFIED Requirements",
			wantDeltaType: "MODIFIED",
			wantOk:        true,
		},
		{
			name:          "REMOVED",
			input:         "## REMOVED Requirements",
			wantDeltaType: "REMOVED",
			wantOk:        true,
		},
		{
			name:          "RENAMED",
			input:         "## RENAMED Requirements",
			wantDeltaType: "RENAMED",
			wantOk:        true,
		},
		{
			name:          "with trailing space",
			input:         "## ADDED Requirements ",
			wantDeltaType: "ADDED",
			wantOk:        true,
		},
		{
			name:          "lowercase delta type",
			input:         "## added Requirements",
			wantDeltaType: "",
			wantOk:        false,
		},
		{
			name:          "plain Requirements section",
			input:         "## Requirements",
			wantDeltaType: "",
			wantOk:        false,
		},
		{
			name:          "invalid delta type",
			input:         "## UPDATED Requirements",
			wantDeltaType: "",
			wantOk:        false,
		},
		{
			name:          "missing Requirements suffix",
			input:         "## ADDED",
			wantDeltaType: "",
			wantOk:        false,
		},
		{
			name:          "extra content after Requirements",
			input:         "## ADDED Requirements extra",
			wantDeltaType: "",
			wantOk:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeltaType, gotOk := MatchH2DeltaSection(tt.input)
			if gotDeltaType != tt.wantDeltaType {
				t.Errorf(
					"MatchH2DeltaSection() deltaType = %q, want %q",
					gotDeltaType,
					tt.wantDeltaType,
				)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchH2DeltaSection() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchH3Requirement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "simple requirement",
			input:    "### Requirement: User Login",
			wantName: "User Login",
			wantOk:   true,
		},
		{
			name:     "requirement with special chars",
			input:    "### Requirement: API-v2.0 Endpoint",
			wantName: "API-v2.0 Endpoint",
			wantOk:   true,
		},
		{
			name:     "requirement no space after colon",
			input:    "### Requirement:NoSpace",
			wantName: "NoSpace",
			wantOk:   true,
		},
		{
			name:     "requirement with extra space after colon",
			input:    "### Requirement:  Double Space",
			wantName: "Double Space",
			wantOk:   true,
		},
		{
			name:     "not a requirement - scenario",
			input:    "#### Scenario: Test Case",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a requirement - H2",
			input:    "## Requirement: Not Valid",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a requirement - missing colon",
			input:    "### Requirement User Login",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a requirement - wrong keyword",
			input:    "### Spec: Some Spec",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "empty name",
			input:    "### Requirement:",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "requirement with only spaces after colon",
			input:    "### Requirement:   ",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchH3Requirement(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchH3Requirement() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchH3Requirement() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchH4Scenario(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "simple scenario",
			input:    "#### Scenario: Happy Path",
			wantName: "Happy Path",
			wantOk:   true,
		},
		{
			name:     "scenario with numbers",
			input:    "#### Scenario: Test Case 123",
			wantName: "Test Case 123",
			wantOk:   true,
		},
		{
			name:     "scenario no space after colon",
			input:    "#### Scenario:NoSpace",
			wantName: "NoSpace",
			wantOk:   true,
		},
		{
			name:     "not a scenario - requirement",
			input:    "### Requirement: Something",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a scenario - H5",
			input:    "##### Scenario: Too Deep",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a scenario - wrong keyword",
			input:    "#### Test: Something",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "empty name",
			input:    "#### Scenario:",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchH4Scenario(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchH4Scenario() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchH4Scenario() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
