//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package archive

import (
	"fmt"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parser"
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
	deltaPlan, err := parsers.ParseDeltaSpec(deltaSpecPath)
	if err != nil {
		return "", counts, fmt.Errorf("parse delta spec: %w", err)
	}

	if !deltaPlan.HasDeltas() {
		return "", counts, fmt.Errorf("delta spec has no operations")
	}

	// If spec doesn't exist, create skeleton and only allow ADDED operations
	if !specExists {
		if len(deltaPlan.Modified) > 0 || len(deltaPlan.Removed) > 0 || len(deltaPlan.Renamed) > 0 {
			return "", counts, fmt.Errorf(
				"target spec does not exist; only ADDED requirements are allowed for new specs",
			)
		}
		skeleton := generateSpecSkeleton(baseSpecPath)
		merged, addCount := applyAdded(skeleton, deltaPlan.Added)
		counts.Added = addCount

		return merged, counts, nil
	}

	// Load existing spec
	baseContent, err := os.ReadFile(baseSpecPath)
	if err != nil {
		return "", counts, fmt.Errorf("read base spec: %w", err)
	}

	// Parse existing requirements
	baseReqs, err := parsers.ParseRequirements(baseSpecPath)
	if err != nil {
		return "", counts, fmt.Errorf("parse base spec: %w", err)
	}

	// Build requirement map (normalized name -> block)
	reqMap := make(map[string]parsers.RequirementBlock)
	for _, req := range baseReqs {
		normalized := parsers.NormalizeRequirementName(req.Name)
		reqMap[normalized] = req
	}

	// Apply operations in order: RENAMED -> REMOVED -> MODIFIED -> ADDED
	reqMap, renameCount := applyRenamed(reqMap, deltaPlan.Renamed)
	counts.Renamed = renameCount

	reqMap, removeCount := applyRemoved(reqMap, deltaPlan.Removed)
	counts.Removed = removeCount

	reqMap, modifyCount := applyModified(reqMap, deltaPlan.Modified)
	counts.Modified = modifyCount

	// ADDED requirements will be appended at the end
	counts.Added = len(deltaPlan.Added)

	// Reconstruct spec
	merged := reconstructSpec(string(baseContent), reqMap, deltaPlan.Added)

	return merged, counts, nil
}

// applyRenamed updates requirement names in the map
func applyRenamed(
	reqMap map[string]parsers.RequirementBlock,
	renames []parsers.RenameOp,
) (map[string]parsers.RequirementBlock, int) {
	count := 0
	for _, op := range renames {
		fromNorm := parsers.NormalizeRequirementName(op.From)
		toNorm := parsers.NormalizeRequirementName(op.To)

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
		normalized := parsers.NormalizeRequirementName(name)
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
		normalized := parsers.NormalizeRequirementName(mod.Name)
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
		result.WriteString(strings.TrimRight(req.Raw, newlineChar))
		result.WriteString("\n\n")
	}

	return result.String(), len(added)
}

// reconstructSpec rebuilds the spec from preamble,
// updated requirements, and added requirements
func reconstructSpec(
	baseContent string,
	reqMap map[string]parsers.RequirementBlock,
	added []parsers.RequirementBlock,
) string {
	// Split spec into: preamble, requirements section, after
	preamble, reqsContent, after := splitSpec(baseContent)

	// Extract original requirement order from base content
	orderedReqs := extractOrderedRequirements(reqsContent, reqMap)

	// Build requirements section
	var reqsBuilder strings.Builder
	for i := range orderedReqs {
		if i > 0 {
			reqsBuilder.WriteString(newlineChar)
		}
		reqsBuilder.WriteString(
			strings.TrimRight(orderedReqs[i].Raw, newlineChar),
		)
		reqsBuilder.WriteString(newlineChar)
	}

	// Add new requirements at the end
	for _, req := range added {
		reqsBuilder.WriteString(newlineChar)
		reqsBuilder.WriteString(strings.TrimRight(req.Raw, newlineChar))
		reqsBuilder.WriteString(newlineChar)
	}

	// Combine all parts
	var result strings.Builder
	result.WriteString(preamble)
	result.WriteString(reqsBuilder.String())
	result.WriteString(after)

	// Normalize blank lines (collapse 3+ newlines to 2)
	output := normalizeBlankLines(result.String())

	return output
}

