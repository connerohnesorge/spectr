# Archive Command Design

## Context

The archive command completes the Spectr change lifecycle by moving deployed changes to a dated archive and applying their delta specifications to the main spec files. This is a critical workflow step that ensures specs stay synchronized with implemented changes.

### Constraints

- Must maintain backward compatibility with existing spec file structure
- Must preserve requirement ordering where possible
- Must validate rigorously before modifying main specs (no partial updates)
- Must support both interactive and non-interactive (CI/CD) usage
- Must match OpenSpec's proven archive command semantics closely

### Stakeholders

- Developers using Spectr for spec-driven development
- CI/CD pipelines that need automated archiving
- Teams maintaining specification consistency

## Goals / Non-Goals

### Goals

1. Provide a complete, safe archive workflow
2. Automate delta spec application to prevent manual errors
3. Validate changes thoroughly before archiving
4. Support both interactive and automated usage
5. Maintain audit trail with dated archives
6. Match OpenSpec archive command behavior

### Non-Goals

- Archive command does NOT deploy changes (deployment happens separately)
- Archive command does NOT modify implementation code
- Archive command does NOT create git commits (user does this separately)
- No automatic rollback of archives (manual git revert if needed)

## Decisions

### Decision 1: Delta Operation Ordering

**Choice**: Apply delta operations in strict order: RENAMED → REMOVED → MODIFIED → ADDED

**Rationale**:

- RENAMED must come first to update requirement names before other operations reference them
- REMOVED must come before MODIFIED/ADDED to prevent conflicts
- MODIFIED must come before ADDED to ensure existing requirements are updated first
- This order matches OpenSpec's implementation and is proven to work

**Alternatives considered**:

- Topological sort based on dependencies: Too complex, operation order is deterministic
- User-specified order: Error-prone, strict order is safer

### Decision 2: Requirement Name Normalization

**Choice**: Normalize requirement names by trimming whitespace and using case-insensitive matching

**Rationale**:

- Prevents duplicate requirements due to minor formatting differences
- Matches OpenSpec behavior
- Users often have slight variations in whitespace/casing

**Implementation**:

```go
func normalizeRequirementName(name string) string {
    return strings.ToLower(strings.TrimSpace(name))
}
```

### Decision 3: Validation Strategy (Multi-Stage)

**Choice**: Validate at three stages:

1. Pre-merge: Validate delta operations and source requirements exist
2. During merge: Detect duplicates and conflicts
3. Post-merge: Validate rebuilt spec structure and scenarios

**Rationale**:

- Fail fast on obvious errors (pre-merge)
- Prevent invalid states during merge (during merge)
- Ensure output is valid (post-merge)
- Never write invalid specs to main specs directory

**Alternatives considered**:

- Single post-merge validation: Misses opportunities to fail fast with better error messages
- No validation: Unacceptable risk of corrupting main specs

### Decision 4: Atomic Spec Updates

**Choice**: Prepare all spec updates first (validation pass), then write all at once

**Rationale**:

- Prevents partial updates if later specs fail validation
- All-or-nothing guarantee for spec consistency
- Matches OpenSpec approach

**Implementation flow**:

```
1. Find all delta specs
2. For each delta spec:
   a. Load base spec (or create skeleton)
   b. Parse delta operations
   c. Apply operations to build updated spec
   d. Validate rebuilt spec
   e. Store in memory (don't write yet)
3. If any validation failed, abort with no writes
4. If all passed, write all specs
```

### Decision 5: New Spec Creation

**Choice**: When archiving creates a new spec (no existing spec.md), only ADDED operations are allowed

**Rationale**:

- MODIFIED/REMOVED/RENAMED require an existing spec with existing requirements
- Prevents confusing errors when users try to modify non-existent requirements
- Clear error message guides users to correct usage

**Error message**:

```
<spec-name>: target spec does not exist; only ADDED requirements are allowed for new specs.
```

### Decision 6: Flag Design

**Flags**:

- `-y, --yes`: Skip all confirmation prompts (for CI/CD)
- `--skip-specs`: Skip spec update operations (for tooling/doc-only changes)
- `--no-validate`: Skip validation (requires confirmation, not recommended)

**Rationale**:

- `--yes`: Common pattern in Unix tools, needed for automation
- `--skip-specs`: Some changes (tooling, CI config) don't modify specs
- `--no-validate`: Emergency escape hatch, but discouraged via confirmation

### Decision 7: Archive Naming Format

**Choice**: `YYYY-MM-DD-<change-id>` (e.g., `2025-11-17-add-archive-command`)

**Rationale**:

- Chronological sorting by date
- Preserves original change ID for reference
- ISO 8601 date format is universal
- Matches OpenSpec format

### Decision 8: Task Completion Checking

**Choice**: Warn on incomplete tasks but allow archiving with confirmation

