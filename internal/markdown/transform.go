//nolint:revive // file-length-limit: transform requires handling all node types
package markdown

// TransformAction signals the intended action for a node during transformation.
// It is returned by TransformVisitor methods to indicate whether to keep,
// replace, or delete the visited node.
type TransformAction uint8

const (
	// ActionKeep indicates the node should remain unchanged.
	// The returned Node value is ignored when this action is specified.
	ActionKeep TransformAction = iota

	// ActionReplace indicates the node should be replaced with the returned Node.
	// The returned Node will take the place of the original in the parent's children.
	ActionReplace

	// ActionDelete indicates the node should be removed from its parent's children.
	// The returned Node value is ignored when this action is specified.
	ActionDelete
)

// String returns a human-readable name for the transform action.
func (a TransformAction) String() string {
	switch a {
	case ActionKeep:
		return "Keep"
	case ActionReplace:
		return "Replace"
	case ActionDelete:
		return "Delete"
	default:
		return "Unknown"
	}
}

// TransformVisitor defines the interface for AST transformation via the visitor pattern.
// Each method receives a typed node and returns (replacement, action, error).
//
// The Transform function uses post-order traversal, so children are transformed
// before their parent. This allows parent transform methods to see the results
// of child transformations.
//
// When ActionKeep is returned, the replacement Node is ignored and the original
// is kept. When ActionReplace is returned, the replacement Node takes the place
// of the original. When ActionDelete is returned, the node is removed from its
// parent's children.
type TransformVisitor interface {
	TransformDocument(
		*NodeDocument,
	) (Node, TransformAction, error)
	TransformSection(
		*NodeSection,
	) (Node, TransformAction, error)
	TransformRequirement(
		*NodeRequirement,
	) (Node, TransformAction, error)
	TransformScenario(
		*NodeScenario,
	) (Node, TransformAction, error)
	TransformParagraph(
		*NodeParagraph,
	) (Node, TransformAction, error)
	TransformList(
		*NodeList,
	) (Node, TransformAction, error)
	TransformListItem(
		*NodeListItem,
	) (Node, TransformAction, error)
	TransformCodeBlock(
		*NodeCodeBlock,
	) (Node, TransformAction, error)
	TransformBlockquote(
		*NodeBlockquote,
	) (Node, TransformAction, error)
	TransformText(
		*NodeText,
	) (Node, TransformAction, error)
	TransformStrong(
		*NodeStrong,
	) (Node, TransformAction, error)
	TransformEmphasis(
		*NodeEmphasis,
	) (Node, TransformAction, error)
	TransformStrikethrough(
		*NodeStrikethrough,
	) (Node, TransformAction, error)
	TransformCode(
		*NodeCode,
	) (Node, TransformAction, error)
	TransformLink(
		*NodeLink,
	) (Node, TransformAction, error)
	TransformLinkDef(
		*NodeLinkDef,
	) (Node, TransformAction, error)
	TransformWikilink(
		*NodeWikilink,
	) (Node, TransformAction, error)
}

// BaseTransformVisitor provides default no-op implementations for all
// TransformVisitor methods. Each method returns (original, ActionKeep, nil)
// by default. Embed this struct in custom transform visitors to only
// override the methods you need.
type BaseTransformVisitor struct{}

