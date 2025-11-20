## MODIFIED Requirements

### Requirement: Automatic Slash Command Installation

When a config-based AI tool is selected during initialization, the system SHALL automatically install the corresponding slash command files for that tool without requiring separate user selection.

Config-based tools include those that create instruction files (e.g., `claude-code` creates `CLAUDE.md`). Slash command files are the workflow command files (e.g., `.claude/commands/spectr/proposal.md`).

The `ToolDefinition` model SHALL NOT include a `ConfigPath` field, as actual file paths are determined by individual configurators. The registry maintains tool metadata (ID, Name, Type, Priority) but delegates file path resolution to configurator implementations. Tool IDs SHALL use a type-safe constant approach to prevent typos and enable compile-time validation.

This automatic installation provides users with complete Spectr integration in a single selection, eliminating the need for redundant tool entries in the wizard.

#### Scenario: Claude Code auto-installs slash commands

- **WHEN** user selects `claude-code` in the init wizard
- **THEN** the system creates `CLAUDE.md` in the project root
- **AND** the system creates `.claude/commands/spectr/proposal.md`
- **AND** the system creates `.claude/commands/spectr/apply.md`
- **AND** the system creates `.claude/commands/spectr/archive.md`
- **AND** all files are tracked in the execution result
- **AND** the completion screen shows all 4 files created

#### Scenario: Multiple tools with slash commands selected

- **WHEN** user selects both `claude-code` and `cursor` in the init wizard
- **THEN** the system creates `CLAUDE.md` and both config + slash commands for Claude
- **AND** the system creates `.cursor/commands/spectr-proposal.md` and slash commands for Cursor
- **AND** all files from both tools are created and tracked separately
- **AND** the completion screen lists all created files grouped by tool

#### Scenario: Slash command files already exist

- **WHEN** user run init and selects `claude-code`
- **AND** `.claude/commands/spectr/proposal.md` already exists
- **THEN** the existing file's content between `<!-- spectr:START -->` and `<!-- spectr:END -->` is updated
- **AND** the file's YAML frontmatter is preserved
- **AND** no error occurs
- **AND** the file is marked as "updated" rather than "created" in execution result

#### Scenario: Config-based tool without slash mapping

- **WHEN** a config-based tool has no slash command equivalent in the mapping
- **THEN** only the config file is created
- **AND** no error occurs
- **AND** the system continues with remaining tool configurations

#### Scenario: Tool mapping is explicit and centralized

- **WHEN** a developer reviews the mapping logic
- **THEN** they find the tool mapping integrated into the tool registry configuration
- **AND** the registry uses data-driven tool definitions with type-safe IDs
- **AND** the mapping can be extended for new tools through configuration

#### Scenario: ToolDefinition structure simplified

- **WHEN** a developer reviews the ToolDefinition struct in `internal/init/models.go`
- **THEN** the struct contains: ID (type-safe ToolID), Name, Type, Priority, and Configured fields
- **AND** the struct does NOT contain a ConfigPath field
- **AND** file paths are determined by configurator implementations, not the registry
- **AND** the configurator implementations use tool-specific configuration data

#### Scenario: Type-safe tool ID usage

- **WHEN** a developer references a tool ID in code
- **THEN** they use a defined ToolID constant (e.g., `ToolClaudeCode`)
- **AND** string literals for tool IDs trigger compiler warnings or errors
- **AND** IDE autocomplete suggests available tool ID constants
- **AND** typos in tool IDs are caught at compile time
