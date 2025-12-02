// Package parsers provides utilities for extracting and counting
// information from markdown specification files.
//
// This file contains helper functions for validating the structure of
// tasks.md files, including parsing numbered sections and detecting
// issues like orphaned tasks, empty sections, and non-sequential numbering.
package parsers

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// newTasksStructureResult creates a new TasksStructureResult with all fields
// initialized to their default values. This includes empty slices for
// sections, empty sections list, and non-sequential gaps.
func newTasksStructureResult() *TasksStructureResult {
	return &TasksStructureResult{
		Sections:          make([]TaskSection, 0),
		OrphanedTasks:     0,
		EmptySections:     make([]string, 0),
		SequentialNumbers: true,
		NonSequentialGaps: make([]int, 0),
	}
}

// parseSectionNumber converts a string representation of a section number
// to an integer. Uses strconv.Atoi for safe conversion.
func parseSectionNumber(numStr string) int {
	num, _ := strconv.Atoi(numStr)

	return num
}

// parseTasksFile scans a file line by line and populates the result
// with discovered sections and task counts. It identifies numbered
// section headers (## N. Title) and task checkboxes (- [ ] or - [x]).
func parseTasksFile(file *os.File, result *TasksStructureResult) {
	// Pattern for numbered section headers like "## 1. Section Name"
	sectPat := regexp.MustCompile(`^##\s+([1-9][0-9]*)\.\s+(.+)$`)
	// Pattern for task checkboxes like "- [ ]" or "- [x]"
	taskPat := regexp.MustCompile(`^\s*-\s*\[([xX ])\]`)

	var currentSection *TaskSection
	lineNumber := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check if line matches a numbered section header
		if matches := sectPat.FindStringSubmatch(line); len(matches) > 2 {
			if currentSection != nil {
				result.Sections = append(result.Sections, *currentSection)
			}

			currentSection = &TaskSection{
				Number:    parseSectionNumber(matches[1]),
				Name:      strings.TrimSpace(matches[2]),
				TaskCount: 0,
				Line:      lineNumber,
			}

			continue
		}

		// Skip non-task lines
		if !taskPat.MatchString(line) {
			continue
		}

		// Count task under current section or as orphaned
		if currentSection != nil {
			currentSection.TaskCount++
		} else {
			result.OrphanedTasks++
		}
	}

	// Save the last section if present
	if currentSection != nil {
		result.Sections = append(result.Sections, *currentSection)
	}
}

// finalizeTasksResult performs post-processing on the parsed result,
// computing the list of empty sections and checking for gaps in
// the section numbering sequence.
func finalizeTasksResult(result *TasksStructureResult) {
	result.EmptySections = findEmptySections(result.Sections)
	result.SequentialNumbers, result.NonSequentialGaps = checkSequentialGaps(
		result.Sections,
	)
}

// findEmptySections iterates through all sections and returns the names
// of sections that contain zero tasks.
func findEmptySections(sections []TaskSection) []string {
	var empty []string
	for _, section := range sections {
		if section.TaskCount == 0 {
			empty = append(empty, section.Name)
		}
	}

	return empty
}

// checkSequentialGaps analyzes section numbers to determine if they form
// a sequential series starting from 1. Returns true if sequential, along
// with a slice of any missing numbers (gaps) in the sequence.
func checkSequentialGaps(sections []TaskSection) (bool, []int) {
	if len(sections) == 0 {
		return true, nil
	}

	// Build a set of existing section numbers and find the maximum
	existingNumbers := make(map[int]bool)
	maxNumber := 0

	for _, section := range sections {
		existingNumbers[section.Number] = true
		if section.Number > maxNumber {
			maxNumber = section.Number
		}
	}

	// Find all missing numbers from 1 to maxNumber
	var gaps []int
	for i := 1; i <= maxNumber; i++ {
		if !existingNumbers[i] {
			gaps = append(gaps, i)
		}
	}

	return len(gaps) == 0, gaps
}
