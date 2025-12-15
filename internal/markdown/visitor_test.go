package markdown

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestWalk_NilNode_ReturnsNil(t *testing.T) {
	v := &BaseVisitor{}
	err := Walk(nil, v)
	if err != nil {
		t.Errorf(
			"Walk(nil, v) = %v, want nil",
			err,
		)
	}
}

func TestWalk_VisitsAllNodesInPreOrder(
	t *testing.T,
) {
	// Build a tree:
	// Document
	//   Section
	//     Paragraph
	//       Text
	//     Paragraph
	//       Text
	textNode1 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("hello")).
		WithStart(0).WithEnd(5).
		Build()

	textNode2 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("world")).
		WithStart(10).WithEnd(15).
		Build()

	para1 := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode1}).
		WithStart(0).WithEnd(5).
		Build()

	para2 := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode2}).
		WithStart(10).WithEnd(15).
		Build()

	section := NewNodeBuilder(NodeTypeSection).
		WithLevel(2).
		WithTitle([]byte("Test")).
		WithChildren([]Node{para1, para2}).
		WithStart(0).WithEnd(15).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{section}).
		WithStart(0).WithEnd(15).
		Build()

	// Track visit order
	var visited []NodeType
	v := &recordingVisitor{
		onVisit: func(nt NodeType) error {
			visited = append(visited, nt)

			return nil
		},
	}

	err := Walk(doc, v)
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}

	expected := []NodeType{
		NodeTypeDocument,
		NodeTypeSection,
		NodeTypeParagraph,
		NodeTypeText,
		NodeTypeParagraph,
		NodeTypeText,
	}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf(
			"Visit order = %v, want %v",
			visited,
			expected,
		)
	}
}

func TestWalk_SkipChildrenSkipsSubtree(
	t *testing.T,
) {
	// Build a tree:
	// Document
	//   Section (will skip children)
	//     Paragraph
	//       Text
	//   Paragraph
	//     Text
	textNode1 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("skipped")).
		WithStart(0).WithEnd(7).
		Build()

	textNode2 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("visited")).
		WithStart(20).WithEnd(27).
		Build()

	para1 := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode1}).
		WithStart(0).WithEnd(7).
		Build()

	para2 := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode2}).
		WithStart(20).WithEnd(27).
		Build()

	section := NewNodeBuilder(NodeTypeSection).
		WithLevel(2).
		WithTitle([]byte("Skip")).
		WithChildren([]Node{para1}).
		WithStart(0).WithEnd(15).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{section, para2}).
		WithStart(0).WithEnd(27).
		Build()

	// Track visit order and skip section children
	var visited []NodeType
	v := &recordingVisitor{
		onVisit: func(nt NodeType) error {
			visited = append(visited, nt)
			if nt == NodeTypeSection {
				return SkipChildren
			}

			return nil
		},
	}

	err := Walk(doc, v)
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}

	// Section's children should be skipped
	expected := []NodeType{
		NodeTypeDocument,
		NodeTypeSection,
		NodeTypeParagraph, // para2
		NodeTypeText,      // textNode2
	}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf(
			"Visit order = %v, want %v",
			visited,
			expected,
		)
	}
}

func TestWalk_ErrorStopsTraversal(t *testing.T) {
	stopError := errors.New("stop here")

	textNode1 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("hello")).
		WithStart(0).WithEnd(5).
		Build()

	textNode2 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("world")).
		WithStart(10).WithEnd(15).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode1, textNode2}).
		WithStart(0).WithEnd(15).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(15).
		Build()

	var visited []NodeType
	v := &recordingVisitor{
		onVisit: func(nt NodeType) error {
			visited = append(visited, nt)
			if nt == NodeTypeText {
				return stopError
			}

			return nil
		},
	}

	err := Walk(doc, v)
	if !errors.Is(err, stopError) {
		t.Errorf(
			"Walk error = %v, want %v",
			err,
			stopError,
		)
	}

	// Should stop at first text node
	expected := []NodeType{
		NodeTypeDocument,
		NodeTypeParagraph,
		NodeTypeText,
	}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf(
			"Visit order = %v, want %v",
			visited,
			expected,
		)
	}
}

func TestWalk_CallsCorrectVisitMethod(
	t *testing.T,
) {
	// Test each node type gets the correct visitor method called
	testCases := []struct {
		name     string
		node     Node
		expected NodeType
	}{
		{
			name: "Document",
			node: NewNodeBuilder(
				NodeTypeDocument,
			).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeDocument,
		},
		{
			name: "Section",
			node: NewNodeBuilder(NodeTypeSection).
				WithLevel(1).
				WithTitle([]byte("Test")).
				WithStart(0).WithEnd(0).
				Build(),
			expected: NodeTypeSection,
		},
		{
			name: "Requirement",
			node: NewNodeBuilder(
				NodeTypeRequirement,
			).
				WithName("TestReq").
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeRequirement,
		},
		{
			name: "Scenario",
			node: NewNodeBuilder(
				NodeTypeScenario,
			).
				WithName("TestScenario").
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeScenario,
		},
		{
			name: "Paragraph",
			node: NewNodeBuilder(
				NodeTypeParagraph,
			).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeParagraph,
		},
		{
			name: "List",
			node: NewNodeBuilder(NodeTypeList).
				WithOrdered(false).
				WithStart(0).WithEnd(0).
				Build(),
			expected: NodeTypeList,
		},
		{
			name: "ListItem",
			node: NewNodeBuilder(
				NodeTypeListItem,
			).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeListItem,
		},
		{
			name: "CodeBlock",
			node: NewNodeBuilder(
				NodeTypeCodeBlock,
			).
				WithLanguage([]byte("go")).
				WithContent([]byte("code")).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeCodeBlock,
		},
		{
			name: "Blockquote",
			node: NewNodeBuilder(
				NodeTypeBlockquote,
			).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeBlockquote,
		},
		{
			name: "Text",
			node: NewNodeBuilder(NodeTypeText).
				WithSource([]byte("text")).
				WithStart(0).WithEnd(4).
				Build(),
			expected: NodeTypeText,
		},
		{
			name: "Strong",
			node: NewNodeBuilder(NodeTypeStrong).
				WithStart(0).WithEnd(0).
				Build(),
			expected: NodeTypeStrong,
		},
		{
			name: "Emphasis",
			node: NewNodeBuilder(
				NodeTypeEmphasis,
			).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeEmphasis,
		},
		{
			name: "Strikethrough",
			node: NewNodeBuilder(
				NodeTypeStrikethrough,
			).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeStrikethrough,
		},
		{
			name: "Code",
			node: NewNodeBuilder(NodeTypeCode).
				WithStart(0).WithEnd(0).
				Build(),
			expected: NodeTypeCode,
		},
		{
			name: "Link",
			node: NewNodeBuilder(NodeTypeLink).
				WithURL([]byte("http://example.com")).
				WithStart(0).WithEnd(0).
				Build(),
			expected: NodeTypeLink,
		},
		{
			name: "LinkDef",
			node: NewNodeBuilder(NodeTypeLinkDef).
				WithURL([]byte("http://example.com")).
				WithStart(0).WithEnd(0).
				Build(),
			expected: NodeTypeLinkDef,
		},
		{
			name: "Wikilink",
			node: NewNodeBuilder(
				NodeTypeWikilink,
			).
				WithTarget([]byte("page")).
				WithStart(0).
				WithEnd(0).
				Build(),
			expected: NodeTypeWikilink,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tv := &typeCheckingVisitor{}
			err := Walk(tc.node, tv)
			if err != nil {
				t.Fatalf(
					"Walk returned error: %v",
					err,
				)
			}

			if tv.lastVisited != tc.expected {
				t.Errorf(
					"Visited node type = %v, want %v",
					tv.lastVisited,
					tc.expected,
				)
			}
		})
	}
}

