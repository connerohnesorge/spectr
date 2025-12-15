package markdown

import (
	"errors"
)

// ErrSkipChildren is a sentinel error that can be returned from a visitor method
// to skip traversal of the current node's children. The traversal will continue
// with the next sibling. This is NOT treated as an actual error.
//
//nolint:errname,revive,staticcheck // keeping legacy name SkipChildren for backward compat
var SkipChildren = errors.New("skip children")

// Visitor defines the interface for AST node visitors.
// Each method receives a typed node and returns an error to control traversal.
// Return nil to continue traversal, SkipChildren to skip children, or any other
// error to stop traversal immediately.
//
//nolint:revive // exported: interface methods are self-documenting
type Visitor interface {
	VisitDocument(*NodeDocument) error
	VisitSection(*NodeSection) error
	VisitRequirement(*NodeRequirement) error
	VisitScenario(*NodeScenario) error
	VisitParagraph(*NodeParagraph) error
	VisitList(*NodeList) error
	VisitListItem(*NodeListItem) error
	VisitCodeBlock(*NodeCodeBlock) error
	VisitBlockquote(*NodeBlockquote) error
	VisitText(*NodeText) error
	VisitStrong(*NodeStrong) error
	VisitEmphasis(*NodeEmphasis) error
	VisitStrikethrough(*NodeStrikethrough) error
	VisitCode(*NodeCode) error
	VisitLink(*NodeLink) error
	VisitLinkDef(*NodeLinkDef) error
	VisitWikilink(*NodeWikilink) error
}

// BaseVisitor provides no-op default implementations for all Visitor methods.
// Embed this struct in custom visitors to only override the methods you need.
type BaseVisitor struct{}

// VisitDocument is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitDocument(
	*NodeDocument,
) error {
	return nil
}

// VisitSection is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitSection(
	*NodeSection,
) error {
	return nil
}

// VisitRequirement is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitRequirement(
	*NodeRequirement,
) error {
	return nil
}

// VisitScenario is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitScenario(
	*NodeScenario,
) error {
	return nil
}

// VisitParagraph is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitParagraph(
	*NodeParagraph,
) error {
	return nil
}

// VisitList is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitList(
	*NodeList,
) error {
	return nil
}

// VisitListItem is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitListItem(
	*NodeListItem,
) error {
	return nil
}

// VisitCodeBlock is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitCodeBlock(
	*NodeCodeBlock,
) error {
	return nil
}

// VisitBlockquote is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitBlockquote(
	*NodeBlockquote,
) error {
	return nil
}

// VisitText is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitText(
	*NodeText,
) error {
	return nil
}

// VisitStrong is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitStrong(
	*NodeStrong,
) error {
	return nil
}

// VisitEmphasis is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitEmphasis(
	*NodeEmphasis,
) error {
	return nil
}

// VisitStrikethrough is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitStrikethrough(
	*NodeStrikethrough,
) error {
	return nil
}

// VisitCode is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitCode(
	*NodeCode,
) error {
	return nil
}

// VisitLink is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitLink(
	*NodeLink,
) error {
	return nil
}

// VisitLinkDef is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitLinkDef(
	*NodeLinkDef,
) error {
	return nil
}

// VisitWikilink is a no-op that returns nil (continue traversal).
func (BaseVisitor) VisitWikilink(
	*NodeWikilink,
) error {
	return nil
}

