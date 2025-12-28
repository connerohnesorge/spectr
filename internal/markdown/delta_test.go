package markdown

import (
	"testing"
)

func TestParseDelta_AddedRequirements(
	t *testing.T,
) {
	content := []byte(`# Change Proposal

## ADDED Requirements

### Requirement: New Feature

This is a new requirement being added.

#### Scenario: Happy path

- **WHEN** user performs action
- **THEN** expected result happens
`)

	delta, errors := ParseDelta(content)

	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	if len(delta.Added) != 1 {
		t.Fatalf(
			"expected 1 added requirement, got %d",
			len(delta.Added),
		)
	}

	if delta.Added[0].Name != "New Feature" {
		t.Errorf(
			"expected name 'New Feature', got '%s'",
			delta.Added[0].Name,
		)
	}

	if len(delta.Added[0].Scenarios) != 1 {
		t.Errorf(
			"expected 1 scenario, got %d",
			len(delta.Added[0].Scenarios),
		)
	}
}

func TestParseDelta_ModifiedRequirements(
	t *testing.T,
) {
	content := []byte(`# Change Proposal

## MODIFIED Requirements

### Requirement: Existing Feature

This requirement is being modified.

#### Scenario: Updated scenario

- **WHEN** modified action
- **THEN** modified result
`)

	delta, errors := ParseDelta(content)

	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	if len(delta.Modified) != 1 {
		t.Fatalf(
			"expected 1 modified requirement, got %d",
			len(delta.Modified),
		)
	}

	if delta.Modified[0].Name != "Existing Feature" {
		t.Errorf(
			"expected name 'Existing Feature', got '%s'",
			delta.Modified[0].Name,
		)
	}
}

func TestParseDelta_RemovedRequirements(
	t *testing.T,
) {
	content := []byte(`# Change Proposal

## REMOVED Requirements

- Deprecated Feature
- Old Functionality
- Legacy Component
`)

	delta, errors := ParseDelta(content)

	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	if len(delta.Removed) != 3 {
		t.Fatalf(
			"expected 3 removed requirements, got %d",
			len(delta.Removed),
		)
	}

	expected := []string{
		"Deprecated Feature",
		"Old Functionality",
		"Legacy Component",
	}
	for i, name := range expected {
		if delta.Removed[i] != name {
			t.Errorf(
				"removed[%d]: expected '%s', got '%s'",
				i,
				name,
				delta.Removed[i],
			)
		}
	}
}

func TestParseDelta_RenamedRequirements_ListFormat(
	t *testing.T,
) {
	content := []byte(`# Change Proposal

## RENAMED Requirements

- FROM: Old Name TO: New Name
- FROM: Previous Feature TO: Current Feature
`)

	delta, errors := ParseDelta(content)

	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	if len(delta.Renamed) != 2 {
		t.Fatalf(
			"expected 2 renamed requirements, got %d",
			len(delta.Renamed),
		)
	}

	if delta.Renamed[0].From != "Old Name" {
		t.Errorf(
			"renamed[0].From: expected 'Old Name', got '%s'",
			delta.Renamed[0].From,
		)
	}
	if delta.Renamed[0].To != "New Name" {
		t.Errorf(
			"renamed[0].To: expected 'New Name', got '%s'",
			delta.Renamed[0].To,
		)
	}

	if delta.Renamed[1].From != "Previous Feature" {
		t.Errorf(
			"renamed[1].From: expected 'Previous Feature', got '%s'",
			delta.Renamed[1].From,
		)
	}
	if delta.Renamed[1].To != "Current Feature" {
		t.Errorf(
			"renamed[1].To: expected 'Current Feature', got '%s'",
			delta.Renamed[1].To,
		)
	}
}

func TestParseDelta_RenamedRequirements_HeaderFormat(
	t *testing.T,
) {
	content := []byte(`# Change Proposal

## RENAMED Requirements

### Requirement: New Feature Name

FROM: Old Feature Name

This requirement was renamed.
`)

	delta, errors := ParseDelta(content)

	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	if len(delta.Renamed) != 1 {
		t.Fatalf(
			"expected 1 renamed requirement, got %d",
			len(delta.Renamed),
		)
	}

	if delta.Renamed[0].To != "New Feature Name" {
		t.Errorf(
			"renamed.To: expected 'New Feature Name', got '%s'",
			delta.Renamed[0].To,
		)
	}

	if delta.Renamed[0].From != "Old Feature Name" {
		t.Errorf(
			"renamed.From: expected 'Old Feature Name', got '%s'",
			delta.Renamed[0].From,
		)
	}
}