func TestBaseVisitor_AllMethodsReturnNil(
	t *testing.T,
) {
	v := BaseVisitor{}

	// Test all visitor methods return nil
	if err := v.VisitDocument(nil); err != nil {
		t.Errorf(
			"VisitDocument() = %v, want nil",
			err,
		)
	}
	if err := v.VisitSection(nil); err != nil {
		t.Errorf(
			"VisitSection() = %v, want nil",
			err,
		)
	}
	if err := v.VisitRequirement(nil); err != nil {
		t.Errorf(
			"VisitRequirement() = %v, want nil",
			err,
		)
	}
	if err := v.VisitScenario(nil); err != nil {
		t.Errorf(
			"VisitScenario() = %v, want nil",
			err,
		)
	}
	if err := v.VisitParagraph(nil); err != nil {
		t.Errorf(
			"VisitParagraph() = %v, want nil",
			err,
		)
	}
	if err := v.VisitList(nil); err != nil {
		t.Errorf(
			"VisitList() = %v, want nil",
			err,
		)
	}
	if err := v.VisitListItem(nil); err != nil {
		t.Errorf(
			"VisitListItem() = %v, want nil",
			err,
		)
	}
	if err := v.VisitCodeBlock(nil); err != nil {
		t.Errorf(
			"VisitCodeBlock() = %v, want nil",
			err,
		)
	}
	if err := v.VisitBlockquote(nil); err != nil {
		t.Errorf(
			"VisitBlockquote() = %v, want nil",
			err,
		)
	}
	if err := v.VisitText(nil); err != nil {
		t.Errorf(
			"VisitText() = %v, want nil",
			err,
		)
	}
	if err := v.VisitStrong(nil); err != nil {
		t.Errorf(
			"VisitStrong() = %v, want nil",
			err,
		)
	}
	if err := v.VisitEmphasis(nil); err != nil {
		t.Errorf(
			"VisitEmphasis() = %v, want nil",
			err,
		)
	}
	if err := v.VisitStrikethrough(nil); err != nil {
		t.Errorf(
			"VisitStrikethrough() = %v, want nil",
			err,
		)
	}
	if err := v.VisitCode(nil); err != nil {
		t.Errorf(
			"VisitCode() = %v, want nil",
			err,
		)
	}
	if err := v.VisitLink(nil); err != nil {
		t.Errorf(
			"VisitLink() = %v, want nil",
			err,
		)
	}
	if err := v.VisitLinkDef(nil); err != nil {
		t.Errorf(
			"VisitLinkDef() = %v, want nil",
			err,
		)
	}
	if err := v.VisitWikilink(nil); err != nil {
		t.Errorf(
			"VisitWikilink() = %v, want nil",
			err,
		)
	}
}

func TestBaseVisitor_EmbeddingWorksCorrectly(
	t *testing.T,
) {
	// Custom visitor that only overrides VisitText
	type countingVisitor struct {
		BaseVisitor
	}

	cv := &countingVisitor{}

	// Override only VisitText - this simulates what users would do
	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	// Using BaseVisitor, all visits should succeed (return nil)
	err := Walk(doc, cv)
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}
}

func TestVisitorContext_ParentReturnsCorrectParent(
	t *testing.T,
) {
	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	type parentRecord struct {
		nodeType   NodeType
		parentType NodeType
		parentNil  bool
	}
	var records []parentRecord

	cv := &parentTrackingVisitor{
		onVisit: func(nt NodeType, ctx *VisitorContext) error {
			pr := parentRecord{
				nodeType:  nt,
				parentNil: ctx.Parent() == nil,
			}
			if !pr.parentNil {
				pr.parentType = ctx.Parent().
					NodeType()
			}
			records = append(records, pr)

			return nil
		},
	}

	err := WalkWithContext(doc, cv)
	if err != nil {
		t.Fatalf(
			"WalkWithContext returned error: %v",
			err,
		)
	}

	expected := []parentRecord{
		{
			nodeType:  NodeTypeDocument,
			parentNil: true,
		},
		{
			nodeType:   NodeTypeParagraph,
			parentType: NodeTypeDocument,
		},
		{
			nodeType:   NodeTypeText,
			parentType: NodeTypeParagraph,
		},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Errorf(
			"Parent records = %+v, want %+v",
			records,
			expected,
		)
	}
}

func TestVisitorContext_DepthReturnsCorrectDepth(
	t *testing.T,
) {
	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	section := NewNodeBuilder(NodeTypeSection).
		WithLevel(2).
		WithTitle([]byte("Test")).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{section}).
		WithStart(0).WithEnd(4).
		Build()

	type depthRecord struct {
		nodeType NodeType
		depth    int
	}
	var records []depthRecord

	cv := &parentTrackingVisitor{
		onVisit: func(nt NodeType, ctx *VisitorContext) error {
			records = append(records, depthRecord{
				nodeType: nt,
				depth:    ctx.Depth(),
			})

			return nil
		},
	}

	err := WalkWithContext(doc, cv)
	if err != nil {
		t.Fatalf(
			"WalkWithContext returned error: %v",
			err,
		)
	}

	expected := []depthRecord{
		{nodeType: NodeTypeDocument, depth: 0},
		{nodeType: NodeTypeSection, depth: 1},
		{nodeType: NodeTypeParagraph, depth: 2},
		{nodeType: NodeTypeText, depth: 3},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Errorf(
			"Depth records = %+v, want %+v",
			records,
			expected,
		)
	}
}

