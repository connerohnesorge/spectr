package markdown

import (
	"testing"
)

func TestMatchRequirementHeader(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantName string
		wantOk   bool
	}{
		{
			name:     "valid requirement header",
			line:     "### Requirement: User Authentication",
			wantName: "User Authentication",
			wantOk:   true,
		},
		{
			name:     "requirement with extra spaces",
			line:     "### Requirement:   Multiple Spaces",
			wantName: "Multiple Spaces",
			wantOk:   true,
		},
		{
			name:     "requirement with tab",
			line:     "### Requirement:\tTab After",
			wantName: "Tab After",
			wantOk:   true,
		},
		{
			name:     "requirement preserves trailing spaces",
			line:     "### Requirement: Name With Trailing  ",
			wantName: "Name With Trailing  ",
			wantOk:   true,
		},
		{
			name:     "not a requirement header - wrong prefix",
			line:     "## Requirement: Wrong Level",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a requirement header - missing colon",
			line:     "### Requirement User Auth",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a requirement header - empty name",
			line:     "### Requirement:",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a requirement header - only spaces after",
			line:     "### Requirement:   ",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a requirement header - plain text",
			line:     "This is not a header",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "h4 requirement not matched",
			line:     "#### Requirement: Too Deep",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchRequirementHeader(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchRequirementHeader(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

func TestMatchScenarioHeader(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantName string
		wantOk   bool
	}{
		{
			name:     "valid scenario header",
			line:     "#### Scenario: User logs in",
			wantName: "User logs in",
			wantOk:   true,
		},
		{
			name:     "scenario with extra spaces",
			line:     "#### Scenario:   Multiple Spaces",
			wantName: "Multiple Spaces",
			wantOk:   true,
		},
		{
			name:     "not a scenario - wrong level",
			line:     "### Scenario: Wrong Level",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a scenario - missing colon",
			line:     "#### Scenario User logs in",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a scenario - empty name",
			line:     "#### Scenario:",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "plain text",
			line:     "Not a scenario header",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchScenarioHeader(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchScenarioHeader(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

func TestIsH2Header(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"## Title", true},
		{"## ", true},
		{"##Title", false},
		{"## Long Title With Spaces", true},
		{"### Not H2", false},
		{"# H1", false},
		{"Plain text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsH2Header(tt.line); got != tt.want {
				t.Errorf(
					"IsH2Header(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIsH3Header(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"### Title", true},
		{"### ", true},
		{"###Title", false},
		{"## Not H3", false},
		{"#### H4", false},
		{"Plain text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsH3Header(tt.line); got != tt.want {
				t.Errorf(
					"IsH3Header(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIsH4Header(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"#### Title", true},
		{"#### ", true},
		{"####Title", false},
		{"### Not H4", false},
		{"##### H5", false},
		{"Plain text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsH4Header(tt.line); got != tt.want {
				t.Errorf(
					"IsH4Header(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestMatchTaskCheckbox(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantStatus rune
		wantOk     bool
	}{
		{
			name:       "unchecked task",
			line:       "- [ ] Unchecked task",
			wantStatus: ' ',
			wantOk:     true,
		},
		{
			name:       "checked task lowercase",
			line:       "- [x] Checked task",
			wantStatus: 'x',
			wantOk:     true,
		},
		{
			name:       "checked task uppercase",
			line:       "- [X] Checked task",
			wantStatus: 'X',
			wantOk:     true,
		},
		{
			name:       "task with leading spaces",
			line:       "  - [ ] Indented task",
			wantStatus: ' ',
			wantOk:     true,
		},
		{
			name:       "task with leading tab",
			line:       "\t- [x] Tab indented",
			wantStatus: 'x',
			wantOk:     true,
		},
		{
			name:       "not a task - wrong bracket",
			line:       "- (x) Wrong brackets",
			wantStatus: 0,
			wantOk:     false,
		},
		{
			name:       "not a task - missing dash",
			line:       "[ ] Missing dash",
			wantStatus: 0,
			wantOk:     false,
		},
		{
			name:       "not a task - invalid state",
			line:       "- [?] Invalid state",
			wantStatus: 0,
			wantOk:     false,
		},
		{
			name:       "not a task - plain text",
			line:       "Plain text",
			wantStatus: 0,
			wantOk:     false,
		},
		{
			name:       "empty line",
			line:       "",
			wantStatus: 0,
			wantOk:     false,
		},
		{
			name:       "too short",
			line:       "- [",
			wantStatus: 0,
			wantOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotOk := MatchTaskCheckbox(
				tt.line,
			)
			if gotStatus != tt.wantStatus ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchTaskCheckbox(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotStatus,
					gotOk,
					tt.wantStatus,
					tt.wantOk,
				)
			}
		})
	}
}

func TestMatchNumberedTask(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantOk  bool
		wantNum string
		wantSec string
		wantSt  rune
		wantCon string
	}{
		{
			name:    "valid numbered task unchecked",
			line:    "- [ ] 1.1 Create the parser",
			wantOk:  true,
			wantNum: "1.1",
			wantSec: "1",
			wantSt:  ' ',
			wantCon: "Create the parser",
		},
		{
			name:    "valid numbered task checked",
			line:    "- [x] 2.3 Implement feature",
			wantOk:  true,
			wantNum: "2.3",
			wantSec: "2",
			wantSt:  'x',
			wantCon: "Implement feature",
		},
		{
			name:    "double digit numbers",
			line:    "- [ ] 12.34 Complex task",
			wantOk:  true,
			wantNum: "12.34",
			wantSec: "12",
			wantSt:  ' ',
			wantCon: "Complex task",
		},
		{
			name:    "uppercase X",
			line:    "- [X] 1.2 Done task",
			wantOk:  true,
			wantNum: "1.2",
			wantSec: "1",
			wantSt:  'X',
			wantCon: "Done task",
		},
		{
			name:   "no number",
			line:   "- [ ] No number here",
			wantOk: false,
		},
		{
			name:   "missing dot in number",
			line:   "- [ ] 123 No dot",
			wantOk: false,
		},
		{
			name:   "wrong format",
			line:   "* [ ] 1.1 Wrong bullet",
			wantOk: false,
		},
		{
			name:   "plain text",
			line:   "Plain text",
			wantOk: false,
		},
		{
			name:   "empty line",
			line:   "",
			wantOk: false,
		},
		{
			name:   "missing content",
			line:   "- [ ] 1.1 ",
			wantOk: false,
		},
		{
			name:   "no space after number",
			line:   "- [ ] 1.1Task",
			wantOk: false,
		},
		// Tests for N. format (without digits after dot)
		{
			name:    "simple format unchecked",
			line:    "- [ ] 1. Create parser",
			wantOk:  true,
			wantNum: "1.",
			wantSec: "1",
			wantSt:  ' ',
			wantCon: "Create parser",
		},
		{
			name:    "simple format checked",
			line:    "- [x] 2. Implement feature",
			wantOk:  true,
			wantNum: "2.",
			wantSec: "2",
			wantSt:  'x',
			wantCon: "Implement feature",
		},
		{
			name:    "simple format double digit",
			line:    "- [ ] 12. Complex task",
			wantOk:  true,
			wantNum: "12.",
			wantSec: "12",
			wantSt:  ' ',
			wantCon: "Complex task",
		},
		{
			name:   "simple format no space after dot",
			line:   "- [ ] 1.Task",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := MatchNumberedTask(
				tt.line,
			)
			if gotOk != tt.wantOk {
				t.Errorf(
					"MatchNumberedTask(%q) ok = %v, want %v",
					tt.line,
					gotOk,
					tt.wantOk,
				)

				return
			}
			if !gotOk {
				return
			}
			if got.Number != tt.wantNum {
				t.Errorf(
					"Number = %q, want %q",
					got.Number,
					tt.wantNum,
				)
			}
			if got.Section != tt.wantSec {
				t.Errorf(
					"Section = %q, want %q",
					got.Section,
					tt.wantSec,
				)
			}
			if got.Status != tt.wantSt {
				t.Errorf(
					"Status = %q, want %q",
					got.Status,
					tt.wantSt,
				)
			}
			if got.Content != tt.wantCon {
				t.Errorf(
					"Content = %q, want %q",
					got.Content,
					tt.wantCon,
				)
			}
		})
	}
}

func TestMatchFlexibleTask(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantMatch   bool
		wantNumber  string
		wantStatus  rune
		wantContent string
	}{
		// Decimal format (1.1, 1.2, etc.)
		{
			name:        "decimal format unchecked",
			line:        "- [ ] 1.1 Create the parser",
			wantMatch:   true,
			wantNumber:  "1.1",
			wantStatus:  ' ',
			wantContent: "Create the parser",
		},
		{
			name:        "decimal format checked",
			line:        "- [x] 2.3 Implement feature",
			wantMatch:   true,
			wantNumber:  "2.3",
			wantStatus:  'x',
			wantContent: "Implement feature",
		},
		// Simple dot format (1., 2., etc.)
		{
			name:        "simple dot format",
			line:        "- [ ] 1. Create something",
			wantMatch:   true,
			wantNumber:  "1.",
			wantStatus:  ' ',
			wantContent: "Create something",
		},
		{
			name:        "simple dot format checked",
			line:        "- [X] 5. Complete task",
			wantMatch:   true,
			wantNumber:  "5.",
			wantStatus:  'X',
			wantContent: "Complete task",
		},
		// Number only format (1, 2, etc.)
		{
			name:        "number only format",
			line:        "- [ ] 1 Create item",
			wantMatch:   true,
			wantNumber:  "1",
			wantStatus:  ' ',
			wantContent: "Create item",
		},
		{
			name:        "number only format double digit",
			line:        "- [x] 12 Another task",
			wantMatch:   true,
			wantNumber:  "12",
			wantStatus:  'x',
			wantContent: "Another task",
		},
		// No number format
		{
			name:        "no number format",
			line:        "- [ ] Create the parser",
			wantMatch:   true,
			wantNumber:  "",
			wantStatus:  ' ',
			wantContent: "Create the parser",
		},
		{
			name:        "no number format checked",
			line:        "- [x] Implement feature",
			wantMatch:   true,
			wantNumber:  "",
			wantStatus:  'x',
			wantContent: "Implement feature",
		},
		// Edge cases
		{
			name:        "complex decimal",
			line:        "- [ ] 12.34 Multi-digit numbers",
			wantMatch:   true,
			wantNumber:  "12.34",
			wantStatus:  ' ',
			wantContent: "Multi-digit numbers",
		},
		{
			name:        "content with backticks",
			line:        "- [ ] Update `cmd/validate.go` to fix",
			wantMatch:   true,
			wantNumber:  "",
			wantStatus:  ' ',
			wantContent: "Update `cmd/validate.go` to fix",
		},
		{
			name:        "numbered with backticks",
			line:        "- [ ] 1.1 Update `cmd/validate.go`",
			wantMatch:   true,
			wantNumber:  "1.1",
			wantStatus:  ' ',
			wantContent: "Update `cmd/validate.go`",
		},
		// Invalid inputs
		{
			name:      "not a task",
			line:      "## Section header",
			wantMatch: false,
		},
		{
			name:      "plain text",
			line:      "Just some text",
			wantMatch: false,
		},
		{
			name:      "empty line",
			line:      "",
			wantMatch: false,
		},
		{
			name:      "checkbox without content",
			line:      "- [ ] ",
			wantMatch: false,
		},
		{
			name:      "invalid checkbox",
			line:      "- [?] Task",
			wantMatch: false,
		},
		{
			name:      "wrong bullet format",
			line:      "* [ ] Task",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, ok := MatchFlexibleTask(tt.line)
			if ok != tt.wantMatch {
				t.Errorf(
					"MatchFlexibleTask() matched = %v, want %v",
					ok,
					tt.wantMatch,
				)

				return
			}
			if !ok {
				return
			}
			if match.Number != tt.wantNumber {
				t.Errorf(
					"Number = %q, want %q",
					match.Number,
					tt.wantNumber,
				)
			}
			if match.Status != tt.wantStatus {
				t.Errorf(
					"Status = %q, want %q",
					match.Status,
					tt.wantStatus,
				)
			}
			if match.Content != tt.wantContent {
				t.Errorf(
					"Content = %q, want %q",
					match.Content,
					tt.wantContent,
				)
			}
		})
	}
}

func TestMatchNumberedSection(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantName string
		wantOk   bool
	}{
		{
			name:     "valid numbered section",
			line:     "## 1. Core Accept Command",
			wantName: "Core Accept Command",
			wantOk:   true,
		},
		{
			name:     "double digit section",
			line:     "## 12. Extended Features",
			wantName: "Extended Features",
			wantOk:   true,
		},
		{
			name:     "not numbered section - no number",
			line:     "## Plain Section",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not numbered section - wrong level",
			line:     "### 1. Wrong Level",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not numbered section - no space after dot",
			line:     "## 1.Section",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not numbered section - no content",
			line:     "## 1. ",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "plain text",
			line:     "Plain text",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchNumberedSection(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchNumberedSection(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

func TestMatchAnySection(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantName   string
		wantNumber string
		wantOk     bool
	}{
		{
			name:       "numbered section",
			line:       "## 1. Setup",
			wantName:   "Setup",
			wantNumber: "1",
			wantOk:     true,
		},
		{
			name:       "double digit numbered section",
			line:       "## 12. Advanced Setup",
			wantName:   "Advanced Setup",
			wantNumber: "12",
			wantOk:     true,
		},
		{
			name:       "unnumbered section",
			line:       "## Implementation",
			wantName:   "Implementation",
			wantNumber: "",
			wantOk:     true,
		},
		{
			name:       "unnumbered section with spaces",
			line:       "## Long Section Name",
			wantName:   "Long Section Name",
			wantNumber: "",
			wantOk:     true,
		},
		{
			name:       "not H2 - wrong level",
			line:       "### 1. Not H2",
			wantName:   "",
			wantNumber: "",
			wantOk:     false,
		},
		{
			name:       "empty after ##",
			line:       "## ",
			wantName:   "",
			wantNumber: "",
			wantOk:     false,
		},
		{
			name:       "plain text",
			line:       "Plain text",
			wantName:   "",
			wantNumber: "",
			wantOk:     false,
		},
		{
			name:       "empty line",
			line:       "",
			wantName:   "",
			wantNumber: "",
			wantOk:     false,
		},
		{
			name:       "numbered but no content",
			line:       "## 1. ",
			wantName:   "",
			wantNumber: "",
			wantOk:     false,
		},
		{
			name:       "number without dot treated as unnumbered",
			line:       "## 1 Setup",
			wantName:   "1 Setup",
			wantNumber: "",
			wantOk:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotNumber, gotOk := MatchAnySection(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotNumber != tt.wantNumber ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchAnySection(%q) = (%q, %q, %v), want (%q, %q, %v)",
					tt.line,
					gotName,
					gotNumber,
					gotOk,
					tt.wantName,
					tt.wantNumber,
					tt.wantOk,
				)
			}
		})
	}
}

func TestExtractHeaderLevel(t *testing.T) {
	tests := []struct {
		line string
		want int
	}{
		{"# H1", 1},
		{"## H2", 2},
		{"### H3", 3},
		{"#### H4", 4},
		{"##### H5", 5},
		{"###### H6", 6},
		{
			"####### Too Many",
			0,
		}, // 7 hashes not valid
		{"#NoSpace", 0},
		{"Plain text", 0},
		{"", 0},
		{"   # Indented", 0},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := ExtractHeaderLevel(tt.line); got != tt.want {
				t.Errorf(
					"ExtractHeaderLevel(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestExtractHeaderText(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{"# Title", "Title"},
		{"## My Section", "My Section"},
		{"### Subsection  ", "Subsection"},
		{"#NoSpace", ""},
		{"Plain text", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := ExtractHeaderText(tt.line); got != tt.want {
				t.Errorf(
					"ExtractHeaderText(%q) = %q, want %q",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIsBlankLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t", true},
		{"\t  \t", true},
		{"text", false},
		{" text ", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsBlankLine(tt.line); got != tt.want {
				t.Errorf(
					"IsBlankLine(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIsListItem(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"- Item", true},
		{"* Item", true},
		{"+ Item", true},
		{"1. Item", true},
		{"23. Item", true},
		{"  - Indented", true},
		{"\t* Tab indented", true},
		{"-NoSpace", false},
		{"1NoSpace", false},
		{"Plain text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsListItem(tt.line); got != tt.want {
				t.Errorf(
					"IsListItem(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIsCodeFence(t *testing.T) {
	tests := []struct {
		line      string
		wantFence bool
		wantChar  rune
	}{
		{"```", true, '`'},
		{"```go", true, '`'},
		{"~~~", true, '~'},
		{"~~~python", true, '~'},
		{"  ```", true, '`'},
		{"\t~~~", true, '~'},
		{"``", false, 0},
		{"~~", false, 0},
		{"Plain text", false, 0},
		{"", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotFence, gotChar := IsCodeFence(
				tt.line,
			)
			if gotFence != tt.wantFence ||
				gotChar != tt.wantChar {
				t.Errorf(
					"IsCodeFence(%q) = (%v, %q), want (%v, %q)",
					tt.line,
					gotFence,
					gotChar,
					tt.wantFence,
					tt.wantChar,
				)
			}
		})
	}
}

func TestIsBlockquote(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"> Quote", true},
		{">Quote", true},
		{"  > Indented", true},
		{"\t> Tab", true},
		{"Plain text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsBlockquote(tt.line); got != tt.want {
				t.Errorf(
					"IsBlockquote(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestMatchH2SectionHeader(t *testing.T) {
	tests := []struct {
		line     string
		wantName string
		wantOk   bool
	}{
		{"## Purpose", "Purpose", true},
		{"## Requirements", "Requirements", true},
		{
			"## Long Section Name",
			"Long Section Name",
			true,
		},
		{"## ", "", false},
		{"### Not H2", "", false},
		{"##NoSpace", "", false},
		{"Plain", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotName, gotOk := MatchH2SectionHeader(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchH2SectionHeader(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

func TestMatchH2DeltaSection(t *testing.T) {
	tests := []struct {
		line     string
		wantType string
		wantOk   bool
	}{
		{"## ADDED Requirements", "ADDED", true},
		{
			"## MODIFIED Requirements",
			"MODIFIED",
			true,
		},
		{
			"## REMOVED Requirements",
			"REMOVED",
			true,
		},
		{
			"## RENAMED Requirements",
			"RENAMED",
			true,
		},
		{
			"## ADDED Requirements  ",
			"ADDED",
			true,
		}, // trailing spaces trimmed
		{
			"## Added Requirements",
			"",
			false,
		}, // case sensitive
		{
			"## ADDED Reqs",
			"",
			false,
		}, // must be "Requirements"
		{
			"### ADDED Requirements",
			"",
			false,
		}, // wrong level
		{
			"## Requirements",
			"",
			false,
		}, // no delta type
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotType, gotOk := MatchH2DeltaSection(
				tt.line,
			)
			if gotType != tt.wantType ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchH2DeltaSection(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotType,
					gotOk,
					tt.wantType,
					tt.wantOk,
				)
			}
		})
	}
}

func TestIsTaskChecked(t *testing.T) {
	tests := []struct {
		state rune
		want  bool
	}{
		{'x', true},
		{'X', true},
		{' ', false},
		{0, false},
		{'?', false},
	}

	for _, tt := range tests {
		t.Run(
			string(tt.state),
			func(t *testing.T) {
				if got := IsTaskChecked(tt.state); got != tt.want {
					t.Errorf(
						"IsTaskChecked(%q) = %v, want %v",
						tt.state,
						got,
						tt.want,
					)
				}
			},
		)
	}
}

func TestExtractListMarker(t *testing.T) {
	tests := []struct {
		line       string
		wantMarker string
		wantOk     bool
	}{
		{"- Item", "-", true},
		{"* Item", "*", true},
		{"+ Item", "+", true},
		{"1. Item", "1.", true},
		{"23. Item", "23.", true},
		{"  - Indented", "-", true},
		{"-NoSpace", "", false},
		{"1NoSpace", "", false},
		{"Plain text", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotMarker, gotOk := ExtractListMarker(
				tt.line,
			)
			if gotMarker != tt.wantMarker ||
				gotOk != tt.wantOk {
				t.Errorf(
					"ExtractListMarker(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotMarker,
					gotOk,
					tt.wantMarker,
					tt.wantOk,
				)
			}
		})
	}
}

func TestCountLeadingSpaces(t *testing.T) {
	tests := []struct {
		line string
		want int
	}{
		{"text", 0},
		{" text", 1},
		{"  text", 2},
		{"\ttext", 1},
		{"  \ttext", 3},
		{"", 0},
		{"   ", 3},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := CountLeadingSpaces(tt.line); got != tt.want {
				t.Errorf(
					"CountLeadingSpaces(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestTrimLeadingHashes(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{"# Title", "Title"},
		{"## Section", "Section"},
		{"### Sub", "Sub"},
		{"#NoSpace", "NoSpace"},
		{"Plain text", "Plain text"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := TrimLeadingHashes(tt.line); got != tt.want {
				t.Errorf(
					"TrimLeadingHashes(%q) = %q, want %q",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestContainsKeyword(t *testing.T) {
	tests := []struct {
		line        string
		wantKeyword string
		wantOk      bool
	}{
		{"- **WHEN** user clicks", "WHEN", true},
		{"- **THEN** result shows", "THEN", true},
		{"- **AND** another thing", "AND", true},
		{
			"- **GIVEN** initial state",
			"GIVEN",
			true,
		},
		{"- WHEN no bold", "", false},
		{"- Plain text", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotKeyword, gotOk := ContainsKeyword(
				tt.line,
			)
			if gotKeyword != tt.wantKeyword ||
				gotOk != tt.wantOk {
				t.Errorf(
					"ContainsKeyword(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotKeyword,
					gotOk,
					tt.wantKeyword,
					tt.wantOk,
				)
			}
		})
	}
}

func TestIsHorizontalRule(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"---", true},
		{"***", true},
		{"___", true},
		{"----", true},
		{"- - -", true},
		{"* * *", true},
		{"_ _ _", true},
		{"--", false},
		{"**", false},
		{"__", false},
		{"-*-", false},
		{"text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsHorizontalRule(tt.line); got != tt.want {
				t.Errorf(
					"IsHorizontalRule(%q) = %v, want %v",
					tt.line,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestMatchRenamedFromAlt(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantName string
		wantOk   bool
	}{
		{
			name:     "valid FROM without backticks",
			line:     "- FROM: ### Requirement: Old Name",
			wantName: "Old Name",
			wantOk:   true,
		},
		{
			name:     "FROM with leading spaces",
			line:     "  - FROM: ### Requirement: Old Name",
			wantName: "Old Name",
			wantOk:   true,
		},
		{
			name:     "FROM case insensitive",
			line:     "- from: ### Requirement: Old Name",
			wantName: "Old Name",
			wantOk:   true,
		},
		{
			name:     "FROM with extra space after colon",
			line:     "- FROM:   ### Requirement: Name Here",
			wantName: "Name Here",
			wantOk:   true,
		},
		{
			name:     "not a FROM - backtick wrapped",
			line:     "- FROM: `### Requirement: Old Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a FROM - wrong prefix",
			line:     "FROM: ### Requirement: Old Name",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a FROM - missing requirement",
			line:     "- FROM: ### Something Else",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a FROM - empty name",
			line:     "- FROM: ### Requirement:",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "plain text",
			line:     "Plain text",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchRenamedFromAlt(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchRenamedFromAlt(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

func TestMatchRenamedToAlt(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantName string
		wantOk   bool
	}{
		{
			name:     "valid TO without backticks",
			line:     "- TO: ### Requirement: New Name",
			wantName: "New Name",
			wantOk:   true,
		},
		{
			name:     "TO with leading spaces",
			line:     "  - TO: ### Requirement: New Name",
			wantName: "New Name",
			wantOk:   true,
		},
		{
			name:     "TO case insensitive",
			line:     "- to: ### Requirement: New Name",
			wantName: "New Name",
			wantOk:   true,
		},
		{
			name:     "TO with extra space after colon",
			line:     "- TO:   ### Requirement: Name Here",
			wantName: "Name Here",
			wantOk:   true,
		},
		{
			name:     "not a TO - backtick wrapped",
			line:     "- TO: `### Requirement: New Name`",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a TO - wrong prefix",
			line:     "TO: ### Requirement: New Name",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a TO - missing requirement",
			line:     "- TO: ### Something Else",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "not a TO - empty name",
			line:     "- TO: ### Requirement:",
			wantName: "",
			wantOk:   false,
		},
		{
			name:     "plain text",
			line:     "Plain text",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchRenamedToAlt(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchRenamedToAlt(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

func TestMatchAnyRenamedFrom(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantName string
		wantOk   bool
	}{
		{
			name:     "backtick format",
			line:     "- FROM: `### Requirement: Old Name`",
			wantName: "Old Name",
			wantOk:   true,
		},
		{
			name:     "non-backtick format",
			line:     "- FROM: ### Requirement: Old Name",
			wantName: "Old Name",
			wantOk:   true,
		},
		{
			name:     "neither format",
			line:     "- FROM: Something else",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchAnyRenamedFrom(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchAnyRenamedFrom(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

func TestMatchAnyRenamedTo(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantName string
		wantOk   bool
	}{
		{
			name:     "backtick format",
			line:     "- TO: `### Requirement: New Name`",
			wantName: "New Name",
			wantOk:   true,
		},
		{
			name:     "non-backtick format",
			line:     "- TO: ### Requirement: New Name",
			wantName: "New Name",
			wantOk:   true,
		},
		{
			name:     "neither format",
			line:     "- TO: Something else",
			wantName: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := MatchAnyRenamedTo(
				tt.line,
			)
			if gotName != tt.wantName ||
				gotOk != tt.wantOk {
				t.Errorf(
					"MatchAnyRenamedTo(%q) = (%q, %v), want (%q, %v)",
					tt.line,
					gotName,
					gotOk,
					tt.wantName,
					tt.wantOk,
				)
			}
		})
	}
}

// TestCompatibilityWithRegexPackage tests that the compat functions
// produce equivalent results to the old regex package for common cases.
func TestCompatibilityWithRegexPackage(
	t *testing.T,
) {
	t.Run(
		"requirement header compatibility",
		func(t *testing.T) {
			// These are the patterns the old regex would match
			testCases := []string{
				"### Requirement: User Authentication",
				"### Requirement: Data Validation",
				"### Requirement:  Extra Space",
			}

			for _, tc := range testCases {
				name, ok := MatchRequirementHeader(
					tc,
				)
				if !ok {
					t.Errorf(
						"MatchRequirementHeader(%q) should match",
						tc,
					)

					continue
				}
				if name == "" {
					t.Errorf(
						"MatchRequirementHeader(%q) should extract name",
						tc,
					)
				}
			}
		},
	)

	t.Run(
		"scenario header compatibility",
		func(t *testing.T) {
			testCases := []string{
				"#### Scenario: User logs in successfully",
				"#### Scenario: Invalid credentials",
			}

			for _, tc := range testCases {
				name, ok := MatchScenarioHeader(
					tc,
				)
				if !ok {
					t.Errorf(
						"MatchScenarioHeader(%q) should match",
						tc,
					)

					continue
				}
				if name == "" {
					t.Errorf(
						"MatchScenarioHeader(%q) should extract name",
						tc,
					)
				}
			}
		},
	)

	t.Run(
		"task checkbox compatibility",
		func(t *testing.T) {
			testCases := []struct {
				line    string
				checked bool
			}{
				{"- [ ] Unchecked", false},
				{"- [x] Checked", true},
				{"- [X] Also checked", true},
				{"  - [ ] Indented", false},
			}

			for _, tc := range testCases {
				status, ok := MatchTaskCheckbox(
					tc.line,
				)
				if !ok {
					t.Errorf(
						"MatchTaskCheckbox(%q) should match",
						tc.line,
					)

					continue
				}
				if IsTaskChecked(
					status,
				) != tc.checked {
					t.Errorf(
						"MatchTaskCheckbox(%q) checked = %v, want %v",
						tc.line,
						IsTaskChecked(status),
						tc.checked,
					)
				}
			}
		},
	)

	t.Run(
		"numbered task compatibility",
		func(t *testing.T) {
			testCases := []struct {
				line   string
				wantID string
				wantOk bool
			}{
				{
					"- [ ] 1.1 Create the regex package",
					"1.1",
					true,
				},
				{
					"- [x] 2.3 Done task",
					"2.3",
					true,
				},
				{
					"- [ ] Plain task without number",
					"",
					false,
				},
			}

			for _, tc := range testCases {
				match, ok := MatchNumberedTask(
					tc.line,
				)
				if ok != tc.wantOk {
					t.Errorf(
						"MatchNumberedTask(%q) ok = %v, want %v",
						tc.line,
						ok,
						tc.wantOk,
					)

					continue
				}
				if ok &&
					match.Number != tc.wantID {
					t.Errorf(
						"MatchNumberedTask(%q) ID = %q, want %q",
						tc.line,
						match.Number,
						tc.wantID,
					)
				}
			}
		},
	)
}
