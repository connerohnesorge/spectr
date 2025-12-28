//nolint:revive // unchecked-type-assertion - panics acceptable in tests
package markdown

import (
	"bytes"
	"testing"
)

const (
	deltaCAdded    = "ADDED"
	deltaWhen      = "WHEN"
	deltaModified  = "MODIFIED"
	testExampleURL = "https://example.com"
)

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}

func TestNodeInterface_Document(t *testing.T) {
	source := []byte("# Test Document")
	node := NewNodeBuilder(NodeTypeDocument).
		WithStart(0).
		WithEnd(15).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeDocument {
		t.Errorf(
			"expected NodeType Document, got %v",
			node.NodeType(),
		)
	}

	start, end := node.Span()
	if start != 0 || end != 15 {
		t.Errorf(
			"expected Span (0, 15), got (%d, %d)",
			start,
			end,
		)
	}

	if string(
		node.Source(),
	) != "# Test Document" {
		t.Errorf(
			"expected Source '# Test Document', got '%s'",
			string(node.Source()),
		)
	}
}

func TestNodeInterface_Section(t *testing.T) {
	source := []byte("## My Section")
	node := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(13).
		WithSource(source).
		WithLevel(2).
		WithTitle([]byte("My Section")).
		WithDeltaType(deltaCAdded).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeSection {
		t.Errorf(
			"expected NodeType Section, got %v",
			node.NodeType(),
		)
	}

	section := node.(*NodeSection)
	if section.Level() != 2 {
		t.Errorf(
			"expected Level 2, got %d",
			section.Level(),
		)
	}
	if string(section.Title()) != "My Section" {
		t.Errorf(
			"expected Title 'My Section', got '%s'",
			string(section.Title()),
		)
	}
	if section.DeltaType() != deltaCAdded {
		t.Errorf(
			"expected DeltaType 'ADDED', got '%s'",
			section.DeltaType(),
		)
	}
}

func TestNodeInterface_Requirement(t *testing.T) {
	source := []byte(
		"### Requirement: ValidateInput",
	)
	node := NewNodeBuilder(NodeTypeRequirement).
		WithStart(0).
		WithEnd(31).
		WithSource(source).
		WithName("ValidateInput").
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeRequirement {
		t.Errorf(
			"expected NodeType Requirement, got %v",
			node.NodeType(),
		)
	}

	req := node.(*NodeRequirement)
	if req.Name() != "ValidateInput" {
		t.Errorf(
			"expected Name 'ValidateInput', got '%s'",
			req.Name(),
		)
	}
}

func TestNodeInterface_Scenario(t *testing.T) {
	source := []byte("#### Scenario: User Login")
	node := NewNodeBuilder(NodeTypeScenario).
		WithStart(0).
		WithEnd(25).
		WithSource(source).
		WithName("User Login").
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeScenario {
		t.Errorf(
			"expected NodeType Scenario, got %v",
			node.NodeType(),
		)
	}

	scenario := node.(*NodeScenario)
	if scenario.Name() != "User Login" {
		t.Errorf(
			"expected Name 'User Login', got '%s'",
			scenario.Name(),
		)
	}
}

func TestNodeInterface_Paragraph(t *testing.T) {
	source := []byte("This is a paragraph.")
	node := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(20).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeParagraph {
		t.Errorf(
			"expected NodeType Paragraph, got %v",
			node.NodeType(),
		)
	}
}

func TestNodeInterface_List(t *testing.T) {
	source := []byte("1. Item one\n2. Item two")
	node := NewNodeBuilder(NodeTypeList).
		WithStart(0).
		WithEnd(23).
		WithSource(source).
		WithOrdered(true).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeList {
		t.Errorf(
			"expected NodeType List, got %v",
			node.NodeType(),
		)
	}

	list := node.(*NodeList)
	if !list.Ordered() {
		t.Error("expected Ordered to be true")
	}
}

func TestNodeInterface_ListItem(t *testing.T) {
	source := []byte("- [ ] Task item")
	node := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(15).
		WithSource(source).
		WithChecked(boolPtr(false)).
		WithKeyword(deltaWhen).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeListItem {
		t.Errorf(
			"expected NodeType ListItem, got %v",
			node.NodeType(),
		)
	}

	item := node.(*NodeListItem)
	isChecked, hasCheckbox := item.Checked()
	if !hasCheckbox {
		t.Error("expected hasCheckbox to be true")
	}
	if isChecked {
		t.Error("expected isChecked to be false")
	}
	if item.Keyword() != deltaWhen {
		t.Errorf(
			"expected Keyword 'WHEN', got '%s'",
			item.Keyword(),
		)
	}
}

func TestNodeInterface_ListItem_NoCheckbox(
	t *testing.T,
) {
	source := []byte("- Regular item")
	node := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(14).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	item := node.(*NodeListItem)
	isChecked, hasCheckbox := item.Checked()
	if hasCheckbox {
		t.Error(
			"expected hasCheckbox to be false",
		)
	}
	if isChecked {
		t.Error("expected isChecked to be false")
	}
}

func TestNodeInterface_ListItem_CheckedCheckbox(
	t *testing.T,
) {
	source := []byte("- [x] Done item")
	node := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(15).
		WithSource(source).
		WithChecked(boolPtr(true)).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	item := node.(*NodeListItem)
	isChecked, hasCheckbox := item.Checked()
	if !hasCheckbox {
		t.Error("expected hasCheckbox to be true")
	}
	if !isChecked {
		t.Error("expected isChecked to be true")
	}
}

func TestNodeInterface_CodeBlock(t *testing.T) {
	source := []byte(
		"```go\nfmt.Println(\"Hello\")\n```",
	)
	node := NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(0).
		WithEnd(31).
		WithSource(source).
		WithLanguage([]byte("go")).
		WithContent([]byte("fmt.Println(\"Hello\")")).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeCodeBlock {
		t.Errorf(
			"expected NodeType CodeBlock, got %v",
			node.NodeType(),
		)
	}

	codeBlock := node.(*NodeCodeBlock)
	if string(codeBlock.Language()) != "go" {
		t.Errorf(
			"expected Language 'go', got '%s'",
			string(codeBlock.Language()),
		)
	}
	if string(
		codeBlock.Content(),
	) != "fmt.Println(\"Hello\")" {
		t.Errorf(
			"expected Content 'fmt.Println(\"Hello\")', got '%s'",
			string(codeBlock.Content()),
		)
	}
}

