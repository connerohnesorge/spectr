# Spectr Package Knowledge Base

Spec-driven development workflow for Spectr CLI itself. Instructions for
managing specs, changes, and proposals.

## OVERVIEW

Spectr enforces propose → validate → archive workflow. `spectr/specs/` is
current truth, `spectr/changes/` are proposed deltas. Archive merges deltas
preserving history.

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
| Create proposals | spectr/changes/ | Scaffold with proposal.md, tasks.md |
| Validate changes | spectr validate | Enforce scenarios, formatting rules |
| Accept proposals | spectr accept | Convert tasks.md → tasks.jsonc |
| Execute next task | /spectr:next | Auto-execute next pending task |
| Archive changes | spectr archive | Merge deltas into specs/ |
| View status | spectr view | Interactive dashboard |

## CONVENTIONS

- **Verb-led IDs**: Use `add-`, `update-`, `remove-`, `refactor-` prefixes

- **Delta operations**: Use `## ADDED`, `## MODIFIED`, `## REMOVED`,
  `## RENAMED Requirements`

- **Scenario format**: Use `#### Scenario:` (4 hashtags) with WHEN/THEN bullets

- **tasks.md + tasks.jsonc**: Both coexist, update status in .jsonc after accept

## UNIQUE PATTERNS

- **Spec-driven development**: All features tracked in specs/ before implementation

- **Three-stage workflow**: Propose → Validate (pre-implementation) → Archive (post-deployment)

- **Dual task files**: tasks.md (human-readable) + tasks.jsonc (machine-readable)

## ANTI-PATTERNS

- **NEVER implement without proposal**: Features need change in spectr/changes/

- **DON'T skip validation**: Always run `spectr validate` before `spectr accept`

- **NO partial MODIFIED**: MODIFIED requirements must include complete content
  (header + all scenarios)

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

## CHAINED PROPOSALS

Proposals can declare dependencies on other proposals using YAML frontmatter:

```yaml
---
id: feat-dashboard
requires:
  - id: feat-auth
    reason: "needs user model and session management"
  - id: feat-db
    reason: "needs schema migrations for dashboard tables"
enables:
  - id: feat-analytics
    reason: "unlocks event tracking on dashboard"
---

```

**Fields:**

- `id`: Optional explicit ID for this proposal

- `requires`: List of proposals that must be archived before this can be accepted

- `enables`: Informational list of proposals this unlocks (not enforced)

**Behavior:**

- `spectr validate`: Warns if required proposals aren't archived, errors on cycles

- `spectr accept`: Hard fails if any `requires` entries aren't archived

- `spectr graph`: Visualizes the proposal dependency DAG

**Commands:**

```bash

# View dependency graph (ASCII)
spectr graph

# View graph for specific proposal
spectr graph <change-id>

# Output in Graphviz DOT format
spectr graph --dot

# Output in JSON format
spectr graph --json

```

**Example Graph Output:**

```text
feat-dashboard (⧖)
├── requires: feat-auth ✓
├── requires: feat-db ⧖
└── enables: feat-analytics

```

Legend: ✓ = archived, ⧖ = active/pending

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

### CLI Commands

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

### AI Slash Commands

```bash

# Execute next pending task (AI agents only)
/spectr:next [change-id]

```

**How /spectr:next works:**

1. Discovers change proposal directory in `spectr/changes/`
2. Parses `tasks.jsonc` to find first task with `status: "pending"`
3. Executes the task based on its description
4. Updates task status: `pending` → `in_progress` → `completed`
5. Reports progress and suggests next steps

**Example workflow:**

```text
User: /spectr:next add-feature-x

AI Agent:
→ Found next pending task: #3.2 "Implement validation logic"
→ Marked task #3.2 as in_progress
→ Implementing validation in internal/validation/validator.go
→ Updated task status to completed
→ Next task: #3.3 "Add unit tests for validation"

```

**Best practices:**

- Always verify the task description before execution

- Update status immediately when starting work

- Report progress clearly and suggest next steps

- Use with `/spectr:apply` for full change lifecycle

## VALIDATION RULES

| Rule | Severity | Description |
|------|----------|-------------|
| RequirementScenarios | Error | Every requirement MUST have ≥1 scenario |
| ScenarioFormatting | Error | Scenarios MUST use `#### Scenario:` (4 hashtags) |
| PurposeLength | Warning | Purpose sections MUST be ≥50 chars |
| ModifiedComplete | Error | MODIFIED requirements include full content |
| DeltaPresence | Error | Changes MUST have ≥1 delta spec |

## MULTI-REPO DISCOVERY

Spectr supports mono-repo setups with nested git repositories, each with their
own `spectr/` directory.

### Discovery Behavior

- **Automatic discovery**: Spectr walks up from the current working directory to
  find all `spectr/` directories, stopping at `.git` boundaries
- **Git isolation**: Each git repository is isolated; discovery stops at the git
  root
- **Aggregated results**: Commands like `list`, `validate`, `view` aggregate
  results from all discovered roots
- **Root prefix**: In multi-root scenarios, items are prefixed with their
  relative path: `[../project] add-feature`

### SPECTR_ROOT Environment Variable

Override automatic discovery by setting `SPECTR_ROOT`:

```bash
# Use explicit spectr root
SPECTR_ROOT=/path/to/project spectr list

# Relative paths work too
SPECTR_ROOT=../other-project spectr validate --all
```

**Behavior:**

- When set, uses ONLY the specified root (skips automatic discovery)
- Errors if the path doesn't contain a `spectr/` directory

### Clipboard Copy Paths

When selecting items in interactive mode (Enter key), Spectr copies the full
path relative to your cwd:

- Single root: `spectr/changes/add-feature/proposal.md`
- Nested root: `../project/spectr/changes/add-feature/proposal.md`

This enables direct navigation with `@` file references in AI tools.

## NOTES

- **Validation is strict**: All issues treated as errors (no warnings in strict mode)

- **Archive is atomic**: Spec merge + move happens together or not at all

- **tasks.md vs tasks.jsonc**: Both coexist after accept. Update statuses in
  .jsonc, regenerate from .md on structure changes

- **Git worktrees**: PR commands create isolated worktrees, cleanup
  automatically

- **Multi-platform PRs**: Supports GitHub (gh), GitLab (glab), Gitea/Forgejo
  (tea), Bitbucket (manual)

- **Multi-repo support**: Aggregates from all discovered `spectr/` directories
  within the git boundary
