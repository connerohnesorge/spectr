package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// validateMixedIndentation checks mixed indentation test case
func validateMixedIndentation(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if !strings.Contains(tasks[0].Description, "Sub-item one") {
		t.Error("Failed to capture sub-item one")
	}
	if !strings.Contains(tasks[0].Description, "Nested item") {
		t.Error("Failed to capture nested item")
	}
	if !strings.Contains(tasks[0].Description, "Sub-item two") {
		t.Error("Failed to capture sub-item two")
	}
}

// validateTabsAndSpaces checks tabs and spaces mixed test case
func validateTabsAndSpaces(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "Tab-indented") {
		t.Error("Failed to capture tab-indented sub-item")
	}
	if !strings.Contains(desc, "Space-indented") {
		t.Error("Failed to capture space-indented sub-item")
	}
}

// validateEmptyLines checks empty lines within continuation test case
func validateEmptyLines(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if !strings.Contains(tasks[0].Description, "Item one") {
		t.Error("Failed to capture Item one")
	}
	if strings.Contains(tasks[0].Description, "Item two") {
		t.Error("Item two should not be captured (after blank line)")
	}
}

// validateSpecialChars checks special characters test case
func validateSpecialChars(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "Backslash") {
		t.Error("Failed to capture backslash item")
	}
	if !strings.Contains(desc, "Quote") {
		t.Error("Failed to capture quote item")
	}
	if !strings.Contains(desc, "\\test") {
		t.Error("Failed to preserve backslash in content")
	}
}

// validateUnicode checks unicode in continuation lines test case
func validateUnicode(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "🚀") {
		t.Error("Failed to preserve emoji")
	}
	if !strings.Contains(desc, "你好") {
		t.Error("Failed to preserve Chinese characters")
	}
	if !strings.Contains(desc, "مرحبا") {
		t.Error("Failed to preserve Arabic characters")
	}
	if !strings.Contains(desc, "Привет") {
		t.Error("Failed to preserve Russian characters")
	}
}

// validateLongContent checks very long continuation lines test case
func validateLongContent(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if len(tasks[0].Description) < 200 {
		t.Errorf(
			"Description too short (%d chars), expected to preserve long content",
			len(tasks[0].Description),
		)
	}
}

// validateCodeBlocks checks code blocks in continuation test case
func validateCodeBlocks(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "fmt.Println") {
		t.Error("Failed to preserve code example")
	}
	if !strings.Contains(desc, "[a-zA-Z0-9]") {
		t.Error("Failed to preserve regex")
	}
}

// validateMultipleSections checks multiple sections with multi-line tasks test case
func validateMultipleSections(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
	if !strings.Contains(tasks[0].Description, "Detail A") {
		t.Error("Task 1 missing Detail A")
	}
	if !strings.Contains(tasks[1].Description, "Detail X") {
		t.Error("Task 2 missing Detail X")
	}
	if !strings.Contains(tasks[2].Description, "Detail 3") {
		t.Error("Task 3 missing Detail 3")
	}
	if strings.Contains(tasks[0].Description, "Detail X") {
		t.Error("Task 1 incorrectly includes Task 2 details")
	}
}

// validateMixedLists checks numbered and unnumbered sub-items test case
func validateMixedLists(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "Bullet item one") {
		t.Error("Failed to capture bullet item")
	}
	if !strings.Contains(desc, "Numbered item") {
		t.Error("Failed to capture numbered item")
	}
}

// validateMinimalTask checks task with only continuation test case
func validateMinimalTask(t *testing.T, tasks []parsers.Task) {
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if !strings.Contains(tasks[0].Description, "minimal base text") {
		t.Error("Failed to capture continuation for minimal task")
	}
}

