# Change: Require TTY for Archive Command

## Why

The `spectr archive` command is a critical operation that moves completed changes
to the archive directory and merges spec deltas into the main specifications.
This operation should only be performed by humans in an interactive terminal
session, not by automated systems or AI agents. Currently, there is no check to
prevent non-human execution of this command, which could lead to unintended
archive operations.

By requiring a TTY (terminal) for the archive command, we ensure that only
humans can execute this critical operation, preventing automated systems from
accidentally or incorrectly archiving changes.

## What Changes

- Add TTY detection check to `spectr archive` command execution
- Return clear error when archive is attempted in non-TTY environment
- Follow existing pattern used in interactive validation mode

## Impact

- Affected specs: `archive-workflow`
- Affected code: `internal/archive/archiver.go`, `internal/archive/cmd.go`
- **BREAKING**: Automated scripts calling `spectr archive` will fail with TTY
  error
- CI/CD pipelines and automated workflows will not be able to run `spectr
  archive`
