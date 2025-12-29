# ADDED Requirements
## ADDED Requirements

### Requirement: Domain Package

The system SHALL define a `internal/domain` package containing shared domain types to break import cycles.

#### Scenario: TemplateRef in domain package

- **WHEN** code needs to reference a template
- **THEN** it SHALL use `domain.TemplateRef` from `internal/domain`
- **AND** `TemplateRef` SHALL have the following structure:

```go
type TemplateRef struct {
    Name     string             // template file name (e.g., "instruction-pointer.md.tmpl")
    Template *template.Template // pre-parsed template
}
```

- **AND** rendering SHALL be performed by TemplateManager, not by TemplateRef
- **AND** TemplateRef is a lightweight typed handle without rendering logic
- **AND** initializers SHALL call `tm.Render(templateRef.Name, ctx)` to render templates

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
- **AND** the struct SHALL have the following definition:

```go
// TemplateContext holds path-related template variables for dynamic directory names.
// Created via templateContextFromConfig(cfg) in the executor, not via defaults.
type TemplateContext struct {
    BaseDir     string // e.g., "spectr" (from cfg.SpectrDir)
    SpecsDir    string // e.g., "spectr/specs" (from cfg.SpecsDir())
    ChangesDir  string // e.g., "spectr/changes" (from cfg.ChangesDir())
    ProjectFile string // e.g., "spectr/project.md" (from cfg.ProjectFile())
    AgentsFile  string // e.g., "spectr/AGENTS.md" (from cfg.AgentsFile())
}
```

- **AND** TemplateContext instances SHALL be created via `templateContextFromConfig(cfg)`, not via a DefaultTemplateContext function
- **AND** all path values SHALL be derived from Config.SpectrDir, not hardcoded defaults

#### Scenario: TemplateContext derived from Config

- **WHEN** the executor needs to create a TemplateContext from Config
- **THEN** it SHALL use `templateContextFromConfig(cfg *Config)` to derive the values
- **AND** `BaseDir` SHALL equal `cfg.SpectrDir`
- **AND** `SpecsDir` SHALL equal `cfg.SpecsDir()`
- **AND** `ChangesDir` SHALL equal `cfg.ChangesDir()`
- **AND** `ProjectFile` SHALL equal `cfg.ProjectFile()`
- **AND** `AgentsFile` SHALL equal `cfg.AgentsFile()`

### Requirement: Provider Interface

The system SHALL define a `Provider` interface that returns a list of initializers.

```go
type Provider interface {
    // Initializers returns the list of initializers for this provider.
    // Receives TemplateManager to allow passing TemplateRef directly to initializers.
    Initializers(ctx context.Context, tm *TemplateManager) []Initializer
}
```

#### Scenario: Provider returns initializers

- **WHEN** a provider's `Initializers(ctx context.Context, tm *TemplateManager)` method is called
- **THEN** it SHALL receive a context.Context for cancellation and deadlines
- **AND** it SHALL receive a TemplateManager for resolving template references
- **AND** it SHALL return a slice of `Initializer` implementations
- **AND** the initializers MAY be empty if the provider requires no setup

### Requirement: Initializer Interface

The system SHALL define an `Initializer` interface with `Init` and `IsSetup` methods.

```go
type Initializer interface {
    // Init creates or updates files. Returns result with file changes and error if initialization fails.
    // Must be idempotent (safe to run multiple times).
    // Receives both filesystems - initializer decides which to use based on its type.
    Init(ctx context.Context, projectFs, homeFs afero.Fs, cfg *Config, tm *TemplateManager) (ExecutionResult, error)

    // IsSetup returns true if this initializer's artifacts already exist.
    // Receives both filesystems - initializer checks the appropriate one.
    // PURPOSE: Used by the setup wizard to show which providers are already configured.
    // NOT used to skip initializers during execution - Init() always runs (idempotent).
    IsSetup(projectFs, homeFs afero.Fs, cfg *Config) bool
}
```

#### Scenario: Initializer setup check

