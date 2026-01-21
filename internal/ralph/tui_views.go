// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// tui_views.go contains view rendering functions for the Bubble Tea TUI.
// It separates presentation logic from the TUI model to improve
// maintainability and testability.
//
// View components:
// - Task list rendering with status icons and colors
// - Agent output pane with scrolling and ANSI support
// - Help bar with context-sensitive commands
// - Header with change information
// - Footer with status summary
//
// All views use Lipgloss for styling and layout.

const (
	// defaultInitialCapacity is the default initial capacity for slices
	defaultInitialCapacity = 5

	// maxVisibleTasks is the maximum number of tasks to show in the task list
	maxVisibleTasks = 8

	// maxTerminalWidth is the maximum terminal width to use for rendering
	maxTerminalWidth = 80

	// newlineString is the newline string used in rendering
	newlineString = "\n"

	// quitInstructions is the quit instruction shown to users
	quitInstructions = "[q] quit"
)

// Styling constants
var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1)

	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("229"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Bold(true)
)

// renderNormalView renders the main orchestration view
func (m *TUIModel) renderNormalView() string {
	sections := []string{
		m.renderHeader(),
		m.renderTaskList(),
		m.renderAgentOutput(),
		m.renderHelpBar(),
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderInteractiveView renders the interactive task selection view
func (m *TUIModel) renderInteractiveView() string {
	header := headerStyle.Render("Interactive Task Selection")
	instructions := helpStyle.Render(
		"Select tasks to run. Space to toggle, Enter to confirm, Esc to cancel.",
	)

	sections := make([]string, 0, defaultInitialCapacity)
	sections = append(sections, header, instructions, "")

	// Task list with selection checkboxes
	for i, task := range m.tasks {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		checkbox := "[ ]"
		if m.selectedTasks[task.ID] {
			checkbox = "[x]"
		}

		icon := getStatusIcon(task.Status)
		line := fmt.Sprintf("%s%s %s %s - %s", cursor, checkbox, icon, task.ID, task.Section)

		if i == m.cursor {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Bold(true).
				Render(line)
		}

		sections = append(sections, line)
	}

	// Help
	sections = append(sections, "")
	help := helpStyle.Render("[↑↓] navigate  [space] select  [enter] confirm  [esc] cancel")
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderHeader renders the header with change ID and progress
func (m *TUIModel) renderHeader() string {
	completed := m.GetCompletedCount()
	total := len(m.tasks)
	progress := fmt.Sprintf("[%d/%d tasks]", completed, total)

	title := fmt.Sprintf("spectr ralph: %s", m.changeID)

	// Add mode indicator
	modeIndicator := ""
	switch m.mode {
	case ModePaused:
		modeIndicator = warningStyle.Render(" [PAUSED]")
	case ModeFailure:
		modeIndicator = errorStyle.Render(" [FAILED]")
	case ModeInteractive:
		modeIndicator = " [INTERACTIVE]"
	case ModeNormal:
		// No indicator for normal mode
	}

	header := fmt.Sprintf("%s%s %s", title, modeIndicator, progress)

	// Create separator line
	separator := strings.Repeat("─", minInt(m.width, maxTerminalWidth))

	return headerStyle.Render(header) + newlineString + separator
}

// renderTaskList renders the task list with status indicators
//
//nolint:revive // View rendering function complexity is intentional for comprehensive layout logic
func (m *TUIModel) renderTaskList() string {
	if len(m.tasks) == 0 {
		return sectionTitleStyle.Render(
			"Tasks",
		) + newlineString + helpStyle.Render(
			"No tasks found",
		)
	}

	lines := []string{sectionTitleStyle.Render("Tasks")}
	startIdx := 0
	currentIdx := -1

	// Find current task index
	for i, task := range m.tasks {
		if task.ID == m.currentTaskID {
			currentIdx = i

			break
		}
	}

	// Calculate visible window
	if currentIdx >= 0 {
		startIdx = currentIdx - 2
		if startIdx < 0 {
			startIdx = 0
		}
		if startIdx+maxVisibleTasks > len(m.tasks) {
			startIdx = len(m.tasks) - maxVisibleTasks
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}

	endIdx := startIdx + maxVisibleTasks
	if endIdx > len(m.tasks) {
		endIdx = len(m.tasks)
	}

	// Show tasks in window
	for i := startIdx; i < endIdx; i++ {
		task := m.tasks[i]
		isCurrent := task.ID == m.currentTaskID
		line := FormatTaskLine(task, isCurrent)
		lines = append(lines, "  "+line)
	}

	// Show indicator if there are more tasks
	if endIdx < len(m.tasks) {
		remaining := len(m.tasks) - endIdx
		lines = append(lines, helpStyle.Render(fmt.Sprintf("  ... %d more tasks", remaining)))
	}

	content := strings.Join(lines, newlineString)

	// Add border
	bordered := borderStyle.
		Width(minInt(m.width-2, 78)).
		Render(content)

	return bordered
}

// renderAgentOutput renders the agent output pane
func (m *TUIModel) renderAgentOutput() string {
	title := "Agent Output"
	if m.currentTaskID != "" {
		title = fmt.Sprintf("Agent Output (task %s)", m.currentTaskID)
	}

	header := sectionTitleStyle.Render(title)

	// Render viewport content
	content := m.output.View()

	// Add border around output
	bordered := borderStyle.
		Width(minInt(m.width-2, 78)).
		Height(m.getOutputHeight()).
		Render(content)

	return header + newlineString + bordered
}

// renderHelpBar renders the help bar with available commands
func (m *TUIModel) renderHelpBar() string {
	var commands []string

	switch m.mode {
	case ModeFailure:
		commands = []string{
			"[r] retry",
			"[s] skip",
			quitInstructions,
		}

	case ModePaused:
		commands = []string{
			"[p] resume",
			quitInstructions,
		}

	case ModeInteractive:
		commands = []string{
			"[↑↓/jk] navigate",
			"[space] select",
			"[enter] confirm",
			"[esc] cancel",
		}

	case ModeNormal:
		commands = []string{
			quitInstructions,
			"[p] pause",
			"[i] interactive",
		}

	default:
		commands = []string{
			quitInstructions,
		}
	}

	help := strings.Join(commands, "  ")

	// Create separator line
	separator := strings.Repeat("─", minInt(m.width, maxTerminalWidth))

	return separator + newlineString + helpStyle.Render(help)
}

// RenderTaskSummary renders a summary of task completion
func RenderTaskSummary(tasks []*Task) string {
	completed := 0
	failed := 0
	pending := 0

	for _, task := range tasks {
		switch task.Status {
		case StatusCompleted:
			completed++
		case StatusFailed:
			failed++
		default:
			pending++
		}
	}

	total := len(tasks)

	lines := []string{
		headerStyle.Render("Task Summary"),
		"",
		successStyle.Render(fmt.Sprintf("Completed: %d", completed)),
	}

	if failed > 0 {
		lines = append(lines, errorStyle.Render(fmt.Sprintf("Failed: %d", failed)))
	}

	if pending > 0 {
		lines = append(lines, helpStyle.Render(fmt.Sprintf("Pending: %d", pending)))
	}

	lines = append(lines, "", fmt.Sprintf("Total: %d", total))

	return strings.Join(lines, newlineString)
}
