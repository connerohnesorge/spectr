# Archive Workflow Specification

## Purpose

This specification defines the workflow for archiving completed changes, including moving them to dated archive directories, validating change structure, parsing delta operations, merging deltas into base specs, and ensuring atomic spec updates.

## Requirements

### Requirement: Change Archive Directory Structure
The system SHALL archive completed changes to `spectr/changes/archive/YYYY-MM-DD-<change-id>/` where YYYY-MM-DD is the current date in ISO 8601 format.

#### Scenario: Archive with date prefix
- **WHEN** archiving a change on 2025-11-17
- **THEN** the change is moved to `spectr/changes/archive/2025-11-17-<change-id>/`

#### Scenario: Prevent duplicate archives
- **WHEN** an archive with the same name already exists
- **THEN** the system returns an error and does not overwrite the existing archive

### Requirement: Pre-Archive Validation
The system SHALL validate changes before archiving to ensure spec consistency.

#### Scenario: Validate proposal structure
- **WHEN** validating a change for archive
- **THEN** the system checks proposal.md structure and reports informational warnings

#### Scenario: Validate delta specs strictly
- **WHEN** validating delta specs
- **THEN** the system checks for required delta operations, scenario formatting, and blocks archive on errors

#### Scenario: Detect duplicate requirements within sections
- **WHEN** a delta spec has duplicate requirement names in the same section
- **THEN** the system returns a validation error with the duplicate requirement name

#### Scenario: Detect cross-section conflicts
- **WHEN** a requirement appears in multiple delta sections (e.g., ADDED and MODIFIED)
- **THEN** the system returns a validation error indicating the conflict

#### Scenario: Validate scenarios are properly formatted
- **WHEN** a requirement lacks properly formatted scenarios
- **THEN** the system returns an error requiring at least one `#### Scenario:` block per requirement

### Requirement: Task Completion Checking
The system SHALL check task completion status and warn users before archiving. The system SHALL read from `tasks.json` when present, falling back to `tasks.md`.

#### Scenario: Display task status from JSON
- **WHEN** archiving a change with `tasks.json`
- **THEN** the system reads task status from JSON file
- **AND** displays task completion status (e.g., "3/5 complete")

#### Scenario: Display task status from Markdown
- **WHEN** archiving a change with only `tasks.md`
- **THEN** the system reads task status from Markdown file
- **AND** displays task completion status (e.g., "3/5 complete")

#### Scenario: Warn on incomplete tasks
- **WHEN** a change has incomplete tasks
- **THEN** the system warns the user and requires confirmation to proceed (unless --yes flag is provided)

#### Scenario: Proceed with incomplete tasks when confirmed
- **WHEN** user confirms archiving despite incomplete tasks
- **THEN** the system proceeds with the archive operation

### Requirement: Delta Spec Discovery
The system SHALL find all delta specifications in a change directory for processing.

#### Scenario: Find delta specs in change
- **WHEN** preparing to archive a change
- **THEN** the system scans `spectr/changes/<id>/specs/*/spec.md` for delta specifications

#### Scenario: Identify corresponding main specs
- **WHEN** delta specs are found
- **THEN** the system maps them to `spectr/specs/*/spec.md` paths based on capability directory name

#### Scenario: Determine create vs update status
- **WHEN** mapping delta specs to main specs
- **THEN** the system checks if each main spec exists and marks it as "create" or "update"

### Requirement: Delta Operation Parsing
The system SHALL parse delta operations from change spec files following a strict format.

#### Scenario: Parse ADDED requirements
- **WHEN** a delta spec contains `## ADDED Requirements` section
- **THEN** the system extracts all requirement blocks with headers and scenarios

#### Scenario: Parse MODIFIED requirements
- **WHEN** a delta spec contains `## MODIFIED Requirements` section
- **THEN** the system extracts complete modified requirement blocks

#### Scenario: Parse REMOVED requirements
- **WHEN** a delta spec contains `## REMOVED Requirements` section
- **THEN** the system extracts requirement names to be removed

#### Scenario: Parse RENAMED requirements
- **WHEN** a delta spec contains `## RENAMED Requirements` section with FROM/TO pairs
- **THEN** the system extracts the old and new requirement names

#### Scenario: Require at least one delta operation
- **WHEN** a delta spec has no ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** the system returns an error indicating no delta operations were found

### Requirement: Delta Operation Application Order
The system SHALL apply delta operations in the order: RENAMED → REMOVED → MODIFIED → ADDED to ensure correct merging.

#### Scenario: Apply RENAMED first
- **WHEN** applying delta operations
- **THEN** RENAMED operations are applied first to update requirement names before other operations

#### Scenario: Apply REMOVED second
- **WHEN** applying delta operations
- **THEN** REMOVED operations are applied after RENAMED to delete requirements

#### Scenario: Apply MODIFIED third
- **WHEN** applying delta operations
- **THEN** MODIFIED operations are applied after REMOVED to update existing requirements

#### Scenario: Apply ADDED last
- **WHEN** applying delta operations
- **THEN** ADDED operations are applied last to insert new requirements

### Requirement: Requirement Name Normalization
The system SHALL normalize requirement names by trimming whitespace and using case-insensitive matching to prevent duplicates.

#### Scenario: Normalize whitespace
- **WHEN** comparing requirement names
- **THEN** the system trims leading and trailing whitespace

#### Scenario: Case-insensitive matching
- **WHEN** comparing requirement names
- **THEN** the system uses lowercase comparison to match requirements

#### Scenario: Prevent duplicates due to formatting
- **WHEN** requirements differ only in whitespace or case
- **THEN** the system treats them as the same requirement

