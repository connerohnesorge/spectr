# Documentation Specification Changes

## ADDED Requirements

### Requirement: Automated Changelog Integration
The documentation site SHALL integrate with GitHub releases to automatically display version history and release notes through the starlight-changelogs plugin, providing users with an in-site changelog experience without requiring manual maintenance.

#### Scenario: User views changelog overview
- **WHEN** a user navigates to `/changelog/` on the documentation site
- **THEN** they SHALL see a paginated list of all versions with release dates
- **AND** each version entry SHALL link to its detailed changelog page

#### Scenario: User views specific version changelog
- **WHEN** a user navigates to `/changelog/versions/<version>/` or clicks a version from the overview
- **THEN** they SHALL see the full changelog content for that version as published in GitHub releases
- **AND** the content SHALL be rendered in the Starlight theme matching the rest of the documentation

#### Scenario: User searches changelog content
- **WHEN** a user enters a search term that appears in a release note
- **THEN** the relevant changelog page SHALL appear in search results via Pagefind
- **AND** clicking the search result SHALL navigate to the changelog page

#### Scenario: New release is published
- **WHEN** a new GitHub release is published in the connerohnesorge/spectr repository
- **THEN** the changelog SHALL automatically include it on the next site build
- **AND** no manual documentation update SHALL be required

### Requirement: GitHub Releases Provider Configuration
The documentation site SHALL use the GitHub provider for the starlight-changelogs plugin to fetch release data from the connerohnesorge/spectr repository, with proper configuration for pagination, search indexing, and URL structure.

#### Scenario: Plugin fetches releases from correct repository
- **WHEN** the documentation site is built
- **THEN** the plugin SHALL fetch releases from owner `connerohnesorge` and repo `spectr`
- **AND** it SHALL use the GitHub API without requiring authentication for public repository access

#### Scenario: Changelog uses consistent URL structure
- **WHEN** changelog pages are generated
- **THEN** the overview page SHALL be available at `/changelog/`
- **AND** individual version pages SHALL follow the pattern `/changelog/versions/<version>/`
- **AND** URLs SHALL account for the site base path `/spectr/`

#### Scenario: Changelog pages are searchable
- **WHEN** Pagefind indexes the documentation site
- **THEN** changelog pages SHALL be included in the search index
- **AND** users SHALL be able to search for release-specific content

#### Scenario: Changelog displays reasonable pagination
- **WHEN** the changelog overview page renders
- **THEN** it SHALL display 10 versions per page by default
- **AND** provide pagination controls when there are more than 10 versions

### Requirement: Documentation Site Plugin Integration
The Astro documentation site SHALL configure the starlight-changelogs plugin in both the Astro config and content collections, following Starlight's plugin architecture and Astro's content collections pattern.

#### Scenario: Plugin is registered in Starlight
- **WHEN** the Astro configuration is loaded
- **THEN** `starlightChangelogs` SHALL be imported from the `starlight-changelogs` package
- **AND** it SHALL be added to the Starlight plugins array alongside existing plugins

#### Scenario: Changelogs content collection is defined
- **WHEN** Astro loads content collections
- **THEN** a `changelogs` collection SHALL be defined in `src/content.config.ts`
- **AND** it SHALL use the `changelogsLoader` with GitHub provider configuration

#### Scenario: Plugin dependency is tracked
- **WHEN** the docs project dependencies are reviewed
- **THEN** `starlight-changelogs` SHALL appear in the dependencies section of `docs/package.json`
- **AND** the version SHALL be the latest compatible release
