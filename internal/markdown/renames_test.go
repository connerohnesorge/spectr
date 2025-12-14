package markdown

import (
	"testing"
)

func TestMatchRenamedFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "standard format",
			input:    "- FROM: `### Requirement: Old Name`",
			wantName: "Old Name",
			wantOk:   true,
		},
		{
			name:     "name with special chars",
			input:    "- FROM: `### Requirement: API-v2.0 Endpoint`",
			wantName: "API-v2.0 Endpoint",
			wantOk:   true,
		},
		{
			name:     "without backticks",
			input:    "- FROM: ### Requirement: Name",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "TO instead of FROM",
			input:    "- TO: `### Requirement: Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "missing backticks",
			input:    "- FROM: ### Requirement: Name",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "plain text",
			input:    "Some text",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "empty line",
			input:    "",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchRenamedFrom(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchRenamedFrom() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchRenamedFrom() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchRenamedFromAlt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "standard format",
			input:    "- FROM: ### Requirement: Old Name",
			wantName: "Old Name",
			wantOk:   true,
		},
		{
			name:     "with leading space",
			input:    "  - FROM: ### Requirement: Name",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "no space after Requirement colon",
			input:    "- FROM: ### Requirement:Name",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "with backticks (wrong format for alt)",
			input:    "- FROM: `### Requirement: Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "missing Requirement keyword",
			input:    "- FROM: ### Something Else",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "TO line",
			input:    "- TO: ### Requirement: Name",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchRenamedFromAlt(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchRenamedFromAlt() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchRenamedFromAlt() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchRenamedTo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "standard format",
			input:    "- TO: `### Requirement: New Name`",
			wantName: "New Name",
			wantOk:   true,
		},
		{
			name:     "name with special chars",
			input:    "- TO: `### Requirement: API-v3.0 Endpoint`",
			wantName: "API-v3.0 Endpoint",
			wantOk:   true,
		},
		{
			name:     "FROM instead of TO",
			input:    "- FROM: `### Requirement: Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "without backticks",
			input:    "- TO: ### Requirement: Name",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchRenamedTo(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchRenamedTo() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchRenamedTo() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchRenamedToAlt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "standard format",
			input:    "- TO: ### Requirement: New Name",
			wantName: "New Name",
			wantOk:   true,
		},
		{
			name:     "with leading space",
			input:    "  - TO: ### Requirement: Name",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "no space after Requirement colon",
			input:    "- TO: ### Requirement:Name",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "with backticks (wrong format for alt)",
			input:    "- TO: `### Requirement: Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "FROM line",
			input:    "- FROM: ### Requirement: Name",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchRenamedToAlt(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchRenamedToAlt() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchRenamedToAlt() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchAnyRenamedFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "backtick format",
			input:    "- FROM: `### Requirement: Name`",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "non-backtick format",
			input:    "- FROM: ### Requirement: Name",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "TO line",
			input:    "- TO: `### Requirement: Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "invalid format",
			input:    "- FROM: Something else",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchAnyRenamedFrom(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchAnyRenamedFrom() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchAnyRenamedFrom() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchAnyRenamedTo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "backtick format",
			input:    "- TO: `### Requirement: Name`",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "non-backtick format",
			input:    "- TO: ### Requirement: Name",
			wantName: "Name",
			wantOk:   true,
		},
		{
			name:     "FROM line",
			input:    "- FROM: `### Requirement: Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "invalid format",
			input:    "- TO: Something else",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchAnyRenamedTo(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchAnyRenamedTo() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchAnyRenamedTo() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
