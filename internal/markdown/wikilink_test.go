package markdown

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestProject creates a temporary project structure for testing wikilink resolution.
func setupTestProject(t *testing.T) string {
	t.Helper()

	// Create temp directory
	tempDir, err := os.MkdirTemp(
		"",
		"wikilink-test-*",
	)
	if err != nil {
		t.Fatalf(
			"failed to create temp dir: %v",
			err,
		)
	}

	// Create spectr directory structure
	spectrDir := filepath.Join(tempDir, "spectr")
	specsDir := filepath.Join(spectrDir, "specs")
	changesDir := filepath.Join(
		spectrDir,
		"changes",
	)

	// Create directories
	dirs := []string{
		filepath.Join(specsDir, "validation"),
		filepath.Join(specsDir, "cli-interface"),
		filepath.Join(
			specsDir,
			"naming-conventions",
		),
		filepath.Join(changesDir, "my-change"),
		filepath.Join(
			changesDir,
			"another-change",
		),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf(
				"failed to create dir %s: %v",
				dir,
				err,
			)
		}
	}

	// Create spec files with content
	validationSpec := `# Validation Spec

## Requirements

### Requirement: Spec File Validation

The system must validate spec files.

#### Scenario: Valid spec file

- **WHEN** a valid spec file is provided
- **THEN** validation passes

### Requirement: Another Requirement

Another requirement description.

## Testing Section

Some content here.
`

	cliSpec := `# CLI Interface

## Commands

### Requirement: Init Command

The init command creates a new spec.
`

	namingSpec := `# Naming Conventions

## Overview

Guidelines for naming things.

### Requirement: Lowercase Names

Names should be lowercase.
`

	// Create spec files
	specFiles := map[string]string{
		filepath.Join(specsDir, "validation", "spec.md"):         validationSpec,
		filepath.Join(specsDir, "cli-interface", "spec.md"):      cliSpec,
		filepath.Join(specsDir, "naming-conventions", "spec.md"): namingSpec,
	}

	for path, content := range specFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf(
				"failed to write %s: %v",
				path,
				err,
			)
		}
	}

	// Create proposal files
	myChangeProposal := `# My Change

## Summary

This is a change proposal.

### Requirement: New Feature

A new feature requirement.

#### Scenario: Feature works

- **WHEN** feature is used
- **THEN** it works
`

	anotherChangeProposal := `# Another Change

## Summary

Another change proposal.
`

	proposalFiles := map[string]string{
		filepath.Join(changesDir, "my-change", "proposal.md"):      myChangeProposal,
		filepath.Join(changesDir, "another-change", "proposal.md"): anotherChangeProposal,
	}

	for path, content := range proposalFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf(
				"failed to write %s: %v",
				path,
				err,
			)
		}
	}

	return tempDir
}

func TestResolveWikilink(t *testing.T) {
	projectRoot := setupTestProject(t)
	defer func() { _ = os.RemoveAll(projectRoot) }()

	tests := []struct {
		name       string
		target     string
		wantExists bool
		wantPath   string // relative to projectRoot
	}{
		{
			name:       "spec target - validation",
			target:     "validation",
			wantExists: true,
			wantPath:   "spectr/specs/validation/spec.md",
		},
		{
			name:       "spec target - cli-interface",
			target:     "cli-interface",
			wantExists: true,
			wantPath:   "spectr/specs/cli-interface/spec.md",
		},
		{
			name:       "spec target - naming-conventions",
			target:     "naming-conventions",
			wantExists: true,
			wantPath:   "spectr/specs/naming-conventions/spec.md",
		},
		{
			name:       "change target - explicit prefix",
			target:     "changes/my-change",
			wantExists: true,
			wantPath:   "spectr/changes/my-change/proposal.md",
		},
		{
			name:       "change target - another-change",
			target:     "changes/another-change",
			wantExists: true,
			wantPath:   "spectr/changes/another-change/proposal.md",
		},
		{
			name:       "nonexistent spec",
			target:     "nonexistent",
			wantExists: false,
			wantPath:   "spectr/specs/nonexistent/spec.md",
		},
		{
			name:       "nonexistent change",
			target:     "changes/nonexistent",
			wantExists: false,
			wantPath:   "spectr/changes/nonexistent/proposal.md",
		},
		{
			name:       "empty target",
			target:     "",
			wantExists: false,
			wantPath:   "",
		},
		{
			name:       "target with anchor stripped",
			target:     "validation#Requirement: Spec File Validation",
			wantExists: true,
			wantPath:   "spectr/specs/validation/spec.md",
		},
		{
			name:       "explicit specs prefix",
			target:     "specs/validation",
			wantExists: true,
			wantPath:   "spectr/specs/validation/spec.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotExists := ResolveWikilink(
				tt.target,
				projectRoot,
			)

			if gotExists != tt.wantExists {
				t.Errorf(
					"ResolveWikilink(%q) exists = %v, want %v",
					tt.target,
					gotExists,
					tt.wantExists,
				)
			}

			if tt.wantPath == "" {
				return
			}

			wantFullPath := filepath.Join(
				projectRoot,
				tt.wantPath,
			)
			if gotPath != wantFullPath {
				t.Errorf(
					"ResolveWikilink(%q) path = %q, want %q",
					tt.target,
					gotPath,
					wantFullPath,
				)
			}
		})
	}
}

