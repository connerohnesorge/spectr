# Change: Add `ls` alias for `spectr list` command

## Why

Users familiar with Unix command-line conventions often expect `ls` to list
items. Providing a shorthand `spectr ls` reduces typing and improves CLI
ergonomics, consistent with the existing pattern where PR subcommands have
single-letter aliases (e.g., `spectr pr a` for `spectr pr archive`).

## What Changes

- Add `ls` as an alias for the `list` command
- Users can invoke `spectr ls` as equivalent to `spectr list`
- All existing flags (`--specs`, `--all`, `--long`, `--json`, `--interactive`)
  work with the alias

## Dependencies

None.

## Impact

- Affected specs: `cli-framework`
- Affected code: `cmd/root.go` (single line change to add `aliases:"ls"` tag)
