# Documentation Specification

## Purpose

Comprehensive documentation enables users and developers to understand and use Spectr effectively. Clear guides, command references, and examples reduce onboarding friction and support self-service learning for all user personas.

## Requirements

### Requirement: Comprehensive README with Multiple Sections
The system SHALL provide a comprehensive README.md file that serves both end users and developers, including installation instructions, usage guide, command reference, architecture overview, and contribution guidelines.

#### Scenario: User finds installation instructions
- **WHEN** a new user visits the repository
- **THEN** they SHALL find clear instructions for installing via Nix, building from source, or using pre-built binaries

#### Scenario: Developer understands architecture
- **WHEN** a developer reads the README
- **THEN** they SHALL find an architecture overview explaining the clean separation of concerns and package structure

#### Scenario: Contributor knows how to contribute
- **WHEN** someone wants to contribute
- **THEN** they SHALL find guidelines for code style, testing, commit conventions, and PR process

### Requirement: Quick Start Workflow Guide
The system SHALL provide a quick-start guide demonstrating the core workflow: creating a change, validating it, implementing it, and archiving it.

#### Scenario: User follows workflow example
- **WHEN** a user reads the quick start section
- **THEN** they SHALL see a concrete example of `spectr init`, `spectr list`, `spectr validate`, and `spectr archive` commands in sequence

#### Scenario: User understands file structure
- **WHEN** a user completes the quick start
- **THEN** they SHALL understand the distinction between `specs/`, `changes/`, and `archive/` directories

### Requirement: Complete Command Reference
The system SHALL document all CLI commands with flags, examples, and expected output. Documentation SHALL only reference commands that actually exist in the CLI. Documentation SHALL distinguish between instructions for human users and AI agents.

#### Scenario: User learns init command usage
- **WHEN** a user reads the init command documentation
- **THEN** they SHALL see all available flags (`--path`, `--tools`, `--non-interactive`) with explanations and examples

#### Scenario: User learns list command options
- **WHEN** a user reads the list command documentation
- **THEN** they SHALL understand the `--specs`, `--json`, and `--long` flags with example outputs

#### Scenario: User learns validate command options
- **WHEN** a user reads the validate command documentation
- **THEN** they SHALL see how to use `--strict` flag and understand what validation rules are enforced

#### Scenario: User learns archive command
- **WHEN** a user reads the archive command documentation
- **THEN** they SHALL understand the archiving workflow and `--skip-specs` flag usage

#### Scenario: User learns view command options
- **WHEN** a user reads the view command documentation
- **THEN** they SHALL understand the `--json` flag for dashboard output

#### Scenario: Documentation accuracy
- **WHEN** a user or AI assistant reads any documentation file
- **THEN** all referenced CLI commands SHALL exist and work as documented
- **AND** no nonexistent commands (such as `spectr show`) SHALL be referenced

#### Scenario: AI agent documentation uses direct file reading
- **WHEN** an AI agent reads documentation for viewing specs or changes
- **THEN** the documentation SHALL instruct agents to read files directly (e.g., `spectr/specs/<capability>/spec.md` or `spectr/changes/<change-id>/proposal.md`)
- **AND** the documentation SHALL NOT instruct agents to use CLI commands like `spectr view` for reading content
- **AND** the `spectr view` command SHALL only be documented for human users

### Requirement: Development Setup Guide
The system SHALL provide clear instructions for setting up a development environment and running tests.

#### Scenario: Developer sets up environment with Nix
- **WHEN** a developer reads the development setup section
- **THEN** they SHALL see instructions to run `nix develop` and what tools are available

#### Scenario: Developer runs tests
- **WHEN** a developer reads the testing section
- **THEN** they SHALL know how to run `go test ./...` and understand test organization

### Requirement: Spec-Driven Development Explanation
The system SHALL explain the three-stage workflow and key concepts for users unfamiliar with spec-driven development.

#### Scenario: User understands change proposals
- **WHEN** a user reads about spec-driven development
- **THEN** they SHALL understand that changes are proposals separate from current specs