- **WHEN** `IsSetup(projectFs, homeFs, cfg)` is called on an initializer
- **THEN** it SHALL receive both project and home filesystems
- **AND** it SHALL return `true` if and only if ALL of the initializer's artifacts already exist
- **AND** it SHALL return `false` if ANY artifact is missing or setup is needed
- **AND** the initializer SHALL decide internally which filesystem to check based on its configuration

#### Scenario: Initializer execution

- **WHEN** `Init(ctx, projectFs, homeFs, cfg, tm)` is called on an initializer
- **THEN** it SHALL receive both project and home filesystems
- **AND** it SHALL decide internally which filesystem to use based on its configuration
- **AND** it SHALL create or update the necessary files in the appropriate filesystem
- **AND** it SHALL return an `ExecutionResult` containing created and updated file paths
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

### Requirement: Registration Struct

The system SHALL define a `Registration` struct containing provider metadata and implementation.

```go
type Registration struct {
    ID       string   // Unique identifier (kebab-case, e.g., "claude-code")
    Name     string   // Human-readable name (e.g., "Claude Code")
    Priority int      // Display order (lower = higher priority)
    Provider Provider // Implementation
}
```

#### Scenario: Registration fields

- **WHEN** a Registration is created
- **THEN** it SHALL have an `ID` field containing a unique kebab-case identifier
- **AND** it SHALL have a `Name` field containing a human-readable name
- **AND** it SHALL have a `Priority` field containing an integer for display ordering (lower values = higher priority)
- **AND** it SHALL have a `Provider` field containing the Provider implementation

### Requirement: Provider Registration (Explicit, No init())

The system SHALL support registering providers explicitly from a central location, not via init() functions.

#### Scenario: Register provider with metadata

- **WHEN** a provider is registered via `RegisterProvider(reg Registration) error`
- **THEN** the registration SHALL use the Registration struct with ID, Name, Priority, and Provider
- **AND** the system SHALL reject duplicate provider IDs with a clear error
- **AND** the function SHALL return an error (not panic) for invalid registrations

#### Scenario: RegisterAllProviders at startup

- **WHEN** the application starts
- **THEN** it SHALL call `RegisterAllProviders()` explicitly from `cmd/root.go` or `main()`
- **AND** the function SHALL register all built-in providers in one place
- **AND** the function SHALL return an error if any registration fails
- **AND** if any registration fails, successfully registered providers SHALL remain registered (no rollback)
- **AND** individual provider files SHALL NOT contain `init()` functions for registration

#### Scenario: Retrieve registered providers

- **WHEN** providers are queried via `RegisteredProviders() []Registration`
- **THEN** the system SHALL return all registered providers sorted by priority (lower first)
- **AND** the function SHALL be callable after `RegisterAllProviders()` completes

### Requirement: TemplateManager Interface

The system SHALL provide a `TemplateManager` for resolving and rendering templates.

```go
type TemplateManager struct {
    templates *template.Template
}

// Render renders a template by name with the given data.
func (tm *TemplateManager) Render(templateName string, data interface{}) (string, error)

// InstructionPointer returns the instruction-pointer.md.tmpl template reference.
func (tm *TemplateManager) InstructionPointer() domain.TemplateRef

// Agents returns the AGENTS.md.tmpl template reference.
func (tm *TemplateManager) Agents() domain.TemplateRef

// SlashCommand returns a Markdown template reference for the given slash command type.
// Used by SlashCommandsInitializer, HomeSlashCommandsInitializer, PrefixedSlashCommandsInitializer, and HomePrefixedSlashCommandsInitializer.
func (tm *TemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef

// TOMLSlashCommand returns a TOML template reference for the given slash command type.
// Used by TOMLSlashCommandsInitializer (Gemini only).
func (tm *TemplateManager) TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef
```

#### Scenario: TemplateManager rendering

- **WHEN** a template is rendered via `TemplateManager.Render(templateName, data)`
- **THEN** the template SHALL be looked up by name
- **AND** the template SHALL be executed with the provided data
- **AND** the rendered string SHALL be returned
- **AND** an error SHALL be returned if the template is not found or rendering fails

