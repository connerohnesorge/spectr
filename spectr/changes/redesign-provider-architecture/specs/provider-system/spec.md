## ADDED Requirements

### Requirement: Provider Interface
The system SHALL define a `Provider` interface that returns a list of initializers.

#### Scenario: Provider returns initializers
- **WHEN** a provider is queried for its initializers
- **THEN** it SHALL return a slice of `Initializer` implementations
- **AND** the initializers MAY be empty if the provider requires no setup

### Requirement: Initializer Interface
The system SHALL define an `Initializer` interface with `Init`, `IsSetup`, `Path`, and `IsGlobal` methods.

#### Scenario: Initializer setup check
- **WHEN** `IsSetup(fs, cfg)` is called on an initializer
- **THEN** it SHALL return `true` if the initializer's artifacts already exist
- **AND** it SHALL return `false` if setup is needed

#### Scenario: Initializer execution
- **WHEN** `Init(ctx, fs, cfg, tm)` is called on an initializer
- **THEN** it SHALL create or update the necessary files in the filesystem
- **AND** it SHALL return an error if initialization fails
- **AND** it SHALL be idempotent (safe to run multiple times)

#### Scenario: Initializer path identification
- **WHEN** `Path()` is called on an initializer
- **THEN** it SHALL return the file or directory path this initializer manages
- **AND** the path SHALL be used for deduplication

#### Scenario: Initializer global flag
- **WHEN** `IsGlobal()` is called on an initializer
- **THEN** it SHALL return `true` if the initializer uses the global filesystem
- **AND** it SHALL return `false` if the initializer uses the project filesystem

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

### Requirement: Provider Registration
The system SHALL support registering providers with metadata at registration time.

#### Scenario: Register provider with metadata
- **WHEN** a provider is registered
- **THEN** the registration SHALL include ID, Name, Priority, and Provider implementation
- **AND** the system SHALL reject duplicate provider IDs

#### Scenario: Retrieve registered providers
- **WHEN** providers are queried
- **THEN** the system SHALL return providers sorted by priority (lower first)

### Requirement: Filesystem Abstraction
The system SHALL use `afero.Fs` rooted at project directory for all file operations.

#### Scenario: Project-relative paths
- **WHEN** an initializer accesses files
- **THEN** all paths SHALL be relative to the project root
- **AND** the filesystem SHALL be created via `afero.NewBasePathFs(osFs, projectPath)`

### Requirement: ConfigFile Initializer
The system SHALL provide a built-in `ConfigFileInitializer` for marker-based file updates.

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

### Requirement: SlashCommands Initializer
The system SHALL provide a built-in `SlashCommandsInitializer` for creating slash commands.

#### Scenario: Create proposal command
- **WHEN** the initializer runs
- **THEN** it SHALL create a `proposal` command file in the specified directory
- **AND** the file SHALL use the specified format (Markdown or TOML)

#### Scenario: Create apply command
- **WHEN** the initializer runs
- **THEN** it SHALL create an `apply` command file in the specified directory
- **AND** the file SHALL use the specified format (Markdown or TOML)

### Requirement: Directory Initializer
The system SHALL provide a built-in `DirectoryInitializer` for creating directories.

#### Scenario: Create directories
- **WHEN** the initializer runs
- **THEN** it SHALL create all specified directories if they do not exist
- **AND** it SHALL create parent directories as needed

### Requirement: Initializer Deduplication
The system SHALL deduplicate initializers by file path when multiple providers are configured.

#### Scenario: Shared initializer deduplication
- **WHEN** multiple providers return initializers with the same path
- **THEN** the system SHALL run the initializer only once
- **AND** deduplication SHALL be based on the `Path()` return value

#### Scenario: Different paths run separately
- **WHEN** providers return initializers with different paths
- **THEN** all initializers SHALL run

### Requirement: Initializer Ordering
The system SHALL execute initializers in a guaranteed order by type.

#### Scenario: Directory initializers run first
- **WHEN** initializers are collected for execution
- **THEN** `DirectoryInitializer` SHALL run before `ConfigFileInitializer`
- **AND** `ConfigFileInitializer` SHALL run before `SlashCommandsInitializer`

#### Scenario: Ordering is guaranteed
- **WHEN** documentation describes initializer ordering
- **THEN** it SHALL be a documented API guarantee
- **AND** implementers MAY rely on this ordering

### Requirement: Git Repository Requirement
The system SHALL require a git repository for initialization.

#### Scenario: Git repo check at start
- **WHEN** `spectr init` is executed
- **THEN** the system SHALL check for a git repository before proceeding
- **AND** if not a git repo, it SHALL fail with a clear error message
- **AND** the error SHALL instruct the user to run `git init` first

### Requirement: Dual Filesystem Support
The system SHALL provide two filesystem instances for project and global paths.

#### Scenario: Project filesystem
- **WHEN** an initializer has `IsGlobal() == false`
- **THEN** it SHALL receive the project filesystem rooted at the project directory

#### Scenario: Global filesystem
- **WHEN** an initializer has `IsGlobal() == true`
- **THEN** it SHALL receive the global filesystem rooted at the user's home directory

### Requirement: Partial Failure Handling
The system SHALL handle partial initialization failures gracefully.

#### Scenario: Initializer failure
- **WHEN** an initializer fails during execution
- **THEN** the system SHALL report which initializers failed
- **AND** the system SHALL NOT rollback successful initializers
- **AND** the user SHALL be able to re-run `spectr init` to retry