func TestParseDelta_MultipleSections(
	t *testing.T,
) {
	content := []byte(`# Change Proposal

## ADDED Requirements

### Requirement: Brand New

New requirement.

## MODIFIED Requirements

### Requirement: Updated

Modified requirement.

## REMOVED Requirements

- Obsolete Feature

## RENAMED Requirements

- FROM: OldName TO: NewName
`)

	delta, errors := ParseDelta(content)

	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	if len(delta.Added) != 1 {
		t.Errorf(
			"expected 1 added, got %d",
			len(delta.Added),
		)
	}
	if len(delta.Modified) != 1 {
		t.Errorf(
			"expected 1 modified, got %d",
			len(delta.Modified),
		)
	}
	if len(delta.Removed) != 1 {
		t.Errorf(
			"expected 1 removed, got %d",
			len(delta.Removed),
		)
	}
	if len(delta.Renamed) != 1 {
		t.Errorf(
			"expected 1 renamed, got %d",
			len(delta.Renamed),
		)
	}
}

func TestFindDeltaSection(t *testing.T) {
	content := []byte(`# Change Proposal

## ADDED Requirements

Content of added section.

## MODIFIED Requirements

Content of modified section.
`)

	tests := []struct {
		deltaType DeltaType
		wantEmpty bool
	}{
		{DeltaAdded, false},
		{DeltaModified, false},
		{DeltaRemoved, true},
		{DeltaRenamed, true},
	}

	for _, tt := range tests {
		result := FindDeltaSection(
			content,
			tt.deltaType,
		)
		gotEmpty := result == ""

		if gotEmpty != tt.wantEmpty {
			t.Errorf(
				"FindDeltaSection(%s): wantEmpty=%v, gotEmpty=%v",
				tt.deltaType,
				tt.wantEmpty,
				gotEmpty,
			)
		}

		if gotEmpty ||
			tt.deltaType != DeltaAdded {
			continue
		}

		if result == "" {
			t.Error(
				"expected non-empty content for ADDED section",
			)
		}
	}
}

func TestFindAllDeltaSections(t *testing.T) {
	content := []byte(`# Change

## ADDED Requirements
Added content.

## MODIFIED Requirements
Modified content.

## REMOVED Requirements
- Item
`)

	sections := FindAllDeltaSections(content)

	if len(sections) != 3 {
		t.Fatalf(
			"expected 3 sections, got %d",
			len(sections),
		)
	}

	if _, ok := sections[DeltaAdded]; !ok {
		t.Error("missing ADDED section")
	}
	if _, ok := sections[DeltaModified]; !ok {
		t.Error("missing MODIFIED section")
	}
	if _, ok := sections[DeltaRemoved]; !ok {
		t.Error("missing REMOVED section")
	}
}

func TestParseRenamedListItem(t *testing.T) {
	tests := []struct {
		input    string
		wantFrom string
		wantTo   string
	}{
		{
			"- FROM: OldName TO: NewName",
			"OldName",
			"NewName",
		},
		{
			"- FROM: Old Name TO: New Name",
			"Old Name",
			"New Name",
		},
		{"* from: old to: new", "old", "new"},
		{
			"- FROM OldName TO NewName",
			"OldName",
			"NewName",
		},
		{"- Just some text", "", ""},
		{"", "", ""},
		{"- FROM: OnlyFrom", "", ""},
		{"- TO: OnlyTo", "", ""},
	}

	for _, tt := range tests {
		from, to := parseRenamedListItem(tt.input)
		if from != tt.wantFrom {
			t.Errorf(
				"parseRenamedListItem(%q): from=%q, want %q",
				tt.input,
				from,
				tt.wantFrom,
			)
		}
		if to != tt.wantTo {
			t.Errorf(
				"parseRenamedListItem(%q): to=%q, want %q",
				tt.input,
				to,
				tt.wantTo,
			)
		}
	}
}

