package markdown

import (
	"bytes"
	"strings"

	"github.com/russross/blackfriday/v2"
)

// Checkbox text length constants.
const (
	checkboxWithSpaceLen = 4 // Length of "[ ] " or "[x] "
	checkboxNoSpaceLen   = 3 // Length of "[ ]" or "[x]"
)

// extractTasks extracts task checkboxes from the AST.
// Tasks are list items with checkbox markers `- [ ]` or `- [x]`/`- [X]`.
// Preserves hierarchical structure through Children field.
func extractTasks(
	node *blackfriday.Node,
	source []byte,
	lineIndex *lineIndex,
) []Task {
	var tasks []Task

	node.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		// Only process list items when entering
		if !entering || n.Type != blackfriday.List {
			return blackfriday.GoToNext
		}

		// Process this list's items
		listTasks := extractTasksFromList(n, source, lineIndex, 0)
		tasks = append(tasks, listTasks...)

		// Skip children since we processed them
		return blackfriday.SkipChildren
	})

	return tasks
}

// extractTasksFromList extracts tasks from a list node and nested lists.
func extractTasksFromList(
	listNode *blackfriday.Node,
	source []byte,
	lineIndex *lineIndex,
	depth int,
) []Task {
	var tasks []Task

	for item := listNode.FirstChild; item != nil; item = item.Next {
		if item.Type != blackfriday.Item {
			continue
		}

		task, isTask := parseTaskItem(item, source, lineIndex)
		if !isTask {
			nested := extractNestedTasks(item, source, lineIndex, depth)
			tasks = append(tasks, nested...)

			continue
		}

		task.Children = findChildTasks(item, source, lineIndex, depth)
		tasks = append(tasks, task)
	}

	return tasks
}

// extractNestedTasks finds tasks in nested lists of a non-task item.
func extractNestedTasks(
	item *blackfriday.Node,
	source []byte,
	lineIndex *lineIndex,
	depth int,
) []Task {
	var tasks []Task
	for child := item.FirstChild; child != nil; child = child.Next {
		if child.Type == blackfriday.List {
			nested := extractTasksFromList(child, source, lineIndex, depth+1)
			tasks = append(tasks, nested...)
		}
	}

	return tasks
}

// findChildTasks finds child tasks in nested lists of a task item.
func findChildTasks(
	item *blackfriday.Node,
	source []byte,
	lineIndex *lineIndex,
	depth int,
) []Task {
	for child := item.FirstChild; child != nil; child = child.Next {
		if child.Type == blackfriday.List {
			return extractTasksFromList(child, source, lineIndex, depth+1)
		}
	}

	return nil
}

// parseTaskItem checks if a list item is a task checkbox.
func parseTaskItem(
	item *blackfriday.Node,
	source []byte,
	lineIndex *lineIndex,
) (Task, bool) {
	firstText := findFirstTextNode(item)
	if firstText == nil || len(firstText.Literal) == 0 {
		return Task{}, false
	}

	text := string(firstText.Literal)
	hasCheckbox, checked := parseCheckbox(text)
	if !hasCheckbox {
		return Task{}, false
	}

	lineNum, fullLine := findTaskLocation(source, lineIndex, text, checked)

	return Task{
		Line:     fullLine,
		Checked:  checked,
		LineNum:  lineNum,
		Children: nil,
	}, true
}

// findFirstTextNode finds the first text node in a list item.
func findFirstTextNode(item *blackfriday.Node) *blackfriday.Node {
	for child := item.FirstChild; child != nil; child = child.Next {
		if child.Type == blackfriday.Paragraph {
			return findTextInParagraph(child)
		}
		if child.Type == blackfriday.Text {
			return child
		}
	}

	return nil
}

// findTextInParagraph finds the first text node in a paragraph.
func findTextInParagraph(para *blackfriday.Node) *blackfriday.Node {
	for node := para.FirstChild; node != nil; node = node.Next {
		if node.Type == blackfriday.Text {
			return node
		}
	}

	return nil
}

