// Package track provides automatic git commit tracking for spectr tasks.
//
// This package watches tasks.jsonc for status changes and auto-commits
// related file modifications when a task transitions to "in_progress"
// or "completed" status. It streamlines development workflow by ensuring
// consistent commit history without manual intervention.
//
// # Main Components
//
// The package is organized into three main components:
//
//   - Watcher: fsnotify-based file watcher that monitors tasks.jsonc.
//     Implements debouncing to handle rapid writes from editors.
//
//   - Committer: Git operations handler that stages modified files
//     (excluding task files) and creates commits. Handles edge cases
//     like no files to commit and git failures.
//
//   - Tracker: Main event loop that coordinates Watcher and Committer.
//     Processes task status transitions, manages graceful shutdown,
//     and terminates when all tasks complete or a git commit fails.
//
// # Commit Message Format
//
// Commits follow a consistent message format:
//
//	spectr(<change-id>): start task <task-id>   (in_progress)
//	spectr(<change-id>): complete task <task-id> (completed)
//
// Each commit message includes an automated footer:
//
//	[Automated by spectr track]
//
// # Usage
//
// The track package is invoked via `spectr track [change-id]`:
//
//  1. Locates the tasks.jsonc file for the specified change
//  2. Starts watching for file modifications
//  3. On task status change, stages and commits related changes
//  4. Continues until all tasks complete or user interrupts (Ctrl+C)
//
// # Error Handling
//
// The tracker stops immediately on git commit failure to allow manual
// intervention. When no files are modified besides task files, a warning
// is displayed but tracking continues.
//
// Related error types are defined in the specterrs package:
//
//   - NoTasksFileError: tasks.jsonc not found for change
//   - TasksAlreadyCompleteError: all tasks already completed
//   - TrackInterruptedError: tracking stopped by user interrupt
//   - GitCommitError: git commit operation failed
package track
