## ADDED Requirements

### Requirement: Domain Package
The system SHALL define a `internal/domain` package containing shared domain types to break import cycles.

#### Scenario: TemplateRef in domain package
- **WHEN** code needs to reference a template
- **THEN** it SHALL use `domain.TemplateRef` from `internal/domain`
- **AND** `TemplateRef` SHALL have `Name` and `Template` fields
- **AND** `TemplateRef` SHALL have a `Render(ctx TemplateContext) (string, error)` method

#### Scenario: SlashCommand in domain package
- **WHEN** code needs to reference a slash command type
- **THEN** it SHALL use `domain.SlashCommand` from `internal/domain`
- **AND** `SlashCommand` SHALL be a typed constant (`SlashProposal`, `SlashApply`)
- **AND** `SlashCommand` SHALL have a `String()` method for debugging
- **AND** Markdown templates SHALL be accessed via `TemplateManager.SlashCommand(cmd)`
- **AND** TOML templates SHALL be accessed via `TemplateManager.TOMLSlashCommand(cmd)`

#### Scenario: TemplateContext in domain package
- **WHEN** code needs template context with path variables
- **THEN** it SHALL use `domain.TemplateContext` from `internal/domain`
- **AND** `domain.DefaultTemplateContext()` SHALL return default path values

### Requirement: Provider Interface
The system SHALL define a `Provider` interface that returns a list of initializers.

#### Scenario: Provider returns initializers
- **WHEN** a provider's `Initializers(ctx, tm *TemplateManager)` method is called
- **THEN** it SHALL receive a TemplateManager for resolving template references
- **AND** it SHALL return a slice of `Initializer` implementations
- **AND** the initializers MAY be empty if the provider requires no setup

### Requirement: Initializer Interface
The system SHALL define an `Initializer` interface with `Init` and `IsSetup` methods.

#### Scenario: Initializer setup check
- **WHEN** `IsSetup(projectFs, globalFs, cfg)` is called on an initializer
- **THEN** it SHALL receive both project and global filesystems
- **AND** it SHALL return `true` if the initializer's artifacts already exist
- **AND** it SHALL return `false` if setup is needed
- **AND** the initializer SHALL decide internally which filesystem to check based on its configuration

#### Scenario: Initializer execution
- **WHEN** `Init(ctx, projectFs, globalFs, cfg, tm)` is called on an initializer
- **THEN** it SHALL receive both project and global filesystems
- **AND** it SHALL decide internally which filesystem to use based on its configuration
- **AND** it SHALL create or update the necessary files in the appropriate filesystem
- **AND** it SHALL return an `InitResult` containing created and updated file paths
- **AND** it SHALL return an error if initialization fails
- **AND** it SHALL be idempotent (safe to run multiple times)

### Requirement: Config Struct
The system SHALL provide a `Config` struct containing initialization configuration.

#### Scenario: Config with SpectrDir
- **WHEN** a Config is created
- **THEN** it SHALL have a `SpectrDir` field specifying the spectr directory path
- **AND** the path SHALL be relative to the filesystem root

#### Scenario: Config derived paths
- **WHEN** derived path methods are called on Config
- **THEN** `SpecsDir()` SHALL return `SpectrDir + "/specs"`
- **AND** `ChangesDir()` SHALL return `SpectrDir + "/changes"`
- **AND** `ProjectFile()` SHALL return `SpectrDir + "/project.md"`
- **AND** `AgentsFile()` SHALL return `SpectrDir + "/AGENTS.md"`

### Requirement: Provider Registration (Explicit, No init())
The system SHALL support registering providers explicitly from a central location, not via init() functions.

#### Scenario: Register provider with metadata
- **WHEN** a provider is registered via `RegisterProvider(reg Registration) error`
- **THEN** the registration SHALL include ID, Name, Priority, and Provider implementation
- **AND** the system SHALL reject duplicate provider IDs with a clear error
- **AND** the function SHALL return an error (not panic) for invalid registrations

#### Scenario: RegisterAllProviders at startup
- **WHEN** the application starts
- **THEN** it SHALL call `RegisterAllProviders()` explicitly from `cmd/root.go` or `main()`
- **AND** the function SHALL register all built-in providers in one place
- **AND** the function SHALL return an error if any registration fails
- **AND** individual provider files SHALL NOT contain `init()` functions for registration

