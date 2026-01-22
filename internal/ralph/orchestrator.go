// Package ralph provides task orchestration for Spectr change proposals.
// It automates the execution of tasks from tasks.jsonc files by coordinating
// agent CLI sessions with dependency-aware parallel execution.
package ralph

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// orchestrator.go contains the main orchestration loop that coordinates
// task execution across the dependency graph. It integrates the task graph,
// prompt generation, PTY subprocess management, and status file watching
// to drive automated task completion.
//
// The orchestrator handles:
// - Sequential execution of dependent tasks
// - Parallel execution of independent tasks
// - Retry logic with configurable limits
// - User interaction for retry/skip/abort decisions
// - Session state persistence for resumable workflows

const (
	// DefaultMaxRetries is the default number of times to retry a failed task.
	DefaultMaxRetries = 3

	// DefaultTaskTimeout is the maximum time to wait for a task to complete.
	DefaultTaskTimeout = 30 * time.Minute

	// DefaultStatusPollInterval is how often to check tasks.jsonc for status changes.
	DefaultStatusPollInterval = 2 * time.Second
)

var (
	// ErrTaskFailed indicates a task failed and all retries were exhausted.
	ErrTaskFailed = errors.New("task failed after maximum retries")

	// ErrOrchestrationAborted indicates the user aborted the orchestration.
	ErrOrchestrationAborted = errors.New("orchestration aborted by user")

	// ErrTaskTimeout indicates a task exceeded its timeout duration.
	ErrTaskTimeout = errors.New("task execution timeout exceeded")

	// ErrNoTasksToExecute indicates there are no pending tasks to run.
	ErrNoTasksToExecute = errors.New("no tasks to execute")
)

// UserAction represents actions the user can take when a task fails.
type UserAction string

const (
	// UserActionRetry retries the failed task.
	UserActionRetry UserAction = "retry"

	// UserActionSkip skips the failed task and continues to the next.
	UserActionSkip UserAction = "skip"

	// UserActionAbort aborts the entire orchestration.
	UserActionAbort UserAction = "abort"
)

// TaskResult represents the outcome of a task execution.
type TaskResult struct {
	TaskID   string
	Success  bool
	ExitCode int
	Duration time.Duration
	Error    error
	Output   string
	Retries  int
}

// Orchestrator coordinates the execution of tasks from a change proposal.
// It manages the lifecycle of task execution, including:
//   - Loading or resuming sessions
//   - Executing tasks in topological order with parallel execution
//   - Monitoring task status via file watching
//   - Handling failures with retry logic
//   - Persisting session state for resume capability
type Orchestrator struct {
	// changeID is the unique identifier for the change proposal.
	changeID string

	// changeDir is the path to the change directory (e.g., spectr/changes/<id>).
	changeDir string

	// provider is the Ralpher implementation used to invoke agent CLIs.
	provider Ralpher

	// maxRetries is the maximum number of times to retry a failed task.
	maxRetries int

	// taskTimeout is the maximum time to wait for a task to complete.
	taskTimeout time.Duration

	// graph is the parsed task dependency graph.
	graph *TaskGraph

	// session tracks orchestration progress for resume capability.
	session *SessionState

	// watcher monitors tasks.jsonc files for status changes.
	watcher *StatusWatcher

	// ctx is the context for cancellation and timeout control.
	ctx context.Context

	// cancel allows cancelling the orchestration.
	cancel context.CancelFunc

	// onUserAction is a callback for handling user actions on failure.
	// If nil, defaults to abort on failure.
	onUserAction func(task *Task, result *TaskResult) UserAction

	// onTaskStart is called when a task starts execution.
	onTaskStart func(task *Task)

	// onTaskComplete is called when a task completes successfully.
	onTaskComplete func(task *Task, result *TaskResult)

	// onTaskFail is called when a task fails.
	onTaskFail func(task *Task, result *TaskResult)

	// mu protects concurrent access to session state.
	mu sync.Mutex
}

