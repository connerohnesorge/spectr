# Design Document: Starlight Icons Plugin Integration

## Context

The Spectr documentation site currently uses Astro Starlight as its
documentation framework with two existing plugins: starlight-site-graph and
starlight-llms-txt. The site lacks visual iconography for navigation and content
organization. The starlight-plugin-icons plugin can add this functionality, but
it requires a significant architectural addition: UnoCSS as a build-time CSS
framework.

### Current State

- **Framework**: Astro 5.16.0 + Starlight 0.36.3
- **Build Tools**: Bun for package management
- **Existing Plugins**: starlight-site-graph, starlight-llms-txt
- **CSS Approach**: Starlight's built-in styling (no additional CSS framework)

### Constraints

- Must maintain compatibility with existing Starlight plugins
- Must not break current documentation pages or styling
- Must support build-time icon processing for optimal performance
- Must be compatible with GitHub Pages deployment (static site only)

## Goals / Non-Goals

### Goals

- Enable icon support in sidebar navigation with minimal configuration
- Provide automatic icon display for code blocks based on language
- Access to 200,000+ Iconify icons via UnoCSS integration
- Maintain existing plugin functionality (site graph, llms.txt)
- Keep configuration simple and maintainable

### Non-Goals

- Replacing all existing styling with UnoCSS (only using for icons)
- Adding icons to every sidebar entry immediately (icons are optional)
- Creating custom icon sets or collections
- Migrating away from Starlight's default theme

## Decisions

### Decision 1: Use UnoCSS for Icon Processing

**Rationale**: The starlight-plugin-icons plugin requires UnoCSS to render icons
from Iconify at build time. This is a hard dependency, not optional.

**Alternatives Considered**:

1. **Use Starlight's built-in icon support**: Limited to specific icon sets and
  requires manual SVG imports. Does not support Iconify's vast library.
2. **Use a different icon plugin**: No other Starlight plugins provide
  comparable icon coverage with sidebar, code block, and component integration.
3. **Manual icon implementation**: Would require significant custom development
  and maintenance.

**Trade-offs**:

- ✅ **Pro**: Access to 200,000+ icons from Iconify with minimal effort
- ✅ **Pro**: Build-time processing ensures optimal runtime performance
- ✅ **Pro**: Official Starlight plugin with good documentation
- ❌ **Con**: Adds UnoCSS as a new build dependency (~450KB of dev dependencies)
- ❌ **Con**: Introduces another CSS processing layer (though scoped to icons
  only)
- ❌ **Con**: Slight increase in build complexity

**Decision**: Accept UnoCSS as a dependency. The benefits of comprehensive icon
support outweigh the cost of an additional build tool, especially since it's
scoped to icon processing only.

### Decision 2: Use Wrapper Pattern for Starlight Integration

**Rationale**: The plugin uses a wrapper pattern where `Icons()` wraps the
entire `starlight()` configuration, rather than being added to the `plugins`
array. This is the plugin's required integration method.

**Configuration Pattern**:

```javascript
// Standard Starlight plugin pattern (what we're NOT using)
starlight({
  title: 'Spectr',
  plugins: [somePlugin()]
})

// Icons wrapper pattern (what we MUST use)
Icons({
  starlight: {
    title: 'Spectr',
    plugins: [somePlugin()]
  }
})
```text

**Why This Pattern**:

- The plugin needs to intercept and enhance sidebar configuration with icon
  support
- Wrapper pattern allows the plugin to process the entire Starlight config
  before initialization
- Required for automatic icon extraction from sidebar entries

**Trade-offs**:

- ✅ **Pro**: Enables automatic icon extraction from config
- ✅ **Pro**: Seamless sidebar icon integration
- ✅ **Pro**: Maintains compatibility with existing plugins in the plugins array
- ❌ **Con**: Non-standard integration pattern compared to other Starlight
  plugins
- ❌ **Con**: Requires refactoring existing config structure

**Decision**: Use the wrapper pattern as required by the plugin. The integration
benefits justify the configuration restructuring.

### Decision 3: Use Phosphor Icons as Initial Collection

**Rationale**: Icon collections are installed separately as `@iconify-json/*`
packages. We need to choose an initial collection for the implementation.

**Options Evaluated**:

1. **Material Design Icons** (`@iconify-json/mdi`): ~7,000 icons, comprehensive
  but large
2. **Phosphor Icons** (`@iconify-json/ph`): ~6,000 icons, modern design system,
  duotone variants
3. **Heroicons** (`@iconify-json/heroicons`): ~300 icons, minimal but limited
4. **Lucide** (`@iconify-json/lucide`): ~1,000 icons, popular in React ecosystem

