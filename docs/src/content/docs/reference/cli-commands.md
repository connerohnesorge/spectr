---
title: CLI Commands
description: Complete reference for all Spectr command-line interface commands.
---

This page documents all available Spectr CLI commands and their options.

## Essential Commands

### List Changes and Specs

List all active changes and specifications:

```bash
spectr list
```

List specifications with detailed information:

```bash
spectr list --specs
spectr spec list --long
```

**Options:**
- `--specs` - Show specifications instead of changes
- `--long` - Show detailed information
- `--json` - Machine-readable JSON output

### Show Details

Display details about a specific change or spec:

```bash
spectr show <change-id>
spectr show <spec-id> --type spec
```

**Examples:**
```bash
# Show change details
spectr show add-two-factor-auth

# Show spec details
spectr show auth --type spec

# JSON output with delta details
spectr show add-two-factor-auth --json --deltas-only
```

**Options:**
- `--type change|spec` - Specify what to show (required if ambiguous)
- `--json` - Machine-readable JSON output
- `--deltas-only` - Show only spec deltas

### Validate Changes and Specs

Validate a change or specification:

```bash
spectr validate <change-id>
spectr validate --strict
```

**Examples:**
```bash
# Validate single change
spectr validate add-two-factor-auth --strict

# Bulk validation (interactive)
spectr validate

# Validate all specs
spectr validate --specs
```

**Options:**
- `--strict` - Comprehensive validation with all checks
- `--specs` - Validate specs instead of changes
- `--no-interactive` - Disable prompts

### Archive a Change

Move a completed change to archive and merge specs:

```bash
spectr archive <change-id> --yes
```

**Examples:**
```bash
# Interactive archiving
spectr archive add-two-factor-auth

# Non-interactive (use --yes flag)
spectr archive add-two-factor-auth --yes

# Archive without updating specs
spectr archive add-two-factor-auth --skip-specs --yes
```

**Options:**
- `--yes`, `-y` - Skip confirmation prompts
- `--skip-specs` - Archive without merging specs
- `--no-interactive` - Disable all prompts

## Project Management

### Initialize Spectr

Initialize a new Spectr project:

```bash
spectr init [path]
```

**Examples:**
```bash
# Initialize in current directory
spectr init

# Initialize in specific directory
spectr init ./docs
```

### Update Instructions

Update Spectr instruction files:

```bash
spectr update [path]
```

This updates the instruction markdown files in the project.

## Global Options

Most commands support these global options:

```bash
spectr [command] [options]
```

- `--help` - Show command help
- `--version` - Show Spectr version
- `--json` - Output as JSON (where supported)
- `--no-interactive` - Disable interactive prompts

## Interactive Mode

Some commands support interactive mode:

```bash
# Interactive spec selection
spectr show

# Interactive validation
spectr validate

# Interactive archiving
spectr archive
```

## Output Formats

### Text Output (Default)

Human-readable text format:

```bash
spectr list
```

Output:
```
Active Changes:
├─ add-two-factor-auth (created 2025-01-15)
│  ├─ status: pending
│  └─ affected: auth, notifications
└─ update-api-versioning (created 2025-01-10)
```

### JSON Output

Machine-readable JSON format:

```bash
spectr list --json
```

Output:
```json
{
  "changes": [
    {
      "id": "add-two-factor-auth",
      "created": "2025-01-15T10:30:00Z",
      "status": "pending",
      "affectedSpecs": ["auth", "notifications"]
    }
  ]
}
```

## Common Workflows

### Create and Validate a Change

```bash
# 1. Create change directory and files
mkdir -p spectr/changes/add-feature/specs/auth

# 2. Write proposal.md, tasks.md, and spec.md files

# 3. Validate the change
spectr validate add-feature --strict

# 4. Fix any issues shown in validation

# 5. Validate again
spectr validate add-feature --strict
```

### Implement and Archive

```bash
# 1. Check change details
spectr show add-feature

# 2. Implement according to tasks.md

# 3. Mark all tasks complete in tasks.md

# 4. Validate before archiving
spectr validate add-feature --strict

# 5. Archive the change
spectr archive add-feature --yes

# 6. Verify specs were updated
spectr validate --strict
```

### Explore Project State

```bash
# 1. List all active changes
spectr list

# 2. List all specs
spectr list --specs

# 3. Show details of a change
spectr show add-feature

# 4. Show details of a spec
spectr show auth --type spec

# 5. Validate everything
spectr validate --strict
```

## Troubleshooting

### Command not found

If you get "command not found", ensure Spectr is installed:

```bash
# Check installation
which spectr

# Install if needed (see Installation guide)
```

### Ambiguous item

If you get "ambiguous item" error, specify the type:

```bash
# Error: Could be change or spec
spectr show auth

# Solution: Specify type
spectr show auth --type spec
```

### Permission denied

If you get permission errors when creating files:

```bash
# Check directory permissions
ls -la spectr/

# Ensure you have write permission
chmod -R u+w spectr/
```

## Tips and Tricks

### Piping JSON

Use JSON output with `jq` for filtering:

```bash
# Get all change IDs
spectr list --json | jq '.changes[].id'

# Get specific change details
spectr show add-feature --json | jq '.proposal'
```

### Quick Validation Loop

Validate multiple times while editing:

```bash
# Watch validation results
while true; do
  clear
  spectr validate add-feature --strict
  sleep 2
done
```

### List Specs by Capability

```bash
# Show all specs with details
spectr spec list --long

# Find specs matching a pattern
spectr spec list --json | jq '.specs[] | select(.id | contains("auth"))'
```

## Further Reading

- [Creating Changes](/guides/creating-changes)
- [Archiving Workflow](/guides/archiving-workflow)
- [Configuration](/reference/configuration)
