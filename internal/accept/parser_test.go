package accept

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// =============================================================================
// Section Parsing Tests
// =============================================================================

func TestParseSectionHeader_Numbered(t *testing.T) {
	autoNum := 0
	section := parseSectionHeader("## 1. Implementation", &autoNum)

	if section == nil {
		t.Fatal("Expected section to be parsed, got nil")
	}
	if section.Number != 1 {
		t.Errorf("Expected section number 1, got %d", section.Number)
	}
	if section.Name != "Implementation" {
		t.Errorf("Expected section name 'Implementation', got %q", section.Name)
	}
	if autoNum != 0 {
		t.Errorf("Auto number should not be incremented for numbered sections, got %d", autoNum)
	}
}

func TestParseSectionHeader_Plain(t *testing.T) {
	autoNum := 0
	section := parseSectionHeader("## Validation", &autoNum)

	if section == nil {
		t.Fatal("Expected section to be parsed, got nil")
	}
	if section.Number != 1 {
		t.Errorf("Expected section number 1 (auto-numbered), got %d", section.Number)
	}
	if section.Name != "Validation" {
		t.Errorf("Expected section name 'Validation', got %q", section.Name)
	}
	if autoNum != 1 {
		t.Errorf("Auto number should be incremented to 1, got %d", autoNum)
	}
}

func TestParseSectionHeader_MultiDigit(t *testing.T) {
	autoNum := 0
	section := parseSectionHeader("## 10. Large Section", &autoNum)

	if section == nil {
		t.Fatal("Expected section to be parsed, got nil")
	}
	if section.Number != 10 {
		t.Errorf("Expected section number 10, got %d", section.Number)
	}
	if section.Name != "Large Section" {
		t.Errorf("Expected section name 'Large Section', got %q", section.Name)
	}
}

func TestParseSectionHeader_NotSection(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{"Regular text", "This is not a section"},
		{"H1 heading", "# Main Title"},
		{"H3 heading", "### Subsection"},
		{"Task line", "- [ ] Task"},
		{"Empty line", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			autoNum := 0
			section := parseSectionHeader(tt.line, &autoNum)
			if section != nil {
				t.Errorf("Expected nil for %q, got section %+v", tt.line, section)
			}
		})
	}
}

// =============================================================================
// Task Parsing Tests
// =============================================================================

func TestParseTask_SimpleIncomplete(t *testing.T) {
	input := `## 1. Implementation
- [ ] 1.1 Create schema`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}
	if len(sections[0].Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(sections[0].Tasks))
	}

	task := sections[0].Tasks[0]
	if task.ID != "1.1" {
		t.Errorf("Expected task ID '1.1', got %q", task.ID)
	}
	if task.Description != "Create schema" {
		t.Errorf("Expected description 'Create schema', got %q", task.Description)
	}
	if task.Completed {
		t.Error("Expected task to be incomplete")
	}
}

func TestParseTask_Completed(t *testing.T) {
	input := `## 2. Tasks
- [x] 2.1 Done task`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 || len(sections[0].Tasks) != 1 {
		t.Fatal("Expected 1 section with 1 task")
	}

	task := sections[0].Tasks[0]
	if !task.Completed {
		t.Error("Expected task to be completed")
	}
	if task.ID != "2.1" {
		t.Errorf("Expected task ID '2.1', got %q", task.ID)
	}
}

func TestParseTask_CompletedUppercase(t *testing.T) {
	input := `## 2. Tasks
- [X] 2.2 Done`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 || len(sections[0].Tasks) != 1 {
		t.Fatal("Expected 1 section with 1 task")
	}

	task := sections[0].Tasks[0]
	if !task.Completed {
		t.Error("Expected task to be completed (uppercase X)")
	}
}

func TestParseTask_WithoutID(t *testing.T) {
	input := `## 1. Tasks
- [ ] Task no ID`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 || len(sections[0].Tasks) != 1 {
		t.Fatal("Expected 1 section with 1 task")
	}

	task := sections[0].Tasks[0]
	if task.ID == "" {
		t.Error("Expected auto-generated ID, got empty string")
	}
	// Auto-generated ID should be section.taskNum format
	if task.ID != "1.1" {
		t.Errorf("Expected auto-generated ID '1.1', got %q", task.ID)
	}
	if task.Description != "Task no ID" {
		t.Errorf("Expected description 'Task no ID', got %q", task.Description)
	}
}

