package ralph

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// mockRalpher implements the Ralpher interface for testing.
type mockRalpher struct {
	binary       string
	invokeFunc   func(context.Context, *Task, string) (*exec.Cmd, error)
	shouldFail   bool
	failCount    int
	currentFails int
}

func (m *mockRalpher) Binary() string {
	return m.binary
}

func (m *mockRalpher) InvokeTask(
	ctx context.Context,
	task *Task,
	prompt string,
) (*exec.Cmd, error) {
	if m.invokeFunc != nil {
		return m.invokeFunc(ctx, task, prompt)
	}

	// Default: create a simple echo command that exits successfully
	if m.shouldFail && m.currentFails < m.failCount {
		m.currentFails++
		// Return a command that will fail
		cmd := exec.CommandContext(ctx, "false")

		return cmd, nil
	}

	// Return a command that will succeed
	cmd := exec.CommandContext(ctx, "true")

	return cmd, nil
}

var _ Ralpher = (*mockRalpher)(nil)

// setupTestOrchestrator creates a test orchestrator with a mock provider.
//
//nolint:revive // Intentionally modifies provider parameter to provide default for nil
func setupTestOrchestrator(t *testing.T, changeDir string, provider Ralpher) *Orchestrator {
	t.Helper()

	if provider == nil {
		newProvider := &mockRalpher{binary: "test-cli"}
		provider = newProvider
	}

	config := OrchestratorConfig{
		ChangeID:    "test-change",
		ChangeDir:   changeDir,
		Provider:    provider,
		MaxRetries:  2,
		TaskTimeout: 5 * time.Second,
	}

	orch, err := NewOrchestrator(&config)
	if err != nil {
		t.Fatalf("failed to create orchestrator: %v", err)
	}

	return orch
}

// setupTestChangeDir creates a temporary change directory with tasks.jsonc.
func setupTestChangeDir(t *testing.T, tasksJSON string) string {
	t.Helper()

	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Write tasks.jsonc
	tasksPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := os.WriteFile(tasksPath, []byte(tasksJSON), 0o644); err != nil {
		t.Fatalf("failed to write tasks.jsonc: %v", err)
	}

	// Write proposal.md (required by prompt generation)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	proposalContent := "# Test Proposal\n\nThis is a test proposal."
	if err := os.WriteFile(proposalPath, []byte(proposalContent), 0o644); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	return changeDir
}

func TestNewOrchestrator(t *testing.T) {
	tests := []struct {
		name    string
		config  OrchestratorConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: OrchestratorConfig{
				ChangeID:  "test-change",
				ChangeDir: "/tmp/test",
				Provider:  &mockRalpher{binary: "test"},
			},
			wantErr: false,
		},
		{
			name: "missing changeID",
			config: OrchestratorConfig{
				ChangeDir: "/tmp/test",
				Provider:  &mockRalpher{binary: "test"},
			},
			wantErr: true,
			errMsg:  "changeID cannot be empty",
		},
		{
			name: "missing changeDir",
			config: OrchestratorConfig{
				ChangeID: "test",
				Provider: &mockRalpher{binary: "test"},
			},
			wantErr: true,
			errMsg:  "changeDir cannot be empty",
		},
		{
			name: "missing provider",
			config: OrchestratorConfig{
				ChangeID:  "test",
				ChangeDir: "/tmp/test",
			},
			wantErr: true,
			errMsg:  "provider cannot be nil",
		},
		{
			name: "defaults applied",
			config: OrchestratorConfig{
				ChangeID:  "test",
				ChangeDir: "/tmp/test",
				Provider:  &mockRalpher{binary: "test"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orch, err := NewOrchestrator(&tt.config)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			// Check defaults
			if orch.maxRetries == 0 {
				t.Error("maxRetries not set to default")
			}
			if orch.taskTimeout == 0 {
				t.Error("taskTimeout not set to default")
			}
		})
	}
}

