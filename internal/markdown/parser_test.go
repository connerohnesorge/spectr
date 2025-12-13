package markdown

import (
	"testing"

	"github.com/connerohnesorge/spectr/internal/specterrs"
)

func TestParseDocument_Headers(t *testing.T) {
	content := []byte(`# Main Title

Some intro text.

## Requirements

### Requirement: Auth

#### Scenario: Valid Login
`)

	doc, err := ParseDocument(content)
	if err != nil {
		t.Fatalf("ParseDocument failed: %v", err)
	}

	if len(doc.Headers) != 4 {
		t.Errorf("Expected 4 headers, got %d", len(doc.Headers))
	}

	expected := []struct {
		level int
		text  string
	}{
		{1, "Main Title"},
		{2, "Requirements"},
		{3, "Requirement: Auth"},
		{4, "Scenario: Valid Login"},
	}

	for i, exp := range expected {
		if i >= len(doc.Headers) {
			break
		}
		h := doc.Headers[i]
		if h.Level != exp.level {
			t.Errorf("Header %d: expected level %d, got %d", i, exp.level, h.Level)
		}
		if h.Text != exp.text {
			t.Errorf("Header %d: expected text %q, got %q", i, exp.text, h.Text)
		}
		if h.Line < 1 {
			t.Errorf("Header %d: line should be >= 1, got %d", i, h.Line)
		}
	}
}

func TestParseDocument_Sections(t *testing.T) {
	content := []byte(`# Main Title

Intro text.

## Requirements

Requirement content here.

## Other Section

Other content.
`)

	doc, err := ParseDocument(content)
	if err != nil {
		t.Fatalf("ParseDocument failed: %v", err)
	}

	if len(doc.Sections) == 0 {
		t.Fatal("Expected sections to be extracted")
	}

	// Check that Requirements section exists and has content
	reqSection, ok := doc.Sections["Requirements"]
	if !ok {
		t.Error("Expected 'Requirements' section")
	} else {
		if reqSection.Header.Level != 2 {
			t.Errorf("Requirements section header level: expected 2, got %d", reqSection.Header.Level)
		}
		if len(reqSection.Content) == 0 {
			t.Error("Requirements section content should not be empty")
		}
	}
}

func TestParseDocument_Tasks(t *testing.T) {
	content := []byte(`# Tasks

- [ ] Task 1
- [x] Task 2 completed
- [X] Task 3 also completed
- [ ] Task 4
`)

	doc, err := ParseDocument(content)
	if err != nil {
		t.Fatalf("ParseDocument failed: %v", err)
	}

	if len(doc.Tasks) < 3 {
		t.Errorf("Expected at least 3 tasks, got %d", len(doc.Tasks))
	}

	// Check first task is unchecked
	if len(doc.Tasks) > 0 && doc.Tasks[0].Checked {
		t.Error("First task should be unchecked")
	}

	// Check second task is checked
	if len(doc.Tasks) > 1 && !doc.Tasks[1].Checked {
		t.Error("Second task should be checked")
	}
}

func TestParseDocument_NestedTasks(t *testing.T) {
	content := []byte(`# Tasks

- [ ] Parent task
  - [ ] Child task 1
  - [x] Child task 2
- [ ] Another parent
`)

	doc, err := ParseDocument(content)
	if err != nil {
		t.Fatalf("ParseDocument failed: %v", err)
	}

	// Should have parent tasks with children
	if len(doc.Tasks) < 1 {
		t.Fatal("Expected at least 1 task")
	}

	// Check that nested tasks are captured (either as children or flat list)
	totalTasks := countTasks(doc.Tasks)
	if totalTasks < 3 {
		t.Errorf("Expected at least 3 total tasks (including nested), got %d", totalTasks)
	}
}

func countTasks(tasks []Task) int {
	count := len(tasks)
	for _, t := range tasks {
		count += countTasks(t.Children)
	}

	return count
}

func TestParseDocument_EmptyContent(t *testing.T) {
	_, err := ParseDocument([]byte(""))
	if err == nil {
		t.Error("Expected error for empty content")
	}

	_, ok := err.(*specterrs.EmptyContentError)
	if !ok {
		t.Errorf("Expected EmptyContentError, got %T", err)
	}
}

func TestParseDocument_WhitespaceOnly(t *testing.T) {
	_, err := ParseDocument([]byte("   \n\t\n  "))
	if err == nil {
		t.Error("Expected error for whitespace-only content")
	}

	_, ok := err.(*specterrs.EmptyContentError)
	if !ok {
		t.Errorf("Expected EmptyContentError, got %T", err)
	}
}

func TestParseDocument_BinaryContent(t *testing.T) {
	_, err := ParseDocument([]byte("Hello\x00World"))
	if err == nil {
		t.Error("Expected error for binary content")
	}

	_, ok := err.(*specterrs.BinaryContentError)
	if !ok {
		t.Errorf("Expected BinaryContentError, got %T", err)
	}
}

func TestParseDocument_LineNumbers(t *testing.T) {
	content := []byte(`# Title

## Section One

## Section Two
`)

	doc, err := ParseDocument(content)
	if err != nil {
		t.Fatalf("ParseDocument failed: %v", err)
	}

	if len(doc.Headers) < 3 {
		t.Fatal("Expected 3 headers")
	}

	// First header should be on line 1
	if doc.Headers[0].Line != 1 {
		t.Errorf("First header should be on line 1, got %d", doc.Headers[0].Line)
	}

	// Each subsequent header should be on a later line
	for i := 1; i < len(doc.Headers); i++ {
		if doc.Headers[i].Line <= doc.Headers[i-1].Line {
			t.Errorf("Header %d (line %d) should be after header %d (line %d)",
				i, doc.Headers[i].Line, i-1, doc.Headers[i-1].Line)
		}
	}
}
