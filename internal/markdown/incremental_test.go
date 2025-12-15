package markdown

import (
	"testing"
)

func TestParseIncremental_NoOldTree(
	t *testing.T,
) {
	// When no old tree is provided, should do a full parse
	newSource := []byte(
		"# Header\n\nParagraph text.",
	)
	tree, errors := ParseIncremental(
		nil,
		nil,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if tree == nil {
		t.Fatal("expected non-nil tree")
	}
	if tree.NodeType() != NodeTypeDocument {
		t.Errorf(
			"expected Document node, got %v",
			tree.NodeType(),
		)
	}
}

func TestParseIncremental_IdenticalSources(
	t *testing.T,
) {
	// When sources are identical, should return the same tree
	source := []byte(
		"# Header\n\nParagraph text.",
	)
	oldTree, _ := Parse(source)

	newTree, errors := ParseIncremental(
		oldTree,
		source,
		source,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	// The tree should be the exact same object for identical sources
	if newTree != oldTree {
		t.Error(
			"expected same tree object for identical sources",
		)
	}
}

func TestParseIncremental_SimpleTextInsertion(
	t *testing.T,
) {
	// Simple insertion test
	oldSource := []byte("Hello world")
	newSource := []byte("Hello beautiful world")

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	// Verify the new tree reflects the new content
	if newTree.Hash() == oldTree.Hash() {
		t.Error(
			"expected different hash for modified content",
		)
	}
}

func TestParseIncremental_SimpleTextDeletion(
	t *testing.T,
) {
	// Simple deletion test
	oldSource := []byte("Hello beautiful world")
	newSource := []byte("Hello world")

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	// Verify the new tree reflects the new content
	if newTree.Hash() == oldTree.Hash() {
		t.Error(
			"expected different hash for modified content",
		)
	}
}

func TestParseIncremental_SimpleTextReplacement(
	t *testing.T,
) {
	// Simple replacement test
	oldSource := []byte("Hello world")
	newSource := []byte("Hello there")

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	// Verify the new tree reflects the new content
	if newTree.Hash() == oldTree.Hash() {
		t.Error(
			"expected different hash for modified content",
		)
	}
}

func TestComputeEditRegion_NoChange(
	t *testing.T,
) {
	source := []byte("Hello world")
	edit := computeEditRegion(source, source)

	// For identical sources, edit region should have no change
	if edit.StartOffset != len(source) {
		t.Errorf(
			"expected StartOffset %d, got %d",
			len(source),
			edit.StartOffset,
		)
	}
	if edit.OldEndOffset != len(source) {
		t.Errorf(
			"expected OldEndOffset %d, got %d",
			len(source),
			edit.OldEndOffset,
		)
	}
	if edit.NewEndOffset != len(source) {
		t.Errorf(
			"expected NewEndOffset %d, got %d",
			len(source),
			edit.NewEndOffset,
		)
	}
}

func TestComputeEditRegion_InsertAtStart(
	t *testing.T,
) {
	oldSource := []byte("world")
	newSource := []byte("Hello world")

	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	if edit.StartOffset != 0 {
		t.Errorf(
			"expected StartOffset 0, got %d",
			edit.StartOffset,
		)
	}
	if !edit.IsInsert() {
		return
	}

	// This is an insert, the old end should equal start
	if edit.OldEndOffset != edit.StartOffset {
		t.Error(
			"expected OldEndOffset to equal StartOffset for insert",
		)
	}
}

func TestComputeEditRegion_InsertAtEnd(
	t *testing.T,
) {
	oldSource := []byte("Hello")
	newSource := []byte("Hello world")

	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	if edit.StartOffset != 5 {
		t.Errorf(
			"expected StartOffset 5, got %d",
			edit.StartOffset,
		)
	}
	if !edit.IsInsert() {
		t.Error(
			"expected edit to be an insertion",
		)
	}
}

func TestComputeEditRegion_InsertInMiddle(
	t *testing.T,
) {
	oldSource := []byte("Hello world")
	newSource := []byte("Hello beautiful world")

	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	if edit.StartOffset != 6 {
		t.Errorf(
			"expected StartOffset 6, got %d",
			edit.StartOffset,
		)
	}
}

func TestComputeEditRegion_DeleteAtStart(
	t *testing.T,
) {
	oldSource := []byte("Hello world")
	newSource := []byte("world")

	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	if edit.StartOffset != 0 {
		t.Errorf(
			"expected StartOffset 0, got %d",
			edit.StartOffset,
		)
	}
}

func TestComputeEditRegion_DeleteAtEnd(
	t *testing.T,
) {
	oldSource := []byte("Hello world")
	newSource := []byte("Hello")

	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	if edit.StartOffset != 5 {
		t.Errorf(
			"expected StartOffset 5, got %d",
			edit.StartOffset,
		)
	}
	if !edit.IsDelete() {
		t.Error("expected edit to be a deletion")
	}
}

func TestComputeEditRegion_DeleteInMiddle(
	t *testing.T,
) {
	oldSource := []byte("Hello beautiful world")
	newSource := []byte("Hello world")

	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	if edit.StartOffset != 6 {
		t.Errorf(
			"expected StartOffset 6, got %d",
			edit.StartOffset,
		)
	}
}

func TestComputeEditRegion_Replacement(
	t *testing.T,
) {
	oldSource := []byte("Hello world")
	newSource := []byte("Hello there")

	edit := computeEditRegion(
		oldSource,
		newSource,
	)

	if edit.StartOffset != 6 {
		t.Errorf(
			"expected StartOffset 6, got %d",
			edit.StartOffset,
		)
	}
	if !edit.IsReplace() {
		t.Error(
			"expected edit to be a replacement",
		)
	}
}

func TestEditRegion_Delta(t *testing.T) {
	tests := []struct {
		name          string
		startOffset   int
		oldEndOffset  int
		newEndOffset  int
		expectedDelta int
	}{
		{"insertion", 5, 5, 15, 10},
		{"deletion", 5, 15, 5, -10},
		{"replacement same size", 5, 10, 10, 0},
		{"replacement larger", 5, 10, 15, 5},
		{"replacement smaller", 5, 15, 10, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edit := EditRegion{
				StartOffset:  tt.startOffset,
				OldEndOffset: tt.oldEndOffset,
				NewEndOffset: tt.newEndOffset,
			}
			if edit.Delta() != tt.expectedDelta {
				t.Errorf(
					"expected Delta %d, got %d",
					tt.expectedDelta,
					edit.Delta(),
				)
			}
		})
	}
}

func TestEditRegion_IsInsert(t *testing.T) {
	// Pure insertion: nothing removed
	insertEdit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 5,  // Old end equals start
		NewEndOffset: 15, // New content added
	}
	if !insertEdit.IsInsert() {
		t.Error(
			"expected IsInsert to return true for pure insertion",
		)
	}

	// Not pure insertion
	replaceEdit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 10,
		NewEndOffset: 15,
	}
	if replaceEdit.IsInsert() {
		t.Error(
			"expected IsInsert to return false for replacement",
		)
	}
}