func TestWalkWithContext_NilNode_ReturnsNil(
	t *testing.T,
) {
	cv := &BaseContextVisitor{}
	err := WalkWithContext(nil, cv)
	if err != nil {
		t.Errorf(
			"WalkWithContext(nil, v) = %v, want nil",
			err,
		)
	}
}

func TestWalkWithContext_SkipChildrenSkipsSubtree(
	t *testing.T,
) {
	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	var visited []NodeType
	cv := &parentTrackingVisitor{
		onVisit: func(nt NodeType, _ *VisitorContext) error {
			visited = append(visited, nt)
			if nt == NodeTypeParagraph {
				return SkipChildren
			}

			return nil
		},
	}

	err := WalkWithContext(doc, cv)
	if err != nil {
		t.Fatalf(
			"WalkWithContext returned error: %v",
			err,
		)
	}

	expected := []NodeType{
		NodeTypeDocument,
		NodeTypeParagraph,
		// Text should be skipped
	}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf(
			"Visit order = %v, want %v",
			visited,
			expected,
		)
	}
}

func TestWalkEnterLeave_CallsEnterBeforeChildren(
	t *testing.T,
) {
	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	var events []string
	ev := &enterLeaveRecordingVisitor{
		onEnter: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("enter:%s", nt),
			)

			return nil
		},
		onLeave: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("leave:%s", nt),
			)

			return nil
		},
	}

	err := WalkEnterLeave(doc, ev)
	if err != nil {
		t.Fatalf(
			"WalkEnterLeave returned error: %v",
			err,
		)
	}

	expected := []string{
		"enter:Document",
		"enter:Paragraph",
		"enter:Text",
		"leave:Text",
		"leave:Paragraph",
		"leave:Document",
	}

	if !reflect.DeepEqual(events, expected) {
		t.Errorf(
			"Events = %v, want %v",
			events,
			expected,
		)
	}
}

func TestWalkEnterLeave_CallsLeaveAfterChildren(
	t *testing.T,
) {
	// Same test as above but verify leave is called after all children
	textNode1 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("a")).
		WithStart(0).WithEnd(1).
		Build()

	textNode2 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("b")).
		WithStart(2).WithEnd(3).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode1, textNode2}).
		WithStart(0).WithEnd(3).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(3).
		Build()

	var events []string
	ev := &enterLeaveRecordingVisitor{
		onEnter: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("enter:%s", nt),
			)

			return nil
		},
		onLeave: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("leave:%s", nt),
			)

			return nil
		},
	}

	err := WalkEnterLeave(doc, ev)
	if err != nil {
		t.Fatalf(
			"WalkEnterLeave returned error: %v",
			err,
		)
	}

	// Both text nodes should be visited before paragraph leave
	expected := []string{
		"enter:Document",
		"enter:Paragraph",
		"enter:Text",
		"leave:Text",
		"enter:Text",
		"leave:Text",
		"leave:Paragraph",
		"leave:Document",
	}

	if !reflect.DeepEqual(events, expected) {
		t.Errorf(
			"Events = %v, want %v",
			events,
			expected,
		)
	}
}

func TestWalkEnterLeave_SkipChildrenInEnterSkipsChildrenButCallsLeave(
	t *testing.T,
) {
	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	var events []string
	ev := &enterLeaveRecordingVisitor{
		onEnter: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("enter:%s", nt),
			)
			if nt == NodeTypeParagraph {
				return SkipChildren
			}

			return nil
		},
		onLeave: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("leave:%s", nt),
			)

			return nil
		},
	}

	err := WalkEnterLeave(doc, ev)
	if err != nil {
		t.Fatalf(
			"WalkEnterLeave returned error: %v",
			err,
		)
	}

	// Children skipped but leave still called
	expected := []string{
		"enter:Document",
		"enter:Paragraph",
		"leave:Paragraph", // Leave called even though children skipped
		"leave:Document",
	}

	if !reflect.DeepEqual(events, expected) {
		t.Errorf(
			"Events = %v, want %v",
			events,
			expected,
		)
	}
}

func TestWalkEnterLeave_ErrorInEnterStopsWithoutCallingLeave(
	t *testing.T,
) {
	stopError := errors.New("stop in enter")

	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	var events []string
	ev := &enterLeaveRecordingVisitor{
		onEnter: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("enter:%s", nt),
			)
			if nt == NodeTypeParagraph {
				return stopError
			}

			return nil
		},
		onLeave: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("leave:%s", nt),
			)

			return nil
		},
	}

	err := WalkEnterLeave(doc, ev)
	if !errors.Is(err, stopError) {
		t.Errorf(
			"WalkEnterLeave error = %v, want %v",
			err,
			stopError,
		)
	}

	// Leave should NOT be called for paragraph or document
	expected := []string{
		"enter:Document",
		"enter:Paragraph",
		// No leave events - error stopped traversal
	}

	if !reflect.DeepEqual(events, expected) {
		t.Errorf(
			"Events = %v, want %v",
			events,
			expected,
		)
	}
}

func TestWalkEnterLeave_NilNode_ReturnsNil(
	t *testing.T,
) {
	ev := &BaseEnterLeaveVisitor{}
	err := WalkEnterLeave(nil, ev)
	if err != nil {
		t.Errorf(
			"WalkEnterLeave(nil, v) = %v, want nil",
			err,
		)
	}
}

func TestWalkEnterLeave_ErrorInLeaveStopsTraversal(
	t *testing.T,
) {
	stopError := errors.New("stop in leave")

	textNode := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(0).WithEnd(4).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{textNode}).
		WithStart(0).WithEnd(4).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(4).
		Build()

	var events []string
	ev := &enterLeaveRecordingVisitor{
		onEnter: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("enter:%s", nt),
			)

			return nil
		},
		onLeave: func(nt NodeType) error {
			events = append(
				events,
				fmt.Sprintf("leave:%s", nt),
			)
			if nt == NodeTypeText {
				return stopError
			}

			return nil
		},
	}

	err := WalkEnterLeave(doc, ev)
	if !errors.Is(err, stopError) {
		t.Errorf(
			"WalkEnterLeave error = %v, want %v",
			err,
			stopError,
		)
	}

	expected := []string{
		"enter:Document",
		"enter:Paragraph",
		"enter:Text",
		"leave:Text",
		// Stops here due to error
	}

	if !reflect.DeepEqual(events, expected) {
		t.Errorf(
			"Events = %v, want %v",
			events,
			expected,
		)
	}
}

