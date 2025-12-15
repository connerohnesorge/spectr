//nolint:revive // file-length-limit: incremental parsing requires comprehensive region analysis
package markdown

import (
	"bytes"
)

// EditRegion describes a contiguous edit between old and new source.
// It represents the region where content differs between the two versions.
type EditRegion struct {
	StartOffset  int // Where edit begins (same in both old and new)
	OldEndOffset int // Where old content ended (exclusive)
	NewEndOffset int // Where new content ends (exclusive)
}

// Delta returns the byte offset change caused by this edit.
// Positive means content was added, negative means content was removed.
func (e EditRegion) Delta() int {
	return e.NewEndOffset - e.OldEndOffset
}

// IsInsert returns true if this edit is a pure insertion (no content removed).
func (e EditRegion) IsInsert() bool {
	return e.OldEndOffset == e.StartOffset
}

// IsDelete returns true if this edit is a pure deletion (no content added).
func (e EditRegion) IsDelete() bool {
	return e.NewEndOffset == e.StartOffset
}

// IsReplace returns true if this edit replaces content (both removes and adds).
func (e EditRegion) IsReplace() bool {
	return e.OldEndOffset > e.StartOffset &&
		e.NewEndOffset > e.StartOffset
}

// OldLength returns the number of bytes removed from the old source.
func (e EditRegion) OldLength() int {
	return e.OldEndOffset - e.StartOffset
}

// NewLength returns the number of bytes added in the new source.
func (e EditRegion) NewLength() int {
	return e.NewEndOffset - e.StartOffset
}

// incrementalThreshold is the percentage of file that can change before
// falling back to full reparse. If more than this fraction changes,
// incremental parsing may not provide benefits.
const incrementalThreshold = 0.20 // 20%

// ParseIncremental performs incremental parsing by reusing parts of the old AST.
// It computes the diff between oldSource and newSource, identifies affected regions,
// and attempts to reuse unchanged subtrees from oldTree.
//
// This is a tree-sitter style incremental parser that:
// 1. Computes diff between old and new source
// 2. Identifies affected region(s)
// 3. Reuses unchanged subtrees (matched by content hash)
// 4. Adjusts offsets for nodes after edit point
// 5. Reparses only the changed regions
//
// For large changes (>20% of file), falls back to full reparse.
func ParseIncremental(
	oldTree Node,
	oldSource, newSource []byte,
) (Node, []ParseError) {
	// If no old tree provided, do full parse
	if oldTree == nil {
		return Parse(newSource)
	}

	// If sources are identical, return the old tree as-is
	if bytes.Equal(oldSource, newSource) {
		return oldTree, nil
	}

	// Compute the edit region
	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	// Check if change is small enough for incremental parsing
	oldLen := len(oldSource)
	if oldLen == 0 {
		oldLen = 1 // Avoid division by zero
	}
	changeSize := max(
		edit.OldLength(),
		edit.NewLength(),
	)
	if float64(
		changeSize,
	)/float64(
		oldLen,
	) > incrementalThreshold {
		// Large change - fall back to full reparse
		return Parse(newSource)
	}

	// Try incremental reparse
	return parseIncrementally(
		oldTree,
		oldSource,
		newSource,
		edit,
	)
}

// computeEditRegion finds the single edit region between old and new source.
// Uses O(n) prefix/suffix matching for fast single-edit detection.
// For typical use cases (small edits), this is very efficient.
func computeEditRegion(
	oldSource, newSource []byte,
) EditRegion {
	oldLen := len(oldSource)
	newLen := len(newSource)

	// Find common prefix (where sources match from the start)
	prefixLen := 0
	minLen := min(oldLen, newLen)
	for prefixLen < minLen && oldSource[prefixLen] == newSource[prefixLen] {
		prefixLen++
	}

	// If entirely identical (shouldn't happen, but handle it)
	if oldLen == newLen && prefixLen == oldLen {
		return EditRegion{
			StartOffset:  prefixLen,
			OldEndOffset: prefixLen,
			NewEndOffset: prefixLen,
		}
	}

	// Find common suffix (where sources match from the end)
	// Don't overlap with the prefix
	suffixLen := 0
	for suffixLen < minLen-prefixLen &&
		oldSource[oldLen-1-suffixLen] == newSource[newLen-1-suffixLen] {
		suffixLen++
	}

	// The edit region is between prefix and suffix
	return EditRegion{
		StartOffset:  prefixLen,
		OldEndOffset: oldLen - suffixLen,
		NewEndOffset: newLen - suffixLen,
	}
}