func TestOrchestrator_SimpleExecution(t *testing.T) {
	tasksJSON := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "First task",
				"status": "pending"
			},
			{
				"id": "1.2",
				"section": "Test",
				"description": "Second task",
				"status": "pending"
			}
		]
	}`

	changeDir := setupTestChangeDir(t, tasksJSON)

	// Mock provider that updates task status when invoked
	provider := &mockRalpher{
		binary: "test-cli",
		invokeFunc: func(ctx context.Context, task *Task, _ string) (*exec.Cmd, error) {
			// Update the task status to "completed" in tasks.jsonc
			go func() {
				time.Sleep(100 * time.Millisecond)
				updateTaskStatus(t, changeDir, task.ID, "completed")
			}()

			// Return a command that will succeed
			cmd := exec.CommandContext(ctx, "sleep", "0.2")

			return cmd, nil
		},
	}

	orch := setupTestOrchestrator(t, changeDir, provider)

	// Track task execution
	var executedTasks []string
	orch.onTaskStart = func(task *Task) {
		executedTasks = append(executedTasks, task.ID)
	}

	err := orch.Run()
	if err != nil {
		t.Fatalf("orchestration failed: %v", err)
	}

	// Verify all tasks were executed
	if len(executedTasks) != 2 {
		t.Errorf("expected 2 tasks executed, got %d", len(executedTasks))
	}

	// Verify session was cleaned up
	sessionPath := filepath.Join(changeDir, ".ralph-session.json")
	if _, err := os.Stat(sessionPath); !os.IsNotExist(err) {
		t.Error("session file should be deleted after successful completion")
	}
}

func TestOrchestrator_RetryOnFailure(t *testing.T) {
	tasksJSON := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Flaky task",
				"status": "pending"
			}
		]
	}`

	changeDir := setupTestChangeDir(t, tasksJSON)

	// Provider that fails twice, then succeeds
	failCount := 0
	provider := &mockRalpher{
		binary: "test-cli",
		invokeFunc: func(ctx context.Context, task *Task, _ string) (*exec.Cmd, error) {
			failCount++

			var cmd *exec.Cmd
			if failCount <= 2 {
				// Fail the first two attempts
				cmd = exec.CommandContext(ctx, "false")
			} else {
				// Succeed on third attempt
				go func() {
					time.Sleep(100 * time.Millisecond)
					updateTaskStatus(t, changeDir, task.ID, "completed")
				}()
				cmd = exec.CommandContext(ctx, "sleep", "0.2")
			}

			return cmd, nil
		},
	}

	orch := setupTestOrchestrator(t, changeDir, provider)
	orch.maxRetries = 3 // Allow 3 retries

	var failCount2 int
	orch.onTaskFail = func(_ *Task, _ *TaskResult) {
		failCount2++
	}

	err := orch.Run()
	if err != nil {
		t.Fatalf("orchestration failed: %v", err)
	}

	// Verify the task failed twice before succeeding
	if failCount2 != 2 {
		t.Errorf("expected 2 failures, got %d", failCount2)
	}
}

func TestOrchestrator_SessionPersistence(t *testing.T) {
	tasksJSON := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "First task",
				"status": "pending"
			},
			{
				"id": "1.2",
				"section": "Test",
				"description": "Second task",
				"status": "pending"
			}
		]
	}`

	changeDir := setupTestChangeDir(t, tasksJSON)

	// First orchestration: complete task 1.1, then abort
	taskCount := 0
	provider := &mockRalpher{
		binary: "test-cli",
		invokeFunc: func(ctx context.Context, task *Task, _ string) (*exec.Cmd, error) {
			taskCount++
			if taskCount == 1 {
				// Complete first task
				go func() {
					time.Sleep(100 * time.Millisecond)
					updateTaskStatus(t, changeDir, task.ID, "completed")
				}()
				cmd := exec.CommandContext(ctx, "sleep", "0.2")

				return cmd, nil
			}
			// Abort on second task
			return nil, errors.New("simulated abort")
		},
	}

	orch1 := setupTestOrchestrator(t, changeDir, provider)
	_ = orch1.Run() // Ignore error (expected to fail)

	// Verify session was saved
	sessionPath := filepath.Join(changeDir, ".ralph-session.json")
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		t.Fatal("session file should exist after interruption")
	}

	// Load session and verify state
	session, err := LoadSession(changeDir)
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}

	if len(session.CompletedIDs) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(session.CompletedIDs))
	}

	const testTaskOne = testTaskIDOne
	if session.CompletedIDs[0] != testTaskOne {
		t.Errorf("expected completed task 1.1, got %s", session.CompletedIDs[0])
	}
}

func TestOrchestrator_ParallelExecution(t *testing.T) {
	tasksJSON := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Task from tree 1",
				"status": "pending"
			},
			{
				"id": "2.1",
				"section": "Test",
				"description": "Task from tree 2",
				"status": "pending"
			}
		]
	}`

	changeDir := setupTestChangeDir(t, tasksJSON)

	// Track concurrent execution
	executing := make(map[string]bool)
	var mu sync.Mutex
	maxConcurrent := 0
	currentConcurrent := 0

	provider := &mockRalpher{
		binary: "test-cli",
		invokeFunc: func(ctx context.Context, tsk *Task, _ string) (*exec.Cmd, error) {
			// Track concurrency with proper synchronization
			t.Logf("Starting task %s", tsk.ID)
			mu.Lock()
			executing[tsk.ID] = true
			currentConcurrent++
			if currentConcurrent > maxConcurrent {
				maxConcurrent = currentConcurrent
			}
			mu.Unlock()

			// Update status during command execution, decrement after command completes
			taskID := tsk.ID
			go func() {
				// Update status partway through command execution
				time.Sleep(150 * time.Millisecond)
				updateTaskStatus(t, changeDir, taskID, "completed")
			}()

			// Decrement concurrency after command completes
			go func() {
				time.Sleep(350 * time.Millisecond) // After 300ms command + buffer
				mu.Lock()
				delete(executing, taskID)
				currentConcurrent--
				mu.Unlock()
			}()

			cmd := exec.CommandContext(ctx, "sleep", "0.3")

			return cmd, nil
		},
	}

	orch := setupTestOrchestrator(t, changeDir, provider)

	err := orch.Run()
	if err != nil {
		t.Fatalf("orchestration failed: %v", err)
	}

	// Verify tasks ran in parallel (maxConcurrent should be 2)
	if maxConcurrent < 2 {
		t.Errorf("expected parallel execution (max concurrent >= 2), got %d", maxConcurrent)
	}
}

