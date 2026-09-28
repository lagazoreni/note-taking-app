import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
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

	it('filters questions, reports selections, and labels unlinked questions', async () => {
		const onSelect = vi.fn();
		const linkedQuestion = {
			id: '1',
			workspaceId: 'w',
			questionText: 'Which source?',
			status: 'unanswered' as const,
			priority: 'none' as const,
			dueDate: null,
			version: 1,
			answerMarkdown: null,
			linkedNotes: [{ id: 'n', title: 'Research', displayMode: 'collapsed' as const }],
			tagIds: [],
			createdAt: '',
			updatedAt: ''
		};
		const unlinkedQuestion = {
			...linkedQuestion,
			id: '2',
			questionText: 'What remains unlinked?',
			linkedNotes: []
		};

		render(QuestionPicker, {
			questions: [linkedQuestion, unlinkedQuestion],
			onSelect
		});

		expect(screen.getByText('What remains unlinked?')).toBeInTheDocument();
		expect(screen.getByText(/Unlinked/)).toBeInTheDocument();

		const search = screen.getByRole('textbox', { name: 'Link an existing question' });
		await fireEvent.input(search, { target: { value: 'unlinked' } });
		expect(screen.queryByText('Which source?')).not.toBeInTheDocument();
		expect(screen.getByText('What remains unlinked?')).toBeInTheDocument();

		await fireEvent.click(screen.getByRole('button', { name: /What remains unlinked\?/ }));
		expect(onSelect).toHaveBeenCalledOnce();
		expect(onSelect).toHaveBeenCalledWith(unlinkedQuestion);
	});
});
