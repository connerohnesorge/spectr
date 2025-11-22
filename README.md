# spectr

<img src="https://github.com/connerohnesorge/spectr/blob/main/assets/logo.png" alt="Logo" width="95">

**Validatable spec-driven development (inspired by openspec and kiro)**

Tired of your specs disappearing like a ghost? `spectr archive` is your friend - it merges your change deltas into spec files so nothing gets lost.

Built with Go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](https://go.dev/)

---

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Installation](#installation)
  - [Using Nix Flakes](#using-nix-flakes)
  - [Building from Source](#building-from-source)
  - [Requirements](#requirements)
- [Quick Start](#quick-start)
  - [Initialize a Project](#initialize-a-project)
  - [Create Your First Change](#create-your-first-change)
  - [File Structure](#file-structure)
- [Command Reference](#command-reference)
  - [spectr init](#spectr-init)
  - [spectr list](#spectr-list)
  - [spectr validate](#spectr-validate)
  - [spectr archive](#spectr-archive)
  - [spectr view](#spectr-view)
- [Architecture & Development](#architecture--development)
  - [Architecture Overview](#architecture-overview)
  - [Package Structure](#package-structure)
  - [Development Setup](#development-setup)
  - [Testing Strategy](#testing-strategy)
- [Contributing](#contributing)
  - [Contribution Workflow](#contribution-workflow)
  - [Code Style Guidelines](#code-style-guidelines)
  - [Commit Conventions](#commit-conventions)
  - [Testing Requirements](#testing-requirements)
- [Advanced Topics](#advanced-topics)
  - [Spec-Driven Development](#spec-driven-development)
  - [Delta Specifications](#delta-specifications)
  - [Validation Rules](#validation-rules)
  - [Archiving Workflow](#archiving-workflow)
- [Troubleshooting](#troubleshooting)
  - [Common Issues](#common-issues)
  - [FAQ](#faq)
- [Links & Resources](#links--resources)

---

## Overview

**Spectr** is a CLI tool for validatable spec-driven development. It helps teams manage specifications and changes through a structured three-stage workflow:

1. **Creating Changes**: Write proposals with delta specs showing what SHOULD change
2. **Implementing Changes**: Follow the implementation checklist in `tasks.md`
3. **Archiving Changes**: Merge deltas into specs, preserving history

Spectr enforces a clear separation between current truth (`specs/` - what IS built) and proposed changes (`changes/` - what SHOULD change), ensuring all modifications are intentional, documented, and validated.

## Why "spectr"?

You might wonder: does the name of your specs folder and CLI tool actually matter? We thought the same thing, so we tested it systematically across multiple AI models (Claude, GPT, and others) to find what works best.

### Alternatives Evaluated

We considered several naming approaches:
- `specs/` - Common but generic; easily confused with documentation folders
- `specifications/` - Descriptive but verbose; slower to type in CLI workflows
- `requirements/` - Often associated with waterfall methodologies; less distinct
- `docs/specs/` - Nested approach; less CLI-friendly
- **`spectr/` - What we chose**

### Why spectr Won

Testing across AI models revealed that `spectr` excels in several dimensions:

1. **Brevity**: Fast to type, easy to remember, works well as a CLI command
2. **Distinctiveness**: Unique enough to avoid conflicts with common folder names
3. **AI Compatibility**: Models consistently recognize and use it correctly in generated code
4. **Ergonomic**: Short enough to feel natural in commands like `spectr validate` and `spectr archive`

The name reflects the tool's core purpose—providing clear visibility into specifications and changes, like a spectrum analyzer revealing what IS and what SHOULD BE.

## Key Features

- **Structured Workflow**: Propose, validate, implement, and archive changes systematically
- **Delta Specifications**: Track proposed changes separately from current specs
- **Strict Validation**: Enforce requirements format, scenarios, and spec consistency
- **Interactive TUI**: Beautiful terminal UI for wizards and selection flows
- **Archive Merging**: Automatically merge change deltas into spec files with `spectr archive`
- **Clean Architecture**: Well-organized codebase with clear separation of concerns
- **Comprehensive Testing**: Table-driven tests with high coverage
- **Nix Integration**: First-class Nix flake support for reproducible builds

---

## Installation

### Using Nix Flakes

The recommended way to install Spectr is via Nix flakes:

```bash
# Run directly without installing
nix run github:connerohnesorge/spectr

# Install to your profile
nix profile install github:connerohnesorge/spectr

# Add to your flake.nix inputs
{
  inputs.spectr.url = "github:connerohnesorge/spectr";
}
```

### Building from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/connerohnesorge/spectr.git
cd spectr

# Build with Go
go build -o spectr

# Or use Nix
nix build

# Install to your PATH
mv spectr /usr/local/bin/  # or any directory in your PATH
```

### Requirements

- **Go 1.25+** (if building from source)
- **Nix with flakes enabled** (optional, for Nix installation)
- **Git** (for project version control)

---

## Quick Start

### Initialize a Project

Start by initializing Spectr in your project:

```bash
# Initialize with interactive wizard
spectr init

# Or specify a path
spectr init /path/to/project

# Non-interactive mode with defaults
spectr init --non-interactive
```

This creates the following structure:

```
your-project/
└── spectr/
    ├── project.md        # Project conventions and context
    ├── specs/            # Current specifications (truth)
    │   └── [capability]/ # One directory per capability
    │       ├── spec.md   # Requirements and scenarios
    │       └── design.md # Technical patterns (optional)
    └── changes/          # Proposed changes
        └── archive/      # Completed changes
```

### Create Your First Change

Let's create a simple "Hello World" change:

```bash
# 1. List current state
spectr list              # See active changes
spectr list --specs      # See existing capabilities

# 2. Create a change directory
mkdir -p spectr/changes/add-hello-world/specs/greeting

# 3. Write a proposal
cat > spectr/changes/add-hello-world/proposal.md << 'EOF'
# Change: Add Hello World Greeting

## Why
We need a simple greeting capability to welcome users.

## What Changes
- Add new `greeting` capability with hello world functionality

## Impact
- Affected specs: greeting (new)
- Affected code: None (example)
EOF

# 4. Create delta spec
cat > spectr/changes/add-hello-world/specs/greeting/spec.md << 'EOF'
## ADDED Requirements

### Requirement: Hello World Greeting
The system SHALL provide a greeting function that returns "Hello, World!".

#### Scenario: Greet successfully
- **WHEN** the greeting function is called
- **THEN** it SHALL return "Hello, World!"
EOF

# 5. Create tasks checklist
cat > spectr/changes/add-hello-world/tasks.md << 'EOF'
## 1. Implementation
- [ ] 1.1 Create greeting.go file
- [ ] 1.2 Implement HelloWorld() function
- [ ] 1.3 Write tests for greeting
- [ ] 1.4 Update documentation
EOF

# 6. Validate the change
spectr validate add-hello-world --strict

# 7. After implementation, archive it
spectr archive add-hello-world
```

### File Structure

Understanding the directory structure is crucial:

```
spectr/
├── project.md              # Project-wide conventions
├── specs/                  # CURRENT TRUTH - what IS built
│   └── [capability]/
│       ├── spec.md         # Requirements with scenarios
│       └── design.md       # Technical patterns (optional)
├── changes/                # PROPOSALS - what SHOULD change
│   ├── [change-id]/
│   │   ├── proposal.md     # Why, what, impact
│   │   ├── tasks.md        # Implementation checklist
│   │   ├── design.md       # Technical decisions (optional)
│   │   └── specs/          # Delta changes
│   │       └── [capability]/
│   │           └── spec.md # ADDED/MODIFIED/REMOVED requirements
│   └── archive/            # Completed changes (history)
│       └── YYYY-MM-DD-[change-id]/
```

**Key Concepts:**
- **specs/**: The source of truth for what's currently built
- **changes/**: Proposed modifications, kept separate until approved
- **archive/**: Historical record of all changes with timestamps
- **Delta Specs**: Use `## ADDED`, `## MODIFIED`, `## REMOVED`, or `## RENAMED Requirements` headers

---

## Command Reference

### spectr init

Initialize Spectr in a project directory.

**Usage:**
```bash
spectr init [PATH] [FLAGS]
```

**Flags:**
- `--tools <tools>`: Comma-separated list of tools to include (e.g., `git,github`)
- `--non-interactive`: Skip interactive wizard, use defaults
- `--path <path>`: Project directory (default: current directory)

**Examples:**
```bash
# Interactive initialization (recommended)
spectr init

# Initialize specific directory with Git integration
spectr init /path/to/project --tools git

# Non-interactive with defaults
spectr init --non-interactive
```

**Output:**
```
✓ Created spectr/ directory
✓ Created specs/ directory
✓ Created changes/ directory
✓ Created project.md
✓ Spectr initialized successfully!
```

### spectr list

List active changes or specifications.

**Usage:**
```bash
spectr list [FLAGS]
```

**Flags:**
- `--specs`: List specifications instead of changes
- `--json`: Output in JSON format
- `--long`: Show detailed information
- `--no-interactive`: Disable interactive selection

**Examples:**
```bash
# List all active changes
spectr list

# List all specifications
spectr list --specs

# Get detailed JSON output
spectr list --json --long

# List specs with full details
spectr list --specs --long
```

**Example Output:**
```
Active Changes:
  add-two-factor-auth    Add 2FA authentication support
  refactor-validation    Improve validation error messages

Run 'spectr show <change>' for details
```

### spectr validate

Validate changes or specifications against rules.

**Usage:**
```bash
spectr validate [ITEM] [FLAGS]
```

**Flags:**
- `--strict`: Enable strict validation (warnings become errors)
- `--type <change|spec>`: Disambiguate when name conflicts exist
- `--json`: Output validation results as JSON
- `--no-interactive`: Skip interactive mode

**Examples:**
```bash
# Validate a specific change (strict mode recommended)
spectr validate add-two-factor-auth --strict

# Validate all changes interactively
spectr validate

# Validate a specification
spectr validate auth --type spec

# Get JSON validation results
spectr validate add-2fa --json
```

**Validation Rules:**
- Every requirement MUST have at least one scenario
- Scenarios MUST use `#### Scenario:` format (4 hashtags)
- Purpose sections MUST be at least 50 characters
- MODIFIED requirements MUST include complete updated content
- Change directories MUST contain at least one delta spec

**Example Output:**
```
Validating change: add-two-factor-auth

✓ Proposal file exists
✓ Delta specs found
✓ All requirements have scenarios
✓ Scenario formatting correct
✓ All validations passed!
```

### spectr archive

Archive a completed change, merging deltas into specs.

**Usage:**
```bash
spectr archive <CHANGE-ID> [FLAGS]
```

**Flags:**
- `--skip-specs`: Archive without updating specs (for tooling-only changes)
- `--yes` / `-y`: Skip confirmation prompts (non-interactive)
- `--no-interactive`: Disable interactive mode

**Examples:**
```bash
# Archive with interactive confirmation
spectr archive add-two-factor-auth

# Archive without updating specs
spectr archive fix-typo --skip-specs

# Non-interactive archive (for CI/CD)
spectr archive add-feature --yes
```

**What It Does:**
1. Validates the change before archiving
2. Merges delta specs into `specs/` (unless `--skip-specs`)
3. Moves `changes/[name]` → `changes/archive/YYYY-MM-DD-[name]`
4. Preserves complete history in archive

**Example Output:**
```
Archiving change: add-two-factor-auth

✓ Validation passed
✓ Merging deltas into specs/auth/spec.md
  - Added 2 requirements
  - Modified 1 requirement
✓ Moving to archive/2025-11-18-add-two-factor-auth/
✓ Archive complete!
```

### spectr view

Display detailed information about a change or spec.

**Usage:**
```bash
spectr view [ITEM] [FLAGS]
```

**Flags:**
- `--type <change|spec>`: Specify item type
- `--json`: Output in JSON format
- `--deltas-only`: Show only delta specifications (changes only)

**Examples:**
```bash
# View a change interactively
spectr view

# View specific change
spectr view add-two-factor-auth

# View spec details
spectr view auth --type spec

# Debug delta parsing
spectr view add-2fa --json --deltas-only
```

**Example Output:**
```
Change: add-two-factor-auth
Status: Active

Proposal:
  Add two-factor authentication support via OTP

Affected Specs:
  - auth
  - notifications

Tasks: 4 total, 2 completed

Delta Summary:
  auth:
    - ADDED: 2 requirements
    - MODIFIED: 1 requirement
```

---

## Architecture & Development

### Architecture Overview

Spectr follows **Clean Architecture** principles with clear separation of concerns:

```
spectr/
├── cmd/                    # CLI command definitions (thin layer)
│   ├── root.go            # Kong CLI framework setup
│   ├── init.go            # Init command handler
│   ├── list.go            # List command handler
│   ├── validate.go        # Validate command handler
│   ├── archive.go         # Archive command handler
│   └── view.go            # View command handler
├── internal/              # Core business logic (not importable externally)
│   ├── init/             # Initialization wizard and setup
│   ├── validation/       # Spec and change validation rules
│   ├── parsers/          # Requirement and delta parsing
│   ├── archive/          # Archive workflow and spec merging
│   ├── list/             # Listing and formatting logic
│   ├── discovery/        # File discovery utilities
│   └── view/             # Display and formatting
├── main.go               # Application entry point
└── testdata/             # Test fixtures and integration tests
```

**Design Principles:**
- **Thin CLI Layer**: Commands delegate to internal packages
- **No Circular Dependencies**: Strict dependency flow from cmd → internal
- **Single Responsibility**: Each package has one focused purpose
- **Testability**: Logic separated from I/O for easy testing

### Package Structure

| Package | Purpose | Key Types |
|---------|---------|-----------|
| `cmd/` | CLI command handlers using Kong framework | Command structs |
| `internal/init/` | Project initialization wizard and templates | `Wizard`, `Executor`, `Templates` |
| `internal/validation/` | Validation rules for specs and changes | `Validator`, `Rule`, `ValidationResult` |
| `internal/parsers/` | Parse requirements, scenarios, and deltas | `RequirementParser`, `DeltaParser` |
| `internal/archive/` | Archive changes and merge deltas into specs | `Archiver`, `SpecMerger` |
| `internal/list/` | List changes and specs with formatting | `Lister`, `Formatter` |
| `internal/discovery/` | Discover spec and change files | `Discoverer`, `FileInfo` |
| `internal/view/` | Display detailed information with TUI | `Dashboard`, `ProgressTracker` |

### Development Setup

#### Using Nix (Recommended)

```bash
# Clone the repository
git clone https://github.com/connerohnesorge/spectr.git
cd spectr

# Enter development shell (provides all tools)
nix develop

# Available tools:
# - go_1_25: Go compiler and runtime
# - air: Live reload during development
# - gopls: Language server for IDE integration
# - golangci-lint: Comprehensive linting
# - gotestsum: Enhanced test output
# - delve: Debugger
```

#### Without Nix

```bash
# Install Go 1.25+
# Download from https://go.dev/dl/

# Clone repository
git clone https://github.com/connerohnesorge/spectr.git
cd spectr

# Install dependencies
go mod download

# Build
go build -o spectr

# Run
./spectr --help
```

### Testing Strategy

Spectr uses **table-driven tests** with high coverage:

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run with race detector
go test ./... -race

# Run with enhanced output (if gotestsum installed)
gotestsum --format testname

# Run specific package tests
go test ./internal/validation/...

# Run with verbose output
go test -v ./internal/parsers/...
```

**Test Organization:**
- **Unit Tests**: Co-located with source files (`*_test.go`)
- **Table-Driven**: Subtests with `t.Run()` for different scenarios
- **Integration Tests**: Located in `testdata/integration/`
- **Test Fixtures**: Stored in `testdata/` directory

**Example Test Structure:**
```go
func TestValidator_ValidateSpec(t *testing.T) {
    tests := []struct {
        name    string
        spec    *Spec
        wantErr bool
    }{
        {"valid spec", validSpec, false},
        {"missing scenarios", specNoScenarios, true},
        {"invalid format", malformedSpec, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.ValidateSpec(tt.spec)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, wantErr = %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## Contributing

We welcome contributions! Please follow these guidelines to ensure smooth collaboration.

### Contribution Workflow

1. **Fork the Repository**
   ```bash
   # Click "Fork" on GitHub, then clone your fork
   git clone https://github.com/YOUR-USERNAME/spectr.git
   cd spectr
   ```

2. **Create a Feature Branch**
   ```bash
   git checkout -b add-new-feature
   ```

3. **Make Changes**
   - Follow code style guidelines
   - Write tests for new functionality
   - Update documentation as needed

4. **Run Tests and Linting**
   ```bash
   go test ./...
   golangci-lint run
   ```

5. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "Add new validation rule for scenarios"
   ```

6. **Push and Create Pull Request**
   ```bash
   git push origin add-new-feature
   # Create PR on GitHub
   ```

### Code Style Guidelines

- **Formatting**: Use `gofmt` (or `gofumpt` for stricter formatting)
- **Linting**: All code must pass `golangci-lint run`
- **Naming Conventions**:
  - Packages: lowercase, single-word (e.g., `validation`, `parsers`)
  - Interfaces: Descriptive nouns (e.g., `Validator`, `Parser`)
  - Exported functions: Clear, verb-led names (e.g., `ValidateSpec`, `ParseRequirement`)
- **Comments**: All exported types and functions MUST have doc comments
- **Error Handling**: Use explicit error returns with context via `fmt.Errorf` wrapping

**Example:**
```go
// ValidateSpec checks if a specification meets all validation rules.
// It returns an error if any rule is violated.
func ValidateSpec(spec *Spec) error {
    if spec == nil {
        return fmt.Errorf("spec cannot be nil")
    }
    // validation logic...
}
```

### Commit Conventions

Use clear, descriptive commit messages:

```
<type>: <short summary>

<optional body>

<optional footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `docs`: Documentation changes
- `chore`: Maintenance tasks

**Examples:**
```
feat: add strict validation mode

Implement --strict flag for validate command that treats
warnings as errors. Useful for CI/CD pipelines.

fix: correct scenario header parsing

Scenarios with extra whitespace were not being recognized.
Updated regex to trim whitespace before matching.

docs: update README with archive examples
```

### Testing Requirements

All contributions MUST include appropriate tests:

- **New Features**: Add unit tests covering success and error cases
- **Bug Fixes**: Add regression test demonstrating the fix
- **Refactoring**: Ensure existing tests still pass
- **Test Coverage**: Aim for high coverage (current: >80%)

**Running Tests Before PR:**
```bash
# Run all tests
go test ./...

# Check coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Run linting
golangci-lint run

# Format code
go fmt ./...
```

---

## Advanced Topics

### Spec-Driven Development

Spectr implements a **three-stage workflow** for managing changes:

#### Stage 1: Creating Changes
Create a proposal when you need to:
- Add features or functionality
- Make breaking changes (API, schema)
- Change architecture or patterns
- Optimize performance (changes behavior)
- Update security patterns

**Skip proposals for:**
- Bug fixes (restore intended behavior)
- Typos, formatting, comments
- Dependency updates (non-breaking)
- Tests for existing behavior

#### Stage 2: Implementing Changes
1. Read `proposal.md` - Understand what's being built
2. Read `design.md` (if exists) - Review technical decisions
3. Read `tasks.md` - Get implementation checklist
4. Implement tasks sequentially
5. Mark tasks complete with `- [x]` after implementation
6. **Approval gate**: Do not implement until proposal is approved

#### Stage 3: Archiving Changes
After deployment:
1. Run `spectr validate <change> --strict` to ensure quality
2. Run `spectr archive <change>` to merge deltas into specs
3. Changes move to `archive/YYYY-MM-DD-<change>/`
4. Specs in `specs/` are updated with merged requirements

### Delta Specifications

**Delta specs** describe proposed changes using operation headers:

```markdown
## ADDED Requirements
### Requirement: New Feature
The system SHALL provide new functionality.

#### Scenario: Success case
- **WHEN** condition occurs
- **THEN** expected result

## MODIFIED Requirements
### Requirement: Existing Feature
[Complete modified requirement with all scenarios]

## REMOVED Requirements
### Requirement: Deprecated Feature
**Reason**: Why removing
**Migration**: How to handle existing usage

## RENAMED Requirements
- FROM: `### Requirement: Old Name`
- TO: `### Requirement: New Name`
```

**Key Rules:**
- **ADDED**: New capabilities that stand alone
- **MODIFIED**: Changes to existing requirements (include FULL updated content)
- **REMOVED**: Deprecated features (provide reason and migration path)
- **RENAMED**: Name-only changes (use with MODIFIED if behavior changes too)

### Validation Rules

Spectr enforces strict validation rules to maintain quality:

| Rule | Description | Severity |
|------|-------------|----------|
| Requirement Scenarios | Every requirement MUST have ≥1 scenario | Error |
| Scenario Format | Scenarios MUST use `#### Scenario:` (4 hashtags) | Error |
| Purpose Length | Purpose sections MUST be ≥50 characters | Warning |
| MODIFIED Complete | MODIFIED requirements MUST be complete, not partial | Error |
| Delta Presence | Changes MUST have ≥1 delta spec | Error |
| Scenario Structure | Scenarios SHOULD have WHEN/THEN bullets | Warning |
| Header Matching | Operation headers use trim() - whitespace ignored | Info |

**Strict Mode:**
```bash
# Treat warnings as errors
spectr validate <change> --strict
```

**Debugging Validation:**
```bash
# See detailed validation output
spectr validate <change> --json | jq '.errors'

# Check delta parsing
spectr view <change> --json --deltas-only
```

### Archiving Workflow

The `spectr archive` command performs an atomic operation:

1. **Pre-Archive Validation**
   - Validates change structure
   - Checks all delta specs
   - Ensures requirements have scenarios

2. **Delta Merging**
   - Reads each delta spec in `changes/<id>/specs/`
   - For each capability:
     - **ADDED**: Appends to `specs/<capability>/spec.md`
     - **MODIFIED**: Replaces entire requirement block
     - **REMOVED**: Removes requirement (keeps comment)
     - **RENAMED**: Updates requirement header

3. **Archive Move**
   - Creates `changes/archive/YYYY-MM-DD-<id>/`
   - Moves entire change directory
   - Preserves all history (proposal, tasks, design, deltas)

4. **Post-Archive Verification**
   - Validates updated specs
   - Ensures merge was successful
   - Reports summary

**Example Archive:**
```bash
$ spectr archive add-two-factor-auth

Archiving change: add-two-factor-auth
✓ Validation passed
✓ Merging deltas:
  - specs/auth/spec.md: +2 ADDED, 1 MODIFIED
  - specs/notifications/spec.md: +1 ADDED
✓ Moved to archive/2025-11-18-add-two-factor-auth/
✓ Archive complete!
```

---

## Troubleshooting

### Common Issues

#### "Change must have at least one delta"

**Problem**: Validation fails because no delta specs found.

**Solution:**
1. Ensure `changes/<name>/specs/` directory exists
2. Create at least one `.md` file with delta operations
3. Verify files have `## ADDED`, `## MODIFIED`, `## REMOVED`, or `## RENAMED Requirements` headers

```bash
# Check delta structure
ls -la changes/my-change/specs/
cat changes/my-change/specs/*/spec.md | grep "^## "
```

#### "Requirement must have at least one scenario"

**Problem**: Requirement found without scenarios.

**Solution:**
Use the exact format for scenarios (4 hashtags, specific text):

```markdown
### Requirement: My Feature
The system SHALL do something.

#### Scenario: Success case
- **WHEN** user does X
- **THEN** system does Y
```

**Common Mistakes:**
- Using `###` (3 hashtags) instead of `####` (4 hashtags)
- Using bold `**Scenario:**` instead of header `####`
- Using bullets `- Scenario:` instead of header

#### Validation Errors in Strict Mode

**Problem**: `--strict` flag causes warnings to fail.

**Solution:**
1. Review warning messages carefully
2. Fix underlying issues (often scenario structure or purpose length)
3. Use non-strict mode during development: `spectr validate <change>`
4. Use strict mode before archiving: `spectr validate <change> --strict`

#### Archive Merge Conflicts

**Problem**: Multiple changes modify the same requirement.

**Solution:**
1. Archive changes sequentially, not in parallel
2. Resolve conflicts manually in `specs/` after first archive
3. Validate the second change after first is archived
4. Consider combining related changes into a single proposal

### FAQ

#### Do I need approval before implementing changes?

**Yes**. The approval gate is intentional. Changes should be reviewed and approved before implementation begins. This prevents wasted effort on changes that may be rejected or need significant revision.

#### How do I handle multiple capabilities in one change?

Create multiple delta specs, one per capability:

```
changes/add-2fa-notifications/
├── proposal.md
├── tasks.md
└── specs/
    ├── auth/
    │   └── spec.md       # Auth-related deltas
    └── notifications/
        └── spec.md       # Notification-related deltas
```

#### What's the difference between design.md in specs/ vs changes/?

- **specs/[capability]/design.md**: Current technical patterns for a capability
- **changes/[name]/design.md**: Design decisions for a proposed change

The change's `design.md` explains new architectural decisions. After archiving, relevant design details may be added to capability design docs.

#### Can I modify specs directly without a change?

**For minor fixes only**: typos, formatting, clarifications that don't change meaning.

**For everything else**: Create a change proposal. This ensures:
- Changes are reviewed and approved
- History is preserved in archive
- Validation catches errors before merging

#### How do I debug silent scenario parsing failures?

Use JSON output to see parsed structure:

```bash
# Check what was parsed
spectr view <change> --json --deltas-only | jq '.deltas[].requirements[].scenarios'

# Verify scenario count
spectr validate <change> --json | jq '.errors[] | select(.rule == "RequirementScenarios")'
```

#### What happens to archive/ directory over time?

Archives accumulate but remain organized by date:

```
changes/archive/
├── 2025-11-15-add-auth/
├── 2025-11-16-fix-validation/
├── 2025-11-18-add-notifications/
└── ...
```

Periodically, you may:
- Compress old archives
- Move ancient archives to separate storage
- Keep 6-12 months in active repository

---

## Links & Resources

- **GitHub Repository**: [github.com/connerohnesorge/spectr](https://github.com/connerohnesorge/spectr)
- **Specification Documentation**: See `spectr/specs/` for detailed capability specs
  - [CLI Interface](spectr/specs/cli-interface/spec.md)
  - [Validation Rules](spectr/specs/validation/spec.md)
  - [Archive Workflow](spectr/specs/archive-workflow/spec.md)
  - [CLI Framework](spectr/specs/cli-framework/spec.md)
- **AI Agents Documentation**: See [spectr/AGENTS.md](spectr/AGENTS.md) for AI assistant instructions
- **Project Conventions**: See [spectr/project.md](spectr/project.md)
- **Issue Tracker**: [GitHub Issues](https://github.com/connerohnesorge/spectr/issues)
- **Discussions**: [GitHub Discussions](https://github.com/connerohnesorge/spectr/discussions)