#### Scenario: User understands requirements and scenarios
- **WHEN** a user reads about key concepts
- **THEN** they SHALL know what requirements, scenarios, and delta specs mean

### Requirement: Troubleshooting and FAQ Section
The system SHALL provide solutions for common issues and answer frequently asked questions.

#### Scenario: User encounters validation error
- **WHEN** a user reads the troubleshooting section
- **THEN** they SHALL find explanations of common validation errors and how to fix them

#### Scenario: User has question about workflow
- **WHEN** a user reads the FAQ
- **THEN** they SHALL find answers to questions like "Do I need approval before implementing?" or "How do I handle merge conflicts?"

### Requirement: Visual CLI Demonstrations
The system SHALL provide visual demonstrations of core CLI workflows using VHS-generated GIF recordings to help users quickly understand Spectr's capabilities.

#### Scenario: User sees initialization demo
- **WHEN** a user reads the Quick Start section of the README
- **THEN** they SHALL see a GIF demonstrating the `spectr init` command and resulting directory structure

#### Scenario: User sees validation demo
- **WHEN** a user reads about validation in the documentation
- **THEN** they SHALL see a GIF showing validation errors and how to fix them

#### Scenario: User sees complete workflow demo
- **WHEN** a user visits the getting-started guide
- **THEN** they SHALL see a GIF demonstrating the complete workflow from proposal to archive

### Requirement: Reproducible Demo Source Files
The system SHALL maintain VHS tape files as version-controlled source for all demo GIFs to enable easy regeneration when the CLI changes.

#### Scenario: Developer regenerates outdated GIF
- **WHEN** a developer updates a CLI command
- **THEN** they SHALL be able to run the corresponding VHS tape file to regenerate an accurate GIF

#### Scenario: Developer creates new demo
- **WHEN** a developer wants to add a new demo
- **THEN** they SHALL find existing tape files as examples in `assets/vhs/` directory

#### Scenario: Contributor finds demo standards
- **WHEN** a contributor reads the development documentation
- **THEN** they SHALL find guidelines for VHS tape configuration (theme, size, typing speed)

### Requirement: Demo Asset Organization
The system SHALL organize demo assets with clear separation between source files (VHS tapes) and generated outputs (GIFs).

#### Scenario: Developer locates tape source
- **WHEN** a developer needs to update a demo
- **THEN** they SHALL find VHS tape files in `assets/vhs/` directory

#### Scenario: Documentation references generated GIF
- **WHEN** the README or docs site needs to embed a demo
- **THEN** they SHALL reference GIF files from `assets/gifs/` directory

#### Scenario: Developer regenerates all demos
- **WHEN** a developer runs the regeneration command
- **THEN** all GIFs SHALL be generated from their corresponding tape files and placed in `assets/gifs/`

### Requirement: Core Workflow Coverage
The system SHALL provide demo GIFs covering all essential Spectr workflows: initialization, listing, validation, and archiving.

#### Scenario: User learns initialization
- **WHEN** a user views the init demo GIF
- **THEN** they SHALL see `spectr init` being run and the resulting `spectr/` directory structure

#### Scenario: User learns listing
- **WHEN** a user views the list demo GIF
- **THEN** they SHALL see both `spectr list` (changes) and `spectr list --specs` (specifications) commands

#### Scenario: User learns validation
- **WHEN** a user views the validate demo GIF
- **THEN** they SHALL see `spectr validate` catching an error, the error being fixed, and validation passing

#### Scenario: User learns archiving
- **WHEN** a user views the archive demo GIF
- **THEN** they SHALL see `spectr archive` merging deltas into specs and moving the change to the archive directory

#### Scenario: User sees end-to-end workflow
- **WHEN** a user views the workflow demo GIF
- **THEN** they SHALL see the complete three-stage workflow from creating a change through archiving it