**Decision**: Use Phosphor Icons (`@iconify-json/ph`) as the initial collection.

**Reasoning**:

- Modern, well-designed icon system with consistent visual language
- Includes duotone variants for visual richness (e.g.,
  `i-ph:rocket-launch-duotone`)
- Large enough library (~6,000 icons) for most documentation needs
- Used in the plugin's own documentation, so well-supported
- Additional collections can be installed on-demand later

**Trade-offs**:

- ✅ **Pro**: High-quality, modern design system
- ✅ **Pro**: Duotone variants add visual interest
- ✅ **Pro**: Good balance of coverage vs package size
- ❌ **Con**: Not as comprehensive as Material Design Icons
- ⚪ **Neutral**: Other collections can be added later without migration

### Decision 4: Enable Automatic Safelist Extraction

**Rationale**: The plugin offers two approaches for UnoCSS to discover icon
classes:

1. Manual UnoCSS content configuration to scan `astro.config.mjs`
2. Automatic safelist extraction with `extractSafelist: true`

**Decision**: Use automatic safelist extraction (`extractSafelist: true`).

**Reasoning**:

- Simpler configuration (one flag vs complex content pipeline config)
- Plugin automatically generates safelist of icon classes used in sidebar
- Less prone to configuration errors
- Recommended approach in plugin documentation
- Generates `.starlight-icons` cache directory for faster rebuilds

**Trade-offs**:

- ✅ **Pro**: Simpler, more maintainable configuration
- ✅ **Pro**: Automatic icon discovery without manual content config
- ✅ **Pro**: Faster rebuilds with cache
- ❌ **Con**: Creates cache directory that needs .gitignore entry
- ❌ **Con**: Slightly less explicit than manual configuration

### Decision 5: Minimal Initial Configuration

**Rationale**: The plugin supports extensive customization (custom icon scaling,
transformations, styling), but we'll start with minimal configuration.

**Initial Configuration**:

```javascript
Icons({
  sidebar: true,
  extractSafelist: true,
  starlight: {
    // existing starlight config
  }
})
```text

**Not Configuring Initially**:

- Custom icon sizes or scaling
- Custom transformations per icon
- Custom CSS variables for icon styling
- Additional icon collections beyond Phosphor

**Reasoning**:

- Start simple, add complexity only when needed
- Plugin defaults are well-designed for most use cases
- Easier to understand and maintain minimal config
- Can add customizations incrementally based on actual needs

**Trade-offs**:

- ✅ **Pro**: Simple, understandable configuration
- ✅ **Pro**: Easy to maintain and debug
- ✅ **Pro**: Follows "simplicity first" project convention
- ❌ **Con**: May need configuration updates if defaults don't meet needs
- ⚪ **Neutral**: Customization can be added later without migration

## Architecture Impact

### Build Process Changes

**Before**:

```text
Astro build → Starlight processing → Static site output
```text

**After**:

```text
Astro build → UnoCSS processing → Starlight processing (wrapped by Icons)
→ Static site output
```

### Dependency Graph

```text
docs/
├── astro (framework)
├── @astrojs/starlight (docs theme)
├── unocss (CSS engine for icons) ← NEW
│   └── @iconify-json/ph (icon collection) ← NEW
├── starlight-plugin-icons (icon integration) ← NEW
├── starlight-site-graph (existing plugin)
└── starlight-llms-txt (existing plugin)
```text

### File Structure Changes

**New Files**:

- `docs/uno.config.ts` - UnoCSS configuration with presetStarlightIcons()
- `docs/.starlight-icons` - Cache directory (git-ignored)

**Modified Files**:

- `docs/package.json` - Add 3 new devDependencies
- `docs/astro.config.mjs` - Restructure with wrapper pattern
- `docs/.gitignore` - Add .starlight-icons

**No Changes**:

- `docs/src/**` - Content files unchanged
- `docs/public/**` - Static assets unchanged
- Existing documentation pages - No content migration needed

## Risks / Trade-offs

### Risk 1: UnoCSS Conflicts with Starlight Styling

**Risk**: UnoCSS and Starlight's built-in styling might conflict or create
specificity issues.

**Likelihood**: Low - UnoCSS is scoped to icon classes only

**Mitigation**:

- Plugin uses `presetStarlightIcons()` which is designed for Starlight
  compatibility
- UnoCSS only processes icon classes (i-* pattern), not general utility classes
- Thorough testing during implementation to catch any style conflicts

**Impact if occurs**: Visual styling issues in documentation site