#### Scenario: TemplateManager SlashCommand accessor

- **WHEN** `TemplateManager.SlashCommand(cmd)` is called with a SlashCommand
- **THEN** it SHALL return a `domain.TemplateRef` for the corresponding Markdown template
- **AND** `SlashProposal` SHALL map to `slash-proposal.md.tmpl`
- **AND** `SlashApply` SHALL map to `slash-apply.md.tmpl`

#### Scenario: TemplateManager TOMLSlashCommand accessor

- **WHEN** `TemplateManager.TOMLSlashCommand(cmd)` is called with a SlashCommand
- **THEN** it SHALL return a `domain.TemplateRef` for the corresponding TOML template
- **AND** `SlashProposal` SHALL map to `slash-proposal.toml.tmpl`
- **AND** `SlashApply` SHALL map to `slash-apply.toml.tmpl`

#### Scenario: Template resolution

- **WHEN** templates are resolved
- **THEN** templates from `internal/initialize/templates` SHALL be available (AGENTS.md.tmpl, instruction-pointer.md.tmpl)
- **AND** templates from `internal/domain` SHALL be available (slash-*.md.tmpl, slash-*.toml.tmpl)
- **AND** if duplicate template names exist, the last-wins precedence SHALL apply

### Requirement: Filesystem Abstraction

The system SHALL use `afero.Fs` rooted at project directory for all file operations.

#### Scenario: Project-relative paths

- **WHEN** an initializer accesses files
- **THEN** all paths SHALL be relative to the project root
- **AND** the filesystem SHALL be created via `afero.NewBasePathFs(osFs, projectPath)`

#### Scenario: Home filesystem root

- **WHEN** the home filesystem is created
- **THEN** it SHALL be rooted at the user's home directory
- **AND** the home directory SHALL be obtained via `os.UserHomeDir()`
- **AND** if `os.UserHomeDir()` returns an error, initialization SHALL fail entirely
- **AND** the filesystem SHALL be created as `afero.NewBasePathFs(afero.NewOsFs(), homeDir)` to create an afero.Fs instance

### Requirement: ConfigFile Initializer

The system SHALL provide a built-in `ConfigFileInitializer` for marker-based file updates.

```go
// ConfigFileInitializer creates or updates a config file with marker-based content.
type ConfigFileInitializer struct {
    path     string           // target file path (e.g., "CLAUDE.md", "AGENTS.md")
    template domain.TemplateRef // template to render for content between markers
}

// NewConfigFileInitializer creates a ConfigFileInitializer for the given path and template.
func NewConfigFileInitializer(path string, template domain.TemplateRef) *ConfigFileInitializer
```

#### Scenario: ConfigFileInitializer construction

- **WHEN** a ConfigFileInitializer is created via `NewConfigFileInitializer(path, template)`
- **THEN** it SHALL receive a file path (e.g., "CLAUDE.md", "AGENTS.md")
- **AND** it SHALL receive a TemplateRef directly (not a function)
- **AND** the TemplateRef SHALL be resolved at provider construction time when Initializers() is called
- **AND** the initializer SHALL use `projectFs` for all file operations

#### Scenario: Create new config file

- **WHEN** the config file does not exist
- **THEN** the initializer SHALL create it with the instruction content between markers

#### Scenario: Update existing config file

- **WHEN** the config file exists with markers
- **THEN** the initializer SHALL replace content between markers
- **AND** it SHALL preserve content outside markers

#### Scenario: Config file markers

- **WHEN** content is written to a config file
- **THEN** it SHALL be wrapped with `<!-- spectr:start -->` and `<!-- spectr:end -->` markers (lowercase)
- **NOTE**: All markdown markers use lowercase `start`/`end` for consistency when writing

#### Scenario: Case-insensitive marker matching

- **WHEN** searching for markers in existing files
- **THEN** the search SHALL be case-insensitive for backward compatibility
- **AND** both `<!-- spectr:START -->` (uppercase) and `<!-- spectr:start -->` (lowercase) SHALL be recognized
- **AND** when writing new markers, the system SHALL always use lowercase
- **AND** this ensures behavioral equivalence with files created by older versions

