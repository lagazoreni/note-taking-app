import { describe, expect, it } from 'vitest';
import { rejectRawHtml, renderMarkdown, renderNoteHtml } from '$lib/editor/markdown';
import { directiveIds } from '$lib/editor/directives';

const id = '550e8400-e29b-41d4-a716-446655440000';

describe('markdown security', () => {
	it('rejects raw HTML and keeps directives structured', () => {
		expect(() => rejectRawHtml('<img src=x onerror=alert(1)>')).toThrow();
		expect(directiveIds(`{{question:${id}}}`)).toEqual([id]);
		expect(renderMarkdown('**safe**')).not.toContain('<script>');
	});

	it('never renders a raw question directive token', () => {
		const wrapped = renderNoteHtml(`Hello {{question:${id}}}world{{/question}}`);
		const bare = renderNoteHtml(`See {{question:${id}}} later`);
		expect(wrapped).not.toContain('{{question:');
		expect(bare).not.toContain('{{question:');
		expect(wrapped).toContain('data-annotation-id');
	});
});
