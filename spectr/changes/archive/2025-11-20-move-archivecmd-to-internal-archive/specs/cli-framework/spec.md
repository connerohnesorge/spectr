## MODIFIED Requirements

### Requirement: Archive Command
The CLI SHALL provide an `archive` command that moves completed changes to a dated archive directory and applies delta specifications to main specs.

#### Scenario: Archive with change ID
- **WHEN** user runs `spectr archive <change-id>`
- **THEN** the system archives the specified change without prompting

#### Scenario: Interactive archive selection
- **WHEN** user runs `spectr archive` without specifying a change ID
- **THEN** the system displays a list of active changes and prompts for selection

#### Scenario: Non-interactive archiving with yes flag
- **WHEN** user runs `spectr archive <change-id> --yes`
- **THEN** the system archives without any confirmation prompts

#### Scenario: Skip spec updates for tooling changes
- **WHEN** user runs `spectr archive <change-id> --skip-specs`
- **THEN** the system archives the change without updating main specs

#### Scenario: Skip validation with confirmation
- **WHEN** user runs `spectr archive <change-id> --no-validate`
- **THEN** the system warns about skipping validation and requires confirmation unless --yes flag is also provided