func TestResolveWikilinkWithAnchor(t *testing.T) {
	projectRoot := setupTestProject(t)
	defer func() { _ = os.RemoveAll(projectRoot) }()

	tests := []struct {
		name         string
		target       string
		anchor       string
		wantAnchorOK bool
		wantErr      bool
	}{
		{
			name:         "valid requirement anchor",
			target:       "validation",
			anchor:       "Requirement: Spec File Validation",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "valid requirement anchor - case insensitive",
			target:       "validation",
			anchor:       "requirement: spec file validation",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "valid scenario anchor",
			target:       "validation",
			anchor:       "Scenario: Valid spec file",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "valid section anchor",
			target:       "validation",
			anchor:       "Testing Section",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "invalid anchor",
			target:       "validation",
			anchor:       "Nonexistent Section",
			wantAnchorOK: false,
			wantErr:      false,
		},
		{
			name:         "no anchor - should be valid",
			target:       "validation",
			anchor:       "",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "nonexistent target",
			target:       "nonexistent",
			anchor:       "Some Anchor",
			wantAnchorOK: false,
			wantErr:      false,
		},
		{
			name:         "requirement in change proposal",
			target:       "changes/my-change",
			anchor:       "Requirement: New Feature",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "scenario in change proposal",
			target:       "changes/my-change",
			anchor:       "Scenario: Feature works",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "partial match - Requirements section",
			target:       "validation",
			anchor:       "Requirements",
			wantAnchorOK: true,
			wantErr:      false,
		},
		{
			name:         "another requirement name",
			target:       "validation",
			anchor:       "Requirement: Another Requirement",
			wantAnchorOK: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotAnchorOK, gotErr := ResolveWikilinkWithAnchor(
				tt.target,
				tt.anchor,
				projectRoot,
			)

			if (gotErr != nil) != tt.wantErr {
				t.Errorf(
					"ResolveWikilinkWithAnchor(%q, %q) error = %v, wantErr %v",
					tt.target,
					tt.anchor,
					gotErr,
					tt.wantErr,
				)

				return
			}

			if gotAnchorOK != tt.wantAnchorOK {
				t.Errorf(
					"ResolveWikilinkWithAnchor(%q, %q) anchorValid = %v, want %v",
					tt.target,
					tt.anchor,
					gotAnchorOK,
					tt.wantAnchorOK,
				)
			}
		})
	}
}