func TestParseTask_MultipleWithoutID(t *testing.T) {
	input := `## 3. Tasks
- [ ] First task
- [ ] Second task
- [x] Third task`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 || len(sections[0].Tasks) != 3 {
		t.Fatalf("Expected 1 section with 3 tasks, got %d sections with %d tasks",
			len(sections), len(sections[0].Tasks))
	}

	expectedIDs := []string{"3.1", "3.2", "3.3"}
	for i, task := range sections[0].Tasks {
		if task.ID != expectedIDs[i] {
			t.Errorf("Task %d: expected ID %q, got %q", i, expectedIDs[i], task.ID)
		}
	}
}

// =============================================================================
// Nested Subtask Tests
// =============================================================================

func TestParseTask_NestedSubtask(t *testing.T) {
	input := `## 1. Implementation
- [ ] 1.1 Parent task
- [ ] 1.1.1 Child subtask`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}

	// Should have 1 top-level task with 1 subtask
	if len(sections[0].Tasks) != 1 {
		t.Fatalf("Expected 1 top-level task, got %d", len(sections[0].Tasks))
	}

	parent := sections[0].Tasks[0]
	if parent.ID != "1.1" {
		t.Errorf("Expected parent ID '1.1', got %q", parent.ID)
	}
	if len(parent.Subtasks) != 1 {
		t.Fatalf("Expected 1 subtask, got %d", len(parent.Subtasks))
	}

	child := parent.Subtasks[0]
	if child.ID != "1.1.1" {
		t.Errorf("Expected subtask ID '1.1.1', got %q", child.ID)
	}
	if child.Description != "Child subtask" {
		t.Errorf("Expected subtask description 'Child subtask', got %q", child.Description)
	}
}

func TestParseTask_DeeplyNested(t *testing.T) {
	input := `## 1. Deep Nesting
- [ ] 1 Level 1
- [ ] 1.1 Level 2
- [ ] 1.1.1 Level 3
- [ ] 1.1.1.1 Level 4`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}

	// Should have 1 top-level task
	if len(sections[0].Tasks) != 1 {
		t.Fatalf("Expected 1 top-level task, got %d", len(sections[0].Tasks))
	}

	// Level 1: ID "1"
	level1 := sections[0].Tasks[0]
	if level1.ID != "1" {
		t.Errorf("Expected level 1 ID '1', got %q", level1.ID)
	}
	if len(level1.Subtasks) != 1 {
		t.Fatalf("Expected 1 subtask at level 1, got %d", len(level1.Subtasks))
	}

	// Level 2: ID "1.1"
	level2 := level1.Subtasks[0]
	if level2.ID != "1.1" {
		t.Errorf("Expected level 2 ID '1.1', got %q", level2.ID)
	}
	if len(level2.Subtasks) != 1 {
		t.Fatalf("Expected 1 subtask at level 2, got %d", len(level2.Subtasks))
	}

	// Level 3: ID "1.1.1"
	level3 := level2.Subtasks[0]
	if level3.ID != "1.1.1" {
		t.Errorf("Expected level 3 ID '1.1.1', got %q", level3.ID)
	}
	if len(level3.Subtasks) != 1 {
		t.Fatalf("Expected 1 subtask at level 3, got %d", len(level3.Subtasks))
	}

	// Level 4: ID "1.1.1.1"
	level4 := level3.Subtasks[0]
	if level4.ID != "1.1.1.1" {
		t.Errorf("Expected level 4 ID '1.1.1.1', got %q", level4.ID)
	}
	if level4.Description != "Level 4" {
		t.Errorf("Expected level 4 description 'Level 4', got %q", level4.Description)
	}
}

