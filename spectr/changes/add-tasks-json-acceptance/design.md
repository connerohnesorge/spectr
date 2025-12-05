## Context

This change implements learnings from Anthropic's research on effective harnesses for long-running agents. The core insight is that JSON-formatted task lists are more stable than Markdown when worked on by AI agents across multiple context windows. Agents working incrementally are less likely to accidentally modify or overwrite JSON files.

**Stakeholders**: Developers using Spectr with AI coding assistants (Claude Code, Cursor, etc.)

**Constraints**:
- Must preserve existing tasks.md format for human authoring during proposal phase
- Must be backwards compatible (changes without tasks.json should still work)
- Must integrate cleanly with existing archive workflow

## Goals / Non-Goals

**Goals**:
- Enable stable task tracking for long-running AI agent sessions
- Preserve full task information (sections, IDs, descriptions, status)
- Provide clear acceptance gate in workflow
- Support incremental task completion by agents
- Maintain human-readable proposal experience

**Non-Goals**:
- Automatic conversion without explicit acceptance step
- Changing the tasks.md authoring format
- Requiring tasks.json for all workflows
- Complex status tracking beyond completed/pending

## Decisions

### Decision 1: Two-Phase Task Lifecycle

**What**: Tasks exist as `tasks.md` during proposal phase, converted to `tasks.json` at acceptance.

**Why**:
- tasks.md is easier for humans to author during proposal creation
- tasks.json is more stable for AI agents during implementation
- Clear lifecycle boundary at "acceptance" moment

**Alternatives considered**:
- Always use JSON: Harder for humans to author, less readable in proposals
- Keep both: Risk of drift between formats
- Convert at archive: Too late, agents need JSON during implementation

### Decision 2: `spectr accept` as Explicit Command

**What**: New CLI subcommand that performs conversion and removes original.

**Why**:
- Explicit acceptance gate aligns with proposal workflow
- CLI command can be validated and tested
- Provides clear point of no return for format transition

**Alternatives considered**:
- Magic conversion on first task completion: Too implicit, confusing
- Conversion during archive: Doesn't help agents during implementation

### Decision 3: tasks.json Schema

**What**: Structured JSON preserving sections and task hierarchy with unlimited nesting:
```json
{
  "version": "1.0",
  "changeId": "add-feature",
  "acceptedAt": "2025-12-05T10:30:00Z",
  "sections": [
    {
      "name": "Implementation",
      "number": 1,
      "tasks": [
        {
          "id": "1.1",
          "description": "Create database schema\n- Parse requirement headers\n- Extract requirement name",
          "completed": true,
          "subtasks": [
            {
              "id": "1.1.1",
              "description": "Add users table",
              "completed": false,
              "subtasks": [
                {
                  "id": "1.1.1.1",
                  "description": "Add primary key constraint",
                  "completed": false,
                  "subtasks": []
                }
              ]
            }
          ]
        }
      ]
    }
  ],
  "summary": {
    "total": 4,
    "completed": 1
  }
}
```

**Why**:
- Preserves section grouping for logical organization
- Task IDs enable precise references (e.g., "complete task 1.3")
- Unlimited recursive nesting via `subtasks` array supports deep task hierarchies
- Indented detail lines appended to description with newlines
- Summary enables quick progress checks without parsing (includes all nested tasks)
- Version field allows future schema evolution
- acceptedAt provides audit trail

**Alternatives considered**:
- Flat task list: Loses section context
- Minimal schema: Harder to extend later

### Decision 4: Remove tasks.md After Conversion

**What**: Delete tasks.md after successful tasks.json creation.

**Why**:
- Prevents drift between two sources of truth
- Forces agents to use JSON format
- Clear signal that proposal is now "accepted"

**Alternatives considered**:
- Keep both with warning: Risk of agents updating wrong file
- Rename to tasks.md.bak: Adds clutter, unclear lifecycle

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Users accidentally accept before tasks are finalized | Clear warning prompt, --yes flag for automation |
| JSON corruption by agent errors | Validate JSON before writing, atomic writes |
| Loss of tasks.md without backup | Prompt user, keep in git history |
| Parser fails on unusual task formats | Extensive test fixtures from archive |

## Migration Plan

1. **No migration required** - New feature, existing changes continue to work
2. **Existing tasks.md files remain valid** - Archive workflow unchanged
3. **Apply slash command updated** - Will prompt to run `spectr accept` if tasks.json missing

## Resolved Questions

1. **Should `spectr accept` work on partial task completion?**
   - **Decision**: Yes, preserve existing `[x]` markers as `"completed": true` in tasks.json

2. **Should we support nested subtasks (e.g., 1.1.1)?**
   - **Decision**: Unlimited recursive nesting - 1.1.1.1 is subtask of 1.1.1 which is subtask of 1.1

3. **How should indented detail lines under tasks be handled?**
   - **Decision**: Append to description - concatenate indented lines into the task description with newlines

4. **How should the apply slash command integrate?**
   - **Decision**: Check for tasks.json; if missing, instruct agent to run `spectr accept <change-id>` and use tasks.json for tracking from then on

5. **Should tasks.json be excluded from validation if tasks.md exists?**
   - **Decision**: Yes, only one format should exist at a time
