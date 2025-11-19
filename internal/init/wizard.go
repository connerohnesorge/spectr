//nolint:revive // file-length-limit,receiver-naming,unused-receiver,add-constant,early-return - UI code prioritizes readability
package init

//nolint:revive // line-length-limit,add-constant - readability over strict limits

//nolint:revive // file-length-limit, comments-density - UI code is cohesive

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	newline       = "\n"
	doubleNewline = "\n\n"
)

// WizardStep represents the current step in the wizard
type WizardStep int

const (
	StepIntro WizardStep = iota
	StepSelect
	StepReview
	StepExecute
	StepComplete
)

// WizardModel is the Bubbletea model for the init wizard
type WizardModel struct {
	step            WizardStep
	projectPath     string
	registry        *ToolRegistry
	selectedTools   map[string]bool // tool ID -> selected
	cursor          int             // cursor position in list
	executing       bool
	executionResult *ExecutionResult
	err             error
	allTools        []*ToolDefinition // sorted tools for display
}

// ExecutionResult holds the result of initialization
type ExecutionResult struct {
	CreatedFiles []string
	UpdatedFiles []string
	Errors       []string
}

// ExecutionCompleteMsg is sent when execution finishes
type ExecutionCompleteMsg struct {
	result *ExecutionResult
	err    error
}

// Lipgloss styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	dimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// NewWizardModel creates a new wizard model
func NewWizardModel(projectPath string) (*WizardModel, error) {
	registry := NewRegistry()
	allTools := registry.GetAllTools()

	// Sort tools by priority
	sort.Slice(allTools, func(i, j int) bool {
		return allTools[i].Priority < allTools[j].Priority
	})

	return &WizardModel{
		step:          StepIntro,
		projectPath:   projectPath,
		registry:      registry,
		selectedTools: make(map[string]bool),
		cursor:        0,
		allTools:      allTools,
	}, nil
}

// Init is the Bubbletea Init function
func (WizardModel) Init() tea.Cmd {
	return nil
}

// Update is the Bubbletea Update function
func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case StepIntro:
			return m.handleIntroKeys(typedMsg)
		case StepSelect:
			return m.handleSelectKeys(typedMsg)
		case StepReview:
			return m.handleReviewKeys(typedMsg)
		case StepExecute:
			// Execution in progress, no input
			return m, nil
		case StepComplete:
			return m.handleCompleteKeys(typedMsg)
		}
	case ExecutionCompleteMsg:
		m.executing = false
		m.executionResult = typedMsg.result
		m.err = typedMsg.err
		m.step = StepComplete

		return m, nil
	}

	return m, nil
}

// View is the Bubbletea View function
func (m WizardModel) View() string {
	switch m.step {
	case StepIntro:
		return m.renderIntro()
	case StepSelect:
		return m.renderSelect()
	case StepReview:
		return m.renderReview()
	case StepExecute:
		return m.renderExecute()
	case StepComplete:
		return m.renderComplete()
	}

	return ""
}

// ============================================================================
// Keyboard handlers
// ============================================================================

func (m WizardModel) handleIntroKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.step = StepSelect

		return m, nil
	}

	return m, nil
}

func (m WizardModel) handleSelectKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.allTools)-1 {
			m.cursor++
		}
	case " ":
		// Toggle selection
		if m.cursor < len(m.allTools) {
			tool := m.allTools[m.cursor]
			m.selectedTools[tool.ID] = !m.selectedTools[tool.ID]
		}
	case "enter":
		// Confirm and move to review
		m.step = StepReview

		return m, nil
	case "a":
		// Select all
		for _, tool := range m.allTools {
			m.selectedTools[tool.ID] = true
		}
	case "n":
		// Deselect all
		m.selectedTools = make(map[string]bool)
	}

	return m, nil
}

func (m WizardModel) handleReviewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "backspace", "esc":
		// Go back to selection
		m.step = StepSelect

		return m, nil
	case "enter":
		// Execute initialization
		m.step = StepExecute
		m.executing = true

		return m, executeInit(m.projectPath, m.getSelectedToolIDs())
	}

	return m, nil
}

func (m WizardModel) handleCompleteKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case keyQuit, keyCtrlC, keyEnter:
		return m, tea.Quit
	}

	return m, nil
}

