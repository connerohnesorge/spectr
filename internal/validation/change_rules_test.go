package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// Helper function to create a change directory with spec files
// Returns both the changeDir and spectrRoot paths
func createChangeDir(
	t *testing.T,
	specs map[string]string,
) (changeDir, spectrRoot string) {
	t.Helper()

	// Create project root
	projectRoot := t.TempDir()
	spectrRoot = filepath.Join(projectRoot, "spectr")
	changesRoot := filepath.Join(spectrRoot, "changes")
	specsRoot := filepath.Join(spectrRoot, "specs")

	// Create change directory
	changeDir = filepath.Join(changesRoot, "test-change")
	specsDir := filepath.Join(changeDir, "specs")

	// Create necessary directories
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}
	if err := os.MkdirAll(specsRoot, 0755); err != nil {
		t.Fatalf("Failed to create spectr/specs directory: %v", err)
	}

	// Create delta spec files
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

	return changeDir, spectrRoot
}

// Helper function to create a base spec file in the spectr/specs directory
func createBaseSpec(t *testing.T, spectrRoot, capability, content string) {
	t.Helper()

	specPath := filepath.Join(spectrRoot, "specs", capability, "spec.md")
	dir := filepath.Dir(specPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(specPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write base spec %s: %v", specPath, err)
	}
}

func TestValidateChangeDeltaSpecs_ValidAddedRequirements(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication functionality.

#### Scenario: Successful login
- **WHEN** user provides valid credentials
- **THEN** user is authenticated

#### Scenario: Failed login
- **WHEN** user provides invalid credentials
- **THEN** authentication fails
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid && len(report.Issues) == 0 {
		return
	}
	t.Errorf("Expected valid report, got invalid with %d errors", report.Summary.Errors)
	for _, issue := range report.Issues {
		t.Logf("  %s: %s - %s", issue.Level, issue.Path, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_ValidModifiedRequirements(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## MODIFIED Requirements

### Requirement: User Authentication
The system MUST provide enhanced user authentication functionality.

#### Scenario: Two-factor authentication
- **WHEN** user provides valid credentials and OTP
- **THEN** user is authenticated with 2FA
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)

	// Create base spec with the requirement that will be modified
	createBaseSpec(t, spectrRoot, "auth", `## Requirements

### Requirement: User Authentication
The system SHALL provide user authentication functionality.

#### Scenario: Successful login
- **WHEN** user provides valid credentials
- **THEN** user is authenticated
`)

	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		return
	}
	t.Error("Expected valid report, got invalid")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s - %s", issue.Level, issue.Path, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_ValidRemovedRequirements(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## REMOVED Requirements

### Requirement: Legacy Authentication
**Reason**: Replaced by modern OAuth flow
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)

	// Create base spec with the requirement that will be removed
	createBaseSpec(t, spectrRoot, "auth", `## Requirements

### Requirement: Legacy Authentication
The system SHALL provide legacy authentication.

#### Scenario: Legacy login
- **WHEN** user uses legacy credentials
- **THEN** user is authenticated via legacy method
`)

	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		return
	}
	t.Error("Expected valid report, got invalid")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s - %s", issue.Level, issue.Path, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_ValidRenamedRequirements(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## RENAMED Requirements

- FROM: ### Requirement: Login
- TO: ### Requirement: User Authentication
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		return
	}
	t.Error("Expected valid report, got invalid")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s - %s", issue.Level, issue.Path, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_MultipleSpecFiles(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: Two-Factor Authentication
The system SHALL provide two-factor authentication.

#### Scenario: OTP required
- **WHEN** valid credentials are provided
- **THEN** OTP challenge is required
`,
		"notifications/spec.md": `## ADDED Requirements

### Requirement: Email Notifications
The system MUST send email notifications for authentication events.

#### Scenario: Login notification
- **WHEN** user logs in
- **THEN** email notification is sent
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		return
	}
	t.Error("Expected valid report, got invalid")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s - %s", issue.Level, issue.Path, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_NoDeltas(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `# Some content without delta sections

