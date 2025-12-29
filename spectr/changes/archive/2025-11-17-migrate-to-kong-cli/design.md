# Design: Kong CLI Migration

## Context

Spectr currently uses Cobra for CLI command handling, which requires manual command registration, flag binding, and extensive boilerplate code. As the CLI grows, this pattern becomes harder to maintain and extend. Kong offers a declarative alternative that reduces boilerplate and improves type safety.

## Goals / Non-Goals

### Goals

- Migrate from Cobra to Kong with zero breaking changes to user-facing CLI syntax
- Reduce boilerplate code in command definitions
- Improve type safety and maintainability of CLI code
- Maintain all existing functionality (interactive/non-interactive init, flag handling)
- Make future command additions simpler and more consistent

### Non-Goals

- Change the user-facing CLI interface or command syntax
- Add new commands or features (pure refactoring)
- Modify the internal init wizard logic or execution
- Change behavior of any existing flags or arguments

## Decisions

### Decision: Struct-Based Command Definition

Instead of Cobra's builder pattern with `cobra.Command` structs and `AddCommand` calls, we'll use Kong's struct tags to define the CLI structure declaratively.

**Before (Cobra):**

```go
var initCmd = &cobra.Command{
    Use:   "init [path]",
    Short: "Initialize Spectr",
    RunE:  runInit,
}

func init() {
    rootCmd.AddCommand(initCmd)
    initCmd.Flags().StringVarP(&initPath, "path", "p", "", "Project path")
}
```

**After (Kong):**

```go
type CLI struct {
    Init InitCmd `cmd:"" help:"Initialize Spectr in a project"`
}

type InitCmd struct {
    Path           string   `arg:"" optional:"" help:"Project path"`
    PathFlag       string   `name:"path" short:"p" help:"Project path (alternative to positional)"`
    Tools          []string `name:"tools" short:"t" help:"Tools to configure"`
    NonInteractive bool     `name:"non-interactive" help:"Run in non-interactive mode"`
}

func (c *InitCmd) Run() error {
    // Implementation here
}
```

**Rationale:** This approach reduces boilerplate, makes command structure more readable, and enables better IDE support for refactoring.

### Decision: Preserve Existing Flag Names and Behavior

All flag names, short flags, and argument positions will remain identical to maintain backward compatibility.

**Alternatives Considered:**

- Restructure flags for better consistency → Rejected to avoid breaking changes
- Mix Kong and Cobra temporarily → Rejected due to complexity and dependency bloat

### Decision: Single-Phase Migration

Migrate all commands in one change rather than incrementally.

**Rationale:** With only one command (`init`), incremental migration adds unnecessary complexity. A single-phase migration is cleaner and easier to review.

## Technical Approach

### CLI Structure

```go
// cmd/root.go
type CLI struct {
    Init InitCmd `cmd:"" help:"Initialize Spectr in a project"`
    // Future commands will be added as struct fields here
}

// main.go
func main() {
    cli := &cmd.CLI{}
    ctx := kong.Parse(cli,
        kong.Name("spectr"),
        kong.Description("Validatable spec-driven development"),
        kong.UsageOnError(),
    )
    err := ctx.Run()
    ctx.FatalIfErrorf(err)
}
```

### Flag Mapping Strategy

| Cobra Flag | Kong Struct Tag | Type |
|------------|----------------|------|
| `--path, -p` | `name:"path" short:"p"` | `string` |
| `--tools, -t` | `name:"tools" short:"t"` | `[]string` |
| `--non-interactive` | `name:"non-interactive"` | `bool` |
| Positional `[path]` | `arg:"" optional:""` | `string` |

### Method Dispatch

Kong will automatically call `Run()` method on the selected command struct, eliminating the need for manual RunE function wiring.

## Risks / Trade-offs

### Risk: Behavior Differences

Kong and Cobra may parse flags slightly differently.

**Mitigation:** Comprehensive testing of all flag combinations and edge cases before deployment.

### Risk: Help Text Formatting Changes

Kong's help output format differs from Cobra's.

**Mitigation:** This is cosmetic and unlikely to impact users. Can be customized if needed using Kong's template options.

### Trade-off: Community Size

Cobra has a larger community than Kong.

**Mitigation:** Kong is well-maintained, stable, and simpler to understand. Less community support is offset by reduced need for support due to simplicity.

## Migration Plan

### Phase 1: Setup

1. Add Kong dependency
2. Remove Cobra and pflag dependencies
3. Run `go mod tidy`

### Phase 2: Implementation

1. Create new CLI struct in `cmd/root.go`
2. Convert init command to Kong struct
3. Update main.go to use Kong parser
4. Remove old Cobra code

### Phase 3: Validation

1. Run all existing tests
2. Manual testing of all command variations
3. Verify help text is acceptable
4. Check edge cases (invalid flags, missing args)

### Rollback Plan

If critical issues are discovered:

1. Revert commits
2. Restore Cobra dependencies
3. Original code is preserved in git history

## Open Questions

None - the scope is clear and well-defined.
