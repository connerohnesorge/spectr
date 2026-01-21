package ralph

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSessionStateSaveLoad tests round-trip persistence: save → load → verify identical.
func TestSessionStateSaveLoad(t *testing.T) {
	tests := []struct {
		name    string
		session *SessionState
	}{
		{
			name: "basic session state",
			session: &SessionState{
				ChangeID:      "add-feature-x",
				StartedAt:     time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
				LastUpdated:   time.Date(2025, 1, 15, 11, 45, 0, 0, time.UTC),
				CompletedIDs:  []string{"1.1", "1.2", "2.1"},
				FailedIDs:     []string{"3.1"},
				RetryCount:    map[string]int{"2.2": 1, "3.1": 3},
				CurrentTaskID: "2.2",
			},
		},
		{
			name: "empty session",
			session: &SessionState{
				ChangeID:     "empty-change",
				StartedAt:    time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
				LastUpdated:  time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
				CompletedIDs: make([]string, 0),
				FailedIDs:    make([]string, 0),
				RetryCount:   make(map[string]int),
			},
		},
		{
			name: "session with special characters in task IDs",
			session: &SessionState{
				ChangeID:     "feature-with-special-chars",
				StartedAt:    time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
				LastUpdated:  time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
				CompletedIDs: []string{"1.1.1", "1.2.3.4", "10.1"},
				FailedIDs:    make([]string, 0),
				RetryCount:   map[string]int{"1.1.2": 2},
			},
		},
		{
			name: "session with large retry counts",
			session: &SessionState{
				ChangeID:     "retry-heavy-change",
				StartedAt:    time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
				LastUpdated:  time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				CompletedIDs: []string{"1.1"},
				FailedIDs:    []string{"1.2", "1.3"},
				RetryCount:   map[string]int{"1.2": 5, "1.3": 10, "1.4": 1},
			},
		},
		{
			name: "session with many completed tasks",
			session: &SessionState{
				ChangeID:    "large-change",
				StartedAt:   time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
				LastUpdated: time.Date(2025, 1, 15, 18, 0, 0, 0, time.UTC),
				CompletedIDs: []string{
					"1.1",
					"1.2",
					"1.3",
					"2.1",
					"2.2",
					"3.1",
					"3.2",
					"4.1",
					"4.2",
					"4.3",
				},
				FailedIDs:     []string{"5.1"},
				RetryCount:    map[string]int{"5.1": 3},
				CurrentTaskID: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Save the session
			err := tt.session.Save(tmpDir)
			if err != nil {
				t.Fatalf("Save() error = %v", err)
			}

			// Verify file exists
			sessionPath := filepath.Join(tmpDir, ".ralph-session.json")
			if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
				t.Fatalf("session file not created at %s", sessionPath)
			}

			// Load the session
			loaded, err := LoadSession(tmpDir)
			if err != nil {
				t.Fatalf("LoadSession() error = %v", err)
			}

			// Verify all fields match
			if loaded.ChangeID != tt.session.ChangeID {
				t.Errorf("ChangeID = %s, want %s", loaded.ChangeID, tt.session.ChangeID)
			}

			if !loaded.StartedAt.Equal(tt.session.StartedAt) {
				t.Errorf("StartedAt = %v, want %v", loaded.StartedAt, tt.session.StartedAt)
			}

			if !loaded.LastUpdated.Equal(tt.session.LastUpdated) {
				t.Errorf("LastUpdated = %v, want %v", loaded.LastUpdated, tt.session.LastUpdated)
			}

			if !stringSlicesEqual(loaded.CompletedIDs, tt.session.CompletedIDs) {
				t.Errorf("CompletedIDs = %v, want %v", loaded.CompletedIDs, tt.session.CompletedIDs)
			}

			if !stringSlicesEqual(loaded.FailedIDs, tt.session.FailedIDs) {
				t.Errorf("FailedIDs = %v, want %v", loaded.FailedIDs, tt.session.FailedIDs)
			}

			if !intMapsEqual(loaded.RetryCount, tt.session.RetryCount) {
				t.Errorf("RetryCount = %v, want %v", loaded.RetryCount, tt.session.RetryCount)
			}

			if loaded.CurrentTaskID != tt.session.CurrentTaskID {
				t.Errorf(
					"CurrentTaskID = %s, want %s",
					loaded.CurrentTaskID,
					tt.session.CurrentTaskID,
				)
			}
		})
	}
}

