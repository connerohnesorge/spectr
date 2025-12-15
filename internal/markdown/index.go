package markdown

import (
	"sort"
)

// PositionIndex provides efficient O(log n) queries for finding AST nodes
// at a given source position. It uses a sorted interval structure internally
// and is built lazily on the first query.
//
// The index becomes stale when the AST is modified. Use Rebuild() to
// update the index with a new AST root, or create a new index.
type PositionIndex struct {
	root      Node           // AST root node
	rootHash  uint64         // Hash of root for staleness detection
	intervals []nodeInterval // Sorted intervals for queries
	lineIndex *LineIndex     // For line/column calculations
	source    []byte         // Source text for LineIndex
	built     bool           // Whether the index has been built
}

// nodeInterval represents a node's position range in the source.
// Intervals are stored sorted by start position, with maxEnd tracking
// the maximum end position in the subtree for efficient pruning.
type nodeInterval struct {
	node  Node
	start int
	end   int
	depth int // Nesting depth (root = 0)
}

// NewPositionIndex creates a new PositionIndex for the given AST root.
// The index is built lazily on the first query.
// If root is nil, the index will be empty and queries return nil/empty.
func NewPositionIndex(
	root Node,
	source []byte,
) *PositionIndex {
	var rootHash uint64
	if root != nil {
		rootHash = root.Hash()
	}

	return &PositionIndex{
		root:     root,
		rootHash: rootHash,
		source:   source,
		built:    false,
	}
}

// build constructs the interval index from the AST.
// This is called lazily on the first query.
func (pi *PositionIndex) build() {
	if pi.built {
		return
	}

	// Build LineIndex lazily as well
	if pi.lineIndex == nil && pi.source != nil {
		pi.lineIndex = NewLineIndex(pi.source)
	}

	// Collect all intervals from the AST
	pi.intervals = nil
	if pi.root != nil {
		pi.collectIntervals(pi.root, 0)
	}

	// Sort by start position, then depth (deeper nodes last for tie-breaking)
	sort.Slice(
		pi.intervals,
		func(i, j int) bool {
			if pi.intervals[i].start != pi.intervals[j].start {
				return pi.intervals[i].start < pi.intervals[j].start
			}
			// For same start, put deeper (more specific) nodes later
			return pi.intervals[i].depth < pi.intervals[j].depth
		},
	)

	pi.built = true
}

// collectIntervals recursively collects all node intervals from the AST.
func (pi *PositionIndex) collectIntervals(
	node Node,
	depth int,
) {
	if node == nil {
		return
	}

	start, end := node.Span()
	pi.intervals = append(
		pi.intervals,
		nodeInterval{
			node:  node,
			start: start,
			end:   end,
			depth: depth,
		},
	)

	// Recursively collect from children
	for _, child := range node.Children() {
		pi.collectIntervals(child, depth+1)
	}
}

// IsStale returns true if the index may be stale due to AST modifications.
// This checks if the root node's hash has changed since the index was created.
func (pi *PositionIndex) IsStale(
	root Node,
) bool {
	if root == nil {
		return pi.root != nil
	}

	return root.Hash() != pi.rootHash
}

// Rebuild discards the current index and rebuilds it for the given root.
// This should be called when the AST has been modified.
func (pi *PositionIndex) Rebuild(
	root Node,
	source []byte,
) {
	var rootHash uint64
	if root != nil {
		rootHash = root.Hash()
	}

	pi.root = root
	pi.rootHash = rootHash
	pi.source = source
	pi.intervals = nil
	pi.lineIndex = nil
	pi.built = false
}

// NodeAt returns the innermost (most specific) node at the given offset.
// Returns nil if the offset is outside all nodes or if the AST is empty.
//
// Query time is O(log n + k) where k is the number of nodes at that position.
func (pi *PositionIndex) NodeAt(
	offset int,
) Node {
	pi.build()

	if len(pi.intervals) == 0 {
		return nil
	}

	// Find all nodes containing this offset
	var best Node
	var bestDepth = -1

	for i := range pi.intervals {
		iv := &pi.intervals[i]
		// Prune: if start > offset, no more intervals can contain offset
		if iv.start > offset {
			break
		}
		// Check if this interval contains the offset
		if iv.start > offset || offset >= iv.end {
			continue
		}

		// Prefer deeper (more specific) nodes
		if iv.depth > bestDepth {
			best = iv.node
			bestDepth = iv.depth
		}
	}

	return best
}

// NodesAt returns all nodes containing the given offset, ordered from
// outermost (root/document) to innermost (leaf/text).
// Returns an empty slice if the offset is outside all nodes.
//
// Query time is O(log n + k) where k is the number of matching nodes.
func (pi *PositionIndex) NodesAt(
	offset int,
) []Node {
	pi.build()

	if len(pi.intervals) == 0 {
		return nil
	}

	// Collect all nodes containing this offset
	var result []Node

	for i := range pi.intervals {
		iv := &pi.intervals[i]
		// Prune: if start > offset, no more intervals can contain offset
		if iv.start > offset {
			break
		}
		// Check if this interval contains the offset
		if iv.start <= offset && offset < iv.end {
			result = append(result, iv.node)
		}
	}

	// Sort by depth (outermost first)
	sort.Slice(result, func(i, j int) bool {
		// Find depths for each node
		var depthI, depthJ int
		for _, iv := range pi.intervals {
			if iv.node == result[i] {
				depthI = iv.depth
			}
			if iv.node == result[j] {
				depthJ = iv.depth
			}
		}

		return depthI < depthJ
	})

	return result
}