func TestNodeInterface_Blockquote(t *testing.T) {
	source := []byte("> Quoted text")
	node := NewNodeBuilder(NodeTypeBlockquote).
		WithStart(0).
		WithEnd(13).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeBlockquote {
		t.Errorf(
			"expected NodeType Blockquote, got %v",
			node.NodeType(),
		)
	}
}

func TestNodeInterface_Text(t *testing.T) {
	source := []byte("Plain text content")
	node := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(18).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeText {
		t.Errorf(
			"expected NodeType Text, got %v",
			node.NodeType(),
		)
	}

	text := node.(*NodeText)
	if text.Text() != "Plain text content" {
		t.Errorf(
			"expected Text 'Plain text content', got '%s'",
			text.Text(),
		)
	}
}

func TestNodeInterface_Strong(t *testing.T) {
	source := []byte("**bold**")
	node := NewNodeBuilder(NodeTypeStrong).
		WithStart(0).
		WithEnd(8).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeStrong {
		t.Errorf(
			"expected NodeType Strong, got %v",
			node.NodeType(),
		)
	}
}

func TestNodeInterface_Emphasis(t *testing.T) {
	source := []byte("*italic*")
	node := NewNodeBuilder(NodeTypeEmphasis).
		WithStart(0).
		WithEnd(8).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeEmphasis {
		t.Errorf(
			"expected NodeType Emphasis, got %v",
			node.NodeType(),
		)
	}
}

func TestNodeInterface_Strikethrough(
	t *testing.T,
) {
	source := []byte("~~struck~~")
	node := NewNodeBuilder(NodeTypeStrikethrough).
		WithStart(0).
		WithEnd(10).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeStrikethrough {
		t.Errorf(
			"expected NodeType Strikethrough, got %v",
			node.NodeType(),
		)
	}
}

func TestNodeInterface_Code(t *testing.T) {
	source := []byte("`inline code`")
	node := NewNodeBuilder(NodeTypeCode).
		WithStart(0).
		WithEnd(13).
		WithSource(source).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeCode {
		t.Errorf(
			"expected NodeType Code, got %v",
			node.NodeType(),
		)
	}

	code := node.(*NodeCode)
	if code.Code() != "`inline code`" {
		t.Errorf(
			"expected Code '`inline code`', got '%s'",
			code.Code(),
		)
	}
}

func TestNodeInterface_Link(t *testing.T) {
	source := []byte(
		"[text](https://example.com \"Title\")",
	)
	node := NewNodeBuilder(NodeTypeLink).
		WithStart(0).
		WithEnd(35).
		WithSource(source).
		WithURL([]byte(testExampleURL)).
		WithLinkTitle([]byte("Title")).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeLink {
		t.Errorf(
			"expected NodeType Link, got %v",
			node.NodeType(),
		)
	}

	link := node.(*NodeLink)
	if string(
		link.URL(),
	) != testExampleURL {
		t.Errorf(
			"expected URL 'https://example.com', got '%s'",
			string(link.URL()),
		)
	}
	if string(link.Title()) != "Title" {
		t.Errorf(
			"expected Title 'Title', got '%s'",
			string(link.Title()),
		)
	}
}

func TestNodeInterface_LinkDef(t *testing.T) {
	source := []byte(
		"[ref]: https://example.com \"Title\"",
	)
	node := NewNodeBuilder(NodeTypeLinkDef).
		WithStart(0).
		WithEnd(34).
		WithSource(source).
		WithURL([]byte(testExampleURL)).
		WithLinkTitle([]byte("Title")).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeLinkDef {
		t.Errorf(
			"expected NodeType LinkDef, got %v",
			node.NodeType(),
		)
	}

	linkDef := node.(*NodeLinkDef)
	if string(
		linkDef.URL(),
	) != testExampleURL {
		t.Errorf(
			"expected URL 'https://example.com', got '%s'",
			string(linkDef.URL()),
		)
	}
	if string(linkDef.Title()) != "Title" {
		t.Errorf(
			"expected Title 'Title', got '%s'",
			string(linkDef.Title()),
		)
	}
}

func TestNodeInterface_Wikilink(t *testing.T) {
	source := []byte("[[target|display#anchor]]")
	node := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(25).
		WithSource(source).
		WithTarget([]byte("target")).
		WithDisplay([]byte("display")).
		WithAnchor([]byte("anchor")).
		Build()

	if node == nil {
		t.Fatal(
			"expected node to be built, got nil",
		)
	}

	if node.NodeType() != NodeTypeWikilink {
		t.Errorf(
			"expected NodeType Wikilink, got %v",
			node.NodeType(),
		)
	}

	wikilink := node.(*NodeWikilink)
	if string(wikilink.Target()) != "target" {
		t.Errorf(
			"expected Target 'target', got '%s'",
			string(wikilink.Target()),
		)
	}
	if string(wikilink.Display()) != "display" {
		t.Errorf(
			"expected Display 'display', got '%s'",
			string(wikilink.Display()),
		)
	}
	if string(wikilink.Anchor()) != "anchor" {
		t.Errorf(
			"expected Anchor 'anchor', got '%s'",
			string(wikilink.Anchor()),
		)
	}
}

func TestHash_SameContentSameHash(t *testing.T) {
	source := []byte("Same content")

	node1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(12).
		WithSource(source).
		Build()

	node2 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(12).
		WithSource(source).
		Build()

	if node1.Hash() != node2.Hash() {
		t.Errorf(
			"expected same hash for identical nodes, got %d and %d",
			node1.Hash(),
			node2.Hash(),
		)
	}
}