// normalizeBlankLines collapses 3+ consecutive newlines into 2.
//
// This normalization ensures consistent spacing throughout the spec
// without using regex. It iterates through the string and tracks
// consecutive newlines, outputting at most 2 in a row.
func normalizeBlankLines(content string) string {
	if content == "" {
		return content
	}

	var result strings.Builder
	result.Grow(len(content)) // Pre-allocate for efficiency

	consecutiveNewlines := 0

	//nolint:intrange // Keep compatible with Go <1.22
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			consecutiveNewlines++
			// Only write newlines if we haven't exceeded the limit
			if consecutiveNewlines <= 2 {
				result.WriteByte('\n')
			}
		} else {
			consecutiveNewlines = 0
			result.WriteByte(content[i])
		}
	}

	return result.String()
}

// splitSpec splits spec into preamble, requirements section content, and after
func splitSpec(content string) (preamble, requirements, after string) {
	// Parse the document
	doc, err := parser.Parse(content)
	if err != nil {
		// On parse error, return everything as preamble
		return content, "", ""
	}

	// Find the "## Requirements" header
	reqHeaders := parser.FindHeaders(doc, func(h *parser.Header) bool {
		return h.Level == 2 && strings.TrimSpace(h.Text) == "Requirements"
	})

	if len(reqHeaders) == 0 {
		// No requirements section, return everything as preamble
		return content, "", ""
	}

	reqHeader := reqHeaders[0]
	reqHeaderPos := reqHeader.Pos()

	// Find the next level-2 header after Requirements
	nextHeaders := parser.FindHeaders(doc, func(h *parser.Header) bool {
		return h.Level == 2 && h.Pos().Offset > reqHeaderPos.Offset
	})

	// Extract preamble (everything up to and including "## Requirements\n\n")
	preambleEndOffset := reqHeaderPos.Offset + len("## Requirements")
	// Find the end of the requirements header line
	for i := preambleEndOffset; i < len(content); i++ {
		if content[i] == '\n' {
			preambleEndOffset = i + 1

			break
		}
	}
	preamble = content[:preambleEndOffset] + "\n"

	// Extract requirements section and after
	if len(nextHeaders) > 0 {
		// There's a section after Requirements
		nextHeaderPos := nextHeaders[0].Pos()
		requirements = content[preambleEndOffset:nextHeaderPos.Offset]
		after = content[nextHeaderPos.Offset:]
	} else {
		// Requirements section goes to end of document
		requirements = content[preambleEndOffset:]
		after = ""
	}

	return preamble, requirements, after
}

// extractOrderedRequirements preserves requirement ordering
// from original content
func extractOrderedRequirements(
	reqsContent string,
	reqMap map[string]parsers.RequirementBlock,
) []parsers.RequirementBlock {
	var ordered []parsers.RequirementBlock

	// Parse the requirements content
	doc, err := parser.Parse(reqsContent)
	if err != nil {
		// On parse error, return requirements from map in arbitrary order
		for _, req := range reqMap {
			ordered = append(ordered, req)
		}

		return ordered
	}

	// Find requirement headers in order
	reqHeaders := parser.FindHeaders(doc, func(h *parser.Header) bool {
		return h.Level == 3 && strings.HasPrefix(strings.TrimSpace(h.Text), "Requirement:")
	})

	for _, header := range reqHeaders {
		// Extract requirement name
		name := strings.TrimPrefix(strings.TrimSpace(header.Text), "Requirement:")
		name = strings.TrimSpace(name)
		normalized := parsers.NormalizeRequirementName(name)

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

	return ordered
}

// generateSpecSkeleton creates a new spec skeleton for a capability
func generateSpecSkeleton(targetPath string) string {
	// Extract capability name from path
	// (e.g., "spectr/specs/archive-workflow/spec.md" ->
	// "Archive-Workflow")
	parts := strings.Split(targetPath, "/")
	capability := "Capability"
	if len(parts) >= 2 {
		capability = formatCapabilityName(parts[len(parts)-2])
	}

	var skeleton strings.Builder
	skeleton.WriteString(fmt.Sprintf("# %s Specification\n\n", capability))
	skeleton.WriteString("## Purpose\n\n")
	skeleton.WriteString("TODO: Add purpose description\n\n")
	skeleton.WriteString("## Requirements\n")

	return skeleton.String()
}

// formatCapabilityName converts kebab-case to Title Case
func formatCapabilityName(kebab string) string {
	words := strings.Split(kebab, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
