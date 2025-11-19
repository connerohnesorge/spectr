# Design: VHS Demo Assets

## Context

Spectr's documentation currently lacks visual demonstrations of the CLI in action. Users must mentally simulate command execution from text examples, which creates friction during onboarding. VHS (https://github.com/charmbracelet/vhs) provides a solution: scriptable terminal recordings that can be version-controlled, automated, and regenerated as the CLI evolves.

**Background:**
- VHS uses `.tape` files (text scripts) to automate terminal recording
- Outputs can be GIF, MP4, WebM, or PNG sequences
- Tape files are deterministic and can be committed to version control
- GIFs can be embedded in Markdown (GitHub, Astro docs)

**Stakeholders:**
- New users evaluating Spectr (need quick visual understanding)
- Contributors learning the workflow (need accurate examples)
- Maintainers updating docs (need reproducible demos)

## Goals / Non-Goals

**Goals:**
- Provide visual demonstrations of core Spectr workflows
- Make demos reproducible and version-controlled via VHS tapes
- Enhance README and docs site with embedded GIFs
- Enable easy regeneration of demos when CLI changes
- Keep GIF file sizes reasonable (< 5MB each)

**Non-Goals:**
- Video tutorials or narrated screencasts (out of scope)
- Interactive demos or web-based playgrounds
- Demonstrating every possible flag combination (focus on common paths)
- Replacing text documentation (GIFs complement, don't replace)

## Decisions

### Directory Structure
```
assets/
├── logo.png              # Existing
├── vhs/                  # NEW: VHS tape source files
│   ├── init.tape
│   ├── list.tape
│   ├── validate.tape
│   ├── archive.tape
│   └── workflow.tape
└── gifs/                 # NEW: Generated GIF outputs
    ├── init.gif
    ├── list.gif
    ├── validate.gif
    ├── archive.gif
    └── workflow.gif
```

**Rationale:**
- Keep source (`.tape`) separate from generated (`.gif`) for clarity
- All assets in `assets/` directory maintains consistency
- Both tapes and GIFs committed to repo (tapes for source, GIFs for users without VHS)

### Demo Coverage

**Essential Demos** (MVP):
1. **init.tape**: `spectr init` with wizard → show created structure
2. **list.tape**: `spectr list` and `spectr list --specs` → show output
3. **validate.tape**: `spectr validate` showing error → fix → success
4. **archive.tape**: `spectr archive` showing merge and move
5. **workflow.tape**: Complete flow from proposal to archive

**Rationale:**
- Covers the three-stage workflow comprehensively
- Shows both success and error cases (validate.tape)
- Demonstrates interactive and non-interactive modes
- Balances coverage with maintainability (5 tapes is manageable)

### VHS Configuration

**Standard Settings** (applied to all tapes):
```elixir
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 10
Set Theme "Catppuccin Mocha"  # Dark theme, good readability
Set TypingSpeed 50ms
```

**Rationale:**
- 1200x600 is readable but not too large (file size)
- Catppuccin Mocha theme is popular, modern, and accessible
- 50ms typing speed feels realistic without being slow
- Consistent theme across all demos creates professional appearance

### Regeneration Workflow

**Makefile Target:**
```makefile
.PHONY: gifs
gifs:
	@echo "Generating GIFs from VHS tapes..."
	@command -v vhs >/dev/null 2>&1 || { echo "VHS not installed. See: https://github.com/charmbracelet/vhs"; exit 1; }
	@for tape in assets/vhs/*.tape; do \
		vhs $$tape; \
	done
	@mv assets/vhs/*.gif assets/gifs/
	@echo "✓ GIFs generated in assets/gifs/"
```

**Alternatives Considered:**
- **Shell script**: More portable than Makefile, but Makefile is more conventional for Go projects
- **Go task runner**: Overkill for simple file generation
- **Manual execution**: Error-prone and doesn't scale

**Decision:** Use Makefile target for simplicity and consistency with Go ecosystem conventions.

### README Integration

**Placement Strategy:**
- **Hero section**: Large workflow.gif at the top (after logo)
- **Quick Start**: Inline GIFs for each step (init, validate, archive)
- **Command Reference**: One GIF per command section

**Markdown Format:**
```markdown
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="assets/gifs/init.gif">
  <source media="(prefers-color-scheme: light)" srcset="assets/gifs/init.gif">
  <img width="800" alt="Spectr init demo" src="assets/gifs/init.gif">
</picture>
```

**Rationale:**
- `<picture>` tag supports theme-aware images (future: light theme GIFs)
- Width constraint (800px) prevents oversized rendering
- Alt text improves accessibility
- Currently using same GIF for dark/light (can add light theme GIFs later)

### Docs Site Integration

**Pages to Update:**
- `/getting-started/installation` → Add init.gif
- `/getting-started/quick-start` → Add workflow.gif
- `/concepts/spec-driven-development` → Add workflow.gif
- `/reference/*` → Add command-specific GIFs

**Astro Integration:**
```mdx
import { Image } from 'astro:assets';
import workflowGif from '../../../assets/gifs/workflow.gif';

<Image src={workflowGif} alt="Spectr workflow demo" width={800} />
```

**Rationale:**
- Astro's Image component optimizes assets
- Maintains consistency with docs site architecture
- GIFs in `assets/` are accessible from docs via relative paths

## Risks / Trade-offs

### Risk: Large Binary Files in Git
**Impact**: GIFs bloat repository size over time
**Mitigation**:
- Optimize GIFs before committing (gifsicle, ImageOptim)
- Target < 5MB per GIF
- Monitor repo size; consider Git LFS if needed (future)
**Trade-off**: Accepted for ease of use (users don't need VHS)

### Risk: Demos Become Outdated
**Impact**: CLI changes make GIFs inaccurate
**Mitigation**:
- Tape files are version-controlled (easy to regenerate)
- Add `make gifs` to release checklist
- Consider CI check that tapes still execute (future)
**Trade-off**: Manual regeneration burden vs. accuracy

### Risk: VHS Not Installed for Contributors
**Impact**: Contributors can't regenerate GIFs
**Mitigation**:
- Document VHS installation clearly
- Generated GIFs are committed (most contributors don't need to regenerate)
- Make VHS optional dependency (only needed for demo changes)
**Trade-off**: Acceptable; demos updated infrequently

### Risk: Theme/Style Consistency
**Impact**: Different tapes might use different themes/settings
**Mitigation**:
- Document standard VHS settings in this design doc
- Create shared `config.tape` that other tapes can source
- Code review checks for consistency
**Trade-off**: None; discipline is sufficient

## Migration Plan

**N/A** - This is a purely additive change with no migration required.

## Open Questions

1. **Should we add light theme GIFs?**
   - **Status**: Deferred to future iteration
   - **Reasoning**: Dark theme covers most users; light theme can be added later
   - **Action**: Use same GIF for both `prefers-color-scheme` for now

2. **Should we add CI verification of tapes?**
   - **Status**: Optional (task 6.1)
   - **Reasoning**: Adds value but not critical for MVP
   - **Action**: Implement if time permits; otherwise defer

3. **What about MP4/WebM for docs site?**
   - **Status**: Deferred
   - **Reasoning**: GIFs work everywhere; video formats need player UI
   - **Action**: Stick with GIFs for simplicity

4. **Should we document the tape scripts?**
   - **Status**: Yes, via comments in `.tape` files
   - **Reasoning**: Tape files are code; comments help maintainability
   - **Action**: Add descriptive comments to each tape file