func TestValidateWikilinks(t *testing.T) {
	projectRoot := setupTestProject(t)
	defer func() { _ = os.RemoveAll(projectRoot) }()

	tests := []struct {
		name       string
		content    string
		wantErrors int
	}{
		{
			name:       "valid wikilink to spec",
			content:    "See [[validation]] for details.",
			wantErrors: 0,
		},
		{
			name:       "valid wikilink with display text",
			content:    "See [[validation|the validation spec]] for details.",
			wantErrors: 0,
		},
		{
			name:       "valid wikilink with anchor",
			content:    "See [[validation#Requirement: Spec File Validation]].",
			wantErrors: 0,
		},
		{
			name:       "valid wikilink to change",
			content:    "See [[changes/my-change]] for the proposal.",
			wantErrors: 0,
		},
		{
			name:       "invalid wikilink - nonexistent target",
			content:    "See [[nonexistent]] for details.",
			wantErrors: 1,
		},
		{
			name:       "invalid wikilink - nonexistent anchor",
			content:    "See [[validation#Nonexistent Anchor]].",
			wantErrors: 1,
		},
		{
			name:       "multiple valid wikilinks",
			content:    "See [[validation]] and [[cli-interface]].",
			wantErrors: 0,
		},
		{
			name:       "mix of valid and invalid wikilinks",
			content:    "See [[validation]] and [[nonexistent]].",
			wantErrors: 1,
		},
		{
			name:       "multiple invalid wikilinks",
			content:    "See [[nonexistent1]] and [[nonexistent2]].",
			wantErrors: 2,
		},
		{
			name:       "no wikilinks",
			content:    "This is plain text with no wikilinks.",
			wantErrors: 0,
		},
		{
			name:       "valid wikilink with valid anchor in change",
			content:    "See [[changes/my-change#Requirement: New Feature]].",
			wantErrors: 0,
		},
		{
			name:       "valid wikilink with invalid anchor in change",
			content:    "See [[changes/my-change#Requirement: Invalid]].",
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, _ := Parse([]byte(tt.content))
			errors := ValidateWikilinks(
				root,
				[]byte(tt.content),
				projectRoot,
			)

			if len(errors) == tt.wantErrors {
				return
			}

			t.Errorf(
				"ValidateWikilinks() got %d errors, want %d",
				len(errors),
				tt.wantErrors,
			)
			for i, err := range errors {
				t.Logf(
					"  error %d: %s",
					i,
					err.Message,
				)
			}
		})
	}
}