// ============================================================================
// Render functions for each step
// ============================================================================

func (m WizardModel) renderIntro() string {
	var b strings.Builder

	// ASCII art banner
	b.WriteString(headerStyle.Render(asciiArt))
	b.WriteString(newlineDouble)

	// Welcome message
	b.WriteString(titleStyle.Render("Welcome to Spectr Initialization"))
	b.WriteString(newlineDouble)

	b.WriteString(
		"This wizard will help you initialize Spectr in " +
			"your project.\n\n",
	)
	b.WriteString("Spectr provides a structured approach to:\n")
	b.WriteString("  • Creating and managing change proposals\n")
	b.WriteString(
		"  • Documenting project architecture and " +
			"specifications\n",
	)
	b.WriteString("  • Integrating with AI coding assistants\n\n")

	b.WriteString(
		infoStyle.Render(
			fmt.Sprintf("Project path: %s", m.projectPath),
		),
	)
	b.WriteString(newlineDouble)

	// Instructions
	b.WriteString(
		subtleStyle.Render(
			"Press Enter to continue, or 'q' to quit" + newline,
		),
	)

	return b.String()
}

func (m WizardModel) renderSelect() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Select AI Tools to Configure"))
	b.WriteString(newlineDouble)

	b.WriteString(
		"Choose which AI coding tools you want to configure " +
			"with Spectr.\n",
	)
	b.WriteString("You can come back later to add more tools.\n\n")

	// Render all tools in a single flat list
	b.WriteString(m.renderToolGroup(m.allTools, 0))

	// Instructions
	b.WriteString(doubleNewline)
	b.WriteString(
		subtleStyle.Render(
			"↑/↓: Navigate  Space: Toggle  a: All  n: None  " +
				"Enter: Continue  q: Quit\n",
		),
	)

	return b.String()
}

func (m WizardModel) renderToolGroup(tools []*ToolDefinition, offset int) string {
	var b strings.Builder

	for i, tool := range tools {
		actualIndex := offset + i
		cursor := " "
		if m.cursor == actualIndex {
			cursor = cursorStyle.Render("▸")
		}

		checkbox := "[ ]"
		if m.selectedTools[tool.ID] {
			checkbox = selectedStyle.Render("[✓]")
		}

		line := fmt.Sprintf("  %s %s %s", cursor, checkbox, tool.Name)

		switch {
		case m.cursor == actualIndex:
			b.WriteString(cursorStyle.Render(line))
		case m.selectedTools[tool.ID]:
			b.WriteString(selectedStyle.Render(line))
		default:
			b.WriteString(dimmedStyle.Render(line))
		}

		b.WriteString("\n")
	}

	return b.String()
}

func (m WizardModel) renderReview() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Review Your Selections"))
	b.WriteString("\n\n")

	selectedCount := len(m.getSelectedToolIDs())
	m.renderSelectedTools(&b, selectedCount)
	m.renderCreationPlan(&b, selectedCount)

	b.WriteString(newline)
	b.WriteString(
		subtleStyle.Render(
			"Press Enter to initialize, Backspace to go back, " +
				"or 'q' to quit\n",
		),
	)

	return b.String()
}