func TestHash_DifferentContentDifferentHash(
	t *testing.T,
) {
	node1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(12).
		WithSource([]byte("Content one")).
		Build()

	node2 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(12).
		WithSource([]byte("Content two")).
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash for different content",
		)
	}
}

func TestHash_DifferentNodeTypeDifferentHash(
	t *testing.T,
) {
	source := []byte("Same content")

	node1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(12).
		WithSource(source).
		Build()

	node2 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(12).
		WithSource(source).
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash for different node types",
		)
	}
}

func TestHash_IncludesChildren(t *testing.T) {
	child1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(5).
		WithSource([]byte("Child")).
		Build()

	child2 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(8).
		WithSource([]byte("DiffChild")).
		Build()

	parent1 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("Parent")).
		WithChildren([]Node{child1}).
		Build()

	parent2 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("Parent")).
		WithChildren([]Node{child2}).
		Build()

	if parent1.Hash() == parent2.Hash() {
		t.Error(
			"expected different hash when children differ",
		)
	}
}

func TestHash_IncludesTypeSpecificFields_Section(
	t *testing.T,
) {
	source := []byte("## Section")

	node1 := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(10).
		WithSource(source).
		WithLevel(2).
		WithTitle([]byte("Section")).
		WithDeltaType(deltaCAdded).
		Build()

	node2 := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(10).
		WithSource(source).
		WithLevel(2).
		WithTitle([]byte("Section")).
		WithDeltaType("REMOVED").
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash when DeltaType differs",
		)
	}
}

func TestHash_IncludesTypeSpecificFields_List(
	t *testing.T,
) {
	source := []byte("1. Item")

	node1 := NewNodeBuilder(NodeTypeList).
		WithStart(0).
		WithEnd(7).
		WithSource(source).
		WithOrdered(true).
		Build()

	node2 := NewNodeBuilder(NodeTypeList).
		WithStart(0).
		WithEnd(7).
		WithSource(source).
		WithOrdered(false).
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash when Ordered differs",
		)
	}
}

func TestHash_IncludesTypeSpecificFields_ListItem(
	t *testing.T,
) {
	source := []byte("- [ ] Item")

	node1 := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(10).
		WithSource(source).
		WithChecked(boolPtr(false)).
		Build()

	node2 := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(10).
		WithSource(source).
		WithChecked(boolPtr(true)).
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash when Checked state differs",
		)
	}
}

func TestHash_IncludesTypeSpecificFields_CodeBlock(
	t *testing.T,
) {
	source := []byte("```go\ncode\n```")

	node1 := NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(0).
		WithEnd(14).
		WithSource(source).
		WithLanguage([]byte("go")).
		WithContent([]byte("code")).
		Build()

	node2 := NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(0).
		WithEnd(14).
		WithSource(source).
		WithLanguage([]byte("python")).
		WithContent([]byte("code")).
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash when Language differs",
		)
	}
}

func TestHash_IncludesTypeSpecificFields_Link(
	t *testing.T,
) {
	source := []byte("[text](url)")

	node1 := NewNodeBuilder(NodeTypeLink).
		WithStart(0).
		WithEnd(11).
		WithSource(source).
		WithURL([]byte("url1")).
		Build()

	node2 := NewNodeBuilder(NodeTypeLink).
		WithStart(0).
		WithEnd(11).
		WithSource(source).
		WithURL([]byte("url2")).
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash when URL differs",
		)
	}
}

func TestHash_IncludesTypeSpecificFields_Wikilink(
	t *testing.T,
) {
	source := []byte("[[target]]")

	node1 := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(10).
		WithSource(source).
		WithTarget([]byte("target1")).
		Build()

	node2 := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(10).
		WithSource(source).
		WithTarget([]byte("target2")).
		Build()

	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different hash when Target differs",
		)
	}
}

func TestImmutability_ChildrenDefensiveCopy(
	t *testing.T,
) {
	child := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(4).
		WithSource([]byte("text")).
		Build()

	parent := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren([]Node{child}).
		Build()

	// Get children and modify the returned slice
	children := parent.Children()
	if len(children) != 1 {
		t.Fatalf(
			"expected 1 child, got %d",
			len(children),
		)
	}

	// Set the returned slice element to nil
	children[0] = nil

	// Original children should be unaffected
	originalChildren := parent.Children()
	if len(originalChildren) != 1 {
		t.Fatalf(
			"expected 1 original child, got %d",
			len(originalChildren),
		)
	}
	if originalChildren[0] == nil {
		t.Error(
			"modifying returned Children() slice should not affect node's internal children",
		)
	}
}

func TestImmutability_SourceNotModifiable(
	t *testing.T,
) {
	source := []byte("original source")
	node := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(15).
		WithSource(source).
		Build()

	// Modify the original source slice
	source[0] = 'X'

	// Note: Since source is a reference, this WILL affect the node's source.
	// This is documented behavior - Source() returns a zero-copy view.
	// The test demonstrates this expected behavior.
	nodeSource := node.Source()
	if nodeSource[0] != 'X' {
		t.Error(
			"expected source to share underlying array with original (zero-copy)",
		)
	}
}

func TestImmutability_NilChildren(t *testing.T) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(4).
		WithSource([]byte("text")).
		Build()

	children := node.Children()
	if children != nil {
		t.Errorf(
			"expected nil children for leaf node, got %v",
			children,
		)
	}
}

