## ADDED Requirements

### Requirement: Project Configuration File

The system SHALL support an optional `spectr.yaml` configuration file at the project root that allows customization of Spectr settings, with sensible defaults when no configuration file exists.

#### Scenario: Default behavior without config file
- **WHEN** user runs any spectr command in a project without `spectr.yaml`
- **THEN** the system uses `spectr/` as the root directory name
- **AND** all functionality works identically to current behavior
- **AND** no error or warning is displayed about missing config

#### Scenario: Custom root directory via config
- **WHEN** `spectr.yaml` exists with `root_dir: specs`
- **AND** user runs any spectr command
- **THEN** the system uses `specs/` instead of `spectr/` for all operations
- **AND** changes are discovered in `specs/changes/`
- **AND** specifications are discovered in `specs/specs/`
- **AND** archives are stored in `specs/changes/archive/`

#### Scenario: Config file format
- **WHEN** user creates a `spectr.yaml` file
- **THEN** the file uses YAML format
- **AND** the `root_dir` field is optional (defaults to `spectr`)
- **AND** comments are supported for documentation
- **AND** unknown fields are ignored with a warning

### Requirement: Config File Discovery

The system SHALL discover the configuration file by walking up the directory tree from the current working directory, enabling commands to work from any subdirectory within a project.

#### Scenario: Config discovery from project root
- **WHEN** user runs spectr command from project root directory
- **AND** `spectr.yaml` exists in the current directory
- **THEN** the config file is found immediately
- **AND** settings are applied

#### Scenario: Config discovery from subdirectory
- **WHEN** user runs spectr command from `src/components/` subdirectory
- **AND** `spectr.yaml` exists at the project root (two levels up)
- **THEN** the system walks up the directory tree
- **AND** finds and uses the config file at project root
- **AND** the project root is determined as the directory containing `spectr.yaml`

#### Scenario: Config discovery stops at filesystem root
- **WHEN** user runs spectr command outside any project
- **AND** no `spectr.yaml` exists in any parent directory up to filesystem root
- **THEN** the system uses default configuration
- **AND** looks for `spectr/` directory in current working directory

#### Scenario: Multiple config files in hierarchy
- **WHEN** `spectr.yaml` exists at both `/project/spectr.yaml` and `/project/subdir/spectr.yaml`
- **AND** user runs command from `/project/subdir/`
- **THEN** the nearest config file (`/project/subdir/spectr.yaml`) is used
- **AND** parent config files are ignored

### Requirement: Config Validation

The system SHALL validate configuration settings and provide clear error messages for invalid configurations.

#### Scenario: Root directory does not exist
- **WHEN** `spectr.yaml` specifies `root_dir: myspecs`
- **AND** the `myspecs/` directory does not exist
- **THEN** display error: "configured root directory 'myspecs/' does not exist"
- **AND** suggest running `spectr init` or creating the directory

#### Scenario: Invalid YAML syntax
- **WHEN** `spectr.yaml` contains invalid YAML syntax
- **THEN** display error with line number and column
- **AND** include snippet of problematic content
- **AND** suggest how to fix the syntax

#### Scenario: Invalid root_dir value
- **WHEN** `root_dir` contains invalid characters (e.g., `/`, `..`, `*`)
- **THEN** display error: "root_dir must be a simple directory name without path separators"
- **AND** list the invalid characters found
