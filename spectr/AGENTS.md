# Spectr Package Knowledge Base

Spec-driven development workflow for Spectr CLI itself. Instructions for managing specs, changes, and proposals.

## OVERVIEW

Spectr enforces propose → validate → archive workflow. `spectr/specs/` is current truth, `spectr/changes/` are proposed deltas. Archive merges deltas preserving history.



## STRUCTURE

```text
spectr/
├── project.md           # Project-wide conventions
├── specs/               # Current truth - what IS built
│   └── [capability]/  # Single focused capability
│       ├── spec.md         # Requirements + scenarios
│       └── design.md       # Technical patterns (optional)
├── changes/              # Proposals - what SHOULD change
│   └── [change-id]/    # Verb-led, kebab-case IDs
│       ├── proposal.md     # Why, what, impact
│       ├── tasks.md        # Implementation checklist
│       ├── tasks.jsonc    # Machine-readable task status (post-accept)
│       ├── design.md       # Technical decisions (optional)
│       └── specs/          # Delta changes (ADDED/MODIFIED/REMOVED/RENAMED)
└── changes/archive/     # Completed changes with timestamps
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Create proposals | spectr/changes/ | Scaffold with proposal.md, tasks.md, delta specs |
| Validate changes | spectr validate | Enforce scenarios, formatting rules |
| Accept proposals | spectr accept | Convert tasks.md → tasks.jsonc |
| Archive changes | spectr archive | Merge deltas into specs/ |
| View status | spectr view | Interactive dashboard |

## CONVENTIONS

- **Verb-led IDs**: Use `add-`, `update-`, `remove-`, `refactor-` prefixes
- **Delta operations**: Use `## ADDED`, `## MODIFIED`, `## REMOVED`, `## RENAMED Requirements`
- **Scenario format**: Use `#### Scenario:` (4 hashtags) with WHEN/THEN bullets
- **tasks.md + tasks.jsonc**: Both coexist, update status in .jsonc after accept

## UNIQUE PATTERNS

- **Spec-driven development**: All features tracked in specs/ before implementation
- **Three-stage workflow**: Propose → Validate (pre-implementation) → Archive (post-deployment)
- **Dual task files**: tasks.md (human-readable) + tasks.jsonc (machine-readable)

## ANTI-PATTERNS

- **NEVER implement without proposal**: Features need change in spectr/changes/
- **DON'T skip validation**: Always run `spectr validate` before `spectr accept`
- **NO partial MODIFIED**: MODIFIED requirements must include complete content (header + all scenarios)
- **NO scenario-less requirements**: Every requirement must have at least one scenario

## DELTA SPEC FORMAT

```markdown
## ADDED Requirements
### Requirement: New Feature
The system SHALL...

#### Scenario: Success case
- **WHEN** condition
- **THEN** expected result

## MODIFIED Requirements
### Requirement: Existing Feature
[Complete updated requirement with all scenarios]
```

## CHANGE PROPOSAL STRUCTURE

```markdown
# Change: [Brief description]

## Why
[1-2 sentences on problem/opportunity]

## What Changes
- [Bullet list of changes]
- [Mark breaking changes with **BREAKING**]

## Impact
- Affected specs: [list capabilities]
- Affected code: [key files/systems]
```

## TASKS FILE FORMAT

```markdown
## 1. Implementation
- [ ] 1.1 Create database schema
- [ ] 1.2 Implement API endpoint
- [ ] 1.3 Add frontend component
- [ ] 1.4 Write tests
```

## tasks.jsonc FORMAT

```json
{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Implementation",
      "description": "Create database schema",
      "status": "pending"
    }
  ]
}
```

## COMMANDS

```bash
# Initialize
spectr init [path]

# List changes
spectr list
spectr list --specs

# Validate
spectr validate <change-id>

# Accept proposal
spectr accept <change-id>

# Archive change
spectr pr archive <change-id>

# View details
spectr view <change-id>

# Pull request
spectr pr new <change-id>
```

## VALIDATION RULES
| Rule | Severity | Description |
|------|----------|-------------|
| RequirementScenarios | Error | Every requirement MUST have ≥1 scenario |
| ScenarioFormatting | Error | Scenarios MUST use `#### Scenario:` (4 hashtags) |
| PurposeLength | Warning | Purpose sections MUST be ≥50 chars |
| ModifiedComplete | Error | MODIFIED requirements MUST include full updated content |
| DeltaPresence | Error | Changes MUST have ≥1 delta spec |

## NOTES

- **Validation is strict**: All issues treated as errors (no warnings in strict mode)
- **Archive is atomic**: Spec merge + move happens together or not at all
- **tasks.md vs tasks.jsonc**: Both coexist after accept. Update statuses in .jsonc, regenerate from .md on structure changes
- **Git worktrees**: PR commands create isolated worktrees, cleanup automatically
- **Multi-platform PRs**: Supports GitHub (gh), GitLab (glab), Gitea/Forgejo (tea), Bitbucket (manual)
