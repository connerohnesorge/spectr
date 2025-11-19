//nolint:revive // line-length-limit - comments need clarity

package init

// ToolType represents the type of tool configuration
type ToolType string

const (
	// ToolTypeConfig represents tools using config files
	ToolTypeConfig ToolType = "config"
	// ToolTypeSlash represents tools using slash commands
	ToolTypeSlash ToolType = "slash"
)

// ToolDefinition defines the configuration for a tool integration
type ToolDefinition struct {
	// ID is the unique identifier for the tool
	ID string
	// Name is the human-readable name of the tool
	Name string
	// Type indicates whether this is a config or slash tool
	Type ToolType
	// SlashCommand is the slash command syntax (for slash-command based tools)
	SlashCommand string
	// Priority determines the display order (lower numbers first)
	Priority int
	// Configured indicates whether the tool has been configured by the user
	Configured bool
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

// SpecContext holds template variables for rendering spec.md
type SpecContext struct {
	// CapabilityName is the name of the capability/spec
	CapabilityName string
}

// ProposalContext holds template variables for rendering proposal.md
type ProposalContext struct {
	// ChangeName is the name of the change proposal
	ChangeName string
}
