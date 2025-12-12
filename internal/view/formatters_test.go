package view

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestFormatDashboardText_FullDashboard tests the complete dashboard output
func TestFormatDashboardText_FullDashboard(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        5,
			TotalRequirements: 42,
			ActiveChanges:     2,
			CompletedChanges:  1,
			TotalTasks:        15,
			CompletedTasks:    8,
		},
		ActiveChanges: []ChangeProgress{
			{
				ID:    "add-view-command",
				Title: "Add View Command",
				Progress: ProgressMetrics{
					Total:      8,
					Completed:  3,
					Percentage: 37,
				},
			},
			{
				ID:    "add-validate-command",
				Title: "Add Validate Command",
				Progress: ProgressMetrics{
					Total:      5,
					Completed:  3,
					Percentage: 60,
				},
			},
		},
		CompletedChanges: []CompletedChange{
			{
				ID:    "add-list-command",
				Title: "Add List Command",
			},
		},
		Specs: []SpecInfo{
			{
				ID:               "cli-framework",
				Title:            "CLI Framework",
				RequirementCount: 15,
			},
			{
				ID:               "validation",
				Title:            "Validation",
				RequirementCount: 12,
			},
		},
	}

	output := FormatDashboardText(data)

	// Verify key elements are present
	expectedElements := []string{
		"Spectr Dashboard",
		"════════════════════════════════════════════════════════════",
		"Summary:",
		"5 specs, 42 requirements",
		"Active Changes: 2 in progress",
		"Completed Changes: 1",
		"Task Progress: 8/15 (53% complete)",
		"Active Changes",
		"────────────────────────────────────────────────────────────",
		"add-view-command",
		"add-validate-command",
		"Completed Changes",
		"add-list-command",
		"Specifications",
		"cli-framework",
		"validation",
		"15 requirements",
		"12 requirements",
		"Use spectr list --changes or spectr list --specs for detailed views",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf(
				"Expected output to contain %q, but it didn't.\nFull output:\n%s",
				expected,
				output,
			)
		}
	}
}

// TestFormatDashboardText_EmptyDashboard tests dashboard with no data
func TestFormatDashboardText_EmptyDashboard(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        0,
			TotalRequirements: 0,
			ActiveChanges:     0,
			CompletedChanges:  0,
			TotalTasks:        0,
			CompletedTasks:    0,
		},
		ActiveChanges: make(
			[]ChangeProgress,
			0,
		),
		CompletedChanges: make(
			[]CompletedChange,
			0,
		),
		Specs: make([]SpecInfo, 0),
	}

	output := FormatDashboardText(data)

	// Should still have header and footer
	if !strings.Contains(
		output,
		"Spectr Dashboard",
	) {
		t.Error("Expected dashboard title")
	}
	if !strings.Contains(
		output,
		"Use spectr list --changes or spectr list --specs",
	) {
		t.Error("Expected footer hint")
	}

	// Should NOT have section headers for empty sections
	if strings.Contains(
		output,
		"Active Changes\n────",
	) {
		t.Error(
			"Should not show Active Changes section when empty",
		)
	}
	if strings.Contains(
		output,
		"Completed Changes\n────",
	) {
		t.Error(
			"Should not show Completed Changes section when empty",
		)
	}
	if strings.Contains(
		output,
		"Specifications\n────",
	) {
		t.Error(
			"Should not show Specifications section when empty",
		)
	}

	// Should have summary with zero values
	if !strings.Contains(
		output,
		"0 specs, 0 requirements",
	) {
		t.Error(
			"Expected zero specs/requirements in summary",
		)
	}
}

