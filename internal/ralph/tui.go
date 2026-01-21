//nolint:revive // TUI implementation file intentionally exceeds 300 lines for cohesion

// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tui.go implements the Bubble Tea TUI model for Ralph orchestration.
// It provides an interactive terminal interface showing:
// - Task list with real-time status indicators
// - Live agent output pane with ANSI rendering
// - Keyboard controls for user interaction
// - Help bar with available commands
//
// The TUI model integrates with the orchestrator to provide
// visual feedback during task execution and enable user control
// over the workflow (retry, skip, abort, pause).

const (
	// minHeaderLines is the minimum number of lines for the header
	minHeaderLines = 3

	// minHelpBarLines is the minimum number of lines for the help bar
	minHelpBarLines = 2

	// minOutputHeight is the minimum height for the output pane
	minOutputHeight = 5

	// maxTaskListHeight is the maximum height for the task list
	maxTaskListHeight = 10

	// taskListBaseLines is the base number of lines for the task list
	taskListBaseLines = 2

	// maxSectionNameLength is the maximum length for section names
	maxSectionNameLength = 30

	// quitKey is the key for quitting
	quitKey = "q"
)

// TUIMode represents the current mode of the TUI
type TUIMode int

const (
	// ModeNormal is the default mode showing task execution
	ModeNormal TUIMode = iota
	// ModeInteractive is the task selection mode
	ModeInteractive
	// ModePaused is when orchestration is paused
	ModePaused
	// ModeFailure is when a task has failed and awaiting user action
	ModeFailure
)

// TUIModel is the Bubble Tea model for the Ralph TUI
type TUIModel struct {
	// changeID is the change proposal identifier
	changeID string

	// tasks is the list of all tasks in the graph
	tasks []*Task

	// currentTaskID is the ID of the currently executing task
	currentTaskID string

	// output is the viewport for agent output
	output viewport.Model

	// outputLines stores the accumulated output lines
	outputLines []string

	// mode is the current TUI mode
	mode TUIMode

	// cursor is the cursor position in interactive mode
	cursor int

	// selectedTasks is the set of selected task IDs in interactive mode
	selectedTasks map[string]bool

	// width and height are the terminal dimensions
	width  int
	height int

	// mu protects concurrent access to output
	mu sync.Mutex

	// quitting indicates if the TUI is shutting down
	quitting bool

	// paused indicates if orchestration is paused
	paused bool

	// userAction is set when user triggers an action (retry, skip, abort)
	userAction UserAction

	// onUserAction is called when user makes a decision on failure
	onUserAction func(action UserAction) tea.Cmd

	// onQuit is called when user quits
	onQuit func() tea.Cmd

	// ready indicates if the viewport is initialized
	ready bool
}

// TUIConfig holds configuration for creating a TUI model
type TUIConfig struct {
	ChangeID      string
	Tasks         []*Task
	OnUserAction  func(action UserAction) tea.Cmd
	OnQuit        func() tea.Cmd
	Interactive   bool
	InitialWidth  int
	InitialHeight int
}

// NewTUIModel creates a new TUI model
func NewTUIModel(config *TUIConfig) *TUIModel {
	mode := ModeNormal
	if config.Interactive {
		mode = ModeInteractive
	}

	return &TUIModel{
		changeID:      config.ChangeID,
		tasks:         config.Tasks,
		mode:          mode,
		selectedTasks: make(map[string]bool),
		outputLines:   make([]string, 0),
		width:         config.InitialWidth,
		height:        config.InitialHeight,
		onUserAction:  config.OnUserAction,
		onQuit:        config.OnQuit,
	}
}

// PTYOutputMsg is a message containing output from the PTY
type PTYOutputMsg struct {
	Line string
}

// TaskStartMsg is sent when a task starts
type TaskStartMsg struct {
	TaskID string
}

// TaskCompleteMsg is sent when a task completes
type TaskCompleteMsg struct {
	TaskID  string
	Success bool
}

// TaskFailMsg is sent when a task fails
type TaskFailMsg struct {
	TaskID string
	Error  error
}

// Init initializes the TUI model
func (*TUIModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
//
//nolint:revive // Type switch intentionally shadows msg parameter
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// Initialize viewport on first window size message
			m.output = viewport.New(msg.Width, m.getOutputHeight())
			m.output.YPosition = 0
			m.ready = true
			m.updateViewportContent()
		} else {
			m.output.Width = msg.Width
			m.output.Height = m.getOutputHeight()
			m.updateViewportContent()
		}

	case PTYOutputMsg:
		m.mu.Lock()
		m.outputLines = append(m.outputLines, msg.Line)
		m.updateViewportContentLocked()
		m.mu.Unlock()

	case TaskStartMsg:
		m.currentTaskID = msg.TaskID
		m.clearOutput()

	case TaskCompleteMsg:
		// Update task status in our local copy
		for _, task := range m.tasks {
			if task.ID == msg.TaskID {
				if msg.Success {
					task.Status = StatusCompleted
				} else {
					task.Status = StatusFailed
				}

				break
			}
		}

	case TaskFailMsg:
		m.mode = ModeFailure
		for _, task := range m.tasks {
			if task.ID == msg.TaskID {
				task.Status = StatusFailed

				break
			}
		}
	}

	// Update viewport
	if m.ready {
		var cmd tea.Cmd
		m.output, cmd = m.output.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleKeyPress handles keyboard input
func (m *TUIModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit keys
	if msg.Type == tea.KeyCtrlC {
		m.quitting = true
		if m.onQuit != nil {
			return m, m.onQuit()
		}

		return m, tea.Quit
	}

	switch m.mode {
	case ModeInteractive:
		return m.handleInteractiveMode(msg)
	case ModeFailure:
		return m.handleFailureMode(msg)
	case ModePaused:
		return m.handlePausedMode(msg)
	case ModeNormal:
		return m.handleNormalMode(msg)
	default:
		return m.handleNormalMode(msg)
	}
}

// handleNormalMode handles keys in normal mode
func (m *TUIModel) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case quitKey:
		m.quitting = true
		if m.onQuit != nil {
			return m, m.onQuit()
		}

		return m, tea.Quit

	case "p":
		m.paused = !m.paused
		if m.paused {
			m.mode = ModePaused
		} else {
			m.mode = ModeNormal
		}

	case "i":
		m.mode = ModeInteractive
		m.cursor = 0
	}

	return m, nil
}

