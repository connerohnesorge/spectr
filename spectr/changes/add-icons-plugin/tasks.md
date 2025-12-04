# Implementation Tasks

## 1. Prerequisites & Dependencies
- [ ] 1.1 Install UnoCSS Astro integration: `cd docs && npm install -D unocss`
- [ ] 1.2 Install starlight-plugin-icons: `npm install -D starlight-plugin-icons`
- [ ] 1.3 Install an icon collection (Phosphor recommended): `npm install -D @iconify-json/ph`
- [ ] 1.4 Verify all packages appear in docs/package.json devDependencies

## 2. UnoCSS Configuration
- [ ] 2.1 Create `docs/uno.config.ts` configuration file
- [ ] 2.2 Import and configure `presetStarlightIcons()` preset in uno.config.ts
- [ ] 2.3 Configure UnoCSS content pipeline to include astro.config.mjs for icon class extraction
- [ ] 2.4 Add `.starlight-icons` cache directory to `docs/.gitignore`

## 3. Astro Configuration Refactor
- [ ] 3.1 Import `UnoCSS` from 'unocss/astro' in docs/astro.config.mjs
- [ ] 3.2 Import `Icons` from 'starlight-plugin-icons' in docs/astro.config.mjs
- [ ] 3.3 Add `UnoCSS()` to the integrations array before Starlight
- [ ] 3.4 Wrap existing `starlight()` configuration with `Icons()` function
- [ ] 3.5 Move all Starlight config (title, social, sidebar, plugins) into Icons({ starlight: { ... } })
- [ ] 3.6 Verify existing plugins (starlightSiteGraph, starlightLlmsTxt) remain in plugins array

## 4. Feature Enablement
- [ ] 4.1 Enable sidebar icons by adding `sidebar: true` to Icons configuration
- [ ] 4.2 Enable safelist extraction with `extractSafelist: true` for automatic icon discovery
- [ ] 4.3 Verify configuration follows wrapper pattern with starlight config nested under Icons()

## 5. Build & Development Testing
- [ ] 5.1 Run `npm run dev` in docs/ directory to start development server
- [ ] 5.2 Verify documentation site builds without errors or warnings
- [ ] 5.3 Check browser console for any UnoCSS or icon-related errors
- [ ] 5.4 Verify existing documentation pages render correctly with no regressions

## 6. Feature Verification
- [ ] 6.1 Verify UnoCSS integration is working (check Network tab for UnoCSS styles)
- [ ] 6.2 Verify icon classes can be used (add test icon to a sidebar entry as proof-of-concept)
- [ ] 6.3 Test that code blocks display with automatic language icons
- [ ] 6.4 Verify FileTree component renders with file type icons
- [ ] 6.5 Check that enhanced components (Card, Aside, IconLink) are available

## 7. Production Build Test
- [ ] 7.1 Run `npm run build` in docs/ directory
- [ ] 7.2 Verify build completes without errors
- [ ] 7.3 Run `npm run preview` to test production build
- [ ] 7.4 Verify icons render correctly in production build

## 8. Documentation & Validation
- [ ] 8.1 Document any configuration decisions or customizations made
- [ ] 8.2 Run `spectr validate add-icons-plugin --strict` to ensure proposal validity
- [ ] 8.3 Resolve any validation errors
- [ ] 8.4 Update tasks.md to mark all completed tasks