### Requirement: CI Integration Documentation
The system SHALL provide documentation explaining how to integrate Spectr validation into CI/CD pipelines using the spectr-action GitHub Action.

#### Scenario: User finds spectr-action repository
- **WHEN** a user reads the README
- **THEN** they SHALL find a link to the connerohnesorge/spectr-action repository in the Links & Resources section

#### Scenario: User adds CI validation to their project
- **WHEN** a user reads the CI Integration section
- **THEN** they SHALL see a complete example of adding the spectr-action to a GitHub Actions workflow
- **AND** the example SHALL include the action reference, checkout step, and proper configuration

#### Scenario: User understands CI validation benefits
- **WHEN** a user reads the CI Integration section
- **THEN** they SHALL understand that the action provides automated validation on push and pull request events
- **AND** they SHALL know that it fails fast to provide rapid feedback on specification violations

### Requirement: Pre-made Example Projects for VHS Demos

The system SHALL provide pre-made spectr project examples in the `examples/` directory that VHS tape files use for demonstrations by executing commands directly within the example directories, keeping demos focused on spectr commands rather than file creation boilerplate.

#### Scenario: Developer creates clean demo
- **WHEN** a VHS tape file needs a spectr project for demonstration
- **THEN** it SHALL execute spectr commands directly in the `examples/<demo-name>/` directory
- **AND** the demo output SHALL focus on spectr commands, not temporary directory copying

#### Scenario: Demo runs in example directory
- **WHEN** a VHS tape demonstrates a spectr command
- **THEN** the tape SHALL use `Hide`/`Show` to conceal the directory change
- **AND** spectr commands SHALL run directly in the example directory without copying to `_demo`

#### Scenario: Developer maintains example project
- **WHEN** a change to demo content is needed
- **THEN** the developer SHALL edit the pre-made example in `examples/` directory
- **AND** the change will automatically apply to any tape using that example

#### Scenario: Demo cleanup is minimal
- **WHEN** a tape completes execution
- **THEN** no `rm -rf _demo` cleanup commands SHALL be needed
- **AND** the example directory remains unchanged

### Requirement: Batch GIF Generation Command
The system SHALL provide a `generate-gif` command (via Nix flake) to generate all demo GIFs in one command, supporting both full regeneration and single-demo regeneration.

#### Scenario: Developer regenerates all GIFs
- **WHEN** a developer runs `generate-gif` in the nix develop shell
- **THEN** all VHS tape files SHALL be processed
- **AND** GIFs SHALL be output to `assets/gifs/` directory

#### Scenario: Developer regenerates single GIF
- **WHEN** a developer runs `generate-gif <demo-name>`
- **THEN** only the specified demo's GIF SHALL be regenerated

#### Scenario: Developer gets command usage help
- **WHEN** a developer runs `generate-gif --help`
- **THEN** they SHALL see available demo names and usage instructions

### Requirement: VHS Tape Output Clarity

VHS tape files SHALL NOT contain typed echo statements that display section headers or commentary. Demos SHALL let the spectr commands and their output speak for themselves. Comments within the tape file (lines starting with `#`) SHOULD be used to document sections for maintainers, but these are not displayed in the recording.

#### Scenario: No typed echo section headers
- **WHEN** a VHS tape file is reviewed
- **THEN** it SHALL contain no `Type "echo ..."` commands for section headers
- **AND** visual context SHALL be provided through VHS comments (starting with `#`) which are not recorded

#### Scenario: No useless echo statements
- **WHEN** a VHS tape file is reviewed
- **THEN** it SHALL contain no `Type "echo ''"` commands
- **AND** visual spacing SHALL be achieved through `Sleep` commands only

#### Scenario: Commands are self-documenting
- **WHEN** a user views a demo GIF
- **THEN** they SHALL see only the actual spectr commands being typed
- **AND** they SHALL NOT see preparatory echo statements being typed out

### Requirement: Page Action Buttons for Documentation
The system SHALL provide page action buttons on documentation pages enabling users to quickly copy markdown content and open pages in AI chat services through the starlight-page-actions plugin.