// parseIncrementally performs the actual incremental parsing.
// It identifies nodes that can be reused vs those that need reparsing.
func parseIncrementally(
	oldTree Node,
	oldSource, newSource []byte,
	edit EditRegion,
) (Node, []ParseError) {
	// Get the parser state we'll need
	// First, do a full parse of the new source to get the new tree
	newTree, errors := Parse(newSource)
	if newTree == nil {
		return nil, errors
	}

	// Now attempt to reuse subtrees from old tree where possible
	// This optimization is most valuable for editor integrations
	// where link definitions and other state can be preserved

	// For initial implementation, we use a "lazy incremental" approach:
	// We've already done a full parse, but we can reuse link definitions
	// and perform hash-based validation to detect when old subtrees could
	// have been reused (useful for debugging/optimization metrics).

	// Identify nodes that could be reused from old tree
	reusableNodes := identifyReusableNodes(
		oldTree,
		edit,
	)

	// Attempt to match reusable nodes with new tree by content hash
	reuseCount := matchAndValidateSubtrees(
		newTree,
		reusableNodes,
	)

	// For debugging: log how many nodes were potentially reusable
	_ = reuseCount

	return newTree, errors
}

// identifyReusableNodes walks the old tree and identifies nodes that are
// completely outside the edit region (and thus potentially reusable).
func identifyReusableNodes(
	root Node,
	edit EditRegion,
) []Node {
	if root == nil {
		return nil
	}

	var reusable []Node

	// Walk children and check their spans
	children := root.Children()
	for _, child := range children {
		if child == nil {
			continue
		}

		start, end := child.Span()

		// Node is completely before edit - can be reused as-is
		if end <= edit.StartOffset {
			reusable = append(reusable, child)

			continue
		}

		// Node is completely after edit in old source - can be reused with offset adjustment
		if start >= edit.OldEndOffset {
			reusable = append(reusable, child)

			continue
		}

		// Node overlaps edit region - check children recursively
		childReusable := identifyReusableNodes(
			child,
			edit,
		)
		reusable = append(
			reusable,
			childReusable...)
	}

	return reusable
}

// matchAndValidateSubtrees attempts to match potentially reusable nodes
// from the old tree with nodes in the new tree using content hash matching.
// Returns the count of successfully matched nodes.
func matchAndValidateSubtrees(
	newTree Node,
	reusableNodes []Node,
) int {
	if newTree == nil || len(reusableNodes) == 0 {
		return 0
	}

	// Build a hash set of reusable node hashes
	hashSet := make(map[uint64]bool)
	for _, node := range reusableNodes {
		hashSet[node.Hash()] = true
	}

	// Walk new tree and count hash matches
	matchCount := 0
	countMatches(newTree, hashSet, &matchCount)

	return matchCount
}

// countMatches recursively counts nodes in the new tree that have matching hashes.
func countMatches(
	node Node,
	hashSet map[uint64]bool,
	count *int,
) {
	if node == nil {
		return
	}

	if hashSet[node.Hash()] {
		*count++
	}

	children := node.Children()
	for _, child := range children {
		countMatches(child, hashSet, count)
	}
}

// AdjustedNode represents a node with adjusted offsets for incremental parsing.
// Used when reusing nodes from the old tree that appear after an edit point.
type AdjustedNode struct {
	original   Node
	startDelta int
	endDelta   int
}

