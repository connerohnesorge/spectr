# Delta Specification

## MODIFIED Requirements

### Requirement: Initialization Next Steps Message

The `spectr init` command SHALL display a formatted "Next steps" message after
successful initialization that provides users with clear, actionable guidance
for getting started with Spectr.

The message SHALL include:

1. Three progressive steps with copy-paste ready prompts for AI assistants
2. Visual separators to make the message stand out
3. References to key Spectr files and documentation
4. Placeholder text that users can customize (e.g., "[YOUR FEATURE HERE]")

The init command SHALL NOT automatically create project files outside the
`spectr/` directory (such as README.md). Users maintain full control over their
project's root-level documentation.

#### Scenario: Interactive mode initialization succeeds

- **WHEN** a user completes initialization via the interactive TUI wizard
- **THEN** the completion screen SHALL display the next steps message
- **AND** the message SHALL appear after the list of created/updated files
- **AND** the message SHALL be visually distinct with a separator line
- **AND** the message SHALL provide three numbered steps with specific prompts

#### Scenario: Non-interactive mode initialization succeeds

- **WHEN** a user runs `spectr init --non-interactive` and initialization
  succeeds
- **THEN** the command output SHALL display the next steps message
- **AND** the message SHALL appear after the list of created/updated files
- **AND** the message SHALL be formatted consistently with the interactive mode
- **AND** the message SHALL include the same three progressive steps

#### Scenario: Initialization fails with errors

- **WHEN** initialization fails with errors
- **THEN** the next steps message SHALL NOT be displayed
- **AND** only error messages SHALL be shown

#### Scenario: Next steps message content

- **WHEN** the next steps message is displayed
- **THEN** step 1 SHALL guide users to populate spectr/project.md
- **AND** step 2 SHALL guide users to create their first change proposal
- **AND** step 3 SHALL guide users to learn the Spectr workflow from
  spectr/AGENTS.md
- **AND** each step SHALL include a complete, copy-paste ready prompt in quotes
- **AND** the message SHALL include a visual separator using dashes or similar
  characters

#### Scenario: Init does not create README

- **WHEN** a user runs `spectr init` on a project without a README.md
- **THEN** the init command SHALL NOT create a README.md file
- **AND** only files within the `spectr/` directory SHALL be created
- **AND** tool-specific files (e.g., CLAUDE.md, .cursor/) SHALL be created as
  configured
