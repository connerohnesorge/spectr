## MODIFIED Requirements

### Requirement: Title Extraction
The system SHALL extract titles from proposal and spec markdown files using goldmark AST parsing to find the first level-1 heading and removing "Change:" or "Spec:" prefix if present, correctly handling CommonMark edge cases.

#### Scenario: Extract title from proposal
- **WHEN** the system reads a `proposal.md` file with heading `# Change: Add Feature`
- **THEN** it extracts the title as "Add Feature" using AST traversal for H1 nodes

#### Scenario: Extract title from spec
- **WHEN** the system reads a `spec.md` file with heading `# CLI Framework`
- **THEN** it extracts the title as "CLI Framework" by finding first H1 in AST

#### Scenario: Fallback to ID when title not found
- **WHEN** the system cannot extract a title from a markdown file (no H1 nodes in AST)
- **THEN** it uses the directory name (ID) as the title

#### Scenario: Title in code block ignored
- **WHEN** a markdown file has `# Title` within a fenced code block before the real title
- **THEN** the system SHALL ignore the code block content using AST-based parsing
- **AND** SHALL extract the first structural H1 heading outside code blocks

#### Scenario: Title with inline code and formatting
- **WHEN** the H1 heading contains inline code or emphasis (e.g., `# Title with \`code\` and **bold**`)
- **THEN** the system SHALL extract text content from AST nodes
- **AND** SHALL produce clean title text: "Title with code and bold"

### Requirement: Task Counting
The system SHALL count tasks in `tasks.md` files by parsing markdown AST to find list items with checkbox patterns `- [ ]` or `- [x]` (case-insensitive), with completed tasks marked by `[x]`.

#### Scenario: Count completed and total tasks
- **WHEN** the system reads a `tasks.md` file with 3 tasks, 2 marked `[x]` and 1 marked `[ ]`
- **THEN** it parses the markdown to AST and identifies list items
- **AND** it reports `taskStatus` as `{ total: 3, completed: 2 }`

#### Scenario: Handle missing tasks file
- **WHEN** the system cannot find or read a `tasks.md` file for a change
- **THEN** it reports `taskStatus` as `{ total: 0, completed: 0 }`
- **AND** continues processing without error

#### Scenario: Nested task lists
- **WHEN** a `tasks.md` contains nested lists (e.g., sub-tasks indented under main tasks)
- **THEN** the system SHALL parse using AST to correctly identify all checkbox items
- **AND** SHALL count both parent and child checkboxes as separate tasks
- **AND** SHALL handle arbitrary nesting depth

#### Scenario: Tasks in code blocks ignored
- **WHEN** a `tasks.md` contains checkbox patterns within fenced code blocks
- **THEN** the system SHALL use AST parsing to distinguish code from structure
- **AND** SHALL NOT count checkboxes in code blocks as tasks
- **AND** SHALL only count actual markdown list items with checkboxes

## ADDED Requirements

### Requirement: AST-Based Requirement Parsing
The system SHALL parse requirement blocks from spec files using goldmark AST traversal to extract requirement headers, scenarios, and content with accurate source position tracking.

#### Scenario: Parse requirements from spec AST
- **WHEN** the system parses a spec.md file
- **THEN** it SHALL create goldmark AST from the file content
- **AND** SHALL traverse AST to find all `### Requirement:` headings (H3 nodes)
- **AND** SHALL collect requirement content until next H2 or H3 heading

#### Scenario: Extract requirement with scenarios
- **WHEN** a requirement block contains `#### Scenario:` headings (H4 nodes)
- **THEN** the system SHALL include all scenario content in the requirement block
- **AND** SHALL extract scenario names from H4 text content
- **AND** SHALL preserve scenario bodies including bullet points and paragraphs

#### Scenario: Track source positions
- **WHEN** parsing requirement blocks from AST
- **THEN** each RequirementBlock SHALL include SourcePos with start/end line:column
- **AND** SourcePos SHALL be derived from AST node segment positions
- **AND** SHALL enable precise error reporting (e.g., "spec.md:42:3: Requirement missing scenario")

#### Scenario: Handle requirements with code blocks
- **WHEN** a requirement body contains fenced code blocks with pseudo-headers
- **THEN** AST parsing SHALL correctly treat code as content, not structure
- **AND** SHALL NOT terminate requirement parsing at fake headers in code
- **AND** SHALL include entire code block in requirement Raw content