// TransformDocument returns the document unchanged.
func (BaseTransformVisitor) TransformDocument(
	n *NodeDocument,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformSection returns the section unchanged.
func (BaseTransformVisitor) TransformSection(
	n *NodeSection,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformRequirement returns the requirement unchanged.
func (BaseTransformVisitor) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformScenario returns the scenario unchanged.
func (BaseTransformVisitor) TransformScenario(
	n *NodeScenario,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformParagraph returns the paragraph unchanged.
func (BaseTransformVisitor) TransformParagraph(
	n *NodeParagraph,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformList returns the list unchanged.
func (BaseTransformVisitor) TransformList(
	n *NodeList,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformListItem returns the list item unchanged.
func (BaseTransformVisitor) TransformListItem(
	n *NodeListItem,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformCodeBlock returns the code block unchanged.
func (BaseTransformVisitor) TransformCodeBlock(
	n *NodeCodeBlock,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformBlockquote returns the blockquote unchanged.
func (BaseTransformVisitor) TransformBlockquote(
	n *NodeBlockquote,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformText returns the text unchanged.
func (BaseTransformVisitor) TransformText(
	n *NodeText,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformStrong returns the strong unchanged.
func (BaseTransformVisitor) TransformStrong(
	n *NodeStrong,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformEmphasis returns the emphasis unchanged.
func (BaseTransformVisitor) TransformEmphasis(
	n *NodeEmphasis,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformStrikethrough returns the strikethrough unchanged.
func (BaseTransformVisitor) TransformStrikethrough(
	n *NodeStrikethrough,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformCode returns the code unchanged.
func (BaseTransformVisitor) TransformCode(
	n *NodeCode,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformLink returns the link unchanged.
func (BaseTransformVisitor) TransformLink(
	n *NodeLink,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformLinkDef returns the link definition unchanged.
func (BaseTransformVisitor) TransformLinkDef(
	n *NodeLinkDef,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// TransformWikilink returns the wikilink unchanged.
func (BaseTransformVisitor) TransformWikilink(
	n *NodeWikilink,
) (Node, TransformAction, error) {
	return n, ActionKeep, nil
}

// Transform applies a TransformVisitor to an AST using post-order traversal.
// Children are transformed before their parent, so parent transform methods
// see the results of child transformations.
//
// The original AST is not modified; Transform returns a new tree with the
// transformations applied.
//
// If the root node is deleted (ActionDelete returned for root), Transform
// returns (nil, nil).
//
// If any transform method returns a non-nil error, Transform stops immediately
// and returns the error. Partial results are not returned.
func Transform(
	root Node,
	v TransformVisitor,
) (Node, error) {
	if root == nil {
		return nil, nil
	}

	result, action, err := transformNode(root, v)
	if err != nil {
		return nil, err
	}

	if action == ActionDelete {
		return nil, nil
	}

	return result, nil
}

// transformNode recursively transforms a node and its children in post-order.
// Returns the transformed node, the action to take, and any error.
//
//nolint:revive // function-length - transform logic requires handling all action types
func transformNode(
	node Node,
	v TransformVisitor,
) (Node, TransformAction, error) {
	if node == nil {
		return nil, ActionKeep, nil
	}

	// First, transform all children (post-order: children before parent)
	children := node.Children()
	var newChildren []Node
	childrenChanged := false

	for _, child := range children {
		transformedChild, action, err := transformNode(
			child,
			v,
		)
		if err != nil {
			return nil, ActionKeep, err
		}

		switch action {
		case ActionKeep:
			if newChildren != nil {
				newChildren = append(
					newChildren,
					child,
				)
			}
		case ActionReplace:
			if newChildren == nil {
				// Copy children up to this point
				newChildren = make(
					[]Node,
					0,
					len(children),
				)
				for _, prev := range children {
					if prev == child {
						break
					}
					newChildren = append(
						newChildren,
						prev,
					)
				}
			}
			newChildren = append(
				newChildren,
				transformedChild,
			)
			childrenChanged = true
		case ActionDelete:
			if newChildren == nil {
				// Copy children up to this point, excluding this child
				newChildren = make(
					[]Node,
					0,
					len(children)-1,
				)
				for _, prev := range children {
					if prev == child {
						break
					}
					newChildren = append(
						newChildren,
						prev,
					)
				}
			}
			childrenChanged = true
			// Don't append deleted child
		}
	}

	// If children changed, create a new node with updated children
	var nodeToTransform Node
	if childrenChanged {
		nodeToTransform = rebuildWithChildren(
			node,
			newChildren,
		)
	} else {
		nodeToTransform = node
	}

	// Now apply the transform to this node
	return applyTransform(nodeToTransform, v)
}

// applyTransform calls the appropriate transform method based on node type.
func applyTransform(
	node Node,
	v TransformVisitor,
) (Node, TransformAction, error) {
	switch n := node.(type) {
	case *NodeDocument:
		return v.TransformDocument(n)
	case *NodeSection:
		return v.TransformSection(n)
	case *NodeRequirement:
		return v.TransformRequirement(n)
	case *NodeScenario:
		return v.TransformScenario(n)
	case *NodeParagraph:
		return v.TransformParagraph(n)
	case *NodeList:
		return v.TransformList(n)
	case *NodeListItem:
		return v.TransformListItem(n)
	case *NodeCodeBlock:
		return v.TransformCodeBlock(n)
	case *NodeBlockquote:
		return v.TransformBlockquote(n)
	case *NodeText:
		return v.TransformText(n)
	case *NodeStrong:
		return v.TransformStrong(n)
	case *NodeEmphasis:
		return v.TransformEmphasis(n)
	case *NodeStrikethrough:
		return v.TransformStrikethrough(n)
	case *NodeCode:
		return v.TransformCode(n)
	case *NodeLink:
		return v.TransformLink(n)
	case *NodeLinkDef:
		return v.TransformLinkDef(n)
	case *NodeWikilink:
		return v.TransformWikilink(n)
	default:
		// Unknown node type - keep as-is
		return node, ActionKeep, nil
	}
}

// rebuildWithChildren creates a new node with the same properties but new children.
func rebuildWithChildren(
	node Node,
	newChildren []Node,
) Node {
	builder := nodeToBuilder(node)
	if builder == nil {
		return node
	}
	builder.WithChildren(newChildren)
	result := builder.Build()
	if result == nil {
		return node
	}

	return result
}

// Compose creates a TransformVisitor that applies t1 first, then t2.
// The result of t1 is passed to t2, allowing sequential transformations.
func Compose(
	t1, t2 TransformVisitor,
) TransformVisitor {
	return &composedTransform{t1: t1, t2: t2}
}

type composedTransform struct {
	BaseTransformVisitor
	t1, t2 TransformVisitor
}

func (c *composedTransform) TransformDocument(
	n *NodeDocument,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformDocument,
		c.t2.TransformDocument,
	)
}

func (c *composedTransform) TransformSection(
	n *NodeSection,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformSection,
		c.t2.TransformSection,
	)
}

func (c *composedTransform) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformRequirement,
		c.t2.TransformRequirement,
	)
}

func (c *composedTransform) TransformScenario(
	n *NodeScenario,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformScenario,
		c.t2.TransformScenario,
	)
}

func (c *composedTransform) TransformParagraph(
	n *NodeParagraph,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformParagraph,
		c.t2.TransformParagraph,
	)
}

func (c *composedTransform) TransformList(
	n *NodeList,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformList,
		c.t2.TransformList,
	)
}

func (c *composedTransform) TransformListItem(
	n *NodeListItem,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformListItem,
		c.t2.TransformListItem,
	)
}

func (c *composedTransform) TransformCodeBlock(
	n *NodeCodeBlock,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformCodeBlock,
		c.t2.TransformCodeBlock,
	)
}

func (c *composedTransform) TransformBlockquote(
	n *NodeBlockquote,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformBlockquote,
		c.t2.TransformBlockquote,
	)
}

func (c *composedTransform) TransformText(
	n *NodeText,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformText,
		c.t2.TransformText,
	)
}

func (c *composedTransform) TransformStrong(
	n *NodeStrong,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformStrong,
		c.t2.TransformStrong,
	)
}

func (c *composedTransform) TransformEmphasis(
	n *NodeEmphasis,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformEmphasis,
		c.t2.TransformEmphasis,
	)
}

func (c *composedTransform) TransformStrikethrough(
	n *NodeStrikethrough,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformStrikethrough,
		c.t2.TransformStrikethrough,
	)
}

func (c *composedTransform) TransformCode(
	n *NodeCode,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformCode,
		c.t2.TransformCode,
	)
}

