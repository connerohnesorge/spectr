//nolint:revive // line-length-limit, early-return - task parsing patterns
package markdown

import (
	"bufio"
	"regexp"
	"strings"

	bf "github.com/russross/blackfriday/v2"
)

// taskCheckboxPattern matches task checkbox syntax from AST-parsed item text.
// The blackfriday parser strips the leading "- " from list items, so we match
// the checkbox directly: [ ], [x], [X]
var taskCheckboxPattern = regexp.MustCompile(`^\s*\[([xX ])\]\s*(.*)$`)

// ExtractTasks extracts all task checkbox items from the markdown AST.
// Tasks are identified by "- [ ]" (unchecked) or "- [x]" (checked) syntax.
func ExtractTasks(node *bf.Node) []Task {
	tasks := make([]Task, 0)
	if node == nil {
		return tasks
	}

	// Walk through all list items looking for task checkboxes
	node.Walk(func(n *bf.Node, entering bool) bf.WalkStatus {
		if !entering || n.Type != bf.Item {
			return bf.GoToNext
		}

		// Get the text content of the list item
		text := extractItemText(n)

		// Check if it matches task checkbox pattern
		matches := taskCheckboxPattern.FindStringSubmatch(text)
		if matches == nil {
			return bf.GoToNext
		}

		checkbox := matches[1]
		taskText := strings.TrimSpace(matches[2])

		tasks = append(tasks, Task{
			Checked: strings.ToLower(checkbox) == "x",
			Text:    taskText,
			Line:    nodeLineNumber(n),
		})

		return bf.GoToNext
	})

	return tasks
}

// extractItemText gets the raw text content of a list item,
// preserving the checkbox syntax if present.
func extractItemText(item *bf.Node) string {
	if item == nil {
		return ""
	}

	var result strings.Builder

	for child := item.FirstChild; child != nil; child = child.Next {
		if child.Type == bf.Paragraph {
			// Get the first text content which should contain the checkbox
			for pChild := child.FirstChild; pChild != nil; pChild = pChild.Next {
				//nolint:exhaustive // default case handles all other node types
				switch pChild.Type {
				case bf.Text:
					result.WriteString(string(pChild.Literal))
				case bf.Code:
					result.WriteString("`")
					result.WriteString(string(pChild.Literal))
					result.WriteString("`")
				case bf.Softbreak, bf.Hardbreak:
					result.WriteString(" ")
				default:
					result.WriteString(extractText(pChild))
				}
			}
		}
	}

	return result.String()
}

// CountTasks counts total and completed tasks in the markdown content.
// Returns (total, completed).
func CountTasks(node *bf.Node) (total, completed int) {
	tasks := ExtractTasks(node)
	total = len(tasks)
	for _, task := range tasks {
		if task.Checked {
			completed++
		}
	}

	return total, completed
}

// CountTasksFromContent parses content and counts tasks using line-by-line parsing.
// This approach is more lenient than AST-based parsing and will find task checkboxes
// even in non-standard markdown formatting (e.g., without blank lines before lists).
// Returns (total, completed).
func CountTasksFromContent(content string) (total, completed int) {
	// Use line-by-line parsing for backwards compatibility with the original
	// regex-based implementation, which matched task patterns anywhere in text.
	scanner := bufio.NewScanner(strings.NewReader(content))
	taskPattern := regexp.MustCompile(`^\s*-\s*\[([xX ])\]`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := taskPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		total++
		checkbox := strings.ToLower(strings.TrimSpace(matches[1]))
		if checkbox == "x" {
			completed++
		}
	}

	return total, completed
}

// ExtractTasksFromContent parses content and extracts tasks.
func ExtractTasksFromContent(content string) []Task {
	node := Parse([]byte(content))

	return ExtractTasks(node)
}

// TaskWithID represents a task with section context and ID.
// This is used for tasks.md parsing in the accept command.
type TaskWithID struct {
	ID          string // Task ID (e.g., "1.1")
	Section     string // Section name
	Description string // Task description
	Checked     bool   // Completion status
}

// taskWithIDPattern matches tasks with IDs: - [ ] 1.1 Description
var taskWithIDPattern = regexp.MustCompile(`^-\s+\[([ xX])\]\s+(\d+\.\d+)\s+(.+)$`)

// sectionHeaderPattern matches section headers: ## 1. Section Name
var sectionHeaderPattern = regexp.MustCompile(`^##\s+\d+\.\s+(.+)$`)

// ExtractTasksWithIDs extracts tasks with section context and IDs from tasks.md format.
// This matches the format expected by the accept command:
// ## 1. Section Name
// - [ ] 1.1 Task description
// - [x] 1.2 Completed task
func ExtractTasksWithIDs(content string) []TaskWithID {
	tasks := make([]TaskWithID, 0)

	var currentSection string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Check for section header
		if matches := sectionHeaderPattern.FindStringSubmatch(line); matches != nil {
			currentSection = strings.TrimSpace(matches[1])

			continue
		}

		// Check for task with ID
		if matches := taskWithIDPattern.FindStringSubmatch(line); matches != nil {
			checkbox := matches[1]
			taskID := matches[2]
			description := strings.TrimSpace(matches[3])

			tasks = append(tasks, TaskWithID{
				ID:          taskID,
				Section:     currentSection,
				Description: description,
				Checked:     strings.ToLower(strings.TrimSpace(checkbox)) == "x",
			})
		}
	}

	return tasks
}

// ExtractTasksWithIDsFromFile reads a file and extracts tasks with IDs.
func ExtractTasksWithIDsFromFile(path string) ([]TaskWithID, error) {
	node, err := ParseFile(path)
	if err != nil {
		return nil, err
	}

	// For tasks.md format, we need to use line-by-line parsing
	// because the AST doesn't preserve the exact checkbox format
	content := renderFullDocument(node)

	return ExtractTasksWithIDs(content), nil
}

// renderFullDocument renders the entire document back to text.
func renderFullDocument(node *bf.Node) string {
	if node == nil {
		return ""
	}
	var result strings.Builder
	for child := node.FirstChild; child != nil; child = child.Next {
		result.WriteString(renderNode(child))
	}

	return result.String()
}
