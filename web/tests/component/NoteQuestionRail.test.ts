import { fireEvent, render, screen, within } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import NoteQuestionRail from '$lib/components/NoteQuestionRail.svelte';
import type { Question } from '$lib/types/question';

function makeQuestion(
	id: string,
	questionText: string,
	status: Question['status'],
	kind: Question['kind'] = 'question'
): Question {
	return {
		id,
		workspaceId: '11111111-1111-4111-8111-111111111111',
		questionText,
		answerMarkdown: status === 'answered' ? 'Already answered' : null,
		status,
		priority: 'none',
		dueDate: null,
		reminder: null,
		tagIds: [],
		linkedNotes: [],
		createdAt: '2026-08-09T00:00:00.000Z',
		updatedAt: '2026-08-09T00:00:00.000Z',
		version: 1,
		kind
	};
}

describe('NoteQuestionRail', () => {
	it('lists only open questions in directive order and handles selection and next', async () => {
		const onSelect = vi.fn();
		const onNext = vi.fn();
		const questions = [
			makeQuestion('annotation', 'A note about this passage', 'unanswered', 'annotation'),
			makeQuestion('answered', 'This question is done', 'answered'),
			makeQuestion('deferred', 'The deferred question', 'deferred'),
			makeQuestion('in-progress', 'The in-progress question', 'in_progress'),
			makeQuestion('unanswered', 'The unanswered question', 'unanswered')
		];

		render(NoteQuestionRail, {
			questions,
			orderedIds: ['unanswered', 'deferred', 'in-progress'],
			onSelect,
			onNext
		});

		const rail = screen.getByRole('complementary', { name: 'Open questions in this note' });
		const nextButton = within(rail).getByRole('button', { name: 'Next unanswered' });
		const itemButtons = within(rail)
			.getAllByRole('button')
			.filter((button) => button !== nextButton);

		expect(itemButtons).toHaveLength(3);
		expect(itemButtons[0]).toHaveTextContent('The unanswered question');
		expect(itemButtons[0]).toHaveTextContent('unanswered');
		expect(itemButtons[1]).toHaveTextContent('The deferred question');
		expect(itemButtons[1]).toHaveTextContent('deferred');
		expect(itemButtons[2]).toHaveTextContent('The in-progress question');
		expect(itemButtons[2]).toHaveTextContent('in_progress');
		expect(screen.queryByText('A note about this passage')).not.toBeInTheDocument();
		expect(screen.queryByText('This question is done')).not.toBeInTheDocument();

		await fireEvent.click(itemButtons[1]);
		expect(onSelect).toHaveBeenCalledWith('deferred');

		await fireEvent.click(nextButton);
		expect(onNext).toHaveBeenCalledOnce();
	});

	it('shows the empty state when no open questions remain', () => {
		render(NoteQuestionRail, {
			questions: [
				makeQuestion('answered', 'This question is done', 'answered'),
				makeQuestion('annotation', 'A note about this passage', 'unanswered', 'annotation')
			],
			orderedIds: ['answered', 'annotation']
		});

		const rail = screen.getByRole('complementary', { name: 'Open questions in this note' });
		expect(within(rail).getByText('No open questions in this note')).toBeInTheDocument();

		const nextButton = within(rail).queryByRole('button', { name: 'Next unanswered' });
		if (nextButton) {
			expect(nextButton).toBeDisabled();
		}
	});
});
