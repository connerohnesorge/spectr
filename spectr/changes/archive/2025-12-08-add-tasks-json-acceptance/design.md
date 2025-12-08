## Context

AI agents working with Spectr proposals frequently need to update task status as they complete implementation work. The current `tasks.md` Markdown format is human-readable but presents challenges for agent stability:

1. **Parsing fragility**: Regex-based parsing of `- [ ]` and `- [x]` patterns can fail on edge cases
2. **Overwrite risk**: Agents may accidentally rewrite task descriptions when updating status
3. **Format drift**: Multiple agents or manual edits can introduce inconsistent formatting
4. **Silent corruption**: Markdown doesn't fail loudly when the structure is wrong

Anthropic's research on effective harnesses for long-running agents recommends structured formats like JSON for state that agents need to read and write repeatedly.

## Goals / Non-Goals

**Goals:**
- Provide a stable, machine-readable format for task tracking during implementation
- Prevent accidental task list corruption by agents
- Maintain human-readability during the proposal/review phase
- Ensure backward compatibility with existing changes using `tasks.md`

**Non-Goals:**
- Forcing JSON format for the initial proposal phase (Markdown is better for human review)
- Breaking existing tooling that reads `tasks.md`
- Automatic migration of all existing changes

## Decisions

### Decision 1: Two-phase task format

**What**: Use `tasks.md` during proposal creation/review, convert to `tasks.json` at acceptance time.

**Why**:
- Markdown is better for human review (GitHub diffs, PRs, comments)
- JSON is better for agent manipulation during implementation
- Clear transition point (acceptance) provides semantic meaning

**Alternatives considered**:
- JSON-only: Would hurt proposal readability and review experience
- Both formats simultaneously: Creates drift risk, violates single source of truth
- YAML: Slightly more readable than JSON but similar parsing complexity

### Decision 2: Remove tasks.md after conversion

**What**: The `accept` command removes `tasks.md` after successfully creating `tasks.json`.

**Why**:
- Prevents drift between two files
- Clear signal to agents that JSON is the authoritative source
- Archived changes retain `tasks.json` for historical record

**Alternatives considered**:
- Keep both: Creates confusion about which is authoritative
- Rename to `tasks.md.bak`: Clutters directory, agents might still find it

### Decision 3: Fallback to tasks.md for backward compatibility

**What**: All tooling checks for `tasks.json` first, falls back to `tasks.md`.

**Why**:
- Existing changes continue to work
- Gradual migration path
- No breaking changes for users who don't use `accept`

### Decision 4: Require accept before apply

**What**: The `apply` slash command instructions require agents to run `spectr accept` first.

**Why**:
- Ensures agents work with stable JSON format
- Validates change before implementation begins
- Creates explicit approval gate

## JSON Schema

```json
{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Implementation",
      "description": "Create database schema",
      "status": "pending"
    },
    {
      "id": "1.2",
      "section": "Implementation",
      "description": "Implement API endpoint",
      "status": "completed"
    }
  ]
}
```

**Fields:**
- `version`: Schema version for future compatibility
- `tasks[]`: Array of task objects
- `tasks[].id`: Original task ID (e.g., "1.1", "2.3")
- `tasks[].section`: Section header (e.g., "Implementation", "Testing")
- `tasks[].description`: Full task description text
- `tasks[].status`: One of `"pending"`, `"in_progress"`, `"completed"`

**Status values:**
- `pending`: Not started (default from `- [ ]`)
- `in_progress`: Currently being worked on (agents can set this)
- `completed`: Done (from `- [x]`)

Note: `in_progress` is a new status that JSON enables but Markdown couldn't easily represent.

## Risks / Trade-offs

### Risk: Agent confusion with status values
**Mitigation**: Clear documentation in AGENTS.md and slash command templates.

### Risk: Incomplete accept leaving orphan tasks.json
**Mitigation**: Accept is atomic - either fully succeeds or rolls back.

### Trade-off: Slightly more complex implementation
**Accepted**: The stability benefits outweigh the implementation cost.

## Migration Plan

1. **Phase 1**: Add `accept` command and JSON support to all tooling
2. **Phase 2**: Update slash commands to recommend `spectr accept` before implementation
3. **Phase 3**: (Future, not in this change) Consider auto-accept on first task update

## Open Questions

- Should `spectr archive` auto-accept if `tasks.md` still exists? (Proposed: Yes, with warning)
- Should there be a `spectr reject` to revert to `tasks.md`? (Proposed: No, just delete `tasks.json` and restore from git)
