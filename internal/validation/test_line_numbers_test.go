package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLineNumbersAreTracked(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system provides user authentication.

##### Scenario: Wrong number of hashtags
- **WHEN** user logs in
- **THEN** user is authenticated
`,
	}

	// Create test structure
	projectRoot := t.TempDir()
	spectrRoot := filepath.Join(projectRoot, "spectr")
	changesRoot := filepath.Join(spectrRoot, "changes")
	changeDir := filepath.Join(changesRoot, "test-change")
	specsDir := filepath.Join(changeDir, "specs")

	// Create directories
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Create spec files
	for path, content := range specs {
		fullPath := filepath.Join(specsDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	// Validate
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	// Check that issues have line numbers
	foundIssues := 0
	for _, issue := range report.Issues {
		if issue.Line > 0 {
			foundIssues++
			t.Logf("Issue at line %d: %s - %s", issue.Line, issue.Level, issue.Message)
		} else {
			t.Errorf("Issue missing line number: %s - %s", issue.Level, issue.Message)
		}
	}

	if foundIssues == 0 {
		t.Error("Expected to find issues with line numbers")
	}

	// Verify specific line numbers for known issues
	for _, issue := range report.Issues {
		if strings.Contains(issue.Message, "SHALL or MUST") {
			// Line 3 is "### Requirement: User Authentication"
			if issue.Line != 3 {
				t.Errorf("SHALL/MUST error should be at line 3, got %d", issue.Line)
			}

			continue
		}

		if !strings.Contains(issue.Message, "#### Scenario:") {
			continue
		}

		// Line 6 is "##### Scenario: Wrong number of hashtags"
		if issue.Line != 6 {
			t.Errorf("Malformed scenario error should be at line 6, got %d", issue.Line)
		}
	}
}