func TestNodeBuilder_BuildAllNodeTypes(
	t *testing.T,
) {
	testCases := []struct {
		name     string
		nodeType NodeType
		builder  func() *NodeBuilder
	}{
		{
			"Document",
			NodeTypeDocument,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeDocument,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("doc"))
			},
		},
		{
			"Section",
			NodeTypeSection,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeSection,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("## sec")).
					WithLevel(2).
					WithTitle([]byte("sec"))
			},
		},
		{
			"Requirement",
			NodeTypeRequirement,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeRequirement,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("### Req")).
					WithName("Req")
			},
		},
		{
			"Scenario",
			NodeTypeScenario,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeScenario,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("#### Scn")).
					WithName("Scn")
			},
		},
		{
			"Paragraph",
			NodeTypeParagraph,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeParagraph,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("para"))
			},
		},
		{
			"List",
			NodeTypeList,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeList,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("- item")).
					WithOrdered(false)
			},
		},
		{
			"ListItem",
			NodeTypeListItem,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeListItem,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("- item")).
					WithKeyword("THEN")
			},
		},
		{
			"CodeBlock",
			NodeTypeCodeBlock,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeCodeBlock,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("```\ncode\n```")).
					WithLanguage([]byte("go")).
					WithContent([]byte("code"))
			},
		},
		{
			"Blockquote",
			NodeTypeBlockquote,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeBlockquote,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("> quote"))
			},
		},
		{
			"Text",
			NodeTypeText,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeText,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("text"))
			},
		},
		{
			"Strong",
			NodeTypeStrong,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeStrong,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("**bold**"))
			},
		},
		{
			"Emphasis",
			NodeTypeEmphasis,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeEmphasis,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("*italic*"))
			},
		},
		{
			"Strikethrough",
			NodeTypeStrikethrough,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeStrikethrough,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("~~strike~~"))
			},
		},
		{
			"Code",
			NodeTypeCode,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeCode,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("`code`"))
			},
		},
		{
			"Link",
			NodeTypeLink,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeLink,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("[t](url)")).
					WithURL([]byte("url"))
			},
		},
		{
			"LinkDef",
			NodeTypeLinkDef,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeLinkDef,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("[ref]: url")).
					WithURL([]byte("url"))
			},
		},
		{
			"Wikilink",
			NodeTypeWikilink,
			func() *NodeBuilder {
				return NewNodeBuilder(
					NodeTypeWikilink,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("[[target]]")).
					WithTarget([]byte("target"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.builder().Build()
			if node == nil {
				t.Fatal(
					"expected node to be built, got nil",
				)
			}
			if node.NodeType() != tc.nodeType {
				t.Errorf(
					"expected NodeType %v, got %v",
					tc.nodeType,
					node.NodeType(),
				)
			}
		})
	}
}

func TestNodeBuilder_ValidationStartLessThanEnd(
	t *testing.T,
) {
	builder := NewNodeBuilder(NodeTypeText).
		WithStart(10).
		WithEnd(5) // End < Start

	err := builder.Validate()
	if err == nil {
		t.Fatal(
			"expected validation error for Start > End",
		)
	}

	if _, ok := err.(*BuilderValidationError); !ok {
		t.Errorf(
			"expected BuilderValidationError, got %T",
			err,
		)
	}

	node := builder.Build()
	if node != nil {
		t.Error(
			"expected Build() to return nil when validation fails",
		)
	}
}

func TestNodeBuilder_ValidationStartEqualsEnd(
	t *testing.T,
) {
	builder := NewNodeBuilder(NodeTypeText).
		WithStart(5).
		WithEnd(5)
		// Start == End is valid (empty span)

	err := builder.Validate()
	if err != nil {
		t.Errorf(
			"expected no validation error for Start == End, got %v",
			err,
		)
	}

	node := builder.Build()
	if node == nil {
		t.Error(
			"expected Build() to succeed for Start == End",
		)
	}
}

func TestNodeBuilder_ValidationSectionLevel(
	t *testing.T,
) {
	// Valid levels
	for level := 1; level <= 6; level++ {
		builder := NewNodeBuilder(
			NodeTypeSection,
		).
			WithStart(0).
			WithEnd(10).
			WithSource([]byte("## sec")).
			WithLevel(level)

		err := builder.Validate()
		if err != nil {
			t.Errorf(
				"expected no error for level %d, got %v",
				level,
				err,
			)
		}
	}

	// Invalid levels
	invalidLevels := []int{0, 7, -1, 100}
	for _, level := range invalidLevels {
		builder := NewNodeBuilder(
			NodeTypeSection,
		).
			WithStart(0).
			WithEnd(10).
			WithSource([]byte("## sec")).
			WithLevel(level)

		err := builder.Validate()
		if err == nil {
			t.Errorf(
				"expected validation error for invalid level %d",
				level,
			)
		}
	}
}

func TestNodeBuilder_ValidationChildrenWithinParentSpan(
	t *testing.T,
) {
	child := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(5).
		WithSource([]byte("child")).
		Build()

	// Child within parent span - should succeed
	validParent := NewNodeBuilder(
		NodeTypeParagraph,
	).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren([]Node{child})

	err := validParent.Validate()
	if err != nil {
		t.Errorf(
			"expected no error for child within parent span, got %v",
			err,
		)
	}

	// Child outside parent span - should fail
	invalidParent := NewNodeBuilder(
		NodeTypeParagraph,
	).
		WithStart(6).
		// Parent starts after child
		WithEnd(10).
		WithSource([]byte("para")).
		WithChildren([]Node{child})

	err = invalidParent.Validate()
	if err == nil {
		t.Error(
			"expected validation error for child outside parent span",
		)
	}
}

func TestNodeBuilder_ValidationNilChild(
	t *testing.T,
) {
	builder := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph"))

	// Manually set children with nil (bypassing WithChildren which makes a copy)
	builder.children = []Node{nil}

	err := builder.Validate()
	if err == nil {
		t.Error(
			"expected validation error for nil child",
		)
	}
}

func TestNodeBuilder_WithChildren(t *testing.T) {
	child1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(5).
		WithSource([]byte("child")).
		Build()

	child2 := NewNodeBuilder(NodeTypeText).
		WithStart(5).
		WithEnd(10).
		WithSource([]byte("child")).
		Build()

	parent := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren([]Node{child1, child2}).
		Build()

	if parent == nil {
		t.Fatal(
			"expected parent to be built, got nil",
		)
	}

	children := parent.Children()
	if len(children) != 2 {
		t.Errorf(
			"expected 2 children, got %d",
			len(children),
		)
	}
}

