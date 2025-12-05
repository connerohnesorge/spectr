## MODIFIED Requirements

### Requirement: Gemini Slash Commands
The provider SHALL create slash commands in `.gemini/commands/spectr/` directory using TOML format.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it creates `.gemini/commands/spectr/proposal.toml`
- **AND** it creates `.gemini/commands/spectr/apply.toml`

#### Scenario: TOML command format
- **WHEN** slash command files are created
- **THEN** they use TOML format
- **AND** include `description` field with command description
- **AND** include `prompt` field with command content
