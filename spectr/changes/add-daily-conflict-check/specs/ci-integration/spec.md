## ADDED Requirements

### Requirement: Daily Conflict Check Workflow
The system SHALL provide a GitHub Action workflow that runs daily at 5 AM UTC to detect conflicts between pending change proposals.

#### Scenario: Scheduled conflict check execution
- **WHEN** the cron schedule triggers at 5 AM UTC
- **THEN** the conflict-check workflow executes automatically
- **AND** the workflow checks out the repository with full history
- **AND** the workflow runs conflict detection against all pending changes

#### Scenario: Manual workflow trigger
- **WHEN** a maintainer manually triggers the conflict-check workflow via workflow_dispatch
- **THEN** the workflow executes immediately with the same logic as scheduled runs
- **AND** the workflow accepts optional inputs to scope the check

### Requirement: Change-to-Change Conflict Detection
The system SHALL detect when multiple pending changes modify the same capability or requirement, indicating potential merge conflicts.

#### Scenario: Capability-level conflict detection
- **WHEN** two or more pending changes contain delta specs for the same capability directory
- **THEN** the system SHALL report a capability-level conflict
- **AND** the report SHALL list all change IDs affecting that capability
- **AND** the report SHALL indicate the conflict type as "capability"

#### Scenario: Requirement-level conflict detection
- **WHEN** two or more pending changes contain delta operations (ADDED/MODIFIED/REMOVED/RENAMED) for the same requirement name
- **THEN** the system SHALL report a requirement-level conflict
- **AND** the report SHALL list the requirement name and all changes touching it
- **AND** the report SHALL indicate what operation each change performs

#### Scenario: No conflicts detected
- **WHEN** all pending changes target distinct capabilities and requirements
- **THEN** the system SHALL report zero conflicts
- **AND** no GitHub issue SHALL be created

### Requirement: Conflict Check CLI Command
The system SHALL provide a `spectr conflicts` CLI command to detect change-to-change conflicts locally or in CI.

#### Scenario: Run conflict detection
- **WHEN** user executes `spectr conflicts`
- **THEN** the system scans all directories in `spectr/changes/` (excluding archive)
- **AND** the system parses delta specs from each change
- **AND** the system compares all changes for overlapping capabilities and requirements
- **AND** results are displayed in human-readable format

#### Scenario: JSON output for CI
- **WHEN** user executes `spectr conflicts --json`
- **THEN** the output SHALL be valid JSON
- **AND** the JSON SHALL include a conflicts array with type, capability, requirement, changes, and operations
- **AND** the JSON SHALL include a summary object with total_conflicts, affected_capabilities, and affected_changes

#### Scenario: Exit code on conflicts
- **WHEN** conflicts are detected
- **THEN** the command SHALL exit with code 1
- **AND** the exit code indicates conflicts were found (not a runtime error)

#### Scenario: Exit code on no conflicts
- **WHEN** no conflicts are detected
- **THEN** the command SHALL exit with code 0

### Requirement: Automatic GitHub Issue Creation
The system SHALL automatically create GitHub issues when the daily conflict check detects conflicts between pending changes.

#### Scenario: Create issue for detected conflicts
- **WHEN** the conflict-check workflow detects one or more conflicts
- **THEN** a GitHub issue SHALL be created with title "Spectr: Conflicting change proposals detected"
- **AND** the issue body SHALL list all detected conflicts with details
- **AND** the issue SHALL be labeled with "spectr-conflict"

#### Scenario: Issue deduplication
- **WHEN** the conflict-check workflow detects conflicts
- **AND** an open issue with label "spectr-conflict" and matching conflict signature already exists
- **THEN** a new issue SHALL NOT be created
- **AND** the workflow SHALL update the existing issue with fresh conflict details

#### Scenario: Issue auto-close on resolution
- **WHEN** the conflict-check workflow runs
- **AND** previously detected conflicts are no longer present
- **AND** an open issue exists for those conflicts
- **THEN** a comment SHALL be added indicating conflicts resolved
- **AND** the issue MAY be automatically closed (configurable)

### Requirement: Conflict Check Workflow Permissions
The system SHALL configure the conflict-check workflow with minimal required permissions to create and update issues.

#### Scenario: Issue write permission
- **WHEN** the conflict-check workflow needs to create or update issues
- **THEN** it SHALL have `issues: write` permission
- **AND** it SHALL have `contents: read` permission for repository checkout
- **AND** no additional secrets SHALL be required beyond GITHUB_TOKEN

### Requirement: Conflict Check Concurrency Management
The system SHALL prevent concurrent conflict-check workflow runs to avoid race conditions in issue management.

#### Scenario: Single concurrent run
- **WHEN** a conflict-check workflow is triggered while another is running
- **THEN** the newer run SHALL wait for the previous run to complete
- **AND** the concurrency group SHALL be scoped to the workflow name