#### Scenario: Orphaned start marker handling

- **WHEN** a config file contains a start marker but the end marker is missing immediately after
- **THEN** the initializer SHALL search for an end marker anywhere after the start position using `strings.Index` on a slice starting from the position after the start marker
- **AND** if an end marker is found after the start marker, the initializer SHALL use it to perform the update
- **AND** if no end marker exists anywhere after the start, the initializer SHALL replace content from the start marker onward with the new block (start + content + end)
- **AND** the initializer SHALL NOT append a duplicate block that leaves orphaned markers

#### Scenario: Missing end marker recovery

- **WHEN** start marker exists at position X but no end marker exists anywhere after position X
- **THEN** the initializer SHALL trim content from position X onward
- **AND** insert the complete new block (startMarker + newContent + endMarker)
- **AND** this prevents duplicate marker blocks and orphaned start markers

#### Scenario: Missing markers in existing file

- **WHEN** ConfigFileInitializer finds an existing file but markers are missing
- **THEN** the initializer SHALL insert start and end markers at the end of the file
- **AND** insert the content between the newly created markers
- **AND** preserve all existing file content

#### Scenario: Orphaned end marker (end without start)

- **WHEN** a config file contains an end marker without a preceding start marker
- **THEN** the initializer SHALL return an error indicating corrupted marker structure
- **AND** the error message SHALL indicate the orphaned end marker position

#### Scenario: Nested markers (start before previous end)

- **WHEN** a config file contains a start marker before the previous start marker's end marker
- **THEN** the initializer SHALL return an error indicating nested markers are not supported
- **AND** the error message SHALL indicate both marker positions

#### Scenario: Multiple start markers without end

- **WHEN** a config file contains multiple start markers without end markers between them
- **THEN** the initializer SHALL return an error indicating multiple unpaired start markers
- **AND** the error message SHALL indicate the positions of the duplicate start markers

### Requirement: SlashCommands Initializer

The system SHALL provide built-in slash command initializers with separate types for filesystem and format.

#### Scenario: Create project Markdown slash commands

- **WHEN** `SlashCommandsInitializer` runs
- **THEN** it SHALL create `proposal.md` and `apply.md` command files in the project filesystem
- **AND** it SHALL use `slash-proposal.md.tmpl` and `slash-apply.md.tmpl` templates

#### Scenario: Create home Markdown slash commands

- **WHEN** `HomeSlashCommandsInitializer` runs
- **THEN** it SHALL create `proposal.md` and `apply.md` command files in the home filesystem (user home)
- **AND** it SHALL use `slash-proposal.md.tmpl` and `slash-apply.md.tmpl` templates

#### Scenario: Create prefixed Markdown slash commands

- **WHEN** `PrefixedSlashCommandsInitializer` runs with prefix `spectr-`
- **THEN** it SHALL create `spectr-proposal.md` and `spectr-apply.md` command files in the project filesystem
- **AND** it SHALL use `slash-proposal.md.tmpl` and `slash-apply.md.tmpl` templates
- **NOTE**: Used by Antigravity for non-standard path patterns in project directory

#### Scenario: Create home prefixed Markdown slash commands

- **WHEN** `HomePrefixedSlashCommandsInitializer` runs with prefix `spectr-`
- **THEN** it SHALL create `spectr-proposal.md` and `spectr-apply.md` command files in the home filesystem
- **AND** it SHALL use `slash-proposal.md.tmpl` and `slash-apply.md.tmpl` templates
- **NOTE**: Used by Codex for home directory paths with prefixed filenames (e.g., `~/.codex/prompts/spectr-proposal.md`)

#### Scenario: Create TOML slash commands

- **WHEN** `TOMLSlashCommandsInitializer` runs
- **THEN** it SHALL create `proposal.toml` and `apply.toml` command files in the project filesystem
- **AND** it SHALL use `slash-proposal.toml.tmpl` and `slash-apply.toml.tmpl` templates
- **AND** the templates SHALL produce TOML files with `description` and `prompt` fields
- **NOTE**: Only Gemini uses this initializer type

