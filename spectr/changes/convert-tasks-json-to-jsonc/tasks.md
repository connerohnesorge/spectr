## 1. Core Parser Changes

- [ ] 1.1 Add `stripJSONComments` function to `internal/parsers/parsers.go`
- [ ] 1.2 Update `ReadTasksJson` to call `stripJSONComments` before unmarshalling
- [ ] 1.3 Update `CountTasks` to check `tasks.jsonc` first, then fall back to `tasks.md` only (ignore `tasks.json`)

## 2. Accept Command Changes

- [ ] 2.1 Add `tasksJSONHeader` constant with comprehensive usage guide (status values, transitions, workflow)
- [ ] 2.2 Rename `writeTasksJSON` to `writeTasksJSONC` and prepend header to output
- [ ] 2.3 Update file path from `tasks.json` to `tasks.jsonc` in `processChange`
- [ ] 2.4 Update output messages to reference `tasks.jsonc`

## 3. Testing

- [ ] 3.1 Add `TestStripJSONComments` for comment stripping edge cases
- [ ] 3.2 Add `TestReadTasksJsonWithComments` for JSONC parsing
- [ ] 3.3 Add `TestCountTasks_JsoncPreferredOverMarkdown` to verify priority order
- [ ] 3.4 Add `TestCountTasks_IgnoresLegacyJson` to verify tasks.json is ignored
- [ ] 3.5 Update `cmd/accept_test.go` references from `tasks.json` to `tasks.jsonc`

## 4. Documentation Updates

- [ ] 4.1 Update `spectr/AGENTS.md` references
- [ ] 4.2 Update spec files (`cli-interface`, `archive-workflow`, `agent-instructions`)
- [ ] 4.3 Update template files in `internal/initialize/templates/`
- [ ] 4.4 Update agent command files (`.claude/`, `.agent/`, `.gemini/`, `.opencode/`)
- [ ] 4.5 Update root documentation (`README.md`, `AGENTS.md`, `CLAUDE.md`, `CRUSH.md`)

## 5. Validation

- [ ] 5.1 Run `spectr validate convert-tasks-json-to-jsonc --strict`
- [ ] 5.2 Run all tests with `go test ./...`
