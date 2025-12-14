package regex

import "regexp"

// Match group indices for NumberedTask pattern.
const (
	numberedTaskMatchLen  = 4
	numberedTaskDescIndex = 3
)

// Task patterns - all pre-compiled at package init.
var (
	// TaskCheckbox matches "- [ ]" or "- [x]" task items and captures
	// the checkbox state. The captured group contains 'x', 'X', or ' '
	// indicating completion status.
	TaskCheckbox = regexp.MustCompile(`^\s*-\s*\[([xX ])\]`)

	// NumberedTask matches "- [ ] 1.1 Description" format used in
	// tasks.md files. Captures: [1] checkbox state, [2] task ID, [3] desc.
	NumberedTask = regexp.MustCompile(`^-\s+\[([ xX])\]\s+(\d+\.\d+)\s+(.+)$`)

	// NumberedSection matches "## 1. Section Name" format used in
	// tasks.md files. Captures the section name (without the number prefix).
	NumberedSection = regexp.MustCompile(`^##\s+\d+\.\s+(.+)$`)
)

// MatchTaskCheckbox checks if a line contains a task checkbox and
// extracts the state. Returns the checkbox state as a rune ('x', 'X',
// or ' ') and true if matched, or 0 and false otherwise.
func MatchTaskCheckbox(line string) (state rune, ok bool) {
	matches := TaskCheckbox.FindStringSubmatch(line)
	if len(matches) < 2 {
		return 0, false
	}
	// The match is a single character: 'x', 'X', or ' '
	return rune(matches[1][0]), true
}

// NumberedTaskMatch holds the parsed result of a numbered task line.
type NumberedTaskMatch struct {
	Checkbox    string
	ID          string
	Description string
}

// MatchNumberedTask parses a numbered task line from tasks.md format.
// Returns the parsed task match and true if matched, or nil and false.
//
// Example input: "- [ ] 1.1 Create the regex package"
// Returns: &NumberedTaskMatch{" ", "1.1", "Create the regex package"}, true
func MatchNumberedTask(line string) (*NumberedTaskMatch, bool) {
	matches := NumberedTask.FindStringSubmatch(line)
	if len(matches) < numberedTaskMatchLen {
		return nil, false
	}

	return &NumberedTaskMatch{
		Checkbox:    matches[1],
		ID:          matches[2],
		Description: matches[numberedTaskDescIndex],
	}, true
}

// MatchNumberedSection parses a numbered section header from tasks.md format.
// Returns the section name (without number prefix) and true if matched,
// or empty string and false otherwise.
//
// Example input: "## 1. Core Accept Command"
// Returns: "Core Accept Command", true
func MatchNumberedSection(line string) (name string, ok bool) {
	matches := NumberedSection.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// IsTaskChecked returns true if the checkbox state indicates completion.
// Accepts 'x' or 'X' as checked states.
func IsTaskChecked(state rune) bool {
	return state == 'x' || state == 'X'
}
