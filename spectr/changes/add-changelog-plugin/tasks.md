# Implementation Tasks

## Package Installation
- [x] Navigate to `docs/` directory
- [x] Install `starlight-changelogs` package using Bun: `bun add starlight-changelogs`
- [x] Verify package appears in `docs/package.json` dependencies

## Configuration Updates
- [x] Add import statement for `starlight-changelogs` to `docs/astro.config.mjs`
- [x] Add `starlightChangelogs()` to plugins array in Starlight configuration
- [x] Import `changelogsLoader` from `starlight-changelogs/loader` in `docs/src/content.config.ts`
- [x] Add `changelogs` collection to exports with `changelogsLoader` configuration
- [x] Configure GitHub provider with:
  - `provider: 'github'`
  - `base: 'changelog'`
  - `owner: 'connerohnesorge'`
  - `repo: 'spectr'`
  - `title: 'Changelog'`
  - `pageSize: 10` (default)
  - `pagefind: true` (default, explicit for clarity)

## Verification
- [x] Run `bun run dev` in docs/ directory to start development server
- [x] Navigate to `http://localhost:4321/spectr/changelog/` (accounting for base path)
- [x] Verify changelog overview page loads with paginated version list
- [x] Click on a version entry to verify individual changelog page renders correctly
- [x] Verify changelog pages appear in site search (test search for a release keyword)
- [x] Check that GitHub releases are being fetched and displayed correctly

## Documentation
- [x] Verify proposal.md accurately reflects the implemented changes
- [x] No additional user-facing documentation needed (changelog page is self-explanatory)
