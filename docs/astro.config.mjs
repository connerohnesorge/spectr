// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightSiteGraph from 'starlight-site-graph';
import starlightLlmsTxt from 'starlight-llms-txt';
import starlightChangelogs from 'starlight-changelogs';
import starlightPageActions from 'starlight-page-actions';

// https://astro.build/config
export default defineConfig({
	site: 'https://connerohnesorge.github.io',
	base: 'spectr',
	integrations: [
		starlight({
			title: 'Spectr',
			social: [
				{
					label: 'GitHub',
					href: 'https://github.com/connerohnesorge/spectr',
					icon: 'github',
				},
			],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Quick Start', slug: 'getting-started/quick-start' },
					],
				},
				{
					label: 'Core Concepts',
					items: [
						{ label: 'Spec-Driven Development', slug: 'concepts/spec-driven-development' },
						{ label: 'Delta Specifications', slug: 'concepts/delta-specifications' },
						{ label: 'Validation Rules', slug: 'concepts/validation-rules' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Creating Changes', slug: 'guides/creating-changes' },
						{ label: 'Archiving Workflow', slug: 'guides/archiving-workflow' },
					],
				},
				{
					label: 'Reference',
					items: [
						{ label: 'CLI Commands', slug: 'reference/cli-commands' },
						{ label: 'Configuration', slug: 'reference/configuration' },
					],
				},
			],
			plugins: [
				starlightSiteGraph(),
				starlightLlmsTxt(),
				starlightChangelogs(),
				starlightPageActions({
					llmstxt: false, // Disable to avoid conflict with starlight-llms-txt
				}),
			],
		}),
	],
});