// TestFormatDashboardText_OnlyActiveChanges tests dashboard with only active changes
func TestFormatDashboardText_OnlyActiveChanges(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        0,
			TotalRequirements: 0,
			ActiveChanges:     1,
			CompletedChanges:  0,
			TotalTasks:        10,
			CompletedTasks:    5,
		},
		ActiveChanges: []ChangeProgress{
			{
				ID:    "test-change",
				Title: "Test Change",
				Progress: ProgressMetrics{
					Total:      10,
					Completed:  5,
					Percentage: 50,
				},
			},
		},
		CompletedChanges: make(
			[]CompletedChange,
			0,
		),
		Specs: make([]SpecInfo, 0),
	}

	output := FormatDashboardText(data)

	// Should have active changes section
	if !strings.Contains(
		output,
		"Active Changes",
	) {
		t.Error("Expected Active Changes section")
	}
	if !strings.Contains(output, "test-change") {
		t.Error("Expected test-change in output")
	}

	// Should NOT have completed changes or specs sections
	if strings.Contains(
		output,
		"Completed Changes\n────",
	) {
		t.Error(
			"Should not show Completed Changes section when empty",
		)
	}
	if strings.Contains(
		output,
		"Specifications\n────",
	) {
		t.Error(
			"Should not show Specifications section when empty",
		)
	}
}

// TestFormatDashboardText_OnlyCompletedChanges tests dashboard with only completed changes
func TestFormatDashboardText_OnlyCompletedChanges(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        0,
			TotalRequirements: 0,
			ActiveChanges:     0,
			CompletedChanges:  2,
			TotalTasks:        0,
			CompletedTasks:    0,
		},
		ActiveChanges: make([]ChangeProgress, 0),
		CompletedChanges: []CompletedChange{
			{
				ID:    "change-one",
				Title: "Change One",
			},
			{
				ID:    "change-two",
				Title: "Change Two",
			},
		},
		Specs: make([]SpecInfo, 0),
	}

	output := FormatDashboardText(data)

	// Should have completed changes section
	if !strings.Contains(
		output,
		"Completed Changes",
	) {
		t.Error(
			"Expected Completed Changes section",
		)
	}
	if !strings.Contains(output, "change-one") {
		t.Error("Expected change-one in output")
	}
	if !strings.Contains(output, "change-two") {
		t.Error("Expected change-two in output")
	}

	// Should NOT have active changes or specs sections
	if strings.Contains(
		output,
		"Active Changes\n────",
	) {
		t.Error(
			"Should not show Active Changes section when empty",
		)
	}
	if strings.Contains(
		output,
		"Specifications\n────",
	) {
		t.Error(
			"Should not show Specifications section when empty",
		)
	}
}

// TestFormatDashboardText_OnlySpecs tests dashboard with only specifications
func TestFormatDashboardText_OnlySpecs(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        2,
			TotalRequirements: 20,
			ActiveChanges:     0,
			CompletedChanges:  0,
			TotalTasks:        0,
			CompletedTasks:    0,
		},
		ActiveChanges: make(
			[]ChangeProgress,
			0,
		),
		CompletedChanges: make(
			[]CompletedChange,
			0,
		),
		Specs: []SpecInfo{
			{
				ID:               "spec-one",
				Title:            "Spec One",
				RequirementCount: 10,
			},
			{
				ID:               "spec-two",
				Title:            "Spec Two",
				RequirementCount: 10,
			},
		},
	}

	output := FormatDashboardText(data)

	// Should have specs section
	if !strings.Contains(
		output,
		"Specifications",
	) {
		t.Error("Expected Specifications section")
	}
	if !strings.Contains(output, "spec-one") {
		t.Error("Expected spec-one in output")
	}
	if !strings.Contains(output, "spec-two") {
		t.Error("Expected spec-two in output")
	}
	if !strings.Contains(
		output,
		"10 requirements",
	) {
		t.Error(
			"Expected requirement count in output",
		)
	}

	// Should NOT have active or completed changes sections
	if strings.Contains(
		output,
		"Active Changes\n────",
	) {
		t.Error(
			"Should not show Active Changes section when empty",
		)
	}
	if strings.Contains(
		output,
		"Completed Changes\n────",
	) {
		t.Error(
			"Should not show Completed Changes section when empty",
		)
	}
}

