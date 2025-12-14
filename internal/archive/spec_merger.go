//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package archive

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

const (
	newlineChar = "\n"
)

// MergeSpec applies delta operations from a delta spec to a base spec
// Returns the merged spec content and operation counts
//
//nolint:revive // specExists is a legitimate control parameter
func MergeSpec(
	baseSpecPath, deltaSpecPath string,
	specExists bool,
) (string, OperationCounts, error) {
	counts := OperationCounts{}

	// Parse delta operations
	deltaPlan, err := parsers.ParseDeltaSpec(
		deltaSpecPath,
	)
	if err != nil {
		return "", counts, fmt.Errorf(
			"parse delta spec: %w",
			err,
		)
	}

	if !deltaPlan.HasDeltas() {
		return "", counts, fmt.Errorf(
			"delta spec has no operations",
		)
	}

	// If spec doesn't exist, create skeleton and only allow ADDED operations
	if !specExists {
		if len(deltaPlan.Modified) > 0 ||
			len(deltaPlan.Removed) > 0 ||
			len(deltaPlan.Renamed) > 0 {
			return "", counts, fmt.Errorf(
				"target spec does not exist; only ADDED requirements are allowed for new specs",
			)
		}
		skeleton := generateSpecSkeleton(
			baseSpecPath,
		)
		merged, addCount := applyAdded(
			skeleton,
			deltaPlan.Added,
		)
		counts.Added = addCount

		return merged, counts, nil
	}

	// Load existing spec
	baseContent, err := os.ReadFile(baseSpecPath)
	if err != nil {
		return "", counts, fmt.Errorf(
			"read base spec: %w",
			err,
		)
	}

	// Parse existing requirements
	baseReqs, err := parsers.ParseRequirements(
		baseSpecPath,
	)
	if err != nil {
		return "", counts, fmt.Errorf(
			"parse base spec: %w",
			err,
		)
	}

	// Build requirement map (normalized name -> block)
	reqMap := make(
		map[string]parsers.RequirementBlock,
	)
	for _, req := range baseReqs {
		normalized := parsers.NormalizeRequirementName(
			req.Name,
		)
		reqMap[normalized] = req
	}

	// Apply operations in order: RENAMED -> REMOVED -> MODIFIED -> ADDED
	reqMap, renameCount := applyRenamed(
		reqMap,
		deltaPlan.Renamed,
	)
	counts.Renamed = renameCount

	reqMap, removeCount := applyRemoved(
		reqMap,
		deltaPlan.Removed,
	)
	counts.Removed = removeCount

	reqMap, modifyCount := applyModified(
		reqMap,
		deltaPlan.Modified,
	)
	counts.Modified = modifyCount

	// ADDED requirements will be appended at the end
	counts.Added = len(deltaPlan.Added)

	// Reconstruct spec
	merged, err := reconstructSpec(
		string(baseContent),
		reqMap,
		deltaPlan.Added,
	)
	if err != nil {
		return "", counts, fmt.Errorf("reconstruct spec: %w", err)
	}

	return merged, counts, nil
}

// applyRenamed updates requirement names in the map
func applyRenamed(
	reqMap map[string]parsers.RequirementBlock,
	renames []parsers.RenameOp,
) (map[string]parsers.RequirementBlock, int) {
	count := 0
	for _, op := range renames {
		fromNorm := parsers.NormalizeRequirementName(
			op.From,
		)
		toNorm := parsers.NormalizeRequirementName(
			op.To,
		)

		req, exists := reqMap[fromNorm]
		if !exists {
			continue
		}

		// Update the header line
		req.HeaderLine = "### Requirement: " + op.To
		// Update the name
		req.Name = op.To
		// Update the raw content (first line)
		lines := strings.Split(req.Raw, "\n")
		if len(lines) > 0 {
			lines[0] = req.HeaderLine
			req.Raw = strings.Join(lines, "\n")
		}
		// Remove old key and add with new key
		delete(reqMap, fromNorm)
		reqMap[toNorm] = req
		count++
	}

	return reqMap, count
}

// applyRemoved removes requirements from the map
func applyRemoved(
	reqMap map[string]parsers.RequirementBlock,
	removed []string,
) (map[string]parsers.RequirementBlock, int) {
	count := 0
	for _, name := range removed {
		normalized := parsers.NormalizeRequirementName(
			name,
		)
		if _, exists := reqMap[normalized]; exists {
			delete(reqMap, normalized)
			count++
		}
	}

	return reqMap, count
}

// applyModified replaces requirements in the map
func applyModified(
	reqMap map[string]parsers.RequirementBlock,
	modified []parsers.RequirementBlock,
) (map[string]parsers.RequirementBlock, int) {
	count := 0
	for _, mod := range modified {
		normalized := parsers.NormalizeRequirementName(
			mod.Name,
		)
		if _, exists := reqMap[normalized]; exists {
			reqMap[normalized] = mod
			count++
		}
	}

	return reqMap, count
}

// applyAdded adds new requirements to spec skeleton
func applyAdded(
	skeleton string,
	added []parsers.RequirementBlock,
) (string, int) {
	if len(added) == 0 {
		return skeleton, 0
	}

	var result strings.Builder
	result.WriteString(skeleton)
	result.WriteString("\n")

	for _, req := range added {
		result.WriteString(
			strings.TrimRight(
				req.Raw,
				newlineChar,
			),
		)
		result.WriteString("\n\n")
	}

	return result.String(), len(added)
}

