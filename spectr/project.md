# spectr Context

## Purpose
Spectr is a CLI tool for validatable spec-driven development, inspired by OpenSpec and Kiro. It helps teams manage specifications and changes through a structured workflow: creating change proposals with delta specs, validating them against strict rules, implementing the changes, and archiving them after deployment. The tool enforces clear separation between current truth (specs/) and proposed changes (changes/), ensuring all modifications are intentional, documented, and validated.

## Tech Stack
- **Language**: Go 1.25.0
- **CLI Framework**: Kong (github.com/alecthomas/kong v1.13.0)
- **TUI Framework**: Bubbletea (github.com/charmbracelet/bubbletea v1.3.10)
- **Styling**: Lipgloss (github.com/charmbracelet/lipgloss v1.1.0)
- **Linting**: golangci-lint with comprehensive linters (asasalint, exhaustive, bidichk, gocritic, staticcheck, revive, etc.)
- **Build/Release**: GoReleaser
- **Package Manager**: Nix (with flakes and direnv)

## Project Conventions

### Code Style
- **Formatting**: Standard Go formatting (gofmt)
- **Linting**: Strict linting via golangci-lint with severity set to "error"
- **Naming**:
  - Packages: lowercase, single-word (e.g., `validation`, `parsers`, `archive`)
  - Interfaces: Descriptive nouns (e.g., `Validator`, `Parser`)
  - Change IDs: kebab-case with verb-led prefixes (`add-`, `update-`, `remove-`, `refactor-`)
- **Comments**: Exported types and functions must have doc comments
- **Error Handling**: Explicit error returns with context using `fmt.Errorf` wrapping

### Architecture Patterns
- **Clean Architecture**: Clear separation between cmd/, internal/, and domain logic
- **Packages**:
  - `cmd/`: CLI command definitions and handlers (thin layer)
  - `internal/init/`: Initialization wizard and setup logic
  - `internal/validation/`: Spec and change validation rules
  - `internal/parsers/`: Requirement and delta parsing
  - `internal/archive/`: Archive workflow and spec merging
  - `internal/list/`: Listing and formatting logic
  - `internal/discovery/`: File discovery utilities
- **Dependency Flow**: cmd → internal packages (no circular dependencies)
- **Interactive Mode**: Bubbletea TUI for wizards, with non-interactive fallback
- **Validation**: Strict mode available for comprehensive checks
- **Single Responsibility**: Each package has a focused purpose

### Testing Strategy
- **Unit Tests**: Table-driven tests for all core logic
- **Test Files**: Co-located with source files (`*_test.go`)
- **Coverage**: Aim for high coverage (current: 9814 lines in coverage.out)
- **Test Organization**:
  - Subtests with `t.Run()` for different scenarios
  - Descriptive test names (e.g., `TestValidator_ValidateSpec_ValidSpec`)
  - Setup/teardown using temp directories for filesystem tests
- **Validation**: Tests must verify actual behavior, not just pass
- **Test Data**: Store in `testdata/` directory
- **Assertions**: Explicit error checking, no test helpers that obscure failures

### Git Workflow
- **Main Branch**: `main`
- **Feature Development**: Work in branches, merge to main
- **Commit Style**: Descriptive commits with context
- **Status**: Track changes via git status (currently has staged and untracked files)
- **Release**: Automated via GitHub Actions and GoReleaser
- **Hooks**: Pre-commit validation recommended

## Domain Context

### Spec-Driven Development
Spectr implements a three-stage workflow:

1. **Creating Changes** (`spectr/changes/`):
   - Scaffold proposals with `proposal.md`, `tasks.md`, optional `design.md`
   - Write delta specs using `## ADDED|MODIFIED|REMOVED|RENAMED Requirements`
   - Each requirement MUST have `#### Scenario:` with WHEN/THEN format
   - Validate with `spectr validate <change-id>`

2. **Implementing Changes**:
   - Follow tasks.md checklist sequentially
   - Mark tasks complete with `- [x]` after implementation
   - Approval gate: no implementation without proposal review

3. **Archiving Changes**:
   - Move to `changes/archive/YYYY-MM-DD-[name]/`
   - Merge deltas into `specs/`
   - Use `spectr archive <change-id>` with optional `--skip-specs` flag

### Key Concepts
- **Specs**: Current truth - what IS built (`spectr/specs/`)
- **Changes**: Proposals - what SHOULD change (`spectr/changes/`)
- **Capabilities**: Single focused feature areas (e.g., `cli-framework`, `validation`)
- **Deltas**: Proposed changes to specs (ADDED/MODIFIED/REMOVED/RENAMED)
- **Requirements**: SHALL/MUST statements with scenarios
- **Scenarios**: WHEN/THEN test cases for requirements

## Important Constraints

### Validation Rules
- Every requirement MUST have at least one scenario
- Scenarios MUST use `#### Scenario:` format (4 hashtags, not bullets)
- MODIFIED requirements MUST include full updated content
- Delta operations MUST use exact header formats
- Strict mode treats warnings as errors

### File Structure
- Change directories MUST have at least one delta spec
- Spec files MUST have `## Requirements` section
- Delta specs MUST be in `changes/<id>/specs/<capability>/spec.md`
- Archive preserves history in `archive/YYYY-MM-DD-<id>/`

### Development
- Simplicity first: default to <100 lines of new code
- Single-file implementations until proven insufficient
- Avoid frameworks without clear justification
- 10-minute understandability rule for capabilities

## External Dependencies
- **None** - Spectr is a self-contained CLI tool with no external API dependencies
- All dependencies are Go libraries managed via `go.mod`
- Nix flake manages development environment setup
- GoReleaser handles binary distribution