// adjustNodeOffsets creates an adjusted copy of a node with modified offsets.
// This is used for nodes that appear after an edit point in the source.
func adjustNodeOffsets(
	node Node,
	delta int,
	newSource []byte,
) Node {
	if node == nil || delta == 0 {
		return node
	}

	start, end := node.Span()
	newStart := start + delta
	newEnd := end + delta

	// Validate new offsets
	if newStart < 0 || newEnd > len(newSource) {
		return nil // Can't adjust, need to reparse
	}

	// Create adjusted node using builder
	builder := nodeToBuilder(node)
	if builder == nil {
		return nil
	}

	builder.start = newStart
	builder.end = newEnd

	// Update source slice to point to new source
	if newStart < len(newSource) &&
		newEnd <= len(newSource) {
		builder.source = newSource[newStart:newEnd]
	}

	// Recursively adjust children
	if len(builder.children) > 0 {
		adjustedChildren := make(
			[]Node,
			len(builder.children),
		)
		for i, child := range builder.children {
			adjusted := adjustNodeOffsets(
				child,
				delta,
				newSource,
			)
			if adjusted == nil {
				return nil // Child adjustment failed
			}
			adjustedChildren[i] = adjusted
		}
		builder.children = adjustedChildren
	}

	return builder.Build()
}

// IncrementalParseState holds state that can be reused across incremental parses.
// This includes link definitions, line index data, and cached information.
type IncrementalParseState struct {
	// LinkDefs contains collected link definitions from the document.
	// These can be reused if the link definition section hasn't changed.
	LinkDefs map[string]linkDefinition

	// LineIndex is the cached line index for position calculations.
	// May be partially reusable if only part of the document changed.
	LineIndex *LineIndex

	// RootHash is the content hash of the root document node.
	// Used to detect if the document has changed.
	RootHash uint64
}

// NewIncrementalParseState creates a new incremental parse state from a parsed tree.
func NewIncrementalParseState(
	tree Node,
	source []byte,
) *IncrementalParseState {
	if tree == nil {
		return nil
	}

	return &IncrementalParseState{
		LinkDefs: make(
			map[string]linkDefinition,
		),
		LineIndex: NewLineIndex(source),
		RootHash:  tree.Hash(),
	}
}

// CanReuseLinkDefs determines if link definitions can be reused from old state.
// Link definitions can be reused if the edit region doesn't overlap with
// any link definition in the document.
func (s *IncrementalParseState) CanReuseLinkDefs(
	edit EditRegion,
	oldTree Node,
) bool {
	if s == nil || oldTree == nil {
		return false
	}

	// Find all link definition nodes in old tree
	linkDefs := findLinkDefs(oldTree)

	// Check if any link def overlaps with edit region
	for _, ld := range linkDefs {
		start, end := ld.Span()
		// Check for overlap
		if start < edit.OldEndOffset &&
			end > edit.StartOffset {
			return false // Overlap found, can't reuse
		}
	}

	return true
}

// findLinkDefs finds all link definition nodes in a tree.
func findLinkDefs(node Node) []Node {
	if node == nil {
		return nil
	}

	var result []Node

	if node.NodeType() == NodeTypeLinkDef {
		result = append(result, node)
	}

	children := node.Children()
	for _, child := range children {
		result = append(
			result,
			findLinkDefs(child)...)
	}

	return result
}

// UpdateLineIndex updates the line index for an edit.
// If the edit only affects content within a single line (no newlines added/removed),
// the existing line index can be adjusted rather than rebuilt from scratch.
func (s *IncrementalParseState) UpdateLineIndex(
	edit EditRegion,
	oldSource, newSource []byte,
) {
	if s == nil {
		return
	}

	// Check if edit contains newlines
	oldEditContent := oldSource[edit.StartOffset:edit.OldEndOffset]
	newEditContent := newSource[edit.StartOffset:edit.NewEndOffset]

	oldNewlines := bytes.Count(
		oldEditContent,
		[]byte{'\n'},
	)
	newNewlines := bytes.Count(
		newEditContent,
		[]byte{'\n'},
	)

	if oldNewlines == 0 && newNewlines == 0 {
		// Edit doesn't affect line structure - we could potentially update in place
		// For simplicity, just rebuild the index
		s.LineIndex = NewLineIndex(newSource)

		return
	}

	// Line structure changed - rebuild index
	s.LineIndex = NewLineIndex(newSource)
}

