# Design: Browser-Based Markdown Renderer for Spectr

## Context

Spectr specifications and change proposals are complex Markdown documents that benefit from rich rendering with images, math notation, and proper typography. This design document outlines the technical approach for a browser-based renderer using Goldmark and MathJax.

**Background:**
- Specs contain requirements with scenarios in structured Markdown format
- Changes contain multiple files (proposal.md, tasks.md, design.md, delta specs)
- Users need to review specs visually with rendered images and math
- Current workflow requires manual conversion or external tools

**Constraints:**
- Must work on Linux, macOS, and Windows
- Must not require external dependencies beyond Go toolchain
- Must handle relative image paths correctly
- Must render MathJax without network access issues (CDN fallback)

**Stakeholders:**
- Developers reviewing specs and changes
- Product managers and stakeholders reviewing proposals
- Technical writers authoring specifications

## Goals / Non-Goals

**Goals:**
- Render specs and changes as polished HTML in browser
- Support MathJax for mathematical notation (`$inline$` and `$$display$$`)
- Resolve image paths relative to spec/change directory and project root
- Start ephemeral HTTP server with auto-assigned port
- Auto-open browser to rendered content
- Clean, minimal CSS for readability

**Non-Goals:**
- Live reload on file changes (future enhancement)
- PDF export (future enhancement)
- Editing specs in browser (read-only view)
- Custom themes or CSS customization (future enhancement)
- Offline MathJax bundling (CDN is acceptable)

## Decisions

### Decision 1: Use Goldmark for Markdown Parsing

**Chosen:** Goldmark (`github.com/yuin/goldmark`)

**Rationale:**
- De facto standard Markdown parser in Go ecosystem
- CommonMark compliant with extensions support
- Used by Hugo, GitHub, and other major projects
- Excellent performance and correctness
- Supports tables, strikethrough, task lists via extensions

**Alternatives considered:**
- `blackfriday`: Older, less actively maintained
- `gomarkdown/markdown`: Fewer features, less ecosystem adoption

### Decision 2: HTTP Server with Auto-Assigned Port

**Chosen:** Ephemeral HTTP server with port auto-assignment (`:0`)

**Rationale:**
- Avoids port conflicts (user may have multiple specs open)
- Simple stdlib implementation with `net/http`
- Allows MathJax and images to load without CORS issues (unlike `file://` protocol)
- Server URL printed to terminal for manual browser access

**Alternatives considered:**
- Static HTML file with `file://` protocol: CORS issues prevent MathJax from loading
- Fixed port (e.g., 8080): Port conflicts when multiple instances run

### Decision 3: MathJax from CDN

