//nolint:revive // file-length-limit: parser requires comprehensive markdown handling
package markdown

import (
	"bytes"
	"strings"
	"sync"
	"unicode"
)

// DefaultMaxErrors is the maximum number of parse errors before aborting.
const DefaultMaxErrors = 100

// ParseError represents an error encountered during parsing.
// It contains the byte offset where the error occurred, a human-readable message,
// and optionally a list of expected token types.
type ParseError struct {
	Offset   int         // Byte offset where error occurred
	Message  string      // Human-readable error description
	Expected []TokenType // What tokens would have been valid (may be nil)
}

// Error implements the error interface.
func (e ParseError) Error() string {
	if e.Offset >= 0 {
		return "offset " + itoa(
			e.Offset,
		) + ": " + e.Message
	}

	return e.Message
}

// Position converts the byte offset to a Position using the provided LineIndex.
func (e ParseError) Position(
	idx *LineIndex,
) Position {
	return idx.PositionAt(e.Offset)
}

// itoa converts an integer to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}

// linkDefinition stores a collected link definition.
type linkDefinition struct {
	url   []byte
	title []byte
}

// parser holds the internal state during a single parse operation.
// This struct is NOT exported - the public API is the stateless Parse function.
type parser struct {
	source      []byte
	tokens      []Token
	pos         int // Current token position
	errors      []ParseError
	maxErrors   int
	linkDefs    map[string]linkDefinition // Case-insensitive label -> definition
	lineIndex   *LineIndex
	inlineState *inlineParser
}

// delimiter represents an emphasis delimiter on the stack.
type delimiter struct {
	token     Token     // The delimiter token
	count     int       // Number of delimiter characters
	canOpen   bool      // Can this delimiter open emphasis?
	canClose  bool      // Can this delimiter close emphasis?
	active    bool      // Is this delimiter still active?
	textStart int       // Start position of text after this delimiter
	delimType TokenType // TokenAsterisk or TokenUnderscore
}

// inlineParser handles inline content parsing with delimiter stack.
type inlineParser struct {
	source     []byte
	tokens     []Token
	pos        int
	start      int // Start offset of inline content
	end        int // End offset of inline content
	delimiters []delimiter
	linkDefs   map[string]linkDefinition
	errors     *[]ParseError
}

// Object pools for parser internals
var (
	parserPool = sync.Pool{
		New: func() interface{} {
			return &parser{
				linkDefs: make(
					map[string]linkDefinition,
				),
				errors: make(
					[]ParseError,
					0,
					8,
				),
			}
		},
	}

	tokenSlicePool = sync.Pool{
		New: func() interface{} {
			s := make([]Token, 0, 256)

			return &s
		},
	}

	nodeSlicePool = sync.Pool{
		New: func() interface{} {
			s := make([]Node, 0, 16)

			return &s
		},
	}
)

// Parse transforms source bytes into an immutable AST.
// It returns the root document node and any parse errors encountered.
// This function is stateless and safe for concurrent calls.
//
//nolint:revive // function-length: parse entry point requires setup/teardown
func Parse(source []byte) (Node, []ParseError) {
	// Get parser from pool
	p, ok := parserPool.Get().(*parser)
	if !ok {
		p = &parser{
			linkDefs: make(
				map[string]linkDefinition,
			),
		}
	}
	defer func() {
		// Clear and return to pool
		p.source = nil
		p.tokens = nil
		p.pos = 0
		p.errors = p.errors[:0]
		for k := range p.linkDefs {
			delete(p.linkDefs, k)
		}
		p.lineIndex = nil
		p.inlineState = nil
		parserPool.Put(p)
	}()

	// Initialize parser state
	p.source = source
	p.maxErrors = DefaultMaxErrors
	p.lineIndex = NewLineIndex(source)

	// Tokenize
	lex := newLexer(source)
	tokensPtr, ok := tokenSlicePool.Get().(*[]Token)
	if !ok {
		slice := make([]Token, 0, 256)
		tokensPtr = &slice
	}
	tokens := (*tokensPtr)[:0]
	for {
		tok := lex.Next()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	p.tokens = tokens
	defer func() {
		*tokensPtr = tokens[:0]
		tokenSlicePool.Put(tokensPtr)
	}()

	// First pass: collect link definitions
	p.collectLinkDefinitions()

	// Second pass: parse document structure
	p.pos = 0
	doc := p.parseDocument()

	// Copy errors before returning parser to pool
	var errors []ParseError
	if len(p.errors) > 0 {
		errors = make([]ParseError, len(p.errors))
		copy(errors, p.errors)
	}

	return doc, errors
}

// addError adds a parse error and returns true if parsing should continue.
//
//nolint:unused // Utility function for error recovery in future enhancements
func (p *parser) addError(
	offset int,
	message string,
	expected ...TokenType,
) bool {
	p.errors = append(p.errors, ParseError{
		Offset:   offset,
		Message:  message,
		Expected: expected,
	})

	return len(p.errors) < p.maxErrors
}

// current returns the current token without advancing.
func (p *parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{
			Type:  TokenEOF,
			Start: len(p.source),
			End:   len(p.source),
		}
	}

	return p.tokens[p.pos]
}

// peek returns the token at offset from current position without advancing.
func (p *parser) peek(offset int) Token {
	idx := p.pos + offset
	if idx < 0 || idx >= len(p.tokens) {
		return Token{
			Type:  TokenEOF,
			Start: len(p.source),
			End:   len(p.source),
		}
	}

	return p.tokens[idx]
}

// advance moves to the next token and returns the previous current token.
func (p *parser) advance() Token {
	tok := p.current()
	if p.pos < len(p.tokens) {
		p.pos++
	}

	return tok
}

// skipWhitespace skips whitespace tokens (not newlines).
func (p *parser) skipWhitespace() {
	for p.current().Type == TokenWhitespace {
		p.advance()
	}
}

// atLineStart returns true if we're at the start of a line.
func (p *parser) atLineStart() bool {
	if p.pos == 0 {
		return true
	}
	// Check if previous non-whitespace token was a newline
	for i := p.pos - 1; i >= 0; i-- {
		if p.tokens[i].Type == TokenNewline {
			return true
		}
		if p.tokens[i].Type != TokenWhitespace {
			return false
		}
	}

	return true
}

// skipToSyncPoint advances to the next synchronization point for error recovery.
// Sync points are: blank line (two newlines), header start, list marker.
//
//nolint:unused // reserved for future error recovery implementation
func (p *parser) skipToSyncPoint() {
	prevNewline := false
	for p.current().Type != TokenEOF {
		tok := p.current()

		// Blank line: two consecutive newlines
		if tok.Type == TokenNewline {
			if prevNewline {
				p.advance() // Skip second newline

				return
			}
			prevNewline = true
			p.advance()

			continue
		}

		// After newline, check for sync points
		if prevNewline {
			// Header start
			if tok.Type == TokenHash {
				return
			}
			// List marker
			if tok.Type == TokenDash ||
				tok.Type == TokenPlus ||
				tok.Type == TokenAsterisk ||
				tok.Type == TokenNumber {
				return
			}
			// Blockquote
			if tok.Type == TokenGreaterThan {
				return
			}
			// Code fence
			if tok.Type == TokenBacktick ||
				tok.Type == TokenTilde {
				return
			}
		}

		prevNewline = false
		p.advance()
	}
}