func (c *composedTransform) TransformLink(
	n *NodeLink,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformLink,
		c.t2.TransformLink,
	)
}

func (c *composedTransform) TransformLinkDef(
	n *NodeLinkDef,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformLinkDef,
		c.t2.TransformLinkDef,
	)
}

func (c *composedTransform) TransformWikilink(
	n *NodeWikilink,
) (Node, TransformAction, error) {
	return composeTransform(
		n,
		c.t1.TransformWikilink,
		c.t2.TransformWikilink,
	)
}

// composeTransform applies two transforms in sequence.
func composeTransform[T Node](
	n T,
	f1, f2 func(T) (Node, TransformAction, error),
) (Node, TransformAction, error) {
	result1, action1, err1 := f1(n)
	if err1 != nil {
		return nil, ActionKeep, err1
	}

	// If t1 deletes, don't run t2
	if action1 == ActionDelete {
		return nil, ActionDelete, nil
	}

	// Get the node to pass to t2
	var nodeForT2 T
	if action1 == ActionReplace {
		var ok bool
		nodeForT2, ok = result1.(T)
		if !ok {
			// Type changed, just return t1's result
			return result1, action1, nil
		}
	} else {
		nodeForT2 = n
	}

	result2, action2, err2 := f2(nodeForT2)
	if err2 != nil {
		return nil, ActionKeep, err2
	}

	// Combine actions: delete takes precedence, then replace, then keep
	if action2 == ActionDelete {
		return nil, ActionDelete, nil
	}
	if action2 == ActionReplace {
		return result2, ActionReplace, nil
	}
	if action1 == ActionReplace {
		return result1, ActionReplace, nil
	}

	return n, ActionKeep, nil
}