func TestParseTask_SubtaskWithoutParent(t *testing.T) {
	input := `## 1. Orphans
- [ ] 1.1.1 Orphan subtask`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}

	// Orphan should be treated as top-level
	if len(sections[0].Tasks) != 1 {
		t.Fatalf(
			"Expected orphan to be treated as top-level task, got %d tasks",
			len(sections[0].Tasks),
		)
	}

	task := sections[0].Tasks[0]
	if task.ID != "1.1.1" {
		t.Errorf("Expected task ID '1.1.1', got %q", task.ID)
	}
	if task.Description != "Orphan subtask" {
		t.Errorf("Expected description 'Orphan subtask', got %q", task.Description)
	}
}

func TestParseTask_MultipleSiblings(t *testing.T) {
	input := `## 1. Tests
- [ ] 1.3 Add tests
  - [ ] 1.3.1 Unit tests
  - [ ] 1.3.2 Integration tests`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}

	// Should have 1 top-level task with 2 subtasks
	if len(sections[0].Tasks) != 1 {
		t.Fatalf("Expected 1 top-level task, got %d", len(sections[0].Tasks))
	}

	parent := sections[0].Tasks[0]
	if len(parent.Subtasks) != 2 {
		t.Fatalf("Expected 2 subtasks, got %d", len(parent.Subtasks))
	}

	if parent.Subtasks[0].ID != "1.3.1" {
		t.Errorf("Expected first subtask ID '1.3.1', got %q", parent.Subtasks[0].ID)
	}
	if parent.Subtasks[1].ID != "1.3.2" {
		t.Errorf("Expected second subtask ID '1.3.2', got %q", parent.Subtasks[1].ID)
	}
}

// =============================================================================
// Detail Line Tests
// =============================================================================

func TestParseTask_WithDetailLines(t *testing.T) {
	input := `## 1. Implementation
- [ ] 1.1 Create database schema
  - Parse requirement headers
  - Extract requirement name`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 || len(sections[0].Tasks) != 1 {
		t.Fatal("Expected 1 section with 1 task")
	}

	task := sections[0].Tasks[0]
	expectedDesc := "Create database schema\n- Parse requirement headers\n- Extract requirement name"
	if task.Description != expectedDesc {
		t.Errorf("Expected description:\n%q\nGot:\n%q", expectedDesc, task.Description)
	}
}

func TestParseTask_DetailNotTask(t *testing.T) {
	// Indented task line should be parsed as a task, not a detail line
	input := `## 1. Tasks
- [ ] 1.1 Parent
  - [ ] 1.1.1 This is a subtask, not detail`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}

	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}

	// The nested task should be recognized as a task and become a subtask
	parent := sections[0].Tasks[0]
	if len(parent.Subtasks) != 1 {
		t.Fatalf("Expected 1 subtask (nested task), got %d", len(parent.Subtasks))
	}

	subtask := parent.Subtasks[0]
	if subtask.ID != "1.1.1" {
		t.Errorf("Expected subtask ID '1.1.1', got %q", subtask.ID)
	}
	if subtask.Description != "This is a subtask, not detail" {
		t.Errorf("Expected subtask description, got %q", subtask.Description)
	}
}

func TestParseTask_TabIndentedDetail(t *testing.T) {
	input := "## 1. Tasks\n- [ ] 1.1 Task with tab detail\n\tTabbed detail line"

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}

	task := sections[0].Tasks[0]
	if !strings.Contains(task.Description, "Tabbed detail line") {
		t.Errorf("Expected tab-indented detail to be appended, got %q", task.Description)
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestParseTasks_EmptyInput(t *testing.T) {
	sections, err := ParseTasks(strings.NewReader(""))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 0 {
		t.Errorf("Expected nil or empty sections for empty input, got %d sections", len(sections))
	}
}

func TestParseTasks_EmptySection(t *testing.T) {
	input := `## 1. Empty Section

## 2. Another Section`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 2 {
		t.Fatalf("Expected 2 sections, got %d", len(sections))
	}

	if len(sections[0].Tasks) != 0 {
		t.Errorf("Expected empty tasks for first section, got %d", len(sections[0].Tasks))
	}
	if sections[0].Name != "Empty Section" {
		t.Errorf("Expected section name 'Empty Section', got %q", sections[0].Name)
	}
}

func TestParseTasks_BlankLines(t *testing.T) {
	input := `## 1. Tasks

- [ ] 1.1 First task

- [ ] 1.2 Second task

`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}
	if len(sections[0].Tasks) != 2 {
		t.Fatalf("Expected 2 tasks (blank lines ignored), got %d", len(sections[0].Tasks))
	}
}