// TestFormatSummarySection tests summary formatting
func TestFormatSummarySection(t *testing.T) {
	summary := SummaryMetrics{
		TotalSpecs:        3,
		TotalRequirements: 25,
		ActiveChanges:     2,
		CompletedChanges:  1,
		TotalTasks:        10,
		CompletedTasks:    5,
	}

	output := formatSummarySection(summary)

	expectedElements := []string{
		"Summary:",
		"3 specs, 25 requirements",
		"Active Changes: 2 in progress",
		"Completed Changes: 1",
		"Task Progress: 5/10 (50% complete)",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf(
				"Expected summary to contain %q, but it didn't.\nFull output:\n%s",
				expected,
				output,
			)
		}
	}
}

// TestFormatActiveChangesSection tests active changes formatting
func TestFormatActiveChangesSection(
	t *testing.T,
) {
	changes := []ChangeProgress{
		{
			ID:    "change-one",
			Title: "Change One",
			Progress: ProgressMetrics{
				Total:      10,
				Completed:  3,
				Percentage: 30,
			},
		},
		{
			ID:    "change-two",
			Title: "Change Two",
			Progress: ProgressMetrics{
				Total:      5,
				Completed:  5,
				Percentage: 100,
			},
		},
	}

	output := formatActiveChangesSection(changes)

	expectedElements := []string{
		"Active Changes",
		"────────────────────────────────────────────────────────────",
		"change-one",
		"change-two",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf(
				"Expected active changes to contain %q, but it didn't.\nFull output:\n%s",
				expected,
				output,
			)
		}
	}
}

// TestFormatCompletedChangesSection tests completed changes formatting
func TestFormatCompletedChangesSection(
	t *testing.T,
) {
	changes := []CompletedChange{
		{
			ID:    "completed-one",
			Title: "Completed One",
		},
		{
			ID:    "completed-two",
			Title: "Completed Two",
		},
	}

	output := formatCompletedChangesSection(
		changes,
	)

	expectedElements := []string{
		"Completed Changes",
		"────────────────────────────────────────────────────────────",
		"completed-one",
		"completed-two",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf(
				"Expected completed changes to contain %q, but it didn't.\nFull output:\n%s",
				expected,
				output,
			)
		}
	}
}

// TestFormatSpecsSection tests specifications formatting
func TestFormatSpecsSection(t *testing.T) {
	specs := []SpecInfo{
		{
			ID:               "spec-alpha",
			Title:            "Spec Alpha",
			RequirementCount: 15,
		},
		{
			ID:               "spec-beta",
			Title:            "Spec Beta",
			RequirementCount: 8,
		},
	}

	output := formatSpecsSection(specs)

	expectedElements := []string{
		"Specifications",
		"────────────────────────────────────────────────────────────",
		"spec-alpha",
		"spec-beta",
		"15 requirements",
		"8 requirements",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf(
				"Expected specs to contain %q, but it didn't.\nFull output:\n%s",
				expected,
				output,
			)
		}
	}
}

