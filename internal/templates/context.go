// Package templates manages embedded templates for Spectr initialization.
package templates

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
