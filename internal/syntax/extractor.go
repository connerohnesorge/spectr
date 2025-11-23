package syntax

import (
	"strings"
)

// ExtractedRequirement represents a requirement extracted from the AST.
type ExtractedRequirement struct {
	Name      string
	Raw       string
	Scenarios []string
}

// ExtractRequirements extracts requirements and their scenarios from the AST.
func ExtractRequirements(doc *Document) []ExtractedRequirement {
	var reqs []ExtractedRequirement
	var currReq *ExtractedRequirement

	for _, node := range doc.Nodes {
		switch n := node.(type) {
		case *Header:
			if n.Level == 3 && strings.HasPrefix(n.Text, "Requirement:") {
				// Start new requirement
				if currReq != nil {
					reqs = append(reqs, *currReq)
				}
				name := strings.TrimSpace(strings.TrimPrefix(n.Text, "Requirement:"))
				currReq = &ExtractedRequirement{
					Name: name,
					Raw:  "",
				}
				currReq.Raw += n.Raw
			} else if n.Level == 2 {
				// Level 2 header ends a requirement section
				if currReq != nil {
					reqs = append(reqs, *currReq)
					currReq = nil
				}
			} else if n.Level == 4 && strings.HasPrefix(n.Text, "Scenario:") && currReq != nil {
				// Scenario inside requirement
				name := strings.TrimSpace(strings.TrimPrefix(n.Text, "Scenario:"))
				currReq.Scenarios = append(currReq.Scenarios, name)
				currReq.Raw += n.Raw
			} else if currReq != nil {
				currReq.Raw += n.Raw
			}
		case *Text:
			if currReq != nil {
				currReq.Raw += n.Raw
			}
		case *CodeBlock:
			if currReq != nil {
				currReq.Raw += n.Raw
			}
		}
	}

	if currReq != nil {
		reqs = append(reqs, *currReq)
	}
	return reqs
}

// RenameOp represents a requirement rename operation.
type RenameOp struct {
	From string
	To   string
}

// DeltaPlan represents all delta operations for a spec.
type DeltaPlan struct {
	Added    []ExtractedRequirement
	Modified []ExtractedRequirement
	Removed  []string
	Renamed  []RenameOp
}

// ExtractDelta extracts delta operations from the AST.
func ExtractDelta(doc *Document) *DeltaPlan {
	plan := &DeltaPlan{
		Added:    make([]ExtractedRequirement, 0),
		Modified: make([]ExtractedRequirement, 0),
		Removed:  make([]string, 0),
		Renamed:  make([]RenameOp, 0),
	}

	var currentSection string // "ADDED", "MODIFIED", "REMOVED", "RENAMED"
	var currentReq *ExtractedRequirement
	var currentFrom string

	for _, node := range doc.Nodes {
		switch n := node.(type) {
		case *Header:
			if n.Level == 2 {
				// Flush currentReq before switching sections
				if currentReq != nil {
					if currentSection == "ADDED" {
						plan.Added = append(plan.Added, *currentReq)
					} else if currentSection == "MODIFIED" {
						plan.Modified = append(plan.Modified, *currentReq)
					}
					currentReq = nil
				}

				// Section Header
				text := strings.TrimSpace(n.Text)
				if strings.Contains(text, "ADDED Requirements") {
					currentSection = "ADDED"
				} else if strings.Contains(text, "MODIFIED Requirements") {
					currentSection = "MODIFIED"
				} else if strings.Contains(text, "REMOVED Requirements") {
					currentSection = "REMOVED"
				} else if strings.Contains(text, "RENAMED Requirements") {
					currentSection = "RENAMED"
				} else {
					currentSection = ""
				}

				currentFrom = ""

			} else if n.Level == 3 && strings.HasPrefix(n.Text, "Requirement:") {
				// Flush previous req
				if currentReq != nil {
					if currentSection == "ADDED" {
						plan.Added = append(plan.Added, *currentReq)
					} else if currentSection == "MODIFIED" {
						plan.Modified = append(plan.Modified, *currentReq)
					}
				}

				if currentSection == "ADDED" || currentSection == "MODIFIED" {
					name := strings.TrimSpace(strings.TrimPrefix(n.Text, "Requirement:"))
					currentReq = &ExtractedRequirement{
						Name: name,
						Raw:  "",
					}
					currentReq.Raw += n.Raw
				} else if currentSection == "REMOVED" {
					name := strings.TrimSpace(strings.TrimPrefix(n.Text, "Requirement:"))
					plan.Removed = append(plan.Removed, name)
					currentReq = nil // REMOVED reqs don't have body we care about?
					// Actually, parseRemovedSection in original code just extracts names.
				} else {
					currentReq = nil
				}
			} else if n.Level == 4 && strings.HasPrefix(n.Text, "Scenario:") && currentReq != nil {
				name := strings.TrimSpace(strings.TrimPrefix(n.Text, "Scenario:"))
				currentReq.Scenarios = append(currentReq.Scenarios, name)
				currentReq.Raw += n.Raw
			} else if currentReq != nil {
				currentReq.Raw += n.Raw
			}

		case *List:
			if currentSection == "RENAMED" {
				content := n.Content
				if strings.Contains(content, "FROM:") {
					start := strings.Index(content, "Requirement:")
					if start != -1 {
						name := strings.TrimSpace(content[start+len("Requirement:"):])
						name = strings.Trim(name, "`")
						currentFrom = name
					}
				} else if strings.Contains(content, "TO:") && currentFrom != "" {
					start := strings.Index(content, "Requirement:")
					if start != -1 {
						name := strings.TrimSpace(content[start+len("Requirement:"):])
						name = strings.Trim(name, "`")
						plan.Renamed = append(plan.Renamed, RenameOp{
							From: currentFrom,
							To:   name,
						})
						currentFrom = ""
					}
				}
			} else if currentReq != nil {
				currentReq.Raw += n.Raw
			}

		case *Text:
			if currentReq != nil {
				currentReq.Raw += n.Raw
			}
		case *CodeBlock:
			if currentReq != nil {
				currentReq.Raw += n.Raw
			}
		}
	}

	// Flush last req
	if currentReq != nil {
		if currentSection == "ADDED" {
			plan.Added = append(plan.Added, *currentReq)
		} else if currentSection == "MODIFIED" {
			plan.Modified = append(plan.Modified, *currentReq)
		}
	}

	return plan
}
