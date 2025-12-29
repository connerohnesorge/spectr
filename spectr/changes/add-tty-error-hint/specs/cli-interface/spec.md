# ADDED Requirements
## ADDED Requirements

### Requirement: TTY Error Hint

The system SHALL provide a helpful hint when the interactive wizard fails due to TTY unavailability.

#### Scenario: TTY unavailable error shows hint

- WHEN a user runs `spectr init` in an environment without a TTY (CI, Docker, piped input)
- AND the Bubbletea TUI fails with a TTY-related error
- THEN the error message SHALL include the original error
- AND the error message SHALL suggest using `--non-interactive` flag
- AND the error message SHALL provide an example command: `spectr init --non-interactive --tools <tool1,tool2>`

#### Scenario: Non-TTY errors remain unchanged

- WHEN the interactive wizard fails for reasons unrelated to TTY access
- THEN the original error message SHALL be displayed without modification
