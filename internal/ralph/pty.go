// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/creack/pty"
)

// pty.go handles PTY (pseudo-terminal) subprocess management for agent CLIs.
// It uses the creack/pty library for cross-platform PTY support.
//
// Features:
// - Spawning agent CLI processes in a PTY environment
// - Streaming output to the TUI in real-time
// - Handling PTY resize events from terminal size changes
// - Graceful process termination on skip/abort
// - Process lifecycle management (start, wait, kill)
//
// The PTY interface ensures agent CLIs receive proper terminal
// emulation and can use interactive features like progress bars.

// PTYSession represents a running process in a PTY with lifecycle management.
type PTYSession struct {
	PTY     *os.File    // PTY file handle for reading/writing
	Process *os.Process // Process handle for signaling
	Command *exec.Cmd   // Original command for reference
}

// SpawnPTY starts a command in a PTY and returns the PTY file handle.
//
// The PTY provides full terminal emulation, supporting ANSI colors,
// progress bars, and interactive prompts. The returned file handle
// can be read like any file to stream output to the TUI.
//
// Example:
//
//	cmd := exec.Command("claude", "code", "--prompt", "Fix bug in handler.go")
//	ptyFile, err := SpawnPTY(cmd)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer ptyFile.Close()
//	io.Copy(os.Stdout, ptyFile) // Stream output
//
// Returns error if PTY creation or process start fails.
func SpawnPTY(cmd *exec.Cmd) (*os.File, error) {
	// Start the command with a PTY
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return ptyFile, nil
}

// StartPTYSession starts a command with PTY and returns a session for lifecycle management.
//
// The returned PTYSession provides methods for waiting, killing, and cleaning up
// the process and PTY resources.
//
// Example:
//
//	cmd := exec.Command("claude", "code", "--prompt", "Implement feature X")
//	session, err := StartPTYSession(cmd)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer session.Close()
//
//	// Stream output in a goroutine
//	go io.Copy(os.Stdout, session.PTY)
//
//	// Wait for completion
//	if err := session.Wait(); err != nil {
//	    log.Printf("Process failed: %v", err)
//	}
//
// Returns error if PTY creation or process start fails.
func StartPTYSession(cmd *exec.Cmd) (*PTYSession, error) {
	ptyFile, err := SpawnPTY(cmd)
	if err != nil {
		return nil, err
	}

	return &PTYSession{
		PTY:     ptyFile,
		Process: cmd.Process,
		Command: cmd,
	}, nil
}

// Wait blocks until the process completes and returns its exit status.
//
// Returns nil if the process exited successfully (exit code 0).
// Returns an error if the process was terminated by a signal or
// exited with a non-zero status.
//
// This method should be called after reading all output from the PTY
// to avoid deadlocks.
func (s *PTYSession) Wait() error {
	if s.Command == nil {
		return nil
	}

	return s.Command.Wait()
}

// Terminate terminates the process gracefully with a timeout.
//
// This method first attempts graceful shutdown by sending SIGTERM to the process,
// allowing it to clean up resources and exit normally. If the process does not
// exit within the specified timeout, it sends SIGKILL for forceful termination.
//
// Parameters:
//   - timeout: Maximum time to wait for graceful shutdown before forcing termination.
//     Use 0 for immediate forceful termination (equivalent to Kill()).
//
// Returns error if:
//   - Process is nil (returns nil, no-op)
//   - SIGTERM signal cannot be sent
//   - SIGKILL fails after timeout (only if SIGTERM was successful)
//
// Example:
//
//	// Give process 5 seconds to shutdown gracefully
//	if err := session.Terminate(5 * time.Second); err != nil {
//	    log.Printf("Failed to terminate process: %v", err)
//	}
//
// Use cases:
//   - User skips a task: Terminate(5*time.Second) for clean agent shutdown
//   - User aborts orchestration: Terminate(5*time.Second) then Close()
//   - Timeout exceeded: Terminate(0) for immediate forceful termination
func (s *PTYSession) Terminate(timeout time.Duration) error {
	if s.Process == nil {
		return nil
	}

	// If timeout is 0, skip graceful shutdown and kill immediately
	if timeout == 0 {
		return s.Kill()
	}

	// Send SIGTERM for graceful shutdown
	// Note: We send to the process group (negative PID) to ensure child processes receive it
	if err := syscall.Kill(-s.Process.Pid, syscall.SIGTERM); err != nil {
		// If process group kill fails, try individual process
		if err := s.Process.Signal(syscall.SIGTERM); err != nil {
			// If SIGTERM fails (e.g., process already exited), return the error
			return err
		}
	}

	// Set up channels for waiting with timeout
	done := make(chan error, 1)

	// Wait for process to exit in a goroutine
	go func() {
		done <- s.Wait()
	}()

	// Wait for either process exit or timeout
	select {
	case <-done:
		// Process exited (gracefully or with error), we're done
		return nil
	case <-time.After(timeout):
		// Timeout expired, force kill
		return s.Kill()
	}
}

// Kill terminates the process forcefully and immediately.
//
// Sends SIGKILL to the process, which cannot be caught or ignored.
// This is immediate and does not give the process a chance to clean up.
//
// For graceful shutdown with a timeout, use Terminate() instead.
//
// Returns error if the process has already exited or cannot be killed.
func (s *PTYSession) Kill() error {
	if s.Process == nil {
		return nil
	}

	return s.Process.Kill()
}

// Close cleans up PTY resources.
//
// This method closes the PTY file handle but does not terminate the process.
// Call Kill() first if you need to stop the process before cleanup.
//
// Returns error if the PTY cannot be closed.
func (s *PTYSession) Close() error {
	if s.PTY == nil {
		return nil
	}

	return s.PTY.Close()
}

// ReadOutput reads all output from the PTY until EOF.
//
// This is a convenience method for testing and simple use cases.
// For interactive TUI, use direct reading from s.PTY with bufio.Scanner.
//
// Returns all output as a byte slice, or error if reading fails.
func (s *PTYSession) ReadOutput() ([]byte, error) {
	if s.PTY == nil {
		return nil, nil
	}

	return io.ReadAll(s.PTY)
}

// Resize updates the PTY dimensions to match terminal window size changes.
//
// When the terminal window is resized, the TUI must notify the PTY so that
// running commands can adapt their output (e.g., adjust progress bars,
// text wrapping, or table formatting).
//
// Parameters:
//   - rows: number of rows (must be > 0)
//   - cols: number of columns (must be > 0)
//
// Example:
//
//	session, _ := StartPTYSession(cmd)
//	// Terminal resized to 80x24
//	if err := session.Resize(24, 80); err != nil {
//	    log.Printf("Failed to resize PTY: %v", err)
//	}
//
// Returns error if:
//   - PTY is nil
//   - rows or cols are <= 0
//   - PTY resize operation fails
func (s *PTYSession) Resize(rows, cols int) error {
	if s.PTY == nil {
		return nil
	}

	// Validate dimensions
	if rows <= 0 || cols <= 0 {
		return nil
	}

	// Set PTY window size
	size := &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	}

	return pty.Setsize(s.PTY, size)
}