// collectLinkDefinitions performs first pass to collect all link definitions.
// Link definitions have the format: [label]: url "optional title"
func (p *parser) collectLinkDefinitions() {
	p.pos = 0
	for p.current().Type != TokenEOF {
		// Link definitions must start at beginning of line
		if !p.atLineStart() {
			p.advance()

			continue
		}

		p.skipWhitespace()

		// Check for [label]:
		if p.current().Type != TokenBracketOpen {
			p.skipToNextLine()

			continue
		}

		// Try to parse link definition
		startPos := p.pos
		if !p.tryParseLinkDefinition() {
			// Not a link definition, restore position and skip line
			p.pos = startPos
			p.skipToNextLine()
		}
	}
}

// tryParseLinkDefinition attempts to parse a link definition.
// Returns true if successful, false otherwise.
//
//nolint:revive // function-length: link definition parsing is inherently complex
func (p *parser) tryParseLinkDefinition() bool {
	// [
	if p.current().Type != TokenBracketOpen {
		return false
	}
	p.advance()

	// Collect label
	var labelParts [][]byte
	for p.current().Type != TokenBracketClose && p.current().Type != TokenEOF &&
		p.current().Type != TokenNewline {
		labelParts = append(
			labelParts,
			p.current().Source,
		)
		p.advance()
	}

	// ]
	if p.current().Type != TokenBracketClose {
		return false
	}
	p.advance()

	// :
	if p.current().Type != TokenColon {
		return false
	}
	p.advance()

	// Optional whitespace
	p.skipWhitespace()

	// URL
	var urlParts [][]byte
	for p.current().Type != TokenEOF && p.current().Type != TokenNewline &&
		p.current().Type != TokenWhitespace {
		urlParts = append(
			urlParts,
			p.current().Source,
		)
		p.advance()
	}

	if len(urlParts) == 0 {
		return false
	}

	url := bytes.Join(urlParts, nil)

	// Optional title
	var title []byte
	p.skipWhitespace()
	if p.current().Type == TokenText {
		text := string(p.current().Source)
		if len(text) >= 2 &&
			(text[0] == '"' || text[0] == '\'') {
			title = p.current().Source[1 : len(p.current().Source)-1]
			p.advance()
		}
	}

	// Build label (case-insensitive)
	label := strings.ToLower(
		string(bytes.Join(labelParts, nil)),
	)
	label = strings.TrimSpace(label)

	// Store definition (first definition wins)
	if _, exists := p.linkDefs[label]; !exists {
		p.linkDefs[label] = linkDefinition{
			url:   url,
			title: title,
		}
	}

	return true
}

// skipToNextLine advances to the start of the next line.
func (p *parser) skipToNextLine() {
	for p.current().Type != TokenEOF {
		if p.current().Type == TokenNewline {
			p.advance()

			return
		}
		p.advance()
	}
}

// parseDocument parses the entire document and returns the root node.
func (p *parser) parseDocument() Node {
	startOffset := 0
	if len(p.tokens) > 0 {
		startOffset = p.tokens[0].Start
	}

	nodesPtr, ok := nodeSlicePool.Get().(*[]Node)
	if !ok {
		slice := make([]Node, 0, 32)
		nodesPtr = &slice
	}
	children := (*nodesPtr)[:0]
	defer func() {
		*nodesPtr = children[:0]
		nodeSlicePool.Put(nodesPtr)
	}()

	for p.current().Type != TokenEOF {
		// Skip blank lines
		if p.current().Type == TokenNewline {
			p.advance()

			continue
		}

		node := p.parseBlock()
		if node != nil {
			children = append(children, node)
		}

		// Check for too many errors
		if len(p.errors) >= p.maxErrors {
			break
		}
	}

	endOffset := len(p.source)

	// Build document node
	childrenCopy := make([]Node, len(children))
	copy(childrenCopy, children)

	return NewNodeBuilder(NodeTypeDocument).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(p.source).
		WithChildren(childrenCopy).
		Build()
}

// parseBlock parses a single block-level element.
// Block detection order: code fence, header, blockquote, list item, paragraph
func (p *parser) parseBlock() Node {
	p.skipWhitespace()

	tok := p.current()
	if tok.Type == TokenEOF {
		return nil
	}

	// Check for code fence (3+ backticks or tildes at line start)
	if tok.Type == TokenBacktick ||
		tok.Type == TokenTilde {
		if node := p.tryParseCodeFence(); node != nil {
			return node
		}
	}

	// Check for header (1-6 # at line start)
	if tok.Type == TokenHash {
		if node := p.parseHeader(); node != nil {
			return node
		}
	}

	// Check for blockquote (> at line start)
	if tok.Type == TokenGreaterThan {
		return p.parseBlockquote()
	}

	// Check for list item
	if tok.Type == TokenDash ||
		tok.Type == TokenPlus ||
		tok.Type == TokenAsterisk {
		return p.parseList(false)
	}
	if tok.Type == TokenNumber {
		next := p.peek(1)
		if next.Type == TokenDot {
			return p.parseList(true)
		}
	}

	// Default: paragraph
	return p.parseParagraph()
}

// tryParseCodeFence attempts to parse a fenced code block.
// Returns nil if not a valid code fence.
//
//nolint:revive // function-length: code fence parsing is inherently complex
func (p *parser) tryParseCodeFence() Node {
	startPos := p.pos
	startOffset := p.current().Start

	// Count fence characters
	fenceChar := p.current().Type
	fenceCount := 0
	for p.current().Type == fenceChar {
		fenceCount++
		p.advance()
	}

	// Need at least 3 fence characters
	if fenceCount < 3 {
		p.pos = startPos

		return nil
	}

	// Get optional language identifier
	p.skipWhitespace()
	var language []byte
	if p.current().Type == TokenText {
		language = p.current().Source
		p.advance()
	}

	// Skip to end of line
	p.skipToNextLine()

	// Collect content until closing fence
	var contentParts [][]byte

	for p.current().Type != TokenEOF {
		lineStart := p.pos

		// Check for closing fence
		closingCount := 0
		for p.current().Type == fenceChar {
			closingCount++
			p.advance()
		}

		if closingCount >= fenceCount {
			// Verify rest of line is whitespace/empty
			p.skipWhitespace()
			if p.current().Type == TokenNewline ||
				p.current().Type == TokenEOF {
				if p.current().Type == TokenNewline {
					p.advance()
				}

				break
			}
			// Not a valid closing fence, treat as content
			p.pos = lineStart
		} else if closingCount > 0 {
			// Partial fence, treat as content
			p.pos = lineStart
		}

		// Collect line content
		lineContent := p.collectLineContent()
		if len(lineContent) > 0 ||
			p.current().Type != TokenEOF {
			contentParts = append(
				contentParts,
				lineContent,
				[]byte{'\n'},
			)
		}

		if p.current().Type == TokenNewline {
			p.advance()
		}
	}

	endOffset := p.current().Start
	if p.pos > 0 {
		endOffset = p.tokens[p.pos-1].End
	}

	// Build content
	var content []byte
	if len(contentParts) > 0 {
		content = bytes.Join(contentParts, nil)
		// Trim trailing newline
		if len(content) > 0 &&
			content[len(content)-1] == '\n' {
			content = content[:len(content)-1]
		}
	}

	return NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(p.source[startOffset:endOffset]).
		WithLanguage(language).
		WithContent(content).
		Build()
}

