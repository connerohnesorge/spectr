//nolint:revive // file-length-limit,receiver-naming,unused-receiver,add-constant,early-return - UI code prioritizes readability
package initialize

//nolint:revive // line-length-limit,add-constant - readability over strict limits

//nolint:revive // file-length-limit, comments-density - UI code is cohesive

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
	"github.com/connerohnesorge/spectr/internal/tui"
	"github.com/spf13/afero"
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
	step                 WizardStep
	projectPath          string
	selectedProviders    map[string]bool // provider ID -> selected
	configuredProviders  map[string]bool // provider ID -> is configured
	cursor               int             // cursor position in list
	executing            bool
	executionResult      *ExecutionResult
	err                  error
	allProviders         []providers.Registration // sorted providers for display
	ciWorkflowEnabled    bool                     // whether user wants CI workflow created
	ciWorkflowConfigured bool                     // whether .github/workflows/spectr-ci.yml already exists
	// Search mode state
	searchMode        bool                     // whether search mode is active
	searchQuery       string                   // current search query
	searchInput       textinput.Model          // text input for search
	filteredProviders []providers.Registration // providers matching search query
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
func NewWizardModel(
	cmd *InitCmd,
) (*WizardModel, error) {
	// Use the resolved path from InitCmd
	projectPath := cmd.Path
	if projectPath == "" {
		return nil, fmt.Errorf(
			"project path is required",
		)
	}

	allProviders := providers.All()

	projectFs := afero.NewBasePathFs(
		afero.NewOsFs(),
		projectPath,
	)
	globalFs := afero.NewOsFs()
	cfg := &types.Config{SpectrDir: "spectr"}

	configuredProviders := make(map[string]bool)
	selectedProviders := make(map[string]bool)

	for _, reg := range allProviders {
		isConfigured := true
		inits := reg.Provider.Initializers()
		if len(inits) == 0 {
			isConfigured = false
		} else {
			for _, ini := range inits {
				setup, err := ini.IsSetup(projectFs, globalFs, cfg)
				if err != nil || !setup {
					isConfigured = false

					break
				}
			}
		}

		configuredProviders[reg.ID] = isConfigured

		// Pre-select already-configured providers
		if isConfigured {
			selectedProviders[reg.ID] = true
		}
	}

	// Detect if CI workflow is already configured
	ciWorkflowPath := filepath.Join(
		projectPath,
		".github",
		"workflows",
		"spectr-ci.yml",
	)
	ciWorkflowConfigured := FileExists(
		ciWorkflowPath,
	)

	// Pre-select CI workflow if already configured
	ciWorkflowEnabled := ciWorkflowConfigured

	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Type to search..."
	searchInput.CharLimit = 50
	searchInput.Width = 30

	return &WizardModel{
		step:                 StepIntro,
		projectPath:          projectPath,
		selectedProviders:    selectedProviders,
		configuredProviders:  configuredProviders,
		cursor:               0,
		allProviders:         allProviders,
		ciWorkflowEnabled:    ciWorkflowEnabled,
		ciWorkflowConfigured: ciWorkflowConfigured,
		searchInput:          searchInput,
		// Initially show all providers
		filteredProviders: allProviders,
	}, nil
}

// Init is the Bubbletea Init function
func (WizardModel) Init() tea.Cmd {
	return nil
}

// Update is the Bubbletea Update function
func (m WizardModel) Update(
	msg tea.Msg,
) (tea.Model, tea.Cmd) {
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

func (m WizardModel) handleIntroKeys(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.step = StepSelect

		return m, nil
	}

	return m, nil
}

func (m WizardModel) handleSelectKeys(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	// Handle search mode input
	if m.searchMode {
		return m.handleSearchModeInput(msg)
	}

	// Normal mode key handling
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.filteredProviders)-1 {
			m.cursor++
		}
	case " ":
		// Toggle selection on filtered list
		if m.cursor < len(m.filteredProviders) {
			provider := m.filteredProviders[m.cursor]
			m.selectedProviders[provider.ID] = !m.selectedProviders[provider.ID]
		}
	case "enter":
		// Confirm and move to review
		m.step = StepReview

		return m, nil
	case "a":
		// Select all (from full list, not just filtered)
		for _, provider := range m.allProviders {
			m.selectedProviders[provider.ID] = true
		}
	case "n":
		// Deselect all
		m.selectedProviders = make(
			map[string]bool,
		)
	case "/":
		// Enter search mode
		m.searchMode = true
		m.searchInput.Focus()

		return m, nil
	}

	return m, nil
}