// TestFormatDashboardJSON_FullDashboard tests JSON formatting with all sections
func TestFormatDashboardJSON_FullDashboard(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        5,
			TotalRequirements: 42,
			ActiveChanges:     2,
			CompletedChanges:  1,
			TotalTasks:        15,
			CompletedTasks:    8,
		},
		ActiveChanges: []ChangeProgress{
			{
				ID:    "add-view-command",
				Title: "Add View Command",
				Progress: ProgressMetrics{
					Total:      8,
					Completed:  3,
					Percentage: 37,
				},
			},
			{
				ID:    "add-validate-command",
				Title: "Add Validate Command",
				Progress: ProgressMetrics{
					Total:      5,
					Completed:  3,
					Percentage: 60,
				},
			},
		},
		CompletedChanges: []CompletedChange{
			{
				ID:    "add-list-command",
				Title: "Add List Command",
			},
		},
		Specs: []SpecInfo{
			{
				ID:               "cli-framework",
				Title:            "CLI Framework",
				RequirementCount: 15,
			},
			{
				ID:               "validation",
				Title:            "Validation",
				RequirementCount: 12,
			},
		},
	}

	output, err := FormatDashboardJSON(data)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}

	// Verify it's valid JSON by checking for key elements
	expectedElements := []string{
		`"summary"`,
		`"totalSpecs": 5`,
		`"totalRequirements": 42`,
		`"activeChanges": 2`,
		`"completedChanges": 1`,
		`"totalTasks": 15`,
		`"completedTasks": 8`,
		`"activeChanges"`,
		`"id": "add-view-command"`,
		`"title": "Add View Command"`,
		`"progress"`,
		`"total": 8`,
		`"completed": 3`,
		`"percentage": 37`,
		`"id": "add-validate-command"`,
		`"completedChanges"`,
		`"id": "add-list-command"`,
		`"specs"`,
		`"id": "cli-framework"`,
		`"requirementCount": 15`,
		`"id": "validation"`,
		`"requirementCount": 12`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf(
				"Expected JSON to contain %q, but it didn't.\nFull output:\n%s",
				expected,
				output,
			)
		}
	}
}

// TestFormatDashboardJSON_EmptyDashboard tests JSON formatting with no data
func TestFormatDashboardJSON_EmptyDashboard(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        0,
			TotalRequirements: 0,
			ActiveChanges:     0,
			CompletedChanges:  0,
			TotalTasks:        0,
			CompletedTasks:    0,
		},
		ActiveChanges: make(
			[]ChangeProgress,
			0,
		),
		CompletedChanges: make(
			[]CompletedChange,
			0,
		),
		Specs: make([]SpecInfo, 0),
	}

	output, err := FormatDashboardJSON(data)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}

	// Verify empty arrays are represented as []
	expectedElements := []string{
		`"summary"`,
		`"totalSpecs": 0`,
		`"totalRequirements": 0`,
		`"activeChanges": []`,
		`"completedChanges": []`,
		`"specs": []`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf(
				"Expected JSON to contain %q, but it didn't.\nFull output:\n%s",
				expected,
				output,
			)
		}
	}
}

// TestFormatDashboardJSON_ValidJSON tests that output is valid JSON
func TestFormatDashboardJSON_ValidJSON(
	t *testing.T,
) {
	data := &DashboardData{
		Summary: SummaryMetrics{
			TotalSpecs:        1,
			TotalRequirements: 5,
			ActiveChanges:     1,
			CompletedChanges:  0,
			TotalTasks:        3,
			CompletedTasks:    1,
		},
		ActiveChanges: []ChangeProgress{
			{
				ID:    "test-change",
				Title: "Test Change",
				Progress: ProgressMetrics{
					Total:      3,
					Completed:  1,
					Percentage: 33,
				},
			},
		},
		CompletedChanges: make(
			[]CompletedChange,
			0,
		),
		Specs: []SpecInfo{
			{
				ID:               "test-spec",
				Title:            "Test Spec",
				RequirementCount: 5,
			},
		},
	}

	output, err := FormatDashboardJSON(data)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}

	// Unmarshal back to verify it's valid JSON
	var parsedData DashboardData
	if err := json.Unmarshal([]byte(output), &parsedData); err != nil {
		t.Fatalf(
			"Failed to parse output as JSON: %v\nOutput:\n%s",
			err,
			output,
		)
	}

	// Verify data integrity
	if parsedData.Summary.TotalSpecs != 1 {
		t.Errorf(
			"Expected TotalSpecs=1, got %d",
			parsedData.Summary.TotalSpecs,
		)
	}
	if len(parsedData.ActiveChanges) != 1 {
		t.Errorf(
			"Expected 1 active change, got %d",
			len(parsedData.ActiveChanges),
		)
	}
	if parsedData.ActiveChanges[0].ID != "test-change" {
		t.Errorf(
			"Expected ID='test-change', got '%s'",
			parsedData.ActiveChanges[0].ID,
		)
	}
}
