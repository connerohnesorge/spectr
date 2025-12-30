# Delta Specification

## ADDED Requirements

### Requirement: List Command Alias

The `spectr list` command SHALL support `ls` as a shorthand alias, allowing
users to invoke `spectr ls` as equivalent to `spectr list`.

#### Scenario: User runs spectr ls shorthand

- **WHEN** user runs `spectr ls`
- **THEN** the system displays the list of changes identically to `spectr list`
- **AND** all flags (`--specs`, `--all`, `--long`, `--json`, `--interactive`)
  work with the alias

#### Scenario: User runs spectr ls with flags

- **WHEN** user runs `spectr ls --specs --long`
- **THEN** the command behaves identically to `spectr list --specs --long`
- **AND** specs are displayed in long format

#### Scenario: Help text shows list alias

- **WHEN** user runs `spectr --help`
- **THEN** the help text displays `list` with its `ls` alias
- **AND** the alias is shown in parentheses or as comma-separated alternatives
