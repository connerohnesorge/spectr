import { defineConfig } from 'unocss';
import { presetStarlightIcons } from 'starlight-plugin-icons/uno';

export default defineConfig({
	presets: [presetStarlightIcons()],
	content: {
		pipeline: {
			include: [
				// Include astro.config.mjs to extract icon classes from sidebar configuration
				/astro\.config\.mjs$/,
			],
		},
	},
});
