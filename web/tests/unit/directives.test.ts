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

	it('rejects an empty selection with a stable error', () => {
		expect(() => wrapSelection('Hello world', '', id)).toThrow('selection is empty');
	});

	it('wraps the first exact occurrence and preserves directive order', () => {
		const markdown = `before {{question:${other}}}existing{{/question}} after after`;
		expect(wrapSelection(markdown, 'after', id)).toBe(
			`before {{question:${other}}}existing{{/question}} {{question:${id}}}after{{/question}} after`
		);
		expect(directiveIds(wrapSelection(markdown, 'after', id))).toEqual([other, id]);
		expect(wrapSelection('Hello world today', 'world', id)).toBe(
			`Hello {{question:${id}}}world{{/question}} today`
		);
		expect(() => wrapSelection(`{{question:${other}}}`, 'ignored', id)).toThrow();
	});

	it('matches collapsed whitespace while preserving the original spaced span', () => {
		const markdown = 'hello   world';
		expect(wrapSelection(markdown, 'hello world', id)).toBe(
			`{{question:${id}}}hello   world{{/question}}`
		);
	});

	it('matches rendered text inside Markdown markers while preserving the markers', () => {
		const markdown = '**bold**';
		expect(wrapSelection(markdown, 'bold', id)).toBe(
			`{{question:${id}}}**bold**{{/question}}`
		);
	});

	it('reports when a selection cannot be found in the note', () => {
		expect(() => wrapSelection('hello world', 'zzz', id)).toThrow(
			'Could not find that passage in the note. Try selecting plain text.'
		);
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