//
//nolint:revive // cyclomatic: test function necessarily tests many methods
func TestBaseEnterLeaveVisitor_AllMethodsReturnNil(
	t *testing.T,
) {
	v := BaseEnterLeaveVisitor{}

	// Test all Enter methods
	if err := v.EnterDocument(nil); err != nil {
		t.Errorf(
			"EnterDocument() = %v, want nil",
			err,
		)
	}
	if err := v.EnterSection(nil); err != nil {
		t.Errorf(
			"EnterSection() = %v, want nil",
			err,
		)
	}
	if err := v.EnterRequirement(nil); err != nil {
		t.Errorf(
			"EnterRequirement() = %v, want nil",
			err,
		)
	}
	if err := v.EnterScenario(nil); err != nil {
		t.Errorf(
			"EnterScenario() = %v, want nil",
			err,
		)
	}
	if err := v.EnterParagraph(nil); err != nil {
		t.Errorf(
			"EnterParagraph() = %v, want nil",
			err,
		)
	}
	if err := v.EnterList(nil); err != nil {
		t.Errorf(
			"EnterList() = %v, want nil",
			err,
		)
	}
	if err := v.EnterListItem(nil); err != nil {
		t.Errorf(
			"EnterListItem() = %v, want nil",
			err,
		)
	}
	if err := v.EnterCodeBlock(nil); err != nil {
		t.Errorf(
			"EnterCodeBlock() = %v, want nil",
			err,
		)
	}
	if err := v.EnterBlockquote(nil); err != nil {
		t.Errorf(
			"EnterBlockquote() = %v, want nil",
			err,
		)
	}
	if err := v.EnterText(nil); err != nil {
		t.Errorf(
			"EnterText() = %v, want nil",
			err,
		)
	}
	if err := v.EnterStrong(nil); err != nil {
		t.Errorf(
			"EnterStrong() = %v, want nil",
			err,
		)
	}
	if err := v.EnterEmphasis(nil); err != nil {
		t.Errorf(
			"EnterEmphasis() = %v, want nil",
			err,
		)
	}
	if err := v.EnterStrikethrough(nil); err != nil {
		t.Errorf(
			"EnterStrikethrough() = %v, want nil",
			err,
		)
	}
	if err := v.EnterCode(nil); err != nil {
		t.Errorf(
			"EnterCode() = %v, want nil",
			err,
		)
	}
	if err := v.EnterLink(nil); err != nil {
		t.Errorf(
			"EnterLink() = %v, want nil",
			err,
		)
	}
	if err := v.EnterLinkDef(nil); err != nil {
		t.Errorf(
			"EnterLinkDef() = %v, want nil",
			err,
		)
	}
	if err := v.EnterWikilink(nil); err != nil {
		t.Errorf(
			"EnterWikilink() = %v, want nil",
			err,
		)
	}

	// Test all Leave methods
	if err := v.LeaveDocument(nil); err != nil {
		t.Errorf(
			"LeaveDocument() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveSection(nil); err != nil {
		t.Errorf(
			"LeaveSection() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveRequirement(nil); err != nil {
		t.Errorf(
			"LeaveRequirement() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveScenario(nil); err != nil {
		t.Errorf(
			"LeaveScenario() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveParagraph(nil); err != nil {
		t.Errorf(
			"LeaveParagraph() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveList(nil); err != nil {
		t.Errorf(
			"LeaveList() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveListItem(nil); err != nil {
		t.Errorf(
			"LeaveListItem() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveCodeBlock(nil); err != nil {
		t.Errorf(
			"LeaveCodeBlock() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveBlockquote(nil); err != nil {
		t.Errorf(
			"LeaveBlockquote() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveText(nil); err != nil {
		t.Errorf(
			"LeaveText() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveStrong(nil); err != nil {
		t.Errorf(
			"LeaveStrong() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveEmphasis(nil); err != nil {
		t.Errorf(
			"LeaveEmphasis() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveStrikethrough(nil); err != nil {
		t.Errorf(
			"LeaveStrikethrough() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveCode(nil); err != nil {
		t.Errorf(
			"LeaveCode() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveLink(nil); err != nil {
		t.Errorf(
			"LeaveLink() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveLinkDef(nil); err != nil {
		t.Errorf(
			"LeaveLinkDef() = %v, want nil",
			err,
		)
	}
	if err := v.LeaveWikilink(nil); err != nil {
		t.Errorf(
			"LeaveWikilink() = %v, want nil",
			err,
		)
	}
}

func TestIntegration_CollectAllRequirements(
	t *testing.T,
) {
	// Build a document with multiple requirements
	req1 := NewNodeBuilder(NodeTypeRequirement).
		WithName("Auth").
		WithStart(0).WithEnd(10).
		Build()

	req2 := NewNodeBuilder(NodeTypeRequirement).
		WithName("Validation").
		WithStart(20).WithEnd(30).
		Build()

	req3 := NewNodeBuilder(NodeTypeRequirement).
		WithName("Storage").
		WithStart(40).WithEnd(50).
		Build()

	section := NewNodeBuilder(NodeTypeSection).
		WithLevel(2).
		WithTitle([]byte("Features")).
		WithChildren([]Node{req1, req2, req3}).
		WithStart(0).WithEnd(50).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{section}).
		WithStart(0).WithEnd(50).
		Build()

	// Collect requirements using visitor
	collector := &requirementCollector{}
	err := Walk(doc, collector)
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}

	expected := []string{
		"Auth",
		"Validation",
		"Storage",
	}
	if !reflect.DeepEqual(
		collector.requirements,
		expected,
	) {
		t.Errorf(
			"Requirements = %v, want %v",
			collector.requirements,
			expected,
		)
	}
}

func TestIntegration_CountNodesByType(
	t *testing.T,
) {
	// Build a tree with various node types
	text1 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("hello")).
		WithStart(0).WithEnd(5).
		Build()

	text2 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("world")).
		WithStart(10).WithEnd(15).
		Build()

	text3 := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("test")).
		WithStart(20).WithEnd(24).
		Build()

	para1 := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{text1}).
		WithStart(0).WithEnd(5).
		Build()

	para2 := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{text2, text3}).
		WithStart(10).WithEnd(24).
		Build()

	section := NewNodeBuilder(NodeTypeSection).
		WithLevel(2).
		WithTitle([]byte("Test")).
		WithChildren([]Node{para1, para2}).
		WithStart(0).WithEnd(24).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{section}).
		WithStart(0).WithEnd(24).
		Build()

	counter := &nodeCounter{
		counts: make(map[NodeType]int),
	}
	err := Walk(doc, counter)
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}

	if counter.counts[NodeTypeDocument] != 1 {
		t.Errorf(
			"Document count = %d, want 1",
			counter.counts[NodeTypeDocument],
		)
	}
	if counter.counts[NodeTypeSection] != 1 {
		t.Errorf(
			"Section count = %d, want 1",
			counter.counts[NodeTypeSection],
		)
	}
	if counter.counts[NodeTypeParagraph] != 2 {
		t.Errorf(
			"Paragraph count = %d, want 2",
			counter.counts[NodeTypeParagraph],
		)
	}
	if counter.counts[NodeTypeText] != 3 {
		t.Errorf(
			"Text count = %d, want 3",
			counter.counts[NodeTypeText],
		)
	}
}

func TestIntegration_BuildPathFromRootToNode(
	t *testing.T,
) {
	// Build a tree
	text := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("target")).
		WithStart(0).WithEnd(6).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{text}).
		WithStart(0).WithEnd(6).
		Build()

	section := NewNodeBuilder(NodeTypeSection).
		WithLevel(2).
		WithTitle([]byte("Test")).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(6).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{section}).
		WithStart(0).WithEnd(6).
		Build()

	// Find path to text node using context visitor
	pathFinder := &pathFinder{
		targetType: NodeTypeText,
	}

	err := WalkWithContext(doc, pathFinder)
	if err != nil &&
		!errors.Is(err, errFoundTarget) {
		t.Fatalf(
			"WalkWithContext returned error: %v",
			err,
		)
	}

	expectedPath := []NodeType{
		NodeTypeDocument,
		NodeTypeSection,
		NodeTypeParagraph,
		NodeTypeText,
	}

	if !reflect.DeepEqual(
		pathFinder.path,
		expectedPath,
	) {
		t.Errorf(
			"Path = %v, want %v",
			pathFinder.path,
			expectedPath,
		)
	}
}

func TestIntegration_BuildHTMLLikeOutput(
	t *testing.T,
) {
	// Test using enter/leave to build output
	text := NewNodeBuilder(NodeTypeText).
		WithSource([]byte("Hello")).
		WithStart(0).WithEnd(5).
		Build()

	strong := NewNodeBuilder(NodeTypeStrong).
		WithChildren([]Node{text}).
		WithStart(0).WithEnd(5).
		Build()

	para := NewNodeBuilder(NodeTypeParagraph).
		WithChildren([]Node{strong}).
		WithStart(0).WithEnd(5).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithChildren([]Node{para}).
		WithStart(0).WithEnd(5).
		Build()

	builder := &htmlBuilder{}
	err := WalkEnterLeave(doc, builder)
	if err != nil {
		t.Fatalf(
			"WalkEnterLeave returned error: %v",
			err,
		)
	}

	expected := "<doc><p><strong><text>Hello</text></strong></p></doc>"
	if builder.output != expected {
		t.Errorf(
			"Output = %q, want %q",
			builder.output,
			expected,
		)
	}
}

func TestBaseContextVisitor_AllMethodsReturnNil(
	t *testing.T,
) {
	v := BaseContextVisitor{}
	ctx := &VisitorContext{}

	if err := v.VisitDocumentWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitDocumentWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitSectionWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitSectionWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitRequirementWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitRequirementWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitScenarioWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitScenarioWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitParagraphWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitParagraphWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitListWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitListWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitListItemWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitListItemWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitCodeBlockWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitCodeBlockWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitBlockquoteWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitBlockquoteWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitTextWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitTextWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitStrongWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitStrongWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitEmphasisWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitEmphasisWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitStrikethroughWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitStrikethroughWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitCodeWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitCodeWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitLinkWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitLinkWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitLinkDefWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitLinkDefWithContext() = %v, want nil",
			err,
		)
	}
	if err := v.VisitWikilinkWithContext(nil, ctx); err != nil {
		t.Errorf(
			"VisitWikilinkWithContext() = %v, want nil",
			err,
		)
	}
}

// recordingVisitor records node types as they are visited
type recordingVisitor struct {
	BaseVisitor
	onVisit func(NodeType) error
}

func (v *recordingVisitor) VisitDocument(
	_ *NodeDocument,
) error {
	return v.onVisit(NodeTypeDocument)
}

func (v *recordingVisitor) VisitSection(
	_ *NodeSection,
) error {
	return v.onVisit(NodeTypeSection)
}

func (v *recordingVisitor) VisitRequirement(
	_ *NodeRequirement,
) error {
	return v.onVisit(NodeTypeRequirement)
}

func (v *recordingVisitor) VisitScenario(
	_ *NodeScenario,
) error {
	return v.onVisit(NodeTypeScenario)
}

func (v *recordingVisitor) VisitParagraph(
	_ *NodeParagraph,
) error {
	return v.onVisit(NodeTypeParagraph)
}

func (v *recordingVisitor) VisitList(
	_ *NodeList,
) error {
	return v.onVisit(NodeTypeList)
}

func (v *recordingVisitor) VisitListItem(
	_ *NodeListItem,
) error {
	return v.onVisit(NodeTypeListItem)
}

func (v *recordingVisitor) VisitCodeBlock(
	_ *NodeCodeBlock,
) error {
	return v.onVisit(NodeTypeCodeBlock)
}

func (v *recordingVisitor) VisitBlockquote(
	_ *NodeBlockquote,
) error {
	return v.onVisit(NodeTypeBlockquote)
}

func (v *recordingVisitor) VisitText(
	_ *NodeText,
) error {
	return v.onVisit(NodeTypeText)
}

func (v *recordingVisitor) VisitStrong(
	_ *NodeStrong,
) error {
	return v.onVisit(NodeTypeStrong)
}

func (v *recordingVisitor) VisitEmphasis(
	_ *NodeEmphasis,
) error {
	return v.onVisit(NodeTypeEmphasis)
}

func (v *recordingVisitor) VisitStrikethrough(
	_ *NodeStrikethrough,
) error {
	return v.onVisit(NodeTypeStrikethrough)
}

func (v *recordingVisitor) VisitCode(
	_ *NodeCode,
) error {
	return v.onVisit(NodeTypeCode)
}

func (v *recordingVisitor) VisitLink(
	_ *NodeLink,
) error {
	return v.onVisit(NodeTypeLink)
}

func (v *recordingVisitor) VisitLinkDef(
	_ *NodeLinkDef,
) error {
	return v.onVisit(NodeTypeLinkDef)
}

func (v *recordingVisitor) VisitWikilink(
	_ *NodeWikilink,
) error {
	return v.onVisit(NodeTypeWikilink)
}

// typeCheckingVisitor records the last visited node type
type typeCheckingVisitor struct {
	BaseVisitor
	lastVisited NodeType
}

func (v *typeCheckingVisitor) VisitDocument(
	_ *NodeDocument,
) error {
	v.lastVisited = NodeTypeDocument

	return nil
}

func (v *typeCheckingVisitor) VisitSection(
	_ *NodeSection,
) error {
	v.lastVisited = NodeTypeSection

	return nil
}

func (v *typeCheckingVisitor) VisitRequirement(
	_ *NodeRequirement,
) error {
	v.lastVisited = NodeTypeRequirement

	return nil
}

func (v *typeCheckingVisitor) VisitScenario(
	_ *NodeScenario,
) error {
	v.lastVisited = NodeTypeScenario

	return nil
}

func (v *typeCheckingVisitor) VisitParagraph(
	_ *NodeParagraph,
) error {
	v.lastVisited = NodeTypeParagraph

	return nil
}

func (v *typeCheckingVisitor) VisitList(
	_ *NodeList,
) error {
	v.lastVisited = NodeTypeList

	return nil
}

func (v *typeCheckingVisitor) VisitListItem(
	_ *NodeListItem,
) error {
	v.lastVisited = NodeTypeListItem

	return nil
}

func (v *typeCheckingVisitor) VisitCodeBlock(
	_ *NodeCodeBlock,
) error {
	v.lastVisited = NodeTypeCodeBlock

	return nil
}

func (v *typeCheckingVisitor) VisitBlockquote(
	_ *NodeBlockquote,
) error {
	v.lastVisited = NodeTypeBlockquote

	return nil
}

func (v *typeCheckingVisitor) VisitText(
	_ *NodeText,
) error {
	v.lastVisited = NodeTypeText

	return nil
}

func (v *typeCheckingVisitor) VisitStrong(
	_ *NodeStrong,
) error {
	v.lastVisited = NodeTypeStrong

	return nil
}

func (v *typeCheckingVisitor) VisitEmphasis(
	_ *NodeEmphasis,
) error {
	v.lastVisited = NodeTypeEmphasis

	return nil
}

func (v *typeCheckingVisitor) VisitStrikethrough(
	_ *NodeStrikethrough,
) error {
	v.lastVisited = NodeTypeStrikethrough

	return nil
}

func (v *typeCheckingVisitor) VisitCode(
	_ *NodeCode,
) error {
	v.lastVisited = NodeTypeCode

	return nil
}

func (v *typeCheckingVisitor) VisitLink(
	_ *NodeLink,
) error {
	v.lastVisited = NodeTypeLink

	return nil
}

func (v *typeCheckingVisitor) VisitLinkDef(
	_ *NodeLinkDef,
) error {
	v.lastVisited = NodeTypeLinkDef

	return nil
}

func (v *typeCheckingVisitor) VisitWikilink(
	_ *NodeWikilink,
) error {
	v.lastVisited = NodeTypeWikilink

	return nil
}

// parentTrackingVisitor is a context visitor that tracks parent/depth
type parentTrackingVisitor struct {
	BaseContextVisitor
	onVisit func(NodeType, *VisitorContext) error
}

func (v *parentTrackingVisitor) VisitDocumentWithContext(
	_ *NodeDocument,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeDocument, ctx)
}

