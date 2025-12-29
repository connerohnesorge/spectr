# Documentation Specification - Delta

## Purpose

Enhances the documentation site with advanced search capabilities by integrating the Star Warp plugin, enabling users to search Spectr docs directly from their browser's address bar and providing quick navigation to search results.

## ADDED Requirements

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
