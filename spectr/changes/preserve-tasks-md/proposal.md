# Proposal: Preserve tasks.md During Accept

## Problem

Currently, the `spectr accept` command converts `tasks.md` → `tasks.jsonc` and
then **deletes** the original `tasks.md` file. This creates several problems:

1. **Information Loss**: Markdown formatting (backticks, bold, italic, links) in
  task descriptions is stripped during conversion
2. **No Source of Truth**: Once tasks.md is deleted, there's no way to recover
  the original human-readable format
3. **Comments Lost**: Any explanatory comments or notes in tasks.md disappear
4. **Workflow Friction**: Users cannot easily revert from tasks.jsonc back to
  tasks.md if needed
5. **Manual Editing Harder**: Editing tasks.jsonc is less ergonomic than editing
  markdown

### Current Behavior

```bash
# Before accept
spectr/changes/my-change/
├── proposal.md
└── tasks.md

# After accept
spectr/changes/my-change/
├── proposal.md
└── tasks.jsonc  # tasks.md is deleted!
```

### Example Information Loss

**Original tasks.md:**

```markdown
## 1. Implementation

<!-- This section handles the core accept logic -->
- [ ] 1.1 Update `AcceptCmd` struct in `cmd/accept.go`
- [ ] 1.2 Add **validation** for `tasks.jsonc` format
- [ ] 1.3 See [issue #123](https://github.com/example/issue/123)
```

**Converted tasks.jsonc (loses formatting):**

```jsonc
{
  "tasks": [
    {
      "id": "1.1",
      "description": "Update AcceptCmd struct in cmd/accept.go",  // Lost backticks
      "status": "pending"
    },
    {
      "id": "1.2",
      "description": "Add validation for tasks.jsonc format",  // Lost bold
      "status": "pending"
    },
    {
      "id": "1.3",
      "description": "See issue #123",  // Lost link
      "status": "pending"
    }
  ]
}
// Lost comment about section purpose
```

## Solution

**Keep both `tasks.md` and `tasks.jsonc` after accept.**

### Proposed Behavior

```bash
# After accept
spectr/changes/my-change/
├── proposal.md
├── tasks.md      # Preserved!
└── tasks.jsonc   # Generated
```

### Benefits

1. **No Information Loss**: Original markdown formatting preserved
2. **Human-Readable Source**: tasks.md remains the canonical human-readable
  version
3. **Easy Editing**: Users can edit tasks.md and re-run accept if needed
4. **Backward Compatible**: Existing code already prefers tasks.jsonc when both
  exist
5. **No New Commands**: No need for `unaccept` subcommand

### Implementation Strategy

1. Modify `cmd/accept.go` to NOT delete `tasks.md` after conversion
2. Update documentation to clarify that both files will exist
3. Add validation to warn if tasks.md and tasks.jsonc diverge
4. Ensure `spectr list` and `spectr view` prefer tasks.jsonc when both exist

### Trade-offs

**Pros:**

- Simple implementation (remove file deletion)
- Preserves all original information
- Flexible workflow (users can edit either file)

**Cons:**

- Two files to maintain (could diverge)
- Slightly more complex mental model (which file is source of truth?)

**Mitigation:**

- Document that tasks.jsonc is the runtime source of truth
- Add validation warning if files diverge
- Consider `spectr sync-tasks` command in future to sync tasks.md ← tasks.jsonc

## Alternatives Considered

### Alternative 1: Add `unaccept` Subcommand

**Rejected**: More complex, still doesn't preserve original formatting

### Alternative 2: Enhance tasks.jsonc with Metadata

**Rejected**: Too complex, adds significant schema changes

### Alternative 3: Best-Effort Reconstruction

**Rejected**: Cannot recover lost formatting (backticks, links, comments)

## Success Criteria

1. After `spectr accept`, both tasks.md and tasks.jsonc exist
2. Existing code continues to use tasks.jsonc as primary source
3. All tests pass
4. Documentation updated to reflect new behavior
5. No breaking changes to existing workflows

## Related Changes

- None (standalone change)

## Migration Path

**Existing changes**: No action needed. Changes that already have tasks.jsonc
(without tasks.md) continue to work.

**New changes**: After this change, `spectr accept` will preserve tasks.md
automatically.