// collectLineContent collects all tokens on the current line as content.
func (p *parser) collectLineContent() []byte {
	var parts [][]byte
	for p.current().Type != TokenNewline && p.current().Type != TokenEOF {
		parts = append(parts, p.current().Source)
		p.advance()
	}

	return bytes.Join(parts, nil)
}

// parseHeader parses an ATX-style header.
//
//nolint:revive // function-length: header parsing handles multiple cases
func (p *parser) parseHeader() Node {
	startOffset := p.current().Start

	// Count # characters
	level := 0
	for p.current().Type == TokenHash && level < 6 {
		level++
		p.advance()
	}

	// Must be followed by whitespace (or end of line for empty header)
	if p.current().Type != TokenWhitespace &&
		p.current().Type != TokenNewline &&
		p.current().Type != TokenEOF {
		// Not a valid header, treat as paragraph
		p.pos -= level

		return p.parseParagraph()
	}

	p.skipWhitespace()

	// Collect header content
	var titleParts [][]byte
	contentStart := p.pos
	for p.current().Type != TokenNewline && p.current().Type != TokenEOF {
		titleParts = append(
			titleParts,
			p.current().Source,
		)
		p.advance()
	}

	title := bytes.TrimSpace(
		bytes.Join(titleParts, nil),
	)

	// Skip trailing newline
	if p.current().Type == TokenNewline {
		p.advance()
	}

	endOffset := p.current().Start
	if p.pos > 0 && p.pos <= len(p.tokens) {
		endOffset = p.tokens[p.pos-1].End
	}

	// Check for Spectr-specific headers
	titleStr := string(title)

	// Check for Requirement: header (level 3)
	if level == 3 &&
		strings.HasPrefix(
			titleStr,
			"Requirement:",
		) {
		name := strings.TrimSpace(
			strings.TrimPrefix(
				titleStr,
				"Requirement:",
			),
		)

		return NewNodeBuilder(
			NodeTypeRequirement,
		).
			WithStart(startOffset).
			WithEnd(endOffset).
			WithSource(p.source[startOffset:endOffset]).
			WithName(name).
			Build()
	}

	// Check for Scenario: header (level 4)
	if level == 4 &&
		strings.HasPrefix(titleStr, "Scenario:") {
		name := strings.TrimSpace(
			strings.TrimPrefix(
				titleStr,
				"Scenario:",
			),
		)

		return NewNodeBuilder(NodeTypeScenario).
			WithStart(startOffset).
			WithEnd(endOffset).
			WithSource(p.source[startOffset:endOffset]).
			WithName(name).
			Build()
	}

	// Check for delta section headers (level 2)
	var deltaType string
	if level == 2 {
		deltaType = detectDeltaType(titleStr)
	}

	// Parse inline content for title
	contentEnd := contentStart
	for i := contentStart; i < p.pos; i++ {
		if p.tokens[i].Type == TokenNewline {
			break
		}
		contentEnd = i + 1
	}
	titleChildren := p.parseInlineContent(
		contentStart,
		contentEnd,
	)

	return NewNodeBuilder(NodeTypeSection).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(p.source[startOffset:endOffset]).
		WithLevel(level).
		WithTitle(title).
		WithDeltaType(deltaType).
		WithChildren(titleChildren).
		Build()
}

// detectDeltaType checks if the header text indicates a delta section.
// Returns "ADDED", "MODIFIED", "REMOVED", "RENAMED", or empty string.
func detectDeltaType(title string) string {
	title = strings.TrimSpace(title)
	upper := strings.ToUpper(title)

	if strings.HasPrefix(upper, "ADDED") {
		return "ADDED"
	}
	if strings.HasPrefix(upper, "MODIFIED") {
		return "MODIFIED"
	}
	if strings.HasPrefix(upper, "REMOVED") {
		return "REMOVED"
	}
	if strings.HasPrefix(upper, "RENAMED") {
		return "RENAMED"
	}

	return ""
}

// parseBlockquote parses a blockquote (lines starting with >).
//
//nolint:revive // function-length: blockquote parsing requires multiple passes
func (p *parser) parseBlockquote() Node {
	startOffset := p.current().Start

	nodesPtr, ok := nodeSlicePool.Get().(*[]Node)
	if !ok {
		slice := make([]Node, 0, 32)
		nodesPtr = &slice
	}
	children := (*nodesPtr)[:0]
	defer func() {
		*nodesPtr = children[:0]
		nodeSlicePool.Put(nodesPtr)
	}()

	for p.current().Type != TokenEOF {
		// Skip >
		if p.current().Type != TokenGreaterThan {
			break
		}
		p.advance()

		// Optional space after >
		if p.current().Type == TokenWhitespace {
			p.advance()
		}

		// Parse the content as a block
		if p.current().Type != TokenNewline &&
			p.current().Type != TokenEOF {
			block := p.parseBlock()
			if block != nil {
				children = append(children, block)
			}
		} else if p.current().Type == TokenNewline {
			p.advance()
		}

		// Check if next line continues the blockquote
		p.skipWhitespace()
		if p.current().Type != TokenGreaterThan {
			break
		}
	}

	endOffset := p.current().Start
	if p.pos > 0 {
		endOffset = p.tokens[p.pos-1].End
	}

	childrenCopy := make([]Node, len(children))
	copy(childrenCopy, children)

	return NewNodeBuilder(NodeTypeBlockquote).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(p.source[startOffset:endOffset]).
		WithChildren(childrenCopy).
		Build()
}

