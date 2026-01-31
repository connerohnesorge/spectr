# Design: Size-Based Task Splitting

## Context

AI agents (Claude Code, Cursor, Aider, etc.) have practical limits on Read
operations:

- Typical context windows allow ~100-150 lines per read
- Larger files require offset/limit parameters, losing context
- Reading files in chunks forces agents to make multiple reads and mentally
  stitch results

The current flat `tasks.jsonc` format works well for small changes
(10-30 tasks,
<100 lines) but breaks down for large changes. For example:

- A change with 60 tasks across 9 sections → 180+ line tasks.jsonc
- Agents hit truncation limits or must request multiple reads
- Loss of context between reads leads to missing tasks or incomplete
  understanding

The previous `add-hierarchical-tasks` proposal attempted to solve this by
tying splitting to delta spec directories. However, this created unnecessary
coupling and complexity. The new approach is simpler: **split purely based on
file size, using natural section boundaries**.

## Goals / Non-Goals

**Goals:**

- Enable automatic splitting of large tasks.md files (>100 lines) during accept
- Use section boundaries as natural split points (preserve sections together)
- Support subsection splitting when a single section exceeds 100 lines
- Generate hierarchical structure with root + child files
- Maintain backwards compatibility with flat version 1 files
- Keep implementation simple: hardcoded 100-line threshold, no configuration
- Preserve human-friendly single-file authoring (tasks.md remains unchanged)

**Non-Goals:**

- Delta spec directory integration (keep splitting logic independent)
- Deeply nested hierarchies (max 2 levels: root → child)
- Complex dependency graphs between tasks
- Configurable thresholds (hardcode 100 lines for simplicity)
- Automatic sync (users run `spectr accept` again to regenerate)

## Decisions

### Decision 1: 100-Line Threshold

**Choice:** Hardcode 100 lines as the splitting threshold.

**Rationale:**

- 100 lines fits comfortably in most AI agent Read operations
- Provides headroom for JSONC comments and formatting
- Aligns with typical terminal viewport height (~50-80 lines visible)
- Simple to implement and reason about (no configuration complexity)
- Can be adjusted later if needed, but start with proven value

**Trade-offs:**

- Less flexible than configurable threshold
- May split some files that could fit in a single read
- But: simplicity wins over premature optimization

### Decision 2: Section-Based Splitting

**Choice:** Split on top-level section boundaries (lines starting with `## N.`),
then subsections if needed when files exceed limits.

**Splitting algorithm:**

1. Count total lines in tasks.md
2. If ≤100 lines: generate version 1 flat file (no splitting)
3. If >100 lines: enter splitting mode
   - Parse tasks.md into sections (by `## N. Section Name`)
   - For each section:
     - If section <100 lines: candidate for splitting
     - If section >100 lines: split by subsections (task ID prefixes
       like "1.1", "1.2")
   - Generate root tasks.jsonc with reference tasks
   - Generate tasks-N.jsonc files for each split section/subsection
     group

**Rationale:**

- Sections are natural conceptual boundaries (Implementation, Testing,
  Documentation)
- Preserves related tasks together (better for agent understanding)
- Subsection splitting handles edge case of single large section
- No artificial chunking that breaks logical groupings

**Example:**

tasks.md (160 lines):

```markdown
## 1. Foundation (40 lines)
- [ ] 1.1 Task...
- [ ] 1.2 Task...
...

## 2. Implementation (80 lines)
- [ ] 2.1 Subsection A...
- [ ] 2.2 Subsection A...
- [ ] 2.3 Subsection B...
- [ ] 2.4 Subsection B...
...

## 3. Testing (30 lines)
- [ ] 3.1 Task...
...
```

Generated files:

- `tasks.jsonc` (root: ~30 lines with 3 reference tasks)
- `tasks-1.jsonc` (Foundation: 40 lines)
- `tasks-2.jsonc` (Implementation subsection A: 40 lines)
- `tasks-3.jsonc` (Implementation subsection B: 40 lines)
- `tasks-4.jsonc` (Testing: 30 lines)

### Decision 3: File Naming Convention

**Choice:** Use `tasks-N.jsonc` suffixes where N is sequential (1, 2, 3...).

**Rationale:**

- Simple, predictable naming
- Avoids filename conflicts
- Sequential numbering matches section order from tasks.md
- Easy for agents to discover and read in order
- Glob pattern `tasks-*.jsonc` catches all files

**Alternatives considered:**

- `tasks-<section-name>.jsonc` → too long, requires normalization, conflicts
- `tasks/<N>.jsonc` → extra directory, complicates file discovery
- `tasks.part<N>.jsonc` → uglier, less clear

### Decision 4: Version 2 Schema

**Choice:** Introduce `version: 2` for hierarchical files with new fields.

**New fields:**

Root tasks.jsonc:

```jsonc
{
  "version": 2,
  "tasks": [
    {
      "id": "1",
      "section": "Foundation",
      "description": "Implement core features",
      "status": "pending",
      "children": "$ref:tasks-1.jsonc"  // NEW
    }
  ],
  "includes": ["tasks-*.jsonc"]  // NEW
}
```

Child tasks-N.jsonc:

