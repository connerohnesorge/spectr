## ADDED Requirements

### Requirement: Chained Proposal Support

The system SHALL support declaration of dependencies between proposals using
YAML frontmatter in proposal.md files.

#### Scenario: Valid frontmatter with dependencies

- **WHEN** a proposal.md file contains YAML frontmatter with requires field
- **THEN** the system parses the dependencies
- **AND** validates that required proposals exist or are archived

#### Scenario: Accept with unmet dependencies

- **WHEN** a user runs `spectr accept` on a proposal with unmet dependencies
- **THEN** the command fails with a clear error message
- **AND** lists the unmet dependencies