// parseList parses a list (ordered or unordered).
//
//nolint:revive // function-length: list parsing handles nested structures
func (p *parser) parseList(ordered bool) Node {
	startOffset := p.current().Start

	nodesPtr, ok := nodeSlicePool.Get().(*[]Node)
	if !ok {
		slice := make([]Node, 0, 32)
		nodesPtr = &slice
	}
	items := (*nodesPtr)[:0]
	defer func() {
		*nodesPtr = items[:0]
		nodeSlicePool.Put(nodesPtr)
	}()

	// Track indentation level
	baseIndent := p.countLeadingWhitespace()

	for p.current().Type != TokenEOF {
		// Check indentation
		currentIndent := p.countLeadingWhitespace()
		p.skipWhitespace()

		// Check for list marker
		isListItem := false
		if ordered {
			if p.current().Type == TokenNumber {
				next := p.peek(1)
				if next.Type == TokenDot {
					isListItem = true
				}
			}
		} else {
			if p.current().Type == TokenDash || p.current().Type == TokenPlus ||
				p.current().Type == TokenAsterisk {
				isListItem = true
			}
		}

		if !isListItem {
			break
		}

		// Nested list check
		if currentIndent > baseIndent {
			// This is a nested list item - parse as nested list
			nestedList := p.parseList(ordered)
			if nestedList != nil &&
				len(items) > 0 {
				// Add nested list to last item's children
				// For now, just add as separate item
				items = append(items, nestedList)
			}

			continue
		}

		item := p.parseListItem(ordered)
		if item != nil {
			items = append(items, item)
		}

		// Skip blank lines between items
		for p.current().Type == TokenNewline {
			p.advance()
			// Check if next line is blank (paragraph break)
			if p.current().Type == TokenNewline {
				p.advance()

				break
			}
		}
	}

	endOffset := p.current().Start
	if p.pos > 0 {
		endOffset = p.tokens[p.pos-1].End
	}

	itemsCopy := make([]Node, len(items))
	copy(itemsCopy, items)

	return NewNodeBuilder(NodeTypeList).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(p.source[startOffset:endOffset]).
		WithOrdered(ordered).
		WithChildren(itemsCopy).
		Build()
}

// countLeadingWhitespace returns the number of whitespace characters at current position.
func (p *parser) countLeadingWhitespace() int {
	count := 0
	savedPos := p.pos
	for p.current().Type == TokenWhitespace {
		count += len(p.current().Source)
		p.advance()
	}
	p.pos = savedPos

	return count
}

// parseListItem parses a single list item.
//
//nolint:revive // function-length: list item parsing handles checkboxes/keywords
func (p *parser) parseListItem(
	ordered bool,
) Node {
	startOffset := p.current().Start

	// Skip list marker
	if ordered {
		p.advance() // number
		p.advance() // dot
	} else {
		p.advance() // -, +, or *
	}

	// Optional space after marker
	p.skipWhitespace()

	// Check for checkbox [ ] or [x]
	var checked *bool
	if p.current().Type == TokenBracketOpen {
		next1 := p.peek(1)
		next2 := p.peek(2)
		if next2.Type == TokenBracketClose {
			switch next1.Type { //nolint:exhaustive // Only care about checkbox tokens
			case TokenWhitespace:
				// [ ] - unchecked
				f := false
				checked = &f
				p.advance() // [
				p.advance() // space
				p.advance() // ]
				p.skipWhitespace()
			case TokenX:
				// [x] - checked
				t := true
				checked = &t
				p.advance() // [
				p.advance() // x
				p.advance() // ]
				p.skipWhitespace()
			default:
				// Not a valid checkbox pattern
			}
		}
	}

	// Collect item content
	contentStart := p.pos
	for p.current().Type != TokenNewline && p.current().Type != TokenEOF {
		p.advance()
	}
	contentEnd := p.pos

	// Parse inline content
	children := p.parseInlineContent(
		contentStart,
		contentEnd,
	)

	// Detect WHEN/THEN/AND keyword
	keyword := p.detectKeyword(
		contentStart,
		contentEnd,
	)

	// Skip trailing newline
	if p.current().Type == TokenNewline {
		p.advance()
	}

	endOffset := p.current().Start
	if p.pos > 0 {
		endOffset = p.tokens[p.pos-1].End
	}

	return NewNodeBuilder(NodeTypeListItem).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(p.source[startOffset:endOffset]).
		WithChecked(checked).
		WithKeyword(keyword).
		WithChildren(children).
		Build()
}

// detectKeyword checks if the list item content starts with **WHEN**, **THEN**, or **AND**.
func (p *parser) detectKeyword(
	start, end int,
) string {
	if start >= end || start >= len(p.tokens) {
		return ""
	}

	// Look for **WHEN**, **THEN**, or **AND** pattern
	pos := start
	for pos < end && pos < len(p.tokens) {
		tok := p.tokens[pos]
		if tok.Type == TokenWhitespace {
			pos++

			continue
		}

		// Check for ** opening
		if tok.Type == TokenAsterisk {
			if pos+1 < end &&
				p.tokens[pos+1].Type == TokenAsterisk {
				// Found **, now check for keyword
				if pos+2 < end &&
					p.tokens[pos+2].Type == TokenText {
					text := strings.ToUpper(
						string(
							p.tokens[pos+2].Source,
						),
					)
					if text == "WHEN" ||
						text == "THEN" ||
						text == "AND" {
						// Check for closing **
						if pos+3 < end &&
							p.tokens[pos+3].Type == TokenAsterisk &&
							pos+4 < end &&
							p.tokens[pos+4].Type == TokenAsterisk {
							return text
						}
					}
				}
			}
		}

		break
	}

	return ""
}

// parseParagraph parses a paragraph (consecutive non-blank lines).
//
//nolint:revive // function-length: paragraph parsing handles block transitions
func (p *parser) parseParagraph() Node {
	startOffset := p.current().Start
	contentStart := p.pos

	// Collect content until blank line or block element
	for p.current().Type != TokenEOF {
		tok := p.current()

		// Check for blank line (paragraph terminator)
		if tok.Type == TokenNewline {
			next := p.peek(1)
			if next.Type == TokenNewline ||
				next.Type == TokenEOF {
				p.advance() // consume first newline

				break
			}
			// Check if next line starts a block element
			p.advance() // consume newline
			p.skipWhitespace()
			nextTok := p.current()
			if nextTok.Type == TokenHash ||
				nextTok.Type == TokenGreaterThan ||
				nextTok.Type == TokenBacktick ||
				nextTok.Type == TokenTilde {
				break
			}
			if nextTok.Type == TokenDash ||
				nextTok.Type == TokenPlus {
				// Could be list marker
				break
			}
			if nextTok.Type == TokenNumber {
				next := p.peek(1)
				if next.Type == TokenDot {
					break
				}
			}

			continue
		}

		p.advance()
	}

	contentEnd := p.pos
	endOffset := p.current().Start
	if p.pos > 0 &&
		p.tokens[p.pos-1].Type != TokenNewline {
		endOffset = p.tokens[p.pos-1].End
	} else if p.pos > 1 {
		endOffset = p.tokens[p.pos-1].End
	}

	// Parse inline content
	children := p.parseInlineContent(
		contentStart,
		contentEnd,
	)

	return NewNodeBuilder(NodeTypeParagraph).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(p.source[startOffset:endOffset]).
		WithChildren(children).
		Build()
}

// parseInlineContent parses inline content from the token range [start, end).
func (p *parser) parseInlineContent(
	start, end int,
) []Node {
	if start >= end || start >= len(p.tokens) {
		return nil
	}

	// Extract tokens for inline parsing
	tokens := p.tokens[start:end]
	if len(tokens) == 0 {
		return nil
	}

	// Filter out trailing newlines
	for len(tokens) > 0 && tokens[len(tokens)-1].Type == TokenNewline {
		tokens = tokens[:len(tokens)-1]
	}
	if len(tokens) == 0 {
		return nil
	}

	// Create inline parser
	ip := &inlineParser{
		source:     p.source,
		tokens:     tokens,
		pos:        0,
		start:      tokens[0].Start,
		end:        tokens[len(tokens)-1].End,
		delimiters: make([]delimiter, 0, 8),
		linkDefs:   p.linkDefs,
		errors:     &p.errors,
	}

	return ip.parse()
}

