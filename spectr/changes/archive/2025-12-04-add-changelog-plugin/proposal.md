# Change: Add Changelog Plugin to Documentation Site

## Why

The Spectr documentation site currently lacks an automated changelog display that shows version history and release notes to users. GitHub releases serve as the authoritative source of version information, but users must leave the documentation site to view them. By adding the starlight-changelogs plugin configured to pull from GitHub releases, we provide users with an integrated, searchable changelog experience directly within the documentation site, improving discoverability of new features and changes.

## What Changes

- Install `starlight-changelogs` npm package in docs/ project
- Configure plugin in `docs/astro.config.mjs` to integrate with Starlight
- Add `changelogs` collection to `docs/src/content.config.ts` using GitHub provider
- Configure GitHub provider to fetch releases from `connerohnesorge/spectr` repository
- Set changelog base path to `/changelog/` for consistent URL structure
- Enable pagefind indexing for changelog pages to make them searchable

## Impact

- **Affected specs**: `documentation` (MODIFIED - adds changelog integration requirements)
- **Affected files**:
  - Modified `docs/package.json` - Add starlight-changelogs dependency
  - Modified `docs/astro.config.mjs` - Import and configure starlightChangelogs plugin
  - Modified `docs/src/content.config.ts` - Add changelogs collection with GitHub provider
- **User-visible changes**:
  - New `/changelog/` page displaying paginated list of all versions
  - Individual changelog pages at `/changelog/versions/<version>/`
  - Changelog content sourced from GitHub releases
  - Changelog pages indexed in site search

## Benefits

- **Integrated experience**: Users can view changelog without leaving documentation site
- **Automated updates**: Changelog automatically pulls from GitHub releases, no manual maintenance
- **Searchable history**: Changelog pages indexed by Pagefind for easy search
- **Version navigation**: Paginated interface makes it easy to browse version history
- **Consistent design**: Changelog pages match Starlight's design system and site navigation
