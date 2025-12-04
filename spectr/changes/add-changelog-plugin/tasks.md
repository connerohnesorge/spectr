# Implementation Tasks

## Package Installation
- [ ] Navigate to `docs/` directory
- [ ] Install `starlight-changelogs` package using Bun: `bun add starlight-changelogs`
- [ ] Verify package appears in `docs/package.json` dependencies

## Configuration Updates
- [ ] Add import statement for `starlight-changelogs` to `docs/astro.config.mjs`
- [ ] Add `starlightChangelogs()` to plugins array in Starlight configuration
- [ ] Import `changelogsLoader` from `starlight-changelogs/loader` in `docs/src/content.config.ts`
- [ ] Add `changelogs` collection to exports with `changelogsLoader` configuration
- [ ] Configure GitHub provider with:
  - `provider: 'github'`
  - `base: 'changelog'`
  - `owner: 'connerohnesorge'`
  - `repo: 'spectr'`
  - `title: 'Changelog'`
  - `pageSize: 10` (default)
  - `pagefind: true` (default, explicit for clarity)

## Verification
- [ ] Run `bun run dev` in docs/ directory to start development server
- [ ] Navigate to `http://localhost:4321/spectr/changelog/` (accounting for base path)
- [ ] Verify changelog overview page loads with paginated version list
- [ ] Click on a version entry to verify individual changelog page renders correctly
- [ ] Verify changelog pages appear in site search (test search for a release keyword)
- [ ] Check that GitHub releases are being fetched and displayed correctly

## Documentation
- [ ] Verify proposal.md accurately reflects the implemented changes
- [ ] No additional user-facing documentation needed (changelog page is self-explanatory)
