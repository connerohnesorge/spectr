## MODIFIED Requirements

### Requirement: Provider Interface

The init system SHALL define a `Provider` interface that all AI CLI tool integrations implement, with one provider per tool handling both instruction files and slash commands.

#### Scenario: Provider interface methods

- **WHEN** a new provider is created
- **THEN** it SHALL implement `ID() string` returning a unique kebab-case identifier
- **AND** it SHALL implement `Name() string` returning the human-readable name
- **AND** it SHALL implement `Priority() int` returning display order
- **AND** it SHALL implement `ConfigFile() string` returning instruction file path or empty string
- **AND** it SHALL implement `GetProposalCommandPath() string` returning relative path for proposal command or empty string
- **AND** it SHALL implement `GetArchiveCommandPath() string` returning relative path for archive command or empty string
- **AND** it SHALL implement `GetApplyCommandPath() string` returning relative path for apply command or empty string
- **AND** it SHALL implement `CommandFormat() CommandFormat` returning Markdown or TOML
- **AND** it SHALL implement `Configure(projectPath, spectrDir string) error` for configuration
- **AND** it SHALL implement `IsConfigured(projectPath string) bool` for status checks

#### Scenario: Single provider per tool

- **WHEN** a tool has both an instruction file and slash commands
- **THEN** one provider SHALL handle both (e.g., ClaudeProvider handles CLAUDE.md and .claude/commands/)
- **AND** there SHALL NOT be separate config and slash providers for the same tool

#### Scenario: Flexible command paths

- **WHEN** a provider returns paths from command path methods
- **THEN** each method SHALL return a relative path including directory and filename
- **AND** paths MAY have different directories for each command type
- **AND** paths MAY have different file extensions based on CommandFormat
- **AND** empty string indicates the provider does not support that command

#### Scenario: HasSlashCommands detection

- **WHEN** code calls `HasSlashCommands()` on a provider
- **THEN** it SHALL return true if ANY command path method returns a non-empty string
- **AND** it SHALL return false only if ALL command path methods return empty strings