// OrchestratorConfig holds configuration for creating an Orchestrator.
type OrchestratorConfig struct {
	ChangeID       string
	ChangeDir      string
	Provider       Ralpher
	MaxRetries     int
	TaskTimeout    time.Duration
	OnUserAction   func(task *Task, result *TaskResult) UserAction
	OnTaskStart    func(task *Task)
	OnTaskComplete func(task *Task, result *TaskResult)
	OnTaskFail     func(task *Task, result *TaskResult)
}

// NewOrchestrator creates a new Orchestrator with the given configuration.
//
// The orchestrator is initialized but not started. Call Run() to begin execution.
//
// Parameters:
//   - config: Configuration for the orchestrator
//
// Returns:
//   - *Orchestrator: The configured orchestrator
//   - error: Any error encountered during initialization
func NewOrchestrator(config *OrchestratorConfig) (*Orchestrator, error) {
	if config.ChangeID == "" {
		return nil, errors.New("changeID cannot be empty")
	}
	if config.ChangeDir == "" {
		return nil, errors.New("changeDir cannot be empty")
	}
	if config.Provider == nil {
		return nil, errors.New("provider cannot be nil")
	}

	// Set defaults
	if config.MaxRetries == 0 {
		config.MaxRetries = DefaultMaxRetries
	}
	if config.TaskTimeout == 0 {
		config.TaskTimeout = DefaultTaskTimeout
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Orchestrator{
		changeID:       config.ChangeID,
		changeDir:      config.ChangeDir,
		provider:       config.Provider,
		maxRetries:     config.MaxRetries,
		taskTimeout:    config.TaskTimeout,
		ctx:            ctx,
		cancel:         cancel,
		onUserAction:   config.OnUserAction,
		onTaskStart:    config.OnTaskStart,
		onTaskComplete: config.OnTaskComplete,
		onTaskFail:     config.OnTaskFail,
	}, nil
}

// Run starts the orchestration process.
//
// This method:
//  1. Loads or resumes the session
//  2. Parses the task graph
//  3. Executes tasks in topological order
//  4. Handles retries and failures
//  5. Saves session state on interruption
//  6. Cleans up on completion
//
// Returns:
//   - error: Any error that caused orchestration to fail or abort
func (o *Orchestrator) Run() error {
	// Load task graph
	graph, err := ParseTaskGraph(o.changeDir)
	if err != nil {
		return fmt.Errorf("failed to parse task graph: %w", err)
	}
	o.graph = graph

	// Load or create session
	if err := o.loadOrCreateSession(); err != nil {
		return fmt.Errorf("failed to initialize session: %w", err)
	}

	// Initialize status watcher
	if err := o.initializeWatcher(); err != nil {
		return fmt.Errorf("failed to initialize watcher: %w", err)
	}
	defer o.stopWatcher()

	// Get topological execution order
	stages, err := o.graph.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to compute execution order: %w", err)
	}

	// Execute stages
	for stageIdx, stage := range stages {
		// Filter out already completed or failed tasks
		pendingTasks := o.filterPendingTasks(stage)
		if len(pendingTasks) == 0 {
			continue // All tasks in this stage are done
		}

		// Execute tasks in this stage (potentially in parallel)
		if err := o.executeStage(stageIdx, pendingTasks); err != nil {
			// Save session before returning
			_ = o.saveSession()

			return err
		}
	}

	// All tasks completed successfully
	return o.cleanup()
}

// Stop gracefully stops the orchestration and saves session state.
func (o *Orchestrator) Stop() error {
	o.cancel()

	return o.saveSession()
}

// loadOrCreateSession loads an existing session or creates a new one.
func (o *Orchestrator) loadOrCreateSession() error {
	session, err := LoadSession(o.changeDir)
	if err != nil {
		if os.IsNotExist(err) {
			// No existing session, create new one
			o.session = &SessionState{
				ChangeID:     o.changeID,
				StartedAt:    time.Now(),
				LastUpdated:  time.Now(),
				CompletedIDs: nil,
				FailedIDs:    nil,
				RetryCount:   make(map[string]int),
			}

			return nil
		}

		return fmt.Errorf("failed to load session: %w", err)
	}

	o.session = session

	return nil
}