func TestParseTasks_MultipleSections(t *testing.T) {
	input := `## 1. Implementation
- [ ] 1.1 First impl task
- [x] 1.2 Second impl task

## 2. Validation
- [ ] 2.1 Validate build
- [ ] 2.2 Run tests`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 2 {
		t.Fatalf("Expected 2 sections, got %d", len(sections))
	}

	// First section
	if sections[0].Number != 1 {
		t.Errorf("Expected first section number 1, got %d", sections[0].Number)
	}
	if sections[0].Name != "Implementation" {
		t.Errorf("Expected first section name 'Implementation', got %q", sections[0].Name)
	}
	if len(sections[0].Tasks) != 2 {
		t.Errorf("Expected 2 tasks in first section, got %d", len(sections[0].Tasks))
	}

	// Second section
	if sections[1].Number != 2 {
		t.Errorf("Expected second section number 2, got %d", sections[1].Number)
	}
	if sections[1].Name != "Validation" {
		t.Errorf("Expected second section name 'Validation', got %q", sections[1].Name)
	}
	if len(sections[1].Tasks) != 2 {
		t.Errorf("Expected 2 tasks in second section, got %d", len(sections[1].Tasks))
	}
}

func TestParseTasks_MixedNumberedAndPlainSections(t *testing.T) {
	input := `## 1. Numbered Section
- [ ] 1.1 Task

## Plain Section
- [ ] Task`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}
	if len(sections) != 2 {
		t.Fatalf("Expected 2 sections, got %d", len(sections))
	}

	if sections[0].Number != 1 {
		t.Errorf("Expected first section number 1, got %d", sections[0].Number)
	}
	if sections[1].Number != 1 {
		t.Errorf("Expected second section (plain) to auto-number to 1, got %d", sections[1].Number)
	}
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestParseTasks_FullExample(t *testing.T) {
	input := `## 1. Implementation
- [ ] 1.1 Create database schema
  - Parse requirement headers
  - Extract requirement name
- [x] 1.2 Implement API endpoint
- [ ] 1.3 Add tests
  - [ ] 1.3.1 Unit tests
  - [ ] 1.3.2 Integration tests

## 2. Validation
- [ ] 2.1 Run build
- [x] 2.2 Run tests`

	sections, err := ParseTasks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseTasks failed: %v", err)
	}

	// Verify section count
	if len(sections) != 2 {
		t.Fatalf("Expected 2 sections, got %d", len(sections))
	}

	// Section 1: Implementation
	section1 := sections[0]
	if section1.Number != 1 {
		t.Errorf("Section 1 number: expected 1, got %d", section1.Number)
	}
	if section1.Name != "Implementation" {
		t.Errorf("Section 1 name: expected 'Implementation', got %q", section1.Name)
	}
	if len(section1.Tasks) != 3 {
		t.Fatalf("Section 1: expected 3 top-level tasks, got %d", len(section1.Tasks))
	}

	// Task 1.1 with detail lines
	task11 := section1.Tasks[0]
	if task11.ID != "1.1" {
		t.Errorf("Task 1.1 ID: expected '1.1', got %q", task11.ID)
	}
	if !strings.HasPrefix(task11.Description, "Create database schema") {
		t.Errorf("Task 1.1 should start with 'Create database schema', got %q", task11.Description)
	}
	if !strings.Contains(task11.Description, "Parse requirement headers") {
		t.Error("Task 1.1 should contain detail line 'Parse requirement headers'")
	}
	if task11.Completed {
		t.Error("Task 1.1 should be incomplete")
	}

	// Task 1.2 completed
	task12 := section1.Tasks[1]
	if task12.ID != "1.2" {
		t.Errorf("Task 1.2 ID: expected '1.2', got %q", task12.ID)
	}
	if !task12.Completed {
		t.Error("Task 1.2 should be completed")
	}

	// Task 1.3 with subtasks
	task13 := section1.Tasks[2]
	if task13.ID != "1.3" {
		t.Errorf("Task 1.3 ID: expected '1.3', got %q", task13.ID)
	}
	if len(task13.Subtasks) != 2 {
		t.Fatalf("Task 1.3: expected 2 subtasks, got %d", len(task13.Subtasks))
	}
	if task13.Subtasks[0].ID != "1.3.1" {
		t.Errorf("Subtask ID: expected '1.3.1', got %q", task13.Subtasks[0].ID)
	}
	if task13.Subtasks[1].ID != "1.3.2" {
		t.Errorf("Subtask ID: expected '1.3.2', got %q", task13.Subtasks[1].ID)
	}

	// Section 2: Validation
	section2 := sections[1]
	if section2.Number != 2 {
		t.Errorf("Section 2 number: expected 2, got %d", section2.Number)
	}
	if section2.Name != "Validation" {
		t.Errorf("Section 2 name: expected 'Validation', got %q", section2.Name)
	}
	if len(section2.Tasks) != 2 {
		t.Fatalf("Section 2: expected 2 tasks, got %d", len(section2.Tasks))
	}

	// Task 2.1 incomplete
	if section2.Tasks[0].Completed {
		t.Error("Task 2.1 should be incomplete")
	}

	// Task 2.2 completed
	if !section2.Tasks[1].Completed {
		t.Error("Task 2.2 should be completed")
	}
}

