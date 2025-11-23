package parser

import (
	"strings"
	"testing"
)

// BenchmarkLexer_SmallDocument benchmarks lexing a small document.
func BenchmarkLexer_SmallDocument(b *testing.B) {
	input := `# Header

Some text here.

## Subheader

- List item 1
- List item 2

` + "```go\nfunc main() {}\n```"

	b.ResetTimer()
	for range b.N {
		l := NewLexer(input)
		_ = l.Lex()
	}
}

// BenchmarkLexer_MediumDocument benchmarks lexing a medium-sized document.
func BenchmarkLexer_MediumDocument(b *testing.B) {
	var sb strings.Builder

	sb.WriteString("# Main Title\n\n")
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 50; i++ {
		sb.WriteString("## Section ")
		sb.WriteString(strings.Repeat("A", 10))
		sb.WriteString("\n\nSome content.\n\n")
		sb.WriteString("```go\ncode\n```\n\n")
		sb.WriteString("- Item\n")
	}

	input := sb.String()

	b.ResetTimer()
	for range b.N {
		l := NewLexer(input)
		_ = l.Lex()
	}
}

// BenchmarkParser_SmallDocument benchmarks parsing a small document.
func BenchmarkParser_SmallDocument(b *testing.B) {
	input := `# Header

Some text here.

## Subheader

- List item 1
- List item 2

` + "```go\nfunc main() {}\n```"

	b.ResetTimer()
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input)
	}
}

// BenchmarkParser_MediumDocument benchmarks parsing a medium-sized document.
func BenchmarkParser_MediumDocument(b *testing.B) {
	var sb strings.Builder

	sb.WriteString("# Main Title\n\n")
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 50; i++ {
		sb.WriteString("## Section ")
		sb.WriteString(strings.Repeat("A", 10))
		sb.WriteString("\n\nSome content.\n\n")
		sb.WriteString("```go\ncode\n```\n\n")
		sb.WriteString("- Item\n")
	}

	input := sb.String()

	b.ResetTimer()
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input)
	}
}

// BenchmarkExtractor_Requirements benchmarks requirement extraction.
func BenchmarkExtractor_Requirements(b *testing.B) {
	var sb strings.Builder

	sb.WriteString("## Requirements\n\n")
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 20; i++ {
		sb.WriteString("### Requirement: Feature\n")
		sb.WriteString("The system SHALL provide functionality.\n\n")
		sb.WriteString("#### Scenario: Test\n")
		sb.WriteString("- **WHEN** action\n")
		sb.WriteString("- **THEN** result\n\n")
	}

	input := sb.String()
	doc, _ := Parse(input)

	b.ResetTimer()
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < b.N; i++ {
		_, _ = ExtractRequirements(doc)
	}
}

// BenchmarkExtractor_Deltas benchmarks delta extraction.
func BenchmarkExtractor_Deltas(b *testing.B) {
	var sb strings.Builder

	sb.WriteString("## ADDED Requirements\n\n")
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 10; i++ {
		sb.WriteString("### Requirement: New Feature\n")
		sb.WriteString("The system SHALL provide new feature.\n\n")
	}

	sb.WriteString("## MODIFIED Requirements\n\n")
	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < 10; i++ {
		sb.WriteString("### Requirement: Updated Feature\n")
		sb.WriteString("The system SHALL update feature.\n\n")
	}

	input := sb.String()
	doc, _ := Parse(input)

	b.ResetTimer()
	for range b.N {
		_, _ = ExtractDeltas(doc)
	}
}

// BenchmarkEndToEnd benchmarks the complete pipeline.
func BenchmarkEndToEnd(b *testing.B) {
	input := `# Specification

## Requirements

### Requirement: User Authentication
The system SHALL authenticate users.

#### Scenario: Valid login
- **WHEN** valid credentials provided
- **THEN** user is authenticated

### Requirement: Session Management
The system SHALL manage sessions.

#### Scenario: Session creation
- **WHEN** user logs in
- **THEN** session is created
`

	b.ResetTimer()
	for range b.N {
		doc, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
		_, err = ExtractRequirements(doc)
		if err != nil {
			b.Fatal(err)
		}
	}
}
