// Package archive provides utilities for archiving completed changes,
// including merging spec deltas and moving change directories to archive.
//
// The archive workflow consists of several stages:
// 1. Validation - Ensures the change passes all spec validation rules
// 2. Task checking - Verifies implementation tasks are complete
// 3. Spec merging - Applies delta specs to main specifications
// 4. Archiving - Moves change to archive with timestamp prefix
//
// This package enforces the Spectr workflow where changes move from
// spectr/changes/ to spectr/changes/archive/YYYY-MM-DD-<name>/ after
// deployment, with their delta specs merged into spectr/specs/.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package archive

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/conneroisu/spectr/internal/parsers"
)

// Archiver manages the archive workflow for completed changes.
// It coordinates validation, spec merging, and directory movement.
type Archiver struct {
	yes         bool // Skip confirmation prompts
	skipSpecs   bool // Skip spec delta merging
	noValidate  bool // Skip pre-archive validation
	interactive bool // Use interactive table mode
}

// NewArchiver creates a new Archiver with the given flags.
// The flags control which steps of the archive workflow are executed:
// - yes: Auto-confirm all prompts (non-interactive mode)
// - skipSpecs: Skip merging delta specs into main specs
// - noValidate: Skip validation checks before archiving
// - interactive: Use interactive table mode for change selection
func NewArchiver(yes, skipSpecs, noValidate, interactive bool) (*Archiver, error) {
	return &Archiver{
		yes:         yes,
		skipSpecs:   skipSpecs,
		noValidate:  noValidate,
		interactive: interactive,
	}, nil
}

// Archive archives a change by validating, applying specs, and moving to archive directory
//
//nolint:revive // changeID parameter needs to be reassigned when empty
func (a *Archiver) Archive(changeID string) error {
	// Get current working directory as project root
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Check if spectr directory exists
	spectrRoot := filepath.Join(projectRoot, "spectr")
	if _, err := os.Stat(spectrRoot); os.IsNotExist(err) {
		return fmt.Errorf("spectr directory not found in %s", projectRoot)
	}

	// If no change ID provided, use interactive selection
	if changeID == "" {
		selectedID, err := a.selectChange(projectRoot, spectrRoot)
		if err != nil {
			return fmt.Errorf("select change: %w", err)
		}
		if selectedID == "" {
			return fmt.Errorf("no change selected")
		}
		changeID = selectedID
	}

	changeDir := filepath.Join(spectrRoot, "changes", changeID)

	// Check if change exists
	if _, err := os.Stat(changeDir); os.IsNotExist(err) {
		return fmt.Errorf("change not found: %s", changeID)
	}

	fmt.Printf("Archiving change: %s\n\n", changeID)

	// Validation workflow
	if !a.noValidate {
		if err := a.runValidation(changeDir); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	} else {
		if !a.yes {
			if !a.confirm("Validation is disabled. Continue anyway?") {
				return errors.New("archive cancelled")
			}
		}
		fmt.Println("⚠️  Skipping validation")
	}

	// Task checking
	if err := a.checkTasks(changeDir); err != nil {
		return fmt.Errorf("task check failed: %w", err)
	}

	// Spec update workflow
	if !a.skipSpecs {
		if err := a.updateSpecs(changeDir, spectrRoot); err != nil {
			return fmt.Errorf("spec update failed: %w", err)
		}
	} else {
		fmt.Println("⚠️  Skipping spec updates")
	}

	// Archive operation
	if err := a.moveToArchive(changeDir, changeID, spectrRoot); err != nil {
		return fmt.Errorf("move to archive failed: %w", err)
	}

	fmt.Printf("\n✓ Successfully archived: %s\n", changeID)

	return nil
}

// selectChange prompts user to select a change interactively
func (a *Archiver) selectChange(projectRoot, spectrRoot string) (string, error) {
	// Use interactive table mode if enabled
	if a.interactive {
		return a.selectChangeInteractive(projectRoot)
	}

	// Fallback to numbered list selection
	return a.selectChangeBasic(spectrRoot)
}

// selectChangeInteractive uses the interactive table for change selection
func (*Archiver) selectChangeInteractive(projectRoot string) (string, error) {
	// Import list package functions
	// Note: This will be done at the package level
	lister := newListerForArchive(projectRoot)
	changes, err := lister.ListChanges()
	if err != nil {
		return "", fmt.Errorf("list changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("No changes found.")

		return "", nil
	}

	selectedID, err := runInteractiveArchiveForArchiver(changes, projectRoot)
	if err != nil {
		return "", fmt.Errorf("interactive selection: %w", err)
	}

	return selectedID, nil
}

// selectChangeBasic uses basic numbered list for change selection
func (*Archiver) selectChangeBasic(spectrRoot string) (string, error) {
	changesDir := filepath.Join(spectrRoot, "changes")

	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return "", fmt.Errorf("read changes directory: %w", err)
	}

	var changes []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "archive" {
			changes = append(changes, entry.Name())
		}
	}

	if len(changes) == 0 {
		return "", errors.New("no changes found")
	}

	fmt.Println("Available changes:")
	for i, change := range changes {
		fmt.Printf("  %d. %s\n", i+1, change)
	}

	fmt.Print("\nSelect change number: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}

	var selection int
	trimmed := strings.TrimSpace(input)
	if _, err := fmt.Sscanf(trimmed, "%d", &selection); err != nil {
		return "", fmt.Errorf("invalid selection: %w", err)
	}

	if selection < 1 || selection > len(changes) {
		return "", errors.New("selection out of range")
	}

	return changes[selection-1], nil
}

