// Package config provides configuration loading for Spectr projects.
// This file handles loading and parsing of spectr.yaml configuration files.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultAppendTasksSection is the default section name for appended tasks
// when append_tasks.section is not specified in the config.
const DefaultAppendTasksSection = "Automated Tasks"

// ErrConfigMalformed is returned when the config file contains invalid YAML.
var ErrConfigMalformed = errors.New(
	"config file is malformed",
)

// Config represents the root configuration structure for spectr.yaml.
type Config struct {
	// AppendTasks defines tasks to automatically append during accept.
	AppendTasks *AppendTasksConfig `yaml:"append_tasks"`
	// RefsAlwaysPrepend defines tasks to prepend to each child task file (v2 format).
	RefsAlwaysPrepend *RefsTasksConfig `yaml:"refs_always_prepend"`
	// RefsAlwaysAppend defines tasks to append to each child task file (v2 format).
	RefsAlwaysAppend *RefsTasksConfig `yaml:"refs_always_append"`
}

// AppendTasksConfig defines the configuration for auto-appending tasks.
type AppendTasksConfig struct {
	// Section is the name of the section for appended tasks.
	// Defaults to "Automated Tasks" if not specified.
	Section string `yaml:"section"`
	// Tasks is the list of task descriptions to append.
	Tasks []string `yaml:"tasks"`
}

// GetSection returns the section name, using the default if not specified.
func (c *AppendTasksConfig) GetSection() string {
	if c == nil || c.Section == "" {
		return DefaultAppendTasksSection
	}

	return c.Section
}

// HasTasks returns true if there are tasks to append.
func (c *AppendTasksConfig) HasTasks() bool {
	if c == nil {
		return false
	}

	return len(c.Tasks) > 0
}

// RefsTasksConfig defines tasks to prepend or append to each child task file.
// Used for refs_always_prepend and refs_always_append configuration.
type RefsTasksConfig struct {
	// Tasks is the list of task descriptions to prepend/append to each child file.
	Tasks []string `yaml:"tasks"`
}

// HasTasks returns true if there are tasks configured.
func (c *RefsTasksConfig) HasTasks() bool {
	if c == nil {
		return false
	}

	return len(c.Tasks) > 0
}

// LoadConfig searches for and loads spectr.yaml from the given directory
// or any parent directory. Returns nil config (not an error) if no config
// file is found.
func LoadConfig(
	startDir string,
) (*Config, error) {
	configPath, err := findConfigFile(startDir)
	if err != nil {
		return nil, err
	}
	if configPath == "" {
		// No config file found - this is not an error
		return nil, nil
	}

	return parseConfigFile(configPath)
}

// findConfigFile walks up from startDir to find spectr.yaml.
// Returns empty string if not found (not an error).
func findConfigFile(
	startDir string,
) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get absolute path: %w",
			err,
		)
	}

	for {
		configPath := filepath.Join(
			dir,
			"spectr.yaml",
		)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return "", nil
}

// parseConfigFile reads and parses the config file at the given path.
func parseConfigFile(
	path string,
) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read config file: %w",
			err,
		)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf(
			"%w: %s: %v",
			ErrConfigMalformed,
			path,
			err,
		)
	}

	return &cfg, nil
}
