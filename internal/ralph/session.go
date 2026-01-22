// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// sessionFilePermissions is the file permission mode for session files (rw-r--r--)
const sessionFilePermissions = 0o644

// session.go manages persistent session state for resumable workflows.
// It handles:
// - Saving session state on interruption or quit
// - Loading saved sessions on startup
// - Prompting users to resume or restart
// - Cleaning up session files on completion
//
// Session state includes:
// - Current task index
// - Completed tasks list
// - Failed tasks with retry counts
// - Timestamp and change ID
//
// Sessions are stored at: spectr/changes/<change-id>/.ralph-session.json

// SessionState persists orchestration progress for resume after interruption.
// The session state enables users to resume orchestration after Ctrl+C, crashes,
// or other interruptions without losing progress. It tracks which tasks have
// been completed, which have failed, and how many retries have been attempted.
//
// Session files are stored at: spectr/changes/<change-id>/.ralph-session.json
//
// The session state is:
//   - Saved automatically on interruption (Ctrl+C) or quit
//   - Loaded on startup to prompt the user to resume or restart
//   - Updated after each task completion or failure
//   - Cleaned up when the entire change orchestration completes successfully
//
// Example usage:
//
//	session := &SessionState{
//	    ChangeID:      "add-feature-x",
//	    StartedAt:     time.Now(),
//	    LastUpdated:   time.Now(),
//	    CompletedIDs:  []string{},
//	    FailedIDs:     []string{},
//	    RetryCount:    make(map[string]int),
//	}
//
//	// After completing a task
//	session.MarkTaskCompleted("1.1")
//
//	// After a task fails
//	session.MarkTaskFailed("1.2")
//	session.IncrementRetry("1.2")
type SessionState struct {
	// ChangeID is the unique identifier for the change proposal being orchestrated.
	// Example: "add-feature-x", "fix-bug-123"
	// This corresponds to the directory name in spectr/changes/<change-id>/
	ChangeID string `json:"change_id"`

	// StartedAt is the timestamp when the orchestration session was first started.
	// This is set once when the session is created and never changed.
	// Used to calculate total elapsed time and for session history.
	StartedAt time.Time `json:"started_at"`

	// LastUpdated is the timestamp of the most recent session state update.
	// This is updated after every task completion, failure, or retry.
	// Used to detect stale sessions and track activity.
	LastUpdated time.Time `json:"last_updated"`

	// CompletedIDs is the list of task IDs that have been successfully completed.
	// Tasks are added to this list when:
	//   1. The agent process exits with code 0, AND
	//   2. The task status in tasks.jsonc changed to "completed"
	//
	// This list is used to skip completed tasks on session resume.
	// Example: ["1.1", "1.2", "2.1"]
	CompletedIDs []string `json:"completed_ids"`

	// FailedIDs is the list of task IDs that have failed and exhausted all retries.
	// Tasks are added to this list when:
	//   1. The agent process exits with non-zero code or times out, AND
	//   2. The retry count reaches maxRetries
	//
	// Failed tasks are skipped on resume unless the user explicitly retries them.
	// Example: ["1.3", "2.2"]
	FailedIDs []string `json:"failed_ids"`

	// RetryCount tracks how many times each task has been retried.
	// Key: task ID (e.g., "1.2")
	// Value: number of retry attempts (0 = first attempt, 1 = first retry, etc.)
	//
	// When a task fails:
	//   1. Increment the retry count
	//   2. If count < maxRetries, retry the task
	//   3. If count >= maxRetries, add to FailedIDs and prompt user
	//
	// This map is cleared for a task when it completes successfully.
	// Example: {"1.3": 2, "2.1": 1}
	RetryCount map[string]int `json:"retry_count"`

	// CurrentTaskID is the ID of the task currently being executed.
	// This is set when a task starts and cleared when it completes or fails.
	// Used to identify which task was interrupted if the session crashes.
	//
	// Empty string means no task is currently in progress (between tasks).
	// Example: "1.3"
	CurrentTaskID string `json:"current_task_id,omitempty"`
}

// MarkTaskCompleted adds a task to the CompletedIDs list and updates the timestamp.
// This method should be called when:
//   - The agent process exits successfully with code 0
//   - The task status in tasks.jsonc is "completed"
//
// If the task was previously in the retry count map, it is removed since it
// has now completed successfully.
//
// This method also clears CurrentTaskID if it matches the completed task.
func (s *SessionState) MarkTaskCompleted(taskID string) {
	// Add to completed list if not already present
	if !contains(s.CompletedIDs, taskID) {
		s.CompletedIDs = append(s.CompletedIDs, taskID)
	}

	// Remove from retry count since task completed successfully
	delete(s.RetryCount, taskID)

	// Clear current task if it matches
	if s.CurrentTaskID == taskID {
		s.CurrentTaskID = ""
	}

	// Update timestamp
	s.UpdateTimestamp()
}

