//nolint:revive // file-length-limit: wikilink validation requires comprehensive helpers
package markdown

import (
	"os"
	"path/filepath"
	"strings"
)

// WikilinkError represents an error found during wikilink validation.
// It contains details about the broken wikilink and where it was found.
type WikilinkError struct {
	// Target is the wikilink target that failed to resolve.
	Target string

	// Display is the optional display text of the wikilink.
	Display string

	// Anchor is the optional anchor/fragment of the wikilink.
	Anchor string

	// Offset is the byte offset where the wikilink starts in the source.
	Offset int

	// Message describes why the wikilink is invalid.
	Message string
}

// Error implements the error interface.
func (e WikilinkError) Error() string {
	if e.Offset >= 0 {
		return "offset " + itoa(
			e.Offset,
		) + ": " + e.Message
	}

	return e.Message
}

// ResolveWikilink resolves a wikilink target to a file path within the project.
// It follows the Spectr resolution rules:
//  1. First check spectr/specs/{target}/spec.md
//  2. Then check spectr/changes/{target}/proposal.md
//  3. If target contains "/", treat first segment as directory type
//
// Returns the resolved path and whether the file exists.
// The projectRoot should be the root directory containing the spectr/ folder.
func ResolveWikilink(
	target, projectRoot string,
) (path string, exists bool) {
	if target == "" {
		return "", false
	}

	// Strip any anchor from the target for path resolution
	cleanTarget := target
	if idx := strings.Index(target, "#"); idx >= 0 {
		cleanTarget = target[:idx]
	}

	// Handle targets that explicitly specify a directory type
	if strings.HasPrefix(
		cleanTarget,
		"changes/",
	) {
		// Explicit change target: changes/my-change -> proposal.md
		changeName := strings.TrimPrefix(
			cleanTarget,
			"changes/",
		)
		path = filepath.Join(
			projectRoot,
			"spectr",
			"changes",
			changeName,
			"proposal.md",
		)
		exists = fileExists(path)

		return path, exists
	}

	if strings.HasPrefix(cleanTarget, "specs/") {
		// Explicit spec target: specs/validation -> spec.md
		specName := strings.TrimPrefix(
			cleanTarget,
			"specs/",
		)
		path = filepath.Join(
			projectRoot,
			"spectr",
			"specs",
			specName,
			"spec.md",
		)
		exists = fileExists(path)

		return path, exists
	}

	// Default resolution order: specs first, then changes

	// Try spectr/specs/{target}/spec.md
	specPath := filepath.Join(
		projectRoot,
		"spectr",
		"specs",
		cleanTarget,
		"spec.md",
	)
	if fileExists(specPath) {
		return specPath, true
	}

	// Try spectr/changes/{target}/proposal.md
	changePath := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		cleanTarget,
		"proposal.md",
	)
	if fileExists(changePath) {
		return changePath, true
	}

	// Return the spec path as the "expected" path even though it doesn't exist
	return specPath, false
}

// ResolveWikilinkWithAnchor resolves a wikilink target and validates anchor.
// It first resolves the target path using ResolveWikilink, then parses the
// target file to validate that the anchor exists as a header.
//
// Returns:
//   - path: the resolved file path
//   - anchorValid: true if anchor exists in target (or if no anchor given)
//   - err: any error encountered during resolution or file reading
//
// Anchor matching rules:
//   - Anchors are matched case-insensitively
//   - "Requirement: Name" matches ### Requirement: Name headers
//   - "Scenario: Name" matches #### Scenario: Name headers
//   - Plain text matches any header containing that text
func ResolveWikilinkWithAnchor(
	target, anchor, projectRoot string,
) (path string, anchorValid bool, err error) {
	// Resolve the target path
	path, exists := ResolveWikilink(
		target,
		projectRoot,
	)
	if !exists {
		return path, false, nil
	}

	// If no anchor, resolution is complete
	if anchor == "" {
		return path, true, nil
	}

	// Read and parse the target file to find the anchor
	content, err := os.ReadFile(path)
	if err != nil {
		return path, false, err
	}

	// Check if the anchor exists in the file
	anchorValid = anchorExistsInContent(
		content,
		anchor,
	)

	return path, anchorValid, nil
}

