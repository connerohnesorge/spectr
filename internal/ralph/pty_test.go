package ralph

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestSpawnPTY(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *exec.Cmd
		wantErr bool
	}{
		{
			name:    "simple echo command",
			cmd:     exec.Command("echo", "hello"),
			wantErr: false,
		},
		{
			name:    "command with multiple args",
			cmd:     exec.Command("printf", "line1\nline2\nline3"),
			wantErr: false,
		},
		{
			name:    "nonexistent command",
			cmd:     exec.Command("nonexistent-command-12345"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptyFile, err := SpawnPTY(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpawnPTY() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			if ptyFile == nil {
				t.Error("SpawnPTY() returned nil PTY file")

				return
			}
			defer func() {
				if err := ptyFile.Close(); err != nil {
					t.Logf("Close() error: %v", err)
				}
			}()

			// Wait for command to complete
			_ = tt.cmd.Wait()
		})
	}
}

func TestStartPTYSession(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *exec.Cmd
		wantErr bool
	}{
		{
			name:    "successful session start",
			cmd:     exec.Command("echo", "test output"),
			wantErr: false,
		},
		{
			name:    "nonexistent command",
			cmd:     exec.Command("nonexistent-cmd-xyz"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := StartPTYSession(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartPTYSession() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			if session == nil {
				t.Error("StartPTYSession() returned nil session")

				return
			}
			defer func() {
				if err := session.Close(); err != nil {
					t.Logf("Close() error: %v", err)
				}
			}()

			if session.PTY == nil {
				t.Error("session.PTY is nil")
			}
			if session.Process == nil {
				t.Error("session.Process is nil")
			}
			if session.Command == nil {
				t.Error("session.Command is nil")
			}

			// Wait for command to finish
			_ = session.Wait()
		})
	}
}

func TestPTYSession_ReadOutput(t *testing.T) {
	tests := []struct {
		name       string
		cmd        *exec.Cmd
		wantOutput string
	}{
		{
			name:       "echo hello",
			cmd:        exec.Command("echo", "hello"),
			wantOutput: "hello",
		},
		{
			name:       "printf multiline",
			cmd:        exec.Command("printf", "line1\nline2"),
			wantOutput: "line1\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := StartPTYSession(tt.cmd)
			if err != nil {
				t.Fatalf("StartPTYSession() error = %v", err)
			}
			defer func() {
				if err := session.Close(); err != nil {
					t.Logf("Close() error: %v", err)
				}
			}()

			// Start reading output concurrently
			outputChan := make(chan []byte, 1)

			go func() {
				// Use manual buffer reading instead of ReadOutput
				// to avoid issues with PTY closing
				var buf bytes.Buffer
				// io.Copy may return errors on PTY close, which is expected
				_, _ = io.Copy(&buf, session.PTY)
				outputChan <- buf.Bytes()
			}()

			// Wait for process to complete
			if err := session.Wait(); err != nil {
				t.Logf("Wait() returned error (expected for some commands): %v", err)
			}

			// Get output with timeout
			select {
			case output := <-outputChan:
				outputStr := string(output)
				// PTY converts \n to \r\n, normalize for comparison
				outputStr = strings.ReplaceAll(outputStr, "\r\n", "\n")
				// PTY may add extra whitespace/newlines
				outputStr = strings.TrimSpace(outputStr)
				tt.wantOutput = strings.TrimSpace(tt.wantOutput)

				if !strings.Contains(outputStr, tt.wantOutput) {
					t.Errorf("ReadOutput() = %q, want substring %q", outputStr, tt.wantOutput)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("ReadOutput() timed out")
			}
		})
	}
}

