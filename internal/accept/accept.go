// Package accept provides the accept command implementation for converting
// tasks.md to tasks.json, marking a change as ready for implementation.
package accept

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/discovery"
)

// Error variables for accept operations.
var (
	ErrChangeNotFound  = errors.New("change not found")
	ErrNoTasksMd       = errors.New("no tasks.md found")
	ErrAlreadyAccepted = errors.New("already accepted")
	ErrUserCancelled   = errors.New("user cancelled")
)

// AcceptCmd represents the accept command configuration.
type AcceptCmd struct {
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change ID"`
	Yes      bool   `name:"yes" short:"y" help:"Skip confirmation prompt"`
}

// Accept accepts a change by converting tasks.md to tasks.json.
//
// The projectPath parameter allows operating in a different working
// directory. If empty, the current working directory is used.
func Accept(cmd *AcceptCmd, projectPath string) error {
	projectRoot, err := getProjectRoot(projectPath)
	if err != nil {
		return err
	}

	spectrRoot := filepath.Join(projectRoot, "spectr")
	if err = validateSpectrDir(spectrRoot, projectRoot); err != nil {
		return err
	}

	changeID, err := resolveChangeID(cmd.ChangeID, projectRoot)
	if err != nil {
		return err
	}
	if changeID == "" {
		return nil // User cancelled
	}

	changeDir := filepath.Join(spectrRoot, "changes", changeID)
	tasksMdPath, tasksJSONPath, err := validateChangePaths(changeDir, changeID)
	if err != nil {
		return err
	}

	fmt.Printf("Accepting change: %s\n\n", changeID)

	if !cmd.Yes && !confirm("Convert tasks.md to tasks.json?") {
		return ErrUserCancelled
	}

	return processAcceptance(tasksMdPath, tasksJSONPath, changeID)
}

// validateSpectrDir checks if spectr directory exists.
func validateSpectrDir(spectrRoot, projectRoot string) error {
	if _, err := os.Stat(spectrRoot); os.IsNotExist(err) {
		return fmt.Errorf("spectr directory not found in %s", projectRoot)
	}

	return nil
}

// validateChangePaths validates change directory and returns paths.
func validateChangePaths(
	changeDir, changeID string,
) (tasksMdPath, tasksJSONPath string, err error) {
	if _, err = os.Stat(changeDir); os.IsNotExist(err) {
		return "", "", fmt.Errorf("change '%s' not found", changeID)
	}

	tasksMdPath = filepath.Join(changeDir, "tasks.md")
	if _, err = os.Stat(tasksMdPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("no tasks.md found for change '%s'", changeID)
	}

	tasksJSONPath = filepath.Join(changeDir, "tasks.json")
	if _, err = os.Stat(tasksJSONPath); err == nil {
		return "", "", fmt.Errorf("change '%s' already accepted", changeID)
	}

	return tasksMdPath, tasksJSONPath, nil
}

// processAcceptance parses tasks.md and writes tasks.json.
func processAcceptance(tasksMdPath, tasksJSONPath, changeID string) error {
	sections, err := ParseTasksFile(tasksMdPath)
	if err != nil {
		return fmt.Errorf("parse tasks.md: %w", err)
	}

	err = WriteTasksJSON(tasksJSONPath, changeID, sections)
	if err != nil {
		return fmt.Errorf("write tasks.json: %w", err)
	}

	err = os.Remove(tasksMdPath)
	if err != nil {
		return fmt.Errorf("remove tasks.md: %w", err)
	}

	fmt.Printf("\n Successfully accepted change: %s\n", changeID)
	fmt.Printf("  Created: %s\n", tasksJSONPath)

	return nil
}

// getProjectRoot returns the project root directory.
// If projectPath is empty, returns the current working directory.
func getProjectRoot(projectPath string) (string, error) {
	if projectPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}

		return cwd, nil
	}

	return projectPath, nil
}

// resolveChangeID resolves a change ID, using interactive selection if empty.
// Returns empty string and nil error if user cancelled.
func resolveChangeID(changeID, projectRoot string) (string, error) {
	if changeID == "" {
		selectedID, selectErr := selectChangeInteractive(projectRoot)
		if errors.Is(selectErr, ErrUserCancelled) {
			return "", nil // User cancelled, exit gracefully
		}
		if selectErr != nil {
			return "", fmt.Errorf("select change: %w", selectErr)
		}

		return selectedID, nil
	}

	// Resolve partial ID to full change ID
	result, resolveErr := discovery.ResolveChangeID(changeID, projectRoot)
	if resolveErr != nil {
		return "", fmt.Errorf("change '%s' not found", changeID)
	}
	if result.PartialMatch {
		fmt.Printf("Resolved '%s' -> '%s'\n\n", changeID, result.ChangeID)
	}

	return result.ChangeID, nil
}

// selectChangeInteractive uses interactive selection for change.
// Returns ErrUserCancelled if the user cancels the selection.
func selectChangeInteractive(projectRoot string) (string, error) {
	acceptableChanges, err := getAcceptableChanges(projectRoot)
	if err != nil {
		return "", err
	}

	if len(acceptableChanges) == 0 {
		fmt.Println("No changes with tasks.md ready for acceptance.")

		return "", ErrUserCancelled
	}

	displayChangeList(acceptableChanges)

	return promptForSelection(acceptableChanges)
}

// getAcceptableChanges returns changes that have tasks.md but not tasks.json.
func getAcceptableChanges(projectRoot string) ([]string, error) {
	changes, err := discovery.GetActiveChangeIDs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("list changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("No changes found.")

		return nil, ErrUserCancelled
	}

	var acceptableChanges []string
	for _, change := range changes {
		if isChangeAcceptable(projectRoot, change) {
			acceptableChanges = append(acceptableChanges, change)
		}
	}

	return acceptableChanges, nil
}

// isChangeAcceptable checks if a change has tasks.md but not tasks.json.
func isChangeAcceptable(projectRoot, change string) bool {
	changeDir := filepath.Join(projectRoot, "spectr", "changes", change)
	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	tasksJSONPath := filepath.Join(changeDir, "tasks.json")

	if _, err := os.Stat(tasksMdPath); err != nil {
		return false
	}

	_, err := os.Stat(tasksJSONPath)

	return os.IsNotExist(err)
}

// displayChangeList prints the list of acceptable changes.
func displayChangeList(changes []string) {
	fmt.Println("Changes ready for acceptance:")
	for i, change := range changes {
		fmt.Printf("  %d. %s\n", i+1, change)
	}
	fmt.Println()
}

// promptForSelection prompts for and validates user selection.
func promptForSelection(changes []string) (string, error) {
	fmt.Print("Enter change number (or 'q' to quit): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", ErrUserCancelled
	}

	response = strings.TrimSpace(response)
	if response == "q" || response == "Q" || response == "" {
		return "", ErrUserCancelled
	}

	selection, err := parseSelectionNumber(response)
	if err != nil {
		return "", err
	}

	if selection < 1 || selection > len(changes) {
		return "", fmt.Errorf("selection out of range: %d", selection)
	}

	return changes[selection-1], nil
}

// parseSelectionNumber parses a numeric selection string.
func parseSelectionNumber(response string) (int, error) {
	var selection int
	for _, c := range response {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid selection: %s", response)
		}
		selection = selection*decimalBase + int(c-'0')
	}

	return selection, nil
}

// confirm prompts user for yes/no confirmation.
func confirm(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}
