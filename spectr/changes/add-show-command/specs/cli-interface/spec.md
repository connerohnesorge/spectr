# Cli Interface Specification Delta

## ADDED Requirements

### Requirement: Browser-Based Markdown Rendering
The show command SHALL render Markdown files to HTML using Goldmark parser with CommonMark compliance and GitHub Flavored Markdown extensions.

#### Scenario: Render spec Markdown to HTML
- **WHEN** rendering a spec's `spec.md` file
- **THEN** the system SHALL parse Markdown using Goldmark
- **AND** SHALL convert to HTML with proper semantic structure
- **AND** SHALL preserve heading hierarchy (h1, h2, h3, h4)
- **AND** SHALL render lists, links, and blockquotes correctly

#### Scenario: Support GitHub Flavored Markdown extensions
- **WHEN** rendering Markdown content
- **THEN** tables SHALL render as HTML tables with proper borders and styling
- **AND** strikethrough text SHALL render with `<del>` tags
- **AND** task lists SHALL render as checkboxes (checked or unchecked)
- **AND** autolinks SHALL convert URLs to clickable links

#### Scenario: Syntax highlighting for code blocks
- **WHEN** rendering fenced code blocks with language identifiers
- **THEN** the system SHALL apply syntax highlighting using Chroma or similar
- **AND** SHALL support common languages (Go, JavaScript, Python, Bash, Markdown)
- **AND** SHALL render with readable color scheme on light background

### Requirement: MathJax Integration for Mathematical Notation
The show command SHALL render LaTeX-style mathematical notation using MathJax 3.x from CDN.

#### Scenario: Inline math rendering
- **WHEN** Markdown content contains inline math delimited by single dollar signs (`$...$`)
- **THEN** MathJax SHALL render the math inline with surrounding text
- **AND** SHALL use proper font sizing and alignment
- **EXAMPLE**: `The formula $E = mc^2$ shows` renders as "The formula E = mc² shows"

#### Scenario: Display math rendering
- **WHEN** Markdown content contains display math delimited by double dollar signs (`$$...$$`)
- **THEN** MathJax SHALL render the math as a centered block
- **AND** SHALL use larger font size than inline math
- **EXAMPLE**: `$$\int_0^\infty e^{-x^2} dx$$` renders as centered integral

#### Scenario: MathJax CDN loading
- **WHEN** the HTML page is rendered
- **THEN** MathJax SHALL load from `https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js`
- **AND** SHALL configure tex input with inlineMath: `[['$', '$']]` and displayMath: `[['$$', '$$']]`
- **AND** SHALL allow time for CDN to load before user views page

#### Scenario: MathJax unavailable graceful degradation
- **WHEN** MathJax CDN is unreachable (offline environment)
- **THEN** raw LaTeX notation SHALL display as plain text
- **AND** dollar signs SHALL remain visible
- **AND** no rendering errors SHALL block page display

### Requirement: Image Path Resolution with Fallback
The show command's HTTP server SHALL resolve image paths relative to the spec/change directory first, then fall back to project root.

#### Scenario: Serve spec-local images
- **WHEN** Markdown contains image reference `![diagram](./diagram.png)`
- **AND** file exists at `spectr/specs/<spec-id>/diagram.png`
- **THEN** the HTTP server SHALL serve the file from the spec directory
- **AND** SHALL set correct Content-Type header (image/png, image/jpeg, image/svg+xml)

#### Scenario: Serve project-root images as fallback
- **WHEN** Markdown contains image reference `![logo](assets/logo.png)`
- **AND** file does NOT exist in spec/change directory
- **AND** file exists at `./assets/logo.png` (project root)
- **THEN** the HTTP server SHALL serve the file from project root
- **AND** SHALL set correct Content-Type header

#### Scenario: Image not found returns 404
- **WHEN** Markdown contains image reference `![missing](./missing.png)`
- **AND** file does NOT exist in spec/change directory or project root
- **THEN** the HTTP server SHALL return HTTP 404 Not Found
- **AND** browser SHALL display broken image placeholder

#### Scenario: HTTP route for images
- **WHEN** browser requests image at `/images/<path>`
- **THEN** the server SHALL resolve `<path>` using fallback strategy
- **AND** SHALL rewrite Markdown image URLs to use `/images/` prefix during rendering

### Requirement: Tabbed Interface for Changes
The show command SHALL render changes with tabbed sections for Proposal, Tasks, Design, and Deltas.

#### Scenario: Render change with all files
- **WHEN** showing a change that has proposal.md, tasks.md, design.md, and delta specs
- **THEN** the HTML SHALL display 4 tabs: "Proposal", "Tasks", "Design", "Deltas"
- **AND** the Proposal tab SHALL be active by default
- **AND** each tab SHALL display the rendered Markdown content for that file

#### Scenario: Switch between tabs
- **WHEN** user clicks a tab button
- **THEN** the previously active tab content SHALL hide
- **AND** the clicked tab content SHALL display
- **AND** the clicked tab button SHALL gain "active" visual styling
- **AND** tab switching SHALL work via JavaScript without page reload

#### Scenario: Handle missing design.md
- **WHEN** showing a change that lacks design.md
- **THEN** the Design tab SHALL NOT be displayed
- **AND** only Proposal, Tasks, and Deltas tabs SHALL appear
- **AND** tab layout SHALL adjust accordingly

