# Proposal: Add Star Warp Plugin to Documentation Site

## Summary

Integrate the @inox-tools/star-warp Starlight plugin into the Spectr
documentation site to enable "jump to first search result" functionality and
OpenSearch browser integration.

## Motivation

Star Warp enhances the documentation site's search experience by:

- Enabling users to search the Spectr docs directly from their browser's address
  bar
- Providing quick navigation with `/warp?q=search term` URLs that jump to the
  first result
- Supporting OpenSearch protocol for native browser integration
- Working completely statically without requiring SSR

This improves documentation accessibility and aligns with modern documentation
UX patterns used by Astro and other developer tools.

## Scope

- Install `@inox-tools/star-warp` package in the docs directory
- Configure Star Warp as a Starlight plugin in `astro.config.mjs` with default
  settings
- Enable OpenSearch Description generation (plugin handles head tag injection
  automatically)
- Verify `/warp` route functionality through local testing
- Test browser integration (Chrome/Safari) with manual registration

## Implementation Decisions

- **Head tag handling**: Star Warp plugin automatically injects OpenSearch link
  tag
- **Configuration**: Use all defaults (path: `/warp`, title: `Spectr`,
  description: `Search Spectr`)
- **Testing**: Local browser testing including /warp route and OpenSearch
  registration
- **Documentation**: No explicit user documentation; feature discoverable via
  browser

## Affected Files

- `docs/package.json` - Add Star Warp dependency
- `docs/astro.config.mjs` - Configure plugin with OpenSearch enabled

## Success Criteria

1. Star Warp plugin is installed and configured
2. Navigating to `/spectr/warp?q=<term>` redirects to first search result
3. OpenSearch Description is generated at `/spectr/warp.xml`
4. Browser can register the Spectr docs as a searchable site
5. Documentation build succeeds without errors
