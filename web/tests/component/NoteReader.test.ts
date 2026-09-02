import { fireEvent, render, screen, within } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
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
	beforeEach(() => {
		vi.restoreAllMocks();
	});

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

	it('opens the matching highlight card when focusQuestionId is set', () => {
		render(NoteReader, {
			markdown: `Lead {{question:${id}}}selected passage{{/question}} tail`,
			questions: [question],
			focusQuestionId: id
		});

		const dialog = screen.getByRole('dialog', { name: /question/i });
		expect(dialog).toBeInTheDocument();
		expect(within(dialog).getByText('selected passage')).toBeInTheDocument();
		expect(within(dialog).getByText('What causes this?')).toBeInTheDocument();
	});

	it('shows capture errors and keeps the composer open', async () => {
		const errorMessage = 'Could not find that passage in the note. Try selecting plain text.';
		const onCapture = vi.fn().mockRejectedValue(new Error(errorMessage));
		render(NoteReader, { markdown: 'Select this passage.', questions: [], onCapture });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Ask a question' }));
		const input = screen.getByLabelText('Question text');
		await fireEvent.input(input, { target: { value: 'Why does this happen?' } });
		await fireEvent.click(screen.getByRole('button', { name: 'Save question' }));

		expect(await screen.findByRole('alert')).toHaveTextContent(errorMessage);
		expect(screen.getByLabelText('Question text')).toBeInTheDocument();
		expect(screen.queryByRole('dialog', { name: /question/i })).not.toBeInTheDocument();
		expect(onCapture).toHaveBeenCalledOnce();
	});

	it('dismisses selection actions and the composer without capturing', async () => {
		const onCapture = vi.fn();
		render(NoteReader, { markdown: 'Select this passage.', questions: [], onCapture });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
		expect(screen.queryByRole('toolbar', { name: 'Selection actions' })).not.toBeInTheDocument();

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Ask a question' }));
		expect(screen.getByLabelText('Question text')).toBeInTheDocument();
		await fireEvent.keyDown(window, { key: 'Escape' });
		expect(screen.queryByLabelText('Question text')).not.toBeInTheDocument();
		expect(onCapture).not.toHaveBeenCalled();

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Add annotation' }));
		await fireEvent.click(document.body);
		expect(screen.queryByLabelText('Annotation')).not.toBeInTheDocument();
		expect(onCapture).not.toHaveBeenCalled();
	});
});