```jsonc
{
  "version": 2,
  "parent": "1",  // NEW
  "tasks": [
    {
      "id": "1.1",  // Hierarchical ID
      "section": "Foundation",
      "description": "Create database schema",
      "status": "pending"
    }
  ]
}
```

**Rationale:**

- `version` field enables detection and backwards compatibility
- `children` field makes parent-child relationship explicit
- `$ref:` prefix signals file reference (extensible for future use)
- `parent` field in child enables standalone reading
- `includes` glob enables auto-discovery by tooling

### Decision 5: Task ID Schema

**Choice:** Use dot notation for hierarchical IDs (e.g., "5", "5.1", "5.1.1").

**Rules:**

- Root tasks: single number or decimal (e.g., "1", "2", "5")
- Child tasks: parent ID + dot + child number (e.g., "5.1", "5.2")
- Nested children: additional dots (e.g., "5.1.1", "5.1.2")

**Rationale:**

- Standard hierarchical notation (familiar from outlines, TOCs)
- Preserves existing task IDs from tasks.md
- Parent ID becomes prefix for children (clear relationship)
- Enables quick filtering/sorting by ID prefix

### Decision 6: Status Aggregation

**Choice:** Compute parent reference task status from children.

**Aggregation rules:**

- All children "pending" → parent "pending"
- Any child "in_progress" → parent "in_progress"
- All children "completed" → parent "completed"
- Mixed "pending" + "completed" (no "in_progress") → parent "in_progress"

**Rationale:**

- Provides quick status overview in root file
- "in_progress" acts as catch-all for partial completion
- Agents can see progress without reading all child files
- Status computed on write (during accept), not on read

### Decision 7: Regeneration from tasks.md

**Choice:** Re-running `spectr accept` regenerates all tasks.jsonc files from
tasks.md.

**Behavior:**

1. Delete existing `tasks-*.jsonc` files
2. Parse tasks.md from scratch
3. Apply splitting logic
4. Generate new root + child files
5. Preserve task statuses from old files where IDs match

**Rationale:**

- tasks.md remains the source of truth for structure
- Avoids complex sync logic and partial updates
- Idempotent operation (safe to run multiple times)
- Matches existing accept behavior (regenerate, not update)
- Status preservation maintains work progress

### Decision 8: Child File Headers

**Choice:** Include full header with origin information in child files.

**Header format:**

```jsonc
// Generated by: spectr accept add-size-based-task-splitting
// Parent change: add-size-based-task-splitting
// Parent task: 5
// Status values: "pending" | "in_progress" | "completed"
// Valid transitions: pending -> in_progress -> completed
// [rest of standard header]
```

**Rationale:**

- Provides context when agent reads child file standalone
- Documents origin for debugging and understanding
- Full header ensures child files are self-documenting
- Parent task ID enables quick lookup back to root file

## Implementation Plan

### Phase 1: Schema Updates

1. Update `internal/parsers/types.go`:
   - Add `Children` field to Task struct: `Children string
     json:"children,omitempty"`
   - Add `Parent` field to TasksFile struct: `Parent string
     json:"parent,omitempty"`
   - Add `Includes` field to TasksFile struct: `Includes []string
     json:"includes,omitempty"`

2. Update version constant to support version 2

### Phase 2: Splitting Logic

1. In `cmd/accept.go`:
   - Add `detectSplitting()` function: count lines in tasks.md, return bool
   - Add `parseSections()` function: extract sections with line ranges
   - Add `shouldSplitSection()` function: check if section >100 lines
   - Add `splitBySubsections()` function: group tasks by ID prefixes

2. In `cmd/accept_writer.go`:
   - Add `writeHierarchicalTasks()` function: orchestrate root + child file
     generation
   - Add `writeChildTasksJSONC()` function: write tasks-N.jsonc with full header
   - Add `computeAggregateStatus()` function: implement status aggregation rules
   - Update existing `writeTasksJSONC()` to call hierarchical path when needed

### Phase 3: Task ID Generation

1. Add `assignHierarchicalIDs()` function:
   - Parse existing IDs from tasks.md
   - For child tasks, prepend parent ID with dot notation
   - Validate ID uniqueness

### Phase 4: Status Preservation

1. Add `loadExistingStatuses()` function:
   - Read existing tasks.jsonc (root + all children)
   - Build map of ID → status
   - Return for use during regeneration

2. Update `writeTasksJSONC()` to merge statuses

### Phase 5: Testing

1. Add unit tests for splitting logic
2. Add integration test with 150-line tasks.md
3. Test status preservation across regeneration
4. Test backwards compatibility with version 1 files

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Agents confused by new format | Version field enables detection |
| Loss of context across files | Keep root minimal (<40 lines) |
| Status desync between files | tasks.jsonc is source of truth |
| Hardcoded 100-line threshold | Start with 100, adjust if needed |
| Breaking existing workflows | Version 1 fully supported |
| Complexity in splitting logic | Extensive testing; edge cases handled |

## Open Questions

None - all questions resolved during user interview.

## Success Metrics

- Agents can read any tasks.jsonc file in a single Read operation
- Root file stays under 40 lines for quick overview
- Child files stay under 100 lines (typically 40-60 lines)
- Zero breaking changes to existing flat files
- Implementation completed in <500 lines of new code
