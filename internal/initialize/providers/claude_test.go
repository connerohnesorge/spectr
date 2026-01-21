package providers

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/ralph"
)

const (
	testProviderName = "claude"
	testPromptText   = "Test prompt"
)

// TestClaudeProvider_Binary tests ClaudeProvider.Binary() returns testProviderName.
func TestClaudeProvider_Binary(t *testing.T) {
	provider := &ClaudeProvider{}

	got := provider.Binary()
	want := testProviderName

	if got != want {
		t.Errorf("ClaudeProvider.Binary() = %q, want %q", got, want)
	}
}

// TestClaudeProvider_InvokeTask tests ClaudeProvider.InvokeTask() creates a valid command.
func TestClaudeProvider_InvokeTask(t *testing.T) {
	tests := []struct {
		name         string
		task         ralph.Task
		prompt       string
		wantErr      bool
		checkContext bool
	}{
		{
			name: "basic task with simple prompt",
			task: ralph.Task{
				ID:          "1.1",
				Section:     "Implementation",
				Description: "Add new feature",
				Status:      "pending",
			},
			prompt:  "# Task: 1.1 - Implementation\n\nAdd new feature",
			wantErr: false,
		},
		{
			name: "task with empty prompt",
			task: ralph.Task{
				ID:          "2.3",
				Section:     "Testing",
				Description: "Write tests",
				Status:      "in_progress",
			},
			prompt:  "",
			wantErr: false,
		},
		{
			name: "task with multiline prompt",
			task: ralph.Task{
				ID:          "3.5",
				Section:     "Refactoring",
				Description: "Clean up code",
				Status:      "pending",
			},
			prompt:  "# Task: 3.5 - Refactoring\n\nClean up code.\n\nContext:\n- Remove old code\n- Add comments\n- Update tests",
			wantErr: false,
		},
		{
			name: "task with special characters in prompt",
			task: ralph.Task{
				ID:          "4.2",
				Section:     "Documentation",
				Description: "Update README",
				Status:      "pending",
			},
			prompt:  "# Task: 4.2\n\nUpdate README with `code` and \"quotes\" and 'apostrophes'",
			wantErr: false,
		},
		{
			name: "task with unicode in prompt",
			task: ralph.Task{
				ID:          "5.1",
				Section:     "Internationalization",
				Description: "Add i18n",
				Status:      "pending",
			},
			prompt:  "# Task: 5.1\n\nAdd support for 日本語, العربية, and emoji 🚀",
			wantErr: false,
		},
		{
			name: "task with very long prompt",
			task: ralph.Task{
				ID:          "6.1",
				Section:     "Analysis",
				Description: "Large context",
				Status:      "pending",
			},
			prompt:       strings.Repeat("This is a very long prompt. ", 1000),
			wantErr:      false,
			checkContext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &ClaudeProvider{}
			ctx := context.Background()

			cmd, err := provider.InvokeTask(ctx, &tt.task, tt.prompt)

			if (err != nil) != tt.wantErr {
				t.Errorf("ClaudeProvider.InvokeTask() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			if cmd == nil {
				t.Fatal("ClaudeProvider.InvokeTask() returned nil command")
			}

			// Verify command is for testProviderName binary
			if len(cmd.Args) == 0 {
				t.Error("ClaudeProvider.InvokeTask() returned command with no args")

				return
			}

			// The first arg should be the command name (claude)
			if !strings.Contains(cmd.Args[0], testProviderName) {
				t.Errorf(
					"ClaudeProvider.InvokeTask() command args[0] = %q, want containing \"claude\"",
					cmd.Args[0],
				)
			}

			// Verify stdin is configured
			if cmd.Stdin == nil {
				t.Error("ClaudeProvider.InvokeTask() command has nil Stdin")

				return
			}

			// Read from stdin to verify prompt was set correctly
			stdinData, err := io.ReadAll(cmd.Stdin)
			if err != nil {
				t.Errorf("Failed to read from command Stdin: %v", err)

				return
			}

			if string(stdinData) != tt.prompt {
				t.Errorf(
					"ClaudeProvider.InvokeTask() stdin = %q, want %q",
					string(stdinData),
					tt.prompt,
				)
			}

			// Verify command is not started
			if cmd.Process != nil {
				t.Error("ClaudeProvider.InvokeTask() command was already started (Process != nil)")
			}

			// Verify stdout/stderr are not set (orchestrator will attach PTY)
			if cmd.Stdout != nil {
				t.Error(
					"ClaudeProvider.InvokeTask() command has Stdout set (should be nil for PTY attachment)",
				)
			}
			if cmd.Stderr != nil {
				t.Error(
					"ClaudeProvider.InvokeTask() command has Stderr set (should be nil for PTY attachment)",
				)
			}
		})
	}
}

