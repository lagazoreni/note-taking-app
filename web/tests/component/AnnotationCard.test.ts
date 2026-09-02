import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import AnnotationCard from '$lib/components/AnnotationCard.svelte';
import type { Question } from '$lib/types/question';

function makeItem(
	kind: 'question' | 'annotation',
	text: string,
	overrides: Partial<Question> = {}
): Question {
	return {
		id: '550e8400-e29b-41d4-a716-446655440000',
		workspaceId: '11111111-1111-4111-8111-111111111111',
		questionText: text,
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
		kind,
		...overrides
	};
}

describe('AnnotationCard', () => {
	it('shows question lifecycle chrome for questions', () => {
		render(AnnotationCard, {
			question: makeItem('question', 'What causes this?'),
			passage: 'selected passage'
		});
		expect(screen.getByRole('dialog', { name: /question/i })).toBeInTheDocument();
		expect(screen.getByText('selected passage')).toBeInTheDocument();
		expect(screen.getByText('What causes this?')).toBeInTheDocument();
		expect(screen.getByLabelText('Answer')).toBeInTheDocument();
		expect(screen.getByLabelText('Status')).toBeInTheDocument();
	});

	it('shows an insert button for answered questions and invokes it once', async () => {
		const onInsertAnswer = vi.fn();
		render(AnnotationCard, {
			question: makeItem('question', 'What causes this?', {
				status: 'answered',
				answerMarkdown: 'The answer'
			}),
			passage: 'selected passage',
			onInsertAnswer
		});

		const button = screen.getByRole('button', { name: 'Insert answer into note' });
		expect(button).toBeInTheDocument();
		await fireEvent.click(button);
		expect(onInsertAnswer).toHaveBeenCalledTimes(1);
	});

	it('does not show an insert button for unanswered questions or annotations', () => {
		render(AnnotationCard, {
			question: makeItem('question', 'What causes this?'),
			passage: 'selected passage'
		});
		expect(screen.queryByRole('button', { name: 'Insert answer into note' })).not.toBeInTheDocument();

		render(AnnotationCard, {
			question: makeItem('annotation', 'Remember this later', {
				status: 'answered',
				answerMarkdown: 'An annotation answer'
			}),
			passage: 'selected passage'
		});
		expect(screen.queryByRole('button', { name: 'Insert answer into note' })).not.toBeInTheDocument();
	});
});