This file doesn't have any ADDED, MODIFIED, REMOVED, or RENAMED sections.
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to no deltas")
	}

	if report.Summary.Errors != 1 {
		t.Errorf("Expected 1 error, got %d", report.Summary.Errors)
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "at least one delta") {
			found = true

			break
		}
	}
	if !found {
		t.Error("Expected error about missing deltas")
	}
}

func TestValidateChangeDeltaSpecs_EmptyDeltaSections(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

## MODIFIED Requirements
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to empty delta sections")
	}

	// Should have 2 errors: one for empty ADDED, one for empty MODIFIED
	if report.Summary.Errors == 2 {
		return
	}
	t.Errorf("Expected 2 errors, got %d", report.Summary.Errors)
	for _, issue := range report.Issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_AddedWithoutShallMust(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system provides user authentication functionality.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to missing SHALL/MUST")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "SHALL or MUST") {
			found = true

			break
		}
	}
	if !found {
		t.Error("Expected error about missing SHALL or MUST")
	}
}

func TestValidateChangeDeltaSpecs_AddedWithoutScenarios(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication functionality.
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to missing scenarios")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "at least one scenario") {
			found = true

			break
		}
	}
	if !found {
		t.Error("Expected error about missing scenarios")
	}
}

func TestValidateChangeDeltaSpecs_ModifiedWithoutShallMust(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## MODIFIED Requirements

### Requirement: User Authentication
The system provides enhanced authentication.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to missing SHALL/MUST")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "SHALL or MUST") {
			found = true

			break
		}
	}
	if !found {
		t.Error("Expected error about missing SHALL or MUST")
	}
}

func TestValidateChangeDeltaSpecs_ModifiedWithoutScenarios(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## MODIFIED Requirements

### Requirement: User Authentication
The system SHALL provide enhanced authentication.
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to missing scenarios")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "at least one scenario") {
			found = true

			break
		}
	}
	if !found {
		t.Error("Expected error about missing scenarios")
	}
}

func TestValidateChangeDeltaSpecs_DuplicateRequirementNames(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated

### Requirement: User Authentication
The system SHALL also do something else.

#### Scenario: Another scenario
- **WHEN** something happens
- **THEN** something occurs
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to duplicate requirement names")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(issue.Message, "Duplicate requirement name") {
			found = true

			break
		}
	}
	if !found {
		t.Error("Expected error about duplicate requirement names")
	}
}

func TestValidateChangeDeltaSpecs_CrossSectionConflicts(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated

## MODIFIED Requirements

### Requirement: User Authentication
The system SHALL provide enhanced authentication.

#### Scenario: Enhanced login
- **WHEN** user logs in
- **THEN** user gets enhanced authentication
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to cross-section conflicts")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "both ADDED and MODIFIED") {
			found = true

			break
		}
	}
	if found {
		return
	}
	t.Error("Expected error about cross-section conflict")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_MalformedRenamedFormat(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## RENAMED Requirements

- FROM: ### Requirement: Old Name
This is missing the TO line
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to malformed RENAMED format")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(issue.Message, "Malformed RENAMED requirement") {
			found = true

			break
		}
	}
	if found {
		return
	}
	t.Error("Expected error about malformed RENAMED format")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_MissingSpecsDirectory(t *testing.T) {
	// Create a simple temp dir structure
	projectRoot := t.TempDir()
	spectrRoot := filepath.Join(projectRoot, "spectr")
	changesRoot := filepath.Join(spectrRoot, "changes")
	changeDir := filepath.Join(changesRoot, "test-change")

	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}
	// Don't create specs directory

	_, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err == nil {
		t.Error("Expected error for missing specs directory")
	}

	if !strings.Contains(err.Error(), "specs directory not found") {
		t.Errorf("Expected error about missing specs directory, got: %v", err)
	}
}

func TestValidateChangeDeltaSpecs_NoSpecFiles(t *testing.T) {
	// Create a simple temp dir structure
	projectRoot := t.TempDir()
	spectrRoot := filepath.Join(projectRoot, "spectr")
	changesRoot := filepath.Join(spectrRoot, "changes")
	changeDir := filepath.Join(changesRoot, "test-change")
	specsDir := filepath.Join(changeDir, "specs")

	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Create empty specs directory
	_, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err == nil {
		t.Error("Expected error for no spec.md files")
	}

	if !strings.Contains(err.Error(), "no spec.md files found") {
		t.Errorf("Expected error about no spec files, got: %v", err)
	}
}

