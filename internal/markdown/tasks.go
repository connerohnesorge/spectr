package markdown

import (
	"strings"
)

// Task parsing constants
const (
	minNumberedTaskLen = 6 // Minimum length for "- [ ] X"
	afterCheckboxIdx   = 5 // Index after checkbox for rest of line
)

// Task Line Matchers
//
// These functions check individual lines for task patterns.
// They are designed for efficient line-by-line iteration.
//
// For bulk extraction, use ParseDocument + CountTasks instead.

// MatchTaskCheckbox checks if a line contains a task checkbox and
// extracts the state. Returns the checkbox state as a rune ('x', 'X',
// or ' ') and true if matched, or 0 and false otherwise.
func MatchTaskCheckbox(line string) (state rune, ok bool) {
	trimmed := strings.TrimLeft(line, " \t")

	if !strings.HasPrefix(trimmed, "- [") {
		return 0, false
	}

	if len(trimmed) < minCheckboxLen {
		return 0, false
	}

	if trimmed[closeBracketIdx] != ']' {
		return 0, false
	}

	return rune(trimmed[checkboxCharIdx]), true
}

// IsTaskChecked returns true if the checkbox state indicates completion.
// Accepts 'x' or 'X' as checked states.
func IsTaskChecked(state rune) bool {
	return state == 'x' || state == 'X'
}

// MatchNumberedTask parses a numbered task line from tasks.md format.
// Returns the parsed task and true if matched, or nil and false.
//
// Example input: "- [ ] 1.1 Create the regex package"
// Returns: &NumberedTask{ID: "1.1", Description: "Create...", Checked: false}
func MatchNumberedTask(line string) (*NumberedTask, bool) {
	trimmed := strings.TrimSpace(line)

	if !strings.HasPrefix(trimmed, "- [") {
		return nil, false
	}

	if len(trimmed) < minNumberedTaskLen {
		return nil, false
	}

	if trimmed[closeBracketIdx] != ']' {
		return nil, false
	}

	checkbox := trimmed[checkboxCharIdx]
	checked := checkbox == 'x' || checkbox == 'X'

	rest := strings.TrimSpace(trimmed[afterCheckboxIdx:])

	id, description := parseTaskID(rest)
	if id == "" {
		return nil, false
	}

	return &NumberedTask{
		ID:          id,
		Description: description,
		Checked:     checked,
	}, true
}

// parseTaskID extracts the task ID (e.g., "1.1") from the start of a string.
func parseTaskID(s string) (id, description string) {
	i := 0

	// Parse first number
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}

	if i == 0 || i >= len(s) || s[i] != '.' {
		return "", ""
	}

	dotPos := i
	i++

	// Parse second number
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}

	if i == dotPos+1 {
		return "", ""
	}

	if i < len(s) && s[i] != ' ' {
		return "", ""
	}

	id = s[:i]
	if i < len(s) {
		description = strings.TrimSpace(s[i:])
	}

	return id, description
}

// MatchNumberedSection parses a numbered section header from tasks.md format.
// Returns the section name (without number prefix) and true if matched.
//
// Example input: "## 1. Core Accept Command"
// Returns: "Core Accept Command", true
func MatchNumberedSection(line string) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)

	if !strings.HasPrefix(trimmed, "## ") {
		return "", false
	}

	rest := strings.TrimPrefix(trimmed, "## ")

	i := 0
	for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
		i++
	}

	if i == 0 || i >= len(rest) || rest[i] != '.' {
		return "", false
	}

	i++
	rest = strings.TrimSpace(rest[i:])

	if rest == "" {
		return "", false
	}

	return rest, true
}

// CountTasksInContent counts task checkboxes in content.
// Returns total count and completed count.
//
// For multiple queries on the same content, prefer:
//
//	doc, _ := ParseDocument([]byte(content))
//	total, completed := doc.CountTasks()
func CountTasksInContent(content string) (total, completed int) {
	doc, err := ParseDocument([]byte(content))
	if err != nil {
		return 0, 0
	}

	return doc.CountTasks()
}

// Internal task extraction functions (used by ParseDocument)

// extractTasks extracts task checkboxes from lines.
func extractTasks(doc *Document) {
	var tasks []Task

	for i, line := range doc.Lines {
		task, ok := parseTaskLine(line, i+1)
		if ok {
			tasks = append(tasks, task)
		}
	}

	doc.Tasks = buildTaskHierarchy(tasks)
}

// parseTaskLine parses a single line for task checkbox.
func parseTaskLine(line string, lineNum int) (Task, bool) {
	indent := 0

	for _, ch := range line {
		switch ch {
		case ' ':
			indent++
		case '\t':
			indent += tabWidth
		default:
			goto doneCount
		}
	}

doneCount:
	trimmed := strings.TrimSpace(line)

	if !strings.HasPrefix(trimmed, "- [") {
		return Task{}, false
	}

	if len(trimmed) < minCheckboxLen {
		return Task{}, false
	}

	checkChar := trimmed[checkboxCharIdx]
	if trimmed[closeBracketIdx] != ']' {
		return Task{}, false
	}

	checked := checkChar == 'x' || checkChar == 'X'

	return Task{
		Text:     line,
		Checked:  checked,
		Line:     lineNum,
		Indent:   indent,
		Children: nil,
	}, true
}

// buildTaskHierarchy organizes flat tasks into hierarchical structure.
func buildTaskHierarchy(tasks []Task) []Task {
	if len(tasks) == 0 {
		return tasks
	}

	var result []Task
	var stack []*Task

	for i := range tasks {
		task := &tasks[i]

		for len(stack) > 0 && stack[len(stack)-1].Indent >= task.Indent {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			result = append(result, *task)
			stack = append(stack, &result[len(result)-1])
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, *task)
			stack = append(stack, &parent.Children[len(parent.Children)-1])
		}
	}

	return result
}
