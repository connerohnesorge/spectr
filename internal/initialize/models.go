package initialize

// InitCmd represents the init command with all its flags
type InitCmd struct {
	Path           string   `arg:"" optional:"" help:"Project path"`
	PathFlag       string   `                   help:"Alt project path"        name:"path"            short:"p"` //nolint:lll,revive
	Tools          []string `                   help:"Tools list"              name:"tools"           short:"t"` //nolint:lll,revive
	NonInteractive bool     `                   help:"Non-interactive"         name:"non-interactive"`           //nolint:lll,revive
	CIWorkflow     bool     `                   help:"Create CI workflow file" name:"ci-workflow"`               //nolint:lll,revive
}

// ProjectConfig holds the overall project configuration during init
type ProjectConfig struct {
	// ProjectPath is the absolute path to the project directory
	ProjectPath string
	// SelectedTools is the list of tools the user has selected to configure
	SelectedTools []string
	// SpectrEnabled indicates whether Spectr framework should be initialized
	SpectrEnabled bool
}

// InitState represents the current state of the initialization process
type InitState int

const (
	// StateSelectTools is the tool selection screen
	StateSelectTools InitState = iota
	// StateConfigureTools is the tool configuration screen
	StateConfigureTools
	// StateConfirmation is the final confirmation screen
	StateConfirmation
	// StateComplete is the completion state
	StateComplete
)