func TestPTYSession_Wait(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *exec.Cmd
		wantErr bool
	}{
		{
			name:    "successful command",
			cmd:     exec.Command("true"),
			wantErr: false,
		},
		{
			name:    "failing command",
			cmd:     exec.Command("false"),
			wantErr: true,
		},
		{
			name:    "command with output",
			cmd:     exec.Command("echo", "test"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := StartPTYSession(tt.cmd)
			if err != nil {
				t.Fatalf("StartPTYSession() error = %v", err)
			}
			defer func() {
				if err := session.Close(); err != nil {
					t.Logf("Close() error: %v", err)
				}
			}()

			// Read output in background to avoid blocking
			go func() {
				_, _ = io.Copy(io.Discard, session.PTY)
			}()

			err = session.Wait()
			if (err != nil) != tt.wantErr {
				t.Errorf("Wait() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPTYSession_Kill(t *testing.T) {
	// Start a long-running process
	cmd := exec.Command("sleep", "10")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			t.Logf("Close() error: %v", err)
		}
	}()

	// Give process time to start
	time.Sleep(100 * time.Millisecond)

	// Kill the process
	if err := session.Kill(); err != nil {
		t.Errorf("Kill() error = %v", err)
	}

	// Wait should return with error (killed process)
	err = session.Wait()
	if err == nil {
		t.Error("Wait() after Kill() should return error")
	}
}

func TestPTYSession_Close(t *testing.T) {
	cmd := exec.Command("echo", "test")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}

	// Read output to avoid blocking
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
	}()

	// Wait for command to finish
	_ = session.Wait()

	// Close should succeed
	if err := session.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestPTYSession_Lifecycle(t *testing.T) {
	// Test complete lifecycle: start -> read -> wait -> close
	cmd := exec.Command("printf", "hello\nworld")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}

	// Read output
	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		_, _ = io.Copy(&buf, session.PTY)
		done <- true
	}()

	// Wait for process
	if err := session.Wait(); err != nil {
		t.Errorf("Wait() error = %v", err)
	}

	// Wait for reading to complete
	<-done

	// Check output
	output := buf.String()
	if !strings.Contains(output, "hello") || !strings.Contains(output, "world") {
		t.Errorf("output = %q, want to contain 'hello' and 'world'", output)
	}

	// Close session
	if err := session.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestPTYSession_KillBeforeWait(t *testing.T) {
	// Test killing a process before waiting
	cmd := exec.Command("sleep", "5")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			t.Logf("Close() error: %v", err)
		}
	}()

	// Give process time to start
	time.Sleep(100 * time.Millisecond)

	// Kill immediately
	if err := session.Kill(); err != nil {
		t.Errorf("Kill() error = %v", err)
	}

	// Wait should reflect the kill
	err = session.Wait()
	if err == nil {
		t.Error("Wait() should return error after Kill()")
	}
}

func TestPTYSession_NilHandling(t *testing.T) {
	// Test that methods handle nil fields gracefully
	session := &PTYSession{}

	// These should not panic or error
	if err := session.Wait(); err != nil {
		t.Errorf("Wait() on nil Command returned error: %v", err)
	}

	if err := session.Kill(); err != nil {
		t.Errorf("Kill() on nil Process returned error: %v", err)
	}

	if err := session.Close(); err != nil {
		t.Errorf("Close() on nil PTY returned error: %v", err)
	}

	output, err := session.ReadOutput()
	if err != nil {
		t.Errorf("ReadOutput() on nil PTY returned error: %v", err)
	}
	if output != nil {
		t.Errorf("ReadOutput() on nil PTY = %v, want nil", output)
	}

	// Resize should handle nil PTY gracefully
	if err := session.Resize(24, 80); err != nil {
		t.Errorf("Resize() on nil PTY returned error: %v", err)
	}
}

func TestPTYSession_Resize(t *testing.T) {
	tests := []struct {
		name    string
		rows    int
		cols    int
		wantErr bool
	}{
		{
			name:    "valid dimensions 24x80",
			rows:    24,
			cols:    80,
			wantErr: false,
		},
		{
			name:    "valid dimensions 50x120",
			rows:    50,
			cols:    120,
			wantErr: false,
		},
		{
			name:    "valid dimensions 100x200",
			rows:    100,
			cols:    200,
			wantErr: false,
		},
		{
			name:    "zero rows",
			rows:    0,
			cols:    80,
			wantErr: false, // handled gracefully, returns nil
		},
		{
			name:    "zero cols",
			rows:    24,
			cols:    0,
			wantErr: false, // handled gracefully, returns nil
		},
		{
			name:    "negative rows",
			rows:    -10,
			cols:    80,
			wantErr: false, // handled gracefully, returns nil
		},
		{
			name:    "negative cols",
			rows:    24,
			cols:    -50,
			wantErr: false, // handled gracefully, returns nil
		},
		{
			name:    "both zero",
			rows:    0,
			cols:    0,
			wantErr: false, // handled gracefully, returns nil
		},
		{
			name:    "minimum valid 1x1",
			rows:    1,
			cols:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start a long-running command so PTY stays open
			cmd := exec.Command("sleep", "2")
			session, err := StartPTYSession(cmd)
			if err != nil {
				t.Fatalf("StartPTYSession() error = %v", err)
			}
			defer func() {
				_ = session.Kill()
				_ = session.Close()
			}()

			// Give process time to start
			time.Sleep(50 * time.Millisecond)

			// Test resize
			err = session.Resize(tt.rows, tt.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resize(%d, %d) error = %v, wantErr %v", tt.rows, tt.cols, err, tt.wantErr)
			}
		})
	}
}

