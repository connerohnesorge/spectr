# Implementation Tasks

## 1. Core Accept Command

- [x] 1.1 Create `cmd/accept.go` with `AcceptCmd` struct following Kong patterns
- [x] 1.2 Add `AcceptCmd` to `CLI` struct in `cmd/root.go`
- [x] 1.3 Implement `Run()` method that validates change exists
- [x] 1.4 Implement `parseTasksMd()` function to parse tasks.md into structured
  format
- [x] 1.5 Implement `writeTasksJson()` function to write tasks.json with proper
  schema
- [x] 1.6 Add change validation step before conversion (reuse existing
  validation)
- [x] 1.7 Remove tasks.md after successful tasks.json creation
- [x] 1.8 Add `--dry-run` flag to preview conversion without writing files

## 2. JSON Schema and Types

- [x] 2.1 Define `TasksFile` struct with version field and tasks array
- [x] 2.2 Define `Task` struct with id, section, description, and status fields
- [x] 2.3 Add JSON struct tags for proper serialization
- [x] 2.4 Add `TaskStatus` type with `pending`, `in_progress`, `completed`
  values

## 3. Parser Updates

- [x] 3.1 Update `parsers.CountTasks()` to check for `tasks.json` first
- [x] 3.2 Add `parsers.ReadTasksJson()` function to read and parse tasks.json
- [x] 3.3 Ensure backward compatibility - fall back to tasks.md if no JSON
  exists
- [x] 3.4 Update `parsers.TaskStatus` to include in_progress count if needed

## 4. Integration with Existing Commands

- [x] 4.1 Update `internal/archive/archiver.go` to use new task reading logic
- [x] 4.2 Update `internal/list/lister.go` to use new task reading logic
- [x] 4.3 Update `internal/view/dashboard.go` to use new task reading logic
- [x] 4.4 Add auto-accept step to archive command when tasks.md still exists
  (with warning)

## 5. Slash Command Updates

- [x] 5.1 Update `.claude/commands/spectr/apply.md` to require `spectr accept`
  first
- [x] 5.2 Update `.gemini/commands/spectr/apply.toml` with accept requirement
- [x] 5.3 Update template files in
  `internal/initialize/templates/tools/slash-apply.md.tmpl`
- [x] 5.4 Add `spectr accept` instructions to AGENTS.md Stage 2 workflow

## 6. Testing

- [x] 6.1 Add unit tests for `parseTasksMd()` with various task formats
- [x] 6.2 Add unit tests for `writeTasksJson()` with expected output
- [x] 6.3 Add unit tests for `parsers.ReadTasksJson()`
- [x] 6.4 Add integration test for full accept workflow
- [x] 6.5 Add test for backward compatibility (tasks.md fallback)
- [x] 6.6 Add test for dry-run mode

## 7. Documentation

- [x] 7.1 Update README with `spectr accept` command documentation
- [x] 7.2 Update CLI reference docs with accept command
- [x] 7.3 Add example tasks.json format to documentation
