## ADDED Requirements

### Requirement: Provider Interface
The init system SHALL define a `Provider` interface that all AI CLI tool integrations implement, with one provider per tool handling both instruction files and slash commands.

#### Scenario: Provider interface methods
- **WHEN** a new provider is created
- **THEN** it SHALL implement `ID() string` returning a unique kebab-case identifier
- **AND** it SHALL implement `Name() string` returning the human-readable name
- **AND** it SHALL implement `Priority() int` returning display order
- **AND** it SHALL implement `ConfigFile() string` returning instruction file path or empty string
- **AND** it SHALL implement `SlashDir() string` returning slash commands directory or empty string
- **AND** it SHALL implement `CommandFormat() CommandFormat` returning Markdown or TOML
- **AND** it SHALL implement `Configure(projectPath, spectrDir string) error` for configuration
- **AND** it SHALL implement `IsConfigured(projectPath string) bool` for status checks

#### Scenario: Single provider per tool
- **WHEN** a tool has both an instruction file and slash commands
- **THEN** one provider SHALL handle both (e.g., ClaudeProvider handles CLAUDE.md and .claude/commands/)
- **AND** there SHALL NOT be separate config and slash providers for the same tool

### Requirement: Provider Registry
The init system SHALL provide a `Registry` that manages registration and lookup of providers using a registry pattern similar to `database/sql`.

#### Scenario: Register provider
- **WHEN** a provider calls `Register(provider Provider)`
- **THEN** the registry SHALL store the provider by its ID
- **AND** duplicate registration SHALL panic with a descriptive message

#### Scenario: Get provider by ID
- **WHEN** code calls `Get(id string) (Provider, bool)`
- **THEN** the registry SHALL return the provider and true if found
- **AND** SHALL return nil and false if not found

#### Scenario: List all providers
- **WHEN** code calls `All() []Provider`
- **THEN** the registry SHALL return all registered providers
- **AND** providers SHALL be sorted by Priority ascending

### Requirement: Per-Provider File Organization
The init system SHALL organize provider implementations as separate Go files under `internal/init/providers/`, with one file per provider.

#### Scenario: Provider file structure
- **WHEN** a provider file is created
- **THEN** it SHALL be named `{provider-id}.go` (e.g., `claude.go`, `gemini.go`)
- **AND** it SHALL contain an `init()` function that registers its provider
- **AND** it SHALL be self-contained with all provider-specific configuration

#### Scenario: Adding a new provider
- **WHEN** a developer adds a new AI CLI provider
- **THEN** they SHALL create a single file under `internal/init/providers/`
- **AND** the file SHALL implement the `Provider` interface
- **AND** the file SHALL call `Register()` in its `init()` function
- **AND** no other files SHALL require modification

### Requirement: Init Function Registration
The init system SHALL use Go's `init()` function pattern for automatic provider registration at startup.

#### Scenario: Auto-registration at startup
- **WHEN** the program starts
- **THEN** all provider `init()` functions SHALL execute before `main()`
- **AND** all providers SHALL be registered in the global registry
- **AND** registration order SHALL not affect functionality

### Requirement: Command Format Support
The init system SHALL support multiple command file formats through the `CommandFormat` type.

#### Scenario: Markdown command format
- **WHEN** a provider returns `FormatMarkdown` from `CommandFormat()`
- **THEN** slash command files SHALL be generated as `.md` files
- **AND** files SHALL use frontmatter and spectr markers

#### Scenario: TOML command format
- **WHEN** a provider returns `FormatTOML` from `CommandFormat()`
- **THEN** slash command files SHALL be generated as `.toml` files
- **AND** the TOML SHALL include `description` field with command description
- **AND** the TOML SHALL include `prompt` field with the command prompt content

## MODIFIED Requirements

### Requirement: Struct-Based Command Definition
The CLI framework SHALL use Go struct types with struct tags to declaratively define command structure, subcommands, flags, and arguments. Provider configuration SHALL be retrieved from the `Registry` interface rather than static global maps.

#### Scenario: Root command definition
- **WHEN** the CLI is initialized
- **THEN** it SHALL use a root struct with subcommand fields tagged with `cmd` for command definitions
- **AND** each subcommand SHALL be a nested struct type with appropriate tags

#### Scenario: Subcommand registration
- **WHEN** a new subcommand is added to the CLI
- **THEN** it SHALL be defined as a struct field on the parent command struct
- **AND** it SHALL use `cmd` tag to indicate it is a subcommand
- **AND** it SHALL include a `help` tag describing the command purpose

#### Scenario: Tool configuration lookup
- **WHEN** the executor needs tool configuration
- **THEN** it SHALL query the `Registry` via `Get(id)` method
- **AND** it SHALL NOT use hardcoded global maps