**Rationale**:

- Sometimes tasks are marked for future work or cancelled
- Blocking on incomplete tasks is too strict
- Warning + confirmation balances safety and flexibility
- `--yes` flag allows automation to proceed

## Risks / Trade-offs

### Risk 1: Spec Corruption

**Risk**: Bug in delta merging could corrupt main specs

**Mitigation**:

1. Comprehensive unit tests for all delta operations
2. Multi-stage validation prevents invalid output
3. Atomic updates (all-or-nothing)
4. Users should commit main specs to git before archiving
5. Integration tests with real-world delta scenarios

**Accepted trade-off**:

- Could add `--dry-run` flag in future if needed
- Current mitigation is sufficient for v1

### Risk 2: Complex Merge Conflicts

**Risk**: Users might create conflicting delta operations across sections

**Mitigation**:

1. Pre-merge validation detects cross-section conflicts
2. Clear error messages guide users to fix conflicts
3. Spectr AGENTS.md documentation provides examples

**Accepted trade-off**:

- No automatic conflict resolution (too complex, error-prone)
- Users must resolve conflicts manually

### Risk 3: Archive Name Collision

**Risk**: Archiving the same change twice on the same day creates collision

**Mitigation**:

1. Check if archive already exists before moving
2. Error if collision detected
3. Users can manually rename/remove old archive if needed

**Accepted trade-off**:

- No automatic suffix numbering (rare edge case)
- Clear error message is sufficient

### Risk 4: Partial Task Completion

**Risk**: Users might archive changes with incomplete implementation

**Mitigation**:

1. Parse tasks.md and count incomplete tasks
2. Warn user and require confirmation
3. Display task status before archiving

**Accepted trade-off**:

- Doesn't prevent archiving incomplete work (user decision)
- Warning is informational, not blocking

## Migration Plan

### Phase 1: Implementation (this change)

1. Implement archive command per tasks.md
2. Add comprehensive tests
3. Manual testing with sample changes
4. Archive this change itself as integration test

### Phase 2: Documentation

1. Update AGENTS.md with archive command usage
2. Add examples of delta spec patterns
3. Document flag usage and best practices

### Phase 3: Rollout

1. Archive existing active changes (if any are complete)
2. Monitor for issues in first few archives
3. Gather user feedback

### Rollback Plan

If critical bug found:

1. Restore main specs from git history
2. Restore archived change to active changes (reverse date prefix)
3. Fix bug and re-archive

## Implementation Notes

### Package Structure

```
cmd/
  archive.go          # CLI command implementation

internal/archive/
  archiver.go         # Main archive orchestration
  spec_merger.go      # Delta spec parsing and merging logic
  validator.go        # Pre/post merge validation
  types.go           # Data structures (SpecUpdate, DeltaPlan, etc.)

internal/parsers/
  requirement_parser.go  # Parse requirement blocks
  delta_parser.go        # Parse delta sections
```

### Key Data Structures

```go
// RequirementBlock represents a requirement with its header and content
type RequirementBlock struct {
    HeaderLine string  // "### Requirement: <name>"
    Name       string  // Extracted requirement name
    Raw        string  // Full block content (header + scenarios)
}

// DeltaPlan represents all delta operations for a spec
type DeltaPlan struct {
    Added    []RequirementBlock
    Modified []RequirementBlock
    Removed  []string // Just names
    Renamed  []RenameOp
}

// RenameOp represents a requirement rename
type RenameOp struct {
    From string
    To   string
}

// SpecUpdate represents a spec file to update
type SpecUpdate struct {
    Source string // Path to delta spec in change
    Target string // Path to main spec
    Exists bool   // Does target already exist?
}
```

### Error Handling Patterns

1. Validation errors: Collect all errors, display together, then abort
2. File I/O errors: Fail immediately with context (file path, operation)
3. Parse errors: Include file path, line number, and what was expected

### Testing Strategy

1. **Unit tests**: Each parser, each delta operation, each validation rule
2. **Integration tests**: Full archive workflow with realistic delta specs
3. **Edge cases**: Empty specs, new specs, large specs, all delta types combined
4. **Error cases**: Duplicates, conflicts, missing requirements, invalid scenarios

## Open Questions

1. **Q**: Should we add a `--dry-run` flag to preview spec updates without applying?
   - **A**: Defer to future enhancement. Current validation is sufficient for v1.

2. **Q**: Should validation be strict by default (block on warnings)?
   - **A**: Yes, match OpenSpec behavior. Use `--no-validate` as escape hatch.

3. **Q**: Should we support archiving multiple changes at once?
   - **A**: No, keep it simple. Archive one change at a time to reduce complexity.

4. **Q**: Should we create a git commit automatically after archiving?
   - **A**: No, git operations are user responsibility. Archive command focuses on spec management only.