// anchorExistsInContent checks if an anchor exists as a header in the content.
// It parses the content and searches for matching headers.
func anchorExistsInContent(
	content []byte,
	anchor string,
) bool {
	root, _ := Parse(content)
	if root == nil {
		return false
	}

	// Normalize anchor for case-insensitive comparison
	normalizedAnchor := strings.ToLower(
		strings.TrimSpace(anchor),
	)

	// Check for "Requirement: Name" pattern
	if strings.HasPrefix(
		normalizedAnchor,
		"requirement:",
	) {
		reqName := strings.TrimSpace(
			strings.TrimPrefix(
				normalizedAnchor,
				"requirement:",
			),
		)

		return requirementExists(root, reqName)
	}

	// Check for "Scenario: Name" pattern
	if strings.HasPrefix(
		normalizedAnchor,
		"scenario:",
	) {
		scenarioName := strings.TrimSpace(
			strings.TrimPrefix(
				normalizedAnchor,
				"scenario:",
			),
		)

		return scenarioExists(root, scenarioName)
	}

	// Generic header search - match any section/header containing anchor text
	return headerExists(root, normalizedAnchor)
}

// requirementExists checks if a requirement with the given name exists.
func requirementExists(
	root Node,
	name string,
) bool {
	normalizedName := strings.ToLower(
		name,
	) //nolint:revive // modifies-parameter
	found := false

	_ = Walk(root, &requirementFinder{
		targetName: normalizedName,
		found:      &found,
	})

	return found
}

// requirementFinder is a visitor that finds requirements by name.
type requirementFinder struct {
	BaseVisitor
	targetName string
	found      *bool
}

func (f *requirementFinder) VisitRequirement(
	n *NodeRequirement,
) error {
	if strings.ToLower(n.Name()) == f.targetName {
		*f.found = true

		return SkipChildren
	}

	return nil
}

// scenarioExists checks if a scenario with the given name exists.
func scenarioExists(root Node, name string) bool {
	normalizedName := strings.ToLower(
		name,
	) //nolint:revive // modifies-parameter
	found := false

	_ = Walk(root, &scenarioFinderForAnchor{
		targetName: normalizedName,
		found:      &found,
	})

	return found
}

// scenarioFinderForAnchor is a visitor that finds scenarios by name.
type scenarioFinderForAnchor struct {
	BaseVisitor
	targetName string
	found      *bool
}

func (f *scenarioFinderForAnchor) VisitScenario(
	n *NodeScenario,
) error {
	if strings.ToLower(n.Name()) == f.targetName {
		*f.found = true

		return SkipChildren
	}

	return nil
}

// headerExists checks if any header contains the given text.
func headerExists(root Node, text string) bool {
	text = strings.ToLower(text)
	found := false

	_ = Walk(root, &headerFinder{
		targetText: text,
		found:      &found,
	})

	return found
}

// headerFinder is a visitor that finds sections containing specific text.
type headerFinder struct {
	BaseVisitor
	targetText string
	found      *bool
}

func (f *headerFinder) VisitSection(
	n *NodeSection,
) error {
	title := strings.ToLower(string(n.Title()))
	if title == f.targetText ||
		strings.Contains(title, f.targetText) {
		*f.found = true

		return SkipChildren
	}

	return nil
}

func (f *headerFinder) VisitRequirement(
	n *NodeRequirement,
) error {
	name := strings.ToLower(n.Name())
	if name == f.targetText ||
		strings.Contains(name, f.targetText) {
		*f.found = true

		return SkipChildren
	}

	return nil
}

