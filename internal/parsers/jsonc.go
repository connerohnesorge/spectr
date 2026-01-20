package parsers

// ToJSON strips out comments and trailing commas and convert the input to a
// valid JSON per the official spec: https://tools.ietf.org/html/rfc8259
//
// The resulting JSON will always be the same length as the input and it will
// include all of the same line breaks at matching offsets. This is to ensure
// the result can be later processed by a external parser and that that
// parser will report messages or errors with the correct offsets.
func JSONCToJSON(src []byte) []byte {
	return toJSON(src, nil)
}

// ToJSONInPlace is the same as ToJSON, but this method reuses the input json
// buffer to avoid allocations. Do not use the original bytes slice upon return.
func JSONCToJSONInPlace(src []byte) []byte {
	return toJSON(src, src)
}

// toJSON strips JSONC from the provided bytes and writes the result to dst.
// It preserves line breaks at the same offsets to maintain error reporting accuracy.
func toJSON(src, dst []byte) []byte {
	out := dst[:0]

	for i := 0; i < len(src); i++ {
		if src[i] == '/' && i < len(src)-1 {
			if src[i+1] == '/' {
				out, i = skipSingleLineComment(src, out, i)

				continue
			}

			if src[i+1] == '*' {
				out, i = skipMultiLineComment(src, out, i)

				continue
			}
		}

		out = append(out, src[i])

		switch src[i] {
		case '"':
			out, i = copyStringLiteral(src, out, i)
		case '}', ']':
			out = removeTrailingComma(out)
		}
	}

	return out
}

// skipSingleLineComment processes a single-line comment (//) starting at pos
// and returns the updated destination buffer and final index.
func skipSingleLineComment(src, dst []byte, pos int) ([]byte, int) {
	out := dst
	out = append(out, ' ', ' ')
	idx := pos + 2

	for ; idx < len(src); idx++ {
		if src[idx] == '\n' {
			out = append(out, '\n')

			break
		}

		if src[idx] == '\t' || src[idx] == '\r' {
			out = append(out, src[idx])
		} else {
			out = append(out, ' ')
		}
	}

	return out, idx
}

// skipMultiLineComment processes a multi-line comment (/* */) starting at pos
// and returns the updated destination buffer and final index.
func skipMultiLineComment(src, dst []byte, pos int) ([]byte, int) {
	out := dst
	out = append(out, ' ', ' ')
	idx := pos + 2

	for ; idx < len(src)-1; idx++ {
		if src[idx] == '*' && src[idx+1] == '/' {
			out = append(out, ' ', ' ')
			idx++

			break
		}

		if src[idx] == '\n' || src[idx] == '\t' || src[idx] == '\r' {
			out = append(out, src[idx])
		} else {
			out = append(out, ' ')
		}
	}

	return out, idx
}

// copyStringLiteral copies a string literal (including escape sequences)
// starting at pos and returns the updated destination buffer and final index.
func copyStringLiteral(src, dst []byte, pos int) ([]byte, int) {
	out := dst
	idx := pos + 1

	for ; idx < len(src); idx++ {
		out = append(out, src[idx])

		if src[idx] == '"' && !isEscapedQuote(src, idx) {
			break
		}
	}

	return out, idx
}

// isEscapedQuote checks if the quote at position idx is escaped by counting
// preceding backslashes. Returns true if the quote is escaped.
func isEscapedQuote(src []byte, idx int) bool {
	backslashCount := 0
	for j := idx - 1; j >= 0 && src[j] == '\\'; j-- {
		backslashCount++
	}

	return backslashCount%2 != 0
}

// removeTrailingComma removes a trailing comma before a closing brace/bracket
// by replacing it with a space.
func removeTrailingComma(dst []byte) []byte {
	for j := len(dst) - 2; j >= 0; j-- {
		if dst[j] <= ' ' {
			continue
		}

		if dst[j] == ',' {
			dst[j] = ' '
		}

		break
	}

	return dst
}
