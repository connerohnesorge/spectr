package ralph

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Test constants to avoid goconst violations
const (
	testTaskIDOne    = "1.1"
	testStatusComp   = "completed"
	testStatusInProg = "in_progress"
	testStatusFail   = "failed"
)

func TestNewTUIModel(t *testing.T) {
	tasks := []*Task{
		{ID: testTaskIDOne, Section: "Setup", Description: "Init", Status: "pending"},
		{ID: "1.2", Section: "Setup", Description: "Config", Status: "pending"},
	}

	config := TUIConfig{
		ChangeID:      "test-change",
		Tasks:         tasks,
		Interactive:   false,
		InitialWidth:  80,
		InitialHeight: 24,
	}

	model := NewTUIModel(&config)

	if model.changeID != "test-change" {
		t.Errorf("changeID = %q, want %q", model.changeID, "test-change")
	}

	if len(model.tasks) != 2 {
		t.Errorf("len(tasks) = %d, want 2", len(model.tasks))
	}

	if model.mode != ModeNormal {
		t.Errorf("mode = %v, want ModeNormal", model.mode)
	}

	if model.width != 80 {
		t.Errorf("width = %d, want 80", model.width)
	}
}

func TestNewTUIModel_Interactive(t *testing.T) {
	config := TUIConfig{
		ChangeID:    "test",
		Tasks:       make([]*Task, 0),
		Interactive: true,
	}

	model := NewTUIModel(&config)

	if model.mode != ModeInteractive {
		t.Errorf("mode = %v, want ModeInteractive", model.mode)
	}
}

func TestTUIModel_Init(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init() returned non-nil cmd")
	}
}

func TestTUIModel_HandleKeyPress_Quit(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	// Test 'q' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := model.Update(msg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if !model.quitting {
		t.Error("Expected quitting to be true after 'q' key")
	}

	if cmd == nil {
		t.Error("Expected quit command")
	}
}