func TestOrchestrator_UserActionAbort(t *testing.T) {
	tasksJSON := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Failing task",
				"status": "pending"
			}
		]
	}`

	changeDir := setupTestChangeDir(t, tasksJSON)

	// Provider that always fails
	provider := &mockRalpher{
		binary: "test-cli",
		invokeFunc: func(ctx context.Context, _ *Task, _ string) (*exec.Cmd, error) {
			cmd := exec.CommandContext(ctx, "false")

			return cmd, nil
		},
	}

	orch := setupTestOrchestrator(t, changeDir, provider)
	orch.maxRetries = 1

	// User chooses to abort
	orch.onUserAction = func(_ *Task, _ *TaskResult) UserAction {
		return UserActionAbort
	}

	err := orch.Run()
	if !errors.Is(err, ErrOrchestrationAborted) {
		t.Errorf("expected ErrOrchestrationAborted, got %v", err)
	}

	// Verify session was saved
	sessionPath := filepath.Join(changeDir, ".ralph-session.json")
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		t.Error("session file should exist after abort")
	}
}

func TestOrchestrator_UserActionSkip(t *testing.T) {
	tasksJSON := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "Failing task",
				"status": "pending"
			},
			{
				"id": "1.2",
				"section": "Test",
				"description": "Next task",
				"status": "pending"
			}
		]
	}`

	changeDir := setupTestChangeDir(t, tasksJSON)

	// Provider that fails first task, succeeds second
	provider := &mockRalpher{
		binary: "test-cli",
		invokeFunc: func(ctx context.Context, task *Task, _ string) (*exec.Cmd, error) {
			if task.ID == "1.1" {
				// Always fail first task (including retries)
				cmd := exec.CommandContext(ctx, "false")

				return cmd, nil
			}

			// Succeed second task
			go func() {
				time.Sleep(100 * time.Millisecond)
				updateTaskStatus(t, changeDir, task.ID, "completed")
			}()
			cmd := exec.CommandContext(ctx, "sleep", "0.2")

			return cmd, nil
		},
	}

	orch := setupTestOrchestrator(t, changeDir, provider)
	orch.maxRetries = 1

	// User chooses to skip failed tasks
	orch.onUserAction = func(_ *Task, _ *TaskResult) UserAction {
		return UserActionSkip
	}

	err := orch.Run()
	if err != nil {
		t.Fatalf("orchestration failed: %v", err)
	}

	// Verify first task is in failed list
	session := orch.GetSession()
	if len(session.FailedIDs) != 1 || session.FailedIDs[0] != "1.1" {
		t.Errorf("expected task 1.1 in failed list, got %v", session.FailedIDs)
	}

	// Verify second task completed
	if len(session.CompletedIDs) != 1 || session.CompletedIDs[0] != "1.2" {
		t.Errorf("expected task 1.2 in completed list, got %v", session.CompletedIDs)
	}
}

// updateTaskStatus updates a task status in tasks.jsonc for testing.
func updateTaskStatus(t *testing.T, changeDir, taskID, status string) {
	t.Helper()

	graph, err := ParseTaskGraph(changeDir)
	if err != nil {
		t.Logf("failed to parse task graph: %v", err)

		return
	}

	task, exists := graph.Tasks[taskID]
	if !exists {
		t.Logf("task %s not found", taskID)

		return
	}

	task.Status = status

	// Write back to tasks.jsonc
	// For simplicity, we'll use a naive approach
	// In a real implementation, you'd want to preserve formatting and comments
	tasksPath := filepath.Join(changeDir, "tasks.jsonc")
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Logf("failed to read tasks.jsonc: %v", err)

		return
	}

	// Simple string replacement (fragile, but okay for tests)
	oldStatus := `"status": "pending"`
	newStatus := `"status": "` + status + `"`

	// Find the task section and replace its status
	content := string(data)
	// This is a simplified approach - in production you'd want proper JSON manipulation
	updatedContent := bytes.ReplaceAll([]byte(content), []byte(oldStatus), []byte(newStatus))

	if err := os.WriteFile(tasksPath, updatedContent, 0o644); err != nil {
		t.Logf("failed to write tasks.jsonc: %v", err)
	}
}