// TestSessionStateMarkTaskCompleted tests marking tasks as completed.
func TestSessionStateMarkTaskCompleted(t *testing.T) {
	tests := []struct {
		name            string
		initial         *SessionState
		taskID          string
		wantCompleted   []string
		wantRetryCount  map[string]int
		wantCurrentTask string
	}{
		{
			name: "mark first task completed",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  make([]string, 0),
				FailedIDs:     make([]string, 0),
				RetryCount:    make(map[string]int),
				CurrentTaskID: "1.1",
			},
			taskID:          "1.1",
			wantCompleted:   []string{"1.1"},
			wantRetryCount:  make(map[string]int),
			wantCurrentTask: "",
		},
		{
			name: "mark task with retries completed (clears retry count)",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  []string{"1.1"},
				FailedIDs:     make([]string, 0),
				RetryCount:    map[string]int{"1.2": 2},
				CurrentTaskID: "1.2",
			},
			taskID:          "1.2",
			wantCompleted:   []string{"1.1", "1.2"},
			wantRetryCount:  make(map[string]int),
			wantCurrentTask: "",
		},
		{
			name: "mark duplicate task completed (idempotent)",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  []string{"1.1"},
				FailedIDs:     make([]string, 0),
				RetryCount:    make(map[string]int),
				CurrentTaskID: "",
			},
			taskID:          "1.1",
			wantCompleted:   []string{"1.1"},
			wantRetryCount:  make(map[string]int),
			wantCurrentTask: "",
		},
		{
			name: "mark task completed doesn't affect other retry counts",
			initial: &SessionState{
				ChangeID:     "test-change",
				StartedAt:    time.Now(),
				LastUpdated:  time.Now(),
				CompletedIDs: make([]string, 0),
				FailedIDs:    make([]string, 0),
				RetryCount:   map[string]int{"1.2": 1, "2.1": 3},
			},
			taskID:         "1.1",
			wantCompleted:  []string{"1.1"},
			wantRetryCount: map[string]int{"1.2": 1, "2.1": 3},
		},
		{
			name: "mark task completed when current task is different",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  make([]string, 0),
				FailedIDs:     make([]string, 0),
				RetryCount:    make(map[string]int),
				CurrentTaskID: "2.1",
			},
			taskID:          "1.1",
			wantCompleted:   []string{"1.1"},
			wantRetryCount:  make(map[string]int),
			wantCurrentTask: "2.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := tt.initial
			beforeUpdate := session.LastUpdated

			// Mark task completed
			session.MarkTaskCompleted(tt.taskID)

			// Verify CompletedIDs
			if !stringSlicesEqual(session.CompletedIDs, tt.wantCompleted) {
				t.Errorf("CompletedIDs = %v, want %v", session.CompletedIDs, tt.wantCompleted)
			}

			// Verify RetryCount
			if !intMapsEqual(session.RetryCount, tt.wantRetryCount) {
				t.Errorf("RetryCount = %v, want %v", session.RetryCount, tt.wantRetryCount)
			}

			// Verify CurrentTaskID
			if session.CurrentTaskID != tt.wantCurrentTask {
				t.Errorf("CurrentTaskID = %s, want %s", session.CurrentTaskID, tt.wantCurrentTask)
			}

			// Verify LastUpdated was updated
			if !session.LastUpdated.After(beforeUpdate) {
				t.Error("LastUpdated was not updated")
			}
		})
	}
}

