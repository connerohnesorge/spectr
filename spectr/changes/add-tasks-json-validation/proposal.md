# Change: Add tasks.json Validation

## Why

Go's `json.Unmarshal` accepts any string value for `TaskStatusValue` fields because it's a string type alias. This means invalid status values like `"done"` (instead of `"completed"`) silently pass through parsing without any error. The `countTasksFromJson` function won't count these invalid statuses (they fall through the switch statement), leading to incorrect task counts and no user feedback about malformed files.

## What Changes

- Add `IsValid()` method to `TaskStatusValue` type for status validation
- Add `ValidateTasksJson()` function in parsers package to validate tasks.json files
- Integrate tasks.json validation into `spectr validate` command for changes
- Report invalid status values as validation errors with task ID context

## Impact

- Affected specs: `cli-framework` (validate command)
- Affected code:
  - `internal/parsers/types.go` - Add validation method
  - `internal/parsers/parsers.go` - Add validation function
  - `internal/validation/change_rules.go` - Integrate validation
  - `internal/parsers/parsers_test.go` - Add tests