func TestNodeBuilder_WithChildrenNil(
	t *testing.T,
) {
	parent := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren(nil).
		Build()

	if parent == nil {
		t.Fatal(
			"expected parent to be built, got nil",
		)
	}

	children := parent.Children()
	if children != nil {
		t.Errorf(
			"expected nil children, got %v",
			children,
		)
	}
}

func TestNodeBuilder_ToBuilder(t *testing.T) {
	original := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(15).
		WithSource([]byte("## My Section")).
		WithLevel(2).
		WithTitle([]byte("My Section")).
		WithDeltaType(deltaCAdded).
		Build()

	section := original.(*NodeSection)
	builder := section.ToBuilder()

	if builder == nil {
		t.Fatal(
			"expected builder from ToBuilder, got nil",
		)
	}

	// Modify the builder
	builder.WithDeltaType(deltaModified)

	// Build new node
	modified := builder.Build()
	if modified == nil {
		t.Fatal("expected modified node, got nil")
	}

	modifiedSection := modified.(*NodeSection)
	if modifiedSection.DeltaType() != deltaModified {
		t.Errorf(
			"expected DeltaType 'MODIFIED', got '%s'",
			modifiedSection.DeltaType(),
		)
	}

	// Original should be unchanged
	if section.DeltaType() != deltaCAdded {
		t.Errorf(
			"expected original DeltaType 'ADDED', got '%s'",
			section.DeltaType(),
		)
	}
}

func TestNodeBuilder_ToBuilder_AllTypes(
	t *testing.T,
) {
	// Test ToBuilder works for all node types with type-specific fields
	testCases := []struct {
		name    string
		builder func() Node
		verify  func(t *testing.T, original, rebuilt Node)
	}{
		{
			"Document",
			func() Node {
				return NewNodeBuilder(
					NodeTypeDocument,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("doc")).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				if original.Hash() != rebuilt.Hash() {
					t.Error("hashes should match")
				}
			},
		},
		{
			"Section",
			func() Node {
				return NewNodeBuilder(
					NodeTypeSection,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("## sec")).
					WithLevel(2).
					WithTitle([]byte("sec")).
					WithDeltaType(deltaCAdded).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeSection)
				r := rebuilt.(*NodeSection)
				if o.Level() != r.Level() ||
					!bytes.Equal(o.Title(), r.Title()) ||
					o.DeltaType() != r.DeltaType() {
					t.Error(
						"Section fields should match",
					)
				}
			},
		},
		{
			"Requirement",
			func() Node {
				return NewNodeBuilder(
					NodeTypeRequirement,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("### Req")).
					WithName("MyReq").
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeRequirement)
				r := rebuilt.(*NodeRequirement)
				if o.Name() != r.Name() {
					t.Error(
						"Requirement name should match",
					)
				}
			},
		},
		{
			"Scenario",
			func() Node {
				return NewNodeBuilder(
					NodeTypeScenario,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("#### Scn")).
					WithName("MyScenario").
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeScenario)
				r := rebuilt.(*NodeScenario)
				if o.Name() != r.Name() {
					t.Error(
						"Scenario name should match",
					)
				}
			},
		},
		{
			"List",
			func() Node {
				return NewNodeBuilder(
					NodeTypeList,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("1. item")).
					WithOrdered(true).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeList)
				r := rebuilt.(*NodeList)
				if o.Ordered() != r.Ordered() {
					t.Error(
						"List ordered should match",
					)
				}
			},
		},
		{
			"ListItem",
			func() Node {
				return NewNodeBuilder(
					NodeTypeListItem,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("- [ ] item")).
					WithChecked(boolPtr(false)).
					WithKeyword(deltaWhen).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeListItem)
				r := rebuilt.(*NodeListItem)
				oc, oh := o.Checked()
				rc, rh := r.Checked()
				if oc != rc || oh != rh ||
					o.Keyword() != r.Keyword() {
					t.Error(
						"ListItem fields should match",
					)
				}
			},
		},
		{
			"CodeBlock",
			func() Node {
				return NewNodeBuilder(
					NodeTypeCodeBlock,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("```\ncode\n```")).
					WithLanguage([]byte("go")).
					WithContent([]byte("code")).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeCodeBlock)
				r := rebuilt.(*NodeCodeBlock)
				if !bytes.Equal(o.Language(), r.Language()) ||
					!bytes.Equal(o.Content(), r.Content()) {
					t.Error(
						"CodeBlock fields should match",
					)
				}
			},
		},
		{
			"Link",
			func() Node {
				return NewNodeBuilder(
					NodeTypeLink,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("[t](url)")).
					WithURL([]byte(testExampleURL)).
					WithLinkTitle([]byte("Title")).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeLink)
				r := rebuilt.(*NodeLink)
				if !bytes.Equal(o.URL(), r.URL()) ||
					!bytes.Equal(o.Title(), r.Title()) {
					t.Error(
						"Link fields should match",
					)
				}
			},
		},
		{
			"LinkDef",
			func() Node {
				return NewNodeBuilder(
					NodeTypeLinkDef,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("[ref]: url")).
					WithURL([]byte(testExampleURL)).
					WithLinkTitle([]byte("Title")).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeLinkDef)
				r := rebuilt.(*NodeLinkDef)
				if !bytes.Equal(o.URL(), r.URL()) ||
					!bytes.Equal(o.Title(), r.Title()) {
					t.Error(
						"LinkDef fields should match",
					)
				}
			},
		},
		{
			"Wikilink",
			func() Node {
				return NewNodeBuilder(
					NodeTypeWikilink,
				).
					WithStart(0).
					WithEnd(10).
					WithSource([]byte("[[t|d#a]]")).
					WithTarget([]byte("target")).
					WithDisplay([]byte("display")).
					WithAnchor([]byte("anchor")).
					Build()
			},
			func(t *testing.T, original, rebuilt Node) {
				o := original.(*NodeWikilink)
				r := rebuilt.(*NodeWikilink)
				if !bytes.Equal(o.Target(), r.Target()) ||
					!bytes.Equal(o.Display(), r.Display()) ||
					!bytes.Equal(o.Anchor(), r.Anchor()) {
					t.Error(
						"Wikilink fields should match",
					)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := tc.builder()
			if original == nil {
				t.Fatal("original node is nil")
			}

			// Get builder and rebuild
			var builder *NodeBuilder
			switch n := original.(type) {
			case *NodeDocument:
				builder = n.ToBuilder()
			case *NodeSection:
				builder = n.ToBuilder()
			case *NodeRequirement:
				builder = n.ToBuilder()
			case *NodeScenario:
				builder = n.ToBuilder()
			case *NodeList:
				builder = n.ToBuilder()
			case *NodeListItem:
				builder = n.ToBuilder()
			case *NodeCodeBlock:
				builder = n.ToBuilder()
			case *NodeLink:
				builder = n.ToBuilder()
			case *NodeLinkDef:
				builder = n.ToBuilder()
			case *NodeWikilink:
				builder = n.ToBuilder()
			default:
				t.Fatalf("unhandled type %T", original)
			}

			if builder == nil {
				t.Fatal("builder is nil")
			}

			rebuilt := builder.Build()
			if rebuilt == nil {
				t.Fatal("rebuilt node is nil")
			}

			tc.verify(t, original, rebuilt)
		})
	}
}

