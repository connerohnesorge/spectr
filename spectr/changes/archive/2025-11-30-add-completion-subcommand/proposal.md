# Change: Add Shell Completion Subcommand

## Why

Shell completion improves CLI usability by allowing users to tab-complete
commands, flags, and arguments. This reduces typing errors and helps users
discover available options without consulting documentation. The kong-completion
library provides a clean integration with Kong's struct-based command
definitions.

## What Changes

- Add `github.com/jotaen/kong-completion` dependency
- Modify `main.go` to use the kong-completion registration pattern (init Kong,
  register completions, then parse)
- Add `Completion` subcommand to the CLI struct in `cmd/root.go`
- Implement custom predictors for dynamic arguments (change IDs, spec IDs, item
  types)
- Support bash, zsh, and fish completions (PowerShell not supported by
  kong-completion)

## Impact

- Affected specs: `cli-framework`
- Affected code: `main.go`, `cmd/root.go`, new `cmd/completion.go`
- Non-breaking: adds new subcommand without affecting existing functionality
