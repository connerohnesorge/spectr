# Change: Convert tasks.json to tasks.jsonc with Header Comments

## Why

The current `tasks.json` format lacks context for users unfamiliar with Spectr. Adding header comments explains available task status values and usage instructions, improving the developer experience without requiring external documentation lookup.

## What Changes

- The `spectr accept` command outputs `tasks.jsonc` instead of `tasks.json`
- Output includes header comments documenting task status values (`pending`, `in_progress`, `completed`)
- The parser strips JSONC comments before unmarshalling for backward compatibility
- Existing `tasks.json` files remain readable via fallback logic

## Impact

- Affected specs: cli-interface
- Affected code:
  - `cmd/accept.go` - Write tasks.jsonc with header comments
  - `internal/parsers/parsers.go` - Add comment stripping, update file discovery
  - `internal/parsers/parsers_test.go` - Add JSONC parsing tests
  - `cmd/accept_test.go` - Update tests for new file extension
  - Documentation and agent instruction files referencing tasks.json