func (v *parentTrackingVisitor) VisitSectionWithContext(
	_ *NodeSection,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeSection, ctx)
}

func (v *parentTrackingVisitor) VisitRequirementWithContext(
	_ *NodeRequirement,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeRequirement, ctx)
}

func (v *parentTrackingVisitor) VisitScenarioWithContext(
	_ *NodeScenario,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeScenario, ctx)
}

func (v *parentTrackingVisitor) VisitParagraphWithContext(
	_ *NodeParagraph,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeParagraph, ctx)
}

func (v *parentTrackingVisitor) VisitListWithContext(
	_ *NodeList,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeList, ctx)
}

func (v *parentTrackingVisitor) VisitListItemWithContext(
	_ *NodeListItem,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeListItem, ctx)
}

func (v *parentTrackingVisitor) VisitCodeBlockWithContext(
	_ *NodeCodeBlock,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeCodeBlock, ctx)
}

func (v *parentTrackingVisitor) VisitBlockquoteWithContext(
	_ *NodeBlockquote,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeBlockquote, ctx)
}

func (v *parentTrackingVisitor) VisitTextWithContext(
	_ *NodeText,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeText, ctx)
}

func (v *parentTrackingVisitor) VisitStrongWithContext(
	_ *NodeStrong,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeStrong, ctx)
}