// TestSessionStateMarkTaskFailed tests marking tasks as failed.
func TestSessionStateMarkTaskFailed(t *testing.T) {
	tests := []struct {
		name            string
		initial         *SessionState
		taskID          string
		wantFailed      []string
		wantCurrentTask string
	}{
		{
			name: "mark first task failed",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  make([]string, 0),
				FailedIDs:     make([]string, 0),
				RetryCount:    make(map[string]int),
				CurrentTaskID: "1.1",
			},
			taskID:          "1.1",
			wantFailed:      []string{"1.1"},
			wantCurrentTask: "",
		},
		{
			name: "mark additional task failed",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  []string{"1.1"},
				FailedIDs:     []string{"2.1"},
				RetryCount:    make(map[string]int),
				CurrentTaskID: "2.2",
			},
			taskID:          "2.2",
			wantFailed:      []string{"2.1", "2.2"},
			wantCurrentTask: "",
		},
		{
			name: "mark duplicate task failed (idempotent)",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  make([]string, 0),
				FailedIDs:     []string{"1.1"},
				RetryCount:    make(map[string]int),
				CurrentTaskID: "",
			},
			taskID:          "1.1",
			wantFailed:      []string{"1.1"},
			wantCurrentTask: "",
		},
		{
			name: "mark task failed when current task is different",
			initial: &SessionState{
				ChangeID:      "test-change",
				StartedAt:     time.Now(),
				LastUpdated:   time.Now(),
				CompletedIDs:  make([]string, 0),
				FailedIDs:     make([]string, 0),
				RetryCount:    make(map[string]int),
				CurrentTaskID: "2.1",
			},
			taskID:          "1.1",
			wantFailed:      []string{"1.1"},
			wantCurrentTask: "2.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := tt.initial
			beforeUpdate := session.LastUpdated

			// Mark task failed
			session.MarkTaskFailed(tt.taskID)

			// Verify FailedIDs
			if !stringSlicesEqual(session.FailedIDs, tt.wantFailed) {
				t.Errorf("FailedIDs = %v, want %v", session.FailedIDs, tt.wantFailed)
			}

			// Verify CurrentTaskID
			if session.CurrentTaskID != tt.wantCurrentTask {
				t.Errorf("CurrentTaskID = %s, want %s", session.CurrentTaskID, tt.wantCurrentTask)
			}

			// Verify LastUpdated was updated
			if !session.LastUpdated.After(beforeUpdate) {
				t.Error("LastUpdated was not updated")
			}
		})
	}
}

// TestSessionStateIncrementRetry tests incrementing retry counts.
func TestSessionStateIncrementRetry(t *testing.T) {
	tests := []struct {
		name           string
		initial        *SessionState
		taskID         string
		wantRetryCount int
	}{
		{
			name: "first retry",
			initial: &SessionState{
				ChangeID:     "test-change",
				StartedAt:    time.Now(),
				LastUpdated:  time.Now(),
				CompletedIDs: make([]string, 0),
				FailedIDs:    make([]string, 0),
				RetryCount:   make(map[string]int),
			},
			taskID:         "1.1",
			wantRetryCount: 1,
		},
		{
			name: "increment existing retry",
			initial: &SessionState{
				ChangeID:     "test-change",
				StartedAt:    time.Now(),
				LastUpdated:  time.Now(),
				CompletedIDs: make([]string, 0),
				FailedIDs:    make([]string, 0),
				RetryCount:   map[string]int{"1.1": 2},
			},
			taskID:         "1.1",
			wantRetryCount: 3,
		},
		{
			name: "multiple retries",
			initial: &SessionState{
				ChangeID:     "test-change",
				StartedAt:    time.Now(),
				LastUpdated:  time.Now(),
				CompletedIDs: make([]string, 0),
				FailedIDs:    make([]string, 0),
				RetryCount:   map[string]int{"1.1": 9},
			},
			taskID:         "1.1",
			wantRetryCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := tt.initial
			beforeUpdate := session.LastUpdated

			// Increment retry
			session.IncrementRetry(tt.taskID)

			// Verify retry count
			if session.RetryCount[tt.taskID] != tt.wantRetryCount {
				t.Errorf("RetryCount[%s] = %d, want %d",
					tt.taskID, session.RetryCount[tt.taskID], tt.wantRetryCount)
			}

			// Verify LastUpdated was updated
			if !session.LastUpdated.After(beforeUpdate) {
				t.Error("LastUpdated was not updated")
			}
		})
	}
}

