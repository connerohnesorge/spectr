## ADDED Requirements

### Requirement: Provider Composition Architecture
The init system SHALL use composable provider interfaces to separate memory file configuration from slash command configuration.

#### Scenario: Tool provider composition
- **WHEN** a tool supports both memory files and slash commands
- **THEN** it SHALL implement ToolProvider interface using embedded provider fields
- **AND** SHALL compose MemoryFileProvider and SlashCommandProvider implementations
- **AND** SHALL allow nil providers for tools that only support one interface

#### Scenario: Memory file provider separation
- **WHEN** configuring memory files (CLAUDE.md, CLINE.md, etc.)
- **THEN** each tool SHALL have a dedicated MemoryFileProvider implementation
- **AND** each provider SHALL manage only its specific memory file
- **AND** providers SHALL use marker-based updates for idempotent configuration

#### Scenario: Slash command provider separation
- **WHEN** configuring slash commands (.claude/commands/, etc.)
- **THEN** each tool SHALL have a dedicated SlashCommandProvider implementation
- **AND** each provider SHALL manage only its specific slash command directory
- **AND** providers SHALL render templates for proposal, apply, and archive commands

### Requirement: SpectrAgentsUpdater Cross-File Updates
The init system SHALL provide a reusable provider that updates spectr/AGENTS.md with generic Spectr usage instructions for all memory file-based tools.

#### Scenario: Memory file tools update spectr/AGENTS.md
- **WHEN** configuring a memory file-based tool (Claude Code, Cline, Qoder, etc.)
- **THEN** the tool provider SHALL compose SpectrAgentsUpdater
- **AND** SpectrAgentsUpdater SHALL update spectr/AGENTS.md using marker-based updates
- **AND** SHALL use the existing spectr/AGENTS.md.tmpl template
- **AND** SHALL ensure all memory file tools reference Spectr instructions consistently

#### Scenario: SpectrAgentsUpdater reusability
- **WHEN** creating a composite tool provider for memory file tools
- **THEN** it SHALL include SpectrAgentsUpdater as an embedded field
- **AND** executor SHALL call SpectrAgentsUpdater.ConfigureMemoryFile()
- **AND** updates to spectr/AGENTS.md SHALL be idempotent across multiple tool configurations

### Requirement: Provider Interface Definitions
The init system SHALL define clear provider interfaces for tool integration patterns.

#### Scenario: MemoryFileProvider interface
- **WHEN** implementing memory file configuration
- **THEN** provider SHALL implement ConfigureMemoryFile(projectPath string) error
- **AND** SHALL implement IsMemoryFileConfigured(projectPath string) bool
- **AND** SHALL use marker-based file updates for idempotency

#### Scenario: SlashCommandProvider interface
- **WHEN** implementing slash command configuration
- **THEN** provider SHALL implement ConfigureSlashCommands(projectPath string) error
- **AND** SHALL implement AreSlashCommandsConfigured(projectPath string) bool
- **AND** SHALL configure proposal, apply, and archive command files

#### Scenario: ToolProvider interface composition
- **WHEN** implementing a complete tool provider
- **THEN** provider SHALL implement GetName() string
- **AND** SHALL implement GetMemoryFileProvider() MemoryFileProvider (may return nil)
- **AND** SHALL implement GetSlashCommandProvider() SlashCommandProvider (may return nil)
- **AND** SHALL use embedded fields for automatic interface implementation

### Requirement: Embedded Field Composition Pattern
The init system SHALL use Go embedded fields for provider composition to enable automatic interface implementation without wrapper methods.

#### Scenario: Composite provider structure
- **WHEN** defining a composite tool provider
- **THEN** it SHALL embed provider implementations as anonymous fields
- **AND** embedded fields SHALL be pointers to provider structs
- **AND** Go SHALL automatically delegate interface methods to embedded fields