### Requirement: Directory Initializer

The system SHALL provide built-in directory initializers with separate types for project vs home filesystem.

#### Scenario: Create project directories

- **WHEN** `DirectoryInitializer` runs
- **THEN** it SHALL create all specified directories in the project filesystem if they do not exist
- **AND** it SHALL recursively create parent directories as needed (like `os.MkdirAll`)
- **AND** it SHALL succeed silently if the directory already exists

#### Scenario: Create home directories

- **WHEN** `HomeDirectoryInitializer` runs
- **THEN** it SHALL create all specified directories in the home filesystem (user home) if they do not exist
- **AND** it SHALL recursively create parent directories as needed (like `os.MkdirAll`)
- **AND** it SHALL succeed silently if the directory already exists

### Requirement: Initializer Deduplication

The system SHALL deduplicate initializers by type and path when multiple providers are configured.

```go
// deduplicatable is an optional interface for initializers that support deduplication.
// Initializers that implement this interface can be deduplicated based on their key.
// Note: lowercase name indicates this is a private/internal interface.
type deduplicatable interface {
    // dedupeKey returns a unique key for deduplication.
    // The key represents the "scope" or "territory" that the initializer operates on.
    // Format varies by initializer category:
    //
    // Directory-based (creates/ensures directory exists):
    //   Format: "<TypeName>:<directory>"
    //   Examples:
    //     - "DirectoryInitializer:.claude/commands/spectr"
    //     - "HomeDirectoryInitializer:.codex/prompts"
    //
    // File-based (creates/updates specific file):
    //   Format: "<TypeName>:<file_path>"
    //   Examples:
    //     - "ConfigFileInitializer:CLAUDE.md"
    //     - "ConfigFileInitializer:AGENTS.md"
    //
    // Directory+files (creates multiple files in directory):
    //   Format: "<TypeName>:<directory>"
    //   Examples:
    //     - "SlashCommandsInitializer:.claude/commands/spectr"
    //     - "HomeSlashCommandsInitializer:.codex/prompts"
    //     - "TOMLSlashCommandsInitializer:.gemini/commands/spectr"
    //
    // Prefixed (creates multiple prefixed files in directory):
    //   Format: "<TypeName>:<directory>:<prefix>"
    //   Examples:
    //     - "PrefixedSlashCommandsInitializer:.agent/workflows:spectr-"
    //     - "HomePrefixedSlashCommandsInitializer:.codex/prompts:spectr-"
    //
    // Note: Paths are normalized with filepath.Clean before key generation.
    dedupeKey() string
}
```

#### Scenario: Optional deduplicatable interface

- **WHEN** initializers are collected for execution
- **THEN** the system SHALL check if each initializer implements the optional `deduplicatable` interface
- **AND** initializers implementing `deduplicatable` SHALL provide a `dedupeKey() string` method
- **AND** initializers NOT implementing `deduplicatable` SHALL always run

#### Scenario: Deduplication timing

- **WHEN** initializers are prepared for execution
- **THEN** the system SHALL first deduplicate initializers (remove duplicates by key)
- **AND** then sort the deduplicated list by type priority
- **AND** then execute in the resulting order

#### Scenario: Shared initializer deduplication

- **WHEN** multiple providers return initializers with the same dedup key
- **THEN** the system SHALL run the initializer only once
- **AND** the dedup key SHALL include the type name (e.g., "DirectoryInitializer:.claude/commands/spectr")
- **AND** separate types (`DirectoryInitializer` vs `HomeDirectoryInitializer`) SHALL have different keys
- **AND** paths SHALL be normalized (filepath.Clean) before generating keys
- **AND** the dedup key format SHALL be: `<TypeName>:<path>` where TypeName is the concrete type name

#### Scenario: Different configurations run separately

- **WHEN** providers return initializers with different paths or different types
- **THEN** all initializers SHALL run

### Requirement: Initializer Ordering

The system SHALL execute initializers in a guaranteed order by type.

#### Scenario: Directory initializers run first

