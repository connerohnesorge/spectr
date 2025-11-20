# Implementation Tasks: Add `spectr show` Command

## 1. Setup and Dependencies

- [ ] 1.1 Add `github.com/yuin/goldmark` to go.mod
- [ ] 1.2 Add goldmark extensions (tables, strikethrough, task lists, syntax highlighting)
- [ ] 1.3 Add `github.com/pkg/browser` for cross-platform browser opening
- [ ] 1.4 Create `internal/show/` package directory

## 2. Markdown Renderer

- [ ] 2.1 Create `internal/show/renderer.go` with Goldmark setup
- [ ] 2.2 Implement `RenderMarkdown(content []byte) (string, error)` function
- [ ] 2.3 Configure Goldmark extensions (tables, strikethrough, task lists)
- [ ] 2.4 Add syntax highlighting with Chroma extension
- [ ] 2.5 Write unit tests for Markdown rendering with fixtures

## 3. HTML Templates

- [ ] 3.1 Create `internal/show/templates.go` with HTML template
- [ ] 3.2 Implement template with MathJax 3.x CDN integration
- [ ] 3.3 Add clean, readable CSS with typography and spacing
- [ ] 3.4 Implement tab-based layout for changes (proposal/tasks/design/deltas)
- [ ] 3.5 Add JavaScript for tab switching functionality
- [ ] 3.6 Write unit tests for template rendering

## 4. Item Loader

- [ ] 4.1 Create `internal/show/item_loader.go` with ItemLoader struct
- [ ] 4.2 Implement `LoadSpec(specID, projectPath) (*Item, error)` for specs
- [ ] 4.3 Implement `LoadChange(changeID, projectPath) (*Item, error)` for changes
- [ ] 4.4 Handle missing files gracefully with clear error messages
- [ ] 4.5 Combine multiple files for changes (proposal + tasks + design + deltas)
- [ ] 4.6 Write unit tests with testdata fixtures

## 5. HTTP Server

- [ ] 5.1 Create `internal/show/server.go` with Server struct
- [ ] 5.2 Implement `StartServer(html, itemDir, projectRoot) (port int, error)` with auto-port
- [ ] 5.3 Add `GET /` handler serving rendered HTML
- [ ] 5.4 Add `GET /images/<path>` handler with fallback path resolution
- [ ] 5.5 Implement graceful shutdown on Ctrl+C (signal handling)
- [ ] 5.6 Add retry logic for port conflicts (up to 3 attempts)
- [ ] 5.7 Write integration tests for HTTP routes

## 6. Image Path Resolution

- [ ] 6.1 Implement `resolveImagePath(imagePath, itemDir, projectRoot) string` in server.go
- [ ] 6.2 Try spec/change directory first (e.g., `spectr/specs/auth/diagram.png`)
- [ ] 6.3 Fallback to project root (e.g., `./assets/diagram.png`)
- [ ] 6.4 Return 404 if neither path exists
- [ ] 6.5 Write unit tests for path resolution with various scenarios

## 7. Browser Opening

- [ ] 7.1 Create `internal/show/browser.go` wrapper for browser opening
- [ ] 7.2 Implement `OpenBrowser(url string) error` with pkg/browser
- [ ] 7.3 Print URL to terminal before opening
- [ ] 7.4 Handle browser opening failures gracefully with fallback message
- [ ] 7.5 Test on Linux, macOS, and Windows (manual testing)

## 8. CLI Command

- [ ] 8.1 Create `cmd/show.go` with ShowCmd struct
- [ ] 8.2 Add command flags: `--type` (spec|change) for disambiguation
- [ ] 8.3 Implement `Run()` method with item type detection
- [ ] 8.4 Validate item exists before rendering
- [ ] 8.5 Wire up renderer + server + browser opening flow
- [ ] 8.6 Add helpful error messages for common issues (item not found, port conflict)
- [ ] 8.7 Register ShowCmd in `cmd/root.go` CLI struct
- [ ] 8.8 Write integration tests for show command

## 9. MathJax Integration

- [ ] 9.1 Add MathJax 3.x script tag to HTML template
- [ ] 9.2 Configure inline math delimiters (`$...$`)
- [ ] 9.3 Configure display math delimiters (`$$...$$`)
- [ ] 9.4 Test with sample spec containing math notation
- [ ] 9.5 Document CDN requirement and offline behavior

## 10. Testing and Validation

- [ ] 10.1 Create testdata fixtures (sample specs and changes with images and math)
- [ ] 10.2 Write end-to-end test for spec rendering
- [ ] 10.3 Write end-to-end test for change rendering with tabs
- [ ] 10.4 Test image resolution (spec-local and project-root)
- [ ] 10.5 Test MathJax rendering in browser (manual)
- [ ] 10.6 Test graceful degradation when CDN unreachable
- [ ] 10.7 Run linter and fix all issues (golangci-lint)

## 11. Documentation

- [ ] 11.1 Add `spectr show` to CLI help text
- [ ] 11.2 Update README with `spectr show` examples
- [ ] 11.3 Document image path resolution strategy
- [ ] 11.4 Document MathJax syntax and limitations
- [ ] 11.5 Add usage examples to spectr/AGENTS.md if relevant

## 12. Edge Cases and Polish

- [ ] 12.1 Handle missing design.md gracefully (skip tab for changes)
- [ ] 12.2 Handle empty delta specs (show message instead of empty tab)
- [ ] 12.3 Add breadcrumb showing item type and ID in HTML header
- [ ] 12.4 Test with very large specs (>1MB) and add warning if slow
- [ ] 12.5 Test with Unicode filenames and paths
- [ ] 12.6 Ensure proper content-type headers for images (PNG, JPG, SVG)
