# Pathological Test Corpus

## Purpose
This file tests edge cases that commonly break regex-based parsers.

## Requirements

### Requirement: Code Block with Markdown Syntax
The parser MUST correctly handle code blocks containing markdown-like text.

#### Scenario: Code block contains requirement header
- **WHEN** code block contains text like "### Requirement: Foo"
- **THEN** parser treats it as code, not a requirement
- **AND** does not extract it as a requirement

Here's an example that breaks regex parsers:

```markdown
# Example Spec

## Requirements

### Requirement: Example Feature
This shows what a requirement looks like.

#### Scenario: Basic test
- **WHEN** something happens
- **THEN** something else occurs
```

The above code block should NOT be parsed as actual requirements.

### Requirement: Deeply Nested Lists
The parser SHALL handle deeply nested list structures.

#### Scenario: Multi-level nesting
- **WHEN** list has many levels
  - Level 2 item
    - Level 3 item
      - Level 4 item
        - Level 5 item
- **THEN** parser preserves structure
- **AND** maintains correct hierarchy

### Requirement: Special Characters in Names
Requirements MAY contain special characters in their names.

#### Scenario: Unicode in requirement name
- **WHEN** requirement name contains émojis 🎉 or ünicode
- **THEN** parser handles it correctly
- **AND** preserves exact text

### Requirement: Escaped Characters
The parser SHALL handle escaped markdown characters.

#### Scenario: Escaped hash symbols
- **WHEN** text contains \### or \## or \#
- **THEN** these are treated as literal characters
- **AND** not as heading markers

#### Scenario: Backticks in text
- **WHEN** text contains `inline code` with backticks
- **THEN** parser preserves formatting
- **AND** doesn't confuse with code blocks

### Requirement: Adjacent Code Blocks
Multiple code blocks MAY appear consecutively.

#### Scenario: Back-to-back code blocks
- **WHEN** two code blocks are adjacent
- **THEN** both are parsed correctly

```go
func First() {
    // This is the first block
}
```

```python
def second():
    # This is the second block
    pass
```

### Requirement: Code Block Without Language
Code blocks MAY omit language specification.

#### Scenario: Plain code fence
- **WHEN** code block has no language
- **THEN** parser handles it correctly

```
This is a code block without a language
It should still work correctly
### Requirement: This is not a real requirement
```

### Requirement: Mixed Content Types
Requirements MAY contain various content types.

#### Scenario: Requirement with code and lists
- **WHEN** requirement contains mixed content
- **THEN** all content is preserved

Example implementation:

```typescript
interface User {
    id: string;
    name: string;
}

// Usage notes:
// 1. Always validate user input
// 2. Check permissions
// 3. Log access attempts
```

Additional notes:
- Point one about something
- Point two about another thing

### Requirement: Blank Lines in Requirements
Blank lines MAY appear within requirement content.

#### Scenario: Multiple blank lines
- **WHEN** requirement has several blank lines


- **THEN** they are preserved correctly
- **AND** requirement boundaries are clear

### Requirement: Long Scenario Names
Scenario names MAY be very long and contain complex punctuation.

#### Scenario: This is an extremely long scenario name that might cause issues with some parsers, especially those that make assumptions about line length or content, and it includes special characters like: colons, semicolons; parentheses (like these), brackets [like these], and even quotes "like these"
- **WHEN** scenario name is very long
- **THEN** parser handles it without truncation
- **AND** preserves all characters

### Requirement: HTML-Like Content
Content MAY contain HTML-like text that should be preserved.

#### Scenario: HTML tags in content
- **WHEN** content contains <div> or <span> tags
- **THEN** parser treats them as literal text
- **AND** does not interpret as HTML

Example: Use `<component>` tags in your XML configuration.

### Requirement: Inline Links and Images
Markdown links and images SHALL be preserved.

