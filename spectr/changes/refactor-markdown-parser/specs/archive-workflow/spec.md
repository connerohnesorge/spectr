# Archive Workflow Delta Spec

## MODIFIED Requirements

### Requirement: Delta Operation Parsing

The system SHALL parse delta operations from change spec files using a lexer/parser architecture that correctly handles markdown edge cases.

#### Scenario: Parse ADDED requirements

- **WHEN** a delta spec contains `## ADDED Requirements` section
- **THEN** the system SHALL extract all requirement blocks with headers and scenarios using AST-based extraction
- **AND** SHALL correctly ignore requirement headers that appear inside code blocks
- **AND** SHALL preserve formatting and structure of each requirement

#### Scenario: Parse MODIFIED requirements

- **WHEN** a delta spec contains `## MODIFIED Requirements` section
- **THEN** the system SHALL extract complete modified requirement blocks using AST-based extraction
- **AND** SHALL correctly handle MODIFIED requirements containing code blocks
- **AND** SHALL match requirement names using normalized comparison (whitespace-insensitive)

#### Scenario: Parse REMOVED requirements

- **WHEN** a delta spec contains `## REMOVED Requirements` section
- **THEN** the system SHALL extract requirement names to be removed
- **AND** SHALL correctly ignore requirement headers inside code blocks
- **AND** SHALL extract only actual requirement names, not code examples

#### Scenario: Parse RENAMED requirements

- **WHEN** a delta spec contains `## RENAMED Requirements` section with FROM/TO pairs
- **THEN** the system SHALL extract the old and new requirement names using AST-based parsing
- **AND** SHALL match the format: `- FROM: \`### Requirement: Old Name\`` and `- TO: \`### Requirement: New Name\``
- **AND** SHALL handle whitespace variations in the FROM/TO syntax

#### Scenario: Require at least one delta operation

- **WHEN** a delta spec has no ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** the system SHALL return an error indicating no delta operations were found
- **AND** the error SHALL include position information for debugging

#### Scenario: Handle code blocks in delta specs

- **WHEN** a delta spec contains code blocks with markdown syntax examples
- **THEN** the parser SHALL NOT extract requirement headers from inside code blocks
- **AND** SHALL NOT extract scenario headers from inside code blocks
- **AND** SHALL preserve code block content exactly as written in the delta spec

#### Scenario: Report parsing errors with position

- **WHEN** a delta spec contains malformed markdown (e.g., unclosed code fence)
- **THEN** the system SHALL report a ParseError with line and column information
- **AND** the error message SHALL include context showing the problematic markdown
- **AND** SHALL prevent archiving until the malformed markdown is fixed
