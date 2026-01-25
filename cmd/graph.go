// Package cmd provides command-line interface implementations.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/validation"
)

// GraphCmd represents the graph command for visualizing proposal dependencies.
type GraphCmd struct {
	// ChangeID is the optional change ID to show the graph for.
	// If not specified, shows all proposals.
	//nolint:lll // Kong struct tag requires long line for help text
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Show graph for specific change (optional)"`

	// Dot outputs in Graphviz DOT format
	Dot bool `name:"dot" help:"Output in Graphviz DOT format"`

	// JSON outputs in JSON format
	JSON bool `name:"json" help:"Output in JSON format"`
}

// GraphNode represents a proposal in the dependency graph.
type GraphNode struct {
	ID       string      `json:"id"`
	Status   string      `json:"status"` // "archived", "active", "unknown"
	Requires []GraphEdge `json:"requires,omitempty"`
	Enables  []GraphEdge `json:"enables,omitempty"`
}

// GraphEdge represents a dependency relationship.
type GraphEdge struct {
	ID     string `json:"id"`
	Reason string `json:"reason,omitempty"`
	Status string `json:"status"` // "archived", "active", "unknown"
}

// GraphOutput is the JSON output structure.
type GraphOutput struct {
	Nodes []GraphNode `json:"nodes"`
}

