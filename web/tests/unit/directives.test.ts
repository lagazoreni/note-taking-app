import { describe, expect, it } from 'vitest';
import {
	directiveIds,
	stripDirectives,
	tokenizeDirectives,
	wrapSelection
} from '$lib/editor/directives';

const id = '550e8400-e29b-41d4-a716-446655440000';
const other = '7f877d74-6a53-442a-9a16-c1c85f029fe8';

describe('wrapped highlight directives', () => {
	it('tokenizes wrapped highlights and keeps the passage as a snippet', () => {
		const markdown = `Lead {{question:${id}}}selected passage{{/question}} tail`;
		const tokens = tokenizeDirectives(markdown);
		const question = tokens.find((token) => token.kind === 'question');
		expect(question?.directive?.id).toBe(id);
		expect(question?.directive?.snippet).toBe('selected passage');
		expect(question?.directive?.wrapped).toBe(true);
		expect(tokens.some((token) => token.value.includes('{{/question}}'))).toBe(false);
		expect(directiveIds(markdown)).toEqual([id]);
	});

	it('tokenizes bare legacy tokens without treating them as visible prose', () => {
		const markdown = `See {{question:${id}}} later`;
		const tokens = tokenizeDirectives(markdown);
		const question = tokens.find((token) => token.kind === 'question');
		expect(question?.directive?.id).toBe(id);
		expect(question?.directive?.snippet ?? '').toBe('');
		expect(question?.directive?.wrapped).toBe(false);
		expect(directiveIds(markdown)).toEqual([id]);
	});

	it('wraps the selected passage with a highlight directive', () => {
		expect(wrapSelection('Hello world today', 'world', id)).toBe(
			`Hello {{question:${id}}}world{{/question}} today`
		);
		expect(wrapSelection(`Keep {{question:${other}}}world{{/question}}`, 'Keep', id)).toBe(
			`{{question:${id}}}Keep{{/question}} {{question:${other}}}world{{/question}}`
		);
		expect(() => wrapSelection(`{{question:${other}}}`, 'ignored', id)).toThrow();
	});

	it('strips wrap tokens while keeping the snippet, and drops bare tokens', () => {
		expect(stripDirectives(`Hello {{question:${id}}}world{{/question}} today`)).toBe(
			'Hello world today'
		);
		expect(stripDirectives(`See {{question:${id}}} later`).replace(/\s+/g, ' ').trim()).toBe(
			'See later'
		);
	});
});
