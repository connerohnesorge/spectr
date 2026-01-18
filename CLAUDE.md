# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

All commands run inside `nix develop` shell (or use `nix develop -c '<command>'`):

```bash
# Build
go build -o spectr              # Build binary
nix build                       # Build via Nix

# Lint & Format
nix develop -c lint             # Run golangci-lint + markdownlint on spectr/
golangci-lint run               # Go linting only

# Test
nix develop -c tests            # Run tests with race detector (gotestsum)
go test ./...                   # Basic test run
go test -v ./internal/validation/...  # Specific package

# Single test
go test -run TestValidateSpec ./internal/validation/

# Format (via treefmt)
nix fmt                         # Format Go + Nix files
```

## Architecture Overview

Spectr is a CLI tool for spec-driven development. Key concepts:
- specs/ - Current truth: what IS built (requirements + scenarios)
- changes/ - Proposals: what SHOULD change (deltas against specs)
- archive/ - Completed changes with timestamps

### Code Structure

```
cmd/                    # CLI commands (thin layer using Kong framework)
├── root.go            # CLI setup with kong.Context
├── init.go            # spectr init
├── list.go            # spectr list
├── validate.go        # spectr validate
├── accept.go          # spectr accept (converts tasks.md → tasks.jsonc)
├── pr.go              # spectr pr archive|new (git worktree + PR creation)
└── view.go            # spectr view

internal/               # Business logic (not importable externally)
├── validation/        # Validation rules for specs and changes
├── parsers/           # Requirement and delta parsing from markdown
├── archive/           # Archive workflow and spec merging
├── discovery/         # File discovery utilities
├── initialize/        # Init wizard and AI tool templates
├── tui/               # Interactive terminal UI (Bubble Tea)
├── domain/            # Core domain types (Spec, Change, Requirement)
├── pr/                # Pull request creation via git worktree
└── git/               # Git operations
```

### Key Dependencies
- Kong: CLI framework (`github.com/alecthomas/kong`)
- Bubble Tea/Bubbles/Lipgloss: TUI framework (Charmbracelet)
- Afero: Filesystem abstraction for testing

### Data Flow
1. Commands in `cmd/` parse flags and call `internal/` packages
2. `discovery/` finds spec/change files in `spectr/` directory
3. `parsers/` extracts requirements and scenarios from markdown
4. `validation/` enforces rules (scenarios required, format checks)
5. `archive/` merges delta specs into main specs

## Spectr Workflow

See `spectr/AGENTS.md` for detailed spec-driven development instructions:
- Create proposals in `spectr/changes/<id>/` with `proposal.md`, `tasks.md`, delta specs
- Run `spectr validate <id>` before implementation
- Run `spectr accept <id>` to convert tasks.md to tasks.jsonc
- Track task status in `tasks.jsonc` during implementation
- Archive completed changes with `spectr pr archive <id>`

## Delta Spec Format

```markdown
## ADDED Requirements
### Requirement: Feature Name
The system SHALL...

#### Scenario: Success case
- WHEN condition
- THEN result

## MODIFIED Requirements
### Requirement: Existing Feature
[Complete updated requirement with all scenarios]

## REMOVED Requirements
### Requirement: Old Feature
Reason: Why removing
Migration: How to handle
```

## Testing Patterns

Tests use table-driven style with `t.Run()`:
```go
tests := []struct {
    name    string
    input   string
    wantErr bool
}{...}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {...})
}
```

Test fixtures in `testdata/`. TUI tests use `charmbracelet/x/exp/teatest`.

### Round-Trip Validation
For JSONC generation, tests verify data integrity using round-trip validation:
```go
// Marshal to JSONC
jsonData, _ := json.MarshalIndent(original, "", "  ")

// Parse back
stripped := parsers.StripJSONComments(jsonData)
var parsed TasksFile
json.Unmarshal(stripped, &parsed)

// Verify identical
if parsed.Field != original.Field {
    t.Errorf("round-trip failed")
}
```

Special character testing systematically covers JSON escape sequences (`\`, `"`, `\n`, `\t`, etc.).

## Orchestration Model

This project uses an orchestrator pattern for complex tasks:
- Orchestrator (you): Maintains big picture, creates todo lists, delegates
- coder agent: Implements ONE specific todo item (`.claude/agents/coder.md`)
- tester agent: Verifies implementations with Playwright
- stuck agent: Escalates to human when blocked

Workflow: Create todos → delegate to coder → verify with tester → mark complete

<!-- spectr:start -->
# Spectr Instructions

These instructions are for AI assistants working in this project.

Always open `@/spectr/AGENTS.md` when the request:

- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big
  performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/spectr/AGENTS.md` to learn:

- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

When delegating tasks from a change proposal to subagents:

- Provide the proposal path: `spectr/changes/<id>/proposal.md`
- Include task context: `spectr/changes/<id>/tasks.jsonc`
- Reference delta specs: `spectr/changes/<id>/specs/<capability>/spec.md`

<!-- spectr:end -->