// initializeWatcher creates and starts the status watcher.
func (o *Orchestrator) initializeWatcher() error {
	// Find all tasks.jsonc files
	tasksPattern := o.changeDir + "/tasks*.jsonc"
	matches, err := findTasksFiles(o.changeDir)
	if err != nil {
		return fmt.Errorf("failed to find tasks files: %w", err)
	}

	_ = tasksPattern // Avoid unused warning

	// Create watcher
	o.watcher = NewStatusWatcher(
		matches,
		DefaultStatusPollInterval,
		o.onStatusChange,
	)

	// Start watching
	return o.watcher.Start()
}

// findTasksFiles finds all tasks*.jsonc files in the change directory.
func findTasksFiles(changeDir string) ([]string, error) {
	pattern := changeDir + "/tasks*.jsonc"
	matches, err := findFiles(pattern)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, errors.New("no tasks*.jsonc files found")
	}

	return matches, nil
}

// findFiles is a wrapper around filepath.Glob for testing.
var findFiles = filepath.Glob

// stopWatcher stops the status watcher if it's running.
func (o *Orchestrator) stopWatcher() {
	if o.watcher != nil {
		_ = o.watcher.Stop()
	}
}

// onStatusChange is called when a task status changes in tasks.jsonc.
func (*Orchestrator) onStatusChange(_, _ string) {
	// This callback is used by the watcher to notify us of status changes
	// We don't need to do anything here since we poll for completion in waitForCompletion
}

// filterPendingTasks filters a stage to only include tasks that need execution.
func (o *Orchestrator) filterPendingTasks(stage []string) []string {
	o.mu.Lock()
	defer o.mu.Unlock()

	pending := make([]string, 0, len(stage))
	for _, taskID := range stage {
		if !o.session.IsCompleted(taskID) && !o.session.IsFailed(taskID) {
			pending = append(pending, taskID)
		}
	}

	return pending
}

// executeStage executes all tasks in a stage.
// Tasks in the same stage can potentially run in parallel if they're independent.
func (o *Orchestrator) executeStage(_ int, taskIDs []string) error {
	// Identify parallel groups within this stage
	parallelGroups := o.identifyParallelGroups(taskIDs)

	// Execute each parallel group
	for _, group := range parallelGroups {
		if err := o.executeParallelGroup(group); err != nil {
			return err
		}
	}

	return nil
}