func (v *parentTrackingVisitor) VisitEmphasisWithContext(
	_ *NodeEmphasis,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeEmphasis, ctx)
}

func (v *parentTrackingVisitor) VisitStrikethroughWithContext(
	_ *NodeStrikethrough,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeStrikethrough, ctx)
}

func (v *parentTrackingVisitor) VisitCodeWithContext(
	_ *NodeCode,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeCode, ctx)
}

func (v *parentTrackingVisitor) VisitLinkWithContext(
	_ *NodeLink,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeLink, ctx)
}

func (v *parentTrackingVisitor) VisitLinkDefWithContext(
	_ *NodeLinkDef,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeLinkDef, ctx)
}

func (v *parentTrackingVisitor) VisitWikilinkWithContext(
	_ *NodeWikilink,
	ctx *VisitorContext,
) error {
	return v.onVisit(NodeTypeWikilink, ctx)
}

// enterLeaveRecordingVisitor records enter/leave events
type enterLeaveRecordingVisitor struct {
	BaseEnterLeaveVisitor
	onEnter func(NodeType) error
	onLeave func(NodeType) error
}

func (v *enterLeaveRecordingVisitor) EnterDocument(
	_ *NodeDocument,
) error {
	return v.onEnter(NodeTypeDocument)
}

func (v *enterLeaveRecordingVisitor) LeaveDocument(
	_ *NodeDocument,
) error {
	return v.onLeave(NodeTypeDocument)
}