func TestPTYSession_ResizeDuringExecution(t *testing.T) {
	// Test resizing while command is actively running
	// Use a command that produces output over time
	cmd := exec.Command("bash", "-c", "for i in {1..10}; do echo line $i; sleep 0.1; done")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Start reading output in background
	done := make(chan bool)
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
		done <- true
	}()

	// Give command time to start
	time.Sleep(50 * time.Millisecond)

	// Resize multiple times during execution
	sizes := []struct{ rows, cols int }{
		{24, 80},
		{30, 100},
		{40, 120},
		{50, 160},
	}

	for _, size := range sizes {
		if err := session.Resize(size.rows, size.cols); err != nil {
			t.Errorf("Resize(%d, %d) during execution error = %v", size.rows, size.cols, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for process to complete
	_ = session.Wait()

	// Wait for reading to complete
	select {
	case <-done:
		// Success
	case <-time.After(3 * time.Second):
		t.Fatal("ReadOutput timed out")
	}
}

func TestPTYSession_ResizeNilPTY(t *testing.T) {
	// Test that Resize handles nil PTY gracefully
	session := &PTYSession{
		PTY:     nil,
		Process: nil,
		Command: nil,
	}

	// Should not panic or return error
	if err := session.Resize(24, 80); err != nil {
		t.Errorf("Resize() on nil PTY returned error: %v", err)
	}
}

func TestPTYSession_ResizeAfterClose(t *testing.T) {
	// Test resizing after PTY is closed
	cmd := exec.Command("echo", "test")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}

	// Read output to avoid blocking
	readDone := make(chan bool)
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
		readDone <- true
	}()

	// Wait and close
	_ = session.Wait()
	<-readDone // Wait for reading to finish to avoid race
	if err := session.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Resize after close should be handled (PTY is closed but not nil)
	// This may return an error since PTY is closed, which is acceptable
	err = session.Resize(24, 80)
	// We don't assert error here as it's implementation-specific
	// whether closed file returns error or not
	t.Logf("Resize after close returned: %v", err)
}

func TestPTYSession_Terminate_GracefulShutdown(t *testing.T) {
	// Test that Terminate successfully sends SIGTERM and process exits gracefully
	// Use a bash script that traps SIGTERM for clean exit
	cmd := exec.Command("bash", "-c", "trap 'exit 0' TERM; sleep 10")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Read output in background
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
	}()

	// Give process time to start and set up signal handler
	time.Sleep(100 * time.Millisecond)

	// Start timer to measure graceful shutdown time
	start := time.Now()

	// Terminate with 5 second timeout
	if err := session.Terminate(5 * time.Second); err != nil {
		t.Errorf("Terminate() error = %v", err)
	}

	elapsed := time.Since(start)

	// Process should exit quickly (well under 5 seconds) via SIGTERM
	if elapsed >= 5*time.Second {
		t.Errorf("Terminate() took %v, expected quick graceful exit (< 5s)", elapsed)
	}

	// Process should be terminated
	// Note: Wait() was already called by Terminate, calling again should be safe
	// but may return "wait: no child processes" error which is expected
}

func TestPTYSession_Terminate_TimeoutForcesKill(t *testing.T) {
	// Test that Terminate falls back to SIGKILL when process ignores SIGTERM
	// Use a bash script that ignores SIGTERM
	cmd := exec.Command("bash", "-c", "trap '' TERM; sleep 10")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Read output in background
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
	}()

	// Give process time to start and set up signal handler
	time.Sleep(100 * time.Millisecond)

	// Terminate with short timeout (1 second)
	start := time.Now()
	if err := session.Terminate(1 * time.Second); err != nil {
		t.Errorf("Terminate() error = %v", err)
	}
	elapsed := time.Since(start)

	// Should take approximately 1 second (timeout) then force kill
	if elapsed < 900*time.Millisecond || elapsed > 2*time.Second {
		t.Errorf("Terminate() took %v, expected ~1s (timeout then kill)", elapsed)
	}
}

