# Change: Validate Empty tasks.md Files

## Why

Currently, when a `tasks.md` file exists in a change directory but contains no task items (`- [ ]` or `- [x]`), validation passes silently. This allows proposals with empty or incomplete task files to be accepted, which undermines the purpose of requiring task tracking for changes.

## What Changes

- Add validation to detect when `tasks.md` exists but has no task items
- Report an ERROR level validation issue when empty tasks.md is found
- Missing `tasks.md` files are allowed (not an error)

## Impact

- Affected specs: `validation`
- Affected code: `internal/validation/change_rules.go`, `internal/validation/change_rules_test.go`
