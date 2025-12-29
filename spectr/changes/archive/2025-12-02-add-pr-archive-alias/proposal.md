# Change: Add `a` alias for `spectr pr archive` subcommand

## Why

Users frequently run `spectr pr archive <id>` to create PRs for completed changes. Providing a shorthand `spectr pr a <id>` reduces typing and improves CLI ergonomics, consistent with common CLI patterns where single-letter aliases exist for frequently-used subcommands.

## What Changes

- Add `a` as an alias for the `archive` subcommand under `spectr pr`
- Users can invoke `spectr pr a <id>` as equivalent to `spectr pr archive <id>`
- All existing flags (`--base`, `--draft`, `--force`, `--dry-run`, `--skip-specs`) work with the alias

## Dependencies

- Requires `add-pr-subcommand` to be archived first (PR Command Structure must exist in base specs)

## Impact

- Affected specs: `cli-interface`
- Affected code: `cmd/pr.go` (single line change to add `aliases:"a"` tag)
