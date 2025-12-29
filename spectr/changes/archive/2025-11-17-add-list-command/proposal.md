# Change: Add List Command for Enumerating Specs and Changes

## Why

Spectr currently lacks a way to enumerate existing specs and changes, forcing users to manually browse directories or use generic filesystem tools. Without a list command, users cannot easily discover what specs exist, what changes are in progress, or get a quick overview of project status. OpenSpec provides a clean, flexible list command that serves this need—Spectr should have similar functionality to improve discoverability and workflow efficiency.

## What Changes

- Add new `list` command to the CLI that enumerates changes by default
- Add `--specs` flag to list specifications instead of changes
- Add `--long` flag to show detailed output (ID, title, and counts)
- Add `--json` flag for machine-readable structured output
- Implement change discovery that finds active changes (excluding archive directory)
- Implement spec discovery that finds specs with valid spec.md files
- Display changes with delta counts and task completion status in long format
- Display specs with requirement counts in long format
- Sort all output alphabetically by ID for consistency

## Impact

- **Affected specs**: `cli-framework` (extends existing)
- **Affected code**:
  - `cmd/root.go` - Add ListCmd struct to CLI
  - New `cmd/list.go` - List command implementation
  - New `internal/list/` - List functionality package
    - `lister.go` - Core listing logic for changes and specs
    - `types.go` - Data structures for list output
    - `formatters.go` - Text and JSON output formatting
  - New `internal/discovery/` - Item discovery utilities (or extend if created by validate command)
    - `discovery.go` - Find changes and specs in project
  - New `internal/parsers/` - Basic markdown parsing
    - `proposal_parser.go` - Extract title from proposal.md
    - `spec_parser.go` - Extract title from spec.md
    - `task_parser.go` - Count tasks in tasks.md

## Benefits

- **Discoverability**: Users can quickly see all available specs and changes
- **Workflow efficiency**: No need to manually navigate directories
- **Consistency**: Matches OpenSpec's proven list command UX
- **Machine-readable**: JSON output enables scripting and automation
- **Status visibility**: Task completion progress visible in long format for changes