#### Scenario: Combine delta specs in Deltas tab
- **WHEN** rendering the Deltas tab
- **THEN** all delta spec files under `specs/<capability>/spec.md` SHALL be combined
- **AND** each delta SHALL be prefixed with a heading showing the capability name
- **AND** deltas SHALL be sorted alphabetically by capability ID

### Requirement: Clean, Readable HTML Styling
The show command SHALL apply CSS styling for readability and professional appearance.

#### Scenario: Typography and layout
- **WHEN** rendering HTML
- **THEN** body text SHALL use system font stack (system-ui, sans-serif)
- **AND** line height SHALL be at least 1.6 for readability
- **AND** content SHALL be centered with max-width of 900px
- **AND** body SHALL have left/right margin of 40px

#### Scenario: Heading hierarchy styling
- **WHEN** rendering headings
- **THEN** h1, h2, h3, h4 SHALL have increased top margin (1.5em) for visual separation
- **AND** SHALL use bold font weight
- **AND** SHALL maintain relative font sizes (h1 largest, h4 smallest)

#### Scenario: Code block styling
- **WHEN** rendering code blocks
- **THEN** background SHALL be light gray (#f5f5f5)
- **AND** padding SHALL be 1em for spacing
- **AND** horizontal overflow SHALL scroll if content is too wide
- **AND** font SHALL be monospace with readable size

#### Scenario: Table styling
- **WHEN** rendering tables
- **THEN** borders SHALL be visible with 1px solid (#ddd)
- **AND** cells SHALL have padding of 8px
- **AND** header row SHALL have distinct background (#f5f5f5)
- **AND** table SHALL use full width of content area

#### Scenario: Image responsive sizing
- **WHEN** rendering images
- **THEN** max-width SHALL be 100% to prevent overflow
- **AND** height SHALL auto-scale to maintain aspect ratio
- **AND** images SHALL not exceed content width

### Requirement: Auto-Open Browser
The show command SHALL automatically open the rendered page in the user's default browser.

#### Scenario: Successful browser opening
- **WHEN** HTTP server starts successfully
- **THEN** the system SHALL invoke browser.OpenURL() with server URL
- **AND** SHALL use cross-platform browser opening library (pkg/browser)
- **AND** SHALL respect user's default browser preference (system setting)

#### Scenario: Browser opening failure
- **WHEN** browser.OpenURL() fails (e.g., headless environment)
- **THEN** the system SHALL print error message to terminal
- **AND** SHALL display fallback message "Please open manually: http://localhost:<port>"
- **AND** SHALL continue running server (not exit on browser failure)

#### Scenario: Cross-platform browser support
- **WHEN** running on Linux
- **THEN** the system SHALL use xdg-open or equivalent to open browser
- **WHEN** running on macOS
- **THEN** the system SHALL use open command to open browser
- **WHEN** running on Windows
- **THEN** the system SHALL use start command to open browser

### Requirement: Ephemeral HTTP Server
The show command SHALL start a local HTTP server on an auto-assigned port that runs until user terminates it.

#### Scenario: Auto-assign port with :0
- **WHEN** starting HTTP server without `--port` flag
- **THEN** the system SHALL bind to address `:0` for OS-level port assignment
- **AND** SHALL retrieve the assigned port number from the listener
- **AND** SHALL print the assigned port in terminal output

#### Scenario: Server runs until interrupted
- **WHEN** HTTP server is running
- **THEN** it SHALL continue serving requests until Ctrl+C is pressed
- **AND** SHALL handle SIGINT signal for graceful shutdown
- **AND** SHALL close listener and exit cleanly on interrupt

#### Scenario: Serve index route
- **WHEN** browser requests `GET /`
- **THEN** the server SHALL respond with rendered HTML content
- **AND** SHALL set Content-Type header to `text/html; charset=utf-8`
- **AND** SHALL respond with HTTP 200 OK

#### Scenario: Serve images route
- **WHEN** browser requests `GET /images/<path>`
- **THEN** the server SHALL resolve image path with fallback strategy
- **AND** SHALL serve file with appropriate Content-Type
- **AND** SHALL respond with HTTP 200 OK for found files, 404 for missing files

### Requirement: Terminal Output for Show Command
The show command SHALL provide clear terminal output indicating server status and user actions.

#### Scenario: Startup messages
- **WHEN** show command starts successfully
- **THEN** terminal SHALL display "Rendering <item-type>: <item-id>"
- **AND** SHALL display "Starting HTTP server..."
- **AND** SHALL display "Opening browser at http://localhost:<port>"
- **AND** SHALL display "Press Ctrl+C to stop server"

#### Scenario: Error messages
- **WHEN** item does not exist
- **THEN** terminal SHALL display "Error: <item-type> not found: <item-id>"
- **AND** SHALL suggest running `spectr list --specs` or `spectr list` to see available items
- **WHEN** server fails to start
- **THEN** terminal SHALL display "Error: Failed to start server: <error details>"

#### Scenario: Shutdown message
- **WHEN** user presses Ctrl+C to stop server
- **THEN** terminal SHALL display "Shutting down server..."
- **AND** SHALL exit cleanly with status code 0
