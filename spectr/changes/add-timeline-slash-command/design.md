# Timeline Slash Command: Design Notes

## Architecture Overview

The `/spectr:timeline` feature is implemented as a skill (agent-executable
command) rather than a CLI command. This makes it accessible to AI agents
without requiring changes to the Go codebase and gives it access to filesystem
and Python tooling.

### Design Decisions

1. **Skill-based, not CLI**: Implemented as `.agents/skills/spectr-timeline/`
   rather than `cmd/timeline.go`. Advantages:
   - AI agents can execute directly
   - Can leverage external tools if needed
   - Doesn't require Go compilation or testing
   - Easier iteration and debugging

2. **Dependency analysis reuses existing patterns**: The chained-proposals
   feature already parses and validates dependencies. The timeline skill reads
   the same frontmatter format, so no new parsing is required—just different
   output.

3. **Phase-based output for planning**: Rather than just listing dependencies,
   compute implementation phases (batches of changes that can run in parallel).
   This makes the timeline immediately actionable for sprint planning.

4. **JSON output for machine readability**: Uses structured JSON to allow
   downstream tools to consume the timeline (dashboards, scripts, reports).
   Still human-friendly with good formatting and descriptive field names.

5. **Circular dependency detection**: Hard error if cycles are detected. Better
   to fail early than generate nonsensical timelines.

### Timeline Algorithm

```text
1. Discover all active changes
   - List directories in spectr/changes/ (not archive/)
   - Read proposal.md from each

2. Parse metadata
   - Extract requires/enables relationships
   - Build adjacency list (change -> depends on [changes])

3. Detect cycles
   - Topological sort or DFS coloring
   - Fail if cycle found

4. Calculate phases
   - Phase 0: changes with no dependencies
   - Phase N: changes whose dependencies are all in phases 0..N-1
   - Repeat until all changes assigned

5. Generate timeline.json
   - Include summary metadata
   - Include full dependency graph
   - Include phase-based timeline
   - Add implementation notes

6. Output location
   - Write to ./spectr/timeline.json (project root)
   - Pretty-print with 2-space indentation
```

### Data Structure

```json
{
  "generated": "2025-01-31T15:30:00Z",
  "summary": {
    "total_changes": 15,
    "total_phases": 4,
    "active_count": 5,
    "archived_count": 10
  },
  "dependency_graph": {
    "feat-dashboard": {
      "requires": [
        { "id": "feat-auth", "reason": "needs user model and session" },
        { "id": "feat-db", "reason": "needs migrations for dashboard" }
      ],
      "enables": [
        { "id": "feat-analytics", "reason": "unlocks event tracking" }
      ]
    }
  },
  "timeline": [
    {
      "phase": 1,
      "parallel": true,
      "changes": [
        {
          "id": "feat-auth",
          "title": "User Authentication",
          "description": "OAuth and session management",
          "tasks": { "total": 8, "completed": 3 },
          "blocked_by": [],
          "notes": "Critical path item. Blocks 3 downstream changes."
        }
      ]
    }
  ]
}
```

### File Organization

```text
.agents/skills/spectr-timeline/
├── SKILL.md           # Skill definition with instructions
└── references/
    └── timeline-schema.md  # Detailed JSON schema documentation
```

### Error Handling

1. **Circular dependencies**: Report cycle with all nodes involved, fail loudly
2. **Malformed frontmatter**: Log warning, skip that change, continue
3. **Missing proposal.md**: Log warning, skip that change
4. **Empty changes directory**: Success, generate empty timeline
5. **Parse errors**: Report file + line + error message

### Integration Points

1. **Reads from**: `spectr/changes/*/proposal.md` (frontmatter parsing)
2. **Writes to**: `./spectr/timeline.json`
3. **No modifications** to existing files or specs
4. **Optional**: Future CLI command could reuse analysis logic

### Reusability

The skill can be invoked in different contexts:

- `/spectr:timeline` - analyze all changes
- `/spectr:timeline feat-auth` - analyze single change + dependencies only
- Called programmatically by planning tools

## Validation Criteria

1. Timeline correctly identifies implementation phases
2. Circular dependencies are detected and reported
3. JSON is well-formed and pretty-printed
4. All active changes are included
5. Dependency reasons are preserved
6. Phase assignment is correct (dependencies satisfied before use)
7. Skill is discoverable by claude-code CLI