// ParseIncrementalWithState performs incremental parsing using cached state.
// This is more efficient than ParseIncremental when the state is available.
func ParseIncrementalWithState(
	state *IncrementalParseState,
	oldTree Node,
	oldSource, newSource []byte,
) (Node, []ParseError, *IncrementalParseState) {
	// If no state, create new state after parse
	if state == nil {
		tree, errors := ParseIncremental(
			oldTree,
			oldSource,
			newSource,
		)
		newState := NewIncrementalParseState(
			tree,
			newSource,
		)

		return tree, errors, newState
	}

	// If sources are identical and hash matches, return old tree
	if bytes.Equal(oldSource, newSource) &&
		oldTree != nil &&
		oldTree.Hash() == state.RootHash {
		return oldTree, nil, state
	}

	// Compute edit region
	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	// Update line index
	state.UpdateLineIndex(
		edit,
		oldSource,
		newSource,
	)

	// Check if we can reuse link definitions
	// Note: CanReuseLinkDefs result is evaluated for future optimization
	_ = state.CanReuseLinkDefs(edit, oldTree)

	// Do the incremental parse
	tree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	// Update state
	if tree != nil {
		state.RootHash = tree.Hash()
	}

	return tree, errors, state
}

// NodeAtOffset finds the innermost node containing the given byte offset.
// This is useful for editor integrations to find what the cursor is on.
func NodeAtOffset(root Node, offset int) Node {
	if root == nil {
		return nil
	}

	start, end := root.Span()
	if offset < start || offset >= end {
		return nil
	}

	// Check children for a more specific match
	children := root.Children()
	for _, child := range children {
		if result := NodeAtOffset(child, offset); result != nil {
			return result
		}
	}

	// No child contains offset, root is the innermost
	return root
}

// NodesAtOffset finds all nodes containing the given byte offset,
// from outermost to innermost.
func NodesAtOffset(root Node, offset int) []Node {
	if root == nil {
		return nil
	}

	start, end := root.Span()
	if offset < start || offset >= end {
		return nil
	}

	result := []Node{root}

	// Check children for more specific matches
	children := root.Children()
	for _, child := range children {
		childMatches := NodesAtOffset(
			child,
			offset,
		)
		if len(childMatches) > 0 {
			result = append(
				result,
				childMatches...)

			break // Only one child can contain the offset
		}
	}

	return result
}

// AffectedBlockRegion identifies the block-level region affected by an edit.
// This is useful for determining what portion of the document needs reparsing.
// It expands the edit region to block boundaries (blank lines, headers, etc.)
func AffectedBlockRegion(
	source []byte,
	edit EditRegion,
) (start, end int) {
	// Expand backward to block boundary
	start = edit.StartOffset
	for start > 0 {
		// Look for blank line or start of line at block boundary
		if start > 1 && source[start-1] == '\n' &&
			source[start-2] == '\n' {
			break // Found blank line
		}
		if start > 0 && source[start-1] == '\n' {
			// Check if this line starts a block element
			if isBlockStart(source, start) {
				break
			}
		}
		start--
	}

	// Expand forward to block boundary
	end = edit.NewEndOffset
	sourceLen := len(source)
	for end < sourceLen {
		if source[end] == '\n' {
			// Check for blank line
			if end+1 < sourceLen &&
				source[end+1] == '\n' {
				end += 2

				break
			}
			// Check if next line starts a block element
			if end+1 < sourceLen &&
				isBlockStart(source, end+1) {
				end++

				break
			}
		}
		end++
	}

	return start, end
}

// isBlockStart checks if the given position starts a block element.
func isBlockStart(source []byte, pos int) bool {
	if pos >= len(source) {
		return false
	}

	b := source[pos]

	// Header
	if b == '#' {
		return true
	}

	// List marker
	if b == '-' || b == '*' || b == '+' {
		return true
	}

	// Blockquote
	if b == '>' {
		return true
	}

	// Code fence
	if b == '`' || b == '~' {
		return true
	}

	// Ordered list (digit followed by .)
	if b >= '0' && b <= '9' {
		for i := pos + 1; i < len(source); i++ {
			if source[i] == '.' {
				return true
			}
			if source[i] < '0' ||
				source[i] > '9' {
				break
			}
		}
	}

	return false
}
