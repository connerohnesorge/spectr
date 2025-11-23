# Implementation Tasks

## 1. Dependencies and Setup
- [ ] 1.1 Add `gopkg.in/yaml.v3` to `go.mod`
- [ ] 1.2 Run `go mod tidy` to update dependencies
- [ ] 1.3 Verify YAML library imports correctly

## 2. Create `internal/config` Package Structure
- [ ] 2.1 Create `internal/config/` directory
- [ ] 2.2 Create `internal/config/types.go` with config struct types
- [ ] 2.3 Create `internal/config/loader.go` for loading YAML files
- [ ] 2.4 Create `internal/config/resolver.go` for Kong resolver implementation
- [ ] 2.5 Create `internal/config/paths.go` for XDG directory handling
- [ ] 2.6 Create `internal/config/validator.go` for config validation logic

## 3. Implement Configuration File Loading
- [ ] 3.1 Implement `GetUserConfigPath()` with XDG Base Directory support
- [ ] 3.2 Implement `LoadConfig(path string)` to read and parse YAML
- [ ] 3.3 Add error handling for missing files (return empty config, not error)
- [ ] 3.4 Add error handling for malformed YAML (return error with line numbers)
- [ ] 3.5 Add support for nested command-specific configuration
- [ ] 3.6 Implement home directory expansion (`~`) in path handling

## 4. Implement Kong Resolver
- [ ] 4.1 Create `userConfigResolver` struct implementing `kong.Resolver` interface
- [ ] 4.2 Implement `Resolve()` method to look up flag values from config
- [ ] 4.3 Handle command path to config key mapping (e.g., `validate.strict`)
- [ ] 4.4 Return `nil` when value not found (allows fallback to next resolver)
- [ ] 4.5 Implement type coercion from YAML values to Go types
- [ ] 4.6 Add unit tests for resolver with various config structures

## 5. Add Environment Variable Support
- [ ] 5.1 Add `envar` tags to all command struct fields in `cmd/` files
- [ ] 5.2 Update `cmd/root.go` - CLI struct fields
- [ ] 5.3 Update `cmd/list.go` - ListCmd struct fields
- [ ] 5.4 Update `cmd/validate.go` - ValidateCmd struct fields
- [ ] 5.5 Update `cmd/archive.go` - ArchiveCmd struct fields
- [ ] 5.6 Update `cmd/view.go` - ViewCmd struct fields
- [ ] 5.7 Verify environment variable naming follows `SPECTR_` prefix convention

## 6. Wire Kong Resolvers in main.go
- [ ] 6.1 Import `internal/config` package in `main.go`
- [ ] 6.2 Load user config before calling `kong.Parse()`
- [ ] 6.3 Create user config resolver instance
- [ ] 6.4 Add `kong.Resolvers()` option with config resolver
- [ ] 6.5 Handle config loading errors gracefully (log and continue with defaults)
- [ ] 6.6 Test that CLI flags still override config values

## 7. Implement `spectr config init` Command
- [ ] 7.1 Create `cmd/config.go` with ConfigCmd struct
- [ ] 7.2 Add InitCmd subcommand struct with `--force` flag
- [ ] 7.3 Implement `Run()` method to scaffold default config file
- [ ] 7.4 Create config template with commented examples for all commands
- [ ] 7.5 Add directory creation logic (`~/.config/spectr/`)
- [ ] 7.6 Add overwrite confirmation prompt (unless `--force`)
- [ ] 7.7 Write success message with file path
- [ ] 7.8 Add unit tests for init command

## 8. Implement `spectr config show` Command
- [ ] 8.1 Add ShowCmd subcommand struct with `--json` and `--command` flags
- [ ] 8.2 Implement `Run()` method to display merged configuration
- [ ] 8.3 Merge config from all sources (defaults, file, env vars, CLI flags)
- [ ] 8.4 Annotate each value with its source (default/config/env/flag)
- [ ] 8.5 Format output as human-readable text by default
- [ ] 8.6 Add JSON output mode with `--json` flag
- [ ] 8.7 Add command filtering with `--command` flag
- [ ] 8.8 Add unit tests for show command

## 9. Implement `spectr config edit` Command
- [ ] 9.1 Add EditCmd subcommand struct with `--editor` flag
- [ ] 9.2 Implement `Run()` method to open config file in editor
- [ ] 9.3 Read `$EDITOR` environment variable (fallback to `vi`)
- [ ] 9.4 Support `--editor` flag to override `$EDITOR`
- [ ] 9.5 Handle missing config file (prompt to init first)
- [ ] 9.6 Use `os/exec` to launch editor and wait for exit
- [ ] 9.7 Add error handling for editor launch failures
- [ ] 9.8 Add unit tests for edit command (mock editor execution)

