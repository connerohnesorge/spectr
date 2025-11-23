package mdparser

import "strings"

// lexStart is the initial state function.
func lexStart(l *Lexer) stateFn {
	if l.pos >= len(l.input) {
		return nil // EOF
	}

	// If we're in a code block, only look for code fence or code content
	if l.inCodeBlock {
		return lexCodeBlockContent
	}

	// Check if at start of line
	if !l.atStartOfLine() {
		// Not at start of line, consume as text
		return lexText
	}

	line := l.peekLine()

	// Check for code fence (``` or more backticks)
	if strings.HasPrefix(line, "```") {
		return lexCodeFence
	}

	// Check for header (# at start of line)
	if len(line) > 0 && line[0] == '#' {
		return lexHeader
	}

	// Check for list item (- or * followed by space)
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
		return lexListItem
	}

	// Check for ordered list item
	if l.checkOrderedListItem(line) {
		return lexListItem
	}

	// Check for blank line
	if strings.TrimSpace(line) == "" {
		return lexBlankLine
	}

	// Default to text
	return lexText
}

// lexHeader handles markdown headers.
func lexHeader(l *Lexer) stateFn {
	// Count leading # symbols
	count := 0
	for count < maxHeaderLevel && l.peek() == '#' {
		l.next()
		count++
	}

	// Skip whitespace after #
	if l.peek() == ' ' {
		l.next()
	}

	// Consume rest of line
	l.skipToEndOfLine()

	l.emit(TokenHeader)

	// Consume newline if present
	if l.peek() == '\n' {
		l.next()
		l.ignore()
	}

	return lexStart
}

// lexText handles plain text lines.
func lexText(l *Lexer) stateFn {
	// Consume until end of line or special character at start of line
	for {
		r := l.peek()
		if r == 0 || r == '\n' {
			break
		}
		l.next()
	}

	if l.pos > l.start {
		l.emit(TokenText)
	}

	// Consume newline if present
	if l.peek() == '\n' {
		l.next()
		l.ignore()
	}

	return lexStart
}

// lexCodeFence handles code fence opening/closing.
func lexCodeFence(l *Lexer) stateFn {
	// Count backticks
	backtickCount := 0
	for l.peek() == '`' {
		l.next()
		backtickCount++
	}

	// Consume rest of line (language identifier or closing fence)
	l.skipToEndOfLine()

	l.emit(TokenCodeFence)

	// Consume newline if present
	if l.peek() == '\n' {
		l.next()
		l.ignore()
	}

	// Toggle code block state
	l.inCodeBlock = !l.inCodeBlock

	return lexStart
}

// lexCodeBlockContent handles content inside code blocks.
func lexCodeBlockContent(l *Lexer) stateFn {
	// Check if this line is a closing fence
	if l.atStartOfLine() {
		line := l.peekLine()
		if strings.HasPrefix(line, "```") {
			return lexCodeFence
		}
	}

	// Consume entire line as code content
	l.skipToEndOfLine()

	if l.pos > l.start {
		l.emit(TokenCodeContent)
	}

	// Consume newline if present
	if l.peek() == '\n' {
		l.next()
		l.ignore()
	}

	return lexStart
}

// lexListItem handles list items.
func lexListItem(l *Lexer) stateFn {
	// Skip bullet/number
	if l.peek() == '-' || l.peek() == '*' {
		l.next()
	} else {
		// Skip digits and period for ordered lists
		for l.peek() >= '0' && l.peek() <= '9' {
			l.next()
		}
		if l.peek() == '.' {
			l.next()
		}
	}

	// Skip space after bullet/number
	if l.peek() == ' ' {
		l.next()
	}

	// Consume rest of line
	l.skipToEndOfLine()

	l.emit(TokenListItem)

	// Consume newline if present
	if l.peek() == '\n' {
		l.next()
		l.ignore()
	}

	return lexStart
}

// lexBlankLine handles blank lines.
func lexBlankLine(l *Lexer) stateFn {
	// Consume whitespace on this line
	for l.peek() == ' ' || l.peek() == '\t' {
		l.next()
	}

	// Consume newline
	if l.peek() == '\n' {
		l.next()
	}

	l.emit(TokenBlankLine)

	return lexStart
}
