---
title: Configuration
description: Configure Spectr projects with project.md and other settings.
---

This page documents how to configure Spectr projects.

## Project Configuration

Spectr projects are configured using `spectr/project.md`, which defines conventions and guidelines for your project.

### Default Configuration

When you initialize a Spectr project with `spectr init`, it creates a default `spectr/project.md`:

```bash
spectr init /path/to/project
```

This creates:
- `spectr/project.md` - Project conventions
- `spectr/AGENTS.md` - AI assistant instructions
- `spectr/specs/` - Master specifications directory
- `spectr/changes/` - Changes and proposals directory

## Project.md Structure

The `project.md` file defines how changes and specs should be created:

```markdown
# Project: Spectr

## Capabilities

Spectr uses the following capability structure:
- Each capability has a name: `auth`, `api-design`, `webhooks`
- Each capability has one spec file: `spectr/specs/<capability>/spec.md`
- Each capability is independently versioned and modified

## Change Organization

Changes are organized in `spectr/changes/<change-id>/`:
- `proposal.md` - Why and what
- `tasks.md` - Implementation checklist
- `design.md` (optional) - Technical decisions
- `specs/` - Spec deltas for affected capabilities

## Naming Conventions

### Change IDs
- Use kebab-case: `add-two-factor-auth`
- Start with verb: `add-`, `update-`, `remove-`, `refactor-`
- Keep concise and descriptive
- Ensure uniqueness across all changes and archives

### Capability Names
- Use kebab-case: `auth`, `api-design`, `webhook-config`
- Single purpose per capability
- 10-minute understandability rule

## File Formats

### Spec Files
Specs use standard Markdown with structured requirements:

```markdown
### Requirement: Feature Name
Description of requirement.

#### Scenario: Success case
- **WHEN** user performs action
- **THEN** expected result
```

### Change Deltas
Deltas use section headers to indicate changes:

- `## ADDED Requirements` - New capabilities
- `## MODIFIED Requirements` - Changed behavior
- `## REMOVED Requirements` - Deprecated features
- `## RENAMED Requirements` - Name changes
```

## Customizing Your Project

### Change ID Naming Rules

Edit `project.md` to define change ID requirements:

```markdown
## Change ID Naming
- Format: kebab-case
- Must start with: add-, update-, remove-, refactor-
- Length: 20-50 characters
- Pattern: [verb]-[noun-phrase]
```

### Capability Naming Rules

```markdown
## Capability Naming
- Format: kebab-case
- Prefix: [domain]-[feature]
- Examples: auth-mfa, api-pagination, webhooks-retry
```

### Validation Rules

```markdown
## Validation Rules
- Every requirement must have at least 1 scenario
- Every scenario must use #### Scenario: format
- Change must affect at least 1 spec
- Spec deltas must use operation headers
```

## Directory Structure

After initialization, your project looks like:

```
your-project/
├── CLAUDE.md
├── spectr/
│   ├── project.md              # Configuration (you edit this)
│   ├── AGENTS.md               # AI instructions (you edit this)
│   ├── specs/                  # Master specifications
│   │   ├── auth/
│   │   │   └── spec.md
│   │   ├── api-design/
│   │   │   └── spec.md
│   │   └── webhooks/
│   │       └── spec.md
│   └── changes/                # Proposed and archived changes
│       ├── add-two-factor-auth/
│       │   ├── proposal.md
│       │   ├── tasks.md
│       │   └── specs/
│       │       └── auth/spec.md
│       └── archive/            # Completed changes
│           └── 2025-01-15-add-two-factor-auth/
└── ...rest of project...
```

## Best Practices

### Project Configuration

1. **Document conventions** - Make your naming rules explicit in `project.md`
2. **Define scope** - Clarify what capabilities your project manages
3. **Set validation rules** - Document what makes a valid spec or change
4. **Provide examples** - Include example change proposals

### Capability Design

1. **Single purpose** - Each capability does one thing well
2. **Clear boundaries** - No overlap with other capabilities
3. **Consistent naming** - Follow your naming conventions
4. **Meaningful names** - Someone should understand it in 10 minutes

### Change Organization

1. **Focused proposals** - One major feature per change
2. **Complete deltas** - All affected specs must be updated
3. **Clear requirements** - Write requirements that are testable
4. **Comprehensive tasks** - Break implementation into clear steps

## Environment Variables

Spectr respects these environment variables:

```bash
# Set working directory
export SPECTR_DIR=/path/to/project/spectr

# Enable debug logging
export DEBUG=spectr:*

# Set output format
export SPECTR_FORMAT=json
```

## Integration with Workflows

### CI/CD Integration

Validate changes in CI/CD:

```bash
# In your CI/CD pipeline
spectr validate --strict --no-interactive
```

### Git Hooks

Validate before commit:

```bash
#!/bin/bash
# .git/hooks/pre-commit
spectr validate --strict || exit 1
```

### GitHub Actions

Example GitHub Actions workflow:

```yaml
name: Validate Spectr Changes

on: [pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install Spectr
        run: go install github.com/conneroisu/spectr@latest
      - name: Validate
        run: spectr validate --strict --no-interactive
```

## Troubleshooting Configuration

### Invalid configuration format

Ensure `project.md` follows Markdown conventions:

```bash
# Check syntax
spectr validate --strict
```

### Capability not found

If a capability doesn't exist:

1. Check spelling matches exactly
2. Verify file exists: `ls spectr/specs/<capability>/spec.md`
3. Ensure directory structure is correct

### Naming conflicts

If change IDs conflict:

1. List all changes: `spectr list`
2. List archived changes: `ls spectr/changes/archive/`
3. Use unique ID or add suffix: `add-feature-2`

## Advanced Configuration

### Multi-Repository Setup

For monorepos or multiple Spectr instances:

```bash
# Initialize multiple Spectr projects
spectr init ./services/auth
spectr init ./services/payments
spectr init ./services/webhooks
```

Each has independent:
- `project.md` - Own conventions
- `specs/` - Own capabilities
- `changes/` - Own proposals

### Custom Validation Rules

Extend validation in `project.md`:

```markdown
## Validation Rules
- Scenarios must include both WHEN and THEN
- Requirements must reference related issues
- Change proposals require team approval
- All tasks must include effort estimates
```

## Configuration Examples

### Minimal Project

```markdown
# Project: MyApp

## Capabilities
- auth
- api
- frontend

## Change Naming
Format: `add-`, `update-`, `remove-` prefix with kebab-case name
```

### Comprehensive Project

```markdown
# Project: Enterprise Platform

## Capabilities by Domain
### Auth Domain
- user-authentication
- multi-factor-auth
- session-management

### API Domain
- api-versioning
- rate-limiting
- oauth-integration

## Validation Rules
- Every requirement: 1+ scenario
- Every scenario: WHEN + THEN format
- Change scope: 1-3 capabilities
- Tasks: Estimated hours required

## Naming Conventions
- Changes: [verb]-[noun-phrase] (20-50 chars)
- Capabilities: [domain]-[feature]
- Requirements: Action-oriented verbs

## Approval Process
- Changes: Require technical review
- Specs: Require architecture team sign-off
- Archives: Require deployment confirmation
```

## Further Reading

- [Spec-Driven Development](/concepts/spec-driven-development)
- [Creating Changes](/guides/creating-changes)
- [CLI Commands](/reference/cli-commands)