// runValidation validates the change before archiving
func (*Archiver) runValidation(changeDir string) error {
	fmt.Println("Validating change...")

	report, err := ValidatePreArchive(changeDir, true)
	if err != nil {
		return err
	}

	if !report.Valid {
		fmt.Printf("❌ Validation failed: %d error(s), %d warning(s)\n",
			report.Summary.Errors, report.Summary.Warnings)

		for _, issue := range report.Issues {
			fmt.Printf(
				"  [%s] %s: %s\n",
				issue.Level,
				issue.Path,
				issue.Message,
			)
		}

		return errors.New("validation errors must be fixed before archiving")
	}

	if report.Summary.Warnings > 0 {
		fmt.Printf(
			"⚠️  Validation passed with %d warning(s)\n",
			report.Summary.Warnings,
		)
	} else {
		fmt.Println("✓ Validation passed")
	}

	return nil
}

// checkTasks checks task completion status
func (a *Archiver) checkTasks(changeDir string) error {
	tasksPath := filepath.Join(changeDir, "tasks.md")
	status, err := parsers.CountTasks(tasksPath)
	if err != nil {
		// tasks.md is optional
		return nil
	}

	if status.Total == 0 {
		return nil
	}

	incomplete := status.Total - status.Completed
	fmt.Printf("Tasks: %d/%d completed", status.Completed, status.Total)

	if incomplete > 0 {
		fmt.Printf(" (%d incomplete)\n", incomplete)
		if !a.yes {
			if !a.confirm("Archive with incomplete tasks?") {
				return errors.New("archive cancelled due to incomplete tasks")
			}
		}
	} else {
		fmt.Println()
	}

	return nil
}

// updateSpecs applies delta specs to main specs
func (a *Archiver) updateSpecs(changeDir, spectrRoot string) error {
	specsDir := filepath.Join(changeDir, "specs")
	deltaSpecs, err := a.findAndValidateDeltaSpecs(specsDir)
	if err != nil {
		return err
	}

	if len(deltaSpecs) == 0 {
		fmt.Println("No spec deltas found")

		return nil
	}

	updates, err := a.buildUpdatePlan(deltaSpecs, specsDir, spectrRoot)
	if err != nil {
		return err
	}

	a.displayUpdatePlan(updates)

	if !a.yes && !a.confirm("\nApply spec updates?") {
		return errors.New("archive cancelled")
	}

	totalCounts, mergedSpecs, err := a.processMerges(updates)
	if err != nil {
		return err
	}

	if err := a.writeSpecs(mergedSpecs); err != nil {
		return err
	}

	a.displaySummary(totalCounts)

	return nil
}

// findAndValidateDeltaSpecs finds delta specs in the given directory
func (a *Archiver) findAndValidateDeltaSpecs(
	specsDir string,
) ([]string, error) {
	if _, err := os.Stat(specsDir); os.IsNotExist(err) {
		fmt.Println("No spec deltas found")

		return nil, nil
	}

	deltaSpecs, err := a.findDeltaSpecs(specsDir)
	if err != nil {
		return nil, fmt.Errorf("find delta specs: %w", err)
	}

	return deltaSpecs, nil
}

// buildUpdatePlan creates spec update plan from delta specs
func (*Archiver) buildUpdatePlan(
	deltaSpecs []string,
	specsDir, spectrRoot string,
) ([]SpecUpdate, error) {
	updates := make([]SpecUpdate, 0, len(deltaSpecs))

	for _, deltaPath := range deltaSpecs {
		relPath, err := filepath.Rel(specsDir, deltaPath)
		if err != nil {
			return nil, fmt.Errorf("get relative path: %w", err)
		}

		capabilityDir := filepath.Dir(relPath)
		targetPath := filepath.Join(
			spectrRoot,
			"specs",
			capabilityDir,
			"spec.md",
		)

		exists := false
		if _, err := os.Stat(targetPath); err == nil {
			exists = true
		}

		updates = append(updates, SpecUpdate{
			Source: deltaPath,
			Target: targetPath,
			Exists: exists,
		})
	}

	return updates, nil
}

