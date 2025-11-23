# Change: Add User Configuration File Support Using Kong

## Why

Currently, Spectr requires users to specify flags for common preferences on every command invocation (e.g., `--strict`, `--json`, `--interactive`). This creates friction for users who have consistent preferences across commands. Additionally, there's no support for environment variables, which limits integration with CI/CD pipelines and scripting environments. Adding configuration file support with proper precedence handling will improve the developer experience by allowing users to set defaults once and override them as needed.

## What Changes

- Add YAML-based configuration file loading from `~/.config/spectr/config.yaml` (XDG-compliant)
- Implement environment variable support with `SPECTR_*` prefix (e.g., `SPECTR_STRICT=true`)
- Integrate Kong resolvers to establish precedence: CLI flags > Environment variables > User config file
- Add `spectr config` command with subcommands: `show`, `edit`, `init`, `validate`
- Create new `internal/config` package for configuration loading, merging, and validation
- Update `main.go` to wire Kong resolvers for configuration precedence
- Add YAML dependency: `gopkg.in/yaml.v3`

## Impact

- Affected specs: `cli-framework` (MODIFIED for Kong resolver integration), `config-management` (NEW capability)
- Affected code:
  - `main.go` - Kong initialization with resolvers
  - `cmd/root.go` - CLI struct may need config-aware fields
  - New `cmd/config.go` - Config command implementation
  - New `internal/config/` - Configuration package (loader, resolver, types, validator)
- Breaking changes: **None** - This is purely additive. All existing CLI behavior remains unchanged. Configuration files and environment variables are optional and only override defaults when provided.
