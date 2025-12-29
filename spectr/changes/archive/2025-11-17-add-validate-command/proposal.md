# Change: Add Validate Command for Spec and Change Validation

## Why

Spectr currently lacks a mechanism to validate spec and change documents before they are committed or deployed. Without validation, users can create malformed specs with missing required sections, requirements without scenarios, or changes without deltas—leading to confusion, broken tooling, and wasted development time. A comprehensive validate command will catch these errors early, provide actionable feedback, and ensure consistency across all Spectr projects.

## What Changes

- Add new `validate` command to the CLI with support for validating individual items, bulk validation, and interactive selection
- Implement validation engine for spec files (checking Purpose, Requirements, scenarios, SHALL/MUST keywords)
- Implement validation engine for change delta specs (checking ADDED/MODIFIED/REMOVED/RENAMED sections, duplicates, conflicts)
- Add strict mode flag (`--strict`) that treats warnings as errors
- Add JSON output flag (`--json`) for machine-readable validation reports
- Add bulk validation flags (`--all`, `--changes`, `--specs`) for validating multiple items
- Add type disambiguation flag (`--type`) for resolving ambiguous item names
- Implement parallel validation with configurable concurrency for performance
- Provide helpful error messages with remediation guidance

## Impact

- **Affected specs**: `validation` (new capability), `cli-framework` (extends existing)
- **Affected code**:
  - `cmd/root.go` - Add ValidateCmd struct to CLI
  - New `cmd/validate.go` - Validate command implementation
  - New `internal/validation/` - Validation engine package
    - `validator.go` - Core validation logic
    - `types.go` - Validation types and interfaces
    - `spec_validator.go` - Spec-specific validation
    - `change_validator.go` - Change delta validation
    - `parser.go` - Markdown parsing for validation
  - New `internal/discovery/` - Item discovery utilities
    - `discovery.go` - Find changes and specs in project
  - `go.mod` - May need markdown parsing library

## Benefits

- **Quality assurance**: Catch structural and semantic errors before they cause problems
- **Developer experience**: Clear, actionable error messages guide users to fix issues
- **Automation-friendly**: JSON output enables CI/CD integration
- **Performance**: Parallel validation handles large projects efficiently
- **Consistency**: Enforces Spectr conventions across all projects
