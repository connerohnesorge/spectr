# Delta Specification

## ADDED Requirements

### Requirement: Completion Command Structure

The CLI SHALL provide a `completion` subcommand that outputs shell completion
scripts for supported shells using the kong-completion library.

#### Scenario: Completion command registration

- **WHEN** the CLI is initialized
- **THEN** it SHALL include a Completion field using `kongcompletion.Completion`
  type
- **AND** the command SHALL be accessible via `spectr completion`
- **AND** help text SHALL describe shell completion functionality

#### Scenario: Bash completion output

- **WHEN** user runs `spectr completion bash`
- **THEN** the system outputs a valid bash completion script
- **AND** the script can be sourced directly or added to bash-completion.d

#### Scenario: Zsh completion output

- **WHEN** user runs `spectr completion zsh`
- **THEN** the system outputs a valid zsh completion script
- **AND** the script can be sourced or placed in $fpath

#### Scenario: Fish completion output

- **WHEN** user runs `spectr completion fish`
- **THEN** the system outputs a valid fish completion script
- **AND** the script can be saved to fish completions directory

### Requirement: Custom Predictors for Dynamic Arguments

The completion system SHALL provide context-aware suggestions for arguments that
accept dynamic values like change IDs or spec IDs.

#### Scenario: Change ID completion

- **WHEN** user types `spectr archive <TAB>` or `spectr validate <TAB>`
- **AND** the argument expects a change ID
- **THEN** completion suggests all active change IDs from `spectr/changes/`
- **AND** excludes archived changes

#### Scenario: Spec ID completion

- **WHEN** user types `spectr validate --type spec <TAB>`
- **AND** the argument expects a spec ID
- **THEN** completion suggests all spec IDs from `spectr/specs/`

#### Scenario: Item type completion

- **WHEN** user types `spectr validate --type <TAB>`
- **THEN** completion suggests `change` and `spec`

### Requirement: Kong-Completion Integration Pattern

The CLI initialization SHALL follow the kong-completion pattern where Kong is
initialized, completions are registered, and then arguments are parsed.

#### Scenario: Initialization order

- **WHEN** the program starts
- **THEN** `kong.Must()` is called first to create the Kong application
- **AND** `kongcompletion.Register()` is called before parsing
- **AND** `app.Parse()` is called after completion registration
- **AND** this order ensures completions work correctly

#### Scenario: Predictor registration

- **WHEN** custom predictors are defined
- **THEN** they SHALL be registered via `kongcompletion.WithPredictor()`
- **AND** struct fields SHALL reference predictors using `predictor:"name"` tag
