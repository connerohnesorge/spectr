## ADDED Requirements

### Requirement: List Command for Changes

The system SHALL provide a `list` command that enumerates all active changes in the project, displaying their IDs by default.

#### Scenario: List changes with IDs only

- **WHEN** user runs `spectr list` without flags
- **THEN** the system displays change IDs, one per line, sorted alphabetically
- **AND** excludes archived changes in the `archive/` directory

#### Scenario: List changes with details

- **WHEN** user runs `spectr list --long`
- **THEN** the system displays each change with format: `{id}: {title} [deltas {count}] [tasks {completed}/{total}]`
- **AND** sorts output alphabetically by ID

#### Scenario: List changes as JSON

- **WHEN** user runs `spectr list --json`
- **THEN** the system outputs a JSON array of objects with fields: `id`, `title`, `deltaCount`, `taskStatus` (with `total` and `completed`)
- **AND** sorts the array by ID

#### Scenario: No changes found

- **WHEN** user runs `spectr list` and no active changes exist
- **THEN** the system displays "No items found"

### Requirement: List Command for Specs

The system SHALL support a `--specs` flag that switches the list command to enumerate specifications instead of changes.

#### Scenario: List specs with IDs only

- **WHEN** user runs `spectr list --specs` without other flags
- **THEN** the system displays spec IDs, one per line, sorted alphabetically
- **AND** only includes directories with valid `spec.md` files

#### Scenario: List specs with details

- **WHEN** user runs `spectr list --specs --long`
- **THEN** the system displays each spec with format: `{id}: {title} [requirements {count}]`
- **AND** sorts output alphabetically by ID

#### Scenario: List specs as JSON

- **WHEN** user runs `spectr list --specs --json`
- **THEN** the system outputs a JSON array of objects with fields: `id`, `title`, `requirementCount`
- **AND** sorts the array by ID

#### Scenario: No specs found

- **WHEN** user runs `spectr list --specs` and no specs exist
- **THEN** the system displays "No items found"

### Requirement: Change Discovery

The system SHALL discover active changes by scanning the `spectr/changes/` directory and identifying subdirectories that contain a `proposal.md` file, excluding the `archive/` directory.

#### Scenario: Find active changes

- **WHEN** the system scans for changes
- **THEN** it includes all subdirectories of `spectr/changes/` that contain `proposal.md`
- **AND** excludes the `spectr/changes/archive/` directory and its contents
- **AND** excludes hidden directories (starting with `.`)

### Requirement: Spec Discovery

The system SHALL discover specs by scanning the `spectr/specs/` directory and identifying subdirectories that contain a `spec.md` file.

#### Scenario: Find specs

- **WHEN** the system scans for specs
- **THEN** it includes all subdirectories of `spectr/specs/` that contain `spec.md`
- **AND** excludes hidden directories (starting with `.`)

### Requirement: Title Extraction

The system SHALL extract titles from proposal and spec markdown files by finding the first level-1 heading and removing the "Change:" or "Spec:" prefix if present.

#### Scenario: Extract title from proposal

- **WHEN** the system reads a `proposal.md` file with heading `# Change: Add Feature`
- **THEN** it extracts the title as "Add Feature"

#### Scenario: Extract title from spec

- **WHEN** the system reads a `spec.md` file with heading `# CLI Framework`
- **THEN** it extracts the title as "CLI Framework"

#### Scenario: Fallback to ID when title not found

- **WHEN** the system cannot extract a title from a markdown file
- **THEN** it uses the directory name (ID) as the title

### Requirement: Task Counting

The system SHALL count tasks in `tasks.md` files by identifying lines matching the pattern `- [ ]` or `- [x]` (case-insensitive), with completed tasks marked by `[x]`.

#### Scenario: Count completed and total tasks

- **WHEN** the system reads a `tasks.md` file with 3 tasks, 2 marked `[x]` and 1 marked `[ ]`
- **THEN** it reports `taskStatus` as `{ total: 3, completed: 2 }`

#### Scenario: Handle missing tasks file

- **WHEN** the system cannot find or read a `tasks.md` file for a change
- **THEN** it reports `taskStatus` as `{ total: 0, completed: 0 }`
- **AND** continues processing without error