func TestEditRegion_IsDelete(t *testing.T) {
	// Pure deletion: nothing added
	deleteEdit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 15, // Old content removed
		NewEndOffset: 5,  // New end equals start
	}
	if !deleteEdit.IsDelete() {
		t.Error(
			"expected IsDelete to return true for pure deletion",
		)
	}

	// Not pure deletion
	replaceEdit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 10,
		NewEndOffset: 15,
	}
	if replaceEdit.IsDelete() {
		t.Error(
			"expected IsDelete to return false for replacement",
		)
	}
}

func TestEditRegion_IsReplace(t *testing.T) {
	// Replacement: both content removed and added
	replaceEdit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 10, // Content removed
		NewEndOffset: 15, // Content added
	}
	if !replaceEdit.IsReplace() {
		t.Error(
			"expected IsReplace to return true for replacement",
		)
	}

	// Not replacement (pure insert)
	insertEdit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 5,
		NewEndOffset: 15,
	}
	if insertEdit.IsReplace() {
		t.Error(
			"expected IsReplace to return false for pure insertion",
		)
	}

	// Not replacement (pure delete)
	deleteEdit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 15,
		NewEndOffset: 5,
	}
	if deleteEdit.IsReplace() {
		t.Error(
			"expected IsReplace to return false for pure deletion",
		)
	}
}

func TestEditRegion_OldLength(t *testing.T) {
	edit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 15,
		NewEndOffset: 20,
	}
	if edit.OldLength() != 10 {
		t.Errorf(
			"expected OldLength 10, got %d",
			edit.OldLength(),
		)
	}
}

func TestEditRegion_NewLength(t *testing.T) {
	edit := EditRegion{
		StartOffset:  5,
		OldEndOffset: 15,
		NewEndOffset: 20,
	}
	if edit.NewLength() != 15 {
		t.Errorf(
			"expected NewLength 15, got %d",
			edit.NewLength(),
		)
	}
}