// TestSessionStateQueryMethods tests IsCompleted, IsFailed, and GetRetryCount.
func TestSessionStateQueryMethods(t *testing.T) {
	session := &SessionState{
		ChangeID:     "test-change",
		StartedAt:    time.Now(),
		LastUpdated:  time.Now(),
		CompletedIDs: []string{"1.1", "1.2", "2.1"},
		FailedIDs:    []string{"3.1", "3.2"},
		RetryCount:   map[string]int{"2.2": 1, "3.1": 3, "4.1": 2},
	}

	// Test IsCompleted
	completedTests := []struct {
		taskID string
		want   bool
	}{
		{"1.1", true},
		{"1.2", true},
		{"2.1", true},
		{"2.2", false},
		{"3.1", false},
		{"9.9", false},
	}

	for _, tt := range completedTests {
		t.Run("IsCompleted/"+tt.taskID, func(t *testing.T) {
			got := session.IsCompleted(tt.taskID)
			if got != tt.want {
				t.Errorf("IsCompleted(%s) = %v, want %v", tt.taskID, got, tt.want)
			}
		})
	}

	// Test IsFailed
	failedTests := []struct {
		taskID string
		want   bool
	}{
		{"3.1", true},
		{"3.2", true},
		{"1.1", false},
		{"2.2", false},
		{"9.9", false},
	}

	for _, tt := range failedTests {
		t.Run("IsFailed/"+tt.taskID, func(t *testing.T) {
			got := session.IsFailed(tt.taskID)
			if got != tt.want {
				t.Errorf("IsFailed(%s) = %v, want %v", tt.taskID, got, tt.want)
			}
		})
	}

	// Test GetRetryCount
	retryTests := []struct {
		taskID string
		want   int
	}{
		{"2.2", 1},
		{"3.1", 3},
		{"4.1", 2},
		{"1.1", 0}, // completed, not in retry map
		{"9.9", 0}, // never retried
	}

	for _, tt := range retryTests {
		t.Run("GetRetryCount/"+tt.taskID, func(t *testing.T) {
			got := session.GetRetryCount(tt.taskID)
			if got != tt.want {
				t.Errorf("GetRetryCount(%s) = %d, want %d", tt.taskID, got, tt.want)
			}
		})
	}
}

// TestLoadSessionNotExists tests loading a non-existent session.
func TestLoadSessionNotExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Attempt to load non-existent session
	_, err := LoadSession(tmpDir)

	// Should return os.ErrNotExist
	if !os.IsNotExist(err) {
		t.Errorf("LoadSession() error = %v, want os.ErrNotExist", err)
	}
}

// TestLoadSessionCorruptedJSON tests loading a corrupted session file.
func TestLoadSessionCorruptedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, ".ralph-session.json")

	// Write invalid JSON
	invalidJSON := `{
		"change_id": "test",
		"started_at": "invalid-date-format",
		invalid syntax here
	}`

	err := os.WriteFile(sessionPath, []byte(invalidJSON), 0o644)
	if err != nil {
		t.Fatalf("failed to write corrupted session file: %v", err)
	}

	// Attempt to load
	_, err = LoadSession(tmpDir)

	// Should return an error
	if err == nil {
		t.Error("LoadSession() expected error for corrupted JSON, got nil")
	}
}

// TestDeleteSession tests deleting session files.
func TestDeleteSession(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  bool
		wantErr    bool
		expectGone bool
	}{
		{
			name:       "delete existing session",
			setupFile:  true,
			wantErr:    false,
			expectGone: true,
		},
		{
			name:       "delete non-existent session (idempotent)",
			setupFile:  false,
			wantErr:    false,
			expectGone: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			sessionPath := filepath.Join(tmpDir, ".ralph-session.json")

			// Setup: create session file if needed
			if tt.setupFile {
				session := &SessionState{
					ChangeID:     "test-change",
					StartedAt:    time.Now(),
					LastUpdated:  time.Now(),
					CompletedIDs: []string{"1.1"},
					FailedIDs:    make([]string, 0),
					RetryCount:   make(map[string]int),
				}
				if err := session.Save(tmpDir); err != nil {
					t.Fatalf("failed to setup session file: %v", err)
				}

				// Verify file exists before deletion
				if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
					t.Fatal("setup failed: session file not created")
				}
			}

			// Delete the session
			err := DeleteSession(tmpDir)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify file is gone
			if !tt.expectGone {
				return
			}
			if _, err := os.Stat(sessionPath); !os.IsNotExist(err) {
				t.Error("session file still exists after deletion")
			}
		})
	}
}

// TestSessionStateUpdateTimestamp tests the UpdateTimestamp method.
func TestSessionStateUpdateTimestamp(t *testing.T) {
	session := &SessionState{
		ChangeID:     "test-change",
		StartedAt:    time.Now(),
		LastUpdated:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CompletedIDs: nil,
		FailedIDs:    nil,
		RetryCount:   make(map[string]int),
	}

	oldTimestamp := session.LastUpdated

	// Wait a tiny bit to ensure time has passed
	time.Sleep(10 * time.Millisecond)

	// Update timestamp
	session.UpdateTimestamp()

	// Verify timestamp was updated
	if !session.LastUpdated.After(oldTimestamp) {
		t.Error("UpdateTimestamp() did not update LastUpdated field")
	}

	// Verify StartedAt was not changed
	if !session.StartedAt.Before(session.LastUpdated) {
		t.Error("UpdateTimestamp() should not change StartedAt")
	}
}

