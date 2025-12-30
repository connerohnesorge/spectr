# Delta Specification

## ADDED Requirements

### Requirement: Local Cleanup After PR Worktree Operations

The system SHALL clean up the local change directory after successfully
completing a PR worktree operation (archive or remove), preventing orphan files
from remaining in the user's working directory after the PR is merged.

#### Scenario: PR archive cleans local directory

- **WHEN** `spectr pr archive <change-id>` successfully creates a PR
- **THEN** the system SHALL remove the local change directory at
  `spectr/changes/<change-id>/`
- **AND** all files including untracked files (e.g., `tasks.json` from `spectr
  accept`) are removed

#### Scenario: PR remove cleans local directory

- **WHEN** `spectr pr rm <change-id>` successfully creates a PR
- **THEN** the system SHALL remove the local change directory at
  `spectr/changes/<change-id>/`
- **AND** all files including untracked files are removed

#### Scenario: Cleanup skipped on error

- **WHEN** a PR operation fails at any step
- **THEN** the local change directory is NOT removed
- **AND** no cleanup message is displayed

#### Scenario: Cleanup failure is non-fatal

- **WHEN** local cleanup fails (e.g., permission error)
- **THEN** a warning is displayed
- **AND** the command still exits successfully
- **AND** the PR is still created
