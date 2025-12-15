//nolint:revive // file-length-limit: object pools require comprehensive type-safe implementations
package markdown

import (
	"sync"
	"sync/atomic"
)

// Pool provides object pooling for tokens, nodes, and children slices.
// Using sync.Pool reduces GC pressure during parsing by reusing allocations.
// Each node type has its own pool for type safety.

// Token pool for reusing Token structs
var tokenPool = sync.Pool{
	New: func() any {
		return new(Token)
	},
}

// GetToken retrieves a Token from the pool.
// The returned token's fields are zeroed and ready for use.
func GetToken() *Token {
	if statsEnabled {
		atomic.AddUint64(&poolStats.TokenGets, 1)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return tokenPool.Get().(*Token)
}

// PutToken returns a Token to the pool after clearing its fields.
// The token should not be used after this call.
func PutToken(t *Token) {
	if t == nil {
		return
	}
	if statsEnabled {
		atomic.AddUint64(&poolStats.TokenPuts, 1)
	}
	// Clear all fields to prevent memory leaks
	t.Type = 0
	t.Start = 0
	t.End = 0
	t.Source = nil
	t.Message = ""
	tokenPool.Put(t)
}

// Node pools - one per node type for type safety

var documentPool = sync.Pool{
	New: func() any {
		return new(NodeDocument)
	},
}

var sectionPool = sync.Pool{
	New: func() any {
		return new(NodeSection)
	},
}

var requirementPool = sync.Pool{
	New: func() any {
		return new(NodeRequirement)
	},
}

var scenarioPool = sync.Pool{
	New: func() any {
		return new(NodeScenario)
	},
}

var paragraphPool = sync.Pool{
	New: func() any {
		return new(NodeParagraph)
	},
}

var listPool = sync.Pool{
	New: func() any {
		return new(NodeList)
	},
}

var listItemPool = sync.Pool{
	New: func() any {
		return new(NodeListItem)
	},
}

var codeBlockPool = sync.Pool{
	New: func() any {
		return new(NodeCodeBlock)
	},
}

var blockquotePool = sync.Pool{
	New: func() any {
		return new(NodeBlockquote)
	},
}

var textPool = sync.Pool{
	New: func() any {
		return new(NodeText)
	},
}

var strongPool = sync.Pool{
	New: func() any {
		return new(NodeStrong)
	},
}

var emphasisPool = sync.Pool{
	New: func() any {
		return new(NodeEmphasis)
	},
}

var strikethroughPool = sync.Pool{
	New: func() any {
		return new(NodeStrikethrough)
	},
}

var codePool = sync.Pool{
	New: func() any {
		return new(NodeCode)
	},
}

var linkPool = sync.Pool{
	New: func() any {
		return new(NodeLink)
	},
}

var linkDefPool = sync.Pool{
	New: func() any {
		return new(NodeLinkDef)
	},
}

var wikilinkPool = sync.Pool{
	New: func() any {
		return new(NodeWikilink)
	},
}

// GetDocument retrieves a NodeDocument from the pool.
func GetDocument() *NodeDocument {
	if statsEnabled {
		incrementNodeGets(NodeTypeDocument)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return documentPool.Get().(*NodeDocument)
}

// GetSection retrieves a NodeSection from the pool.
func GetSection() *NodeSection {
	if statsEnabled {
		incrementNodeGets(NodeTypeSection)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return sectionPool.Get().(*NodeSection)
}

// GetRequirement retrieves a NodeRequirement from the pool.
func GetRequirement() *NodeRequirement {
	if statsEnabled {
		incrementNodeGets(NodeTypeRequirement)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return requirementPool.Get().(*NodeRequirement)
}

// GetScenario retrieves a NodeScenario from the pool.
func GetScenario() *NodeScenario {
	if statsEnabled {
		incrementNodeGets(NodeTypeScenario)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return scenarioPool.Get().(*NodeScenario)
}

// GetParagraph retrieves a NodeParagraph from the pool.
func GetParagraph() *NodeParagraph {
	if statsEnabled {
		incrementNodeGets(NodeTypeParagraph)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return paragraphPool.Get().(*NodeParagraph)
}

// GetList retrieves a NodeList from the pool.
func GetList() *NodeList {
	if statsEnabled {
		incrementNodeGets(NodeTypeList)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return listPool.Get().(*NodeList)
}

// GetListItem retrieves a NodeListItem from the pool.
func GetListItem() *NodeListItem {
	if statsEnabled {
		incrementNodeGets(NodeTypeListItem)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return listItemPool.Get().(*NodeListItem)
}

// GetCodeBlock retrieves a NodeCodeBlock from the pool.
func GetCodeBlock() *NodeCodeBlock {
	if statsEnabled {
		incrementNodeGets(NodeTypeCodeBlock)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return codeBlockPool.Get().(*NodeCodeBlock)
}

// GetBlockquote retrieves a NodeBlockquote from the pool.
func GetBlockquote() *NodeBlockquote {
	if statsEnabled {
		incrementNodeGets(NodeTypeBlockquote)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return blockquotePool.Get().(*NodeBlockquote)
}

// GetText retrieves a NodeText from the pool.
func GetText() *NodeText {
	if statsEnabled {
		incrementNodeGets(NodeTypeText)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return textPool.Get().(*NodeText)
}

// GetStrong retrieves a NodeStrong from the pool.
func GetStrong() *NodeStrong {
	if statsEnabled {
		incrementNodeGets(NodeTypeStrong)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return strongPool.Get().(*NodeStrong)
}

// GetEmphasis retrieves a NodeEmphasis from the pool.
func GetEmphasis() *NodeEmphasis {
	if statsEnabled {
		incrementNodeGets(NodeTypeEmphasis)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return emphasisPool.Get().(*NodeEmphasis)
}

// GetStrikethrough retrieves a NodeStrikethrough from the pool.
func GetStrikethrough() *NodeStrikethrough {
	if statsEnabled {
		incrementNodeGets(NodeTypeStrikethrough)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return strikethroughPool.Get().(*NodeStrikethrough)
}

// GetCode retrieves a NodeCode from the pool.
func GetCode() *NodeCode {
	if statsEnabled {
		incrementNodeGets(NodeTypeCode)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return codePool.Get().(*NodeCode)
}

// GetLink retrieves a NodeLink from the pool.
func GetLink() *NodeLink {
	if statsEnabled {
		incrementNodeGets(NodeTypeLink)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return linkPool.Get().(*NodeLink)
}

// GetLinkDef retrieves a NodeLinkDef from the pool.
func GetLinkDef() *NodeLinkDef {
	if statsEnabled {
		incrementNodeGets(NodeTypeLinkDef)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return linkDefPool.Get().(*NodeLinkDef)
}

// GetWikilink retrieves a NodeWikilink from the pool.
func GetWikilink() *NodeWikilink {
	if statsEnabled {
		incrementNodeGets(NodeTypeWikilink)
	}
	//nolint:revive // unchecked-type-assertion - pool always returns correct type
	return wikilinkPool.Get().(*NodeWikilink)
}

// clearBaseNode clears the common baseNode fields.
func clearBaseNode(b *baseNode) {
	b.nodeType = 0
	b.hash = 0
	b.start = 0
	b.end = 0
	b.source = nil
	b.children = nil
}

// PutNode returns a node to the appropriate pool based on its type.
// The node should not be used after this call.
//
//nolint:revive // function-length - type switch covers all node types
func PutNode(n Node) {
	if n == nil {
		return
	}

	nodeType := n.NodeType()
	if statsEnabled {
		incrementNodePuts(nodeType)
	}

	switch node := n.(type) {
	case *NodeDocument:
		clearBaseNode(&node.baseNode)
		documentPool.Put(node)

	case *NodeSection:
		clearBaseNode(&node.baseNode)
		node.level = 0
		node.title = nil
		node.deltaType = ""
		sectionPool.Put(node)

	case *NodeRequirement:
		clearBaseNode(&node.baseNode)
		node.name = ""
		requirementPool.Put(node)

	case *NodeScenario:
		clearBaseNode(&node.baseNode)
		node.name = ""
		scenarioPool.Put(node)

	case *NodeParagraph:
		clearBaseNode(&node.baseNode)
		paragraphPool.Put(node)

	case *NodeList:
		clearBaseNode(&node.baseNode)
		node.ordered = false
		listPool.Put(node)

	case *NodeListItem:
		clearBaseNode(&node.baseNode)
		node.checked = nil
		node.keyword = ""
		listItemPool.Put(node)

	case *NodeCodeBlock:
		clearBaseNode(&node.baseNode)
		node.language = nil
		node.content = nil
		codeBlockPool.Put(node)

	case *NodeBlockquote:
		clearBaseNode(&node.baseNode)
		blockquotePool.Put(node)

	case *NodeText:
		clearBaseNode(&node.baseNode)
		textPool.Put(node)

	case *NodeStrong:
		clearBaseNode(&node.baseNode)
		strongPool.Put(node)

	case *NodeEmphasis:
		clearBaseNode(&node.baseNode)
		emphasisPool.Put(node)

	case *NodeStrikethrough:
		clearBaseNode(&node.baseNode)
		strikethroughPool.Put(node)

	case *NodeCode:
		clearBaseNode(&node.baseNode)
		codePool.Put(node)

	case *NodeLink:
		clearBaseNode(&node.baseNode)
		node.url = nil
		node.title = nil
		linkPool.Put(node)

	case *NodeLinkDef:
		clearBaseNode(&node.baseNode)
		node.url = nil
		node.title = nil
		linkDefPool.Put(node)

	case *NodeWikilink:
		clearBaseNode(&node.baseNode)
		node.target = nil
		node.display = nil
		node.anchor = nil
		wikilinkPool.Put(node)
	}
}

// Children slice pools - bucketed by capacity for efficient reuse
// Bucket sizes: 4 (small), 16 (medium), 64 (large)

const (
	childrenSmallCap  = 4
	childrenMediumCap = 16
	childrenLargeCap  = 64
)

var childrenSmallPool = sync.Pool{
	New: func() any {
		s := make([]Node, 0, childrenSmallCap)

		return &s
	},
}

var childrenMediumPool = sync.Pool{
	New: func() any {
		s := make([]Node, 0, childrenMediumCap)

		return &s
	},
}

var childrenLargePool = sync.Pool{
	New: func() any {
		s := make([]Node, 0, childrenLargeCap)

		return &s
	},
}

// GetChildren retrieves a children slice from the appropriate pool based on capacity.
// The returned slice has length 0 and capacity >= the requested capacity.
func GetChildren(capacity int) []Node {
	if statsEnabled {
		atomic.AddUint64(
			&poolStats.ChildrenGets,
			1,
		)
	}

	var slicePtr *[]Node

	switch {
	case capacity <= childrenSmallCap:
		//nolint:revive // unchecked-type-assertion - pool always returns correct type
		slicePtr = childrenSmallPool.Get().(*[]Node)
	case capacity <= childrenMediumCap:
		//nolint:revive // unchecked-type-assertion - pool always returns correct type
		slicePtr = childrenMediumPool.Get().(*[]Node)
	case capacity <= childrenLargeCap:
		//nolint:revive // unchecked-type-assertion - pool always returns correct type
		slicePtr = childrenLargePool.Get().(*[]Node)
	default:
		// For very large capacities, allocate directly
		s := make([]Node, 0, capacity)

		return s
	}

	// Reset length to 0 while preserving capacity
	*slicePtr = (*slicePtr)[:0]

	return *slicePtr
}

// PutChildren returns a children slice to the appropriate pool.
// The slice is cleared before being returned to prevent memory leaks.
func PutChildren(c []Node) {
	if c == nil {
		return
	}

	if statsEnabled {
		atomic.AddUint64(
			&poolStats.ChildrenPuts,
			1,
		)
	}

	sliceCap := cap(c)

	// Clear the slice to allow GC of referenced nodes
	for i := range c {
		c[i] = nil
	}
	c = c[:0]

	// Return to appropriate pool based on capacity
	switch sliceCap {
	case childrenSmallCap:
		childrenSmallPool.Put(&c)
	case childrenMediumCap:
		childrenMediumPool.Put(&c)
	case childrenLargeCap:
		childrenLargePool.Put(&c)
		// Slices with non-standard capacities are not returned to pools
	}
}

// Pool statistics tracking (optional, for debugging/profiling)

// PoolStats contains statistics about pool usage.
// All fields are accessed atomically when stats are enabled.
type PoolStats struct {
	TokenGets    uint64
	TokenPuts    uint64
	ChildrenGets uint64
	ChildrenPuts uint64
	// NodeGets and NodePuts are tracked per node type
	// Access via GetNodeStats()
}

// nodeStats tracks gets and puts per node type
type nodeStats struct {
	gets [18]uint64 // One per NodeType (0-17)
	puts [18]uint64
}

var (
	statsEnabled bool
	poolStats    PoolStats
	nodeStatsVal nodeStats
	statsMu      sync.RWMutex
)

// EnablePoolStats enables pool statistics tracking.
// This adds some overhead to pool operations.
func EnablePoolStats() {
	statsMu.Lock()
	statsEnabled = true
	statsMu.Unlock()
}

// DisablePoolStats disables pool statistics tracking.
func DisablePoolStats() {
	statsMu.Lock()
	statsEnabled = false
	statsMu.Unlock()
}

// GetPoolStats returns current pool statistics.
// Statistics are only collected when enabled via EnablePoolStats.
func GetPoolStats() PoolStats {
	statsMu.RLock()
	defer statsMu.RUnlock()

	return PoolStats{
		TokenGets: atomic.LoadUint64(
			&poolStats.TokenGets,
		),
		TokenPuts: atomic.LoadUint64(
			&poolStats.TokenPuts,
		),
		ChildrenGets: atomic.LoadUint64(
			&poolStats.ChildrenGets,
		),
		ChildrenPuts: atomic.LoadUint64(
			&poolStats.ChildrenPuts,
		),
	}
}

// NodeTypeStats contains get/put counts for a specific node type.
type NodeTypeStats struct {
	Type NodeType
	Gets uint64
	Puts uint64
}

// GetNodeStats returns statistics for all node types.
func GetNodeStats() []NodeTypeStats {
	statsMu.RLock()
	defer statsMu.RUnlock()

	result := make([]NodeTypeStats, 18)
	for i := range 18 {
		result[i] = NodeTypeStats{
			Type: NodeType(i),
			Gets: atomic.LoadUint64(
				&nodeStatsVal.gets[i],
			),
			Puts: atomic.LoadUint64(
				&nodeStatsVal.puts[i],
			),
		}
	}

	return result
}

// ResetPoolStats resets all pool statistics to zero.
func ResetPoolStats() {
	statsMu.Lock()
	defer statsMu.Unlock()

	atomic.StoreUint64(&poolStats.TokenGets, 0)
	atomic.StoreUint64(&poolStats.TokenPuts, 0)
	atomic.StoreUint64(&poolStats.ChildrenGets, 0)
	atomic.StoreUint64(&poolStats.ChildrenPuts, 0)

	for i := range 18 {
		atomic.StoreUint64(
			&nodeStatsVal.gets[i],
			0,
		)
		atomic.StoreUint64(
			&nodeStatsVal.puts[i],
			0,
		)
	}
}

// incrementNodeGets atomically increments the get counter for a node type.
func incrementNodeGets(t NodeType) {
	if int(t) < len(nodeStatsVal.gets) {
		atomic.AddUint64(&nodeStatsVal.gets[t], 1)
	}
}

// incrementNodePuts atomically increments the put counter for a node type.
func incrementNodePuts(t NodeType) {
	if int(t) < len(nodeStatsVal.puts) {
		atomic.AddUint64(&nodeStatsVal.puts[t], 1)
	}
}
