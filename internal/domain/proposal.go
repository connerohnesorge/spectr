// Package domain contains shared domain types used across the Spectr codebase.
package domain

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Dependency represents a dependency on another proposal.
type Dependency struct {
	ID     string `yaml:"id"`
	Reason string `yaml:"reason,omitempty"`
}

// ProposalMetadata contains metadata extracted from proposal.md frontmatter.
type ProposalMetadata struct {
	// ID is the optional explicit ID for this proposal
	ID string `yaml:"id,omitempty"`
	// Requires lists proposals that must be archived before this can be accepted
	Requires []Dependency `yaml:"requires,omitempty"`
	// Enables lists proposals that this proposal unlocks (informational only)
	Enables []Dependency `yaml:"enables,omitempty"`
}

// HasDependencies returns true if the proposal has any requires dependencies.
func (m *ProposalMetadata) HasDependencies() bool {
	return len(m.Requires) > 0
}

// RequiredIDs returns a slice of all required proposal IDs.
func (m *ProposalMetadata) RequiredIDs() []string {
	ids := make([]string, len(m.Requires))
	for i, dep := range m.Requires {
		ids[i] = dep.ID
	}

	return ids
}

// EnabledIDs returns a slice of all enabled proposal IDs.
func (m *ProposalMetadata) EnabledIDs() []string {
	ids := make([]string, len(m.Enables))
	for i, dep := range m.Enables {
		ids[i] = dep.ID
	}

	return ids
}

// ParseProposalFrontmatter extracts and parses YAML frontmatter from proposal.md content.
// Returns empty ProposalMetadata if no frontmatter is present (backward compatible).
// Returns an error if frontmatter is present but contains invalid YAML.
func ParseProposalFrontmatter(content []byte) (*ProposalMetadata, error) {
	fm, err := ExtractFrontmatter(content)
	if err != nil {
		return nil, err
	}

	if fm == nil {
		// No frontmatter - return empty metadata (backward compatible)
		return &ProposalMetadata{}, nil
	}

	var meta ProposalMetadata
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return nil, fmt.Errorf("invalid YAML in frontmatter: %w", err)
	}

	return &meta, nil
}

// ParseProposalFrontmatterFromFile reads a proposal.md file and parses its frontmatter.
func ParseProposalFrontmatterFromFile(path string) (*ProposalMetadata, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseProposalFrontmatter(content)
}

const frontmatterDelimiter = "---"

// ExtractFrontmatter extracts the YAML frontmatter block from markdown content.
// Returns nil if no frontmatter is present (content doesn't start with "---").
// Returns an error if frontmatter starts but is not properly closed.
func ExtractFrontmatter(content []byte) ([]byte, error) {
	reader := bufio.NewReader(bytes.NewReader(content))

	// Read first line
	firstLine, err := reader.ReadString('\n')
	if err != nil && firstLine == "" {
		// Empty file or read error on first line
		return nil, nil
	}

	// Check if file starts with frontmatter delimiter
	if strings.TrimSpace(firstLine) != frontmatterDelimiter {
		return nil, nil
	}

	// Read until closing delimiter
	var fmContent bytes.Buffer
	for {
		line, readErr := reader.ReadString('\n')
		if readErr != nil {
			if line == "" {
				// EOF without closing delimiter
				return nil, errors.New(
					"frontmatter not closed: missing closing '---'",
				)
			}
			// Handle last line without newline
			if strings.TrimSpace(line) == frontmatterDelimiter {
				break
			}
			fmContent.WriteString(line)

			return nil, errors.New(
				"frontmatter not closed: missing closing '---'",
			)
		}

		if strings.TrimSpace(line) == frontmatterDelimiter {
			break
		}

		fmContent.WriteString(line)
	}

	return fmContent.Bytes(), nil
}

// ValidateProposalMetadata validates the metadata for common errors.
// Returns an error if the proposal references itself in requires.
func ValidateProposalMetadata(meta *ProposalMetadata, proposalID string) error {
	// Check for self-reference in requires
	for _, dep := range meta.Requires {
		if dep.ID == proposalID {
			return fmt.Errorf("proposal cannot require itself: %s", proposalID)
		}
	}

	// Check for self-reference in enables (less critical but still odd)
	for _, dep := range meta.Enables {
		if dep.ID == proposalID {
			return fmt.Errorf("proposal cannot enable itself: %s", proposalID)
		}
	}

	return nil
}
