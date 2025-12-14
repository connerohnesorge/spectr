// Package markdown provides AST-based markdown parsing using blackfriday v2.
//
// This package provides a single entry point for parsing markdown content
// and extracting structural elements needed by Spectr:
//   - Headers (H1-H4) with source line numbers
//   - Sections (content between headers) as raw markdown text
//   - Tasks (checkbox items) with hierarchical structure
//
// # Usage
//
// Parse a document once and extract all needed information:
//
//	doc, err := markdown.ParseDocument(content)
//	if err != nil {
//	    // Handle error
//	}
//	// Use doc.Headers, doc.Sections, doc.Tasks as needed
//
// # Design Principles
//
//   - Parse once, never reparse: The architecture enforces single parse
//   - Type-safe AST traversal: All blackfriday internals are hidden
//   - Line numbers: All types include source line numbers for error messages
//   - Raw content: Section content preserves original markdown formatting
package markdown
