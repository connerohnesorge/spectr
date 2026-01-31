# Implementation Design for /spectr:next Command

## Overview

The `/spectr:next` slash command will automate task execution for AI agents
working with spectr proposals. It finds the next pending task, executes it,
and updates statuses automatically.

## Key Components

### 1. Slash Command Definition

**File**: `internal/domain/slashcmd.go`

Add new enum value:

```go
const (
  SlashProposal SlashCommand = iota
  SlashApply
  SlashNext  // NEW: Execute next pending task
)
```

Update String() method to include "next".

### 2. Task Discovery Logic

**Responsibilities**:

- Parse tasks.jsonc (handle both flat version 1 and hierarchical version 2)
- Follow $ref links to child task files
- Find first task with status "pending"
- Return task details: ID, description, section, any children

**Implementation approach**:

- Reuse existing parsers from `internal/parsers/`
- Create new `taskdiscovery` package or add to existing `discovery`
- Handle circular reference detection
- Support both flat and hierarchical task schemas

### 3. Task Execution Engine

**Responsibilities**:

- Map task descriptions to actions
- Support common task types:
  - Code implementation (single file changes)
  - Test execution
  - Documentation updates
  - File creation
  - Validation tasks
- Provide extensibility for custom task types

**Simple Version 1 approach**:

- Parse task description for keywords
- Use heuristics to determine action type
- For unknown tasks, provide clear instructions to AI agent
- Delegate actual implementation to specialized sub-agents

### 4. Status Management

**Responsibilities**:

- Update task status: pending → in_progress before execution
- Update task status: in_progress → completed after execution
- Handle failures gracefully (in_progress → pending on error)
- Save updated tasks.jsonc atomically

**Implementation**:

- Reuse existing task sync logic from `internal/sync/`
- Ensure atomic writes with proper error handling
- Support both hierarchical and flat file structures

### 5. Provider Template Updates

**Files to update**:

- `internal/initialize/templates.go` - add SlashNext template
- All provider files in `internal/initialize/providers/` - add SlashNext to
  command maps

**Template structure**:

- Include discovery logic
- Include execution engine
- Include status management
- Provide clear instructions for AI agent

## Implementation Sequence

1. **Core domain changes**:
   - Add SlashNext to domain.SlashCommand
   - Update all provider maps to include SlashNext

2. **Task discovery**:
   - Create task discovery package
   - Implement tasks.jsonc parsing with $ref support
   - Test with both v1 and v2 schemas

3. **Template creation**:
   - Create SlashNext template in templates.go
   - Include basic task execution patterns
   - Add status management logic

4. **Provider updates**:
   - Update all 15+ provider initializers
   - Generate slash command files

5. **Testing**:
   - Unit tests for task discovery
   - Integration tests with sample proposals
   - Test with hierarchical task files

## Design Decisions

### Heuristic-based vs. Explicit Task Types

**Chosen**: Heuristic-based approach for simplicity

- Parse task description keywords ("Implement", "Test", "Update", "Create")
- Map to common actions
- Provide fallback for unknown tasks

**Alternative considered**: Explicit task types in schema

- Add "type" field to tasks.jsonc
- More precise but requires spec changes
- Can be added later if needed

### Single Command vs. Multi-Step

**Chosen**: Single `/spectr:next` command

- Simple interface for AI agents
- Automatically handles discovery → execution → status update

**Alternative**: Separate commands (`/spectr:discover`, `/spectr:execute`,
`/spectr:complete`)

- More flexible but more complex
- Requires agent to manage state

### Template Implementation Language

**Chosen**: Keep templates in Go with shell commands

- Consistent with existing slash commands
- AI agents can read and understand the logic
- Easy to extend and customize per provider

## Future Extensions

1. **Task Type Registry**: Explicit task types with custom handlers
2. **Progress Reporting**: Detailed execution logs and metrics
3. **Parallel Execution**: Run independent tasks concurrently
4. **Dependency Resolution**: Honor task dependencies before execution
5. **Custom Task Handlers**: Allow projects to define their own task executors