// parse performs inline parsing and returns the resulting nodes.
//
//nolint:revive // function-length: inline parsing handles multiple token types
func (ip *inlineParser) parse() []Node {
	nodesPtr := nodeSlicePool.Get().(*[]Node)
	nodes := (*nodesPtr)[:0]
	defer func() {
		*nodesPtr = nodes[:0]
		nodeSlicePool.Put(nodesPtr)
	}()

	textStart := -1

	for ip.pos < len(ip.tokens) {
		tok := ip.tokens[ip.pos]

		switch tok.Type {
		case TokenBacktick:
			// Flush pending text
			if textStart >= 0 {
				nodes = append(
					nodes,
					ip.buildTextNode(
						textStart,
						ip.pos,
					),
				)
				textStart = -1
			}
			// Parse inline code
			if node := ip.parseInlineCode(); node != nil {
				nodes = append(nodes, node)
			} else {
				// Not inline code, treat as text
				if textStart < 0 {
					textStart = ip.pos
				}
				ip.pos++
			}

		case TokenAsterisk, TokenUnderscore:
			// Flush pending text
			if textStart >= 0 {
				nodes = append(
					nodes,
					ip.buildTextNode(
						textStart,
						ip.pos,
					),
				)
				textStart = -1
			}
			// Handle emphasis delimiter
			ip.handleEmphasisDelimiter()

		case TokenTilde:
			// Check for strikethrough (~~)
			if ip.pos+1 < len(ip.tokens) &&
				ip.tokens[ip.pos+1].Type == TokenTilde {
				// Flush pending text
				if textStart >= 0 {
					nodes = append(
						nodes,
						ip.buildTextNode(
							textStart,
							ip.pos,
						),
					)
					textStart = -1
				}
				if node := ip.parseStrikethrough(); node != nil {
					nodes = append(nodes, node)
				} else {
					if textStart < 0 {
						textStart = ip.pos
					}
					ip.pos++
				}
			} else {
				if textStart < 0 {
					textStart = ip.pos
				}
				ip.pos++
			}

		case TokenBracketOpen:
			// Flush pending text
			if textStart >= 0 {
				nodes = append(
					nodes,
					ip.buildTextNode(
						textStart,
						ip.pos,
					),
				)
				textStart = -1
			}
			// Check for wikilink [[...]] or link [...]
			if ip.pos+1 < len(ip.tokens) &&
				ip.tokens[ip.pos+1].Type == TokenBracketOpen {
				if node := ip.parseWikilink(); node != nil {
					nodes = append(nodes, node)
				} else {
					if textStart < 0 {
						textStart = ip.pos
					}
					ip.pos++
				}
			} else if node := ip.parseLink(); node != nil {
				nodes = append(nodes, node)
			} else {
				if textStart < 0 {
					textStart = ip.pos
				}
				ip.pos++
			}

		case TokenNewline:
			// Newlines in inline content become spaces
			if textStart >= 0 {
				nodes = append(
					nodes,
					ip.buildTextNode(
						textStart,
						ip.pos,
					),
				)
				textStart = -1
			}
			ip.pos++

		case TokenEOF,
			TokenWhitespace,
			TokenText,
			TokenError,
			TokenHash,
			TokenDash,
			TokenPlus,
			TokenDot,
			TokenColon,
			TokenPipe,
			TokenBracketClose,
			TokenParenOpen,
			TokenParenClose,
			TokenGreaterThan,
			TokenNumber,
			TokenX:
			// Accumulate text for all other token types
			if textStart < 0 {
				textStart = ip.pos
			}
			ip.pos++

		default:
			// Accumulate text for any unhandled token type
			if textStart < 0 {
				textStart = ip.pos
			}
			ip.pos++
		}
	}

	// Flush remaining text
	if textStart >= 0 {
		nodes = append(
			nodes,
			ip.buildTextNode(textStart, ip.pos),
		)
	}

	// Process delimiter stack for emphasis
	nodes = ip.processDelimiters(nodes)

	// Copy result
	result := make([]Node, len(nodes))
	copy(result, nodes)

	return result
}

// buildTextNode creates a text node from tokens in range [start, end).
func (ip *inlineParser) buildTextNode(
	start, end int,
) Node {
	if start >= end || start >= len(ip.tokens) {
		return nil
	}

	startOffset := ip.tokens[start].Start
	endOffset := ip.tokens[end-1].End
	if end > len(ip.tokens) {
		endOffset = ip.tokens[len(ip.tokens)-1].End
	}

	return NewNodeBuilder(NodeTypeText).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(ip.source[startOffset:endOffset]).
		Build()
}

// parseInlineCode parses inline code (`code` or “code with `backtick` “).
func (ip *inlineParser) parseInlineCode() Node {
	if ip.pos >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenBacktick {
		return nil
	}

	startOffset := ip.tokens[ip.pos].Start

	// Count opening backticks
	openCount := 0
	for ip.pos < len(ip.tokens) && ip.tokens[ip.pos].Type == TokenBacktick {
		openCount++
		ip.pos++
	}

	// Collect content until matching closing backticks
	contentStart := ip.pos
	for ip.pos < len(ip.tokens) {
		// Look for matching backtick sequence
		if ip.tokens[ip.pos].Type == TokenBacktick {
			closeCount := 0
			closeStart := ip.pos
			for ip.pos < len(ip.tokens) && ip.tokens[ip.pos].Type == TokenBacktick {
				closeCount++
				ip.pos++
			}
			if closeCount == openCount {
				// Found matching close
				endOffset := ip.tokens[ip.pos-1].End
				// contentEnd is closeStart, used below for unused variable suppression
				_ = closeStart

				return NewNodeBuilder(
					NodeTypeCode,
				).
					WithStart(startOffset).
					WithEnd(endOffset).
					WithSource(ip.source[startOffset:endOffset]).
					Build()
			}
			// Not matching, continue looking
		} else {
			ip.pos++
		}
	}

	// No closing backticks found - reset and return nil
	ip.pos = contentStart - openCount

	return nil
}

// handleEmphasisDelimiter handles * or _ delimiter for emphasis.
func (ip *inlineParser) handleEmphasisDelimiter() {
	if ip.pos >= len(ip.tokens) {
		return
	}

	tok := ip.tokens[ip.pos]
	delimType := tok.Type

	// Count consecutive delimiters
	count := 0
	delimStart := ip.pos
	for ip.pos < len(ip.tokens) && ip.tokens[ip.pos].Type == delimType {
		count++
		ip.pos++
	}

	// Determine if left-flanking and/or right-flanking
	canOpen, canClose := ip.isFlankingDelimiter(
		delimStart,
		delimType,
		count,
	)

	// For underscore, apply intraword restriction
	if delimType == TokenUnderscore {
		canOpen, canClose = ip.applyUnderscoreRestriction(
			delimStart,
			canOpen,
			canClose,
		)
	}

	// Push onto delimiter stack
	ip.delimiters = append(
		ip.delimiters,
		delimiter{
			token:     tok,
			count:     count,
			canOpen:   canOpen,
			canClose:  canClose,
			active:    true,
			textStart: ip.tokens[delimStart].Start,
			delimType: delimType,
		},
	)
}