func TestParseFromAnnotation(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"FROM: OldName", "OldName"},
		{"FROM: Old Name", "Old Name"},
		{"from: lower case", "lower case"},
		{
			"Some text FROM: Name Here",
			"Name Here",
		},
		{"FROM: Name.\n", "Name"},
		{"No from here", ""},
		{"", ""},
	}

	for _, tt := range tests {
		got := parseFromAnnotation(tt.input)
		if got != tt.want {
			t.Errorf(
				"parseFromAnnotation(%q) = %q, want %q",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestExtractListItemText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"- item text", "item text"},
		{"* starred item", "starred item"},
		{"+ plus item", "plus item"},
		{"1. numbered item", "numbered item"},
		{"  - indented item", "indented item"},
		{"- [ ] unchecked", "unchecked"},
		{"- [x] checked", "checked"},
		{"plain text", "plain text"},
	}

	for _, tt := range tests {
		got := extractListItemText(tt.input)
		if got != tt.want {
			t.Errorf(
				"extractListItemText(%q) = %q, want %q",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestGetDeltaRequirementNames(t *testing.T) {
	delta := &Delta{
		Added: []*Requirement{
			{Name: "AddedOne"},
			{Name: "AddedTwo"},
		},
		Modified: []*Requirement{
			{Name: "ModifiedOne"},
		},
		Removed: []string{
			"RemovedOne",
			"RemovedTwo",
		},
		Renamed: []*RenamedRequirement{
			{From: "OldName", To: "NewName"},
		},
	}

	names := GetDeltaRequirementNames(delta)

	expected := []string{
		"AddedOne",
		"AddedTwo",
		"ModifiedOne",
		"RemovedOne",
		"RemovedTwo",
		"OldName",
		"NewName",
	}

	if len(names) != len(expected) {
		t.Fatalf(
			"expected %d names, got %d",
			len(expected),
			len(names),
		)
	}

	// Create a set for comparison
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	for _, exp := range expected {
		if !nameSet[exp] {
			t.Errorf(
				"missing expected name: %s",
				exp,
			)
		}
	}
}

func TestHasDeltaSection(t *testing.T) {
	content := []byte(`# Change
## ADDED Requirements
Content
`)

	if !HasDeltaSection(content, DeltaAdded) {
		t.Error("expected to find ADDED section")
	}

	if HasDeltaSection(content, DeltaModified) {
		t.Error(
			"should not find MODIFIED section",
		)
	}
}

func TestCountDeltaChanges(t *testing.T) {
	delta := &Delta{
		Added: []*Requirement{
			{Name: "A"},
			{Name: "B"},
		},
		Modified: []*Requirement{{Name: "C"}},
		Removed:  []string{"D"},
		Renamed: []*RenamedRequirement{
			{From: "E", To: "F"},
		},
	}

	count := CountDeltaChanges(delta)
	if count != 5 {
		t.Errorf(
			"expected 5 changes, got %d",
			count,
		)
	}
}

func TestValidateRenamed(t *testing.T) {
	delta := &Delta{
		Renamed: []*RenamedRequirement{
			{From: "Old", To: "New"}, // Valid
			{
				From: "",
				To:   "OnlyTo",
			}, // Missing From
			{
				From: "OnlyFrom",
				To:   "",
			}, // Missing To
		},
	}

	incomplete := ValidateRenamed(delta)

	if len(incomplete) != 2 {
		t.Fatalf(
			"expected 2 incomplete, got %d",
			len(incomplete),
		)
	}
}

func TestMergeDelta(t *testing.T) {
	merged := &Delta{
		Added: []*Requirement{{Name: "A"}},
	}

	delta := &Delta{
		Added:    []*Requirement{{Name: "B"}},
		Modified: []*Requirement{{Name: "C"}},
		Removed:  []string{"D"},
	}

	MergeDelta(merged, delta)

	if len(merged.Added) != 2 {
		t.Errorf(
			"expected 2 added, got %d",
			len(merged.Added),
		)
	}
	if len(merged.Modified) != 1 {
		t.Errorf(
			"expected 1 modified, got %d",
			len(merged.Modified),
		)
	}
	if len(merged.Removed) != 1 {
		t.Errorf(
			"expected 1 removed, got %d",
			len(merged.Removed),
		)
	}
}

func TestFindRenamedPairs(t *testing.T) {
	content := []byte(`# Change
## RENAMED Requirements
- FROM: Old1 TO: New1
- FROM: Old2 TO: New2
`)

	pairs := FindRenamedPairs(content)

	if len(pairs) != 2 {
		t.Fatalf(
			"expected 2 pairs, got %d",
			len(pairs),
		)
	}

	if pairs["Old1"] != "New1" {
		t.Errorf(
			"expected Old1 -> New1, got Old1 -> %s",
			pairs["Old1"],
		)
	}

	if pairs["Old2"] != "New2" {
		t.Errorf(
			"expected Old2 -> New2, got Old2 -> %s",
			pairs["Old2"],
		)
	}
}

func TestFindAddedRequirements(t *testing.T) {
	content := []byte(`# Change
## ADDED Requirements
### Requirement: Feature A
Description.
### Requirement: Feature B
Description.
`)

	names := FindAddedRequirements(content)

	if len(names) != 2 {
		t.Fatalf(
			"expected 2 names, got %d",
			len(names),
		)
	}

	if names[0] != "Feature A" {
		t.Errorf(
			"expected 'Feature A', got '%s'",
			names[0],
		)
	}
	if names[1] != "Feature B" {
		t.Errorf(
			"expected 'Feature B', got '%s'",
			names[1],
		)
	}
}

func TestGetDeltaSummary(t *testing.T) {
	delta := &Delta{
		Added: []*Requirement{{Name: "New"}},
		Modified: []*Requirement{
			{Name: "Updated"},
		},
		Removed: []string{"Old"},
		Renamed: []*RenamedRequirement{
			{From: "Before", To: "After"},
		},
	}

	summary := GetDeltaSummary(delta)

	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Check that all sections are mentioned
	if !containsSubstring(summary, "Added:") {
		t.Error("summary missing Added section")
	}
	if !containsSubstring(summary, "Modified:") {
		t.Error(
			"summary missing Modified section",
		)
	}
	if !containsSubstring(summary, "Removed:") {
		t.Error("summary missing Removed section")
	}
	if !containsSubstring(summary, "Renamed:") {
		t.Error("summary missing Renamed section")
	}
}

func TestParseDelta_EmptyContent(t *testing.T) {
	delta, errors := ParseDelta(make([]byte, 0))

	if delta == nil {
		t.Fatal("expected non-nil delta")
	}

	if len(errors) > 0 {
		t.Errorf(
			"unexpected errors for empty content: %v",
			errors,
		)
	}

	if len(delta.Added) != 0 ||
		len(delta.Modified) != 0 ||
		len(delta.Removed) != 0 ||
		len(delta.Renamed) != 0 {
		t.Error(
			"expected empty delta for empty content",
		)
	}
}

func TestParseDelta_NoDeltaSections(
	t *testing.T,
) {
	content := []byte(`# Regular Document

## Introduction

This is not a delta file.

### Requirement: Some Requirement

Regular requirement.
`)

	delta, errors := ParseDelta(content)

	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	// Requirements outside delta sections should not be categorized
	if len(delta.Added) != 0 {
		t.Errorf(
			"expected 0 added, got %d",
			len(delta.Added),
		)
	}
	if len(delta.Modified) != 0 {
		t.Errorf(
			"expected 0 modified, got %d",
			len(delta.Modified),
		)
	}
}

func TestIsDeltaSectionHeader(t *testing.T) {
	tests := []struct {
		title    string
		wantType DeltaType
		wantOk   bool
	}{
		{"ADDED Requirements", DeltaAdded, true},
		{"Added Requirements", DeltaAdded, true},
		{"added requirements", DeltaAdded, true},
		{
			"MODIFIED Requirements",
			DeltaModified,
			true,
		},
		{
			"REMOVED Requirements",
			DeltaRemoved,
			true,
		},
		{
			"RENAMED Requirements",
			DeltaRenamed,
			true,
		},
		{"Regular Section", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		gotType, gotOk := isDeltaSectionHeader(
			tt.title,
		)
		if gotType != tt.wantType ||
			gotOk != tt.wantOk {
			t.Errorf(
				"isDeltaSectionHeader(%q) = (%q, %v), want (%q, %v)",
				tt.title,
				gotType,
				gotOk,
				tt.wantType,
				tt.wantOk,
			)
		}
	}
}

// Helper function to check if string contains substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || s != "" && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(
	s, substr string,
) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