// Walk traverses the AST in pre-order depth-first order, calling the appropriate
// visitor method for each node. It handles the traversal logic including child
// recursion and error handling.
//
// If a visitor method returns SkipChildren, the children of that node are skipped
// but traversal continues with the next sibling.
//
// If a visitor method returns any other non-nil error, traversal stops immediately
// and that error is returned.
//
// Walk safely handles nil nodes by returning nil without calling any visitor methods.
//
//nolint:revive // function-length - visitor dispatch requires handling all node types
func Walk(node Node, v Visitor) error {
	if node == nil {
		return nil
	}

	// Call the appropriate visitor method based on node type
	var err error
	switch n := node.(type) {
	case *NodeDocument:
		err = v.VisitDocument(n)
	case *NodeSection:
		err = v.VisitSection(n)
	case *NodeRequirement:
		err = v.VisitRequirement(n)
	case *NodeScenario:
		err = v.VisitScenario(n)
	case *NodeParagraph:
		err = v.VisitParagraph(n)
	case *NodeList:
		err = v.VisitList(n)
	case *NodeListItem:
		err = v.VisitListItem(n)
	case *NodeCodeBlock:
		err = v.VisitCodeBlock(n)
	case *NodeBlockquote:
		err = v.VisitBlockquote(n)
	case *NodeText:
		err = v.VisitText(n)
	case *NodeStrong:
		err = v.VisitStrong(n)
	case *NodeEmphasis:
		err = v.VisitEmphasis(n)
	case *NodeStrikethrough:
		err = v.VisitStrikethrough(n)
	case *NodeCode:
		err = v.VisitCode(n)
	case *NodeLink:
		err = v.VisitLink(n)
	case *NodeLinkDef:
		err = v.VisitLinkDef(n)
	case *NodeWikilink:
		err = v.VisitWikilink(n)
	default:
		// Unknown node type - skip it
		return nil
	}

	// Handle visitor result
	if err != nil {
		if errors.Is(err, SkipChildren) {
			// Skip children but continue with siblings
			return nil
		}
		// Return any other error to stop traversal
		return err
	}

	// Recursively visit children
	for _, child := range node.Children() {
		if err := Walk(child, v); err != nil {
			return err
		}
	}

	return nil
}

// VisitorContext provides context information during traversal,
// including access to the parent node and current depth.
type VisitorContext struct {
	parent Node
	depth  int
}

// Parent returns the parent node of the current node being visited.
// Returns nil for the root node.
func (c *VisitorContext) Parent() Node {
	return c.parent
}

// Depth returns the current depth in the tree.
// The root node has depth 0.
func (c *VisitorContext) Depth() int {
	return c.depth
}

// ContextVisitor is a visitor interface that receives context information
// during traversal, including parent node access.
//
//nolint:revive // exported: interface methods are self-documenting
type ContextVisitor interface {
	VisitDocumentWithContext(
		*NodeDocument,
		*VisitorContext,
	) error
	VisitSectionWithContext(
		*NodeSection,
		*VisitorContext,
	) error
	VisitRequirementWithContext(
		*NodeRequirement,
		*VisitorContext,
	) error
	VisitScenarioWithContext(
		*NodeScenario,
		*VisitorContext,
	) error
	VisitParagraphWithContext(
		*NodeParagraph,
		*VisitorContext,
	) error
	VisitListWithContext(
		*NodeList,
		*VisitorContext,
	) error
	VisitListItemWithContext(
		*NodeListItem,
		*VisitorContext,
	) error
	VisitCodeBlockWithContext(
		*NodeCodeBlock,
		*VisitorContext,
	) error
	VisitBlockquoteWithContext(
		*NodeBlockquote,
		*VisitorContext,
	) error
	VisitTextWithContext(
		*NodeText,
		*VisitorContext,
	) error
	VisitStrongWithContext(
		*NodeStrong,
		*VisitorContext,
	) error
	VisitEmphasisWithContext(
		*NodeEmphasis,
		*VisitorContext,
	) error
	VisitStrikethroughWithContext(
		*NodeStrikethrough,
		*VisitorContext,
	) error
	VisitCodeWithContext(
		*NodeCode,
		*VisitorContext,
	) error
	VisitLinkWithContext(
		*NodeLink,
		*VisitorContext,
	) error
	VisitLinkDefWithContext(
		*NodeLinkDef,
		*VisitorContext,
	) error
	VisitWikilinkWithContext(
		*NodeWikilink,
		*VisitorContext,
	) error
}

// BaseContextVisitor provides no-op defaults for all ContextVisitor methods.
type BaseContextVisitor struct{}

// VisitDocumentWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitDocumentWithContext(
	*NodeDocument,
	*VisitorContext,
) error {
	return nil
}

// VisitSectionWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitSectionWithContext(
	*NodeSection,
	*VisitorContext,
) error {
	return nil
}

// VisitRequirementWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitRequirementWithContext(
	*NodeRequirement,
	*VisitorContext,
) error {
	return nil
}

// VisitScenarioWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitScenarioWithContext(
	*NodeScenario,
	*VisitorContext,
) error {
	return nil
}

// VisitParagraphWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitParagraphWithContext(
	*NodeParagraph,
	*VisitorContext,
) error {
	return nil
}

// VisitListWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitListWithContext(
	*NodeList,
	*VisitorContext,
) error {
	return nil
}

// VisitListItemWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitListItemWithContext(
	*NodeListItem,
	*VisitorContext,
) error {
	return nil
}

// VisitCodeBlockWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitCodeBlockWithContext(
	*NodeCodeBlock,
	*VisitorContext,
) error {
	return nil
}

// VisitBlockquoteWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitBlockquoteWithContext(
	*NodeBlockquote,
	*VisitorContext,
) error {
	return nil
}

// VisitTextWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitTextWithContext(
	*NodeText,
	*VisitorContext,
) error {
	return nil
}

// VisitStrongWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitStrongWithContext(
	*NodeStrong,
	*VisitorContext,
) error {
	return nil
}

// VisitEmphasisWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitEmphasisWithContext(
	*NodeEmphasis,
	*VisitorContext,
) error {
	return nil
}

// VisitStrikethroughWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitStrikethroughWithContext(
	*NodeStrikethrough,
	*VisitorContext,
) error {
	return nil
}

// VisitCodeWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitCodeWithContext(
	*NodeCode,
	*VisitorContext,
) error {
	return nil
}

// VisitLinkWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitLinkWithContext(
	*NodeLink,
	*VisitorContext,
) error {
	return nil
}

// VisitLinkDefWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitLinkDefWithContext(
	*NodeLinkDef,
	*VisitorContext,
) error {
	return nil
}

// VisitWikilinkWithContext is a no-op that returns nil.
func (BaseContextVisitor) VisitWikilinkWithContext(
	*NodeWikilink,
	*VisitorContext,
) error {
	return nil
}

// WalkWithContext traverses the AST like Walk but provides context information
// including parent node access to the visitor.
func WalkWithContext(
	node Node,
	v ContextVisitor,
) error {
	return walkWithContextInternal(
		node,
		v,
		nil,
		0,
	)
}

//nolint:revive // function-length: visitor dispatch requires handling all node types
func walkWithContextInternal(
	node Node,
	v ContextVisitor,
	parent Node,
	depth int,
) error {
	if node == nil {
		return nil
	}

	ctx := &VisitorContext{
		parent: parent,
		depth:  depth,
	}

	// Call the appropriate visitor method based on node type
	var err error
	switch n := node.(type) {
	case *NodeDocument:
		err = v.VisitDocumentWithContext(n, ctx)
	case *NodeSection:
		err = v.VisitSectionWithContext(n, ctx)
	case *NodeRequirement:
		err = v.VisitRequirementWithContext(n, ctx)
	case *NodeScenario:
		err = v.VisitScenarioWithContext(n, ctx)
	case *NodeParagraph:
		err = v.VisitParagraphWithContext(n, ctx)
	case *NodeList:
		err = v.VisitListWithContext(n, ctx)
	case *NodeListItem:
		err = v.VisitListItemWithContext(n, ctx)
	case *NodeCodeBlock:
		err = v.VisitCodeBlockWithContext(n, ctx)
	case *NodeBlockquote:
		err = v.VisitBlockquoteWithContext(n, ctx)
	case *NodeText:
		err = v.VisitTextWithContext(n, ctx)
	case *NodeStrong:
		err = v.VisitStrongWithContext(n, ctx)
	case *NodeEmphasis:
		err = v.VisitEmphasisWithContext(n, ctx)
	case *NodeStrikethrough:
		err = v.VisitStrikethroughWithContext(n, ctx)
	case *NodeCode:
		err = v.VisitCodeWithContext(n, ctx)
	case *NodeLink:
		err = v.VisitLinkWithContext(n, ctx)
	case *NodeLinkDef:
		err = v.VisitLinkDefWithContext(n, ctx)
	case *NodeWikilink:
		err = v.VisitWikilinkWithContext(n, ctx)
	default:
		return nil
	}

	// Handle visitor result
	if err != nil {
		if errors.Is(err, SkipChildren) {
			return nil
		}

		return err
	}

	// Recursively visit children
	for _, child := range node.Children() {
		if err := walkWithContextInternal(child, v, node, depth+1); err != nil {
			return err
		}
	}

	return nil
}

