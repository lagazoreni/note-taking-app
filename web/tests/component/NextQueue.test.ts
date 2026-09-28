import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import NextQueue from '$lib/components/NextQueue.svelte';
import type { Question } from '$lib/types/question';

function makeQuestion(id: string): Question {
	return {
		id,
		workspaceId: '11111111-1111-4111-8111-111111111111',
		questionText: `Question ${id}`,
		status: 'unanswered',
		priority: 'none',
		dueDate: null,
		answerMarkdown: null,
		reminder: null,
		tagIds: [],
		linkedNotes: [],
		createdAt: '2026-04-01T00:00:00.000Z',
		updatedAt: '2026-04-01T00:00:00.000Z',
		version: 1,
		kind: 'question'
	};
}

describe('NextQueue', () => {
	it('shows guidance when there are no queue sections', () => {
		render(NextQueue, { sections: [] });

		expect(
			screen.getByText('Nothing in Next. Capture a question from a note, or set a due date.')
		).toBeInTheDocument();
	});

	it('renders the heading for each non-empty queue section', () => {
		const sections = [
			{ id: 'overdue', title: 'Overdue', items: [makeQuestion('overdue')] },
			{ id: 'due-today', title: 'Due today', items: [makeQuestion('due-today')] },
			{ id: 'in-progress', title: 'In progress', items: [makeQuestion('in-progress')] },
			{ id: 'deferred-ready', title: 'Deferred ready', items: [makeQuestion('deferred-ready')] },
			{ id: 'high-priority', title: 'High priority', items: [makeQuestion('high-priority')] }
		];

		render(NextQueue, { sections });

		for (const title of sections.map((section) => section.title)) {
			expect(screen.getByRole('heading', { name: title })).toBeInTheDocument();
		}
	});
});
