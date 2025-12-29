# Documentation Specification Delta

## ADDED Requirements

### Requirement: Starlight Icons Plugin Integration with UnoCSS

The documentation site SHALL integrate the starlight-plugin-icons plugin with UnoCSS to enable icon support in sidebar navigation, code blocks, and file tree components using the Iconify icon library through a wrapper-based integration pattern.

#### Scenario: Icons plugin dependencies are installed

- **WHEN** the documentation site dependencies are installed
- **THEN** the starlight-plugin-icons package SHALL be present in docs/package.json
- **AND** the unocss package SHALL be present as a prerequisite dependency
- **AND** at least one @iconify-json/* icon collection package SHALL be installed
- **AND** the plugin SHALL be configured in docs/astro.config.mjs using the wrapper pattern

#### Scenario: UnoCSS is configured for icon processing

- **WHEN** the documentation site is built
- **THEN** a uno.config.ts file SHALL exist with presetStarlightIcons() configured
- **AND** UnoCSS SHALL be added to the Astro integrations array before Starlight
- **AND** UnoCSS SHALL extract icon classes from astro.config.mjs for safelist generation

#### Scenario: Starlight integration uses Icons wrapper pattern

- **WHEN** the astro.config.mjs is read
- **THEN** the Starlight configuration SHALL be wrapped by the Icons() function
- **AND** all Starlight options (title, social, sidebar, plugins) SHALL be nested under Icons({ starlight: { ... } })
- **AND** existing Starlight plugins SHALL remain functional within the wrapped configuration

#### Scenario: Sidebar icons feature is enabled

- **WHEN** sidebar icons are configured
- **THEN** the Icons configuration SHALL include sidebar: true
- **AND** sidebar entries MAY include an icon field with format i-<collection>:<name>
- **AND** icon classes SHALL be automatically extracted via extractSafelist: true

#### Scenario: Icons are available for documentation content

- **WHEN** documentation pages are viewed
- **THEN** icons from installed Iconify collections SHALL be available for use
- **AND** icons MAY be used in sidebar navigation entries via icon field
- **AND** code blocks SHALL automatically display Material Icons based on language or title
- **AND** FileTree components SHALL render with file type icons

#### Scenario: Documentation builds with icons plugin and UnoCSS

- **WHEN** the documentation site is built
- **THEN** the build SHALL complete successfully with UnoCSS and icons plugin enabled
- **AND** icons SHALL render correctly on all documentation pages
- **AND** no conflicts SHALL occur with existing Starlight plugins (starlightSiteGraph, starlightLlmsTxt)

#### Scenario: Icon cache directory is ignored by version control

- **WHEN** the .gitignore file is checked
- **THEN** the .starlight-icons cache directory SHALL be listed in docs/.gitignore
- **AND** icon cache files SHALL not be committed to the repository

### Requirement: Enhanced Components with Icon Support

The documentation site SHALL provide enhanced Starlight components (Card, Aside, FileTree, IconLink) with icon support through the starlight-plugin-icons integration.

#### Scenario: FileTree displays file type icons

- **WHEN** a FileTree component is rendered in documentation
- **THEN** file entries SHALL display appropriate icons based on file extensions
- **AND** folder entries SHALL display folder icons
- **AND** icons SHALL use Material Icon set by default

#### Scenario: Code blocks display language icons

- **WHEN** a code block is rendered with a language identifier
- **THEN** the code block SHALL display an appropriate icon for the language
- **AND** the icon SHALL be from the Material Icon set
- **AND** icons SHALL appear automatically without manual configuration

#### Scenario: Enhanced components are available for use

- **WHEN** documentation content is authored
- **THEN** Card components SHALL support icon attributes
- **AND** Aside components SHALL support icon attributes
- **AND** IconLink components SHALL be available for icon-labeled links
- **AND** all enhanced components SHALL use icons from installed Iconify collections
