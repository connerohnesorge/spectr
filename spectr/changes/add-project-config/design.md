## Context

Spectr currently hardcodes `"spectr"` as the root directory name in ~50 locations across the codebase. Users have requested the ability to rename this folder to match their project conventions (e.g., `specs/`, `.spectr/`, `openspec/`).

The configuration system must:
1. Be optional - existing projects work without any config file
2. Be discoverable - CLI finds config from any subdirectory
3. Be minimal - only settings that genuinely need customization

## Goals / Non-Goals

**Goals:**
- Allow renaming the spectr root directory
- Maintain backward compatibility with unconfigured projects
- Keep configuration minimal and optional
- Support discovery from any working directory within the project

**Non-Goals:**
- Complex configuration schema with many options
- Per-command configuration overrides
- Environment variable configuration (keep it simple)
- Remote/network configuration sources

## Decisions

### Decision: Config file format and location

**Choice:** Single `spectr.yaml` file at project root (next to `spectr/` or custom root)

**Rationale:**
- YAML is human-readable and familiar to developers
- Single file location is predictable and simple
- Project root placement follows convention (like `.gitignore`, `package.json`)

**Alternatives considered:**
- `.spectr.yaml` (hidden) - Less discoverable, harder to edit
- `spectr/config.yaml` (inside root) - Chicken-and-egg problem: can't find root without knowing root
- TOML format - Less familiar, Go support less mature than YAML

### Decision: Config schema

**Choice:** Minimal schema with only `root_dir` field initially

```yaml
# spectr.yaml
root_dir: specs  # Optional, defaults to "spectr"
```

**Rationale:**
- Start minimal, expand when needed
- Only include settings users have requested
- Avoid over-engineering

### Decision: Config discovery strategy

**Choice:** Walk up directory tree from current working directory looking for `spectr.yaml`

**Algorithm:**
1. Check current directory for `spectr.yaml`
2. If found, use `root_dir` value (or default)
3. If not found, check parent directory
4. Repeat until filesystem root
5. If no `spectr.yaml` found, use default `"spectr"` directory name
6. Verify the root directory exists (validation)

**Rationale:**
- Works from any subdirectory in project
- Matches behavior of tools like git, npm, cargo
- Falls back gracefully when no config exists

### Decision: Caching strategy

**Choice:** No caching - load config on each command invocation

**Rationale:**
- Config file is small and fast to parse
- Avoids cache invalidation complexity
- Commands are short-lived, not long-running processes

### Decision: Dependency for YAML parsing

**Choice:** `gopkg.in/yaml.v3`

**Rationale:**
- Standard Go YAML library
- Already widely used in the ecosystem
- Good error messages for invalid YAML

## Risks / Trade-offs

**Risk:** Config file discovery adds overhead to every command
**Mitigation:** File stat operations are fast; measure impact and optimize if needed

**Risk:** Users might create config with invalid root_dir
**Mitigation:** Validate that root directory exists; clear error messages

**Trade-off:** Simplicity vs flexibility
**Choice:** Start simple with just `root_dir`; add more settings only when requested

## Migration Plan

1. Add config package with backward-compatible defaults
2. Migrate one module at a time (discovery first)
3. Run full test suite after each migration
4. No breaking changes to existing projects

## Open Questions

- Should `spectr init` prompt for custom root directory name? (Defer: can add later)
- Should config support comments/documentation? (Yes, YAML supports comments natively)
