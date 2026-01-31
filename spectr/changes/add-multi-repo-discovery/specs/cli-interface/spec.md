# Delta Spec: CLI Interface - Multi-Repo Discovery

## ADDED Requirements

### Requirement: Multi-Root Discovery

The system SHALL discover all `spectr/` directories by walking up the directory
tree from the current working directory, stopping at git repository boundaries.

#### Scenario: Single spectr root in cwd

- **WHEN** user runs spectr from a directory with `spectr/` as direct child
- **THEN** discover that single root
- **AND** behave identically to current (no prefix in output)

#### Scenario: Multiple spectr roots in parent chain

- **WHEN** user runs spectr from `mono/project/src/`
- **AND** `mono/spectr/` and `mono/project/spectr/` both exist
- **THEN** discover both roots
- **AND** aggregate results from all roots

#### Scenario: Git boundary stops discovery

- **WHEN** walking up from cwd
- **AND** a `.git` directory is encountered
- **THEN** stop discovery at that git root
- **AND** do not traverse beyond the git repository boundary

#### Scenario: No spectr root found

- **WHEN** no `spectr/` directory exists in the path up to git root
- **THEN** return empty list (silent, no error)
- **AND** commands proceed with empty data (current behavior preserved)

### Requirement: SPECTR_ROOT Environment Variable

The system SHALL support a `SPECTR_ROOT` environment variable that overrides
automatic discovery with an explicit spectr root path.

#### Scenario: Env var set to valid path

- **WHEN** `SPECTR_ROOT` is set to `/path/to/project`
- **AND** `/path/to/project/spectr/` exists
- **THEN** use only that root (skip automatic discovery)

#### Scenario: Env var set to invalid path

- **WHEN** `SPECTR_ROOT` is set to a path without `spectr/` directory
- **THEN** emit error "SPECTR_ROOT path does not contain spectr/ directory"
- **AND** exit non-zero

#### Scenario: Env var not set

- **WHEN** `SPECTR_ROOT` is not set
- **THEN** use automatic multi-root discovery

### Requirement: Aggregated Command Output

The system SHALL aggregate results from all discovered spectr roots in command
output, prefixing items with their source root when multiple roots exist.

#### Scenario: List with multiple roots

- **WHEN** `spectr list` runs with multiple discovered roots
- **THEN** display all changes/specs from all roots
- **AND** prefix each item with relative path: `[project] add-feature`

#### Scenario: List with single root

- **WHEN** `spectr list` runs with single discovered root
- **THEN** display items without prefix (backward compatible)

#### Scenario: View dashboard with multiple roots

- **WHEN** `spectr view` runs with multiple discovered roots
- **THEN** show aggregated summary stats
- **AND** group items by root in display

#### Scenario: Validate across roots

- **WHEN** `spectr validate` runs with multiple discovered roots
- **THEN** validate items in all roots
- **AND** prefix validation results with root path

## MODIFIED Requirements

### Requirement: Clipboard Copy on Selection

The system SHALL copy the item path relative to cwd on Enter key press in
interactive list mode.

#### Scenario: Copy relative path

- **WHEN** user presses Enter on a selected item in TUI
- **THEN** copy path relative to cwd to clipboard
- **AND** path format: `spectr/changes/<id>/proposal.md` for changes
- **AND** path format: `spectr/specs/<id>/spec.md` for specs

#### Scenario: Copy path for nested project

- **WHEN** user is in `mono/project/src/`
- **AND** selects item from `mono/project/spectr/`
- **THEN** copy `../spectr/changes/<id>/proposal.md`

#### Scenario: Copy path for parent project

- **WHEN** user is in `mono/project/src/`
- **AND** selects item from `mono/spectr/`
- **THEN** copy `../../spectr/changes/<id>/proposal.md`
