// Package regex provides pre-compiled regular expression patterns for parsing
// markdown specification files.
//
// This package consolidates all markdown-related regex patterns used across
// Spectr for parsing:
//   - H2 section headers (## Requirements, ## ADDED Requirements, etc.)
//   - H3 requirement headers (### Requirement: Name)
//   - H4 scenario headers (#### Scenario: Name)
//   - Task checkboxes (- [ ], - [x])
//   - Delta operations (ADDED, MODIFIED, REMOVED, RENAMED sections)
//   - RENAMED requirement FROM/TO pairs
//
// All patterns are pre-compiled at package initialization using
// regexp.MustCompile, ensuring single compilation and efficient
// matching throughout the application.
//
// # Pattern Organization
//
// Patterns are organized by category:
//   - headers.go: H2, H3, H4 header patterns and matchers
//   - tasks.go: Task checkbox and numbered task patterns
//   - renames.go: RENAMED section FROM/TO patterns
//   - sections.go: Section content extraction helpers
//
// # Helper Functions
//
// Each category provides helper functions that return (value, ok bool) pairs
// for idiomatic Go usage:
//
//	if name, ok := regex.MatchH3Requirement(line); ok {
//	    // name contains the extracted requirement name
//	}
//
// # Migration Note
//
// This package is a prerequisite for the replace-regex-with-blackfriday change.
// After consolidation, blackfriday migration can replace this package with
// internal/markdown/ for proper AST-based parsing.
package regex
