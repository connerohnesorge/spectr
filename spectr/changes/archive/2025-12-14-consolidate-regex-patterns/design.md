# Design: Consolidated Regex Package

## Context

Spectr parses markdown files to extract structural elements:

- H2 section headers (`## Requirements`, `## ADDED Requirements`)
- H3 requirement headers (`### Requirement: Name`)
- H4 scenario headers (`#### Scenario: Name`)
- Task checkboxes (`- [ ]`, `- [x]`)
- Delta operations (ADDED/MODIFIED/REMOVED/RENAMED sections)

Currently, regex patterns are defined locally in each file that needs them, leading to:

- 10+ instances of requirement header pattern
- 6+ instances of section header pattern
- 4+ instances of scenario pattern
- Patterns compiled on every function call (no caching)

## Goals

- Consolidate all markdown-related regex patterns into one package
- Pre-compile patterns at package initialization
- Provide consistent matching API
- Preserve exact existing behavior (no parsing changes)
- Simplify future blackfriday migration

## Non-Goals

- Changing parsing behavior or output format
- Consolidating non-markdown regex (git URLs)
- Consolidating simple utility patterns (`\s+`, `\n{3,}`)
- Adding new parsing capabilities

## Decisions

### Decision 1: Package Structure

**What**: Create `internal/regex/` with files split by category.

**Structure**:

```
internal/regex/
├── doc.go           # Package documentation
├── headers.go       # H2, H3, H4 header patterns and matchers
├── headers_test.go
├── tasks.go         # Task checkbox and numbered task patterns
├── tasks_test.go
├── renames.go       # RENAMED section FROM/TO patterns
├── renames_test.go
├── sections.go      # Section content extraction helpers
└── sections_test.go
```

**Why**:

- Clear separation by semantic category
- Easier to navigate and maintain
- Test files co-located with source
- Each file stays focused and small (~50-80 lines)

### Decision 2: Pattern Organization

**What**: Split patterns into separate files by category. Export both raw patterns and helper functions.

**headers.go**:

```go
package regex

import "regexp"

// Header patterns - all pre-compiled at package init
var (
    // H2SectionHeader matches "## Section Name"
    H2SectionHeader = regexp.MustCompile(`^##\s+(.+)$`)

    // H2DeltaSection matches "## ADDED|MODIFIED|REMOVED|RENAMED Requirements"
    H2DeltaSection = regexp.MustCompile(`^##\s+(ADDED|MODIFIED|REMOVED|RENAMED)\s+Requirements\s*$`)

    // H2RequirementsSection matches exactly "## Requirements"
    H2RequirementsSection = regexp.MustCompile(`(?m)^##\s+Requirements\s*$`)

    // H2NextSection matches any "## " header (for finding section boundaries)
    H2NextSection = regexp.MustCompile(`(?m)^##\s+`)

    // H3Requirement matches "### Requirement: Name"
    H3Requirement = regexp.MustCompile(`^###\s+Requirement:\s*(.+)$`)

    // H3AnyHeader matches any "### " header
    H3AnyHeader = regexp.MustCompile(`^###\s+`)

    // H4Scenario matches "#### Scenario: Name"
    H4Scenario = regexp.MustCompile(`^####\s+Scenario:\s*(.+)$`)
)

// Helper functions for headers
func MatchH2SectionHeader(line string) (name string, ok bool) { ... }
func MatchH2DeltaSection(line string) (deltaType string, ok bool) { ... }
func MatchH3Requirement(line string) (name string, ok bool) { ... }
func MatchH4Scenario(line string) (name string, ok bool) { ... }
```

**tasks.go**:

```go
package regex

import "regexp"

var (
    // TaskCheckbox matches "- [ ]" or "- [x]" task items
    TaskCheckbox = regexp.MustCompile(`^\s*-\s*\[([xX ])\]`)

    // NumberedTask matches "- [ ] 1.1 Description" format
    NumberedTask = regexp.MustCompile(`^-\s+\[([ xX])\]\s+(\d+\.\d+)\s+(.+)$`)

    // NumberedSection matches "## 1. Section Name" format
    NumberedSection = regexp.MustCompile(`^##\s+\d+\.\s+(.+)$`)
)

func MatchTaskCheckbox(line string) (state rune, ok bool) { ... }
func MatchNumberedTask(line string) (checkbox, id, desc string, ok bool) { ... }
func MatchNumberedSection(line string) (name string, ok bool) { ... }
```

**renames.go**:

```go
package regex

import "regexp"

// Both backtick and non-backtick variants exported separately
var (
    // RenamedFrom matches "- FROM: `### Requirement: Name`" (with backticks)
    RenamedFrom = regexp.MustCompile(`^-\s*FROM:\s*` + "`" + `###\s+Requirement:\s*(.+?)` + "`" + `\s*$`)

    // RenamedTo matches "- TO: `### Requirement: Name`" (with backticks)
    RenamedTo = regexp.MustCompile(`^-\s*TO:\s*` + "`" + `###\s+Requirement:\s*(.+?)` + "`" + `\s*$`)

    // RenamedFromAlt matches "- FROM: ### Requirement: Name" (without backticks)
    RenamedFromAlt = regexp.MustCompile(`^\s*-\s*FROM:\s*###\s*Requirement:\s*(.+)$`)

    // RenamedToAlt matches "- TO: ### Requirement: Name" (without backticks)
    RenamedToAlt = regexp.MustCompile(`^\s*-\s*TO:\s*###\s*Requirement:\s*(.+)$`)
)

