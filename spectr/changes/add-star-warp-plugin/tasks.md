# Tasks: Add Star Warp Plugin to Documentation Site

## Implementation Tasks

- [ ] **Install Star Warp package**: Run `cd docs && npm install @inox-tools/star-warp` to add the dependency
- [ ] **Configure plugin with defaults**: Add `starWarp({ openSearch: { enabled: true } })` to Starlight plugins array in `astro.config.mjs`
- [ ] **Build documentation**: Run `cd docs && npm run build` to verify no build errors
- [ ] **Start dev server**: Run `cd docs && npm run dev` for local testing
- [ ] **Test warp route**: Navigate to `http://localhost:4321/spectr/warp?q=validation` and verify redirect to first result
- [ ] **Verify OpenSearch XML**: Check `http://localhost:4321/spectr/warp.xml` contains correct metadata (title: "Spectr")
- [ ] **Verify head tag injection**: Inspect any docs page source to confirm OpenSearch link tag is automatically injected by plugin
- [ ] **Test browser registration (Chrome)**: Visit docs page, check chrome://settings/searchEngines for Spectr in "Inactive shortcuts", activate it
- [ ] **Test browser registration (Safari)**: Visit docs page, verify search suggestion appears when typing after domain
- [ ] **Test production build**: Run `cd docs && npm run preview` and verify warp route works in production mode

## Verification Checklist

- [ ] `@inox-tools/star-warp` appears in `docs/package.json` dependencies
- [ ] Plugin configured in `docs/astro.config.mjs` plugins array with `openSearch.enabled: true`
- [ ] Documentation builds without errors
- [ ] Warp route `/spectr/warp?q=<term>` redirects to first search result
- [ ] OpenSearch XML at `/spectr/warp.xml` has title "Spectr" and correct base path
- [ ] OpenSearch link tag automatically present in page `<head>` (no manual injection needed)
- [ ] Chrome can detect and activate Spectr as search engine
- [ ] Safari shows search option when typing after domain
- [ ] Production preview build works correctly