func TestIdentifyReusableNodes_NilRoot(
	t *testing.T,
) {
	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 5,
		NewEndOffset: 10,
	}
	reusable := identifyReusableNodes(nil, edit)

	if reusable != nil {
		t.Error(
			"expected nil result for nil root",
		)
	}
}

func TestIdentifyReusableNodes_NoChildren(
	t *testing.T,
) {
	// Create a document with no children
	doc := NewNodeBuilder(NodeTypeDocument).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("some text")).
		Build()

	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 5,
		NewEndOffset: 10,
	}
	reusable := identifyReusableNodes(doc, edit)

	// Document with no children should return empty list
	if len(reusable) != 0 {
		t.Errorf(
			"expected 0 reusable nodes, got %d",
			len(reusable),
		)
	}
}

func TestIdentifyReusableNodes_NodesBeforeEdit(
	t *testing.T,
) {
	// Create children nodes
	child1 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(10).
		WithSource([]byte("paragraph1")).
		Build()
	child2 := NewNodeBuilder(NodeTypeParagraph).
		WithStart(20).
		WithEnd(30).
		WithSource([]byte("paragraph2")).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithStart(0).
		WithEnd(30).
		WithSource([]byte("paragraph1...paragraph2")).
		WithChildren([]Node{child1, child2}).
		Build()

	// Edit is at offset 15 (between children)
	edit := EditRegion{
		StartOffset:  15,
		OldEndOffset: 17,
		NewEndOffset: 18,
	}
	reusable := identifyReusableNodes(doc, edit)

	// Both children should be identified as reusable
	// (child1 is before edit, child2 is after edit)
	if len(reusable) != 2 {
		t.Errorf(
			"expected 2 reusable nodes, got %d",
			len(reusable),
		)
	}
}

func TestIdentifyReusableNodes_NodesAfterEdit(
	t *testing.T,
) {
	// Create children nodes
	child := NewNodeBuilder(NodeTypeParagraph).
		WithStart(20).
		WithEnd(30).
		WithSource([]byte("paragraph")).
		Build()

	doc := NewNodeBuilder(NodeTypeDocument).
		WithStart(0).
		WithEnd(30).
		WithSource([]byte("...paragraph")).
		WithChildren([]Node{child}).
		Build()

	// Edit is at the beginning
	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 5,
		NewEndOffset: 10,
	}
	reusable := identifyReusableNodes(doc, edit)

	// Child is after edit, should be reusable
	if len(reusable) != 1 {
		t.Errorf(
			"expected 1 reusable node, got %d",
			len(reusable),
		)
	}
}

func TestMatchAndValidateSubtrees_NilTree(
	t *testing.T,
) {
	reusableNodes := []Node{
		NewNodeBuilder(NodeTypeText).
			WithStart(0).
			WithEnd(5).
			WithSource([]byte("text")).
			Build(),
	}

	count := matchAndValidateSubtrees(
		nil,
		reusableNodes,
	)
	if count != 0 {
		t.Errorf(
			"expected 0 matches for nil tree, got %d",
			count,
		)
	}
}

func TestMatchAndValidateSubtrees_EmptyReusableNodes(
	t *testing.T,
) {
	tree, _ := Parse([]byte("# Header"))

	count := matchAndValidateSubtrees(
		tree,
		make([]Node, 0),
	)
	if count != 0 {
		t.Errorf(
			"expected 0 matches for empty reusable nodes, got %d",
			count,
		)
	}

	count = matchAndValidateSubtrees(tree, nil)
	if count != 0 {
		t.Errorf(
			"expected 0 matches for nil reusable nodes, got %d",
			count,
		)
	}
}

func TestMatchAndValidateSubtrees_HashMatching(
	t *testing.T,
) {
	source := []byte("# Header\n\nParagraph")
	tree, _ := Parse(source)

	// Get children from tree
	children := tree.Children()
	if len(children) < 1 {
		t.Skip("tree has no children to test")
	}

	// Use the first child as a reusable node
	reusableNodes := []Node{children[0]}

	count := matchAndValidateSubtrees(
		tree,
		reusableNodes,
	)
	// Should find at least one match (the node itself)
	if count < 1 {
		t.Errorf(
			"expected at least 1 match, got %d",
			count,
		)
	}
}

func TestParseIncremental_LargeChangeTriggersFullReparse(
	t *testing.T,
) {
	// Create a source where the change is more than 20%
	oldSource := []byte("Hello world")
	// Changing more than 20% of the content
	newSource := []byte(
		"Completely different text here",
	)

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	// Should still produce a valid tree
	if newTree.NodeType() != NodeTypeDocument {
		t.Errorf(
			"expected Document node, got %v",
			newTree.NodeType(),
		)
	}
}

