# Change: Add Starlight Icons Plugin to Documentation Site

## Why

The Spectr documentation site currently lacks visual iconography in the sidebar
navigation, code blocks, and file tree components. The starlight-plugin-icons
plugin provides access to over 200,000 icons from the Iconify library via UnoCSS
integration, enabling enhanced visual navigation and improved content
organization through icon-based visual cues. This plugin extends Starlight's
functionality by wrapping the Starlight integration and adding build-time icon
processing capabilities.

## What Changes

- Install `starlight-plugin-icons` package and UnoCSS dependencies in docs/
  project
- Install UnoCSS Astro integration (`unocss`) as a prerequisite
- Create `uno.config.ts` configuration file with `presetStarlightIcons()`
- Modify `docs/astro.config.mjs` to wrap Starlight integration with Icons()
  function
- Install initial icon collection package (e.g., `@iconify-json/ph` for Phosphor
  icons)
- Add `.starlight-icons` cache directory to `.gitignore`
- Enable sidebar icons feature through plugin configuration
- Configure UnoCSS to extract icon classes from `astro.config.mjs`
- Verify icons render correctly on documentation pages
- Test icon functionality with basic functional testing

## Impact

- **Affected specs**: `documentation` (MODIFIED - adds icon integration
  requirements)
- **Affected files**:
  - Modified `docs/package.json` - Add starlight-plugin-icons, unocss, and icon
    collection dependencies
  - Created `docs/uno.config.ts` - UnoCSS configuration with Starlight Icons
    preset
  - Modified `docs/astro.config.mjs` - Replace Starlight integration with Icons
    wrapper
  - Modified `docs/.gitignore` - Add `.starlight-icons` cache directory
- **Architecture changes**:
  - Introduces UnoCSS as a new build-time dependency
  - Changes Starlight integration pattern from direct plugin to wrapper function
  - Adds icon class extraction and safelist generation to build process
- **User-visible changes**:
  - Enhanced sidebar navigation with icon support (icons added per-entry via
    `icon` field)
  - Automatic code block icons matching Material Icon sets by language or title
  - Enhanced FileTree, Card, Aside, and IconLink components with icon support
  - Access to 200,000+ Iconify icons for documentation content

## Benefits

- **Visual navigation**: Icons provide quick visual cues for navigation
  structure
- **Improved scannability**: Code blocks and file trees become easier to scan
  with automatic language icons
- **Consistency**: Unified icon system across all documentation components via
  Iconify
- **Extensibility**: Large Iconify library enables future icon enhancements by
  installing additional icon collections
- **Build-time processing**: Icons are processed at build time for optimal
  performance
- **Flexible configuration**: Supports custom icon scaling, transformations, and
  collections per project needs

## Technical Context

### Integration Pattern

Unlike typical Starlight plugins that use the `plugins` array,
starlight-plugin-icons uses a **wrapper pattern** where the `Icons()` function
wraps the entire Starlight configuration. This allows the plugin to intercept
and enhance the sidebar configuration with icon support.

**Before:**

```javascript
starlight({
  title: 'Spectr',
  plugins: [/* other plugins */]
})
```

**After:**

```javascript
Icons({
  starlight: {
    title: 'Spectr',
    plugins: [/* other plugins */]
  }
})
```

### UnoCSS Requirement

The plugin requires UnoCSS to be installed and configured as it uses UnoCSS's
icon preset to render icons from Iconify at build time. The
`presetStarlightIcons()` preset provides optimized configuration for Starlight's
specific needs.

### Icon Collection Management

Icons are provided through installable `@iconify-json/*` packages. Each icon
collection (e.g., Material Design, Phosphor, Heroicons) is a separate npm
package that can be installed on-demand. The plugin uses the format
`i-<collection>:<name>` to reference icons (e.g., `i-ph:rocket-launch-duotone`).

### Prerequisites

- Astro 4+ (currently using 5.16.0 ✓)
- Starlight 0.35+ (currently using 0.36.3 ✓)
- UnoCSS Astro integration (to be installed)

## Out of Scope

- Custom icon styling beyond default configuration (using plugin defaults)
- Adding icons to all sidebar entries immediately (icons are optional per entry)
- Creating custom icon collections (using Iconify's existing collections)
- Replacing existing Starlight plugins (Icons wrapper is compatible with
  existing plugins)

## References

- [starlight-plugin-icons
  Documentation](https://docs.rettend.me/starlight-plugin-icons)
- [starlight-plugin-icons GitHub
  Repository](https://github.com/Rettend/starlight-plugin-icons)
- [Iconify Icon Gallery](https://icones.js.org/)
- [UnoCSS Astro Integration](https://unocss.dev/integrations/astro)
- [Starlight Plugins &
  Integrations](https://starlight.astro.build/resources/plugins/)
