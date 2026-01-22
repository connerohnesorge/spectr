// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"errors"
	"path/filepath"
	"sync"
	"time"
)

// watcher.go implements status file watching for task completion detection.
// The StatusWatcher monitors tasks.jsonc (or tasks-*.jsonc split files)
// for changes to task status fields.
//
// Features:
// - Polling-based file watching with configurable interval (default 2s)
// - Detection of status transitions (pending -> in_progress -> completed)
// - Support for split task files (tasks-*.jsonc glob pattern)
// - Event emission for task status changes
//
// The watcher enables the orchestrator to detect when agents have
// completed tasks and automatically proceed to the next task.

// StatusWatcher polls tasks.jsonc files for status changes and notifies
// the orchestrator when tasks transition between states.
type StatusWatcher struct {
	// paths are the file paths to monitor (e.g., tasks.jsonc, tasks-1.jsonc)
	paths []string

	// interval is the polling interval (e.g., 2s)
	interval time.Duration

	// onChange is the callback function called when a task status changes.
	// Parameters: taskID, newStatus
	onChange func(taskID string, newStatus string)

	// lastState caches the most recent status for each task ID
	lastState map[string]string

	// stopChan is used to signal the watcher to stop polling
	stopChan chan struct{}

	// doneChan signals when the watcher has fully stopped
	doneChan chan struct{}

	// mu protects concurrent access to lastState
	mu sync.RWMutex

	// running indicates whether the watcher is currently active
	running bool
}

// NewStatusWatcher creates a new StatusWatcher that monitors the given file paths.
//
// Parameters:
//   - paths: List of tasks.jsonc file paths to monitor
//   - interval: How often to poll the files (e.g., 2*time.Second)
//   - onChange: Callback function invoked when task status changes
//
// The watcher does not start automatically. Call Start() to begin polling.
func NewStatusWatcher(
	paths []string,
	interval time.Duration,
	onChange func(string, string),
) *StatusWatcher {
	return &StatusWatcher{
		paths:     paths,
		interval:  interval,
		onChange:  onChange,
		lastState: make(map[string]string),
		stopChan:  make(chan struct{}),
		doneChan:  make(chan struct{}),
		running:   false,
	}
}

// Start begins polling the tasks.jsonc files in a background goroutine.
// Returns an error if the watcher is already running.
func (w *StatusWatcher) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return errors.New("watcher is already running")
	}

	w.running = true

	// Start the polling goroutine
	go w.pollLoop()

	return nil
}

// Stop gracefully stops the watcher and waits for the polling goroutine to exit.
func (w *StatusWatcher) Stop() error {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()

		return errors.New("watcher is not running")
	}
	w.mu.Unlock()

	// Signal stop
	close(w.stopChan)

	// Wait for the polling goroutine to finish
	<-w.doneChan

	w.mu.Lock()
	w.running = false
	w.mu.Unlock()

	return nil
}

// pollLoop is the main polling loop that runs in a background goroutine.
func (w *StatusWatcher) pollLoop() {
	defer close(w.doneChan)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Do an initial poll immediately
	_ = w.poll()

	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			_ = w.poll()
		}
	}
}

// poll reads all task files and detects status changes.
// It compares the current state with lastState and calls onChange for each change.
func (w *StatusWatcher) poll() error {
	// If no paths to monitor, return early
	if len(w.paths) == 0 {
		return nil
	}

	// ParseTaskGraph expects a directory, so we'll use the directory
	// containing the first path. All tasks*.jsonc files should be in the same directory.
	changeDir := filepath.Dir(w.paths[0])

	// Parse all tasks from the directory
	graph, err := ParseTaskGraph(changeDir)
	if err != nil {
		// If we can't parse, just return the error but don't crash
		// This handles cases where files don't exist yet or are temporarily invalid
		return err
	}

	// Collect all current task statuses
	allTasks := make(map[string]string) // taskID -> status
	for taskID, task := range graph.Tasks {
		allTasks[taskID] = task.Status
	}

	// Detect changes
	w.mu.Lock()
	defer w.mu.Unlock()

	for taskID, newStatus := range allTasks {
		oldStatus, exists := w.lastState[taskID]

		// If the status changed, call the onChange callback
		if !exists || oldStatus != newStatus {
			// Only notify if the status actually changed (not initial discovery)
			if exists && w.onChange != nil {
				w.onChange(taskID, newStatus)
			}

			// Update the cache
			w.lastState[taskID] = newStatus
		}
	}

	return nil
}