- **WHEN** initializers are collected for execution
- **THEN** `DirectoryInitializer` and `HomeDirectoryInitializer` SHALL run before `ConfigFileInitializer`
- **AND** `ConfigFileInitializer` SHALL run before `SlashCommandsInitializer`, `HomeSlashCommandsInitializer`, `PrefixedSlashCommandsInitializer`, `HomePrefixedSlashCommandsInitializer`, and `TOMLSlashCommandsInitializer`

#### Scenario: Ordering within same category

- **WHEN** multiple initializers of the same type exist (e.g., multiple SlashCommandsInitializer instances)
- **THEN** the order of execution within that category is unspecified
- **AND** implementations MAY use any stable ordering (e.g., registration order, alphabetical)
- **AND** implementations SHALL NOT rely on a specific order within the same type category

#### Scenario: Ordering is guaranteed

- **WHEN** documentation describes initializer ordering
- **THEN** it SHALL be a documented API guarantee
- **AND** implementers MAY rely on this ordering

### Requirement: ExecutionResult Type

The system SHALL define an `ExecutionResult` type for initialization results.

```go
type ExecutionResult struct {
    CreatedFiles []string // files created
    UpdatedFiles []string // files updated
}
// Note: Error is returned separately, not stored in this struct
```

#### Scenario: ExecutionResult from initializer

- **WHEN** `Init()` is called on an initializer
- **THEN** it SHALL return an `ExecutionResult` containing created and updated files
- **AND** errors SHALL be returned separately (second return value)

#### Scenario: ExecutionResult from executor

- **WHEN** multiple initializers are executed
- **THEN** the executor SHALL combine all results into a single `ExecutionResult`
- **AND** it SHALL concatenate all `CreatedFiles` slices
- **AND** it SHALL concatenate all `UpdatedFiles` slices
- **AND** on error, it SHALL return partial results from initializers that succeeded

### Requirement: Dual Filesystem Support

The system SHALL provide two filesystem instances to all initializers.

#### Scenario: Filesystem provision

- **WHEN** an initializer's `Init()` or `IsSetup()` method is called
- **THEN** it SHALL receive both `projectFs` (rooted at project directory) and `homeFs` (rooted at user's home directory)
- **AND** the initializer SHALL decide internally which filesystem to use based on its type

#### Scenario: Initializer configuration

- **WHEN** an initializer is constructed
- **THEN** it MAY be configured to use either the project or home filesystem
- **AND** this configuration is determined by the initializer type (Home* types use homeFs)

#### Scenario: Filesystem selection by type

- **WHEN** an initializer determines which filesystem to use
- **THEN** the choice SHALL be based on the initializer's type
- **AND** `HomeDirectoryInitializer`, `HomeSlashCommandsInitializer`, and `HomePrefixedSlashCommandsInitializer` SHALL use `homeFs`
- **AND** `DirectoryInitializer`, `SlashCommandsInitializer`, `PrefixedSlashCommandsInitializer`, `ConfigFileInitializer`, and `TOMLSlashCommandsInitializer` SHALL use `projectFs`
- **AND** initializers receive both filesystems but use only the appropriate one based on their type

### Requirement: Fail-Fast Error Handling

The system SHALL stop on the first initialization error.

#### Scenario: Initializer failure

- **WHEN** an initializer fails during execution
- **THEN** the system SHALL stop immediately (fail-fast)
- **AND** the system SHALL return partial results (files created before failure) in ExecutionResult
- **AND** the system SHALL return the error separately (not stored in ExecutionResult)
- **AND** the system SHALL NOT rollback successful initializers - files created before the error SHALL remain on disk
- **AND** the user SHALL be able to fix the issue and re-run `spectr init`

#### Scenario: Partial results persistence

- **WHEN** execution fails partway through (e.g., 2 files created successfully, 3rd file fails)
- **THEN** all files created before the error SHALL remain on disk
- **AND** the returned ExecutionResult.CreatedFiles SHALL list those files
- **AND** the returned error SHALL describe the failure
- **AND** no rollback or cleanup SHALL occur automatically
- **AND** the user can inspect the partial state to diagnose the issue

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
