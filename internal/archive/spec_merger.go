//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package archive

import (
	"fmt"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/mdparser"
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
	// Parse the base content into AST
	doc, err := mdparser.Parse(baseContent)
	if err != nil {
		// Fallback to empty doc on parse error
		doc = &mdparser.Document{Children: []mdparser.Node{}}
	}

	// Find Requirements section index and extract ordering
	reqsSectionIdx, orderedReqs := extractRequirementsSection(doc, reqMap)
	missingRequirementsSection := reqsSectionIdx == -1
	if missingRequirementsSection {
		orderedReqs = collectRequirementsWithoutSection(doc, reqMap)
	}

	// Rebuild document with updated requirements
	var result strings.Builder

	preambleEnd := reqsSectionIdx
	if missingRequirementsSection {
		preambleEnd = len(doc.Children)
	}

	// Write everything before Requirements section
	for i := 0; i < preambleEnd && i < len(doc.Children); i++ {
		renderNode(&result, doc.Children[i])
	}

	// Write Requirements header if we found one
	if reqsSectionIdx >= 0 && reqsSectionIdx < len(doc.Children) {
		if header, ok := doc.Children[reqsSectionIdx].(*mdparser.Header); ok {
			result.WriteString(strings.Repeat("#", header.Level))
			result.WriteString(" ")
			result.WriteString(header.Text)
			result.WriteString("\n\n")
		}
	} else {
		// No Requirements section found, add one
		result.WriteString("## Requirements\n\n")
	}

	// Write ordered requirements
	for i, req := range orderedReqs {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(strings.TrimRight(req.Raw, newlineChar))
		result.WriteString("\n")
	}

	// Add new requirements at the end
	for _, req := range added {
		result.WriteString("\n")
		result.WriteString(strings.TrimRight(req.Raw, newlineChar))
		result.WriteString("\n")
	}

	// Write everything after Requirements section
	afterReqsIdx := findNextSectionAfterRequirements(doc, reqsSectionIdx)
	if afterReqsIdx >= 0 {
		result.WriteString("\n")
		for i := afterReqsIdx; i < len(doc.Children); i++ {
			renderNode(&result, doc.Children[i])
		}
	}

	// Normalize blank lines using AST-aware approach
	return normalizeBlankLines(result.String())
}

// extractRequirementsSection finds the Requirements section and extracts requirements in order
// Returns: (sectionIndex, orderedRequirements)
func extractRequirementsSection(
	doc *mdparser.Document,
	reqMap map[string]parsers.RequirementBlock,
) (int, []parsers.RequirementBlock) {
	// Find Requirements section header (H2)
	reqsSectionIdx := -1
	for i, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != 2 {
			continue
		}
		if strings.TrimSpace(header.Text) == "Requirements" {
			reqsSectionIdx = i

			break
		}
	}

	if reqsSectionIdx == -1 {
		// No Requirements section found
		return -1, nil
	}

	// Find the end of the Requirements section (next H2 or end of document)
	endIdx := len(doc.Children)
	for i := reqsSectionIdx + 1; i < len(doc.Children); i++ {
		if header, ok := doc.Children[i].(*mdparser.Header); ok && header.Level == 2 {
			endIdx = i

			break
		}
	}

	// Extract requirements in order from the section
	var ordered []parsers.RequirementBlock
	for i := reqsSectionIdx + 1; i < endIdx; i++ {
		header, ok := doc.Children[i].(*mdparser.Header)
		if !ok || header.Level != 3 {
			continue
		}

		// Check if this is a requirement header
		if !strings.HasPrefix(header.Text, "Requirement: ") {
			continue
		}

		// Extract requirement name
		name := strings.TrimPrefix(header.Text, "Requirement: ")
		name = strings.TrimSpace(name)
		normalized := parsers.NormalizeRequirementName(name)

		// Look up in reqMap
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

	return reqsSectionIdx, ordered
}

// collectRequirementsWithoutSection orders requirements using document order when the
// Requirements section is missing, allowing us to still merge updated requirement blocks.
func collectRequirementsWithoutSection(
	doc *mdparser.Document,
	reqMap map[string]parsers.RequirementBlock,
) []parsers.RequirementBlock {
	var ordered []parsers.RequirementBlock

	for _, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != 3 {
			continue
		}
		if !strings.HasPrefix(header.Text, "Requirement: ") {
			continue
		}

		name := strings.TrimSpace(strings.TrimPrefix(header.Text, "Requirement: "))
		normalized := parsers.NormalizeRequirementName(name)

		req, exists := reqMap[normalized]
		if !exists {
			continue
		}

		ordered = append(ordered, req)
		delete(reqMap, normalized)
	}

	for _, req := range reqMap {
		ordered = append(ordered, req)
	}

	return ordered
}

// findNextSectionAfterRequirements finds the index of the next H2 section after Requirements
func findNextSectionAfterRequirements(doc *mdparser.Document, reqsSectionIdx int) int {
	if reqsSectionIdx < 0 {
		return -1
	}

	for i := reqsSectionIdx + 1; i < len(doc.Children); i++ {
		if header, ok := doc.Children[i].(*mdparser.Header); ok && header.Level == 2 {
			return i
		}
	}

	return -1
}

// renderNode converts an AST node back to markdown text
func renderNode(sb *strings.Builder, node mdparser.Node) {
	switch n := node.(type) {
	case *mdparser.Header:
		sb.WriteString(strings.Repeat("#", n.Level))
		sb.WriteString(" ")
		sb.WriteString(n.Text)
		sb.WriteString("\n")

	case *mdparser.Paragraph:
		for _, line := range n.Lines {
			sb.WriteString(line)
			sb.WriteString("\n")
		}

	case *mdparser.CodeBlock:
		sb.WriteString("```")
		sb.WriteString(n.Language)
		sb.WriteString("\n")
		for _, line := range n.Lines {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("```\n")

	case *mdparser.List:
		for _, item := range n.Items {
			if n.Ordered {
				sb.WriteString("1. ")
			} else {
				sb.WriteString("- ")
			}
			sb.WriteString(item.Text)
			sb.WriteString("\n")
		}

	case *mdparser.BlankLine:
		// Render blank lines but cap at 2 consecutive
		count := n.Count
		if count > 2 {
			count = 2
		}

		for range count {
			sb.WriteString("\n")
		}
	}
}

// normalizeBlankLines collapses 3+ consecutive newlines to 2
func normalizeBlankLines(content string) string {
	// Parse to understand structure
	doc, err := mdparser.Parse(content)
	if err != nil {
		// If parsing fails, use simple string replacement
		lines := strings.Split(content, "\n")
		var result []string
		blankCount := 0

		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				blankCount++
				if blankCount <= 2 {
					result = append(result, line)
				}
			} else {
				blankCount = 0
				result = append(result, line)
			}
		}

		return strings.Join(result, "\n")
	}

	// Render with normalized blank lines
	var sb strings.Builder
	for _, node := range doc.Children {
		renderNode(&sb, node)
	}

	return sb.String()
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