// isFlankingDelimiter determines if a delimiter run is left-flanking and/or right-flanking.
// Per CommonMark spec section 6.2:
// - Left-flanking: not followed by whitespace, and (not followed by punctuation OR preceded by whitespace/punctuation)
// - Right-flanking: not preceded by whitespace, and (not preceded by punctuation OR followed by whitespace/punctuation)
func (ip *inlineParser) isFlankingDelimiter(
	pos int,
	delimType TokenType,
	count int,
) (canOpen, canClose bool) {
	// Get character before delimiter run
	charBefore := ' ' // Default to space (start of text)
	if pos > 0 {
		prevTok := ip.tokens[pos-1]
		if len(prevTok.Source) > 0 {
			charBefore = rune(
				prevTok.Source[len(prevTok.Source)-1],
			)
		}
	}

	// Get character after delimiter run
	charAfter := ' ' // Default to space (end of text)
	afterPos := pos + count
	if afterPos < len(ip.tokens) {
		nextTok := ip.tokens[afterPos]
		if len(nextTok.Source) > 0 {
			charAfter = rune(nextTok.Source[0])
		}
	}

	// Check flanking conditions
	beforeIsWhitespace := unicode.IsSpace(
		charBefore,
	)
	beforeIsPunctuation := unicode.IsPunct(
		charBefore,
	)
	afterIsWhitespace := unicode.IsSpace(
		charAfter,
	)
	afterIsPunctuation := unicode.IsPunct(
		charAfter,
	)

	// Left-flanking delimiter run
	leftFlanking := !afterIsWhitespace &&
		(!afterIsPunctuation || beforeIsWhitespace || beforeIsPunctuation)

	// Right-flanking delimiter run
	rightFlanking := !beforeIsWhitespace &&
		(!beforeIsPunctuation || afterIsWhitespace || afterIsPunctuation)

	canOpen = leftFlanking
	canClose = rightFlanking

	return canOpen, canClose
}

// applyUnderscoreRestriction applies the underscore intraword restriction.
// foo_bar_baz should NOT be parsed as emphasis.
func (ip *inlineParser) applyUnderscoreRestriction(
	pos int,
	canOpen, canClose bool,
) (newCanOpen, newCanClose bool) {
	// Get character before
	charBefore := ' '
	if pos > 0 {
		prevTok := ip.tokens[pos-1]
		if len(prevTok.Source) > 0 {
			charBefore = rune(
				prevTok.Source[len(prevTok.Source)-1],
			)
		}
	}

	// Get character after (after the delimiter run)
	charAfter := ' '
	// Find end of delimiter run
	endPos := pos
	for endPos < len(ip.tokens) && ip.tokens[endPos].Type == TokenUnderscore {
		endPos++
	}
	if endPos < len(ip.tokens) {
		nextTok := ip.tokens[endPos]
		if len(nextTok.Source) > 0 {
			charAfter = rune(nextTok.Source[0])
		}
	}

	// If both sides are alphanumeric, underscore cannot open or close
	beforeIsAlnum := unicode.IsLetter(
		charBefore,
	) ||
		unicode.IsDigit(charBefore)
	afterIsAlnum := unicode.IsLetter(charAfter) ||
		unicode.IsDigit(charAfter)

	if beforeIsAlnum && afterIsAlnum {
		return false, false
	}

	// If preceded by alphanumeric, cannot open
	if beforeIsAlnum {
		canOpen = false
	}

	// If followed by alphanumeric, cannot close
	if afterIsAlnum {
		canClose = false
	}

	return canOpen, canClose
}

// processDelimiters processes the delimiter stack and creates emphasis nodes.
// Implements the CommonMark emphasis processing algorithm (section 6.4):
// https://spec.commonmark.org/0.30/#emphasis-and-strong-emphasis
func (ip *inlineParser) processDelimiters(
	nodes []Node,
) []Node {
	// Process emphasis according to CommonMark spec section 6.4
	// We iterate through the delimiter stack, processing emphasis in multiple passes.
	// Each pass tries to match delimiters and create emphasis/strong nodes.

	result := nodes
	for {
		newResult, found := ip.processEmphasisPass(result)
		if !found {
			break
		}
		result = newResult
	}

	return result
}

// processEmphasisPass makes a single pass through the delimiter stack,
// processing emphasis and strong emphasis. Returns the updated node list and
// a boolean indicating if any emphasis was processed.
func (ip *inlineParser) processEmphasisPass(
	nodes []Node,
) ([]Node, bool) {
	if len(ip.delimiters) == 0 {
		return nodes, false
	}

	// Process delimiters from left to right
	for i := range len(ip.delimiters) {
		opener := &ip.delimiters[i]

		// Skip inactive delimiters and those that can't open
		if !opener.active || !opener.canOpen {
			continue
		}

		// Look for matching closer to the right
		for j := i + 1; j < len(ip.delimiters); j++ {
			closer := &ip.delimiters[j]

			// Skip inactive closers and those that can't close
			if !closer.active || !closer.canClose {
				continue
			}

			// Must be same type (asterisk or underscore)
			if opener.delimType != closer.delimType {
				continue
			}

			// Found a potential match
			// Try to match strong emphasis first (use 2 delimiters)
			if opener.count >= 2 && closer.count >= 2 {
				result, success := ip.createEmphasisNode(
					nodes,
					i,
					j,
					2,
					NodeTypeStrong,
				)
				if success {
					return result, true
				}
			}

			// Try regular emphasis (use 1 delimiter)
			if opener.count >= 1 && closer.count >= 1 {
				result, success := ip.createEmphasisNode(
					nodes,
					i,
					j,
					1,
					NodeTypeEmphasis,
				)
				if success {
					return result, true
				}
			}
		}
	}

	return nodes, false
}