func (f *headerFinder) VisitScenario(
	n *NodeScenario,
) error {
	name := strings.ToLower(n.Name())
	if name == f.targetText ||
		strings.Contains(name, f.targetText) {
		*f.found = true

		return SkipChildren
	}

	return nil
}

// ValidateWikilinks checks all wikilinks in a parsed document and returns errors
// for any that cannot be resolved or have invalid anchors.
//
// Parameters:
//   - root: the root node of the parsed document
//   - source: the original source bytes (reserved for future use)
//   - projectRoot: the project root directory containing spectr/
//
// Returns a slice of WikilinkError for each invalid wikilink found.
func ValidateWikilinks(
	root Node,
	_ []byte,
	projectRoot string,
) []WikilinkError {
	if root == nil {
		return nil
	}

	validator := &wikilinkValidator{
		projectRoot: projectRoot,
		errors:      make([]WikilinkError, 0),
	}

	_ = Walk(root, validator)

	return validator.errors
}

// wikilinkValidator is a visitor that validates wikilinks.
type wikilinkValidator struct {
	BaseVisitor
	projectRoot string
	errors      []WikilinkError
}

func (v *wikilinkValidator) VisitWikilink(
	n *NodeWikilink,
) error {
	target := string(n.Target())
	display := string(n.Display())
	anchor := string(n.Anchor())
	start, _ := n.Span()

	// Resolve the wikilink target
	path, exists := ResolveWikilink(
		target,
		v.projectRoot,
	)

	if !exists {
		msg := "wikilink target not found: " + target + " (expected at " + path + ")" //nolint:revive // line-length-limit
		v.errors = append(v.errors, WikilinkError{
			Target:  target,
			Display: display,
			Anchor:  anchor,
			Offset:  start,
			Message: msg,
		})

		return nil
	}

	// If there's an anchor, validate it
	if anchor != "" {
		_, anchorValid, err := ResolveWikilinkWithAnchor(
			target,
			anchor,
			v.projectRoot,
		)
		if err != nil {
			v.errors = append(
				v.errors,
				WikilinkError{
					Target:  target,
					Display: display,
					Anchor:  anchor,
					Offset:  start,
					Message: "error reading target file: " + err.Error(),
				},
			)

			return nil
		}

		if !anchorValid {
			v.errors = append(
				v.errors,
				WikilinkError{
					Target:  target,
					Display: display,
					Anchor:  anchor,
					Offset:  start,
					Message: "anchor not found in target: #" + anchor,
				},
			)
		}
	}

	return nil
}

// ValidateWikilinkTarget checks if a single wikilink target is valid.
// This is a convenience function for validating individual targets.
func ValidateWikilinkTarget(
	target, projectRoot string,
) error {
	_, exists := ResolveWikilink(
		target,
		projectRoot,
	)
	if !exists {
		return WikilinkError{
			Target:  target,
			Offset:  -1,
			Message: "wikilink target not found: " + target,
		}
	}

	return nil
}

// GetWikilinkTargetType returns the type of target a wikilink refers to.
// Returns "spec", "change", or "unknown".
func GetWikilinkTargetType(
	target, projectRoot string,
) string {
	path, exists := ResolveWikilink(
		target,
		projectRoot,
	)
	if !exists {
		return "unknown"
	}

	if strings.Contains(
		path,
		filepath.Join("spectr", "specs"),
	) {
		return "spec"
	}
	if strings.Contains(
		path,
		filepath.Join("spectr", "changes"),
	) {
		return "change"
	}

	return "unknown"
}

// ListWikilinkTargets extracts all wikilink targets from content.
// Returns a slice of unique target strings.
func ListWikilinkTargets(
	content []byte,
) []string {
	wikilinks := ExtractWikilinks(content)

	// Use a map to deduplicate
	seen := make(map[string]bool)
	targets := make([]string, 0, len(wikilinks))

	for _, wl := range wikilinks {
		if !seen[wl.Target] {
			seen[wl.Target] = true
			targets = append(targets, wl.Target)
		}
	}

	return targets
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
