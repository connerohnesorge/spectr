# SUPERSEDED - Not Implemented

**Date Archived**: 2026-01-18

**Superseded By**: `add-size-based-task-splitting`

## Why This Change Was Not Implemented

This proposal introduced hierarchical task splitting tied to delta spec directories (`specs/<capability>/tasks.jsonc`). While the goal of managing large task files was valid, the implementation approach proved too complex:

1. **Tight coupling to delta specs**: Required tasks to map to specific delta spec capabilities, creating unnecessary coupling between task organization and spec structure
2. **Complex mapping logic**: Auto-splitting by section name matching to capabilities added fragile heuristics
3. **Unclear navigation**: Agents would need to discover and navigate capability-specific task files, increasing cognitive load

## The Better Approach

The `add-size-based-task-splitting` change replaced this with a simpler approach:

1. **Size-based splitting only**: Split when `tasks.md` exceeds 100 lines, regardless of delta spec structure
2. **Section-based organization**: Split by markdown sections (## N. headers) rather than capabilities
3. **Flat structure with includes**: All task files live alongside the root `tasks.jsonc` (e.g., `tasks-1.jsonc`, `tasks-2.jsonc`) with an includes field for discovery
4. **Status aggregation**: Parent tasks automatically aggregate child task statuses (pending/in_progress/completed)

This keeps the benefits (readable task files, single Read operations) while avoiding the complexity of capability-based organization.

## Key Learnings

- Task organization should be independent of spec organization
- Size-based thresholds are more objective than semantic mapping
- Simpler file discovery patterns (glob includes) beat complex navigation logic
- Keep concerns separated: specs describe requirements, tasks describe work breakdown

## Original Proposal

See `proposal.md`, `design.md`, `tasks.md`, and `specs/` in this archive for the full original proposal.
