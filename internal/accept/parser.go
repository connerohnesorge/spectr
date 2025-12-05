// Package accept provides parsing and processing for tasks.md files,
// converting them into structured Task/Section objects for change acceptance.
package accept

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// Constants for regex match group counts.
const (
	numberedSectionMatchCount = 3 // Full match + 2 groups (number, name)
	plainSectionMatchCount    = 2 // Full match + 1 group (name)
	taskLineMatchCount        = 4 // Full match + 3 groups (marker, id, desc)
	decimalBase               = 10
)

var (
	// sectionHeaderNumbered matches "## N. Section Name" with explicit number
	sectionHeaderNumbered = regexp.MustCompile(`^##\s+(\d+)\.\s+(.+)$`)

	// sectionHeaderPlain matches "## Section Name" without number
	sectionHeaderPlain = regexp.MustCompile(`^##\s+(.+)$`)

	// taskLine matches task lines with optional ID:
	// - [ ] 1.2.3 Description or - [x] Description
	taskLine = regexp.MustCompile(
		`^\s*-\s*\[([xX ])\]\s*(?:(\d+(?:\.\d+)*)\s+)?(.+)$`,
	)

	// detailLine matches indented lines that are not task lines
	// (2+ spaces or tab at the start, not starting with - [ ] or - [x])
	detailLinePrefix = regexp.MustCompile(`^(\s{2,}|\t)`)
)

// ParseTasksFile parses a tasks.md file and returns structured sections.
func ParseTasksFile(filePath string) ([]Section, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	return ParseTasks(file)
}

// parseState holds mutable state during parsing.
type parseState struct {
	sections       []Section
	currentSection *Section
	currentTask    *Task
	flatTasks      []Task
	autoSectionNum int
	autoTaskNum    int
}

// ParseTasks parses tasks from an io.Reader and returns structured sections.
func ParseTasks(reader io.Reader) ([]Section, error) {
	scanner := bufio.NewScanner(reader)
	state := &parseState{}

	for scanner.Scan() {
		processLine(scanner.Text(), state)
	}

	// Finalize last section
	if state.currentSection != nil {
		state.currentSection.Tasks = organizeTaskHierarchy(state.flatTasks)
		state.sections = append(state.sections, *state.currentSection)
	}

	return state.sections, scanner.Err()
}

// processLine handles a single line during parsing.
func processLine(line string, state *parseState) {
	// Check for section header
	if sec := parseSectionHeader(line, &state.autoSectionNum); sec != nil {
		finalizeSection(state)
		state.currentSection = sec
		state.flatTasks = nil
		state.currentTask = nil
		state.autoTaskNum = 0

		return
	}

	// Check for task line
	task := parseTaskLine(line, state.currentSection, &state.autoTaskNum)
	if task != nil {
		state.flatTasks = append(state.flatTasks, *task)
		state.currentTask = &state.flatTasks[len(state.flatTasks)-1]

		return
	}

	// Check for detail line (indented continuation)
	if state.currentTask != nil && isDetailLine(line) {
		appendDetailLine(state.currentTask, line)
		state.flatTasks[len(state.flatTasks)-1] = *state.currentTask
	}
	// Blank lines are intentionally ignored
}

// finalizeSection adds the current section to sections list if present.
func finalizeSection(state *parseState) {
	if state.currentSection != nil {
		state.currentSection.Tasks = organizeTaskHierarchy(state.flatTasks)
		state.sections = append(state.sections, *state.currentSection)
	}
}

// parseSectionHeader parses a section header line and returns a Section or nil.
func parseSectionHeader(line string, autoNum *int) *Section {
	// Try numbered section first: "## 1. Section Name"
	matches := sectionHeaderNumbered.FindStringSubmatch(line)
	if len(matches) == numberedSectionMatchCount {
		num := 0
		for _, c := range matches[1] {
			num = num*decimalBase + int(c-'0')
		}

		return &Section{
			Number: num,
			Name:   strings.TrimSpace(matches[2]),
		}
	}

	// Try plain section: "## Section Name"
	matches = sectionHeaderPlain.FindStringSubmatch(line)
	if len(matches) == plainSectionMatchCount {
		*autoNum++

		return &Section{
			Number: *autoNum,
			Name:   strings.TrimSpace(matches[1]),
		}
	}

	return nil
}

// parseTaskLine parses a task line and returns a Task or nil.
func parseTaskLine(
	line string,
	currentSection *Section,
	autoTaskNum *int,
) *Task {
	matches := taskLine.FindStringSubmatch(line)
	if len(matches) < taskLineMatchCount {
		return nil
	}

	marker := strings.ToLower(strings.TrimSpace(matches[1]))
	completed := marker == "x"
	taskID := matches[2]
	description := strings.TrimSpace(matches[3])

	// Auto-generate ID if not provided
	if taskID == "" {
		*autoTaskNum++
		if currentSection != nil {
			taskID = formatTaskID(currentSection.Number, *autoTaskNum)
		} else {
			taskID = formatTaskID(0, *autoTaskNum)
		}
	}

	return &Task{
		ID:          taskID,
		Description: description,
		Completed:   completed,
		Subtasks:    nil,
	}
}

// formatTaskID creates a task ID from section number and task number.
func formatTaskID(sectionNum, taskNum int) string {
	if sectionNum > 0 {
		return intToString(sectionNum) + "." + intToString(taskNum)
	}

	return intToString(taskNum)
}

// intToString converts an int to a string without using strconv.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToString(-n)
	}

	return positiveIntToString(n)
}

// positiveIntToString converts a positive int to a string.
func positiveIntToString(n int) string {
	var digits []byte
	for val := n; val > 0; val /= decimalBase {
		digits = append([]byte{byte('0' + val%decimalBase)}, digits...)
	}

	return string(digits)
}

// isDetailLine checks if a line is an indented detail line (not a task line).
func isDetailLine(line string) bool {
	// Must start with indentation
	if !detailLinePrefix.MatchString(line) {
		return false
	}

	// Must not be a task line
	trimmed := strings.TrimSpace(line)

	return !strings.HasPrefix(trimmed, "- [")
}

// appendDetailLine appends a trimmed detail line to the task's description.
func appendDetailLine(task *Task, line string) {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return
	}
	if task.Description != "" {
		task.Description += "\n" + trimmedLine
	} else {
		task.Description = trimmedLine
	}
}