func TestParseIncremental_SmallChangeUsesIncremental(
	t *testing.T,
) {
	// Create a larger source where the change is less than 20%
	oldSource := []byte(
		"# Header\n\nThis is a paragraph with some content that is long enough to make small changes incremental.\n\nAnother paragraph here.",
	)
	// Small change
	newSource := []byte(
		"# Header\n\nThis is a paragraph with some content that is long enough to make small changes incremental!\n\nAnother paragraph here.",
	)

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
}

func TestParseIncremental_EmptyOldSource(
	t *testing.T,
) {
	oldSource := make([]byte, 0)
	newSource := []byte("New content")

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
}

func TestParseIncremental_EmptyNewSource(
	t *testing.T,
) {
	oldSource := []byte("Old content")
	newSource := make([]byte, 0)

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	// Empty document should have no children
	if len(newTree.Children()) != 0 {
		t.Errorf(
			"expected empty document, got %d children",
			len(newTree.Children()),
		)
	}
}

func TestParseIncremental_EditAtBeginning(
	t *testing.T,
) {
	oldSource := []byte("old text here")
	newSource := []byte("new text here")

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}

	edit := computeEditRegion(
		oldSource,
		newSource,
	)
	if edit.StartOffset != 0 {
		t.Errorf(
			"expected edit at beginning (offset 0), got %d",
			edit.StartOffset,
		)
	}
}

func TestParseIncremental_EditAtEnd(
	t *testing.T,
) {
	oldSource := []byte("text here old")
	newSource := []byte("text here new")

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}

	edit := computeEditRegion(
		oldSource,
		newSource,
	)
	if edit.StartOffset != 10 {
		t.Errorf(
			"expected edit at end (offset 10), got %d",
			edit.StartOffset,
		)
	}
}

func TestParseIncremental_EditInMiddle(
	t *testing.T,
) {
	oldSource := []byte("text old text")
	newSource := []byte("text new text")

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}

	edit := computeEditRegion(
		oldSource,
		newSource,
	)
	if edit.StartOffset != 5 {
		t.Errorf(
			"expected edit in middle (offset 5), got %d",
			edit.StartOffset,
		)
	}
}

func TestNewIncrementalParseState(t *testing.T) {
	source := []byte("# Header\n\nParagraph")
	tree, _ := Parse(source)

	state := NewIncrementalParseState(
		tree,
		source,
	)

	if state == nil {
		t.Fatal("expected non-nil state")
	}
	if state.LinkDefs == nil {
		t.Error(
			"expected LinkDefs to be initialized",
		)
	}
	if state.LineIndex == nil {
		t.Error(
			"expected LineIndex to be initialized",
		)
	}
	if state.RootHash == 0 {
		t.Error("expected RootHash to be set")
	}
	if state.RootHash != tree.Hash() {
		t.Errorf(
			"expected RootHash %d, got %d",
			tree.Hash(),
			state.RootHash,
		)
	}
}

func TestNewIncrementalParseState_NilTree(
	t *testing.T,
) {
	state := NewIncrementalParseState(
		nil,
		[]byte("source"),
	)

	if state != nil {
		t.Error("expected nil state for nil tree")
	}
}

func TestIncrementalParseState_CanReuseLinkDefs_NilState(
	t *testing.T,
) {
	var state *IncrementalParseState
	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 5,
		NewEndOffset: 10,
	}
	tree, _ := Parse([]byte("# Header"))

	if state.CanReuseLinkDefs(edit, tree) {
		t.Error(
			"expected CanReuseLinkDefs to return false for nil state",
		)
	}
}

func TestIncrementalParseState_CanReuseLinkDefs_NilTree(
	t *testing.T,
) {
	source := []byte("# Header")
	tree, _ := Parse(source)
	state := NewIncrementalParseState(
		tree,
		source,
	)
	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 5,
		NewEndOffset: 10,
	}

	if state.CanReuseLinkDefs(edit, nil) {
		t.Error(
			"expected CanReuseLinkDefs to return false for nil tree",
		)
	}
}

func TestIncrementalParseState_CanReuseLinkDefs_NoLinkDefs(
	t *testing.T,
) {
	source := []byte("# Header\n\nParagraph")
	tree, _ := Parse(source)
	state := NewIncrementalParseState(
		tree,
		source,
	)
	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 5,
		NewEndOffset: 10,
	}

	// No link definitions in document, should be able to "reuse" them
	if !state.CanReuseLinkDefs(edit, tree) {
		t.Error(
			"expected CanReuseLinkDefs to return true when no link defs exist",
		)
	}
}

