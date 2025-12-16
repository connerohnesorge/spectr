# Provider System Design

## Overview
This document outlines the architectural changes for the Provider system refactor. The goal is to move from an inheritance-heavy, behavior-rich `Provider` interface to a clean, data-driven design.

## Current vs Proposed
- **Current**: `Provider` interface includes `Configure(Project, TemplateRenderer) error`. `BaseProvider` implements the logic.
- **Proposed**: `Provider` interface returns data structs (`Info`, `Capabilities`). An `Enactor` service implements the logic.

## Type Definitions

### 1. Provider Interface
The `Provider` interface is reduced to returning metadata and capabilities.

```go
package providers

// Provider describes an AI tool integration.
type Provider interface {
	// Info returns the provider's metadata.
	Info() ProviderInfo

	// Capabilities returns the set of capabilities this provider supports.
	Capabilities() []Capability
}

type ProviderInfo struct {
	ID       string // e.g. "claude-code"
	Name     string // e.g. "Claude Code"
	Priority int    // e.g. 1
}
```

### 2. Capabilities
Capabilities are data structs that implement a marker interface.

```go
// Capability is a marker interface for provider features.
type Capability interface {
	isCapability()
}

// InstructionFileCapability declares that the provider needs an instruction file (e.g. CLAUDE.md).
type InstructionFileCapability struct {
	Filename string
    // If true, the Enactor will look for existing content variants or strictly overwrite?
    // For now, we assume the standard "UpdateFileWithMarkers" behavior.
}

func (c InstructionFileCapability) isCapability() {}

// SlashCommandsCapability declares that the provider needs slash commands.
type SlashCommandsCapability struct {
	Commands []SlashCommand
}

func (c SlashCommandsCapability) isCapability() {}

// SlashCommand defines a single slash command file.
type SlashCommand struct {
	Name        string // e.g. "proposal", "apply" - used for template selection
	Path        string // e.g. ".claude/commands/spectr/proposal.md"
	Description string // Used for the YAML frontmatter
}
```

### 3. The Enactor
The `Enactor` (or `Applier`) handles the execution.

```go
type Enactor struct {
    tm TemplateRenderer
}

func NewEnactor(tm TemplateRenderer) *Enactor {
    return &Enactor{tm: tm}
}

// Apply configures a provider in the given project.
func (e *Enactor) Apply(p Provider, projectPath string) ([]string, error) {
    var createdFiles []string
    
    for _, cap := range p.Capabilities() {
        switch c := cap.(type) {
        case InstructionFileCapability:
            // 1. Render instruction pointer template
            // 2. Write/Update file at projectPath/c.Filename
            // 3. Add to createdFiles
            
        case SlashCommandsCapability:
            // 1. For each command in c.Commands:
            //    a. Render command template (using c.Name)
            //    b. Inject c.Description into frontmatter
            //    c. Write to projectPath/c.Path
            //    d. Add to createdFiles
        }
    }
    return createdFiles, nil
}
```

### 4. IsConfigured Check
We also need to replace `IsConfigured`.

```go
func (e *Enactor) IsConfigured(p Provider, projectPath string) bool {
    for _, cap := range p.Capabilities() {
        switch c := cap.(type) {
        case InstructionFileCapability:
            if !FileExists(filepath.Join(projectPath, c.Filename)) {
                return false
            }
        case SlashCommandsCapability:
            for _, cmd := range c.Commands {
                if !FileExists(filepath.Join(projectPath, cmd.Path)) {
                    return false
                }
            }
        }
    }
    return true
}
```

## Migration Strategy
1.  Define the new types in `internal/initialize/providers/types.go` (new file).
2.  Implement `Enactor` in `internal/initialize/providers/enactor.go`.
3.  Rewrite `claude.go` to implement the new `Provider` interface.
4.  Updates `registry.go` to store the new `Provider` type.
5.  Update `InitExecutor` to use `Enactor`.
6.  Delete `BaseProvider` and old methods.