// MarkTaskFailed adds a task to the FailedIDs list and updates the timestamp.
// This method should be called when:
//   - The agent process exits with non-zero code or times out
//   - The retry count has reached maxRetries
//
// Tasks in FailedIDs are skipped during session resume unless the user
// explicitly chooses to retry them.
//
// This method also clears CurrentTaskID if it matches the failed task.
func (s *SessionState) MarkTaskFailed(taskID string) {
	// Add to failed list if not already present
	if !contains(s.FailedIDs, taskID) {
		s.FailedIDs = append(s.FailedIDs, taskID)
	}

	// Clear current task if it matches
	if s.CurrentTaskID == taskID {
		s.CurrentTaskID = ""
	}

	// Update timestamp
	s.UpdateTimestamp()
}

// IncrementRetry increments the retry count for a task and updates the timestamp.
// This method should be called when:
//   - A task fails (non-zero exit, timeout, or status not updated)
//   - Before deciding whether to retry or mark as permanently failed
//
// The orchestrator should check the retry count after incrementing to decide
// whether to retry the task or mark it as failed:
//
//	session.IncrementRetry(taskID)
//	if session.RetryCount[taskID] >= maxRetries {
//	    session.MarkTaskFailed(taskID)
//	} else {
//	    // Retry the task
//	}
func (s *SessionState) IncrementRetry(taskID string) {
	s.RetryCount[taskID]++
	s.UpdateTimestamp()
}

// UpdateTimestamp updates the LastUpdated field to the current time.
// This method is called automatically by other state-modifying methods
// (MarkTaskCompleted, MarkTaskFailed, IncrementRetry), but can also be
// called manually when other session state changes occur.
func (s *SessionState) UpdateTimestamp() {
	s.LastUpdated = time.Now()
}

// Save writes the session state to disk at the standard location:
// spectr/changes/<change-id>/.ralph-session.json
//
// The session file is saved with pretty-printed JSON for human readability.
// If the directory doesn't exist, it will be created.
//
// This method should be called:
//   - After each task completes or fails
//   - On interruption (Ctrl+C)
//   - On orchestration quit
//
// Returns an error if the file cannot be written.
func (s *SessionState) Save(changeDir string) error {
	sessionPath := filepath.Join(changeDir, ".ralph-session.json")

	// Marshal to pretty-printed JSON
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session state: %w", err)
	}

	// Write to file (create if doesn't exist, overwrite if exists)
	if err := os.WriteFile(sessionPath, data, sessionFilePermissions); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// Load reads the session state from disk at the standard location:
// spectr/changes/<change-id>/.ralph-session.json
//
// This method should be called on orchestration startup to check if a
// previous session exists. If the session file exists, the orchestrator
// should prompt the user to resume or restart.
//
// Returns:
//   - *SessionState: The loaded session state
//   - error: os.ErrNotExist if no session file exists (first run)
//   - error: Parse error if the session file is corrupted
func LoadSession(changeDir string) (*SessionState, error) {
	sessionPath := filepath.Join(changeDir, ".ralph-session.json")

	// Read the file
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return nil, err // Will be os.ErrNotExist if file doesn't exist
	}

	// Unmarshal into SessionState
	var session SessionState
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session state: %w", err)
	}

	// Initialize nil slices to empty slices for consistent behavior
	if session.CompletedIDs == nil {
		session.CompletedIDs = make([]string, 0)
	}
	if session.FailedIDs == nil {
		session.FailedIDs = make([]string, 0)
	}
	if session.RetryCount == nil {
		session.RetryCount = make(map[string]int)
	}

	return &session, nil
}

// Delete removes the session file from disk.
// This method should be called when:
//   - The orchestration completes successfully (all tasks done)
//   - The user chooses to restart instead of resume
//
// Returns an error if the file cannot be deleted (except os.ErrNotExist,
// which is ignored since the goal is to ensure the file doesn't exist).
func DeleteSession(changeDir string) error {
	sessionPath := filepath.Join(changeDir, ".ralph-session.json")

	// Delete the file (ignore if it doesn't exist)
	if err := os.Remove(sessionPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete session file: %w", err)
	}

	return nil
}

// IsCompleted returns true if the specified task is in the CompletedIDs list.
// This is used during session resume to skip already-completed tasks.
func (s *SessionState) IsCompleted(taskID string) bool {
	return contains(s.CompletedIDs, taskID)
}

// IsFailed returns true if the specified task is in the FailedIDs list.
// This is used during session resume to skip tasks that have failed and
// exhausted all retries (unless the user explicitly retries them).
func (s *SessionState) IsFailed(taskID string) bool {
	return contains(s.FailedIDs, taskID)
}

// GetRetryCount returns the number of retries attempted for the specified task.
// Returns 0 if the task has never been retried.
func (s *SessionState) GetRetryCount(taskID string) int {
	return s.RetryCount[taskID]
}

// contains checks if a string slice contains a specific value.
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}
