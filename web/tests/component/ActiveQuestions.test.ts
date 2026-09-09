import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import QuestionList from '$lib/components/QuestionList.svelte';
import ActiveQuestionsPage from '../../src/routes/questions/+page.svelte';
import { questionsApi } from '$lib/api/questions';
import { currentWorkspaceId } from '$lib/stores/workspace';
import type { Question } from '$lib/types/question';

afterEach(() => {
	currentWorkspaceId.set(null);
	vi.restoreAllMocks();
});

describe('QuestionList', () => {
	it('shows an empty state when there are no questions', () => {
		render(QuestionList, { questions: [] });
		expect(screen.getByText(/no questions/i)).toBeInTheDocument();
	});
});

describe('Active Questions page', () => {
	it('creates an unlinked question with only question fields', async () => {
		const workspaceId = '11111111-1111-4111-8111-111111111111';
		const createdQuestion = {} as Question;
		vi.spyOn(questionsApi, 'list').mockResolvedValue({ items: [] });
		const create = vi.spyOn(questionsApi, 'create').mockResolvedValue(createdQuestion);
		currentWorkspaceId.set(workspaceId);

		render(ActiveQuestionsPage);
		await waitFor(() =>
			expect(screen.getByRole('button', { name: 'New question' })).toBeInTheDocument()
		);
		await fireEvent.click(screen.getByRole('button', { name: 'New question' }));
		const input = screen.getByRole('textbox', { name: /question text/i });
		await fireEvent.input(input, { target: { value: '  Why does this matter?  ' } });
		await fireEvent.click(screen.getByRole('button', { name: /create question|save question/i }));

		await waitFor(() => expect(create).toHaveBeenCalledOnce());
		expect(create).toHaveBeenCalledWith({
			workspaceId,
			questionText: 'Why does this matter?',
			kind: 'question',
			status: 'unanswered',
			priority: 'none',
			tagIds: []
		});
	});
});