#### Scenario: Retrieve registered providers
- **WHEN** providers are queried via `RegisteredProviders() []Registration`
- **THEN** the system SHALL return all registered providers sorted by priority (lower first)
- **AND** the function SHALL be callable after `RegisterAllProviders()` completes

### Requirement: Filesystem Abstraction
The system SHALL use `afero.Fs` rooted at project directory for all file operations.

#### Scenario: Project-relative paths
- **WHEN** an initializer accesses files
- **THEN** all paths SHALL be relative to the project root
- **AND** the filesystem SHALL be created via `afero.NewBasePathFs(osFs, projectPath)`

### Requirement: ConfigFile Initializer
The system SHALL provide a built-in `ConfigFileInitializer` for marker-based file updates.

#### Scenario: ConfigFileInitializer construction
- **WHEN** a ConfigFileInitializer is created
- **THEN** it SHALL receive a TemplateRef directly (not a function)
- **AND** the TemplateRef SHALL be resolved at provider construction time when Initializers() is called

#### Scenario: Create new config file
- **WHEN** the config file does not exist
- **THEN** the initializer SHALL create it with the instruction content between markers

#### Scenario: Update existing config file
- **WHEN** the config file exists with markers
- **THEN** the initializer SHALL replace content between markers
- **AND** it SHALL preserve content outside markers