// NodesInRange returns all nodes overlapping the given [start, end) range.
// A node overlaps if any part of its range intersects with [start, end).
//
// The result includes nodes that:
// - Fully contain the range
// - Are fully contained within the range
// - Partially overlap the range
//
// Query time is O(log n + k) where k is the number of overlapping nodes.
func (pi *PositionIndex) NodesInRange(
	start, end int,
) []Node {
	pi.build()

	if len(pi.intervals) == 0 || start >= end {
		return nil
	}

	var result []Node

	// Use binary search to find the first interval that might overlap
	searchStart := sort.Search(
		len(pi.intervals),
		func(i int) bool {
			return pi.intervals[i].end > start
		},
	)

	for i := searchStart; i < len(pi.intervals); i++ {
		iv := &pi.intervals[i]

		// If interval starts at or after our end, we're done
		// (since intervals are sorted by start)
		if iv.start >= end {
			break
		}

		// Check for overlap: intervals overlap if one starts before other ends
		// [iv.start, iv.end) overlaps [start, end) if:
		// iv.start < end AND iv.end > start
		if iv.start < end && iv.end > start {
			result = append(result, iv.node)
		}
	}

	return result
}

// PositionAt returns the Position (line, column, offset) for the given byte offset.
// This uses the integrated LineIndex for efficient conversion.
//
// If no source was provided, returns a Position with just the offset.
//
//nolint:revive // line-length-limit: method description needs full context
func (pi *PositionIndex) PositionAt(
	offset int,
) Position {
	pi.build()

	if pi.lineIndex == nil {
		return Position{
			Line:   0,
			Column: 0,
			Offset: offset,
		}
	}

	return pi.lineIndex.PositionAt(offset)
}

// NodePosition returns the Position for the start of the given node.
// This is a convenience method equivalent to PositionAt(node.Start()).
func (pi *PositionIndex) NodePosition(
	node Node,
) Position {
	if node == nil {
		return Position{}
	}
	start, _ := node.Span()

	return pi.PositionAt(start)
}

// EnclosingSection returns the NodeSection containing the given offset.
// Returns nil if the offset is not within any section.
//
// This finds the innermost section if sections are nested.
func (pi *PositionIndex) EnclosingSection(
	offset int,
) *NodeSection {
	pi.build()

	if len(pi.intervals) == 0 {
		return nil
	}

	var best *NodeSection
	var bestDepth = -1

	for i := range pi.intervals {
		iv := &pi.intervals[i]
		// Prune: if start > offset, no more intervals can contain offset
		if iv.start > offset {
			break
		}
		// Check if this interval contains the offset and is a section
		if iv.start > offset || offset >= iv.end {
			continue
		}

		section, ok := iv.node.(*NodeSection)
		if !ok || iv.depth <= bestDepth {
			continue
		}

		best = section
		bestDepth = iv.depth
	}

	return best
}

// EnclosingRequirement returns the NodeRequirement containing the given offset.
// Returns nil if the offset is not within any requirement.
//
// This finds the innermost requirement if requirements are nested.
func (pi *PositionIndex) EnclosingRequirement(
	offset int,
) *NodeRequirement {
	pi.build()

	if len(pi.intervals) == 0 {
		return nil
	}

	var best *NodeRequirement
	var bestDepth = -1

	for i := range pi.intervals {
		iv := &pi.intervals[i]
		// Prune: if start > offset, no more intervals can contain offset
		if iv.start > offset {
			break
		}
		// Check if this interval contains the offset and is a requirement
		if iv.start > offset || offset >= iv.end {
			continue
		}

		req, ok := iv.node.(*NodeRequirement)
		if !ok || iv.depth <= bestDepth {
			continue
		}

		best = req
		bestDepth = iv.depth
	}

	return best
}

// EnclosingScenario returns the NodeScenario containing the given offset.
// Returns nil if the offset is not within any scenario.
func (pi *PositionIndex) EnclosingScenario(
	offset int,
) *NodeScenario {
	pi.build()

	if len(pi.intervals) == 0 {
		return nil
	}

	var best *NodeScenario
	var bestDepth = -1

	for i := range pi.intervals {
		iv := &pi.intervals[i]
		// Prune: if start > offset, no more intervals can contain offset
		if iv.start > offset {
			break
		}
		// Check if this interval contains the offset and is a scenario
		if iv.start > offset || offset >= iv.end {
			continue
		}

		scenario, ok := iv.node.(*NodeScenario)
		if !ok || iv.depth <= bestDepth {
			continue
		}

		best = scenario
		bestDepth = iv.depth
	}

	return best
}

// NodeCount returns the total number of nodes in the index.
// This triggers index building if not already built.
func (pi *PositionIndex) NodeCount() int {
	pi.build()

	return len(pi.intervals)
}

// Root returns the root node of the indexed AST.
func (pi *PositionIndex) Root() Node {
	return pi.root
}

// LineIndex returns the integrated LineIndex for line/column calculations.
// Returns nil if no source was provided.
func (pi *PositionIndex) LineIndex() *LineIndex {
	pi.build()

	return pi.lineIndex
}
