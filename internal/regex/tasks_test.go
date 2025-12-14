package regex

import (
	"testing"
)

func TestMatchTaskCheckbox(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState rune
		wantOk    bool
	}{
		{
			name:      "unchecked task",
			input:     "- [ ] Some task",
			wantState: ' ',
			wantOk:    true,
		},
		{
			name:      "checked task lowercase",
			input:     "- [x] Completed task",
			wantState: 'x',
			wantOk:    true,
		},
		{
			name:      "checked task uppercase",
			input:     "- [X] Completed task",
			wantState: 'X',
			wantOk:    true,
		},
		{
			name:      "task with leading spaces",
			input:     "  - [ ] Indented task",
			wantState: ' ',
			wantOk:    true,
		},
		{
			name:      "task with tabs",
			input:     "\t- [x] Tab indented",
			wantState: 'x',
			wantOk:    true,
		},
		{
			name:      "no space after dash",
			input:     "-[ ] No space",
			wantState: ' ',
			wantOk:    true,
		},
		{
			name:      "not a task - regular list",
			input:     "- Regular list item",
			wantState: 0,
			wantOk:    false,
		},
		{
			name:      "not a task - numbered list",
			input:     "1. Numbered item",
			wantState: 0,
			wantOk:    false,
		},
		{
			name:      "not a task - plain text",
			input:     "Some text",
			wantState: 0,
			wantOk:    false,
		},
		{
			name:      "empty line",
			input:     "",
			wantState: 0,
			wantOk:    false,
		},
		{
			name:      "invalid checkbox state",
			input:     "- [?] Invalid",
			wantState: 0,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotState, gotOk := MatchTaskCheckbox(tt.input)
			if gotState != tt.wantState {
				t.Errorf("MatchTaskCheckbox() state = %q, want %q", gotState, tt.wantState)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchTaskCheckbox() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMatchNumberedTask(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantCheckbox string
		wantID       string
		wantDesc     string
		wantOk       bool
	}{
		{
			name:         "unchecked task",
			input:        "- [ ] 1.1 Create the regex package",
			wantCheckbox: " ",
			wantID:       "1.1",
			wantDesc:     "Create the regex package",
			wantOk:       true,
		},
		{
			name:         "checked task",
			input:        "- [x] 2.3 Implement feature",
			wantCheckbox: "x",
			wantID:       "2.3",
			wantDesc:     "Implement feature",
			wantOk:       true,
		},
		{
			name:         "checked uppercase",
			input:        "- [X] 10.15 Large task number",
			wantCheckbox: "X",
			wantID:       "10.15",
			wantDesc:     "Large task number",
			wantOk:       true,
		},
		{
			name:         "task with special characters",
			input:        "- [ ] 1.2 Create `internal/regex/` package",
			wantCheckbox: " ",
			wantID:       "1.2",
			wantDesc:     "Create `internal/regex/` package",
			wantOk:       true,
		},
		{
			name:         "missing checkbox",
			input:        "- 1.1 No checkbox",
			wantCheckbox: "",
			wantID:       "",
			wantDesc:     "",
			wantOk:       false,
		},
		{
			name:         "missing task ID",
			input:        "- [ ] No ID here",
			wantCheckbox: "",
			wantID:       "",
			wantDesc:     "",
			wantOk:       false,
		},
		{
			name:         "invalid task ID format",
			input:        "- [ ] 1 Single number",
			wantCheckbox: "",
			wantID:       "",
			wantDesc:     "",
			wantOk:       false,
		},
		{
			name:         "plain list item",
			input:        "- Regular item",
			wantCheckbox: "",
			wantID:       "",
			wantDesc:     "",
			wantOk:       false,
		},
		{
			name:         "section header",
			input:        "## 1. Core Section",
			wantCheckbox: "",
			wantID:       "",
			wantDesc:     "",
			wantOk:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, gotOk := MatchNumberedTask(tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("MatchNumberedTask() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if !gotOk {
				if match != nil {
					t.Error("MatchNumberedTask() match should be nil when ok=false")
				}

				return
			}
			if match.Checkbox != tt.wantCheckbox {
				t.Errorf("checkbox = %q, want %q", match.Checkbox, tt.wantCheckbox)
			}
			if match.ID != tt.wantID {
				t.Errorf("id = %q, want %q", match.ID, tt.wantID)
			}
			if match.Description != tt.wantDesc {
				t.Errorf("desc = %q, want %q", match.Description, tt.wantDesc)
			}
		})
	}
}

func TestMatchNumberedSection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantOk   bool
	}{
		{
			name:     "simple section",
			input:    "## 1. Core Accept Command",
			wantName: "Core Accept Command",
			wantOk:   true,
		},
		{
			name:     "double digit section",
			input:    "## 10. Large Section Number",
			wantName: "Large Section Number",
			wantOk:   true,
		},
		{
			name:     "section with special chars",
			input:    "## 2. API & Integration",
			wantName: "API & Integration",
			wantOk:   true,
		},
		{
			name:     "not numbered section",
			input:    "## Requirements",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "H3 numbered",
			input:    "### 1. Not H2",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "missing period after number",
			input:    "## 1 No Period",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "task line",
			input:    "- [ ] 1.1 Task",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchNumberedSection(tt.input)
			if gotName != tt.wantName {
				t.Errorf("MatchNumberedSection() name = %q, want %q", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("MatchNumberedSection() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestIsTaskChecked(t *testing.T) {
	tests := []struct {
		name  string
		state rune
		want  bool
	}{
		{"lowercase x", 'x', true},
		{"uppercase X", 'X', true},
		{"space", ' ', false},
		{"other char", '?', false},
		{"zero rune", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTaskChecked(tt.state); got != tt.want {
				t.Errorf("IsTaskChecked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskCheckboxPattern(t *testing.T) {
	// Test the raw pattern behavior
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"matches unchecked", "- [ ] task", true},
		{"matches checked", "- [x] task", true},
		{"matches at start", "- [X] task", true},
		{"matches with indent", "   - [ ] task", true},
		{"no match for regular list", "- item", false},
		{"no match for wrong bracket", "- ( ) item", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TaskCheckbox.MatchString(tt.input)
			if got != tt.want {
				t.Errorf("TaskCheckbox.MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNumberedTaskPattern(t *testing.T) {
	// Test the raw pattern behavior
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"matches standard format", "- [ ] 1.1 Description", true},
		{"matches larger numbers", "- [x] 12.34 Description", true},
		{"no match without ID", "- [ ] Description", false},
		{"no match without checkbox", "- 1.1 Description", false},
		{"no match for section", "## 1. Section", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NumberedTask.MatchString(tt.input)
			if got != tt.want {
				t.Errorf("NumberedTask.MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
