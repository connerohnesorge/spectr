# Proposal: Add Starlight Page Actions Plugin to Documentation Site

## Summary

Integrate the starlight-page-actions Starlight plugin into the Spectr
documentation site to enable page action buttons including "Copy Markdown"
functionality and AI chat service integration via "Open" dropdown.

## Motivation

Starlight Page Actions enhances the documentation site's usability and
AI-friendliness by:

- Enabling users to quickly copy raw markdown content for AI chat services
- Providing one-click access to open documentation in various AI chat interfaces
- Improving the documentation workflow for users who leverage AI tools

This aligns with modern documentation patterns that recognize AI assistants as
first-class documentation consumers and supports developer workflows that
integrate AI tools.

## Scope

- Install `starlight-page-actions` package in the docs directory
- Configure Starlight Page Actions as a Starlight plugin in `astro.config.mjs`
- Disable llms.txt generation in starlight-page-actions (keeping existing
  starlight-llms-txt plugin)
- Use default AI service list in "Open" dropdown provided by the plugin
- Verify page action buttons render correctly on documentation pages
- Test "Copy Markdown" and "Open" dropdown functionality with basic functional
  testing

## Affected Files

- `docs/package.json` - Add starlight-page-actions dependency
- `docs/astro.config.mjs` - Configure plugin in Starlight plugins array

## Implementation Details

### Configuration

The plugin will be configured in `docs/astro.config.mjs` as follows:

```js
import starlightPageActions from 'starlight-page-actions';

export default defineConfig({
  integrations: [
    starlight({
      plugins: [
        starlightSiteGraph(),
        starlightLlmsTxt(),
        starlightPageActions({
          llmstxt: false  // Disable to avoid conflict with starlight-llms-txt
        })
      ],
      // ... rest of config
    })
  ]
});
```

### Key Configuration Decisions

1. **llms.txt Generation**: Disabled in starlight-page-actions (`llmstxt:
  false`) to prevent conflicts with existing starlight-llms-txt plugin
2. **AI Services**: Using plugin defaults (no custom service list configuration)
3. **Styling**: Using plugin defaults (no custom styling configuration)

## Success Criteria

1. Starlight Page Actions plugin is installed and configured with llms.txt
  generation disabled
2. Page action buttons (Copy Markdown, Open dropdown) appear on documentation
  pages
3. "Copy Markdown" button successfully copies raw markdown content
4. "Open" dropdown provides links to default AI chat services (ChatGPT, Claude,
  Gemini, etc.)
5. Existing starlight-llms-txt plugin continues to work without conflicts
6. Documentation build and preview work without errors

## Out of Scope

- Custom styling of page action buttons (using plugin defaults)
- Removal of existing starlight-llms-txt plugin (keeping both for their
  respective features)
- Deployment testing on GitHub Pages (implementation uses local testing only)