// displayUpdatePlan prints the update plan to console
func (*Archiver) displayUpdatePlan(updates []SpecUpdate) {
	fmt.Printf("\nSpec updates (%d):\n", len(updates))
	for _, update := range updates {
		capability := filepath.Base(filepath.Dir(update.Target))
		status := "update"
		if !update.Exists {
			status = "create"
		}
		fmt.Printf("  [%s] %s\n", status, capability)
	}
}

// processMerges validates and merges all spec updates
func (*Archiver) processMerges(
	updates []SpecUpdate,
) (OperationCounts, map[string]string, error) {
	totalCounts := OperationCounts{}
	mergedSpecs := make(map[string]string)

	for _, update := range updates {
		merged, counts, err := processOneMerge(update)
		if err != nil {
			return totalCounts, nil, err
		}

		mergedSpecs[update.Target] = merged
		totalCounts.Added += counts.Added
		totalCounts.Modified += counts.Modified
		totalCounts.Removed += counts.Removed
		totalCounts.Renamed += counts.Renamed
	}

	return totalCounts, mergedSpecs, nil
}

// processOneMerge handles validation and merging for a single update
func processOneMerge(
	update SpecUpdate,
) (string, OperationCounts, error) {
	deltaPlan, err := parsers.ParseDeltaSpec(update.Source)
	if err != nil {
		return "", OperationCounts{},
			fmt.Errorf("parse delta spec %s: %w", update.Source, err)
	}

	err = CheckDuplicatesAndConflicts(deltaPlan)
	if err != nil {
		return "", OperationCounts{},
			fmt.Errorf(
				"delta validation failed for %s: %w",
				update.Source,
				err,
			)
	}

	err = ValidatePreMerge(update.Target, deltaPlan, update.Exists)
	if err != nil {
		return "", OperationCounts{},
			fmt.Errorf(
				"pre-merge validation failed for %s: %w",
				update.Source,
				err,
			)
	}

	merged, counts, err := MergeSpec(
		update.Target,
		update.Source,
		update.Exists,
	)
	if err != nil {
		return "", OperationCounts{},
			fmt.Errorf("merge spec %s: %w", update.Source, err)
	}

	if err := ValidatePostMerge(merged, update.Target); err != nil {
		return "", OperationCounts{},
			fmt.Errorf(
				"post-merge validation failed for %s: %w",
				update.Source,
				err,
			)
	}

	return merged, counts, nil
}

// writeSpecs writes all merged specs to disk
func (*Archiver) writeSpecs(mergedSpecs map[string]string) error {
	for targetPath, content := range mergedSpecs {
		if err := os.MkdirAll(
			filepath.Dir(targetPath),
			dirPerm,
		); err != nil {
			return fmt.Errorf("create spec directory: %w", err)
		}

		if err := os.WriteFile(
			targetPath,
			[]byte(content),
			filePerm,
		); err != nil {
			return fmt.Errorf("write spec %s: %w", targetPath, err)
		}
	}

	return nil
}

// displaySummary prints operation summary to console
func (*Archiver) displaySummary(totalCounts OperationCounts) {
	fmt.Println("\nSpec operations applied:")
	if totalCounts.Added > 0 {
		fmt.Printf("  + %d added\n", totalCounts.Added)
	}
	if totalCounts.Modified > 0 {
		fmt.Printf("  ~ %d modified\n", totalCounts.Modified)
	}
	if totalCounts.Removed > 0 {
		fmt.Printf("  - %d removed\n", totalCounts.Removed)
	}
	if totalCounts.Renamed > 0 {
		fmt.Printf("  → %d renamed\n", totalCounts.Renamed)
	}
	fmt.Printf("  = %d total\n", totalCounts.Total())
}

// findDeltaSpecs recursively finds all spec.md files in a directory
func (*Archiver) findDeltaSpecs(dir string) ([]string, error) {
	var specs []string

	err := filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() == "spec.md" {
				specs = append(specs, path)
			}

			return nil
		},
	)

	return specs, err
}

// moveToArchive moves the change directory to archive with date prefix
func (*Archiver) moveToArchive(
	changeDir, changeID, spectrRoot string,
) error {
	// Create archive directory if it doesn't exist
	archiveDir := filepath.Join(spectrRoot, "changes", "archive")
	if err := os.MkdirAll(archiveDir, dirPerm); err != nil {
		return fmt.Errorf("create archive directory: %w", err)
	}

	// Generate archive name with date
	date := time.Now().Format("2006-01-02")
	archiveName := fmt.Sprintf("%s-%s", date, changeID)
	archivePath := filepath.Join(archiveDir, archiveName)

	// Check if archive already exists
	if _, err := os.Stat(archivePath); err == nil {
		return fmt.Errorf("archive already exists: %s", archiveName)
	}

	// Move change to archive
	if err := os.Rename(changeDir, archivePath); err != nil {
		return fmt.Errorf("move to archive: %w", err)
	}

	fmt.Printf("\nMoved to: changes/archive/%s\n", archiveName)

	return nil
}

// confirm prompts user for yes/no confirmation
func (*Archiver) confirm(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}
