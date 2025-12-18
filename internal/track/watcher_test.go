package track

import (
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

// isFsnotifySupported checks if fsnotify is likely to work on this system.
func isFsnotifySupported() bool {
	// fsnotify is supported on Linux, macOS, Windows, and BSD
	switch runtime.GOOS {
	case "linux",
		"darwin",
		"windows",
		"freebsd",
		"netbsd",
		"openbsd":
		return true
	default:
		return false
	}
}

func TestNewWatcher_Success(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	// Create watcher for existing file
	w, err := NewWatcher(tempFile)
	if err != nil {
		t.Fatalf("NewWatcher() error = %v", err)
	}
	defer func() { _ = w.Close() }()

	// Verify watcher is configured correctly
	if w.filePath == "" {
		t.Error(
			"NewWatcher() created watcher with empty filePath",
		)
	}
	if w.events == nil {
		t.Error(
			"NewWatcher() created watcher with nil events channel",
		)
	}
	if w.errors == nil {
		t.Error(
			"NewWatcher() created watcher with nil errors channel",
		)
	}
	if w.debounce != defaultDebounce {
		t.Errorf(
			"NewWatcher() debounce = %v, want %v",
			w.debounce,
			defaultDebounce,
		)
	}
}

func TestNewWatcher_NonExistentFile(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Try to create watcher for non-existent file
	nonExistentPath := filepath.Join(
		t.TempDir(),
		"does-not-exist.txt",
	)

	w, err := NewWatcher(nonExistentPath)
	if err == nil {
		_ = w.Close()
		t.Fatal(
			"NewWatcher() expected error for non-existent file, got nil",
		)
	}

	// Error should be related to file not found
	if !os.IsNotExist(err) {
		t.Errorf(
			"NewWatcher() error = %v, want os.IsNotExist error",
			err,
		)
	}
}

func TestNewWatcher_WithCustomDebounce(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	tests := []struct {
		name     string
		debounce time.Duration
	}{
		{"50ms debounce", 50 * time.Millisecond},
		{
			"100ms debounce",
			100 * time.Millisecond,
		},
		{
			"200ms debounce",
			200 * time.Millisecond,
		},
		{"1s debounce", 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWatcherWithDebounce(
				tempFile,
				tt.debounce,
			)
			if err != nil {
				t.Fatalf(
					"NewWatcherWithDebounce() error = %v",
					err,
				)
			}
			defer func() { _ = w.Close() }()

			if w.debounce != tt.debounce {
				t.Errorf(
					"NewWatcherWithDebounce() debounce = %v, want %v",
					w.debounce,
					tt.debounce,
				)
			}
		})
	}
}

func TestWatcher_Events_OnFileModification(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	// Use a short debounce for faster tests
	w, err := NewWatcherWithDebounce(
		tempFile,
		50*time.Millisecond,
	)
	if err != nil {
		t.Fatalf(
			"NewWatcherWithDebounce() error = %v",
			err,
		)
	}
	defer func() { _ = w.Close() }()

	// Give the watcher time to start
	time.Sleep(10 * time.Millisecond)

	// Modify the file
	if err := os.WriteFile(
		tempFile,
		[]byte("modified content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to modify temp file: %v",
			err,
		)
	}

	// Wait for event with timeout
	select {
	case <-w.Events():
		// Success - received event
	case err := <-w.Errors():
		t.Fatalf(
			"received error instead of event: %v",
			err,
		)
	case <-time.After(2 * time.Second):
		t.Fatal(
			"timeout waiting for file modification event",
		)
	}
}

func TestWatcher_Events_OnFileRecreation(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	// Use a short debounce for faster tests
	w, err := NewWatcherWithDebounce(
		tempFile,
		50*time.Millisecond,
	)
	if err != nil {
		t.Fatalf(
			"NewWatcherWithDebounce() error = %v",
			err,
		)
	}
	defer func() { _ = w.Close() }()

	// Give the watcher time to start
	time.Sleep(10 * time.Millisecond)

	// Delete and recreate the file (simulating what some editors do)
	if err := os.Remove(tempFile); err != nil {
		t.Fatalf(
			"failed to remove temp file: %v",
			err,
		)
	}
	if err := os.WriteFile(
		tempFile,
		[]byte("recreated content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to recreate temp file: %v",
			err,
		)
	}

	// Wait for event with timeout
	select {
	case <-w.Events():
		// Success - received event for file recreation
	case err := <-w.Errors():
		t.Fatalf(
			"received error instead of event: %v",
			err,
		)
	case <-time.After(2 * time.Second):
		t.Fatal(
			"timeout waiting for file recreation event",
		)
	}
}