// createEmphasisNode attempts to create an emphasis node by matching delimiters.
// Returns the updated node list and true if successful.
func (ip *inlineParser) createEmphasisNode(
	nodes []Node,
	openerIdx, closerIdx int,
	delimCount int,
	nodeType NodeType,
) ([]Node, bool) {
	if openerIdx >= len(ip.delimiters) || closerIdx >= len(ip.delimiters) {
		return nodes, false
	}

	opener := &ip.delimiters[openerIdx]
	closer := &ip.delimiters[closerIdx]

	// Reduce delimiter counts
	opener.count -= delimCount
	closer.count -= delimCount

	// Mark as inactive if no more delimiters
	if opener.count == 0 {
		opener.active = false
	}
	if closer.count == 0 {
		closer.active = false
	}

	// Find the index of the first node in the emphasis range
	var startIdx int
	var endIdx int

	// Find start index: first node at or after opener position
	startIdx = len(nodes)
	for k := range nodes {
		start, _ := nodes[k].Span()
		if start >= opener.token.Start {
			startIdx = k

			break
		}
	}

	// Find end index: last node at or before closer position
	endIdx = -1
	for k := len(nodes) - 1; k >= 0; k-- {
		_, end := nodes[k].Span()
		if end <= closer.token.End {
			endIdx = k

			break
		}
	}

	if startIdx > endIdx || startIdx >= len(nodes) {
		return nodes, false
	}

	// Extract children (nodes between opener and closer, inclusive)
	var children []Node
	if startIdx <= endIdx {
		children = nodes[startIdx : endIdx+1]
	}

	// Create emphasis node
	node := NewNodeBuilder(nodeType).
		WithStart(opener.token.Start).
		WithEnd(closer.token.End).
		WithSource(ip.source[opener.token.Start:closer.token.End]).
		WithChildren(children).
		Build()

	if node == nil {
		return nodes, false
	}

	// Build new node list with emphasis node replacing the matched range
	newNodes := make([]Node, 0, len(nodes))
	newNodes = append(newNodes, nodes[:startIdx]...)
	newNodes = append(newNodes, node)
	if endIdx+1 < len(nodes) {
		newNodes = append(newNodes, nodes[endIdx+1:]...)
	}

	// Mark all delimiters between opener and closer as inactive
	for k := openerIdx + 1; k < closerIdx; k++ {
		ip.delimiters[k].active = false
	}

	return newNodes, true
}

// parseStrikethrough parses ~~strikethrough~~ content.
func (ip *inlineParser) parseStrikethrough() Node {
	if ip.pos+1 >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenTilde ||
		ip.tokens[ip.pos+1].Type != TokenTilde {
		return nil
	}

	startOffset := ip.tokens[ip.pos].Start
	ip.pos += 2 // Skip opening ~~

	// Find closing ~~
	contentStart := ip.pos
	for ip.pos < len(ip.tokens) {
		if ip.pos+1 < len(ip.tokens) &&
			ip.tokens[ip.pos].Type == TokenTilde &&
			ip.tokens[ip.pos+1].Type == TokenTilde {
			// Found closing ~~
			contentEnd := ip.pos
			ip.pos += 2 // Skip closing ~~

			endOffset := ip.tokens[ip.pos-1].End

			// Parse content as inline
			var children []Node
			if contentStart < contentEnd {
				subParser := &inlineParser{
					source: ip.source,
					tokens: ip.tokens[contentStart:contentEnd],
					pos:    0,
					start:  ip.tokens[contentStart].Start,
					end:    ip.tokens[contentEnd-1].End,
					delimiters: make(
						[]delimiter,
						0,
					),
					linkDefs: ip.linkDefs,
					errors:   ip.errors,
				}
				children = subParser.parse()
			}

			return NewNodeBuilder(
				NodeTypeStrikethrough,
			).
				WithStart(startOffset).
				WithEnd(endOffset).
				WithSource(ip.source[startOffset:endOffset]).
				WithChildren(children).
				Build()
		}
		ip.pos++
	}

	// No closing found, reset
	ip.pos = contentStart - 2

	return nil
}

// parseWikilink parses [[target|display#anchor]] wikilinks.
func (ip *inlineParser) parseWikilink() Node {
	if ip.pos+1 >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenBracketOpen ||
		ip.tokens[ip.pos+1].Type != TokenBracketOpen {
		return nil
	}

	startOffset := ip.tokens[ip.pos].Start
	ip.pos += 2 // Skip [[

	// Collect content until ]]
	var parts [][]byte
	contentStart := ip.pos

	for ip.pos < len(ip.tokens) {
		if ip.pos+1 < len(ip.tokens) &&
			ip.tokens[ip.pos].Type == TokenBracketClose &&
			ip.tokens[ip.pos+1].Type == TokenBracketClose {
			// Found ]]
			contentEnd := ip.pos
			ip.pos += 2 // Skip ]]

			endOffset := ip.tokens[ip.pos-1].End

			// Build content
			for i := contentStart; i < contentEnd; i++ {
				parts = append(
					parts,
					ip.tokens[i].Source,
				)
			}
			content := bytes.Join(parts, nil)

			// Parse content: target|display#anchor
			target, display, anchor := parseWikilinkContent(
				content,
			)

			return NewNodeBuilder(
				NodeTypeWikilink,
			).
				WithStart(startOffset).
				WithEnd(endOffset).
				WithSource(ip.source[startOffset:endOffset]).
				WithTarget(target).
				WithDisplay(display).
				WithAnchor(anchor).
				Build()
		}

		// Don't allow newlines in wikilinks
		if ip.tokens[ip.pos].Type == TokenNewline {
			ip.pos = contentStart - 2

			return nil
		}

		ip.pos++
	}

	// No closing found
	ip.pos = contentStart - 2

	return nil
}

// parseWikilinkContent parses the content of a wikilink.
// Format: target or target|display or target#anchor or target|display#anchor
func parseWikilinkContent(
	content []byte,
) (target, display, anchor []byte) {
	// Check for | separator
	pipeIdx := bytes.IndexByte(content, '|')
	hashIdx := bytes.IndexByte(content, '#')

	switch { //nolint:gocritic // complex conditions don't fit tagged switch
	case pipeIdx >= 0 && (hashIdx < 0 || pipeIdx < hashIdx):
		// Has display text
		target = bytes.TrimSpace(
			content[:pipeIdx],
		)
		remaining := content[pipeIdx+1:]

		// Check for # in remaining
		hashIdx = bytes.IndexByte(remaining, '#')
		if hashIdx >= 0 {
			display = bytes.TrimSpace(
				remaining[:hashIdx],
			)
			anchor = bytes.TrimSpace(
				remaining[hashIdx+1:],
			)
		} else {
			display = bytes.TrimSpace(remaining)
		}
	case hashIdx >= 0:
		// Has anchor but no display
		target = bytes.TrimSpace(
			content[:hashIdx],
		)
		anchor = bytes.TrimSpace(
			content[hashIdx+1:],
		)
	default:
		// Just target
		target = bytes.TrimSpace(content)
	}

	return target, display, anchor
}

