// Package initializers provides built-in initializers for the provider system.
package initializers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
)

// Compile-time interface satisfaction checks.
var (
	_ providers.ClaudeInitializerFactory = (*Factory)(
		nil,
	)
	_ providers.GeminiInitializerFactory = (*Factory)(
		nil,
	)
	_ providers.InitializerFactory = (*Factory)(
		nil,
	)
)

// Factory implements both ClaudeInitializerFactory and
// GeminiInitializerFactory interfaces. It creates initializers for
// providers without causing import cycles.
type Factory struct{}

// NewFactory creates a new Factory instance.
func NewFactory() *Factory {
	return &Factory{}
}

// CreateDirectoryInitializer creates a DirectoryInitializer for the
// given paths.
func (*Factory) CreateDirectoryInitializer(
	paths ...string,
) providers.Initializer {
	return NewDirectoryInitializer(paths...)
}

// CreateConfigFileInitializer creates a ConfigFileInitializer for the
// given path and template.
func (*Factory) CreateConfigFileInitializer(
	path, template string,
) providers.Initializer {
	return NewConfigFileInitializer(
		path,
		template,
	)
}

// CreateSlashCommandsInitializer creates a SlashCommandsInitializer.
func (*Factory) CreateSlashCommandsInitializer(
	dir, ext string,
	format providers.CommandFormat,
	renderer providers.TemplateRenderer,
) providers.Initializer {
	return NewSlashCommandsInitializer(
		dir,
		ext,
		format,
		renderer,
	)
}
