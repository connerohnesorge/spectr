# Change: Add /spectr:next Slash Command for Task Execution

## Why

AI agents working with spectr proposals currently need to manually:
1. Read tasks.jsonc to find the next pending task
2. Understand the task requirements
3. Execute the task manually
4. Update the task status

This creates friction and requires agents to implement task iteration logic themselves. By providing a `/spectr:next` slash command that automatically executes the next pending task in the current change proposal, we streamline the agent workflow and ensure consistent task execution order.

## What Changes

- **ADDED**: New `SlashNext` constant to `domain.SlashCommand` enum
- **ADDED**: `/spectr:next` slash command template that:
  - Discovers the current change proposal directory
  - Parses tasks.jsonc to find the first pending task
  - Reads the task description and requirements
  - Executes the appropriate action based on task type
  - Updates task status from pending → in_progress → completed
- **MODIFIED**: `domain/slashcmd.go` - add SlashNext to enum
- **MODIFIED**: All provider initializers to include SlashNext command
- **MODIFIED**: Slash command templates to support the new command

## Impact

- **Affected specs**: `provider-system`
- **Affected code**:
  - `internal/domain/slashcmd.go` - add SlashNext enum value
  - `internal/initialize/providers/` - update all provider initializers
  - `internal/initialize/templates.go` - add SlashNext template
- **Breaking changes**: None - purely additive
- **New functionality**: AI agents can now use `/spectr:next` to automatically work through task lists

## Examples

### Before (manual task execution)

```
Human: Work on the add-feature-x proposal

AI Agent:
1. Read spectr/changes/add-feature-x/tasks.jsonc
2. Find first pending task (task #3)
3. Understand task requirements
4. Manually execute task
5. Update tasks.jsonc status
6. Report completion
```

### After (automated with /spectr:next)

```
Human: /spectr:next add-feature-x

AI Agent:
- Automatically executes next pending task
- Updates status automatically
- Reports progress without manual intervention
```

### Slash Command File Structure

```
.claude/commands/spectr/
├── proposal.md    # Create new proposal
├── apply.md       # Apply a change
└── next.md        # Execute next task [NEW]
```

### Command Behavior

The `/spectr:next` command will:
1. Find the current change directory in spectr/changes/
2. Read tasks.jsonc (following $ref links if hierarchical)
3. Locate the first task with status "pending"
4. Execute the task based on its description
5. Update status: pending → in_progress before execution
6. Update status: in_progress → completed after execution
7. Report what was done and what's next

### Example Task Execution Flow

```
User: /spectr:next

AI Agent:
→ Found next pending task: #5 "Update API documentation"
→ Marked task #5 as in_progress
→ Updating API documentation...
→ Completed task #5
→ Next task: #6 "Run tests and fix any failures"
```