func (v *enterLeaveRecordingVisitor) EnterSection(
	_ *NodeSection,
) error {
	return v.onEnter(NodeTypeSection)
}

func (v *enterLeaveRecordingVisitor) LeaveSection(
	_ *NodeSection,
) error {
	return v.onLeave(NodeTypeSection)
}

func (v *enterLeaveRecordingVisitor) EnterRequirement(
	_ *NodeRequirement,
) error {
	return v.onEnter(NodeTypeRequirement)
}

func (v *enterLeaveRecordingVisitor) LeaveRequirement(
	_ *NodeRequirement,
) error {
	return v.onLeave(NodeTypeRequirement)
}

func (v *enterLeaveRecordingVisitor) EnterScenario(
	_ *NodeScenario,
) error {
	return v.onEnter(NodeTypeScenario)
}

func (v *enterLeaveRecordingVisitor) LeaveScenario(
	_ *NodeScenario,
) error {
	return v.onLeave(NodeTypeScenario)
}

func (v *enterLeaveRecordingVisitor) EnterParagraph(
	_ *NodeParagraph,
) error {
	return v.onEnter(NodeTypeParagraph)
}

func (v *enterLeaveRecordingVisitor) LeaveParagraph(
	_ *NodeParagraph,
) error {
	return v.onLeave(NodeTypeParagraph)
}

func (v *enterLeaveRecordingVisitor) EnterList(
	_ *NodeList,
) error {
	return v.onEnter(NodeTypeList)
}

func (v *enterLeaveRecordingVisitor) LeaveList(
	_ *NodeList,
) error {
	return v.onLeave(NodeTypeList)
}

func (v *enterLeaveRecordingVisitor) EnterListItem(
	_ *NodeListItem,
) error {
	return v.onEnter(NodeTypeListItem)
}

func (v *enterLeaveRecordingVisitor) LeaveListItem(
	_ *NodeListItem,
) error {
	return v.onLeave(NodeTypeListItem)
}

func (v *enterLeaveRecordingVisitor) EnterCodeBlock(
	_ *NodeCodeBlock,
) error {
	return v.onEnter(NodeTypeCodeBlock)
}

func (v *enterLeaveRecordingVisitor) LeaveCodeBlock(
	_ *NodeCodeBlock,
) error {
	return v.onLeave(NodeTypeCodeBlock)
}

func (v *enterLeaveRecordingVisitor) EnterBlockquote(
	_ *NodeBlockquote,
) error {
	return v.onEnter(NodeTypeBlockquote)
}

func (v *enterLeaveRecordingVisitor) LeaveBlockquote(
	_ *NodeBlockquote,
) error {
	return v.onLeave(NodeTypeBlockquote)
}

func (v *enterLeaveRecordingVisitor) EnterText(
	_ *NodeText,
) error {
	return v.onEnter(NodeTypeText)
}

func (v *enterLeaveRecordingVisitor) LeaveText(
	_ *NodeText,
) error {
	return v.onLeave(NodeTypeText)
}

func (v *enterLeaveRecordingVisitor) EnterStrong(
	_ *NodeStrong,
) error {
	return v.onEnter(NodeTypeStrong)
}

func (v *enterLeaveRecordingVisitor) LeaveStrong(
	_ *NodeStrong,
) error {
	return v.onLeave(NodeTypeStrong)
}

func (v *enterLeaveRecordingVisitor) EnterEmphasis(
	_ *NodeEmphasis,
) error {
	return v.onEnter(NodeTypeEmphasis)
}

func (v *enterLeaveRecordingVisitor) LeaveEmphasis(
	_ *NodeEmphasis,
) error {
	return v.onLeave(NodeTypeEmphasis)
}

func (v *enterLeaveRecordingVisitor) EnterStrikethrough(
	_ *NodeStrikethrough,
) error {
	return v.onEnter(NodeTypeStrikethrough)
}

func (v *enterLeaveRecordingVisitor) LeaveStrikethrough(
	_ *NodeStrikethrough,
) error {
	return v.onLeave(NodeTypeStrikethrough)
}

func (v *enterLeaveRecordingVisitor) EnterCode(
	_ *NodeCode,
) error {
	return v.onEnter(NodeTypeCode)
}

func (v *enterLeaveRecordingVisitor) LeaveCode(
	_ *NodeCode,
) error {
	return v.onLeave(NodeTypeCode)
}

func (v *enterLeaveRecordingVisitor) EnterLink(
	_ *NodeLink,
) error {
	return v.onEnter(NodeTypeLink)
}

func (v *enterLeaveRecordingVisitor) LeaveLink(
	_ *NodeLink,
) error {
	return v.onLeave(NodeTypeLink)
}

func (v *enterLeaveRecordingVisitor) EnterLinkDef(
	_ *NodeLinkDef,
) error {
	return v.onEnter(NodeTypeLinkDef)
}

func (v *enterLeaveRecordingVisitor) LeaveLinkDef(
	_ *NodeLinkDef,
) error {
	return v.onLeave(NodeTypeLinkDef)
}

func (v *enterLeaveRecordingVisitor) EnterWikilink(
	_ *NodeWikilink,
) error {
	return v.onEnter(NodeTypeWikilink)
}

func (v *enterLeaveRecordingVisitor) LeaveWikilink(
	_ *NodeWikilink,
) error {
	return v.onLeave(NodeTypeWikilink)
}

// requirementCollector collects requirement names
type requirementCollector struct {
	BaseVisitor
	requirements []string
}

func (v *requirementCollector) VisitRequirement(
	n *NodeRequirement,
) error {
	v.requirements = append(
		v.requirements,
		n.Name(),
	)

	return nil
}

// nodeCounter counts nodes by type
type nodeCounter struct {
	BaseVisitor
	counts map[NodeType]int
}

func (v *nodeCounter) VisitDocument(
	_ *NodeDocument,
) error {
	v.counts[NodeTypeDocument]++

	return nil
}

