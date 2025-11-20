# Change: Add `spectr show` Browser Renderer with Goldmark and MathJax

## Why

Spectr specifications and change proposals are written in Markdown, which often includes complex formatting, diagrams (images), and mathematical notation. While Markdown is human-readable in plain text, rendering it in a browser provides:
- **Better readability**: Proper typography, syntax highlighting, and layout
- **Rich media support**: Inline images and diagrams render visually
- **Mathematical notation**: MathJax renders LaTeX-style math (`$...$` and `$$...$$`) for technical specifications
- **Improved collaboration**: Stakeholders can review specs in a polished format without learning Markdown

Currently, users must manually convert Markdown to HTML or use external tools. A `spectr show` command that auto-renders specs and changes in the browser streamlines the review process and improves the developer/stakeholder experience.

## What Changes

- **NEW**: Add `spectr show <item>` command that:
  - Accepts a spec ID (e.g., `spectr show auth --type spec`) or change ID (e.g., `spectr show add-2fa --type change`)
  - Renders Markdown content to HTML using **Goldmark** (Go's standard Markdown parser)
  - Starts a local HTTP server (e.g., `http://localhost:8080`) with auto-assigned port
  - Serves rendered HTML with:
    - **MathJax** integration for math rendering (`$inline$` and `$$display$$`)
    - **Image support** with path resolution (spec/change directory → project root fallback)
    - Clean, readable CSS styling with proper typography
    - Syntax highlighting for code blocks
  - Automatically opens the rendered page in the user's default browser
  - Keeps server running until user terminates (Ctrl+C)
  - Displays server URL in terminal for manual access
- **For specs**: Renders `spectr/specs/<id>/spec.md`
- **For changes**: Renders a combined view of `proposal.md`, `tasks.md`, `design.md`, and all delta specs
- **BREAKING**: None

## Impact

- **Affected specs**: `cli-interface`, `cli-framework`
- **Affected code**:
  - `cmd/` - Add `ShowCmd` command handler (wait, there's already a view command, so we need to check if "show" is taken)
  - `internal/` - New `show` package with:
    - Markdown rendering with Goldmark
    - HTTP server with static asset serving
    - MathJax template integration
    - Image path resolution
  - `go.mod` - Add `github.com/yuin/goldmark` dependency
- **User-visible changes**: One new command
- **Dependencies**:
  - `goldmark` for Markdown parsing
  - Go's `net/http` for local server (stdlib)
  - MathJax served from CDN (no installation required)
