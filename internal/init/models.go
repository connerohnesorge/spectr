package init

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

// ProjectContext holds template variables for rendering project.md
type ProjectContext struct {
	// ProjectName is the name of the project
	ProjectName string
	// Description is the project description/purpose
	Description string
	// TechStack is the list of technologies used
	TechStack []string
	// Conventions are the project conventions (unused in template currently)
	Conventions string
}

// InitCmd represents the init command with all its flags
type InitCmd struct {
	Path           string   `arg:"" optional:"" help:"Project path"`
	PathFlag       string   `name:"path" short:"p" help:"Alt project path"`
	Tools          []string `name:"tools" short:"t" help:"Tools list"`
	NonInteractive bool     `name:"non-interactive" help:"Non-interactive"`
}