## 10. Implement `spectr config validate` Command
- [ ] 10.1 Add ValidateConfigCmd subcommand struct with `--strict` flag
- [ ] 10.2 Implement `Run()` method to validate config file
- [ ] 10.3 Check YAML syntax and report line numbers for errors
- [ ] 10.4 Validate keys against known CLI flags from Kong structs
- [ ] 10.5 Validate value types match expected flag types
- [ ] 10.6 Warn about unknown keys (error if `--strict` enabled)
- [ ] 10.7 Report all validation issues with clear messages
- [ ] 10.8 Exit with non-zero code if validation fails
- [ ] 10.9 Add unit tests for validate command

## 11. Wire Config Command into Root CLI
- [ ] 11.1 Add `Config ConfigCmd` field to CLI struct in `cmd/root.go`
- [ ] 11.2 Add `cmd` and `help` struct tags for config command
- [ ] 11.3 Verify `spectr config --help` displays subcommands
- [ ] 11.4 Verify all subcommands appear in help text

## 12. Write Unit Tests
- [ ] 12.1 Test config file loading with valid YAML
- [ ] 12.2 Test config file loading with invalid YAML
- [ ] 12.3 Test config file loading with missing file
- [ ] 12.4 Test XDG_CONFIG_HOME environment variable override
- [ ] 12.5 Test Kong resolver with nested command config
- [ ] 12.6 Test Kong resolver with global config
- [ ] 12.7 Test precedence: CLI flag > env var > config file
- [ ] 12.8 Test type coercion for booleans, strings, slices
- [ ] 12.9 Test error handling for type mismatches
- [ ] 12.10 Test config path discovery and expansion

## 13. Write Integration Tests
- [ ] 13.1 Create integration test with temp config file and commands
- [ ] 13.2 Test `spectr config init` creates file with correct content
- [ ] 13.3 Test `spectr config show` displays merged config
- [ ] 13.4 Test `spectr validate --strict` with config file setting `strict: true`
- [ ] 13.5 Test environment variable overrides config file
- [ ] 13.6 Test CLI flag overrides environment variable
- [ ] 13.7 Test config validation catches unknown keys
- [ ] 13.8 Test config validation catches type mismatches

## 14. Update Documentation
- [ ] 14.1 Add configuration section to main README.md
- [ ] 14.2 Document YAML config file format and location
- [ ] 14.3 Document environment variable naming convention
- [ ] 14.4 Document precedence order (CLI > Env > Config > Default)
- [ ] 14.5 Add examples for common configuration scenarios
- [ ] 14.6 Document all `spectr config` subcommands
- [ ] 14.7 Update `--help` text for commands to mention env vars and config

## 15. Manual Testing and Validation
- [ ] 15.1 Test `spectr config init` on clean system
- [ ] 15.2 Test editing config file with `spectr config edit`
- [ ] 15.3 Test `spectr config show` displays correct values
- [ ] 15.4 Test `spectr config validate` catches errors
- [ ] 15.5 Test precedence by setting same flag in config, env, and CLI
- [ ] 15.6 Test with missing config file (should work with defaults)
- [ ] 15.7 Test with malformed YAML (should show clear error)
- [ ] 15.8 Test XDG_CONFIG_HOME override
- [ ] 15.9 Run `spectr validate add-config-file-support --strict`
- [ ] 15.10 Verify backward compatibility (no breaking changes to existing CLI usage)

## 16. Code Quality and Linting
- [ ] 16.1 Run `go fmt` on all new code
- [ ] 16.2 Run `golangci-lint run` and fix all issues
- [ ] 16.3 Add doc comments to all exported types and functions
- [ ] 16.4 Ensure error messages are clear and actionable
- [ ] 16.5 Check test coverage for `internal/config` package (aim for >80%)

## 17. Final Validation
- [ ] 17.1 Run all unit tests: `go test ./...`
- [ ] 17.2 Run integration tests
- [ ] 17.3 Build binary: `go build -o spectr`
- [ ] 17.4 Test binary manually with various config scenarios
- [ ] 17.5 Verify `spectr validate add-config-file-support --strict` passes
- [ ] 17.6 Review all spec requirements are met
- [ ] 17.7 Update this checklist with `[x]` for completed tasks
