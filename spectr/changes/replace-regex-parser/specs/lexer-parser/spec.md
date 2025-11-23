# Lexer/Parser Capability

## ADDED Requirements

### Requirement: Parse Markdown Structure
The system SHALL parse markdown documents into a structured Abstract Syntax Tree (AST) that accurately represents headers, paragraphs, code blocks, and lists.

#### Scenario: Parsing basic structure
WHEN the parser processes a document with headers and paragraphs
THEN it produces an AST with corresponding Header and Paragraph nodes
AND the hierarchy is preserved

#### Scenario: Parsing code blocks
WHEN the parser processes a document containing code blocks
THEN the content inside the code blocks is treated as literal text
AND markdown syntax inside code blocks is NOT parsed as structure

### Requirement: Extract Spectr Elements
The system SHALL extract Spectr-specific elements (Requirements, Scenarios, Delta sections) by traversing the AST.

#### Scenario: Extracting valid requirement
WHEN the extractor traverses an AST with a "### Requirement: Foo" header
THEN it identifies a Requirement named "Foo"
AND captures the following content as the requirement body

#### Scenario: Ignoring syntax in code blocks
WHEN the extractor encounters "### Requirement: Foo" inside a code block
THEN it ignores it
AND does NOT create a Requirement entity

### Requirement: Report Parsing Errors
The system SHALL report parsing errors with precise line and column information.

#### Scenario: Reporting malformed input
WHEN the parser encounters invalid syntax
THEN it returns an error
AND the error includes the line and column number of the failure
