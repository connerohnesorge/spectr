# Implementation Tasks

Complete the implementation of the `/spectr:next` slash command for automated task execution.

## 1. Core Domain Changes

- [ ] 1.1 Add `SlashNext SlashCommand = iota` to `internal/domain/slashcmd.go`
- [ ] 1.2 Update String() method to return "next" for SlashNext
- [ ] 1.3 Ensure SlashNext has unique value distinct from SlashProposal and SlashApply
- [ ] 1.4 Run domain tests: `go test ./internal/domain/...`

## 2. Task Discovery Implementation

- [ ] 2.1 Create `internal/discovery/task_discovery.go`
- [ ] 2.2 Implement `FindNextPendingTask(changeDir string) (*Task, error)`
- [ ] 2.3 Add support for parsing v1 flat tasks.jsonc files
- [ ] 2.4 Add support for v2 hierarchical files with $ref links
- [ ] 2.5 Implement $ref resolution to follow child task files
- [ ] 2.6 Add circular reference detection
- [ ] 2.7 Create comprehensive tests with sample task files

## 3. Status Management Implementation

- [ ] 3.1 Create `internal/taskexec/status.go`
- [ ] 3.2 Implement `UpdateTaskStatus(changeDir string, taskID string, status string) error`
- [ ] 3.3 Handle both flat v1 and hierarchical v2 file structures
- [ ] 3.4 Ensure atomic writes to tasks.jsonc files
- [ ] 3.5 Implement parent task status aggregation logic
- [ ] 3.6 Add tests for status updates with various task structures

## 4. Template System Updates

- [ ] 4.1 Add `SlashNext()` method to TemplateManager in `internal/initialize/templates.go`
- [ ] 4.2 Create SlashNext template file in `internal/initialize/templates/`
- [ ] 4.3 Include task discovery logic in template
- [ ] 4.4 Include status management workflow
- [ ] 4.5 Add execution reporting format
- [ ] 4.6 Include error handling and recovery guidance

## 5. Provider Integration - Claude Ecosystem

- [ ] 5.1 Update `internal/initialize/providers/claude.go` with SlashNext
- [ ] 5.2 Update `internal/initialize/providers/claude_code.go` with SlashNext
- [ ] 5.3 Update `internal/initialize/providers/cursor.go` with SlashNext
- [ ] 5.4 Update `internal/initialize/providers/codex.go` with SlashNext
- [ ] 5.5 Test `spectr init` generates next.md for each provider

## 6. Provider Integration - Other Providers (Batch 1)

- [ ] 6.1 Update `internal/initialize/providers/continue.go` with SlashNext
- [ ] 6.2 Update `internal/initialize/providers/cline.go` with SlashNext
- [ ] 6.3 Update `internal/initialize/providers/kimi.go` with SlashNext
- [ ] 6.4 Update `internal/initialize/providers/windsurf.go` with SlashNext
- [ ] 6.5 Update `internal/initialize/providers/qoder.go` with SlashNext
- [ ] 6.6 Update `internal/initialize/providers/qwen.go` with SlashNext

## 7. Provider Integration - Other Providers (Batch 2)

- [ ] 7.1 Update `internal/initialize/providers/gemini.go` with SlashNext
- [ ] 7.2 Update `internal/initialize/providers/aider.go` with SlashNext
- [ ] 7.3 Update `internal/initialize/providers/opencode.go` with SlashNext
- [ ] 7.4 Update `internal/initialize/providers/costrict.go` with SlashNext
- [ ] 7.5 Update `internal/initialize/providers/crush.go` with SlashNext
- [ ] 7.6 Update `internal/initialize/providers/antigravity.go` with SlashNext
- [ ] 7.7 Update `internal/initialize/providers/kilocode.go` with SlashNext

## 8. Integration Testing

- [ ] 8.1 Create test proposal with v1 flat tasks.jsonc
- [ ] 8.2 Execute SlashNext and verify correct task discovery
- [ ] 8.3 Verify status updates from pending → in_progress → completed
- [ ] 8.4 Create test proposal with v2 hierarchical tasks.jsonc
- [ ] 8.5 Test $ref resolution to child task files
- [ ] 8.6 Verify parent task status aggregation when children complete

## 9. Error Handling Tests

- [ ] 9.1 Test with all tasks completed (no pending tasks)
- [ ] 9.2 Test with malformed tasks.jsonc
- [ ] 9.3 Test with missing $ref child files
- [ ] 9.4 Verify clear error messages in each scenario
- [ ] 9.5 Verify graceful recovery and status rollback on error

## 10. Documentation and Validation

- [ ] 10.1 Run `spectr validate add-slash-next-command`
- [ ] 10.2 Fix any spec validation errors
- [ ] 10.3 Update provider system spec documentation
- [ ] 10.4 Create usage examples with sample proposals
- [ ] 10.5 Document SlashNext behavior in AGENTS.md

## 11. Final Integration

- [ ] 11.1 Run full test suite: `go test ./...`
- [ ] 11.2 Test `spectr init` in fresh project
- [ ] 11.3 Verify next.md generated for multiple providers
- [ ] 11.4 Review generated command files for accuracy
- [ ] 11.5 Create demo video/GIF of SlashNext usage

## 12. Cleanup and Archive

- [ ] 12.1 Commit all changes with clear messages
- [ ] 12.2 Create PR following spectr workflow
- [ ] 12.3 Include test results in PR description
- [ ] 12.4 Add screenshots of generated slash command files
- [ ] 12.5 Run `spectr pr archive add-slash-next-command` on completion