// =============================================================================
// ParseTasksFile Tests
// =============================================================================

func TestParseTasksFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")

	content := `## 1. Test Section
- [ ] 1.1 Test task
- [x] 1.2 Completed task`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sections, err := ParseTasksFile(filePath)
	if err != nil {
		t.Fatalf("ParseTasksFile failed: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(sections))
	}
	if len(sections[0].Tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(sections[0].Tasks))
	}
}

func TestParseTasksFile_NotFound(t *testing.T) {
	_, err := ParseTasksFile("/nonexistent/path/tasks.md")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

// =============================================================================
// Helper Function Tests
// =============================================================================

func TestFindParentID(t *testing.T) {
	tests := []struct {
		taskID   string
		expected string
	}{
		{"1.2.3", "1.2"},
		{"1.2", "1"},
		{"1", ""},
		{"10.20.30", "10.20"},
		{"1.1.1.1", "1.1.1"},
	}

	for _, tt := range tests {
		t.Run(tt.taskID, func(t *testing.T) {
			result := findParentID(tt.taskID)
			if result != tt.expected {
				t.Errorf("findParentID(%q): expected %q, got %q", tt.taskID, tt.expected, result)
			}
		})
	}
}

func TestIsDetailLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"Two spaces indent", "  Some detail", true},
		{"Tab indent", "\tSome detail", true},
		{"Four spaces", "    More detail", true},
		{"No indent", "Not a detail", false},
		{"Task line indented", "  - [ ] Task", false},
		{"Task line with x indented", "  - [x] Task", false},
		{"Empty line", "", false},
		{"Single space", " Not enough indent", false},
		{"Bullet not task", "  - Regular bullet", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDetailLine(tt.line)
			if result != tt.expected {
				t.Errorf("isDetailLine(%q): expected %v, got %v", tt.line, tt.expected, result)
			}
		})
	}
}

func TestIntToString(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{10, "10"},
		{123, "123"},
		{-5, "-5"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := intToString(tt.input)
			if result != tt.expected {
				t.Errorf("intToString(%d): expected %q, got %q", tt.input, tt.expected, result)
			}
		})
	}
}

// =============================================================================
// Table-Driven Tests
// =============================================================================

