# Change: Convert tasks.json to tasks.jsonc with Header Comments

## Why

The current `tasks.json` format lacks context for agents working with Spectr.
Adding comprehensive header comments explains task status values, valid
transitions, and the expected workflow, enabling agents to correctly update
tasks without consulting external documentation.

## What Changes

- The `spectr accept` command outputs `tasks.jsonc` instead of `tasks.json`
- Output includes header comments with a full usage guide:
  - Valid status values: `pending`, `in_progress`, `completed`
  - Status transitions: `pending` → `in_progress` → `completed`
  - Workflow instructions for when/how agents should update task status
- The parser strips JSONC comments before unmarshalling
- Legacy `tasks.json` files are silently ignored (hard break, no backward
  compatibility)

## Impact

- Affected specs: cli-interface
- Affected code:
  - `cmd/accept.go` - Write tasks.jsonc with comprehensive header comments
  - `internal/parsers/parsers.go` - Add comment stripping, read only tasks.jsonc
  - `internal/parsers/parsers_test.go` - Add JSONC parsing tests
  - `cmd/accept_test.go` - Update tests for new file extension
  - Documentation and agent instruction files referencing tasks.json

## Breaking Change

This is a **hard break** from `tasks.json`. Existing `tasks.json` files will be
silently ignored. Projects must re-run `spectr accept` to generate the new
`tasks.jsonc` format.
