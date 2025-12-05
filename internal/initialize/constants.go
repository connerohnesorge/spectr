package initialize

const (
	// File and directory permission constants
	filePerm = 0o644

	// UI control keys
	keyQuit  = "q"
	keyEnter = "enter"
	keyCtrlC = "ctrl+c"
	keyCopy  = "c"

	// Common strings
	newlineDouble = "\n\n"

	// Marker constants for managing config file updates
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"

	// PopulateContextPrompt is the suggested prompt for users to populate
	// their project context.
	PopulateContextPrompt = "Review spectr/project.md and help me fill in " +
		"our project's tech stack, conventions, and description. " +
		"Ask me questions to understand the codebase."
)
