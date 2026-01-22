package ralph

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

const (
	// testTasksJSONTemplate is the template for test tasks.jsonc content
	testTasksJSONTemplate = `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Test task",
				"status": "pending"
			}
		]
	}`

	// testTaskID12 is the task ID "1.2" used in tests
	testTaskID12 = "1.2"

	// testStatusPend is the status "pending"
	testStatusPend = "pending"
)

// TestStatusWatcherCreation tests NewStatusWatcher constructor
func TestStatusWatcherCreation(t *testing.T) {
	paths := []string{"/path/to/tasks.jsonc"}
	interval := 2 * time.Second
	onChange := func(_ string, _ string) {
		// Callback for testing
	}

	watcher := NewStatusWatcher(paths, interval, onChange)

	if watcher == nil {
		t.Fatal("NewStatusWatcher returned nil")
	}

	if len(watcher.paths) != 1 || watcher.paths[0] != paths[0] {
		t.Errorf("paths not set correctly, got %v, want %v", watcher.paths, paths)
	}

	if watcher.interval != interval {
		t.Errorf("interval not set correctly, got %v, want %v", watcher.interval, interval)
	}

	if watcher.lastState == nil {
		t.Error("lastState map not initialized")
	}

	if watcher.stopChan == nil {
		t.Error("stopChan not initialized")
	}

	if watcher.doneChan == nil {
		t.Error("doneChan not initialized")
	}

	if watcher.running {
		t.Error("watcher should not be running initially")
	}
}

// TestStatusWatcherStartStop tests Start and Stop methods
func TestStatusWatcherStartStop(t *testing.T) {
	// Create a temp directory with a valid tasks.jsonc file
	tmpDir := t.TempDir()
	tasksFile := filepath.Join(tmpDir, "tasks.jsonc")

	if err := os.WriteFile(tasksFile, []byte(testTasksJSONTemplate), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file: %v", err)
	}

	watcher := NewStatusWatcher([]string{tasksFile}, 100*time.Millisecond, func(string, string) {})

	// Test starting the watcher
	if err := watcher.Start(); err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	// Verify running status
	watcher.mu.RLock()
	running := watcher.running
	watcher.mu.RUnlock()

	if !running {
		t.Error("Watcher should be running after Start()")
	}

	// Test double start
	if err := watcher.Start(); err == nil {
		t.Error("Start() should return error when already running")
	}

	// Test stopping the watcher
	if err := watcher.Stop(); err != nil {
		t.Fatalf("Failed to stop watcher: %v", err)
	}

	// Verify stopped status
	watcher.mu.RLock()
	running = watcher.running
	watcher.mu.RUnlock()

	if running {
		t.Error("Watcher should not be running after Stop()")
	}

	// Test double stop
	if err := watcher.Stop(); err == nil {
		t.Error("Stop() should return error when not running")
	}
}

// TestPollDetectsStatusChange tests that poll() detects status changes
func TestPollDetectsStatusChange(t *testing.T) {
	// Create a temp directory with a tasks.jsonc file
	tmpDir := t.TempDir()
	tasksFile := filepath.Join(tmpDir, "tasks.jsonc")

	initialContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Test task",
				"status": "pending"
			}
		]
	}`

	if err := os.WriteFile(tasksFile, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file: %v", err)
	}

	// Track onChange calls
	var mu sync.Mutex
	var changes []struct {
		taskID string
		status string
	}

	onChange := func(taskID string, status string) {
		mu.Lock()
		defer mu.Unlock()
		changes = append(changes, struct {
			taskID string
			status string
		}{taskID, status})
	}

	watcher := NewStatusWatcher([]string{tasksFile}, 100*time.Millisecond, onChange)

	// First poll - should initialize lastState but not call onChange
	if err := watcher.poll(); err != nil {
		t.Fatalf("First poll failed: %v", err)
	}

	// Verify initial state is cached
	watcher.mu.RLock()
	if len(watcher.lastState) != 1 {
		t.Errorf("Expected 1 task in lastState, got %d", len(watcher.lastState))
	}
	if watcher.lastState[testTaskIDOne] != testStatusPend {
		t.Errorf(
			"Expected task 1.1 status to be 'pending', got '%s'",
			watcher.lastState[testTaskIDOne],
		)
	}
	watcher.mu.RUnlock()

	// No changes should be detected yet
	mu.Lock()
	if len(changes) != 0 {
		t.Errorf("Expected 0 changes on initial poll, got %d", len(changes))
	}
	mu.Unlock()

	// Update the tasks file with a status change
	updatedContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Test task",
				"status": "in_progress"
			}
		]
	}`

	if err := os.WriteFile(tasksFile, []byte(updatedContent), 0o644); err != nil {
		t.Fatalf("Failed to update tasks file: %v", err)
	}

	// Second poll - should detect the change
	if err := watcher.poll(); err != nil {
		t.Fatalf("Second poll failed: %v", err)
	}

	// Verify onChange was called with correct parameters
	mu.Lock()
	if len(changes) != 1 {
		t.Fatalf("Expected 1 change detected, got %d", len(changes))
	}
	if changes[0].taskID != testTaskIDOne {
		t.Errorf("Expected taskID '%s', got '%s'", testTaskIDOne, changes[0].taskID)
	}
	if changes[0].status != testStatusInProg {
		t.Errorf("Expected status '%s', got '%s'", testStatusInProg, changes[0].status)
	}
	mu.Unlock()

	// Verify lastState was updated
	watcher.mu.RLock()
	if watcher.lastState["1.1"] != testStatusInProg {
		t.Errorf("Expected task 1.1 status to be 'in_progress', got '%s'", watcher.lastState["1.1"])
	}
	watcher.mu.RUnlock()
}