func TestTUIModel_HandleKeyPress_CtrlC(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	// Test Ctrl+C
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := model.Update(msg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if !model.quitting {
		t.Error("Expected quitting to be true after Ctrl+C")
	}

	if cmd == nil {
		t.Error("Expected quit command")
	}
}

func TestTUIModel_HandleKeyPress_Pause(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	// Test 'p' key to pause
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	updatedModel, _ := model.Update(msg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if !model.paused {
		t.Error("Expected paused to be true")
	}

	if model.mode != ModePaused {
		t.Errorf("mode = %v, want ModePaused", model.mode)
	}

	// Test 'p' again to unpause
	updatedModel, _ = model.Update(msg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.paused {
		t.Error("Expected paused to be false")
	}

	if model.mode != ModeNormal {
		t.Errorf("mode = %v, want ModeNormal", model.mode)
	}
}

func TestTUIModel_HandleKeyPress_Interactive(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	// Test 'i' key to enter interactive mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
	updatedModel, _ := model.Update(msg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.mode != ModeInteractive {
		t.Errorf("mode = %v, want ModeInteractive", model.mode)
	}

	if model.cursor != 0 {
		t.Errorf("cursor = %d, want 0", model.cursor)
	}
}

func TestTUIModel_InteractiveMode_Navigation(t *testing.T) {
	tasks := []*Task{
		{ID: testTaskIDOne, Section: "A", Description: "Task A", Status: "pending"},
		{ID: "1.2", Section: "B", Description: "Task B", Status: "pending"},
		{ID: "1.3", Section: "C", Description: "Task C", Status: "pending"},
	}

	model := NewTUIModel(&TUIConfig{
		ChangeID:    "test",
		Tasks:       tasks,
		Interactive: true,
	})

	// Test down navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(downMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.cursor != 1 {
		t.Errorf("cursor = %d, want 1", model.cursor)
	}

	// Test up navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = model.Update(upMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.cursor != 0 {
		t.Errorf("cursor = %d, want 0", model.cursor)
	}

	// Test 'j' (vim down)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ = model.Update(jMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.cursor != 1 {
		t.Errorf("cursor = %d, want 1", model.cursor)
	}

	// Test 'k' (vim up)
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedModel, _ = model.Update(kMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.cursor != 0 {
		t.Errorf("cursor = %d, want 0", model.cursor)
	}
}

func TestTUIModel_InteractiveMode_Selection(t *testing.T) {
	tasks := []*Task{
		{ID: testTaskIDOne, Section: "A", Description: "Task A", Status: "pending"},
		{ID: "1.2", Section: "B", Description: "Task B", Status: "pending"},
	}

	model := NewTUIModel(&TUIConfig{
		ChangeID:    "test",
		Tasks:       tasks,
		Interactive: true,
	})

	// Select first task
	spaceMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	updatedModel, _ := model.Update(spaceMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if !model.selectedTasks[testTaskIDOne] {
		t.Errorf("Expected task %s to be selected", testTaskIDOne)
	}

	// Toggle selection (deselect)
	updatedModel, _ = model.Update(spaceMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.selectedTasks[testTaskIDOne] {
		t.Errorf("Expected task %s to be deselected", testTaskIDOne)
	}
}

func TestTUIModel_InteractiveMode_ExitToNormal(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID:    "test",
		Tasks:       make([]*Task, 0),
		Interactive: true,
	})

	// Test 'esc' to exit
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(escMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.mode != ModeNormal {
		t.Errorf("mode = %v, want ModeNormal", model.mode)
	}

	// Re-enter interactive mode
	model.mode = ModeInteractive

	// Test 'q' to exit
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, _ = model.Update(qMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.mode != ModeNormal {
		t.Errorf("mode = %v, want ModeNormal", model.mode)
	}
}

func TestTUIModel_FailureMode_Retry(t *testing.T) {
	actionCalled := false
	var capturedAction UserAction

	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
		OnUserAction: func(action UserAction) tea.Cmd {
			actionCalled = true
			capturedAction = action

			return nil
		},
	})

	model.mode = ModeFailure

	// Test 'r' for retry
	rMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updatedModel, _ := model.Update(rMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if !actionCalled {
		t.Error("Expected onUserAction to be called")
	}

	if capturedAction != UserActionRetry {
		t.Errorf("action = %v, want UserActionRetry", capturedAction)
	}

	if model.mode != ModeNormal {
		t.Errorf("mode = %v, want ModeNormal", model.mode)
	}
}

func TestTUIModel_FailureMode_Skip(t *testing.T) {
	var capturedAction UserAction

	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
		OnUserAction: func(action UserAction) tea.Cmd {
			capturedAction = action

			return nil
		},
	})

	model.mode = ModeFailure

	// Test 's' for skip
	sMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updatedModel, _ := model.Update(sMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if capturedAction != UserActionSkip {
		t.Errorf("action = %v, want UserActionSkip", capturedAction)
	}

	if model.mode != ModeNormal {
		t.Errorf("mode = %v, want ModeNormal", model.mode)
	}
}

func TestTUIModel_FailureMode_Abort(t *testing.T) {
	var capturedAction UserAction

	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
		OnUserAction: func(action UserAction) tea.Cmd {
			capturedAction = action

			return tea.Quit
		},
	})

	model.mode = ModeFailure

	// Test 'q' for abort
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := model.Update(qMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if capturedAction != UserActionAbort {
		t.Errorf("action = %v, want UserActionAbort", capturedAction)
	}

	if !model.quitting {
		t.Error("Expected quitting to be true")
	}

	if cmd == nil {
		t.Error("Expected quit command")
	}
}

func TestTUIModel_TaskMessages(t *testing.T) {
	tasks := []*Task{
		{ID: testTaskIDOne, Section: "A", Description: "Task A", Status: "pending"},
	}

	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    tasks,
	})

	// Initialize viewport
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ := model.Update(sizeMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	// Test TaskStartMsg
	startMsg := TaskStartMsg{TaskID: testTaskIDOne}
	updatedModel, _ = model.Update(startMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.currentTaskID != testTaskIDOne {
		t.Errorf("currentTaskID = %q, want %q", model.currentTaskID, testTaskIDOne)
	}

	// Test TaskCompleteMsg
	completeMsg := TaskCompleteMsg{TaskID: testTaskIDOne, Success: true}
	updatedModel, _ = model.Update(completeMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.tasks[0].Status != testStatusComp {
		t.Errorf("task status = %q, want %q", model.tasks[0].Status, testStatusComp)
	}

	// Test TaskFailMsg
	model.tasks[0].Status = testStatusInProg
	failMsg := TaskFailMsg{TaskID: testTaskIDOne}
	updatedModel, _ = model.Update(failMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if model.tasks[0].Status != testStatusFail {
		t.Errorf("task status = %q, want %q", model.tasks[0].Status, testStatusFail)
	}

	if model.mode != ModeFailure {
		t.Errorf("mode = %v, want ModeFailure", model.mode)
	}
}

func TestTUIModel_PTYOutput(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	// Initialize viewport
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ := model.Update(sizeMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	// Add output
	outputMsg := PTYOutputMsg{Line: "Test output line"}
	updatedModel, _ = model.Update(outputMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	if len(model.outputLines) != 1 {
		t.Errorf("len(outputLines) = %d, want 1", len(model.outputLines))
	}

	if model.outputLines[0] != "Test output line" {
		t.Errorf("outputLines[0] = %q, want %q", model.outputLines[0], "Test output line")
	}
}

func TestTUIModel_View_NotReady(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	view := model.View()

	if view != "Initializing..." {
		t.Errorf("view = %q, want %q", view, "Initializing...")
	}
}

func TestTUIModel_View_Quitting(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	model.quitting = true
	view := model.View()

	if view != "" {
		t.Errorf("view = %q, want empty string", view)
	}
}

func TestTUIModel_View_Ready(t *testing.T) {
	tasks := []*Task{
		{ID: testTaskIDOne, Section: "Setup", Description: "Init", Status: testStatusComp},
	}

	model := NewTUIModel(&TUIConfig{
		ChangeID: "test-change",
		Tasks:    tasks,
	})

	// Initialize with window size
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ := model.Update(sizeMsg)
	if m, ok := updatedModel.(*TUIModel); ok {
		model = m
	} else {
		t.Fatal("Update() returned wrong type")
	}

	view := model.View()

	// Check for key elements
	if !strings.Contains(view, "test-change") {
		t.Error("View should contain change ID")
	}

	if !strings.Contains(view, "Tasks") {
		t.Error("View should contain 'Tasks' section")
	}

	if !strings.Contains(view, "Agent Output") {
		t.Error("View should contain 'Agent Output' section")
	}
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{testStatusComp, statusCompleted.String()},
		{testStatusInProg, statusInProgress.String()},
		{testStatusFail, statusFailed.String()},
		{"pending", statusPending.String()},
		{"unknown", statusPending.String()},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := getStatusIcon(tt.status)
			if got != tt.want {
				t.Errorf("getStatusIcon(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestFormatTaskLine(t *testing.T) {
	task := &Task{
		ID:      testTaskIDOne,
		Section: "Setup",
		Status:  testStatusComp,
	}

	// Test non-current task
	line := FormatTaskLine(task, false)
	if !strings.Contains(line, testTaskIDOne) {
		t.Error("Line should contain task ID")
	}
	if !strings.Contains(line, "Setup") {
		t.Error("Line should contain section")
	}

	// Test current task (should have highlighting)
	lineCurrent := FormatTaskLine(task, true)
	if !strings.Contains(lineCurrent, testTaskIDOne) {
		t.Error("Current line should contain task ID")
	}
}

func TestTUIModel_GetCompletedCount(t *testing.T) {
	tasks := []*Task{
		{ID: testTaskIDOne, Status: testStatusComp},
		{ID: "1.2", Status: testStatusComp},
		{ID: "1.3", Status: "pending"},
		{ID: "1.4", Status: testStatusInProg},
	}

	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    tasks,
	})

	count := model.GetCompletedCount()
	if count != 2 {
		t.Errorf("GetCompletedCount() = %d, want 2", count)
	}
}

func TestTUIModel_AddOutput(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	// Add output (thread-safe method)
	model.AddOutput("Line 1")
	model.AddOutput("Line 2")

	if len(model.outputLines) != 2 {
		t.Errorf("len(outputLines) = %d, want 2", len(model.outputLines))
	}

	if model.outputLines[0] != "Line 1" {
		t.Errorf("outputLines[0] = %q, want %q", model.outputLines[0], "Line 1")
	}

	if model.outputLines[1] != "Line 2" {
		t.Errorf("outputLines[1] = %q, want %q", model.outputLines[1], "Line 2")
	}
}

func TestTUIModel_SetCurrentTask(t *testing.T) {
	model := NewTUIModel(&TUIConfig{
		ChangeID: "test",
		Tasks:    nil,
	})

	model.SetCurrentTask(testTaskIDOne)

	if model.currentTaskID != testTaskIDOne {
		t.Errorf("currentTaskID = %q, want %q", model.currentTaskID, testTaskIDOne)
	}
}

func TestRenderTaskSummary(t *testing.T) {
	tasks := []*Task{
		{ID: testTaskIDOne, Status: testStatusComp},
		{ID: "1.2", Status: testStatusComp},
		{ID: "1.3", Status: testStatusFail},
		{ID: "1.4", Status: "pending"},
	}

	summary := RenderTaskSummary(tasks)

	if !strings.Contains(summary, "Completed: 2") {
		t.Error("Summary should contain completed count")
	}

	if !strings.Contains(summary, "Failed: 1") {
		t.Error("Summary should contain failed count")
	}

	if !strings.Contains(summary, "Pending: 1") {
		t.Error("Summary should contain pending count")
	}

	if !strings.Contains(summary, "Total: 4") {
		t.Error("Summary should contain total count")
	}
}

func TestMinInt(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{10, 10, 10},
		{-1, 5, -1},
	}

	for _, tt := range tests {
		got := min(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
