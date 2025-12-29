# Validation Delta Spec

## MODIFIED Requirements

### Requirement: Validation Report Structure

The validation system SHALL produce structured validation reports containing issue details and summary statistics, always treating warnings as errors.

#### Scenario: Report always strict

- WHEN validation encounters WARNING level issues
- THEN the report SHALL treat warnings as errors
- AND valid SHALL be false if errors OR warnings exist
- AND exit code SHALL be non-zero for warnings
- AND there is no opt-in strict mode flag

#### Scenario: Report with errors and warnings

- WHEN validation encounters both ERROR and WARNING level issues
- THEN the report SHALL list all issues with level, path, and message
- AND the summary SHALL count errors (including promoted warnings), warnings (always 0), and info separately
- AND valid SHALL be false if any errors exist

#### Scenario: JSON output format

- WHEN validation is invoked with --json flag
- THEN the output SHALL be valid JSON
- AND SHALL include items array with per-item results
- AND SHALL include summary with totals and byType breakdowns
- AND SHALL include version field for format compatibility
