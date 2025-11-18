// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'Spectr',
			social: {
				github: 'https://github.com/conneroisu/spectr',
			},
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
		}),
	],
});