func TestNodeBuilder_ToBuilder_Nil(t *testing.T) {
	builder := nodeToBuilder(nil)
	if builder != nil {
		t.Error(
			"expected nil builder for nil node",
		)
	}
}

func TestEqual_SameNodes(t *testing.T) {
	source := []byte("test text")
	node1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(9).
		WithSource(source).
		Build()

	node2 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(9).
		WithSource(source).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected equal nodes to return true",
		)
	}
	if !node2.Equal(node1) {
		t.Error("expected Equal to be symmetric")
	}
}

func TestEqual_DifferentSource(t *testing.T) {
	node1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(9).
		WithSource([]byte("text one")).
		Build()

	node2 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(9).
		WithSource([]byte("text two")).
		Build()

	if node1.Equal(node2) {
		t.Error(
			"expected different source nodes to return false",
		)
	}
}

func TestEqual_DifferentType(t *testing.T) {
	source := []byte("content")
	node1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(7).
		WithSource(source).
		Build()

	node2 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(7).
		WithSource(source).
		Build()

	if node1.Equal(node2) {
		t.Error(
			"expected different type nodes to return false",
		)
	}
}

func TestEqual_WithChildren(t *testing.T) {
	child := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(4).
		WithSource([]byte("text")).
		Build()

	node1 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren([]Node{child}).
		Build()

	node2 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren([]Node{child}).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected nodes with same children to be equal",
		)
	}
}

func TestEqual_DifferentChildren(t *testing.T) {
	child1 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(5).
		WithSource([]byte("child")).
		Build()

	child2 := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(8).
		WithSource([]byte("diffchld")).
		Build()

	node1 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren([]Node{child1}).
		Build()

	node2 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph")).
		WithChildren([]Node{child2}).
		Build()

	if node1.Equal(node2) {
		t.Error(
			"expected nodes with different children to not be equal",
		)
	}
}

func TestEqual_NilComparison(t *testing.T) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(4).
		WithSource([]byte("text")).
		Build()

	if node.Equal(nil) {
		t.Error(
			"expected Equal(nil) to return false",
		)
	}
}

