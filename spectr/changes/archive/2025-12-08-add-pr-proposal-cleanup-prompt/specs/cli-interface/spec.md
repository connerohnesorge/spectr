## ADDED Requirements

### Requirement: PR Proposal Yes Flag

The `spectr pr proposal` command SHALL support a `--yes` or `-y` flag that enables non-interactive mode for automated usage, skipping the cleanup confirmation prompt after successful PR creation.

#### Scenario: Non-interactive proposal with yes flag

- **WHEN** user runs `spectr pr proposal <change-id> --yes`
- **AND** PR creation succeeds
- **THEN** the cleanup confirmation prompt is NOT displayed
- **AND** the local change directory is kept (safe default behavior)
- **AND** the command exits successfully without user interaction

#### Scenario: Yes flag short form

- **WHEN** user runs `spectr pr proposal <change-id> -y`
- **THEN** the behavior is identical to `--yes`
- **AND** the PR workflow proceeds non-interactively

#### Scenario: Yes flag in CI/CD pipelines

- **WHEN** the command is run in a non-TTY environment (e.g., CI pipeline)
- **AND** `--yes` flag is provided
- **THEN** the command completes successfully without prompting
- **AND** no TTY-dependent operations are attempted