func MatchRenamedFrom(line string) (name string, ok bool) { ... }      // tries backtick version
func MatchRenamedFromAlt(line string) (name string, ok bool) { ... }   // tries non-backtick version
func MatchRenamedTo(line string) (name string, ok bool) { ... }
func MatchRenamedToAlt(line string) (name string, ok bool) { ... }
```

**Why**:

- Package-level `var` ensures single compilation
- Exported patterns allow direct use for complex matching needs
- Exported helpers provide clean API for common use cases
- Both RENAMED variants kept separate to preserve exact existing behavior
- Files organized by what they match, not by consumer

### Decision 3: Section Content Extraction

**What**: Provide `sections.go` with helpers for extracting content between markdown headers.

**sections.go**:

```go
package regex

import (
    "fmt"
    "regexp"
)

// FindSectionContent extracts content between a specific H2 section and the next H2.
// The sectionHeader parameter is the exact text after "## " (e.g., "Requirements").
// Returns empty string if section not found.
func FindSectionContent(content, sectionHeader string) string {
    pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^##\s+%s\s*$`, regexp.QuoteMeta(sectionHeader)))
    matches := pattern.FindStringIndex(content)
    if matches == nil {
        return ""
    }

    sectionStart := matches[1]
    nextMatches := H2NextSection.FindStringIndex(content[sectionStart:])

    if nextMatches != nil {
        return content[sectionStart : sectionStart+nextMatches[0]]
    }
    return content[sectionStart:]
}

// FindDeltaSectionContent extracts content from a delta section (ADDED, MODIFIED, etc.).
// Convenience wrapper around FindSectionContent for delta specs.
func FindDeltaSectionContent(content, deltaType string) string {
    return FindSectionContent(content, deltaType+" Requirements")
}

// FindRequirementsSection extracts the "## Requirements" section content.
func FindRequirementsSection(content string) string {
    return FindSectionContent(content, "Requirements")
}
```

**Why**:

- Common operation duplicated in 3+ files (validation, archive, parsers)
- Consolidates section boundary detection logic
- Specialized helpers reduce boilerplate for common cases
- Uses pre-compiled H2NextSection pattern for efficiency

### Decision 4: Helper Function Return Style

**What**: All helper functions use `(value, ok bool)` return pattern.

```go
// Single value extraction
func MatchH3Requirement(line string) (name string, ok bool)
func MatchH4Scenario(line string) (name string, ok bool)
func MatchH2DeltaSection(line string) (deltaType string, ok bool)

// Multi-value extraction
func MatchNumberedTask(line string) (checkbox, id, desc string, ok bool)

// Special case: rune for checkbox state
func MatchTaskCheckbox(line string) (state rune, ok bool)
```

**Why**:

- Idiomatic Go pattern - consistent with stdlib (e.g., map lookup)
- Clear success/failure without nil checks
- Callers use `if name, ok := regex.MatchH3Requirement(line); ok { ... }`
- Values are trimmed/normalized before return

### Decision 5: Migration Strategy

**What**: Migrate files one at a time, running tests after each.

**Order**:

1. Create `internal/regex/` package with split files (headers.go, tasks.go, renames.go, sections.go)
2. Migrate `internal/parsers/parsers.go` (simplest, 3 patterns)
3. Migrate `internal/parsers/requirement_parser.go` (2 patterns)
4. Migrate `internal/parsers/delta_parser.go` (most complex, 10 patterns)
5. Migrate `internal/validation/parser.go` (5 patterns)
6. Migrate `internal/archive/spec_merger.go` (3 patterns)
7. Migrate `cmd/accept.go` (2 patterns)
8. Remove any now-unused local patterns

**Why**:

- Incremental migration reduces risk
- Tests validate each step
- Order goes from simple to complex

### Decision 6: What NOT to Consolidate

**What**: Leave these patterns inline:

1. **`\s+`** (whitespace normalization) - Too generic, used everywhere
2. **`\n{3,}`** (newline collapsing) - Single use, utility purpose
3. **`(?i)\b(shall|must)\b`** (normative language) - Validation-specific, not structural parsing
4. **Git URL patterns** in `internal/git/platform.go` - Not markdown-related

**Why**:

- These are utilities, not markdown structure patterns
- Consolidating would create unnecessary dependencies
- Single-use patterns don't benefit from sharing

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Hidden behavioral differences | Medium - subtle bugs | Comprehensive test coverage; diff existing vs new output |
| Import cycles | Low | `internal/regex/` has no internal dependencies |
| Over-abstraction | Low | Keep simple; only patterns, not scanning logic |

## Migration Plan

1. Create regex package with all patterns
2. Add unit tests for each pattern
3. Migrate consumers one file at a time
4. Verify full test suite passes after each migration
5. Remove unused local patterns
6. Final verification with `go test ./...` and `golangci-lint run`

**Rollback**: Revert to commit before change; no external state involved.

## Open Questions

None - design is straightforward consolidation without behavior changes.