// parseLink parses [text](url "title") or [text][ref] links.
func (ip *inlineParser) parseLink() Node {
	if ip.pos >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenBracketOpen {
		return nil
	}

	startOffset := ip.tokens[ip.pos].Start
	ip.pos++ // Skip [

	// Collect link text until ]
	textStart := ip.pos
	bracketDepth := 1
	for ip.pos < len(ip.tokens) && bracketDepth > 0 {
		switch ip.tokens[ip.pos].Type { //nolint:exhaustive // Only care about bracket tokens
		case TokenBracketOpen:
			bracketDepth++
		case TokenBracketClose:
			bracketDepth--
		default:
			// Other tokens don't affect bracket depth
		}
		if bracketDepth > 0 {
			ip.pos++
		}
	}

	if bracketDepth != 0 {
		ip.pos = textStart - 1

		return nil
	}

	textEnd := ip.pos
	ip.pos++ // Skip ]

	// Check what follows: ( for inline link, [ for reference link, or shortcut
	if ip.pos >= len(ip.tokens) {
		// Shortcut reference: [text] uses text as reference
		return ip.parseShortcutLink(
			startOffset,
			textStart,
			textEnd,
		)
	}

	switch ip.tokens[ip.pos].Type { //nolint:exhaustive // Only care about link-related tokens
	case TokenParenOpen:
		// Inline link: [text](url "title")
		return ip.parseInlineLink(
			startOffset,
			textStart,
			textEnd,
		)
	case TokenBracketOpen:
		// Reference link: [text][ref]
		return ip.parseReferenceLink(
			startOffset,
			textStart,
			textEnd,
		)
	default:
		// Shortcut reference: [text] uses text as reference
		return ip.parseShortcutLink(
			startOffset,
			textStart,
			textEnd,
		)
	}
}

// parseInlineLink parses [text](url "title").
//
//nolint:revive // function-length: inline link parsing handles URL/title
func (ip *inlineParser) parseInlineLink(
	startOffset int,
	textStart, textEnd int,
) Node {
	if ip.pos >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenParenOpen {
		return nil
	}
	ip.pos++ // Skip (

	// Collect URL
	var urlParts [][]byte
	for ip.pos < len(ip.tokens) {
		tok := ip.tokens[ip.pos]
		if tok.Type == TokenParenClose ||
			tok.Type == TokenWhitespace ||
			tok.Type == TokenNewline {
			break
		}
		urlParts = append(urlParts, tok.Source)
		ip.pos++
	}

	url := bytes.Join(urlParts, nil)

	// Optional whitespace and title
	var title []byte
	for ip.pos < len(ip.tokens) && ip.tokens[ip.pos].Type == TokenWhitespace {
		ip.pos++
	}

	// Check for title in quotes
	if ip.pos < len(ip.tokens) &&
		ip.tokens[ip.pos].Type == TokenText {
		text := ip.tokens[ip.pos].Source
		if len(text) >= 2 &&
			(text[0] == '"' || text[0] == '\'') &&
			text[len(text)-1] == text[0] {
			title = text[1 : len(text)-1]
			ip.pos++
		}
	}

	// Skip whitespace
	for ip.pos < len(ip.tokens) && ip.tokens[ip.pos].Type == TokenWhitespace {
		ip.pos++
	}

	// Must end with )
	if ip.pos >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenParenClose {
		ip.pos = textStart - 1

		return nil
	}
	ip.pos++ // Skip )

	endOffset := ip.tokens[ip.pos-1].End

	// Parse text content
	var children []Node
	if textStart < textEnd {
		subParser := &inlineParser{
			source:     ip.source,
			tokens:     ip.tokens[textStart:textEnd],
			pos:        0,
			start:      ip.tokens[textStart].Start,
			end:        ip.tokens[textEnd-1].End,
			delimiters: make([]delimiter, 0),
			linkDefs:   ip.linkDefs,
			errors:     ip.errors,
		}
		children = subParser.parse()
	}

	return NewNodeBuilder(NodeTypeLink).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(ip.source[startOffset:endOffset]).
		WithURL(url).
		WithLinkTitle(title).
		WithChildren(children).
		Build()
}

// parseReferenceLink parses [text][ref].
//
//nolint:revive // function-length: reference link parsing handles label lookup
func (ip *inlineParser) parseReferenceLink(
	startOffset int,
	textStart, textEnd int,
) Node {
	if ip.pos >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenBracketOpen {
		return nil
	}
	ip.pos++ // Skip [

	// Collect reference label
	var labelParts [][]byte
	for ip.pos < len(ip.tokens) && ip.tokens[ip.pos].Type != TokenBracketClose {
		if ip.tokens[ip.pos].Type == TokenNewline {
			// Invalid
			ip.pos = textStart - 1

			return nil
		}
		labelParts = append(
			labelParts,
			ip.tokens[ip.pos].Source,
		)
		ip.pos++
	}

	if ip.pos >= len(ip.tokens) ||
		ip.tokens[ip.pos].Type != TokenBracketClose {
		ip.pos = textStart - 1

		return nil
	}
	ip.pos++ // Skip ]

	endOffset := ip.tokens[ip.pos-1].End

	// Look up reference
	var label string
	if len(labelParts) == 0 {
		// Empty reference [text][] uses text as label
		var textParts [][]byte
		for i := textStart; i < textEnd; i++ {
			textParts = append(
				textParts,
				ip.tokens[i].Source,
			)
		}
		label = strings.ToLower(
			string(bytes.Join(textParts, nil)),
		)
	} else {
		label = strings.ToLower(string(bytes.Join(labelParts, nil)))
	}
	label = strings.TrimSpace(label)

	def, found := ip.linkDefs[label]
	if !found {
		// Reference not found
		ip.pos = textStart - 1

		return nil
	}

	// Parse text content
	var children []Node
	if textStart < textEnd {
		subParser := &inlineParser{
			source:     ip.source,
			tokens:     ip.tokens[textStart:textEnd],
			pos:        0,
			start:      ip.tokens[textStart].Start,
			end:        ip.tokens[textEnd-1].End,
			delimiters: make([]delimiter, 0),
			linkDefs:   ip.linkDefs,
			errors:     ip.errors,
		}
		children = subParser.parse()
	}

	return NewNodeBuilder(NodeTypeLink).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(ip.source[startOffset:endOffset]).
		WithURL(def.url).
		WithLinkTitle(def.title).
		WithChildren(children).
		Build()
}

// parseShortcutLink parses [text] as shortcut reference link.
func (ip *inlineParser) parseShortcutLink(
	startOffset int,
	textStart, textEnd int,
) Node {
	endOffset := ip.tokens[ip.pos-1].End

	// Use text as label
	var textParts [][]byte
	for i := textStart; i < textEnd; i++ {
		textParts = append(
			textParts,
			ip.tokens[i].Source,
		)
	}
	label := strings.ToLower(
		string(bytes.Join(textParts, nil)),
	)
	label = strings.TrimSpace(label)

	def, found := ip.linkDefs[label]
	if !found {
		// Not a valid shortcut link, reset
		ip.pos = textStart - 1

		return nil
	}

	// Parse text content
	var children []Node
	if textStart < textEnd {
		subParser := &inlineParser{
			source:     ip.source,
			tokens:     ip.tokens[textStart:textEnd],
			pos:        0,
			start:      ip.tokens[textStart].Start,
			end:        ip.tokens[textEnd-1].End,
			delimiters: make([]delimiter, 0),
			linkDefs:   ip.linkDefs,
			errors:     ip.errors,
		}
		children = subParser.parse()
	}

	return NewNodeBuilder(NodeTypeLink).
		WithStart(startOffset).
		WithEnd(endOffset).
		WithSource(ip.source[startOffset:endOffset]).
		WithURL(def.url).
		WithLinkTitle(def.title).
		WithChildren(children).
		Build()
}