// identifyParallelGroups groups tasks by their root prefix for parallel execution.
// Tasks with different root prefixes can run in parallel.
func (o *Orchestrator) identifyParallelGroups(taskIDs []string) [][]string {
	// Group tasks by root prefix
	groups := make(map[string][]string)
	for _, taskID := range taskIDs {
		root := o.graph.GetRootPrefix(taskID)
		groups[root] = append(groups[root], taskID)
	}

	// Convert map to slice
	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// executeParallelGroup executes a group of tasks that can run in parallel.
func (o *Orchestrator) executeParallelGroup(taskIDs []string) error {
	if len(taskIDs) == 1 {
		// Single task, execute directly
		return o.executeTaskWithRetry(taskIDs[0])
	}

	// Multiple tasks that can run in parallel
	var wg sync.WaitGroup
	errChan := make(chan error, len(taskIDs))

	for _, taskID := range taskIDs {
		wg.Add(1)
		go func(tid string) {
			defer wg.Done()
			if err := o.executeTaskWithRetry(tid); err != nil {
				errChan <- fmt.Errorf("task %s: %w", tid, err)
			}
		}(taskID)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// executeTaskWithRetry executes a task with retry logic.
func (o *Orchestrator) executeTaskWithRetry(taskID string) error {
	task := o.graph.Tasks[taskID]
	if task == nil {
		return fmt.Errorf("task %s not found in graph", taskID)
	}

	// Mark task as current
	o.mu.Lock()
	o.session.CurrentTaskID = taskID
	o.mu.Unlock()

	for {
		// Execute the task
		result, err := o.executeTask(task)

		if err == nil && result.Success {
			// Task succeeded
			o.handleTaskSuccess(task, result)

			return nil
		}

		// Task failed, handle failure
		action := o.handleTaskFailure(task, result)

		switch action {
		case UserActionRetry:
			// Retry the task
			continue

		case UserActionSkip:
			// Skip this task and continue
			o.mu.Lock()
			o.session.MarkTaskFailed(taskID)
			_ = o.saveSessionLocked()
			o.mu.Unlock()

			return nil

		case UserActionAbort:
			// Abort orchestration
			return ErrOrchestrationAborted

		default:
			// Unknown action, abort
			return fmt.Errorf("unknown user action: %s", action)
		}
	}
}

// executeTask executes a single task and waits for completion.
func (o *Orchestrator) executeTask(task *Task) (*TaskResult, error) {
	startTime := time.Now()

	// Notify start
	if o.onTaskStart != nil {
		o.onTaskStart(task)
	}

	// Generate prompt
	prompt, err := GeneratePrompt(task, o.changeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to generate prompt: %w", err)
	}

	// Create task context with timeout
	taskCtx, taskCancel := context.WithTimeout(o.ctx, o.taskTimeout)
	defer taskCancel()

	// Invoke the provider to get the command
	cmd, err := o.provider.InvokeTask(taskCtx, task, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke provider: %w", err)
	}

	// Start PTY session
	session, err := StartPTYSession(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to start PTY session: %w", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Wait for completion with timeout
	exitCode := 0
	waitErr := session.Wait()
	if waitErr != nil {
		exitCode = 1
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	// Check if context was cancelled (timeout)
	select {
	case <-taskCtx.Done():
		_ = session.Terminate(5 * time.Second)

		return &TaskResult{
			TaskID:   task.ID,
			Success:  false,
			ExitCode: exitCode,
			Duration: time.Since(startTime),
			Error:    ErrTaskTimeout,
		}, ErrTaskTimeout
	default:
	}

	// Determine if task succeeded
	success := exitCode == 0 && o.checkTaskStatusCompleted(task.ID)

	result := &TaskResult{
		TaskID:   task.ID,
		Success:  success,
		ExitCode: exitCode,
		Duration: time.Since(startTime),
		Error:    waitErr,
	}

	return result, nil
}

// checkTaskStatusCompleted checks if the task status in tasks.jsonc is "completed".
func (o *Orchestrator) checkTaskStatusCompleted(taskID string) bool {
	// Re-parse the task graph to get the latest status
	graph, err := ParseTaskGraph(o.changeDir)
	if err != nil {
		return false
	}

	task, exists := graph.Tasks[taskID]
	if !exists {
		return false
	}

	return task.Status == "completed"
}

// handleTaskSuccess handles successful task completion.
func (o *Orchestrator) handleTaskSuccess(task *Task, result *TaskResult) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.session.MarkTaskCompleted(task.ID)
	_ = o.saveSessionLocked()

	if o.onTaskComplete != nil {
		o.onTaskComplete(task, result)
	}
}

// handleTaskFailure handles task failure and determines the next action.
func (o *Orchestrator) handleTaskFailure(task *Task, result *TaskResult) UserAction {
	o.mu.Lock()
	retryCount := o.session.GetRetryCount(task.ID)
	o.session.IncrementRetry(task.ID)
	_ = o.saveSessionLocked()
	o.mu.Unlock()

	if o.onTaskFail != nil {
		o.onTaskFail(task, result)
	}

	// Check if we've exhausted retries
	if retryCount >= o.maxRetries {
		// Ask user what to do
		if o.onUserAction != nil {
			return o.onUserAction(task, result)
		}
		// Default to abort if no user action handler
		return UserActionAbort
	}

	// Auto-retry if we haven't hit the limit
	return UserActionRetry
}

// saveSession saves the current session state to disk.
// This method acquires the mutex before saving.
func (o *Orchestrator) saveSession() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.saveSessionLocked()
}

// saveSessionLocked saves the current session state to disk without acquiring the mutex.
// The caller must hold the mutex before calling this method.
func (o *Orchestrator) saveSessionLocked() error {
	if o.session == nil {
		return nil
	}

	return o.session.Save(o.changeDir)
}

// cleanup performs cleanup after successful completion.
func (o *Orchestrator) cleanup() error {
	// Delete session file
	if err := DeleteSession(o.changeDir); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// GetSession returns the current session state.
// This is useful for testing and monitoring.
func (o *Orchestrator) GetSession() *SessionState {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.session
}
