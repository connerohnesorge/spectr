## 1. Task Parser Implementation

- [ ] 1.1 Create `internal/accept/types.go` with Task, Section, and TasksJSON struct definitions
- [ ] 1.2 Create `internal/accept/parser.go` to parse tasks.md into structured format
  - Parse section headers (`## N. Section Name`)
  - Parse task lines (`- [ ] N.N Description`)
  - Handle unlimited nested subtasks recursively (`- [ ] N.N.N.N...`)
  - Parse indented detail lines and append to task description
  - Preserve completion status from `[x]` markers
- [ ] 1.3 Add comprehensive tests in `internal/accept/parser_test.go`
  - Test section parsing
  - Test task ID extraction
  - Test completion status detection (preserve [x] markers)
  - Test unlimited recursive nesting (1.1.1.1.1...)
  - Test indented detail line parsing and description concatenation
  - Test edge cases (empty sections, no tasks, malformed lines)

## 2. Test Fixtures from Archive

- [ ] 2.1 Identify diverse tasks.md examples from `spectr/changes/archive/`
- [ ] 2.2 Copy selected fixtures to `testdata/tasks/` directory
  - Include simple, medium, and complex examples
  - Include examples with nested tasks
  - Include examples with all tasks completed
- [ ] 2.3 Create expected JSON outputs for each fixture

## 3. JSON Schema Implementation

- [ ] 3.1 Define tasks.json schema with version, metadata, sections, and summary
- [ ] 3.2 Create `internal/accept/writer.go` for JSON generation
  - Atomic write with temp file rename
  - Pretty-print with 2-space indentation
  - Calculate summary totals automatically
- [ ] 3.3 Add writer tests with golden file comparison

## 4. Accept Command Implementation

- [x] 4.1 Create `cmd/accept.go` with AcceptCmd struct
  - ChangeID positional argument (optional for interactive selection)
  - `--yes` flag to skip confirmation prompt
- [x] 4.2 Add AcceptCmd to CLI struct in `cmd/root.go`
- [x] 4.3 Create `internal/accept/accept.go` with core workflow
  - Validate change exists
  - Check tasks.md exists
  - Check tasks.json does not already exist
  - Parse tasks.md
  - Generate tasks.json
  - Remove tasks.md
- [x] 4.4 Add confirmation prompt before conversion (unless --yes)
- [x] 4.5 Add success/error output messages

## 5. Apply Slash Command Integration

- [ ] 5.1 Update `.claude/commands/spectr/apply.md` to integrate acceptance workflow
  - Check if tasks.json exists for the change
  - If missing, instruct agent: "Run `spectr accept <change-id>` first"
  - Update task tracking instructions to use tasks.json instead of tasks.md
  - Clarify that tasks.json is the source of truth after acceptance
- [ ] 5.2 Document the acceptance workflow in AGENTS.md reference

## 6. Archive Workflow Compatibility

- [ ] 6.1 Update `internal/archive/archiver.go` checkTasks() to handle both formats
  - Check for tasks.json first, fall back to tasks.md
  - Parse tasks.json for completion status
- [ ] 6.2 Add JSON task parser to `internal/parsers/parsers.go`
  - CountTasksJSON() function
- [ ] 6.3 Update `internal/list/` and `internal/view/` to support tasks.json

## 7. Validation

- [ ] 7.1 Run `go build` to verify compilation
- [ ] 7.2 Run `go test ./...` to verify all tests pass
- [ ] 7.3 Run `golangci-lint run` to verify no linting errors
- [ ] 7.4 Test `spectr accept` command manually with a test change
- [ ] 7.5 Run `spectr validate add-tasks-json-acceptance --strict`

## 8. Documentation

- [ ] 8.1 Add help text to accept command describing purpose and workflow
- [ ] 8.2 Update CLI help to show accept command
