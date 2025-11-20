// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightSiteGraph from 'starlight-site-graph';
import starlightLlmsTxt from 'starlight-llms-txt';

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
			plugins: [starlightSiteGraph(), starlightLlmsTxt()],
		}),
	],
});