// TestPollMultipleChanges tests detection of multiple status changes in one poll
func TestPollMultipleChanges(t *testing.T) {
	tmpDir := t.TempDir()
	tasksFile := filepath.Join(tmpDir, "tasks.jsonc")

	initialContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Task 1",
				"status": "pending"
			},
			{
				"id": "1.2",
				"section": "Test",
				"description": "Task 2",
				"status": "pending"
			},
			{
				"id": "2.1",
				"section": "Test",
				"description": "Task 3",
				"status": "completed"
			}
		]
	}`

	if err := os.WriteFile(tasksFile, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file: %v", err)
	}

	var mu sync.Mutex
	var changes []struct {
		taskID string
		status string
	}

	onChange := func(taskID string, status string) {
		mu.Lock()
		defer mu.Unlock()
		changes = append(changes, struct {
			taskID string
			status string
		}{taskID, status})
	}

	watcher := NewStatusWatcher([]string{tasksFile}, 100*time.Millisecond, onChange)

	// First poll - initialize state
	if err := watcher.poll(); err != nil {
		t.Fatalf("First poll failed: %v", err)
	}

	// Update multiple tasks
	updatedContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Task 1",
				"status": "in_progress"
			},
			{
				"id": "1.2",
				"section": "Test",
				"description": "Task 2",
				"status": "completed"
			},
			{
				"id": "2.1",
				"section": "Test",
				"description": "Task 3",
				"status": "completed"
			}
		]
	}`

	if err := os.WriteFile(tasksFile, []byte(updatedContent), 0o644); err != nil {
		t.Fatalf("Failed to update tasks file: %v", err)
	}

	// Second poll - should detect multiple changes
	if err := watcher.poll(); err != nil {
		t.Fatalf("Second poll failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	// Should detect 2 changes (1.1 and 1.2, but not 2.1 which stayed completed)
	if len(changes) != 2 {
		t.Fatalf("Expected 2 changes detected, got %d: %+v", len(changes), changes)
	}

	// Check that both expected changes are present (order doesn't matter)
	foundTask11 := false
	foundTask12 := false

	for _, change := range changes {
		if change.taskID == testTaskIDOne && change.status == testStatusInProg {
			foundTask11 = true
		}
		if change.taskID == testTaskID12 && change.status == testStatusComp {
			foundTask12 = true
		}
	}

	if !foundTask11 {
		t.Error("Expected to find change for task 1.1 -> in_progress")
	}
	if !foundTask12 {
		t.Error("Expected to find change for task 1.2 -> completed")
	}
}

// TestPollNoFalsePositives tests that unchanged tasks don't trigger onChange
func TestPollNoFalsePositives(t *testing.T) {
	tmpDir := t.TempDir()
	tasksFile := filepath.Join(tmpDir, "tasks.jsonc")

	if err := os.WriteFile(tasksFile, []byte(testTasksJSONTemplate), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file: %v", err)
	}

	var mu sync.Mutex
	changeCount := 0

	onChange := func(_ string, _ string) {
		mu.Lock()
		defer mu.Unlock()
		changeCount++
	}

	watcher := NewStatusWatcher([]string{tasksFile}, 100*time.Millisecond, onChange)

	// First poll - initialize state
	if err := watcher.poll(); err != nil {
		t.Fatalf("First poll failed: %v", err)
	}

	// Second poll - same content, should detect no changes
	if err := watcher.poll(); err != nil {
		t.Fatalf("Second poll failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if changeCount != 0 {
		t.Errorf("Expected 0 changes for unchanged tasks, got %d", changeCount)
	}
}

// TestPollEmptyPaths tests poll() with no paths to monitor
func TestPollEmptyPaths(t *testing.T) {
	watcher := NewStatusWatcher(nil, 100*time.Millisecond, func(string, string) {})

	// poll() should return nil without error when no paths are configured
	if err := watcher.poll(); err != nil {
		t.Errorf("poll() with empty paths should return nil, got: %v", err)
	}
}

// TestPollNonexistentFile tests poll() gracefully handles missing files
func TestPollNonexistentFile(t *testing.T) {
	watcher := NewStatusWatcher(
		[]string{"/nonexistent/tasks.jsonc"},
		100*time.Millisecond,
		func(string, string) {},
	)

	// poll() should return error but not crash
	err := watcher.poll()
	if err == nil {
		t.Error("poll() should return error for nonexistent file")
	}
}

// TestPollInvalidJSON tests poll() gracefully handles malformed JSON
func TestPollInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tasksFile := filepath.Join(tmpDir, "tasks.jsonc")

	invalidContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1"
				// Missing comma - invalid JSON
				"status": "pending"
			}
		]
	}`

	if err := os.WriteFile(tasksFile, []byte(invalidContent), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file: %v", err)
	}

	watcher := NewStatusWatcher([]string{tasksFile}, 100*time.Millisecond, func(string, string) {})

	// poll() should return error for malformed JSON
	err := watcher.poll()
	if err == nil {
		t.Error("poll() should return error for invalid JSON")
	}
}

// TestWatcherIntegration tests the full watcher lifecycle with Start/Stop
func TestWatcherIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	tasksFile := filepath.Join(tmpDir, "tasks.jsonc")

	initialContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Test task",
				"status": "pending"
			}
		]
	}`

	if err := os.WriteFile(tasksFile, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file: %v", err)
	}

	var mu sync.Mutex
	var changes []struct {
		taskID string
		status string
	}

	onChange := func(taskID string, status string) {
		mu.Lock()
		defer mu.Unlock()
		changes = append(changes, struct {
			taskID string
			status string
		}{taskID, status})
	}

	// Create watcher with short interval for faster testing
	watcher := NewStatusWatcher([]string{tasksFile}, 50*time.Millisecond, onChange)

	// Start the watcher
	if err := watcher.Start(); err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	// Give the watcher time to do the initial poll
	time.Sleep(100 * time.Millisecond)

	// Update the file
	updatedContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Test task",
				"status": "completed"
			}
		]
	}`

	if err := os.WriteFile(tasksFile, []byte(updatedContent), 0o644); err != nil {
		t.Fatalf("Failed to update tasks file: %v", err)
	}

	// Wait for watcher to detect the change
	time.Sleep(150 * time.Millisecond)

	// Stop the watcher
	if err := watcher.Stop(); err != nil {
		t.Fatalf("Failed to stop watcher: %v", err)
	}

	// Verify the change was detected
	mu.Lock()
	defer mu.Unlock()

	if len(changes) == 0 {
		t.Fatal("Expected at least 1 change to be detected")
	}

	// Find the change for task 1.1
	found := false
	for _, change := range changes {
		if change.taskID == testTaskIDOne && change.status == testStatusComp {
			found = true

			break
		}
	}

	if !found {
		t.Errorf("Expected change for task 1.1 -> completed not found in %+v", changes)
	}
}

// TestPollMultipleFiles tests polling with multiple tasks*.jsonc files
func TestPollMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	tasksFile1 := filepath.Join(tmpDir, "tasks.jsonc")
	tasksFile2 := filepath.Join(tmpDir, "tasks-2.jsonc")

	content1 := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Task 1",
				"status": "pending"
			}
		]
	}`

	content2 := `{
		"version": 1,
		"tasks": [
			{
				"id": "2.1",
				"section": "Test",
				"description": "Task 2",
				"status": "pending"
			}
		]
	}`

	if err := os.WriteFile(tasksFile1, []byte(content1), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file 1: %v", err)
	}
	if err := os.WriteFile(tasksFile2, []byte(content2), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file 2: %v", err)
	}

	var mu sync.Mutex
	var changes []struct {
		taskID string
		status string
	}

	onChange := func(taskID string, status string) {
		mu.Lock()
		defer mu.Unlock()
		changes = append(changes, struct {
			taskID string
			status string
		}{taskID, status})
	}

	// Note: ParseTaskGraph expects all tasks*.jsonc files in the same directory,
	// so we pass just one path and it will discover all files via glob
	watcher := NewStatusWatcher([]string{tasksFile1}, 100*time.Millisecond, onChange)

	// First poll - initialize state
	if err := watcher.poll(); err != nil {
		t.Fatalf("First poll failed: %v", err)
	}

	// Verify both tasks are in lastState
	watcher.mu.RLock()
	if len(watcher.lastState) != 2 {
		t.Errorf("Expected 2 tasks in lastState, got %d", len(watcher.lastState))
	}
	watcher.mu.RUnlock()

	// Update task in second file
	updatedContent2 := `{
		"version": 1,
		"tasks": [
			{
				"id": "2.1",
				"section": "Test",
				"description": "Task 2",
				"status": "completed"
			}
		]
	}`

	if err := os.WriteFile(tasksFile2, []byte(updatedContent2), 0o644); err != nil {
		t.Fatalf("Failed to update tasks file 2: %v", err)
	}

	// Second poll - should detect change in second file
	if err := watcher.poll(); err != nil {
		t.Fatalf("Second poll failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change detected, got %d", len(changes))
	}
	if changes[0].taskID != "2.1" || changes[0].status != testStatusComp {
		t.Errorf("Expected change for task 2.1 -> completed, got %+v", changes[0])
	}
}