### Risk 2: Build Time Increase

**Risk**: Adding UnoCSS and icon processing could slow down builds.

**Likelihood**: Medium - Additional build step adds some overhead

**Mitigation**:

- UnoCSS is known for fast build times (on-demand CSS generation)
- Icon safelist cache (`.starlight-icons`) speeds up rebuilds
- Build-time processing is better than runtime icon loading for user experience

**Measurement**: Track build times before/after during implementation

**Impact if occurs**: Slower development feedback loop, longer CI/CD times

### Risk 3: Plugin Compatibility Issues

**Risk**: The wrapper pattern might cause issues with existing or future
Starlight plugins.

**Likelihood**: Low - Existing plugins remain in plugins array

**Mitigation**:

- Plugin documentation shows compatibility with standard Starlight plugins
- Existing plugins (site-graph, llms-txt) remain in their current location
- Test all existing functionality after integration

**Rollback**: Can revert to standard Starlight integration if compatibility
issues arise

**Impact if occurs**: Loss of existing plugin functionality, need to choose
between plugins

### Risk 4: Icon Cache Directory Commits

**Risk**: `.starlight-icons` cache directory could be accidentally committed to
git.

**Likelihood**: Medium without proper .gitignore

**Mitigation**:

- Add `.starlight-icons` to `.gitignore` as part of implementation
- Document this requirement in tasks.md

**Impact if occurs**: Repository bloat, unnecessary cache files in version
control

## Migration Plan

This is a new feature addition, not a migration, but there are integration
steps:

### Phase 1: Dependency Installation

1. Install UnoCSS, starlight-plugin-icons, and initial icon collection
2. Verify packages in package.json
3. No user-facing changes yet

### Phase 2: Configuration Setup

1. Create uno.config.ts with preset
2. Add .gitignore entry for cache
3. Still no user-facing changes

### Phase 3: Integration Refactor

1. Restructure astro.config.mjs with wrapper pattern
2. Enable sidebar icons and safelist extraction
3. Test build and preview

### Phase 4: Verification

1. Verify existing pages render correctly (no regressions)
2. Verify existing plugins work (site-graph, llms.txt)
3. Test that icons are available for use
4. Verify production build works

### Rollback Procedure

If issues arise during implementation:

1. **Revert package.json** to remove new dependencies
2. **Revert astro.config.mjs** to original Starlight configuration
3. **Delete uno.config.ts**
4. **Run `npm install`** to clean up node_modules
5. **Test that documentation site builds successfully**

Time to rollback: ~5 minutes

## Open Questions

### Q1: Should we add icons to sidebar entries immediately?

**Status**: Deferred

**Context**: Icons in sidebar are optional per entry. We could add them during
initial implementation or incrementally.

**Options**:

- Add icons to all sidebar entries immediately for consistent visual language
- Leave sidebar without icons initially, add them incrementally based on benefit
- Add icons only to top-level sections, not individual pages

**Decision Needed**: During or after implementation

**Recommendation**: Start without sidebar icons, add them incrementally. This:

- Reduces initial implementation scope
- Allows time to choose appropriate icons thoughtfully
- Demonstrates the feature works without requiring full migration
- Can be done in follow-up PRs as enhancement

### Q2: Do we need additional icon collections?

**Status**: Deferred

**Context**: We're starting with Phosphor icons. We might need Material Design
Icons for code block language icons.

**Research Needed**: Verify whether Material Design Icons are needed for
automatic code block icons, or if the plugin includes them internally.

**Decision Point**: During testing phase when verifying code block icons

**If needed**: Install `@iconify-json/mdi` as additional collection

### Q3: Should we customize icon sizes or styling?

**Status**: Deferred

**Context**: Plugin supports custom CSS variables for icon sizing and spacing.

**Default values**:

```css
:root {
  --spi-sidebar-icon-size: 1.25rem;
  --spi-sidebar-icon-gap: 0.25rem;
}
```text

**Decision**: Use defaults initially, customize only if user feedback indicates
issues

**Review point**: After initial deployment to staging/production

## References

- [starlight-plugin-icons
  Documentation](https://docs.rettend.me/starlight-plugin-icons)
- [starlight-plugin-icons
  GitHub](https://github.com/Rettend/starlight-plugin-icons)
- [UnoCSS Astro Integration](https://unocss.dev/integrations/astro)
- [Iconify Icon Collections](https://icones.js.org/)
- [Phosphor Icons](https://phosphoricons.com/)
- [Starlight Configuration
  Reference](https://starlight.astro.build/reference/configuration/)
