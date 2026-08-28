import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import NoteReader from '$lib/editor/NoteReader.svelte';
import type { Question } from '$lib/types/question';

const id = '550e8400-e29b-41d4-a716-446655440000';

const question: Question = {
	id,
	workspaceId: '11111111-1111-4111-8111-111111111111',
	questionText: 'What causes this?',
	answerMarkdown: null,
	status: 'unanswered',
	priority: 'none',
	dueDate: null,
	reminder: null,
	tagIds: [],
	linkedNotes: [],
	createdAt: '2026-08-09T00:00:00.000Z',
	updatedAt: '2026-08-09T00:00:00.000Z',
	version: 1,
	kind: 'question'
};

function mockSelection(text: string) {
	vi.spyOn(window, 'getSelection').mockReturnValue({
		toString: () => text,
		rangeCount: text ? 1 : 0,
		isCollapsed: !text,
		removeAllRanges: () => undefined,
		addRange: () => undefined,
		getRangeAt: () => ({ collapse: () => undefined })
	} as unknown as Selection);
}

describe('NoteReader', () => {
	it('shows selection menu actions after text is selected', async () => {
		render(NoteReader, { markdown: 'The mitochondria is the powerhouse.', questions: [] });
		mockSelection('The mitochondria is the powerhouse.');
		await fireEvent.mouseUp(screen.getByRole('article', { name: 'Reading note' }));
		expect(screen.getByRole('button', { name: 'Ask a question' })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Add annotation' })).toBeInTheDocument();
	});

	it('opens a card when a highlight is clicked', async () => {
		render(NoteReader, {
			markdown: `Lead {{question:${id}}}selected passage{{/question}} tail`,
			questions: [question]
		});
		expect(screen.queryByText('{{question:')).not.toBeInTheDocument();
		await fireEvent.click(screen.getByText('selected passage'));
		expect(screen.getByRole('dialog', { name: /question/i })).toBeInTheDocument();
		expect(screen.getByText('What causes this?')).toBeInTheDocument();
	});
});