#### Scenario: User copies markdown content
- **WHEN** a user visits any documentation page
- **THEN** they SHALL see a "Copy Markdown" button
- **AND** clicking it SHALL copy the raw markdown content to clipboard

#### Scenario: User opens page in AI chat service
- **WHEN** a user clicks the "Open" dropdown button
- **THEN** they SHALL see options to open the page in default AI chat services (ChatGPT, Claude, Gemini, etc.)
- **AND** selecting an option SHALL open the documentation in the chosen service

### Requirement: Starlight Page Actions Plugin Configuration
The system SHALL configure the starlight-page-actions plugin in the Astro configuration with llms.txt generation disabled to avoid conflicts with existing starlight-llms-txt plugin while enabling page action functionality with default AI service list.

#### Scenario: Plugin is installed
- **WHEN** dependencies are installed
- **THEN** the starlight-page-actions package SHALL be present in package.json
- **AND** it SHALL be importable in astro.config.mjs

#### Scenario: Plugin is configured with options
- **WHEN** Starlight is initialized
- **THEN** starlightPageActions() SHALL be included in the plugins array with a configuration object
- **AND** the configuration SHALL set `llmstxt: false` to disable llms.txt generation
- **AND** the plugin SHALL use default AI services for the Open dropdown

#### Scenario: No conflict with existing llms.txt plugin
- **WHEN** the documentation site is built with both starlight-page-actions and starlight-llms-txt plugins
- **THEN** the build SHALL complete without conflicts
- **AND** llms.txt SHALL be generated only by the starlight-llms-txt plugin
- **AND** page action buttons SHALL render correctly without interfering with llms.txt generation

#### Scenario: Plugin renders page actions
- **WHEN** a documentation page is loaded
- **THEN** the plugin SHALL render page action buttons
- **AND** the buttons SHALL function correctly with the configured options
### Requirement: Browser-Integrated Search with Star Warp
The documentation site SHALL integrate the @inox-tools/star-warp Starlight plugin to enable quick navigation to search results and native browser search integration via OpenSearch protocol.

#### Scenario: User searches via warp URL
- **WHEN** a user navigates to `/spectr/warp?q=validation`
- **THEN** they SHALL be redirected to the first search result matching "validation"
- **AND** the redirect SHALL work statically without requiring server-side rendering

#### Scenario: User searches from browser address bar
- **WHEN** a user has registered the Spectr docs as a search engine in their browser
- **AND** they type the site shortcut followed by a search query in the address bar
- **THEN** the browser SHALL submit the query to the Spectr docs warp endpoint
- **AND** they SHALL be navigated to the first matching result

#### Scenario: OpenSearch description is available
- **WHEN** a user visits any documentation page
- **THEN** the page SHALL include an OpenSearch link tag in the `<head>` section
- **AND** the link tag SHALL be automatically injected by the Star Warp plugin
- **AND** the OpenSearch XML SHALL be available at `/spectr/warp.xml`
- **AND** the description SHALL identify the site as "Spectr" for browser display

#### Scenario: Browser registers search engine
- **WHEN** a user visits the Spectr documentation site
- **THEN** their browser (Chrome, Safari, Firefox) SHALL automatically detect the OpenSearch description
- **AND** the browser SHALL allow the user to activate "Spectr" as a custom search engine
- **AND** typing queries after the domain SHALL trigger Spectr documentation search

#### Scenario: Star Warp configuration with defaults
- **WHEN** the Star Warp plugin is configured in `astro.config.mjs`
- **THEN** it SHALL be added to the Starlight plugins array (not the root integrations array)
- **AND** OpenSearch SHALL be enabled with `openSearch.enabled: true`
- **AND** the path SHALL use the default `/warp` value
- **AND** the OpenSearch title SHALL default to "Spectr" from the Starlight config
- **AND** the OpenSearch description SHALL default to "Search Spectr"
- **AND** the warp route SHALL respect the project's base path `/spectr`