**Chosen:** MathJax 3.x from CDN (https://cdn.jsdelivr.net/npm/mathjax@3/...)

**Rationale:**
- No bundling required (keeps binary size small)
- Always up-to-date with latest MathJax version
- Widely available CDN with high reliability
- Graceful degradation if CDN unreachable (shows raw LaTeX)

**Alternatives considered:**
- Bundle MathJax in binary: 10+ MB increase in binary size
- KaTeX: Faster but less feature-complete for complex math

**Configuration:**
```javascript
MathJax = {
  tex: {
    inlineMath: [['$', '$']],
    displayMath: [['$$', '$$']]
  }
};
```

### Decision 4: Image Path Resolution Strategy

**Chosen:** Multi-pass resolution with fallback

1. Resolve relative to spec/change directory (e.g., `spectr/specs/auth/diagram.png`)
2. Fallback to project root (e.g., `./assets/diagram.png`)
3. Return 404 if neither exists

**Rationale:**
- Supports both spec-local images and shared project assets
- Predictable path resolution for users
- Allows organizing images alongside specs or centrally

**Implementation:**
```go
// HTTP handler for /images/<path>
func (s *Server) serveImage(w http.ResponseWriter, r *http.Request, imagePath string) {
    // Try spec/change directory first
    fullPath := filepath.Join(s.itemDir, imagePath)
    if fileExists(fullPath) {
        http.ServeFile(w, r, fullPath)
        return
    }
    // Fallback to project root
    fullPath = filepath.Join(s.projectRoot, imagePath)
    if fileExists(fullPath) {
        http.ServeFile(w, r, fullPath)
        return
    }
    http.NotFound(w, r)
}
```

### Decision 5: Combined View for Changes

**Chosen:** Render all change files in tabbed sections

**Structure:**
- **Proposal** tab: Rendered `proposal.md`
- **Tasks** tab: Rendered `tasks.md` with checkboxes
- **Design** tab: Rendered `design.md` (if exists)
- **Deltas** tab: All delta specs combined with capability headers

**Rationale:**
- Single-page view eliminates need to open multiple files
- Tab navigation provides organized access to each section
- Capability headers in deltas tab clarify which spec is affected

**HTML structure:**
```html
<div class="tabs">
  <button class="tab active" onclick="showTab('proposal')">Proposal</button>
  <button class="tab" onclick="showTab('tasks')">Tasks</button>
  <button class="tab" onclick="showTab('design')">Design</button>
  <button class="tab" onclick="showTab('deltas')">Deltas</button>
</div>
<div id="proposal" class="tab-content active">...</div>
<div id="tasks" class="tab-content">...</div>
<div id="design" class="tab-content">...</div>
<div id="deltas" class="tab-content">...</div>
```

### Decision 6: Auto-Open Browser

**Chosen:** Use `browser.OpenURL()` from `pkg/browser` or similar

**Rationale:**
- Cross-platform browser opening without platform-specific code
- Respects user's default browser preference
- Graceful fallback if browser opening fails (print URL)

**Implementation:**
```go
url := fmt.Sprintf("http://localhost:%d", port)
fmt.Printf("Opening browser at %s\n", url)
if err := browser.OpenURL(url); err != nil {
    fmt.Printf("Failed to open browser: %v\n", err)
    fmt.Printf("Please open manually: %s\n", url)
}
```

**Library:** `github.com/pkg/browser` or equivalent

## Architecture

### Package Structure

```
internal/show/
├── renderer.go       # Goldmark setup, Markdown → HTML conversion
├── server.go         # HTTP server with /images handler
├── templates.go      # HTML template with MathJax, CSS
├── item_loader.go    # Load spec or change files
└── browser.go        # Browser opening utility
```

### Data Flow

```
User runs: spectr show auth --type spec
                ↓
         cmd/show.go (ShowCmd.Run)
                ↓
    Determine item type and paths
                ↓
    show.LoadItem(itemID, itemType, projectPath)
                ↓
  Read Markdown files (spec.md or proposal.md + tasks.md + ...)
                ↓
    show.RenderMarkdown(content) using Goldmark
                ↓
    show.WrapInTemplate(html) with MathJax + CSS
                ↓
    show.StartServer(html, itemDir, projectRoot)
                ↓
  HTTP server starts on :0 (auto-assigned port)
                ↓
    browser.OpenURL(http://localhost:<port>)
                ↓
  Server serves /index.html and /images/<path>
                ↓
  User views rendered content in browser
                ↓
  User presses Ctrl+C to stop server
```

### HTTP Routes

- `GET /` → Rendered HTML content
- `GET /images/<path>` → Image files (with fallback resolution)

### HTML Template

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <script src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
    <script>
        MathJax = {
            tex: {
                inlineMath: [['$', '$']],
                displayMath: [['$$', '$$']]
            }
        };
    </script>
    <style>
        /* Clean, readable CSS */
        body { max-width: 900px; margin: 40px auto; font-family: system-ui; line-height: 1.6; }
        h1, h2, h3, h4 { margin-top: 1.5em; }
        pre { background: #f5f5f5; padding: 1em; overflow-x: auto; }
        code { background: #f5f5f5; padding: 0.2em 0.4em; border-radius: 3px; }
        img { max-width: 100%; height: auto; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background: #f5f5f5; }
        .tabs { border-bottom: 1px solid #ddd; margin-bottom: 20px; }
        .tab { padding: 10px 20px; border: none; background: none; cursor: pointer; }
        .tab.active { border-bottom: 2px solid #0066cc; color: #0066cc; }
        .tab-content { display: none; }
        .tab-content.active { display: block; }
    </style>
</head>
<body>
    {{if .HasTabs}}
        <div class="tabs">
            {{range .Tabs}}
                <button class="tab {{if .Active}}active{{end}}" onclick="showTab('{{.ID}}')">{{.Label}}</button>
            {{end}}
        </div>
        {{range .TabContents}}
            <div id="{{.ID}}" class="tab-content {{if .Active}}active{{end}}">
                {{.HTML}}
            </div>
        {{end}}
        <script>
            function showTab(tabID) {
                document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
                document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
                event.target.classList.add('active');
                document.getElementById(tabID).classList.add('active');
            }
        </script>
    {{else}}
        {{.HTML}}
    {{end}}
</body>
</html>
```

## Risks / Trade-offs

### Risk: CDN Unavailability

**Description:** MathJax CDN may be unreachable in air-gapped environments

**Impact:** Math notation displays as raw LaTeX (readable but not rendered)

**Mitigation:**
- Document requirement for internet access in README
- Future enhancement: `--offline` flag to use bundled MathJax

### Risk: Port Conflicts

**Description:** Auto-assigned port may still conflict if OS assigns already-used port

**Impact:** Server fails to start

**Mitigation:**
- Use `:0` for OS-level port assignment (avoids most conflicts)
- Retry with new port if bind fails (up to 3 attempts)
- Print clear error message with retry instructions

### Risk: Browser Opening Failures

**Description:** `browser.OpenURL()` may fail on headless systems or restricted environments

**Impact:** User must manually copy URL to browser

**Mitigation:**
- Always print URL to terminal as fallback
- Graceful degradation with clear instructions
- Document usage for SSH/remote environments

### Risk: Large Markdown Files

**Description:** Very large specs (>10MB) may be slow to render

**Impact:** Browser may hang or be sluggish

**Mitigation:**
- Goldmark is fast enough for typical specs (<1MB)
- Add warning if file size >5MB (future enhancement)
- Consider pagination for large changes (future enhancement)

## Migration Plan

No migration required. This is a new command with no breaking changes.

**Rollout:**
1. Implement `spectr show` command
2. Add tests for renderer, server, and path resolution
3. Update documentation with examples
4. Add to CLI help text and README

**Testing:**
- Unit tests for Markdown rendering
- Integration tests for HTTP server
- Manual testing with sample specs containing images and math

**Rollback:**
- If issues arise, command can be disabled or removed without affecting existing workflows

## Open Questions

1. **Syntax highlighting**: Should we add syntax highlighting for code blocks?
   - **Answer**: Yes, use Goldmark's syntax highlighting extension with Chroma
2. **Table of contents**: Should we generate TOC for long specs?
   - **Answer**: Future enhancement, not in initial implementation
3. **Dark mode**: Should we support dark mode CSS?
   - **Answer**: Future enhancement, start with light mode only
4. **Spec linking**: Should we support links between specs (e.g., `[auth spec](spec://auth)`)?
   - **Answer**: Future enhancement, use standard Markdown links for now