// EnterLeaveVisitor defines the interface for visitors that need to perform
// actions both before (Enter) and after (Leave) visiting a node's children.
// This is useful for operations that need to maintain state or perform
// cleanup, such as building output or managing a stack.
//
//nolint:revive // exported: interface methods are self-documenting
type EnterLeaveVisitor interface {
	EnterDocument(*NodeDocument) error
	LeaveDocument(*NodeDocument) error
	EnterSection(*NodeSection) error
	LeaveSection(*NodeSection) error
	EnterRequirement(*NodeRequirement) error
	LeaveRequirement(*NodeRequirement) error
	EnterScenario(*NodeScenario) error
	LeaveScenario(*NodeScenario) error
	EnterParagraph(*NodeParagraph) error
	LeaveParagraph(*NodeParagraph) error
	EnterList(*NodeList) error
	LeaveList(*NodeList) error
	EnterListItem(*NodeListItem) error
	LeaveListItem(*NodeListItem) error
	EnterCodeBlock(*NodeCodeBlock) error
	LeaveCodeBlock(*NodeCodeBlock) error
	EnterBlockquote(*NodeBlockquote) error
	LeaveBlockquote(*NodeBlockquote) error
	EnterText(*NodeText) error
	LeaveText(*NodeText) error
	EnterStrong(*NodeStrong) error
	LeaveStrong(*NodeStrong) error
	EnterEmphasis(*NodeEmphasis) error
	LeaveEmphasis(*NodeEmphasis) error
	EnterStrikethrough(*NodeStrikethrough) error
	LeaveStrikethrough(*NodeStrikethrough) error
	EnterCode(*NodeCode) error
	LeaveCode(*NodeCode) error
	EnterLink(*NodeLink) error
	LeaveLink(*NodeLink) error
	EnterLinkDef(*NodeLinkDef) error
	LeaveLinkDef(*NodeLinkDef) error
	EnterWikilink(*NodeWikilink) error
	LeaveWikilink(*NodeWikilink) error
}

// BaseEnterLeaveVisitor provides no-op default implementations for all
// EnterLeaveVisitor methods. Embed this struct in custom visitors to only
// override the methods you need.
type BaseEnterLeaveVisitor struct{}

// EnterDocument is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterDocument(
	*NodeDocument,
) error {
	return nil
}

// LeaveDocument is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveDocument(
	*NodeDocument,
) error {
	return nil
}

// EnterSection is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterSection(
	*NodeSection,
) error {
	return nil
}

// LeaveSection is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveSection(
	*NodeSection,
) error {
	return nil
}

// EnterRequirement is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterRequirement(
	*NodeRequirement,
) error {
	return nil
}

// LeaveRequirement is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveRequirement(
	*NodeRequirement,
) error {
	return nil
}

// EnterScenario is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterScenario(
	*NodeScenario,
) error {
	return nil
}

// LeaveScenario is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveScenario(
	*NodeScenario,
) error {
	return nil
}

// EnterParagraph is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterParagraph(
	*NodeParagraph,
) error {
	return nil
}

// LeaveParagraph is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveParagraph(
	*NodeParagraph,
) error {
	return nil
}

// EnterList is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterList(
	*NodeList,
) error {
	return nil
}

// LeaveList is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveList(
	*NodeList,
) error {
	return nil
}

// EnterListItem is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterListItem(
	*NodeListItem,
) error {
	return nil
}

// LeaveListItem is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveListItem(
	*NodeListItem,
) error {
	return nil
}

// EnterCodeBlock is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterCodeBlock(
	*NodeCodeBlock,
) error {
	return nil
}

// LeaveCodeBlock is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveCodeBlock(
	*NodeCodeBlock,
) error {
	return nil
}

// EnterBlockquote is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterBlockquote(
	*NodeBlockquote,
) error {
	return nil
}

// LeaveBlockquote is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveBlockquote(
	*NodeBlockquote,
) error {
	return nil
}