func TestParseTaskLine_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		sectionNum int
		wantID     string
		wantDesc   string
		wantDone   bool
		wantNil    bool
	}{
		{
			name:       "Simple incomplete task",
			line:       "- [ ] 1.1 Create schema",
			sectionNum: 1,
			wantID:     "1.1",
			wantDesc:   "Create schema",
			wantDone:   false,
		},
		{
			name:       "Completed task lowercase",
			line:       "- [x] 2.1 Done task",
			sectionNum: 2,
			wantID:     "2.1",
			wantDesc:   "Done task",
			wantDone:   true,
		},
		{
			name:       "Completed task uppercase",
			line:       "- [X] 3.1 Also done",
			sectionNum: 3,
			wantID:     "3.1",
			wantDesc:   "Also done",
			wantDone:   true,
		},
		{
			name:       "Task without ID",
			line:       "- [ ] Task without ID",
			sectionNum: 4,
			wantID:     "4.1",
			wantDesc:   "Task without ID",
			wantDone:   false,
		},
		{
			name:       "Deep nested ID",
			line:       "- [ ] 1.2.3.4 Deeply nested",
			sectionNum: 1,
			wantID:     "1.2.3.4",
			wantDesc:   "Deeply nested",
			wantDone:   false,
		},
		{
			name:       "Indented task",
			line:       "  - [ ] 2.1.1 Indented subtask",
			sectionNum: 2,
			wantID:     "2.1.1",
			wantDesc:   "Indented subtask",
			wantDone:   false,
		},
		{
			name:    "Not a task - regular text",
			line:    "This is not a task",
			wantNil: true,
		},
		{
			name:    "Not a task - regular bullet",
			line:    "- Regular bullet",
			wantNil: true,
		},
		{
			name:    "Not a task - empty line",
			line:    "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var section *Section
			if !tt.wantNil {
				section = &Section{Number: tt.sectionNum}
			}
			autoTaskNum := 0

			task := parseTaskLine(tt.line, section, &autoTaskNum)

			if tt.wantNil {
				if task != nil {
					t.Errorf("Expected nil task, got %+v", task)
				}

				return
			}

			if task == nil {
				t.Fatal("Expected task, got nil")
			}
			if task.ID != tt.wantID {
				t.Errorf("ID: expected %q, got %q", tt.wantID, task.ID)
			}
			if task.Description != tt.wantDesc {
				t.Errorf("Description: expected %q, got %q", tt.wantDesc, task.Description)
			}
			if task.Completed != tt.wantDone {
				t.Errorf("Completed: expected %v, got %v", tt.wantDone, task.Completed)
			}
		})
	}
}

func TestParseSectionHeader_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantNum  int
		wantName string
		wantNil  bool
		autoIncr bool // whether auto number should increment
	}{
		{
			name:     "Numbered section",
			line:     "## 1. Implementation",
			wantNum:  1,
			wantName: "Implementation",
		},
		{
			name:     "Multi-digit numbered section",
			line:     "## 10. Large Section",
			wantNum:  10,
			wantName: "Large Section",
		},
		{
			name:     "Plain section",
			line:     "## Validation",
			wantNum:  1,
			wantName: "Validation",
			autoIncr: true,
		},
		{
			name:     "Section with extra whitespace",
			line:     "##   5.   Spaced Section   ",
			wantNum:  5,
			wantName: "Spaced Section",
		},
		{
			name:    "H1 heading",
			line:    "# Main Title",
			wantNil: true,
		},
		{
			name:    "H3 heading",
			line:    "### Subsection",
			wantNil: true,
		},
		{
			name:    "Regular text",
			line:    "Regular text",
			wantNil: true,
		},
		{
			name:    "Empty line",
			line:    "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			autoNum := 0
			section := parseSectionHeader(tt.line, &autoNum)

			if tt.wantNil {
				if section != nil {
					t.Errorf("Expected nil section, got %+v", section)
				}

				return
			}

			if section == nil {
				t.Fatal("Expected section, got nil")
			}
			if section.Number != tt.wantNum {
				t.Errorf("Number: expected %d, got %d", tt.wantNum, section.Number)
			}
			if section.Name != tt.wantName {
				t.Errorf("Name: expected %q, got %q", tt.wantName, section.Name)
			}
			if tt.autoIncr && autoNum != 1 {
				t.Errorf("Auto number should have incremented to 1, got %d", autoNum)
			}
			if !tt.autoIncr && autoNum != 0 {
				t.Errorf("Auto number should not have changed, got %d", autoNum)
			}
		})
	}
}