func TestValidateChangeDeltaSpecs_MultipleFilesWithConflicts(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated
`,
		"security/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide secure authentication.

#### Scenario: Secure login
- **WHEN** user logs in securely
- **THEN** user is authenticated securely
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to duplicate requirement across files")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "ADDED in multiple files") {
			found = true

			break
		}
	}
	if found {
		return
	}
	t.Error("Expected error about duplicate requirement across files")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_MalformedScenarios(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication.

##### Scenario: Wrong number of hashtags
- **WHEN** user logs in
- **THEN** user is authenticated
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to malformed scenarios")
	}

	// Should have 2 errors: missing scenario (since malformed ones don't count) + malformed scenario format
	foundMissingScenario := false
	foundMalformedFormat := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError && strings.Contains(issue.Message, "at least one scenario") {
			foundMissingScenario = true
		}
		if issue.Level == LevelError && strings.Contains(issue.Message, "#### Scenario:") {
			foundMalformedFormat = true
		}
	}

	if foundMissingScenario && foundMalformedFormat {
		return
	}
	t.Errorf(
		"Expected both missing scenario and malformed format errors. Found missing=%v, found malformed=%v",
		foundMissingScenario,
		foundMalformedFormat,
	)
	for _, issue := range report.Issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_StrictMode(t *testing.T) {
	// In strict mode, any warnings would be converted to errors
	// Since change validation uses errors by default, this test
	// ensures strict mode doesn't break anything
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, true)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if !report.Valid {
		t.Error("Expected valid report in strict mode for valid change")
		for _, issue := range report.Issues {
			t.Logf("  %s: %s", issue.Level, issue.Message)
		}
	}

	// Verify no warnings (all should be converted to errors if any exist)
	if report.Summary.Warnings != 0 {
		t.Errorf("Expected 0 warnings in strict mode, got %d", report.Summary.Warnings)
	}
}

func TestValidateChangeDeltaSpecs_AllDeltaTypes(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## ADDED Requirements

### Requirement: Two-Factor Authentication
The system SHALL provide two-factor authentication.

#### Scenario: OTP required
- **WHEN** user logs in
- **THEN** OTP is required

## MODIFIED Requirements

### Requirement: Password Policy
The system MUST enforce stronger password policies.

#### Scenario: Password strength
- **WHEN** user sets password
- **THEN** password meets strength requirements

## REMOVED Requirements

### Requirement: Legacy Login
**Reason**: Deprecated in favor of OAuth

## RENAMED Requirements

- FROM: ### Requirement: User Login
- TO: ### Requirement: User Authentication
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)

	// Create base spec with requirements for MODIFIED, REMOVED, and RENAMED
	createBaseSpec(t, spectrRoot, "auth", `## Requirements

### Requirement: Password Policy
The system SHALL enforce password policies.

#### Scenario: Basic password
- **WHEN** user sets password
- **THEN** password is validated

### Requirement: Legacy Login
The system SHALL provide legacy login.

#### Scenario: Legacy auth
- **WHEN** user logs in with legacy method
- **THEN** user is authenticated

### Requirement: User Login
The system SHALL provide user login.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated
`)

	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if !report.Valid {
		t.Error("Expected valid report with all delta types")
		for _, issue := range report.Issues {
			t.Logf("  %s: %s", issue.Level, issue.Message)
		}
	}

	if len(report.Issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(report.Issues))
	}
}

func TestValidateChangeDeltaSpecs_DuplicateRenamedFromNames(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## RENAMED Requirements

- FROM: ### Requirement: Old Name
- TO: ### Requirement: New Name One

- FROM: ### Requirement: Old Name
- TO: ### Requirement: New Name Two
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to duplicate FROM names")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(issue.Message, "Duplicate FROM requirement name") {
			found = true

			break
		}
	}
	if found {
		return
	}
	t.Error("Expected error about duplicate FROM names")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateChangeDeltaSpecs_DuplicateRenamedToNames(t *testing.T) {
	specs := map[string]string{
		"auth/spec.md": `## RENAMED Requirements

- FROM: ### Requirement: Old Name One
- TO: ### Requirement: New Name

- FROM: ### Requirement: Old Name Two
- TO: ### Requirement: New Name
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to duplicate TO names")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level != LevelError ||
			!strings.Contains(issue.Message, "Duplicate TO requirement name") {
			continue
		}
		found = true
		if issue.Line != 7 {
			t.Fatalf("expected duplicate TO issue at line 7, got %d", issue.Line)
		}

		break
	}
	if found {
		return
	}
	t.Error("Expected error about duplicate TO names")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s (line %d)", issue.Level, issue.Message, issue.Line)
	}
}

