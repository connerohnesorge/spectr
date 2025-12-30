# Documentation Specification Delta

## ADDED Requirements

### Requirement: Page Action Buttons for Documentation

The system SHALL provide page action buttons on documentation pages enabling
users to quickly copy markdown content and open pages in AI chat services
through the starlight-page-actions plugin.

#### Scenario: User copies markdown content

- **WHEN** a user visits any documentation page
- **THEN** they SHALL see a "Copy Markdown" button
- **AND** clicking it SHALL copy the raw markdown content to clipboard

#### Scenario: User opens page in AI chat service

- **WHEN** a user clicks the "Open" dropdown button
- **THEN** they SHALL see options to open the page in default AI chat services
  (ChatGPT, Claude, Gemini, etc.)
- **AND** selecting an option SHALL open the documentation in the chosen service

### Requirement: Starlight Page Actions Plugin Configuration

The system SHALL configure the starlight-page-actions plugin in the Astro
configuration with llms.txt generation disabled to avoid conflicts with existing
starlight-llms-txt plugin while enabling page action functionality with default
AI service list.

#### Scenario: Plugin is installed

- **WHEN** dependencies are installed
- **THEN** the starlight-page-actions package SHALL be present in package.json
- **AND** it SHALL be importable in astro.config.mjs

#### Scenario: Plugin is configured with options

- **WHEN** Starlight is initialized
- **THEN** starlightPageActions() SHALL be included in the plugins array with a
  configuration object
- **AND** the configuration SHALL set `llmstxt: false` to disable llms.txt
  generation
- **AND** the plugin SHALL use default AI services for the Open dropdown

#### Scenario: No conflict with existing llms.txt plugin

- **WHEN** the documentation site is built with both starlight-page-actions and
  starlight-llms-txt plugins
- **THEN** the build SHALL complete without conflicts
- **AND** llms.txt SHALL be generated only by the starlight-llms-txt plugin
- **AND** page action buttons SHALL render correctly without interfering with
  llms.txt generation

#### Scenario: Plugin renders page actions

- **WHEN** a documentation page is loaded
- **THEN** the plugin SHALL render page action buttons
- **AND** the buttons SHALL function correctly with the configured options
