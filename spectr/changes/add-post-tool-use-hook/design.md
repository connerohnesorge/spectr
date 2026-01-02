# Design: postToolUse Hook System

## Overview

The postToolUse hook enables read-only observation of Claude Code tool results.
This document captures key design decisions.

## Architecture

```
Claude Code Tool Execution
         │
         ▼
   Tool completes
         │
         ▼
┌─────────────────────┐
│ conclaude intercepts│
│  via hook system    │
└─────────────────────┘
         │
         ▼
   Match tool filters
         │
         ├── No match → Skip
         │
         ▼
   Set environment vars
   CONCLAUDE_TOOL_NAME
   CONCLAUDE_TOOL_INPUT
   CONCLAUDE_TOOL_OUTPUT
   CONCLAUDE_TOOL_TIMESTAMP
         │
         ▼
   Execute hook command
         │
         ├── Sync: Wait for completion
         │
         └── Async: Fire and forget
         │
         ▼
   Claude continues
   (original output preserved)
```

## Key Decisions

### Environment Variables vs JSON stdin

**Decision**: Environment variables

**Rationale**:
- Simpler for shell scripts - no need to parse stdin
- Works with any command/script language
- Matches existing preToolUse pattern in conclaude
- Trade-off: Large outputs may hit shell limits (~128KB typical)

**Mitigation**: For large outputs, users can:
1. Configure to skip large-output tools
2. Use a script that reads from a temp file (future enhancement)

### Read-Only vs Modifiable

**Decision**: Read-only (hooks cannot modify output)

**Rationale**:
- Simpler mental model - hooks are observers
- Avoids complex error handling for modification failures
- Prevents accidental corruption of tool outputs
- Matches the "logging" use case perfectly
- Future: Could add separate `transformToolOutput` hook if needed

### Sequential vs Parallel Hooks

**Decision**: Sequential with async option

**Rationale**:
- Sequential by default ensures deterministic logging order
- `async: true` option for fire-and-forget scenarios
- Matches behavior of other conclaude hooks
- Prevents race conditions in file appending

## Integration with Spectr

Spectr can leverage postToolUse for:

1. **Q&A Documentation**: Log AskUserQuestion interactions during proposal
   creation
2. **Audit Trail**: Track all tool usage during spec implementation
3. **Search Caching**: Log WebSearch results for offline reference

Example `.conclaude.yaml` for Spectr projects:

```yaml
postToolUse:
  commands:
    # Log all Q&A interactions for documentation
    - tool: "AskUserQuestion"
      run: ".claude/scripts/log-qa.sh"

    # Optional: Log search results
    - tool: "*Search*"
      run: ".claude/scripts/log-search.sh"
      async: true
```

## Future Considerations

1. **Temp file mode**: For large outputs, write to temp file and pass path
2. **Transform hooks**: Separate hook type that CAN modify outputs
3. **Agent-specific hooks**: Filter by subagent type (like subagentStop)
4. **Output streaming**: Hook receives output as it streams (for long tools)