// TestParseTasksMdFuzzMultilineVariations tests various multi-line task description patterns
func TestParseTasksMdFuzzMultilineVariations(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		validate func(t *testing.T, tasks []parsers.Task)
	}{
		{
			name: "mixed indentation levels",
			markdown: `## 1. Mixed Indentation

- [ ] 1.1 Task with mixed indentation
  - Sub-item one
    - Nested item
  - Sub-item two
- [ ] 1.2 Next task
`,
			validate: validateMixedIndentation,
		},
		{
			name: "tabs and spaces mixed",
			markdown: `## 1. Tab and Space Mix

- [ ] 1.1 Task with tab indentation
	- Tab-indented sub-item
  - Space-indented sub-item
	  - Mixed indent sub-item
`,
			validate: validateTabsAndSpaces,
		},
		{
			name: "empty lines within continuation",
			markdown: `## 1. Empty Lines

- [ ] 1.1 Task with gaps
  - Item one

  - Item two (after blank line)
- [ ] 1.2 Next task
`,
			validate: validateEmptyLines,
		},
		{
			name: "special characters in sub-items",
			markdown: `## 1. Special Characters

- [ ] 1.1 Test special chars
  - Backslash: \test
  - Quote: "quoted"
  - Single: 'single'
  - Bracket: [nested]
  - Brace: {object}
`,
			validate: validateSpecialChars,
		},
		{
			name: "unicode in continuation lines",
			markdown: `## 1. Unicode Test

- [ ] 1.1 Unicode task
  - Emoji: 🚀🔧🐛
  - Chinese: 你好世界
  - Arabic: مرحبا بالعالم
  - Russian: Привет мир
`,
			validate: validateUnicode,
		},
		{
			name: "very long continuation lines",
			markdown: `## 1. Long Content

- [ ] 1.1 Task with long description
  - Very long line: ` + strings.Repeat("Lorem ipsum dolor sit amet, ", 10) + `
  - Another long line with numbers: ` + strings.Repeat("0123456789", 20),
			validate: validateLongContent,
		},
		{
			name: "code blocks in continuation",
			markdown: `## 1. Code Example

- [ ] 1.1 Task with code
  - Example: ` + "`fmt.Println(\"hello\")`" + `
  - Regex: ` + "`[a-zA-Z0-9]+`" + `
  - JSON: ` + "`{\"key\": \"value\"}`" + `
`,
			validate: validateCodeBlocks,
		},
		{
			name: "multiple sections with multi-line tasks",
			markdown: `## 1. Section One

- [ ] 1.1 Task one
  - Detail A
  - Detail B

## 2. Section Two

- [ ] 2.1 Task two
  - Detail X
  - Detail Y

## 3. Section Three

- [ ] 3.1 Task three
  - Detail 1
  - Detail 2
  - Detail 3
`,
			validate: validateMultipleSections,
		},
		{
			name: "numbered and unnumbered sub-items",
			markdown: `## 1. Mixed Lists

- [ ] 1.1 Task with mixed list
  - Bullet item one
  1. Numbered item one
  - Bullet item two
  2. Numbered item two
`,
			validate: validateMixedLists,
		},
		{
			name: "task with only continuation no base description",
			markdown: `## 1. Minimal Task

- [ ] 1.1
  - This task has minimal base text
  - All content in continuations
- [ ] 1.2 Normal task
`,
			validate: validateMinimalTask,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tasksMdPath := filepath.Join(tmpDir, "tasks.md")

			if err := os.WriteFile(tasksMdPath, []byte(tt.markdown), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := parseTasksMd(tasksMdPath)
			if err != nil {
				t.Fatalf("parseTasksMd() error = %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}

// validateOnlySections checks file with only sections no tasks
func validateOnlySections(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

// validateOnlyTasks checks file with only tasks no sections
func validateOnlyTasks(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

// validateContinuationAtEnd checks continuation lines at end of file
func validateContinuationAtEnd(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "Continuation three") {
		t.Error("Failed to capture last continuation line")
	}
}

// validateDeepNesting checks deeply nested indentation
func validateDeepNesting(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "Level 5") {
		t.Error("Failed to capture deeply nested items")
	}
}

// validateMarkdownFormatting checks task descriptions with markdown formatting
func validateMarkdownFormatting(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	desc := tasks[0].Description
	if !strings.Contains(desc, "**bold**") {
		t.Error("Failed to preserve markdown bold")
	}
	if !strings.Contains(desc, "_underscore_") {
		t.Error("Failed to preserve underscore emphasis")
	}
}

// validateWhitespaceOnly checks whitespace-only continuation lines treated as stop
func validateWhitespaceOnly(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if strings.Contains(tasks[0].Description, "Item two") {
		t.Error("Item two should not be captured (after blank line)")
	}
}

// validateUnusualIDs checks task IDs with unusual formats
func validateUnusualIDs(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) < 3 {
		t.Errorf("Expected at least 3 tasks, got %d", len(tasks))
	}
}

// validateRapidAlternation checks rapid alternation between tasks and continuations
func validateRapidAlternation(t *testing.T, tasks []parsers.Task, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tasks) != 4 {
		t.Errorf("Expected 4 tasks, got %d", len(tasks))
	}
	for i, task := range tasks {
		if !strings.Contains(task.Description, "Detail") {
			t.Errorf("Task %d missing detail", i+1)
		}
	}
}

// TestParseTasksMdFuzzEdgeCases tests edge cases and boundary conditions
func TestParseTasksMdFuzzEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		validate func(t *testing.T, tasks []parsers.Task, err error)
	}{
		{
			name: "file with only sections no tasks",
			markdown: `## 1. Section One

## 2. Section Two

## 3. Section Three
`,
			validate: validateOnlySections,
		},
		{
			name: "file with only tasks no sections",
			markdown: `- [ ] 1.1 Task without section
- [ ] 1.2 Another task
- [ ] 1.3 Third task
`,
			validate: validateOnlyTasks,
		},
		{
			name: "continuation lines at end of file",
			markdown: `## 1. Last Section

- [ ] 1.1 Last task
  - Continuation one
  - Continuation two
  - Continuation three`,
			validate: validateContinuationAtEnd,
		},
		{
			name: "deeply nested indentation",
			markdown: `## 1. Deep Nesting

- [ ] 1.1 Task with deep nesting
  - Level 1
    - Level 2
      - Level 3
        - Level 4
          - Level 5
`,
			validate: validateDeepNesting,
		},
		{
			name: "task descriptions with markdown formatting",
			markdown: `## 1. Formatted

- [ ] 1.1 Task with **bold** and *italic*
  - Item with _underscore_ emphasis
  - Link-like text: [description]
  - Code reference: ` + "`function()`" + `
`,
			validate: validateMarkdownFormatting,
		},
		{
			name: "whitespace-only continuation lines treated as stop",
			markdown: `## 1. Whitespace Test

- [ ] 1.1 Task
  - Item one

  - Item two (after blank line should not be captured)
`,
			validate: validateWhitespaceOnly,
		},
		{
			name: "task IDs with unusual formats",
			markdown: `## 1. Unusual IDs

- [ ] 1.1 Standard ID
- [ ] 1.10 Double digit
- [ ] 1.99 High number
- [ ] 1 Single digit
- [ ] a.1 Letter prefix
`,
			validate: validateUnusualIDs,
		},
		{
			name: "rapid alternation between tasks and continuations",
			markdown: `## 1. Rapid

- [ ] 1.1 Task A
  - Detail
- [ ] 1.2 Task B
  - Detail
- [ ] 1.3 Task C
  - Detail
- [ ] 1.4 Task D
  - Detail
`,
			validate: validateRapidAlternation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tasksMdPath := filepath.Join(tmpDir, "tasks.md")

			if err := os.WriteFile(tasksMdPath, []byte(tt.markdown), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := parseTasksMd(tasksMdPath)

			if tt.validate != nil {
				tt.validate(t, got, err)
			}
		})
	}
}

// testTaskParsed checks basic task parsing
func testTaskParsed(t *testing.T, task *parsers.Task, expectedID, expectedContent string) {
	if task.ID != expectedID || task.Description != expectedContent {
		t.Errorf("Task %s corrupted: %+v", expectedID, task)
	}
}

// testEscapedCharacters checks task with escaped characters
func testEscapedCharacters(t *testing.T, task *parsers.Task) {
	if !strings.Contains(task.Description, "TestJSONCValidation_SpecialCharacters") {
		t.Error("Task 1.2 lost title")
	}
	if !strings.Contains(task.Description, "\\") {
		t.Error("Task 1.2 lost backslash")
	}
	if !strings.Contains(task.Description, "\\n") {
		t.Error("Task 1.2 lost newline escape")
	}
	if !strings.Contains(task.Description, "\\t") {
		t.Error("Task 1.2 lost tab escape")
	}
}

// testUnicodeContent checks task with unicode
func testUnicodeContent(t *testing.T, task *parsers.Task) {
	if !strings.Contains(task.Description, "TestJSONCValidation_Unicode") {
		t.Error("Task 1.3 lost title")
	}
	if !strings.Contains(task.Description, "🚀") {
		t.Error("Task 1.3 lost emoji")
	}
	if !strings.Contains(task.Description, "你好") {
		t.Error("Task 1.3 lost Chinese")
	}
}

// testSubItems checks task with multiple sub-items
func testSubItems(t *testing.T, task *parsers.Task) {
	if task.ID != "1.4" {
		t.Errorf("Task 1.4 has wrong ID: %s", task.ID)
	}
	if !strings.Contains(task.Description, "Sub-item one") {
		t.Error("Task 1.4 lost Sub-item one")
	}
	if !strings.Contains(task.Description, "Sub-item three") {
		t.Error("Task 1.4 lost Sub-item three")
	}
}

// testNoCrossContamination checks no cross-contamination between tasks
func testNoCrossContamination(t *testing.T, tasks []parsers.Task) {
	if strings.Contains(tasks[0].Description, "Special") {
		t.Error("Task 1.1 incorrectly includes Task 1.2 content")
	}
	if strings.Contains(tasks[1].Description, "Unicode") {
		t.Error("Task 1.2 incorrectly includes Task 1.3 content")
	}
}

// testJSONSerializationIntegrity checks JSON output preserves content
func testJSONSerializationIntegrity(t *testing.T, jsonStr string) {
	if !strings.Contains(jsonStr, "TestJSONCValidation_SpecialCharacters") {
		t.Error("Special characters task title lost in JSON output")
	}
	if !strings.Contains(jsonStr, "Backslash") {
		t.Error("Backslash item lost in JSON output")
	}
	if !strings.Contains(jsonStr, "TestJSONCValidation_Unicode") {
		t.Error("Unicode task title lost in JSON output")
	}
	if !strings.Contains(jsonStr, "🚀") {
		t.Error("Emoji lost in JSON output")
	}
	if !strings.Contains(jsonStr, "Sub-item one") {
		t.Error("Sub-item content lost in JSON output")
	}
	if !strings.Contains(jsonStr, "Sub-item three") {
		t.Error("Last sub-item lost in JSON output")
	}
}

// TestParseTasksMdContinuationIntegrity ensures multi-line parsing doesn't corrupt data
func TestParseTasksMdContinuationIntegrity(t *testing.T) {
	markdown := `## 1. Property-Based Testing

- [ ] 1.1 Create test infrastructure
- [ ] 1.2 Implement TestJSONCValidation_SpecialCharacters with test cases for:
  - Backslash ` + "`\\`" + `
  - Quote ` + "`\"`" + `
  - Newline ` + "`\\n`" + `
  - Tab ` + "`\\t`" + `
- [ ] 1.3 Implement TestJSONCValidation_Unicode with test cases for:
  - Emoji (🚀, 💻, ✅)
  - Non-ASCII (你好, مرحبا)
- [ ] 1.4 Simple task
  - Sub-item one
  - Sub-item two
  - Sub-item three
`

	tmpDir := t.TempDir()
	tasksMdPath := filepath.Join(tmpDir, "tasks.md")

	if err := os.WriteFile(tasksMdPath, []byte(markdown), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	got, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("parseTasksMd() error = %v", err)
	}

	if len(got) != 4 {
		t.Fatalf("Expected 4 tasks, got %d", len(got))
	}

	testTaskParsed(t, &got[0], "1.1", "Create test infrastructure")
	testEscapedCharacters(t, &got[1])
	testUnicodeContent(t, &got[2])
	testSubItems(t, &got[3])
	testNoCrossContamination(t, got)

	tasksJSONPath := filepath.Join(tmpDir, "tasks.jsonc")
	if err := writeTasksJSONC(tasksJSONPath, got, nil, nil); err != nil {
		t.Fatalf("writeTasksJSONC() error = %v", err)
	}

	jsonContent, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read generated JSON: %v", err)
	}

	jsonStr := string(jsonContent)
	testJSONSerializationIntegrity(t, jsonStr)
}

// BenchmarkParseTasksMdMultiline benchmarks multi-line task parsing performance
func BenchmarkParseTasksMdMultiline(b *testing.B) {
	markdown := `## 1. Section One

- [ ] 1.1 Task one
  - Detail line 1
  - Detail line 2
  - Detail line 3
  - Detail line 4
  - Detail line 5

- [ ] 1.2 Task two
  - Sub A
  - Sub B

## 2. Section Two

- [ ] 2.1 Task three
  - Line 1
  - Line 2
  - Line 3

- [ ] 2.2 Task four with very long description that contains lots of text
  - Item 1
  - Item 2
  - Item 3
  - Item 4
  - Item 5
  - Item 6
  - Item 7
  - Item 8
  - Item 9
  - Item 10
`

	tmpDir := b.TempDir()
	tasksMdPath := filepath.Join(tmpDir, "tasks.md")

	if err := os.WriteFile(tasksMdPath, []byte(markdown), 0o644); err != nil {
		b.Fatalf("failed to write test file: %v", err)
	}

	b.ResetTimer()
	for range b.N {
		_, err := parseTasksMd(tasksMdPath)
		if err != nil {
			b.Fatalf("parseTasksMd() error = %v", err)
		}
	}
}
