import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import QuestionPicker from '$lib/components/QuestionPicker.svelte';

describe('QuestionPicker', () => {
	it('shows context for existing questions', () => {
		render(QuestionPicker, {
			questions: [
				{
					id: '1',
					workspaceId: 'w',
					questionText: 'Which source?',
					status: 'unanswered',
					priority: 'none',
					dueDate: null,
					version: 1,
					answerMarkdown: null,
					reminder: null,
					linkedNotes: [{ id: 'n', title: 'Research', displayMode: 'collapsed' }],
					tagIds: [],
					createdAt: '',
					updatedAt: ''
				}
			]
		});
		expect(screen.getByText('Which source?')).toBeInTheDocument();
		expect(screen.getByText(/Research/)).toBeInTheDocument();
	});
});