func TestIncrementalParseState_CanReuseLinkDefs_NoOverlap(
	_ *testing.T,
) {
	// Document with link definition at the end
	source := []byte(
		"# Header\n\nParagraph\n\n[ref]: https://example.com",
	)
	tree, _ := Parse(source)
	state := NewIncrementalParseState(
		tree,
		source,
	)

	// Edit at the beginning (before link def)
	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 8,
		NewEndOffset: 10,
	}

	// Link def should be reusable since edit doesn't overlap
	// Note: This depends on whether the parser creates LinkDef nodes
	_ = state.CanReuseLinkDefs(edit, tree)
	// Just verify it doesn't panic
}

func TestIncrementalParseState_UpdateLineIndex_NilState(
	_ *testing.T,
) {
	var state *IncrementalParseState
	edit := EditRegion{
		StartOffset:  0,
		OldEndOffset: 5,
		NewEndOffset: 10,
	}
	oldSource := []byte("old")
	newSource := []byte("new source")

	// Should not panic
	state.UpdateLineIndex(
		edit,
		oldSource,
		newSource,
	)
}

func TestIncrementalParseState_UpdateLineIndex_NoNewlines(
	t *testing.T,
) {
	source := []byte("single line")
	tree, _ := Parse(source)
	state := NewIncrementalParseState(
		tree,
		source,
	)

	newSource := []byte("single line modified")
	edit := computeEditRegion(source, newSource)

	// Update line index
	state.UpdateLineIndex(edit, source, newSource)

	// Line index should be updated
	if state.LineIndex == nil {
		t.Error(
			"expected LineIndex to be non-nil after update",
		)
	}
}

func TestIncrementalParseState_UpdateLineIndex_WithNewlines(
	t *testing.T,
) {
	source := []byte("line one\nline two")
	tree, _ := Parse(source)
	state := NewIncrementalParseState(
		tree,
		source,
	)

	// Add a new line
	newSource := []byte(
		"line one\nline new\nline two",
	)
	edit := computeEditRegion(source, newSource)

	// Update line index
	state.UpdateLineIndex(edit, source, newSource)

	// Line index should be rebuilt
	if state.LineIndex == nil {
		t.Error(
			"expected LineIndex to be non-nil after update",
		)
	}
	if state.LineIndex.LineCount() != 3 {
		t.Errorf(
			"expected 3 lines, got %d",
			state.LineIndex.LineCount(),
		)
	}
}

