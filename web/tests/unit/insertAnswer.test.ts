import { describe, expect, it } from 'vitest';
import { insertAnswerAfterDirective } from '$lib/editor/directives';

const id = '550e8400-e29b-41d4-a716-446655440000';
const other = '7f877d74-6a53-442a-9a16-c1c85f029fe8';

describe('insertAnswerAfterDirective', () => {
	it('inserts a trimmed answer as a blockquote after a wrapped directive', () => {
		const markdown = `Before {{question:${id}}}selected passage{{/question}} after`;

		expect(insertAnswerAfterDirective(markdown, id, '  The answer  ')).toBe(
			`Before {{question:${id}}}selected passage{{/question}}\n\n> The answer\n after`
		);
	});

	it('prefixes every line of a multiline answer', () => {
		const markdown = `{{question:${id}}}selected passage{{/question}}`;

		expect(insertAnswerAfterDirective(markdown, id, 'first line\nsecond line')).toBe(
			`{{question:${id}}}selected passage{{/question}}\n\n> first line\n> second line\n`
		);
	});

	it('is idempotent when the answer blockquote is already present', () => {
		const markdown = `{{question:${id}}}selected passage{{/question}}`;
		const inserted = insertAnswerAfterDirective(markdown, id, 'The answer');

		expect(insertAnswerAfterDirective(inserted, id, 'The answer')).toBe(inserted);
	});

	it('leaves markdown unchanged when the directive id is unknown', () => {
		const markdown = `{{question:${id}}}selected passage{{/question}}`;

		expect(insertAnswerAfterDirective(markdown, other, 'The answer')).toBe(markdown);
	});

	it('throws for an empty answer without changing the markdown', () => {
		const markdown = `{{question:${id}}}selected passage{{/question}}`;

		expect(() => insertAnswerAfterDirective(markdown, id, ' \n\t ')).toThrow();
		expect(markdown).toBe(`{{question:${id}}}selected passage{{/question}}`);
	});
});