### Requirement: Spec Merging Algorithm
The system SHALL merge delta operations into base specs while preserving requirement ordering.

#### Scenario: Preserve original requirement order
- **WHEN** merging modified or renamed requirements
- **THEN** the system maintains their original position in the requirements list

#### Scenario: Append new requirements at end
- **WHEN** adding new requirements
- **THEN** the system appends them after all existing requirements

#### Scenario: Rebuild requirements section
- **WHEN** merging is complete
- **THEN** the system reconstructs the spec with proper markdown structure and spacing

### Requirement: New Spec Creation
The system SHALL create new spec files when a delta spec has no corresponding main spec.

#### Scenario: Generate spec skeleton for new specs
- **WHEN** creating a new spec file
- **THEN** the system generates a skeleton with title, purpose placeholder, and requirements section

#### Scenario: Restrict operations for new specs
- **WHEN** creating a new spec
- **THEN** only ADDED operations are allowed (MODIFIED/REMOVED/RENAMED return an error)

#### Scenario: Error on invalid operations for new specs
- **WHEN** a delta spec for a new capability contains MODIFIED/REMOVED/RENAMED operations
- **THEN** the system returns an error indicating only ADDED is allowed for new specs

### Requirement: Post-Merge Validation
The system SHALL validate rebuilt specs before writing them to ensure correctness.

#### Scenario: Validate rebuilt spec structure
- **WHEN** a spec has been rebuilt from delta operations
- **THEN** the system validates the spec structure before writing

#### Scenario: Prevent writing invalid specs
- **WHEN** post-merge validation fails
- **THEN** the system aborts the archive operation without writing any specs

#### Scenario: Display validation errors
- **WHEN** post-merge validation fails
- **THEN** the system displays detailed validation errors with file paths and issue descriptions

### Requirement: Atomic Spec Updates
The system SHALL prepare all spec updates before writing to ensure all-or-nothing consistency.

#### Scenario: Prepare all updates first
- **WHEN** processing multiple delta specs
- **THEN** the system builds and validates all updated specs before writing any files

#### Scenario: Abort on first validation error
- **WHEN** any spec fails validation during preparation
- **THEN** the system aborts without writing any spec files

#### Scenario: Write all specs after validation
- **WHEN** all specs pass validation
- **THEN** the system writes all updated specs to their main spec locations

### Requirement: Spec Update Display
The system SHALL display spec update operations with clear formatting and counts.

#### Scenario: Display specs to be updated
- **WHEN** spec updates are prepared
- **THEN** the system lists each spec with "create" or "update" status

#### Scenario: Display operation counts per spec
- **WHEN** applying spec updates
- **THEN** the system displays counts: "+ N added, ~ N modified, - N removed, → N renamed" per spec

#### Scenario: Display total operation counts
- **WHEN** all spec updates are complete
- **THEN** the system displays totals across all specs

### Requirement: Spec Update Confirmation
The system SHALL prompt for confirmation before applying spec updates unless --yes flag is provided.

#### Scenario: Prompt for spec update confirmation
- **WHEN** spec updates are ready to apply
- **THEN** the system prompts "Proceed with spec updates?" and waits for user response

#### Scenario: Skip prompt with yes flag
- **WHEN** --yes flag is provided
- **THEN** the system proceeds with spec updates without prompting

#### Scenario: Skip updates when user declines
- **WHEN** user declines spec update prompt
- **THEN** the system skips spec updates but continues with archive operation

### Requirement: Archive Success Reporting
The system SHALL display clear success messages after archiving.

#### Scenario: Display archive success
- **WHEN** archive operation completes successfully
- **THEN** the system displays "Change '<change-id>' archived as '<archive-name>'"

#### Scenario: Display spec update success
- **WHEN** spec updates are applied
- **THEN** the system displays "Specs updated successfully" after showing operation counts

### Requirement: Auto-Accept on Archive
The system SHALL automatically convert `tasks.md` to `tasks.json` during archive if not already accepted, ensuring archived changes have stable task format.

#### Scenario: Archive triggers auto-accept
- **WHEN** archiving a change that has `tasks.md` but no `tasks.json`
- **THEN** the system displays a warning that auto-acceptance will occur
- **AND** the system converts `tasks.md` to `tasks.json` before archiving
- **AND** the system removes `tasks.md` after successful conversion

#### Scenario: Archive with existing tasks.json
- **WHEN** archiving a change that already has `tasks.json`
- **THEN** the system proceeds normally without conversion
- **AND** the archived change contains `tasks.json`

#### Scenario: Auto-accept failure blocks archive
- **WHEN** auto-acceptance fails during archive (e.g., invalid tasks.md format)
- **THEN** the system displays the conversion error
- **AND** the system aborts the archive operation
- **AND** no files are modified

### Requirement: Local Cleanup After PR Worktree Operations
The system SHALL clean up the local change directory after successfully completing a PR worktree operation (archive or remove), preventing orphan files from remaining in the user's working directory after the PR is merged.

#### Scenario: PR archive cleans local directory
- **WHEN** `spectr pr archive <change-id>` successfully creates a PR
- **THEN** the system SHALL remove the local change directory at `spectr/changes/<change-id>/`
- **AND** all files including untracked files (e.g., `tasks.json` from `spectr accept`) are removed

#### Scenario: PR remove cleans local directory
- **WHEN** `spectr pr rm <change-id>` successfully creates a PR
- **THEN** the system SHALL remove the local change directory at `spectr/changes/<change-id>/`
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