#### Scenario: Example composite provider
- **WHEN** implementing ClaudeCodeToolProvider
- **THEN** it SHALL embed *ClaudeMemoryFileProvider for CLAUDE.md management
- **AND** SHALL embed *ClaudeSlashCommandProvider for .claude/commands/ management
- **AND** SHALL embed *SpectrAgentsUpdater for spectr/AGENTS.md updates
- **AND** executor SHALL call each provider's methods in sequence

## MODIFIED Requirements

### Requirement: Automatic Slash Command Installation
When a config-based AI tool is selected during initialization, the system SHALL automatically install the corresponding slash command files for that tool using the ToolProvider composition pattern instead of legacy Configurator implementations.

Config-based tools include those that create instruction files (e.g., `claude-code` creates `CLAUDE.md`). Slash command files are the workflow command files (e.g., `.claude/commands/spectr/proposal.md`).

The `ToolDefinition` model SHALL NOT include a `ConfigPath` field, as actual file paths are determined by individual providers. The registry maintains tool metadata (ID, Name, Type, Priority) but delegates file path resolution to provider implementations.

This automatic installation provides users with complete Spectr integration in a single selection, eliminating the need for redundant tool entries in the wizard.

#### Scenario: Claude Code auto-installs slash commands
- **WHEN** user selects `claude-code` in the init wizard
- **THEN** the system retrieves ClaudeCodeToolProvider via getToolProvider()
- **AND** calls ClaudeMemoryFileProvider.ConfigureMemoryFile() to create `CLAUDE.md`
- **AND** calls SpectrAgentsUpdater.ConfigureMemoryFile() to update `spectr/AGENTS.md`
- **AND** calls ClaudeSlashCommandProvider.ConfigureSlashCommands() to create slash commands
- **AND** the system creates `.claude/commands/spectr/proposal.md`
- **AND** the system creates `.claude/commands/spectr/apply.md`
- **AND** the system creates `.claude/commands/spectr/archive.md`
- **AND** all files are tracked in the execution result
- **AND** the completion screen shows all created/updated files

#### Scenario: Multiple tools with slash commands selected
- **WHEN** user selects both `claude-code` and `cursor` in the init wizard
- **THEN** the system retrieves ClaudeCodeToolProvider and creates CLAUDE.md + slash commands
- **AND** retrieves CursorToolProvider and creates .cursor/commands/spectr-proposal.md + slash commands
- **AND** all files from both tools are created and tracked separately
- **AND** the completion screen lists all created files grouped by tool

#### Scenario: Slash command files already exist
- **WHEN** user runs init and selects `claude-code`
- **AND** `.claude/commands/spectr/proposal.md` already exists
- **THEN** the ClaudeSlashCommandProvider updates content between markers
- **AND** the file's YAML frontmatter is preserved
- **AND** no error occurs
- **AND** the file is marked as "updated" rather than "created" in execution result

#### Scenario: Config-based tool without slash mapping
- **WHEN** a config-based tool has no slash command provider
- **THEN** GetSlashCommandProvider() returns nil
- **AND** only the memory file provider is invoked
- **AND** no error occurs
- **AND** the system continues with remaining tool configurations

#### Scenario: Tool provider retrieval
- **WHEN** executor needs to configure a tool
- **THEN** it calls getToolProvider(toolID) to retrieve ToolProvider
- **AND** checks if GetMemoryFileProvider() is non-nil before invoking
- **AND** checks if GetSlashCommandProvider() is non-nil before invoking
- **AND** aggregates files created/updated from all providers

#### Scenario: Sequential provider invocation
- **WHEN** configuring a composite tool provider
- **THEN** executor invokes MemoryFileProvider.ConfigureMemoryFile() first
- **AND** invokes SpectrAgentsUpdater.ConfigureMemoryFile() if present
- **AND** invokes SlashCommandProvider.ConfigureSlashCommands() last
- **AND** tracks created/updated files from each provider separately
- **AND** aggregates errors from all provider invocations