// handleSearchModeInput handles keyboard input when in search mode
func (m WizardModel) handleSearchModeInput(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	//nolint:exhaustive // Only handling specific keys, default handles text input
	switch msg.Type {
	case tea.KeyEsc:
		// Exit search mode, clear query and restore all providers
		m.searchMode = false
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m.filteredProviders = m.allProviders
		m.cursor = 0

		return m, nil
	case tea.KeyEnter:
		// Exit search mode but keep filter applied, proceed to review
		m.searchMode = false
		m.step = StepReview

		return m, nil
	case tea.KeyUp:
		// Allow navigation while searching
		if m.cursor > 0 {
			m.cursor--
		}

		return m, nil
	case tea.KeyDown:
		// Allow navigation while searching
		if m.cursor < len(m.filteredProviders)-1 {
			m.cursor++
		}

		return m, nil
	case tea.KeySpace:
		// Toggle selection on filtered list while in search mode
		if m.cursor < len(m.filteredProviders) {
			provider := m.filteredProviders[m.cursor]
			m.selectedProviders[provider.ID] = !m.selectedProviders[provider.ID]
		}

		return m, nil
	default:
		// Handle text input for search
		m.searchInput, cmd = m.searchInput.Update(
			msg,
		)
		m.searchQuery = m.searchInput.Value()
		m.applyProviderFilter()

		return m, cmd
	}
}

// applyProviderFilter filters providers based on the current search query
func (m *WizardModel) applyProviderFilter() {
	query := strings.ToLower(m.searchQuery)

	if query == "" {
		m.filteredProviders = m.allProviders
	} else {
		m.filteredProviders = make([]providers.Registration, 0)
		for _, provider := range m.allProviders {
			if strings.Contains(strings.ToLower(provider.Name), query) {
				m.filteredProviders = append(m.filteredProviders, provider)
			}
		}
	}

	// Adjust cursor position to stay within bounds
	if len(m.filteredProviders) > 0 {
		if m.cursor >= len(m.filteredProviders) {
			m.cursor = len(
				m.filteredProviders,
			) - 1
		}
	} else {
		m.cursor = 0
	}
}

func (m WizardModel) handleReviewKeys(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "backspace", "esc":
		// Go back to selection
		m.step = StepSelect

		return m, nil
	case " ":
		// Toggle CI workflow option
		m.ciWorkflowEnabled = !m.ciWorkflowEnabled

		return m, nil
	case "enter":
		// Execute initialization
		m.step = StepExecute
		m.executing = true

		return m, executeInit(
			m.projectPath,
			m.getSelectedProviderIDs(),
			m.ciWorkflowEnabled,
		)
	}

	return m, nil
}

func (m WizardModel) handleCompleteKeys(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case keyQuit, keyCtrlC, keyEnter:
		return m, tea.Quit
	case keyCopy:
		// Only allow copy on success screen (no init error)
		if m.err == nil {
			// Copy the populate context prompt to clipboard
			// CopyToClipboard uses OSC 52 fallback, so it never fails
			_ = tui.CopyToClipboard(
				PopulateContextPrompt,
			)

			return m, tea.Quit
		}
	}

	return m, nil
}