// TestClaudeProvider_InvokeTask_ContextCancellation tests context cancellation handling.
func TestClaudeProvider_InvokeTask_ContextCancellation(t *testing.T) {
	provider := &ClaudeProvider{}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	task := ralph.Task{
		ID:          "1.1",
		Section:     "Test",
		Description: "Test task",
		Status:      "pending",
	}
	prompt := testPromptText

	cmd, err := provider.InvokeTask(ctx, &task, prompt)

	// The command should still be created (context is for the command execution, not creation)
	if err != nil {
		t.Errorf("ClaudeProvider.InvokeTask() with cancelled context error = %v, want nil", err)
	}

	if cmd == nil {
		t.Fatal("ClaudeProvider.InvokeTask() returned nil command")
	}

	// Verify the command uses the cancelled context
	if cmd.Cancel == nil {
		t.Error("ClaudeProvider.InvokeTask() command does not have Cancel function set")
	}

	// Starting the command should fail because context is cancelled
	err = cmd.Start()
	if err != nil {
		return
	}

	// Clean up if it somehow started
	_ = cmd.Process.Kill()
	_ = cmd.Wait()
	t.Error("ClaudeProvider.InvokeTask() command started despite cancelled context")
}

// TestClaudeProvider_InvokeTask_CommandConfiguration tests command configuration details.
func TestClaudeProvider_InvokeTask_CommandConfiguration(t *testing.T) {
	provider := &ClaudeProvider{}
	ctx := context.Background()

	task := ralph.Task{
		ID:          "1.1",
		Section:     "Test",
		Description: "Test configuration",
		Status:      "pending",
	}
	prompt := "Configuration test prompt"

	cmd, err := provider.InvokeTask(ctx, &task, prompt)
	if err != nil {
		t.Fatalf("ClaudeProvider.InvokeTask() error = %v, want nil", err)
	}

	// Verify command path is set
	if cmd.Path == "" {
		// Path might not be set until Start() is called, check Args instead
		if len(cmd.Args) == 0 {
			t.Error("ClaudeProvider.InvokeTask() command has no Path or Args")
		}
	}

	// Verify command can be started (if claude binary is available)
	// We don't actually start it in tests, but we can verify the structure is correct
	if len(cmd.Args) == 0 {
		t.Error("ClaudeProvider.InvokeTask() command has nil or empty Args")
	}

	// Verify stdin is a strings.Reader (or compatible)
	if cmd.Stdin == nil {
		t.Fatal("ClaudeProvider.InvokeTask() command has nil Stdin")
	}

	// Verify we can read the prompt from stdin
	stdinBytes, err := io.ReadAll(cmd.Stdin)
	if err != nil {
		t.Fatalf("Failed to read from command Stdin: %v", err)
	}

	if string(stdinBytes) != prompt {
		t.Errorf(
			"ClaudeProvider.InvokeTask() stdin content = %q, want %q",
			string(stdinBytes),
			prompt,
		)
	}
}

