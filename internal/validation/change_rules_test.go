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