// Run executes the graph command.
func (c *GraphCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Build the dependency graph
	graph, err := validation.BuildDependencyGraph(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// If no proposals have dependencies, show a message
	if len(graph.Nodes) == 0 {
		fmt.Println("No proposals found")

		return nil
	}

	// Filter to specific change if requested
	if c.ChangeID != "" {
		result, resolveErr := discovery.ResolveChangeID(c.ChangeID, projectRoot)
		if resolveErr != nil {
			return resolveErr
		}
		c.ChangeID = result.ChangeID
	}

	// Output in requested format
	switch {
	case c.Dot:
		return c.outputDot(graph, projectRoot)
	case c.JSON:
		return c.outputJSON(graph, projectRoot)
	default:
		return c.outputASCII(graph, projectRoot)
	}
}

// outputASCII outputs the dependency graph in ASCII tree format.
//
//nolint:revive // cognitive-complexity: graph rendering requires nested iteration
func (c *GraphCmd) outputASCII(
	graph *validation.DependencyGraph,
	projectRoot string,
) error {
	changeIDs := c.getFilteredChangeIDs(graph)
	sort.Strings(changeIDs)

	if len(changeIDs) == 0 {
		if c.ChangeID != "" {
			fmt.Printf("Change '%s' not found in dependency graph\n", c.ChangeID)
		}

		return nil
	}

	for i, id := range changeIDs {
		c.printProposalASCII(id, graph, projectRoot)

		// Add blank line between entries
		if i < len(changeIDs)-1 {
			fmt.Println()
		}
	}

	return nil
}

// getFilteredChangeIDs returns the list of change IDs to display.
func (c *GraphCmd) getFilteredChangeIDs(
	graph *validation.DependencyGraph,
) []string {
	changeIDs := make([]string, 0, len(graph.Nodes))
	for id := range graph.Nodes {
		if c.ChangeID == "" || id == c.ChangeID {
			changeIDs = append(changeIDs, id)
		}
	}

	return changeIDs
}

// printProposalASCII prints a single proposal in ASCII format.
//
//nolint:revive // cognitive-complexity: graph rendering requires nested iteration for dependencies
func (*GraphCmd) printProposalASCII(
	id string,
	graph *validation.DependencyGraph,
	projectRoot string,
) {
	meta := graph.Nodes[id]

	// Print change ID with status
	status := getStatusSymbol(discovery.ChangeStatusActive)
	fmt.Printf("%s (%s)\n", id, status)

	// Print requires
	if meta != nil && len(meta.Requires) > 0 {
		for j, dep := range meta.Requires {
			depStatus, _ := discovery.GetChangeStatus(dep.ID, projectRoot)
			symbol := getStatusSymbol(depStatus)

			prefix := "├──"
			if j == len(meta.Requires)-1 && len(meta.Enables) == 0 {
				prefix = "└──"
			}

			reason := ""
			if dep.Reason != "" {
				reason = fmt.Sprintf(" (%s)", dep.Reason)
			}
			fmt.Printf("%s requires: %s %s%s\n", prefix, dep.ID, symbol, reason)
		}
	}

	// Print enables
	if meta != nil && len(meta.Enables) > 0 {
		for j, dep := range meta.Enables {
			depStatus, _ := discovery.GetChangeStatus(dep.ID, projectRoot)
			symbol := getStatusSymbol(depStatus)

			prefix := "├──"
			if j == len(meta.Enables)-1 {
				prefix = "└──"
			}

			reason := ""
			if dep.Reason != "" {
				reason = fmt.Sprintf(" (%s)", dep.Reason)
			}
			fmt.Printf("%s enables: %s %s%s\n", prefix, dep.ID, symbol, reason)
		}
	}
}

// outputDot outputs the dependency graph in Graphviz DOT format.
func (c *GraphCmd) outputDot(
	graph *validation.DependencyGraph,
	projectRoot string,
) error {
	var sb strings.Builder

	sb.WriteString("digraph proposals {\n")
	sb.WriteString("  rankdir=BT;\n")
	sb.WriteString("  node [shape=box];\n")
	sb.WriteString("\n")

	// Collect all nodes and their statuses
	allNodes := c.collectAllNodes(graph, projectRoot)

	// Write node definitions with styling
	sortedNodes := make([]string, 0, len(allNodes))
	for id := range allNodes {
		sortedNodes = append(sortedNodes, id)
	}
	sort.Strings(sortedNodes)

	for _, id := range sortedNodes {
		status := allNodes[id]
		style := getNodeStyle(status)
		sb.WriteString(fmt.Sprintf("  %q%s;\n", id, style))
	}

	sb.WriteString("\n")

	// Write edges
	c.writeEdges(&sb, graph)

	sb.WriteString("}\n")

	fmt.Print(sb.String())

	return nil
}

// collectAllNodes collects all nodes and their statuses for DOT output.
func (c *GraphCmd) collectAllNodes(
	graph *validation.DependencyGraph,
	projectRoot string,
) map[string]discovery.ChangeStatus {
	allNodes := make(map[string]discovery.ChangeStatus)

	for id := range graph.Nodes {
		if c.ChangeID != "" && id != c.ChangeID {
			continue
		}
		status, _ := discovery.GetChangeStatus(id, projectRoot)
		allNodes[id] = status

		meta := graph.Nodes[id]
		if meta != nil {
			for _, dep := range meta.Requires {
				depStatus, _ := discovery.GetChangeStatus(dep.ID, projectRoot)
				allNodes[dep.ID] = depStatus
			}
			for _, dep := range meta.Enables {
				depStatus, _ := discovery.GetChangeStatus(dep.ID, projectRoot)
				allNodes[dep.ID] = depStatus
			}
		}
	}

	return allNodes
}

// getNodeStyle returns the DOT style for a node based on its status.
func getNodeStyle(status discovery.ChangeStatus) string {
	switch status {
	case discovery.ChangeStatusArchived:
		return " [style=filled, fillcolor=lightgreen]"
	case discovery.ChangeStatusActive:
		return " [style=filled, fillcolor=lightyellow]"
	case discovery.ChangeStatusUnknown:
		return " [style=dashed]"
	default:
		return " [style=dashed]"
	}
}

// writeEdges writes the edge definitions for DOT output.
func (c *GraphCmd) writeEdges(
	sb *strings.Builder,
	graph *validation.DependencyGraph,
) {
	for id, meta := range graph.Nodes {
		if c.ChangeID != "" && id != c.ChangeID {
			continue
		}
		if meta == nil {
			continue
		}

		for _, dep := range meta.Requires {
			label := "requires"
			if dep.Reason != "" {
				label = dep.Reason
			}
			_, _ = fmt.Fprintf(sb, "  %q -> %q [label=%q];\n", dep.ID, id, label)
		}

		for _, dep := range meta.Enables {
			label := "enables"
			if dep.Reason != "" {
				label = dep.Reason
			}
			_, _ = fmt.Fprintf(sb, "  %q -> %q [label=%q, style=dashed];\n", id, dep.ID, label)
		}
	}
}

// outputJSON outputs the dependency graph in JSON format.
func (c *GraphCmd) outputJSON(
	graph *validation.DependencyGraph,
	projectRoot string,
) error {
	output := GraphOutput{
		Nodes: make([]GraphNode, 0),
	}

	changeIDs := c.getFilteredChangeIDs(graph)
	sort.Strings(changeIDs)

	for _, id := range changeIDs {
		meta := graph.Nodes[id]
		status, _ := discovery.GetChangeStatus(id, projectRoot)

		node := GraphNode{
			ID:       id,
			Status:   status.String(),
			Requires: make([]GraphEdge, 0),
			Enables:  make([]GraphEdge, 0),
		}

		if meta != nil {
			for _, dep := range meta.Requires {
				depStatus, _ := discovery.GetChangeStatus(dep.ID, projectRoot)
				node.Requires = append(node.Requires, GraphEdge{
					ID:     dep.ID,
					Reason: dep.Reason,
					Status: depStatus.String(),
				})
			}

			for _, dep := range meta.Enables {
				depStatus, _ := discovery.GetChangeStatus(dep.ID, projectRoot)
				node.Enables = append(node.Enables, GraphEdge{
					ID:     dep.ID,
					Reason: dep.Reason,
					Status: depStatus.String(),
				})
			}
		}

		output.Nodes = append(output.Nodes, node)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return encoder.Encode(output)
}

// getStatusSymbol returns a symbol indicating the status of a change.
func getStatusSymbol(status discovery.ChangeStatus) string {
	switch status {
	case discovery.ChangeStatusArchived:
		return "✓"
	case discovery.ChangeStatusActive:
		return "⧖"
	case discovery.ChangeStatusUnknown:
		return "?"
	default:
		return "?"
	}
}
