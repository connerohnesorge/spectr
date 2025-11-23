package syntax

import (
	"testing"
)

func TestParser(t *testing.T) {
	input := `# Header
Some text
` + "```go" + `
code
` + "```" + `
- List item
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(doc.Nodes) != 5 {
		t.Fatalf("expected 5 nodes, got %d", len(doc.Nodes))
	}

	h, ok := doc.Nodes[0].(*Header)
	if !ok {
		t.Errorf("expected node 0 to be Header, got %T", doc.Nodes[0])
	} else {
		if h.Level != 1 {
			t.Errorf("expected header level 1, got %d", h.Level)
		}
		if h.Text != "Header" {
			t.Errorf("expected header text 'Header', got %q", h.Text)
		}
	}

	txt, ok := doc.Nodes[1].(*Text)
	if !ok {
		t.Errorf("expected node 1 to be Text, got %T", doc.Nodes[1])
	} else {
		if txt.Content != "Some text\n" {
			t.Errorf("expected text content 'Some text\\n', got %q", txt.Content)
		}
	}

	cb, ok := doc.Nodes[2].(*CodeBlock)
	if !ok {
		t.Errorf("expected node 2 to be CodeBlock, got %T", doc.Nodes[2])
	} else {
		if cb.Language != "go" {
			t.Errorf("expected language 'go', got %q", cb.Language)
		}
		if cb.Content != "code" {
			t.Errorf("expected content 'code', got %q", cb.Content)
		}
	}

	// Node 3 is the newline text
	txt2, ok := doc.Nodes[3].(*Text)
	if !ok {
		t.Errorf("expected node 3 to be Text, got %T", doc.Nodes[3])
	} else {
		if txt2.Content != "\n" {
			t.Errorf("expected text content '\\n', got %q", txt2.Content)
		}
	}

	lst, ok := doc.Nodes[4].(*List)
	if !ok {
		t.Errorf("expected node 4 to be List, got %T", doc.Nodes[4])
	} else {
		if lst.Marker != "-" {
			t.Errorf("expected marker '-', got %q", lst.Marker)
		}
		// Content should be "List item" without newline because parseList uses TrimSpace
		if lst.Content != "List item" {
			t.Errorf("expected content 'List item', got %q", lst.Content)
		}
	}
}
