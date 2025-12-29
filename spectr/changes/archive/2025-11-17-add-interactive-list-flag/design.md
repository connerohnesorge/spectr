# Design: Interactive List with Clipboard Support

## Context

The `spectr list` command currently outputs changes and specs as plain text,
long-format text, or JSON. Users frequently need to reference IDs for subsequent
commands, requiring manual copy-paste or retyping. An interactive table
interface would improve the user experience by allowing visual browsing and
one-key clipboard copying.

**Constraints:**

- Must maintain backward compatibility with existing output formats
- Should work across platforms (Linux, macOS, Windows)
- Must work in SSH sessions and remote terminals
- Should use existing dependencies where possible (bubbletea, lipgloss already
  present)

**Stakeholders:**

- CLI users who frequently run `spectr show`, `spectr validate`, etc.
- Power users working in SSH sessions
- CI/CD systems (unaffected - they use `--json` flag)

## Goals / Non-Goals

**Goals:**

- Provide interactive table interface with `-I` flag for list command
- Enable single-keypress clipboard copy of item IDs
- Support standard keyboard navigation (arrows, j/k, Enter, q, Ctrl+C)
- Work across all major platforms and terminal environments
- Maintain existing text/JSON output modes unchanged

**Non-Goals:**

- Multi-select or batch clipboard operations (future enhancement)
- Inline editing or modification of items
- Integration with other commands beyond `list`
- Custom key bindings or configuration file
- Mouse support (keyboard-only for simplicity)

## Decisions

### Decision 1: Use Bubbletea Table Bubble

**Choice:** Use `github.com/charmbracelet/bubbles/table` for interactive table

**Rationale:**

- Already in dependencies (bubbletea is present)
- Provides built-in navigation, styling, and focus management
- Follows The Elm Architecture pattern (clean state management)
- Well-maintained and commonly used in Go CLI tools
- Reduces custom code for table rendering and navigation

**Alternatives considered:**

1. Custom table implementation with bubbletea primitives
   - Pro: Full control over behavior
   - Con: More code to maintain, reinventing the wheel
   - Con: More testing surface area
2. tview or other TUI framework
   - Pro: More feature-rich
   - Con: Heavier dependency, overkill for single table view
   - Con: Not already in dependencies

### Decision 2: Clipboard Strategy - Dual Approach

**Choice:** Use `atotto/clipboard` for desktop, fallback to OSC 52 for SSH

**Rationale:**

- Desktop environments (X11, Wayland, macOS, Windows) need native clipboard APIs
- SSH sessions without X11 forwarding benefit from OSC 52 terminal escape
  sequences
- `atotto/clipboard` is lightweight, cross-platform, and widely used
- `termenv` (already a transitive dependency) supports OSC 52
- Graceful degradation: try native first, fall back to OSC 52, then error

**Alternatives considered:**

1. OSC 52 only
   - Pro: Works everywhere with modern terminals
   - Con: Not all terminals support OSC 52
   - Con: Some users find terminal escape sequences less intuitive
2. Native clipboard only
   - Pro: Simpler implementation
   - Con: Fails in SSH sessions, frustrating for remote users
3. No clipboard, just print to stdout
   - Pro: Simplest implementation
   - Con: Doesn't solve the core UX problem

**Implementation:**

```go
func copyToClipboard(text string) error {
    // Try native clipboard first
    err := clipboard.WriteAll(text)
    if err == nil {
        return nil
    }

    // Fallback to OSC 52 for SSH sessions
    osc52 := "\x1b]52;c;" + base64.StdEncoding.EncodeToString([]byte(text)) + "\x07"
    fmt.Print(osc52)
    return nil // OSC 52 doesn't report errors
}
```

### Decision 3: Copy ID on Enter, Not Title

**Choice:** Pressing Enter copies the item's ID (kebab-case) to clipboard

**Rationale:**

- IDs are used for subsequent CLI commands (`spectr show <id>`)
- IDs are machine-readable and unambiguous
- Titles are human-readable but not actionable in CLI context
- User explicitly requested ID copying based on clarification

**Alternatives considered:**

1. Copy full row (ID + Title)
   - Pro: More information
   - Con: Not directly usable in CLI commands
2. Copy title only
   - Pro: Human-readable
   - Con: Not the data users need for commands