// TestClaudeProvider_InvokeTask_TaskParameterUnused tests that the task parameter can be anything.
func TestClaudeProvider_InvokeTask_TaskParameterUnused(t *testing.T) {
	provider := &ClaudeProvider{}
	ctx := context.Background()
	prompt := testPromptText

	// Test with various task configurations to ensure task parameter is truly unused
	tasks := []ralph.Task{
		{ID: "", Section: "", Description: "", Status: ""},
		{ID: "1.1", Section: "Test", Description: "Normal task", Status: "pending"},
		{
			ID:          "999.999",
			Section:     "Very Long Section Name Here",
			Description: strings.Repeat("x", 1000),
			Status:      "completed",
		},
		{ID: "special!@#$", Section: "特殊文字", Description: "Unicode 🚀", Status: "invalid"},
	}

	for i, task := range tasks {
		t.Run(task.ID, func(t *testing.T) {
			cmd, err := provider.InvokeTask(ctx, &task, prompt)
			if err != nil {
				t.Errorf("Test %d: ClaudeProvider.InvokeTask() error = %v, want nil", i, err)

				return
			}

			if cmd == nil {
				t.Errorf("Test %d: ClaudeProvider.InvokeTask() returned nil command", i)

				return
			}

			// Verify command is the same regardless of task parameters
			if !strings.Contains(cmd.Args[0], testProviderName) {
				t.Errorf(
					"Test %d: ClaudeProvider.InvokeTask() command = %v, want claude command",
					i,
					cmd.Args[0],
				)
			}
		})
	}
}

// TestClaudeProvider_InvokeTask_PromptEdgeCases tests prompt edge cases.
func TestClaudeProvider_InvokeTask_PromptEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		prompt string
	}{
		{"empty prompt", ""},
		{"whitespace only", "   \t\n  "},
		{"null byte", "prompt\x00with\x00nulls"},
		{"binary data", "\x00\x01\x02\x03\xff\xfe\xfd"},
		{"very long line", strings.Repeat("a", 100000)},
		{"many newlines", strings.Repeat("\n", 1000)},
		{"mixed line endings", "line1\nline2\r\nline3\rline4"},
		{"unicode normalization", "café vs café"}, // Different unicode representations
		{"emoji", "🚀🎉💻🔥⚡"},
		{"right-to-left", "مرحبا بك في البرنامج"},
		{"control characters", "\t\r\n\x1b[31mRed\x1b[0m"},
	}

	provider := &ClaudeProvider{}
	ctx := context.Background()
	task := ralph.Task{
		ID:          "1.1",
		Section:     "Test",
		Description: "Edge case test",
		Status:      "pending",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := provider.InvokeTask(ctx, &task, tt.prompt)
			if err != nil {
				t.Errorf("ClaudeProvider.InvokeTask() error = %v, want nil", err)

				return
			}

			if cmd == nil {
				t.Fatal("ClaudeProvider.InvokeTask() returned nil command")
			}

			// Verify stdin contains the exact prompt
			stdinData, err := io.ReadAll(cmd.Stdin)
			if err != nil {
				t.Errorf("Failed to read stdin: %v", err)

				return
			}

			if string(stdinData) != tt.prompt {
				t.Errorf(
					"ClaudeProvider.InvokeTask() stdin = %q, want %q",
					string(stdinData),
					tt.prompt,
				)
			}
		})
	}
}

