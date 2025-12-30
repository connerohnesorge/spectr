# Delta Specification

## ADDED Requirements

### Requirement: PR Proposal Interactive Selection Filters Unmerged Changes

The `spectr pr proposal` command's interactive selection mode SHALL only display
changes that do not already exist on the target branch (main/master), ensuring
users only see changes that genuinely need proposal PRs.

#### Scenario: Interactive list excludes changes on main

- **WHEN** user runs `spectr pr proposal` without a change ID argument
- **AND** some changes in `spectr/changes/` already exist on `origin/main`
- **THEN** only changes NOT present on `origin/main` are displayed in the
  interactive list
- **AND** changes that exist on main are filtered out before display

#### Scenario: All changes already on main

- **WHEN** user runs `spectr pr proposal` without a change ID argument
- **AND** all active changes already exist on `origin/main`
- **THEN** a message is displayed: "No unmerged proposals found. All changes
  already exist on main."
- **AND** the command exits gracefully without entering interactive mode

#### Scenario: No changes exist at all

- **WHEN** user runs `spectr pr proposal` without a change ID argument
- **AND** no changes exist in `spectr/changes/`
- **THEN** a message is displayed: "No changes found."
- **AND** the command exits gracefully

#### Scenario: Explicit change ID bypasses filter

- **WHEN** user runs `spectr pr proposal <change-id>` with an explicit argument
- **THEN** the filter is NOT applied
- **AND** the command proceeds with the specified change ID
- **AND** existing behavior is preserved for direct invocation

#### Scenario: Archive command unaffected

- **WHEN** user runs `spectr pr archive` without a change ID argument
- **THEN** all active changes are displayed in the interactive list
- **AND** no filtering based on main branch presence is applied
- **AND** existing archive behavior is preserved

#### Scenario: Detection uses git ls-tree

- **WHEN** the system checks if a change exists on main
- **THEN** it uses `git ls-tree` to check if `spectr/changes/<change-id>` path
  exists on `origin/main`
- **AND** the check is performed before displaying the interactive list
- **AND** fetch is performed first to ensure refs are current

#### Scenario: Custom base branch respected

- **WHEN** user runs `spectr pr proposal --base develop` without a change ID
- **THEN** the filter checks against `origin/develop` instead of `origin/main`
- **AND** only changes not present on `origin/develop` are displayed
