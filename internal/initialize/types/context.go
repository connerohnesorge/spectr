package types

// TemplateContext holds path-related template variables for dynamic
// directory names.
type TemplateContext struct {
	// BaseDir is the base directory for spectr files (default: "spectr")
	BaseDir string
	// SpecsDir is the directory for spec files (default: "spectr/specs")
	SpecsDir string
	// ChangesDir is the directory for change proposals
	// (default: "spectr/changes")
	ChangesDir string
	// ProjectFile is the path to the project configuration file
	// (default: "spectr/project.md")
	ProjectFile string
	// AgentsFile is the path to the agents file (default: "spectr/AGENTS.md")
	AgentsFile string
}

// DefaultTemplateContext returns a TemplateContext with default values.
func DefaultTemplateContext() TemplateContext {
	return TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}
}