// TestClaudeProvider_InvokeTask_ConcurrentCalls tests concurrent invocations.
func TestClaudeProvider_InvokeTask_ConcurrentCalls(t *testing.T) {
	provider := &ClaudeProvider{}
	ctx := context.Background()

	const numGoroutines = 10
	const callsPerGoroutine = 100

	done := make(chan error, numGoroutines)

	for i := range numGoroutines {
		go func(_ int) {
			for range callsPerGoroutine {
				task := ralph.Task{
					ID:          "1.1",
					Section:     "Concurrent",
					Description: "Concurrent test",
					Status:      "pending",
				}
				prompt := "Concurrent test prompt"

				cmd, err := provider.InvokeTask(ctx, &task, prompt)
				if err != nil {
					done <- err

					return
				}

				if cmd == nil {
					done <- context.Canceled

					return
				}

				// Verify stdin is readable
				if cmd.Stdin == nil {
					done <- context.Canceled

					return
				}
			}
			done <- nil
		}(i)
	}

	// Wait for all goroutines
	for range numGoroutines {
		if err := <-done; err != nil {
			t.Errorf("Goroutine failed: %v", err)
		}
	}
}

// TestClaudeProvider_InvokeTask_ReturnsNewCommand tests that each call returns a new command instance.
func TestClaudeProvider_InvokeTask_ReturnsNewCommand(t *testing.T) {
	provider := &ClaudeProvider{}
	ctx := context.Background()
	task := ralph.Task{
		ID:          "1.1",
		Section:     "Test",
		Description: "Instance test",
		Status:      "pending",
	}
	prompt := testPromptText

	cmd1, err1 := provider.InvokeTask(ctx, &task, prompt)
	if err1 != nil {
		t.Fatalf("First InvokeTask() error = %v, want nil", err1)
	}

	cmd2, err2 := provider.InvokeTask(ctx, &task, prompt)
	if err2 != nil {
		t.Fatalf("Second InvokeTask() error = %v, want nil", err2)
	}

	// Verify they are different instances
	if cmd1 == cmd2 {
		t.Error(
			"ClaudeProvider.InvokeTask() returned same command instance, want new instance each call",
		)
	}

	// Verify both have independent stdin
	stdin1Data, err := io.ReadAll(cmd1.Stdin)
	if err != nil {
		t.Fatalf("Failed to read cmd1 stdin: %v", err)
	}

	stdin2Data, err := io.ReadAll(cmd2.Stdin)
	if err != nil {
		t.Fatalf("Failed to read cmd2 stdin: %v", err)
	}

	if string(stdin1Data) != prompt {
		t.Errorf("cmd1 stdin = %q, want %q", string(stdin1Data), prompt)
	}

	if string(stdin2Data) != prompt {
		t.Errorf("cmd2 stdin = %q, want %q", string(stdin2Data), prompt)
	}
}

// BenchmarkClaudeProvider_InvokeTask benchmarks the InvokeTask method.
func BenchmarkClaudeProvider_InvokeTask(b *testing.B) {
	provider := &ClaudeProvider{}
	ctx := context.Background()
	task := ralph.Task{
		ID:          "1.1",
		Section:     "Benchmark",
		Description: "Performance test",
		Status:      "pending",
	}
	prompt := "# Task: 1.1 - Benchmark\n\nPerformance test"

	b.ResetTimer()
	for range b.N {
		cmd, err := provider.InvokeTask(ctx, &task, prompt)
		if err != nil {
			b.Fatalf("InvokeTask() error = %v", err)
		}
		if cmd == nil {
			b.Fatal("InvokeTask() returned nil command")
		}
	}
}

// BenchmarkClaudeProvider_InvokeTask_LargePrompt benchmarks with a large prompt.
func BenchmarkClaudeProvider_InvokeTask_LargePrompt(b *testing.B) {
	provider := &ClaudeProvider{}
	ctx := context.Background()
	task := ralph.Task{
		ID:          "1.1",
		Section:     "Benchmark",
		Description: "Large prompt test",
		Status:      "pending",
	}
	// Create a 1MB prompt
	prompt := strings.Repeat("This is a large prompt for benchmarking. ", 25000)

	b.ResetTimer()
	for range b.N {
		cmd, err := provider.InvokeTask(ctx, &task, prompt)
		if err != nil {
			b.Fatalf("InvokeTask() error = %v", err)
		}
		if cmd == nil {
			b.Fatal("InvokeTask() returned nil command")
		}
	}
}
