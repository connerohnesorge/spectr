## ADDED Requirements

### Requirement: FileInitializer Interface
The init system SHALL define a `FileInitializer` interface that represents an atomic unit responsible for creating or updating a single file type.

#### Scenario: FileInitializer interface methods
- **WHEN** a new file initializer is created
- **THEN** it SHALL implement `ID() string` returning a unique identifier in format `{type}:{path}`
- **AND** it SHALL implement `FilePath() string` returning the relative file path it manages
- **AND** it SHALL implement `Configure(projectPath string, tm TemplateRenderer) error` for file operations
- **AND** it SHALL implement `IsConfigured(projectPath string) bool` for status checks

#### Scenario: Initializer path expansion
- **WHEN** an initializer's `FilePath()` returns a path containing `~`
- **THEN** the `Configure()` method SHALL expand `~` to the user's home directory internally
- **AND** external callers SHALL NOT need to pre-expand paths

### Requirement: InstructionFileInitializer
The init system SHALL provide an `InstructionFileInitializer` that creates and updates instruction markdown files with marker-based content management.

#### Scenario: Create instruction file
- **WHEN** `Configure()` is called and the instruction file does not exist
- **THEN** the initializer SHALL create the file at the specified path
- **AND** the file SHALL contain content between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Update existing instruction file
- **WHEN** `Configure()` is called and the instruction file exists
- **THEN** the initializer SHALL update content between existing markers
- **AND** user content outside markers SHALL be preserved

#### Scenario: Initializer ID format
- **WHEN** `ID()` is called on an InstructionFileInitializer
- **THEN** it SHALL return `instruction:{path}` (e.g., `instruction:CLAUDE.md`)

### Requirement: MarkdownSlashCommandInitializer
The init system SHALL provide a `MarkdownSlashCommandInitializer` that creates and updates markdown-format slash command files with YAML frontmatter.

#### Scenario: Create markdown slash command
- **WHEN** `Configure()` is called and the command file does not exist
- **THEN** the initializer SHALL create the directory structure if needed
- **AND** the initializer SHALL create the file with YAML frontmatter
- **AND** the file SHALL contain command content between spectr markers

#### Scenario: Generic command name support
- **WHEN** a MarkdownSlashCommandInitializer is created
- **THEN** it SHALL accept a command name parameter (e.g., "proposal", "apply", "sync")
- **AND** the initializer SHALL render the appropriate template for that command name

#### Scenario: Initializer ID format
- **WHEN** `ID()` is called on a MarkdownSlashCommandInitializer
- **THEN** it SHALL return `markdown-cmd:{path}` (e.g., `markdown-cmd:.claude/commands/spectr/proposal.md`)

### Requirement: TOMLSlashCommandInitializer
The init system SHALL provide a `TOMLSlashCommandInitializer` that creates TOML-format slash command files for providers like Gemini CLI.

#### Scenario: Create TOML slash command
- **WHEN** `Configure()` is called for a TOML command
- **THEN** the initializer SHALL create a `.toml` file
- **AND** the TOML SHALL include `description` field with command description
- **AND** the TOML SHALL include `prompt` field with command content

#### Scenario: Initializer ID format
- **WHEN** `ID()` is called on a TOMLSlashCommandInitializer
- **THEN** it SHALL return `toml-cmd:{path}` (e.g., `toml-cmd:.gemini/commands/spectr/proposal.toml`)

### Requirement: Initializer Helper Functions
The init system SHALL provide helper functions for common operations on initializer lists, enabling providers to implement interface methods with minimal boilerplate.

#### Scenario: Configure all initializers
- **WHEN** `ConfigureInitializers(inits []FileInitializer, projectPath string, tm TemplateRenderer)` is called
- **THEN** it SHALL call `Configure()` on each initializer in order
- **AND** it SHALL stop on the first error and return that error (fail-fast)

#### Scenario: Check all initializers configured
- **WHEN** `AreInitializersConfigured(inits []FileInitializer, projectPath string)` is called
- **THEN** it SHALL return true only if ALL initializers return true for `IsConfigured()`
- **AND** it SHALL return false if ANY initializer returns false

#### Scenario: Get all initializer paths
- **WHEN** `GetInitializerPaths(inits []FileInitializer)` is called
- **THEN** it SHALL return a slice containing each initializer's `FilePath()`
- **AND** duplicate paths SHALL be deduplicated

## MODIFIED Requirements

### Requirement: Provider Interface
The init system SHALL define a minimal `Provider` interface with 6 methods that all AI CLI tool integrations implement, composing functionality from FileInitializers.

#### Scenario: Provider interface methods
- **WHEN** a new provider is created
- **THEN** it SHALL implement `ID() string` returning a unique kebab-case identifier
- **AND** it SHALL implement `Name() string` returning the human-readable name
- **AND** it SHALL implement `Priority() int` returning display order (lower = higher priority)
- **AND** it SHALL implement `Initializers() []FileInitializer` returning its file initializers
- **AND** it SHALL implement `IsConfigured(projectPath string) bool` for status checks
- **AND** it SHALL implement `GetFilePaths() []string` returning all managed file paths

#### Scenario: Provider configuration via helper
- **WHEN** the system needs to configure a provider
- **THEN** it SHALL call `ConfigureInitializers(provider.Initializers(), projectPath, tm)`
- **AND** the Provider interface SHALL NOT include a `Configure()` method

#### Scenario: Provider implements interface directly
- **WHEN** a new provider is created
- **THEN** it SHALL implement the Provider interface directly without embedding BaseProvider
- **AND** it SHALL use helper functions for `IsConfigured()` and `GetFilePaths()` implementations

### Requirement: Per-Provider File Organization
The init system SHALL organize provider implementations as separate Go files under `internal/initialize/providers/`, with one file per provider implementing the interface directly.

#### Scenario: Provider file structure
- **WHEN** a provider file is created
- **THEN** it SHALL be named `{provider-id}.go` (e.g., `claude.go`, `gemini.go`)
- **AND** it SHALL contain an `init()` function that registers its provider
- **AND** it SHALL implement all 6 Provider interface methods
- **AND** it SHALL NOT embed BaseProvider (BaseProvider is removed)

#### Scenario: Adding a new provider
- **WHEN** a developer adds a new AI CLI provider
- **THEN** they SHALL create a single file under `internal/initialize/providers/`
- **AND** the file SHALL implement the 6-method Provider interface
- **AND** the file SHALL compose FileInitializers for its file types
- **AND** the file SHALL use helper functions for IsConfigured and GetFilePaths

## REMOVED Requirements

### Requirement: Command Format Support
The init system previously defined a `CommandFormat` type with `FormatMarkdown` and `FormatTOML` values, and providers returned this via `CommandFormat()` method.

**Reason**: Command format is now implicit in the initializer type used. `MarkdownSlashCommandInitializer` creates markdown files; `TOMLSlashCommandInitializer` creates TOML files. No explicit format enum needed.

**Migration**: Remove `CommandFormat` type and `CommandFormat()` method from Provider interface. Providers compose the appropriate initializer type instead.
