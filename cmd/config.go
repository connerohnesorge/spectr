package cmd

import (
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/config"
)

// ConfigCmd displays the current configuration and path
type ConfigCmd struct{}

// Run executes the config command
func (*ConfigCmd) Run() error {
	configPath := config.ConfigPath()

	// Check if config file exists
	_, err := os.Stat(configPath)
	fileExists := err == nil

	// Display config path
	fmt.Printf("Configuration Path: %s\n", configPath)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Display status
	if fileExists {
		fmt.Println("Status: Loaded")
	} else {
		fmt.Println("Status: No config file found, using defaults")
	}

	fmt.Println()

	// Display theme configuration
	fmt.Println("Current Theme:")
	fmt.Printf("  accent:    %s\n", derefString(cfg.Theme.Accent))
	fmt.Printf("  error:     %s\n", derefString(cfg.Theme.Error))
	fmt.Printf("  success:   %s\n", derefString(cfg.Theme.Success))
	fmt.Printf("  border:    %s\n", derefString(cfg.Theme.Border))
	fmt.Printf("  help:      %s\n", derefString(cfg.Theme.Help))
	fmt.Printf("  selected:  %s\n", derefString(cfg.Theme.Selected))
	fmt.Printf("  highlight: %s\n", derefString(cfg.Theme.Highlight))
	fmt.Printf("  header:    %s\n", derefString(cfg.Theme.Header))

	return nil
}

// derefString safely dereferences a string pointer,
// returning empty string if nil
func derefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
