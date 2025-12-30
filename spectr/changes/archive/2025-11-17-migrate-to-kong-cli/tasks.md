# Implementation Tasks

## 1. Dependency Management

- [x] 1.1 Add `github.com/alecthomas/kong` to go.mod
- [x] 1.2 Remove `github.com/spf13/cobra` from go.mod
- [x] 1.3 Remove `github.com/spf13/pflag` from go.mod (Cobra dependency)
- [x] 1.4 Run `go mod tidy` to clean dependencies

## 2. Define Kong CLI Structure

- [x] 2.1 Create CLI struct in `cmd/root.go` with Kong struct tags
- [x] 2.2 Define global flags using struct fields and tags
- [x] 2.3 Add Init subcommand struct with appropriate tags
- [x] 2.4 Configure help text and descriptions using Kong tags

## 3. Implement Init Command with Kong

- [x] 3.1 Convert init command flags to struct fields with Kong tags
- [x] 3.2 Implement Run method for init command
- [x] 3.3 Preserve existing init functionality (interactive/non-interactive
  modes)
- [x] 3.4 Ensure path handling and validation works identically

## 4. Update Main Entry Point

- [x] 4.1 Replace Cobra Execute() with Kong parser initialization
- [x] 4.2 Add Kong parser options (help formatting, etc.)
- [x] 4.3 Implement proper error handling for Kong parsing errors
- [x] 4.4 Ensure exit codes are preserved

## 5. Testing and Validation

- [x] 5.1 Test `spectr` command with no arguments (should show help)
- [x] 5.2 Test `spectr init` in interactive mode
- [x] 5.3 Test `spectr init --non-interactive --tools all`
- [x] 5.4 Test `spectr init /path/to/project`
- [x] 5.5 Verify help text displays correctly
- [x] 5.6 Run existing tests to ensure no regressions
- [x] 5.7 Update any affected tests

## 6. Documentation

- [x] 6.1 Update README.md if CLI examples exist
- [x] 6.2 Verify help text is clear and accurate
- [x] 6.3 Document any breaking changes in command syntax (if any)