// Pipeline creates a TransformVisitor that applies all transforms in order.
// This is a convenience wrapper around multiple Compose calls.
func Pipeline(
	transforms ...TransformVisitor,
) TransformVisitor {
	if len(transforms) == 0 {
		return &BaseTransformVisitor{}
	}
	if len(transforms) == 1 {
		return transforms[0]
	}

	result := transforms[0]
	for i := 1; i < len(transforms); i++ {
		result = Compose(result, transforms[i])
	}

	return result
}

// When creates a conditional TransformVisitor that only applies the given
// transform when the predicate returns true. Nodes not matching the predicate
// pass through unchanged (ActionKeep).
func When(
	pred func(Node) bool,
	transform TransformVisitor,
) TransformVisitor {
	return &conditionalTransform{
		pred:      pred,
		transform: transform,
	}
}

type conditionalTransform struct {
	BaseTransformVisitor
	pred      func(Node) bool
	transform TransformVisitor
}

func (c *conditionalTransform) TransformDocument(
	n *NodeDocument,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformDocument(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformSection(
	n *NodeSection,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformSection(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformRequirement(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformScenario(
	n *NodeScenario,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformScenario(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformParagraph(
	n *NodeParagraph,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformParagraph(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformList(
	n *NodeList,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformList(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformListItem(
	n *NodeListItem,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformListItem(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformCodeBlock(
	n *NodeCodeBlock,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformCodeBlock(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformBlockquote(
	n *NodeBlockquote,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformBlockquote(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformText(
	n *NodeText,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformText(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformStrong(
	n *NodeStrong,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformStrong(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformEmphasis(
	n *NodeEmphasis,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformEmphasis(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformStrikethrough(
	n *NodeStrikethrough,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformStrikethrough(
			n,
		)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformCode(
	n *NodeCode,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformCode(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformLink(
	n *NodeLink,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformLink(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformLinkDef(
	n *NodeLinkDef,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformLinkDef(n)
	}

	return n, ActionKeep, nil
}

func (c *conditionalTransform) TransformWikilink(
	n *NodeWikilink,
) (Node, TransformAction, error) {
	if c.pred(n) {
		return c.transform.TransformWikilink(n)
	}

	return n, ActionKeep, nil
}

// Map creates a TransformVisitor that applies the given function to every node.
// If f returns the same node (by pointer equality), it is treated as ActionKeep.
// Otherwise, it is treated as ActionReplace with the returned node.
func Map(f func(Node) Node) TransformVisitor {
	return &mapTransform{f: f}
}

type mapTransform struct {
	BaseTransformVisitor
	f func(Node) Node
}

func (m *mapTransform) applyMap(
	n Node,
) (Node, TransformAction, error) {
	result := m.f(n)
	if result == n {
		return n, ActionKeep, nil
	}

	return result, ActionReplace, nil
}

func (m *mapTransform) TransformDocument(
	n *NodeDocument,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformSection(
	n *NodeSection,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformScenario(
	n *NodeScenario,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformParagraph(
	n *NodeParagraph,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformList(
	n *NodeList,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformListItem(
	n *NodeListItem,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformCodeBlock(
	n *NodeCodeBlock,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformBlockquote(
	n *NodeBlockquote,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformText(
	n *NodeText,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformStrong(
	n *NodeStrong,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformEmphasis(
	n *NodeEmphasis,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformStrikethrough(
	n *NodeStrikethrough,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformCode(
	n *NodeCode,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformLink(
	n *NodeLink,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformLinkDef(
	n *NodeLinkDef,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

func (m *mapTransform) TransformWikilink(
	n *NodeWikilink,
) (Node, TransformAction, error) {
	return m.applyMap(n)
}

// Filter creates a TransformVisitor that deletes nodes where the predicate
// returns false. Nodes matching the predicate (returns true) are kept.
func Filter(
	pred func(Node) bool,
) TransformVisitor {
	return &filterTransform{pred: pred}
}

type filterTransform struct {
	BaseTransformVisitor
	pred func(Node) bool
}

func (f *filterTransform) applyFilter(
	n Node,
) (Node, TransformAction, error) {
	if f.pred(n) {
		return n, ActionKeep, nil
	}

	return nil, ActionDelete, nil
}

func (f *filterTransform) TransformDocument(
	n *NodeDocument,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformSection(
	n *NodeSection,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformScenario(
	n *NodeScenario,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformParagraph(
	n *NodeParagraph,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformList(
	n *NodeList,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformListItem(
	n *NodeListItem,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformCodeBlock(
	n *NodeCodeBlock,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformBlockquote(
	n *NodeBlockquote,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformText(
	n *NodeText,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformStrong(
	n *NodeStrong,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformEmphasis(
	n *NodeEmphasis,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformStrikethrough(
	n *NodeStrikethrough,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformCode(
	n *NodeCode,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformLink(
	n *NodeLink,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformLinkDef(
	n *NodeLinkDef,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

func (f *filterTransform) TransformWikilink(
	n *NodeWikilink,
) (Node, TransformAction, error) {
	return f.applyFilter(n)
}

// RenameRequirement creates a TransformVisitor that renames requirements
// matching oldName to newName. Only requirements with Name() == oldName
// are affected; other nodes pass through unchanged.
func RenameRequirement(
	oldName, newName string,
) TransformVisitor {
	return &renameRequirementTransform{
		oldName: oldName,
		newName: newName,
	}
}

type renameRequirementTransform struct {
	BaseTransformVisitor
	oldName string
	newName string
}

func (r *renameRequirementTransform) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	if n.Name() != r.oldName {
		return n, ActionKeep, nil
	}

	// Create a new requirement with the new name
	builder := n.ToBuilder()
	builder.WithName(r.newName)
	newNode := builder.Build()
	if newNode == nil {
		return n, ActionKeep, nil
	}

	return newNode, ActionReplace, nil
}

// AddScenario creates a TransformVisitor that adds a scenario to the requirement
// matching reqName. The scenario is appended to the requirement's children.
// Other nodes pass through unchanged.
func AddScenario(
	reqName string,
	scenario *NodeScenario,
) TransformVisitor {
	return &addScenarioTransform{
		reqName:  reqName,
		scenario: scenario,
	}
}

type addScenarioTransform struct {
	BaseTransformVisitor
	reqName  string
	scenario *NodeScenario
}

func (a *addScenarioTransform) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	if n.Name() != a.reqName {
		return n, ActionKeep, nil
	}

	// Create a new requirement with the scenario appended to children
	builder := n.ToBuilder()
	children := n.Children()
	newChildren := make([]Node, len(children)+1)
	copy(newChildren, children)
	newChildren[len(children)] = a.scenario
	builder.WithChildren(newChildren)
	newNode := builder.Build()
	if newNode == nil {
		return n, ActionKeep, nil
	}

	return newNode, ActionReplace, nil
}

// RemoveRequirement creates a TransformVisitor that deletes requirements
// matching the given name. The action is ActionDelete for matching requirements.
// Other nodes pass through unchanged.
func RemoveRequirement(
	name string,
) TransformVisitor {
	return &removeRequirementTransform{name: name}
}

type removeRequirementTransform struct {
	BaseTransformVisitor
	name string
}

func (r *removeRequirementTransform) TransformRequirement(
	n *NodeRequirement,
) (Node, TransformAction, error) {
	if n.Name() == r.name {
		return nil, ActionDelete, nil
	}

	return n, ActionKeep, nil
}