func (v *nodeCounter) VisitSection(
	_ *NodeSection,
) error {
	v.counts[NodeTypeSection]++

	return nil
}

func (v *nodeCounter) VisitRequirement(
	_ *NodeRequirement,
) error {
	v.counts[NodeTypeRequirement]++

	return nil
}

func (v *nodeCounter) VisitScenario(
	_ *NodeScenario,
) error {
	v.counts[NodeTypeScenario]++

	return nil
}

func (v *nodeCounter) VisitParagraph(
	_ *NodeParagraph,
) error {
	v.counts[NodeTypeParagraph]++

	return nil
}

func (v *nodeCounter) VisitList(
	_ *NodeList,
) error {
	v.counts[NodeTypeList]++

	return nil
}

func (v *nodeCounter) VisitListItem(
	_ *NodeListItem,
) error {
	v.counts[NodeTypeListItem]++

	return nil
}

func (v *nodeCounter) VisitCodeBlock(
	_ *NodeCodeBlock,
) error {
	v.counts[NodeTypeCodeBlock]++

	return nil
}

func (v *nodeCounter) VisitBlockquote(
	_ *NodeBlockquote,
) error {
	v.counts[NodeTypeBlockquote]++

	return nil
}

func (v *nodeCounter) VisitText(
	_ *NodeText,
) error {
	v.counts[NodeTypeText]++

	return nil
}

func (v *nodeCounter) VisitStrong(
	_ *NodeStrong,
) error {
	v.counts[NodeTypeStrong]++

	return nil
}

func (v *nodeCounter) VisitEmphasis(
	_ *NodeEmphasis,
) error {
	v.counts[NodeTypeEmphasis]++

	return nil
}

func (v *nodeCounter) VisitStrikethrough(
	_ *NodeStrikethrough,
) error {
	v.counts[NodeTypeStrikethrough]++

	return nil
}

func (v *nodeCounter) VisitCode(
	_ *NodeCode,
) error {
	v.counts[NodeTypeCode]++

	return nil
}

func (v *nodeCounter) VisitLink(
	_ *NodeLink,
) error {
	v.counts[NodeTypeLink]++

	return nil
}

func (v *nodeCounter) VisitLinkDef(
	_ *NodeLinkDef,
) error {
	v.counts[NodeTypeLinkDef]++

	return nil
}

func (v *nodeCounter) VisitWikilink(
	_ *NodeWikilink,
) error {
	v.counts[NodeTypeWikilink]++

	return nil
}

// errFoundTarget is used to stop traversal when target is found
var errFoundTarget = errors.New("found target")

// pathFinder finds path from root to a target node type
type pathFinder struct {
	BaseContextVisitor
	targetType NodeType
	path       []NodeType
	found      bool
}

func (v *pathFinder) buildPath(
	nt NodeType,
	_ *VisitorContext,
) error {
	if v.found {
		return SkipChildren
	}

	v.path = append(v.path, nt)

	if nt == v.targetType {
		v.found = true

		return errFoundTarget
	}

	return nil
}

func (v *pathFinder) VisitDocumentWithContext(
	_ *NodeDocument,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeDocument, ctx)
}

func (v *pathFinder) VisitSectionWithContext(
	_ *NodeSection,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeSection, ctx)
}

func (v *pathFinder) VisitRequirementWithContext(
	_ *NodeRequirement,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeRequirement, ctx)
}

func (v *pathFinder) VisitScenarioWithContext(
	_ *NodeScenario,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeScenario, ctx)
}

func (v *pathFinder) VisitParagraphWithContext(
	_ *NodeParagraph,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeParagraph, ctx)
}

func (v *pathFinder) VisitListWithContext(
	_ *NodeList,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeList, ctx)
}

func (v *pathFinder) VisitListItemWithContext(
	_ *NodeListItem,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeListItem, ctx)
}

func (v *pathFinder) VisitCodeBlockWithContext(
	_ *NodeCodeBlock,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeCodeBlock, ctx)
}

func (v *pathFinder) VisitBlockquoteWithContext(
	_ *NodeBlockquote,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeBlockquote, ctx)
}

func (v *pathFinder) VisitTextWithContext(
	_ *NodeText,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeText, ctx)
}

func (v *pathFinder) VisitStrongWithContext(
	_ *NodeStrong,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeStrong, ctx)
}

func (v *pathFinder) VisitEmphasisWithContext(
	_ *NodeEmphasis,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeEmphasis, ctx)
}

func (v *pathFinder) VisitStrikethroughWithContext(
	_ *NodeStrikethrough,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeStrikethrough, ctx)
}

func (v *pathFinder) VisitCodeWithContext(
	_ *NodeCode,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeCode, ctx)
}

func (v *pathFinder) VisitLinkWithContext(
	_ *NodeLink,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeLink, ctx)
}

func (v *pathFinder) VisitLinkDefWithContext(
	_ *NodeLinkDef,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeLinkDef, ctx)
}

func (v *pathFinder) VisitWikilinkWithContext(
	_ *NodeWikilink,
	ctx *VisitorContext,
) error {
	return v.buildPath(NodeTypeWikilink, ctx)
}

// htmlBuilder builds an HTML-like string representation
type htmlBuilder struct {
	BaseEnterLeaveVisitor
	output string
}

func (v *htmlBuilder) EnterDocument(
	_ *NodeDocument,
) error {
	v.output += "<doc>"

	return nil
}

func (v *htmlBuilder) LeaveDocument(
	_ *NodeDocument,
) error {
	v.output += "</doc>"

	return nil
}

func (v *htmlBuilder) EnterParagraph(
	_ *NodeParagraph,
) error {
	v.output += "<p>"

	return nil
}

func (v *htmlBuilder) LeaveParagraph(
	_ *NodeParagraph,
) error {
	v.output += "</p>"

	return nil
}

func (v *htmlBuilder) EnterStrong(
	_ *NodeStrong,
) error {
	v.output += "<strong>"

	return nil
}

func (v *htmlBuilder) LeaveStrong(
	_ *NodeStrong,
) error {
	v.output += "</strong>"

	return nil
}

func (v *htmlBuilder) EnterText(
	n *NodeText,
) error {
	v.output += "<text>" + string(n.Source())

	return nil
}

func (v *htmlBuilder) LeaveText(
	_ *NodeText,
) error {
	v.output += "</text>"

	return nil
}