#### Scenario: Handle empty requirements
- **WHEN** a requirement has no body content (just header followed by next section)
- **THEN** the system SHALL create RequirementBlock with empty Raw content
- **AND** SHALL still track SourcePos for the header
- **AND** SHALL allow validation to detect and report missing content

### Requirement: AST-Based Delta Parsing
The system SHALL parse delta spec files using goldmark AST to extract ADDED, MODIFIED, REMOVED, and RENAMED delta operations with precise section boundaries.

#### Scenario: Parse delta sections from AST
- **WHEN** the system parses a delta spec file
- **THEN** it SHALL create goldmark AST from file content
- **AND** SHALL find `## ADDED Requirements` heading (H2 node) in AST
- **AND** SHALL extract all nodes between this heading and next H2 heading
- **AND** SHALL repeat for MODIFIED, REMOVED, RENAMED sections

#### Scenario: Extract requirements from delta sections
- **WHEN** processing ADDED or MODIFIED section nodes
- **THEN** the system SHALL traverse section nodes to find H3 requirement headers
- **AND** SHALL build RequirementBlock for each with SourcePos metadata
- **AND** SHALL collect requirement content until next requirement or section boundary

#### Scenario: Parse REMOVED section requirement names
- **WHEN** processing REMOVED Requirements section
- **THEN** the system SHALL extract requirement names from H3 headers
- **AND** SHALL NOT require full requirement bodies or scenarios
- **AND** SHALL track source positions for removed requirement references

#### Scenario: Parse RENAMED section FROM/TO pairs
- **WHEN** processing RENAMED Requirements section
- **THEN** the system SHALL use AST to find list nodes
- **AND** SHALL extract list item text for FROM and TO patterns
- **AND** SHALL apply regex to clean text (not raw markdown) for requirement name extraction
- **AND** SHALL track source positions for each rename pair

#### Scenario: Handle malformed delta sections
- **WHEN** a delta spec has section headers but no content
- **THEN** AST parsing SHALL detect empty sections (no H3 descendants)
- **AND** validation can report meaningful errors with section line numbers
- **AND** SHALL return empty slices for Added/Modified/Removed as appropriate

### Requirement: AST Utility Functions
The system SHALL provide reusable AST utility functions in `internal/parsers/ast_utils.go` for common markdown parsing operations.

#### Scenario: Parse markdown file to AST
- **WHEN** a function needs to parse markdown
- **THEN** it SHALL call `parseMarkdownToAST(filePath string) (ast.Node, []byte, error)`
- **AND** the function SHALL read file content and parse with goldmark
- **AND** SHALL return AST root node and original source bytes
- **AND** SHALL return error if file cannot be read or parsed

#### Scenario: Extract text content from AST nodes
- **WHEN** extracting readable text from AST nodes (headings, paragraphs)
- **THEN** functions SHALL call `extractTextContent(node ast.Node) string`
- **AND** the function SHALL recursively traverse text nodes
- **AND** SHALL concatenate all text, stripping markdown formatting
- **AND** SHALL produce clean human-readable text

#### Scenario: Find headings by level
- **WHEN** searching for headings of specific level (H1, H2, H3)
- **THEN** functions SHALL use `findHeadingsByLevel(doc ast.Node, level int) []*ast.Heading`
- **AND** SHALL traverse AST collecting all heading nodes matching level
- **AND** SHALL return slice of heading nodes with their positions

#### Scenario: Extract section content by heading
- **WHEN** extracting content between two section markers
- **THEN** functions SHALL call `extractSectionByHeading(doc ast.Node, heading string, level int) []ast.Node`
- **AND** SHALL find target heading by text match at specified level
- **AND** SHALL collect all nodes until next heading of same or higher level
- **AND** SHALL return slice of AST nodes representing section content

#### Scenario: Get source position from AST node
- **WHEN** error reporting needs precise locations
- **THEN** functions SHALL call `getSourcePosition(node ast.Node, source []byte) *SourcePosition`
- **AND** SHALL compute line:column from AST segment byte offsets
- **AND** SHALL return SourcePosition with StartLine, StartColumn, EndLine, EndColumn
- **AND** SHALL handle multi-line nodes correctly
