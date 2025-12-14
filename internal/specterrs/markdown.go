package specterrs

import "fmt"

// MarkdownParseError indicates markdown content failed to parse.
type MarkdownParseError struct {
	Path string // File path if known, empty otherwise
	Line int    // Line number if known, 0 otherwise
	Err  error  // Underlying error
}

func (e *MarkdownParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf(
			"failed to parse markdown %s at line %d: %v",
			e.Path,
			e.Line,
			e.Err,
		)
	}

	if e.Path != "" {
		return fmt.Sprintf(
			"failed to parse markdown %s: %v",
			e.Path,
			e.Err,
		)
	}

	return fmt.Sprintf("failed to parse markdown: %v", e.Err)
}

func (e *MarkdownParseError) Unwrap() error {
	return e.Err
}

// EmptyContentError indicates empty or whitespace-only content was provided.
type EmptyContentError struct {
	Path string
}

func (e *EmptyContentError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("markdown file is empty: %s", e.Path)
	}

	return "markdown content is empty"
}

// BinaryContentError indicates binary (non-text) content was provided.
type BinaryContentError struct {
	Path string
}

func (e *BinaryContentError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf(
			"file appears to be binary, not markdown: %s",
			e.Path,
		)
	}

	return "content appears to be binary, not markdown"
}