// parseCheckbox determines if text starts with a checkbox pattern.
func parseCheckbox(text string) (hasCheckbox, checked bool) {
	uncheckedPrefixes := []string{"[ ] ", "[ ]"}
	checkedPrefixes := []string{"[x] ", "[x]", "[X] ", "[X]"}

	for _, prefix := range uncheckedPrefixes {
		if strings.HasPrefix(text, prefix) {
			return true, false
		}
	}
	for _, prefix := range checkedPrefixes {
		if strings.HasPrefix(text, prefix) {
			return true, true
		}
	}

	return false, false
}

// findTaskLocation determines the line number and full line text for a task.
func findTaskLocation(
	source []byte,
	lineIndex *lineIndex,
	text string,
	checked bool,
) (int, string) {
	lineNum := 1
	fullLine := text

	if lineIndex != nil {
		taskLineNum, taskLine := findTaskLine(source, lineIndex, text, checked)
		if taskLineNum > 0 {
			lineNum = taskLineNum
			fullLine = taskLine
		}
	}

	return lineNum, fullLine
}

// findTaskLine searches for a task line in the source.
// Returns the line number (1-indexed) and full text.
// The checked parameter selects between checked/unchecked patterns.
//
//nolint:revive // checked is a dispatch parameter, not control coupling
func findTaskLine(
	source []byte,
	_ *lineIndex,
	text string,
	checked bool,
) (int, string) {
	if checked {
		return findCheckedTaskLine(source, text)
	}

	return findUncheckedTaskLine(source, text)
}

// findCheckedTaskLine searches for a checked task line.
func findCheckedTaskLine(source []byte, text string) (int, string) {
	return searchTaskLine(source, text, checkedPatterns())
}

// findUncheckedTaskLine searches for an unchecked task line.
func findUncheckedTaskLine(source []byte, text string) (int, string) {
	return searchTaskLine(source, text, uncheckedPatterns())
}

// searchTaskLine searches for a task line with given patterns.
func searchTaskLine(
	source []byte,
	text string,
	patterns []string,
) (int, string) {
	lines := bytes.Split(source, []byte("\n"))
	textContent := extractTextContent(text)

	for i, line := range lines {
		lineStr := string(line)
		if match := matchTaskLine(lineStr, patterns, textContent); match {
			return i + 1, lineStr // 1-indexed
		}
	}

	// Fallback: return 1 and the parsed text
	return 1, text
}

// checkedPatterns returns the patterns for checked checkboxes.
func checkedPatterns() []string {
	return []string{
		"- [x]", "- [X]", "* [x]", "* [X]", "+ [x]", "+ [X]",
	}
}

// uncheckedPatterns returns the patterns for unchecked checkboxes.
func uncheckedPatterns() []string {
	return []string{"- [ ]", "* [ ]", "+ [ ]"}
}

// extractTextContent removes the checkbox prefix from text.
func extractTextContent(text string) string {
	textContent := text
	switch {
	case strings.HasPrefix(text, "[ ] "):
		textContent = text[checkboxWithSpaceLen:]
	case strings.HasPrefix(text, "[x] "), strings.HasPrefix(text, "[X] "):
		textContent = text[checkboxWithSpaceLen:]
	case strings.HasPrefix(text, "[ ]"):
		textContent = text[checkboxNoSpaceLen:]
	case strings.HasPrefix(text, "[x]"), strings.HasPrefix(text, "[X]"):
		textContent = text[checkboxNoSpaceLen:]
	}

	return strings.TrimSpace(textContent)
}

// matchTaskLine checks if a line matches any task pattern.
func matchTaskLine(line string, patterns []string, textContent string) bool {
	trimmedLine := strings.TrimLeft(line, " \t")

	for _, pattern := range patterns {
		if !strings.HasPrefix(trimmedLine, pattern) {
			continue
		}
		afterCheckbox := strings.TrimSpace(trimmedLine[len(pattern):])
		if textContent != "" && strings.HasPrefix(afterCheckbox, textContent) {
			return true
		}
		if textContent == "" && afterCheckbox == "" {
			return true
		}
	}

	return false
}
