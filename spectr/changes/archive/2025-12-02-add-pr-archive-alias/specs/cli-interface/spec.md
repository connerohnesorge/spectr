# Delta Specification

## ADDED Requirements

### Requirement: PR Archive Subcommand Alias

The `spectr pr archive` subcommand SHALL support `a` as a shorthand alias,
allowing users to invoke `spectr pr a <id>` as equivalent to `spectr pr archive
<id>`.

#### Scenario: User runs spectr pr a shorthand

- **WHEN** user runs `spectr pr a <change-id>`
- **THEN** the system executes the archive PR workflow identically to `spectr pr
  archive`
- **AND** all flags (`--base`, `--draft`, `--force`, `--dry-run`,
  `--skip-specs`) work with the alias

#### Scenario: User runs spectr pr a with flags

- **WHEN** user runs `spectr pr a my-change --draft --force`
- **THEN** the command behaves identically to `spectr pr archive my-change
  --draft --force`
- **AND** a draft PR is created after deleting any existing branch

#### Scenario: Help text shows archive alias

- **WHEN** user runs `spectr pr --help`
- **THEN** the help text displays `archive` with its `a` alias
- **AND** the alias is shown in parentheses or as comma-separated alternatives