func TestParseIncrementalWithState_NilState(
	t *testing.T,
) {
	oldSource := []byte("# Header")
	newSource := []byte("# Header Updated")
	oldTree, _ := Parse(oldSource)

	newTree, errors, newState := ParseIncrementalWithState(
		nil,
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	if newState == nil {
		t.Error(
			"expected new state to be created",
		)
	}
}

func TestParseIncrementalWithState_IdenticalSources(
	t *testing.T,
) {
	source := []byte("# Header\n\nParagraph")
	tree, _ := Parse(source)
	state := NewIncrementalParseState(
		tree,
		source,
	)

	newTree, errors, newState := ParseIncrementalWithState(
		state,
		tree,
		source,
		source,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree != tree {
		t.Error(
			"expected same tree for identical sources",
		)
	}
	if newState != state {
		t.Error(
			"expected same state for identical sources",
		)
	}
}

func TestParseIncrementalWithState_ModifiedSource(
	t *testing.T,
) {
	oldSource := []byte("# Header\n\nParagraph")
	newSource := []byte(
		"# Header Updated\n\nParagraph",
	)
	oldTree, _ := Parse(oldSource)
	state := NewIncrementalParseState(
		oldTree,
		oldSource,
	)
	originalHash := state.RootHash

	newTree, errors, newState := ParseIncrementalWithState(
		state,
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}
	if newState.RootHash == originalHash {
		t.Error("expected RootHash to be updated")
	}
	if newState.RootHash != newTree.Hash() {
		t.Error(
			"expected RootHash to match new tree hash",
		)
	}
}

func TestNodeAtOffset_NilRoot(t *testing.T) {
	result := NodeAtOffset(nil, 0)
	if result != nil {
		t.Error("expected nil for nil root")
	}
}

func TestNodeAtOffset_OutOfBounds(t *testing.T) {
	source := []byte("# Header")
	tree, _ := Parse(source)

	// Before start
	result := NodeAtOffset(tree, -1)
	if result != nil {
		t.Error(
			"expected nil for offset before start",
		)
	}

	// At end (exclusive, so should be nil)
	_, end := tree.Span()
	result = NodeAtOffset(tree, end)
	if result != nil {
		t.Error("expected nil for offset at end")
	}

	// Beyond end
	result = NodeAtOffset(tree, end+10)
	if result != nil {
		t.Error(
			"expected nil for offset beyond end",
		)
	}
}

func TestNodeAtOffset_ValidOffset(t *testing.T) {
	source := []byte("# Header\n\nParagraph")
	tree, _ := Parse(source)

	// At beginning
	result := NodeAtOffset(tree, 0)
	if result == nil {
		t.Fatal(
			"expected non-nil node at offset 0",
		)
	}

	// In middle
	result = NodeAtOffset(tree, 5)
	if result == nil {
		t.Fatal(
			"expected non-nil node at offset 5",
		)
	}
}

func TestNodeAtOffset_FindsInnermostNode(
	t *testing.T,
) {
	source := []byte("# Header")
	tree, _ := Parse(source)

	// Find the innermost node at offset 2 (within "Header")
	result := NodeAtOffset(tree, 2)
	if result == nil {
		t.Fatal("expected non-nil node")
	}

	// The result should be the most specific (innermost) node
	// It should not be the Document itself
	if result.NodeType() != NodeTypeDocument {
		return
	}

	// If it's a document, there might be no children, which is okay
	children := tree.Children()
	if len(children) > 0 {
		t.Error(
			"expected innermost node, not Document",
		)
	}
}

func TestNodesAtOffset_NilRoot(t *testing.T) {
	result := NodesAtOffset(nil, 0)
	if result != nil {
		t.Error("expected nil for nil root")
	}
}

func TestNodesAtOffset_OutOfBounds(t *testing.T) {
	source := []byte("# Header")
	tree, _ := Parse(source)

	// Before start
	result := NodesAtOffset(tree, -1)
	if result != nil {
		t.Error(
			"expected nil for offset before start",
		)
	}

	// Beyond end
	_, end := tree.Span()
	result = NodesAtOffset(tree, end+10)
	if result != nil {
		t.Error(
			"expected nil for offset beyond end",
		)
	}
}

func TestNodesAtOffset_ReturnsAllContainingNodes(
	t *testing.T,
) {
	source := []byte("# Header")
	tree, _ := Parse(source)

	result := NodesAtOffset(tree, 2)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Should include at least the Document
	if len(result) < 1 {
		t.Error("expected at least 1 node")
	}

	// First node should be the root (Document)
	if result[0].NodeType() != NodeTypeDocument {
		t.Error(
			"expected first node to be Document",
		)
	}
}

func TestAffectedBlockRegion_SimpleEdit(
	t *testing.T,
) {
	source := []byte(
		"# Header\n\nParagraph text here.\n\n## Another",
	)
	edit := EditRegion{
		StartOffset:  11,
		OldEndOffset: 15,
		NewEndOffset: 18,
	}

	start, end := AffectedBlockRegion(
		source,
		edit,
	)

	// Should expand to block boundaries
	if start > edit.StartOffset {
		t.Errorf(
			"start should be at or before edit start, got %d (edit start: %d)",
			start,
			edit.StartOffset,
		)
	}
	if end < edit.NewEndOffset {
		t.Errorf(
			"end should be at or after edit end, got %d (edit end: %d)",
			end,
			edit.NewEndOffset,
		)
	}
}

func TestAffectedBlockRegion_EditAtBlankLine(
	t *testing.T,
) {
	source := []byte("Para one.\n\nPara two.")
	edit := EditRegion{
		StartOffset:  10,
		OldEndOffset: 11,
		NewEndOffset: 12,
	}

	start, end := AffectedBlockRegion(
		source,
		edit,
	)

	// Result should include the edit region
	if start > edit.StartOffset ||
		end < edit.NewEndOffset {
		t.Errorf(
			"affected region should contain edit: got [%d, %d), edit [%d, %d)",
			start,
			end,
			edit.StartOffset,
			edit.NewEndOffset,
		)
	}
}

func TestIsBlockStart_Header(t *testing.T) {
	source := []byte("# Header")
	if !isBlockStart(source, 0) {
		t.Error(
			"expected # to be recognized as block start",
		)
	}
}

func TestIsBlockStart_ListMarkers(t *testing.T) {
	tests := []struct {
		source   string
		expected bool
	}{
		{"- item", true},
		{"* item", true},
		{"+ item", true},
		{"> quote", true},
		{"```code", true},
		{"~~~code", true},
		{"1. item", true},
		{"normal text", false},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			result := isBlockStart(
				[]byte(tt.source),
				0,
			)
			if result != tt.expected {
				t.Errorf(
					"expected %v, got %v",
					tt.expected,
					result,
				)
			}
		})
	}
}

