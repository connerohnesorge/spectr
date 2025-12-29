# Cli Interface Specification Delta

## ADDED Requirements

### Requirement: Validation Output Format

The validate command SHALL display validation issues in a consistent, detailed format for both single-item and bulk validation modes.

#### Scenario: Single item validation with issues

- **WHEN** user runs `spectr validate <item>` and validation finds issues
- **THEN** output SHALL display "✗ <item> has N issue(s):"
- **AND** each issue SHALL be displayed on a separate line with format "  [LEVEL] PATH: MESSAGE"
- **AND** the command SHALL exit with code 1

#### Scenario: Bulk validation with issues

- **WHEN** user runs `spectr validate --all` and validation finds issues in multiple items
- **THEN** output SHALL display "✗ <item> (<type>): N issue(s)" for each failed item
- **AND** immediately following each failed item, all issue details SHALL be displayed
- **AND** each issue SHALL use the format "  [LEVEL] PATH: MESSAGE"
- **AND** a summary line SHALL display "N passed, M failed, T total"
- **AND** the command SHALL exit with code 1

#### Scenario: Bulk validation all passing

- **WHEN** user runs `spectr validate --all` and all items are valid
- **THEN** output SHALL display "✓ <item> (<type>)" for each item
- **AND** a summary line SHALL display "N passed, 0 failed, N total"
- **AND** the command SHALL exit with code 0

#### Scenario: JSON output format

- **WHEN** user provides `--json` flag with any validation command
- **THEN** output SHALL be valid JSON
- **AND** SHALL include full issue details with level, path, and message fields
- **AND** SHALL include per-item results and summary statistics