// EnterText is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterText(
	*NodeText,
) error {
	return nil
}

// LeaveText is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveText(
	*NodeText,
) error {
	return nil
}

// EnterStrong is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterStrong(
	*NodeStrong,
) error {
	return nil
}

// LeaveStrong is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveStrong(
	*NodeStrong,
) error {
	return nil
}

// EnterEmphasis is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterEmphasis(
	*NodeEmphasis,
) error {
	return nil
}

// LeaveEmphasis is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveEmphasis(
	*NodeEmphasis,
) error {
	return nil
}

// EnterStrikethrough is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterStrikethrough(
	*NodeStrikethrough,
) error {
	return nil
}

// LeaveStrikethrough is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveStrikethrough(
	*NodeStrikethrough,
) error {
	return nil
}

// EnterCode is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterCode(
	*NodeCode,
) error {
	return nil
}

// LeaveCode is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveCode(
	*NodeCode,
) error {
	return nil
}

// EnterLink is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterLink(
	*NodeLink,
) error {
	return nil
}

// LeaveLink is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveLink(
	*NodeLink,
) error {
	return nil
}

// EnterLinkDef is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterLinkDef(
	*NodeLinkDef,
) error {
	return nil
}

// LeaveLinkDef is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveLinkDef(
	*NodeLinkDef,
) error {
	return nil
}

// EnterWikilink is a no-op that returns nil.
func (BaseEnterLeaveVisitor) EnterWikilink(
	*NodeWikilink,
) error {
	return nil
}

// LeaveWikilink is a no-op that returns nil.
func (BaseEnterLeaveVisitor) LeaveWikilink(
	*NodeWikilink,
) error {
	return nil
}

// WalkEnterLeave traverses the AST calling Enter methods before visiting children
// and Leave methods after visiting children.
//
// If an Enter method returns SkipChildren, the children are skipped but the
// corresponding Leave method is still called.
//
// If an Enter method returns any other non-nil error, traversal stops immediately
// and the Leave method is NOT called.
//
// If a Leave method returns a non-nil error, traversal stops immediately.
//
// WalkEnterLeave safely handles nil nodes by returning nil.
//
//nolint:revive // function-length: visitor dispatch requires handling all node types
func WalkEnterLeave(
	node Node,
	v EnterLeaveVisitor,
) error {
	if node == nil {
		return nil
	}

	// Call Enter method and determine if we should visit children
	skipChildren := false
	var err error

	switch n := node.(type) {
	case *NodeDocument:
		err = v.EnterDocument(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveDocument(n)

	case *NodeSection:
		err = v.EnterSection(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveSection(n)

	case *NodeRequirement:
		err = v.EnterRequirement(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveRequirement(n)

	case *NodeScenario:
		err = v.EnterScenario(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveScenario(n)

	case *NodeParagraph:
		err = v.EnterParagraph(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveParagraph(n)

	case *NodeList:
		err = v.EnterList(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveList(n)

	case *NodeListItem:
		err = v.EnterListItem(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveListItem(n)

	case *NodeCodeBlock:
		err = v.EnterCodeBlock(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveCodeBlock(n)

	case *NodeBlockquote:
		err = v.EnterBlockquote(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveBlockquote(n)

	case *NodeText:
		err = v.EnterText(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveText(n)

	case *NodeStrong:
		err = v.EnterStrong(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveStrong(n)

	case *NodeEmphasis:
		err = v.EnterEmphasis(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveEmphasis(n)

	case *NodeStrikethrough:
		err = v.EnterStrikethrough(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveStrikethrough(n)

	case *NodeCode:
		err = v.EnterCode(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveCode(n)

	case *NodeLink:
		err = v.EnterLink(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveLink(n)

	case *NodeLinkDef:
		err = v.EnterLinkDef(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveLinkDef(n)

	case *NodeWikilink:
		err = v.EnterWikilink(n)
		if err != nil {
			if errors.Is(err, SkipChildren) {
				skipChildren = true
			} else {
				return err
			}
		}
		if !skipChildren {
			for _, child := range node.Children() {
				if err := WalkEnterLeave(child, v); err != nil {
					return err
				}
			}
		}

		return v.LeaveWikilink(n)

	default:
		return nil
	}
}