// TestSessionStateJSONFormat tests that saved JSON is properly formatted.
func TestSessionStateJSONFormat(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, ".ralph-session.json")

	session := &SessionState{
		ChangeID:      "test-change",
		StartedAt:     time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		LastUpdated:   time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC),
		CompletedIDs:  []string{"1.1", "1.2"},
		FailedIDs:     []string{"2.1"},
		RetryCount:    map[string]int{"2.2": 1},
		CurrentTaskID: "2.2",
	}

	// Save the session
	err := session.Save(tmpDir)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read the raw JSON
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		t.Fatalf("failed to read session file: %v", err)
	}

	// Verify it's valid JSON
	var parsed SessionState
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("saved JSON is not valid: %v", err)
	}

	// Verify it's pretty-printed (contains newlines and indentation)
	jsonStr := string(data)
	if !strings.Contains(jsonStr, "\n") {
		t.Error("JSON is not pretty-printed (no newlines found)")
	}

	if !strings.Contains(jsonStr, "  ") {
		t.Error("JSON is not indented")
	}
}

// TestSessionStateWithSpecialCharacters tests handling special characters in task IDs.
func TestSessionStateWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()

	// Task IDs with various special characters (valid in JSON)
	session := &SessionState{
		ChangeID:      "test-change-with-dashes-and_underscores",
		StartedAt:     time.Now(),
		LastUpdated:   time.Now(),
		CompletedIDs:  []string{"1.1", "1.2.3", "10.1", "2.1.1.1"},
		FailedIDs:     []string{"3-1", "3_2"}, // technically valid task IDs
		RetryCount:    map[string]int{"4.1": 1, "4.2.1": 2},
		CurrentTaskID: "4.2.1",
	}

	// Save and load
	err := session.Save(tmpDir)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := LoadSession(tmpDir)
	if err != nil {
		t.Fatalf("LoadSession() error = %v", err)
	}

	// Verify round-trip
	if !stringSlicesEqual(loaded.CompletedIDs, session.CompletedIDs) {
		t.Errorf("CompletedIDs = %v, want %v", loaded.CompletedIDs, session.CompletedIDs)
	}

	if !stringSlicesEqual(loaded.FailedIDs, session.FailedIDs) {
		t.Errorf("FailedIDs = %v, want %v", loaded.FailedIDs, session.FailedIDs)
	}

	if !intMapsEqual(loaded.RetryCount, session.RetryCount) {
		t.Errorf("RetryCount = %v, want %v", loaded.RetryCount, session.RetryCount)
	}
}

// TestSessionStateEmptyState tests handling empty/nil slices and maps.
func TestSessionStateEmptyState(t *testing.T) {
	tmpDir := t.TempDir()

	session := &SessionState{
		ChangeID:     "empty-state-test",
		StartedAt:    time.Now(),
		LastUpdated:  time.Now(),
		CompletedIDs: nil,
		FailedIDs:    nil, // nil slice
		RetryCount:   nil, // nil map
	}

	// Save and load
	err := session.Save(tmpDir)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := LoadSession(tmpDir)
	if err != nil {
		t.Fatalf("LoadSession() error = %v", err)
	}

	// Verify empty/nil handling
	if loaded.CompletedIDs == nil {
		t.Error("CompletedIDs should not be nil after load")
	}

	if len(loaded.CompletedIDs) != 0 {
		t.Errorf("CompletedIDs = %v, want empty slice", loaded.CompletedIDs)
	}

	// FailedIDs might be nil or empty after JSON round-trip
	if len(loaded.FailedIDs) != 0 {
		t.Errorf("FailedIDs = %v, want nil or empty", loaded.FailedIDs)
	}

	// RetryCount might be nil or empty after JSON round-trip
	if len(loaded.RetryCount) != 0 {
		t.Errorf("RetryCount = %v, want nil or empty", loaded.RetryCount)
	}
}

// Helper functions

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func intMapsEqual(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, exists := b[key]
		if !exists || valA != valB {
			return false
		}
	}

	return true
}
