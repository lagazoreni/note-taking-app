import { describe, expect, it } from 'vitest';
import { renderNoteHtml } from '$lib/editor/markdown';

const id = '550e8400-e29b-41d4-a716-446655440000';

describe('renderNoteHtml status metadata', () => {
	it('adds status and kind attributes to highlight marks without exposing directives', () => {
		const html = renderNoteHtml(
			`Hello {{question:${id}}}selected passage{{/question}}`,
			[{ id, status: 'answered', kind: 'question' }]
		);
		const mark = html.match(/<mark[^>]*>/)?.[0] ?? '';

		expect(mark).toContain(`data-annotation-id="${id}"`);
		expect(mark).toContain('data-status="answered"');
		expect(mark).toContain('data-kind="question"');
		expect(html).not.toContain('{{question:');
	});
});
