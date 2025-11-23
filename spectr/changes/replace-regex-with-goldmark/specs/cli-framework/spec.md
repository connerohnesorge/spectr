## MODIFIED Requirements

### Requirement: Title Extraction
The system SHALL extract titles from proposal and spec markdown files by parsing the markdown document with goldmark AST and locating the first level-1 heading node, then removing the "Change:" or "Spec:" prefix if present.

#### Scenario: Extract title from proposal using AST
- **WHEN** the system reads a `proposal.md` file with heading `# Change: Add Feature`
- **THEN** it parses the file into a goldmark AST
- **AND** traverses the AST to find the first H1 (heading level 1) node
- **AND** extracts the title as "Add Feature" after removing the "Change:" prefix

#### Scenario: Extract title from spec using AST
- **WHEN** the system reads a `spec.md` file with heading `# CLI Framework`
- **THEN** it parses the file into a goldmark AST
- **AND** traverses the AST to find the first H1 node
- **AND** extracts the title as "CLI Framework"

#### Scenario: Fallback to ID when title not found
- **WHEN** the system cannot extract a title from a markdown file (no H1 node found)
- **THEN** it uses the directory name (ID) as the title

#### Scenario: Handle edge cases in code blocks
- **WHEN** the markdown file contains `# Heading` text inside code blocks
- **THEN** the AST walker SHALL ignore heading-like text within code block nodes
- **AND** only extract titles from actual heading AST nodes

### Requirement: Task Counting
The system SHALL count tasks in `tasks.md` files by parsing the markdown document with goldmark AST and traversing list item nodes to identify task checkbox patterns `- [ ]` or `- [x]` (case-insensitive), with completed tasks marked by `[x]`.

#### Scenario: Count completed and total tasks using AST
- **WHEN** the system reads a `tasks.md` file with 3 tasks, 2 marked `[x]` and 1 marked `[ ]`
- **THEN** it parses the file into a goldmark AST
- **AND** traverses list item nodes in the AST
- **AND** identifies task checkbox patterns in list item text
- **AND** reports `taskStatus` as `{ total: 3, completed: 2 }`

#### Scenario: Handle missing tasks file
- **WHEN** the system cannot find or read a `tasks.md` file for a change
- **THEN** it reports `taskStatus` as `{ total: 0, completed: 0 }`
- **AND** continues processing without error

#### Scenario: Handle nested lists correctly
- **WHEN** the tasks.md file contains nested task lists
- **THEN** the AST walker SHALL traverse all list item nodes recursively
- **AND** count tasks at all nesting levels
- **AND** correctly identify completion status for each task

#### Scenario: Ignore checkbox-like text in code blocks
- **WHEN** `tasks.md` contains fenced code blocks with lines that look like tasks (e.g., "- [x] not a task")
- **THEN** the AST walker SHALL treat the checkbox-like text as code content
- **AND** SHALL NOT count those lines toward total or completed task counts

### Requirement: Data Reuse from Discovery and Parsers
The view command SHALL reuse existing discovery and parsing infrastructure to avoid code duplication, utilizing goldmark AST-based parsers with enhanced internal data structures.

#### Scenario: Discover changes and specs
- **WHEN** building dashboard data
- **THEN** use `internal/discovery` package functions to find changes
- **AND** use `internal/discovery` package functions to find specs
- **AND** exclude archived changes from active/completed lists

#### Scenario: Parse titles and counts with goldmark
- **WHEN** extracting metadata from markdown files
- **THEN** use `internal/parsers` package goldmark-based parsers to parse proposal.md for titles
- **AND** use `internal/parsers` package goldmark-based parsers to parse spec.md for titles and requirement counts
- **AND** use `internal/parsers` package goldmark-based parsers to parse tasks.md for task counts
- **AND** leverage goldmark AST traversal for robust parsing

#### Scenario: Handle parser internal data structures
- **WHEN** parsers return metadata
- **THEN** internal data structures MAY include AST metadata fields (ASTNode, LineNumber)
- **AND** view command SHALL use public parser API without depending on internal AST details
- **AND** parser implementation details remain encapsulated within internal/parsers package
