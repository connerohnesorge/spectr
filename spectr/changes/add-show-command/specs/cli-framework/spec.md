# Cli Framework Specification Delta

## ADDED Requirements

### Requirement: Show Command Structure
The CLI SHALL provide a `show` command that renders specifications or changes as HTML in a browser with support for images and mathematical notation.

#### Scenario: Show command registration
- **WHEN** the CLI is initialized
- **THEN** it SHALL include a ShowCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr show`
- **AND** help text SHALL describe browser rendering functionality

#### Scenario: Show spec invocation
- **WHEN** user runs `spectr show <spec-id> --type spec`
- **THEN** the system renders the spec's Markdown content to HTML
- **AND** starts an HTTP server on an auto-assigned port
- **AND** opens the rendered page in the user's default browser
- **AND** keeps the server running until Ctrl+C is pressed

#### Scenario: Show change invocation
- **WHEN** user runs `spectr show <change-id> --type change`
- **THEN** the system renders all change files (proposal.md, tasks.md, design.md, deltas) to HTML
- **AND** displays content in tabbed sections (Proposal, Tasks, Design, Deltas)
- **AND** starts an HTTP server and opens browser as with specs

#### Scenario: Auto-detect item type
- **WHEN** user runs `spectr show <item-id>` without `--type` flag
- **THEN** the system checks if item exists as a spec in `spectr/specs/<item-id>/`
- **AND** checks if item exists as a change in `spectr/changes/<item-id>/`
- **AND** renders whichever is found
- **AND** displays error if neither exists or if ambiguous (both exist)

#### Scenario: Server URL printed to terminal
- **WHEN** the HTTP server starts successfully
- **THEN** the terminal displays message "Opening browser at http://localhost:<port>"
- **AND** displays message "Press Ctrl+C to stop server"
- **AND** the server URL is shown for manual browser access

#### Scenario: Server start failure
- **WHEN** the HTTP server fails to start (e.g., port binding error)
- **THEN** the system retries with a new port up to 3 times
- **AND** displays error message if all retries fail
- **AND** exits with non-zero status code

### Requirement: Show Command Flags
The show command SHALL support flags for controlling rendering behavior.

#### Scenario: Type disambiguation flag
- **WHEN** user provides `--type spec` or `--type change`
- **THEN** the command SHALL treat the item as the specified type
- **AND** SHALL skip auto-detection
- **AND** SHALL error if item does not exist as that type

#### Scenario: Port flag for custom port
- **WHEN** user provides `--port <number>` flag
- **THEN** the server SHALL attempt to bind to the specified port
- **AND** SHALL fail if port is already in use (no auto-retry for explicit port)
- **AND** SHALL validate port is in valid range (1-65535)

### Requirement: Show Command Help Text
The show command SHALL provide comprehensive help documentation.

#### Scenario: Command help display
- **WHEN** user invokes `spectr show --help`
- **THEN** help text SHALL describe browser rendering purpose
- **AND** SHALL list available flags (--type, --port) with descriptions
- **AND** SHALL show usage examples for specs and changes
- **AND** SHALL indicate required positional argument (item ID)

### Requirement: Positional Argument for Item ID
The show command SHALL accept a required positional argument for the item to render.

#### Scenario: Required item ID argument
- **WHEN** show command is defined
- **THEN** it SHALL have an ItemID field tagged with `arg:"" required:"true"`
- **AND** omitting the argument SHALL display usage error
- **AND** the argument SHALL be a string representing spec ID or change ID
