import eslint from '@eslint/js';
import globals from 'globals';
import tseslint from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import svelte from 'eslint-plugin-svelte';

export default [
	eslint.configs.recommended,
	...svelte.configs['flat/recommended'],
	{
		ignores: [
			'node_modules/**',
			'.svelte-kit/**',
			'build/**',
			'dist/**',
			'coverage/**',
			'*.min.js',
			'src/service-worker.ts',
			'*.config.js',
			'*.config.ts'
		]
	},
	{
		files: ['**/*.ts'],
		languageOptions: {
			parser: tsParser,
			parserOptions: { project: './tsconfig.json', projectService: false },
			globals: { ...globals.browser, ...globals.es2022 }
		},
		plugins: { '@typescript-eslint': tseslint },
		rules: {
			'no-unused-vars': 'off',
			'@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }]
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			globals: { ...globals.browser, ...globals.es2022 },
			parserOptions: { parser: tsParser }
		},
		rules: {
			'svelte/no-at-html-tags': 'error',
			'svelte/no-navigation-without-resolve': 'off',
			'svelte/require-each-key': 'off',
			'svelte/no-reactive-functions': 'off',
			'svelte/no-immutable-reactive-statements': 'off',
			'svelte/infinite-reactive-loop': 'off',
			'no-unused-vars': 'off'
		}
	}
];