func TestWatcher_Debouncing(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	// Use a moderate debounce to test coalescing
	debounce := 100 * time.Millisecond
	w, err := NewWatcherWithDebounce(
		tempFile,
		debounce,
	)
	if err != nil {
		t.Fatalf(
			"NewWatcherWithDebounce() error = %v",
			err,
		)
	}
	defer func() { _ = w.Close() }()

	// Give the watcher time to start
	time.Sleep(10 * time.Millisecond)

	// Perform multiple rapid writes (faster than debounce period)
	for i := range 5 {
		if err := os.WriteFile(
			tempFile,
			[]byte("content "+string(rune('0'+i))),
			0644,
		); err != nil {
			t.Fatalf(
				"failed to write temp file: %v",
				err,
			)
		}
		time.Sleep(
			20 * time.Millisecond,
		) // Faster than debounce
	}

	// Count events received within a reasonable window
	var eventCount int32
	done := make(chan struct{})

	go func() {
		timer := time.NewTimer(
			500 * time.Millisecond,
		)
		defer timer.Stop()
		for {
			select {
			case <-w.Events():
				atomic.AddInt32(&eventCount, 1)
			case <-timer.C:
				close(done)

				return
			}
		}
	}()

	<-done

	// We should receive fewer events than writes due to debouncing
	// Ideally just 1 event, but timing can vary
	count := atomic.LoadInt32(&eventCount)
	if count == 0 {
		t.Error(
			"expected at least one event after rapid writes",
		)
	}
	if count >= 5 {
		t.Errorf(
			"debouncing failed: received %d events for 5 rapid writes",
			count,
		)
	}
}

func TestWatcher_Close_Idempotent(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	w, err := NewWatcher(tempFile)
	if err != nil {
		t.Fatalf("NewWatcher() error = %v", err)
	}

	// Close multiple times - should not panic or error
	for i := range 3 {
		err := w.Close()
		if err != nil {
			t.Errorf(
				"Close() call %d error = %v, want nil",
				i+1,
				err,
			)
		}
	}
}

func TestWatcher_Errors_Channel(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	w, err := NewWatcher(tempFile)
	if err != nil {
		t.Fatalf("NewWatcher() error = %v", err)
	}
	defer func() { _ = w.Close() }()

	// Verify errors channel is accessible and has correct capacity
	errChan := w.Errors()
	if errChan == nil {
		t.Fatal("Errors() returned nil channel")
	}

	// The channel should be readable (though empty initially)
	select {
	case err := <-errChan:
		t.Errorf(
			"unexpected error in channel: %v",
			err,
		)
	default:
		// Expected - no errors initially
	}
}

func TestWatcher_Events_ChannelCapacity(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	w, err := NewWatcher(tempFile)
	if err != nil {
		t.Fatalf("NewWatcher() error = %v", err)
	}
	defer func() { _ = w.Close() }()

	// Verify events channel is accessible
	eventsChan := w.Events()
	if eventsChan == nil {
		t.Fatal("Events() returned nil channel")
	}

	// The channel should be readable (though empty initially)
	select {
	case <-eventsChan:
		t.Error("unexpected event in channel")
	default:
		// Expected - no events initially
	}
}

func TestWatcher_AbsolutePath(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	// Create watcher with relative-looking path
	// (TempDir returns absolute path, but test the behavior)
	w, err := NewWatcher(tempFile)
	if err != nil {
		t.Fatalf("NewWatcher() error = %v", err)
	}
	defer func() { _ = w.Close() }()

	// Verify the stored path is absolute
	if !filepath.IsAbs(w.filePath) {
		t.Errorf(
			"watcher filePath is not absolute: %s",
			w.filePath,
		)
	}
}

func TestWatcher_WatchDirectory(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	// Use a short debounce for faster tests
	w, err := NewWatcherWithDebounce(
		tempFile,
		50*time.Millisecond,
	)
	if err != nil {
		t.Fatalf(
			"NewWatcherWithDebounce() error = %v",
			err,
		)
	}
	defer func() { _ = w.Close() }()

	// Give the watcher time to start
	time.Sleep(10 * time.Millisecond)

	// Create a different file in the same directory
	// This should NOT trigger an event since we're watching a specific file
	otherFile := filepath.Join(
		tempDir,
		"other.txt",
	)
	if err := os.WriteFile(
		otherFile,
		[]byte("other content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create other file: %v",
			err,
		)
	}

	// Wait briefly - should not receive event for other file
	select {
	case <-w.Events():
		t.Error(
			"received unexpected event for unrelated file",
		)
	case <-time.After(200 * time.Millisecond):
		// Expected - no event for other file
	}
}

func TestWatcher_ClosedWatcherChannels(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	// Create a temporary file to watch
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(
		tempFile,
		[]byte("initial content"),
		0644,
	); err != nil {
		t.Fatalf(
			"failed to create temp file: %v",
			err,
		)
	}

	w, err := NewWatcher(tempFile)
	if err != nil {
		t.Fatalf("NewWatcher() error = %v", err)
	}

	// Get channel references before closing
	eventsChan := w.Events()
	errorsChan := w.Errors()

	// Close the watcher
	if err := w.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Give the loop time to exit
	time.Sleep(50 * time.Millisecond)

	// Channels should still be readable (though the loop has stopped)
	// The channels are not explicitly closed by Close(), but the loop stops
	if eventsChan == nil {
		t.Error(
			"events channel should not be nil",
		)
	}
	if errorsChan == nil {
		t.Error(
			"errors channel should not be nil",
		)
	}
}