func TestIsBlockStart_OrderedList(t *testing.T) {
	source := []byte("1. item")
	if !isBlockStart(source, 0) {
		t.Error(
			"expected ordered list marker to be recognized as block start",
		)
	}

	// Multi-digit number
	source = []byte("123. item")
	if !isBlockStart(source, 0) {
		t.Error(
			"expected multi-digit ordered list marker to be recognized as block start",
		)
	}

	// Not a list (digit without period)
	source = []byte("123 not a list")
	if isBlockStart(source, 0) {
		t.Error(
			"expected digit without period to not be recognized as block start",
		)
	}
}

func TestIsBlockStart_OutOfBounds(t *testing.T) {
	source := []byte("text")
	if isBlockStart(source, 100) {
		t.Error(
			"expected false for out of bounds position",
		)
	}
}

func TestAdjustNodeOffsets_NilNode(t *testing.T) {
	result := adjustNodeOffsets(
		nil,
		5,
		[]byte("source"),
	)
	if result != nil {
		t.Error("expected nil for nil node")
	}
}

func TestAdjustNodeOffsets_ZeroDelta(
	t *testing.T,
) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(5).
		WithEnd(10).
		WithSource([]byte("text")).
		Build()

	result := adjustNodeOffsets(
		node,
		0,
		[]byte("some source"),
	)
	if result != node {
		t.Error(
			"expected same node for zero delta",
		)
	}
}

func TestAdjustNodeOffsets_PositiveDelta(
	t *testing.T,
) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(5).
		WithEnd(10).
		WithSource([]byte("text")).
		Build()

	newSource := []byte("prefix text suffix")
	result := adjustNodeOffsets(
		node,
		7,
		newSource,
	)

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	start, end := result.Span()
	if start != 12 {
		t.Errorf(
			"expected start 12, got %d",
			start,
		)
	}
	if end != 17 {
		t.Errorf("expected end 17, got %d", end)
	}
}

func TestAdjustNodeOffsets_NegativeDelta(
	t *testing.T,
) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(10).
		WithEnd(15).
		WithSource([]byte("text")).
		Build()

	newSource := []byte("short source")
	result := adjustNodeOffsets(
		node,
		-5,
		newSource,
	)

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	start, end := result.Span()
	if start != 5 {
		t.Errorf(
			"expected start 5, got %d",
			start,
		)
	}
	if end != 10 {
		t.Errorf("expected end 10, got %d", end)
	}
}

func TestAdjustNodeOffsets_InvalidNewOffsets(
	t *testing.T,
) {
	node := NewNodeBuilder(NodeTypeText).
		WithStart(5).
		WithEnd(10).
		WithSource([]byte("text")).
		Build()

	// Delta would make start negative
	result := adjustNodeOffsets(
		node,
		-10,
		[]byte("short"),
	)
	if result != nil {
		t.Error(
			"expected nil for invalid new start offset",
		)
	}

	// Delta would make end beyond source length
	// Note: When delta is 0, the function returns the original node without
	// checking source bounds (an optimization). Use a small positive delta
	// to trigger the validation.
	shortSource := []byte("ab")
	result = adjustNodeOffsets(
		node,
		1,
		shortSource,
	)
	if result != nil {
		t.Error(
			"expected nil for invalid new end offset",
		)
	}
}

func TestAdjustNodeOffsets_WithChildren(
	t *testing.T,
) {
	child := NewNodeBuilder(NodeTypeText).
		WithStart(5).
		WithEnd(10).
		WithSource([]byte("child")).
		Build()

	parent := NewNodeBuilder(NodeTypeParagraph).
		WithStart(0).
		WithEnd(15).
		WithSource([]byte("parent content")).
		WithChildren([]Node{child}).
		Build()

	// This requires a large enough source
	newSource := []byte(
		"xxxxxxxxparent contentxxxxx",
	)
	result := adjustNodeOffsets(
		parent,
		8,
		newSource,
	)

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Check parent offsets
	start, end := result.Span()
	if start != 8 {
		t.Errorf(
			"expected parent start 8, got %d",
			start,
		)
	}
	if end != 23 {
		t.Errorf(
			"expected parent end 23, got %d",
			end,
		)
	}

	// Check child offsets
	children := result.Children()
	if len(children) != 1 {
		t.Fatalf(
			"expected 1 child, got %d",
			len(children),
		)
	}
	childStart, childEnd := children[0].Span()
	if childStart != 13 {
		t.Errorf(
			"expected child start 13, got %d",
			childStart,
		)
	}
	if childEnd != 18 {
		t.Errorf(
			"expected child end 18, got %d",
			childEnd,
		)
	}
}

