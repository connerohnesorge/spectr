package regex

import (
	"strings"
	"testing"
)

func TestFindSectionContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		sectionHeader string
		want          string
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
			want:          "\nThis is the requirements content.\n\n",
		},
		{
			name: "section at end",
			content: `## Purpose
Purpose content.

## Requirements
Final section content.`,
			sectionHeader: "Requirements",
			want:          "\nFinal section content.",
		},
		{
			name:          "section not found",
			content:       "## Purpose\nContent",
			sectionHeader: "Requirements",
			want:          "",
		},
		{
			name: "delta section",
			content: `## ADDED Requirements
Added content here.

## MODIFIED Requirements
Modified content here.`,
			sectionHeader: "ADDED Requirements",
			want:          "\nAdded content here.\n\n",
		},
		{
			name: "section with special chars in name",
			content: `## API & Integration
API content.

## Other`,
			sectionHeader: "API & Integration",
			want:          "\nAPI content.\n\n",
		},
		{
			name: "section with trailing space",
			content: `## Requirements
Content here.

## Next`,
			sectionHeader: "Requirements",
			want:          "\nContent here.\n\n",
		},
		{
			name:          "empty content",
			content:       "",
			sectionHeader: "Requirements",
			want:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSectionContent(tt.content, tt.sectionHeader)
			if got != tt.want {
				t.Errorf("FindSectionContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFindDeltaSectionContent(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		deltaType string
		want      string
	}{
		{
			name: "ADDED section",
			content: `## ADDED Requirements
Added content.

## MODIFIED Requirements
Modified content.`,
			deltaType: "ADDED",
			want:      "\nAdded content.\n\n",
		},
		{
			name: "MODIFIED section",
			content: `## ADDED Requirements
Added.

## MODIFIED Requirements
Modified content here.

## REMOVED Requirements
Removed.`,
			deltaType: "MODIFIED",
			want:      "\nModified content here.\n\n",
		},
		{
			name: "REMOVED section at end",
			content: `## ADDED Requirements
Added.

## REMOVED Requirements
Removed content.`,
			deltaType: "REMOVED",
			want:      "\nRemoved content.",
		},
		{
			name: "RENAMED section",
			content: `## RENAMED Requirements
- FROM: ### Requirement: Old
- TO: ### Requirement: New`,
			deltaType: "RENAMED",
			want:      "\n- FROM: ### Requirement: Old\n- TO: ### Requirement: New",
		},
		{
			name:      "missing section",
			content:   "## ADDED Requirements\nContent",
			deltaType: "MODIFIED",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDeltaSectionContent(tt.content, tt.deltaType)
			if got != tt.want {
				t.Errorf("FindDeltaSectionContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFindRequirementsSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
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
			want: "\n### Requirement: First\nFirst content.\n\n### Requirement: Second\nSecond content.\n\n",
		},
		{
			name: "no requirements section",
			content: `# Title

## Purpose
Content.`,
			want: "",
		},
		{
			name: "requirements at end",
			content: `## Purpose
Content.

## Requirements
Req content.`,
			want: "\nReq content.",
		},
		{
			name: "empty requirements",
			content: `## Requirements

## Notes`,
			want: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRequirementsSection(tt.content)
			if got != tt.want {
				t.Errorf("FindRequirementsSection() = %q, want %q", got, tt.want)
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
		{
			name: "requirements with special chars",
			content: `### Requirement: API-v2.0
### Requirement: Feature_Test
### Requirement: Name (v1)`,
			want: []string{"API-v2.0", "Feature_Test", "Name (v1)"},
		},
		{
			name:    "empty content",
			content: "",
			want:    nil,
		},
		{
			name: "scenarios not matched",
			content: `### Requirement: MyReq
#### Scenario: Test
Content`,
			want: []string{"MyReq"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindAllH3Requirements(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("FindAllH3Requirements() len = %d, want %d", len(got), len(tt.want))
				t.Errorf("got: %v, want: %v", got, tt.want)

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

func TestFindSectionIndex(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		sectionHeader string
		wantFound     bool
		wantStart     int
	}{
		{
			name:          "section found",
			content:       "## Requirements\nContent",
			sectionHeader: "Requirements",
			wantFound:     true,
			wantStart:     0,
		},
		{
			name:          "section in middle",
			content:       "## Purpose\n\n## Requirements\nContent",
			sectionHeader: "Requirements",
			wantFound:     true,
			wantStart:     12, // Position after "## Purpose\n\n"
		},
		{
			name:          "section not found",
			content:       "## Purpose\nContent",
			sectionHeader: "Requirements",
			wantFound:     false,
			wantStart:     -1,
		},
		{
			name:          "empty content",
			content:       "",
			sectionHeader: "Requirements",
			wantFound:     false,
			wantStart:     -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSectionIndex(tt.content, tt.sectionHeader)
			if tt.wantFound {
				if got == nil {
					t.Error("FindSectionIndex() = nil, want non-nil")

					return
				}
				if got[0] != tt.wantStart {
					t.Errorf("[0] = %d, want %d", got[0], tt.wantStart)
				}
			} else if got != nil {
				t.Errorf("FindSectionIndex() = %v, want nil", got)
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

	t.Run("section name with backslash", func(t *testing.T) {
		content := `## Path\To\Section
Content.

## Other`
		got := FindSectionContent(content, `Path\To\Section`)
		if !strings.Contains(got, "Content.") {
			t.Errorf("FindSectionContent() should handle backslashes, got %q", got)
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
