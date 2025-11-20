package init

const (
	// File and directory permission constants
	// Consolidated from filesystem.go and configurator.go
	dirPerm         = 0o755
	filePerm        = 0o644
	defaultFilePerm = 0o644

	// UI control keys
	keyQuit  = "q"
	keyEnter = "enter"
	keyCtrlC = "ctrl+c"

	// Common strings
	newlineDouble = "\n\n"

	// Marker constants for managing config file updates
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"
)
