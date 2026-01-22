// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"context"
	"os/exec"
)

// Ralpher is an interface that providers implement to support task orchestration via spectr ralph.
//
// Providers that implement this interface enable automated task execution where
// the orchestrator feeds tasks from tasks.jsonc files to the provider's CLI,
// managing the execution lifecycle, status polling, and error handling.
//
// The Ralpher interface abstracts the invocation details of different AI agent
// CLIs (Claude, Gemini, Cursor, etc.), allowing the orchestrator to work with
// any provider that can:
//   - Accept a prompt via stdin, file, or command argument
//   - Execute the task and modify the tasks.jsonc status field
//   - Run in a PTY for full terminal emulation
type Ralpher interface {
	// InvokeTask creates an exec.Cmd configured to run the agent CLI with the given task context.
	//
	// The command should be configured to accept the prompt according to the provider's input method:
	//   - Stdin: Configure cmd.Stdin to read from a pipe or buffer
	//   - File: Write prompt to a temp file and pass path as argument
	//   - Argument: Pass prompt content directly as command argument
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - task: The task being executed (ID, section, description, status)
	//   - prompt: Full prompt content with injected context
	//
	// Returns:
	//   - *exec.Cmd: A configured command ready for PTY attachment (not started)
	//   - error: Returns error if the provider cannot be invoked (binary not found, etc.)
	InvokeTask(ctx context.Context, task *Task, prompt string) (*exec.Cmd, error)

	// Binary returns the CLI binary name for the provider (e.g., "claude", "gemini").
	// This is used for display in TUI, error messages, and binary detection.
	Binary() string
}
