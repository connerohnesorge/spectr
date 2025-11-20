## MODIFIED Requirements

### Requirement: Backward-Compatible CLI Interface
The CLI framework SHALL maintain the same command syntax and flag names as the previous implementation.

The init command structure and related types SHALL be organized within the internal/init package to support cleaner separation of concerns and enable the init package to accept command configurations directly.

#### Scenario: Init command compatibility
- **WHEN** users invoke `spectr init` with existing flag combinations
- **THEN** the behavior SHALL be identical to the previous implementation
- **AND** all flag names SHALL remain unchanged
- **AND** short flag aliases SHALL remain unchanged
- **AND** positional argument handling SHALL remain unchanged
- **AND** the `InitCmd` struct is defined in `internal/init/models.go`
- **AND** `InitCmd` is exported with proper documentation
- **AND** the struct contains fields: Path, PathFlag, Tools, NonInteractive
- **AND** the CLI root struct imports `InitCmd` from `github.com/conneroisu/spectr/internal/init`
- **AND** no `InitCmd` definition exists in `cmd/init.go`

#### Scenario: Help text accessibility
- **WHEN** users invoke `spectr --help` or `spectr init --help`
- **THEN** help information SHALL be displayed (format may differ from Cobra)
- **AND** all commands and flags SHALL be documented
- **AND** help text matches the current implementation exactly

#### Scenario: Init command organizational structure
- **WHEN** the init command handler is invoked
- **THEN** the `Run()` method receives the complete `InitCmd` struct
- **AND** internal functions can accept `*InitCmd` as a parameter
- **AND** this enables direct passing of command configuration to init package functions

#### Scenario: Clean dependency flow
- **WHEN** the codebase is analyzed for circular imports
- **THEN** `cmd/init.go` imports from `internal/init`
- **AND** `internal/init` does not import from `cmd/`
- **AND** the dependency graph remains: `cmd` → `internal/init` (unidirectional)
- **AND** import uses alias `initpkg "github.com/conneroisu/spectr/internal/init"`
- **AND** types are referenced as `initpkg.InitCmd` to avoid conflicts with `init` keyword