// renderSelectedTools displays the selected tools or warning if none
func (m WizardModel) renderSelectedTools(
	b *strings.Builder,
	count int,
) {
	if count == 0 {
		b.WriteString(errorStyle.Render("⚠ No tools selected"))
		b.WriteString(doubleNewline)
		b.WriteString("You haven't selected any tools to configure.\n")
		b.WriteString(
			"Spectr will still be initialized, but no tool " +
				"integrations will be set up.\n\n",
		)

		return
	}

	fmt.Fprintf(b,
		"You have selected %d tool(s) to configure:\n\n",
		count)

	for _, tool := range m.allTools {
		if !m.selectedTools[tool.ID] {
			continue
		}
		b.WriteString(successStyle.Render("  ✓ "))
		b.WriteString(tool.Name)
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

// renderCreationPlan displays what files will be created
func (m WizardModel) renderCreationPlan(b *strings.Builder, count int) {
	b.WriteString("The following will be created:\n")
	b.WriteString(infoStyle.Render("  • spectr/project.md"))
	b.WriteString(" - Project documentation template" + newline)
	b.WriteString(infoStyle.Render("  • spectr/AGENTS.md"))
	b.WriteString(" - AI agent instructions" + newline)

	if count > 0 {
		b.WriteString(infoStyle.Render(fmt.Sprintf(
			"  • Tool configurations for %d selected tools",
			count,
		)))
		b.WriteString(newline)
	}
}

func (WizardModel) renderExecute() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Initializing Spectr..."))
	b.WriteString(doubleNewline)

	b.WriteString(infoStyle.Render("⏳ Setting up your project..."))
	b.WriteString(doubleNewline)

	b.WriteString("This will only take a moment." + newline)

	return b.String()
}

func (m WizardModel) renderComplete() string {
	var b strings.Builder

	if m.err != nil {
		m.renderError(&b)

		return b.String()
	}

	b.WriteString(successStyle.Render("✓ Spectr Initialized Successfully!"))
	b.WriteString("\n\n")

	if m.executionResult != nil {
		m.renderExecutionResults(&b)
	}

	b.WriteString(FormatNextStepsMessage())
	b.WriteString(newline)
	b.WriteString(subtleStyle.Render("Press 'q' to quit" + newline))

	return b.String()
}

// renderError displays initialization errors
func (m WizardModel) renderError(b *strings.Builder) {
	b.WriteString(errorStyle.Render("✗ Initialization Failed"))
	b.WriteString(doubleNewline)
	b.WriteString(errorStyle.Render(m.err.Error()))
	b.WriteString(doubleNewline)

	if m.executionResult != nil && len(m.executionResult.Errors) > 0 {
		b.WriteString("Errors:" + newline)
		for _, err := range m.executionResult.Errors {
			b.WriteString(errorStyle.Render("  • "))
			b.WriteString(err)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}

	b.WriteString(subtleStyle.Render("Press 'q' to quit\n"))
}

// renderExecutionResults displays created/updated files and warnings
func (m WizardModel) renderExecutionResults(b *strings.Builder) {
	if len(m.executionResult.CreatedFiles) > 0 {
		b.WriteString(successStyle.Render("Created files:"))
		b.WriteString(newline)
		for _, file := range m.executionResult.CreatedFiles {
			b.WriteString(infoStyle.Render("  ✓ "))
			b.WriteString(file)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}

	if len(m.executionResult.UpdatedFiles) > 0 {
		b.WriteString(successStyle.Render("Updated files:"))
		b.WriteString(newline)
		for _, file := range m.executionResult.UpdatedFiles {
			b.WriteString(infoStyle.Render("  ↻ "))
			b.WriteString(file)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}

	if len(m.executionResult.Errors) > 0 {
		b.WriteString(errorStyle.Render("Warnings:"))
		b.WriteString(newline)
		for _, err := range m.executionResult.Errors {
			b.WriteString(errorStyle.Render("  ⚠ "))
			b.WriteString(err)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}
}

// ============================================================================
// Helper functions
// ============================================================================

func (m WizardModel) getSelectedToolIDs() []string {
	var selected []string
	for id, isSelected := range m.selectedTools {
		if isSelected {
			selected = append(selected, id)
		}
	}

	return selected
}

// executeInit runs the initialization and sends result
func executeInit(projectPath string, selectedTools []string) tea.Cmd {
	return func() tea.Msg {
		executor, err := NewInitExecutor(projectPath)
		if err != nil {
			return ExecutionCompleteMsg{
				result: nil,
				err:    fmt.Errorf("failed to create executor: %w", err),
			}
		}

		result, err := executor.Execute(selectedTools)

		return ExecutionCompleteMsg{
			result: result,
			err:    err,
		}
	}
}

// GetError returns the error from the wizard (if any)
func (m WizardModel) GetError() error {
	return m.err
}

// ASCII art for Spectr branding
const asciiArt = `
███████ ██████  ███████  ██████ ███████ ████████
██      ██   ██ ██      ██         █    ██    ██
███████ ██████  █████   ██         █    ████████
     ██ ██      ██      ██         █    ██  ██
███████ ██      ███████  ██████    █    ██    ██
`
