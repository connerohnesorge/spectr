# CLI Interface Delta Spec

## ADDED Requirements

### Requirement: Pre-Command Task Sync Hook

The system SHALL synchronize task statuses from `tasks.jsonc` to `tasks.md`
before every command execution using Kong's BeforeRun hook pattern.

#### Scenario: Sync runs before command execution

- **WHEN** any spectr subcommand is invoked
- **THEN** the system SHALL sync all active changes' task statuses before
  executing the command

#### Scenario: Sync updates only status markers

- **WHEN** syncing tasks.md from tasks.jsonc
- **THEN** the system SHALL update only checkbox markers (`[ ]` to `[x]` or
  vice versa)
- **AND** preserve all other markdown content (comments, links, formatting)

#### Scenario: Sync matches tasks by ID

- **WHEN** matching tasks between files
- **THEN** the system SHALL match by task ID (e.g., `1.1`, `2.3`)
- **AND** handle flexible ID formats (decimal, dot-suffixed, number-only)

#### Scenario: Sync handles missing files gracefully

- **WHEN** tasks.jsonc does not exist for a change
- **THEN** the system SHALL skip sync for that change silently

#### Scenario: Sync handles missing tasks.md gracefully

- **WHEN** tasks.md does not exist but tasks.jsonc does
- **THEN** the system SHALL skip sync for that change silently

#### Scenario: tasks.jsonc is source of truth for status

- **WHEN** a task status differs between files
- **THEN** the system SHALL use the status from tasks.jsonc
- **AND** update tasks.md to match

### Requirement: Global No-Sync Flag

The system SHALL provide a `--no-sync` global flag to disable automatic task
synchronization.

#### Scenario: Skip sync with flag

- **WHEN** `spectr --no-sync <command>` is invoked
- **THEN** the system SHALL skip the pre-command sync subroutine

#### Scenario: Flag applies to all subcommands

- **WHEN** `--no-sync` is provided
- **THEN** the system SHALL skip sync regardless of which subcommand follows

### Requirement: Silent Sync by Default

The system SHALL perform sync operations silently unless verbose mode is
enabled or errors occur.

#### Scenario: No output on successful sync

- **WHEN** sync completes successfully
- **THEN** the system SHALL produce no output

#### Scenario: Error output on sync failure

- **WHEN** sync encounters an error
- **THEN** the system SHALL print the error to stderr
- **AND** continue with command execution (non-blocking)

#### Scenario: Verbose output with flag

- **WHEN** `--verbose` flag is set and sync makes changes
- **THEN** the system SHALL print "Synced N task statuses in [change-id]"

### Requirement: Global Verbose Flag

The system SHALL provide a `--verbose` global flag to enable detailed output
including sync operations.

#### Scenario: Verbose flag registration

- **WHEN** CLI is initialized
- **THEN** the system SHALL register `--verbose` as a global flag

#### Scenario: Verbose sync output

- **WHEN** `spectr --verbose <command>` syncs tasks
- **THEN** the system SHALL print sync details for each affected change

### Requirement: Active Changes Only Sync

The system SHALL only synchronize task files in active (non-archived) changes.

#### Scenario: Exclude archived changes

- **WHEN** discovering changes to sync
- **THEN** the system SHALL exclude `spectr/changes/archive/` subdirectories

#### Scenario: Include all active changes

- **WHEN** discovering changes to sync
- **THEN** the system SHALL include all directories in `spectr/changes/` with
  `tasks.jsonc`

### Requirement: Status Mapping for Sync

The system SHALL map task statuses between jsonc and markdown formats
correctly.

#### Scenario: Pending and in_progress map to unchecked

- **WHEN** status is `pending` or `in_progress` in tasks.jsonc
- **THEN** the system SHALL write `[ ]` in tasks.md

#### Scenario: Completed maps to checked

- **WHEN** status is `completed` in tasks.jsonc
- **THEN** the system SHALL write `[x]` in tasks.md
