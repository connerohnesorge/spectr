package markdown

import (
	"strings"
	"testing"
)

func TestFindSectionContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		sectionHeader string
		wantContains  []string
		wantNotEmpty  bool
	}{
		{
			name: "simple section",
			content: `## Purpose
This is the purpose.

## Requirements
This is the requirements content.

## Notes
These are notes.`,
			sectionHeader: "Requirements",
			wantContains:  []string{"This is the requirements content."},
			wantNotEmpty:  true,
		},
		{
			name: "section at end",
			content: `## Purpose
Purpose content.

## Requirements
Final section content.`,
			sectionHeader: "Requirements",
			wantContains:  []string{"Final section content."},
			wantNotEmpty:  true,
		},
		{
			name:          "section not found",
			content:       "## Purpose\nContent",
			sectionHeader: "Requirements",
			wantNotEmpty:  false,
		},
		{
			name: "delta section",
			content: `## ADDED Requirements
Added content here.

## MODIFIED Requirements
Modified content here.`,
			sectionHeader: "ADDED Requirements",
			wantContains:  []string{"Added content here."},
			wantNotEmpty:  true,
		},
		{
			name:          "empty content",
			content:       "",
			sectionHeader: "Requirements",
			wantNotEmpty:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSectionContent(tt.content, tt.sectionHeader)
			if tt.wantNotEmpty && got == "" {
				t.Error("FindSectionContent() = empty, want non-empty")

				return
			}
			if !tt.wantNotEmpty && got != "" {
				t.Errorf("FindSectionContent() = %q, want empty", got)

				return
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("FindSectionContent() should contain %q, got %q", want, got)
				}
			}
		})
	}
}

func TestFindDeltaSectionContent(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		deltaType    string
		wantContains []string
		wantNotEmpty bool
	}{
		{
			name: "ADDED section",
			content: `## ADDED Requirements
Added content.

## MODIFIED Requirements
Modified content.`,
			deltaType:    "ADDED",
			wantContains: []string{"Added content."},
			wantNotEmpty: true,
		},
		{
			name: "MODIFIED section",
			content: `## ADDED Requirements
Added.

## MODIFIED Requirements
Modified content here.

## REMOVED Requirements
Removed.`,
			deltaType:    "MODIFIED",
			wantContains: []string{"Modified content here."},
			wantNotEmpty: true,
		},
		{
			name: "REMOVED section at end",
			content: `## ADDED Requirements
Added.

## REMOVED Requirements
Removed content.`,
			deltaType:    "REMOVED",
			wantContains: []string{"Removed content."},
			wantNotEmpty: true,
		},
		{
			name: "RENAMED section",
			content: `## RENAMED Requirements
- FROM: ### Requirement: Old
- TO: ### Requirement: New`,
			deltaType:    "RENAMED",
			wantContains: []string{"FROM:", "TO:"},
			wantNotEmpty: true,
		},
		{
			name:         "missing section",
			content:      "## ADDED Requirements\nContent",
			deltaType:    "MODIFIED",
			wantNotEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDeltaSectionContent(tt.content, tt.deltaType)
			if tt.wantNotEmpty && got == "" {
				t.Error("FindDeltaSectionContent() = empty, want non-empty")

				return
			}
			if !tt.wantNotEmpty && got != "" {
				t.Errorf("FindDeltaSectionContent() = %q, want empty", got)

				return
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("FindDeltaSectionContent() should contain %q, got %q", want, got)
				}
			}
		})
	}
}

func TestFindRequirementsSection(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantContains []string
		wantNotEmpty bool
	}{
		{
			name: "standard spec",
			content: `# Spec Title

## Purpose
Purpose content.

## Requirements

### Requirement: First
First content.

### Requirement: Second
Second content.

## Notes
Notes content.`,
			wantContains: []string{"### Requirement: First", "### Requirement: Second"},
			wantNotEmpty: true,
		},
		{
			name: "no requirements section",
			content: `# Title

## Purpose
Content.`,
			wantNotEmpty: false,
		},
		{
			name: "requirements at end",
			content: `## Purpose
Content.

## Requirements
Req content.`,
			wantContains: []string{"Req content."},
			wantNotEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRequirementsSection(tt.content)
			if tt.wantNotEmpty && got == "" {
				t.Error("FindRequirementsSection() = empty, want non-empty")

				return
			}
			if !tt.wantNotEmpty && got != "" {
				t.Errorf("FindRequirementsSection() = %q, want empty", got)

				return
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("FindRequirementsSection() should contain %q, got %q", want, got)
				}
			}
		})
	}
}

func TestFindAllH3Requirements(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "multiple requirements",
			content: `## Requirements

### Requirement: First Req
Content.

### Requirement: Second Req
More content.

### Requirement: Third Req
Even more.`,
			want: []string{"First Req", "Second Req", "Third Req"},
		},
		{
			name:    "no requirements",
			content: "## Purpose\nContent",
			want:    nil,
		},
		{
			name: "single requirement",
			content: `### Requirement: Only One
Content here.`,
			want: []string{"Only One"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindAllH3Requirements(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("FindAllH3Requirements() len = %d, want %d", len(got), len(tt.want))

				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSplitIntoSections(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantSections []string
	}{
		{
			name: "multiple sections",
			content: `## Purpose
Purpose content.

## Requirements
Requirements content.

## Notes
Notes content.`,
			wantSections: []string{"Purpose", "Requirements", "Notes"},
		},
		{
			name:         "no sections",
			content:      "Just plain text",
			wantSections: nil,
		},
		{
			name: "single section",
			content: `## Requirements
Content.`,
			wantSections: []string{"Requirements"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitIntoSections(tt.content)
			if len(got) != len(tt.wantSections) {
				t.Errorf("SplitIntoSections() len = %d, want %d", len(got), len(tt.wantSections))

				return
			}
			for _, section := range tt.wantSections {
				if _, ok := got[section]; !ok {
					t.Errorf("SplitIntoSections() missing section %q", section)
				}
			}
		})
	}
}

func TestFindSectionContentEdgeCases(t *testing.T) {
	t.Run("section with regex special chars", func(t *testing.T) {
		content := `## API (v2.0) [Beta]
Content.

## Other`
		got := FindSectionContent(content, "API (v2.0) [Beta]")
		if !strings.Contains(got, "Content.") {
			t.Errorf("FindSectionContent() should handle regex special chars, got %q", got)
		}
	})

	t.Run("nested H3 headers within section", func(t *testing.T) {
		content := `## Requirements

### Requirement: One
Content one.

### Requirement: Two
Content two.

## Notes`
		got := FindSectionContent(content, "Requirements")
		if !strings.Contains(got, "### Requirement: One") {
			t.Error("FindSectionContent() should include nested H3 headers")
		}
		if !strings.Contains(got, "### Requirement: Two") {
			t.Error("FindSectionContent() should include all nested H3 headers")
		}
		if strings.Contains(got, "## Notes") {
			t.Error("FindSectionContent() should not include next H2 section")
		}
	})
}
