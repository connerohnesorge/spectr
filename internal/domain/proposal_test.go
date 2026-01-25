package domain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProposalFrontmatter_ValidFrontmatter(t *testing.T) {
	content := []byte(`---
id: feat-dashboard
requires:
  - id: feat-auth
    reason: "needs user model"
  - id: feat-db
    reason: "needs schema"
enables:
  - id: feat-analytics
    reason: "unlocks tracking"
---

# My Proposal

Content here.
`)

	meta, err := ParseProposalFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.ID != "feat-dashboard" {
		t.Errorf("expected ID 'feat-dashboard', got %q", meta.ID)
	}

	if len(meta.Requires) != 2 {
		t.Fatalf("expected 2 requires, got %d", len(meta.Requires))
	}

	if meta.Requires[0].ID != "feat-auth" {
		t.Errorf("expected first require ID 'feat-auth', got %q", meta.Requires[0].ID)
	}
	if meta.Requires[0].Reason != "needs user model" {
		t.Errorf(
			"expected first require reason 'needs user model', got %q",
			meta.Requires[0].Reason,
		)
	}

	if meta.Requires[1].ID != "feat-db" {
		t.Errorf("expected second require ID 'feat-db', got %q", meta.Requires[1].ID)
	}

	if len(meta.Enables) != 1 {
		t.Fatalf("expected 1 enables, got %d", len(meta.Enables))
	}

	if meta.Enables[0].ID != "feat-analytics" {
		t.Errorf("expected enables ID 'feat-analytics', got %q", meta.Enables[0].ID)
	}
}

func TestParseProposalFrontmatter_NoFrontmatter(t *testing.T) {
	content := []byte(`# My Proposal

No frontmatter here.
`)

	meta, err := ParseProposalFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}

	if meta.ID != "" {
		t.Errorf("expected empty ID, got %q", meta.ID)
	}

	if len(meta.Requires) != 0 {
		t.Errorf("expected 0 requires, got %d", len(meta.Requires))
	}

	if len(meta.Enables) != 0 {
		t.Errorf("expected 0 enables, got %d", len(meta.Enables))
	}
}

func TestParseProposalFrontmatter_EmptyFile(t *testing.T) {
	content := []byte("")

	meta, err := ParseProposalFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
}

func TestParseProposalFrontmatter_EmptyFrontmatter(t *testing.T) {
	content := []byte(`---
---

# My Proposal
`)

	meta, err := ParseProposalFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}

	if meta.ID != "" {
		t.Errorf("expected empty ID, got %q", meta.ID)
	}
}

func TestParseProposalFrontmatter_MalformedYAML(t *testing.T) {
	content := []byte(`---
requires:
  - id: feat-auth
    reason: [invalid yaml
---

# My Proposal
`)

	_, err := ParseProposalFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestParseProposalFrontmatter_UnclosedFrontmatter(t *testing.T) {
	content := []byte(`---
id: test
requires:
  - id: something
`)

	_, err := ParseProposalFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for unclosed frontmatter")
	}
}

func TestParseProposalFrontmatter_OnlyDelimiter(t *testing.T) {
	content := []byte(`---`)

	_, err := ParseProposalFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for unclosed frontmatter")
	}
}

func TestParseProposalFrontmatter_EmptyRequires(t *testing.T) {
	content := []byte(`---
id: test
requires: []
---

# My Proposal
`)

	meta, err := ParseProposalFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(meta.Requires) != 0 {
		t.Errorf("expected 0 requires, got %d", len(meta.Requires))
	}
}

func TestParseProposalFrontmatter_MissingReason(t *testing.T) {
	content := []byte(`---
requires:
  - id: feat-auth
---

# My Proposal
`)

	meta, err := ParseProposalFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(meta.Requires) != 1 {
		t.Fatalf("expected 1 requires, got %d", len(meta.Requires))
	}

	if meta.Requires[0].Reason != "" {
		t.Errorf("expected empty reason, got %q", meta.Requires[0].Reason)
	}
}

func TestProposalMetadata_HasDependencies(t *testing.T) {
	tests := []struct {
		name     string
		meta     ProposalMetadata
		expected bool
	}{
		{
			name:     "no dependencies",
			meta:     ProposalMetadata{},
			expected: false,
		},
		{
			name: "has requires",
			meta: ProposalMetadata{
				Requires: []Dependency{{ID: "test"}},
			},
			expected: true,
		},
		{
			name: "only enables (not counted as dependencies)",
			meta: ProposalMetadata{
				Enables: []Dependency{{ID: "test"}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.meta.HasDependencies(); got != tt.expected {
				t.Errorf("HasDependencies() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProposalMetadata_RequiredIDs(t *testing.T) {
	meta := ProposalMetadata{
		Requires: []Dependency{
			{ID: "feat-auth", Reason: "reason1"},
			{ID: "feat-db", Reason: "reason2"},
		},
	}

	ids := meta.RequiredIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(ids))
	}

	if ids[0] != "feat-auth" {
		t.Errorf("expected first ID 'feat-auth', got %q", ids[0])
	}
	if ids[1] != "feat-db" {
		t.Errorf("expected second ID 'feat-db', got %q", ids[1])
	}
}

func TestProposalMetadata_EnabledIDs(t *testing.T) {
	meta := ProposalMetadata{
		Enables: []Dependency{
			{ID: "feat-analytics"},
		},
	}

	ids := meta.EnabledIDs()
	if len(ids) != 1 {
		t.Fatalf("expected 1 ID, got %d", len(ids))
	}

	if ids[0] != "feat-analytics" {
		t.Errorf("expected ID 'feat-analytics', got %q", ids[0])
	}
}

func TestValidateProposalMetadata_SelfReference(t *testing.T) {
	tests := []struct {
		name       string
		meta       ProposalMetadata
		proposalID string
		wantErr    bool
	}{
		{
			name: "no self-reference",
			meta: ProposalMetadata{
				Requires: []Dependency{{ID: "other"}},
			},
			proposalID: "test",
			wantErr:    false,
		},
		{
			name: "self-reference in requires",
			meta: ProposalMetadata{
				Requires: []Dependency{{ID: "test"}},
			},
			proposalID: "test",
			wantErr:    true,
		},
		{
			name: "self-reference in enables",
			meta: ProposalMetadata{
				Enables: []Dependency{{ID: "test"}},
			},
			proposalID: "test",
			wantErr:    true,
		},
		{
			name:       "empty metadata",
			meta:       ProposalMetadata{},
			proposalID: "test",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProposalMetadata(&tt.meta, tt.proposalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProposalMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseProposalFrontmatterFromFile(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	proposalPath := filepath.Join(tmpDir, "proposal.md")

	content := `---
id: test-proposal
requires:
  - id: feat-auth
---

# Test Proposal
`
	if err := os.WriteFile(proposalPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	meta, err := ParseProposalFrontmatterFromFile(proposalPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.ID != "test-proposal" {
		t.Errorf("expected ID 'test-proposal', got %q", meta.ID)
	}

	if len(meta.Requires) != 1 {
		t.Fatalf("expected 1 requires, got %d", len(meta.Requires))
	}
}

func TestParseProposalFrontmatterFromFile_NotFound(t *testing.T) {
	_, err := ParseProposalFrontmatterFromFile("/nonexistent/path/proposal.md")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestExtractFrontmatter_ContentWithoutNewlineAtEnd(t *testing.T) {
	// Test content where the closing delimiter has no trailing newline
	content := []byte("---\nid: test\n---")

	fm, err := ExtractFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(fm) != "id: test\n" {
		t.Errorf("expected 'id: test\\n', got %q", string(fm))
	}
}