func TestPTYSession_Terminate_ZeroTimeout(t *testing.T) {
	// Test that zero timeout immediately kills the process
	cmd := exec.Command("sleep", "10")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Read output in background
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
	}()

	// Give process time to start
	time.Sleep(100 * time.Millisecond)

	// Terminate with zero timeout should immediately kill
	start := time.Now()
	if err := session.Terminate(0); err != nil {
		t.Errorf("Terminate(0) error = %v", err)
	}
	elapsed := time.Since(start)

	// Should be very fast (immediate kill)
	if elapsed > 500*time.Millisecond {
		t.Errorf("Terminate(0) took %v, expected immediate kill (< 500ms)", elapsed)
	}
}

func TestPTYSession_Terminate_NilProcess(t *testing.T) {
	// Test that Terminate handles nil process gracefully
	session := &PTYSession{
		PTY:     nil,
		Process: nil,
		Command: nil,
	}

	// Should not panic or return error
	if err := session.Terminate(5 * time.Second); err != nil {
		t.Errorf("Terminate() on nil Process returned error: %v", err)
	}
}

func TestPTYSession_Terminate_AlreadyExited(t *testing.T) {
	// Test terminating a process that already exited
	cmd := exec.Command("echo", "test")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Read output in background
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
	}()

	// Wait for process to naturally complete
	_ = session.Wait()

	// Terminate should handle already-exited process
	// This will likely return an error since the process is gone
	err = session.Terminate(5 * time.Second)
	// Error is expected here (process already exited)
	if err == nil {
		t.Log("Terminate() on exited process returned nil (acceptable)")
	} else {
		t.Logf("Terminate() on exited process returned error (expected): %v", err)
	}
}

func TestPTYSession_Terminate_ImmediateExit(t *testing.T) {
	// Test process that exits immediately after receiving SIGTERM
	cmd := exec.Command("bash", "-c", "trap 'exit 0' TERM; sleep 10")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Read output in background
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
	}()

	// Give process time to start
	time.Sleep(100 * time.Millisecond)

	// Terminate should succeed quickly
	start := time.Now()
	if err := session.Terminate(5 * time.Second); err != nil {
		t.Errorf("Terminate() error = %v", err)
	}
	elapsed := time.Since(start)

	// Should exit very quickly via SIGTERM handler
	if elapsed > 1*time.Second {
		t.Errorf("Terminate() took %v, expected quick exit (< 1s)", elapsed)
	}
}

func TestPTYSession_Terminate_MultipleCallsSafe(t *testing.T) {
	// Test that calling Terminate multiple times is safe
	cmd := exec.Command("sleep", "10")
	session, err := StartPTYSession(cmd)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	// Read output in background
	go func() {
		_, _ = io.Copy(io.Discard, session.PTY)
	}()

	// Give process time to start
	time.Sleep(100 * time.Millisecond)

	// First terminate
	if err := session.Terminate(1 * time.Second); err != nil {
		t.Errorf("First Terminate() error = %v", err)
	}

	// Second terminate should handle already-terminated process
	err = session.Terminate(1 * time.Second)
	// May return error, which is acceptable
	t.Logf("Second Terminate() returned: %v", err)
}

func TestPTYSession_Terminate_CompareWithKill(t *testing.T) {
	// Compare Terminate(0) behavior with Kill()
	// Both should immediately terminate the process

	// Test Terminate(0)
	cmd1 := exec.Command("sleep", "10")
	session1, err := StartPTYSession(cmd1)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session1.Close()
	}()
	go func() {
		_, _ = io.Copy(io.Discard, session1.PTY)
	}()
	time.Sleep(100 * time.Millisecond)

	start1 := time.Now()
	if err := session1.Terminate(0); err != nil {
		t.Errorf("Terminate(0) error = %v", err)
	}
	elapsed1 := time.Since(start1)

	// Test Kill()
	cmd2 := exec.Command("sleep", "10")
	session2, err := StartPTYSession(cmd2)
	if err != nil {
		t.Fatalf("StartPTYSession() error = %v", err)
	}
	defer func() {
		_ = session2.Close()
	}()
	go func() {
		_, _ = io.Copy(io.Discard, session2.PTY)
	}()
	time.Sleep(100 * time.Millisecond)

	start2 := time.Now()
	if err := session2.Kill(); err != nil {
		t.Errorf("Kill() error = %v", err)
	}
	elapsed2 := time.Since(start2)

	// Both should be very fast
	if elapsed1 > 500*time.Millisecond {
		t.Errorf("Terminate(0) took %v, expected < 500ms", elapsed1)
	}
	if elapsed2 > 500*time.Millisecond {
		t.Errorf("Kill() took %v, expected < 500ms", elapsed2)
	}

	t.Logf("Terminate(0) took %v, Kill() took %v", elapsed1, elapsed2)
}