// reconstructSpec rebuilds the spec from preamble,
// updated requirements, and added requirements.
func reconstructSpec(
	baseContent string,
	reqMap map[string]parsers.RequirementBlock,
	added []parsers.RequirementBlock,
) (string, error) {
	// Split spec into: preamble, requirements section, after
	preamble, reqsContent, after, err := splitSpec(
		baseContent,
	)
	if err != nil {
		return "", fmt.Errorf("split spec: %w", err)
	}

	// Extract original requirement order from base content
	orderedReqs, err := extractOrderedRequirements(
		reqsContent,
		reqMap,
	)
	if err != nil {
		return "", fmt.Errorf("extract ordered requirements: %w", err)
	}

	// Build requirements section
	var reqsBuilder strings.Builder
	for i := range orderedReqs {
		if i > 0 {
			reqsBuilder.WriteString(newlineChar)
		}
		reqsBuilder.WriteString(
			strings.TrimRight(
				orderedReqs[i].Raw,
				newlineChar,
			),
		)
		reqsBuilder.WriteString(newlineChar)
	}

	// Add new requirements at the end
	for _, req := range added {
		reqsBuilder.WriteString(newlineChar)
		reqsBuilder.WriteString(
			strings.TrimRight(
				req.Raw,
				newlineChar,
			),
		)
		reqsBuilder.WriteString(newlineChar)
	}

	// Combine all parts
	var result strings.Builder
	result.WriteString(preamble)
	result.WriteString(reqsBuilder.String())
	result.WriteString(after)

	// Normalize blank lines (collapse 3+ newlines to 2)
	output := result.String()
	multiNewline := regexp.MustCompile(`\n{3,}`)
	output = multiNewline.ReplaceAllString(
		output,
		"\n\n",
	)

	return output, nil
}

// splitSpec splits spec into preamble, requirements section content, and after.
// Uses AST-based parsing via markdown.ParseDocument for accurate section detection.
func splitSpec(
	content string,
) (preamble, requirements, after string, err error) {
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		return "", "", "", fmt.Errorf("parse spec content: %w", err)
	}

	reqSection := doc.GetSection("Requirements")
	if reqSection == nil {
		return "", "", "", fmt.Errorf("no Requirements section found")
	}

	lines := doc.Lines
	// Header.Line is 1-indexed, convert to 0-indexed
	requirementsStart := reqSection.Header.Line - 1

	// Find next H2 header after Requirements
	requirementsEnd := len(lines)
	for _, h := range doc.H2Headers {
		// h.Line is 1-indexed
		if h.Line > reqSection.Header.Line {
			requirementsEnd = h.Line - 1

			break
		}
	}

	// Build preamble (everything up to and including Requirements header)
	preamble = strings.Join(lines[:requirementsStart+1], "\n") + "\n\n"

	// Build requirements section content (between header and next section)
	if requirementsStart+1 < requirementsEnd {
		requirements = strings.Join(
			lines[requirementsStart+1:requirementsEnd],
			"\n",
		)
	}

	// Build after section (from next H2 header onwards)
	if requirementsEnd < len(lines) {
		after = strings.Join(lines[requirementsEnd:], "\n")
	}

	return preamble, requirements, after, nil
}

// extractOrderedRequirements preserves requirement ordering from original content.
// Uses AST-based parsing via markdown.ParseDocument for accurate extraction.
func extractOrderedRequirements(
	reqsContent string,
	reqMap map[string]parsers.RequirementBlock,
) ([]parsers.RequirementBlock, error) {
	// Parse section content to get requirement names in order
	doc, err := markdown.ParseDocument([]byte(reqsContent))
	if err != nil {
		return nil, fmt.Errorf("parse requirements content: %w", err)
	}

	names := doc.GetRequirementNames()
	ordered := make(
		[]parsers.RequirementBlock,
		0,
		len(names),
	)

	for _, name := range names {
		normalized := parsers.NormalizeRequirementName(
			strings.TrimSpace(name),
		)
		req, exists := reqMap[normalized]
		if !exists {
			continue
		}

		ordered = append(ordered, req)
		// Remove from map so we don't add duplicates
		delete(reqMap, normalized)
	}

	// Add any remaining requirements from map (shouldn't happen in normal flow)
	for _, req := range reqMap {
		ordered = append(ordered, req)
	}

	return ordered, nil
}

// generateSpecSkeleton creates a new spec skeleton for a capability
func generateSpecSkeleton(
	targetPath string,
) string {
	// Extract capability name from path
	// (e.g., "spectr/specs/archive-workflow/spec.md" ->
	// "Archive-Workflow")
	parts := strings.Split(targetPath, "/")
	capability := "Capability"
	if len(parts) >= 2 {
		capability = formatCapabilityName(
			parts[len(parts)-2],
		)
	}

	var skeleton strings.Builder
	skeleton.WriteString(
		fmt.Sprintf(
			"# %s Specification\n\n",
			capability,
		),
	)
	skeleton.WriteString("## Requirements\n")

	return skeleton.String()
}

// formatCapabilityName converts kebab-case to Title Case
func formatCapabilityName(kebab string) string {
	words := strings.Split(kebab, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(
				word[:1],
			) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
