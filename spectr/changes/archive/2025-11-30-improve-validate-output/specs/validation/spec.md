# Delta Specification

## ADDED Requirements

### Requirement: Bulk Validation Human Output Formatting

The validation system SHALL produce bulk validation human-readable output with
improved spacing, relative paths, file grouping, and color-coded error levels
for easier scanning.

#### Scenario: Visual separation between failed items

- **WHEN** bulk validation encounters multiple failed items in human output mode
- **THEN** output SHALL include a blank line between each failed item's error
  listing
- **AND** passed items SHALL be listed without blank lines between them
- **AND** failed items SHALL be visually distinct from passed items

#### Scenario: Relative path display

- **WHEN** validation issues include file paths in human output mode
- **THEN** paths SHALL be displayed relative to the spectr/ directory
- **AND** paths SHALL NOT include the project root or spectr/ prefix
- **AND** example: `changes/foo/specs/bar/spec.md` instead of
  `/home/user/project/spectr/changes/foo/specs/bar/spec.md`

#### Scenario: Grouping issues by file

- **WHEN** multiple issues exist in the same file in human output mode
- **THEN** the file path SHALL be displayed once as a header
- **AND** issues for that file SHALL be indented below the header
- **AND** this grouping SHALL reduce visual clutter from repeated paths

#### Scenario: Color-coded error levels

- **WHEN** bulk validation output is displayed in a TTY
- **THEN** [ERROR] labels SHALL be styled in red using lipgloss
- **AND** [WARNING] labels SHALL be styled in yellow using lipgloss
- **AND** [INFO] labels SHALL remain unstyled
- **AND** when not in a TTY, labels SHALL display without color codes

#### Scenario: Enhanced summary line

- **WHEN** bulk validation completes with failures
- **THEN** summary SHALL show "X passed, Y failed (E errors, W warnings), Z
  total"
- **AND** summary SHALL only show error/warning breakdown if failures exist
- **AND** example: "22 passed, 2 failed (5 errors, 1 warning), 24 total"

#### Scenario: Item type indicators

- **WHEN** bulk validation results are displayed in human output mode
- **THEN** each item SHALL show a type indicator alongside its name
- **AND** changes SHALL display "(change)" indicator
- **AND** specs SHALL display "(spec)" indicator
