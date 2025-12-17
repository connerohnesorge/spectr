# CLI Interface Delta Spec

## MODIFIED Requirements

### Requirement: Validate Command Flags
The validate command SHALL support flags for controlling validation behavior and output format. Validation always treats warnings as errors.

#### Scenario: Default validation behavior (always strict)
- WHEN user runs `spectr validate <item>` without any strict flag
- THEN validation SHALL treat warnings as errors
- AND exit code SHALL be 1 if warnings or errors exist
- AND validation report SHALL show valid=false for any issues

#### Scenario: JSON output flag
- WHEN user provides `--json` flag
- THEN output SHALL be formatted as JSON
- AND SHALL include items, summary, and version fields
- AND SHALL be parseable by standard JSON tools

#### Scenario: Type disambiguation flag
- WHEN user provides `--type change` or `--type spec`
- THEN the command SHALL treat the item as the specified type
- AND SHALL skip type auto-detection
- AND SHALL error if item does not exist as that type

#### Scenario: All items flag
- WHEN user provides `--all` flag
- THEN the command SHALL validate all changes and all specs
- AND SHALL run in bulk validation mode

#### Scenario: Changes only flag
- WHEN user provides `--changes` flag
- THEN the command SHALL validate all changes only
- AND SHALL skip specs

#### Scenario: Specs only flag
- WHEN user provides `--specs` flag
- THEN the command SHALL validate all specs only
- AND SHALL skip changes

#### Scenario: Non-interactive flag
- WHEN user provides `--no-interactive` flag
- THEN the command SHALL not prompt for input
- AND SHALL print usage hint if no item specified
- AND SHALL exit with code 1

