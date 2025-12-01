## ADDED Requirements

### Requirement: Instruction File Pointer Template

The system SHALL use a short pointer template when injecting Spectr instructions into root-level instruction files (e.g., `CLAUDE.md`, `AGENTS.md` at project root), directing AI assistants to read `spectr/AGENTS.md` for full instructions rather than duplicating the entire content.

#### Scenario: Init creates instruction file with pointer

- **WHEN** user runs `spectr init` and selects an AI tool (e.g., Claude Code)
- **THEN** the root-level instruction file (e.g., `CLAUDE.md`) contains a short pointer between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers
- **AND** the pointer directs AI assistants to read `spectr/AGENTS.md` when handling proposals, specs, or changes
- **AND** the full instructions remain only in `spectr/AGENTS.md`

#### Scenario: Update refreshes instruction file with pointer

- **WHEN** user runs `spectr update` on an already-initialized project
- **THEN** the root-level instruction files are updated with the short pointer content
- **AND** the `spectr/AGENTS.md` file retains the full instructions

#### Scenario: Pointer content is concise

- **WHEN** the instruction pointer template is rendered
- **THEN** the output is less than 20 lines
- **AND** the output explains when to read `spectr/AGENTS.md` (proposals, specs, changes, planning)
- **AND** the output does NOT duplicate the full workflow instructions