func TestEqual_TypeSpecific_Section(
	t *testing.T,
) {
	source := []byte("## Section")

	node1 := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithLevel(2).
		WithTitle([]byte("Section")).
		WithDeltaType(deltaCAdded).
		Build()

	node2 := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithLevel(2).
		WithTitle([]byte("Section")).
		WithDeltaType(deltaCAdded).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical Section nodes to be equal",
		)
	}

	// Different level
	node3 := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithLevel(3).
		WithTitle([]byte("Section")).
		WithDeltaType(deltaCAdded).
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected Section nodes with different levels to not be equal",
		)
	}

	// Different delta type
	node4 := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithLevel(2).
		WithTitle([]byte("Section")).
		WithDeltaType("REMOVED").
		Build()

	if node1.Equal(node4) {
		t.Error(
			"expected Section nodes with different delta types to not be equal",
		)
	}

	// Different title
	node5 := NewNodeBuilder(NodeTypeSection).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithLevel(2).
		WithTitle([]byte("Different")).
		WithDeltaType(deltaCAdded).
		Build()

	if node1.Equal(node5) {
		t.Error(
			"expected Section nodes with different titles to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_Requirement(
	t *testing.T,
) {
	source := []byte("### Requirement: Test")

	node1 := NewNodeBuilder(NodeTypeRequirement).
		WithStart(0).
		WithEnd(21).WithSource(source).
		WithName("Test").
		Build()

	node2 := NewNodeBuilder(NodeTypeRequirement).
		WithStart(0).
		WithEnd(21).WithSource(source).
		WithName("Test").
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical Requirement nodes to be equal",
		)
	}

	node3 := NewNodeBuilder(NodeTypeRequirement).
		WithStart(0).
		WithEnd(21).WithSource(source).
		WithName("Different").
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected Requirement nodes with different names to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_Scenario(
	t *testing.T,
) {
	source := []byte("#### Scenario: Test")

	node1 := NewNodeBuilder(NodeTypeScenario).
		WithStart(0).
		WithEnd(19).WithSource(source).
		WithName("Test").
		Build()

	node2 := NewNodeBuilder(NodeTypeScenario).
		WithStart(0).
		WithEnd(19).WithSource(source).
		WithName("Test").
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical Scenario nodes to be equal",
		)
	}

	node3 := NewNodeBuilder(NodeTypeScenario).
		WithStart(0).
		WithEnd(19).WithSource(source).
		WithName("Different").
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected Scenario nodes with different names to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_List(t *testing.T) {
	source := []byte("1. item")

	node1 := NewNodeBuilder(NodeTypeList).
		WithStart(0).
		WithEnd(7).WithSource(source).
		WithOrdered(true).
		Build()

	node2 := NewNodeBuilder(NodeTypeList).
		WithStart(0).
		WithEnd(7).WithSource(source).
		WithOrdered(true).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical List nodes to be equal",
		)
	}

	node3 := NewNodeBuilder(NodeTypeList).
		WithStart(0).
		WithEnd(7).WithSource(source).
		WithOrdered(false).
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected List nodes with different ordered flags to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_ListItem(
	t *testing.T,
) {
	source := []byte("- [ ] item")

	node1 := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithChecked(boolPtr(false)).
		WithKeyword(deltaWhen).
		Build()

	node2 := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithChecked(boolPtr(false)).
		WithKeyword(deltaWhen).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical ListItem nodes to be equal",
		)
	}

	// Different checked state
	node3 := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithChecked(boolPtr(true)).
		WithKeyword(deltaWhen).
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected ListItem nodes with different checked states to not be equal",
		)
	}

	// No checkbox vs checkbox
	node4 := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithKeyword(deltaWhen).
		Build()

	if node1.Equal(node4) {
		t.Error(
			"expected ListItem nodes with/without checkbox to not be equal",
		)
	}

	// Different keyword
	node5 := NewNodeBuilder(NodeTypeListItem).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithChecked(boolPtr(false)).
		WithKeyword("THEN").
		Build()

	if node1.Equal(node5) {
		t.Error(
			"expected ListItem nodes with different keywords to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_CodeBlock(
	t *testing.T,
) {
	source := []byte("```go\ncode\n```")

	node1 := NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(0).
		WithEnd(14).WithSource(source).
		WithLanguage([]byte("go")).
		WithContent([]byte("code")).
		Build()

	node2 := NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(0).
		WithEnd(14).WithSource(source).
		WithLanguage([]byte("go")).
		WithContent([]byte("code")).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical CodeBlock nodes to be equal",
		)
	}

	// Different language
	node3 := NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(0).
		WithEnd(14).WithSource(source).
		WithLanguage([]byte("python")).
		WithContent([]byte("code")).
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected CodeBlock nodes with different languages to not be equal",
		)
	}

	// Different content
	node4 := NewNodeBuilder(NodeTypeCodeBlock).
		WithStart(0).
		WithEnd(14).WithSource(source).
		WithLanguage([]byte("go")).
		WithContent([]byte("different")).
		Build()

	if node1.Equal(node4) {
		t.Error(
			"expected CodeBlock nodes with different content to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_Link(t *testing.T) {
	source := []byte("[text](url)")

	node1 := NewNodeBuilder(NodeTypeLink).
		WithStart(0).
		WithEnd(11).WithSource(source).
		WithURL([]byte(testExampleURL)).
		WithLinkTitle([]byte("Title")).
		Build()

	node2 := NewNodeBuilder(NodeTypeLink).
		WithStart(0).
		WithEnd(11).WithSource(source).
		WithURL([]byte(testExampleURL)).
		WithLinkTitle([]byte("Title")).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical Link nodes to be equal",
		)
	}

	// Different URL
	node3 := NewNodeBuilder(NodeTypeLink).
		WithStart(0).
		WithEnd(11).WithSource(source).
		WithURL([]byte("https://different.com")).
		WithLinkTitle([]byte("Title")).
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected Link nodes with different URLs to not be equal",
		)
	}

	// Different title
	node4 := NewNodeBuilder(NodeTypeLink).
		WithStart(0).
		WithEnd(11).WithSource(source).
		WithURL([]byte(testExampleURL)).
		WithLinkTitle([]byte("Different")).
		Build()

	if node1.Equal(node4) {
		t.Error(
			"expected Link nodes with different titles to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_LinkDef(
	t *testing.T,
) {
	source := []byte("[ref]: url")

	node1 := NewNodeBuilder(NodeTypeLinkDef).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithURL([]byte(testExampleURL)).
		WithLinkTitle([]byte("Title")).
		Build()

	node2 := NewNodeBuilder(NodeTypeLinkDef).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithURL([]byte(testExampleURL)).
		WithLinkTitle([]byte("Title")).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical LinkDef nodes to be equal",
		)
	}

	node3 := NewNodeBuilder(NodeTypeLinkDef).
		WithStart(0).
		WithEnd(10).WithSource(source).
		WithURL([]byte("https://different.com")).
		WithLinkTitle([]byte("Title")).
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected LinkDef nodes with different URLs to not be equal",
		)
	}
}

func TestEqual_TypeSpecific_Wikilink(
	t *testing.T,
) {
	source := []byte("[[target|display#anchor]]")

	node1 := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(25).WithSource(source).
		WithTarget([]byte("target")).
		WithDisplay([]byte("display")).
		WithAnchor([]byte("anchor")).
		Build()

	node2 := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(25).WithSource(source).
		WithTarget([]byte("target")).
		WithDisplay([]byte("display")).
		WithAnchor([]byte("anchor")).
		Build()

	if !node1.Equal(node2) {
		t.Error(
			"expected identical Wikilink nodes to be equal",
		)
	}

	// Different target
	node3 := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(25).WithSource(source).
		WithTarget([]byte("different")).
		WithDisplay([]byte("display")).
		WithAnchor([]byte("anchor")).
		Build()

	if node1.Equal(node3) {
		t.Error(
			"expected Wikilink nodes with different targets to not be equal",
		)
	}

	// Different display
	node4 := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(25).WithSource(source).
		WithTarget([]byte("target")).
		WithDisplay([]byte("different")).
		WithAnchor([]byte("anchor")).
		Build()

	if node1.Equal(node4) {
		t.Error(
			"expected Wikilink nodes with different displays to not be equal",
		)
	}

	// Different anchor
	node5 := NewNodeBuilder(NodeTypeWikilink).
		WithStart(0).
		WithEnd(25).WithSource(source).
		WithTarget([]byte("target")).
		WithDisplay([]byte("display")).
		WithAnchor([]byte("different")).
		Build()

	if node1.Equal(node5) {
		t.Error(
			"expected Wikilink nodes with different anchors to not be equal",
		)
	}
}

