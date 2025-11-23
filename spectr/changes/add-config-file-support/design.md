# Design: User Configuration File Support

## Context

Spectr is a CLI tool built with Kong v1.13.0 for command parsing. Currently, all configuration happens through command-line flags. Users have requested the ability to set persistent preferences to avoid repeating common flags like `--strict`, `--json`, or `--interactive`. This design adds configuration file support while maintaining Kong's struct-based declarative approach and preserving full backward compatibility.

## Goals / Non-Goals

**Goals:**
- Support YAML configuration files at `~/.config/spectr/config.yaml` (XDG Base Directory compliant)
- Support environment variables with `SPECTR_*` prefix for CI/CD integration
- Establish clear precedence: CLI flags > Env vars > User config file
- Provide `spectr config` command for managing configuration
- Maintain 100% backward compatibility (all config is optional)
- Leverage Kong's built-in resolver mechanism for precedence handling

**Non-Goals:**
- System-wide configuration (`/etc/spectr/`) - not needed initially
- Project-level configuration (`.spectr.yaml`) - not in scope for this change
- TOML/JSON format support - YAML only keeps implementation simple
- Config migration from other tools
- Remote configuration fetching

## Decisions

### 1. Configuration File Format: YAML

**Decision:** Use YAML exclusively with `gopkg.in/yaml.v3`

**Rationale:**
- Most common format in Go CLI ecosystem (golangci-lint, docker-compose, k8s)
- Human-friendly for hand-editing
- Mature library with good error messages
- Single format reduces testing surface and dependency footprint

**Alternatives considered:**
- TOML: Less common in Go ecosystem, requires additional dependency
- JSON: Built-in but less human-friendly (no comments, stricter syntax)
- Multiple formats: Added complexity for minimal user benefit

### 2. Configuration Location: XDG Base Directory Specification

**Decision:** Use `~/.config/spectr/config.yaml` on Linux/macOS

**Rationale:**
- XDG Base Directory Specification is the standard on Linux
- Respects `$XDG_CONFIG_HOME` if set, defaults to `~/.config`
- Keeps user home directory clean
- Matches modern CLI tool conventions (gh, glab, kubectl)

**Implementation:**
```go
func GetUserConfigPath() (string, error) {
    configDir := os.Getenv("XDG_CONFIG_HOME")
    if configDir == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        configDir = filepath.Join(homeDir, ".config")
    }
    return filepath.Join(configDir, "spectr", "config.yaml"), nil
}
```

**Alternatives considered:**
- `~/.spectr.yaml`: Pollutes home directory, less modern
- `~/.spectrrc`: Old Unix convention, less discoverable
- Windows: Will use `%APPDATA%\spectr\config.yaml` for cross-platform consistency

### 3. Precedence Order

**Decision:** CLI flags > Environment variables > User config file > Hard-coded defaults

**Rationale:**
- CLI flags are most immediate and should always win (explicit user intent)
- Env vars enable CI/CD and one-off overrides without file changes
- User config provides persistent defaults across all projects
- Hard-coded defaults ensure the tool works out-of-box

**Kong Implementation:**
Use Kong's `Resolvers()` option with custom resolver that chains lookups:
```go
kong.Parse(&cli,
    kong.Resolvers(
        // 1. Environment variables (checked first by Kong's default resolver)
        envVarResolver(),
        // 2. User config file (checked second)
        userConfigResolver(),
    ),
)
```

Note: CLI flags automatically have highest precedence in Kong's architecture.

**Alternatives considered:**
- Project config before user config: Deferred to future change (not needed yet)
- System config: Deferred to future change (enterprise use case only)

### 4. Environment Variable Naming Convention

**Decision:** Use `SPECTR_` prefix with uppercase flag names, replacing hyphens with underscores

**Examples:**
- `--strict` → `SPECTR_STRICT=true`
- `--no-interactive` → `SPECTR_NO_INTERACTIVE=true`
- `--json` → `SPECTR_JSON=true`

**Rationale:**
- Standard convention in Go CLI tools (DOCKER_, KUBECTL_, GH_)
- Easy to discover with `env | grep SPECTR`
- Clear separation from other environment variables
- Kong supports this pattern with `kong.Vars()` and `envar` tags

**Implementation:**
Add `envar` tags to existing struct fields:
```go
type ValidateCmd struct {
    Strict bool `kong:"help='Treat warnings as errors',envar='SPECTR_STRICT'"`
}
```

### 5. Configuration File Schema

**Decision:** Mirror CLI struct tags exactly, using nested YAML structure matching command hierarchy

**Example config.yaml:**
```yaml
# Global defaults (apply to all commands)
json: false
interactive: true

# Command-specific overrides
validate:
  strict: true
  json: false

list:
  long: true
  specs: false

archive:
  yes: false
  skip-specs: false
```