func TestValidateWikilinkTarget(t *testing.T) {
	projectRoot := setupTestProject(t)
	defer func() { _ = os.RemoveAll(projectRoot) }()

	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{
			name:    "valid spec target",
			target:  "validation",
			wantErr: false,
		},
		{
			name:    "valid change target",
			target:  "changes/my-change",
			wantErr: false,
		},
		{
			name:    "invalid target",
			target:  "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWikilinkTarget(
				tt.target,
				projectRoot,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ValidateWikilinkTarget(%q) error = %v, wantErr %v",
					tt.target,
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestGetWikilinkTargetType(t *testing.T) {
	projectRoot := setupTestProject(t)
	defer func() { _ = os.RemoveAll(projectRoot) }()

	tests := []struct {
		name   string
		target string
		want   string
	}{
		{
			name:   "spec target",
			target: "validation",
			want:   "spec",
		},
		{
			name:   "change target",
			target: "changes/my-change",
			want:   "change",
		},
		{
			name:   "nonexistent target",
			target: "nonexistent",
			want:   "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetWikilinkTargetType(
				tt.target,
				projectRoot,
			)
			if got != tt.want {
				t.Errorf(
					"GetWikilinkTargetType(%q) = %q, want %q",
					tt.target,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestListWikilinkTargets(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "single wikilink",
			content: "See [[validation]].",
			want:    []string{"validation"},
		},
		{
			name:    "multiple unique wikilinks",
			content: "See [[validation]] and [[cli-interface]].",
			want: []string{
				"validation",
				"cli-interface",
			},
		},
		{
			name:    "duplicate wikilinks deduplicated",
			content: "See [[validation]] and [[validation]] again.",
			want:    []string{"validation"},
		},
		{
			name:    "wikilinks with anchors",
			content: "See [[validation#anchor1]] and [[validation#anchor2]].",
			want:    []string{"validation"},
		},
		{
			name:    "no wikilinks",
			content: "Plain text.",
			want:    make([]string, 0),
		},
		{
			name:    "wikilinks with display text",
			content: "See [[validation|the validation spec]].",
			want:    []string{"validation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListWikilinkTargets(
				[]byte(tt.content),
			)

			if len(got) != len(tt.want) {
				t.Errorf(
					"ListWikilinkTargets() got %d targets, want %d",
					len(got),
					len(tt.want),
				)
				t.Logf("got: %v", got)
				t.Logf("want: %v", tt.want)

				return
			}

			// Check that all expected targets are present
			gotSet := make(map[string]bool)
			for _, target := range got {
				gotSet[target] = true
			}

			for _, wantTarget := range tt.want {
				if !gotSet[wantTarget] {
					t.Errorf(
						"ListWikilinkTargets() missing expected target %q",
						wantTarget,
					)
				}
			}
		})
	}
}

func TestWikilinkError(t *testing.T) {
	tests := []struct {
		name string
		err  WikilinkError
		want string
	}{
		{
			name: "error with offset",
			err: WikilinkError{
				Target:  "validation",
				Offset:  42,
				Message: "target not found",
			},
			want: "offset 42: target not found",
		},
		{
			name: "error without offset",
			err: WikilinkError{
				Target:  "validation",
				Offset:  -1,
				Message: "target not found",
			},
			want: "target not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf(
					"WikilinkError.Error() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestAnchorExistsInContent(t *testing.T) {
	content := []byte(`# Document

## First Section

### Requirement: Test Requirement

Description here.

#### Scenario: Happy Path

- **WHEN** something happens
- **THEN** it works

## Second Section

More content.
`)

	tests := []struct {
		name   string
		anchor string
		want   bool
	}{
		{
			name:   "exact section match",
			anchor: "First Section",
			want:   true,
		},
		{
			name:   "requirement match",
			anchor: "Requirement: Test Requirement",
			want:   true,
		},
		{
			name:   "scenario match",
			anchor: "Scenario: Happy Path",
			want:   true,
		},
		{
			name:   "case insensitive match",
			anchor: "FIRST SECTION",
			want:   true,
		},
		{
			name:   "partial match in section",
			anchor: "First",
			want:   true,
		},
		{
			name:   "nonexistent anchor",
			anchor: "Nonexistent",
			want:   false,
		},
		{
			name:   "document title",
			anchor: "Document",
			want:   true, // H1 is treated as a section, so this matches
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := anchorExistsInContent(
				content,
				tt.anchor,
			)
			if got != tt.want {
				t.Errorf(
					"anchorExistsInContent(%q) = %v, want %v",
					tt.anchor,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestValidateWikilinkWithNilRoot(
	t *testing.T,
) {
	errors := ValidateWikilinks(nil, nil, "/tmp")
	if errors != nil {
		t.Errorf(
			"ValidateWikilinks(nil) should return nil, got %v",
			errors,
		)
	}
}

func TestResolveWikilinkResolutionOrder(
	t *testing.T,
) {
	// Test that specs are checked before changes when target doesn't have prefix
	tempDir, err := os.MkdirTemp(
		"",
		"wikilink-order-test-*",
	)
	if err != nil {
		t.Fatalf(
			"failed to create temp dir: %v",
			err,
		)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create both a spec and a change with the same name
	specDir := filepath.Join(
		tempDir,
		"spectr",
		"specs",
		"duplicate",
	)
	changeDir := filepath.Join(
		tempDir,
		"spectr",
		"changes",
		"duplicate",
	)

	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf(
			"failed to create spec dir: %v",
			err,
		)
	}
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf(
			"failed to create change dir: %v",
			err,
		)
	}

	// Create both files
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte("# Spec"), 0644); err != nil {
		t.Fatalf("failed to write spec: %v", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte("# Proposal"), 0644); err != nil {
		t.Fatalf(
			"failed to write proposal: %v",
			err,
		)
	}

	// Resolve without prefix - should get spec
	path, exists := ResolveWikilink(
		"duplicate",
		tempDir,
	)
	if !exists {
		t.Error(
			"ResolveWikilink('duplicate') should exist",
		)
	}

	expectedPath := filepath.Join(
		tempDir,
		"spectr",
		"specs",
		"duplicate",
		"spec.md",
	)
	if path != expectedPath {
		t.Errorf(
			"ResolveWikilink('duplicate') = %q, want %q (spec should take priority)",
			path,
			expectedPath,
		)
	}

	// Explicit change prefix should get change
	changePath, changeExists := ResolveWikilink(
		"changes/duplicate",
		tempDir,
	)
	if !changeExists {
		t.Error(
			"ResolveWikilink('changes/duplicate') should exist",
		)
	}

	expectedChangePath := filepath.Join(
		tempDir,
		"spectr",
		"changes",
		"duplicate",
		"proposal.md",
	)
	if changePath != expectedChangePath {
		t.Errorf(
			"ResolveWikilink('changes/duplicate') = %q, want %q",
			changePath,
			expectedChangePath,
		)
	}
}