func TestValidateChangeDeltaSpecs_RenamedToAcrossFilesLineNumber(t *testing.T) {
	specs := map[string]string{
		"alpha/spec.md": `## RENAMED Requirements

- FROM: ### Requirement: Old Name Alpha
- TO: ### Requirement: Shared Name
`,
		"beta/spec.md": `## RENAMED Requirements

- FROM: ### Requirement: Old Name Beta
- TO: ### Requirement: Shared Name
`,
	}

	changeDir, spectrRoot := createChangeDir(t, specs)
	report, err := ValidateChangeDeltaSpecs(changeDir, spectrRoot, false)
	if err != nil {
		t.Fatalf("ValidateChangeDeltaSpecs returned error: %v", err)
	}

	if report.Valid {
		t.Error("Expected invalid report due to cross-file TO duplicates")
	}

	found := false
	for _, issue := range report.Issues {
		if issue.Level != LevelError ||
			!strings.Contains(issue.Message, "renamed (TO) in multiple files") {
			continue
		}
		found = true
		if issue.Line != 4 {
			t.Fatalf("expected cross-file TO issue at line 4, got %d", issue.Line)
		}

		break
	}

	if found {
		return
	}
	t.Error("Expected error about cross-file TO duplicates")
	for _, issue := range report.Issues {
		t.Logf("  %s: %s (line %d)", issue.Level, issue.Message, issue.Line)
	}
}

func TestFindRenamedPairLine_FindsBulletEntries(t *testing.T) {
	lines := []string{
		"## RENAMED Requirements",
		"",
		"- FROM: ### Requirement: Old Name",
		"- TO: ### Requirement: New Name",
		"",
	}

	line := findRenamedPairLine(lines, "Old Name", 1)
	if line != 3 {
		t.Fatalf("expected line 3 for FROM entry, got %d", line)
	}

	toLine := findRenamedPairLine(lines, "New Name", 1)
	if toLine != 4 {
		t.Fatalf("expected line 4 for TO entry, got %d", toLine)
	}
}

func TestFindPreMergeErrorLine_UsesRenamedBulletLine(t *testing.T) {
	lines := []string{
		"## RENAMED Requirements",
		"",
		"- FROM: ### Requirement: Old Name",
		"- TO: ### Requirement: New Name",
		"",
	}

	fromErr := `RENAMED FROM requirement "Old Name" does not exist in base spec`
	if line := findPreMergeErrorLine(lines, fromErr, &parsers.DeltaPlan{}); line != 3 {
		t.Fatalf("expected FROM error to map to line 3, got %d", line)
	}

	toErr := `RENAMED TO requirement "New Name" already exists in base spec`
	if line := findPreMergeErrorLine(lines, toErr, &parsers.DeltaPlan{}); line != 4 {
		t.Fatalf("expected TO error to map to line 4, got %d", line)
	}
}

// Helper function to create a tasks.md file for testing
func createTasksFile(t *testing.T, dir, content string) string {
	t.Helper()
	tasksPath := filepath.Join(dir, "tasks.md")
	if err := os.WriteFile(tasksPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write tasks.md: %v", err)
	}

	return tasksPath
}