**Rationale:**
- Intuitive mapping from CLI flags to config keys
- Supports command-specific overrides naturally
- Easy to validate against Kong struct definitions
- Allows future growth without breaking changes

**Alternatives considered:**
- Flat structure: Doesn't scale well with command-specific settings
- Separate files per command: Over-engineered for current needs

### 6. Config Command Structure

**Decision:** Implement `spectr config` with subcommands: `init`, `show`, `edit`, `validate`

**Subcommands:**
- `spectr config init` - Generate default config file with comments
- `spectr config show` - Display merged configuration (flags + env + file)
- `spectr config edit` - Open config file in `$EDITOR` (falls back to `vi`)
- `spectr config validate` - Check config file syntax and schema

**Rationale:**
- `init` reduces onboarding friction with scaffolded template
- `show` aids debugging precedence issues
- `edit` provides convenience (users can edit manually too)
- `validate` catches YAML syntax errors before runtime

**Alternatives considered:**
- Single `spectr config` command with flags: Less discoverable, less extensible
- No config command: Forces users to manually create/debug files

### 7. Kong Resolver Architecture

**Decision:** Create custom `kong.Resolver` that loads YAML config and provides values to Kong's parser

**Implementation sketch:**
```go
type userConfigResolver struct {
    config map[string]interface{}
}

func (r *userConfigResolver) Resolve(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
    // 1. Build config key from command path + flag name
    // 2. Look up value in r.config
    // 3. Return value if found, or return nil to try next resolver
}

func NewUserConfigResolver() (kong.Resolver, error) {
    configPath, err := GetUserConfigPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(configPath)
    if os.IsNotExist(err) {
        // No config file is OK, return empty resolver
        return &userConfigResolver{config: make(map[string]interface{})}, nil
    }
    // Parse YAML and return resolver
}
```

**Rationale:**
- Kong's resolver interface is designed for exactly this use case
- Separates config loading from command logic (clean architecture)
- Testable in isolation
- Allows adding more resolvers later (project config, system config)

## Risks / Trade-offs

### Risk: YAML Parsing Errors

**Mitigation:**
- Provide `spectr config validate` for pre-flight checks
- Use `gopkg.in/yaml.v3` strict mode to catch typos
- Include clear error messages with line numbers
- Ship example config with comments

### Risk: Precedence Confusion

**Mitigation:**
- Document precedence clearly in help text and README
- `spectr config show` displays effective merged config with source annotations
- Error messages indicate which source provided conflicting values (future enhancement)

### Risk: Environment Variable Name Collisions

**Mitigation:**
- Use `SPECTR_` prefix to namespace variables
- Document all supported env vars in `spectr config init` template comments
- Kong validates types automatically (e.g., `SPECTR_STRICT=invalid` will error)

### Trade-off: Single Format (YAML Only)

**Benefit:** Simpler implementation, fewer dependencies, easier testing
**Cost:** Users who prefer TOML/JSON must use YAML
**Justification:** YAML is de facto standard for Go CLIs. Can add more formats later if users request it.

### Trade-off: User Config Only (No System/Project Config)

**Benefit:** Smaller scope, faster delivery, less complexity
**Cost:** Teams can't enforce defaults via system config or version-control project config
**Justification:** User config solves 80% of use cases. System/project config can be added in future change without breaking changes.

## Migration Plan

### Phase 1: Implementation (This Change)
1. Add `internal/config` package with loader and resolver
2. Add `cmd/config.go` with all subcommands
3. Update `main.go` to wire resolvers
4. Add `envar` tags to all command structs
5. Write comprehensive tests (unit + integration)

### Phase 2: Documentation
1. Update README with configuration section
2. Add `docs/configuration.md` with examples
3. Update `--help` text to mention env vars and config files

### Phase 3: Rollout
1. Release as minor version bump (backward compatible)
2. Announce in release notes with examples
3. Monitor for issues/feedback

### Rollback Plan
If critical bugs are discovered:
- Config loading is isolated in `internal/config` package
- Can disable resolver in `main.go` without code changes (feature flag)
- Users without config files are unaffected (100% backward compatible)

## Open Questions

1. **Windows support:** Should we support Windows-specific config locations (`%APPDATA%`)?
   - **Decision deferred:** Implement XDG on Unix first, add Windows in follow-up if users request

2. **Config file permissions:** Should we validate that config files aren't world-readable?
   - **Decision deferred:** Not needed for initial release (config doesn't contain secrets)

3. **Config schema versioning:** Should config files include a version field for future migrations?
   - **Decision deferred:** YAGNI - add when we have actual breaking changes

4. **Partial config validation:** Should we warn about unknown keys in config files?
   - **Decision:** Yes - use YAML strict mode and validate against Kong struct tags in `spectr config validate`
