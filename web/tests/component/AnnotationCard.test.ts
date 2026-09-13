import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import AnnotationCard from '$lib/components/AnnotationCard.svelte';
import type { Question } from '$lib/types/question';

type CardPosition = {
	top: number;
	left: number;
};

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

	it('anchors the card beside supplied highlight coordinates', () => {
		const cardPosition: CardPosition = { top: 208, left: 96 };
		const props = {
			question: makeItem('question', 'What causes this?'),
			passage: 'selected passage',
			cardPosition
		};
		render(AnnotationCard, props);

		const card = screen.getByRole('dialog', { name: /question/i });
		expect(['absolute', 'fixed']).toContain(card.style.position);
		expect(card.style.top).toBe(`${cardPosition.top}px`);
		expect(card.style.left).toBe(`${cardPosition.left}px`);
	});

	it('keeps an edge-anchored card as a fixed, viewport-safe overlay on mobile', () => {
		const originalWidth = window.innerWidth;
		const originalHeight = window.innerHeight;
		Object.defineProperty(window, 'innerWidth', { configurable: true, value: 375 });
		Object.defineProperty(window, 'innerHeight', { configurable: true, value: 667 });

		try {
			const props = {
				question: makeItem('question', 'What causes this?'),
				passage: 'selected passage',
				cardPosition: { top: 650, left: 360 }
			};
			render(AnnotationCard, props);

			const card = screen.getByRole('dialog', { name: /question/i });
			expect(getComputedStyle(card).position).toBe('fixed');
			const top = Number.parseFloat(card.style.top);
			const left = Number.parseFloat(card.style.left);
			expect(Number.isFinite(top)).toBe(true);
			expect(Number.isFinite(left)).toBe(true);
			expect(top).toBeGreaterThanOrEqual(0);
			expect(left).toBeGreaterThanOrEqual(0);
			expect(top).toBeLessThanOrEqual(window.innerHeight);
			expect(left).toBeLessThanOrEqual(window.innerWidth);
		} finally {
			Object.defineProperty(window, 'innerWidth', { configurable: true, value: originalWidth });
			Object.defineProperty(window, 'innerHeight', { configurable: true, value: originalHeight });
		}
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
