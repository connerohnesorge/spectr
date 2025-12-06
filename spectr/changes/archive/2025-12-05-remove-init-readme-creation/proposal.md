# Change: Remove automatic README creation during spectr init

## Why

The `spectr init` command currently creates a README.md file automatically when one doesn't exist. This behavior is overly opinionated - most projects either already have a README or have specific preferences about their README content and structure. Auto-generating a README can be annoying for users who want control over their project documentation.

## What Changes

- **Remove automatic README creation** during `spectr init`
- Remove the `createReadmeIfMissing` function from the init executor
- Remove the README creation step from the `Execute` method

## Impact

- Affected specs: `cli-interface`
- Affected code: `internal/initialize/executor.go`
- No breaking changes - users who want a README can create their own
- Simplifies the init process by focusing on Spectr-specific files only
