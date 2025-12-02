## 1. Parser Infrastructure

- [ ] 1.1 Add `ValidateTasksStructure()` function to `internal/parsers/parsers.go`
- [ ] 1.2 Implement numbered section header detection (`## [1-9][0-9]*\.`)
- [ ] 1.3 Track tasks per section and orphaned tasks (tasks not under any section)
- [ ] 1.4 Add unit tests for tasks structure validation

## 2. Validation Integration

- [ ] 2.1 Add `ValidateTasksFile()` to `internal/validation/change_rules.go`
- [ ] 2.2 Integrate tasks validation into `ValidateChangeDeltaSpecs()` workflow
- [ ] 2.3 Report warnings for non-compliant structures (backward compatible)
- [ ] 2.4 Report errors in strict mode for structural violations

## 3. Validation Rules

- [ ] 3.1 Warn if no numbered sections found
- [ ] 3.2 Warn if tasks exist outside numbered sections
- [ ] 3.3 Warn if numbered sections are empty (no tasks)
- [ ] 3.4 Warn if section numbers are not sequential (1, 2, 3...)
- [ ] 3.5 Add tests for each validation rule

## 4. Documentation

- [ ] 4.1 Update `spectr/AGENTS.md` with explicit tasks.md format requirements
- [ ] 4.2 Update `internal/initialize/templates/` if tasks.md template needs adjustment

## 5. Final Validation

- [ ] 5.1 Run `spectr validate add-tasks-structure-validation --strict`
- [ ] 5.2 Run existing tests to ensure no regressions
- [ ] 5.3 Run linter (`golangci-lint run`)
