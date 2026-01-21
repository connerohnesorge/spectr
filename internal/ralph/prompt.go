// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	errTaskNil        = errors.New("task cannot be nil")
	errChangeDirEmpty = errors.New("changeDir cannot be empty")
)

const doubleNewline = "\n\n"

// prompt.go handles the generation of prompts for agent CLI sessions.
// It assembles comprehensive task execution prompts by:
// - Loading change context (proposal.md, design.md)
// - Loading relevant delta spec files
// - Formatting task details with dependencies
// - Including file paths and references
//
// The prompt structure ensures agents have all necessary context
// to understand and execute their assigned tasks correctly.

// PromptTemplate holds the data needed to generate a prompt for a task.
type PromptTemplate struct {
	Task        *Task
	Proposal    string
	Design      string
	DeltaSpecs  map[string]string // spec name -> spec content
	TasksJSONCs []string          // paths to tasks.jsonc files for status updates
}

// GeneratePrompt creates a comprehensive prompt for an agent to execute a task.
// It reads the change context (proposal, design, specs) and assembles them into
// a structured prompt that gives the agent all necessary information.
//
// Parameters:
//   - task: The task to generate a prompt for
//   - changeDir: Path to the change directory (e.g., spectr/changes/<id>)
//
// Returns:
//   - string: The complete prompt ready to be passed to the agent CLI
//   - error: Any error encountered during file reading
func GeneratePrompt(task *Task, changeDir string) (string, error) {
	if task == nil {
		return "", errTaskNil
	}
	if changeDir == "" {
		return "", errChangeDirEmpty
	}

	template := &PromptTemplate{
		Task:       task,
		DeltaSpecs: make(map[string]string),
	}

	// Load proposal.md (required)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	proposalContent, err := os.ReadFile(proposalPath)
	if err != nil {
		return "", fmt.Errorf("failed to read proposal.md: %w", err)
	}
	template.Proposal = string(proposalContent)

	// Load design.md (optional)
	designPath := filepath.Join(changeDir, "design.md")
	designContent, err := os.ReadFile(designPath)
	if err == nil {
		template.Design = string(designContent)
	}
	// Ignore error if design.md doesn't exist

	// Load delta specs from specs/ directory
	specsDir := filepath.Join(changeDir, "specs")
	// Ignore errors - specs directory might be empty or missing
	// This is not critical for prompt generation
	_ = loadDeltaSpecs(specsDir, template)

	// Find all tasks.jsonc files for reference
	tasksPattern := filepath.Join(changeDir, "tasks*.jsonc")
	matches, err := filepath.Glob(tasksPattern)
	if err == nil && len(matches) > 0 {
		template.TasksJSONCs = matches
	}

	// Assemble the prompt
	return assemblePrompt(template), nil
}

// loadDeltaSpecs discovers and loads all delta spec markdown files from the specs directory.
// It recursively walks the specs directory and reads all spec.md files.
func loadDeltaSpecs(specsDir string, template *PromptTemplate) error {
	// Check if specs directory exists
	if _, err := os.Stat(specsDir); os.IsNotExist(err) {
		return nil // Not an error, just no specs to load
	}

	// Walk the specs directory to find all spec.md files
	err := filepath.Walk(specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() {
			return nil
		}

		// Only process spec.md files
		if filepath.Base(path) != "spec.md" {
			return nil
		}

		// Read the spec file
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("failed to read %s: %w", path, readErr)
		}

		// Extract the capability name from the directory structure
		// e.g., specs/ralph-orchestration/spec.md -> ralph-orchestration
		relPath, _ := filepath.Rel(specsDir, filepath.Dir(path))
		capabilityName := filepath.Base(relPath)

		template.DeltaSpecs[capabilityName] = string(content)

		return nil
	})

	return err
}

// assemblePrompt builds the final prompt string from the template data.
func assemblePrompt(template *PromptTemplate) string {
	var builder strings.Builder

	// Header
	builder.WriteString(
		fmt.Sprintf("# Task: %s - %s"+doubleNewline, template.Task.ID, template.Task.Section),
	)

	// Task Description
	builder.WriteString("## Task Description\n")
	builder.WriteString(template.Task.Description)
	builder.WriteString(doubleNewline)

	// Change Context
	builder.WriteString("## Change Context" + doubleNewline)

	// Proposal
	builder.WriteString("### Proposal\n")
	builder.WriteString(template.Proposal)
	builder.WriteString(doubleNewline)

	// Design (if exists)
	if template.Design != "" {
		builder.WriteString("### Design\n")
		builder.WriteString(template.Design)
		builder.WriteString(doubleNewline)
	}

	// Relevant Specs (if any)
	if len(template.DeltaSpecs) > 0 {
		builder.WriteString("### Relevant Specs" + doubleNewline)
		for capabilityName, content := range template.DeltaSpecs {
			builder.WriteString(fmt.Sprintf("#### %s"+doubleNewline, capabilityName))
			builder.WriteString(content)
			builder.WriteString(doubleNewline)
		}
	}

	// Instructions
	builder.WriteString("## Instructions\n")
	builder.WriteString(
		"Complete this task and update the task status in tasks.jsonc to \"completed\" ",
	)
	builder.WriteString(
		"when done. If blocked, set status to \"in_progress\" and describe the blocker.\n",
	)

	// Reference to tasks.jsonc files
	if len(template.TasksJSONCs) > 0 {
		builder.WriteString("\nTask status files:\n")
		for _, path := range template.TasksJSONCs {
			builder.WriteString(fmt.Sprintf("- %s\n", path))
		}
	}

	return builder.String()
}