func TestNodeType_String(t *testing.T) {
	testCases := []struct {
		nodeType NodeType
		expected string
	}{
		{NodeTypeDocument, "Document"},
		{NodeTypeSection, "Section"},
		{NodeTypeRequirement, "Requirement"},
		{NodeTypeScenario, "Scenario"},
		{NodeTypeParagraph, "Paragraph"},
		{NodeTypeList, "List"},
		{NodeTypeListItem, "ListItem"},
		{NodeTypeCodeBlock, "CodeBlock"},
		{NodeTypeBlockquote, "Blockquote"},
		{NodeTypeText, "Text"},
		{NodeTypeStrong, "Strong"},
		{NodeTypeEmphasis, "Emphasis"},
		{NodeTypeStrikethrough, "Strikethrough"},
		{NodeTypeCode, "Code"},
		{NodeTypeLink, "Link"},
		{NodeTypeLinkDef, "LinkDef"},
		{NodeTypeWikilink, "Wikilink"},
		{NodeType(255), "Unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if tc.nodeType.String() != tc.expected {
				t.Errorf(
					"expected %s, got %s",
					tc.expected,
					tc.nodeType.String(),
				)
			}
		})
	}
}

func TestBuilderValidationError_Error(
	t *testing.T,
) {
	// Error without index
	err1 := &BuilderValidationError{
		Field:   "Start/End",
		Message: "Start must be less than or equal to End",
	}
	expected1 := "Start/End: Start must be less than or equal to End"
	if err1.Error() != expected1 {
		t.Errorf(
			"expected '%s', got '%s'",
			expected1,
			err1.Error(),
		)
	}

	// Error with index
	err2 := &BuilderValidationError{
		Field:   "Children",
		Message: "Child span outside parent",
		Index:   2,
	}
	expected2 := "Children[2]: Child span outside parent"
	if err2.Error() != expected2 {
		t.Errorf(
			"expected '%s', got '%s'",
			expected2,
			err2.Error(),
		)
	}
}

func TestEdgeCase_EmptySource(t *testing.T) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(0).
		WithSource([]byte{}).
		Build()

	if node == nil {
		t.Fatal(
			"expected node with empty source to be built",
		)
	}

	if len(node.Source()) != 0 {
		t.Error("expected empty source")
	}
}

func TestEdgeCase_DeepNestedChildren(
	t *testing.T,
) {
	// Create a deeply nested structure
	deepChild := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(4).
		WithSource([]byte("text")).Build()

	for range 10 {
		deepChild = NewNodeBuilder(
			NodeTypeParagraph,
		).
			WithStart(0).
			WithEnd(10).
			WithSource([]byte("paragraph")).
			WithChildren([]Node{deepChild}).
			Build()
	}

	if deepChild == nil {
		t.Fatal(
			"expected deeply nested node to be built",
		)
	}

	// Verify we can traverse down
	current := deepChild
	for depth := range 10 {
		children := current.Children()
		if len(children) != 1 {
			t.Fatalf(
				"expected 1 child at depth %d, got %d",
				depth,
				len(children),
			)
		}
		current = children[0]
	}

	if current.NodeType() != NodeTypeText {
		t.Error(
			"expected deepest node to be Text",
		)
	}
}

func TestEdgeCase_ManyChildren(t *testing.T) {
	children := make([]Node, 100)
	for i := range 100 {
		children[i] = NewNodeBuilder(
			NodeTypeText,
		).
			WithStart(i).
			WithEnd(i + 1).
			WithSource([]byte("x")).
			Build()
	}

	parent := NewNodeBuilder(NodeTypeList).
		WithStart(0).
		WithEnd(100).WithSource([]byte("list")).
		WithChildren(children).
		Build()

	if parent == nil {
		t.Fatal(
			"expected node with many children to be built",
		)
	}

	resultChildren := parent.Children()
	if len(resultChildren) != 100 {
		t.Errorf(
			"expected 100 children, got %d",
			len(resultChildren),
		)
	}
}

func TestEdgeCase_BuilderReuse(t *testing.T) {
	// Create a builder and build multiple nodes with modifications
	builder := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(10).WithSource([]byte("initial"))

	node1 := builder.Build()
	if node1 == nil {
		t.Fatal("expected first node to be built")
	}

	// Modify builder and build again
	builder.WithSource([]byte("modified"))
	node2 := builder.Build()
	if node2 == nil {
		t.Fatal(
			"expected second node to be built",
		)
	}

	// Nodes should be different
	if node1.Hash() == node2.Hash() {
		t.Error(
			"expected different nodes to have different hashes",
		)
	}
}

func TestBytesEqual(t *testing.T) {
	testCases := []struct {
		name     string
		a        []byte
		b        []byte
		expected bool
	}{
		{"both nil", nil, nil, true},
		{"a nil", nil, []byte("x"), false},
		{"b nil", []byte("x"), nil, false},
		{
			"equal",
			[]byte("test"),
			[]byte("test"),
			true,
		},
		{
			"different length",
			[]byte("test"),
			[]byte("testing"),
			false,
		},
		{
			"different content",
			[]byte("test"),
			[]byte("best"),
			false,
		},
		{"empty equal", []byte{}, []byte{}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := bytesEqual(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf(
					"bytesEqual(%v, %v) = %v, expected %v",
					tc.a,
					tc.b,
					result,
					tc.expected,
				)
			}
		})
	}
}

func TestEqualNodes_BothNil(t *testing.T) {
	if !equalNodes(nil, nil) {
		t.Error(
			"expected equalNodes(nil, nil) to return true",
		)
	}
}

func TestEqualNodes_OneNil(t *testing.T) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(4).WithSource([]byte("text")).
		Build()

	if equalNodes(nil, node) {
		t.Error(
			"expected equalNodes(nil, node) to return false",
		)
	}

	if equalNodes(node, nil) {
		t.Error(
			"expected equalNodes(node, nil) to return false",
		)
	}
}