func TestValidateTasksFile_ValidStructure(t *testing.T) {
	// Well-formed tasks.md with numbered sections
	tempDir := t.TempDir()
	content := `# Tasks

## 1. Setup

- [ ] Initialize project
- [x] Install dependencies

## 2. Implementation

- [ ] Write core functionality
- [ ] Add error handling

## 3. Testing

- [x] Write unit tests
- [ ] Run integration tests
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	if len(issues) == 0 {
		return
	}

	t.Errorf("Expected 0 issues for valid tasks.md, got %d", len(issues))
	for _, issue := range issues {
		t.Logf("  %s: %s (line %d)", issue.Level, issue.Message, issue.Line)
	}
}

func TestValidateTasksFile_NoSections(t *testing.T) {
	// Tasks without numbered sections
	tempDir := t.TempDir()
	content := `# Tasks

Some intro text without sections.

- [ ] Task one
- [x] Task two
- [ ] Task three
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	// Expect warning about no sections and orphaned tasks
	if len(issues) < 1 {
		t.Errorf("Expected at least 1 issue, got %d", len(issues))
	}

	foundNoSections := false
	for _, issue := range issues {
		if issue.Level == LevelWarning &&
			strings.Contains(issue.Message, "no numbered sections") {
			foundNoSections = true

			break
		}
	}

	if foundNoSections {
		return
	}

	t.Error("Expected warning about no numbered sections")
	for _, issue := range issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateTasksFile_OrphanedTasks(t *testing.T) {
	// Some tasks before first section
	tempDir := t.TempDir()
	content := `# Tasks

- [ ] Orphaned task 1
- [ ] Orphaned task 2

## 1. Setup

- [ ] Task in section
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	foundOrphanedWarning := false
	for _, issue := range issues {
		if issue.Level != LevelWarning {
			continue
		}

		if !strings.Contains(issue.Message, "tasks outside numbered sections") {
			continue
		}

		foundOrphanedWarning = true
		// Verify it says 2 orphaned tasks
		if !strings.Contains(issue.Message, "2 tasks") {
			t.Errorf("Expected warning about 2 orphaned tasks, got: %s", issue.Message)
		}

		break
	}

	if foundOrphanedWarning {
		return
	}

	t.Error("Expected warning about orphaned tasks")
	for _, issue := range issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateTasksFile_EmptySections(t *testing.T) {
	// Sections without tasks
	tempDir := t.TempDir()
	content := `# Tasks

## 1. Setup

- [ ] Setup task

## 2. Empty Section

## 3. Another Empty Section

## 4. Final Section

- [ ] Final task
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	// Count warnings for empty sections
	emptyCount := 0
	for _, issue := range issues {
		if issue.Level == LevelWarning && strings.Contains(issue.Message, "has no tasks") {
			emptyCount++
		}
	}

	// Should have 2 empty section warnings
	if emptyCount != 2 {
		t.Errorf("Expected 2 empty section warnings, got %d", emptyCount)
		for _, issue := range issues {
			t.Logf("  %s: %s (line %d)", issue.Level, issue.Message, issue.Line)
		}
	}

	// Verify specific section names are mentioned
	foundEmpty1 := false
	foundEmpty2 := false
	for _, issue := range issues {
		if strings.Contains(issue.Message, "Empty Section") {
			foundEmpty1 = true
		}
		if strings.Contains(issue.Message, "Another Empty Section") {
			foundEmpty2 = true
		}
	}

	if !foundEmpty1 || !foundEmpty2 {
		t.Error("Expected warnings to mention specific empty section names")
	}
}

func TestValidateTasksFile_NonSequentialSections(t *testing.T) {
	// Sections 1, 3, 5 - missing 2 and 4
	tempDir := t.TempDir()
	content := `# Tasks

## 1. First Section

- [ ] Task 1

## 3. Third Section

- [ ] Task 3

## 5. Fifth Section

- [ ] Task 5
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	foundNonSequential := false
	for _, issue := range issues {
		if issue.Level != LevelWarning {
			continue
		}

		if !strings.Contains(issue.Message, "not sequential") {
			continue
		}

		foundNonSequential = true
		// Check that missing numbers are mentioned
		if !strings.Contains(issue.Message, "2") || !strings.Contains(issue.Message, "4") {
			t.Errorf("Expected warning to mention missing numbers 2 and 4, got: %s", issue.Message)
		}

		break
	}

	if foundNonSequential {
		return
	}

	t.Error("Expected warning about non-sequential section numbers")
	for _, issue := range issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateTasksFile_NonexistentFile(t *testing.T) {
	// File doesn't exist - should return no issues (optional validation)
	tasksPath := filepath.Join(t.TempDir(), "nonexistent", "tasks.md")

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error for nonexistent file: %v", err)
	}

	if len(issues) == 0 {
		return
	}

	t.Errorf("Expected 0 issues for nonexistent tasks.md, got %d", len(issues))
	for _, issue := range issues {
		t.Logf("  %s: %s", issue.Level, issue.Message)
	}
}

func TestValidateTasksFile_MultipleIssues(t *testing.T) {
	// File with multiple problems:
	// - Orphaned tasks
	// - Empty section
	// - Non-sequential sections (1, 3)
	tempDir := t.TempDir()
	content := `# Tasks

- [ ] Orphaned task

## 1. First Section

- [ ] Task in first section

## 3. Third Section (skipped 2)

## 5. Empty Section (skipped 4)
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	// Should have at least 3 issues:
	// 1. Orphaned tasks warning
	// 2. Non-sequential numbers warning
	// 3-4. Empty section warnings (sections 3 and 5 have no tasks)

	if len(issues) < 3 {
		t.Errorf("Expected at least 3 issues for file with multiple problems, got %d", len(issues))
	}

	// Verify all issue types are present
	foundOrphaned := false
	foundNonSequential := false
	foundEmptySection := false

	for _, issue := range issues {
		if strings.Contains(issue.Message, "tasks outside numbered sections") {
			foundOrphaned = true
		}
		if strings.Contains(issue.Message, "not sequential") {
			foundNonSequential = true
		}
		if strings.Contains(issue.Message, "has no tasks") {
			foundEmptySection = true
		}
	}

	if !foundOrphaned {
		t.Error("Expected orphaned tasks warning")
	}
	if !foundNonSequential {
		t.Error("Expected non-sequential numbers warning")
	}
	if !foundEmptySection {
		t.Error("Expected empty section warning")
	}

	// Log all issues for debugging
	t.Log("All issues found:")
	for _, issue := range issues {
		t.Logf("  %s: %s (line %d)", issue.Level, issue.Message, issue.Line)
	}
}

func TestValidateTasksFile_AllWarningsAreLevelWarning(t *testing.T) {
	// Verify that all issues from ValidateTasksFile are warnings, not errors
	tempDir := t.TempDir()
	content := `# Tasks

- [ ] Orphaned task

## 2. Skipped First Section

`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	for _, issue := range issues {
		if issue.Level != LevelWarning {
			t.Errorf("Expected all issues to be warnings, got %s: %s", issue.Level, issue.Message)
		}
	}
}

func TestValidateTasksFile_LineNumbers(t *testing.T) {
	// Verify that line numbers are correctly reported for empty sections
	tempDir := t.TempDir()
	content := `# Tasks

## 1. First Section

## 2. Second Empty Section

## 3. Third Section
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	// Find the empty section warnings and check line numbers
	for _, issue := range issues {
		checkSectionLineNumber(t, issue, "First Section", 3)
		checkSectionLineNumber(t, issue, "Second Empty Section", 5)
		checkSectionLineNumber(t, issue, "Third Section", 7)
	}
}

func checkSectionLineNumber(t *testing.T, issue ValidationIssue, section string, line int) {
	t.Helper()

	if !strings.Contains(issue.Message, section) {
		return
	}

	if !strings.Contains(issue.Message, "has no tasks") {
		return
	}

	if issue.Line != line {
		t.Errorf("Expected %s warning at line %d, got line %d", section, line, issue.Line)
	}
}

func TestValidateTasksFile_PathIncludedInIssues(t *testing.T) {
	// Verify that the tasks.md path is included in all issues
	tempDir := t.TempDir()
	content := `# Tasks

- [ ] Orphaned task
`
	tasksPath := createTasksFile(t, tempDir, content)

	issues, err := ValidateTasksFile(tasksPath)
	if err != nil {
		t.Fatalf("ValidateTasksFile returned error: %v", err)
	}

	for _, issue := range issues {
		if issue.Path != tasksPath {
			t.Errorf("Expected issue path to be %s, got %s", tasksPath, issue.Path)
		}
	}
}