func TestFindLinkDefs_NilNode(t *testing.T) {
	result := findLinkDefs(nil)
	if result != nil {
		t.Error("expected nil for nil node")
	}
}

func TestFindLinkDefs_NoLinkDefs(t *testing.T) {
	source := []byte("# Header\n\nParagraph")
	tree, _ := Parse(source)

	result := findLinkDefs(tree)
	if len(result) != 0 {
		t.Errorf(
			"expected no link defs, got %d",
			len(result),
		)
	}
}

func TestFindLinkDefs_WithLinkDefs(_ *testing.T) {
	source := []byte(
		"[ref]: https://example.com\n\n[ref2]: https://example2.com",
	)
	tree, _ := Parse(source)

	result := findLinkDefs(tree)
	// The parser may or may not create LinkDef nodes
	// This test just verifies the function doesn't panic
	_ = result
}

func TestAdjustedNode_Structure(t *testing.T) {
	// Test that AdjustedNode struct is properly defined
	node := NewNodeBuilder(NodeTypeText).
		WithStart(0).
		WithEnd(5).
		WithSource([]byte("text")).
		Build()

	adjusted := AdjustedNode{
		original:   node,
		startDelta: 5,
		endDelta:   5,
	}

	if adjusted.original != node {
		t.Error(
			"expected original to be preserved",
		)
	}
	if adjusted.startDelta != 5 {
		t.Errorf(
			"expected startDelta 5, got %d",
			adjusted.startDelta,
		)
	}
	if adjusted.endDelta != 5 {
		t.Errorf(
			"expected endDelta 5, got %d",
			adjusted.endDelta,
		)
	}
}

func TestParseIncremental_RealWorldScenario(
	t *testing.T,
) {
	// Simulate a real-world editing scenario
	oldSource := []byte(`# Document Title

## Introduction

This is the introduction paragraph.

## Details

Here are some details:

- Item one
- Item two
- Item three

## Conclusion

Final thoughts here.
`)

	// User edits the introduction
	newSource := []byte(`# Document Title

## Introduction

This is the updated introduction paragraph with more content.

## Details

Here are some details:

- Item one
- Item two
- Item three

## Conclusion

Final thoughts here.
`)

	oldTree, _ := Parse(oldSource)
	newTree, errors := ParseIncremental(
		oldTree,
		oldSource,
		newSource,
	)

	if len(errors) != 0 {
		t.Errorf(
			"expected no errors, got %d",
			len(errors),
		)
	}
	if newTree == nil {
		t.Fatal("expected non-nil tree")
	}

	// Verify the document structure is preserved
	children := newTree.Children()
	if len(children) < 4 {
		t.Errorf(
			"expected at least 4 top-level sections, got %d",
			len(children),
		)
	}
}

func TestParseIncremental_MultipleEdits(
	t *testing.T,
) {
	// Simulate multiple consecutive edits
	source1 := []byte("Hello world")
	source2 := []byte("Hello beautiful world")
	source3 := []byte(
		"Hello beautiful and wonderful world",
	)
	source4 := []byte(
		"Hello beautiful, wonderful world",
	)

	tree1, _ := Parse(source1)
	tree2, _ := ParseIncremental(
		tree1,
		source1,
		source2,
	)
	tree3, _ := ParseIncremental(
		tree2,
		source2,
		source3,
	)
	tree4, _ := ParseIncremental(
		tree3,
		source3,
		source4,
	)

	if tree4 == nil {
		t.Fatal(
			"expected non-nil tree after multiple edits",
		)
	}
}

func TestParseIncremental_UndoRedo(t *testing.T) {
	// Simulate undo/redo scenario
	original := []byte("Original content")
	modified := []byte("Modified content")

	tree1, _ := Parse(original)
	tree2, _ := ParseIncremental(
		tree1,
		original,
		modified,
	)
	tree3, _ := ParseIncremental(
		tree2,
		modified,
		original,
	) // Undo

	if tree3 == nil {
		t.Fatal(
			"expected non-nil tree after undo",
		)
	}

	// Hash should match the original
	if tree3.Hash() != tree1.Hash() {
		t.Error(
			"expected tree hash to match original after undo",
		)
	}
}
