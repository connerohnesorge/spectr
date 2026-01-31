# Spec Delta

## ADDED Requirements

### Requirement: SlashNext Command Support

The system SHALL support a `SlashNext` command identifier for the
`/spectr:next` slash command.

#### Scenario: SlashNext is defined in enum



- GIVEN the domain.SlashCommand type
- WHEN examining the available constants
- THEN SlashNext SHALL be present with a unique value
- AND SlashNext.String() SHALL return "next"

#### Scenario: SlashNext is distinct from existing commands

- GIVEN existing SlashCommand values (SlashProposal, SlashApply)
- WHEN SlashNext is added
- THEN it SHALL have a different integer value
- AND it SHALL not conflict with existing command behavior

### Requirement: Provider SlashNext Integration

All AI assistant provider initializers SHALL include SlashNext in their
slash command maps.

#### Scenario: Claude provider includes SlashNext

- GIVEN the Claude provider initializer
- WHEN constructing command maps
- THEN SlashNext SHALL be mapped to the appropriate template

#### Scenario: Other providers include SlashNext

- GIVEN any provider initializer (Windsurf, Cursor, Continue, etc.)
- WHEN constructing command maps
- THEN SlashNext SHALL be present in their command maps

### Requirement: SlashNext Template Definition

The system SHALL provide a template for the `/spectr:next` slash command
that enables AI agents to execute the next pending task.

#### Scenario: Template includes task discovery logic

- GIVEN a SlashNext template
- WHEN rendered and executed
- THEN it SHALL discover the current change proposal
- AND it SHALL parse tasks.jsonc to find pending tasks

#### Scenario: Template includes status management

- GIVEN the SlashNext template execution
- WHEN beginning task execution
- THEN it SHALL update task status from "pending" to "in_progress"
- AND upon completion, it SHALL update from "in_progress" to "completed"

#### Scenario: Template handles hierarchical tasks

- GIVEN a slash command with hierarchical tasks (version 2 schema)
- WHEN discovering the next task
- THEN it SHALL follow $ref links to child task files
- AND it SHALL search for the first pending task recursively

#### Scenario: Template provides execution guidance

- GIVEN a discovered pending task
- WHEN executing the SlashNext command
- THEN the AI agent SHALL receive clear instructions on what action to take
- AND it SHALL understand the expected outcome

### Requirement: SlashNext Command Generation

The system SHALL generate `/spectr:next` command files during provider
initialization and ensure they are properly formatted.

#### Scenario: SlashNext file is created

- GIVEN a provider with SlashNext in its command map
- WHEN the provider initializer runs
- THEN a file named "next.md" (or prefixed variant) SHALL be created
- AND it SHALL contain executable instructions for the AI agent

#### Scenario: SlashNext file includes proposal context

- GIVEN a SlashNext command file
- WHEN examining its contents
- THEN it SHALL include context about spectr task workflow
- AND it SHALL explain the purpose of the next command

### Requirement: Flat Task File Compatibility

The SlashNext command SHALL support version 1 flat task files without `$ref`
links and without hierarchical structures.

#### Scenario: Flat tasks.jsonc processing

- GIVEN a proposal with version 1 tasks.jsonc (no children, no includes)
- WHEN executing SlashNext
- THEN it SHALL find the first pending task by sequential scan
- AND update the flat file structure correctly

### Requirement: Hierarchical Task File Compatibility

The SlashNext command SHALL support version 2 hierarchical task files with
`$ref` links.

#### Scenario: Hierarchical task discovery

- GIVEN a proposal with version 2 tasks.jsonc (with children and includes)
- WHEN discovering the next task
- THEN it SHALL resolve $ref references to child files
- AND find the first pending task across the hierarchy

#### Scenario: Parent task status aggregation

- GIVEN hierarchical tasks with child task statuses
- WHEN all children of a parent task are completed
- THEN the parent task status SHALL be automatically updated to "completed"

### Requirement: Error Handling

The SlashNext command SHALL handle errors gracefully and provide clear feedback.

#### Scenario: No pending tasks

- GIVEN a proposal where all tasks are completed
- WHEN executing SlashNext
- THEN it SHALL report "No pending tasks found"
- AND provide a summary of completed work

#### Scenario: Invalid task file

- GIVEN a malformed tasks.jsonc file
- WHEN attempting to discover tasks
- THEN it SHALL report the parse error clearly
- AND suggest running `spectr validate`

#### Scenario: Missing referenced file

- GIVEN hierarchical tasks with a $ref to a missing file
- WHEN resolving references
- THEN it SHALL report the missing file
- AND provide the expected file path

### Requirement: Execution Reporting

The SlashNext command SHALL report what task was executed and what comes next.

#### Scenario: Successful task execution

- GIVEN a task was successfully executed
- WHEN SlashNext completes
- THEN it SHALL report:

  - The task ID and description that was executed
  - The new status (completed)
  - The next pending task (if any)
  - Progress summary (e.g., "3 of 7 tasks completed")

#### Scenario: Task execution failure

- GIVEN a task execution failed
- WHEN the failure occurs
- THEN it SHALL report the error
- AND update status from "in_progress" back to "pending"
- AND provide guidance on how to proceed
