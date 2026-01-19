# Provider System - Frontmatter Overrides

## ADDED Requirements

### Requirement: Frontmatter Override Structure

The system SHALL provide a `FrontmatterOverride` type in `internal/domain` for modifying slash command frontmatter.

#### Scenario: Override structure definition

- **WHEN** code needs to specify frontmatter modifications
- **THEN** it SHALL use `domain.FrontmatterOverride` with the following structure:

```go
type FrontmatterOverride struct {
    Set    map[string]interface{} // Fields to add or modify
    Remove []string                // Field names to delete
}
```

- **AND** `Set` values SHALL support any YAML-serializable type (string, bool, []string, etc.)
- **AND** `Remove` SHALL be applied after `Set` to allow replacing fields

#### Scenario: Creating overrides for Claude Code proposal command

- **WHEN** ClaudeProvider needs to customize the proposal slash command
- **THEN** it SHALL create a `FrontmatterOverride` like:

```go
&domain.FrontmatterOverride{
    Set: map[string]interface{}{
        "context": "fork",
    },
    Remove: []string{"agent"},
}
```

- **AND** this SHALL add `context: fork` to the frontmatter
- **AND** this SHALL remove the `agent` field from the frontmatter

### Requirement: Frontmatter Parsing and Rendering

The system SHALL provide utilities in `internal/domain` for YAML frontmatter manipulation.

#### Scenario: Parse frontmatter from template content

- **WHEN** `ParseFrontmatter(content string)` is called with markdown containing YAML frontmatter
- **THEN** it SHALL return a `map[string]interface{}` with parsed YAML values
- **AND** it SHALL return the body content without frontmatter
- **AND** it SHALL return an error if frontmatter is malformed

Example:
```go
content := `---
description: Test
allowed-tools: Read, Write
---
# Body`

fm, body, err := domain.ParseFrontmatter(content)
// fm = map[string]interface{}{"description": "Test", "allowed-tools": "Read, Write"}
// body = "# Body"
// err = nil
```

#### Scenario: Render frontmatter to YAML

- **WHEN** `RenderFrontmatter(fm map[string]interface{}, body string)` is called
- **THEN** it SHALL serialize `fm` to YAML format
- **AND** it SHALL wrap YAML in `---` delimiters
- **AND** it SHALL append the body content
- **AND** it SHALL return the complete markdown with frontmatter

Example:
```go
fm := map[string]interface{}{"context": "fork", "description": "Test"}
body := "# Body"

result := domain.RenderFrontmatter(fm, body)
// result = "---\ncontext: fork\ndescription: Test\n---\n# Body"
```

#### Scenario: Apply overrides to frontmatter

- **WHEN** `ApplyFrontmatterOverrides(base map[string]interface{}, overrides *FrontmatterOverride)` is called
- **THEN** it SHALL copy `base` to avoid mutation
- **AND** it SHALL apply all fields from `overrides.Set`, replacing existing values
- **AND** it SHALL remove all fields listed in `overrides.Remove`
- **AND** it SHALL return the modified frontmatter map

Example:
```go
base := map[string]interface{}{"description": "Test", "agent": "plan", "subtask": false}
overrides := &domain.FrontmatterOverride{
    Set:    map[string]interface{}{"context": "fork"},
    Remove: []string{"agent"},
}

result := domain.ApplyFrontmatterOverrides(base, overrides)
// result = map[string]interface{}{"description": "Test", "context": "fork", "subtask": false}
```

### Requirement: TemplateManager Override Support

The system SHALL extend `TemplateManager` to support frontmatter overrides for slash commands.

#### Scenario: SlashCommandWithOverrides method

- **WHEN** `TemplateManager.SlashCommandWithOverrides(cmd domain.SlashCommand, overrides *FrontmatterOverride)` is called
- **THEN** it SHALL look up the base template for `cmd`
- **AND** it SHALL parse the template's frontmatter
- **AND** it SHALL apply the `overrides` to the frontmatter
- **AND** it SHALL render the modified frontmatter back to YAML
- **AND** it SHALL return a `domain.TemplateRef` that produces the modified content when rendered
- **AND** if `overrides` is `nil`, it SHALL behave identically to `SlashCommand(cmd)`

#### Scenario: Rendering overridden template

- **WHEN** `TemplateRef.Render(ctx)` is called on a template with overrides
- **THEN** the rendered content SHALL contain the modified frontmatter
- **AND** the body content SHALL remain unchanged from the base template

### Requirement: ClaudeProvider Uses Overrides

The `ClaudeProvider` SHALL use frontmatter overrides for the proposal slash command.

#### Scenario: Claude Code proposal command has context: fork

- **WHEN** `ClaudeProvider.Initializers()` is called
- **THEN** it SHALL pass a `FrontmatterOverride` to the proposal slash command initializer
- **AND** the override SHALL add `context: fork`
- **AND** the override SHALL remove the `agent` field
- **AND** the generated `.claude/commands/spectr/proposal.md` SHALL contain `context: fork` in frontmatter
- **AND** the generated file SHALL NOT contain `agent:` in frontmatter

#### Scenario: Claude Code apply command uses defaults

- **WHEN** `ClaudeProvider.Initializers()` is called
- **THEN** the apply slash command SHALL use default frontmatter (no overrides)
- **AND** the generated `.claude/commands/spectr/apply.md` SHALL match the base template exactly
