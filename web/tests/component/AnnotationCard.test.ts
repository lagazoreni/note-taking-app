import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import AnnotationCard from '$lib/components/AnnotationCard.svelte';
import type { Question } from '$lib/types/question';

function makeItem(kind: 'question' | 'annotation', text: string): Question {
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
		kind
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

	it('shows annotation chrome without question lifecycle controls', () => {
		render(AnnotationCard, {
			question: makeItem('annotation', 'Remember this later'),
			passage: 'selected passage'
		});
		expect(screen.getByRole('dialog', { name: /annotation/i })).toBeInTheDocument();
		expect(screen.getByText('Remember this later')).toBeInTheDocument();
		expect(screen.queryByLabelText('Status')).not.toBeInTheDocument();
		expect(screen.getByLabelText('Annotation')).toBeInTheDocument();
	});
});