// handleInteractiveMode handles keys in interactive mode
func (m *TUIModel) handleInteractiveMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case quitKey, "esc":
		m.mode = ModeNormal

		return m, nil

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.tasks)-1 {
			m.cursor++
		}

	case " ":
		// Toggle selection
		if m.cursor < len(m.tasks) {
			taskID := m.tasks[m.cursor].ID
			m.selectedTasks[taskID] = !m.selectedTasks[taskID]
		}

	case "enter":
		// Apply selection and return to normal mode
		m.mode = ModeNormal
		// TODO: Notify orchestrator of selected tasks
	}

	return m, nil
}

// handleFailureMode handles keys when a task has failed
func (m *TUIModel) handleFailureMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		m.userAction = UserActionRetry
		m.mode = ModeNormal
		if m.onUserAction != nil {
			return m, m.onUserAction(UserActionRetry)
		}

	case "s":
		m.userAction = UserActionSkip
		m.mode = ModeNormal
		if m.onUserAction != nil {
			return m, m.onUserAction(UserActionSkip)
		}

	case quitKey:
		m.userAction = UserActionAbort
		m.quitting = true
		if m.onUserAction != nil {
			return m, m.onUserAction(UserActionAbort)
		}

		return m, tea.Quit
	}

	return m, nil
}

// handlePausedMode handles keys when paused
func (m *TUIModel) handlePausedMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "p":
		m.paused = false
		m.mode = ModeNormal

	case quitKey:
		m.quitting = true
		if m.onQuit != nil {
			return m, m.onQuit()
		}

		return m, tea.Quit
	}

	return m, nil
}

// View renders the TUI
func (m *TUIModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return "Initializing..."
	}

	switch m.mode {
	case ModeInteractive:
		return m.renderInteractiveView()
	case ModeNormal, ModePaused, ModeFailure:
		return m.renderNormalView()
	default:
		return m.renderNormalView()
	}
}

// clearOutput clears the output buffer
func (m *TUIModel) clearOutput() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.outputLines = make([]string, 0)
	m.updateViewportContentLocked()
}

// updateViewportContent updates the viewport with current output.
// This method acquires the mutex before updating.
func (m *TUIModel) updateViewportContent() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateViewportContentLocked()
}

// updateViewportContentLocked updates the viewport without acquiring the mutex.
// The caller must hold the mutex before calling this method.
func (m *TUIModel) updateViewportContentLocked() {
	if !m.ready {
		return
	}

	content := strings.Join(m.outputLines, "\n")
	m.output.SetContent(content)

	// Auto-scroll to bottom
	m.output.GotoBottom()
}

// getOutputHeight calculates the height for the output pane
func (m *TUIModel) getOutputHeight() int {
	// Reserve space for header, task list (dynamic), help bar
	taskListHeight := minInt(len(m.tasks)+taskListBaseLines, maxTaskListHeight)
	reserved := minHeaderLines + taskListHeight + minHelpBarLines
	outputHeight := m.height - reserved

	if outputHeight < minOutputHeight {
		outputHeight = minOutputHeight
	}

	return outputHeight
}

// AddOutput adds a line of output to the TUI (thread-safe)
func (m *TUIModel) AddOutput(line string) {
	m.mu.Lock()
	m.outputLines = append(m.outputLines, line)
	m.mu.Unlock()
}

// SetCurrentTask sets the currently executing task
func (m *TUIModel) SetCurrentTask(taskID string) {
	m.currentTaskID = taskID
}

// StreamPTYOutput creates a goroutine that streams PTY output to the TUI
func StreamPTYOutput(pty io.Reader, program *tea.Program) {
	scanner := bufio.NewScanner(pty)
	for scanner.Scan() {
		line := scanner.Text()
		program.Send(PTYOutputMsg{Line: line})
	}
}

// GetCompletedCount returns the number of completed tasks
func (m *TUIModel) GetCompletedCount() int {
	count := 0
	for _, task := range m.tasks {
		if task.Status == StatusCompleted {
			count++
		}
	}

	return count
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Status indicator styles
var (
	statusCompleted  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).SetString("✓")
	statusInProgress = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).SetString("▶")
	statusPending    = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).SetString("○")
	statusFailed     = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).SetString("✗")
)

// getStatusIcon returns the icon for a task status
func getStatusIcon(status string) string {
	switch status {
	case StatusCompleted:
		return statusCompleted.String()
	case "in_progress":
		return statusInProgress.String()
	case "failed":
		return statusFailed.String()
	default:
		return statusPending.String()
	}
}

// FormatTaskLine formats a single task line with status icon
//
//nolint:revive // Boolean parameter is clear from name 'current' and simplifies API
func FormatTaskLine(task *Task, current bool) string {
	icon := getStatusIcon(task.Status)
	taskStr := fmt.Sprintf("%s %s", icon, task.ID)

	if current {
		taskStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Bold(true).
			Render(taskStr)
	}

	// Add section name if not too long
	if len(task.Section) < maxSectionNameLength {
		taskStr += fmt.Sprintf(" - %s", task.Section)
	}

	return taskStr
}
