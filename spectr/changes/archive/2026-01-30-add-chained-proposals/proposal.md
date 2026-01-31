# Change: Add Chained Proposals with Explicit Dependencies

## Why

Large changes are unwieldy to review and implement as monolithic proposals.
Dependencies between proposals are currently implicit—there's no way to express
that proposal X requires proposal Y to be completed first. This leads to review
bottlenecks and confusion about implementation order.

Chained proposals allow breaking large changes into smaller, dependent pieces
with explicit ordering, making dependencies visible and enforceable.

## What Changes

- **ADDED**: YAML frontmatter support in `proposal.md` for declaring
  dependencies
  - `requires:` - list of change IDs that must be archived before this can be
    accepted
  - `enables:` - list of change IDs that this proposal unlocks (informational)
- **ADDED**: Dependency validation in `spectr validate` - warns if required
  proposals aren't archived
- **ADDED**: Hard dependency check in `spectr accept` - fails if any `requires`
  entries aren't in archive
- **ADDED**: New `spectr graph` command to visualize the proposal dependency DAG
- **ADDED**: Circular dependency detection with clear error messages
- **MODIFIED**: Proposal.md parsing to extract frontmatter metadata

## Syntax

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

## Impact

- **Affected specs**: `cli-interface` (new commands, modified accept/validate)
- **Affected code**:
  - `internal/domain/` - new ProposalMetadata type, frontmatter parsing
  - `internal/validation/change_rules.go` - dependency validation rules
  - `cmd/accept.go` - archive check before accepting
  - `cmd/validate.go` - warn on unmet dependencies
  - New `cmd/graph.go` - DAG visualization command
  - `internal/discovery/changes.go` - resolve archived change IDs
- **Breaking changes**: None - existing proposals without frontmatter remain
  valid