#### Scenario: Config file markers
- **WHEN** content is written to a config file
- **THEN** it SHALL be wrapped with `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Orphaned start marker handling
- **WHEN** a config file contains a start marker but the end marker is missing immediately after
- **THEN** the initializer SHALL search for an end marker anywhere after the start position using `strings.LastIndex`
- **AND** if an end marker is found after the start marker, the initializer SHALL use it to perform the update
- **AND** if no end marker exists anywhere after the start, the initializer SHALL replace content from the start marker onward with the new block (start + content + end)
- **AND** the initializer SHALL NOT append a duplicate block that leaves orphaned markers

#### Scenario: Missing end marker recovery
- **WHEN** start marker exists at position X but no end marker exists anywhere after position X
- **THEN** the initializer SHALL trim content from position X onward
- **AND** insert the complete new block (startMarker + newContent + endMarker)
- **AND** this prevents duplicate marker blocks and orphaned start markers

### Requirement: SlashCommands Initializer
The system SHALL provide built-in slash command initializers with separate types for filesystem and format.

#### Scenario: Create project Markdown slash commands
- **WHEN** `SlashCommandsInitializer` runs
- **THEN** it SHALL create `proposal.md` and `apply.md` command files in the project filesystem
- **AND** it SHALL use `slash-proposal.md.tmpl` and `slash-apply.md.tmpl` templates

#### Scenario: Create global Markdown slash commands
- **WHEN** `GlobalSlashCommandsInitializer` runs
- **THEN** it SHALL create `proposal.md` and `apply.md` command files in the global filesystem (user home)
- **AND** it SHALL use `slash-proposal.md.tmpl` and `slash-apply.md.tmpl` templates

#### Scenario: Create TOML slash commands
- **WHEN** `TOMLSlashCommandsInitializer` runs
- **THEN** it SHALL create `proposal.toml` and `apply.toml` command files in the project filesystem
- **AND** it SHALL use `slash-proposal.toml.tmpl` and `slash-apply.toml.tmpl` templates
- **AND** the templates SHALL produce TOML files with `description` and `prompt` fields
- **NOTE**: Only Gemini uses this initializer type

### Requirement: Directory Initializer
The system SHALL provide built-in directory initializers with separate types for local vs global filesystem.

#### Scenario: Create project directories
- **WHEN** `DirectoryInitializer` runs
- **THEN** it SHALL create all specified directories in the project filesystem if they do not exist
- **AND** it SHALL create parent directories as needed

#### Scenario: Create global directories
- **WHEN** `GlobalDirectoryInitializer` runs
- **THEN** it SHALL create all specified directories in the global filesystem (user home) if they do not exist
- **AND** it SHALL create parent directories as needed

### Requirement: Initializer Deduplication
The system SHALL deduplicate initializers by type and path when multiple providers are configured.

#### Scenario: Optional deduplicatable interface
- **WHEN** initializers are collected for execution
- **THEN** the system SHALL check if each initializer implements the optional `deduplicatable` interface
- **AND** initializers implementing `deduplicatable` SHALL provide a `dedupeKey() string` method
- **AND** initializers NOT implementing `deduplicatable` SHALL always run

#### Scenario: Shared initializer deduplication
- **WHEN** multiple providers return initializers with the same dedup key
- **THEN** the system SHALL run the initializer only once
- **AND** the dedup key SHALL include the type name (e.g., "DirectoryInitializer:.claude/commands/spectr")
- **AND** separate types (`DirectoryInitializer` vs `GlobalDirectoryInitializer`) SHALL have different keys

#### Scenario: Different configurations run separately
- **WHEN** providers return initializers with different paths or different types
- **THEN** all initializers SHALL run

### Requirement: Initializer Ordering
The system SHALL execute initializers in a guaranteed order by type.

#### Scenario: Directory initializers run first
- **WHEN** initializers are collected for execution
- **THEN** `DirectoryInitializer` and `GlobalDirectoryInitializer` SHALL run before `ConfigFileInitializer`
- **AND** `ConfigFileInitializer` SHALL run before `SlashCommandsInitializer`, `GlobalSlashCommandsInitializer`, and `TOMLSlashCommandsInitializer`

#### Scenario: Ordering is guaranteed
- **WHEN** documentation describes initializer ordering
- **THEN** it SHALL be a documented API guarantee
- **AND** implementers MAY rely on this ordering

### Requirement: Initialize Result
The system SHALL return file change information from each initializer.

#### Scenario: Initializer returns result
- **WHEN** `Init()` is called on an initializer
- **THEN** it SHALL return an `InitResult` containing created and updated files
- **AND** the `InitResult` SHALL have `CreatedFiles` and `UpdatedFiles` fields

#### Scenario: Result accumulation
- **WHEN** multiple initializers are executed
- **THEN** the executor SHALL accumulate all `InitResult` values
- **AND** the accumulated results SHALL be returned in the `ExecutionResult`

### Requirement: ExecutionResult Type
The system SHALL define an `ExecutionResult` type for aggregated initialization results.

#### Scenario: ExecutionResult structure
- **WHEN** initialization completes
- **THEN** the system SHALL return an `ExecutionResult` containing:
  - `CreatedFiles []string` - all files created across all initializers
  - `UpdatedFiles []string` - all files updated across all initializers
  - `Errors []error` - any errors encountered during initialization

#### Scenario: aggregateResults function (success case)
- **WHEN** all initializers have completed successfully
- **THEN** the `aggregateResults(results []InitResult, errors []error) ExecutionResult` function SHALL combine all results
- **AND** it SHALL concatenate all created files into a single slice
- **AND** it SHALL concatenate all updated files into a single slice
- **AND** the errors parameter SHALL be nil in the success case (due to fail-fast behavior, errors never accumulate)

### Requirement: Dual Filesystem Support
The system SHALL provide two filesystem instances to all initializers.

#### Scenario: Filesystem provision
- **WHEN** an initializer's `Init()` or `IsSetup()` method is called
- **THEN** it SHALL receive both `projectFs` (rooted at project directory) and `globalFs` (rooted at user's home directory)
- **AND** the initializer SHALL decide internally which filesystem to use based on its configuration

#### Scenario: Initializer configuration
- **WHEN** an initializer is constructed
- **THEN** it MAY be configured to use either the project or global filesystem
- **AND** this configuration is internal to the initializer (not exposed via interface methods)

### Requirement: Fail-Fast Error Handling
The system SHALL stop on the first initialization error.

#### Scenario: Initializer failure
- **WHEN** an initializer fails during execution
- **THEN** the system SHALL stop immediately (fail-fast)
- **AND** the system SHALL return partial results (files created before failure)
- **AND** the system SHALL return the error in ExecutionResult.Errors
- **AND** the system SHALL NOT rollback successful initializers
- **AND** the user SHALL be able to fix the issue and re-run `spectr init`

### Requirement: Zero Technical Debt - No Compatibility Shims
The system SHALL NOT provide compatibility shims for deprecated registration patterns.

#### Scenario: Old Register() function removed
- **WHEN** the old provider registration system is removed
- **THEN** the old `Register(p Provider)` function SHALL be completely deleted
- **AND** NO deprecated `Register(_ any)` function SHALL exist that silently swallows calls
- **AND** code attempting to call the old `Register()` SHALL fail to compile with a clear error

#### Scenario: Explicit migration required
- **WHEN** providers are migrated to the new system
- **THEN** all provider `init()` functions SHALL be deleted
- **AND** registration SHALL happen exclusively via `RegisterAllProviders()`
- **AND** compiler errors SHALL enforce complete migration