### Decision 4: Exit After Clipboard Copy

**Choice:** Interactive mode exits after successful Enter press

**Rationale:**

- Single-purpose interaction: user selects one item and is done
- Reduces cognitive load (no need to remember to quit)
- Matches typical CLI workflow (run command, get result, continue)
- 'q' and Ctrl+C still available if user changes mind

**Alternatives considered:**

1. Stay in interactive mode after copy
   - Pro: Allows multiple copies
   - Con: Requires explicit quit step
   - Con: Unclear when user is "done"
2. Add mode toggle (stay/exit)
   - Pro: Flexibility
   - Con: Adds complexity and configuration

### Decision 5: Table Columns

**For Changes:**

- ID (kebab-case identifier)
- Title (extracted from proposal.md)
- Deltas (count of spec changes)
- Tasks (completed/total format)

**For Specs:**

- ID (spec identifier)
- Title (from spec.md)
- Requirements (count)

**Rationale:**

- Matches existing `--long` output format
- Provides enough context for informed selection
- Keeps table width manageable
- All data already available in `ChangeInfo` and `SpecInfo` structs

## Risks / Trade-offs

### Risk 1: Clipboard Failures in Exotic Terminals

**Risk:** Some terminal emulators may not support native clipboard or OSC 52

**Mitigation:**

- Graceful error handling: display copied ID in output even if clipboard fails
- Document supported terminals and clipboard behavior
- Allow command to exit successfully even if clipboard fails (non-fatal error)

**Trade-off:** Accept that some edge-case terminals may not support clipboard,
but still provide value through interactive selection and visual display

### Risk 2: Terminal Size Constraints

**Risk:** Small terminal windows may not render table well

**Mitigation:**

- Use bubbles/table built-in width management
- Truncate long titles with ellipsis
- Set minimum reasonable terminal size requirement (e.g., 80 columns)
- Test with narrow terminals (60-80 cols) during implementation

**Trade-off:** Prioritize common terminal sizes (80+ cols) over edge cases (40
cols)

### Risk 3: Dependency on External Clipboard Library

**Risk:** `atotto/clipboard` adds external dependency and potential
platform-specific bugs

**Mitigation:**

- Library is mature, widely used, and actively maintained
- Fallback to OSC 52 reduces dependency on library
- Consider vendoring if supply chain concerns arise
- Monitor for security advisories

**Trade-off:** Added dependency complexity vs. significant UX improvement

### Risk 4: Interactive Mode Not Scriptable

**Risk:** CI/CD or scripts cannot use `-I` flag in automated contexts

**Mitigation:**

- Interactive mode requires TTY, will fail gracefully in non-interactive
  contexts
- Document that `-I` is for human users, not scripts
- Scripts should continue using `--json` or default text output

**Trade-off:** Interactive mode is explicitly for human interaction; automation
unaffected

## Migration Plan

**No migration required** - this is a new feature addition.

**Deployment steps:**

1. Merge code with `-I` flag implementation
2. Release new version with updated help text
3. Announce feature in release notes with examples
4. Document in README or CLI documentation

**Rollback plan:**

- If critical bugs found, flag can be removed in patch release
- No user data or config files affected
- Existing output modes (`--long`, `--json`) unchanged

**Compatibility:**

- Fully backward compatible
- No breaking changes to existing flags or output
- New flag is opt-in

## Open Questions

1. **Should we add a help overlay in interactive mode?**
   - Show "Press Enter to copy, q to quit" at bottom of screen
   - Decision deferred: Start without, add if users request it

2. **Should table height be configurable?**
   - Allow `--height N` or use full terminal height
   - Decision: Use full terminal height minus header/footer space (standard
     bubbletea behavior)

3. **Should we support mouse clicks in addition to keyboard?**
   - Bubbletea supports mouse events
   - Decision deferred: Keyboard-only for MVP, add mouse later if requested

4. **Should clipboard copy show a temporary success notification?**
   - E.g., "✓ Copied add-archive-command" for 1 second before exit
   - Decision: Yes, print message before exit for clarity

5. **Should we allow configuring which column to copy (ID vs. Title)?**
   - Could add flag or config option
   - Decision: Not for MVP - ID only based on user request