#### Scenario: Link in scenario
- **WHEN** text contains [link](https://example.com)
- **THEN** link is preserved correctly
- **AND** ![image](image.png) is also preserved

### Requirement: Tables in Content
Requirement content MAY include markdown tables.

#### Scenario: Table parsing
- **WHEN** content includes table

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Value A  | Value B  | Value C  |
| Value D  | Value E  | Value F  |

- **THEN** table structure is preserved

### Requirement: Blockquotes
Content MAY include blockquoted text.

#### Scenario: Blockquote handling
- **WHEN** content has blockquotes
- **THEN** they are preserved

> This is a blockquote
> It spans multiple lines
> And should be preserved

### Requirement: Horizontal Rules
Content MAY include horizontal rules.

#### Scenario: Horizontal rule
- **WHEN** content has horizontal rule
- **THEN** it is preserved correctly

---

### Requirement: Multiple Consecutive Hashes
Content with multiple hash symbols SHALL be handled correctly.

#### Scenario: Hash symbol handling
- **WHEN** content has ##### or ###### symbols
- **THEN** parser handles them appropriately
- **AND** doesn't confuse with headers

### Requirement: Empty Requirement Content
Requirements MAY have minimal content between header and scenarios.

#### Scenario: Just a scenario
- **WHEN** requirement has no body text
- **THEN** parser accepts it
- **AND** extracts scenarios correctly

### Requirement: Trailing Whitespace
Content MAY have trailing whitespace on lines.

#### Scenario: Whitespace handling
- **WHEN** lines have trailing spaces or tabs
- **THEN** parser handles consistently
- **AND** doesn't break on whitespace

### Requirement: Complex Code in Code Blocks
Code blocks MAY contain complex nested structures.

#### Scenario: Nested code structures
- **WHEN** code block has complex content
- **THEN** parser preserves everything

```javascript
const complexObject = {
  nested: {
    deeply: {
      structure: {
        with: "many levels",
        and: ["arrays", "too"],
        also: {
          functions: () => {
            // ### This is not a requirement
            // #### This is not a scenario
            return "Just a comment in code";
          }
        }
      }
    }
  }
};

/**
 * ## Documentation comment that looks like markdown
 *
 * ### Requirement: This is not a real requirement
 * It's just inside a doc comment in code.
 *
 * #### Scenario: Also not real
 * - **WHEN** you see this
 * - **THEN** ignore it (it's code)
 */
function testFunction() {
  const markdown = `
    # This is a template string
    ## ADDED Requirements
    ### Requirement: Fake requirement in string
  `;
  return markdown;
}
```

### Requirement: Delta Section Names in Code
Code MAY reference delta section names.

#### Scenario: Delta keywords in code
- **WHEN** code contains "ADDED" or "MODIFIED"
- **THEN** parser treats as code
- **AND** not as delta sections

Example code:

```python
# Constants for change tracking
STATUS_ADDED = "ADDED"
STATUS_MODIFIED = "MODIFIED"
STATUS_REMOVED = "REMOVED"

def track_change(status):
    if status == STATUS_ADDED:
        print("Item was added")
    elif status == STATUS_MODIFIED:
        print("Item was modified")
```

### Requirement: Scenario Steps with Code
Scenario steps MAY contain inline code.

#### Scenario: Code in steps
- **WHEN** user calls `api.authenticate(token)` method
- **THEN** system returns `{status: "ok", user: {...}}` object
- **AND** response contains `Authorization: Bearer <token>` header

### Requirement: Bold and Italic Text
Content MAY use markdown emphasis.

#### Scenario: Text formatting
- **WHEN** text is **bold** or *italic* or ***both***
- **THEN** parser preserves markdown
- **AND** doesn't confuse with WHEN/THEN markers

### Requirement: Numbered Lists in Scenarios
Scenarios MAY use numbered lists instead of bullets.

#### Scenario: Numbered list steps
1. **WHEN** user performs first action
2. **THEN** system responds
3. **AND** performs second action
4. **AND** completes workflow

### Requirement: Mixed List Markers
Lists MAY use different bullet characters.

#### Scenario: Various list markers
- First item with dash
* Second item with asterisk
+ Third item with plus
- **WHEN** all are used
- **THEN** parser handles consistently

### Requirement: Requirements at Document Start
First requirement MAY appear immediately after header.

#### Scenario: No gap before requirement
- **WHEN** requirement is first content
- **THEN** parser detects it correctly
- **AND** doesn't require preceding text

### Requirement: Consecutive Blank Lines
Content MAY have multiple consecutive blank lines.



#### Scenario: Many blank lines above
- **WHEN** many blank lines present
- **THEN** parser is not confused


- **AND** requirement boundaries remain clear
