# Archive Workflow Delta Specification

## MODIFIED Requirements

### Requirement: Auto-Accept on Archive
The system SHALL automatically convert `tasks.md` to `tasks.jsonc` during archive if not already accepted, ensuring archived changes have stable task format.

#### Scenario: Archive triggers auto-accept
- **WHEN** archiving a change that has `tasks.md` but no `tasks.jsonc`
- **THEN** the system displays a warning that auto-acceptance will occur
- **AND** the system converts `tasks.md` to `tasks.jsonc` before archiving
- **AND** the system preserves `tasks.md` alongside `tasks.jsonc` in the archive

#### Scenario: Archive with existing tasks.jsonc
- **WHEN** archiving a change that already has `tasks.jsonc`
- **THEN** the system proceeds normally without conversion
- **AND** the archived change contains both `tasks.jsonc` and `tasks.md` (if present)

#### Scenario: Auto-accept failure blocks archive
- **WHEN** auto-acceptance fails during archive (e.g., invalid tasks.md format)
- **THEN** the system displays the conversion error
- **AND** the system aborts the archive operation
- **AND** no files are modified

#### Scenario: Archive preserves both task files
- **WHEN** archiving a change that has both tasks.md and tasks.jsonc
- **THEN** both files SHALL be moved to the archive directory
- **AND** the archive SHALL contain the human-readable tasks.md for historical reference
- **AND** the archive SHALL contain the machine-readable tasks.jsonc for tooling
