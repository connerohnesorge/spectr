## Context

AI agents (Claude Code, Cursor, etc.) read `tasks.jsonc` files to understand and track implementation progress. The current flat structure works well for small changes but becomes problematic when:

1. A change has 50+ tasks across multiple capabilities
2. Agents have limited context windows (~4-8k tokens for Read operations)
3. Tasks logically belong to specific delta specs but can't be co-located

The `redesign-provider-architecture` change exemplifies this: 60+ tasks, 9 sections, 17 capability-specific delta specs—all tasks stuffed into one file.

## Goals / Non-Goals

**Goals:**

- Enable delta specs to have their own `tasks.jsonc` files
- Provide summary view in root file for quick status overview
- Auto-generate hierarchical structure from `tasks.md` during accept
- Maintain backwards compatibility with flat `tasks.jsonc` files
- Keep the format simple enough for agents to understand and edit

**Non-Goals:**

- Deeply nested hierarchies (max 2 levels: root → capability)
- Complex dependency graphs between tasks
- Task file splitting by arbitrary criteria (only by capability)

## Decisions

### Decision 1: File Structure

Root `tasks.jsonc` contains summary and references:

```jsonc
{
  "version": 2,
  "summary": {
    "total": 60,
    "completed": 12,
    "in_progress": 1,
    "pending": 47
  },
  "tasks": [
    {
      "id": "1",
      "section": "Foundation",
      "description": "Create core interfaces",
      "status": "completed"
    },
    {
      "id": "5",
      "section": "Migrate Providers",
      "description": "Migrate all providers to new interface",
      "status": "in_progress",
      "children": "$ref:specs/support-aider/tasks.jsonc"
    }
  ],
  "includes": [
    "specs/*/tasks.jsonc"
  ]
}
```

Delta spec `specs/support-aider/tasks.jsonc`:

```jsonc
{
  "version": 2,
  "parent": "5",
  "tasks": [
    {
      "id": "5.1",
      "description": "Migrate aider.go to new Provider interface",
      "status": "pending"
    },
    {
      "id": "5.2",
      "description": "Add unit tests for Aider provider",
      "status": "pending"
    }
  ]
}
```

**Rationale:**

- `version: 2` distinguishes new format from legacy `version: 1`
- `summary` enables agents to see progress without reading all files
- `children` syntax is explicit about the relationship
- `includes` glob pattern for auto-discovery
- `parent` in child files allows standalone reading

### Decision 2: Auto-Split Logic in `spectr accept`

When running `spectr accept <change-id>`:

1. Parse `tasks.md` into sections
2. For each section, check if a matching delta spec exists: `specs/<section-kebab>/spec.md`
3. If match found, split those tasks into `specs/<section-kebab>/tasks.jsonc`
4. Root `tasks.jsonc` gets a reference task with `children` pointing to the split file
5. Sections without matching delta specs remain in root as flat tasks

**Matching rules:**

- Section name "5. Migrate Providers" → check `specs/migrate-providers/`
- Section name "Support Aider" → check `specs/support-aider/`
- Exact match required (kebab-case normalized)

### Decision 3: Task ID Schema

Hierarchical IDs follow dot notation:

- Root task: `"5"`
- Child task: `"5.1"`, `"5.2"`
- Nested (if needed): `"5.1.1"`

The parent task ID becomes the prefix for all children. This maintains compatibility with existing ID patterns while enabling hierarchy.

### Decision 4: Resolution Order

When multiple sources exist:

1. Explicit `children` references are resolved first
2. Then `includes` glob patterns are processed
3. Duplicate detection by file path (same file = processed once)

This allows explicit ordering control via `children` while defaulting to glob discovery.

### Decision 5: `spectr tasks` Command

New command for viewing tasks:

```bash
# Show summary (default)
spectr tasks <change-id>
# Output:
# Foundation: 4/4 completed
# Built-in Initializers: 2/6 completed
# Migrate Providers: 0/17 pending
# ...
# Total: 12/60 completed (20%)

# Flatten all tasks into single view
spectr tasks <change-id> --flatten
# Output: Full flat list with hierarchical IDs

# JSON output for tooling
spectr tasks <change-id> --json
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Agents confused by new format | Version field enables detection; agents can handle both |
| Orphaned child files | Validation checks that referenced files exist |
| Complex status aggregation | Summary computed lazily; cache in root on write |
| Breaking existing workflows | Version 1 format remains fully supported |

## Migration Plan

1. **Phase 1**: Add support for reading version 2 format (backwards compatible)
2. **Phase 2**: Update `spectr accept` to generate version 2 when conditions met
3. **Phase 3**: Add `spectr tasks` command for viewing
4. **Phase 4**: Update AGENTS.md documentation

Rollback: None needed—version 1 files continue to work unchanged.

## Open Questions

1. Should child task files include the JSONC header comments, or keep them minimal?
   - **Proposed**: Include abbreviated header (status values only, not full workflow)

2. Should summary be auto-updated when child files change?
   - **Proposed**: Compute on read (no caching complexity)

3. Maximum nesting depth?
   - **Proposed**: 2 levels (root → capability) for simplicity
