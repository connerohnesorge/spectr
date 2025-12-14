package markdown

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
		name        string
		input       string
		wantID      string
		wantDesc    string
		wantChecked bool
		wantOk      bool
	}{
		{
			name:        "unchecked task",
			input:       "- [ ] 1.1 Create the regex package",
			wantID:      "1.1",
			wantDesc:    "Create the regex package",
			wantChecked: false,
			wantOk:      true,
		},
		{
			name:        "checked task",
			input:       "- [x] 2.3 Implement feature",
			wantID:      "2.3",
			wantDesc:    "Implement feature",
			wantChecked: true,
			wantOk:      true,
		},
		{
			name:        "checked uppercase",
			input:       "- [X] 10.15 Large task number",
			wantID:      "10.15",
			wantDesc:    "Large task number",
			wantChecked: true,
			wantOk:      true,
		},
		{
			name:        "task with special characters",
			input:       "- [ ] 1.2 Create `internal/regex/` package",
			wantID:      "1.2",
			wantDesc:    "Create `internal/regex/` package",
			wantChecked: false,
			wantOk:      true,
		},
		{
			name:     "missing checkbox",
			input:    "- 1.1 No checkbox",
			wantID:   "",
			wantDesc: "",
			wantOk:   false,
		},
		{
			name:     "missing task ID",
			input:    "- [ ] No ID here",
			wantID:   "",
			wantDesc: "",
			wantOk:   false,
		},
		{
			name:     "invalid task ID format",
			input:    "- [ ] 1 Single number",
			wantID:   "",
			wantDesc: "",
			wantOk:   false,
		},
		{
			name:     "plain list item",
			input:    "- Regular item",
			wantID:   "",
			wantDesc: "",
			wantOk:   false,
		},
		{
			name:     "section header",
			input:    "## 1. Core Section",
			wantID:   "",
			wantDesc: "",
			wantOk:   false,
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
			if match.ID != tt.wantID {
				t.Errorf("id = %q, want %q", match.ID, tt.wantID)
			}
			if match.Description != tt.wantDesc {
				t.Errorf("desc = %q, want %q", match.Description, tt.wantDesc)
			}
			if match.Checked != tt.wantChecked {
				t.Errorf("checked = %v, want %v", match.Checked, tt.wantChecked)
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

func TestCountTasksInContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		wantTotal     int
		wantCompleted int
	}{
		{
			name: "mixed tasks",
			content: `- [ ] Task 1
- [x] Task 2
- [ ] Task 3
- [X] Task 4`,
			wantTotal:     4,
			wantCompleted: 2,
		},
		{
			name:          "no tasks",
			content:       "Just some text\nand more text",
			wantTotal:     0,
			wantCompleted: 0,
		},
		{
			name: "all completed",
			content: `- [x] Done 1
- [X] Done 2`,
			wantTotal:     2,
			wantCompleted: 2,
		},
		{
			name: "all pending",
			content: `- [ ] Todo 1
- [ ] Todo 2`,
			wantTotal:     2,
			wantCompleted: 0,
		},
		{
			name:          "empty content",
			content:       "",
			wantTotal:     0,
			wantCompleted: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTotal, gotCompleted := CountTasksInContent(tt.content)
			if gotTotal != tt.wantTotal {
				t.Errorf("total = %d, want %d", gotTotal, tt.wantTotal)
			}
			if gotCompleted != tt.wantCompleted {
				t.Errorf("completed = %d, want %d", gotCompleted, tt.wantCompleted)
			}
		})
	}
}
