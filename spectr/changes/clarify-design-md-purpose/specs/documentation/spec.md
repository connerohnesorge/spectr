## MODIFIED Requirements

### Requirement: design.md File Purpose Clarification
The `design.md` file SHALL contain specific implementation details to guide developers and AI agents during implementation. This includes code structures, API signatures, data models, file structures, and concrete examples—not abstract principles or high-level decisions.

#### Scenario: Directory structure comments updated
- **WHEN** reading directory structure comments in documentation
- **THEN** `design.md` SHALL be described as "Implementation details (code structures, APIs, data models)" instead of "Technical patterns" or "Technical decisions"

#### Scenario: File purpose clarification added
- **WHEN** reading file purpose documentation
- **THEN** there SHALL be a clear distinction between:
  - `proposal.md`: High-level why/what/impact
  - `tasks.md`: Implementation checklist
  - `spec.md`: Requirements and acceptance criteria
  - `design.md`: Specific implementation details with code examples

#### Scenario: What to include section added
- **WHEN** reading design.md guidance
- **THEN** there SHALL be explicit examples of what to include:
  - Data structures and type definitions
  - API signatures and interfaces
  - File and directory structures
  - Code snippets showing patterns
  - Configuration schemas

#### Scenario: Before/after examples provided
- **WHEN** reading design.md guidance
- **THEN** there SHALL be concrete examples showing:
  - **VAGUE (avoid)**: "Use clean architecture principles"
  - **DETAILED (preferred)**: Specific struct definitions, function signatures, and composition patterns"

### Requirement: Documentation Files to Update
The following documentation files SHALL be updated to clarify design.md purpose:

#### Scenario: AGENTS.md updated
- **WHEN** reading `spectr/AGENTS.md`
- **THEN** the directory structure comments and skeleton SHALL describe design.md as containing "Implementation details"

#### Scenario: project.md updated
- **WHEN** reading `spectr/project.md`
- **THEN** there SHALL be a File Purposes section clarifying each file type with emphasis on design.md's implementation detail purpose

#### Scenario: README.md updated
- **WHEN** reading `README.md` directory structure and FAQ sections
- **THEN** `design.md` SHALL be described as containing "Implementation details with code examples"

#### Scenario: User guides updated
- **WHEN** reading `docs/src/content/docs/guides/creating-changes.md`
- **THEN** the design document section SHALL include examples of specific implementation details to include