func (m WizardModel) renderIntro() string {
	var b strings.Builder

	// ASCII art banner
	b.WriteString(applyGradient(
		asciiArt,
		lipgloss.Color("99"),
		lipgloss.Color("205"),
	))
	b.WriteString(newlineDouble)

	// Welcome message
	b.WriteString(
		titleStyle.Render(
			"Welcome to Spectr Initialization",
		),
	)
	b.WriteString(newlineDouble)

	b.WriteString(
		"This wizard will help you initialize Spectr in " +
			"your project.\n\n",
	)
	b.WriteString(
		"Spectr provides a structured approach to:\n",
	)
	b.WriteString(
		"  • Creating and managing change proposals\n",
	)
	b.WriteString(
		"  • Documenting project architecture and " +
			"specifications\n",
	)
	b.WriteString(
		"  • Integrating with AI coding assistants\n\n",
	)

	b.WriteString(
		infoStyle.Render(
			fmt.Sprintf(
				"Project path: %s",
				m.projectPath,
			),
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

	b.WriteString(
		titleStyle.Render(
			"Select AI Tools to Configure",
		),
	)
	b.WriteString(newlineDouble)

	b.WriteString(
		"Choose which AI coding tools you want to configure " +
			"with Spectr.\n",
	)
	b.WriteString(
		"You can come back later to add more tools.\n\n",
	)

	// Show search input if search mode is active
	if m.searchMode {
		b.WriteString(fmt.Sprintf(
			"Search: %s\n\n",
			m.searchInput.View(),
		))
	}

	// Render filtered providers or show no match message
	if len(m.filteredProviders) == 0 &&
		m.searchQuery != "" {
		b.WriteString(dimmedStyle.Render(
			fmt.Sprintf(
				"  No providers match '%s'\n",
				m.searchQuery,
			),
		))
	} else {
		b.WriteString(m.renderProviderGroup(m.filteredProviders, 0))
	}

	// Configured indicator explanation
	b.WriteString(doubleNewline)
	b.WriteString(
		subtleStyle.Render(
			"Items marked (configured) are already set up and will be updated.\n",
		),
	)

	// Instructions - show different help text based on search mode
	b.WriteString(newline)
	if m.searchMode {
		b.WriteString(
			subtleStyle.Render(
				"↑/↓: Navigate  Space: Toggle  Esc: Exit search  Enter: Continue\n",
			),
		)
	} else {
		b.WriteString(
			subtleStyle.Render(
				"↑/↓: Navigate  Space: Toggle  a: All  n: None  /: Search  " +
					"Enter: Continue  q: Quit\n",
			),
		)
	}

	return b.String()
}

func (m WizardModel) renderProviderGroup(
	providersList []providers.Registration,
	offset int,
) string {
	var b strings.Builder

	for i, provider := range providersList {
		actualIndex := offset + i
		cursor := " "
		if m.cursor == actualIndex {
			cursor = cursorStyle.Render("▸")
		}

		checkbox := "[ ]"
		if m.selectedProviders[provider.ID] {
			checkbox = selectedStyle.Render("[✓]")
		}

		// Build the base line with provider name
		line := fmt.Sprintf(
			"  %s %s %s",
			cursor,
			checkbox,
			provider.Name,
		)

		// Add configured indicator if provider is already configured
		configuredIndicator := ""
		if m.configuredProviders[provider.ID] {
			configuredIndicator = subtleStyle.Render(
				" (configured)",
			)
		}

		switch {
		case m.cursor == actualIndex:
			b.WriteString(
				cursorStyle.Render(line),
			)
			b.WriteString(configuredIndicator)
		case m.selectedProviders[provider.ID]:
			b.WriteString(
				selectedStyle.Render(line),
			)
			b.WriteString(configuredIndicator)
		default:
			b.WriteString(
				dimmedStyle.Render(line),
			)
			b.WriteString(configuredIndicator)
		}

		b.WriteString("\n")
	}

	return b.String()
}

func (m WizardModel) renderReview() string {
	var b strings.Builder

	b.WriteString(
		titleStyle.Render(
			"Review Your Selections",
		),
	)
	b.WriteString("\n\n")

	selectedCount := len(
		m.getSelectedProviderIDs(),
	)
	m.renderSelectedProviders(&b, selectedCount)

	// CI workflow option
	b.WriteString("Additional options:\n\n")
	checkbox := "[ ]"
	if m.ciWorkflowEnabled {
		checkbox = selectedStyle.Render("[✓]")
	}
	configuredNote := ""
	if m.ciWorkflowConfigured {
		configuredNote = subtleStyle.Render(
			" (configured)",
		)
	}
	b.WriteString(fmt.Sprintf(
		"  %s Spectr CI Validation%s\n",
		checkbox,
		configuredNote,
	))
	b.WriteString(
		subtleStyle.Render(
			"      Creates .github/workflows/spectr-ci.yml for automated validation\n\n",
		),
	)

	m.renderCreationPlan(&b, selectedCount)

	b.WriteString(newline)
	b.WriteString(
		subtleStyle.Render(
			"Space: Toggle option  Enter: Initialize  Backspace: Go back  " +
				"'q': Quit\n",
		),
	)

	return b.String()
}

// renderSelectedProviders displays the selected providers or warning if none
func (m WizardModel) renderSelectedProviders(
	b *strings.Builder,
	count int,
) {
	if count == 0 {
		b.WriteString(
			errorStyle.Render(
				"⚠ No tools selected",
			),
		)
		b.WriteString(doubleNewline)
		b.WriteString(
			"You haven't selected any tools to configure.\n",
		)
		b.WriteString(
			"Spectr will still be initialized, but no tool " +
				"integrations will be set up.\n\n",
		)

		return
	}

	fmt.Fprintf(
		b,
		"You have selected %d tool(s) to configure:\n\n",
		count,
	)

	for _, provider := range m.allProviders {
		if !m.selectedProviders[provider.ID] {
			continue
		}
		b.WriteString(successStyle.Render("  ✓ "))
		b.WriteString(provider.Name)
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

// renderCreationPlan displays what files will be created
func (m WizardModel) renderCreationPlan(
	b *strings.Builder,
	count int,
) {
	b.WriteString(
		"The following will be created:\n",
	)
	b.WriteString(
		infoStyle.Render("  • spectr/project.md"),
	)
	b.WriteString(
		" - Project documentation template" + newline,
	)
	b.WriteString(
		infoStyle.Render("  • spectr/AGENTS.md"),
	)
	b.WriteString(
		" - AI agent instructions" + newline,
	)

	if count > 0 {
		b.WriteString(
			infoStyle.Render(fmt.Sprintf(
				"  • Tool configurations for %d selected tools",
				count,
			)),
		)
		b.WriteString(newline)
	}

	if m.ciWorkflowEnabled {
		status := "Create"
		if m.ciWorkflowConfigured {
			status = "Update"
		}
		b.WriteString(
			infoStyle.Render(
				"  • .github/workflows/spectr-ci.yml",
			),
		)
		fmt.Fprintf(
			b,
			" - %s CI workflow for Spectr validation\n",
			status,
		)
	}
}

func (WizardModel) renderExecute() string {
	var b strings.Builder

	b.WriteString(
		titleStyle.Render(
			"Initializing Spectr...",
		),
	)
	b.WriteString(doubleNewline)

	b.WriteString(
		infoStyle.Render(
			"⏳ Setting up your project...",
		),
	)
	b.WriteString(doubleNewline)

	b.WriteString(
		"This will only take a moment." + newline,
	)

	return b.String()
}

func (m WizardModel) renderComplete() string {
	var b strings.Builder

	if m.err != nil {
		m.renderError(&b)

		return b.String()
	}

	b.WriteString(
		successStyle.Render(
			"✓ Spectr Initialized Successfully!",
		),
	)
	b.WriteString("\n\n")

	if m.executionResult != nil {
		m.renderExecutionResults(&b)
	}

	b.WriteString(FormatNextStepsMessage())
	b.WriteString(newline)
	b.WriteString(
		subtleStyle.Render(
			"c: copy prompt  q: quit" + newline,
		),
	)

	return b.String()
}

// renderError displays initialization errors
func (m WizardModel) renderError(
	b *strings.Builder,
) {
	b.WriteString(
		errorStyle.Render(
			"✗ Initialization Failed",
		),
	)
	b.WriteString(doubleNewline)
	b.WriteString(
		errorStyle.Render(m.err.Error()),
	)
	b.WriteString(doubleNewline)

	if m.executionResult != nil &&
		len(m.executionResult.Errors) > 0 {
		b.WriteString("Errors:" + newline)
		for _, err := range m.executionResult.Errors {
			b.WriteString(
				errorStyle.Render("  • "),
			)
			b.WriteString(err)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}

	b.WriteString(
		subtleStyle.Render("Press 'q' to quit\n"),
	)
}

// renderExecutionResults displays created/updated files and warnings
func (m WizardModel) renderExecutionResults(
	b *strings.Builder,
) {
	if len(m.executionResult.CreatedFiles) > 0 {
		b.WriteString(
			successStyle.Render("Created files:"),
		)
		b.WriteString(newline)
		for _, file := range m.executionResult.CreatedFiles {
			b.WriteString(
				infoStyle.Render("  ✓ "),
			)
			b.WriteString(file)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}

	if len(m.executionResult.UpdatedFiles) > 0 {
		b.WriteString(
			successStyle.Render("Updated files:"),
		)
		b.WriteString(newline)
		for _, file := range m.executionResult.UpdatedFiles {
			b.WriteString(
				infoStyle.Render("  ↻ "),
			)
			b.WriteString(file)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}

	if len(m.executionResult.Errors) > 0 {
		b.WriteString(
			errorStyle.Render("Warnings:"),
		)
		b.WriteString(newline)
		for _, err := range m.executionResult.Errors {
			b.WriteString(
				errorStyle.Render("  ⚠ "),
			)
			b.WriteString(err)
			b.WriteString(newline)
		}
		b.WriteString(newline)
	}
}

func (m WizardModel) getSelectedProviderIDs() []string {
	var selected []string
	for id, isSelected := range m.selectedProviders {
		if isSelected {
			selected = append(selected, id)
		}
	}

	return selected
}

// executeInit runs the initialization and sends result
func executeInit(
	projectPath string,
	selectedProviders []string,
	ciWorkflowEnabled bool,
) tea.Cmd {
	return func() tea.Msg {
		// Create a minimal InitCmd for the executor
		cmd := &InitCmd{
			Path: projectPath,
		}
		executor, err := NewInitExecutor(cmd)
		if err != nil {
			return ExecutionCompleteMsg{
				result: nil,
				err: fmt.Errorf(
					"failed to create executor: %w",
					err,
				),
			}
		}

		result, err := executor.Execute(
			selectedProviders,
			ciWorkflowEnabled,
		)

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
