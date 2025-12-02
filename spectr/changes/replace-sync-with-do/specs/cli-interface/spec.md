## ADDED Requirements

### Requirement: Spectr Do Slash Command

The `spectr:do` slash command SHALL execute all tasks from a specific numbered section of a change's `tasks.md` file, providing a streamlined way to implement approved changes.

#### Scenario: User invokes spectr:do with change-id and section number

- **WHEN** user runs `/spectr:do <change-id> <section-number>`
- **THEN** the command reads `spectr/changes/<change-id>/proposal.md` for context
- **AND** the command reads `spectr/changes/<change-id>/design.md` if it exists
- **AND** the command reads `spectr/changes/<change-id>/specs/*/spec.md` for requirements
- **AND** the command reads `spectr/changes/<change-id>/tasks.md`
- **AND** locates the section matching `## <section-number>.` header pattern
- **AND** executes each uncompleted task (`- [ ]`) in that section sequentially

#### Scenario: Task completion updates tasks.md

- **WHEN** a task is successfully completed
- **THEN** the task checkbox is updated from `- [ ]` to `- [x]` in `tasks.md`
- **AND** the next uncompleted task in the section is processed

#### Scenario: Missing arguments prompts for input

- **WHEN** user runs `/spectr:do` without arguments
- **THEN** the command asks user which change-id and section number to execute
- **AND** waits for user input before proceeding

#### Scenario: Invalid section number

- **WHEN** user provides a section number that does not exist in `tasks.md`
- **THEN** an error message is displayed listing available section numbers
- **AND** the command does not execute any tasks

#### Scenario: All tasks in section already complete

- **WHEN** all tasks in the specified section are already marked `- [x]`
- **THEN** the command reports that all tasks in the section are complete
- **AND** no implementation actions are taken

#### Scenario: Validation after section completion

- **WHEN** all tasks in the section have been executed
- **THEN** the command runs `spectr validate <change-id> --strict`
- **AND** reports any validation issues
- **AND** provides a summary of completed tasks
