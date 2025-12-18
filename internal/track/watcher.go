//nolint:revive // cognitive-complexity is acceptable for event loops
package track

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	// defaultDebounce is the default debounce duration for file events.
	// Editors often perform multiple writes in rapid succession.
	defaultDebounce = 150 * time.Millisecond
)

// Watcher monitors a file for changes using fsnotify with debouncing.
// It handles rapid successive writes from editors by coalescing events
// within a debounce window.
type Watcher struct {
	watcher  *fsnotify.Watcher
	filePath string
	events   chan struct{}
	errors   chan error
	done     chan struct{}
	debounce time.Duration
	mu       sync.Mutex
	closed   bool
}

// NewWatcher creates a new Watcher for the specified file path.
// The file must exist at creation time.
func NewWatcher(
	filePath string,
) (*Watcher, error) {
	return NewWatcherWithDebounce(
		filePath,
		defaultDebounce,
	)
}

// NewWatcherWithDebounce creates a new Watcher with custom debounce.
// The file must exist at creation time.
func NewWatcherWithDebounce(
	filePath string,
	debounce time.Duration,
) (*Watcher, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(absPath); err != nil {
		return nil, err
	}

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(absPath)
	if err := fsWatcher.Add(dir); err != nil {
		_ = fsWatcher.Close()

		return nil, err
	}

	w := &Watcher{
		watcher:  fsWatcher,
		filePath: absPath,
		events:   make(chan struct{}, 1),
		errors:   make(chan error, 1),
		done:     make(chan struct{}),
		debounce: debounce,
	}

	go w.loop()

	return w, nil
}

// Events returns a channel that receives a notification when the watched
// file changes. The channel is buffered with capacity 1, so only the most
// recent event is retained if the consumer is slow.
func (w *Watcher) Events() <-chan struct{} {
	return w.events
}

// Errors returns a channel that receives errors from the underlying
// fsnotify watcher. The channel is buffered with capacity 1.
func (w *Watcher) Errors() <-chan error {
	return w.errors
}

// Close stops the watcher and releases resources.
// It is safe to call Close multiple times.
func (w *Watcher) Close() error {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()

		return nil
	}
	w.closed = true
	w.mu.Unlock()

	close(w.done)

	return w.watcher.Close()
}

// loop is the main event loop that processes fsnotify events.
// It implements debouncing by waiting for a quiet period after events.
func (w *Watcher) loop() {
	var (
		timer     *time.Timer
		timerChan <-chan time.Time
	)

	for {
		select {
		case <-w.done:
			if timer != nil {
				timer.Stop()
			}

			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			timer, timerChan = w.handleEvent(
				event,
				timer,
				timerChan,
			)

		case <-timerChan:
			w.sendEvent()
			timer = nil
			timerChan = nil

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.sendError(err)
		}
	}
}

// handleEvent processes a single fsnotify event.
func (w *Watcher) handleEvent(
	event fsnotify.Event,
	timer *time.Timer,
	timerChan <-chan time.Time,
) (*time.Timer, <-chan time.Time) {
	if !w.isWatchedFile(event.Name) {
		return timer, timerChan
	}

	if !event.Has(fsnotify.Write) &&
		!event.Has(fsnotify.Create) {
		return timer, timerChan
	}

	if timer == nil {
		timer = time.NewTimer(w.debounce)

		return timer, timer.C
	}

	w.resetTimer(timer)

	return timer, timerChan
}

// resetTimer stops and resets the debounce timer.
func (w *Watcher) resetTimer(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(w.debounce)
}

// isWatchedFile checks if the event path matches the watched file.
func (w *Watcher) isWatchedFile(
	eventPath string,
) bool {
	absEventPath, err := filepath.Abs(eventPath)
	if err != nil {
		return false
	}

	return absEventPath == w.filePath
}

// sendEvent sends a notification on the events channel.
// It is non-blocking; if the channel is full, the event is dropped.
func (w *Watcher) sendEvent() {
	select {
	case w.events <- struct{}{}:
	default:
		// Channel full, event coalesced
	}
}

// sendError sends an error on the errors channel.
// It is non-blocking; if the channel is full, the error is dropped.
func (w *Watcher) sendError(err error) {
	select {
	case w.errors <- err:
	default:
		// Channel full, error dropped
	}
}
